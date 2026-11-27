// Command lesson-demo builds the standalone First Signal browser lesson.
package main

import (
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"orbital"
	"orbital/colony"
	"orbital/colony/levels"
)

//go:embed site/index.html site/demo.css site/demo.js
var site embed.FS

type point struct {
	X, Y int
}

type lesson struct {
	Title, Goal string
}

type mission struct {
	Width, Height, Optimal int
	Start, Beacon          point
	English, Arabic        lesson
}

func firstSignal() (mission, error) {
	level := levels.Get("1-1-first-signal")
	if level == nil || level.En == nil {
		return mission{}, errors.New("the bilingual First Signal lesson is missing")
	}
	world := levels.Start(level.ID, io.Discard, 1)
	m := mission{
		Width: world.W, Height: world.H, Optimal: level.Gold,
		Start:   point{world.X, world.Y},
		English: lesson{level.En.Title, level.En.Goal},
		Arabic:  lesson{level.Title, level.Goal},
	}
	beacons := 0
	for y := 0; y < world.H; y++ {
		for x := 0; x < world.W; x++ {
			if world.At(x, y).Entity == colony.Beacon {
				m.Beacon = point{x, y}
				beacons++
			}
		}
	}
	if beacons != 1 || m.Start == m.Beacon || m.Optimal < 1 {
		return mission{}, errors.New("First Signal must have a single target away from the starting cell")
	}
	return m, nil
}

func build(out string) error {
	m, err := firstSignal()
	if err != nil {
		return err
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	files := map[string][]byte{
		"mission.js": []byte(fmt.Sprintf("const mission = %s;\n", data)),
		".nojekyll":  nil,
		"LICENSE":    orbital.License,
	}
	for _, name := range []string{"index.html", "demo.css", "demo.js"} {
		files[name], err = site.ReadFile("site/" + name)
		if err != nil {
			return err
		}
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(out, name), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	out := flag.String("out", "dist", "directory for the static lesson")
	flag.Parse()
	if err := build(*out); err != nil {
		log.Fatal(err)
	}
}
