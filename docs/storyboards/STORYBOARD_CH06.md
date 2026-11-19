# Chapter 6 lesson storyboards — review draft

Short bilingual lesson content for chapter 6. Mission setups, starters, checks, thresholds, and APIs stay unchanged. 6-1 revision 2 is approved; 6-2 to 6-5 use the same short six-step format.

## Style rules

- Short lessons target about 250–320 words per language, counting titles, frame captions, question options/feedback, sandbox labels, and hints; exclude code and repeated application controls.
- Mission: three numbered points, at most about 45 words; name the gold target and what the starter already did.
- One short analogy, then show code and let the player manipulate data. One idea per visual. No paragraph should exceed three lines at 290px; split sentences into short paragraphs where needed. Frame captions: at most about eight words.
- Mission → Idea → Try it → Watch out → Think & check → Your turn. At most two one-line warnings, a short prose plan, three mission questions, one next action. Wrong picks reveal only their feedback until the correct pick or two wrong picks on that question.
- Egyptian Arabic first, short sentences, English technical terms. English has the same meaning and depth. Label invented values as illustrative. No free-form interpreter, complete mission loop, or last-move condition supplied as code.
- All chapter 6 sections below follow these rules.

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

## 6-4 — ترتيب المخزون (`6-4-stock-order`) · revision 1

**Word count:** Arabic 271; English 316, including titles, captions, questions and feedback, sandbox controls, and hints; excluding Go code. The tests count screen-reader card labels too.

### المهمة / Mission

**AR:** 1. رتّب المواقع: الحجم الأكبر أولًا، ثم x الأصغر عند التعادل.
2. اقرأ العرض بـ~colony.WorldSize~ وابعث المواقع بتقرير واحد.
3. الـstarter طبع قياسًا من غير حفظ؛ كلف tick. ممنوع الحصاد.

الذهب: **19 ticks أو أقل**.

**EN:** 1. Rank positions: larger sizes first, then smaller x on ties.
2. Read width with ~colony.WorldSize~; report positions once.
3. The starter printed an unsaved measurement (one tick); do not harvest.

Gold: **19 ticks or fewer**.

### الفكرة / Idea — walkthrough

**AR:** المقارنة بتقول مين يسبق: السالب يقدّم a، والصفر تعادل، والموجب يأخّره.

المثال أسماء ثم أعمار؛ قاعدة المهمة مختلفة.

**EN:** A comparator returns a sign: negative puts a first, zero ties, positive puts it later.

This toy sorts names, then ages; the mission differs.

<!-- go-example id=e64-sort-sign-late imports=cmp,fmt,slices outcome=ok -->
```go
type Player struct { name string; age int }
players := []Player{{"Mira", 8}, {"Ari", 12}, {"Ari", 9}}
before := func(a, b Player) int {
	if byName := cmp.Compare(a.name, b.name); byName != 0 {
		return byName
	}
	return cmp.Compare(a.age, b.age)
}
fmt.Println(before(players[1], players[0]))
fmt.Println(before(players[1], players[2]))
slices.SortFunc(players, before)
fmt.Println(players)
```

| Frame | AR caption | EN caption |
|---|---|---|
| 0 | قبل الكروت. | Before the cards. |
| 1 | نوع بيانات للتجربة. | An illustrative type. |
| 2 | أسماء وأعمار للتوضيح. | Illustrative names and ages. |
| 3 | قاعدة المقارنة محفوظة. | Comparator saved. |
| 9 | Ari يسبق Mira؛ النتيجة -1. | Ari precedes Mira; result -1. |
| 10 | الاسم اتعادل؛ العمر حسمها. | Names tie; age decides. |
| 11 | SortFunc رتّبت القائمة. | SortFunc reorders the slice. |
| 12 | الأسماء ثم الأعمار. | Names, then ages. |

#### Verified e64-sort-sign-late

Real stdout:

```text
-1
1
[{Ari 9} {Ari 12} {Mira 8}]
```

### جرّب / Try it

**AR:** بدّل اتجاه الحجم وكسر التعادل؛ الكروت توضيحية والتعادلات مميزة.
**EN:** Toggle size direction and tie-breaking; cards are illustrative and ties are highlighted.

**Sandbox labels:** illustrative sizes, size direction, ascending/descending, tie-break on/off, reset. Each card has an accessible position and size label. Sample sizes by position: 3, 5, 3, 5, 2.

### خد بالك / Watch out

**AR:** خلي نتيجة المقارنة متسقة: سالب أو صفر أو موجب. ~slices.SortFunc~ مش stable؛ ~slices.SortStableFunc~ بتحافظ على ترتيب العناصر المتعادلة.

~Measure~ و~Move~ (حتى لو فشلت) كل واحدة tick.

~WorldSize~ و~Report~ والحساب بلا ticks. الحصاد tick وممنوع هنا.

**EN:** Keep comparator results consistent: negative, zero, or positive. ~slices.SortFunc~ is not stable; ~slices.SortStableFunc~ keeps equal items in input order.

~Measure~, ~Move~ (even when it fails), and ~Harvest~ each cost one tick.

~WorldSize~, ~Report~, and computation cost none.

### فكّر واختبر / Think & check

**Plan AR:** قيس الكل → احتفظ بالقراءات → رتّب المواقع → ابعت تقريرًا واحدًا.
**Plan EN:** Measure all → keep the readings → rank positions → send one report.

1. **AR:** -1 معناها إيه؟ · a يسبق / تعادل / a يتأخر. **الإجابة:** a يسبق.
   **Feedback:** السالب يقدّم a. / الصفر يعني تعادل. / الموجب يأخر a.
   **EN:** What does -1 mean? · a comes first / Tie / a comes later. **Answer:** a comes first.
   **Feedback:** Negative puts a first. / Zero means a tie. / Positive puts a later.

2. **AR:** هل ~SortFunc~ تحفظ ترتيب المتعادلين؟ · دايمًا / مش مضمون / بس مع أرقام. **الإجابة:** مش مضمون.
   **Feedback:** SortFunc مش stable. / صح؛ SortStableFunc تحفظ ترتيبهم. / نوع العنصر مش فارق هنا.
   **EN:** Does ~SortFunc~ keep tied items in input order? · Always / Not guaranteed / Only for numbers. **Answer:** Not guaranteed.
   **Feedback:** SortFunc is not stable. / Right; SortStableFunc keeps their order. / The element type does not matter here.

3. **AR:** تقرير المهمة يبعت إيه؟ · أحجام / مواقع / طباعة. **الإجابة:** مواقع.
   **Feedback:** المهمة عايزة المواقع. / صح؛ ابعت المواقع المرتبة. / الطباعة للتصحيح.
   **EN:** What does the mission report? · Sizes / Positions / Printed output. **Answer:** Positions.
   **Feedback:** The mission needs positions. / Right; report ranked positions. / Printing is for debugging.

### دورك / Your turn

**AR:** افتح الـstarter، عدّل حاجة صغيرة، وشغّل المهمة وراجع التقرير.
**EN:** Open the starter; make a small change and inspect its report.

**Hints AR:** فكرة: المقارنة دالة قرار؛ نتيجتها تحدد مين يسبق. · بنية: افصل قراءات الحجم عن هويات المواقع قبل اختيار ما ستفرزه. · محدد: راجع معنى نتيجة comparator في ~slices.SortFunc~: السالب والصفر والموجب.
**Hints EN:** Concept: a comparator is a decision rule; its result says which item comes first. · Structure: separate measured sizes from position identities before choosing what to sort. · Specific: review comparator results in ~slices.SortFunc~: negative, zero, and positive.

## 6-5 — شبكة الحرارة (`6-5-heat-grid`) · revision 1

**Word count:** Arabic 311; English 320, including titles, captions, questions and feedback, sandbox controls, and hints; excluding Go code. The tests count screen-reader cell labels too.

### المهمة / Mission

**AR:** دي lab: زر Run بيشغّل tests للدوال، مش حركة للدرون. كمّل ~OverlayHeat~ و~ThroughHeat~ و~FlipHeat~؛ مفيش ~Report~ ولا ticks.
**EN:** Lab: Run tests, not drones. Implement ~OverlayHeat~, ~ThroughHeat~, and ~FlipHeat~; no ~Report~ or ticks.

### الفكرة / Idea — walkthrough

**AR:** الـ2D slice صفوف جوه slice. في ~grid[r][c]~، الأول row والتاني column؛ ~len(grid)~ للصفوف و~len(grid[0])~ للأعمدة.
**EN:** 2D slices hold rows. In ~grid[r][c]~, row comes first, then column; ~len(grid)~ counts rows and ~len(grid[0])~ counts columns.

<!-- go-example id=e65-grid-copy-late imports=fmt,slices outcome=ok -->
```go
grid := make([][]int, 2)
grid[0] = make([]int, 3)
grid[1] = make([]int, 3)
grid[0][1] = 4
grid[1][2] = 5
outer := slices.Clone(grid)
outer[0][1] = 9
rows := [][]int{slices.Clone(grid[0]), slices.Clone(grid[1])}
rows[0][1] = 7
fmt.Println(len(grid), len(grid[0]), grid[1][2])
fmt.Println(grid[0][1], outer[0][1], rows[0][1])
```

| Line | AR caption | EN caption |
|---|---|---|
| 0 | قبل الشبكة. | Before. |
| 1 | اتعمل مكان لصفين. | Two rows. |
| 2 | الصف صفر بقى 3 خانات. | Row 0: three cells. |
| 3 | الصف الأول بقى 3 خانات. | Row 1: three cells. |
| 4 | [0][1] بقت 4. | Cell [0][1] becomes 4. |
| 5 | [1][2] بقت 5. | Cell [1][2] becomes 5. |
| 6 | نسخنا outer بس. | Outer slice cloned. |
| 7 | الشبكتين شافوا 9. | Both show 9. |
| 8 | نسخنا كل صف لوحده. | Clone each row. |
| 9 | النسخة بقت 7 لوحدها. | Copy row now has 7. |
| 10 | عدد الصفوف والأعمدة، ثم خلية. | Rows, columns, cell. |
| 11 | الخرج يوضح الصف المشترك والمنسوخ. | Shared 9; copy 7. |

#### Verified e65-grid-copy-late

Real stdout:

```text
2 3 5
9 9 7
```

### جرّب / Try it

**AR:** اختار خلية وشوف ~grid[r][c]~؛ بدّل النسخ وعدّل النسخة.
**EN:** Select a cell for ~grid[r][c]~; change clone mode and edit the copy.

**Sandbox labels:** illustrative grid, clone outer only, clone each row, edit copy, reset, Original, Copy. The selected row and column headers highlight r and c. Sample grid: [3 7 2] / [8 6 4].

### خد بالك / Watch out

**AR:** ~grid[r][c]~: row ثم column. ~slices.Clone(grid)~ سطحي؛ الصفوف مشتركة.
**EN:** ~grid[r][c]~: row, then column. ~slices.Clone(grid)~ is shallow; rows stay shared.

### فكّر واختبر / Think & check

**Plan AR:** الشكل → العملية → نتيجة مستقلة أو error.
**Plan EN:** Shape → operation → independent result or error.

**Contract AR:** OverlayHeat بتجمع خانات لشبكتين نفس الشكل.

كل خانة في ThroughHeat مجموع ضرب صف في عمود؛ أعمدة A = صفوف B.

FlipHeat بتبدّل المحاور.

المدخل الفاضي أو المتعرج أو أبعاده غلط يرجّع error؛ وممنوع تعدّل المدخلات.

صفوف الناتج مستقلة؛ حتى 1×1 والصف أو العمود الواحد صالحين.
**Contract EN:** OverlayHeat adds matching cells.

Each ThroughHeat cell sums row-by-column products; columns(A)=rows(B).

FlipHeat swaps axes.

Empty, ragged, or incompatible input returns an error; leave inputs unchanged.

Output rows are independent; 1×1 and row or column vectors are valid.

1. **AR:** في شبكة جرّب، ~grid[1][0]~ كام؟ · 7 / 8 / 6. **الإجابة:** 8.
   **Feedback:** 7 في [0][1]. / صح؛ row 1 وcolumn 0. / 6 في [1][2].
   **EN:** In the playground grid, what is ~grid[1][0]~? · 7 / 8 / 6. **Answer:** 8.
   **Feedback:** 7 is at [0][1]. / Right; row 1, column 0. / 6 is at [1][2].

2. **AR:** A 2×3 وB 3×4؛ ThroughHeat صالحة؟ · آه؛ 2×4 / لأ / آه؛ 3×3. **الإجابة:** آه؛ 2×4.
   **Feedback:** صح؛ 3 هو البعد المشترك. / ThroughHeat محتاجة أعمدة A تساوي صفوف B. / 3 هو البعد المشترك، مش الناتج.
   **EN:** A is 2×3 and B is 3×4; is ThroughHeat valid? · Yes; 2×4 / No / Yes; 3×3. **Answer:** Yes; 2×4.
   **Feedback:** Right; 3 is the shared dimension. / ThroughHeat needs A's columns to equal B's rows. / 3 is shared, not the output shape.

3. **AR:** نسخت outer وعدّلت خلية؛ إيه اللي يحصل؟ · المصدر يتغير / النسخة بس / Compile error. **الإجابة:** المصدر يتغير.
   **Feedback:** الصفوف لسه مشتركة. / ده محتاج نسخ كل row. / Go تسمح بالنسخ السطحي.
   **EN:** Clone the outer slice and edit a cell; what happens? · Source changes / Copy only / Compile error. **Answer:** Source changes.
   **Feedback:** Rows are still shared. / That needs each row cloned. / Go allows shallow clones.

### دورك / Your turn

**AR:** شغّل الـstarter، وحدد أول contract فشل، واختار دالة.
**EN:** Run once; note first failed contract, then choose a function.

**Hints AR:** فكرة: ارسم مصدر خلية ناتجة. · بنية: افصل الشكل عن أبعاد العملية. · محدد: ~slices.Clone(grid)~ يترك الصفوف مشتركة.
**Hints EN:** Concept: trace one output cell. · Structure: separate shape and dimension checks. · Specific: ~slices.Clone(grid)~ leaves rows shared.


## English learner copy — 6-1 revision 2

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

The next author-only fixture verifies the eight Add operations; it is not extra learner code. New 6-1 fixtures are verified by `TestSurveyStoryboardExamplesMatchRealGo` in temporary directories outside the worktree.

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
