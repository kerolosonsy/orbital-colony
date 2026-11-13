package main

import (
	"context"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"orbital/colony/econ"
	"orbital/colony/levels"
)

func readColonyFixture(t *testing.T, id string) map[string]string {
	t.Helper()
	dir := filepath.Join("testdata", "colony", id)
	files := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		name, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		body, err := os.ReadFile(path)
		files[filepath.ToSlash(name)] = string(body)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func projectPost(t *testing.T, handler http.HandlerFunc, project colonyProject) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(project)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/", strings.NewReader(string(body)))
	req.Header.Set("X-Orbital-Lang", "en")
	response := httptest.NewRecorder()
	handler(response, req)
	return response
}

func TestColonyProjectValidation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		files map[string]string
		ok    bool
	}{
		{"subpackages", map[string]string{"main.go": "package main", "plan/plan.go": "package plan", "notes.md": "notes"}, true},
		{"missing main", map[string]string{"plan/plan.go": "package plan"}, false},
		{"wrong package", map[string]string{"main.go": "package plan"}, false},
		{"parent", map[string]string{"main.go": "package main", "../outside.go": ""}, false},
		{"absolute", map[string]string{"main.go": "package main", "/outside.go": ""}, false},
		{"backslash", map[string]string{"main.go": "package main", `plan\plan.go`: ""}, false},
		{"module", map[string]string{"main.go": "package main", "go.mod": ""}, false},
		{"proto", map[string]string{"main.go": "package main", "rpc.proto": ""}, false},
		{"alias", map[string]string{"main.go": "package main", "./main.go": ""}, false},
		{"boundary", map[string]string{"main.go": "package main", "notes.txt": strings.Repeat("x", colonyProjectLimit-len("package main"))}, true},
		{"oversize", map[string]string{"main.go": "package main", "notes.txt": strings.Repeat("x", colonyProjectLimit)}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateColonyFiles(tc.files, "en"); (err == nil) != tc.ok {
				t.Fatalf("validation: %v", err)
			}
		})
	}
}

func TestColonyProjectImportPolicyEveryFile(t *testing.T) {
	s := testColony(t)
	for _, pkg := range []string{"os", "net/http", "player/plan", "unsafe", "os/exec", "syscall", "plugin", "example.com/thirdparty", "orbital/colony/econ"} {
		t.Run(pkg, func(t *testing.T) {
			files := map[string]string{"main.go": "package main\nfunc main(){}", "plan/plan.go": "package plan\nimport _ " + `"` + pkg + `"`}
			response := projectPost(t, s.colonyTest, colonyProject{Files: files})
			denied := pkg == "unsafe" || pkg == "os/exec" || pkg == "syscall" || pkg == "plugin" || pkg == "example.com/thirdparty" || pkg == "orbital/colony/econ"
			var reply struct {
				Phase   string
				Missing []string
			}
			if err := json.Unmarshal(response.Body.Bytes(), &reply); err != nil {
				t.Fatal(err)
			}
			if (reply.Phase == "require") != denied || denied && len(reply.Missing) != 1 {
				t.Fatalf("policy: %s", response.Body)
			}
		})
	}
}

func TestColonyProjectThreePackagesRunAndLegacyCode(t *testing.T) {
	s := testColony(t)
	for _, project := range []colonyProject{
		{Files: map[string]string{"main.go": "package main\nimport (\"fmt\";\"player/plan\")\nfunc main(){fmt.Println(plan.Value())}", "plan/plan.go": "package plan\nimport \"player/value\"\nfunc Value() int{return value.Number}", "value/value.go": "package value\nconst Number=42"}},
		{Code: "package main\nfunc main(){}"},
	} {
		response := projectPost(t, s.colonyRun, project)
		if !strings.Contains(response.Body.String(), `"phase":"run"`) {
			t.Fatalf("run: %s", response.Body)
		}
		colonyRequest(t, s.colonyStop, "")
	}
	runs, err := s.colony.readRuns()
	if err != nil || len(runs) != 2 {
		t.Fatalf("history after natural exit or stop: %v %v", runs, err)
	}
}

func TestColonyProjectTestAndVetEndpoints(t *testing.T) {
	s := testColony(t)
	for _, tc := range []struct {
		name         string
		handler      http.HandlerFunc
		source, test string
		ok, compile  bool
	}{
		{"test pass", s.colonyTest, "package main\nfunc main(){}", "package main\nimport \"testing\"\nfunc TestOK(t *testing.T){}", true, false},
		{"test fail", s.colonyTest, "package main\nfunc main(){}", "package main\nimport \"testing\"\nfunc TestFail(t *testing.T){t.Fatal(\"broken\")}", false, false},
		{"test compile", s.colonyTest, "package main\nfunc main(){missing()}", "", false, true},
		{"vet pass", s.colonyVet, "package main\nfunc main(){}", "", true, false},
		{"vet fail", s.colonyVet, "package main\nimport \"fmt\"\nfunc main(){fmt.Printf(\"%d\",\"bad\")}", "", false, true},
		{"vet compile", s.colonyVet, "package main\nfunc main(){missing()}", "", false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{"main.go": tc.source}
			if tc.test != "" {
				files["main_test.go"] = tc.test
			}
			response := projectPost(t, tc.handler, colonyProject{Files: files})
			var reply struct {
				OK            bool
				Output        string
				CompileErrors []compileError
			}
			if err := json.Unmarshal(response.Body.Bytes(), &reply); err != nil {
				t.Fatal(err)
			}
			if reply.OK != tc.ok || (len(reply.CompileErrors) > 0) != tc.compile {
				t.Fatalf("check: %s", response.Body)
			}
		})
	}
	// Checks must not stop an active colony.
	response := projectPost(t, s.colonyRun, colonyProject{Code: levels.Colony.Starter})
	if !strings.Contains(response.Body.String(), `"running":true`) {
		t.Fatal(response.Body)
	}
	for _, handler := range []http.HandlerFunc{s.colonyTest, s.colonyVet} {
		response = projectPost(t, handler, colonyProject{Code: "package main\nfunc main(){}"})
		if !strings.Contains(response.Body.String(), `"ok":true`) {
			t.Fatal(response.Body)
		}
	}
	s.colony.mu.Lock()
	active := s.colony.active
	s.colony.mu.Unlock()
	if active == nil {
		t.Fatal("check stopped colony")
	}
}

func TestColonyProjectRequirementsAcrossFiles(t *testing.T) {
	for _, files := range []map[string]string{
		{"main.go": "package main\nfunc main(){f();f()}", "f.go": "package main\nfunc f(){for i:=0;i<2;i++{}}", "plan/plan.go": "package plan\nfunc F(){for {break}}"},
		{"main.go": "package main\nimport p \"player/plan\"\nfunc main(){for i:=0;i<2;i++{p.F();p.F()}}", "plan/plan.go": "package plan\nfunc F(){for {break}}"},
	} {
		if missing := projectRequirements(files, []string{"for2", "func"}, "en"); len(missing) != 0 {
			t.Fatal(missing)
		}
	}
	if missing := projectRequirements(map[string]string{"main.go": "package main\nfunc main(){}"}, []string{"mutex"}, "en"); len(missing) != 1 {
		t.Fatal(missing)
	}
}

func TestColonyWorkspaceProjectRoundTrip(t *testing.T) {
	runner := testRunner(t)
	ws := &workspace{root: filepath.Join(t.TempDir(), "workspace"), colonyRel: runner.colonyDir}
	rel, err := filepath.Rel(ws.root, runner.colonyDir)
	if err != nil {
		t.Fatal(err)
	}
	ws.colonyRel = rel
	files := map[string]string{"main.go": "package main\nimport \"player/plan\"\nfunc main(){plan.F()}", "plan/plan.go": "package plan\nimport \"player/value\"\nfunc F(){_ = value.Number}", "value/value.go": "package value\nconst Number=42", "notes/readme.md": "player notes"}
	path, err := ws.open(levels.Colony, files)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(files["main.go"]+"\n// edited"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ws.open(levels.Colony, files); err != nil {
		t.Fatal(err)
	}
	got, mtime, err := ws.read(levels.Colony)
	if err != nil || len(got) != len(files) || mtime.IsZero() || !strings.Contains(got["main.go"], "// edited") {
		t.Fatalf("round trip: %v %v %v", got, mtime, err)
	}
	for name, body := range files {
		if name != "main.go" && got[name] != body {
			t.Fatalf("lost %s", name)
		}
	}
	ctx, cancel := contextWithTestTimeout(t)
	defer cancel()
	cmd := runner.goCmd(ctx, filepath.Dir(path), "test", "./...")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("workspace module: %v\n%s", err, out)
	}
}

func TestColonyProjectPlannerMutants(t *testing.T) {
	r := testRunner(t)
	ctx, cancel := contextWithTestTimeout(t)
	defer cancel()
	files := readColonyFixture(t, "c16-tested-planner")
	response, err := r.checkColonyPlanner(withLang(ctx, "en"), files)
	if err != nil || !response.Result.Success {
		t.Fatalf("planner reference: %+v %v", response, err)
	}
	files["plan/plan_test.go"] = "package plan\nimport \"testing\"\nfunc TestWeak(t *testing.T){if Next([]int{4})!=0{t.Fatal(\"single\")}}"
	response, err = r.checkColonyPlanner(withLang(ctx, "en"), files)
	if err != nil || response.Result.Success {
		t.Fatalf("weak tests accepted: %+v %v", response, err)
	}
	for _, step := range response.Steps[1:4] {
		if step.OK {
			t.Fatalf("mutant not rejected: %+v", step)
		}
	}
}

func contextWithTestTimeout(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), 90*time.Second)
}

func TestColonyProjectRejectsInvalidRequestsBeforeRun(t *testing.T) {
	s := testColony(t)
	for _, body := range []string{
		`{"files":{}}`,
		`{"files":{"main.go":"package other"}}`,
		`{"files":{"main.go":"package main","../outside.go":"package main"}}`,
		`{"code":"package main"} {}`,
		`null`,
	} {
		response := httptest.NewRecorder()
		s.colonyRun(response, httptest.NewRequest("POST", "/api/colony/run", strings.NewReader(body)))
		if response.Code != http.StatusBadRequest {
			t.Fatalf("accepted %s: %d %s", body, response.Code, response.Body)
		}
	}
	response := projectPost(t, s.colonyRun, colonyProject{Files: map[string]string{"main.go": "package main\nfunc main(){}", "notes.txt": strings.Repeat("x", colonyProjectLimit)}})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("oversized project: %d %s", response.Code, response.Body)
	}
}

func TestColonyProjectPlannerContractGate(t *testing.T) {
	s := testColony(t)
	s.colony.state = contractFixture(econ.FindContract("c16-tested-planner"))
	files := readColonyFixture(t, "c16-tested-planner")
	files["plan/plan_test.go"] = "package plan"
	response := projectPost(t, s.colonyRun, colonyProject{Files: files})
	var reply runResponse
	if err := json.Unmarshal(response.Body.Bytes(), &reply); err != nil {
		t.Fatal(err)
	}
	if reply.Phase != "lab" || reply.Result.Success || s.colony.active != nil {
		t.Fatalf("untested planner started: %s", response.Body)
	}
}

func TestColonyWorkspaceRejectsSymlinks(t *testing.T) {
	ws := &workspace{root: t.TempDir()}
	outside := filepath.Join(t.TempDir(), "outside.go")
	if err := os.WriteFile(outside, []byte("preserve"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(ws.root, "colony"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(ws.root, "colony", "go.mod")); err != nil {
		t.Fatal(err)
	}
	if _, err := ws.open(levels.Colony, map[string]string{"main.go": "package main"}); err == nil {
		t.Fatal("overwrote symlink")
	}
	if body, err := os.ReadFile(outside); err != nil || string(body) != "preserve" {
		t.Fatalf("external file changed: %s %v", body, err)
	}
}
