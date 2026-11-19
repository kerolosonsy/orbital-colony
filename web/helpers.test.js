"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const source = fs.readFileSync(__dirname + "/app.js", "utf8");
function functions(start, end) { return source.slice(source.indexOf(start), source.indexOf(end, source.indexOf(start))); }

test("server hints name the two common local failures and nothing else", () => {
  const context = vm.createContext({ T: (ar, en) => en });
  vm.runInContext(functions("function serverHint(", "let lastServerHint"), context);
  const hint = (error) => vm.runInContext("serverHint(error)", Object.assign(context, { error }));
  assert.match(hint(new TypeError("Failed to fetch")), /make start/);
  assert.match(hint(new Error("403 missing token")), /Refresh/);
  assert.match(hint(new Error("403 رمز الدخول مفقود أو غير صالح")), /Refresh/);
  assert.equal(hint(new Error("403 forbidden host")), "");
  assert.equal(hint(new Error("500 boom")), "");
});

test("mission links accept the full id or the sheet code", () => {
  const context = vm.createContext({
    S: { levels: [{ id: "1-1-first-signal" }, { id: "3-1-long-row" }, { id: "13-1-tick" }] },
    sheetCode: (l) => l.id.split("-").slice(0, 2).join("-"),
    location: { hash: "" },
  });
  vm.runInContext(functions("function hashLevelID(", "window.addEventListener(\"hashchange\""), context);
  const resolve = (hash) => { context.location.hash = hash; return vm.runInContext("hashLevelID()", context); };
  assert.equal(resolve("#3-1"), "3-1-long-row");
  assert.equal(resolve("#13-1"), "13-1-tick");
  assert.equal(resolve("#1-1-first-signal"), "1-1-first-signal");
  assert.equal(resolve("#nope"), "nope");
});

test("index search: numbers are exact sheet codes, words match titles and chapters", () => {
  class Node {
    constructor(props = {}) { Object.assign(this, { hidden: false, dataset: {}, children: [], open: false }, props); }
    querySelectorAll(sel) { return sel === ".cell" ? this.children : []; }
    querySelector(sel) { return sel === ".quad-name" ? this.name : sel === ".quad-row:not([hidden])" ? this.rows.find((r) => !r.hidden) || null : null; }
    closest() { return null; }
  }
  const cell = (code, title) => new Node({ dataset: { code, search: `${code} ${title}`.toLowerCase() } });
  const rowA = new Node({ name: { textContent: "3 Loops for" }, children: [cell("3-1", "Long row"), cell("3-2", "Snake")] });
  const rowB = new Node({ name: { textContent: "13 Channels" }, children: [cell("13-1", "Pipes")] });
  const group = new Node({ rows: [rowA, rowB] });
  const input = { value: "" }, empty = { hidden: true };
  const context = vm.createContext({
    $: (sel) => (sel === "#mission-search" ? input : sel === "#mission-search-empty" ? empty : { addEventListener() {} }),
    document: { querySelectorAll: (sel) => (sel.endsWith(".quad-row") ? [rowA, rowB] : sel.endsWith(".track-group") ? [group] : []) },
  });
  vm.runInContext(functions("function applyMissionFilter(", "$(\"#mission-search\").addEventListener(\"input\""), context);
  const shown = (query) => { input.value = query; vm.runInContext("applyMissionFilter()", context); return [rowA, rowB].flatMap((r) => r.children).filter((c) => !c.hidden).map((c) => c.dataset.code); };
  assert.deepEqual(shown("3-1"), ["3-1"]);
  assert.deepEqual(shown("3"), ["3-1", "3-2"]);
  assert.deepEqual(shown("snake"), ["3-2"]);
  assert.deepEqual(shown("channels"), ["13-1"]);
  assert.deepEqual(shown("zzz"), []);
  assert.equal(empty.hidden, false);
  assert.equal(shown("").length, 3);
});
