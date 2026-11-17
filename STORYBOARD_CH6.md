# Chapter 6 lesson storyboards — review draft

Content contract, 2026-11-17. Branch: `lessons-ch6`. This document proposes lesson content only; it changes no mission, starter, check, threshold, or API. Arabic is ready for correctness/leakage review. English is supplied for 6-1 only; translation of 6-2–6-5 waits for Arabic approval.

## Reading and implementation conventions

- Four block kinds only: **markdown**, **walkthrough**, **question**, **sandbox**. Every block below has a stable ID, purpose, and complete Arabic learner text. Step IDs are stable too. Concept IDs have one teaching owner; later references are retrieval/application, not another introduction.
- Main path: Mission → Idea → Example → Watch out → Check → Your turn. A prediction may precede its walkthrough inside Example. Tick explanations are optional notes in Watch out. `Optional` means expandable and initially closed, not missing content. No screen requires all its blocks open together.
- Walkthroughs are authored demonstrations, not execution of the player's code. All line numbers refer to the displayed fragment, starting at 1. Frames describe **state after** the highlighted statement, unless marked “before.” `—` means no variable/view/output yet. Output cells are cumulative. A final frame is terminal; panic has no subsequent execution. Next/Back/Reset navigate frames; Reset returns to the initial frame and clears output, without changing editor code or question results.
- Code fragments run inside `func main()` with the imports declared in their example marker. The verification appendix records actual toolchain outputs. Examples never call unimplemented player functions. `fmt`/`slices` imports are shown by the walkthrough UI; fragments are not advertised as complete files.
- All walkthroughs have a textual view equivalent. Cells show numeric indices left to right, regardless of prose direction. Shared backing stores use labels A/B, not real addresses or promises about physical stack/heap placement. Green = valid/current result, red = error, amber = pending/example-only; labels convey the meaning without color.
- Sandbox operations are a fixed menu, never free-form code. They show an **operation label**, state, and explanation; they do not print fictitious console output. Resets/preset changes clear sandbox-local history. Buttons and inputs are keyboard accessible; no operation changes the active mission or editor. Bounds below are deliberate UI limits, not Go limits.
- Shared control copy AR: «التالي»، «السابق»، «ابدأ من الأول»، «جرّب تاني»، «شوف التنفيذ»، «تفاصيل أكتر»، «ارجع للكود»، «خطوة {current} من {total}». EN: “Next”, “Back”, “Reset”, “Try again”, “See execution”, “More detail”, “Return to code”, “Step {current} of {total}”. Question choice is recorded before feedback; retry clears that question only. An output reveal remains available without answering; Run is never gated.
- Wrong answers are learning evidence, not mission failure. Save stable step/question IDs and a content revision separately from stars, best ticks, chapter `quizScores`, and review metadata. Review navigation starts at Mission; review practice results are temporary. No “mastered” claim from a single answer.
- No calculated answer, chosen crystal position, threshold, complete mission loop, multi-key mission comparator, or general matrix-operation implementation is supplied. Toy routes and data never come from the current randomized world.

### Prior knowledge actually present in chapters 1–5

| Source lesson | Reusable knowledge, without assuming later chapters |
|---|---|
| `1-1-first-signal`, `1-2-first-harvest` | Calls run in order; a drone acts at its current position. |
| `2-1-inventory` | Variables/arithmetic; `fmt.Println` is debugging, `colony.Report` is the mission channel. |
| `2-2-survey` | Multiple results, `_`, types including `float64`, `WorldSize`. |
| `2-3-callsigns`, `2-4-units` | `len` of a string counts bytes; zero-based reading IDs; named integer types. |
| `3-1-long-row`, `3-2-full-field` | Counted/nested loops, assignment, boundaries, movement cost, field traversal. |
| `4-1-ripe-only`, `4-2-terrain`, `4-3-checkpoint` | Conditions, comparison signs, cases, inclusive/exclusive bounds. |
| `5-1-functions`, `5-2-storm`, `5-3-fault-log` | Function inputs/results, `error`, `nil`, descriptive `fmt.Errorf`, error checks. |
| `5-4-undock`, `5-5-command-table` | Panic/recover context, anonymous functions, functions as values, closures. |

Sources read: Arabic lessons in `colony/levels/ch01.go` through `ch05.go`; all chapter-6 starters/checks and heat-grid hidden tests in `ch06.go`; API wrappers in `colony/colony.go`; action clock/scans/reports in `colony/internal/sim/sim.go`. Existing early “every command costs a tick” prose is not repeated: chapter-6 costs distinguish actions from free queries/reports.

### Concept ownership

| Owner | Stable concepts |
|---|---|
| 6-1 | `slice.basics`, `slice.len-cap`, `slice.make-append`, `slice.bounds`, `ticks.actions` |
| 6-2 | `range.slice`, `slice.sharing`, `slice.sort`, `slice.search-max`, `ticks.route-order` |
| 6-3 | `map.counts`, `map.presence`, `map.nil`, `map.order` |
| 6-4 | `sort.comparator`, `sort.identity-ties` |
| 6-5 | `grid.shape`, `grid.row-ownership`, `grid.operations`, `lab.contract` |

`slice.sharing` includes subslices and shallow copying; 6-5 applies that model to row slices. `ticks.actions` is explained once in 6-1; later tick notes apply it rather than redefine it. No struct/pointer syntax, generics syntax, concurrency, or new dependencies are prerequisites.

## 6-1 — سجل الكريستال (`6-1-survey`)

**Teaching objective:** keep ordered observations in a slice; distinguish elements from reserved room; use the returned append value; reason about valid indices. Owns the five 6-1 concepts above. Links back to `2-1-inventory`, `2-3-callsigns`, `3-1-long-row`, `3-2-full-field`.

**Starter audit:** creates nil `sizes`, measures the starting crystal once, appends that measurement, and prints the slice and its length. It does not move or report. Do not ask for another initial measurement or another append of the same observation.

### Step `c61.mission-step` — Mission

#### `c61.mission` · markdown
**Purpose:** state the contract and orient the player to the existing starter.
**Links (author metadata; render natural titles):** `2-1-inventory`.

**AR:** بص، القاعدة عايزة سجل فيه أحجام الـ8 كريستالات، بنفس ترتيبهم من الشمال لليمين. الرقم لوحده مش كفاية؛ ترتيب الأرقام جزء من المعلومة.

الكود الموجود قاس أول كريستالة وحفظها فعلًا، ولسه الدرون في مكان البداية. السطر اللي فيه `fmt.Println` بيعرضلك السجل في الـConsole؛ مش بيبعته للقاعدة. الفرق ده عرفناه في درس جرد المخزن (2-1).

المطلوب سجل كامل في slice وتقرير للقاعدة بـ`colony.Report`. مش محتاج تحصد. إنت اللي هتحدد تكمّل الملاحظة والحركة إزاي.

### Step `c61.idea-step` — Idea

#### `c61.idea` · markdown
**Purpose:** analogy before syntax; introduce `slice.basics` and `slice.len-cap`.
**Links (author metadata; render natural titles):** `2-3-callsigns`.

**AR:** تخيل درج فيه خانات مترقّمة، ومعاه بطاقة بتقول مستخدم كام خانة. الـ**slice** هي وصف لجزء من مخزن عناصر: نوع واحد، وعناصر ليها ترتيب. `[]int` يعني slice عناصرها أعداد صحيحة.

الـ**index** رقم العنصر: أول واحد 0، وبعده 1. `len(s)` عدد العناصر الموجودة، و`cap(s)` المساحة المتاحة من بداية الـslice في الـbacking array، يعني مخزن العناصر اللي وراها.

`append` بتضيف في الآخر وبترجّع وصف slice محدث. احتفظ بالناتج؛ الوصف القديم مش بيطوّل لوحده. مش شرط يتعمل مخزن جديد مع كل إضافة.

**Optional AR:** `len` مع slice بتعد عناصر، مش bytes زي النص في درس أسماء النداء (2-3). و`var s []int` بتبدأ nil: طولها وسعتها صفر، وينفع تعملها append. رسم المخزن هنا نموذج للفهم، مش خريطة عناوين ذاكرة حقيقية.

### Step `c61.example-step` — Example + experiment

#### `c61.basic-trace` · walkthrough
**Purpose:** one growing list, actual output after reveal; introduce no mission loop.

**AR:** خلينا نجرّب تذكرتين لمكتبة المحطة، بعيد عن الكريستال. قبل كل زر «التالي»، توقّع السطر هيغيّر إيه. بعدين بص على العناصر والعدد، مش شكل الأقواس بس.

<!-- go-example id=e61-basic imports=fmt outcome=ok -->
```go
var tickets []int
fmt.Println(len(tickets), cap(tickets))
tickets = append(tickets, 7)
fmt.Println(tickets, len(tickets))
tickets = append(tickets, 4)
fmt.Println(tickets[0], len(tickets))
```

| Frame | Highlight | Variables / slice view | Output | AR caption |
|---|---|---|---|---|
| 0 | before 1 | — | — | «لسه مفيش متغير.» |
| 1 | 1 | tickets=nil, len=0, cap=0; no cells | — | «وصف فاضي، ومفيش عناصر.» |
| 2 | 2 | unchanged | `0 0` | «صفر عناصر وصفر سعة.» |
| 3 | 3 | tickets=[7], index 0, len=1; cap not asserted | `0 0` | «حفظنا وصف الإضافة؛ أول index هو 0.» |
| 4 | 4 | unchanged | `0 0\n[7] 1` | «بقت تذكرة واحدة.» |
| 5 | 5 | tickets=[7 4], indices 0,1, len=2; cap not asserted | previous | «الإضافة في الآخر، مش بدل أول عنصر.» |
| 6 terminal | 6 | unchanged | `0 0\n[7] 1\n7 2` | «العنصر الأول 7، والطول 2.» |

#### `c61.make-playground` · sandbox
**Purpose:** directly manipulate length vs capacity and the make-then-append pitfall.

**AR:** جرب بنفسك: اختار بداية، وبعدها اضغط «ضيف 7». هل فيه عناصر موجودة من الأول؟ تابع `len` جنب `cap`.

**Presets (buttons display exactly these Go expressions):** `make([]int, 3)` → len=3, cap=3, cells `[0 0 0]`; `make([]int, 0, 3)` → len=0, cap=3, three reserved cells outside length; `[]int{}` → len=0, cap=0, non-nil empty slice.

**Fixed menu:** «ضيف 7» (`s = append(s, 7)`); «ابدأ من الأول»; the three preset buttons. At six elements disable Append with «وصلنا لـ6 عناصر في التجربة؛ ابدأ من الأول عشان تجرّب بداية تانية.» Preset changes reset all state.

**What operations show:** Append preserves existing cells, adds one final 7, increments len, and marks the returned descriptor. If new len ≤ cap, reuse store A and preserve cap. Otherwise allocate a distinct store label, copy the existing elements, then append; this illustration chooses new cap=max(new len, max(1, 2×old cap)). Show unused cells dashed, labelled «مساحة محجوزة؛ مش عناصر تقدر تفهرسها دلوقتي». This is a display policy, never an assertion about Go's allocator.

**Always-visible AR label:** «نمو السعة هنا مثال توضيحي؛ Go ممكن تختار سعة مختلفة. مفيش سؤال هيطلب منك تتوقع سعة النمو بالظبط.»

**Preset feedback AR:** length-three: «دول 3 عناصر قيمتهم صفر فعلًا. append هتضيف بعدها، مش هتملا مكانها.» Reserved-three: «عندك مكان لـ3 عناصر، بس لسه عدد العناصر صفر.» Empty literal: «بداية فاضية من غير مساحة محجوزة. أول إضافة محتاجة مخزن يساعها.»

**After Append AR:** «عنصر جديد في الآخر. الطول زاد واحد؛ راجع هل المخزن اتغير ولا اتستخدم تاني.» Reset: «رجعنا للبداية اللي اخترتها، ومسحنا إضافات التجربة بس.» No exact-capacity question uses this illustration. Go fixtures `e61-presets` and `e61-many-appends` below verify element/length semantics; allocator choice is not an expected output.

### Step `c61.watch-step` — Watch out

#### `c61.mistakes` · markdown
**Purpose:** explicit misconception/fact/reason/fix-thinking pattern.

**AR:** **الافتراض:** «make بطول 3 بتحجز أماكن فاضية، فهعمل append عشان أملاها.» → **Go بتعمل إيه؟** بتنشئ 3 عناصر صفر، والإضافة تيجي بعدهم. → **ليه؟** الطول بيعد عناصر موجودة، والسعة بتعد المكان المتاح. → **تفكر إزاي؟** قرر: عندك عناصر هتكتب فيها بالـindex، ولا سجل بيبدأ بطول صفر ويكبر؟

**الافتراض:** «append بتزوّد طول المتغير من غير ما أحتفظ بالناتج.» → **Go بتعمل إيه؟** بترجّع وصف محدث؛ ونداء append لوحده من غير استخدام نتيجته مش بيتقبل كسطر مستقل. → **ليه؟** فيه ناتج لازم تستخدمه. → **تفكر إزاي؟** احتفظ بالوصف اللي راجع، في نفس المتغير أو متغير تاني مناسب.

#### `c61.bounds-trace` · walkthrough
**Purpose:** show terminal out-of-range panic, including capacity not making an index valid.

**AR:** خد بالك: رقم العنصر لازم يكون من صفر لحد قبل `len`. وجود مساحة زيادة مش معناه وجود عنصر زيادة. بص ليه الطباعة التانية هنا مش هتتعمل.

<!-- go-example id=e61-bounds imports=fmt outcome=panic -->
```go
s := make([]int, 0, 3)
fmt.Println(len(s), cap(s))
i := 0
fmt.Println(s[i])
fmt.Println("after")
```

| Frame | Highlight | Variables / slice view | Output | AR caption |
|---|---|---|---|---|
| 0 | before 1 | — | — | «قبل البداية.» |
| 1 | 1 | s=[], len=0, cap=3; three reserved cells | — | «مفيش عنصر موجود رغم السعة.» |
| 2 | 2 | unchanged | `0 3` | «الطول صفر.» |
| 3 | 3 | i=0 | `0 3` | «الـindex المطلوب صفر.» |
| 4 terminal panic | 4 | s unchanged; index 0 marked invalid | stdout `0 3`; panic `runtime error: index out of range [0] with length 0` | «البرنامج وقف هنا؛ after مش هتتطبع. الإصلاح يبدأ بفهم الطول، مش إخفاء الـpanic.» |

#### `c61.ticks` · markdown · Optional
**Purpose:** own `ticks.actions`, correct earlier overgeneralization, show a toy estimate.

**AR:** **تقدير لتجربة منفصلة:** ممر تدريب فيه 4 علامات، الدرون عند أول واحدة، وهنقرأ كل علامة مرة ونتحرك بس بين العلامات. قراءة `Measure` واحدة = tick، و`Move` واحدة = tick؛ يبقى 4 قراءات + 3 حركات = **7 ticks**. الإضافة والحساب والطباعة العادية مش بيزوّدوا عدّاد أوامر اللعبة، و`Report` و`WorldSize` من غير tick في المهمات دي.

لو حاولت تتحرك برّه الممر هتدفع tick حتى لو الحركة فشلت. التقدير بيفترض مفيش أوامر زيادة، وإن البرنامج بيقف بعد التقرير الصحيح. الرقم المقاس من التشغيل الأخير بيتعرض بعنوان «مقاس: آخر تشغيل»، مش على إنه نفس خريطة المعاينة أو أفضل رقم محفوظ.

**Author tick audit:** live 6-1 has 8 measurements + 7 moves =15; Par=Gold=15, MaxTicks=60. This verifies the lesson assumptions; the UI should show the existing mission target without adding a solved route. Source: `ch06.go:64–85`, `sim.go:503–544,771–822`.

### Step `c61.check-step` — Check

#### `c61.q-make` · question
**Purpose:** assess `slice.make-append` with fresh toy data.

**AR prompt:** إيه اللي هيتطبع؟ توقّع قبل ما تشوف النتيجة.

<!-- go-example id=e61-make-question imports=fmt outcome=ok -->
```go
cards := make([]int, 2)
cards = append(cards, 9)
fmt.Println(cards, len(cards))
```

Options: A «`[9] 1`»; B «`[0 0 9] 3`»; C «`[9 0] 2`». **Correct: B.**

Feedback A: «make هنا بدأت بعنصرين فعلًا؛ الإضافة مش بتمسحهم. ارجع لبداية الطول في التجربة.»

Feedback B: «كده بصّيت للطول صح: صفرين موجودين، وبعدهم عنصر جديد قيمته 9، فبقوا 3.»

Feedback C: «ده كان محتاج كتابة في index موجود. append بتضيف في الآخر، مش بتختار أول صفر.»

#### `c61.q-append-result` · question
**Purpose:** spot an intentional compile error without giving mission code.

**AR prompt:** طالب نسي يحتفظ بنتيجة append. أي وصف صح للكود؟

<!-- go-example id=e61-unused-append imports=fmt outcome=compile-error -->
```go
s := []int{6}
append(s, 2)
fmt.Println(s)
```

Options: A «هيطبع `[6 2]`»; B «هيطبع `[6]` عادي»; C «مش هيتبني؛ نتيجة append مش مستخدمة». **Correct: C.**

Feedback A: «الإضافة ليها ناتج لازم تستخدمه؛ النداء ده مش بيغيّر طول s تلقائيًا.»

Feedback B: «Go هنا مش بتكتفي بتجاهل النتيجة؛ السطر نفسه مرفوض وقت البناء.»

Feedback C: «بالظبط: ده compile error، قبل أي طباعة. فكّر هتحتفظ بوصف الإضافة فين.»

#### `c61.q-index` · question
**Purpose:** assess `slice.bounds`, without repeating the already revealed panic example.

**AR prompt:** slice فيها 4 عناصر. أنهي index بيختار آخر عنصر؟

Options: A «4»; B «3»; C «-1». **Correct: B.**

Feedback A: «4 هو عدد العناصر، مش رقم آخر عنصر. رقم 4 هيبقى برّه الحدود.»

Feedback B: «أيوه: الأرقام 0 و1 و2 و3. افتكر فرق العدد عن الرقم.»

Feedback C: «Go مش بتستخدم -1 كاختصار لآخر عنصر؛ ده index غير صالح.»

### Step `c61.turn-step` — Your turn

#### `c61.turn` · markdown
**Purpose:** one concrete next action without constructing the traversal.

**AR:** ارجع للكود وشغّل الـstarter زي ما هو. بص على طول `sizes` في الـConsole: أول قراءة اتحفظت فعلًا. حدد في تعليق عندك إيه أول ملاحظة لسه ناقصة، قبل ما تزود أوامر. زر «ارجع للكود» هيفتحلك المحرر.

### Three hints (position labels, not new hint schema)

1. **Concept / الفكرة — `c61.h1`:** «السجل بيحفظ ملاحظات مرتبة؛ مش مجموع الأحجام. افتكر الفرق بين طباعة السجل وإرساله للقاعدة.»
2. **Structure / التنظيم — `c61.h2`:** «قبل ما تكرر أي خطوة، راجع إيه اللي اتسجل بالفعل في الـstarter. خلي لكل ملاحظة مكان واحد في السجل.»
3. **Specific / ركّز هنا — `c61.h3`:** «راجع السطر اللي بيستخدم append: هل محتفظ بالناتج؟ وراجع قيمة التقرير: هل هي السجل نفسه بدل طوله أو آخر قراءة؟»

**Leakage review:** player still constructs traversal, iteration bounds, coordination of observation/storage/movement, and reporting placement. No complete loop, move condition, or mission data is supplied. Starter inspection is a concrete first action, not additional scan code. The toy cost equation is separate from live dimensions.

**Fallback corrections:** replace “append makes a new slice, so it must go in the same variable” with “updated slice value; may reuse storage; use its result.” Add length-versus-capacity and the make-then-append warning. State the full index bound. Keep nil slice append valid. Replace the old solved-loop hint; preserve the starter and correct `Println`/`Report` distinction. Move subslice sharing depth to its owner `slice.sharing` in 6-2; 6-1 fallback can link there instead of introducing an unexplained `nums[1:3]`.

## 6-2 — أغلى 3 (`6-2-top-three`)

**Teaching objective:** inspect stored observations without losing their identity; distinguish a copied value from shared elements; distinguish selecting valuable items from planning travel. Owns `range.slice`, `slice.sharing`, `slice.sort`, `slice.search-max`, `ticks.route-order`. Prerequisites link to 6-1 concepts, counted loops from `3-1-long-row`, `_` from `2-2-survey`, conditions from `4-1-ripe-only`.

**Starter audit:** allocates ten existing elements, measures each crystal once into its matching index, moves East nine times, and prints the sizes. It finishes at x=9, with no harvests. The mission permits three successful harvests and generates distinct sizes from a permutation of 1–20. `range` is required; no particular selection algorithm or `slices` call is required. Do not redo the supplied survey.

### Step `c62.mission-step` — Mission

#### `c62.mission` · markdown
**Purpose:** specify the challenge without a threshold or collection loop.

**AR:** معاك فرصة تشيل 3 كريستالات بس، والمطلوب الأكبر في الصف. الكود قاس الصف كله بالفعل، وحفظ كل حجم عند نفس index بتاع موقعه، والدرون واقف عند x=9. متبدأش المسح من الأول.

فكّر في سؤالين منفصلين: إيه اللي يستاهل يتاخد؟ وإزاي تزور اللي اخترته؟ استخدام `range` جزء من تدريب المهمة. إنت اللي هتبني قرار الاختيار وطريقة الوصول.

### Step `c62.idea-step` — Idea

#### `c62.idea` · markdown
**Purpose:** introduce range and ownership analogy, referencing rather than reteaching 6-1.
**Links (author metadata; render natural titles):** `2-2-survey`.

**AR:** تخيل دفتر تذاكر: على كل سطر رقم مكان وقيمة مكتوبة. `range` بتعدّي على السطور وتديك الـindex والقيمة: الشكل هو `for i, value := range list`. لو مش محتاج الـindex، استخدم `_` اللي عرفناه في درس مسح المنطقة (2-2).

القيمة اللي في متغير اللفة **نسخة من العنصر**، مش اسم تاني لنفس خانته. أما إسناد slice لمتغير تاني فبينسخ وصف الـslice، والعناصر نفسها بتفضل مشتركة. دول مستويين مختلفين من النسخ.

**Optional AR:** في `range` على slice، Go بتحدد مجال الـindices من قيمة الـslice عند بداية الحلقة؛ الإضافات مش بتمد اللفة. ده مش تصوير لكل العناصر مقدمًا: تغيير عنصر لسه دوره مجاش في نفس المخزن ممكن يظهر وقت قراءته. ملحوظة دي عن slice، مش قاعدة مختصرة لكل أشكال range.

### Step `c62.example-step` — Example

#### `c62.q-range-copy` · question
**Purpose:** predict before the range trace; assess `range.slice`.

**AR prompt:** ده دفتر أرقام تدريب. بعد الحلقة، أنهي slice هتتطبع في السطر الأخير؟

<!-- go-example id=e62-range-copy imports=fmt outcome=ok -->
```go
pages := []int{3, 8}
for i, value := range pages {
    value *= 2
    fmt.Println(i, value)
}
fmt.Println(pages)
```

Options: A «`[6 16]`»; B «`[3 8]`»; C «`[0 1]`». **Correct: B.**

Feedback A: «الضرب غيّر value اللي اتنسخت، مش عنصر pages. شوف خانة العنصر وهي ثابتة في التنفيذ.»

Feedback B: «صح: اللي اتغير نسخة القيمة؛ عناصر الدفتر فضلت زي ما هي.»

Feedback C: «0 و1 أرقام الأماكن اللي range بترجعها في i، مش القيم الموجودة في الدفتر.»

#### `c62.range-trace` · walkthrough
**Purpose:** trace `e62-range-copy`; optional second trace establishes fixed bounds without a value snapshot.

**AR:** خلينا نفصل اللي جوه الدفتر عن الورقة اللي في إيد اللفة. «التالي» يبيّن نسخة value وهي بتتغيّر، والدفتر نفسه ثابت.

**Code:** `e62-range-copy` above; no alternate snippet.

| Frame | Highlight | Variables / slice view | Cumulative output | AR caption |
|---|---|---|---|---|
| 0 | before 1 | — | — | «قبل الدفتر.» |
| 1 | 1 | pages=[3 8], indices 0,1 | — | «قيمتين في المخزن A.» |
| 2 | 2 first iteration | i=0, value=3; pages unchanged | — | «استلمنا نسخة أول قيمة.» |
| 3 | 3 | i=0, value=6; pages=[3 8] | — | «نسخة value اتضاعفت بس.» |
| 4 | 4 | unchanged | `0 6` | «ده رقم المكان والنسخة المعدلة.» |
| 5 | 2 second iteration | i=1, value=8 | `0 6` | «نسخة تانية من العنصر التاني.» |
| 6 | 3 | i=1, value=16; pages unchanged | `0 6` | «الدفتر لسه زي ما هو.» |
| 7 | 4 | unchanged | `0 6\n1 16` | «انتهت آخر لفة.» |
| 8 terminal | 6 | i/value out of scope; pages=[3 8] | `0 6\n1 16\n[3 8]` | «الطباعة دي للـslice الأصلية.» |

**Optional AR:** «إيه اللي هيحصل لو زوّدنا عنصر وإحنا بنلف؟ هنا محجوز مكان كفاية من الأول، فالإضافة مش هتغيّر المخزن. تابع حدود اللفة والقيمة اللي لسه دورها مجاش.»

<!-- go-example id=e62-range-bounds imports=fmt outcome=ok -->
```go
s := make([]int, 2, 3)
s[0], s[1] = 2, 5
for i, value := range s {
    if i == 0 {
        s = append(s, 9)
        s[1] = 6
    }
    fmt.Println(i, value)
}
fmt.Println(s)
```

| Frame | Highlight | Variables / slice view | Cumulative output | AR caption |
|---|---|---|---|---|
| 0 | before 1 | — | — | «تجربة مستقلة.» |
| 1 | 1 | s=[0 0], len=2 cap=3; A[2] reserved | — | «عنصرين ومكان زيادة.» |
| 2 | 2 | A=[2 5 reserved] | — | «كتبنا قيمتين.» |
| 3 | 3 | i=0,value=2; initial range length=2 | — | «اللفة هتزور 0 و1.» |
| 4 | 4 true | unchanged | — | «الشرط صحيح لأول لفة.» |
| 5 | 5 | s=[2 5 9], len=3 cap=3; range still length=2 | — | «الإضافة جوه نفس المخزن.» |
| 6 | 6 | s=[2 6 9] | — | «غيرنا قيمة لسه دورها مجاش.» |
| 7 | 8 | i=0,value=2 | `0 2` | «النسخة الحالية لسه 2.» |
| 8 | 3 next | i=1,value=6 | `0 2` | «اتقرت القيمة الحالية من نفس المخزن.» |
| 9 | 4 false | unchanged; branch skipped | `0 2` | «مفيش تعديل في اللفة دي.» |
| 10 | 8 | unchanged | `0 2\n1 6` | «مفيش لفة للـindex 2.» |
| 11 terminal | 10 | i/value out of scope; s=[2 6 9] | `0 2\n1 6\n[2 6 9]` | «حدود ثابتة مش معناها نسخة مجمدة من القيم.» |

#### `c62.sharing-playground` · sandbox
**Purpose:** own `slice.sharing`; show aliasing, descriptor length, Clone, and reallocation with actual user operations.
**Links (author metadata; render natural titles):** `slice.len-cap`.

**AR:** تخيل بطاقتين بيوصفوا نفس درج الأرقام. اختار «b = a»، وبعدها غيّر أول عنصر من b. جرّب «Clone إلى b» وقارن. راجع نموذج الطول والسعة في «الطول والسعة» لو محتاج.

**Presets:** «مكان زيادة» → a=[5 1], len=2 cap=3, store A; «مخزن مليان» → a=[5 1], len=2 cap=2, store A. In both, b is not created yet. Fixed menu: «b = a»; «Clone إلى b» (`b = slices.Clone(a)`); «خلّي b[0] = 9»; «ضيف 7 إلى a»; preset buttons; Reset. No sort operation in this sandbox.

**Operation effects:** `b=a` creates a separate descriptor with a's current length/capacity and the same store/offset. Clone creates independent integer elements in a new store, preserving order; its displayed cap is marked “not specified” rather than asserting len=cap. `b[0]=9` writes to b's store; any alias covering that cell visibly changes. Disable this button before b exists with «اعمل b الأول». Append changes a's descriptor only; with room, a/b may still share elements while b keeps its old length. Without room, a gets a new store, copying elements; b keeps its previous store. Limit a to four elements; show «حد التجربة 4 عناصر؛ ابدأ من الأول عشان تجرّب تاني». Growth is illustrative under the same policy/label as `c61.make-playground`.

**AR status text:** alias: «وصفين منفصلين، نفس مخزن العناصر.» Clone: «الأرقام اتنسخت في مخزن مستقل؛ ده كفاية هنا لأن العناصر int.» Write: «غيّرنا عنصر في المخزن اللي b بتشوفه؛ تابع أي بطاقة تانية بتشوف نفس العنصر.» Append reuse: «a طولها زاد؛ b طولها زي ما هو، والمخزن لسه مشترك.» Append allocation: «a انتقلت لمخزن جديد؛ b لسه بتشوف المخزن القديم.» Reset: «بطاقة a رجعت للبداية، وb اتشالت من التجربة.»

**Optional subslice AR:** «الجزء a[0:1] بيشوف أول عنصر من نفس المخزن؛ الحد الأخير 1 مش داخل الجزء. متفترضش إن اختيار جزء عمل نسخة.» The tiny separate Go fixture `e62-subslice` verifies this; no additional sandbox operations or arbitrary slices.

#### `c62.sort-trace` · walkthrough
**Purpose:** own `slice.sort`; show in-place sorting, without a top-N selection implementation.

**AR:** ترتيب أرقام التذاكر مش بيحرّك أصحابها في الطابور. `slices.Sort` بترتّب العناصر تصاعدي **جوه نفس الـslice**. لازم تضيف `"slices"` للـimports. جرّب تتوقع ترتيب الأرقام دي، وبعدها شوف النتيجة.

<!-- go-example id=e62-sort imports=fmt,slices outcome=ok -->
```go
codes := []int{14, 6, 11}
other := codes
slices.Sort(other)
fmt.Println(codes, other)
```

| Frame | Highlight | Variables / slice view | Output | AR caption |
|---|---|---|---|---|
| 0 | before 1 | — | — | «قبل التذاكر.» |
| 1 | 1 | codes→A=[14 6 11] | — | «ده ترتيب البداية.» |
| 2 | 2 | codes/other→A | — | «الإسناد مش نسخ عناصر.» |
| 3 | 3 after library call | A=[6 11 14], both descriptors unchanged | — | «النداء كله خطوة واحدة؛ مش بنحاكي خوارزمية Go الداخلية.» |
| 4 terminal | 4 | unchanged | `[6 11 14] [6 11 14]` | «الاتنين بيشوفوا الترتيب الجديد.» |

#### `c62.library-notes` · markdown
**Purpose:** introduce `slice.search-max` and connect value identity to the mission.

**AR:** لو الدفتر رابط الـindex بموقع حقيقي، ترتيب قيمه في نفس المخزن بيغيّر الربط ده. العالم نفسه متحركش. نسخة مستقلة ممكن تحافظ على دفتر المواقع الأصلي؛ اختيار التصميم مسؤوليتك.

`slices.Max` بترجع أكبر قيمة، ومحتاجة slice غير فاضية. `slices.Index` بترجع مكان **أول** قيمة مطابقة، أو -1 لو مش موجودة. متستخدمش -1 كـindex. قراءة أكبر قيمة مش هي نفسها تحديد كل المواقع اللي عايزها.

**Optional AR:** الأحجام في مهمة أغلى 3 مختلفة عن بعض. بس في برامج تانية ممكن تتكرر؛ Index مش بتدور على «المطابقة اللي بعدها» لو ناديتها بنفس المدخلات. كمان استبدال قيمة برقم خاص عشان تعتبرها مستخدمة محتاج ضمان إن الرقم ده مستحيل يكون قيمة صحيحة. مفيش sentinel عام صالح لكل البيانات.

#### `c62.max-index-trace` · walkthrough
**Purpose:** observe first match and absence on non-mission data.

**AR:** ده سجل أرقام أبواب؛ الباب 4 مكتوب مرتين. توقّع Max وIndex، وشوف ليه تكرار السؤال عن نفس الرقم مش بيختار باب تاني.

<!-- go-example id=e62-max-index imports=fmt,slices outcome=ok -->
```go
doors := []int{4, 9, 4}
fmt.Println(slices.Max(doors))
fmt.Println(slices.Index(doors, 4))
fmt.Println(slices.Index(doors, 8))
```

| Frame | Highlight | Variables / slice view | Cumulative output | AR caption |
|---|---|---|---|---|
| 0 | before 1 | — | — | «قبل السجل.» |
| 1 | 1 | doors=[4 9 4], indices 0,1,2 | — | «قيمتين 4 في مكانين.» |
| 2 | 2 | unchanged; highlight value 9 | `9` | «دي أكبر قيمة، مش رقم موقعها.» |
| 3 | 3 | unchanged; highlight index 0 | `9\n0` | «أول مطابقة بس.» |
| 4 terminal | 4 | unchanged; no matching cell | `9\n0\n-1` | «8 مش موجودة؛ -1 إشارة غياب.» |

#### `c62.empty-max-trace` · walkthrough · Optional
**Purpose:** complete the Max precondition with a terminal panic, not an invented zero result.

**AR:** خد بالك: سجل فاضي مفيهوش أكبر قيمة. ده خطأ تشغيل، مش إجابة صفر.

<!-- go-example id=e62-empty-max imports=fmt,slices outcome=panic -->
```go
doors := []int{}
fmt.Println(slices.Max(doors))
fmt.Println("after")
```

Frames: 0 before line 1, all absent, output none, caption «قبل البداية». 1 after line 1, doors=[] len=0, view empty, output none, caption «مفيش عناصر نقارنها». 2 terminal at line 2, variables unchanged, no stdout, panic `slices.Max: empty list`, caption «Max وقفت التنفيذ؛ الطباعة نفسها ماكملتش، وafter مش هتتطبع». No line-3 frame.

### Step `c62.watch-step` — Watch out

#### `c62.mistakes` · markdown
**Purpose:** consolidate misconceptions without a mission recipe.

**AR:** **الافتراض:** «b = a عملت قايمة جديدة بعناصر مستقلة.» → **الواقع:** الوصف اتنسخ والعناصر لسه مشتركة. → **ليه؟** slice وصف للمخزن. → **فكّر في الإصلاح:** حدد مين مسموحله يغيّر نفس العناصر، ومين محتاج نسخة مستقلة.

**الافتراض:** «بعد Sort، الـindex لسه هو موقع الحاجة قبل الترتيب.» → **الواقع:** القيم غيرت أماكنها في السجل، مش في العالم. → **ليه؟** الترتيب عدّل العناصر. → **فكّر في الإصلاح:** حافظ على هوية العنصر بدل ما تستنتجها من ترتيب جديد.

**الافتراض:** «range صورت كل القيم أول الحلقة.» → **الواقع:** حدود اللفة ثابتة هنا، لكن عناصر المخزن المشترك ممكن تتغير قبل دور قراءتها. → **ليه؟** الوصف مش نسخة من كل عنصر. → **فكّر في الإصلاح:** اسأل أي مخزن بتقرأه الحلقة، مش بس طول المتغير دلوقتي.

#### `c62.ticks` · markdown · Optional
**Purpose:** apply `ticks.actions`, own route-order reasoning using two preselected toy destinations.

**AR:** تجربة توصيل مستقلة: ممر من 5 خانات x=0..4. افترض المسح اتعمل مرة من أول خانة لآخر واحدة: 5 قراءات و4 حركات = **9 ticks**. الحساب والبحث والترتيب في الذاكرة مش قراءات درون جديدة.

افترض حد اختارلك وجهتين بالفعل، x=1 وx=3، وإنت عند x=4، ومفيش رجوع للقاعدة مطلوب. ترتيب الزيارة 4→3→1 فيه 3 حركات؛ ترتيب 4→1→3 فيه 5. لو كل وجهة محتاجة Harvest واحدة، زوّد 2 في الحالتين. قيمة الوجهة ساعدتك تختارها؛ مكانها بيساعدك تخطط السفر. دي تجربة بوجهتين، مش اختيار أكبر 3 في خريطتك.

**Author tick audit:** real starter scan=19; the implemented Gold is `19+(9-minX)+3` for a westward visit ending at the leftmost selected crystal; Par=Gold+14. Static 31/45 defaults are not every generated run's targets. Do not reveal live minX, selected positions, or this formula as a solved mission route. Show the existing result's dynamic target and measured ticks, labelled with the run; preview can have a different seed. Sources: `ch06.go:214–242`, `runner.go:243–274`.

### Step `c62.check-step` — Check

#### `c62.q-first-match` · question
**Purpose:** assess `slice.search-max` with repeated identity and no selection implementation.

**AR prompt:** في دفتر `[2 6 2]`، ناديت Index للقيمة 2 مرتين، من غير أي تعديل للدفتر. هيرجع إيه بالترتيب؟

Options: A «0 وبعدها 2»; B «0 وبعدها 0»; C «2 وبعدها -1». **Correct: B.**

Feedback A: «Index مبتفتكرش النداء اللي قبله. نفس المدخلات بتبدأ البحث من الأول تاني.»

Feedback B: «صح؛ أول مطابقة عند 0 كل مرة. لو الهوية مهمة، لازم تصميمك يحتفظ بيها.»

Feedback C: «دي مش عملية استهلاك أو مسح عنصر. أول 2 موجودة عند 0 في النداءين.»

#### `c62.q-route` · question
**Purpose:** assess `ticks.route-order`, separate movement from arithmetic.

**AR prompt:** في تجربة التوصيل وإنت عند x=4، الوجهتين معروفين: 3 و1. كام tick للحركة والحصاد في الزيارة 4→3→1، من غير حساب المسح اللي خلص؟

Options: A «3»; B «5»; C «7». **Correct: B.**

Feedback A: «حسبت الحركة بس؛ عندنا كمان عمليتي Harvest، واحدة لكل وجهة.»

Feedback B: «أيوه: حركة واحدة وبعدين حركتين، ومعاهم حصادين؛ الإجمالي 5.»

Feedback C: «7 تناسب السفر الأطول 4→1→3 مع حصادين. ارسم المسافات بدل ما تعد أسماء الوجهات.»

### Step `c62.turn-step` — Your turn

#### `c62.turn` · markdown
**Purpose:** concrete inspection of identity without a selection algorithm.

**AR:** ارجع لنتيجة مسح الـstarter، واختار index واحد في الـConsole. اكتب تعليق عندك يوضح القيمة دي مرتبطة بأنهي x قبل أي ترتيب. بعد كده ابني قرارك؛ الدرس مش هيختار الكريستالات بدل منك.

### Three hints

1. **Concept — `c62.h1`:** «افصل قيمة الحاجة عن مكانها. اختيار الحاجة وخطة الوصول ليها سؤالين مختلفين.»
2. **Structure — `c62.h2`:** «قبل ما تغيّر ترتيب أي سجل، حدد أي معلومات لازم تفضل ثابتة عشان تعرف ترجع للمواقع.»
3. **Specific — `c62.h3`:** «راجع السطر اللي بيغيّر ترتيب البيانات: تقدر تحتفظ بنسخة مستقلة بـ slices.Clone قبل الترتيب، أو تحتفظ بكل قيمة مرتبطة بموقعها في تصميمك. اختار اللي يحافظ على هوية العنصر؛ ترتيب القيم لوحده مش بيحرّك مواقع العالم.»

**Leakage review:** player constructs selection of exactly three distinct crystals, preservation/recovery of positions, actual range use, travel and harvest control, and stopping. There is no threshold index, repeated-max implementation, selected-position list from the real world, or harvesting predicate. Route teaching starts from two already-selected toy destinations; it teaches a cost comparison rather than supplying top-three selection plus travel code.

**Fallback corrections:** add range value-copy and slice-bound semantics; distinguish fixed bounds from frozen values. Qualify Clone as shallow; explain first-match/-1 and empty Max. State Sort changes elements, not world positions. State actual 6-2 sizes are distinct; defer forced-duplicate tie handling to 6-4. Keep dynamic-star wording, clarify preview versus last run, and remove the threshold/copy/sort/west-harvest recipe from all three old hints.

## 6-3 — تعداد الحياة (`6-3-census`)

**Teaching objective:** represent category counts with a map; distinguish absence, stored zero, entry count, and count total; avoid dependence on map iteration order. Owns `map.counts`, `map.presence`, `map.nil`, `map.order`. Links to `2-4-units` for named integer types, `3-2-full-field` for traversal, `slice.basics` for contrast, `range.slice` for the already-known range form.

**Starter audit:** creates a writable `map[colony.Entity]int`, reads the starting tile once, increments its entity count, and prints the map. No movement or report. The world is 6×4. Checks accept entity-keyed or name-keyed maps, with or without empty tiles; moss, lichen, and crystal are guaranteed present. Entity **counts**, not summed crystal sizes, are wanted. Do not restart counting the first tile without accounting for that existing increment.

### Step `c63.mission-step` — Mission

#### `c63.mission` · markdown
**Purpose:** state exactly what is counted and what the starter already recorded.
**Links (author metadata; render natural titles):** `2-4-units`.

**AR:** القاعدة عايزة تعرف **كام خانة** فيها كل نوع: moss وlichen وcrystal. حجم الكريستالة مش عدد الكريستالات؛ هنا بنعدّ وجود النوع.

الخريطة 6×4، والكود سجل نوع خانة البداية مرة فعلًا. `colony.GetEntity` بترجع النوع اللي تحت الدرون دلوقتي. `counts` مفاتيحها من نوع `colony.Entity`، زي فكرة النوع اللي له اسم في درس وحدات القياس (2-4).

ابعت map بالأعداد بـ`colony.Report`؛ الطباعة لوحدها مش تقرير. ينفع تعدّ empty أو تسيبها، والاتنين مقبولين. مش محتاج تحصد حاجة. تخطيط زيارة الخانات من غير تكرار مسؤوليتك.

### Step `c63.idea-step` — Idea

#### `c63.idea` · markdown
**Purpose:** own map construction and key/value model through a cabinet analogy.

**AR:** تخيل دولاب، وعلى كل درج اسم مورد. بدل ما تسأل «العنصر رقم كام؟» زي slice، بتسأل «الدرج اللي مفتاحه water فيه كام؟». الـ**map** بتربط **key**، يعني مفتاح، بـ**value**، يعني قيمة.

`map[string]int{}` بتعمل map جاهزة للكتابة: أسماء مفاتيح وأعداد قيم. `make(map[string]int)` طريقة تانية. القراية لمفتاح مش موجود بترجع صفر هنا لأن نوع القيمة int. الزيادة بـ`++` بتقرا القيمة وتكتبها بعد زيادتها؛ لو المفتاح غايب، بتبدأ من صفر وتضيفه.

**Optional AR:** نوع المفتاح لازم ينفع تقارنه بـ`==`؛ string وint مناسبين، لكن slice وmap مش مفاتيح صالحة. `len` هنا عدد المفاتيح الموجودة، مش مجموع القيم. تشابه اسم len مش معناه إن كل أنواع البيانات بتعد نفس الحاجة.

### Step `c63.example-step` — Example + experiment

#### `c63.counts-playground` · sandbox
**Purpose:** resource buttons that genuinely mutate a bounded map and distinguish len/total/presence.

**AR:** جرب بنفسك: زوّد water كذا مرة. هل عدد المفاتيح زاد كل مرة؟ جرّب تخلي air موجودة بقيمة صفر، وبعدها اقرا مفتاح مش موجود.

**Presets:** «دولاب فاضي» → writable empty map, len=0,total=0; «موجود بصفر» → water=2,air=0, len=2,total=2. Display rows in alphabetical UI order labelled «ترتيب عرض للتوضيح؛ مش ترتيب range».

**Fixed menu:** «زوّد water»; «زوّد air»; «امسح air»; «اقرا air»; «اقرا spare»; preset buttons; Reset. Each increment changes only the chosen value by one, inserting its key if absent. Each resource is capped at 9 only for the demonstration; disable its increment at 9 with «حد التجربة 9؛ ابدأ من الأول عشان تجرّب تاني». No nil preset here; nil is a separate terminal walkthrough.

**Visible state:** key/value table; «عدد المفاتيح len»; «إجمالي الوحدات»; read-result card `(value, ok)`; latest operation. Recompute total from all existing values. Deleting air removes its entry and units; deleting an absent key is harmless. Reads never mutate the map, len, or total. Missing spare always returns (0,false). A stored air=0 reads (0,true). Reset clears result cards and restores the chosen preset.

**AR status:** first increment: «المفتاح اتضاف، وقيمته بقت 1.» Existing increment: «نفس المفتاح، قيمة أعلى؛ len زي ما هي.» Delete existing: «شيلنا المفتاح نفسه؛ len قلت واحد.» Delete absent: «المفتاح مش موجود؛ مفيش تغيير.» Read stored zero: «القيمة صفر، لكن ok=true: المفتاح موجود.» Read absent: «القيمة الصفرية وok=false: المفتاح غايب، والقراية مضافتوش.» Read positive: «دي القيمة الموجودة؛ القراية مش زيادة.» Reset: «رجعنا لحالة الدولاب اللي اخترتها.» Fixture: `e63-sandbox`.

#### `c63.q-presence` · question
**Purpose:** predict before the comma-ok explanation; assess absence versus zero.

**AR prompt:** قريت `counts["lamp"]` ورجعت 0. إيه اللي تقدر تجزم بيه من الرقم لوحده؟

Options: A «lamp أكيد مش موجودة»; B «lamp أكيد موجودة وصفر»; C «محتاج ok عشان أفرق بين الحالتين». **Correct: C.**

Feedback A: «المفتاح الغايب بيرجع صفر، بس المفتاح الموجود ينفع تكون قيمته صفر برضه.»

Feedback B: «الصفر ممكن يكون القيمة الافتراضية للقراية، مش بالضرورة قيمة مخزنة.»

Feedback C: «صح. القراية بنتيجتين بترجع القيمة ومعاها هل المفتاح موجود فعلًا.»

#### `c63.presence-trace` · walkthrough
**Purpose:** observe increment, absence, stored zero, delete, and deterministic fmt output as separate small transitions.
**Links (author metadata; render natural titles):** `2-2-survey`.

**AR:** خلينا نشوف دولاب تدريب. `n, ok := counts["air"]` اسمها **comma-ok**: نتيجتين من قراية واحدة، شبه فكرة النتائج المتعددة في درس مسح المنطقة (2-2). هنا مش بنعمل لفة على خريطة المهمة.

<!-- go-example id=e63-presence imports=fmt outcome=ok -->
```go
counts := map[string]int{"water": 2}
counts["water"]++
fmt.Println(counts["air"], len(counts))
counts["air"] = 0
n, ok := counts["air"]
fmt.Println(n, ok, len(counts))
delete(counts, "water")
fmt.Println(counts)
```

| Frame | Highlight | Variables / map view | Cumulative output | AR caption |
|---|---|---|---|---|
| 0 | before 1 | — | — | «قبل الدولاب.» |
| 1 | 1 | water=2, len=1,total=2 | — | «مفتاح واحد.» |
| 2 | 2 | water=3,len=1,total=3 | — | «الوحدات زادت؛ المفاتيح لأ.» |
| 3 | 3 | unchanged; absent air not inserted | `0 1` | «قراية الغايب مش إضافة.» |
| 4 | 4 | water=3,air=0,len=2,total=3 | `0 1` | «المفتاح موجود بصفر.» |
| 5 | 5 | n=0,ok=true; map unchanged | `0 1` | «ok بتفصل الوجود عن القيمة.» |
| 6 | 6 | unchanged | `0 1\n0 true 2` | «عدد المفاتيح 2.» |
| 7 | 7 | air=0,len=1,total=0 | previous | «delete شالت مدخل water.» |
| 8 terminal | 8 | unchanged | `0 1\n0 true 2\nmap[air:0]` | «المفتاح بصفر لسه بيتطبع.» |

### Step `c63.watch-step` — Watch out

#### `c63.nil-trace` · walkthrough
**Purpose:** demonstrate readable nil map versus forbidden writes, stopping at the actual panic.

**AR:** توقّع الأول: هل القراية والزيادة هينجحوا الاتنين؟ تخيل معاك اسم دولاب، لكن الدولاب نفسه لسه متعملش. `var m map[string]int` بتبدأ nil. القراية ممكنة، لكن الزيادة محتاجة كتابة، فمش هتنفع.

<!-- go-example id=e63-nil-map imports=fmt outcome=panic -->
```go
var m map[string]int
n, ok := m["water"]
fmt.Println(n, ok, len(m))
m["water"]++
fmt.Println("after")
```

Frames: 0 before line 1, absent variables/view/output, caption «قبل التعريف». 1 line 1, m=nil,len=0,total=0,no cabinet, output none, caption «اسم من غير map جاهزة للكتابة». 2 line 2, n=0,ok=false,m unchanged, output none, caption «قراية الغايب مسموحة حتى هنا». 3 line 3, output `0 false 0`, same variables, caption «صفر مش دليل إن الكتابة آمنة». 4 terminal line 4, map still nil, no inserted entry, same stdout, panic `assignment to entry in nil map`, caption «الزيادة محتاجة كتابة؛ وقفنا هنا وafter مش هتتطبع». No recovery or initialization frame after the panic.

#### `c63.mistakes` · markdown
**Purpose:** express misconceptions and order guarantees accurately.

**AR:** **الافتراض:** «len هي إجمالي الموارد.» → **الواقع:** هي عدد المفاتيح. → **ليه؟** المورد الواحد ممكن يكون له قيمة كبيرة في مدخل واحد. → **فكّر في الإصلاح:** حدد هل محتاج عدد الفئات ولا جمع الكميات.

**الافتراض:** «map nil زي slice nil، أقدر أزوّدها عادي.» → **الواقع:** الكتابة للـmap nil بتعمل panic. → **ليه؟** مفيش map جاهزة تخزّن المدخل. → **فكّر في الإصلاح:** راجع إنشاء الـmap قبل أي كتابة، مش أول قراية بس.

**الافتراض:** «الطباعة طلعت مرتبة، يبقى range ماشية بنفس الترتيب.» → **الواقع:** ترتيب map مع range مش مضمون؛ ممكن يتكرر، وممكن يختلف. `fmt.Println` للـmap نفسها بتعرض المفاتيح بترتيب ثابت حسب قيمها. → **ليه؟** طريقة العرض مش عقد التكرار. → **فكّر في الإصلاح:** متخليش صحة العدّ معتمدة على مين ييجي الأول.

**Optional AR:** range على map بترجع key وvalue، مش index وvalue. متستخدمش أول ناتج على إنه رقم خانة. في map مفاتيحها Entity، الطباعة بترتب بالقيم الصحيحة اللي تحت أسماء الأنواع، مش أبجديًا حسب كلمة moss أو crystal.

#### `c63.print-trace` · walkthrough · Optional
**Purpose:** verify sorted fmt display without simulating a guaranteed random iteration sequence.

**AR:** توقّع الطباعة دي، لكن متستخدمهاش عشان تتوقع ترتيب range. إحنا بنطبع الـmap نفسها، مش سطور حلقة.

<!-- go-example id=e63-print imports=fmt outcome=ok -->
```go
stock := map[string]int{"water": 2, "air": 5}
fmt.Println(stock)
```

Frames: 0 before line 1, absent state/output, caption «قبل المخزن». 1 line 1, water=2,air=5,len=2,total=7, output none, caption «ده محتوى، مش ترتيب زيارة». 2 terminal line 2, unchanged state, stdout `map[air:5 water:2]`, caption «fmt رتبت مفاتيح string للعرض». A separate fixture `e63-range-orders` runs a real range repeatedly and records the actual observed sequences; none is presented as the mandatory next sequence.

#### `c63.ticks` · markdown · Optional
**Purpose:** apply action-cost ownership using a 3×2 toy field.
**Links (author metadata; render natural titles):** `ticks.actions`.

**AR:** تقدير لحقل تدريب 3×2: ست خانات، قراية GetEntity واحدة لكل خانة، وخمس حركات تزورهم من غير تكرار؛ يبقى **11 ticks**. زيادة عدّاد map من قيمة قريتها مش scan جديد. لو قريت نفس الخانة تاني فدي قراءة درون زيادة بتكلفة tick؛ ده غير تعديل العدد في الذاكرة.

الافتراض: البداية عند أول خانة، مفيش حصاد أو فحص إضافي أو رجوع، والتقرير نفسه مجاني. راجع «تكلفة أوامر الدرون» في 6-1 للتفاصيل. خليك فاكر إن الـstarter عندك عدّ خانة البداية فعلًا.

**Author tick audit:** live 6×4 minimal one-read-per-tile traversal=24+23=47 (Gold), Par=60, MaxTicks=300. No active field traversal code is supplied. Source: `ch06.go:388–462`.

### Step `c63.check-step` — Check

#### `c63.q-nil` · question
**Purpose:** assess map initialization before writing.

**AR prompt:** إنت محتاج map جديدة عشان تزود عداد water. أنهي بداية جاهزة للكتابة فورًا؟

Options: A «`var m map[string]int`»; B «`m := make(map[string]int)`»; C «اقرا مفتاح من m الأول، ده هيجهزها». **Correct: B.**

Feedback A: «دي nil؛ التعريف لوحده مش إنشاء map قابلة للكتابة.»

Feedback B: «صح؛ make أنشأت map، فتقدر تكتب فيها. الـliteral الفاضي ينفع برضه.»

Feedback C: «القراية مش بتهيّئ map ومش بتضيف مفتاح. راجع الفرق في التجربة.»

#### `c63.q-order` · question
**Purpose:** assess unspecified iteration, without asserting randomness every run.

**AR prompt:** range على map طبعت air وبعدها water. أنهي توقع آمن للفة جديدة على نفس المحتوى؟

Options: A «لازم نفس الترتيب»; B «لازم عكسه»; C «الاتنين ممكنين؛ مفيش ترتيب مضمون». **Correct: C.**

Feedback A: «التكرار ممكن يحصل، لكن مش ضمان تعتمد عليه في برنامجك.»

Feedback B: «Go مش بتتعهد تغيّر الترتيب في كل لفة.»

Feedback C: «أيوه؛ صحة الحساب لازم تفضل واحدة لأي ترتيب مسموح.»

### Step `c63.turn-step` — Your turn

#### `c63.turn` · markdown
**Purpose:** one small inspection action before implementing traversal.

**AR:** شغّل الـstarter مرة، وحدد في الطباعة أنهي مفتاح اتعدّ لخانة البداية. اكتب تعليق يوضح إن الخانة دي متسجلة بالفعل؛ بعدين صمّم زيارتك لباقي الخريطة من غير ما تضيفها بالغلط مرة تانية.

### Three hints

1. **Concept — `c63.h1`:** «المطلوب عدد الخانات لكل نوع، مش حجم المورد ولا مجموع كل الأنواع في رقم واحد.»
2. **Structure — `c63.h2`:** «راجع زيارتك للخانات من الفصل 3: هل كل خانة بتساهم مرة واحدة، بما فيها اللي الـstarter سجلها؟»
3. **Specific — `c63.h3`:** «بص على المفتاح اللي بتزود قيمته: هل هو نوع الخانة الحالية؟ وبص على التقرير: هل بتبعت map الأعداد، مش len بتاعتها؟»

**Leakage review:** player still builds traversal, coordinates observations with updates, handles the existing starting count, and chooses report timing. No snake-route code or exact `counts[colony.GetEntity()]++` expression is supplied by hints. Creating a map and understanding a one-key increment teach the data structure, not the full field algorithm.

**Fallback corrections:** replace “any key type” with equality-comparable key types. Replace “random every time” with unspecified order that may repeat. Keep fmt ordering but clarify Entity numeric order. Make the nil-map example terminal; show initialization as a separate alternative before writing, not an unreachable line after panic. Print/use len in a runnable example rather than a bare `len(counts)` statement. Distinguish entity counts from resource sizes and len from totals. Replace the old route-plus-counting-expression hint.

## 6-4 — ترتيب المخزون (`6-4-stock-order`)

**Teaching objective:** interpret comparator signs; preserve item identity while ordering; make ties explicit. Owns `sort.comparator`, `sort.identity-ties`. Links to `slice.sort`, `slice.sharing`, `slice.basics`, and function values/closures from `5-5-command-table`; these are not introduced again.

**Starter audit:** only prints one `colony.Measure()` result at the starting tile. That measurement has already cost one tick but is **not saved** in a variable. It does not build arrays, move, sort, or report. The player may revise that debugging statement to retain its value or remove it as part of their own implementation; this storyboard changes no starter. Width is 6–10, must come from `WorldSize`; first and last sizes tie by construction. No structs required. Exactly one correct report, no successful harvests; `slices.SortFunc` source feature required.

### Step `c64.mission-step` — Mission

#### `c64.mission` · markdown
**Purpose:** give observable contract, not a comparator implementation.
**Links (author metadata; render natural titles):** `slice.sort`.

**AR:** القاعدة عايزة **المواقع**، مش الأحجام: الحجم الأكبر يسبق، ولو الحجم اتعادل يبقى x الأصغر يسبق. اقرأ عرض الصف من `colony.WorldSize`؛ مش رقم ثابت. ابعت ترتيب المواقع كله بتقرير واحد، من غير حصاد، واستخدم `slices.SortFunc`.

الـstarter طبع قراءة واحدة بس؛ مقاسَتش ببلاش، ومتحفظتش في متغير. خد قرارك هتحتفظ بالقراية دي إزاي قبل ما تبدأ تسجيلك. راجع علاقة الحجم بالموقع من «الترتيب والمواقع» في 6-2.

### Step `c64.idea-step` — Idea

#### `c64.idea` · markdown
**Purpose:** introduce comparator as a consistent decision rule; connect to existing function values.
**Links (author metadata; render natural titles):** `5-5-command-table`, `slice.sort`.

**AR:** تخيل موظف بينظم دور خدمة. كل ما تقوله «مين يسبق من الاتنين دول؟» محتاج قاعدة واحدة واضحة، مش إجابة تتغير كل مرة.

الـ**comparator** دالة مقارنة بترجع int: **سالب** لو a تسبق b، **صفر** لو متساويين حسب القاعدة، **موجب** لو a تيجي بعد b. مش بترجع true أو false. `SortFunc` بتاخد الدالة كقيمة، زي درس لوحة الأوامر (5-5)، وبتعدّل ترتيب عناصر الـslice نفسها زي «الترتيب والمواقع».

**Optional AR:** صفر مش معناه «رتبهم حسب ذوقك»؛ معناه القاعدة بتعتبرهم متساويين. لازم القاعدة تبقى متسقة، والعنصر مع نفسه يرجع صفر. `SortFunc` مش بتضمن الحفاظ على ترتيب العناصر اللي المقارنة اعتبرتها متساوية. قاعدة التعادل الواضحة هي اللي بتحدد الأول هنا، مش الاعتماد على مصادفة ترتيب سابق.

### Step `c64.example-step` — Example

#### `c64.comparator-trace` · walkthrough
**Purpose:** small one-key ascending comparison, not a mission multi-key comparator.

**AR:** في صندوق تذاكر تدريب، عايزين الأرقام الأصغر الأول. توقّع النتيجة قبل «التالي». هنا بنشوف معنى علامة المقارنة؛ المهمة ليها قاعدة مختلفة إنت هتكتبها.

<!-- go-example id=e64-comparator imports=fmt,slices outcome=ok -->
```go
tickets := []int{8, 3, 5}
before := func(a, b int) int {
    if a < b { return -1 }
    if a > b { return 1 }
    return 0
}
fmt.Println(before(3, 8), before(8, 3), before(5, 5))
slices.SortFunc(tickets, before)
fmt.Println(tickets)
```

| Frame | Highlight | Variables / slice view | Cumulative output | AR caption |
|---|---|---|---|---|
| 0 | before 1 | — | — | «قبل التذاكر.» |
| 1 | 1 | tickets=[8 3 5] | — | «ترتيب البداية.» |
| 2 | 2–6 declaration | before=function value; no body executed | — | «تعريف دالة، مش نداء ليها.» |
| 3 | 7 first call →3 | a=3,b=8; result=-1; tickets unchanged | — | «3 يسبق 8.» |
| 4 | 7 second call →4 | a=8,b=3; first condition false; result=1 | — | «8 ييجي بعد 3.» |
| 5 | 7 third call →5 | a=5,b=5; both conditions false; result=0 | — | «المتساوي مع نفسه صفر.» |
| 6 | 7 after calls | a/b out of scope | `-1 1 0` | «دي نتائج النداءات اللي طلبناها صراحة.» |
| 7 | 8 after library call | tickets=[3 5 8]; before unchanged | `-1 1 0` | «SortFunc نداء واحد في الرسم؛ مش بندّعي ترتيب نداءات Go الداخلية.» |
| 8 terminal | 9 | unchanged | `-1 1 0\n[3 5 8]` | «الترتيب بقى تصاعدي.» |

#### `c64.identity-trace` · walkthrough
**Purpose:** distinguish position IDs from the data they identify without supplying sorting of positions.

**AR:** رقم بطاقة الرف غير الكلام المكتوب عليها. السجل ده أسماء رفوف، وورقة تانية فيها أرقام بطاقاتهم. لما أول رقم في الورقة يبقى 2، ده مش معناه قيمة الرف 2؛ دي هويته. توقّع أنهي اسم هنقراه.

<!-- go-example id=e64-identity imports=fmt outcome=ok -->
```go
labels := []string{"ink", "paper", "tape"}
order := []int{2, 0, 1}
id := order[0]
fmt.Println(id, labels[id])
fmt.Println(labels)
```

Frames: 0 before 1, no state/output, caption «قبل الرفوف». 1 line 1, labels=[ink paper tape] in A, caption «كل اسم له رقم ثابت في الدفتر». 2 line 2, order=[2 0 1] in B, labels unchanged, caption «ورقة أرقام، مش نسخة من الأسماء». 3 line 3, id=2, highlight order[0] then labels index 2, caption «أول بطاقة بتشير للرف رقم 2». 4 line 4, stdout `2 tape`, caption «الهوية 2 والاسم tape». 5 terminal line 5, stdout `2 tape\n[ink paper tape]`, unchanged state, caption «ترتيب دفتر الأسماء متغيرش». Output before frame 4 is empty. No ordering algorithm occurs in this example.

### Step `c64.watch-step` — Watch out

#### `c64.mistakes` · markdown
**Purpose:** prevent bool comparator, unstable ties, and accidental identity loss.
**Links (author metadata; render natural titles):** `5-5-command-table`, `c64.identity-trace`, `slice.sharing`.

**AR:** **الافتراض:** «السؤال مين يسبق، فهرجع bool.» → **الواقع:** SortFunc مستنية int له 3 حالات. → **ليه؟** محتاجة تفرق قبل/متساوي/بعد. → **فكّر في الإصلاح:** حدد معنى العلامة قبل كتابة المقارنة.

**الافتراض:** «حجمين متساويين؟ أرجع صفر وخلاص وهما يفضلوا بالترتيب المطلوب.» → **الواقع:** كده لسه ماقلتَش مين يسبق حسب الموقع، وSortFunc مش stable. → **ليه؟** التعادل في الحجم مش كل قاعدة المهمة. → **فكّر في الإصلاح:** اختبر اتنين متساويين في الحجم واسأل عن المعيار التاني؛ إنت اللي هتربطه بالدالة.

**الافتراض:** «أرتب القيم وأعتبر index الجديد مكانها القديم.» → **الواقع:** ضيّعت علاقة الهوية بالقيمة. → **ليه؟** index موضع جوه سجل، مش إثبات موقع في العالم. → **فكّر في الإصلاح:** ارجع لـ«مشاركة مخزن العناصر» و«بطاقة الرف وقيمته» وحدد إيه اللي هيتحرك وإيه اللي لازم يفضل مرجع ثابت.

**Optional AR:** طرح عددين عشان تطلع علامة المقارنة ممكن يعمل overflow لو الأعداد كبيرة. المثال استخدم شروط واضحة بدل تعميم اختصار الطرح. والـclosure تقدر تقرا بيانات برّه جسمها زي ما شوفنا في درس لوحة الأوامر (5-5)؛ ده مفيد لو عنصر الترتيب هو هوية، لكن تصميم المقارنة في المهمة لسه عليك.

#### `c64.ticks` · markdown · Optional
**Purpose:** apply scan cost to variable-width toy data, explicitly accounting for starter debugging.
**Links (author metadata; render natural titles):** `ticks.actions`.

**AR:** تجربة رفوف من 4 علامات: قراءة واحدة لكل علامة و3 حركات بينهم = **7 ticks**. المقارنة والترتيب في الذاكرة مش أوامر درون. لو قريت أول علامة للتجربة وبعدين تجاهلت القراية وبدأت قراءة الكل من جديد، هتزود tick من غير معلومة جديدة.

ده تقدير، مش عرض خريطتك. افترض البداية عند أول علامة، مفيش حصاد ولا أوامر زيادة، وتقرير واحد بعد الحساب. راجع «تكلفة أوامر الدرون»، وراجع `fmt.Println(colony.Measure())` في الـstarter قبل حساب قراءاتك.

**Author tick audit:** live width n=6..10, one-read-per-tile cost `2n-1`=11..19; Par=Gold=19 and MaxTicks=50 are fixed. Do not present 19 as mandatory actual cost or recompute star thresholds as `2n-1`. Source: `ch06.go:510–665`.

### Step `c64.check-step` — Check

#### `c64.q-sign` · question
**Purpose:** assess sign semantics independently of mission values.

**AR prompt:** في قاعدة تدريب بتقول «a يسبق b»، أنهي نتيجة مناسبة من comparator؟

Options: A «-1»; B «0»; C «1». **Correct: A.**

Feedback A: «صح؛ السالب بيقول a يسبق b، ومش لازم يكون -1 بالذات.»

Feedback B: «الصفر بيقول القاعدة بتعتبرهم متساويين، مش إن a يسبق.»

Feedback C: «الموجب بيحط a بعد b؛ اعكس تصورك للعلامة.»

#### `c64.q-bool` · question
**Purpose:** recognize an intentional compile failure rather than silently fix the mission comparator.

**AR prompt:** الكود ده هيرتب التذاكر ولا فيه مشكلة قبل التشغيل؟

<!-- go-example id=e64-bool-error imports=fmt,slices outcome=compile-error -->
```go
tickets := []int{8, 3}
slices.SortFunc(tickets, func(a, b int) int {
    return a < b
})
fmt.Println(tickets)
```

Options: A «هيترتب تصاعدي»; B «Compile error: رجع bool مكان int»; C «Panic بعد أول مقارنة». **Correct: B.**

Feedback A: «ده ينفع كشكل مقارنة في APIs تانية، لكن مش عقد SortFunc هنا.»

Feedback B: «أيوه؛ قبل التشغيل، نوع النتيجة مش النوع المكتوب للدالة.»

Feedback C: «المشكلة معروفة وقت البناء، فمش هنوصل لنداء الدالة أصلًا.»

#### `c64.q-tie` · question
**Purpose:** assess explicit tie rule without the code that enforces it.

**AR prompt:** بطاقتين في تمرين تنظيم: نفس الأولوية، وأرقامهما 7 و2. القاعدة المكتوبة: الأولوية أولًا، وبعدها الرقم الأصغر. هل مساواتهما في المقارنة كفاية لضمان ظهور 2 الأول مع SortFunc؟

Options: A «أيوه، المساواة كفاية»; B «لأ؛ لسه لازم المقارنة تعبر عن فرق الرقم»; C «الترتيب ثابت لو حطيناهم في map». **Correct: B.**

Feedback A: «المساواة تجاهلت المعيار التاني، وSortFunc مش بتضمن الحفاظ على ترتيب المتساويين.»

Feedback B: «صح؛ القاعدة ذات معيارين لازم تميزهم لما المعيار الأول يتعادل.»

Feedback C: «map مش بتضمن ترتيب range. تغيير المخزن مش بديل عن قاعدة مقارنة كاملة.»

### Step `c64.turn-step` — Your turn

#### `c64.turn` · markdown
**Purpose:** one preparatory identity decision, not initialization or comparator code.

**AR:** قبل ما تكتب SortFunc، اكتب تعليق يحدد «العنصر اللي هرتّبه بيمثل إيه؟». جرّب تفسير إجابتك على اتنين حجمهم متساوي؛ لو الموقع الأصلي اختفى من تفكيرك، ارجع لنموذج الهوية.

### Three hints

1. **Concept — `c64.h1`:** «القاعدة بتطلب مواقع مرتبة، مش قايمة أحجام مرتبة. اسأل نفسك مين صاحب كل قيمة.»
2. **Structure — `c64.h2`:** «اختبر قرار المقارنة في حالتين على ورقة: حجمين مختلفين، وحجمين متساويين. هل كل حالة ليها قرار واضح؟»
3. **Specific — `c64.h3`:** «راجع نوع ناتج comparator ومعنى السالب. لو مثال التعادل عندك بيرجع صفر رغم اختلاف الموقعين، قاعدتك لسه ناقصة.»

**Leakage review:** no full position-list initialization, collection loop, lookup-based two-key comparator, or combined “compare sizes then subtract positions” recipe. Player chooses storage, obtains measurements without accidental duplication, implements both criteria and identity preservation, and sends exactly one report. The ascending single-key SortFunc example teaches its actual signature; the shelf example never sorts IDs.

**Fallback corrections:** remove the complete ranks/order two-key comparator, which nearly implements the mission. Keep function-value links to 5-5 and clarify comparator parameters are slice elements, not automatically sort-internal indices. Preserve variable width, descending size/ascending-x tie contract, no harvesting, exactly one report. Prefer explicit comparisons to general subtraction; say SortFunc is not stable. Explain the discarded starter measurement and fixed 19 target. Replace the old near-complete last hint.

## 6-5 — شبكة الحرارة (`6-5-heat-grid`)

**Teaching objective:** represent rectangular numeric grids with independently owned rows; reason about valid shapes and three mathematical operations; use the first lab's function/test contract. Owns `grid.shape`, `grid.row-ownership`, `grid.operations`, `lab.contract`. Links to `slice.make-append`/`slice.bounds`, `slice.sharing`, `range.slice`, nested loops in `3-2-full-field`, float64 in `2-2-survey`, and errors/functions in `5-1-functions`/`5-2-storm`/`5-3-fault-log`.

**Starter audit:** `main.go` already declares `OverlayHeat`, `ThroughHeat`, and `FlipHeat`, returning nil results and descriptive “implement …” errors. `main` is empty. There is no grid allocated or arithmetic implemented. Run invokes tests; no drone traversal or Report is required. Do not ask the player to write tests or insert drone calls. Keep signatures. The lab's supplied hidden tests check rectangularity, operation dimensions, correct numbers, input preservation on success/error, and output-row independence. Tests use relative tolerance 1e-9 with absolute floor 1e-9. No benchmark/bonus/tick target.

### Step `c65.mission-step` — Mission + first-lab orientation

#### `c65.mission` · markdown
**Purpose:** orient to a lab and give the behavioural contract without a solution skeleton.

**AR:** المرة دي Run مش هيحرّك درون: دي أول **lab**، يعني تمرين بيتراجع باختبارات. افتح `main.go` وكمّل الدوال الثلاثة اللي موجودة؛ متغيّرش أسماءها أو المدخلات والنواتج. `main` ينفع تفضل فاضية.

`OverlayHeat` تجمع شبكتين خانة بخانة. `ThroughHeat` تعمل ضرب مصفوفات، مش ضرب كل خانتين قصاد بعض. `FlipHeat` تبدّل الصفوف بالأعمدة. كل نتيجة لازم تكون شبكة جديدة، ومدخلاتك تفضل زي ما كانت حتى لو فيه error.

مش مطلوب `colony.Report` ولا اختبارات تكتبها بنفسك. شوف نتيجة كل خطوة في الـConsole؛ الاختبارات هتقولك أنهي شرط لسه ناقص. مفيش سباق ticks في الـlab دي.

### Step `c65.idea-step` — Idea

#### `c65.idea` · markdown
**Purpose:** introduce slice-of-rows and shape before computation.
**Links (author metadata; render natural titles):** `2-2-survey`, `5-3-fault-log`, `slice.basics`.

**AR:** تخيل دولاب فيه أرفف، وكل رف شايل صف قراءات. `[][]float64` معناها slice عناصرها صفوف، وكل صف `[]float64`. أول index بيختار الصف، والتاني بيختار العمود. float64 عرفناه في درس مسح المنطقة (2-2)، والصف نفسه هو الـslice اللي عرفناها في «الـslice والعناصر».

الـouter slice بتعد الصفوف؛ هي مش ورقة مستطيلة تلقائيًا. ممكن صف يبقى أقصر من التاني، أو فاضي. **rectangular** يعني كل الصفوف بنفس عدد الأعمدة. ده شرط الشبكة في المهمة، مش شرط كل slice من slices في Go.

**Optional AR:** متقراش أول صف قبل ما تعرف إن فيه صفوف، ومتعتبرش طول أول صف دليل على الباقي. الإدخال الغلط هنا يرجع error مفهوم زي درس سجل الأعطال (5-3)، مش panic. الأرقام السالبة والعشرية صالحة؛ فحص الشكل مش مسح للقيم اللي مش عاجباك.

### Step `c65.example-step` — Example

#### `c65.q-index` · question
**Purpose:** predict a row/column selection before visual reveal.

**AR prompt:** شبكة تدريب صفوفها `[3 7 2]` وبعدين `[8 6 4]`. الخانة ذات index `[1][0]` قيمتها كام؟

Options: A «7»; B «8»; C «6». **Correct: B.**

Feedback A: «7 في الصف 0 والعمود 1؛ بدّلت ترتيب المؤشرين.»

Feedback B: «صح؛ الصف رقم 1 هو التاني، والعمود 0 أول خانة فيه.»

Feedback C: «6 عند [1][1]. أول عمود رقمه صفر، مش واحد.»

#### `c65.grid-trace` · walkthrough
**Purpose:** labelled row/column view and lengths, with concrete output.

**AR:** خلينا نشوف نفس لوحة التدريب. خلي عينك على اختيار الصف الأول، وبعده العمود، بدل ما تشوف `[1][0]` كرقم واحد.

<!-- go-example id=e65-grid-index imports=fmt outcome=ok -->
```go
board := [][]float64{{3, 7, 2}, {8, 6, 4}}
fmt.Println(len(board), len(board[0]))
fmt.Println(board[1][0])
```

Frames: 0 before line 1, absent state/output, caption «قبل اللوحة». 1 line 1, outer descriptor O len=2, row descriptors point to A=[3 7 2],B=[8 6 4], each len=3; grid row labels 0,1 and column labels 0,1,2; output none; caption «صفين، كل واحد 3 عناصر». 2 line 2, unchanged view with outer length and row-0 length highlighted; stdout `2 3`; caption «عدد الصفوف، وبعدين طول أول صف». 3 terminal line 3, highlight row 1 then cell column 0; stdout `2 3\n8`; caption «اختارنا 8». Capacities are not asserted; no matrix operation was executed.

#### `c65.shallow-trace` · walkthrough · Optional
**Purpose:** apply `slice.sharing` to row descriptors, exposing shallow outer cloning.
**Links (author metadata; render natural titles):** `slice.sharing`.

**AR:** فاكر Clone في «مشاركة مخزن العناصر»؟ هناك العناصر كانت int. هنا العنصر نفسه وصف صف. توقّع قبل «التالي»: هل تغيير النسخة هيبان في المصدر؟ خلينا نشوف حدود النسخ السطحي بدل ما نفترض إنه نسخ كل الشبكة.

<!-- go-example id=e65-shallow imports=fmt,slices outcome=ok -->
```go
board := [][]float64{{2, 6}, {5, 7}}
copyBoard := slices.Clone(board)
copyBoard[0][1] = 9
fmt.Println(board)
fmt.Println(copyBoard)
```

| Frame | Highlight | Variables / grid view | Cumulative output | AR caption |
|---|---|---|---|---|
| 0 | before 1 | — | — | «قبل الشبكة.» |
| 1 | 1 | board outer O → row stores A=[2 6],B=[5 7] | — | «وصفين صفوف جوه outer slice.» |
| 2 | 2 | copyBoard outer P≠O; P[0]/O[0]→A, P[1]/O[1]→B | — | «الـouter اتنسخت، لكن مخازن الصفوف لأ.» |
| 3 | 3 | A=[2 9],B=[5 7]; both grids see A | — | «التعديل وصل لنفس الصف المشترك.» |
| 4 | 4 | unchanged | `[[2 9] [5 7]]` | «المصدر اتغير.» |
| 5 terminal | 5 | unchanged | `[[2 9] [5 7]]\n[[2 9] [5 7]]` | «Clone للـouter لوحدها مش استقلال شبكة.» |

#### `c65.rows-trace` · walkthrough
**Purpose:** show independent row allocation on a fixed 2×2 blank board, without a general copy/computation loop.

**AR:** زي أرفف منفصلة: تخصيص الدولاب مش هو تخصيص كل رف. دي لوحة فاضية صغيرة، مش تنفيذ أي دالة من المهمة. توقّع هل تغيير أول صف هيغيّر التاني.

<!-- go-example id=e65-rows imports=fmt outcome=ok -->
```go
blank := make([][]float64, 2)
blank[0] = make([]float64, 2)
blank[1] = make([]float64, 2)
blank[0][0] = 7
fmt.Println(blank)
```

Frames: 0 before line 1, absent state/output, caption «قبل التخصيص». 1 line 1, O=[nil nil], two nil row descriptors, no inner cells; output none; caption «صفين لسه من غير عناصر». 2 line 2, O[0]→A=[0 0], O[1]=nil, caption «أول رف اتعمل». 3 line 3, O[1]→B=[0 0], A≠B, caption «رف تاني مستقل». 4 line 4, A=[7 0],B=[0 0], caption «التعديل في أول رف بس». 5 terminal line 5, same variables/grid, stdout `[[7 0] [0 0]]`, caption «الصف التاني فضل صفر». No output before frame 5. This contains no loops or dynamic dimensions.

#### `c65.operations-trace` · walkthrough
**Purpose:** define the mathematics through three separate one-cell examples; no player function implementations.

**AR:** العمليات ليها معاني مختلفة. شوف مثال واحد كل مرة؛ مفيش دالة من دوال مهمتك شغالة هنا.

**Subview A — AR title «جمع خانة بخانة»:** لو شكل الشبكتين واحد، الخانة في الناتج بتجمع القيمتين اللي في نفس المكان. مثال خليتين في أول صف: 1.5 و2.5 ناتجهم 4.

<!-- go-example id=e65-overlay-cell imports=fmt outcome=ok -->
```go
left := []float64{1.5, -2}
right := []float64{2.5, 3}
fmt.Println(left[0] + right[0])
```

Frames A: 0 before 1, absent state/output, «قبل القراءات». 1 line 1, left=[1.5 -2], row view columns 0,1, «صف من المصدر الأول». 2 line 2, right=[2.5 3], both rows visible and unchanged, «صف من المصدر التاني». 3 terminal line 3, highlight both index-0 cells, show scalar result=4 separately (not an allocated output grid), stdout `4`, «حسبنا خانة واحدة بس». No output before the final frame.

**Subview B — AR title «صف في عمود»:** ضرب المصفوفات بيجمع نواتج الضرب على البعد المشترك. صف `[1.5 -2 4]` وعمود `[2 0 -1]` يدوا `1.5×2 + (-2)×0 + 4×(-1) = -1`. مش بنضرب خانتين متناظرتين وبس.

<!-- go-example id=e65-dot-cell imports=fmt outcome=ok -->
```go
row := []float64{1.5, -2, 4}
column := []float64{2, 0, -1}
value := row[0]*column[0] + row[1]*column[1] + row[2]*column[2]
fmt.Println(value)
```

Frames B: 0 before 1, absent state/output, «تجربة جديدة». 1 line 1, row stored A=[1.5 -2 4], «ده الصف». 2 line 2, column stored B=[2 0 -1] rendered vertically with indices 0,1,2, «وده العمود». 3 line 3 after expression, value=-1; paired indices highlighted and terms shown as `3 + 0 + (-4)`; stores unchanged, «كل مكان في البعد المشترك ساهم في نفس الخانة». 4 terminal line 4, stdout `-1`, «الناتج لخانة واحدة». The term highlighting is an explanation of the expression, not invented extra Go statement execution.

**Subview C — AR title «Transpose: بدّل المحورين»:** في شبكة 2×3، الخانة [1][0] بتروح لمكان [0][1] في شبكة 3×2. القيمة نفسها مش بتتغير؛ إحداثياتها هي اللي بتتبدّل. ده اسم **transpose**، مش قلب أرقام الصف من الآخر للأول.

**Code:** reuse the exact `e65-grid-index` fragment, with its original three lines and stdout `2 3\n8`; this view adds a **mathematical destination overlay**, not a call to FlipHeat. Frames C: same 0–3 execution as `c65.grid-trace`; at terminal line 3 highlight source [1][0]=8 and an unallocated destination stencil 3×2 with only [0][1] marked 8 and all other cells “not computed”, caption «ده مكان نقل القيمة حسب تعريف العملية؛ لسه مفيش كود بينفذ النقل». No additional output or claim that Go constructed that stencil.

### Step `c65.watch-step` — Watch out

#### `c65.mistakes` · markdown
**Purpose:** separate shape validation, data mutation, row ownership, and arithmetic misconceptions.
**Links (author metadata; render natural titles):** `slice.bounds`.

**AR:** **الافتراض:** «أول صف طوله 3، يبقى كل الصفوف كده.» → **الواقع:** كل صف slice مستقلة وممكن طولها يختلف. → **ليه؟** النوع مش بيفرض مستطيل. → **فكّر في الإصلاح:** اثبت صلاحية الشكل كله قبل أي فهرسة تفترض عرض ثابت؛ راجع «حدود الفهرسة».

**الافتراض:** «outer جديدة يبقى الشبكة مستقلة.» → **الواقع:** صفوفها ممكن تشير لنفس مخازن المصدر، أو صفين يشاوروا لنفس المخزن. → **ليه؟** عناصر الـouter أوصاف صفوف. → **فكّر في الإصلاح:** راجع مين مالك عناصر كل صف، مش اسم الشبكة بس.

**الافتراض:** «ThroughHeat ضرب كل قيمة في اللي قصادها.» → **الواقع:** كل خانة ناتجة من صف في عمود ومجموع على البعد المشترك. → **ليه؟** ده تعريف ضرب المصفوفات. → **فكّر في الإصلاح:** ارسم مساهمات خانة واحدة قبل تعميمها بحلقاتك.

**Optional AR:** اكتشاف خطأ بعد ما تغيّر المدخلات مش مقبول. الاختبارات بتراجع إن المدخلات متغيرتش حتى في مسار error. متعالجش اختلاف الأبعاد بقص صفوف أو إضافة أصفار من عندك؛ رجّع خطأ بدل تغيير معنى العملية.

### Step `c65.check-step` — Check

#### `c65.q-shape` · question
**Purpose:** assess valid rectangular/dimension reasoning.

**AR prompt:** A شكلها 2×3 وB شكلها 3×2، وكل الصفوف سليمة وغير فاضية. أنهي وصف صح؟

Options: A «OverlayHeat وThroughHeat الاتنين صالحين»; B «ThroughHeat صالحة وناتجها 2×2؛ OverlayHeat ترفض اختلاف الشكل»; C «ThroughHeat ناتجها 3×3». **Correct: B.**

Feedback A: «الجمع خانة بخانة محتاج نفس الشكل، و2×3 مختلفة عن 3×2.»

Feedback B: «صح؛ البعد المشترك 3 متوافق، والنتيجة صفوف A في أعمدة B.»

Feedback C: «3 هو البعد اللي بيتجمع عليه، مش أبعاد الناتج. الناتج صفوف A وأعمدة B.»

#### `c65.q-rows` · question
**Purpose:** assess row sharing without giving a general allocation loop.

**AR prompt:** اتنين صفوف بيشيروا لنفس slice اسمها row. بعد التعديل ده، أنهي طباعة هتحصل؟

<!-- go-example id=e65-shared-rows imports=fmt outcome=ok -->
```go
row := []float64{2, 6}
board := [][]float64{row, row}
board[0][0] = 9
fmt.Println(board)
```

Options: A «`[[9 6] [2 6]]`»; B «`[[9 6] [9 6]]`»; C «Compile error؛ تكرار row ممنوع». **Correct: B.**

Feedback A: «ده كان محتاج مخزنين مستقلين؛ هنا الوصفين بيشوفوا نفس عناصر row.»

Feedback B: «أيوه؛ الكتابة لنفس مخزن العناصر باينة من الصفين. اختلاف الـindex الخارجي مش ضمان الاستقلال.»

Feedback C: «Go بتسمح بالمشاركة؛ استقلال الصفوف شرط المهمة اللي إنت لازم تحققه.»

#### `c65.contract` · markdown
**Purpose:** complete the mathematical/error contract as short main text plus optional details.

**AR:** قبل الحساب، الشبكة لازم فيها صف واحد على الأقل، وكل صف فيه عمود واحد على الأقل، وكل الصفوف نفس الطول. صف فاضي أو صفوف بأطوال مختلفة يتعاملوا كإدخال غلط في المهمة دي.

OverlayHeat محتاجة نفس عدد الصفوف والأعمدة في المدخلين. ThroughHeat محتاجة **أعمدة A = صفوف B**، والناتج **صفوف A × أعمدة B**. FlipHeat بتحوّل R×C إلى C×R، وتنقل [i][j] إلى [j][i]. دي مواصفات النتيجة؛ ترتيب الفحص والتخصيص والحلقات إنت اللي هتبنيه.

**Optional AR:** الشبكات 1×1 وصف أو عمود واحد صالحين. ارفض nil وصفر صفوف وصفر أعمدة والصفوف غير المتساوية، وارفض الأبعاد غير المناسبة للعملية. استخدم error غير nil برسالة مفهومة؛ النص مش لازم يطابق جملة حرفية. كل ناتج مستقل عن المدخلات، وكل صف مستقل عن باقي صفوف الناتج. مقارنة الأرقام في الاختبارات بسماحية نسبية `1e-9` وبحد أدنى مطلق `1e-9`؛ مش مطلوب منك اكتشاف overflow في float64. دي سماحية فحص، مش طلب تقريب أو تعديل الأرقام بإيدك.

### Step `c65.turn-step` — Your turn

#### `c65.turn` · markdown
**Purpose:** one concrete first-lab action without a validation/allocation/loop skeleton.

**AR:** افتح `main.go` وشغّل الـstarter زي ما هو مرة. اختار أول فشل ظاهر، واكتب في تعليق شرط واحد في عقد الدالة يفسره. اشتغل على الشرط ده الأول بدل ما تحاول تكتب الدوال كلها مرة واحدة.

### Three hints

1. **Concept — `c65.h1`:** «كل عملية ليها معنى مختلف. ارسم خانة ناتج واحدة واسأل جت منين، قبل ما تفكر في تعميمها.»
2. **Structure — `c65.h2`:** «افصل في تفكيرك بين صلاحية شكل المدخل وصلاحية الأبعاد للعملية. الاختبار اللي بيقع بيقولك أي وعد لسه متحققش.»
3. **Specific — `c65.h3`:** «راجع سطر make للـouter وسطر إنشاء كل صف، أو نداء slices.Clone لو استخدمته: النسخ السطحي للـouter بينسخ أوصاف الصفوف، مش عناصرها. تغيير خلية في الناتج لازم ما يغيّرش المصدر ولا أي صف تاني.»

**Tick explanation:** none. Learner text already states “مفيش سباق ticks”. Hide estimate/measured tick widgets here. Run is a lab-test result, not a zero-tick achievement. Successful ordinary hidden checks give three stars; no performance target is invented.

**Hint link metadata:** c65.h3 → `c65.shallow-trace`; display the natural walkthrough title, never the internal ID.

**Leakage review:** player still implements full shape/error checking, operation-specific dimensions, dynamic independent result allocation, all traversal/accumulation/index-transfer loops, and input preservation. Fixed 2×2 blank allocation is not a general nested-loop copy skeleton. Mathematical definitions and one scalar worked calculation are necessary specifications; no general matrix function, triple-loop layout, or exact function-specific fix is given. Existing solved grid-doubling loop is deliberately omitted.

**Fallback corrections:** remove the solved general nested-loop allocation/doubling example and loop-role hint. Keep the complete public function/error contract and all shape cases. Explicitly qualify shallow outer Clone and independent output rows; distinguish transpose from row reversal and matrix multiplication from elementwise multiplication. Do not present calls to unimplemented functions as verified outputs. Preserve first-lab guidance, tolerance, and absence of Report/ticks.

## English learner copy — 6-1 only

The IDs, kinds, purposes, code fragments, frame state data, option keys, and correct answers are identical to 6-1 above. This section translates every learner-facing block, optional note, frame caption, sandbox message, and hint; it does not add examples. Common controls use the EN strings in the opening conventions. Step titles: Mission / Idea / Example / Watch out / Check / Your turn. Hint labels by position: Concept / Structure / Specific.

### EN `c61.mission` — markdown

Base needs a log of all eight crystal sizes, in their left-to-right order. Each number matters, and so does its position in the log.

The starter has already measured and stored the first crystal. The drone is still at the starting tile. `fmt.Println` shows you the log in the Console; it does not send it to base. You met that distinction in Inventory (2-1).

Build a complete log in a slice and send it with `colony.Report`. No harvesting is needed. You decide how to continue observing and moving.

### EN `c61.idea` — markdown

Imagine a drawer with numbered compartments and a card saying how many are in use. A **slice** describes part of an element store: one element type, with elements in order. `[]int` means a slice of integers.

An **index** is an element's number: the first is 0, then 1. `len(s)` counts existing elements. `cap(s)` is the room available from the slice's start in its **backing array**, the underlying element store.

`append` adds at the end and returns an updated slice descriptor. Keep that result; the old descriptor does not gain length by itself. Adding an element does not always create a new store.

**Optional:** For a slice, `len` counts elements, rather than bytes as it did for a string in Callsigns (2-3). `var s []int` starts nil, with zero length and capacity, and you can append to it. The store drawing is a model for understanding, not a map of real memory addresses.

### EN `c61.basic-trace` — walkthrough

Let's try two library tickets, away from the crystals. Before each Next, predict what the line will change. Then watch the elements and their count, rather than just the brackets.

Code and frame states: `e61-basic`, unchanged. Frame captions:

| Frame | English caption |
|---|---|
| 0 | “No variable yet.” |
| 1 | “An empty descriptor, with no elements.” |
| 2 | “Zero elements and zero capacity.” |
| 3 | “We kept the returned descriptor. The first index is 0.” |
| 4 | “There is now one ticket.” |
| 5 | “Append adds at the end; it does not replace the first element.” |
| 6 | “The first element is 7, and the length is 2.” |

### EN `c61.make-playground` — sandbox

Try it yourself: choose a starting preset, then press Append 7. Are any elements already present? Watch `len` beside `cap`.

Preset buttons retain their Go expression labels. `make([]int, 3)` starts with three existing zeros and length/capacity 3. `make([]int, 0, 3)` starts with zero elements and capacity 3. `[]int{}` is a non-nil empty slice with length/capacity 0.

Operation buttons: “Append 7”, “Reset”, and the three presets. At six elements, disable Append with “This experiment stops at 6 elements; reset to try another starting point.” Changing presets clears the experiment.

Reserve-cell label: “Reserved room; not elements you can index yet.” The fixed operation effects and store-label changes are the same as the Arabic version. Append keeps the old elements and adds one final 7, using the returned descriptor. With room it reuses storage; without room the illustration allocates a distinct store and copies the elements. The demonstration's growth policy is identical in both languages.

Always-visible growth label: “Capacity growth here is an illustration; Go may choose a different capacity. No question will ask you to predict the exact growth capacity.”

Length-three preset feedback: “These are three real elements whose values are zero. Append adds after them; it does not fill them.”

Reserved-three feedback: “There is room for three elements, but the current element count is zero.”

Empty-literal feedback: “An empty start with no reserved room. The first append needs storage that can hold it.”

After Append: “A new element at the end. Length increased by one; check whether storage changed or was reused.”

After Reset: “Back to your chosen starting preset. Only this experiment's additions were cleared.”

### EN `c61.mistakes` — markdown

**Assumption:** “make with length 3 reserves empty places; I will use append to fill them.” → **What Go does:** creates three zero-valued elements, and appends after them. → **Why:** length counts existing elements; capacity counts available room. → **How to think about a fix:** decide whether you already have elements to write by index, or want a log that starts at length zero and grows.

**Assumption:** “append increases the variable's length without my keeping its result.” → **What Go does:** returns an updated descriptor; a standalone append call that does not use its result is rejected. → **Why:** its returned value must be used. → **How to think about a fix:** keep the returned descriptor in the same variable or another suitable variable.

### EN `c61.bounds-trace` — walkthrough

An element's index must be at least zero and less than `len`. Extra room does not create extra elements. Watch why the second print never completes here.

Code and states: `e61-bounds`, unchanged. Frame captions: 0 “Before we start.”; 1 “There are no elements, despite the capacity.”; 2 “Length is zero.”; 3 “The requested index is zero.”; 4 “Execution stopped here; after is never printed. Start fixing this by understanding length, rather than hiding the panic.” The terminal panic text is the same real Go message in both languages.

### EN `c61.ticks` — markdown, optional

**Estimate for a separate exercise:** a training corridor has four markers. Start at its first marker, read each marker once, and move only between markers. One `Measure` costs one tick, and one `Move` costs one tick: four reads plus three moves gives **7 ticks**. Appending, arithmetic, and ordinary printing do not add game-action ticks; `Report` and `WorldSize` also cost no tick in these missions.

A move outside the corridor still costs a tick even when it fails. This estimate assumes no extra commands and that the program stops after its correct report. Display the last run's scored result as “Measured: last run”; it is neither necessarily the preview's world nor the saved historical best.

### EN `c61.q-make` — question

Prompt: “What will this print? Predict before revealing the result.” Code: `e61-make-question`, unchanged.

Options: A `[9] 1`; B `[0 0 9] 3`; C `[9 0] 2`. Correct: B.

Feedback A: “make started with two real elements. Append does not erase them. Check the starting length in the playground.”

Feedback B: “You tracked length correctly: two existing zeros, then a new element with value 9, making three elements.”

Feedback C: “That would require writing an existing index. Append adds at the end, rather than choosing the first zero.”

### EN `c61.q-append-result` — question

Prompt: “A learner forgot to keep append's result. Which description fits this code?” Code: `e61-unused-append`, unchanged.

Options: A “It prints [6 2].”; B “It simply prints [6].”; C “It does not compile: append's result is unused.” Correct: C.

Feedback A: “Append has a result you must use; this call does not automatically update s's length.”

Feedback B: “Go does not merely ignore the result here; it rejects the statement during compilation.”

Feedback C: “Exactly: this is a compile error, before any printing. Think about where the returned descriptor should go.”

### EN `c61.q-index` — question

Prompt: “A slice contains four elements. Which index selects the last one?”

Options: A `4`; B `3`; C `-1`. Correct: B.

Feedback A: “Four is the number of elements, rather than the last element's index. Index 4 is out of bounds.”

Feedback B: “Yes: the indices are 0, 1, 2, and 3. Keep the count distinct from the index.”

Feedback C: “Go does not use -1 as shorthand for the last element; it is an invalid index.”

### EN `c61.turn` — markdown

Return to the editor and run the starter unchanged. Look at `sizes`' length in the Console: the first observation is already stored. Write yourself a comment identifying the first observation still missing before adding commands. Return to code opens the editor for you.

### English hints

1. **Concept — `c61.h1`:** “The log stores ordered observations, rather than their sum. Remember the difference between printing the log and sending it to base.”
2. **Structure — `c61.h2`:** “Before repeating a step, check what the starter has already recorded. Give each observation exactly one place in the log.”
3. **Specific — `c61.h3`:** “Check the append statement: are you keeping its result? Check the report value too: is it the log itself, rather than its length or its last reading?”

## Verification appendix — sources, fixtures, and real outputs

The following are author-only executable fixtures for sandbox actions and prose examples. They are not extra lesson blocks or suggested mission code. Every fenced Go fragment in this document is registered by a `go-example` marker. The repeated English and transpose views reference canonical snippets rather than creating unverified variants.

### Sandbox/prose fixtures

**`e61-presets`:** all three requested presets, before/after one append; print starting capacity, but not an implementation-dependent grown capacity.

<!-- go-example id=e61-presets imports=fmt outcome=ok -->
```go
a := make([]int, 3)
fmt.Println(a, len(a), cap(a))
a = append(a, 7)
fmt.Println(a, len(a))
b := make([]int, 0, 3)
fmt.Println(b, len(b), cap(b))
b = append(b, 7)
fmt.Println(b, len(b), cap(b))
c := []int{}
fmt.Println(c, len(c), cap(c), c == nil)
c = append(c, 7)
fmt.Println(c, len(c))
```

**`e61-many-appends`:** reach the six-element UI limit from each preset; assert reuse when there was room and adequate capacity after growth, without fixing Go's growth policy.

<!-- go-example id=e61-many-appends imports=fmt outcome=ok -->
```go
presets := [][]int{make([]int, 3), make([]int, 0, 3), {}}
for id, s := range presets {
    valid := true
    for len(s) < 6 {
        oldLen, oldCap := len(s), cap(s)
        s = append(s, 7)
        valid = valid && len(s) == oldLen+1 && cap(s) >= len(s)
        if oldLen < oldCap { valid = valid && cap(s) == oldCap }
    }
    fmt.Println(id, s, len(s), valid)
}
```

**`e62-sharing`:** the two sharing presets and Clone, including append before mutation. All operations reset between scenarios.

<!-- go-example id=e62-sharing imports=fmt,slices outcome=ok -->
```go
a := make([]int, 2, 3)
a[0], a[1] = 5, 1
b := a
a = append(a, 7)
b[0] = 9
fmt.Println("room", a, b, len(a), len(b))
a = []int{5, 1}
b = a
a = append(a, 7)
b[0] = 9
fmt.Println("full", a, b, len(a), len(b))
a = []int{5, 1}
b = slices.Clone(a)
b[0] = 9
fmt.Println("clone", a, b)
```

**`e62-subslice`:** verify the optional end-exclusive shared view.

<!-- go-example id=e62-subslice imports=fmt outcome=ok -->
```go
a := []int{5, 1}
part := a[0:1]
part[0] = 9
fmt.Println(a, part, len(part))
```

**`e62-repeat-index`:** verify the fresh question dataset.

<!-- go-example id=e62-repeat-index imports=fmt,slices outcome=ok -->
```go
book := []int{2, 6, 2}
fmt.Println(slices.Index(book, 2), slices.Index(book, 2))
```

**`e63-sandbox`:** resource increments, a stored zero, missing read without insertion, and repeated deletion. Total is independently recomputed from values in this fixture.

<!-- go-example id=e63-sandbox imports=fmt outcome=ok -->
```go
m := map[string]int{}
m["water"]++
m["water"]++
m["air"] = 0
total := 0
for _, n := range m { total += n }
fmt.Println(m, len(m), total)
n, ok := m["air"]
fmt.Println("air", n, ok)
n, ok = m["spare"]
fmt.Println("spare", n, ok, len(m))
m["air"]++
total = 0
for _, n := range m { total += n }
fmt.Println(m, len(m), total)
delete(m, "air")
delete(m, "air")
fmt.Println(m, len(m))
```

**`e63-range-orders`:** observe actual iteration, without asserting that the order must change or predicting the next iteration. The output transcript is an observation from this verification session only.

<!-- go-example id=e63-range-orders imports=fmt outcome=ok -->
```go
stock := map[string]int{"water": 2, "air": 5}
fmt.Println(stock)
for turn := 0; turn < 4; turn++ {
    line := ""
    for key, value := range stock {
        line += fmt.Sprintf("%s=%d;", key, value)
    }
    fmt.Println(line)
}
```

### Verification method and transcripts

Executed all **28** registered fragments/fixtures outside the worktree in `/private/tmp/orbital-ch6-examples-3torsswg`. Final verification: **23 successful runs, 3 intentional panics, 2 intentional compile errors**, all matching their expectations. No game runs or repository test suites were invoked for this content task.

Compiler: `go version go1.25.7 darwin/arm64`, with matching `GOROOT=/usr/local/go`. Each example has its own temporary module declaring `go 1.24.0` and an explicit `fmt`/`slices` import wrapper. Build command: `/usr/local/go/bin/go build -trimpath -gcflags=-lang=go1.24 -o prog .`; execution uses that built binary, not `go run`. Cache, GOPATH, and TMPDIR are inside the same external temporary root. Network resolution is disabled (GOPROXY/GOSUMDB off, GOTOOLCHAIN local).

**Compatibility limit:** this machine has no installed Go 1.24 compiler. The examples were actually run with Go 1.25.7 in Go 1.24 language mode, not with a claimed 1.24 runtime. Used slices APIs (Sort, SortFunc, Clone, Index, Max) were also checked against the local toolchain's `api/go1.21.txt` exports; no post-1.24 API is used. Unspecified growth and iteration order are not inferred from one runtime observation.

Displayed fragment line numbers differ from compiler wrapper line numbers. Successful examples have empty stderr and build exit 0. For deterministic examples, the following stdout is both the exact expected value and the real observed output, including the final newline. Panic message first lines and pre-panic stdout matched; full real stderr is pasted. Exact stack addresses/line offsets and compiler diagnostic formatting are observations, not portable lesson guarantees. Empty streams are explicitly marked, not fabricated.

**Map-order expectation:** `e63-range-orders` must first print the sorted map, then four legal key/value permutations. Its transcript below is the exact observed output, not a claim that another run will repeat it. This run happened to repeat the same iteration order four times, which is legal.

**Sandbox limit:** fixtures verify Go's element/length/sharing/map semantics; they do not execute a future JavaScript implementation or certify its UI. The illustration's chosen capacity policy remains explicitly labelled and needs implementation tests later.

#### Verified `e61-basic` — ok

Build exit 0; program exit 0.

Real stdout:

```text
0 0
[7] 1
7 2
```

Real stderr: empty.

#### Verified `e61-bounds` — panic

Build exit 0; program exit 2.

Real stdout:

```text
0 3
```

Real stderr (intentional panic):

```text
panic: runtime error: index out of range [0] with length 0

goroutine 1 [running]:
main.main()
	storyboard/e61-bounds/main.go:9 +0x74
```

#### Verified `e61-make-question` — ok

Build exit 0; program exit 0.

Real stdout:

```text
[0 0 9] 3
```

Real stderr: empty.

#### Verified `e61-unused-append` — compile-error

Build exit 1. Intentional compile error; binary not run. Build stdout: empty.

Real compiler stderr:

```text
# storyboard/e61-unused-append
./main.go:7:1: append(s, 2) (value of type []int) is not used
```

#### Verified `e62-range-copy` — ok

Build exit 0; program exit 0.

Real stdout:

```text
0 6
1 16
[3 8]
```

Real stderr: empty.

#### Verified `e62-range-bounds` — ok

Build exit 0; program exit 0.

Real stdout:

```text
0 2
1 6
[2 6 9]
```

Real stderr: empty.

#### Verified `e62-sort` — ok

Build exit 0; program exit 0.

Real stdout:

```text
[6 11 14] [6 11 14]
```

Real stderr: empty.

#### Verified `e62-max-index` — ok

Build exit 0; program exit 0.

Real stdout:

```text
9
0
-1
```

Real stderr: empty.

#### Verified `e62-empty-max` — panic

Build exit 0; program exit 2.

Real stdout: empty.

Real stderr (intentional panic):

```text
panic: slices.Max: empty list

goroutine 1 [running]:
slices.Max[...](...)
	slices/sort.go:97
main.main()
	storyboard/e62-empty-max/main.go:8 +0x30
```

#### Verified `e63-presence` — ok

Build exit 0; program exit 0.

Real stdout:

```text
0 1
0 true 2
map[air:0]
```

Real stderr: empty.

#### Verified `e63-nil-map` — panic

Build exit 0; program exit 2.

Real stdout:

```text
0 false 0
```

Real stderr (intentional panic):

```text
panic: assignment to entry in nil map

goroutine 1 [running]:
main.main()
	storyboard/e63-nil-map/main.go:9 +0xc4
```

#### Verified `e63-print` — ok

Build exit 0; program exit 0.

Real stdout:

```text
map[air:5 water:2]
```

Real stderr: empty.

#### Verified `e64-comparator` — ok

Build exit 0; program exit 0.

Real stdout:

```text
-1 1 0
[3 5 8]
```

Real stderr: empty.

#### Verified `e64-identity` — ok

Build exit 0; program exit 0.

Real stdout:

```text
2 tape
[ink paper tape]
```

Real stderr: empty.

#### Verified `e64-bool-error` — compile-error

Build exit 1. Intentional compile error; binary not run. Build stdout: empty.

Real compiler stderr:

```text
# storyboard/e64-bool-error
./main.go:9:12: cannot use a < b (untyped bool value) as int value in return statement
```

#### Verified `e65-grid-index` — ok

Build exit 0; program exit 0.

Real stdout:

```text
2 3
8
```

Real stderr: empty.

#### Verified `e65-shallow` — ok

Build exit 0; program exit 0.

Real stdout:

```text
[[2 9] [5 7]]
[[2 9] [5 7]]
```

Real stderr: empty.

#### Verified `e65-rows` — ok

Build exit 0; program exit 0.

Real stdout:

```text
[[7 0] [0 0]]
```

Real stderr: empty.

#### Verified `e65-overlay-cell` — ok

Build exit 0; program exit 0.

Real stdout:

```text
4
```

Real stderr: empty.

#### Verified `e65-dot-cell` — ok

Build exit 0; program exit 0.

Real stdout:

```text
-1
```

Real stderr: empty.

#### Verified `e65-shared-rows` — ok

Build exit 0; program exit 0.

Real stdout:

```text
[[9 6] [9 6]]
```

Real stderr: empty.

#### Verified `e61-presets` — ok

Build exit 0; program exit 0.

Real stdout:

```text
[0 0 0] 3 3
[0 0 0 7] 4
[] 0 3
[7] 1 3
[] 0 0 false
[7] 1
```

Real stderr: empty.

#### Verified `e61-many-appends` — ok

Build exit 0; program exit 0.

Real stdout:

```text
0 [0 0 0 7 7 7] 6 true
1 [7 7 7 7 7 7] 6 true
2 [7 7 7 7 7 7] 6 true
```

Real stderr: empty.

#### Verified `e62-sharing` — ok

Build exit 0; program exit 0.

Real stdout:

```text
room [9 1 7] [9 1] 3 2
full [5 1 7] [9 1] 3 2
clone [5 1] [9 1]
```

Real stderr: empty.

#### Verified `e62-subslice` — ok

Build exit 0; program exit 0.

Real stdout:

```text
[9 1] [9] 1
```

Real stderr: empty.

#### Verified `e62-repeat-index` — ok

Build exit 0; program exit 0.

Real stdout:

```text
0 0
```

Real stderr: empty.

#### Verified `e63-sandbox` — ok

Build exit 0; program exit 0.

Real stdout:

```text
map[air:0 water:2] 2 2
air 0 true
spare 0 false 2
map[air:1 water:2] 2 3
map[water:2] 1
```

Real stderr: empty.

#### Verified `e63-range-orders` — ok

Build exit 0; program exit 0.

Real stdout:

```text
map[air:5 water:2]
water=2;air=5;
water=2;air=5;
water=2;air=5;
water=2;air=5;
```

Real stderr: empty.

