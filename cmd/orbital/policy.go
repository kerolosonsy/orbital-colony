package main

import (
	"go/build"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"orbital/colony/levels"
)

// Graded programs may use the lesson API and data-processing packages,
// but cannot open grading files, import the world implementation or use unsafe.
var missionImports = map[string]bool{
	"orbital/colony":        true,
	"orbital/colony/labkit": true,
	"bufio":                 true, "io/fs": true, "path": true, "html": true, "text/template": true, "log/slog": true, "crypto/sha256": true, "crypto/rand": true, "time/tzdata": true,
	"bytes": true, "cmp": true, "container/heap": true,
	"container/list": true, "container/ring": true, "context": true,
	"encoding/base64": true, "encoding/binary": true, "encoding/csv": true,
	"encoding/hex": true, "encoding/json": true, "errors": true, "fmt": true,
	"hash": true, "hash/adler32": true, "hash/crc32": true, "hash/crc64": true,
	"hash/fnv": true, "io": true, "iter": true, "maps": true,
	"math": true, "math/big": true, "math/bits": true, "math/cmplx": true,
	"math/rand": true, "math/rand/v2": true, "regexp": true, "slices": true,
	"sort": true, "strconv": true, "strings": true, "sync": true,
	"sync/atomic": true, "testing": true, "time": true,
	"unicode": true, "unicode/utf8": true, "unicode/utf16": true,
}

var labDenied = map[string]bool{"unsafe": true, "os/exec": true, "syscall": true, "plugin": true}

func restrictedImports(code string, lvl *levels.Level, lang string) []string {
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", code, parser.ImportsOnly)
	if err != nil {
		return nil
	}
	var missing []string
	for _, spec := range f.Imports {
		name, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		allowed := missionImports[name]
		if lvl.Kind == "lab" {
			pkg, err := build.Default.Import(name, "", build.FindOnly)
			allowed = strings.HasPrefix(name, "player/") || name == "orbital/colony" || name == "orbital/colony/labkit" || err == nil && pkg.Goroot && !labDenied[name]
		}
		for _, dep := range lvl.Deps {
			if strings.Contains(strings.Split(name, "/")[0], ".") && (name == dep || strings.HasPrefix(name, dep+"/")) {
				allowed = true
			}
		}
		if !allowed {
			missing = append(missing, tr(lang, "المكتبة ", "The package ")+name+tr(lang, " مش متاحة للمهمة دي.", " is unavailable for this level."))
		}
	}
	return missing
}

func rejectedImports(missing []string, languages ...string) *runResponse {
	lang := "ar"
	if len(languages) > 0 {
		lang = normalizeLang(languages[0])
	}
	return &runResponse{
		Phase: "require", Failure: "check", Missing: missing,
		Result: result{Message: tr(lang, "الكود بيستورد مكتبات غير متاحة في الحلول المُقيّمة.", "The code imports packages that are unavailable in graded solutions.")},
	}
}
