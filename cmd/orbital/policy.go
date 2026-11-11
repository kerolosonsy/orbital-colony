package main

import (
	"go/parser"
	"go/token"
	"strconv"
)

// Graded programs may use the lesson API and data-processing packages,
// but cannot open grading files, import the world implementation or use unsafe.
var playerImports = map[string]bool{
	"orbital/colony": true,
	"bytes":          true, "cmp": true, "container/heap": true,
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

func restrictedImports(code string, languages ...string) []string {
	lang := "ar"
	if len(languages) > 0 {
		lang = normalizeLang(languages[0])
	}
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", code, parser.ImportsOnly)
	if err != nil {
		return nil // Syntax errors are reported by the compiler.
	}
	var missing []string
	for _, spec := range f.Imports {
		name, err := strconv.Unquote(spec.Path.Value)
		if err == nil && !playerImports[name] {
			missing = append(missing, tr(lang, "المكتبة ", "The package ")+name+tr(lang, " مش متاحة في الحلول المُقيّمة. استخدم مكتبات الدرس والحساب ومعالجة البيانات؛ الوصول لملفات النظام وداخل اللعبة مقفول لحماية التقييم.", " is unavailable in graded solutions. Use lesson, math, and data-processing packages; filesystem and game internals are blocked to protect grading."))
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
		Phase: "require", Missing: missing,
		Result: result{Message: tr(lang, "الكود بيستورد مكتبات غير متاحة في الحلول المُقيّمة.", "The code imports packages that are unavailable in graded solutions.")},
	}
}
