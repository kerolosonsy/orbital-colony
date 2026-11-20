package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"orbital/colony/levels"
)

const (
	traceFrameLimit  = 2000
	traceByteLimit   = 256 << 10
	traceStdoutLimit = 64 << 10
	traceValueLimit  = 160
	traceNameLimit   = 96
	traceVarLimit    = 64
)

type traceVariable struct {
	Name      string `json:"name"`
	Formatted string `json:"value"`
}

type traceFrame struct {
	File      string          `json:"file"`
	Line      int             `json:"line"`
	Func      string          `json:"func"`
	Depth     int             `json:"depth"`
	Event     string          `json:"event"`
	Goroutine int             `json:"goroutine"`
	Variables []traceVariable `json:"variables"`
	Returns   []traceVariable `json:"returns,omitempty"`
	Stdout    int64           `json:"stdout"`
	Sim       int64           `json:"sim,omitempty"`
	Events    int             `json:"events"`
	Truncated bool            `json:"truncated,omitempty"`
}

type traceResponse struct {
	Available       bool         `json:"available"`
	Completed       bool         `json:"completed"`
	Frames          []traceFrame `json:"frames"`
	Stdout          string       `json:"stdout"`
	StdoutTruncated bool         `json:"stdoutTruncated,omitempty"`
	Truncated       bool         `json:"truncated,omitempty"`
	Message         string       `json:"message,omitempty"`
	// Synced is true when every frame's Events count refers to the event
	// stream of the same run the client replays on the map.
	Synced bool `json:"synced,omitempty"`
}

type traceRequest struct {
	Level string            `json:"level"`
	Code  string            `json:"code"`
	Files map[string]string `json:"files"`
}

type traceStdoutBuffer struct {
	prefix    bytes.Buffer
	truncated bool
}

func (b *traceStdoutBuffer) Write(part []byte) (int, error) {
	n := len(part)
	remaining := traceStdoutLimit - b.prefix.Len()
	if remaining > 0 {
		_, _ = b.prefix.Write(part[:min(len(part), remaining)])
	}
	if n > remaining {
		b.truncated = true
	}
	return n, nil
}

func (b *traceStdoutBuffer) String() string { return b.prefix.String() }

type traceRecorderSource struct {
	name string
	text string
}

func traceProjectFiles(lvl *levels.Level, code string, edited map[string]string) map[string]string {
	if lvl.Kind != "lab" {
		return map[string]string{"main.go": code}
	}
	files := make(map[string]string, len(lvl.Files))
	for _, file := range lvl.Files {
		files[file.Name] = file.Starter
		if body, ok := edited[file.Name]; ok && !file.ReadOnly && strings.HasSuffix(file.Name, ".go") {
			files[file.Name] = body
		}
	}
	return files
}

func traceMainPresent(files map[string]string) bool {
	for name, body := range files {
		if filepath.Base(name) != name || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), name, body, 0)
		if err != nil || file.Name.Name != "main" {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Recv == nil && fn.Name.Name == "main" {
				return true
			}
		}
	}
	return false
}

func traceImportsAllowed(lvl *levels.Level, code string, project map[string]string, lang string) bool {
	if lvl.Kind != "lab" {
		return len(restrictedImports(code, lvl, lang)) == 0
	}
	for _, file := range lvl.Files {
		if !file.ReadOnly && strings.HasSuffix(file.Name, ".go") && len(restrictedImports(project[file.Name], lvl, lang)) > 0 {
			return false
		}
	}
	return true
}

func (r *runner) Trace(ctx context.Context, lvl *levels.Level, code string, edited map[string]string) (*traceResponse, error) {
	lang := contextLang(ctx)
	if !traceLevelSupported(lvl) {
		return traceUnavailable(lang), nil
	}
	files := traceProjectFiles(lvl, code, edited)
	if !traceMainPresent(files) {
		return traceUnavailable(lang), nil
	}
	return r.runTraceWorkspace(ctx, lvl, files, lang)
}

func traceLevelSupported(lvl *levels.Level) bool {
	return lvl != nil && (lvl.Kind == "" || lvl.Kind == "lab")
}

func (r *runner) runTraceWorkspace(ctx context.Context, lvl *levels.Level, files map[string]string, lang string) (*traceResponse, error) {
	rewritten, recorder, err := rewriteTraceFiles(files)
	if err != nil {
		return traceUnavailable(lang), nil
	}
	dir, err := os.MkdirTemp("", "orbital-trace-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	if err := r.writeTraceProject(dir, lvl, rewritten, recorder); err != nil {
		return nil, err
	}
	return r.buildAndRunTrace(ctx, dir, lvl, lang)
}

func (r *runner) writeTraceProject(dir string, lvl *levels.Level, rewritten map[string]string, recorder traceRecorderSource) error {
	if err := r.writeModule(dir, lvl); err != nil {
		return err
	}
	if err := writeProjectFiles(dir, rewritten); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, recorder.name), []byte(recorder.text), 0o644)
}

func traceUnavailable(lang string) *traceResponse {
	return &traceResponse{Available: false, Frames: []traceFrame{}, Message: tr(lang, "التتبّع مش متاح للكود ده.", "Trace is unavailable for this code.")}
}

func (r *runner) buildAndRunTrace(ctx context.Context, dir string, lvl *levels.Level, lang string) (*traceResponse, error) {
	buildTimeout := 90 * time.Second
	if len(lvl.Deps) > 0 {
		buildTimeout = 240 * time.Second
	}
	bctx, cancel := context.WithTimeout(ctx, buildTimeout)
	defer cancel()
	args := []string{"build"}
	if lvl.Race {
		args = append(args, "-race")
	}
	args = append(args, "-o", "trace-program", ".")
	build := r.goCmd(bctx, dir, args...)
	build.Stdout, build.Stderr = &limitedBuffer{max: 64 << 10, lang: lang}, &limitedBuffer{max: 64 << 10, lang: lang}
	if err := runCommand(build); err != nil {
		return traceUnavailable(lang), nil
	}
	return r.executeTrace(ctx, dir, lvl, lang)
}

func (r *runner) executeTrace(ctx context.Context, dir string, lvl *levels.Level, lang string) (*traceResponse, error) {
	rctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tracePath := filepath.Join(dir, "learner-trace.jsonl")
	simPath := filepath.Join(dir, "colony-events.jsonl")
	traceFile, err := os.OpenFile(tracePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o600)
	if err != nil {
		return traceUnavailable(lang), nil
	}
	defer traceFile.Close()
	stdout := &traceStdoutBuffer{}
	cmd := exec.CommandContext(rctx, filepath.Join(dir, "trace-program"))
	configureProcess(cmd)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "ORBITAL_LANG="+lang, "ORBITAL_LEVEL="+lvl.ID, "ORBITAL_TRACE="+simPath, "ORBITAL_USER_TRACE="+tracePath)
	if seed, ok := ctx.Value(reviewSeedKey{}).(int64); ok {
		cmd.Env = append(cmd.Env, "ORBITAL_SEED="+strconv.FormatInt(seed, 10))
	}
	cmd.Stdout = stdout
	cmd.Stderr = &limitedBuffer{max: 64 << 10, lang: lang}
	runErr := runCommand(cmd)
	resp := readTraceFile(traceFile, lang)
	resp.Completed = runErr == nil
	resp.Stdout = stdout.String()
	resp.StdoutTruncated = stdout.truncated
	return resp, nil
}

func readTraceFile(traceFile *os.File, lang string) *traceResponse {
	if _, err := traceFile.Seek(0, io.SeekStart); err != nil {
		return traceUnavailable(lang)
	}
	traceData, err := io.ReadAll(io.LimitReader(traceFile, traceByteLimit+1))
	if err != nil || len(traceData) > traceByteLimit {
		return traceUnavailable(lang)
	}
	return parseTraceFrames(traceData, lang)
}

func parseTraceFrames(traceData []byte, lang string) *traceResponse {
	resp := &traceResponse{Available: true, Frames: []traceFrame{}}
	for _, line := range bytes.Split(traceData, []byte("\n")) {
		addTraceLine(resp, line)
	}
	if len(resp.Frames) == 0 {
		resp.Message = tr(lang, "مفيش سطور اتنفذت نقدر نعرضها.", "No executed lines were available to show.")
	}
	return resp
}

func addTraceLine(resp *traceResponse, line []byte) {
	if len(bytes.TrimSpace(line)) == 0 {
		return
	}
	if isTraceTruncation(line) {
		resp.Truncated = true
		return
	}
	frame, ok := decodeTraceFrame(line)
	if ok {
		resp.Frames = append(resp.Frames, frame)
		resp.Truncated = resp.Truncated || frame.Truncated
	}
}

func isTraceTruncation(line []byte) bool {
	return bytes.Equal(line, []byte(`{"truncated":true}`))
}

func decodeTraceFrame(line []byte) (traceFrame, bool) {
	var frame traceFrame
	if err := json.Unmarshal(line, &frame); err != nil || frame.Line <= 0 {
		return traceFrame{}, false
	}
	return frame, true
}

type traceInstrumenter struct {
	fset         *token.FileSet
	file         string
	prefix       string
	functionName string
	returnNames  []string
	returnCount  int
	returnID     int
}

func rewriteTraceFiles(files map[string]string) (map[string]string, traceRecorderSource, error) {
	prefix := tracePrefix(files)
	rewritten := make(map[string]string, len(files))
	for _, name := range projectNames(files) {
		body, err := rewriteTraceFile(name, files[name], prefix)
		if err != nil {
			return nil, traceRecorderSource{}, err
		}
		rewritten[name] = body
	}
	recorderName := traceRecorderName(files)
	if recorderName == "" {
		return nil, traceRecorderSource{}, errors.New("no safe trace recorder filename is available")
	}
	return rewritten, traceRecorderSource{name: recorderName, text: traceHelperSource(prefix)}, nil
}

func tracePrefix(files map[string]string) string {
	var source strings.Builder
	for _, name := range projectNames(files) {
		source.WriteString(files[name])
	}
	prefix := "__orbitalTrace"
	for strings.Contains(source.String(), prefix) {
		prefix += "x"
	}
	return prefix
}

func rewriteTraceFile(name, body, prefix string) (string, error) {
	if filepath.Base(name) != name || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
		return body, nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, body, parser.ParseComments)
	if err != nil {
		return "", err
	}
	if file.Name.Name == "main" {
		rewriteTraceOutput(file, prefix)
		instrumentTraceFile(fset, name, prefix, file)
	}
	return formatTraceFile(fset, file)
}

type traceOutputAliases struct {
	fmt map[string]bool
	io  map[string]bool
	os  map[string]bool
}

func rewriteTraceOutput(file *ast.File, prefix string) {
	aliases := collectTraceOutputAliases(file)
	rewriteTraceCalls(file, prefix, aliases)
	removeUnusedOutputImports(file, aliases)
}

func collectTraceOutputAliases(file *ast.File) traceOutputAliases {
	aliases := traceOutputAliases{fmt: map[string]bool{}, io: map[string]bool{}, os: map[string]bool{}}
	for _, imported := range file.Imports {
		addTraceOutputAlias(&aliases, imported)
	}
	return aliases
}

func addTraceOutputAlias(aliases *traceOutputAliases, imported *ast.ImportSpec) {
	path, alias := importPath(imported), traceImportAlias(imported)
	if path == "" || alias == "." || alias == "_" {
		return
	}
	if path == "fmt" {
		aliases.fmt[alias] = true
	} else if path == "io" {
		aliases.io[alias] = true
	} else if path == "os" {
		aliases.os[alias] = true
	}
}

func rewriteTraceCalls(file *ast.File, prefix string, aliases traceOutputAliases) {
	ast.Inspect(file, func(node ast.Node) bool {
		callExpr, isCall := node.(*ast.CallExpr)
		if isCall {
			rewriteTraceCall(callExpr, prefix, aliases)
		}
		return true
	})
}

func rewriteTraceCall(callExpr *ast.CallExpr, prefix string, aliases traceOutputAliases) {
	selector, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	pkg, isPackage := selector.X.(*ast.Ident)
	if isPackage && aliases.fmt[pkg.Name] {
		rewriteFmtOutput(callExpr, selector, prefix)
		return
	}
	if isPackage && aliases.io[pkg.Name] && selector.Sel.Name == "WriteString" && len(callExpr.Args) == 2 {
		callExpr.Fun = &ast.Ident{NamePos: selector.Sel.Pos(), Name: prefix + "WriteString"}
		return
	}
	rewriteStdoutWrite(callExpr, selector, prefix, aliases.os)
}

func rewriteFmtOutput(callExpr *ast.CallExpr, selector *ast.SelectorExpr, prefix string) {
	outputWrapper, supported := map[string]string{"Print": "Print", "Println": "Println", "Printf": "Printf", "Fprint": "Fprint", "Fprintln": "Fprintln", "Fprintf": "Fprintf"}[selector.Sel.Name]
	if supported {
		callExpr.Fun = &ast.Ident{NamePos: selector.Sel.Pos(), Name: prefix + outputWrapper}
	}
}

func rewriteStdoutWrite(callExpr *ast.CallExpr, selector *ast.SelectorExpr, prefix string, osAliases map[string]bool) {
	method, ok := selector.X.(*ast.SelectorExpr)
	if !ok || method.Sel.Name != "Stdout" || len(callExpr.Args) != 1 {
		return
	}
	osPackage, ok := method.X.(*ast.Ident)
	if !ok || !osAliases[osPackage.Name] {
		return
	}
	outputWrapper := map[string]string{"Write": "WriteStdout", "WriteString": "WriteStdoutString"}[selector.Sel.Name]
	if outputWrapper != "" {
		callExpr.Fun = &ast.Ident{NamePos: selector.Sel.Pos(), Name: prefix + outputWrapper}
		callExpr.Args = append([]ast.Expr{method}, callExpr.Args...)
	}
}

func removeUnusedOutputImports(file *ast.File, aliases traceOutputAliases) {
	for _, outputAliases := range []map[string]bool{aliases.fmt, aliases.io} {
		for alias := range outputAliases {
			if !hasPackageSelector(file, alias) {
				removeImport(file, alias)
			}
		}
	}
}

func hasPackageSelector(file *ast.File, alias string) bool {
	used := false
	ast.Inspect(file, func(node ast.Node) bool {
		selector, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		pkg, isPackage := selector.X.(*ast.Ident)
		used = used || isPackage && pkg.Name == alias
		return true
	})
	return used
}

func removeImport(file *ast.File, alias string) {
	file.Decls = filterImportDeclarations(file.Decls, alias)
	file.Imports = collectImportSpecs(file.Decls)
}

func filterImportDeclarations(declarations []ast.Decl, alias string) []ast.Decl {
	filtered := make([]ast.Decl, 0, len(declarations))
	for _, declaration := range declarations {
		gen, ok := declaration.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			filtered = append(filtered, declaration)
			continue
		}
		gen.Specs = filterImportSpecs(gen.Specs, alias)
		if len(gen.Specs) > 0 {
			filtered = append(filtered, gen)
		}
	}
	return filtered
}

func filterImportSpecs(specs []ast.Spec, alias string) []ast.Spec {
	filtered := make([]ast.Spec, 0, len(specs))
	for _, spec := range specs {
		imported, ok := spec.(*ast.ImportSpec)
		if !ok || traceImportAlias(imported) != alias {
			filtered = append(filtered, spec)
		}
	}
	return filtered
}

func collectImportSpecs(declarations []ast.Decl) []*ast.ImportSpec {
	var imports []*ast.ImportSpec
	for _, declaration := range declarations {
		gen, ok := declaration.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			continue
		}
		for _, spec := range gen.Specs {
			if imported, ok := spec.(*ast.ImportSpec); ok {
				imports = append(imports, imported)
			}
		}
	}
	return imports
}

func traceImportAlias(imported *ast.ImportSpec) string {
	if imported.Name != nil {
		return imported.Name.Name
	}
	return filepath.Base(importPath(imported))
}

func importPath(imported *ast.ImportSpec) string {
	path, err := strconv.Unquote(imported.Path.Value)
	if err != nil {
		return ""
	}
	return path
}

func instrumentTraceFile(fset *token.FileSet, name, prefix string, file *ast.File) {
	instrument := &traceInstrumenter{fset: fset, file: name, prefix: prefix}
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Body != nil {
			instrument.function(fn)
		}
	}
}

func formatTraceFile(fset *token.FileSet, file *ast.File) (string, error) {
	var output bytes.Buffer
	if err := format.Node(&output, fset, file); err != nil {
		return "", err
	}
	return output.String(), nil
}

// The recorder file sorts first so learner initializers cannot replace its side channel.
func traceRecorderName(files map[string]string) string {
	for _, recorderName := range []string{"0.go", "0_orbital_trace.go"} {
		if _, exists := files[recorderName]; exists {
			continue
		}
		for name := range files {
			if filepath.Base(name) == name && strings.HasSuffix(name, ".go") && name < recorderName {
				recorderName = ""
				break
			}
		}
		if recorderName != "" {
			return recorderName
		}
	}
	return ""
}

func (t *traceInstrumenter) function(fn *ast.FuncDecl) {
	previousFunction, previousReturns, previousReturnCount := t.functionName, t.returnNames, t.returnCount
	t.functionName = fn.Name.Name
	t.returnNames = namedResultNames(fn.Type.Results)
	t.returnCount = resultFieldCount(fn.Type.Results)
	env := localNames(fn.Recv, fn.Type.Params, fn.Type.Results)
	t.instrumentBlock(fn.Body, env)
	t.appendImplicitReturn(fn.Body, t.line(fn.Body.Rbrace))
	fn.Body.List = append([]ast.Stmt{t.enterAt(t.line(fn.Pos()), env), t.exitCall()}, fn.Body.List...)
	t.functionName, t.returnNames, t.returnCount = previousFunction, previousReturns, previousReturnCount
}

func (t *traceInstrumenter) literal(fn *ast.FuncLit, outer map[string]bool) {
	previousFunction, previousReturns, previousReturnCount := t.functionName, t.returnNames, t.returnCount
	t.functionName = "func@" + strconv.Itoa(t.line(fn.Type.Func))
	t.returnNames = namedResultNames(fn.Type.Results)
	t.returnCount = resultFieldCount(fn.Type.Results)
	env := cloneNames(outer)
	addFields(env, fn.Type.Params)
	addFields(env, fn.Type.Results)
	t.instrumentBlock(fn.Body, env)
	t.appendImplicitReturn(fn.Body, t.line(fn.Body.Rbrace))
	fn.Body.List = append([]ast.Stmt{t.enterAt(t.line(fn.Type.Func), env), t.exitCall()}, fn.Body.List...)
	t.functionName, t.returnNames, t.returnCount = previousFunction, previousReturns, previousReturnCount
}

func (t *traceInstrumenter) instrumentBlock(block *ast.BlockStmt, outer map[string]bool) {
	if block != nil {
		block.List = t.instrumentList(block.List, outer)
	}
}

func (t *traceInstrumenter) instrumentList(statements []ast.Stmt, outer map[string]bool) []ast.Stmt {
	env := cloneNames(outer)
	output := make([]ast.Stmt, 0, len(statements)*2)
	for _, statement := range statements {
		output = append(output, t.recordCall(statement.Pos(), env), statement)
		t.instrumentStatement(statement, env)
		if returned, ok := statement.(*ast.ReturnStmt); ok {
			output[len(output)-1] = t.instrumentReturn(returned, env)
		}
		addStatementNames(env, statement)
	}
	return output
}

func (t *traceInstrumenter) instrumentStatement(statement ast.Stmt, env map[string]bool) {
	switch node := statement.(type) {
	case *ast.BlockStmt:
		t.instrumentBlock(node, env)
	case *ast.IfStmt:
		t.instrumentIf(node, env)
	case *ast.ForStmt:
		t.instrumentFor(node, env)
	case *ast.RangeStmt:
		t.instrumentRange(node, env)
	case *ast.SwitchStmt:
		t.instrumentSwitch(node, env)
	case *ast.TypeSwitchStmt:
		t.instrumentTypeSwitch(node, env)
	case *ast.SelectStmt:
		t.instrumentSelect(node, env)
	case *ast.LabeledStmt:
		t.instrumentLabeled(node, env)
	default:
		t.instrumentLiterals(statement, env)
	}
}

func (t *traceInstrumenter) instrumentLabeled(node *ast.LabeledStmt, env map[string]bool) {
	t.instrumentStatement(node.Stmt, env)
}

func (t *traceInstrumenter) instrumentIf(node *ast.IfStmt, env map[string]bool) {
	t.instrumentLiterals(node.Init, env)
	branch := cloneNames(env)
	addStatementNames(branch, node.Init)
	t.instrumentLiterals(node.Cond, branch)
	t.instrumentBlock(node.Body, branch)
	if nested, ok := node.Else.(*ast.IfStmt); ok {
		node.Else = t.instrumentElseIf(nested, branch)
	} else if node.Else != nil {
		t.instrumentStatement(node.Else, branch)
	}
}

func (t *traceInstrumenter) instrumentElseIf(node *ast.IfStmt, env map[string]bool) ast.Stmt {
	record := t.recordCall(node.Pos(), env)
	t.instrumentIf(node, env)
	return &ast.BlockStmt{List: []ast.Stmt{record, node}}
}

func (t *traceInstrumenter) instrumentFor(node *ast.ForStmt, env map[string]bool) {
	t.instrumentLiterals(node.Init, env)
	loop := cloneNames(env)
	addStatementNames(loop, node.Init)
	t.instrumentLiterals(node.Cond, loop)
	if node.Post != nil {
		t.instrumentLiterals(node.Post, loop)
		node.Post = t.instrumentPost(node.Post, loop)
	}
	t.instrumentBlock(node.Body, loop)
	node.Body.List = append([]ast.Stmt{t.recordCall(node.Pos(), loop)}, node.Body.List...)
}

func (t *traceInstrumenter) instrumentPost(post ast.Stmt, env map[string]bool) ast.Stmt {
	body := &ast.BlockStmt{List: []ast.Stmt{t.recordCall(post.Pos(), env), post}}
	function := &ast.FuncLit{Type: &ast.FuncType{Params: &ast.FieldList{}}, Body: body}
	return &ast.ExprStmt{X: &ast.CallExpr{Fun: function}}
}

func (t *traceInstrumenter) instrumentRange(node *ast.RangeStmt, env map[string]bool) {
	t.instrumentLiterals(node.X, env)
	loop := cloneNames(env)
	if node.Tok == token.DEFINE {
		addIdentifier(loop, node.Key)
		addIdentifier(loop, node.Value)
	}
	t.instrumentBlock(node.Body, loop)
	node.Body.List = append([]ast.Stmt{t.recordCall(node.Pos(), loop)}, node.Body.List...)
}

func (t *traceInstrumenter) instrumentSwitch(node *ast.SwitchStmt, env map[string]bool) {
	t.instrumentLiterals(node.Init, env)
	branch := cloneNames(env)
	addStatementNames(branch, node.Init)
	t.instrumentLiterals(node.Tag, branch)
	t.instrumentCases(node.Body.List, branch)
}

func (t *traceInstrumenter) instrumentTypeSwitch(node *ast.TypeSwitchStmt, env map[string]bool) {
	t.instrumentLiterals(node.Init, env)
	branch := cloneNames(env)
	addStatementNames(branch, node.Init)
	t.instrumentLiterals(node.Assign, branch)
	addStatementNames(branch, node.Assign)
	t.instrumentCases(node.Body.List, branch)
}

func (t *traceInstrumenter) instrumentSelect(node *ast.SelectStmt, env map[string]bool) {
	for _, clause := range node.Body.List {
		comm := clause.(*ast.CommClause)
		t.instrumentLiterals(comm.Comm, env)
		caseEnv := cloneNames(env)
		addStatementNames(caseEnv, comm.Comm)
		comm.Body = t.instrumentList(comm.Body, caseEnv)
	}
}

func (t *traceInstrumenter) instrumentCases(clauses []ast.Stmt, env map[string]bool) {
	for _, statement := range clauses {
		clause := statement.(*ast.CaseClause)
		caseEnv := cloneNames(env)
		for _, expr := range clause.List {
			t.instrumentLiterals(expr, caseEnv)
		}
		clause.Body = t.instrumentList(clause.Body, caseEnv)
	}
}

func (t *traceInstrumenter) instrumentLiterals(node ast.Node, env map[string]bool) {
	if node == nil {
		return
	}
	ast.Inspect(node, func(child ast.Node) bool {
		if child == nil {
			return false
		}
		switch nested := child.(type) {
		case *ast.FuncLit:
			t.literal(nested, env)
			return false
		case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause:
			return false
		default:
			return true
		}
	})
}

func (t *traceInstrumenter) recordCall(pos token.Pos, env map[string]bool) ast.Stmt {
	return t.recordAt(t.line(pos), env)
}

func (t *traceInstrumenter) recordAt(line int, env map[string]bool) ast.Stmt {
	names := sortedNames(env)
	omitted := max(0, len(names)-traceVarLimit)
	args := []ast.Expr{t.pointAt(line, "line", omitted, 0)}
	args = append(args, traceVariableArgs(names, traceVarLimit)...)
	return &ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: t.prefix + "Record"}, Args: args}}
}

func (t *traceInstrumenter) enterAt(line int, env map[string]bool) ast.Stmt {
	names := sortedNames(env)
	omitted := max(0, len(names)-traceVarLimit)
	args := []ast.Expr{t.pointAt(line, "call", omitted, 0)}
	args = append(args, traceVariableArgs(names, traceVarLimit)...)
	return &ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: t.prefix + "Enter"}, Args: args}}
}

func (t *traceInstrumenter) pointAt(line int, event string, omitted, variableCount int) ast.Expr {
	fields := []struct {
		name       string
		expression ast.Expr
	}{
		{"File", stringLiteral(t.file)},
		{"Line", &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(line)}},
		{"Function", stringLiteral(t.functionName)},
		{"Event", stringLiteral(event)},
		{"Omitted", &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(omitted)}},
		{"VariableCount", &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(variableCount)}},
	}
	point := &ast.CompositeLit{Type: &ast.Ident{Name: t.prefix + "Point"}}
	for _, field := range fields {
		point.Elts = append(point.Elts, &ast.KeyValueExpr{Key: &ast.Ident{Name: field.name}, Value: field.expression})
	}
	return point
}

func (t *traceInstrumenter) exitCall() ast.Stmt {
	return &ast.DeferStmt{Call: &ast.CallExpr{Fun: &ast.Ident{Name: t.prefix + "Exit"}}}
}

func (t *traceInstrumenter) implicitReturnAt(line int) ast.Stmt {
	args := []ast.Expr{t.pointAt(line, "return", 0, 0)}
	for _, name := range t.returnNames {
		args = append(args, stringLiteral(name), &ast.Ident{Name: name})
	}
	return &ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: t.prefix + "Return"}, Args: args}}
}

func (t *traceInstrumenter) appendImplicitReturn(block *ast.BlockStmt, line int) {
	if t.returnCount > 0 && len(t.returnNames) == 0 {
		return
	}
	block.List = append(block.List, t.implicitReturnAt(line), &ast.ReturnStmt{})
}

func (t *traceInstrumenter) instrumentReturn(node *ast.ReturnStmt, env map[string]bool) ast.Stmt {
	args := t.returnRecorderArgs(node, env)
	if len(node.Results) == 0 {
		args = append(args, namedReturnArgs(t.returnNames)...)
		return t.recordThenReturn(node, args)
	}
	if len(node.Results) == 1 && t.returnCount > 1 {
		return t.recordThenReturn(node, args)
	}
	return t.recordReturnValues(node, args)
}

func (t *traceInstrumenter) returnRecorderArgs(node *ast.ReturnStmt, env map[string]bool) []ast.Expr {
	locals := sortedNames(env)
	localArgs := traceVariableArgs(locals, traceVarLimit)
	args := []ast.Expr{t.pointAt(t.line(node.Pos()), "return", max(0, len(locals)-traceVarLimit), len(localArgs)/2)}
	return append(args, localArgs...)
}

func namedReturnArgs(names []string) []ast.Expr {
	args := make([]ast.Expr, 0, len(names)*2)
	for _, name := range names {
		args = append(args, stringLiteral(name), &ast.Ident{Name: name})
	}
	return args
}

func (t *traceInstrumenter) recordThenReturn(node *ast.ReturnStmt, args []ast.Expr) ast.Stmt {
	return &ast.BlockStmt{List: []ast.Stmt{t.returnCall(args), node}}
}

func (t *traceInstrumenter) returnCall(args []ast.Expr) ast.Stmt {
	return &ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: t.prefix + "Return"}, Args: args}}
}

func (t *traceInstrumenter) recordReturnValues(node *ast.ReturnStmt, args []ast.Expr) ast.Stmt {
	t.returnID++
	assign, returnVars := t.returnAssignments(node.Results)
	args = append(args, returnRecorderArgs(returnVars)...)
	returnStmt := returnStatement(returnVars)
	return &ast.BlockStmt{List: []ast.Stmt{
		assign,
		&ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: t.prefix + "Return"}, Args: args}},
		returnStmt,
	}}
}

func (t *traceInstrumenter) returnAssignments(expressions []ast.Expr) (*ast.AssignStmt, []*ast.Ident) {
	assign := &ast.AssignStmt{Tok: token.DEFINE}
	returnVars := make([]*ast.Ident, len(expressions))
	for index, expression := range expressions {
		returnVar := &ast.Ident{Name: t.prefix + "Return" + strconv.Itoa(t.returnID) + "Value" + strconv.Itoa(index)}
		returnVars[index] = returnVar
		assign.Lhs = append(assign.Lhs, returnVar)
		assign.Rhs = append(assign.Rhs, expression)
	}
	return assign, returnVars
}

func returnRecorderArgs(returnVars []*ast.Ident) []ast.Expr {
	args := make([]ast.Expr, 0, len(returnVars)*2)
	for index, returnVar := range returnVars {
		args = append(args, stringLiteral("return"+strconv.Itoa(index+1)), returnVar)
	}
	return args
}

func returnStatement(returnVars []*ast.Ident) *ast.ReturnStmt {
	statement := &ast.ReturnStmt{}
	for _, returnVar := range returnVars {
		statement.Results = append(statement.Results, returnVar)
	}
	return statement
}

func traceVariableArgs(names []string, limit int) []ast.Expr {
	args := make([]ast.Expr, 0, min(len(names), limit)*2)
	for _, name := range names[:min(len(names), limit)] {
		args = append(args, stringLiteral(name), &ast.Ident{Name: name})
	}
	return args
}

func namedResultNames(fields *ast.FieldList) []string {
	if fields == nil {
		return nil
	}
	var names []string
	for _, field := range fields.List {
		for _, name := range field.Names {
			if name.Name != "_" {
				names = append(names, name.Name)
			}
		}
	}
	return names
}

func resultFieldCount(fields *ast.FieldList) int {
	if fields == nil {
		return 0
	}
	count := 0
	for _, field := range fields.List {
		if len(field.Names) == 0 {
			count++
		} else {
			count += len(field.Names)
		}
	}
	return count
}

func (t *traceInstrumenter) line(pos token.Pos) int {
	return t.fset.PositionFor(pos, false).Line
}

func stringLiteral(literalText string) ast.Expr {
	return &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(literalText)}
}

func localNames(fields ...*ast.FieldList) map[string]bool {
	env := map[string]bool{}
	for _, list := range fields {
		addFields(env, list)
	}
	return env
}

func addFields(env map[string]bool, fields *ast.FieldList) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		for _, name := range field.Names {
			if name.Name != "_" {
				env[name.Name] = true
			}
		}
	}
}

func cloneNames(names map[string]bool) map[string]bool {
	clone := make(map[string]bool, len(names))
	for name := range names {
		clone[name] = true
	}
	return clone
}

func addStatementNames(env map[string]bool, statement ast.Stmt) {
	switch node := statement.(type) {
	case *ast.DeclStmt:
		addVariableDeclarationNames(env, node)
	case *ast.AssignStmt:
		addShortDeclarationNames(env, node)
	}
}

func addVariableDeclarationNames(env map[string]bool, statement *ast.DeclStmt) {
	decl, ok := statement.Decl.(*ast.GenDecl)
	if !ok || decl.Tok != token.VAR {
		return
	}
	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range valueSpec.Names {
				addIdentifier(env, name)
			}
		}
	}
}

func addShortDeclarationNames(env map[string]bool, statement *ast.AssignStmt) {
	if statement.Tok != token.DEFINE {
		return
	}
	for _, lhs := range statement.Lhs {
		addIdentifier(env, lhs)
	}
}

func addIdentifier(env map[string]bool, expr ast.Expr) {
	if id, ok := expr.(*ast.Ident); ok && id.Name != "_" {
		env[id.Name] = true
	}
}

func sortedNames(env map[string]bool) []string {
	names := make([]string, 0, len(env))
	for name := range env {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

const traceHelperTemplate = `package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type __trace__Frame struct {
	File string ` + "`json:\"file\"`" + `
	Line int ` + "`json:\"line\"`" + `
	Func string ` + "`json:\"func\"`" + `
	Depth int ` + "`json:\"depth\"`" + `
	Event string ` + "`json:\"event\"`" + `
	Goroutine int ` + "`json:\"goroutine\"`" + `
	Variables []__trace__Variable ` + "`json:\"variables\"`" + `
	Returns []__trace__Variable ` + "`json:\"returns,omitempty\"`" + `
	Stdout int64 ` + "`json:\"stdout\"`" + `
	Sim int64 ` + "`json:\"sim,omitempty\"`" + `
	Truncated bool ` + "`json:\"truncated,omitempty\"`" + `
	Paused bool ` + "`json:\"paused,omitempty\"`" + `
}

type __trace__Variable struct {
	Name string ` + "`json:\"name\"`" + `
	Formatted string ` + "`json:\"value\"`" + `
}

type __trace__Point struct {
	File string
	Line int
	Function string
	Event string
	Omitted int
	VariableCount int
}

type __trace__Thread struct {
	Goroutine int
	Depth int
}

var __trace__Stopped atomic.Bool
var __trace__Mu sync.Mutex
var __trace__Live bool
var __trace__File = __trace__Open()
var __trace__Frames int
var __trace__Bytes int
var __trace__DepthMu sync.Mutex
var __trace__Depths = map[int]int{}
var __trace__Stdout int64
var __trace__SimPath = os.Getenv("ORBITAL_TRACE")

func __trace__Record(point __trace__Point, namedLocals ...any) {
	if __trace__Stopped.Load() || __trace__File == nil {
		return
	}
	__trace__RecordAt(point, __trace__CurrentThread(), namedLocals...)
}

func __trace__Enter(point __trace__Point, namedLocals ...any) {
	if __trace__Stopped.Load() || __trace__File == nil {
		return
	}
	point.Event = "call"
	__trace__RecordAt(point, __trace__EnterThread(), namedLocals...)
}

func __trace__EnterThread() __trace__Thread {
	id := __trace__Goroutine()
	__trace__DepthMu.Lock()
	depth, exists := __trace__Depths[id]
	if exists {
		depth++
	}
	__trace__Depths[id] = depth
	__trace__DepthMu.Unlock()
	return __trace__Thread{Goroutine: id, Depth: depth}
}

func __trace__Exit() {
	id := __trace__Goroutine()
	__trace__DepthMu.Lock()
	if depth, exists := __trace__Depths[id]; exists {
		if depth == 0 {
			delete(__trace__Depths, id)
		} else {
			__trace__Depths[id] = depth - 1
		}
	}
	__trace__DepthMu.Unlock()
}

func __trace__Return(point __trace__Point, namedValues ...any) {
	if __trace__Stopped.Load() || __trace__File == nil {
		return
	}
	thread := __trace__CurrentThread()
	if !__trace__Acquire(thread.Goroutine) {
		return
	}
	defer __trace__Release(thread.Goroutine)
	variableValues := namedValues[:min(len(namedValues), point.VariableCount*2)]
	returnValues := namedValues[len(variableValues):]
	returnLimit := min(len(returnValues), __TRACE_VAR_LIMIT__*2)
	vars, truncated := __trace__Variables(variableValues, point.Omitted)
	returns, returnTruncated := __trace__Variables(returnValues[:returnLimit], 0)
	returnTruncated = returnTruncated || returnLimit < len(returnValues)
	__trace__WriteFrame(__trace__Frame{File: point.File, Line: point.Line, Func: point.Function, Depth: thread.Depth, Event: "return", Goroutine: thread.Goroutine, Variables: vars, Returns: returns, Truncated: truncated || returnTruncated})
}


func __trace__CurrentThread() __trace__Thread {
	id := __trace__Goroutine()
	__trace__DepthMu.Lock()
	depth := __trace__Depths[id]
	__trace__DepthMu.Unlock()
	return __trace__Thread{Goroutine: id, Depth: depth}
}


// Formatting a value can call the learner's own String/Error methods, which
// are instrumented too; recording again from inside would recurse forever.
// While a goroutine is formatting, nested record calls are skipped.
var __trace__BusyMu sync.Mutex
var __trace__Busy = map[int]bool{}

func __trace__Acquire(id int) bool {
	__trace__BusyMu.Lock()
	defer __trace__BusyMu.Unlock()
	if __trace__Busy[id] {
		return false
	}
	__trace__Busy[id] = true
	return true
}

func __trace__Release(id int) {
	__trace__BusyMu.Lock()
	delete(__trace__Busy, id)
	__trace__BusyMu.Unlock()
}

func __trace__RecordAt(point __trace__Point, thread __trace__Thread, namedLocals ...any) {
	if !__trace__Acquire(thread.Goroutine) {
		return
	}
	vars, truncated := __trace__Variables(namedLocals, point.Omitted)
	__trace__Release(thread.Goroutine)
	frame := __trace__Frame{File: point.File, Line: point.Line, Func: point.Function, Depth: thread.Depth, Event: point.Event, Goroutine: thread.Goroutine, Variables: vars, Truncated: truncated}
	__trace__WriteFrame(frame)
}

func __trace__WriteFrame(frame __trace__Frame) {
	frame.Stdout = atomic.LoadInt64(&__trace__Stdout)
	if __trace__Live {
		__trace__Gate(&frame)
		return
	}
	if __trace__SimPath != "" {
		// The colony writes each world event straight to this file, so its
		// size says how many events happened before this statement.
		if info, err := os.Stat(__trace__SimPath); err == nil {
			frame.Sim = info.Size()
		}
	}
	encoded, err := json.Marshal(frame)
	if err == nil {
		__trace__Append(encoded)
	}
}

func __trace__Variables(namedLocals []any, omitted int) ([]__trace__Variable, bool) {
	vars := make([]__trace__Variable, 0, len(namedLocals)/2)
	truncated := omitted > 0
	for i := 0; i+1 < len(namedLocals); i += 2 {
		name, ok := namedLocals[i].(string)
		if !ok {
			continue
		}
		if len(name) > __TRACE_NAME_LIMIT__ {
			name = __trace__Prefix(name, __TRACE_NAME_PREFIX_LIMIT__) + "…"
			truncated = true
		}
		formattedValue, cut := __trace__Format(namedLocals[i+1])
		vars = append(vars, __trace__Variable{Name: name, Formatted: formattedValue})
		truncated = truncated || cut
	}
	return vars, truncated
}

func __trace__Append(encoded []byte) {
	__trace__Mu.Lock()
	defer __trace__Mu.Unlock()
	if __trace__Stopped.Load() {
		return
	}
	if __trace__File == nil {
		__trace__Stopped.Store(true)
		return
	}
	if __trace__Frames >= __TRACE_FRAME_LIMIT__ || __trace__Bytes+len(encoded)+1 > __TRACE_BYTE_LIMIT__-32 {
		__trace__Stop()
		return
	}
	if _, err := __trace__File.Write(append(encoded, '\n')); err != nil {
		__trace__Stopped.Store(true)
		return
	}
	__trace__Frames++
	__trace__Bytes += len(encoded) + 1
}

func __trace__Open() *os.File {
	if fd, err := strconv.Atoi(os.Getenv("ORBITAL_DEBUG_FD")); err == nil && fd > 2 {
		if file := os.NewFile(uintptr(fd), "debug"); file != nil {
			__trace__Live = true
			if ctl, err := strconv.Atoi(os.Getenv("ORBITAL_DEBUG_CTL_FD")); err == nil && ctl > 2 {
				if control := os.NewFile(uintptr(ctl), "debug-control"); control != nil {
					go __trace__Control(control)
				}
			}
			return file
		}
	}
	path := os.Getenv("ORBITAL_USER_TRACE")
	if path == "" {
		__trace__Stopped.Store(true)
		return nil
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		__trace__Stopped.Store(true)
		return nil
	}
	return file
}

func __trace__Stop() {
	_, _ = io.WriteString(__trace__File, "{\"truncated\":true}\n")
	__trace__Stopped.Store(true)
}

// ---- live debugging (colony): frames stream out, commands come in ----

const (
	__trace__Running = iota
	__trace__PauseRequested
	__trace__Paused
	__trace__Stepping
)

var __trace__GateMu sync.Mutex
var __trace__GateCond = sync.NewCond(&__trace__GateMu)
var __trace__GateState = __trace__Running
var __trace__StepKind string
var __trace__StepGoroutine, __trace__StepDepth int
var __trace__Breaks = map[string]bool{}
var __trace__PausedAt __trace__Frame
var __trace__Window int64
var __trace__WindowCount int64

// Gate decides whether this statement pauses the program; paused
// goroutines block until a command arrives, and while one goroutine is
// paused every other instrumented goroutine waits at its next statement.
func __trace__Gate(frame *__trace__Frame) {
	__trace__GateMu.Lock()
	for __trace__GateState == __trace__Paused {
		__trace__GateCond.Wait()
	}
	if !__trace__ShouldStop(frame) {
		__trace__GateMu.Unlock()
		__trace__Emit(frame)
		return
	}
	__trace__GateState = __trace__Paused
	__trace__PausedAt = *frame
	frame.Paused = true
	__trace__GateMu.Unlock()
	__trace__Emit(frame)
	__trace__GateMu.Lock()
	for __trace__GateState == __trace__Paused {
		__trace__GateCond.Wait()
	}
	__trace__GateMu.Unlock()
}

func __trace__ShouldStop(frame *__trace__Frame) bool {
	switch __trace__GateState {
	case __trace__PauseRequested:
		return true
	case __trace__Stepping:
		same := frame.Goroutine == __trace__StepGoroutine
		switch __trace__StepKind {
		case "into":
			return true
		case "over":
			if same && frame.Depth <= __trace__StepDepth {
				return true
			}
		case "out":
			if same && frame.Depth < __trace__StepDepth {
				return true
			}
		}
	}
	return frame.Event == "line" && __trace__Breaks[frame.File+":"+strconv.Itoa(frame.Line)]
}

// Emit sends a frame. Running frames are sampled (a busy loop must not
// flood the pipe or slow the colony); paused frames always go out.
func __trace__Emit(frame *__trace__Frame) {
	if !frame.Paused {
		now := time.Now().Unix()
		if atomic.LoadInt64(&__trace__Window) != now {
			atomic.StoreInt64(&__trace__Window, now)
			atomic.StoreInt64(&__trace__WindowCount, 0)
		}
		if atomic.AddInt64(&__trace__WindowCount, 1) > 120 {
			return
		}
	}
	encoded, err := json.Marshal(frame)
	if err != nil {
		return
	}
	__trace__Mu.Lock()
	_, _ = __trace__File.Write(append(encoded, '\n'))
	__trace__Mu.Unlock()
}

func __trace__Control(control *os.File) {
	scanner := bufio.NewScanner(control)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		__trace__GateMu.Lock()
		switch fields[0] {
		case "pause":
			if __trace__GateState == __trace__Running {
				__trace__GateState = __trace__PauseRequested
			}
		case "cont":
			__trace__GateState = __trace__Running
		case "into", "over", "out":
			if __trace__GateState == __trace__Paused {
				__trace__StepKind = fields[0]
				__trace__StepGoroutine, __trace__StepDepth = __trace__PausedAt.Goroutine, __trace__PausedAt.Depth
				__trace__GateState = __trace__Stepping
			}
		case "breaks":
			__trace__Breaks = map[string]bool{}
			for _, key := range fields[1:] {
				__trace__Breaks[key] = true
			}
		}
		__trace__GateCond.Broadcast()
		__trace__GateMu.Unlock()
	}
	// the debugger went away: never leave the program paused
	__trace__GateMu.Lock()
	__trace__GateState = __trace__Running
	__trace__Breaks = map[string]bool{}
	__trace__GateCond.Broadcast()
	__trace__GateMu.Unlock()
}

type __trace__Writer struct {
	text strings.Builder
	cut bool
}

func (w *__trace__Writer) Write(part []byte) (int, error) {
	n := len(part)
	remaining := __TRACE_VALUE_LIMIT__ - w.text.Len()
	if remaining <= 0 {
		w.cut = true
		return n, io.ErrShortWrite
	}
	if n > remaining {
		_, _ = w.text.Write(part[:remaining])
		w.cut = true
		return n, io.ErrShortWrite
	}
	return w.text.Write(part)
}

func __trace__Format(localValue any) (formatted string, cut bool) {
	writer := &__trace__Writer{}
	defer func() {
		if recover() != nil {
			formatted = "<format panic>"
			cut = true
		}
	}()
	_, _ = fmt.Fprintf(writer, "%v", localValue)
	formatted = strings.ToValidUTF8(writer.text.String(), "�")
	cut = writer.cut
	if cut {
		formatted = __trace__Prefix(formatted, __TRACE_VALUE_PREFIX_LIMIT__) + "…"
	}
	return formatted, cut
}

func __trace__Prefix(prefixText string, limit int) string {
	if len(prefixText) <= limit {
		return prefixText
	}
	end := limit
	for end > 0 && end < len(prefixText) && (prefixText[end]&0xc0) == 0x80 {
		end--
	}
	return prefixText[:end]
}

func __trace__Goroutine() int {
	stack := make([]byte, 64)
	n := runtime.Stack(stack, false)
	fields := strings.Fields(string(stack[:n]))
	if len(fields) < 2 {
		return 0
	}
	id, _ := strconv.Atoi(fields[1])
	return id
}

func __trace__WriteOutput(writer io.Writer, formatted string) (int, error) {
	n, err := io.WriteString(writer, formatted)
	if output, ok := writer.(*os.File); ok && output == os.Stdout {
		atomic.AddInt64(&__trace__Stdout, int64(n))
	}
	return n, err
}

func __trace__WriteStdout(writer *os.File, outputBytes []byte) (int, error) {
	n, err := writer.Write(outputBytes)
	if writer == os.Stdout {
		atomic.AddInt64(&__trace__Stdout, int64(n))
	}
	return n, err
}

func __trace__WriteStdoutString(writer *os.File, outputText string) (int, error) {
	return __trace__WriteOutput(writer, outputText)
}

func __trace__WriteString(writer io.Writer, outputText string) (int, error) {
	return __trace__WriteOutput(writer, outputText)
}

func __trace__Print(arguments ...any) (int, error) {
	return __trace__WriteOutput(os.Stdout, fmt.Sprint(arguments...))
}

func __trace__Println(arguments ...any) (int, error) {
	return __trace__WriteOutput(os.Stdout, fmt.Sprintln(arguments...))
}

func __trace__Printf(format string, arguments ...any) (int, error) {
	return __trace__WriteOutput(os.Stdout, fmt.Sprintf(format, arguments...))
}

func __trace__Fprint(writer io.Writer, arguments ...any) (int, error) {
	return __trace__WriteOutput(writer, fmt.Sprint(arguments...))
}

func __trace__Fprintln(writer io.Writer, arguments ...any) (int, error) {
	return __trace__WriteOutput(writer, fmt.Sprintln(arguments...))
}

func __trace__Fprintf(writer io.Writer, format string, arguments ...any) (int, error) {
	return __trace__WriteOutput(writer, fmt.Sprintf(format, arguments...))
}
`

func traceHelperSource(prefix string) string {
	return strings.NewReplacer(
		"__trace__", prefix,
		"__TRACE_FRAME_LIMIT__", strconv.Itoa(traceFrameLimit),
		"__TRACE_BYTE_LIMIT__", strconv.Itoa(traceByteLimit),
		"__TRACE_VALUE_LIMIT__", strconv.Itoa(traceValueLimit),
		"__TRACE_VALUE_PREFIX_LIMIT__", strconv.Itoa(traceValueLimit-3),
		"__TRACE_NAME_LIMIT__", strconv.Itoa(traceNameLimit),
		"__TRACE_NAME_PREFIX_LIMIT__", strconv.Itoa(traceNameLimit-3),
		"__TRACE_VAR_LIMIT__", strconv.Itoa(traceVarLimit),
	).Replace(traceHelperTemplate)
}
