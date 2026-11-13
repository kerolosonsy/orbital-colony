package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orbital/colony/levels"
)

func (ws *workspace) openColony(level *levels.Level, files map[string]string, lang string) (string, error) {
	if len(files) == 0 {
		files = map[string]string{"main.go": level.Starter}
	}
	if err := validateColonyFiles(files, lang); err != nil {
		return "", err
	}
	dir := filepath.Join(ws.root, "colony")
	if err := rejectWorkspaceSymlinks(dir, "go.mod", lang); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if err := ws.colonyModule(dir); err != nil {
		return "", err
	}
	for _, name := range projectNames(files) {
		path := filepath.Join(dir, name)
		if err := rejectWorkspaceSymlinks(dir, name, lang); err != nil {
			return "", err
		}
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return "", err
		}
		if err := writeProjectFiles(dir, map[string]string{name: files[name]}); err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, "main.go"), nil
}

func (ws *workspace) colonyModule(dir string) error {
	colonyDir := filepath.Clean(filepath.Join(ws.root, ws.colonyRel))
	files, err := moduleFiles(colonyDir, levels.Colony)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(dir, colonyDir)
	if err != nil {
		return err
	}
	body := strings.ReplaceAll(files["go.mod"], filepath.ToSlash(colonyDir), filepath.ToSlash(rel))
	return os.WriteFile(filepath.Join(dir, "go.mod"), []byte(body), 0o644)
}

func rejectWorkspaceSymlinks(dir, name, lang string) error {
	for path := filepath.Join(dir, name); ; path = filepath.Dir(path) {
		st, err := os.Lstat(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err == nil && st.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s", tr(lang, "مساحة العمل مينفعش تحتوي روابط رمزية", "Workspace cannot contain symlinks"))
		}
		if path == dir {
			return nil
		}
	}
}

func (ws *workspace) readColony(lang string) (map[string]string, time.Time, error) {
	dir := filepath.Join(ws.root, "colony")
	files := map[string]string{}
	var newest time.Time
	size := 0
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s", tr(lang, "مساحة العمل مينفعش تحتوي روابط رمزية", "Workspace cannot contain symlinks"))
		}
		if entry.IsDir() {
			return nil
		}
		name, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		name = filepath.ToSlash(name)
		if !validColonyName(name) {
			return nil
		}
		st, err := entry.Info()
		if err != nil {
			return err
		}
		if st.Size() > int64(colonyProjectLimit-size) {
			return fmt.Errorf("%s", tr(lang, "حجم ملفات المشروع أكبر من 512 KiB", "Project files exceed 512 KiB"))
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		size += len(body)
		files[name] = string(body)
		if st.ModTime().After(newest) {
			newest = st.ModTime()
		}
		return nil
	})
	if err == nil {
		err = validateColonyFiles(files, lang)
	}
	return files, newest, err
}

func workspaceLang(languages []string) string {
	if len(languages) > 0 {
		return normalizeLang(languages[0])
	}
	return "ar"
}
