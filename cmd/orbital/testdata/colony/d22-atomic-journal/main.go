package main

import (
	"os"
	"path/filepath"

	"orbital/colony"
)

func save(dir string, body []byte) error {
	f, err := os.CreateTemp(dir, ".plan-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(body); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), filepath.Join(dir, "plan.txt"))
}

func main() {
	dir := os.Getenv("ORBITAL_COLONY_DIR")
	oldPlan, err := os.ReadFile(filepath.Join(dir, "plan-old.txt"))
	if err != nil {
		return
	}
	newPlan, err := os.ReadFile(filepath.Join(dir, "plan-new.txt"))
	if err != nil {
		return
	}
	for {
		if save(dir, oldPlan) != nil {
			return
		}
		colony.Wait()
		if save(dir, newPlan) != nil {
			return
		}
		colony.Wait()
	}
}
