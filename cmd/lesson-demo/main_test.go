package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildPublishesOnlyBrowserAssets(t *testing.T) {
	out := t.TempDir()
	if err := build(out); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(out)
	if err != nil || len(entries) != 6 {
		t.Fatalf("expected six static demo files, got %v: %v", entries, err)
	}
	data, err := os.ReadFile(filepath.Join(out, "mission.js"))
	if err != nil || !strings.Contains(string(data), "First Signal") || !strings.Contains(string(data), "إشارة أولى") {
		t.Fatalf("expected bilingual First Signal data, got %s: %v", data, err)
	}
}
