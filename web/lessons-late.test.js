"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const { TestElement } = require("./sandbox/harness.js");
const lessonSource = fs.readFileSync(__dirname + "/lessons.js", "utf8");
const lateSource = fs.readFileSync(__dirname + "/lessons-late.js", "utf8");
const app = fs.readFileSync(__dirname + "/app.js", "utf8");

function fixture({ lang = "en", reduceMotion = false, animations = false } = {}) {
  const animationCalls = [];
  const createElement = (tag) => new TestElement(tag, { animations, animationCalls });
  const roots = Object.fromEntries(["#lesson", "#lesson-panel", "#btn-lesson", "#btn-hint", "#hint-list"].map((id) => [id, createElement("section")]));
  roots["#btn-hint"].appendChild(createElement("span"));
  const editor = createElement("textarea"), body = createElement("body");
  const store = { blob: { progress: {}, prefs: {}, quizScores: {}, review: {}, lessons: {} }, save() {} };
  const context = vm.createContext({
    ResizeObserver: class { observe() {} disconnect() {} },
    document: { body, createElement, querySelector: (selector) => roots["#lesson"].querySelector(selector) },
    $: (selector) => roots[selector], ed: editor, store, S: { review: false },
    T: (arabic, english) => lang === "en" ? english : arabic,
    Sound: { play() {} }, matchMedia: () => ({ matches: reduceMotion }),
  });
  vm.runInContext(app.slice(app.indexOf("const esc ="), app.indexOf("// ---------- files ----------")), context);
  vm.runInContext(lessonSource + "\n" + lateSource + "\nthis.openLessons = Lessons.open;", context);
  return { context, roots, animationCalls };
}

function openSandbox(view, id, spec, language = "en") {
  const ar = language === "ar";
  const level = {
    id, lessonRevision: 1,
    steps: [{ id: "try", title: ar ? "جرّب" : "Try it", blocks: [{ id: "sandbox", kind: "sandbox", text: ar ? "جرّب التفاعل." : "Try the interaction.", sandbox: spec }] }],
  };
  view.context.level = level;
  vm.runInContext("openLessons(level)", view.context);
  return view.roots["#lesson"].querySelector(".lesson-sandbox");
}

function click(button) { assert.ok(button, "button exists"); assert.equal(button.disabled, false); button.onclick(); }

function sortSpec(lang) {
  const ar = lang === "ar";
  return {
    kind: "sort-compare", sequence: [3, 5, 3, 5, 2],
    sourceLabel: ar ? "أحجام توضيحية." : "Illustrative sizes.",
    sortDirectionLabel: ar ? "اتجاه الحجم" : "Size direction", sortTieLabel: ar ? "كسر التعادل" : "Tie-break",
    ascendingLabel: ar ? "تصاعدي" : "Ascending", descendingLabel: ar ? "تنازلي" : "Descending",
    tieOnLabel: ar ? "شغّال" : "On", tieOffLabel: ar ? "مقفول" : "Off", resetLabel: ar ? "من الأول" : "Reset",
  };
}

function gridSpec(lang) {
  const ar = lang === "ar";
  return {
    kind: "grid-index", grid: [[3, 7, 2], [8, 6, 4]],
    sourceLabel: ar ? "شبكة توضيحية." : "Illustrative grid.",
    gridOuterCloneLabel: ar ? "نسخة outer بس" : "Clone outer only",
    gridRowCloneLabel: ar ? "نسخ كل row" : "Clone each row",
    gridEditCopyLabel: ar ? "عدّل النسخة" : "Edit copy", resetLabel: ar ? "من الأول" : "Reset",
  };
}

function positions(row) { return row.querySelectorAll(".lesson-sort-position").map((node) => node.textContent); }

test("sort-compare reorders, highlights ties, toggles both rules, and resets in Arabic and English", () => {
  for (const lang of ["ar", "en"]) {
    const view = fixture({ lang }), spec = sortSpec(lang);
    const box = openSandbox(view, "6-4-stock-order", spec, lang);
    const row = box.querySelector(".lesson-sort-row");
    const cards = () => row.querySelectorAll(".lesson-sort-card");
    assert.deepEqual(positions(row), ["x=4", "x=2", "x=0", "x=3", "x=1"]);
    assert.deepEqual(cards().filter((card) => card.classList.contains("same-size")).map((card) => card.dataset.position), ["2", "0", "3", "1"]);
    const controls = box.querySelector(".lesson-sort-controls").querySelectorAll("button");
    const direction = controls[0], ties = controls[1], reset = controls[2];
    assert.match(direction.textContent, lang === "ar" ? /تصاعدي/ : /Ascending/);
    click(direction);
    assert.deepEqual(positions(row), ["x=3", "x=1", "x=2", "x=0", "x=4"]);
    assert.equal(direction.getAttribute("aria-pressed"), "true");
    click(ties);
    assert.deepEqual(positions(row), ["x=1", "x=3", "x=0", "x=2", "x=4"]);
    assert.equal(ties.getAttribute("aria-pressed"), "true");
    click(reset);
    assert.deepEqual(positions(row), ["x=4", "x=2", "x=0", "x=3", "x=1"]);
    assert.equal(direction.getAttribute("aria-pressed"), "false");
    assert.equal(ties.getAttribute("aria-pressed"), "false");
    assert.match(reset.textContent, lang === "ar" ? /من الأول/ : /Reset/);
  }
});

test("sort animation runs when allowed and is skipped for reduced motion", () => {
  const animated = fixture({ animations: true, reduceMotion: false });
  const box = openSandbox(animated, "6-4-stock-order", sortSpec("en"));
  click(box.querySelector(".lesson-sort-controls").querySelectorAll("button")[0]);
  assert.ok(animated.animationCalls.length > 0);

  const reduced = fixture({ animations: true, reduceMotion: true });
  const reducedBox = openSandbox(reduced, "6-4-stock-order", sortSpec("en"));
  click(reducedBox.querySelector(".lesson-sort-controls").querySelectorAll("button")[0]);
  assert.equal(reduced.animationCalls.length, 0);
});

test("grid-index reports row and column, shows shallow sharing, and resets in Arabic and English", () => {
  for (const lang of ["ar", "en"]) {
    const view = fixture({ lang }), spec = gridSpec(lang);
    const box = openSandbox(view, "6-5-heat-grid", spec, lang);
    assert.match(box.textContent, lang === "ar" ? /الأصل/ : /Original/);
    assert.match(box.textContent, lang === "ar" ? /النسخة/ : /Copy/);
    const views = box.querySelectorAll(".lesson-grid-view");
    const original = views[0], copy = views[1];
    assert.equal(box.querySelector(".lesson-grid-expression").textContent, "grid[0][0] = 3");
    const cells = original.querySelectorAll(".lesson-grid-cell-button");
    click(cells[3]);
    assert.equal(box.querySelector(".lesson-grid-expression").textContent, "grid[1][0] = 8");
    assert.ok(original.querySelectorAll(".lesson-grid-row")[1].classList.contains("selected-axis"));
    assert.ok(original.querySelectorAll(".lesson-grid-column")[0].classList.contains("selected-axis"));

    const modes = box.querySelector(".lesson-grid-modes").querySelectorAll("button");
    const editAndReset = box.querySelector(".lesson-grid-controls").querySelectorAll("button");
    assert.equal(modes[0].getAttribute("aria-pressed"), "true");
    click(editAndReset[0]);
    assert.equal(original.querySelectorAll(".lesson-grid-cell-button")[3].textContent, "99");
    assert.equal(copy.querySelectorAll(".lesson-grid-value")[3].textContent, "99");

    click(modes[1]);
    assert.equal(modes[1].getAttribute("aria-pressed"), "true");
    assert.equal(original.querySelectorAll(".lesson-grid-cell-button")[3].textContent, "8");
    click(editAndReset[0]);
    assert.equal(original.querySelectorAll(".lesson-grid-cell-button")[3].textContent, "8");
    assert.equal(copy.querySelectorAll(".lesson-grid-value")[3].textContent, "99");

    click(editAndReset[1]);
    assert.equal(box.querySelector(".lesson-grid-expression").textContent, "grid[0][0] = 3");
    assert.equal(box.querySelector(".lesson-grid-modes").querySelectorAll("button")[0].getAttribute("aria-pressed"), "true");
    assert.equal(box.querySelectorAll(".lesson-grid-view")[0].querySelectorAll(".lesson-grid-cell-button")[0].textContent, "3");
  }
});
