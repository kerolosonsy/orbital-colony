"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const source = fs.readFileSync(__dirname + "/app.js", "utf8");
function functions(start, end) { return source.slice(source.indexOf(start), source.indexOf(end, source.indexOf(start))); }

test("boss badges stay visible for locked and completed missions in both languages", () => {
  for (const lang of ["ar", "en"]) {
    for (const stars of [0, 3]) {
      class Element {
        constructor() { this.children = []; this.dataset = {}; this.style = { setProperty() {} }; }
        appendChild(child) { this.children.push(child); }
        querySelector() { return null; }
      }
      const list = new Element();
      const context = vm.createContext({ store: { blob: {} },
        S: { chapters: [{ n: 9, track: "Basics" }], levels: [{ id: "plain", title: "Plain", chapter: 9 }, { id: "9-boss-evacuation", title: "Evacuation", chapter: 9, boss: true }], cur: null, running: false },
        document: { createElement: () => new Element() }, $: () => list,
        esc: (s) => s || "", T: (ar, en) => lang === "ar" ? ar : en,
        requestAnimationFrame() {}, rec: () => ({ stars }), unlocked: () => stars > 0,
        sheetCode: (l) => l.id, starHTML: () => "***",
      });
      vm.runInContext(functions("const CHAPTER_UNITS", "function showOverlay("), context);
      vm.runInContext("renderMissions()", context);
      const cells = list.children[0].children[1].children[0].children;
      assert.doesNotMatch(cells[0].innerHTML, /boss-badge/);
      assert.match(cells[1].innerHTML, /boss-badge/);
      assert.ok(cells[1].innerHTML.includes(lang === "en" ? "Track Boss" : "ختام المسار"));
      assert.equal(cells[1].disabled, stars === 0);
    }
  }
});

test("only passing a boss shows track completion, including with old save records", () => {
  for (const lang of ["ar", "en"]) {
    for (const boss of [false, true]) {
      for (const success of [false, true]) {
        const nodes = new Map();
        const node = (selector) => {
          if (!nodes.has(selector)) nodes.set(selector, { innerHTML: "", textContent: "", setAttribute() {}, classList: { toggle() {} } });
          return nodes.get(selector);
        };
        const level = { id: "mission", kind: "lab", boss };
        const record = { stars: 0, best: 0 };
        const context = vm.createContext({ store: { blob: {} },
          S: { cur: level, levels: [level] }, $: node, rec: () => record,
          store: { save() {} }, idx: () => 0, sheetCode: () => "9-boss",
          T: (ar, en) => lang === "ar" ? ar : en, esc: (s) => s,
          banner() {}, log() {}, starHTML: () => "***", updateOverall() {},
          Sound: { play() {} }, setTimeout() {}, showOverlay() {},
          level, result: { success, stars: 2, ticks: 0, message: "lab result" },
        });
        vm.runInContext(functions("function finish(", "// ---------- boot"), context);
        vm.runInContext("finish(level, result)", context);
        const text = node("#result-text").innerHTML;
        assert.equal(text.includes(lang === "en" ? "Track complete!" : "المسار اكتمل!"), boss && success);
      }
    }
  }
});
