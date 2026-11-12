package main

import (
	"strings"

	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strconv"
)

func requirementTexts(lang string) map[string]string {
	return map[string]string{
		"rune":      tr(lang, "استخدم []rune أو range على نص أو utf8.", "Use []rune, range over a string, or utf8."),
		"iota":      tr(lang, "استخدم iota جوه const.", "Use iota inside const."),
		"const":     tr(lang, "عرّف ثابت بـ const.", "Declare a const."),
		"defer":     tr(lang, "استخدم defer.", "Use defer."),
		"recover":   tr(lang, "استخدم recover جوه دالة مؤجلة بـ defer.", "Call recover inside a deferred function literal."),
		"errwrap":   tr(lang, "لف الخطأ بـ %w وافحصه بـ errors.Is أو errors.As.", "Wrap an error with %w and inspect it with errors.Is or errors.As."),
		"errtype":   tr(lang, "عرّف نوع عليه Error() string.", "Define a type with Error() string."),
		"funcvalue": tr(lang, "استخدم دالة كقيمة.", "Use a function as a value."),
		"closure":   tr(lang, "رجّع دالة داخلية من دالة.", "Return a function literal from a function."),
		"embed":     tr(lang, "استخدم field مضمّن من غير اسم.", "Use an embedded struct field."),
		"variadic":  tr(lang, "استخدم parameter بـ ... .", "Use a variadic parameter."),
		"bufchan":   tr(lang, "اعمل channel بسعة محددة.", "Make a buffered channel."),
		"close":     tr(lang, "اقفل channel بـ close.", "Close a channel with close."),
		"recursion": tr(lang, "اكتب دالة بتنادي نفسها.", "Write a recursive function."),
		"fuzz":      tr(lang, "اكتب FuzzXxx بياخد *testing.F.", "Write FuzzXxx taking *testing.F."),
		"subtest":   tr(lang, "استخدم t.Run لتست فرعي.", "Use t.Run for a subtest."),
		"iface":     tr(lang, "عرّف interface بنفسك.", "Define your own interface."),

		"for":         tr(lang, "المهمة دي لازم تستخدم حلقة for.", "This mission requires a for loop."),
		"for2":        tr(lang, "المهمة دي محتاجة حلقتين for (واحدة جوه التانية).", "This mission needs two nested for loops."),
		"if":          tr(lang, "المهمة دي لازم تستخدم if.", "This mission requires if."),
		"switch":      tr(lang, "المهمة دي لازم تستخدم switch.", "This mission requires switch."),
		"func":        tr(lang, "اكتب دالة بنفسك (غير main) ونادِها مرتين على الأقل.", "Write your own function (other than main) and call it at least twice."),
		"errcheck":    tr(lang, "لازم تتعامل مع الخطأ: if err != nil أو errors.Is.", "Handle errors using if err != nil or errors.Is."),
		"slice":       tr(lang, "المهمة دي محتاجة slice: استخدم append أو make.", "This mission needs a slice: use append or make."),
		"range":       tr(lang, "المهمة دي لازم تستخدم for ... range.", "This mission requires for ... range."),
		"map":         tr(lang, "المهمة دي محتاجة map.", "This mission needs a map."),
		"struct":      tr(lang, "عرّف نوع struct بنفسك (type X struct { ... }).", "Define your own struct type (type X struct { ... })."),
		"method":      tr(lang, "اكتب method على نوع إنت عامله: func (x *X) Name() { ... }.", "Write a method on your own type: func (x *X) Name() { ... }."),
		"pointer":     tr(lang, "استخدم pointer: parameter نوعه *int وابعتله &متغير.", "Use a pointer: a parameter of type *int, passed the address of a variable."),
		"ptrrecv":     tr(lang, "محتاج method بـ pointer receiver: func (d *Drone) ...", "Use a method with a pointer receiver: func (d *Drone) ..."),
		"nilcheck":    tr(lang, "اتأكد إن الـ pointer مش nil قبل ما تستخدمه (p != nil).", "Check that the pointer is not nil before using it (p != nil)."),
		"stringer":    tr(lang, "اعمل method اسمها String() string على النوع بتاعك.", "Add a String() string method to your type."),
		"autopilot":   tr(lang, "المهمة دي لازم تستخدم colony.Autopilot.", "This mission requires colony.Autopilot."),
		"typeswitch":  tr(lang, "المهمة دي لازم تستخدم type switch: switch v := s.(type) { ... }.", "This mission requires a type switch: switch v := s.(type) { ... }."),
		"go":          tr(lang, "المهمة دي محتاجة goroutine: go f().", "This mission needs a goroutine: go f()."),
		"waitgroup":   tr(lang, "استخدم sync.WaitGroup عشان main تستنى الـ goroutines.", "Use sync.WaitGroup so main waits for the goroutines."),
		"chan":        tr(lang, "المهمة دي محتاجة channel: make(chan T).", "This mission needs a channel: make(chan T)."),
		"select":      tr(lang, "المهمة دي لازم تستخدم select.", "This mission requires select."),
		"mutex":       tr(lang, "احمي البيانات المشتركة بـ sync.Mutex (أو sync.RWMutex).", "Protect shared data with sync.Mutex (or sync.RWMutex)."),
		"generic":     tr(lang, "اكتب دالة generic: func Name[T ...](...).", "Write a generic function: func Name[T ...](...)."),
		"generictype": tr(lang, "اعمل type generic: type Name[T any] struct { ... }.", "Create a generic type: type Name[T any] struct { ... }."),
		"benchmark":   tr(lang, "اكتب benchmark بنفسك: func BenchmarkXxx(b *testing.B).", "Write your own benchmark: func BenchmarkXxx(b *testing.B)."),
		"tabletest":   tr(lang, "اكتب table-driven test: slice من الحالات وrange عليها.", "Write a table-driven test: a slice of cases iterated with range."),
	}
}

type features struct {
	loops, ifs, switches, errChecks        int
	ranges, slices, maps, structs, methods int
	ptrParams, addrOf, ptrRecv, stringers  int
	typeSwitches, goStmts, selects, chans  int
	genericFuncs, genericTypes, benchmarks int
	extra                                  map[string]bool
	sels                                   map[string]int // every pkg.Name selector
	calls                                  map[string]int // "pkg.Name" -> count
	userFuncs                              map[string]int // name -> call count
}

func scan(code string) (*features, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", code, 0)
	if err != nil {
		return nil, err
	}
	ft := &features{userFuncs: map[string]int{}, calls: map[string]int{}, sels: map[string]int{}}
	imports := map[string]string{}
	for _, spec := range f.Imports {
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, err
		}
		name := importName(importPath)
		if spec.Name != nil {
			name = spec.Name.Name
		}
		imports[name] = importPath
	}
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok {
			if fd.Recv != nil {
				ft.methods++
				if _, ok := fd.Recv.List[0].Type.(*ast.StarExpr); ok {
					ft.ptrRecv++
				}
				if fd.Name.Name == "String" && fd.Type.Params.NumFields() == 0 && fd.Type.Results.NumFields() == 1 {
					ft.stringers++
				}
			} else if fd.Name.Name != "main" && fd.Name.Name != "init" {
				ft.userFuncs[fd.Name.Name] = 0
			}
			if fd.Type.TypeParams.NumFields() > 0 {
				ft.genericFuncs++
			}
			if strings.HasPrefix(fd.Name.Name, "Benchmark") {
				ft.benchmarks++
			}
			for _, p := range fd.Type.Params.List {
				if _, ok := p.Type.(*ast.StarExpr); ok {
					ft.ptrParams++
				}
			}
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ForStmt:
			ft.loops++
		case *ast.RangeStmt:
			ft.loops++
			ft.ranges++
		case *ast.MapType:
			ft.maps++
		case *ast.TypeSpec:
			if _, ok := n.Type.(*ast.StructType); ok {
				ft.structs++
			}
			if n.TypeParams.NumFields() > 0 {
				ft.genericTypes++
			}
		case *ast.IfStmt:
			ft.ifs++
		case *ast.SwitchStmt:
			ft.switches++
		case *ast.TypeSwitchStmt:
			ft.switches++
			ft.typeSwitches++
		case *ast.GoStmt:
			ft.goStmts++
		case *ast.SelectStmt:
			ft.selects++
		case *ast.ChanType:
			ft.chans++
		case *ast.SelectorExpr:
			if pkg, ok := n.X.(*ast.Ident); ok {
				ft.sels[packageName(imports, pkg)+"."+n.Sel.Name]++
			}
		case *ast.UnaryExpr:
			if n.Op == token.AND {
				ft.addrOf++
			}
		case *ast.CallExpr:
			switch fn := n.Fun.(type) {
			case *ast.Ident:
				if fn.Name == "append" || fn.Name == "make" && len(n.Args) > 0 && isSlice(n.Args[0]) {
					ft.slices++
				}
				if c, ok := ft.userFuncs[fn.Name]; ok {
					ft.userFuncs[fn.Name] = c + 1
				}
			case *ast.SelectorExpr:
				if pkg, ok := fn.X.(*ast.Ident); ok {
					name := packageName(imports, pkg)
					ft.calls[name+"."+fn.Sel.Name]++
					if name == "errors" && (fn.Sel.Name == "Is" || fn.Sel.Name == "As") {
						ft.errChecks++
					}
				}
			}
		case *ast.BinaryExpr:
			if n.Op == token.NEQ || n.Op == token.EQL {
				if id, ok := n.Y.(*ast.Ident); ok && id.Name == "nil" {
					ft.errChecks++
				}
			}
		}
		return true
	})
	ft.extra = scanP0(f, imports)
	return ft, nil
}

func packageName(imports map[string]string, id *ast.Ident) string {
	if importPath, ok := imports[id.Name]; ok && id.Obj == nil {
		return importName(importPath)
	}
	return ""
}

// These pinned dependencies declare package names that differ from their paths.
var importNames = map[string]string{
	"github.com/redis/go-redis/v9":  "redis",
	"github.com/segmentio/kafka-go": "kafka",
}

// importName infers package names without loading dependencies or using the network.
func importName(importPath string) string {
	if name, ok := importNames[importPath]; ok {
		return name
	}
	name := path.Base(importPath)
	if dir := path.Dir(importPath); dir != "." && len(name) > 1 && name[0] == 'v' && strings.Trim(name[1:], "0123456789") == "" {
		name = path.Base(dir)
	}
	if strings.HasPrefix(importPath, "gopkg.in/") {
		if base, version, ok := strings.Cut(name, ".v"); ok && version != "" && strings.Trim(version, "0123456789") == "" {
			name = base
		}
	}
	return strings.TrimSuffix(strings.TrimSuffix(strings.TrimPrefix(name, "go-"), "-go"), ".go")
}

// checkRequirements returns a message for every requirement the code misses.
func checkRequirements(code string, reqs []string, languages ...string) []string {
	lang := "ar"
	if len(languages) > 0 {
		lang = normalizeLang(languages[0])
	}
	requirementText := requirementTexts(lang)
	if len(reqs) == 0 {
		return nil
	}
	ft, err := scan(code)
	if err != nil {
		return nil // the compiler already accepted it; don't block on parser quirks
	}
	var missing []string
	for _, r := range reqs {
		ok := true
		switch r {
		case "rune", "iota", "const", "defer", "recover", "errwrap", "errtype", "funcvalue", "closure", "embed", "variadic", "bufchan", "close", "recursion", "fuzz", "subtest", "iface":
			ok = ft.extra[r]
		case "for":
			ok = ft.loops >= 1
		case "for2":
			ok = ft.loops >= 2
		case "if":
			ok = ft.ifs >= 1
		case "switch":
			ok = ft.switches >= 1
		case "errcheck":
			ok = ft.errChecks >= 1
		case "slice":
			ok = ft.slices >= 1
		case "range":
			ok = ft.ranges >= 1
		case "map":
			ok = ft.maps >= 1
		case "struct":
			ok = ft.structs >= 1
		case "method":
			ok = ft.methods >= 1
		case "pointer":
			ok = ft.ptrParams >= 1 && ft.addrOf >= 1
		case "ptrrecv":
			ok = ft.ptrRecv >= 1
		case "nilcheck":
			ok = ft.errChecks >= 1
		case "stringer":
			ok = ft.stringers >= 1
		case "autopilot":
			ok = ft.calls["colony.Autopilot"] >= 1
		case "go":
			ok = ft.goStmts >= 1
		case "waitgroup":
			ok = ft.sels["sync.WaitGroup"] >= 1
		case "chan":
			ok = ft.chans >= 1
		case "select":
			ok = ft.selects >= 1
		case "generic":
			ok = ft.genericFuncs >= 1
		case "generictype":
			ok = ft.genericTypes >= 1
		case "benchmark":
			ok = ft.benchmarks >= 1
		case "tabletest":
			ok = ft.ranges >= 1 && ft.sels["testing.T"] >= 1
		case "mutex":
			ok = ft.sels["sync.Mutex"]+ft.sels["sync.RWMutex"] >= 1
		case "typeswitch":
			ok = ft.typeSwitches >= 1
		case "func":
			ok = false
			for _, c := range ft.userFuncs {
				if c >= 2 {
					ok = true
				}
			}
		}
		msg := requirementText[r]
		if name, isSel := strings.CutPrefix(r, "sel:"); isSel {
			ok = ft.sels[name] >= 1
			msg = tr(lang, "المهمة دي لازم تستخدم ", "This mission requires ") + name + "."
		}
		if !ok {
			missing = append(missing, msg)
		}
	}
	return missing
}

func isSlice(e ast.Expr) bool {
	at, ok := e.(*ast.ArrayType)
	return ok && at.Len == nil
}
