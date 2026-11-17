# Chapter 6 lesson storyboards — review draft

Short bilingual lesson content for branch `lessons-ch6`. Mission setups, starters, checks, thresholds, and APIs stay unchanged. 6-1 revision 2 is approved; 6-2 and 6-3 now use the same short six-step format.

## Style rules

- Short lessons target about 250–320 words per language, counting titles, frame captions, question options/feedback, sandbox labels, and hints; exclude code and repeated application controls.
- Mission: three numbered points, at most about 45 words; name the gold target and what the starter already did.
- One short analogy, then show code and let the player manipulate data. One idea per visual. No paragraph should exceed three lines at 290px; split sentences into short paragraphs where needed. Frame captions: at most about eight words.
- Mission → Idea → Try it → Watch out → Think & check → Your turn. At most two one-line warnings, a short prose plan, three mission questions, one next action. Wrong picks reveal only their feedback until the correct pick or two wrong picks on that question.
- Egyptian Arabic first, short sentences, English technical terms. English has the same meaning and depth. Label invented values as illustrative. No free-form interpreter, complete mission loop, or last-move condition supplied as code.
- 6-2 and 6-3 below follow these rules. The older 6-4 and 6-5 storyboards remain unchanged.

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
| 6-1 revision 2 | `slice.basics`, `slice.append`, `slice.bounds`, `ticks.actions`; len only, no capacity lesson |
| 6-2 | `range.slice`, `slice.sharing`, `slice.sort`, `slice.search-max`, `ticks.route-order` |
| 6-3 | `map.counts`, `map.presence`, `map.nil`, `map.order` |
| 6-4 | `sort.comparator`, `sort.identity-ties` |
| 6-5 | `grid.shape`, `grid.row-ownership`, `grid.operations`, `lab.contract` |

`slice.sharing` includes subslices and shallow copying; 6-5 applies that model to row slices. `ticks.actions` is explained once in 6-1; later tick notes apply it rather than redefine it. No struct/pointer syntax, generics syntax, concurrency, or new dependencies are prerequisites.

## 6-1 — سجل الكريستال (`6-1-survey`) · revision 2

**Word count:** Arabic 275; English 291, including titles, frame captions, every option/feedback, playground labels, and hints; excluding code and shared application controls.

**Objective:** ordered observations, append's returned value, len, zero-based indexes, and the mission's action count. Capacity, nil, storage sharing, make presets, and compile-error questions have been removed from this lesson.

**Starter audit:** already records the first crystal and prints the partial slice; it neither moves nor reports. Starter, setup, checks, Requires, Par=Gold=15, and MaxTicks=60 are unchanged. Revision-1 practice falls back to the new Mission step and does not restore retired answers.

### Step `c61.v2.mission-step` — المهمة

`c61.v2.mission` · markdown

1. قيس الـ8 كريستالات من الشمال لليمين.
2. احفظ أحجامهم بالترتيب في slice واحدة.
3. ابعت الـslice كلها مرة واحدة بـ~colony.Report~.

الـstarter قاس أول كريستالة وحفظها فعلًا.

الذهب: **15 ticks أو أقل**.

### Step `c61.v2.idea-step` — الفكرة

`c61.v2.trace` · walkthrough

تخيل الـslice صندوق، وكل حجم جديد بيتحط في آخره بـ~append~.

<!-- go-example id=e61-crystal-log-v2 imports=fmt outcome=ok -->
```go
sizes := []int{}
sizes = append(sizes, 5)
sizes = append(sizes, 2)
sizes = append(sizes, 8)
fmt.Println(sizes, len(sizes), sizes[0])
```

| Frame / line | Slice (indexed from 0) | Output | Caption |
|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. |
| 1 / 1 | [] | — | سجل فاضي؛ len=0. |
| 2 / 2 | [5] | — | 5 عند index 0. |
| 3 / 3 | [5, 2] | — | 2 بعدها، مش بدلها. |
| 4 / 4 | [5, 2, 8] | — | 8 بقت آخر عنصر. |
| 5 / 5 | [5, 2, 8] | [5 2 8] 3 5 | السجل، طوله، وأول حجم. |

Frames show `sizes`, `len`, and index cells only; no capacity, nil, or store labels. Frame 5 is terminal.

### Step `c61.v2.try-step` — جرّب

`c61.v2.sequence` · sandbox · `slice-append-sequence`

ضيف الأحجام بالترتيب؛ راقب ~len~ والـindex وهما بيتغيروا.

Source row, indexes 0–7: `[5, 2, 8, 4, 1, 7, 3, 6]`. Label: «أحجام للتوضيح، مش بيانات المهمة.».
Fixed operations: «ضيف الحجم الجاي», «من الأول». Initially empty, len=0, counter `0 / 8`. Each Add appends the next source size; only the newly added cell appears. Live len and counter advance together. At 8, disable Add and show «اكتمل السجل: 8 أحجام بالترتيب.». Reset restores the empty slice, counter, and enabled Add; source values stay visible. This sandbox does not move the drone or edit code.

### Step `c61.v2.watch-step` — خد بالك

`c61.v2.watch` · markdown

- احتفظ بنتيجة ~append~: ~sizes = append(...)~.
- أرقام الـindex من 0 لحد ~len(sizes)-1~؛ براها يحصل panic.

### Step `c61.v2.check-step` — فكّر واختبر

`c61.v2.plan` · markdown

التكرار محتاج loop.

لكل كريستالة: قيس، ضيف، واتجه للي بعدها.

التقرير مرة واحدة بعد الـloop؛ أول قراءة موجودة فعلًا.

`c61.v2.q-create` · question

عشان السجل يتجمع، تنشئ الـslice فين؟

1. قبل الـloop → سجل واحد قبل الـloop يحتفظ بكل القراءات. **Correct.**
2. جوه كل لفة → إنشاءها كل لفة يبدأ سجل فاضي تاني.
3. بعد الـloop → السجل لازم يبقى موجود قبل تجميع القراءات.

`c61.v2.q-append` · question

أنهي أداة تضيف حجم من غير مسح الموجود؟

1. len → len بتعد الموجود؛ مش بتضيف قراءة.
2. append → append بتضيف في الآخر؛ احتفظ بالناتج. **Correct.**
3. make → make بتنشىء slice، مش بتضيف للسجل الحالي.

`c61.v2.q-moves` · question

واقف عند أول كريستالة؛ كام حركة توصّلك للتامنة؟

1. 8 → 8 حركات تحاول تعدّي آخر كريستالة.
2. 6 → 6 حركات تسيبك قبل الأخيرة بمسافة.
3. 7 → 8 قياسات + 7 حركات = 15 ticks؛ Report=0 وبافتراض مفيش تكلفة تانية. **Correct.**

### Step `c61.v2.turn-step` — دورك

`c61.v2.turn` · markdown

اكتب أول محاولة بنفسك، شغّلها، وراجع السجل اللي طلع.

Existing Return-to-code action follows this block.

### Three position-labelled hints

1. **Concept:** الـslice سجل بيكبر مع append، مع الاحتفاظ بالأحجام اللي اتسجلت قبله.
2. **Structure:** رتّب loop للكريستالات وابعت تقرير واحد بعده؛ افتكر أول قراءة في الـstarter.
3. **Specific:** راجع سطر ~append~ اللي فيه ~colony.Measure()~ وشرط الحركة: بتسجل مرة وتتحرك لكريستالة تالية؟

**Tick assumptions:** Measure and Move cost one action tick each; Report costs zero. The 7-move question assumes eight total measurements and no other tick costs. Fifteen is the gold target, not the run limit.

**Solution boundary:** this revision deliberately supplies the prose plan and movement-count question requested after playtesting. The player still writes iteration bounds, reconciles the existing first measurement, chooses movement conditions, and places the report; no complete mission loop or move guard is supplied.

**Fallback:** same short mission, illustrative analogy/trace, two warnings, tick assumptions, prose plan, and next action; no capacity/nil material.

## 6-2 — أغلى 3 (`6-2-top-three`) · revision 1

**Source check:** the mission is a 10-crystal row. The setup, starter, and check all use 10; this follows `ch06.go` rather than the brief's “8”. The starter measured each crystal into its matching index, moved east nine times, and leaves the drone at x=9. Sizes are distinct. The drone may harvest three. Gold and Par change with the selected positions; the run limit is 200. No mission code or grading value changed.

**Teaching objective:** read an index/value pair from `range`; see that sorting a copied slice keeps the original position order; reason about choosing before harvesting. Existing fallback markdown does not contradict this lesson and remains unchanged.

**Word count:** Arabic 287; English 312, including titles, captions, options, feedback, sandbox labels, and hints; code is excluded.

### المهمة / Mission · `c62.v1.mission-step`

**AR — `c62.v1.mission`:**

1. الـstarter قاس العشرة وسجل الحجم حسب index؛ الدرون عند x=9.
2. اختار أكبر 3؛ حمولة الدرون 3، وبعدها احصدهم.
3. هدف الذهب بيتغير حسب المواقع؛ شوف المعروض.

**EN:**

1. The starter measured 10 crystals and saved each size by index; drone at x=9.
2. Pick the largest 3; carry limit is 3, then harvest.
3. Gold varies by location; see the shown target.

### الفكرة / Idea · `c62.v1.idea-step`

`c62.v1.trace` · walkthrough · `e62-range-clone-sort`

**AR:** في slice تدريب، range بتدي index وقيمة؛ Clone تنسخ، وSort ترتب.

**EN:** In a sample slice, range gives index and value; Clone copies, Sort orders.

<!-- go-example id=e62-range-clone-sort imports=fmt,slices outcome=ok -->
```go
sizes := []int{4, 1, 3}
for i, size := range sizes {
	fmt.Println(i, size)
}
sorted := slices.Clone(sizes)
slices.Sort(sorted)
fmt.Println(sizes)
fmt.Println(sorted)
```

| Frame / line | State | Cumulative output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 1 | sizes=[4 1 3] | — | أحجام للتوضيح؛ كل index ثابت. | Illustrative sizes; fixed indexes. |
| 2 / 2 | i=0,size=4 | — | range أعطت أول index وقيمة. | First range index and value. |
| 3 / 3 | unchanged | `0 4` | طبعناهم مع بعض. | Printed pair. |
| 4 / 2 | i=1,size=1 | `0 4` | index التالي وقيمته. | Next index/value. |
| 5 / 3 | unchanged | `0 4\n1 1` | آخر عنصر في range. | Last range item. |
| 6 / 2 | i=2,size=3 | `0 4\n1 1` | طبعنا آخر زوج. | Printed last pair. |
| 7 / 3 | unchanged | `0 4\n1 1\n2 3` | اللفة خلصت. | Loop finished. |
| 8 / 5 | sorted=[4 1 3]; sizes unchanged | previous | Clone عملت slice جديدة. | Clone copied the slice. |
| 9 / 6 | sorted=[1 3 4]; sizes unchanged | previous | Sort رتبت النسخة بس. | Sort changed the copy. |
| 10 / 7 | sizes=[4 1 3] | previous + `[4 1 3]` | ترتيب الأصل فضل زي ما هو. | Original order stayed. |
| 11 / 8 | sorted=[1 3 4] | `0 4\n1 1\n2 3\n[4 1 3]\n[1 3 4]` | النسخة المرتبة اتطبعت. | Printed the copy. |

`range` returns each slice index and a copy of its element value. `slices.Clone` copies the integer elements to a separate slice; `slices.Sort` changes that copy in ascending order. The walkthrough shows one loop iteration at a time; it does not describe Go's internal sort calls.

### جرّب / Try it · `c62.v1.try-step`

`c62.v1.sort-playground` · sandbox · `slice-sort-clone`

**AR:** رتّب الأصل أو نسخة؛ راقب الأحجام والـindex.

**EN:** Sort original or copy; watch sizes and indexes.

Illustrative source sizes at indexes 0–7: `[5, 2, 8, 4, 1, 7, 3, 6]`. Source label: AR «أحجام للتوضيح، مش بيانات المهمة.» / EN “Illustrative sizes, not mission data.” Buttons: AR «رتّب نسخة (Clone ثم Sort)»، «رتّب الأصل»، «من الأول» / EN “Sort copy (Clone then Sort)”, “Sort original”, “Reset”. Clone+Sort keeps the source row and shows an ascending copy below. Sort original reorders the source values while the index labels remain 0–7; status says those indexes no longer identify the source crystals. Reset restores the initial row and hides the copy.

### خد بالك / Watch out · `c62.v1.watch-step`

**AR:**

- `b := a` بيشارك العناصر؛ Sort على b يغيّر a.
- Sort تغيّر القيم، مش أماكن الكريستالات.
- Measure وMove وHarvest: tick لكل أمر؛ Report والحساب: صفر.

**EN:**

- `b := a` shares elements; sorting b changes a.
- Sort changes values, not crystal positions.
- Measure, Move, Harvest: 1 tick each; Report and computation: 0.

### فكّر واختبر / Think & check · `c62.v1.check-step`

`c62.v1.plan`: «القياسات جاهزة. اختار الثلاثة، وخطط الزيارة.» / “Measurements are ready. Choose three and plan visits.”

`c62.v1.q-index`: في `[6, 1, 9, 3]`، أنهي index فيه أكبر قيمة؟ / In `[6, 1, 9, 3]`, which index holds the largest value?

- 0 → index 0 فيه 6، مش الأكبر. / Index 0 holds 6.
- 2 → صح، القيمة 9 عند index 2. / Correct: 9 is at index 2. **Correct.**
- 3 → index 3 فيه 3. / Index 3 holds 3.

`c62.v1.q-clone`: ليه Clone قبل Sort لو المواقع مهمة؟ / Why clone before sorting when positions matter?

- عشان ترتيب الأحجام الأصلي يفضل محفوظ. / To preserve the original size order. **Correct.**
- عشان ننقل الكريستالات لمواقع جديدة. / To move crystals to new locations.
- عشان نزود سعة الـslice. / To increase the slice capacity.

Feedback: «صح؛ النسخة تسيب ترتيب الأصل زي ما هو.» / “Correct; a copy preserves the original order.” «Sort تغيّر القيم، مش العالم.» / “Sort changes values, not the world.” «Clone مش بتضمن سعة أكبر.» / “Clone does not guarantee more capacity.”

`c62.v1.q-scan`: إمتى تختار أكبر 3؟ / When choose the largest 3?

- قبل ما تعرف الأحجام. / Before you know the sizes.
- بعد ما تعرف أحجام العشرة. / After you know all ten sizes. **Correct.**
- بعد أول 3 قياسات. / After the first three measurements.

Feedback: «قارن العشرة عشان تعرف الأكبر.» / “Compare all ten sizes to know the largest.” «صح؛ قارن العشرة قبل القرار.» / “Correct: compare all ten before deciding.” «أول 3 مش بالضرورة الأكبر.» / “The first three may not be largest.”

### دورك / Your turn · `c62.v1.turn-step`

`c62.v1.turn`: شغّل الـstarter؛ اكتب خطتك بنفسك. / Run the starter; write your plan.

**Three hints (Concept / Structure / Specific):**

1. «الحجم قيمة، والـindex مكان.» / “Size is value; index is position.”
2. «قارن الأحجام، ثم خطط الزيارة والحصاد.» / “Compare all sizes, then plan visits.”
3. «راجع `slices.Clone` و`slices.Sort`؛ أي ترتيب يفضل مربوط بالمواقع؟» / “Review `slices.Clone` and `slices.Sort`; which order must stay linked to positions?”

**Ticks:** the starter scan uses ten `Measure` calls and nine `Move` calls; each costs one tick. Three `Harvest` calls also cost one tick each. Route cost varies with the chosen locations, so the displayed Gold target is dynamic. `Report` and in-memory comparisons cost zero.

#### Verified `e62-range-clone-sort`

Toolchain: Go 1.24.0; `imports=fmt,slices`; `outcome=ok`.

Real stdout:

```text
0 4
1 1
2 3
[4 1 3]
[1 3 4]
```

## 6-3 — تعداد الحياة (`6-3-census`) · revision 1

**Source check:** the field is 6×4. The starter initializes `map[colony.Entity]int`, reads and counts the starting tile, and prints the map. The check accepts counts with or without `empty`. Gold is 47 ticks, Par is 60, and the run limit is 300. No mission code or grading value changed.

**Teaching objective:** count repeated keys; distinguish a missing key from a stored zero with comma-ok; use `len` for key count; avoid nil-map writes and assumptions about iteration order. Existing fallback markdown does not contradict this lesson and remains unchanged.

**Word count:** Arabic 283; English 302, including titles, captions, options, feedback, sandbox labels, and hints; code is excluded.

### المهمة / Mission · `c63.v1.mission-step`

**AR — `c63.v1.mission`:**

1. عدّ moss وlichen وcrystal في خريطة 6×4 وابعت map الأعداد.
2. عدّ empty اختياري؛ الـstarter سجّل خانة البداية بالفعل.
3. الذهب 47 ticks؛ والـpar هو 60.

**EN:**

1. Count moss, lichen, and crystal on the 6×4 map; report the counts map.
2. Counting empty is optional; the starter already counted its starting tile.
3. Gold is 47 ticks; Par is 60.

### الفكرة / Idea · `c63.v1.idea-step`

`c63.v1.trace` · walkthrough · `e63-map-count`

**AR:** الـmap بتربط key بعدد؛ القراية الغايبة ترجع صفر، وcomma-ok توضّح وجود المفتاح.

**EN:** A map connects keys to counts; a missing read returns zero, and comma-ok shows whether the key exists.

<!-- go-example id=e63-map-count imports=fmt outcome=ok -->
```go
counts := map[string]int{}
counts["moss"]++
counts["moss"]++
counts["lichen"]++
fmt.Println(counts["crystal"])
n, ok := counts["crystal"]
fmt.Println(n, ok, len(counts))
```

| Frame / line | Map summary (displayed alphabetically) | Cumulative output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التنفيذ. | Before execution. |
| 1 / 1 | `{}`; len=0 | — | map فاضية؛ len=0. | Empty map; len=0. |
| 2 / 2 | `{moss: 1}`; len=1 | — | أول عدّ لـmoss بقى 1. | First moss count is 1. |
| 3 / 3 | `{moss: 2}`; len=1 | — | نفس المفتاح؛ len زي ما هي. | Same key; len stays 1. |
| 4 / 4 | `{lichen: 1, moss: 2}`; len=2 | — | مفتاح جديد؛ len بقت 2. | New key; len is 2. |
| 5 / 5 | unchanged | `0` | قراية المفتاح الغايب رجعت صفر. | Missing key reads as zero. |
| 6 / 6 | `n=0`, `ok=false`; map unchanged | `0` | ok بتقول إن المفتاح غايب. | ok says the key is absent. |
| 7 / 7 | len=2 | `0\n0 false 2` | الخرج بيعرض مفتاحين. | Output includes two keys. |

Map summaries in the walkthrough are ordered alphabetically for display. That display choice does not promise a `range` order.

### جرّب / Try it · `c63.v1.try-step`

`c63.v1.map-playground` · sandbox · `map-count`

**AR:** اضغط نوعًا؛ عدّاده يزيد. جرّب قراءة مفتاح مش موجود.

**EN:** Click a kind to increment it; try reading an absent key.

Illustrative chips: `moss`, `lichen`, `crystal`; not mission data. Each click increments that key and updates the key/count table and `len`. “Read absent key: spore” shows `0, false` without inserting `spore`. Reset clears the table and read result. The table is sorted for display only; map iteration order is unspecified and may vary.

### خد بالك / Watch out · `c63.v1.watch-step`

**AR:**

- الكتابة على map nil تعمل panic؛ جهّزها الأول، وrange ترتيبها مش مضمون.
- GetEntity وMove وHarvest كل واحد 1 tick، حتى Move الفاشلة؛ Report وWorldSize والحساب 0.

**EN:**

- A nil map write panics; initialize first; map range order is unspecified.
- GetEntity, Move, Harvest: 1 tick each, even failed Move; Report, WorldSize, computation: 0.

### فكّر واختبر / Think & check · `c63.v1.check-step`

`c63.v1.plan`: «جهّز map، وزُر الخانات وعدّ الأنواع، وبعدها ابعت التقرير.» / “Prepare a map, visit and count the tiles, then report.”

`c63.v1.q-missing`: لو المفتاح `crystal` مش موجود، إيه نتيجة comma-ok؟ / If `crystal` is absent, what is comma-ok?

- `0, false` → صح؛ القيمة صفر وok=false لأن المفتاح غايب. / Correct: a missing read is 0, false. **Correct.**
- `0, true` → ok=true معناها إن المفتاح موجود. / ok=true means the key exists.
- `nil, false` → قيمة int الافتراضية هي صفر، مش nil. / int's zero value is 0.

`c63.v1.q-len`: لو moss=4 وlichen=2، كام `len(counts)`؟ / If moss=4 and lichen=2, what is `len(counts)`?

- 6 → 6 هو مجموع القيم؛ len مش بتجمعها. / 6 is the sum, not len.
- 2 → صح؛ len بتعد المفتاحين. / Correct: len counts two keys. **Correct.**
- 4 → 4 قيمة moss، مش عدد المفاتيح. / 4 is a value, not a key count.

`c63.v1.q-ticks`: بافتراض GetEntity مرة لكل 24 خانة و23 حركة بلا تكلفة أخرى، كام tick؟ / One GetEntity per 24 tiles, 23 moves, no other costs: ticks?

- 47 → صح؛ 24 قراءة + 23 حركة = 47، وReport والحساب مجانيان. / Correct: 24 + 23 = 47; Report and computation are free. **Correct.**
- 48 → في زيارة 24 خانة، فيه 23 حركة بس. / A 24-tile visit needs 23 moves.
- 60 → 60 هو par، مش تقدير المسار المذكور. / 60 is Par, not this estimate.

### دورك / Your turn · `c63.v1.turn-step`

`c63.v1.turn`: راجع طباعة الـstarter وحدد نوع خانة البداية اللي اتعدّت. / Check the starter output and identify the starting tile's counted kind.

**Three hints (Concept / Structure / Specific):**

1. «الـmap بتعد كل نوع؛ `len` بتعد المفاتيح.» / “A map counts each kind; len counts keys.”
2. «ابدأ بعدّ الـstarter، وزُر باقي الخانات مرة.» / “Start from the starter’s count; visit other tiles once.”
3. «راجع الـmap وسطر `colony.GetEntity()`؛ هل المفتاح متوافق؟» / “Check the map and `colony.GetEntity()` line: do key types match?”

#### Verified `e63-map-count`

Toolchain: Go 1.24.0; `imports=fmt`; `outcome=ok`.

Real stdout:

```text
0
0 false 2
```

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

## English learner copy appendix — 6-1 · revision 2

### Step `c61.v2.mission-step` — Mission

`c61.v2.mission` · markdown

1. Measure the 8 crystals from left to right.
2. Keep their sizes in order in one slice.
3. Report the whole slice once with ~colony.Report~.

The starter already saved the first measurement.

Gold: **15 ticks or fewer**.

### Step `c61.v2.idea-step` — Idea

`c61.v2.trace` · walkthrough

A slice is a box: ~append~ adds each size at its end.

Code: canonical `e61-crystal-log-v2` above.

| Frame / line | Slice (indexed from 0) | Output | Caption |
|---|---|---|---|
| 0 / before 1 | — | — | Before execution. |
| 1 / 1 | [] | — | Empty slice; len=0. |
| 2 / 2 | [5] | — | 5 at index 0. |
| 3 / 3 | [5, 2] | — | 2 follows 5. |
| 4 / 4 | [5, 2, 8] | — | 8 comes last. |
| 5 / 5 | [5, 2, 8] | [5 2 8] 3 5 | Log, length, first size. |

Frames show `sizes`, `len`, and index cells only; no capacity, nil, or store labels. Frame 5 is terminal.

### Step `c61.v2.try-step` — Try it

`c61.v2.sequence` · sandbox · `slice-append-sequence`

Add sizes in order; watch ~len~ grow and indexes advance.

Source row, indexes 0–7: `[5, 2, 8, 4, 1, 7, 3, 6]`. Label: «Illustrative sizes, not mission data.».
Fixed operations: «Add next size», «Reset». Initially empty, len=0, counter `0 / 8`. Each Add appends the next source size; only the newly added cell appears. Live len and counter advance together. At 8, disable Add and show «Complete: 8 sizes in order.». Reset restores the empty slice, counter, and enabled Add; source values stay visible. This sandbox does not move the drone or edit code.

### Step `c61.v2.watch-step` — Watch out

`c61.v2.watch` · markdown

- Keep append's result: ~sizes = append(...)~.
- Indexes: 0 through ~len(sizes)-1~; outside means panic.

### Step `c61.v2.check-step` — Think & check

`c61.v2.plan` · markdown

Use a loop for repeated work.

For each crystal: measure, append, move toward the next.

Report once afterward; the starter saved the first reading.

`c61.v2.q-create` · question

Where should this mission's slice be created?

1. Before the loop → One shared log keeps all readings. **Correct.**
2. Inside each round → Each round would restart the log.
3. After the loop → The log must exist before collecting readings.

`c61.v2.q-append` · question

What adds a size without erasing earlier ones?

1. len → len counts; it adds nothing.
2. append → append extends the log; save its result. **Correct.**
3. make → make creates a slice, rather than appending.

`c61.v2.q-moves` · question

Starting on crystal 1, how many moves reach crystal 8?

1. 8 → Eight moves go past the last crystal.
2. 6 → Six moves stop one crystal too soon.
3. 7 → 8 measures + 7 moves = 15 ticks; Report=0, assuming no other tick costs. **Correct.**

### Step `c61.v2.turn-step` — Your turn

`c61.v2.turn` · markdown

Write your first attempt, then run it.

Existing Return-to-code action follows this block.

### Three position-labelled hints

1. **Concept:** append grows the slice while keeping existing sizes.
2. **Structure:** Loop over the crystals; report once afterward, remembering the first saved reading.
3. **Specific:** Check the ~append~ / ~colony.Measure()~ line and move guard: save once; move toward another crystal.

## Verification appendix — sources, fixtures, and real outputs

The following are author-only executable fixtures for sandbox actions and prose examples. They are not extra lesson blocks or suggested mission code. Every fenced Go fragment in this document is registered by a `go-example` marker. The repeated English and transpose views reference canonical snippets rather than creating unverified variants.

### Sandbox/prose fixtures

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

Historical verification of revision 1 executed **28** fragments/fixtures (six retired 6-1 fixtures/transcripts are removed in revision 2) outside the worktree in `/private/tmp/orbital-ch6-examples-3torsswg`. Final verification: **23 successful runs, 3 intentional panics, 2 intentional compile errors**, all matching their expectations. No game runs or repository test suites were invoked for this content task.

Compiler: `go version go1.25.7 darwin/arm64`, with matching `GOROOT=/usr/local/go`. Each example has its own temporary module declaring `go 1.24.0` and an explicit `fmt`/`slices` import wrapper. Build command: `/usr/local/go/bin/go build -trimpath -gcflags=-lang=go1.24 -o prog .`; execution uses that built binary, not `go run`. Cache, GOPATH, and TMPDIR are inside the same external temporary root. Network resolution is disabled (GOPROXY/GOSUMDB off, GOTOOLCHAIN local).

**Compatibility limit:** this machine has no installed Go 1.24 compiler. The examples were actually run with Go 1.25.7 in Go 1.24 language mode, not with a claimed 1.24 runtime. Used slices APIs (Sort, SortFunc, Clone, Index, Max) were also checked against the local toolchain's `api/go1.21.txt` exports; no post-1.24 API is used. Unspecified growth and iteration order are not inferred from one runtime observation.

Displayed fragment line numbers differ from compiler wrapper line numbers. Successful examples have empty stderr and build exit 0. For deterministic examples, the following stdout is both the exact expected value and the real observed output, including the final newline. Panic message first lines and pre-panic stdout matched; full real stderr is pasted. Exact stack addresses/line offsets and compiler diagnostic formatting are observations, not portable lesson guarantees. Empty streams are explicitly marked, not fabricated.

**Map-order expectation:** `e63-range-orders` must first print the sorted map, then four legal key/value permutations. Its transcript below is the exact observed output, not a claim that another run will repeat it. This run happened to repeat the same iteration order four times, which is legal.

**Sandbox limit:** fixtures verify Go's element/length/sharing/map semantics; they do not execute a future JavaScript implementation or certify its UI. The illustration's chosen capacity policy remains explicitly labelled and needs implementation tests later.

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


### Revision-2 sequence fixture and verification

The next author-only fixture verifies the eight Add operations; it is not extra learner code. New chapter-six lesson examples are verified by `TestChapterSixStoryboardExamplesMatchRealGo` in temporary directories.

<!-- go-example id=e61-sequence-v2 imports=fmt outcome=ok -->
```go
sizes := []int{}
for _, size := range []int{5, 2, 8, 4, 1, 7, 3, 6} {
    sizes = append(sizes, size)
    fmt.Println(sizes, len(sizes))
}
```

Compiler for revision 2: `go version go1.27.1 darwin/arm64`; module language `go 1.24.0`, built with `-gcflags=-lang=go1.24`, network dependency resolution disabled. Both examples actually ran outside the worktree under `/private/tmp/orbital-ch61-v2/examples`.

#### Verified `e61-crystal-log-v2` — ok

Build exit 0; program exit 0.

Real stdout:

```text
[5 2 8] 3 5
```

Real stderr: empty.

#### Verified `e61-sequence-v2` — ok

Build exit 0; program exit 0.

Real stdout:

```text
[5] 1
[5 2] 2
[5 2 8] 3
[5 2 8 4] 4
[5 2 8 4 1] 5
[5 2 8 4 1 7] 6
[5 2 8 4 1 7 3] 7
[5 2 8 4 1 7 3 6] 8
```

Real stderr: empty.

