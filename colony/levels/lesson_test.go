package levels

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"unicode"

	"orbital/colony/internal/sim"
)

func TestStructuredLessonsAreValidAndBilingual(t *testing.T) {
	for _, level := range All {
		if len(level.Steps) == 0 {
			wire, err := json.Marshal(level)
			if err != nil || strings.Contains(string(wire), `"steps"`) || strings.Contains(string(wire), `"lessonRevision"`) {
				t.Fatalf("%s: legacy lesson JSON changed: %s (%v)", level.ID, wire, err)
			}
			continue
		}
		t.Run(level.ID, func(t *testing.T) {
			if level.En == nil || level.LessonRevision < 1 || level.En.LessonRevision != level.LessonRevision {
				t.Fatal("missing translated lesson or matching revision")
			}
			if len(level.Hints) != 3 || len(level.En.Hints) != 3 || len(level.Steps) != len(level.En.Steps) {
				t.Fatal("steps and three hints must be bilingual")
			}
			validateLessonSteps(t, level.Steps)
			validateLessonSteps(t, level.En.Steps)
			for i, step := range level.Steps {
				translated := level.En.Steps[i]
				if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
					t.Fatalf("step %s: translated structure differs", step.ID)
				}
				for j, block := range step.Blocks {
					assertLessonParity(t, block, translated.Blocks[j])
				}
			}
			english, _ := json.Marshal(level.En.Steps)
			if strings.ContainsFunc(string(english), func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) {
				t.Error("Arabic text in English steps")
			}
		})
	}
}

func validateLessonSteps(t *testing.T, steps []LessonStep) {
	t.Helper()
	seen := map[string]bool{}
	claimID := func(id string) {
		if strings.TrimSpace(id) == "" || seen[id] {
			t.Fatalf("empty or duplicate lesson ID %q", id)
		}
		seen[id] = true
	}
	for _, step := range steps {
		claimID(step.ID)
		if strings.TrimSpace(step.Title) == "" || len(step.Blocks) == 0 {
			t.Fatalf("%s: empty step", step.ID)
		}
		for _, block := range step.Blocks {
			claimID(block.ID)
			if block.Code != "" && (len(block.Imports) == 0 || block.ExampleID == "") {
				t.Fatalf("%s: executable example missing context or identity", block.ID)
			}
			if strings.TrimSpace(block.Text) == "" {
				t.Fatalf("%s: missing text", block.ID)
			}
			payloads := 0
			for _, present := range []bool{block.Walkthrough != nil, block.Question != nil, block.Sandbox != nil} {
				if present {
					payloads++
				}
			}
			switch block.Kind {
			case "markdown":
				if payloads != 0 || block.Code != "" {
					t.Fatalf("%s: markdown contains another block payload", block.ID)
				}
			case "walkthrough":
				if payloads != 1 || block.Walkthrough == nil || block.Code == "" || len(block.Imports) == 0 {
					t.Fatalf("%s: invalid walkthrough", block.ID)
				}
				validateLessonFrames(t, block)
			case "question":
				if payloads != 1 || block.Question == nil {
					t.Fatalf("%s: invalid question", block.ID)
				}
				q := block.Question
				if len(q.Options) < 2 || q.Answer < 0 || q.Answer >= len(q.Options) || len(q.Feedback) != len(q.Options) {
					t.Fatalf("%s: invalid answers or feedback", block.ID)
				}
				for i := range q.Options {
					if strings.TrimSpace(q.Options[i]) == "" || strings.TrimSpace(q.Feedback[i]) == "" {
						t.Fatalf("%s: empty option or feedback", block.ID)
					}
				}
			case "sandbox":
				if payloads != 1 || block.Sandbox == nil {
					t.Fatalf("%s: invalid sandbox", block.ID)
				}
				validateLessonSandbox(t, block.Sandbox, claimID)
			default:
				t.Fatalf("%s: unknown kind %q", block.ID, block.Kind)
			}
		}
	}
}

func validateLessonFrames(t *testing.T, block LessonBlock) {
	t.Helper()
	frames := block.Walkthrough.Frames
	if len(frames) < 2 || frames[0].Line != 0 {
		t.Fatalf("%s: missing initial state", block.ID)
	}
	lines := strings.Count(block.Code, "\n") + 1
	for i, frame := range frames {
		if frame.Line < 0 || frame.Line > lines || (i > 0 && frame.Line == 0) || frame.Caption == "" {
			t.Fatalf("%s frame %d: invalid line reference/caption", block.ID, i)
		}
		if frame.Panic != "" && i != len(frames)-1 {
			t.Fatalf("%s: execution continues after panic", block.ID)
		}
		if i > 0 && !strings.HasPrefix(frame.Output, frames[i-1].Output) {
			t.Fatalf("%s: output is not cumulative", block.ID)
		}
		if slice := frame.Slice; slice != nil {
			if slice.Name == "" || slice.Length != len(slice.Values) || (slice.Capacity != nil && *slice.Capacity < slice.Length) || (slice.Nil && (slice.Length != 0 || slice.Capacity == nil || *slice.Capacity != 0)) {
				t.Fatalf("%s frame %d: impossible slice state", block.ID, i)
			}
			if slice.ElementsOnly && (slice.Capacity != nil || slice.Store != "" || slice.Nil) {
				t.Fatalf("%s frame %d: introductory view contains storage details", block.ID, i)
			}
		}
	}
}

func validateLessonSandbox(t *testing.T, sandbox *LessonSandbox, claimID func(string)) {
	t.Helper()
	if sandbox.Kind == "sort-compare" {
		if len(sandbox.Sequence) != 5 || len(sandbox.Grid) != 0 || len(sandbox.Presets) != 0 {
			t.Fatal("sort sandbox must have five illustrative sizes")
		}
		for _, label := range []string{sandbox.SourceLabel, sandbox.SortDirectionLabel, sandbox.SortTieLabel, sandbox.AscendingLabel, sandbox.DescendingLabel, sandbox.TieOnLabel, sandbox.TieOffLabel, sandbox.ResetLabel} {
			if strings.TrimSpace(label) == "" {
				t.Fatal("sort sandbox missing a translated label")
			}
		}
		return
	}
	if sandbox.Kind == "grid-index" {
		if len(sandbox.Grid) != 2 || len(sandbox.Grid[0]) != 3 || len(sandbox.Grid[1]) != 3 || len(sandbox.Sequence) != 0 || len(sandbox.Presets) != 0 {
			t.Fatal("grid sandbox must have an illustrative 2-by-3 grid")
		}
		for _, label := range []string{sandbox.SourceLabel, sandbox.GridOuterCloneLabel, sandbox.GridRowCloneLabel, sandbox.GridEditCopyLabel, sandbox.ResetLabel} {
			if strings.TrimSpace(label) == "" {
				t.Fatal("grid sandbox missing a translated label")
			}
		}
		return
	}
	if sandbox.Kind == "slice-append-sequence" {
		if len(sandbox.Sequence) != 8 || len(sandbox.Presets) != 0 || sandbox.Limit != 0 || sandbox.AppendValue != 0 || sandbox.GrowthLabel != "" || sandbox.ReservedLabel != "" || sandbox.LimitLabel != "" || sandbox.AppendText != "" || sandbox.ResetText != "" {
			t.Fatal("sequence sandbox must contain eight sizes and no make/storage controls")
		}
		for _, size := range sandbox.Sequence {
			if size < 1 || size > 9 {
				t.Fatal("invalid illustrative crystal size")
			}
		}
		for _, label := range []string{sandbox.SourceLabel, sandbox.AppendLabel, sandbox.ResetLabel, sandbox.DoneText} {
			if strings.TrimSpace(label) == "" {
				t.Fatal("sequence sandbox missing a label")
			}
		}
		return
	}
	if sandbox.Kind != "slice-make" || len(sandbox.Presets) != 3 || sandbox.Limit != 6 || sandbox.AppendValue != 7 {
		t.Fatal("sandbox must implement the approved bounded make playground")
	}
	for _, text := range []string{sandbox.GrowthLabel, sandbox.ReservedLabel, sandbox.LimitLabel, sandbox.AppendText, sandbox.ResetText} {
		if strings.TrimSpace(text) == "" {
			t.Fatal("sandbox missing explanation or state feedback")
		}
	}
	for _, preset := range sandbox.Presets {
		claimID(preset.ID)
		if preset.Label == "" || preset.Text == "" || preset.Capacity < len(preset.Values) || len(preset.Values) > sandbox.Limit {
			t.Fatal("invalid sandbox preset")
		}
	}
}

func assertLessonParity(t *testing.T, arabic, english LessonBlock) {
	t.Helper()
	if arabic.ID != english.ID || arabic.Kind != english.Kind || arabic.Optional != english.Optional || arabic.Code != english.Code || arabic.ExampleID != english.ExampleID || !reflect.DeepEqual(arabic.Imports, english.Imports) {
		t.Fatalf("%s: block identity/code differs across languages", arabic.ID)
	}
	if arabic.Question != nil {
		if english.Question == nil || arabic.Question.Answer != english.Question.Answer || len(arabic.Question.Options) != len(english.Question.Options) {
			t.Fatalf("%s: translated answers differ", arabic.ID)
		}
	}
	if arabic.Walkthrough != nil {
		if english.Walkthrough == nil || len(arabic.Walkthrough.Frames) != len(english.Walkthrough.Frames) {
			t.Fatalf("%s: frame counts differ", arabic.ID)
		}
		for i, frame := range arabic.Walkthrough.Frames {
			translated := english.Walkthrough.Frames[i]
			frame.Caption, translated.Caption = "", ""
			if !reflect.DeepEqual(frame, translated) {
				t.Fatalf("%s frame %d: translated execution differs", arabic.ID, i)
			}
		}
	}
	if arabic.Sandbox != nil {
		translated := english.Sandbox
		if translated == nil || len(arabic.Sandbox.Presets) != len(translated.Presets) || arabic.Sandbox.Kind != translated.Kind || arabic.Sandbox.Limit != translated.Limit || arabic.Sandbox.AppendValue != translated.AppendValue || !reflect.DeepEqual(arabic.Sandbox.Sequence, translated.Sequence) || !reflect.DeepEqual(arabic.Sandbox.Grid, translated.Grid) {
			t.Fatalf("%s: sandbox behaviour differs", arabic.ID)
		}
		for i, preset := range arabic.Sandbox.Presets {
			other := translated.Presets[i]
			preset.Text, other.Text = "", ""
			if !reflect.DeepEqual(preset, other) {
				t.Fatalf("%s: translated preset state differs", arabic.ID)
			}
		}
	}
}

func TestSurveyLessonStaysShortAndVisual(t *testing.T) {
	l := Get("6-1-survey")
	if l.LessonRevision != 2 {
		t.Fatal("survey must invalidate revision-1 practice")
	}
	for _, copy := range []struct {
		steps []LessonStep
		hints []string
	}{{l.Steps, l.Hints}, {l.En.Steps, l.En.Hints}} {
		words, questions, traces, sandboxes := len(strings.Fields(strings.Join(copy.hints, " "))), 0, 0, 0
		if len(copy.steps) != 6 {
			t.Fatal("survey must have six short steps")
		}
		for i, step := range copy.steps {
			words += len(strings.Fields(step.Title))
			for _, block := range step.Blocks {
				words += len(strings.Fields(block.Text))
				if block.Note != "" || block.Optional {
					t.Fatal("short lesson must not hide extra reading")
				}
				if i == 0 && len(strings.Fields(block.Text)) > 45 {
					t.Fatal("mission exceeds 45 words")
				}
				if block.Walkthrough != nil {
					traces++
					for _, frame := range block.Walkthrough.Frames {
						count := len(strings.Fields(frame.Caption))
						words += count
						if count > 8 || frame.Slice != nil && !frame.Slice.ElementsOnly {
							t.Fatal("trace must use short captions and element-only views")
						}
					}
				}
				if block.Question != nil {
					questions++
					words += len(strings.Fields(strings.Join(block.Question.Options, " ") + " " + strings.Join(block.Question.Feedback, " ")))
				}
				if b := block.Sandbox; b != nil {
					sandboxes++
					if b.Kind != "slice-append-sequence" {
						t.Fatal("survey must use the growing log playground")
					}
					words += len(strings.Fields(b.SourceLabel + " " + b.AppendLabel + " " + b.ResetLabel + " " + b.DoneText))
				}
			}
		}
		if words < 250 || words > 300 || questions != 3 || traces != 1 || sandboxes != 1 {
			t.Fatalf("words=%d questions=%d walkthroughs=%d sandboxes=%d", words, questions, traces, sandboxes)
		}
		t.Logf("learner-authored words, including all feedback and hints: %d", words)
	}
}

func TestSurveyLessonPreservesMissionContract(t *testing.T) {
	l := Get("6-1-survey")
	if l.Par != 15 || l.Gold != 15 || l.MaxTicks != 60 || !reflect.DeepEqual(l.Requires, []string{"slice"}) {
		t.Fatal("lesson changed mission grading")
	}
	for source, hash := range map[string]string{l.Starter: "b4698aef95df13e6c3d72cf4a0fac9e8ffbf1a39da4741b36346a9178dbe1075", l.En.Starter: "5ca3ea052fb17a7fd4191d5186942b8e5cf76ab7fac7a443113effdd5c0a2c82"} {
		if fmt.Sprintf("%x", sha256.Sum256([]byte(source))) != hash {
			t.Fatal("lesson changed starter code")
		}
	}
	for _, seed := range []int64{1, 19, 2026} {
		w := sim.New(io.Discard, seed)
		l.Setup(w)
		if w.W != 8 || w.H != 2 {
			t.Fatal("survey dimensions changed")
		}
		sizes := []int{}
		for x := 0; x < 8; x++ {
			tile := w.At(x, 0)
			if tile.Entity != sim.Crystal || tile.Size < 1 || tile.Size > 9 {
				t.Fatal("survey crystal generation changed")
			}
			sizes = append(sizes, tile.Size)
		}
		w.Reports = []string{"incorrect"}
		if ok, _ := w.Check(w); ok {
			t.Fatal("incorrect report accepted")
		}
		w.Reports = []string{fmt.Sprint(sizes)}
		if ok, _ := w.Check(w); !ok {
			t.Fatal("correct ordered report rejected")
		}
	}
}
