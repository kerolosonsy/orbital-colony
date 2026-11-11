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
	File string `json:"file"`
	Line int    `json:"line"`
	Col  int    `json:"col"`
	Msg  string `json:"msg"`
	Tip  string `json:"tip,omitempty"`
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
	CompileErrors []compileError    `json:"compileErrors,omitempty"`
	CompileOutput string            `json:"compileOutput,omitempty"`
	Missing       []string          `json:"missing,omitempty"`
	Stdout        string            `json:"stdout"`
	Stderr        string            `json:"stderr"`
	Events        []json.RawMessage `json:"events"`
	RuntimeTip    string            `json:"runtimeTip,omitempty"`
	Steps         []labStep         `json:"steps,omitempty"`
	Race          bool              `json:"race,omitempty"`
	Result        result            `json:"result"`
}

type runner struct {
	goBin     string
	colonyDir string
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
	return compileError{m[1], ln, col, m[4], tipFor(m[4], lang)}, true
}

func compileTips(lang string) []struct{ needle, tip string } {
	return []struct{ needle, tip string }{
		{"declared and not used", tr(lang, "عرّفت متغير ومستخدمتوش. Go مبيقبلش متغيرات مش مستخدمة: استخدمه أو امسحه.", "You declared a variable without using it. Go requires variables to be used: use it or remove it.")},
		{"imported and not used", tr(lang, "عملت import لـ package ومستخدمتهاش. امسح السطر أو استخدمها.", "You imported a package without using it. Remove the import or use the package.")},
		{"undefined:", tr(lang, "الاسم ده مش معروف. راجع الإملاء، والحروف الكبيرة والصغيرة بتفرق (Move غير move).", "Unknown name. Check spelling and capitalization (Move differs from move).")},
		{"missing return", tr(lang, "الدالة دي المفروض ترجّع قيمة، بس فيه طريق بيخلص من غير return.", "This function must return a value, but one path ends without a return.")},
		{"syntax error", tr(lang, "خطأ في كتابة الكود. راجع الأقواس { } ( ) وإن كل حاجة مقفولة.", "Syntax error. Check that all braces { } and parentheses ( ) are closed.")},
		{"non-boolean condition", tr(lang, "شرط الـ if أو الـ for لازم يكون true/false. يمكن كتبت = بدل ==؟", "An if or for condition must be true or false. Did you use = instead of ==?")},
		{"cannot use", tr(lang, "النوع مش مناسب. مثلاً بتحط string مكان int.", "The type does not match, such as using a string where an int is expected.")},
		{"mismatched types", tr(lang, "بتعمل عملية على نوعين مختلفين (زي int + string). Go مش بيحوّل لوحده.", "This operation uses different types (such as int + string). Go does not convert them automatically.")},
		{"assignment mismatch", tr(lang, "عدد المتغيرات مش زي عدد القيم اللي الدالة بترجّعها.", "The number of variables does not match the number of returned values.")},
		{"no new variables on left side of :=", tr(lang, "المتغير موجود قبل كده. استخدم = بدل := عشان تغيّر قيمته.", "This variable already exists. Use = instead of := to change its value.")},
		{"is not used", tr(lang, "قيمة اتحسبت ومحدش استخدمها.", "A value was calculated but never used.")},
		{"too many arguments", tr(lang, "بعت قيم أكتر من اللي الدالة مستنياها.", "You passed more arguments than the function expects.")},
		{"not enough arguments", tr(lang, "الدالة مستنية قيم أكتر. راجع الأقواس.", "The function expects more arguments. Check the parentheses.")},
	}
}

func runTips(lang string) []struct{ needle, tip string } {
	return []struct{ needle, tip string }{
		{"WARNING: DATA RACE", tr(lang, "Race detector لقى data race: أكتر من goroutine بيقروا ويكتبوا نفس المتغير في نفس الوقت من غير حماية. التقرير فوق بيقولك السطرين. احميهم بـ sync.Mutex.", "The race detector found a data race: goroutines access the same variable concurrently without protection. The report above shows both lines. Protect access with sync.Mutex.")},
		{"concurrent map writes", tr(lang, "أكتر من goroutine بيكتبوا في نفس الـ map في نفس الوقت. الـ maps في Go مش آمنة للـ concurrency: احميها بـ sync.Mutex.", "Multiple goroutines are writing to the same map at once. Go maps need protection for concurrent access: use sync.Mutex.")},
		{"concurrent map read and map write", tr(lang, "goroutine بتقرا من map وواحدة تانية بتكتب فيها في نفس الوقت. احميها بـ sync.Mutex أو RWMutex.", "One goroutine reads a map while another writes to it. Protect it with sync.Mutex or RWMutex.")},
		{"all goroutines are asleep", tr(lang, "Deadlock: كل الـ goroutines مستنية بعض. دوّر على channel محدش بيستقبل منها، أو محدش قفلها، أو mutex مقفول ومحدش بيفتحه.", "Deadlock: all goroutines are waiting. Check for a channel with no receiver, an unclosed channel, or a mutex that is never unlocked.")},
		{"nil pointer dereference", tr(lang, "استخدمت pointer قيمته nil. اتأكد منه بـ != nil قبل ما تستخدمه.", "You used a nil pointer. Check != nil before using it.")},
		{"index out of range", tr(lang, "حاولت توصل لعنصر برّه حدود الـ slice. اتأكد إن الـ index أصغر من len.", "You accessed an element outside the slice. Check that the index is less than len.")},
		{"send on closed channel", tr(lang, "بعت على channel اتقفلت. اللي بيبعت بس هو اللي يقفل، وبعد ما يخلّص إرسال خالص.", "You sent to a closed channel. Only the sender should close it, after all sends finish.")},
		{"close of closed channel", tr(lang, "قفلت نفس الـ channel مرتين.", "You closed the same channel twice.")},
		{"negative WaitGroup counter", tr(lang, "عملت wg.Done() أكتر من wg.Add(). كل Add(1) يقابلها Done() واحدة بس.", "You called wg.Done() more often than wg.Add(). Each Add(1) needs exactly one Done().")},
	}
}

func runtimeTipFor(stderr string, languages ...string) string {
	lang := "ar"
	if len(languages) > 0 {
		lang = normalizeLang(languages[0])
	}
	for _, t := range runTips(lang) {
		if strings.Contains(stderr, t.needle) {
			return t.tip
		}
	}
	return ""
}

func tipFor(msg string, languages ...string) string {
	lang := "ar"
	if len(languages) > 0 {
		lang = normalizeLang(languages[0])
	}
	for _, t := range compileTips(lang) {
		if strings.Contains(msg, t.needle) {
			return t.tip
		}
	}
	return ""
}

func (r *runner) goCmd(ctx context.Context, dir string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, r.goBin, args...)
	configureProcess(cmd)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "ORBITAL_LANG="+contextLang(ctx), "GOFLAGS=-mod=mod", "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
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
	if err := runCommand(build); err != nil {
		out := strings.ReplaceAll(bout.String(), dir+string(os.PathSeparator), "")
		resp.CompileOutput = dependencyHint(out, lang)
		for _, line := range strings.Split(out, "\n") {
			if ce, ok := parseCompileLine(line, lang); ok {
				resp.CompileErrors = append(resp.CompileErrors, ce)
			}
		}
		resp.Result.Message = tr(lang, "البرنامج مااتبناش. صلّح الأخطاء اللي تحت.", "The program did not compile. Fix the errors below.")
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
	stdout, stderr := &limitedBuffer{max: 64 << 10, lang: contextLang(ctx)}, &limitedBuffer{max: 64 << 10, lang: contextLang(ctx)}
	prog.Stdout, prog.Stderr = stdout, stderr
	runErr := runCommand(prog)
	resp.Stdout = stdout.String()
	resp.Stderr = strings.ReplaceAll(stderr.String(), dir+string(os.PathSeparator), "")

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
	resp.RuntimeTip = runtimeTipFor(resp.Stderr, lang)
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
