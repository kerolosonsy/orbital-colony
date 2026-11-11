// Package i18n picks between the game's two languages. Every player-facing
// string in the Go code is written as T(arabic, english).
package i18n

import (
	"os"
	"sync/atomic"
)

var lang atomic.Value // "ar" or "en"

func init() {
	if os.Getenv("ORBITAL_LANG") == "en" {
		lang.Store("en")
	} else {
		lang.Store("ar")
	}
}

// Lang is the current language: "ar" (default) or "en".
func Lang() string { return lang.Load().(string) }

// SetLang switches the language for this process. Anything other than "en" means Arabic.
func SetLang(l string) {
	if l != "en" {
		l = "ar"
	}
	lang.Store(l)
}

// T returns the Arabic or the English text for the current language.
func T(ar, en string) string {
	if Lang() == "en" {
		return en
	}
	return ar
}
