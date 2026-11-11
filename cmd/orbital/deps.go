package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"orbital/colony/levels"
)

func moduleFiles(colonyDir string, lvl *levels.Level) (map[string]string, error) {
	body := "module player\n\ngo 1.24\n"
	files := map[string]string{}
	if len(lvl.Deps) > 0 {
		dir := filepath.Join(filepath.Dir(colonyDir), "labdeps")
		mod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err != nil {
			return nil, err
		}
		catalog := string(mod)
		for _, dep := range lvl.Deps {
			found := false
			for _, line := range strings.Split(catalog, "\n") {
				fields := strings.Fields(line)
				if len(fields) >= 2 && fields[0] == dep {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("dependency %q is absent from labdeps/go.mod", dep)
			}
		}
		// Keep the complete pinned graph, including indirect versions, for offline MVS.
		_, directives, ok := strings.Cut(catalog, "\n")
		if !ok {
			return nil, fmt.Errorf("invalid labdeps/go.mod")
		}
		body = "module player\n" + directives
		sum, err := os.ReadFile(filepath.Join(dir, "go.sum"))
		if err != nil {
			return nil, err
		}
		files["go.sum"] = string(sum)
	}
	files["go.mod"] = body + fmt.Sprintf("\nrequire orbital/colony v0.0.0\n\nreplace orbital/colony => %s\n", filepath.ToSlash(colonyDir))
	return files, nil
}

func (r *runner) writeModule(dir string, lvl *levels.Level) error {
	files, err := moduleFiles(r.colonyDir, lvl)
	if err != nil {
		return err
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func dependencyHint(out, lang string) string {
	for _, needle := range []string{"missing go.sum entry", "cannot find module", "module lookup disabled", "module download", "GOPROXY=off"} {
		if strings.Contains(out, needle) {
			return out + "\n" + tr(lang, "شغّل `make deps` مرة واحدة عشان تجهّز المكتبات.", "Run `make deps` once to prepare the dependencies.")
		}
	}
	return out
}

func validLabName(name string) bool {
	if !filepath.IsLocal(name) || strings.Contains(name, `\`) {
		return false
	}
	for _, part := range strings.Split(filepath.ToSlash(name), "/") {
		if part == ".." {
			return false
		}
	}
	switch filepath.Ext(name) {
	case ".go", ".sql", ".proto", ".json", ".yaml", ".txt", ".md", ".html":
		return true
	}
	return false
}
