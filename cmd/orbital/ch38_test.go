package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"orbital/colony/levels"
)

func TestModernGoVetStarterBuildsAndReportsAllBugs(t *testing.T) {
	r := testRunner(t)
	lvl := levels.Get("38-6-vet")
	files := map[string]string{"main.go": lvl.Files[0].Starter}
	ok, out, err := r.goTest(context.Background(), lvl, files, "./...")
	if err != nil || !ok {
		t.Fatalf("starter must compile with vet disabled: %v %s", err, out)
	}
	ok, out, err = r.goTest(context.Background(), lvl, files, "-vet=all", "-run", "^$", "./...")
	if err != nil || ok {
		t.Fatalf("vet must reject starter: %v %s", err, out)
	}
	for _, diagnostic := range []string{"passes lock by value", "wrong type", "cancel function"} {
		if !strings.Contains(out, diagnostic) {
			t.Errorf("missing %q diagnostic:\n%s", diagnostic, out)
		}
	}
}

func TestModernGoVetEnforcedWhenBehaviorPasses(t *testing.T) {
	lvl := *levels.Get("38-6-vet")
	body, err := os.ReadFile(filepath.Join("testdata", "solutions", lvl.ID, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	// The extra function is never called: behavioral tests alone cannot catch this lost cancel.
	code := string(body) + "\nfunc lostCancel(parent context.Context) { ctx, _ := context.WithCancel(parent); _ = ctx }\n"
	lab := *lvl.Lab
	lab.Vet = true
	lvl.Lab = &lab
	ok, out, err := testRunner(t).goTest(context.Background(), &lvl, map[string]string{"main.go": code, "zz_hidden_test.go": lab.Hidden}, "-run", "^TestHidden", "./...")
	if err != nil || !ok {
		t.Fatalf("behavior must pass: %v %s", err, out)
	}
	resp, err := testRunner(t).RunLab(withLang(context.Background(), "en"), &lvl, map[string]string{"main.go": code})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Result.Success || resp.Phase != "lab" {
		t.Fatalf("vet escaped: %+v", resp)
	}
	found := false
	for _, step := range resp.Steps {
		if step.Title == "go vet inspection" {
			found = true
			if step.OK || !strings.Contains(step.Output, "cancel function") {
				t.Fatalf("unexpected vet result: %+v", step)
			}
		}
	}
	if !found {
		t.Fatalf("missing vet step: %+v", resp.Steps)
	}
}
