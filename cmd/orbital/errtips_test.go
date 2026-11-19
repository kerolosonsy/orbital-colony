package main

import (
	"context"
	"strings"
	"testing"

	"orbital/colony/levels"
)

func TestErrorTipCatalogueSamples(t *testing.T) {
	for _, entry := range errorTips {
		t.Run(entry.needle, func(t *testing.T) {
			if !strings.Contains(entry.sample, entry.needle) {
				t.Fatalf("sample %q does not contain matcher %q", entry.sample, entry.needle)
			}
			for _, lang := range []string{"ar", "en"} {
				got := errorTipFor(entry.sample, entry.scope, lang)
				if got.Meaning == "" || got.Cause == "" || got.Check == "" {
					t.Fatalf("%s tip is incomplete: %+v", lang, got)
				}
			}
		})
	}
}

func TestGeneratedCompileAndRuntimeErrorTips(t *testing.T) {
	r := testRunner(t)
	lvl := levels.All[0]

	compile, err := r.Run(context.Background(), lvl, `package main
func main() {
	var count int = "four"
	_ = count
}

func TestErrorTipPrefersSpecificDiagnostics(t *testing.T) {
	for _, tc := range []struct {
		message string
		want    string
	}{
		{"invalid operation: a + b (mismatched types int and string)", "العملية بتجمع أو تقارن أنواع مختلفة."},
		{"invalid operation: cannot call count (variable of type int): int is not a function", "الكود بيحاول يستدعي قيمة مش دالة."},
	} {
		got := errorTipFor(tc.message, tipCompile, "ar")
		if got.Meaning != tc.want {
			t.Errorf("errorTipFor(%q).Meaning = %q, want %q", tc.message, got.Meaning, tc.want)
		}
	}
}
`)
	if err != nil {
		t.Fatal(err)
	}
	if compile.Phase != "compile" || len(compile.CompileErrors) == 0 {
		t.Fatalf("expected a compiler diagnostic, got phase %q: %s", compile.Phase, compile.CompileOutput)
	}
	gotCompile := compile.CompileErrors[0]
	if !strings.Contains(gotCompile.Msg, "cannot use") || gotCompile.Tip != "القيمة دي مش من النوع المطلوب." || gotCompile.TipCause == "" || gotCompile.TipCheck == "" {
		t.Fatalf("compiler tip did not explain the generated type error: %+v", gotCompile)
	}

	runtime, err := r.Run(context.Background(), lvl, `package main
import "fmt"
func main() {
	i := 2
	fmt.Println([]int{1}[i])
}
`)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.Phase != "run" || !strings.Contains(runtime.Stderr, "index out of range") {
		t.Fatalf("expected an index panic, got phase %q and stderr %q", runtime.Phase, runtime.Stderr)
	}
	if runtime.RuntimeTip != "الـ index خارج حدود الـ slice أو array." || runtime.RuntimeCause == "" || runtime.RuntimeCheck == "" {
		t.Fatalf("runtime tip is incomplete: %+v", runtime)
	}
}
