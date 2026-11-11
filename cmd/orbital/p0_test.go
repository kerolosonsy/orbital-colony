package main

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"orbital/colony/levels"
)

func TestP0ImportPolicy(t *testing.T) {
	for _, tc := range []struct {
		pkg, kind string
		deps      []string
		allowed   bool
	}{
		{pkg: "net/http"}, {pkg: "net/http", kind: "lab", allowed: true}, {pkg: "os", kind: "lab", allowed: true},
		{pkg: "os/exec", kind: "lab"}, {pkg: "unsafe", kind: "lab"}, {pkg: "syscall", kind: "lab"}, {pkg: "plugin", kind: "lab"},
		{pkg: "orbital/colony/internal/sim", kind: "lab"}, {pkg: "orbital/colony/levels", kind: "lab"},
		{pkg: "github.com/gin-gonic/gin", kind: "lab"},
		{pkg: "github.com/gin-gonic/gin", kind: "lab", deps: []string{"github.com/gin-gonic/gin"}, allowed: true},
		{pkg: "github.com/gin-gonic/gin/binding", kind: "lab", deps: []string{"github.com/gin-gonic/gin"}, allowed: true},
		{pkg: "github.com/gin-gonic/gin-extra", kind: "lab", deps: []string{"github.com/gin-gonic/gin"}},
		{pkg: "orbital/colony/labkit", kind: "lab", allowed: true}, {pkg: "bufio", allowed: true},
	} {
		t.Run(tc.kind+"/"+tc.pkg, func(t *testing.T) {
			missing := restrictedImports("package main\nimport _ "+`"`+tc.pkg+`"`, &levels.Level{Kind: tc.kind, Deps: tc.deps}, "en")
			if (len(missing) == 0) != tc.allowed {
				t.Fatalf("allowed=%v missing=%v", tc.allowed, missing)
			}
		})
	}
}

func TestP0Requirements(t *testing.T) {
	cases := map[string][2]string{
		"rune":      {`func f(s string) { for range s {} }`, `func f(s []int) { for range s {} }`},
		"iota":      {`const x = iota`, `var iota = 0`},
		"const":     {`const x = 1`, `var x = 1`},
		"defer":     {`func f() { defer f() }`, `func f() { f() }`},
		"recover":   {`func f() { defer func(){ recover() }() }`, `func f() { recover() }`},
		"errwrap":   {`import("fmt";"errors"); func f(e error) { _=fmt.Errorf("oops: %w",e); errors.Is(e,e) }`, `import("fmt";"errors"); func f(e error) { _=fmt.Errorf("oops: %%w",e); errors.Is(e,e) }`},
		"errtype":   {`type E struct{}; func(E) Error() string {return ""}`, `type E struct{}; func(E) Error() int {return 0}`},
		"funcvalue": {`var actions map[string]func()`, `func f(){ func(){}() }`},
		"closure":   {`func f() func() int { x:=1; return func() int{return x} }`, `func f() { func(){}() }`},
		"embed":     {`type A struct{ B }`, `type A struct{ b B }`},
		"variadic":  {`func f(xs ...int){}`, `func f(xs []int){}`},
		"bufchan":   {`var ch = make(chan int, 1)`, `var ch = make(chan int)`},
		"close":     {`func f(ch chan int){close(ch)}`, `func f(ch chan int){ _=ch }`},
		"recursion": {`func f(){ f() }`, `func f(){ g() }; func g(){}`},
		"fuzz":      {`import "testing"; func FuzzX(f *testing.F){}`, `import "testing"; func FuzzX(t *testing.T){}`},
		"subtest":   {`import "testing"; func TestX(t *testing.T){ t.Run("x",func(t *testing.T){}) }`, `type X struct{}; func f(t X){ t.Run() }`},
		"iface":     {`type Reader interface{ Read() }`, `type Reader struct{}`},
	}
	for key, pair := range cases {
		t.Run(key, func(t *testing.T) {
			for i, code := range pair {
				code = "package main\n" + code
				if _, err := scan(code); err != nil {
					t.Fatal(err)
				}
				missing := checkRequirements(code, []string{key}, "en")
				if (len(missing) == 0) != (i == 0) {
					t.Fatalf("positive=%v missing=%v", i == 0, missing)
				}
			}
			for _, lang := range []string{"ar", "en"} {
				if requirementTexts(lang)[key] == "" {
					t.Fatal("missing translation")
				}
			}
		})
	}
}

func fixtureLab() *levels.Level {
	return &levels.Level{Kind: "lab", Files: []levels.LabFile{
		{Name: "main.go", Starter: "package main\nfunc Value() int {return 7}"},
		{Name: "main_test.go", Starter: "package main\nimport \"testing\"\nfunc TestValue(t *testing.T) {if Value()!=7 {t.Fatal(Value())}}"},
	}, Lab: &levels.Lab{ReqFile: "main.go"}}
}

func runFixture(t *testing.T, lvl *levels.Level) *runResponse {
	t.Helper()
	resp, err := testRunner(t).RunLab(withLang(context.Background(), "en"), lvl, nil)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

// This fixture is deliberately absent from levels.All: P0 adds infrastructure only.
func TestP0GinLabOffline(t *testing.T) {
	l := fixtureLab()
	l.Deps = []string{"github.com/gin-gonic/gin"}
	l.Files[0].Starter = `package main
import (
 "net/http"
 "os"
 "github.com/gin-gonic/gin"
 "orbital/colony/labkit"
 "player/nested"
)
func Value() int { _=http.MethodGet; _=os.Getenv("PATH"); _=labkit.FreeAddr(); return len(gin.New().Routes())+6+nested.X }
`
	l.Files = append(l.Files, levels.LabFile{Name: "fixtures/note.txt", Starter: "import \"os/exec\"", ReadOnly: true}, levels.LabFile{Name: "nested/data.go", Starter: "package nested\nconst X=1", ReadOnly: true})
	l.Lab.HiddenFiles = map[string]string{"extra_hidden_test.go": `package main
import "testing"
func TestHiddenGin(t *testing.T){if Value()!=7 {t.Fatal(Value())}}`}
	resp := runFixture(t, l)
	if !resp.Result.Success || resp.Result.Stars != 3 {
		t.Fatalf("%+v", resp)
	}
}

func TestP0LabBonus(t *testing.T) {
	for _, tc := range []struct {
		hidden  string
		success bool
		stars   int
	}{
		{`func TestHidden(t *testing.T){}; func TestHiddenBonus(t *testing.T){}`, true, 3},
		{`func TestHidden(t *testing.T){}; func TestHiddenBonus(t *testing.T){t.Fatal("bonus")}`, true, 2},
		{`func TestHidden(t *testing.T){t.Fatal("required")}; func TestHiddenBonus(t *testing.T){}`, false, 0},
	} {
		t.Run(tc.hidden, func(t *testing.T) {
			l := fixtureLab()
			l.Lab.Bonus = true
			l.Lab.Hidden = "package main\nimport \"testing\"\n" + tc.hidden
			resp := runFixture(t, l)
			if resp.Result.Success != tc.success || resp.Result.Stars != tc.stars {
				t.Fatalf("%+v", resp)
			}
		})
	}
}

func TestP0LabRace(t *testing.T) {
	for _, race := range []bool{false, true} {
		t.Run(map[bool]string{false: "clean", true: "race"}[race], func(t *testing.T) {
			l := fixtureLab()
			l.Race = true
			if race {
				l.Files[1].Starter = `package main
import("testing";"sync")
func TestRace(t *testing.T){ var x int; var wg sync.WaitGroup; for i:=0;i<2;i++ {wg.Add(1);go func(){defer wg.Done();for j:=0;j<1000;j++ {x++}}()};wg.Wait();_=x }
`
			}
			resp := runFixture(t, l)
			if resp.Result.Success == race || resp.Race != race || (race && resp.RuntimeTip == "") {
				t.Fatalf("%+v", resp)
			}
		})
	}
}

func TestP0LabFuzz(t *testing.T) {
	for _, broken := range []bool{false, true} {
		t.Run(map[bool]string{false: "pass", true: "fail"}[broken], func(t *testing.T) {
			l := fixtureLab()
			l.Lab.Fuzz = "FuzzValue"
			body := "if Value()!=7 {t.Fatal(Value())}"
			if broken {
				body = "t.Fatal(s)"
			}
			l.Files = append(l.Files, levels.LabFile{Name: "fuzz_test.go", Starter: `package main
import "testing"
func FuzzValue(f *testing.F){ f.Fuzz(func(t *testing.T,s string){ ` + body + ` }) }`})
			resp := runFixture(t, l)
			if resp.Result.Success == broken {
				t.Fatalf("%+v", resp)
			}
			if len(resp.Steps) < 2 || !strings.Contains(resp.Steps[len(resp.Steps)-1].Title, "Fuzz") {
				t.Fatalf("missing fuzz step: %+v", resp)
			}
		})
	}
}

func TestP0LabFileNames(t *testing.T) {
	for _, name := range []string{"nested/file.go", "schema.sql", "api.proto", "fixture.json", "compose.yaml", "readme.txt", "notes.md", "view.html"} {
		if !validLabName(name) {
			t.Error(name)
		}
	}
	for _, name := range []string{"../outside.go", "nested/../main.go", "/abs.go", "..", "go.mod", "file.exe", `nested\main.go`} {
		l := fixtureLab()
		l.Files[0].Name = name
		if _, err := testRunner(t).RunLab(context.Background(), l, nil); err == nil {
			t.Errorf("accepted %q", name)
		}
	}
}

func TestP0WorkspaceDeps(t *testing.T) {
	r := testRunner(t)
	root := t.TempDir()
	rel, err := filepath.Rel(root, r.colonyDir)
	if err != nil {
		t.Fatal(err)
	}
	ws := &workspace{root: root, colonyRel: rel}
	l := fixtureLab()
	l.ID = "fixture"
	l.Deps = []string{"github.com/gin-gonic/gin"}
	l.Files = append(l.Files, levels.LabFile{Name: "assets/schema.sql", Starter: "SELECT 1", ReadOnly: true})
	if _, err := ws.open(l, nil); err != nil {
		t.Fatal(err)
	}
	mod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := os.ReadFile("../../labdeps/go.mod")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(catalog), "\n") {
		if strings.HasPrefix(line, "go ") || strings.Contains(line, "github.com/gin-gonic/gin ") {
			if !strings.Contains(string(mod), line) {
				t.Fatalf("missing %s", line)
			}
		}
	}
	sum, err := os.ReadFile(filepath.Join(root, "go.sum"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("../../labdeps/go.sum")
	if err != nil {
		t.Fatal(err)
	}
	if string(sum) != string(want) {
		t.Fatal("go.sum mismatch")
	}
	files, _, err := ws.read(l)
	if err != nil || files["assets/schema.sql"] != "SELECT 1" {
		t.Fatalf("%v %v", files, err)
	}
}

func TestP0MissingDependencyHint(t *testing.T) {
	for _, out := range []string{"missing go.sum entry", "cannot find module x", "module lookup disabled by GOPROXY=off"} {
		for _, lang := range []string{"ar", "en"} {
			if !strings.Contains(dependencyHint(out, lang), "make deps") {
				t.Fatal(out)
			}
		}
	}
}

func TestP0ServicePrecheckAndEnvironment(t *testing.T) {
	t.Run("live-loopback", func(t *testing.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if errors.Is(err, syscall.EPERM) {
			t.Skip("sandbox prohibits loopback listeners")
		}
		if err != nil {
			t.Fatal(err)
		}
		defer listener.Close()
		t.Setenv("ORBITAL_PG_DSN", "postgres://orbital:orbital@"+listener.Addr().String()+"/orbital")
		if err := checkService(context.Background(), "postgres"); err != nil {
			t.Fatal(err)
		}
	})
	t.Setenv("ORBITAL_PG_DSN", "postgres://127.0.0.1:0/orbital")
	r := testRunner(t)
	env := strings.Join(r.goCmd(context.Background(), t.TempDir(), "version").Env, "\n")
	for name := range serviceDefaults {
		if !strings.Contains(env, name+"="+serviceEnv(name)) {
			t.Fatal("missing " + name)
		}
	}
	// An invalid go binary proves a missing service returns before any build.
	r.goBin = "/no/such/go"
	l := fixtureLab()
	l.Service = "postgres"
	resp, err := r.RunLab(withLang(context.Background(), "en"), l, nil)
	if err != nil || len(resp.Steps) != 1 || resp.Steps[0].OK || !strings.Contains(resp.Result.Message, "make services-up") {
		t.Fatalf("%+v %v", resp, err)
	}
}

func TestP0LabExecRejectedBeforeBuild(t *testing.T) {
	l := fixtureLab()
	l.Files[0].Starter = "package main\nimport _ \"os/exec\""
	r := testRunner(t)
	r.goBin = "/no/such/go"
	resp, err := r.RunLab(context.Background(), l, nil)
	if err != nil || resp.Phase != "require" || len(resp.Missing) == 0 {
		t.Fatalf("%+v %v", resp, err)
	}
}
