package levels_test

import (
	"bytes"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"reflect"
	"strings"
	"testing"

	"orbital/colony/econ"
	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
	"orbital/colony/levels"
)

func assertEnglish(t *testing.T, field, text string) {
	t.Helper()
	for _, r := range text {
		if r >= 0x0600 && r <= 0x06ff || r >= 0x0750 && r <= 0x077f || r >= 0xfb50 && r <= 0xfdff || r >= 0xfe70 && r <= 0xfeff {
			t.Errorf("%s contains Arabic character %U", field, r)
			return
		}
	}
}

func sourceTokens(t *testing.T, source string) []string {
	t.Helper()
	fs := token.NewFileSet()
	if _, err := parser.ParseFile(fs, "starter.go", source, parser.AllErrors); err != nil {
		t.Fatalf("starter does not parse: %v", err)
	}
	var scan scanner.Scanner
	scan.Init(fs.AddFile("tokens.go", -1, len(source)), []byte(source), nil, 0)
	var tokens []string
	for {
		_, kind, literal := scan.Scan()
		if kind == token.EOF {
			return tokens
		}
		if kind == token.STRING {
			literal = ""
		}
		tokens = append(tokens, kind.String()+":"+literal)
	}
}

func assertSameCode(t *testing.T, arabic, english string) {
	t.Helper()
	if !reflect.DeepEqual(sourceTokens(t, arabic), sourceTokens(t, english)) {
		t.Error("English starter changes Go tokens beyond comments and string contents")
	}
}

func TestEnglishLearningContent(t *testing.T) {
	missions := append(append([]*levels.Level{}, levels.All...), levels.Colony)
	for _, mission := range missions {
		t.Run(mission.ID, func(t *testing.T) {
			en := mission.En
			if en == nil {
				t.Fatal("missing English content")
			}
			for field, text := range map[string]string{"title": en.Title, "goal": en.Goal, "lesson": en.Lesson, "starter": en.Starter} {
				if text == "" && (field != "starter" || mission.Starter != "") {
					t.Errorf("English %s is empty", field)
				}
				assertEnglish(t, field, text)
			}
			if len(en.Hints) != len(mission.Hints) {
				t.Errorf("hint count = %d, want %d", len(en.Hints), len(mission.Hints))
			}
			for _, hint := range en.Hints {
				assertEnglish(t, "hint", hint)
			}
			if mission.Starter != "" {
				assertSameCode(t, mission.Starter, en.Starter)
			}
			if mission.Kind == "lab" && len(en.Files) == 0 {
				t.Fatal("English lab files are missing")
			}
			if len(en.Files) != len(mission.Files) {
				t.Fatalf("file count = %d, want %d", len(en.Files), len(mission.Files))
			}
			for i, file := range en.Files {
				original := mission.Files[i]
				if file.Name != original.Name || file.ReadOnly != original.ReadOnly {
					t.Errorf("file %d name/ReadOnly differs from original", i)
				}
				assertEnglish(t, "file name", file.Name)
				assertEnglish(t, file.Name, file.Starter)
				t.Run(file.Name, func(t *testing.T) {
					if strings.HasSuffix(file.Name, ".go") {
						assertSameCode(t, original.Starter, file.Starter)
					} else if original.Starter != file.Starter {
						t.Error("non-Go lab assets differ between languages")
					}
				})
			}
		})
	}
	for _, chapter := range levels.Chapters {
		if chapter.TitleEn == "" || chapter.ConceptEn == "" {
			t.Errorf("chapter %d has incomplete English metadata", chapter.N)
		}
		assertEnglish(t, "chapter title", chapter.TitleEn)
		assertEnglish(t, "chapter concept", chapter.ConceptEn)
	}
	for _, upgrade := range econ.Catalog {
		t.Run(upgrade.ID, func(t *testing.T) {
			for field, text := range map[string]string{"nameEn": upgrade.NameEn, "descEn": upgrade.DescEn, "groupEn": upgrade.GroupEn} {
				if text == "" {
					t.Errorf("%s is empty", field)
				}
				assertEnglish(t, field, text)
			}
		})
	}
}

func TestMissionFeedbackUsesSelectedLanguage(t *testing.T) {
	previous := i18n.Lang()
	t.Cleanup(func() { i18n.SetLang(previous) })
	for _, mission := range append(append([]*levels.Level{}, levels.All...), levels.Colony) {
		t.Run(mission.ID, func(t *testing.T) {
			world := sim.New(io.Discard, 1)
			mission.Setup(world)
			i18n.SetLang("ar")
			arabicDone, arabic := world.Check(world)
			i18n.SetLang("en")
			englishDone, english := world.Check(world)
			if arabicDone != englishDone {
				t.Error("language changes mission completion")
			}
			if english == "" || english == arabic {
				t.Errorf("feedback was not translated: %q", english)
			}
			assertEnglish(t, "mission feedback", english)
			var console bytes.Buffer
			world.Console = &console
			levels.Narrate(world)
			if got, want := console.String(), "🛰️  "+english+"\n"; got != want {
				t.Errorf("narration = %q, want %q", got, want)
			}
		})
	}
}
