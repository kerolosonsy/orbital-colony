package levels

import (
	"fmt"
	"strconv"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

// poolMakespan simulates n workers pulling jobs in order from a shared queue.
func poolMakespan(jobs []sim.Point, n int) int {
	free := make([]int, n)
	at := make([]sim.Point, n)
	for _, j := range jobs {
		w := 0
		for i := range free {
			if free[i] < free[w] {
				w = i
			}
		}
		free[w] += sim.Dist(at[w], j) + 1
		at[w] = j
	}
	end := 0
	for _, f := range free {
		end = max(end, f)
	}
	return end
}

func init() {
	add(&Level{
		En: &LevelText{
			Title: "The Scout",
			Goal:  "Have the scout send locations of large crystals (5 or more) through a channel, then harvest them",
			Lesson: `# Channels: Goroutines Talk to Each Other 📨

> "Don't communicate by sharing memory; share memory by communicating." (Go proverb)

A **channel** is a pipe: one goroutine sends values into one end, and another receives them at the other.

~~~go
found := make(chan int) // A channel carrying int values

go func() {
	found <- 7 // Send 7
}()

x := <-found  // Receive (waits until a value arrives)
~~~

## Sending and receiving wait for each other ⏳
On an ordinary unbuffered channel, the sender **waits** for a receiver, and the receiver waits for a sender. This automatically **synchronizes** them.

## close and range

~~~go
close(found)              // Sender: "I am finished"

for x := range found {    // Receiver: loop until the channel closes
	fmt.Println(x)
}
~~~

- The **sender** closes the channel, not the receiver.
- Sending on a closed channel causes a **panic**.

## Pipeline 🏭
The scout measures while you harvest **at the same time**, following behind it. This is a pipeline: different stages work on different items simultaneously.

## buffered channels
~~~go
ch := make(chan int, 3) // Send 3 values without waiting for a receiver
~~~`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	scout := colony.Spawn()
	found := make(chan int) // Send the x coordinate of each large crystal here

	go func() {
		for x := 0; x < 12; x++ {
			if scout.Measure() >= 5 {
				found <- x
			}
			if x < 11 {
				scout.Move(colony.East)
			}
		}
		close(found) // No more values are coming
	}()

	// Receive from found, move to each x, and harvest
	x := <-found
	fmt.Println("First large crystal at", x)
}
`,
			Hints: []string{
				"~for x := range found { ... }~ receives until the channel closes.",
				"Track your position in ~at := 0~, move East until you reach x, then call ~colony.Harvest()~.",
				"Your drone has no scanner (Measure returns -1), and the scout cannot harvest. They must work together.",
			},
		},
		ID: "11-1-scout", Chapter: 11, Title: "الكشّاف",
		Goal: "الكشّاف يقيس ويبعت مكان الكريستال الكبير (5 فأكتر) في channel، وإنت تحصده",
		Par:  40, Gold: 30, MaxTicks: 300,
		Requires: []string{"chan", "go"},
		Setup: func(w *sim.World) {
			w.Init(12, 2, sim.Regolith)
			w.Realtime, w.MaxDrones, w.ScoutRoles = true, 2, true
			var big, small []sim.Point
			for x := 0; x < 12; x++ {
				t := w.At(x, 0)
				t.Entity, t.Size = sim.Crystal, 1+w.Rand.Intn(9)
				if x%4 == 1 {
					t.Size = 7
				}
				p := sim.Point{X: x}
				if t.Size >= 5 {
					big = append(big, p)
				} else {
					small = append(small, p)
				}
			}
			w.FailCheck = func(w *sim.World) string {
				if remaining(w, small) < len(small) {
					return i18n.T("حصدت كريستال صغير (أقل من 5). الصغير بيتكسر! المهمة فشلت.", "You harvested a small crystal (less than 5). Small crystals shatter! Mission failed.")
				}
				return ""
			}
			w.Check = func(w *sim.World) (bool, string) {
				left := remaining(w, big)
				return left == 0, fmt.Sprintf(i18n.T("كريستال كبير فاضل: %d / %d", "Large crystals remaining: %d / %d"), left, len(big))
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	scout := colony.Spawn()
	found := make(chan int) // هنبعت فيها الـ x بتاع كل كريستال كبير

	go func() {
		for x := 0; x < 12; x++ {
			if scout.Measure() >= 5 {
				found <- x
			}
			if x < 11 {
				scout.Move(colony.East)
			}
		}
		close(found) // خلاص، مفيش حاجة تانية جاية
	}()

	// استقبل من found، وروح احصد كل x
	x := <-found
	fmt.Println("أول كريستال كبير عند", x)
}
`,
		Hints: []string{
			"~for x := range found { ... }~ بتستقبل لحد ما الـ channel تتقفل.",
			"احفظ مكانك في متغير ~at := 0~، واتحرك East لحد ما توصل لـ x، وبعدين ~colony.Harvest()~.",
			"الدرون بتاعك مفيهوش سكانر (Measure بترجع -1)، والكشّاف مبيعرفش يحصد. لازم يشتغلوا مع بعض.",
		},
		Lesson: `# Channels: goroutines بتكلّم بعض 📨

> "Don't communicate by sharing memory; share memory by communicating." (Go proverb)

الـ **channel** ماسورة: goroutine تبعت فيها من ناحية، وواحدة تانية تستقبل من الناحية التانية.

~~~go
found := make(chan int) // channel بتنقل int

go func() {
	found <- 7 // ابعت 7
}()

x := <-found  // استقبل (بيستنى لحد ما حاجة توصل)
~~~

## الإرسال والاستقبال بيستنّوا بعض ⏳
في channel عادية (unbuffered)، اللي بيبعت **بيستنى** لحد ما حد يستقبل، واللي بيستقبل بيستنى لحد ما حد يبعت. ده بيخلّيهم **متزامنين** تلقائياً.

## close و range

~~~go
close(found)              // من ناحية الإرسال: "خلصت"

for x := range found {    // من ناحية الاستقبال: لف لحد ما تتقفل
	fmt.Println(x)
}
~~~

- اللي **بيبعت** هو اللي يقفل، مش اللي بيستقبل.
- الإرسال على channel مقفولة = **panic**.

## Pipeline 🏭
الكشّاف بيقيس، وإنت بتحصد **في نفس الوقت** ورا منه. ده اسمه pipeline: كل مرحلة بتشتغل على حاجة مختلفة في نفس اللحظة.

## buffered channels
~~~go
ch := make(chan int, 3) // تقدر تبعت 3 من غير ما حد يستقبل
~~~`,
	})

	add(&Level{
		En: &LevelText{
			Title: "The Job Queue",
			Goal:  "3 workers take targets from one channel and send back their harvests. Report the total",
			Lesson: `# Worker Pool: A Smart Queue 🏗️

Instead of dividing work manually, put it in a **queue**. Each worker takes the next job as soon as it is free. Faster workers naturally take more jobs.

~~~go
jobs := make(chan colony.Point)
results := make(chan int)

for i := 0; i < 3; i++ {
	go worker(drone, jobs, results) // 3 workers draw from the same channel
}
~~~

## Channel direction in the type
~~~go
func worker(jobs <-chan colony.Point, results chan<- int)
~~~
- ~<-chan T~: receive only.
- ~chan<- T~: send only.

The compiler prevents accidentally using the opposite direction.

## Who closes what? 🔒
1. The job **sender** closes ~jobs~ when done. This lets the workers' ~range jobs~ loops finish.
2. Close ~results~ only after **all** workers finish. Use a WaitGroup:

~~~go
go func() {
	wg.Wait()
	close(results)
}()
for n := range results {
	total += n
}
~~~

## deadlock 💀
If all goroutines are waiting for one another and none can proceed, Go exits with:
~fatal error: all goroutines are asleep - deadlock!~
If you see this, look for a channel with no receiver or one that nobody closed.`,
			Starter: `package main

import (
	"orbital/colony"
	"sync"
)

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
}

// worker takes jobs until jobs closes, sending each harvest to results
func worker(d *colony.Drone, jobs <-chan colony.Point, results chan<- int) {
	for p := range jobs {
		goTo(d, p)
		n, _ := d.TryHarvest()
		results <- n
	}
}

func main() {
	jobs := make(chan colony.Point)
	results := make(chan int)
	var wg sync.WaitGroup

	// 1. Start 3 workers (each with a drone from colony.Spawn)
	// 2. Send all targets into jobs, then close it
	// 3. Sum results and send the total with colony.Report

	_, _, _ = jobs, results, &wg
}
`,
			Hints: []string{
				"For each worker: ~d := colony.Spawn()~, ~wg.Add(1)~, ~go func() { defer wg.Done(); worker(d, jobs, results) }()~.",
				"Sending into jobs blocks main unless someone receives. Send from a goroutine: ~go func() { for _, t := range targets { jobs <- t }; close(jobs) }()~.",
				"Close results once all workers finish: ~go func() { wg.Wait(); close(results) }()~. Then use ~for n := range results { total += n }~.",
			},
		},
		ID: "11-2-worker-pool", Chapter: 11, Title: "طابور الشغل",
		Goal: "3 عمّال بياخدوا الأهداف من channel واحدة، وبيبعتوا اللي جمعوه. ابعت المجموع",
		Par:  70, Gold: 50, MaxTicks: 500,
		Requires: []string{"chan", "go"},
		Setup: func(w *sim.World) {
			w.Init(10, 8, sim.Regolith)
			w.Realtime, w.MaxDrones = true, 4
			w.Targets = randomPoints(w, 12)
			sum := 0
			for _, p := range w.Targets {
				t := w.At(p.X, p.Y)
				t.Entity, t.Size = sim.Crystal, 1+w.Rand.Intn(9)
				sum += t.Size
			}
			want := strconv.Itoa(sum)
			w.Gold = poolMakespan(w.Targets, 3) + 3
			w.Par = w.Gold*3/2 + 5
			w.Check = func(w *sim.World) (bool, string) {
				left := remaining(w, w.Targets)
				rep := reported(w, want)
				msg := fmt.Sprintf(i18n.T("أهداف فاضلة: %d / %d", "Targets remaining: %d / %d"), left, len(w.Targets))
				if left == 0 && !rep {
					msg += i18n.T(" · القاعدة مستنية المجموع", " · Base is waiting for the total")
				}
				return left == 0 && rep, msg
			}
		},
		Starter: `package main

import (
	"orbital/colony"
	"sync"
)

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
}

// worker بياخد شغل من jobs لحد ما تتقفل، ويبعت اللي جمعه في results
func worker(d *colony.Drone, jobs <-chan colony.Point, results chan<- int) {
	for p := range jobs {
		goTo(d, p)
		n, _ := d.TryHarvest()
		results <- n
	}
}

func main() {
	jobs := make(chan colony.Point)
	results := make(chan int)
	var wg sync.WaitGroup

	// 1. شغّل 3 workers (كل واحد بدرون من colony.Spawn)
	// 2. ابعت كل الأهداف في jobs واقفلها
	// 3. اجمع من results وابعت المجموع بـ colony.Report

	_, _, _ = jobs, results, &wg
}
`,
		Hints: []string{
			"لكل worker: ~d := colony.Spawn()~، ~wg.Add(1)~، ~go func() { defer wg.Done(); worker(d, jobs, results) }()~.",
			"الإرسال في jobs هيعلّق main لو محدش بيستقبل، فابعت من goroutine: ~go func() { for _, t := range targets { jobs <- t }; close(jobs) }()~.",
			"اقفل results لما كل العمّال يخلصوا: ~go func() { wg.Wait(); close(results) }()~. وبعدين ~for n := range results { total += n }~.",
		},
		Lesson: `# Worker pool: الطابور الذكي 🏗️

بدل ما تقسّم الشغل بإيدك، خلّي الشغل في **طابور**، وكل عامل ياخد اللي بعده أول ما يفضى. العامل السريع هياخد شغل أكتر لوحده.

~~~go
jobs := make(chan colony.Point)
results := make(chan int)

for i := 0; i < 3; i++ {
	go worker(drone, jobs, results) // 3 عمّال بيسحبوا من نفس الـ channel
}
~~~

## اتجاه الـ channel في النوع
~~~go
func worker(jobs <-chan colony.Point, results chan<- int)
~~~
- ~<-chan T~: استقبال بس.
- ~chan<- T~: إرسال بس.

الـ compiler بيمنعك تعمل العكس بالغلط.

## مين يقفل إيه؟ 🔒
1. اللي **بيبعت** الشغل يقفل ~jobs~ بعد ما يخلص. ده اللي بيخلّي ~range jobs~ عند العمّال تخلص.
2. ~results~ محتاج يتقفل بعد ما **كل** العمّال يخلصوا. استخدم WaitGroup:

~~~go
go func() {
	wg.Wait()
	close(results)
}()
for n := range results {
	total += n
}
~~~

## deadlock 💀
لو كل الـ goroutines مستنية بعض ومحدش بيتحرك، Go بيقفل البرنامج برسالة:
~fatal error: all goroutines are asleep - deadlock!~
لو شفتها، دوّر على channel محدش بيستقبل منها أو محدش قفلها.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "A Storm Is Coming",
			Goal:  "Collect 30 crystal. When warned, return to base before the storm hits",
			Lesson: `# select: Wait for Several Possibilities 🌪️

~select~ is like ~switch~ for channels: it waits until **a channel operation** is ready.

~~~go
select {
case <-warning:
	// The warning arrived
case n := <-results:
	// A result arrived
}
~~~

## default: do not wait

~~~go
select {
case <-warning:
	goHome()
default:
	// Nothing is ready now; keep working
}
~~~

Without ~default~, select **waits**. With it, select checks and moves on immediately. This is a **non-blocking receive**.

## A closed channel signals everyone 📢

~~~go
<-calm // Wait until calm closes
~~~

Receiving from a closed channel returns **immediately** for everyone waiting. That makes ~close~ a useful way to signal several goroutines at once.

## timeouts

~~~go
select {
case r := <-results:
	fmt.Println(r)
case <-time.After(2 * time.Second):
	fmt.Println("Too late")
}
~~~

## Your mission
Work until warned, return to base (0, 0), wait out the storm, then continue until you have 30 crystal. If any drone is outside base when the storm hits, the mission fails.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	warning := colony.StormWarning() // Receives a warning 14 ticks before the storm
	calm := colony.StormClear()      // Closes when the storm ends

	for {
		select {
		case <-warning:
			fmt.Println("Warning! Returning to base")
			// Return to (0, 0) and wait for calm
		default:
			// No warning: keep working
		}
		colony.Harvest()
		colony.Move(colony.East)
	}
	_ = calm
}
`,
			Hints: []string{
				"Back at base, ~<-calm~ waits until the channel closes. Then return to your previous position and continue.",
				"Remember your position so you can return: write ~goTo(p colony.Point)~ and harvest in a snake pattern.",
				"The farthest tile is 12 steps from base; the warning comes 14 ticks before the storm. You have enough time if you return immediately.",
			},
		},
		ID: "11-3-storm", Chapter: 11, Title: "العاصفة جاية",
		Goal: "اجمع 30 كريستال. لما الإنذار ييجي، ارجع للقاعدة قبل ما العاصفة تضرب",
		Par:  140, Gold: 105, MaxTicks: 500,
		Requires: []string{"select"},
		Setup: func(w *sim.World) {
			w.Init(8, 6, sim.Regolith)
			w.Realtime = true
			w.At(0, 0).Entity = sim.Beacon
			w.Each(func(x, y int, t *sim.Tile) {
				if x+y > 0 {
					t.Entity, t.Size = sim.Crystal, 1
				}
			})
			warn := 15 + w.Rand.Intn(20)
			w.Storm = &sim.Storm{
				WarnAt: warn, HitAt: warn + 14, ClearAt: warn + 34,
				Warning: make(chan struct{}, 1), Clear: make(chan struct{}),
			}
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Crystal]
				return n >= 30, fmt.Sprintf(i18n.T("كريستال: %d / 30", "Crystal: %d / 30"), n)
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	warning := colony.StormWarning() // بيوصل فيها إنذار قبل العاصفة بـ 14 tick
	calm := colony.StormClear()      // بتتقفل لما العاصفة تخلص

	for {
		select {
		case <-warning:
			fmt.Println("إنذار! راجعين للقاعدة")
			// ارجع لـ (0, 0) واستنى calm
		default:
			// مفيش إنذار: كمّل شغل
		}
		colony.Harvest()
		colony.Move(colony.East)
	}
	_ = calm
}
`,
		Hints: []string{
			"بعد ما ترجع للقاعدة: ~<-calm~ بتستنى لحد ما الـ channel تتقفل. وبعد كده ارجع مكانك وكمّل.",
			"محتاج تفتكر مكانك عشان ترجعله: اكتب ~goTo(p colony.Point)~ وخلّي الحصاد بتاعك يلف على الخريطة بالثعبان.",
			"أبعد خانة عن القاعدة على بعد 12 خطوة، والإنذار بييجي قبلها بـ 14. يعني عندك وقت لو رجعت فوراً.",
		},
		Lesson: `# select: استنى أكتر من حاجة 🌪️

~select~ زي ~switch~ بس للـ channels: بيستنى **أول** channel تجهز.

~~~go
select {
case <-warning:
	// الإنذار وصل
case n := <-results:
	// نتيجة وصلت
}
~~~

## default: متستناش

~~~go
select {
case <-warning:
	goHome()
default:
	// مفيش حاجة جاهزة دلوقتي، كمّل شغلك
}
~~~

من غير ~default~، الـ select **بتستنى**. ومعاها، بتبص بسرعة وتكمّل. ده اسمه **non-blocking receive**.

## channel مقفولة = إشارة للكل 📢

~~~go
<-calm // بتستنى لحد ما calm تتقفل
~~~

الاستقبال من channel مقفولة بيرجع **فوراً**، لكل اللي مستنيين. عشان كده ~close~ طريقة حلوة تقول لكذا goroutine "يلا" مرة واحدة.

## timeouts

~~~go
select {
case r := <-results:
	fmt.Println(r)
case <-time.After(2 * time.Second):
	fmt.Println("اتأخرت")
}
~~~

## المهمة
اشتغل لحد ما الإنذار ييجي، ارجع للقاعدة (0, 0)، استنى العاصفة تعدّي، وكمّل لحد 30 كريستال. لو أي درون كان برّه القاعدة لما العاصفة تضرب، المهمة تفشل.`,
	})
}
