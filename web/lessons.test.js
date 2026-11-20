"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const { TestElement } = require("./sandbox/harness.js");
const source = fs.readFileSync(__dirname + "/lessons.js", "utf8");
const app = fs.readFileSync(__dirname + "/app.js", "utf8");

// The DOM and persistence are boundaries; all lesson behavior runs unchanged.
function fixture({ blob, review = false, lang = "en" } = {}) {
  let focused = null;
  const observers = [];
  const createElement = (tag) => new TestElement(tag, { onFocus: (focusedElement) => { focused = focusedElement; } });
  const roots = Object.fromEntries(["#lesson", "#lesson-panel", "#btn-lesson", "#btn-hint", "#hint-list"].map((id) => [id, createElement("section")]));
  roots["#btn-hint"].appendChild(createElement("span"));
  const editor = createElement("textarea");
  const body = createElement("body");
  const store = { blob: blob || { progress: {}, prefs: {}, quizScores: {}, review: {}, lessons: {} }, save() {},
    pref(k, v) { if (v === undefined) return this.blob.prefs[k]; this.blob.prefs[k] = v; } };
  const context = vm.createContext({
    ResizeObserver: class {
      constructor(callback) { this.callback = callback; observers.push(this); }
      observe(target) { this.target = target; }
      disconnect() { this.disconnected = true; }
    },
    document: { body, createElement, querySelector: (selector) => roots["#lesson"].querySelector(selector) },
    LearningMetrics: { recordHint() {}, recordStepViewed() {}, recordLessonOpened() {}, recordInteraction() {},
      normalizeMetrics: (metrics) => metrics || {}, mergeMetrics: (a, b) => b || a || {} },
    $: (selector) => roots[selector], ed: editor, store, S: { review }, T: (ar, en) => lang === "en" ? en : ar, Sound: { play() {} },
  });
  vm.runInContext(app.slice(app.indexOf("const esc ="), app.indexOf("// ---------- files ----------")), context);
  vm.runInContext(source + "\nthis.lessonUI = Lessons;", context);
  return { ui: context.lessonUI, article: roots["#lesson"], context, roots, editor, store, body, lang, observers, get focused() { return focused; } };
}

function lesson() {
  return { id: "survey", lessonRevision: 1, lesson: "# Old lesson", steps: [
    { id: "mission", title: "Mission", blocks: [{ id: "intro", kind: "markdown", text: "Read the mission." }] },
    { id: "practice", title: "Practice", blocks: [
      { id: "bounds", kind: "walkthrough", text: "Predict before Next.", imports: ["fmt"], code: "s := make([]int, 0, 3)\nfmt.Println(len(s), cap(s))\ni := 0\nfmt.Println(s[i])", walkthrough: { frames: [
        { line: 0, caption: "Before we start.", output: "" },
        { line: 1, caption: "No elements.", output: "", variables: [{ name: "s", value: "[]" }], slice: { name: "s", values: [], length: 0, capacity: 3, store: "A" } },
        { line: 2, caption: "Length is zero.", output: "0 3", slice: { name: "s", values: [], length: 0, capacity: 3, store: "A" } },
        { line: 4, caption: "Execution stopped.", output: "0 3", panic: "index out of range [0] with length 0" },
      ] } },
      { id: "question", kind: "question", text: "Pick the last index.", question: { options: ["4", "3", "-1"], answer: 1, feedback: ["Four is out of bounds.", "Indices start at zero.", "Negative indices are invalid."] } },
      { id: "playground", kind: "sandbox", text: "Choose a preset.", sandbox: {
        kind: "slice-make", appendValue: 7, limit: 6, growthLabel: "Example growth; Go may choose differently.", reservedLabel: "Reserved room, not elements.", limitLabel: "Six-element experiment limit.", appendText: "Added at the end.", resetText: "Restored the chosen preset.",
        presets: [
          { id: "length", label: "make([]int, 3)", values: [0, 0, 0], capacity: 3, text: "Three existing elements." },
          { id: "capacity", label: "make([]int, 0, 3)", values: [], capacity: 3, text: "Reserved room only." },
          { id: "empty", label: "[]int{}", values: [], capacity: 0, text: "Empty slice." },
        ],
      } },
    ] },
  ] };
}

function click(node) { assert.ok(node, "control exists"); assert.ok(!node.disabled, "control is enabled"); node.onclick(); }
function namedButton(root, label) { return root.querySelectorAll("button").find((b) => b.textContent === label); }
function practice(view, level = lesson()) { view.ui.open(level); click(namedButton(view.article.querySelector("nav"), view.lang === "ar" ? "التالي" : "Next")); }
function summaryParts(root) { return root.querySelector("figcaption").querySelectorAll("bdi").map((part) => part.textContent); }

function sequenceLesson(lang = "en") {
  return { id: "survey", lessonRevision: 2, steps: [
    { id: "c61.v2.mission-step", title: "Mission", blocks: [{ id: "c61.v2.mission", kind: "markdown", text: "The starter saved the first reading." }] },
    { id: "c61.v2.try-step", title: "Try it", blocks: [{ id: "c61.v2.sequence", kind: "sandbox", text: "Watch len grow.", sandbox: {
      kind: "slice-append-sequence", sequence: [5, 2, 8, 4, 1, 7, 3, 6],
      sourceLabel: lang === "ar" ? "أحجام للتوضيح، مش بيانات المهمة." : "Illustrative sizes, not mission data.",
      appendLabel: lang === "ar" ? "ضيف الحجم الجاي" : "Add next size", resetLabel: lang === "ar" ? "من الأول" : "Reset",
      doneText: lang === "ar" ? "اكتمل السجل: 8 أحجام بالترتيب." : "Complete: 8 sizes in order.",
    } }] },
  ] };
}

function sortCloneLesson(lang = "en") {
  return { id: "top-three", lessonRevision: 1, steps: [
    { id: "mission", title: "Mission", blocks: [{ id: "mission-copy", kind: "markdown", text: "The starter measured the crystals." }] },
    { id: "try", title: "Try it", blocks: [{ id: "sort-copy", kind: "sandbox", text: "Compare sorting the original and a copy.", sandbox: {
      kind: "slice-sort-clone", sequence: [5, 2, 8, 4, 1, 7, 3, 6],
      sourceLabel: lang === "ar" ? "أحجام للتوضيح، مش بيانات المهمة." : "Illustrative sizes, not mission data.",
      cloneSortLabel: lang === "ar" ? "رتّب نسخة (Clone ثم Sort)" : "Sort copy (Clone then Sort)",
      originalSortLabel: lang === "ar" ? "رتّب الأصل" : "Sort original", resetLabel: lang === "ar" ? "من الأول" : "Reset",
      cloneSortText: lang === "ar" ? "النسخة المرتبة تحت؛ الأصل فضل زي ما هو." : "Sorted copy below; the original stays unchanged.",
      originalSortText: lang === "ar" ? "الترتيب اتغير؛ الـindex مبقاش يحدد الكريستالة." : "Order changed; indexes no longer identify crystals.",
      resetText: lang === "ar" ? "رجعنا لترتيب البداية." : "Original order restored.",
    } }] },
  ] };
}

function mapCountLesson(lang = "en") {
  return { id: "census", lessonRevision: 1, steps: [
    { id: "mission", title: "Mission", blocks: [{ id: "mission-map", kind: "markdown", text: "Count the illustrative kinds." }] },
    { id: "try", title: "Try it", blocks: [{ id: "map-count", kind: "sandbox", text: "Click kinds to count them.", sandbox: {
      kind: "map-count", kinds: ["moss", "lichen", "crystal"], missingKey: "spore",
      sourceLabel: lang === "ar" ? "أنواع للتوضيح، مش بيانات المهمة." : "Illustrative kinds, not mission data.",
      missingReadLabel: lang === "ar" ? "اقرأ مفتاح غايب: spore" : "Read absent key: spore",
      resetLabel: lang === "ar" ? "من الأول" : "Reset",
      mapOrderLabel: lang === "ar" ? "الجدول للعرض بس؛ ترتيب range على map مش مضمون." : "Sorted for display; map range order is unspecified.",
    } }] },
  ] };
}

test("sequence playground grows one indexed cell per click, stops at eight, and resets in both languages", () => {
  for (const lang of ["ar", "en"]) {
    const view = fixture({ lang }), level = sequenceLesson(lang);
    const spec = level.steps[1].blocks[0].sandbox;
    practice(view, level);
    const box = view.article.querySelector(".lesson-sequence"), row = box.querySelector(".lesson-crystals");
    const add = namedButton(box, spec.appendLabel);
    assert.equal(row.querySelectorAll("li").length, 8);
    assert.deepEqual(row.querySelectorAll("b").map((cell) => cell.textContent), spec.sequence.map(String));
    assert.deepEqual(row.querySelectorAll(".lesson-index").map((cell) => cell.textContent), ["0", "1", "2", "3", "4", "5", "6", "7"]);
    assert.match(row.querySelectorAll("li")[0].attributes["aria-label"], /index 0/);
    assert.equal(box.querySelector(".lesson-position").textContent, "0 / 8");
    assert.equal(box.querySelector(".lesson-caption").hidden, true);
    assert.equal(box.querySelectorAll(".lesson-cell").length, 0);
    assert.deepEqual(summaryParts(box.querySelector(".lesson-slice")), ["sizes:", "len=0"]);
    assert.doesNotMatch(box.textContent, /cap=|Store|nil|مخزن/);
    for (let n = 1; n <= 8; n++) {
      click(add);
      const slice = box.querySelector(".lesson-slice");
      assert.deepEqual(slice.querySelectorAll("b").map((cell) => cell.textContent), spec.sequence.slice(0, n).map(String));
      assert.deepEqual(slice.querySelectorAll(".lesson-index").map((cell) => cell.textContent), Array.from({ length: n }, (_, i) => String(i)));
      assert.deepEqual(summaryParts(slice), ["sizes:", `len=${n}`]);
      assert.equal(box.querySelector(".lesson-position").textContent, `${n} / 8`);
    }
    assert.equal(add.disabled, true);
    assert.equal(box.querySelector(".lesson-caption").hidden, false);
    assert.equal(box.querySelector(".lesson-caption").textContent, spec.doneText);
    click(namedButton(box, spec.resetLabel));
    assert.equal(add.disabled, false);
    assert.equal(box.querySelector(".lesson-caption").hidden, true);
    assert.equal(box.querySelector(".lesson-position").textContent, "0 / 8");
    assert.equal(box.querySelectorAll(".lesson-cell").length, 0);
    assert.deepEqual(spec.sequence, [5, 2, 8, 4, 1, 7, 3, 6]);
    click(add);
    assert.deepEqual(summaryParts(box.querySelector(".lesson-slice")), ["sizes:", "len=1"]);
  }
});

test("sort-copy playground preserves the original or visibly loses its position links, then resets", () => {
  for (const lang of ["ar", "en"]) {
    const view = fixture({ lang }), level = sortCloneLesson(lang);
    const spec = level.steps[1].blocks[0].sandbox;
    practice(view, level);
    const box = view.article.querySelector(".lesson-sort-clone");
    const source = box.querySelector(".lesson-sequence-source");
    const sourceRow = source.querySelector(".lesson-crystals");
    const copyFigure = box.querySelector(".lesson-sort-copy");
    const cloneSort = namedButton(box, spec.cloneSortLabel);
    const originalSort = namedButton(box, spec.originalSortLabel);
    const reset = namedButton(box, spec.resetLabel);
    const initial = spec.sequence.map(String);
    assert.equal(sourceRow.querySelectorAll("li").length, 8);
    assert.deepEqual(sourceRow.querySelectorAll("b").map((cell) => cell.textContent), initial);
    assert.deepEqual(sourceRow.querySelectorAll(".lesson-index").map((cell) => cell.textContent), ["0", "1", "2", "3", "4", "5", "6", "7"]);
    assert.equal(copyFigure.hidden, true);
    assert.equal(box.querySelector(".lesson-caption").hidden, true);
    click(cloneSort);
    assert.deepEqual(sourceRow.querySelectorAll("b").map((cell) => cell.textContent), initial);
    assert.equal(copyFigure.hidden, false);
    assert.deepEqual(copyFigure.querySelectorAll("b").map((cell) => cell.textContent), [...initial].sort((a, b) => Number(a) - Number(b)));
    assert.match(box.querySelector(".lesson-caption").textContent, /unchanged|فضل زي ما هو/);
    click(originalSort);
    assert.deepEqual(sourceRow.querySelectorAll("b").map((cell) => cell.textContent), [...initial].sort((a, b) => Number(a) - Number(b)));
    assert.deepEqual(sourceRow.querySelectorAll(".lesson-index").map((cell) => cell.textContent), ["0", "1", "2", "3", "4", "5", "6", "7"]);
    assert.match(box.querySelector(".lesson-caption").textContent, /no longer identify|مبقاش يحدد/);
    click(reset);
    assert.deepEqual(sourceRow.querySelectorAll("b").map((cell) => cell.textContent), initial);
    assert.equal(copyFigure.hidden, true);
    assert.match(box.querySelector(".lesson-caption").textContent, /restored|رجعنا/);
  }
});

test("map-count playground increments chips, reads absence without insertion, and resets in both languages", () => {
  for (const lang of ["ar", "en"]) {
    const view = fixture({ lang }), level = mapCountLesson(lang);
    const spec = level.steps[1].blocks[0].sandbox;
    practice(view, level);
    const box = view.article.querySelector(".lesson-map-count");
    const rows = () => box.querySelector(".lesson-map-table").querySelector("tbody").querySelectorAll("tr");
    const summary = () => box.querySelector(".lesson-position").textContent;
    assert.match(box.textContent, /range.*مش مضمون|range order is unspecified/);
    assert.equal(summary(), "counts · len=0");
    assert.equal(box.querySelector(".lesson-map-empty").hidden, false);
    assert.equal(rows().length, 0);
    click(namedButton(box, spec.missingReadLabel));
    assert.equal(box.querySelector(".lesson-map-readout").textContent, "spore: 0, false");
    assert.equal(summary(), "counts · len=0");
    click(namedButton(box, "moss")); click(namedButton(box, "moss")); click(namedButton(box, "lichen"));
    assert.equal(summary(), "counts · len=2");
    assert.equal(box.querySelector(".lesson-map-empty").hidden, true);
    assert.deepEqual(rows().map((row) => row.querySelectorAll("td").map((cell) => cell.textContent)), [["lichen", "1"], ["moss", "2"]]);
    click(namedButton(box, spec.missingReadLabel));
    assert.equal(box.querySelector(".lesson-map-readout").textContent, "spore: 0, false");
    assert.equal(summary(), "counts · len=2");
    click(namedButton(box, spec.resetLabel));
    assert.equal(summary(), "counts · len=0");
    assert.equal(rows().length, 0);
    assert.equal(box.querySelector(".lesson-map-empty").hidden, false);
    assert.equal(box.querySelector(".lesson-map-readout").hidden, true);
  }
});

test("revision-1 practice opens the new Mission and cannot restore retired answers or wrong-pick counts", () => {
  const blob = { lessons: { survey: { revision: 1, step: "c61.check-step", answers: { "c61.q-index": 1 }, wrongAttempts: { "c61.q-index": 2 } } } };
  const view = fixture({ blob }), level = sequenceLesson(); view.ui.open(level);
  assert.match(view.article.textContent, /Step 1 of 2/);
  assert.match(view.article.textContent, /starter saved the first reading/);
  assert.equal(view.article.querySelector(".lesson-sequence"), null);
  click(namedButton(view.article.querySelector("nav"), "Next"));
  assert.equal(view.store.blob.lessons.survey.revision, 2);
  assert.equal(view.store.blob.lessons.survey.step, "c61.v2.try-step");
  assert.equal(Object.keys(view.store.blob.lessons.survey.answers).length, 0);
  assert.equal(Object.keys(view.store.blob.lessons.survey.wrongAttempts).length, 0);
});

test("introductory walkthrough frames show only elements and length", () => {
  const level = lesson();
  level.steps[1].blocks[0].walkthrough.frames[1].slice = { name: "sizes", values: [5], length: 1, elementsOnly: true };
  const view = fixture(); practice(view, level);
  const walk = view.article.querySelector(".lesson-walkthrough");
  click(namedButton(walk, "Next line"));
  assert.deepEqual(summaryParts(walk), ["sizes:", "len=1"]);
  assert.equal(walk.querySelectorAll(".reserved").length, 0);
  assert.doesNotMatch(walk.querySelector("figcaption").textContent, /cap|Store|nil/);
});

test("walkthrough starts before execution, advances to terminal panic, goes back and resets", () => {
  const view = fixture(); practice(view);
  const walk = view.article.querySelector(".lesson-walkthrough");
  assert.match(walk.textContent, /Before execution/);
  assert.equal(walk.querySelectorAll(".current").length, 0);
  assert.equal(namedButton(walk, "Previous line").disabled, true);
  click(namedButton(walk, "Next line"));
  assert.equal(walk.querySelectorAll(".current")[0].dataset.line, "1");
  assert.equal(walk.querySelectorAll(".reserved").length, 3);
  assert.deepEqual(summaryParts(walk), ["s:", "len=0", "cap=3", "Store", "A"]);
  click(namedButton(walk, "Next line"));
  assert.equal(walk.querySelector(".lesson-output").textContent, "0 3");
  click(namedButton(walk, "Next line"));
  assert.match(walk.querySelector(".lesson-panic").textContent, /Panic · Execution stopped here/);
  assert.match(walk.textContent, /Current line: 4/);
  assert.equal(namedButton(walk, "Next line").disabled, true);
  click(namedButton(walk, "Previous line"));
  assert.equal(walk.querySelector(".lesson-panic"), null);
  click(namedButton(walk, "Restart"));
  assert.equal(walk.querySelectorAll(".current").length, 0);
  assert.match(walk.querySelector(".lesson-output").textContent, /No output/);
});

test("questions show choice-specific feedback and text labels, persist an answer, and allow retry", () => {
  const view = fixture(); practice(view);
  const question = view.article.querySelector("fieldset");
  assert.equal(question.querySelectorAll(".correct").length, 0);
  click(namedButton(question, "4"));
  assert.match(question.textContent, /Incorrect. Four is out of bounds./);
  assert.doesNotMatch(question.textContent, /Correct answer|Indices start at zero/);
  assert.equal(question.querySelectorAll(".correct").length, 0);
  assert.equal(question.querySelectorAll(".incorrect").length, 1);
  assert.match(question.textContent, /Your answer/);
  assert.equal(view.store.blob.lessons.survey.answers.question, 0);
  assert.equal(view.store.blob.lessons.survey.wrongAttempts.question, 1);
  assert.ok(question.querySelectorAll(".lesson-choice").every((b) => b.disabled));
  click(namedButton(question, "Try again"));
  assert.equal(question.querySelector(".lesson-feedback").textContent, "");
  assert.equal(view.store.blob.lessons.survey.answers.question, undefined);
  click(namedButton(question, "3"));
  assert.match(question.textContent, /Correct! Indices start at zero./);
  assert.match(question.textContent, /Correct answer/);
});

test("two wrong picks reveal only that question and keep the explanation after reload, language change and retry", () => {
  const level = lesson();
  const original = level.steps[1].blocks.find((block) => block.question);
  level.steps[1].blocks.push({ ...original, id: "other-question" });
  const first = fixture(); practice(first, level);
  const [question, other] = first.article.querySelectorAll("fieldset");
  click(namedButton(question, "4"));
  click(namedButton(question, "Try again"));
  click(namedButton(other, "-1"));
  assert.equal(first.article.querySelectorAll(".correct").length, 0);
  // The first wrong pick survives retry and a reload before the second attempt.
  const reload = fixture({ blob: first.store.blob, lang: "ar" }); reload.ui.open(level);
  const [resumed, resumedOther] = reload.article.querySelectorAll("fieldset");
  click(namedButton(resumed, "4"));
  assert.match(resumed.textContent, /الإجابة الصحيحة/);
  assert.match(resumed.querySelector(".lesson-answer-explanation").textContent, /Indices start at zero/);
  assert.equal(resumed.querySelectorAll(".incorrect").length, 1);
  assert.equal(resumedOther.querySelectorAll(".correct").length, 0);
  assert.equal(reload.store.blob.lessons.survey.wrongAttempts.question, 2);
  assert.equal(reload.store.blob.lessons.survey.wrongAttempts["other-question"], 1);
  click(namedButton(resumed, "جرّب تاني"));
  reload.ui.open(level);
  assert.match(reload.article.querySelector("fieldset").textContent, /الإجابة الصحيحة.*Indices start at zero/);
});

test("legacy numeric answers restore without revealing the correct answer after just one wrong pick", () => {
  for (const choice of [0, 1]) {
    const view = fixture({ blob: { lessons: { survey: { revision: 1, step: "practice", answers: { question: choice } } } } });
    view.ui.open(lesson());
    const question = view.article.querySelector("fieldset");
    assert.equal(question.querySelectorAll(".correct").length, choice === 1 ? 1 : 0);
    if (choice === 0) {
      click(namedButton(question, "Try again"));
      click(namedButton(question, "4"));
      assert.equal(question.querySelectorAll(".correct").length, 1);
      assert.match(question.textContent, /Indices start at zero/);
    }
  }
});

test("walkthrough controls are distinct, contextual and placed between the code and state in both languages", () => {
  for (const [lang, labels, prefix] of [["en", ["Previous line", "Next line", "Restart"], "Walkthrough:"], ["ar", ["السطر اللي فات", "السطر الجاي", "من الأول"], "شرح الكود:"]]) {
    const view = fixture({ lang }); practice(view);
    const walk = view.article.querySelector(".lesson-walkthrough");
    const controls = walk.querySelector(".lesson-controls");
    const position = walk.children.indexOf(walk.querySelector(".lesson-code"));
    assert.equal(walk.children[position + 1], controls);
    assert.ok(walk.children[position + 2].querySelector(".lesson-frame"));
    assert.equal(controls.attributes.role, "group");
    for (const label of labels) assert.equal(namedButton(controls, label).attributes["aria-label"], `${prefix} ${label}`);
    assert.equal(walk.querySelectorAll("button").length, 3);
  }
});

test("line changes scroll only the bounded code box, including backwards and restart", () => {
  const view = fixture(); practice(view);
  const walk = view.article.querySelector(".lesson-walkthrough"), code = walk.querySelector(".lesson-code");
  code.clientHeight = 40;
  code.querySelectorAll(".lesson-code-line").forEach((line, index) => { line.offsetTop = index * 20; line.offsetHeight = 20; });
  view.roots["#lesson-panel"].scrollTop = 240; view.body.scrollTop = 80;
  click(namedButton(walk, "Next line"));
  assert.equal(code.scrollTop, 0);
  code.scrollTop = 20;
  click(namedButton(walk, "Next line"));
  assert.equal(code.scrollTop, 20);
  click(namedButton(walk, "Next line"));
  assert.equal(code.scrollTop, 40);
  click(namedButton(walk, "Previous line"));
  assert.equal(code.scrollTop, 20);
  click(namedButton(walk, "Previous line"));
  assert.equal(code.scrollTop, 0);
  code.scrollTop = 20;
  click(namedButton(walk, "Restart"));
  assert.equal(code.scrollTop, 0);
  assert.equal(view.roots["#lesson-panel"].scrollTop, 240);
  assert.equal(view.body.scrollTop, 80);
});

test("sticky navigation offset follows its measured height and is cleaned up for legacy lessons", () => {
  const view = fixture(); view.ui.open(lesson());
  const nav = view.article.querySelector("nav"), panel = view.roots["#lesson-panel"];
  nav.offsetHeight = 164;
  view.observers[0].callback();
  assert.equal(panel.style.getPropertyValue("--lesson-nav-offset"), "180px");
  view.ui.open({ id: "old", lesson: "Old lesson" });
  assert.equal(view.observers[0].disconnected, true);
  assert.equal(panel.style.getPropertyValue("--lesson-nav-offset"), "");
});

test("slice summaries isolate Go runs and translated labels in walkthroughs and sandbox in both languages", () => {
  for (const [lang, unknown, store] of [["ar", "غير محددة", "مخزن"], ["en", "not specified", "Store"]]) {
    const level = lesson(); delete level.steps[1].blocks[0].walkthrough.frames[1].slice.capacity;
    const view = fixture({ lang }); practice(view, level);
    const walk = view.article.querySelector(".lesson-walkthrough");
    click(namedButton(walk, lang === "ar" ? "السطر الجاي" : "Next line"));
    const parts = walk.querySelector("figcaption").querySelectorAll("bdi");
    assert.deepEqual(parts.map((part) => [part.textContent, part.dir]), [["s:", "ltr"], ["len=0", "ltr"], ["cap=", "ltr"], [unknown, "auto"], [store, "auto"], ["A", "ltr"]]);
    const box = view.article.querySelector(".lesson-sandbox");
    click(namedButton(box, lang === "ar" ? "ضيف 7" : "Append 7"));
    assert.deepEqual(summaryParts(box), ["s:", "len=4", "cap=6", store, "B"]);
    assert.equal(box.querySelector("figcaption").querySelectorAll("bdi").find((part) => part.textContent === store).dir, "auto");
  }
});

test("sandbox presets append real elements, reuse room or allocate a new store, and reset at the limit", () => {
  const view = fixture(); practice(view);
  const box = view.article.querySelector(".lesson-sandbox");
  const caption = () => summaryParts(box);
  assert.deepEqual(caption(), ["s:", "len=3", "cap=3", "Store", "A"]);
  click(namedButton(box, "Append 7"));
  assert.deepEqual(caption(), ["s:", "len=4", "cap=6", "Store", "B"]);
  assert.deepEqual(box.querySelectorAll("b").map((n) => n.textContent), ["0", "0", "0", "7", "—", "—"]);
  click(namedButton(box, "make([]int, 0, 3)"));
  click(namedButton(box, "Append 7"));
  assert.deepEqual(caption(), ["s:", "len=1", "cap=3", "Store", "A"]);
  click(namedButton(box, "Append 7")); click(namedButton(box, "Append 7"));
  assert.deepEqual(caption(), ["s:", "len=3", "cap=3", "Store", "A"]);
  click(namedButton(box, "Append 7"));
  assert.deepEqual(caption(), ["s:", "len=4", "cap=6", "Store", "B"]);
  click(namedButton(box, "Append 7")); click(namedButton(box, "Append 7"));
  assert.equal(namedButton(box, "Append 7").disabled, true);
  assert.equal(box.querySelector(".lesson-limit").hidden, false);
  click(namedButton(box, "Reset"));
  assert.deepEqual(caption(), ["s:", "len=0", "cap=3", "Store", "A"]);
  click(namedButton(box, "[]int{}"));
  assert.deepEqual(caption(), ["s:", "len=0", "cap=0"]);
  click(namedButton(box, "Append 7"));
  assert.deepEqual(caption(), ["s:", "len=1", "cap=1", "Store", "A"]);
  assert.match(box.querySelector(".lesson-growth-label").textContent, /Go may choose differently/);
});

test("step and answers restore across language reloads; changed revisions and stale IDs start fresh", () => {
  const first = fixture(); practice(first);
  click(namedButton(first.article.querySelector("fieldset"), "3"));
  const reload = fixture({ blob: first.store.blob, lang: "ar" });
  reload.ui.open(lesson());
  assert.match(reload.article.textContent, /خطوة 2 من 2/);
  assert.match(reload.article.querySelector("fieldset").textContent, /صح!/);
  assert.ok(reload.article.querySelectorAll(".lesson-choice").every((b) => b.disabled));
  const changed = lesson(); changed.lessonRevision++;
  reload.ui.open(changed);
  assert.match(reload.article.textContent, /خطوة 1 من 2/);
  const stale = { revision: 1, step: "removed", answers: { question: 99, obsolete: 0 } };
  const session = first.ui.createSession(lesson(), { saved: stale, save() {}, review: false });
  assert.equal(session.record.step, "mission");
  assert.equal(Object.keys(session.record.answers).length, 0);
});

test("review practice starts at Mission and cannot change durable lesson, quiz or mission records", () => {
  const blob = { progress: { survey: { stars: 3, best: 15, code: "original" } }, quizScores: { 6: 4 }, review: { due: true }, lessons: { survey: { revision: 1, step: "practice", answers: { question: 1 } } } };
  const before = JSON.stringify(blob);
  const view = fixture({ blob, review: true }); practice(view);
  assert.ok(view.article.querySelectorAll(".lesson-choice").every((b) => !b.disabled));
  click(namedButton(view.article.querySelector("fieldset"), "4"));
  click(namedButton(view.article.querySelector("fieldset"), "Try again"));
  click(namedButton(view.article.querySelector("nav"), "Back"));
  assert.equal(JSON.stringify(blob), before);
});

test("legacy levels use exactly the old markdown renderer and Return to code closes the drawer", () => {
  const view = fixture();
  const old = { id: "old", lesson: "# Read\n\n~x < 3~" };
  view.ui.open(old);
  view.context.oldMarkdown = old.lesson;
  assert.equal(view.article.innerHTML, vm.runInContext("md(oldMarkdown)", view.context));
  assert.equal(view.article.children.length, 0);
  view.ui.open(lesson());
  view.body.classList.add("lesson-open");
  click(namedButton(view.article, "Return to code"));
  assert.ok(!view.body.className.includes("lesson-open"));
  assert.equal(view.roots["#btn-lesson"].attributes["aria-expanded"], "false");
  assert.equal(view.focused, view.editor);
});

test("profile loading retains lesson records, defaults legacy saves, and rejects malformed lesson roots", async () => {
  const start = app.indexOf("const store ="), end = app.indexOf("const S =");
  for (const lessons of [undefined, null, "bad", [], { survey: { revision: 1, step: "practice", answers: { question: 0 } } }]) {
    const saved = { progress: { survey: { stars: 2 } }, prefs: {}, quizScores: { 6: 5 }, review: { due: true }, lessons };
    const context = vm.createContext({ PROFILE: "player", api: async () => saved,
      LearningMetrics: { bind() {}, normalizeMetrics: (metrics) => metrics || {} } });
    vm.runInContext(app.slice(start, end), context);
    await vm.runInContext("store.load()", context);
    const loaded = JSON.parse(vm.runInContext("JSON.stringify(store.blob)", context));
    assert.deepEqual(loaded.lessons, lessons && typeof lessons === "object" && !Array.isArray(lessons) ? lessons : {});
    assert.deepEqual(loaded.progress, saved.progress);
    assert.deepEqual(loaded.quizScores, saved.quizScores);
    assert.deepEqual(loaded.review, saved.review);
  }
});

test("hint clicks reveal one position-labelled hint at a time, while legacy hints keep their markup", () => {
  for (const structured of [true, false]) {
    const view = fixture();
    const level = { hints: ["First", "Second", "Third"], ...(structured ? { steps: lesson().steps } : {}) };
    view.context.S.cur = level; view.context.S.hintsShown = 0;
    vm.runInContext(app.slice(app.indexOf('$("#btn-hint").onclick'), app.indexOf('$("#btn-reset").onclick')), view.context);
    for (const [index, label] of ["Concept", "Structure", "Specific"].entries()) {
      click(view.roots["#btn-hint"]);
      assert.equal(view.roots["#hint-list"].children.length, index + 1);
      const markup = view.roots["#hint-list"].children[index].innerHTML;
      if (structured) assert.match(markup, new RegExp(`<strong>${label}</strong>`));
      else assert.equal(markup, level.hints[index]);
    }
    assert.equal(view.roots["#btn-hint"].disabled, true);
    assert.equal(view.store.blob.prefs.hintsSeen?.[level.id], 3, "opened hints are remembered for the next visit");
  }
});

test("client localization carries translated steps and revision without mutating the Arabic level", () => {
  const i18n = fs.readFileSync(__dirname + "/i18n.js", "utf8");
  const level = lesson();
  level.en = { title: "English", steps: [{ id: "mission", title: "Translated", blocks: [] }], lessonRevision: 2 };
  const original = JSON.stringify(level);
  const context = vm.createContext({ LANG: "en", level });
  vm.runInContext(i18n.slice(i18n.indexOf("function localizeLevel"), i18n.indexOf("function localizeChapter")), context);
  const localized = vm.runInContext("localizeLevel(level)", context);
  assert.equal(localized.steps[0].title, "Translated");
  assert.equal(localized.lessonRevision, 2);
  assert.equal(JSON.stringify(level), original);
});

test("practice records from one profile are absent when another profile opens the same mission", () => {
  const first = fixture(); practice(first);
  click(namedButton(first.article.querySelector("fieldset"), "3"));
  const second = fixture(); second.ui.open(lesson());
  assert.match(second.article.textContent, /Step 1 of 2/);
  assert.deepEqual(second.store.blob.lessons, {});
  assert.equal(first.store.blob.lessons.survey.answers.question, 1);
});
