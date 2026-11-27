"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const { codeHash, point, tone, view } = require("./grade.js");
const T = (ar, en) => en;

test("codeHash names one version of the files, whatever the key order", () => {
  assert.equal(codeHash({ "a.go": "x", "b.go": "y" }), codeHash({ "b.go": "y", "a.go": "x" }));
  assert.notEqual(codeHash({ "main.go": "colony.Move()" }), codeHash({ "main.go": "colony.Move() " }));
  assert.notEqual(codeHash({ "a.go": "xy" }), codeHash({ "ax.go": "y" }));
});

test("points escape HTML and mark identifiers as code", () => {
  assert.equal(point("rename `x` to <b>"), 'rename <code dir="ltr">x</code> to &lt;b&gt;');
  assert.deepEqual([tone(95), tone(80), tone(40)], ["high", "mid", "low"]);
});

test("the panel shows the score, opens why or how one at a time, and flags changed code", () => {
  const grade = { score: 82, summary: "Solid.", why: ["clear names"], how: ["name the `3`", "drop the dead loop"] };
  const closed = view({ kind: "done", grade, open: null, stale: false }, T);
  assert.match(closed, /82<small>\/100<\/small>/);
  assert.match(closed, /Why 82\?/);
  assert.match(closed, /How do I reach 100\?/);
  assert.equal((closed.match(/<ul[^>]* hidden>/g) || []).length, 2);
  const how = view({ kind: "done", grade, open: "how", stale: true }, T);
  assert.match(how, /<ul class="grade-list grade-how">/);
  assert.match(how, /<ul class="grade-list grade-why" hidden>/);
  assert.match(how, /The code changed since this score/);
  const perfect = view({ kind: "done", grade: { score: 100, summary: "Nothing to change.", why: ["all good"], how: [] }, open: null }, T);
  assert.doesNotMatch(perfect, /reach 100/);
});

test("other states: ask, busy, setup error, hidden", () => {
  assert.match(view({ kind: "ask" }, T), /Score my code/);
  assert.match(view({ kind: "busy" }, T), /Reading your code/);
  assert.match(view({ kind: "error", setup: true }, T), /codex login/);
  assert.match(view({ kind: "error" }, T), /Try again/);
  assert.equal(view({ kind: "hidden" }, T), "");
});
