package levels

import (
	"context"
	"errors"
	"fmt"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{
			Title: "Recall Order",
			Goal:  "Work with 3 drones until base cancels the mission (ctx.Done), then return everyone to (0, 0) in time",
			Lesson: `# context: Tell Every Task to Stop 📣

When a large operation is canceled, you need to notify **all** its goroutines and the functions they call. That is exactly what ~context.Context~ does.

~~~go
ctx := colony.Mission() // Base will cancel it at an unknown time

select {
case <-ctx.Done():
	// Canceled! ctx.Err() tells you why
	return
default:
	// Continue
}
~~~

- ~ctx.Done()~: a channel that closes when the context is canceled. (Remember that a closed channel signals everyone?)
- ~ctx.Err()~: ~context.Canceled~ or ~context.DeadlineExceeded~.

## The golden rule: ctx is the first parameter

~~~go
func work(ctx context.Context, d *colony.Drone, row int)
~~~

In real Go, functions that can take time (HTTP requests, database queries, file reads...) take ~ctx~ **as their first parameter** and pass it to the functions they call. Cancellation can then travel through the whole chain.

## Do not store ctx in a struct
Pass it with each call. This is the official Go convention.

## Your mission
3 drones are harvesting. Base cancels ctx at a random time, and you have **14 ticks** to get every drone to (0, 0). You must also collect the required amount before then.`,
			Starter: `package main

import (
	"context"
	"fmt"
	"orbital/colony"
	"sync"
)

// work harvests row until ctx is canceled, then returns to base
func work(ctx context.Context, d *colony.Drone, row int) {
	for i := 0; i < row; i++ {
		d.Move(colony.North)
	}
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Drone", d.ID(), "returning:", ctx.Err())
			// Return to (0, 0)
			return
		default:
		}
		d.Harvest()
		d.Move(colony.East)
	}
}

func main() {
	ctx := colony.Mission()
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	for i, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			work(ctx, d, i*2)
		}()
	}
	wg.Wait()
}
`,
			Hints: []string{
				"To return: West until x = 0, then South until y = 0. Use ~d.Position()~.",
				"When a row ends (Move returns false at the wall), move up to the next row or turn back West.",
				"You need a certain crystal count, so do not return early. Work until ctx is actually canceled.",
			},
		},
		ID: "13-1-recall", Chapter: 13, Title: "أمر الرجوع",
		Goal: "اشتغل بـ 3 درونات لحد ما القاعدة تلغي المهمة (ctx.Done)، وبعدها الكل يرجع (0, 0) في الوقت",
		Par:  0, Gold: 0, MaxTicks: 400,
		Requires: []string{"sel:context.Context"},
		Setup: func(w *sim.World) {
			w.Init(9, 5, sim.Regolith)
			w.Realtime, w.MaxDrones = true, 3
			w.At(0, 0).Entity = sim.Beacon
			w.Each(func(x, y int, t *sim.Tile) {
				if x+y > 0 {
					t.Entity, t.Size = sim.Crystal, 1
				}
			})
			at := 30 + w.Rand.Intn(16)
			w.Recall = sim.NewRecall(at, 14)
			need := at * 3 / 2 * 55 / 100
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Crystal]
				msg := fmt.Sprintf(i18n.T("كريستال: %d / %d", "Crystal: %d / %d"), n, need)
				if !w.Recall.Over {
					return false, msg + i18n.T(" · لسه مفيش أمر رجوع", " · No recall order yet")
				}
				return n >= need, msg
			}
		},
		Starter: `package main

import (
	"context"
	"fmt"
	"orbital/colony"
	"sync"
)

// work بيحصد صف row لحد ما ctx يتلغي، وبعدين يرجع للقاعدة
func work(ctx context.Context, d *colony.Drone, row int) {
	for i := 0; i < row; i++ {
		d.Move(colony.North)
	}
	for {
		select {
		case <-ctx.Done():
			fmt.Println("درون", d.ID(), "راجع:", ctx.Err())
			// ارجع لـ (0, 0)
			return
		default:
		}
		d.Harvest()
		d.Move(colony.East)
	}
}

func main() {
	ctx := colony.Mission()
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	for i, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			work(ctx, d, i*2)
		}()
	}
	wg.Wait()
}
`,
		Hints: []string{
			"الرجوع: West لحد x = 0، وبعدين South لحد y = 0. محتاج ~d.Position()~.",
			"لو الصف خلص (Move بترجع false عند الحيطة)، اطلع للصف اللي بعده، أو لف راجع West.",
			"المطلوب عدد كريستال معيّن، فمترجعش بدري. اشتغل لحد ما الـ ctx يتلغي بالظبط.",
		},
		Lesson: `# context: إزاي تقول لكل الشغل "بطّل" 📣

لما مهمة كبيرة اتلغت، محتاج تبلّغ **كل** الـ goroutines اللي شغالة فيها، وكل الدوال اللي هما جوّاها. ده بالظبط شغل ~context.Context~.

~~~go
ctx := colony.Mission() // القاعدة هتلغيه في وقت مش معروف

select {
case <-ctx.Done():
	// اتلغى! ctx.Err() بيقولك ليه
	return
default:
	// كمّل
}
~~~

- ~ctx.Done()~: channel بتتقفل لما الـ context يتلغي. (فاكر إن الـ channel المقفولة إشارة للكل؟)
- ~ctx.Err()~: ~context.Canceled~ أو ~context.DeadlineExceeded~.

## القاعدة الذهبية: ctx أول parameter

~~~go
func work(ctx context.Context, d *colony.Drone, row int)
~~~

في Go الحقيقي، أي دالة ممكن تاخد وقت (طلب HTTP، query على database، قراية ملف…) بتاخد ~ctx~ **كأول parameter**، وبتعدّيه لكل حاجة بتناديها. كده الإلغاء بيوصل لآخر السلسلة.

## متخزّنش الـ ctx في struct
عدّيه في كل استدعاء. ده الـ convention الرسمي في Go.

## المهمة
3 درونات بيحصدوا. في وقت عشوائي القاعدة هتلغي الـ ctx، وقدامك **14 tick** تكون كل الدرونات في (0, 0). ولازم تكون جمعت عدد معيّن قبلها.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "Deep Scan",
			Goal:  "DeepScan each crystal and harvest only those with at least 50% purity. Some scans hang forever!",
			Lesson: `# Timeouts: Do Not Wait Forever ⏱️

Some scans hang. Wait for them and the whole mission stops. In real life this might be an unresponsive server, a failed database, or a slow network. **Operations outside your program need timeouts.**

## context.WithTimeout

~~~go
sctx, cancel := context.WithTimeout(ctx, colony.Ticks(5))
defer cancel() // Always needed, even if the scan finishes early

purity, err := colony.DeepScan(sctx)
if errors.Is(err, context.DeadlineExceeded) {
	// It hung, so it was canceled after 5 ticks
}
~~~

- Returns a **child** context: it is canceled if its parent is canceled **or** its time expires.
- ~cancel~ releases resources held by the context. Forgetting it leaks resources (~go vet~ warns you).

## A context tree 🌳

~~~text
Background
 └── WithTimeout(5 ticks)     ← For each scan
~~~

Canceling a parent cancels all its children; canceling a child does not affect its parent. Each scan can have its own timeout without canceling the mission.

## Choose timeouts thoughtfully ⚖️
- Too short, and you cancel healthy scans that would have finished.
- Too long, and you waste time on stuck scans.

Healthy scans take 2 to 4 ticks.`,
			Starter: `package main

import (
	"context"
	"fmt"
	"orbital/colony"
)

func main() {
	ctx := context.Background()
	for x := 0; x < 10; x++ {
		purity, err := colony.DeepScan(ctx) // ⚠️ May hang forever
		fmt.Println("x =", x, "purity", purity, err)
		if err == nil && purity >= 50 {
			colony.Harvest()
		}
		if x < 9 {
			colony.Move(colony.East)
		}
	}
}
`,
			Hints: []string{
				"A healthy scan takes at most 4 ticks. ~colony.Ticks(5)~ converts ticks to time.Duration.",
				"For each crystal: ~sctx, cancel := context.WithTimeout(ctx, colony.Ticks(5))~, then ~colony.DeepScan(sctx)~, then ~cancel()~.",
				"If ~errors.Is(err, context.DeadlineExceeded)~, the scan hung: skip that crystal and continue.",
			},
		},
		ID: "13-2-timeouts", Chapter: 13, Title: "الفحص العميق",
		Goal: "افحص كل كريستال بـ DeepScan، واحصد بس اللي نقاؤه 50% أو أكتر. وفيه فحوصات بتعلّق للأبد!",
		Par:  70, Gold: 50, MaxTicks: 400,
		Requires: []string{"sel:context.WithTimeout"},
		Setup: func(w *sim.World) {
			w.Init(10, 2, sim.Regolith)
			w.Realtime = true
			w.Scans = map[sim.Point]sim.ScanInfo{}
			var good, bad []sim.Point
			cost := 9
			for x := 0; x < 10; x++ {
				t := w.At(x, 0)
				t.Entity, t.Size = sim.Crystal, 1
				info := sim.ScanInfo{Ticks: 2 + w.Rand.Intn(3), Purity: w.Rand.Intn(100)}
				if x == 0 || x == 7 {
					info.Purity = 60 + w.Rand.Intn(40)
				}
				if x == 3 || x == 8 || w.Rand.Intn(4) == 0 && x != 0 && x != 7 {
					info.Ticks = -1
				}
				p := sim.Point{X: x}
				w.Scans[p] = info
				switch {
				case info.Ticks < 0:
					bad = append(bad, p)
					cost += 5
				case info.Purity >= 50:
					good = append(good, p)
					cost += info.Ticks + 1
				default:
					bad = append(bad, p)
					cost += info.Ticks
				}
			}
			w.Gold = cost + 6 // DeepScan ends on wall-clock time, so allow some drift
			w.Par = cost + 24
			w.FailCheck = func(w *sim.World) string {
				for _, p := range bad {
					if w.At(p.X, p.Y).Entity != sim.Crystal {
						if w.Scans[p].Ticks < 0 {
							return fmt.Sprintf(i18n.T("حصدت الكريستال عند x=%d من غير ما فحصه يخلص. ممكن يكون ملوّث! المهمة فشلت.", "You harvested crystal at x=%d before its scan finished. It might be contaminated! Mission failed."), p.X)
						}
						return fmt.Sprintf(i18n.T("حصدت كريستال نقاؤه %d%% بس عند x=%d. المهمة فشلت.", "You harvested crystal with only %d%% purity at x=%d. Mission failed."), w.Scans[p].Purity, p.X)
					}
				}
				return ""
			}
			w.Check = func(w *sim.World) (bool, string) {
				left := remaining(w, good)
				return left == 0, fmt.Sprintf(i18n.T("كريستال نقي فاضل: %d / %d", "Pure crystals remaining: %d / %d"), left, len(good))
			}
		},
		Starter: `package main

import (
	"context"
	"fmt"
	"orbital/colony"
)

func main() {
	ctx := context.Background()
	for x := 0; x < 10; x++ {
		purity, err := colony.DeepScan(ctx) // ⚠️ ممكن تعلّق للأبد
		fmt.Println("x =", x, "نقاء", purity, err)
		if err == nil && purity >= 50 {
			colony.Harvest()
		}
		if x < 9 {
			colony.Move(colony.East)
		}
	}
}
`,
		Hints: []string{
			"الفحص السليم بياخد 4 tick بالكتير. ~colony.Ticks(5)~ بتحوّل الـ ticks لـ time.Duration.",
			"لكل كريستال: ~sctx, cancel := context.WithTimeout(ctx, colony.Ticks(5))~، وبعدين ~colony.DeepScan(sctx)~، وبعدها ~cancel()~.",
			"لو ~errors.Is(err, context.DeadlineExceeded)~ يبقى الفحص علّق: سيب الكريستال ده وكمّل.",
		},
		Lesson: `# timeouts: متستناش للأبد ⏱️

بعض الفحوصات بتعلّق. لو استنيتها، المهمة كلها هتقف. في الدنيا الحقيقية ده سيرفر مبيردّش، أو database واقعة، أو شبكة بطيئة. **أي حاجة برّه برنامجك لازم يكون ليها timeout.**

## context.WithTimeout

~~~go
sctx, cancel := context.WithTimeout(ctx, colony.Ticks(5))
defer cancel() // لازم دايماً، حتى لو الفحص خلص بدري

purity, err := colony.DeepScan(sctx)
if errors.Is(err, context.DeadlineExceeded) {
	// علّق، فاتلغى بعد 5 tick
}
~~~

- بيرجّع context **ابن**: بيتلغي لو الأب اتلغى، **أو** لو الوقت خلص.
- ~cancel~ بتنضّف الموارد اللي الـ context ماسكها. نسيانها = تسريب (~go vet~ بيحذّرك).

## شجرة contexts 🌳

~~~text
Background
 └── WithTimeout(5 ticks)     ← لكل فحص
~~~

إلغاء الأب بيلغي كل ولاده، لكن إلغاء ابن مش بيأثر على الأب. عشان كده كل فحص ليه timeout خاص بيه من غير ما المهمة كلها تتلغي.

## اختار الـ timeout بذكاء ⚖️
- لو قصير أوي، هتلغي فحوصات سليمة كانت هتخلص.
- لو طويل أوي، هتضيّع وقت على الفحوصات المعلّقة.

الفحص السليم بياخد من 2 لـ 4 tick.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "The Legendary Crystal",
			Goal:  "3 drones search for a size-10 crystal. Once one finds it, the others stop. Return everyone to (0, 0)",
			Lesson: `# WithCancel: First to Find It Wins 🏆

3 drones are searching for the same thing. Once one finds it, the others should **stop immediately** and return, rather than keep searching needlessly.

~~~go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
	if search(ctx, d, row) {
		cancel() // Tell everyone: we found it
	}
}()
~~~

## Check cancellation at every step

~~~go
for x := 0; x < 10; x++ {
	if ctx.Err() != nil {
		return false // Someone else found it
	}
	// ...
}
~~~

~ctx.Err()~ returns nil while the context is active. It is simpler than ~select~ when you do not need to wait for anything.

## cancel is safe
- You can call it **more than once**, from **multiple goroutines**. Later calls do nothing.
- That makes ~defer cancel()~ safe even after an earlier call.

## In real programs: errgroup
~golang.org/x/sync/errgroup~ combines a WaitGroup, a context, and the first error:

~~~go
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { return search(ctx, ...) })
err := g.Wait() // The first error cancels the other tasks
~~~

This is a useful pattern when sending several requests concurrently in real applications.`,
			Starter: `package main

import (
	"context"
	"orbital/colony"
	"sync"
)

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
}

// search explores two rows. It returns true if it finds the legendary crystal
func search(ctx context.Context, d *colony.Drone, row int) bool {
	goTo(d, colony.Point{X: 0, Y: row})
	for x := 0; x < 10; x++ {
		// If ctx is canceled: return false
		if d.Measure() == 10 {
			d.Harvest()
			return true
		}
		if x < 9 {
			d.Move(colony.East)
		}
	}
	// ... The second row
	return false
}

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	for i, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			search(ctx, d, i*2)
			goTo(d, colony.Point{})
		}()
	}
	wg.Wait()
}
`,
			Hints: []string{
				"Use ~ctx, cancel := context.WithCancel(context.Background())~. The drone that finds the crystal calls ~cancel()~.",
				"In search, before each step: ~if ctx.Err() != nil { return false }~.",
				"Complete search to explore ~row+1~ while returning West. Remember ~defer cancel()~ in main.",
			},
		},
		ID: "13-3-first-wins", Chapter: 13, Title: "الكريستال الأسطوري",
		Goal: "3 درونات بيدوّروا على كريستال حجمه 10. أول ما واحد يلاقيه، الباقي يبطّلوا. الكل يرجع (0, 0)",
		Par:  60, Gold: 40, MaxTicks: 400,
		Requires: []string{"sel:context.WithCancel"},
		Setup: func(w *sim.World) {
			w.Init(10, 6, sim.Regolith)
			w.Realtime, w.MaxDrones = true, 3
			w.Each(func(x, y int, t *sim.Tile) { t.Entity, t.Size = sim.Crystal, 1+w.Rand.Intn(5) })
			leg := sim.Point{X: w.Rand.Intn(10), Y: 1 + w.Rand.Intn(5)}
			w.At(leg.X, leg.Y).Size = 10

			// Estimate with the canonical plan: drone i snakes rows 2i and 2i+1.
			paths := make([][]sim.Point, 3)
			for i := range paths {
				for y := 2 * i; y < 2*i+2; y++ {
					for k := 0; k < 10; k++ {
						x := k
						if y%2 == 1 {
							x = 9 - k
						}
						paths[i] = append(paths[i], sim.Point{X: x, Y: y})
					}
				}
			}
			finder, idx := leg.Y/2, 0
			for j, p := range paths[finder] {
				if p == leg {
					idx = j
				}
			}
			tFind := 2*finder + 2*idx + 1 // climb, then measure+move per tile
			end := tFind + 1 + sim.Dist(leg, sim.Point{})
			for i := range paths {
				if i != finder {
					steps := (tFind - 2*i) / 2
					p := sim.Point{X: 0, Y: min(2*i, tFind)}
					if steps > 0 {
						p = paths[i][min(steps, len(paths[i])-1)]
					}
					end = max(end, tFind+1+sim.Dist(p, sim.Point{}))
				}
			}
			w.Gold = end + 3
			w.Par = w.Gold + 18
			w.Check = func(w *sim.World) (bool, string) {
				got := w.At(leg.X, leg.Y).Entity != sim.Crystal
				home, id := w.AllAtBase()
				switch {
				case !got:
					return false, i18n.T("الكريستال الأسطوري لسه متلقاش", "The legendary crystal has not been found yet")
				case !home:
					return false, fmt.Sprintf(i18n.T("لقيناه! مستنيين درون %d يرجع (0, 0)", "Found it! Waiting for drone %d to return to (0, 0)"), id)
				}
				return true, i18n.T("الكريستال الأسطوري في القاعدة", "The legendary crystal is at base")
			}
		},
		Starter: `package main

import (
	"context"
	"orbital/colony"
	"sync"
)

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
}

// search بيدوّر في صفين. بيرجّع true لو لقى الأسطوري
func search(ctx context.Context, d *colony.Drone, row int) bool {
	goTo(d, colony.Point{X: 0, Y: row})
	for x := 0; x < 10; x++ {
		// لو ctx اتلغى: return false
		if d.Measure() == 10 {
			d.Harvest()
			return true
		}
		if x < 9 {
			d.Move(colony.East)
		}
	}
	// ... الصف التاني
	return false
}

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	for i, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			search(ctx, d, i*2)
			goTo(d, colony.Point{})
		}()
	}
	wg.Wait()
}
`,
		Hints: []string{
			"~ctx, cancel := context.WithCancel(context.Background())~. اللي يلاقي الكريستال ينادي ~cancel()~.",
			"جوه search، قبل كل خطوة: ~if ctx.Err() != nil { return false }~.",
			"كمّل search تدوّر في الصف ~row+1~ وهي راجعة West. ومتنساش ~defer cancel()~ في main.",
		},
		Lesson: `# WithCancel: أول واحد يكسب 🏆

3 درونات بيدوّروا على نفس الحاجة. أول ما واحد يلاقيها، الاتنين التانيين لازم **يبطّلوا فوراً** ويرجعوا، بدل ما يكمّلوا دوَر على الفاضي.

~~~go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
	if search(ctx, d, row) {
		cancel() // بلّغ الكل: خلاص لقيناه
	}
}()
~~~

## اتأكد من الإلغاء في كل خطوة

~~~go
for x := 0; x < 10; x++ {
	if ctx.Err() != nil {
		return false // حد تاني لقاه
	}
	// ...
}
~~~

~ctx.Err()~ بيرجّع nil طول ما الـ context شغال. ده أسهل من ~select~ لما مش محتاج تستنى حاجة.

## cancel آمنة
- تقدر تناديها **أكتر من مرة**، ومن **أكتر من goroutine**. التانية والتالتة مش بيعملوا حاجة.
- عشان كده ~defer cancel()~ في الآخر دايماً آمنة.

## في الدنيا الحقيقية: errgroup
المكتبة ~golang.org/x/sync/errgroup~ بتجمع WaitGroup + context + أول error في حاجة واحدة:

~~~go
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { return search(ctx, ...) })
err := g.Wait() // أول error بيلغي الباقي
~~~

ده النمط اللي هتستخدمه في الشغل لما تبعت كذا request في نفس الوقت.`,
	})
	add(&Level{ID: "13-4-chain-of-command", Chapter: 13, Title: "سلسلة القيادة", Goal: "عدّي context عبر القيادة والوسيط والمسبار ورجّع سبب التوقف", Requires: []string{"sel:context.Context", "select"}, Bridge: []int{41}, Par: 50, Gold: 30, MaxTicks: 300,
		Starter: chainOfCommandStarter, Lesson: `# نفس الأمر لآخر السلسلة

الـ context مش إعداد ثابت للجهاز؛ هو عمر عملية واحدة. لو وسيط قطع السلسلة، المسبار هيفضل شغال بعد ما القيادة وقّفت المهمة.

## قبل ما تكتب
- مين مالك cancel ومين بس بيستقبل ctx؟
- لو الأب انتهى قبل الدخول، هل لازم نبدأ شغل أصلاً؟

## عدّي ctx كأول parameter
~~~go
func outer(ctx context.Context) error { return inner(ctx) }
func inner(ctx context.Context) error { return ctx.Err() }
~~~
استورد ~context~. ~context.Context~ عقد بيمرر الإلغاء والموعد. المثال بيرجّع حالة السياق فوراً، مش بينتظر شغل. متحفظش ctx في struct؛ كل استدعاء يستقبل السياق اللي يخصه.

## الانتظار يشارك في الإلغاء
~~~go
select {
case <-ctx.Done():
 return ctx.Err()
case answer := <-answers:
 return answer
}
~~~
~Done~ قناة بتتقفل عند الإلغاء، و~Err~ nil قبل التوقف، وبعده Canceled أو DeadlineExceeded. لو القناتين جاهزين، select تختار أي واحدة؛ متفترضش أولوية. في المهمة مفيش نتيجة شغل، فاستنى Done فقط. السياق الملغي بالفعل يرجع فوراً.

## موعد الأب لا يتمدّد
~~~go
parent,stop := context.WithTimeout(context.Background(),colony.Ticks(8))
defer stop()
child,cancel := context.WithTimeout(parent,colony.Ticks(80))
defer cancel()
deadline,ok := child.Deadline()
fmt.Println(deadline,ok)
~~~
~Deadline()~ بترجّع الموعد وbool لوجوده. رغم طلب80، الابن يقف عند موعد الأب الأقرب. ~WithCancel~ بتدي إلغاء يدوي و~WithTimeout~ بتدي موعد؛ اللي أنشأ السياق لازم ينادي cancel. الدوال الداخلية ترجّع ctx.Err من غير استبدال السبب.

## راقب حدود الاستدعاءات
~~~go
colony.TraceCommand(ctx,"relay.enter")
err := next(ctx)
colony.TraceCommand(ctx,"relay.exit")
return err
~~~
~TraceCommand~ بتسجل الاسم وحالة السياق وموعده في العالم من غير ticks. دي مراقبة لحدود الدالة؛ مش بتلغي ولا بتستنى. خليها في الدوال التلاتة زي الـ starter.

## أخطاء شائعة
~~~go
err := next(context.Background()) // فقد إلغاء الأب وموعده
~~~
ومتعملش context جديد في field أو ترجع nil بعد Done؛ كده القيادة مش هتعرف إن العملية اتوقفت.

## المهمة
كمّل ~command(ctx) → relay(ctx) → probe(ctx)~، كل واحدة تستقبل context.Context كأول parameter. سجّل بالترتيب command.enter ثم relay.enter ثم probe.enter، وبعد رجوع المسبار probe.exit ثم relay.exit ثم command.exit. ممنوع ctx في struct.

نفّذ السلسلة مرتين زي main: أب ملغي بـ WithCancel، وبعدين أب بـ WithTimeout ومدته ~colony.Ticks(colony.Reading(0))~. المسبار ينتظر Done بـ select ويرجع Err، والوسيط والقيادة يرجّعوا نفس الخطأ. سجل كل حدود الدخول والخروج حتى لو ملغي مسبقاً. بعد التحقق إن السببين Canceled وDeadlineExceeded ابعت تقرير واحد بالنص ~canceled deadline~. التحقق بيراجع الأثر والموعد وحالة السياق، مش النص لوحده.
`, Hints: []string{"مينفعش relay تعمل Background؛ عدّي ctx اللي وصلها.", "خلّي probe تستنى Done وتعيد ctx.Err؛ الخروج بيتسجّل بعد الانتظار.", "جرّب ملغي مسبقاً وبعدين موعد من Reading(0)، وارجع نفس الخطأ عبر كل المستويات."},
		En: &LevelText{Title: "Chain of Command", Goal: "Pass context through command, relay, and probe and return the stop reason", Starter: chainOfCommandStarter, Lesson: `# One Order Through the Whole Chain

A context describes one operation's lifetime, rather than permanent device configuration. If a relay breaks the chain, the probe can continue after command has stopped the mission.

## Before you write
- Who owns cancel, and who only receives ctx?
- If the parent is already done on entry, should work even begin?

## Pass ctx as the first parameter
~~~go
func outer(ctx context.Context) error { return inner(ctx) }
func inner(ctx context.Context) error { return ctx.Err() }
~~~
Import ~context~. ~context.Context~ carries cancellation and deadlines. This example immediately returns context state rather than waiting for work. Do not store ctx in a struct; receive the relevant context on each call.

## Waiting participates in cancellation
~~~go
select {
case <-ctx.Done():
 return ctx.Err()
case answer := <-answers:
 return answer
}
~~~
~Done~ closes on cancellation. ~Err~ is nil before stopping, then Canceled or DeadlineExceeded. If both channels are ready, select may choose either; do not assume priority. This mission has no work result, so wait only for Done. An already-canceled context returns immediately.

## A parent's deadline cannot be extended
~~~go
parent,stop := context.WithTimeout(context.Background(),colony.Ticks(8))
defer stop()
child,cancel := context.WithTimeout(parent,colony.Ticks(80))
defer cancel()
deadline,ok := child.Deadline()
fmt.Println(deadline,ok)
~~~
~Deadline()~ returns a time and a presence bool. Despite requesting 80 ticks, the child stops at the parent's earlier deadline. ~WithCancel~ supplies manual cancellation; ~WithTimeout~ supplies a deadline. The creator must call cancel. Inner functions return ctx.Err without replacing its reason.

## Observe call boundaries
~~~go
colony.TraceCommand(ctx,"relay.enter")
err := next(ctx)
colony.TraceCommand(ctx,"relay.exit")
return err
~~~
~TraceCommand~ records the label, context state, and deadline in the world without ticks. It observes boundaries without canceling or waiting. Keep it in all three functions as shown in the starter.

## Common mistakes
~~~go
err := next(context.Background()) // Loses the parent's cancellation and deadline
~~~
Do not put a replacement context in a field or return nil after Done; command would lose the reason the operation stopped.

## Your mission
Complete ~command(ctx) → relay(ctx) → probe(ctx)~, each taking context.Context first. Trace command.enter, relay.enter, probe.enter, then probe.exit, relay.exit, command.exit while unwinding. Do not store ctx in a struct.

Run twice as in main: an already-canceled WithCancel parent, then a WithTimeout parent lasting ~colony.Ticks(colony.Reading(0))~. The probe waits for Done with select and returns Err; relay and command return the same error. Trace every boundary even when canceled before entry. After checking for Canceled and DeadlineExceeded, send exactly one report containing ~canceled deadline~. Verification inspects traces, deadline, and context state as well as the text.
`, Hints: []string{"Relay must pass its received ctx instead of creating Background.", "Have probe wait for Done and return ctx.Err; trace exit after waiting.", "Exercise already-canceled and Reading(0) deadline parents, returning the error through every level."}}, Setup: func(w *sim.World) {
			w.Init(4, 2, sim.Regolith)
			w.Realtime = true
			w.Data.SensorReadings = []int{15 + w.Rand.Intn(10)}
			w.Check = func(w *sim.World) (bool, string) {
				stages := []string{"command.enter", "relay.enter", "probe.enter", "probe.exit", "relay.exit", "command.exit"}
				valid := len(w.Data.CommandLog) == 12
				if valid {
					deadline := w.Data.CommandLog[6].Deadline
					for i, record := range w.Data.CommandLog {
						valid = valid && record.Stage == stages[i%6]
						if i < 6 {
							valid = valid && errors.Is(record.Err, context.Canceled)
						} else {
							valid = valid && record.HasDeadline && record.Deadline.Equal(deadline)
							if i%6 < 3 {
								valid = valid && record.Err == nil
							} else {
								valid = valid && errors.Is(record.Err, context.DeadlineExceeded)
							}
						}
					}
				}
				if valid && len(w.Reports) == 1 && w.Reports[0] == "canceled deadline" {
					return true, i18n.T("الإلغاء والموعد وصلوا للتلات مستويات ورجعوا للقيادة", "Cancellation and deadline reached all three levels and returned to command")
				}
				return false, i18n.T("مستني أثر دخول وخروج المستويات بالـ context الأصلي وتقرير canceled deadline", "Waiting for all entry/exit traces with the original context and a canceled deadline report")
			}
		}})
}

const chainOfCommandStarter = `package main

import (
	"context"
	"orbital/colony"
)

func command(ctx context.Context) error {
	colony.TraceCommand(ctx, "command.enter")
	err := relay(ctx)
	colony.TraceCommand(ctx, "command.exit")
	return err
}
func relay(ctx context.Context) error {
	colony.TraceCommand(ctx, "relay.enter")
	err := probe(context.Background()) // Preserve the caller's context instead.
	colony.TraceCommand(ctx, "relay.exit")
	return err
}
func probe(ctx context.Context) error {
	colony.TraceCommand(ctx, "probe.enter")
	colony.TraceCommand(ctx, "probe.exit")
	return nil
}
func main() {
	canceled, stop := context.WithCancel(context.Background())
	stop()
	_ = command(canceled)
	ctx, cancel := context.WithTimeout(context.Background(), colony.Ticks(colony.Reading(0)))
	defer cancel()
	_ = command(ctx)
	colony.Report("unfinished")
}
`
