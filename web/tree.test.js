"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const context = vm.createContext({});
vm.runInContext(fs.readFileSync(__dirname + "/tree.js", "utf8"), context);
const layoutTree = context.layoutTree;
const layoutBounds = context.layoutBounds;

function assertSpaced(nodes, positions) {
  assert.equal(Object.keys(positions).length, nodes.length);
  for (let i = 0; i < nodes.length; i++) {
    const a = positions[nodes[i].id];
    assert.ok(Number.isFinite(a.x) && Number.isFinite(a.y));
    for (let j = 0; j < i; j++) {
      const b = positions[nodes[j].id];
      assert.ok(Math.abs(a.x - b.x) >= 132 - 1e-8 || Math.abs(a.y - b.y) >= 44 - 1e-8, nodes[i].id + " overlaps " + nodes[j].id);
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

test("sectors use peak ring counts rather than branch totals, in sorted-name order", () => {
  const nodes = [{ id: "base", branch: "", after: [] }];
  for (let i = 0; i < 5; i++) nodes.push({ id: "a" + i, branch: "a", after: [{ id: i ? "a" + (i - 1) : "base" }] });
  for (let i = 0; i < 4; i++) nodes.push({ id: "b" + i, branch: "b", after: [{ id: "base" }] });
  const positions = layoutTree(nodes);
  const angle = (id) => (Math.atan2(positions[id].y, positions[id].x) + 2 * Math.PI) % (2 * Math.PI);
  assert.ok(Math.abs(angle("a0") - Math.PI / 5) < 1e-10);
  assert.ok(Math.abs(angle("b3") - angle("b0") - 24 * Math.PI / 25) < 1e-10);
  assert.deepEqual(layoutTree([...nodes].reverse()), positions);
  assertSpaced(nodes, positions);
});

test("parent-angle ordering shortens edges compared with fixed ID ordering", () => {
  const nodes = [
    { id: "base", branch: "", after: [] },
    ...["a", "b", "c"].map((id) => ({ id, branch: "x", after: [{ id: "base" }] })),
    { id: "first", branch: "x", after: [{ id: "c" }] },
    { id: "middle", branch: "x", after: [{ id: "a" }, { id: "b" }] },
    { id: "last", branch: "x", after: [{ id: "a" }] },
  ];
  const positions = layoutTree(nodes);
  const length = (child, parent) => Math.hypot(child.x - parent.x, child.y - parent.y);
  const edgeLength = (placement) => nodes.slice(4).reduce((sum, node) => sum + node.after.reduce((total, p) => total + length(placement[node.id], positions[p.id]), 0), 0);
  const fixedOrder = { ...positions, first: positions.last, last: positions.first };
  assert.ok(edgeLength(positions) < edgeLength(fixedOrder));
  assert.ok(positions.last.y > 0 && positions.first.y < 0);
  assertSpaced(nodes, positions);
});

test("sparse later rings only move outward as far as card spacing requires", () => {
  const nodes = [
    { id: "base", branch: "", after: [] },
    ...["a", "b", "c", "d"].map((id) => ({ id, branch: "x", after: [{ id: "base" }] })),
    { id: "leaf", branch: "x", after: [{ id: "a" }] },
  ];
  const positions = layoutTree(nodes);
  assertSpaced(nodes, positions);
  const radius = (id) => Math.hypot(positions[id].x, positions[id].y);
  assert.ok(radius("leaf") - radius("a") < 132);
});

test("60 nodes in six uneven branches fit within 3200 by 3200 without card overlaps", () => {
  const nodes = [{ id: "base", branch: "", after: [] }];
  const counts = [[1, 9, 8, 6], [1, 6, 5, 3], [1, 3, 3, 1], [1, 2, 2, 1], [1, 1, 1, 1], [1, 1]];
  counts.forEach((rings, branch) => {
    let parents = ["base"];
    rings.forEach((count, depth) => {
      const ids = Array.from({ length: count }, (_, i) => `b${branch}-d${depth}-${i}`);
      ids.forEach((id, i) => nodes.push({ id, branch: "branch" + branch, after: [{ id: parents[i % parents.length] }] }));
      parents = ids;
    });
  });
  assert.equal(nodes.length, 60);
  const positions = layoutTree(nodes), bounds = layoutBounds(positions);
  assertSpaced(nodes, positions);
  assert.ok(bounds.maxX - bounds.minX <= 3200);
  assert.ok(bounds.maxY - bounds.minY <= 3200);
  assert.deepEqual(layoutTree([...nodes].reverse()), positions);
});

test("bounds include full cards, handle negative coordinates and default empty trees to the origin", () => {
  for (const [positions, expected] of [
    [{}, { minX: 0, minY: 0, maxX: 0, maxY: 0 }],
    [{ root: { x: 0, y: 0, depth: 0 } }, { minX: -66, minY: -22, maxX: 66, maxY: 22 }],
    [{ a: { x: -100, y: 60 }, b: { x: 200, y: -80 } }, { minX: -166, minY: -102, maxX: 266, maxY: 82 }],
  ]) assert.equal(JSON.stringify(layoutBounds(positions)), JSON.stringify(expected));
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
