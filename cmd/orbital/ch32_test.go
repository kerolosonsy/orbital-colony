package main

import (
	"context"
	"testing"

	"orbital/colony/levels"
)

func TestReadOnlyGeneratedSourceIsTrustedAndCannotBeOverridden(t *testing.T) {
	lvl := fixtureLab()
	lvl.Files[0].Starter = "package main\nfunc Value() int {return generatedSize()-1}"
	lvl.Files = append(lvl.Files, levels.LabFile{Name: "wire.pb.go", ReadOnly: true, Starter: `package main
import "unsafe"
func generatedSize() int {return int(unsafe.Sizeof(int64(0)))}`})
	resp, err := testRunner(t).RunLab(withLang(context.Background(), "en"), lvl, map[string]string{
		"wire.pb.go": "package main\nfunc generatedSize() int {return 99}",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Result.Success {
		t.Fatalf("trusted source rejected or replaced: phase=%s missing=%v steps=%+v", resp.Phase, resp.Missing, resp.Steps)
	}
	lvl.Files[0].Starter = "package main\nimport _ \"unsafe\"\nfunc Value() int {return 7}"
	resp = runFixture(t, lvl)
	if resp.Phase != "require" || resp.Result.Success {
		t.Fatalf("editable unsafe import escaped policy: %+v", resp)
	}
}
