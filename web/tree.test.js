"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const context = vm.createContext({});
vm.runInContext(fs.readFileSync(__dirname + "/tree.js", "utf8"), context);
const layoutTree = context.layoutTree;

function assertSpaced(nodes, positions) {
  assert.equal(Object.keys(positions).length, nodes.length);
  for (let i = 0; i < nodes.length; i++) {
    const a = positions[nodes[i].id];
    assert.ok(Number.isFinite(a.x) && Number.isFinite(a.y));
    for (let j = 0; j < i; j++) {
      const b = positions[nodes[j].id];
      assert.ok(Math.hypot(a.x - b.x, a.y - b.y) >= 96 - 1e-8, nodes[i].id + " overlaps " + nodes[j].id);
    }
    for (const p of nodes[i].after) {
      const parent = positions[p.id];
      assert.ok(Math.hypot(a.x, a.y) > Math.hypot(parent.x, parent.y), nodes[i].id + " inside parent " + p.id);
    }
  }
}

test("cross-branch parents set longest-path depths and do not overlap", () => {
  const nodes = [
    { id: "base", branch: "", after: [] },
    { id: "farm", branch: "agriculture", after: [{ id: "base", level: 1 }] },
    { id: "smelter", branch: "manufacturing", after: [{ id: "base", level: 1 }] },
    { id: "generator", branch: "energy", after: [{ id: "farm", level: 1 }] },
    { id: "assembler", branch: "manufacturing", after: [{ id: "smelter", level: 3 }, { id: "generator", level: 1 }] },
  ];
  const positions = layoutTree(nodes);
  assertSpaced(nodes, positions);
  assert.equal(JSON.stringify(positions.base), JSON.stringify({ x: 0, y: 0, depth: 0 }));
  assert.equal(positions.assembler.depth, 3);
  assert.equal(positions.smelter.depth, 1);
  assert.ok(positions.farm.y > 0);
  assert.ok(positions.smelter.y < 0);
});

test("120-node tree expands crowded rings and stays deterministic under shuffling", () => {
  const nodes = [{ id: "base", branch: "", after: [] }];
  for (let i = 1; i < 120; i++) {
    const parent = i <= 80 ? "base" : "n" + (i - 80);
    nodes.push({ id: "n" + i, branch: "branch" + (i % 9), after: [{ id: parent, level: 1 }] });
  }
  const before = JSON.stringify(nodes);
  const positions = layoutTree(nodes);
  assertSpaced(nodes, positions);
  assert.ok(Math.hypot(positions.n1.x, positions.n1.y) > 180);
  const shuffled = [...nodes.filter((_, i) => i % 2), ...nodes.filter((_, i) => i % 2 === 0)].reverse();
  assert.deepEqual(layoutTree(shuffled), positions);
  assert.deepEqual(layoutTree(nodes), positions);
  assert.equal(JSON.stringify(nodes), before);
});

test("single-branch peers spread evenly inside their sector", () => {
  const nodes = [{ id: "base", branch: "", after: [] }];
  for (const id of ["a", "b", "c", "d"]) {
    nodes.push({ id, branch: "agriculture", after: [{ id: "base", level: 1 }] });
  }
  const positions = layoutTree(nodes);
  assertSpaced(nodes, positions);
  const angles = ["a", "b", "c", "d"].map((id) => (Math.atan2(positions[id].y, positions[id].x) + 2 * Math.PI) % (2 * Math.PI));
  for (let i = 1; i < angles.length; i++) assert.ok(Math.abs(angles[i] - angles[i - 1] - 2 * Math.PI / 5) < 1e-10);
});

test("empty and root-only trees have no stray positions", () => {
  assert.equal(JSON.stringify(layoutTree([])), "{}");
  const positions = layoutTree([{ id: "base", branch: "", after: [] }]);
  assert.equal(JSON.stringify(positions), '{"base":{"x":0,"y":0,"depth":0}}');
});

test("invalid graphs fail explicitly", () => {
  for (const nodes of [
    [{ id: "base", branch: "", after: [] }, { id: "a", branch: "x", after: [{ id: "missing" }] }],
    [{ id: "base", branch: "", after: [] }, { id: "a", branch: "x", after: [{ id: "b" }] }, { id: "b", branch: "x", after: [{ id: "a" }] }],
    [{ id: "base", branch: "", after: [] }, { id: "island", branch: "x", after: [] }],
    [{ id: "base", branch: "", after: [] }, { id: "base", branch: "x", after: [] }],
  ]) assert.throws(() => layoutTree(nodes));
});
