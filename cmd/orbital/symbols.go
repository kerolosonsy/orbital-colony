package main

import (
	"bytes"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// apiSymbol is one exported name of the learner package orbital/colony, for
// the editor's autocomplete: what to insert and the first line of its doc.
type apiSymbol struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"` // func | type | const | var | method
	Recv      string `json:"recv,omitempty"`
	Signature string `json:"signature,omitempty"`
	Doc       string `json:"doc,omitempty"`
}

var sameAsFunc = regexp.MustCompile(`^(\w+) is \w+ for (this \w+)\.$`)

var (
	symbolsOnce sync.Once
	symbolsList []apiSymbol
)

// colonySymbols reads the colony package's own source and doc comments once.
func colonySymbols(dir string) []apiSymbol {
	symbolsOnce.Do(func() {
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
			return !strings.HasSuffix(info.Name(), "_test.go")
		}, parser.ParseComments)
		if err != nil || pkgs["colony"] == nil {
			return
		}
		pkg := doc.New(pkgs["colony"], "orbital/colony", 0)
		signature := func(decl *ast.FuncDecl) string {
			var buf bytes.Buffer
			printer.Fprint(&buf, fset, decl.Type)
			return strings.TrimPrefix(buf.String(), "func")
		}
		add := func(s apiSymbol) {
			s.Doc = firstSentence(s.Doc)
			symbolsList = append(symbolsList, s)
		}
		// values lists constants/vars; ones grouped under a type show that type
		// (e.g. East is a Direction), which says more than an empty doc.
		values := func(kind string, groups []*doc.Value, typeName string) {
			for _, group := range groups {
				for _, name := range group.Names {
					if ast.IsExported(name) {
						add(apiSymbol{Name: name, Kind: kind, Signature: typeName, Doc: group.Doc})
					}
				}
			}
		}
		for _, f := range pkg.Funcs {
			add(apiSymbol{Name: f.Name, Kind: "func", Signature: signature(f.Decl), Doc: f.Doc})
		}
		values("const", pkg.Consts, "")
		values("var", pkg.Vars, "")
		for _, t := range pkg.Types {
			add(apiSymbol{Name: t.Name, Kind: "type", Doc: t.Doc})
			values("const", t.Consts, t.Name)
			values("var", t.Vars, t.Name)
			for _, f := range t.Funcs { // constructors such as Spawn() *Drone
				add(apiSymbol{Name: f.Name, Kind: "func", Signature: signature(f.Decl), Doc: f.Doc})
			}
			for _, m := range t.Methods {
				add(apiSymbol{Name: m.Name, Kind: "method", Recv: t.Name, Signature: signature(m.Decl), Doc: m.Doc})
			}
		}
		// "Plant is Plant for this drone." says little in a hint: show the
		// package function's doc for that method instead.
		funcDocs := map[string]string{}
		for _, s := range symbolsList {
			if s.Kind == "func" {
				funcDocs[s.Name] = s.Doc
			}
		}
		for i, s := range symbolsList {
			if s.Kind != "method" {
				continue
			}
			if m := sameAsFunc.FindStringSubmatch(s.Doc); m != nil && funcDocs[m[1]] != "" {
				symbolsList[i].Doc = funcDocs[m[1]] + " (" + m[2] + ")"
			}
		}
		sort.SliceStable(symbolsList, func(i, j int) bool { return symbolsList[i].Name < symbolsList[j].Name })
	})
	return symbolsList
}

func firstSentence(text string) string {
	text = strings.Join(strings.Fields(text), " ")
	if i := strings.Index(text, ". "); i >= 0 {
		return text[:i+1]
	}
	return text
}

// symbols: GET /api/symbols → the colony package's exported API.
func (s *server) symbols(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	writeJSON(w, colonySymbols(s.run.colonyDir))
}
