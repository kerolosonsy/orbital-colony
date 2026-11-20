"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const { stepsFor, bubbleBox, keyAction } = require("./tour.js");

const FIRST = ["1-1-first-signal", "1-2-first-harvest", "1-3-planting"];
const levels = fs.readFileSync(path.join(__dirname, "../colony/levels/ch01.go"), "utf8");
const html = fs.readFileSync(path.join(__dirname, "index.html"), "utf8");

// starter returns the Arabic starter code of a chapter-1 level from ch01.go
function starter(id) {
  const at = levels.indexOf(`ID: "${id}"`);
  const m = /Starter: `([^`]*)`/.exec(levels.slice(at));
  return m[1];
}

test("the first three missions each have a tour in both languages", () => {
  for (const id of FIRST) {
    const steps = stepsFor(id);
    assert.ok(steps?.length >= 5, id);
    for (const step of steps) assert.ok(step.ar && step.en, `${id}: ${step.en}`);
  }
  assert.equal(stepsFor("2-1-anything"), null);
});

test("tour targets exist and typing anchors match the starter code", () => {
  for (const id of FIRST) {
    const lines = starter(id).split("\n");
    for (const step of stepsFor(id)) {
      if (typeof step.at === "string") assert.match(html, new RegExp(`id="${step.at.slice(1)}"|class="${step.at.slice(1)}"`), `${id}: ${step.at}`);
      if (step.at?.line) assert.ok(lines[step.at.line - 1] !== undefined, `${id}: line ${step.at.line}`);
      if (step.act) assert.ok(starter(id).includes(step.act.after), `${id}: anchor ${JSON.stringify(step.act.after)}`);
    }
  }
  // the lines the tour talks about hold what it says
  assert.match(starter(FIRST[0]).split("\n")[5], /colony\.Move/);
  assert.match(starter(FIRST[1]).split("\n")[6], /colony\.Harvest/);
  assert.match(starter(FIRST[2]).split("\n")[6], /colony\.Till/);
});

test("bubbleBox prefers below, then above, then the roomier side, and stays on screen", () => {
  assert.equal(bubbleBox({ left: 100, top: 50, right: 200, bottom: 80, width: 100, height: 30 }, 300, 100, 1000, 800).side, "below");
  assert.equal(bubbleBox({ left: 100, top: 700, right: 200, bottom: 760, width: 100, height: 60 }, 300, 100, 1000, 800).side, "above");
  const tall = bubbleBox({ left: 50, top: 10, right: 300, bottom: 790, width: 250, height: 780 }, 300, 100, 1000, 800);
  assert.equal(tall.side, "right");
  assert.ok(tall.left >= 300);
  const edge = bubbleBox({ left: 980, top: 50, right: 1000, bottom: 70, width: 20, height: 20 }, 300, 100, 1000, 800);
  assert.equal(edge.left, 1000 - 300 - 12);
  assert.equal(bubbleBox({ left: 100, top: 400, right: 200, bottom: 420, width: 100, height: 20 }, 300, 100, 1000, 800, "above").side, "above");
  assert.equal(bubbleBox({ left: 100, top: 40, right: 200, bottom: 60, width: 100, height: 20 }, 300, 100, 1000, 800, "above").side, "below");
  assert.deepEqual(bubbleBox(null, 300, 100, 1000, 800), { left: 350, top: 350, side: "center" });
});

test("keyAction mirrors the arrows in Arabic", () => {
  assert.equal(keyAction("Enter", false), "next");
  assert.equal(keyAction("Escape", true), "close");
  assert.equal(keyAction("ArrowRight", false), "next");
  assert.equal(keyAction("ArrowLeft", true), "next");
  assert.equal(keyAction("ArrowRight", true), "back");
  assert.equal(keyAction("a", false), null);
});
