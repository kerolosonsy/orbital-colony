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
	traceFrameLimit = 400
	traceByteLimit  = 64 << 10
	traceValueLimit = 160
	traceNameLimit  = 96
	traceVarLimit   = 32
)

type traceVariable struct {
	Name      string `json:"name"`
	Formatted string `json:"value"`
}

type traceFrame struct {
	File      string          `json:"file"`
	Line      int             `json:"line"`
	Goroutine int             `json:"goroutine"`
	Variables []traceVariable `json:"variables"`
	Truncated bool            `json:"truncated,omitempty"`
}

type traceResponse struct {
	Available bool         `json:"available"`
	Completed bool         `json:"completed"`
	Frames    []traceFrame `json:"frames"`
	Truncated bool         `json:"truncated,omitempty"`
	Message   string       `json:"message,omitempty"`
}

type traceRequest struct {
	Level string            `json:"level"`
	Code  string            `json:"code"`
	Files map[string]string `json:"files"`
}

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
	cmd := exec.CommandContext(rctx, filepath.Join(dir, "trace-program"))
	configureProcess(cmd)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "ORBITAL_LANG="+lang, "ORBITAL_LEVEL="+lvl.ID, "ORBITAL_TRACE="+simPath, "ORBITAL_USER_TRACE="+tracePath)
	if seed, ok := ctx.Value(reviewSeedKey{}).(int64); ok {
		cmd.Env = append(cmd.Env, "ORBITAL_SEED="+strconv.FormatInt(seed, 10))
	}
	cmd.Stdout = &limitedBuffer{max: 64 << 10, lang: lang}
	cmd.Stderr = &limitedBuffer{max: 64 << 10, lang: lang}
	runErr := runCommand(cmd)
	resp := readTraceFile(traceFile, lang)
	resp.Completed = runErr == nil
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
	fset   *token.FileSet
	file   string
	prefix string
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
		instrumentTraceFile(fset, name, prefix, file)
	}
	return formatTraceFile(fset, file)
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
	env := localNames(fn.Recv, fn.Type.Params, fn.Type.Results)
	t.instrumentBlock(fn.Body, env)
}

func (t *traceInstrumenter) literal(fn *ast.FuncLit, outer map[string]bool) {
	env := cloneNames(outer)
	addFields(env, fn.Type.Params)
	addFields(env, fn.Type.Results)
	t.instrumentBlock(fn.Body, env)
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
	args := []ast.Expr{stringLiteral(t.file), &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(line)}, &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(omitted)}}
	for _, name := range names[:min(len(names), traceVarLimit)] {
		args = append(args, stringLiteral(name), &ast.Ident{Name: name})
	}
	return &ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: t.prefix + "Record"}, Args: args}}
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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type __trace__Frame struct {
	File string ` + "`json:\"file\"`" + `
	Line int ` + "`json:\"line\"`" + `
	Goroutine int ` + "`json:\"goroutine\"`" + `
	Variables []__trace__Variable ` + "`json:\"variables\"`" + `
	Truncated bool ` + "`json:\"truncated,omitempty\"`" + `
}

type __trace__Variable struct {
	Name string ` + "`json:\"name\"`" + `
	Formatted string ` + "`json:\"value\"`" + `
}

var __trace__Stopped atomic.Bool
var __trace__Mu sync.Mutex
var __trace__File = __trace__Open()
var __trace__Frames int
var __trace__Bytes int

func __trace__Record(file string, line int, omitted int, namedLocals ...any) {
	if __trace__Stopped.Load() || __trace__File == nil {
		return
	}
	vars, truncated := __trace__Variables(namedLocals, omitted)
	frame := __trace__Frame{File: file, Line: line, Goroutine: __trace__Goroutine(), Variables: vars, Truncated: truncated}
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
	).Replace(traceHelperTemplate)
}
