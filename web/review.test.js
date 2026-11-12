"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const source = fs.readFileSync(__dirname + "/app.js", "utf8");

test("review editing and passing preserve the original solution and stars", () => {
  const level = { id: "1-1", starter: "" };
  const original = { stars: 1, best: 20, code: "original" };
  const context = vm.createContext({
    S: { cur: level, review: true, reviewDraft: {} }, level,
    rec: () => original, store: { blob: { review: { due: true } } },
    banner() {}, Sound: { play() {} }, T: (ar, en) => en,
  });
  vm.runInContext(source.slice(source.indexOf("function levelFiles("), source.indexOf("let lastErrors")), context);
  assert.equal(vm.runInContext("fileCode(level, levelFiles(level)[0])", context), "");
  vm.runInContext('saveFile(level, "main.go", "review solution")', context);
  assert.equal(vm.runInContext("fileCode(level, levelFiles(level)[0])", context), "review solution");
  vm.runInContext(source.slice(source.indexOf("function finish("), source.indexOf("// ---------- boot")), context);
  vm.runInContext('finish(level, { success: true, stars: 3, ticks: 1 })', context);
  assert.deepEqual(original, { stars: 1, best: 20, code: "original" });
  assert.equal(context.store.blob.review.due, false);
});

test("review card is bilingual, optional, and offers Start and Skip", () => {
  for (const lang of ["ar", "en"]) {
    class Element {
      constructor() { this.children = []; }
      appendChild(child) { this.children.push(child); }
      querySelector() { return null; }
    }
    const list = new Element();
    const context = vm.createContext({
      S: { levels: [{ id: "1-1", title: "Landing" }], chapters: [] },
      store: { blob: { review: { due: true, level: "1-1" } } },
      document: { createElement: () => new Element() }, $: () => list,
      requestAnimationFrame() {}, esc: (s) => s, T: (ar, en) => lang === "en" ? en : ar,
    });
    vm.runInContext(source.slice(source.indexOf("const CHAPTER_UNITS"), source.indexOf("function showOverlay(")), context);
    vm.runInContext("renderMissions()", context);
    assert.match(list.children[0].innerHTML, lang === "en" ? /Review/ : /مراجعة/);
    assert.deepEqual(list.children[0].children.map((b) => b.textContent), lang === "en" ? ["Start", "Skip"] : ["ابدأ", "تخطّي"]);
    context.store.blob.review.due = false;
    list.children = [];
    vm.runInContext("renderMissions()", context);
    assert.equal(list.children.length, 0);
  }
});
