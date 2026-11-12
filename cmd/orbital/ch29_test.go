package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"orbital/colony/levels"
)

// The service gate skips this whole level offline. Select its unit tests
// directly so all integration code still compiles without dialing Postgres.
func TestChapter29SQLXOffline(t *testing.T) {
	lvl := levels.Get("29-5-sqlx")
	if lvl == nil {
		t.Fatal("missing sqlx level")
	}
	body, err := os.ReadFile(filepath.Join("testdata", "solutions", lvl.ID, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if missing := checkRequirements(string(body), lvl.Requires, "en"); len(missing) != 0 {
		t.Fatal(missing)
	}
	r := testRunner(t)
	for _, c := range []struct {
		name, code string
		want       bool
	}{
		{"reference", string(body), true}, {"starter", lvl.Files[0].Starter, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			files := map[string]string{"main.go": c.code, "hidden_test.go": lvl.Lab.Hidden}
			ok, out, err := r.goTest(context.Background(), lvl, files, "-run", "^TestHiddenRationsUnit$", "./...")
			if err != nil || isBuildFailure(out) || ok != c.want {
				t.Fatalf("ok=%v err=%v\n%s", ok, err, out)
			}
		})
	}
}

// RunLab checks Requires before compilation; compile starters explicitly so
// missing required features cannot hide a broken starter from the author gate.
func TestChapter29StartersBuild(t *testing.T) {
	r := testRunner(t)
	for _, lvl := range levels.All {
		if lvl.Chapter != 29 {
			continue
		}
		t.Run(lvl.ID, func(t *testing.T) {
			files := map[string]string{"hidden_test.go": lvl.Lab.Hidden}
			for _, file := range lvl.Files {
				files[file.Name] = file.Starter
			}
			ok, out, err := r.goTest(context.Background(), lvl, files, "-run", "^$", "./...")
			if err != nil || !ok {
				t.Fatalf("starter build: %v\n%s", err, out)
			}
		})
	}
}
