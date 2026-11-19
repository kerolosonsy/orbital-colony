"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const { TestElement } = require("./sandbox/harness.js");
const source = fs.readFileSync(__dirname + "/metrics.js", "utf8");

function metricsContext() {
  let now = 100;
  const documentEvents = {};
  const windowEvents = {};
  const context = vm.createContext({
    document: { visibilityState: "visible", hasFocus: () => true,
      addEventListener: (name, callback) => { documentEvents[name] = callback; },
      createElement: (tag) => new TestElement(tag) },
    performance: { now: () => now },
    addEventListener: (name, callback) => { windowEvents[name] = callback; },
    T: (ar, en) => en,
  });
  vm.runInContext(source, context);
  return { api: context.LearningMetrics, context, setNow: (next) => { now = next; }, documentEvents, windowEvents };
}

test("metrics record run types, lesson activity, active time, and first solve", () => {
  const view = metricsContext();
  const store = { blob: { metrics: {} }, saves: 0, save() { this.saves++; }, push() {} };
  view.api.bind(store);
  view.api.openLevel(store, "1-1");
  view.api.recordLessonOpened(store, "1-1");
  view.api.recordStepViewed(store, "1-1", "intro");
  view.api.recordHint(store, "1-1", 0);
  view.api.recordRun(store, "1-1", { phase: "compile", failure: "compile", result: { success: false } });
  view.setNow(45100);
  view.api.recordRun(store, "1-1", { phase: "run", failure: "runtime", result: { success: false } });
  view.setNow(120100);
  view.api.recordInteraction(store, "1-1");
  view.api.recordRun(store, "1-1", { phase: "run", result: { success: true, ticks: 15 } });
  const entry = store.blob.metrics["1-1"];
  assert.equal(entry.runs, 3);
  assert.equal(entry.failedRuns.compile, 1);
  assert.equal(entry.failedRuns.runtime, 1);
  assert.equal(entry.failedRuns.check, 0);
  assert.equal(entry.lessonsOpened, 1);
  assert.equal(entry.lessonStepsViewed.intro, 1);
  assert.equal(entry.interactions, 1);
  assert.equal(entry.hintCount, 1);
  assert.deepEqual(Array.from(entry.hintsRevealed), ["0"]);
  assert.equal(entry.activeTimeMs, 120000);
  assert.equal(entry.timeToFirstSolveMs, 120000);
  assert.equal(entry.bestTicks, 15);
  assert.equal(entry.solved, true);
  assert.ok(store.saves >= 5);
});

test("metrics view renders ranked struggle levels and chapter summaries", () => {
  const view = metricsContext();
  const roots = { struggles: new TestElement("ol"), chapters: new TestElement("div") };
  const metrics = {
    "1-1": { failedRuns: { compile: 1 }, hintCount: 1, activeTimeMs: 60000 },
    "1-2": { failedRuns: { check: 2 }, hintCount: 1, activeTimeMs: 0 },
    "2-1": { hintCount: 1, activeTimeMs: 180000 },
  };
  const levels = [
    { id: "1-1", title: "Landing", chapter: 1 },
    { id: "1-2", title: "Signal", chapter: 1 },
    { id: "2-1", title: "Harvest", chapter: 2 },
  ];
  const chapters = [{ n: 1, title: "First orbit" }, { n: 2, title: "Harvest" }, { n: 3, title: "Unused" }];
  const ranked = view.api.rankStruggles(metrics, levels);
  assert.deepEqual(Array.from(ranked, (row) => row.id), ["1-2", "1-1", "2-1"]);
  assert.deepEqual(Array.from(ranked, (row) => row.score), [10, 7, 5]);

  view.api.renderStats(roots, metrics, { levels, chapters });
  assert.equal(roots.struggles.children.length, 3);
  assert.match(roots.struggles.children[0].textContent, /1-2.*Signal.*Score 10/);
  assert.match(roots.struggles.children[1].textContent, /1-1.*Landing/);
  assert.equal(roots.chapters.children.length, 2);
  assert.match(roots.chapters.children[0].textContent, /First orbit.*3 failed/);
  assert.match(roots.chapters.children[1].textContent, /Harvest.*1 solved|Harvest.*0 solved/);
});

test("missing metric fields migrate to zero-valued counters", () => {
  const api = metricsContext().api;
  const migrated = api.normalizeMetrics({ "1-1": { runs: 2, hintsRevealed: ["0", "0"] } });
  assert.equal(migrated["1-1"].runs, 2);
  assert.equal(migrated["1-1"].failedRuns.check, 0);
  assert.equal(migrated["1-1"].hintCount, 1);
  assert.equal(migrated["1-1"].activeTimeMs, 0);
  assert.equal(migrated["1-1"].solved, false);
});
