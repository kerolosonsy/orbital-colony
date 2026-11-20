"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const { context, rank, insertion } = require("./complete.js");

const symbols = [
  { name: "Move", kind: "func", signature: "(dir Direction) bool", doc: "Move moves the drone." },
  { name: "Harvest", kind: "func", signature: "() bool" },
  { name: "East", kind: "const" },
  { name: "Drone", kind: "type" },
  { name: "Move", kind: "method", recv: "Drone" },
  { name: "MaxTicks", kind: "func", signature: "() int" },
];

test("context finds the name typed after colony.", () => {
  assert.deepEqual(context("\tcolony.Mo", 10), { start: 8, prefix: "Mo" });
  assert.deepEqual(context("x := colony.", 12), { start: 12, prefix: "" });
  assert.equal(context("mycolony.Mo", 11), null);
  assert.equal(context("a.colony.Mo", 11), null);
  assert.equal(context("// colony.Mo", 12), null);
  assert.deepEqual(context('fmt.Println("//", colony.Ea', 27), { start: 25, prefix: "Ea" });
  assert.equal(context("colony.Move(", 12), null);
});

test("rank prefers prefix matches and hides methods", () => {
  assert.deepEqual(rank(symbols, "m").map((s) => s.name), ["Move", "MaxTicks"]);
  assert.deepEqual(rank([{ name: "Moss", kind: "const" }, ...symbols], "mo").map((s) => s.name), ["Move", "Moss"]);
  assert.deepEqual(rank(symbols, "ar").map((s) => s.name), ["Harvest"]);
  assert.equal(rank(symbols, "").length, 5);
});

test("functions get parentheses, with the caret inside when they take arguments", () => {
  assert.deepEqual(insertion(symbols[0]), { text: "Move()", caret: 5 });
  assert.deepEqual(insertion(symbols[1]), { text: "Harvest()", caret: 9 });
  assert.deepEqual(insertion(symbols[2]), { text: "East", caret: 4 });
});

test("callContext finds the colony call the caret is inside", () => {
  const { callContext } = require("./complete.js");
  assert.equal(callContext("\tcolony.Move(", 13), "Move");
  assert.equal(callContext("colony.Move(colony.Ea", 21), "Move");
  assert.equal(callContext("colony.Plant(f(x), ", 19), "Plant");
  assert.equal(callContext("colony.Move(colony.East)", 24), null);
  assert.equal(callContext("fmt.Println(", 12), null);
  assert.equal(callContext("// colony.Move(", 15), null);
});

test("memberContext follows colony values to their methods", () => {
  const { memberContext, rankMembers } = require("./complete.js");
  const api = [
    { name: "Spawn", kind: "func", signature: "() *Drone" },
    { name: "Here", kind: "func", signature: "() *Building" },
    { name: "Move", kind: "method", recv: "Drone", signature: "(dir Direction) bool" },
    { name: "Harvest", kind: "method", recv: "Drone", signature: "() bool" },
    { name: "Take", kind: "method", recv: "Building", signature: "(e Entity, n int) (int, error)" },
  ];
  const src = "d := colony.Spawn()\nd.Mo";
  assert.deepEqual(memberContext(src, src.length, api), { start: src.length - 2, prefix: "Mo", recv: "Drone" });
  const typed = "func work(b *colony.Building) {\n\tb.";
  assert.equal(memberContext(typed, typed.length, api).recv, "Building");
  assert.equal(memberContext("x := 3\nx.", 9, api), null);
  assert.equal(memberContext("colony.Mo", 9, api), null);
  assert.deepEqual(rankMembers(api, "Drone", "").map((s) => s.name), ["Move", "Harvest"]);
  assert.deepEqual(rankMembers(api, "Drone", "ha").map((s) => s.name), ["Harvest"]);
});

test("callSymbol resolves functions and methods for signature hints", () => {
  const { callSymbol } = require("./complete.js");
  const api = [
    { name: "Spawn", kind: "func", signature: "() *Drone" },
    { name: "Move", kind: "func", signature: "(dir Direction) bool" },
    { name: "Move", kind: "method", recv: "Drone", signature: "(dir Direction) bool", doc: "Move is Move for this drone." },
  ];
  const a = "colony.Move(";
  assert.equal(callSymbol(a, a.length, api).kind, "func");
  const b = "d := colony.Spawn()\nd.Move(colony.";
  assert.equal(callSymbol(b, b.length, api).recv, "Drone");
  assert.equal(callSymbol("x.Move(", 7, api), null);
});
