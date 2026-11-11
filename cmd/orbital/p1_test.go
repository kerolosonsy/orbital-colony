package main

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"orbital/colony/levels"
)

// These plausible mistakes still use the required concepts; the mission's
// state checks must reject them independently of the AST requirements.
func TestP1RejectsIncorrectResults(t *testing.T) {
	r := testRunner(t)
	cases := []struct{ id, from, to string }{
		{"2-3-callsigns", "utf8.RuneCountInString(firstName)", "len(firstName)"},
		{"2-4-units", "crates*CrateFactor", "crates*PackFactor"},
		{"4-3-checkpoint", "n < 20", "n <= 20"},
		{"5-3-fault-log", "all = errors.Join(all, wrapped)", "all = errors.Join(nil, wrapped)"},
		{"5-4-undock", "defer colony.Undock()", "defer func() {}()"},
		{"5-5-command-table", "colony.Wait()", ""},
		{"6-4-stock-order", "return a - b", "return b - a"},
		{"7-4-composite-units", "colony.Report(unit.Shares())", "colony.Report(append(unit.Shares(), 0))"},
		{"7-5-crew-exam", "colony.Answer(i, j)", "colony.Answer(i, (j+1)%len(question.Choices))"},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			body, err := os.ReadFile(filepath.Join("testdata", "solutions", c.id+".go"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(body), c.from) {
				t.Fatalf("mutation does not apply: %q", c.from)
			}
			mutant := strings.Replace(string(body), c.from, c.to, 1)
			resp, err := r.Run(context.Background(), levels.Get(c.id), mutant)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Phase == "compile" || len(resp.Missing) > 0 {
				t.Fatalf("mutant must compile and meet requirements: %+v %s", resp.Missing, resp.CompileOutput)
			}
			if resp.Result.Success {
				t.Fatal("incorrect result passed")
			}
			starter, err := r.Run(context.Background(), levels.Get(c.id), levels.Get(c.id).Starter)
			if err != nil {
				t.Fatal(err)
			}
			if starter.Phase == "compile" || starter.Result.Success {
				t.Fatalf("starter must compile without passing: %+v", starter)
			}
		})
	}
}

// References must not depend on language features that have not been taught.
// The hidden tests themselves are instructor code, not player prerequisites.
func TestP1ReferencesRespectCurriculum(t *testing.T) {
	ids := []string{"2-3-callsigns", "2-4-units", "4-3-checkpoint", "5-3-fault-log", "5-4-undock", "5-5-command-table", "6-4-stock-order", "6-5-heat-grid", "7-4-composite-units", "7-5-crew-exam"}
	for _, id := range ids {
		t.Run(id, func(t *testing.T) {
			level := levels.Get(id)
			name := id + ".go"
			if level.Kind == "lab" {
				name = filepath.Join(id, "main.go")
			}
			body, err := os.ReadFile(filepath.Join("testdata", "solutions", name))
			if err != nil {
				t.Fatal(err)
			}
			source, err := parser.ParseFile(token.NewFileSet(), name, body, 0)
			if err != nil {
				t.Fatal(err)
			}
			ast.Inspect(source, func(n ast.Node) bool {
				forbidden := ""
				switch node := n.(type) {
				case *ast.ForStmt, *ast.IfStmt, *ast.SwitchStmt:
					if level.Chapter < 3 {
						forbidden = "loops or branches before chapter 3/4"
					}
				case *ast.RangeStmt, *ast.ArrayType, *ast.MapType:
					if level.Chapter < 6 {
						forbidden = "collections or range before chapter 6"
					}
				case *ast.StructType, *ast.InterfaceType:
					if level.Chapter < 7 {
						forbidden = "structs or interfaces before chapter 7"
					}
				case *ast.StarExpr:
					if level.Chapter < 8 {
						forbidden = "pointer syntax before chapter 8"
					}
				case *ast.UnaryExpr:
					if node.Op == token.AND && level.Chapter < 8 {
						forbidden = "address-taking before chapter 8"
					}
				case *ast.FuncDecl:
					if node.Recv != nil && level.Chapter < 7 {
						forbidden = "methods before chapter 7"
					}
					if node.Name.Name != "main" && level.Chapter < 5 {
						forbidden = "player-defined functions before chapter 5"
					}
				case *ast.GoStmt, *ast.ChanType, *ast.SelectStmt:
					forbidden = "concurrency before its chapters"
				}
				if forbidden != "" {
					t.Error(forbidden)
				}
				return true
			})
		})
	}
}
