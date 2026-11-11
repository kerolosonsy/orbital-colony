package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"orbital/colony/levels"
)

type labStep struct {
	Title  string `json:"title"`
	OK     bool   `json:"ok"`
	Output string `json:"output,omitempty"`
}

var benchLine = regexp.MustCompile(`(?m)^(BenchmarkHidden\w+)-?\d*\s+\d+\s+([\d.]+) ns/op`)

// goTest writes files into a fresh module and runs go test with args.
func (r *runner) goTest(ctx context.Context, files map[string]string, args ...string) (bool, string, error) {
	dir, err := os.MkdirTemp("", "orbital-lab-*")
	if err != nil {
		return false, "", err
	}
	defer os.RemoveAll(dir)
	gomod := fmt.Sprintf("module player\n\ngo 1.24\n\nrequire orbital/colony v0.0.0\n\nreplace orbital/colony => %s\n", r.colonyDir)
	files["go.mod"] = gomod
	for name, body := range files {
		if strings.ContainsAny(name, `/\`) || !strings.HasSuffix(name, ".go") && name != "go.mod" {
			return false, "", fmt.Errorf("bad file name %q", name)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			return false, "", err
		}
	}
	tctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	out := &limitedBuffer{max: 8000, lang: contextLang(ctx)}
	cmd := r.goCmd(tctx, dir, append([]string{"test", "-count=1"}, args...)...)
	cmd.Stdout, cmd.Stderr = out, out
	runErr := runCommand(cmd)
	text := strings.ReplaceAll(out.String(), dir+string(os.PathSeparator), "")
	return runErr == nil, text, nil
}

func isBuildFailure(out string) bool {
	return strings.Contains(out, "[build failed]") || strings.Contains(out, "[setup failed]")
}

func (r *runner) RunLab(ctx context.Context, lvl *levels.Level, edited map[string]string) (*runResponse, error) {
	lang := contextLang(ctx)
	lvl = levelForLang(lvl, lang)
	resp := &runResponse{Phase: "lab"}
	files := map[string]string{}
	for _, f := range lvl.Files {
		files[f.Name] = f.Starter
		if body, ok := edited[f.Name]; ok && !f.ReadOnly {
			files[f.Name] = body
		}
		if missing := restrictedImports(files[f.Name], lang); len(missing) > 0 {
			return rejectedImports(missing, lang), nil
		}
	}
	with := func(over map[string]string) map[string]string {
		m := map[string]string{}
		for k, v := range files {
			m[k] = v
		}
		for k, v := range over {
			m[k] = v
		}
		return m
	}
	add := func(title string, ok bool, out string) bool {
		resp.Steps = append(resp.Steps, labStep{title, ok, strings.TrimSpace(out)})
		return ok
	}
	lab := lvl.Lab

	// 1. The player's own tests must build and pass on the player's code.
	ok, out, err := r.goTest(ctx, with(nil), "-v", ".")
	if err != nil {
		return nil, err
	}
	if isBuildFailure(out) {
		resp.Phase = "compile"
		resp.CompileOutput = out
		for _, line := range strings.Split(out, "\n") {
			if ce, ok := parseCompileLine(line, lang); ok {
				resp.CompileErrors = append(resp.CompileErrors, ce)
			}
		}
		resp.Result.Message = tr(lang, "الكود مااتبناش. صلّح الأخطاء اللي تحت.", "The code did not compile. Fix the errors below.")
		return resp, nil
	}
	if !add(tr(lang, "التستات بتاعتك بتعدّي على الكود بتاعك", "Your tests pass on your code"), ok, out) {
		resp.Result.Message = tr(lang, "التستات بتاعتك فشلت. لازم تعدّي على الكود الصح الأول.", "Your tests failed. They must pass on the correct code first.")
		return resp, nil
	}
	allOK := true
	if missing := checkRequirements(files[lab.ReqFile], lvl.Requires, lang); len(missing) > 0 {
		resp.Missing = missing
		allOK = false
	}

	// 2. Reproduce the bug: the tests must fail on the original code.
	if lab.Buggy != "" {
		ok, out, err := r.goTest(ctx, with(map[string]string{"main.go": lab.Buggy}), ".")
		if err != nil {
			return nil, err
		}
		allOK = add(tr(lang, "التستات بتاعتك بتمسك الـ bug في الكود الأصلي", "Your tests catch the bug in the original code"), !ok && !isBuildFailure(out), out) && allOK
	}

	// 3. Mutation testing: every broken variant must be caught.
	killed := 0
	for i, m := range lab.Mutants {
		ok, out, err := r.goTest(ctx, with(map[string]string{"main.go": m.Code}), ".")
		if err != nil {
			return nil, err
		}
		caught := !ok && !isBuildFailure(out)
		if caught {
			killed++
			add(fmt.Sprintf(tr(lang, "Mutant %d اتمسك ✓", "Mutant %d caught ✓"), i+1), true, "")
		} else {
			add(fmt.Sprintf(tr(lang, "Mutant %d عدّى من التستات: %s", "Mutant %d escaped your tests: %s"), i+1, mutantDescription(lang, m.Desc)), false, "")
		}
	}
	if len(lab.Mutants) > 0 && killed < len(lab.Mutants) {
		allOK = false
	}

	// 4. Hidden checks of the player's implementation.
	speedup := 0.0
	if lab.Hidden != "" {
		hidden := with(map[string]string{"zz_hidden_test.go": lab.Hidden})
		ok, out, err := r.goTest(ctx, hidden, "-run", "^TestHidden", ".")
		if err != nil {
			return nil, err
		}
		allOK = add(tr(lang, "تستات Mission Control على الكود بتاعك", "Mission Control tests on your code"), ok, out) && allOK
		if ok && lab.Bench != nil {
			_, out, err := r.goTest(ctx, hidden, "-run", "^$", "-bench", "^BenchmarkHidden", "-benchtime", "200ms", ".")
			if err != nil {
				return nil, err
			}
			ns := map[string]float64{}
			for _, m := range benchLine.FindAllStringSubmatch(out, -1) {
				ns[m[1]], _ = strconv.ParseFloat(m[2], 64)
			}
			naive, theirs := ns["BenchmarkHiddenNaive"], ns["BenchmarkHiddenTheirs"]
			if naive > 0 && theirs > 0 {
				speedup = naive / theirs
			}
			fast := speedup >= lab.Bench.Tiers[0]
			allOK = add(fmt.Sprintf(tr(lang, "Benchmark: الكود بتاعك أسرع %.1fx من النسخة البطيئة (المطلوب %.0fx على الأقل)", "Benchmark: your code is %.1fx faster than the slow version (at least %.0fx required)"), speedup, lab.Bench.Tiers[0]), fast, out) && allOK
		}
	}

	resp.Result.Success = allOK
	switch {
	case !allOK && len(resp.Missing) > 0:
		resp.Result.Message = tr(lang, "فاضل حاجات. والمهمة دي محتاجة مفاهيم معيّنة (تحت).", "Some steps remain. This mission also needs the concepts listed below.")
	case !allOK:
		resp.Result.Message = tr(lang, "فيه خطوات لسه مانجحتش. بص على الـ Console.", "Some steps have not passed yet. Check the Console.")
	case lab.Bench != nil:
		resp.Result.Stars = 1
		for i, t := range lab.Bench.Tiers {
			if speedup >= t {
				resp.Result.Stars = i + 1
			}
		}
		resp.Result.Message = fmt.Sprintf(tr(lang, "المختبر عدّى! أسرع %.0fx.", "Lab passed! %.0fx faster."), speedup)
	default:
		resp.Result.Stars = 3
		resp.Result.Message = tr(lang, "المختبر عدّى كله!", "All lab steps passed!")
	}
	return resp, nil
}
