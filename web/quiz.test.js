"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const source = fs.readFileSync(__dirname + "/app.js", "utf8");

// The DOM and persistence are boundaries; quiz scoring runs unchanged.
function fixture(lang, scores = {}) {
  class Element {
    constructor() { this.children = []; this.textContent = ""; this.classList = { add() {}, remove() {} }; }
    appendChild(child) { this.children.push(child); }
    replaceChildren() { this.children = []; }
    setAttribute() {}
    scrollIntoView() {}
    focus() {}
    querySelector() { return this.children.at(-1); }
    querySelectorAll() { return this.children.filter((child) => child.onclick); }
  }
  const elements = new Map();
  const context = vm.createContext({
    esc: (s) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"),
    LANG: lang, T: (ar, en) => lang === "en" ? en : ar,
    document: { createElement: () => new Element() },
    $: (selector) => { if (!elements.has(selector)) elements.set(selector, new Element()); return elements.get(selector); },
    store: { blob: { quizScores: scores }, save() {} }, renderMissions() {},
  });
  vm.runInContext(source.slice(source.indexOf("let quizAttempt"), source.indexOf("const CHAPTER_UNITS")), context);
  context.chapter = { n: 1, title: "Chapter", quiz: Array.from({ length: 5 }, () => ({
    question: "سؤال", options: ["أ", "ب", "ج"], explanation: "شرح", answer: 1,
    en: { question: "Question", options: ["A", "B", "C"], explanation: "Explanation" },
  })) };
  vm.runInContext("openQuiz(chapter)", context);
  return { context, elements };
}

test("quizzes give localized feedback, lock answers, and retain best score across retries", () => {
  for (const lang of ["ar", "en"]) {
    const { context, elements } = fixture(lang);
    const fields = elements.get("#quiz-questions").children;
    assert.equal(fields.length, 5);
    for (const [index, field] of fields.entries()) {
      const buttons = field.querySelectorAll();
      assert.equal(buttons[0].textContent, lang === "en" ? "A" : "أ");
      buttons[index === 0 ? 0 : 1].onclick();
      assert.ok(buttons.every((button) => button.disabled));
      assert.ok(field.children.at(-1).textContent.includes(lang === "en" ? "Explanation" : "شرح"));
      // A second click cannot change a recorded answer.
      buttons[1].onclick();
    }
    assert.equal(context.store.blob.quizScores[1], 4);
    assert.match(elements.get('[data-quiz-chapter="1"]').innerHTML, /4\/5/);
    elements.get("#quiz-retry").onclick();
    for (const field of elements.get("#quiz-questions").children) field.querySelectorAll()[0].onclick();
    assert.equal(context.store.blob.quizScores[1], 4);
    elements.get("#quiz-retry").onclick();
    for (const field of elements.get("#quiz-questions").children) field.querySelectorAll()[1].onclick();
    assert.equal(context.store.blob.quizScores[1], 5);
  }
});

test("loading legacy saves defaults quiz scores and loading new saves preserves them", async () => {
  const start = source.indexOf("const store =");
  const end = source.indexOf("const S =", start);
  for (const saved of [{ progress: {}, prefs: {} }, { progress: {}, prefs: {}, quizScores: { 1: 3, 38: 0 } }]) {
    const context = vm.createContext({ PROFILE: "player", api: async () => saved,
      LearningMetrics: { bind() {}, normalizeMetrics: () => ({}) } });
    vm.runInContext(source.slice(start, end), context);
    await vm.runInContext("store.load()", context);
    const scores = JSON.parse(vm.runInContext("JSON.stringify(store.blob.quizScores)", context));
    assert.deepEqual(scores, saved.quizScores || {});
  }
});
