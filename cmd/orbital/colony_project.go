package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"orbital/colony/levels"
)

const colonyProjectLimit = 512 << 10

type colonyProject struct {
	Code  string            `json:"code"`
	Files map[string]string `json:"files"`
	Race  bool              `json:"race"`
}

func validColonyName(name string) bool {
	return validLabName(name) && filepath.Ext(name) != ".proto" && filepath.ToSlash(filepath.Clean(name)) == name
}

func validateColonyFiles(files map[string]string, lang string) error {
	size := 0
	for _, name := range projectNames(files) {
		if !validColonyName(name) {
			return fmt.Errorf(tr(lang, "اسم ملف غير صالح: %q", "Invalid file name: %q"), name)
		}
		size += len(files[name])
		if size > colonyProjectLimit {
			return fmt.Errorf("%s", tr(lang, "حجم ملفات المشروع أكبر من 512 KiB", "Project files exceed 512 KiB"))
		}
	}
	main, ok := files["main.go"]
	if !ok {
		return fmt.Errorf("%s", tr(lang, "المشروع محتاج main.go في package main", "The project needs main.go in package main"))
	}
	f, _ := parser.ParseFile(token.NewFileSet(), "main.go", main, parser.PackageClauseOnly)
	if f != nil && f.Name != nil && f.Name.Name != "main" {
		return fmt.Errorf("%s", tr(lang, "main.go لازم يكون في package main", "main.go must be in package main"))
	}
	return nil
}

func projectNames(files map[string]string) []string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func readColonyProject(w http.ResponseWriter, r *http.Request) (*colonyProject, bool) {
	var project colonyProject
	// JSON escaping can expand each source byte to six bytes.
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 6*colonyProjectLimit+64<<10))
	if err := decoder.Decode(&project); err != nil {
		requestError(w, r, "bad request", http.StatusBadRequest)
		return nil, false
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		requestError(w, r, "bad request", http.StatusBadRequest)
		return nil, false
	}
	if project.Files == nil {
		project.Files = map[string]string{"main.go": project.Code}
	}
	if err := validateColonyFiles(project.Files, requestLang(r)); err != nil {
		requestError(w, r, err.Error(), http.StatusBadRequest)
		return nil, false
	}
	for _, name := range projectNames(project.Files) {
		if strings.HasSuffix(name, ".go") {
			if missing := restrictedImports(project.Files[name], &levels.Level{Kind: "lab", Deps: colonyControlLevel().Deps}, requestLang(r)); len(missing) > 0 {
				writeJSON(w, rejectedImports(missing, requestLang(r)))
				return nil, false
			}
		}
	}
	return &project, true
}

func writeProjectFiles(dir string, files map[string]string) error {
	for _, name := range projectNames(files) {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(files[name]), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func compileErrorsFor(output, lang string) []compileError {
	errors := []compileError{}
	for _, line := range strings.Split(output, "\n") {
		if ce, ok := parseCompileLine(strings.TrimPrefix(strings.TrimSpace(line), "vet: "), lang); ok {
			errors = append(errors, ce)
		}
	}
	return errors
}

func (s *server) registerColonyProject(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/colony/test", s.colonyTest)
	mux.HandleFunc("POST /api/colony/vet", s.colonyVet)
	mux.HandleFunc("GET /api/colony/runs", s.colonyRuns)
}

func (s *server) colonyTest(w http.ResponseWriter, r *http.Request) {
	s.checkColonyProject(w, r, "test")
}
func (s *server) colonyVet(w http.ResponseWriter, r *http.Request) { s.checkColonyProject(w, r, "vet") }

func (s *server) checkColonyProject(w http.ResponseWriter, r *http.Request, command string) {
	project, ok := readColonyProject(w, r)
	if !ok || !s.checkControlImports(w, r, project) {
		return
	}
	ctx := withLang(r.Context(), requestLang(r))
	var passed bool
	var output string
	var err error
	if command == "test" {
		passed, output, err = s.run.goTest(ctx, colonyControlLevel(), project.Files, "./...")
	} else {
		passed, output, err = s.run.vetColonyProject(ctx, project.Files)
	}
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"ok": passed, "output": output, "compileErrors": compileErrorsFor(output, requestLang(r))})
}

func (r *runner) vetColonyProject(ctx context.Context, files map[string]string) (bool, string, error) {
	dir, err := r.colonyProjectDir(files)
	if err != nil {
		return false, "", err
	}
	defer os.RemoveAll(dir)
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	output := &limitedBuffer{max: 8000, lang: contextLang(ctx)}
	cmd := r.goCmd(ctx, dir, "vet", "./...")
	cmd.Stdout, cmd.Stderr = output, output
	runErr := runCommand(cmd)
	return runErr == nil, stripDir(output.String(), dir), nil
}

func (r *runner) colonyProjectDir(files map[string]string) (string, error) {
	dir, err := runDir("orbital-colony-*")
	if err != nil {
		return "", err
	}
	if err = r.writeModule(dir, colonyControlLevel()); err == nil {
		err = writeProjectFiles(dir, files)
	}
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}
	return dir, nil
}

func (s *server) runColonyProject(w http.ResponseWriter, r *http.Request) {
	project, ok := readColonyProject(w, r)
	if !ok || !s.checkControlImports(w, r, project) {
		return
	}
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.ops.Lock()
	defer c.ops.Unlock()
	if !s.checkColonyContract(w, r, c, project) {
		return
	}
	lang := requestLang(r)
	dir, err := s.run.colonyProjectDir(project.Files)
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	started := false
	defer func() {
		if !started {
			os.RemoveAll(dir)
		}
	}()
	if failure := s.run.buildColonyProject(withLang(r.Context(), lang), dir, project); failure != nil {
		writeJSON(w, failure)
		return
	}
	s.stopColonies()
	details := &colonyRunDetails{lang: lang, record: colonyRunRecord{CodeHash: colonyCodeHash(project.Files)}}
	if err := c.start(dir, lang, details); err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	started = true
	writeJSON(w, map[string]any{"phase": "run", "running": true})
}

func (s *server) checkColonyContract(w http.ResponseWriter, r *http.Request, c *colonyGame, project *colonyProject) bool {
	c.mu.Lock()
	contract, _ := c.state.ActiveContract()
	c.mu.Unlock()
	lang := requestLang(r)
	if contract != nil {
		missing := append(projectRequirements(project.Files, contract.Requires, lang), c.controlRequirements(lang)...)
		if contract.ID == "c12-one-ledger" && !project.Race {
			missing = append(missing, tr(lang, "شغّل فحص السباقات للعقد ده (race:true).", "Enable race mode for this contract (race:true)."))
		}
		if len(missing) > 0 {
			writeJSON(w, map[string]any{"phase": "requires", "missing": missing, "contract": contract.ID})
			return false
		}
		if contract.ID == "c16-tested-planner" {
			response, err := s.run.checkColonyPlanner(withLang(r.Context(), lang), project.Files)
			if err != nil {
				requestError(w, r, err.Error(), 500)
				return false
			}
			if !response.Result.Success {
				writeJSON(w, response)
				return false
			}
		}
	}
	return true
}

func (r *runner) buildColonyProject(ctx context.Context, dir string, project *colonyProject) map[string]any {
	lang := contextLang(ctx)
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	output := &limitedBuffer{max: 64 << 10, lang: lang}
	args := []string{"build", "-o", "prog"}
	if project.Race {
		args = append(args, "-race")
	}
	debugCtx, cancelDebug := context.WithCancel(ctx)
	defer cancelDebug()
	debugBuilt := r.buildColonyDebug(debugCtx, dir, project)
	build := r.goCmd(ctx, dir, append(args, ".")...)
	build.Stdout, build.Stderr = output, output
	err := runCommand(build)
	if err != nil {
		cancelDebug()
	}
	if !<-debugBuilt {
		os.Remove(filepath.Join(dir, colonyDebugProg))
	}
	if err != nil {
		out := stripDir(output.String(), dir)
		return map[string]any{"phase": "compile", "compileErrors": compileErrorsFor(out, lang), "compileOutput": out}
	}
	return nil
}
