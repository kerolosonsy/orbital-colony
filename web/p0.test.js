"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const source = fs.readFileSync(__dirname + "/app.js", "utf8");
function functions(start, end) { return source.slice(source.indexOf(start), source.indexOf(end, source.indexOf(start))); }
function fixture(levels, progress = {}) {
  const context = vm.createContext({ S: { levels }, rec: (id) => progress[id] || {} });
  vm.runInContext(functions("function unlocked(", "function starHTML("), context);
  return context;
}
test("saved progress unlocks inserted missions and service levels do not block", () => {
  const levels = [{ id: "a" }, { id: "inserted" }, { id: "b" }, { id: "service", service: "mongo" }, { id: "after" }];
  let context = fixture(levels, { b: { stars: 2 } });
  assert.equal(vm.runInContext("unlocked(1)", context), true);
  assert.equal(vm.runInContext("unlocked(3)", context), true);
  assert.equal(vm.runInContext("unlocked(4)", context), true);
  context = fixture(levels);
  assert.equal(vm.runInContext("unlocked(0)", context), true);
  assert.equal(vm.runInContext("unlocked(1)", context), false);
  assert.equal(vm.runInContext("unlocked(4)", context), true);
  assert.equal(vm.runInContext("unlocked(-1)", context), false);
});
test("carry-forward respects current saves and read-only non-Go files", () => {
  const progress = { old: { files: { "main.go": "old solution", "schema.sql": "tampered" } }, current: { files: {} } };
  const context = vm.createContext({ rec: (id) => progress[id] || {}, level: { id: "current", carryFrom: "old", files: [
    { name: "main.go", starter: "starter" }, { name: "schema.sql", starter: "SELECT 1" }
  ] } });
  vm.runInContext(functions("function levelFiles(", "function saveFile("), context);
  assert.equal(vm.runInContext("fileCode(level, levelFiles(level)[0])", context), "old solution");
  assert.equal(vm.runInContext("fileCode(level, levelFiles(level)[1])", context), "SELECT 1");
  assert.equal(vm.runInContext("levelFiles(level)[1].readOnly", context), true);
  progress.current.files["main.go"] = "new solution";
  assert.equal(vm.runInContext("fileCode(level, levelFiles(level)[0])", context), "new solution");
});
test("chapters sharing a track stay in one collapsible group", () => {
  class Element {
    constructor(tag) { this.tag = tag; this.children = []; this.dataset = {}; this.style = { setProperty() {} }; }
    appendChild(child) { this.children.push(child); }
    querySelector() { return null; }
  }
  const list = new Element("div");
  const context = vm.createContext({
    S: { chapters: [{ n: 1, track: "Foundations" }, { n: 10, track: "Concurrency" }, { n: 14, track: "Craft" }, { n: 24, track: "Concurrency" }], levels: [], cur: null },
    document: { createElement: (tag) => new Element(tag) },
    $: () => list, esc: (s) => s || "", requestAnimationFrame: () => {},
  });
  vm.runInContext(functions("const CHAPTER_UNITS", "function showOverlay("), context);
  vm.runInContext("renderMissions()", context);
  assert.equal(list.children.length, 3);
  const concurrency = list.children[1];
  assert.equal(concurrency.tag, "details");
  assert.equal(concurrency.open, true);
  assert.equal(concurrency.children[0].tag, "summary");
  assert.equal(concurrency.children.length, 3); // summary and chapters 10 and 24
});
