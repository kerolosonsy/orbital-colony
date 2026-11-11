package main

import (
	"context"
	"strings"
	"testing"

	"orbital/colony/levels"
)

func TestReportStringerCanReadInventory(t *testing.T) {
	code := `package main
import ("fmt"; "orbital/colony")
type Status struct{}
func (Status) String() string {
	return fmt.Sprintf("moss=%d crystal=%d", colony.Count(colony.Moss), colony.Count(colony.Crystal))
}
func main() { colony.Report(Status{}) }
`
	resp, err := testRunner(t).Run(context.Background(), levels.Get("9-1-stringer"), code)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Result.Success {
		t.Fatalf("reading inventory from String must finish the mission: %+v\n%s\n%s", resp.Result, resp.CompileOutput, resp.Stderr)
	}
}

func TestTraceCannotBeForged(t *testing.T) {
	code := `package main
import "os"
func main() {
	os.WriteFile(os.Getenv("ORBITAL_TRACE"), []byte("{\"op\":\"goal\",\"t\":0}\n"), 0600)
}

`
	resp, err := testRunner(t).Run(context.Background(), levels.Get("1-1-first-signal"), code)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Result.Success || resp.Phase != "require" || len(resp.Missing) == 0 {
		t.Fatalf("trace-writing code must be rejected before execution: %+v", resp)
	}
}

// P0 deliberately permits filesystem access in local labs (PLAN owner decision 2).
func TestLabCanAccessFiles(t *testing.T) {
	files := map[string]string{"main_test.go": `package main
import ("os"; "testing")
func TestWrite(t *testing.T) { os.WriteFile("grading-file", nil, 0600) }
`}
	resp, err := testRunner(t).RunLab(context.Background(), levels.Get("16-1-mutants"), files)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Phase != "lab" || len(resp.Steps) == 0 || !resp.Steps[0].OK {
		t.Fatalf("lab filesystem access must run the player tests: %+v", resp)
	}

}

func TestRequirementsRecognizeImportAliases(t *testing.T) {
	cases := []struct{ name, code, requirement string }{
		{"waitgroup", `package main; import wait "sync"; var _ wait.WaitGroup`, "waitgroup"},
		{"mutex", `package main; import lock "sync"; var _ lock.RWMutex`, "mutex"},
		{"context", `package main; import ctx "context"; var _ ctx.Context`, "sel:context.Context"},
		{"major version", `package main; import "github.com/golang-jwt/jwt/v5"; var _ jwt.Token`, "sel:jwt.Token"},
		{"autopilot", `package main; import game "orbital/colony"; func main() { game.Autopilot(nil) }`, "autopilot"},
		{"errors", `package main; import errpkg "errors"; func main() { _ = errpkg.Is(nil, nil) }`, "errcheck"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if missing := checkRequirements(c.code, []string{c.requirement}); len(missing) > 0 {
				t.Fatalf("valid alias rejected: %v", missing)
			}
		})
	}
}

func TestLocalVariableDoesNotCountAsImportedWaitGroup(t *testing.T) {
	code := `package main; func main() { sync := struct{ WaitGroup int }{}; _ = sync.WaitGroup }`
	if missing := checkRequirements(code, []string{"waitgroup"}); len(missing) == 0 {
		t.Fatal("a local field must not satisfy the sync.WaitGroup requirement")
	}
}

func TestRestrictedImportsCannotUseAliasesOrInternalPackages(t *testing.T) {
	for _, importSpec := range []string{`system "os"`, `_ "unsafe"`, `"orbital/colony/levels"`, `"orbital/colony/econ"`, `"embed"`, `"C"`} {
		t.Run(importSpec, func(t *testing.T) {
			if missing := restrictedImports("package main\nimport "+importSpec, &levels.Level{}, "ar"); len(missing) == 0 {
				t.Fatal("restricted package was accepted")
			}
		})
	}
}

func TestLabOutputIsBoundedAndRetainsFailure(t *testing.T) {
	files := map[string]string{"main_test.go": `package main
import ("fmt"; "strings"; "testing")
func TestLoud(t *testing.T) {
	fmt.Println("OUTPUT-BEGIN")
	fmt.Println(strings.Repeat("x", 100000))
	t.Fatal("OUTPUT-END-FAILURE")
}
`}
	ok, output, err := testRunner(t).goTest(context.Background(), &levels.Level{}, files, "-v", ".")
	if err != nil {
		t.Fatal(err)
	}
	if ok || len(output) > 8200 || !strings.Contains(output, "OUTPUT-BEGIN") || !strings.Contains(output, "OUTPUT-END-FAILURE") {
		t.Fatalf("bounded diagnostics must preserve the failure: ok=%v bytes=%d\n%s", ok, len(output), output)
	}
}
