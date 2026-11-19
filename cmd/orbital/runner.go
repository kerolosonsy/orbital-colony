package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"orbital/colony/levels"
)

type compileError struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Col      int    `json:"col"`
	Msg      string `json:"msg"`
	Tip      string `json:"tip,omitempty"`
	TipCause string `json:"tipCause,omitempty"`
	TipCheck string `json:"tipCheck,omitempty"`
}

type result struct {
	Success  bool   `json:"success"`
	Ticks    int    `json:"ticks"`
	Stars    int    `json:"stars"`
	Par      int    `json:"par"`
	Gold     int    `json:"gold"`
	Message  string `json:"message"`
	Progress string `json:"progress,omitempty"`
}

type runResponse struct {
	Phase         string            `json:"phase"` // compile | require | run
	Traceable     bool              `json:"traceable,omitempty"`
	Failure       string            `json:"failure,omitempty"`
	CompileErrors []compileError    `json:"compileErrors,omitempty"`
	CompileOutput string            `json:"compileOutput,omitempty"`
	Missing       []string          `json:"missing,omitempty"`
	Stdout        string            `json:"stdout"`
	Stderr        string            `json:"stderr"`
	Events        []json.RawMessage `json:"events"`
	RuntimeTip    string            `json:"runtimeTip,omitempty"`
	RuntimeCause  string            `json:"runtimeTipCause,omitempty"`
	RuntimeCheck  string            `json:"runtimeTipCheck,omitempty"`
	Steps         []labStep         `json:"steps,omitempty"`
	Race          bool              `json:"race,omitempty"`
	Result        result            `json:"result"`
	Trace         *traceResponse    `json:"trace,omitempty"` // single-run debugger frames
}

type runner struct {
	goBin     string
	colonyDir string
	plainOnly bool // tests: never instrument (single run off)
}

var errLine = regexp.MustCompile(`^\.?/?([\w./-]+\.go):(\d+):(\d+): (.*)$`)

func parseCompileLine(line string, languages ...string) (compileError, bool) {
	lang := "ar"
	if len(languages) > 0 {
		lang = normalizeLang(languages[0])
	}
	m := errLine.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		return compileError{}, false
	}
	ln, _ := strconv.Atoi(m[2])
	col, _ := strconv.Atoi(m[3])
	tip := errorTipFor(m[4], tipCompile, lang)
	return compileError{File: m[1], Line: ln, Col: col, Msg: m[4], Tip: tip.Meaning, TipCause: tip.Cause, TipCheck: tip.Check}, true
}

func (r *runner) goCmd(ctx context.Context, dir string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, r.goBin, args...)
	configureProcess(cmd)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "ORBITAL_LANG="+contextLang(ctx), "GOFLAGS=-mod=mod", "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	if seed, ok := ctx.Value(reviewSeedKey{}).(int64); ok {
		cmd.Env = append(cmd.Env, "ORBITAL_SEED="+strconv.FormatInt(seed, 10))
	}
	for name := range serviceDefaults {
		cmd.Env = append(cmd.Env, name+"="+serviceEnv(name))
	}
	return cmd
}

func (r *runner) Run(ctx context.Context, lvl *levels.Level, code string) (*runResponse, error) {
	lang := contextLang(ctx)
	if missing := restrictedImports(code, lvl, lang); len(missing) > 0 {
		return rejectedImports(missing, lang), nil
	}
	dir, err := os.MkdirTemp("", "orbital-run-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	if err := r.writeModule(dir, lvl); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		return nil, err
	}

	resp := &runResponse{Phase: "compile", Events: []json.RawMessage{}}

	buildTimeout := 90 * time.Second
	if len(lvl.Deps) > 0 {
		buildTimeout = 240 * time.Second
	}
	bctx, cancel := context.WithTimeout(ctx, buildTimeout)
	defer cancel()
	bout := &limitedBuffer{max: 64 << 10, lang: contextLang(ctx)}
	args := []string{"build", "-o", "prog"}
	if lvl.Race {
		args = append(args, "-race")
	}
	build := r.goCmd(bctx, dir, append(args, ".")...)
	build.Stdout, build.Stderr = bout, bout
	// The plain build stays the judge of the learner's code (the recorder
	// reads variables, which would hide "declared and not used"); the
	// instrumented copy builds alongside it for the debugger.
	traced := r.startSingleRunBuild(bctx, lvl, code, args)
	err = runCommand(build)
	instrumentedProg := ""
	if err == nil {
		instrumentedProg = traced.wait()
	} else {
		traced.cancel()
	}
	defer traced.cleanup()
	if err != nil {
		out := strings.ReplaceAll(bout.String(), dir+string(os.PathSeparator), "")
		if errors.Is(bctx.Err(), context.DeadlineExceeded) {
			out += "\n" + context.DeadlineExceeded.Error()
			tip := errorTipFor(out, tipRuntime, lang)
			resp.RuntimeTip, resp.RuntimeCause, resp.RuntimeCheck = tip.Meaning, tip.Cause, tip.Check
		}
		resp.CompileOutput = dependencyHint(out, lang)
		for _, line := range strings.Split(out, "\n") {
			if ce, ok := parseCompileLine(line, lang); ok {
				resp.CompileErrors = append(resp.CompileErrors, ce)
			}
		}
		resp.Result.Message = tr(lang, "البرنامج مااتبناش. صلّح الأخطاء اللي تحت.", "The program did not compile. Fix the errors below.")
		resp.Failure = "compile"
		return resp, nil
	}

	// Missing concepts don't stop the run: seeing what the code actually
	// does (a race, a hang) teaches more than a refusal.
	resp.Missing = checkRequirements(code, lvl.Requires, lang)
	resp.Phase = "run"
	rctx, rcancel := context.WithTimeout(ctx, 10*time.Second)
	defer rcancel()
	trace := filepath.Join(dir, "trace.jsonl")
	prog := exec.CommandContext(rctx, filepath.Join(dir, "prog"))
	configureProcess(prog)
	prog.Dir = dir
	prog.Env = append(os.Environ(), "ORBITAL_LANG="+lang, "ORBITAL_LEVEL="+lvl.ID, "ORBITAL_TRACE="+trace)
	// Both runs get the same world: the plain program is scored, the
	// instrumented twin only supplies the debugger frames.
	seed, seeded := ctx.Value(reviewSeedKey{}).(int64)
	if !seeded && instrumentedProg != "" {
		seed, seeded = time.Now().UnixNano(), true
	}
	if seeded {
		prog.Env = append(prog.Env, "ORBITAL_SEED="+strconv.FormatInt(seed, 10))
	}
	twin := startTracedTwin(ctx, instrumentedProg, prog.Env)
	defer twin.stop()
	stdout, stderr := &limitedBuffer{max: 64 << 10, lang: contextLang(ctx)}, &limitedBuffer{max: 64 << 10, lang: contextLang(ctx)}
	prog.Stdout, prog.Stderr = stdout, stderr
	started := time.Now()
	runErr := runCommand(prog)
	plainTook := time.Since(started)
	resp.Stdout = stdout.String()
	resp.Stderr = strings.ReplaceAll(stderr.String(), dir+string(os.PathSeparator), "")
	if errors.Is(rctx.Err(), context.DeadlineExceeded) {
		if resp.Stderr != "" {
			resp.Stderr += "\n"
		}
		resp.Stderr += context.DeadlineExceeded.Error()
	}

	data, _ := os.ReadFile(trace)
	var last struct {
		Op, Msg, Progress string
		T, Par, Gold      int
		Reached           *bool
	}
	progress := ""
	mapPar, mapGold := 0, 0     // star targets of this generated map, if any
	reached, goalAt := false, 0 // goal state after the last action, and when it became true
	lostAt := 0                 // last tick the goal was true before an action broke it
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		resp.Events = append(resp.Events, json.RawMessage(append([]byte(nil), line...)))
		last.Op, last.Msg, last.Progress, last.Par, last.Gold, last.Reached = "", "", "", 0, 0, nil
		json.Unmarshal(line, &last)
		if last.Progress != "" {
			progress = last.Progress
		}
		if last.Gold > 0 {
			mapPar, mapGold = last.Par, last.Gold
		}
		if last.Reached != nil {
			if *last.Reached && !reached {
				goalAt = last.T
			}
			if !*last.Reached && reached {
				lostAt = goalAt
			}
			reached = *last.Reached
		}
	}
	// The program ended on its own right after reaching the goal: that counts.
	if last.Op != "goal" && last.Op != "fail" && runErr == nil && reached {
		ev, _ := json.Marshal(map[string]any{"op": "goal", "ok": true, "t": goalAt, "msg": tr(lang, "المهمة اكتملت!", "Mission complete!")})
		resp.Events = append(resp.Events, ev)
		last.Op, last.T, last.Par, last.Gold = "goal", goalAt, mapPar, mapGold
	}
	resp.Result.Progress = progress
	resp.Result.Ticks = last.T
	if instrumentedProg != "" {
		resp.Trace = twin.trace(plainTook, data, resp.Stdout, runErr == nil, lang)
	}
	runtimeTip := errorTipFor(resp.Stderr, tipRuntime, lang)
	resp.RuntimeTip, resp.RuntimeCause, resp.RuntimeCheck = runtimeTip.Meaning, runtimeTip.Cause, runtimeTip.Check
	resp.Race = lvl.Race && strings.Contains(resp.Stderr, "WARNING: DATA RACE")

	switch {
	case resp.Race:
		resp.Result.Message = tr(lang, "🏁 الـ Race detector لقى data race! حتى لو الرقم طلع صح، الكود مش آمن. بص على التقرير في الـ Console.", "🏁 The race detector found a data race! Even if the number is correct, the code is unsafe. Check the Console report.")
	case last.Op == "goal":
		par, gold := lvl.Par, lvl.Gold
		if last.Gold > 0 {
			par, gold = last.Par, last.Gold
		}
		resp.Result.Par, resp.Result.Gold = par, gold
		resp.Result.Success = true
		resp.Result.Stars = 1
		if par == 0 && gold == 0 {
			resp.Result.Stars = 3 // mission without a tick target
		}
		if last.T <= par {
			resp.Result.Stars = 2
		}
		if last.T <= gold {
			resp.Result.Stars = 3
		}
		resp.Result.Message = fmt.Sprintf(tr(lang, "المهمة اكتملت في %d tick!", "Mission completed in %d ticks!"), last.T)
	case last.Op == "fail":
		resp.Result.Message = last.Msg
	case errors.Is(rctx.Err(), context.DeadlineExceeded):
		resp.Result.Message = tr(lang, "البرنامج فضل شغال أكتر من 10 ثواني. غالباً عندك حلقة لا نهائية مفيهاش أوامر للدرون.", "The program ran for more than 10 seconds. Check for an endless loop with no drone commands.")
	case runErr != nil && (strings.Contains(resp.Stderr, "panic:") || strings.Contains(resp.Stderr, "fatal error:")):
		resp.Result.Message = tr(lang, "البرنامج وقع (panic). بص على الـ Console.", "The program crashed (panic). Check the Console.")
	case len(resp.Missing) > 0 && last.Op != "fail":
		resp.Result.Message = tr(lang, "البرنامج خلص، بس المهمة لسه متحققتش. وكمان المهمة دي محتاجة مفاهيم معيّنة (تحت).", "The program finished, but the mission is incomplete. It also needs the concepts listed below.")
	default:
		resp.Result.Message = tr(lang, "البرنامج خلص، بس المهمة لسه متحققتش.", "The program finished, but the mission is incomplete.")
		if lostAt > 0 {
			resp.Result.Message += fmt.Sprintf(tr(lang, " الهدف اتحقق عند tick %d، بس الكود كمّل بعدها وبوّظه. خلّي الكود يقف هناك.", " The goal was met at tick %d, but the code kept going and undid it. Make the code stop there."), lostAt)
		}
		if bytes.Contains(data, []byte(`"op":"spawn"`)) {
			resp.Result.Message += tr(lang, " لو main خلصت قبل الدرونات التانية، البرنامج كله بيقفل وهما معاه.", " If main finishes before the other drones, the whole program stops, including those drones.")
		}
	}
	if len(resp.Missing) > 0 && resp.Result.Success {
		resp.Result.Success, resp.Result.Stars = false, 0
		resp.Result.Message = tr(lang, "الدرون خلّص المهمة 👏 بس لازم تحلها باستخدام المفاهيم المطلوبة:", "The drone completed the mission 👏 but your solution must use the required concepts:")
	}
	if !resp.Result.Success {
		if resp.Race || errors.Is(rctx.Err(), context.DeadlineExceeded) || runErr != nil && (strings.Contains(resp.Stderr, "panic:") || strings.Contains(resp.Stderr, "fatal error:")) {
			resp.Failure = "runtime"
		} else {
			resp.Failure = "check"
		}
	}
	return resp, nil
}

type limitedBuffer struct {
	buf       bytes.Buffer
	tail      []byte
	max       int
	truncated bool
	lang      string
}

func (l *limitedBuffer) Write(p []byte) (int, error) {
	n := len(p)
	if !l.truncated && n <= l.max-l.buf.Len() {
		return l.buf.Write(p)
	}
	tailSize := l.max / 4
	if !l.truncated {
		headSize := l.max - tailSize
		if l.buf.Len() < headSize {
			take := min(headSize-l.buf.Len(), len(p))
			l.buf.Write(p[:take])
			p = p[take:]
		}
		l.tail = append(l.tail, l.buf.Bytes()[headSize:]...)
		l.buf.Truncate(headSize)
		l.truncated = true
	}
	if len(p) >= tailSize {
		l.tail = append(l.tail[:0], p[len(p)-tailSize:]...)
	} else {
		l.tail = append(l.tail, p...)
		if len(l.tail) > tailSize {
			l.tail = l.tail[len(l.tail)-tailSize:]
		}
	}
	return n, nil
}

func (l *limitedBuffer) String() string {
	lang := l.lang
	if l.truncated {
		return l.buf.String() + tr(lang, "\n…(الناتج اتقص)\n", "\n…(output truncated)\n") + string(l.tail)
	}
	return l.buf.String()
}
