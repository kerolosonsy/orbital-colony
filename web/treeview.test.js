"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const escapeHTML = (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
const context = vm.createContext({ module: { exports: {} }, LANG: "ar", T: (ar, en) => context.LANG === "en" ? en : ar, esc: escapeHTML });
vm.runInContext(fs.readFileSync(__dirname + "/treeview.js", "utf8"), context);
const { fitView, edgeSatisfied, missingCount, pathToNearestOwned, matchesSearch, buyEnabled, buyLabel, checkMarkup, branchHue } = context.module.exports;

test("fit centres offset bounds with a screen-space margin and never zooms above 100%", () => {
  for (const [bounds, width, height, margin, expected] of [
    [{ minX: -200, minY: -100, maxX: 600, maxY: 300 }, 500, 400, 50, { zoom: 0.5, x: 150, y: 150 }],
    [{ minX: -100, minY: -400, maxX: 100, maxY: 400 }, 500, 400, 40, { zoom: 0.4, x: 250, y: 200 }],
    [{ minX: 100, minY: -20, maxX: 232, maxY: 24 }, 1000, 800, 32, { zoom: 1, x: 334, y: 398 }],
    [{ minX: -2000, minY: -2000, maxX: 2000, maxY: 2000 }, 500, 400, 40, { zoom: 0.08, x: 250, y: 200 }],
    [{ minX: 0, minY: 0, maxX: 0, maxY: 0 }, 500, 400, 32, { zoom: 1, x: 250, y: 200 }],
  ]) {
    const fit = fitView(bounds, width, height, margin);
    assert.equal(JSON.stringify(fit), JSON.stringify(expected));
    assert.ok(bounds.minX * fit.zoom + fit.x >= margin);
    assert.ok(bounds.maxX * fit.zoom + fit.x <= width - margin);
    assert.ok(bounds.minY * fit.zoom + fit.y >= margin);
    assert.ok(bounds.maxY * fit.zoom + fit.y <= height - margin);
  }
});

test("opening and re-centring fit known cards, with a lower zoom floor and overview names", () => {
  const element = () => ({
    attributes: new Map(), events: {}, classes: new Set(),
    setAttribute(name, value) { this.attributes.set(name, value); },
    addEventListener(name, listener) { this.events[name] = listener; },
    closest() { return null; },
  });
  const dialog = element(), canvas = element(), world = element();
  const centre = element(), input = element(), detail = element(), nav = element(), close = element();
  const selectors = { ".tv-canvas": canvas, ".tv-world": world, ".tv-detail": detail, ".tv-branches": nav, input, '[data-action="centre"]': centre, '[data-action="close"]': close };
  dialog.querySelector = (selector) => selectors[selector];
  dialog.querySelectorAll = () => [];
  dialog.showModal = () => { dialog.open = true; };
  dialog.classList = { toggle(name, active) { if (active) dialog.classes.add(name); else dialog.classes.delete(name); } };
  canvas.getBoundingClientRect = () => ({ width: 800, height: 600 });
  const ui = vm.createContext({
    document: { activeElement: null, createElement: () => dialog, body: { appendChild() {} } },
    ResizeObserver: class { observe() {} },
    LANG: "en", T: (_ar, en) => en, esc: escapeHTML,
  });
  vm.runInContext(fs.readFileSync(__dirname + "/tree.js", "utf8"), ui);
  vm.runInContext(fs.readFileSync(__dirname + "/treeview.js", "utf8"), ui);
  const view = vm.runInContext("TreeView", ui);
  const nodes = Array.from({ length: 41 }, (_, i) => ({
    id: "n" + i, branch: i ? "x" : "", after: i ? [{ id: "n" + (i - 1), level: 1 }] : [],
    state: i ? "unknown" : "owned", level: 1, name: "Node " + i, checks: [],
  }));
  const transform = () => {
    const [, x, y, zoom] = world.attributes.get("transform").match(/translate\(([^ ]+) ([^)]+)\) scale\(([^)]+)\)/);
    return { x: Number(x), y: Number(y), zoom: Number(zoom) };
  };
  const key = (key) => canvas.events.keydown({ key, target: canvas, preventDefault() {} });
  view.open({ nodes });
  assert.deepEqual(transform(), { x: 400, y: 300, zoom: 1 });
  for (let i = 0; i < 20; i++) key("-");
  assert.equal(transform().zoom, 0.25);
  assert.ok(dialog.classes.has("is-overview"));
  const known = nodes.map((node) => ({ ...node, state: "owned" }));
  view.update({ nodes: known }); centre.onclick();
  const fit = transform(), positions = ui.layoutTree(known);
  assert.ok(fit.zoom < 0.25);
  for (const point of Object.values(positions)) {
    assert.ok((point.x - 66) * fit.zoom + fit.x >= 32 - 1e-8);
    assert.ok((point.x + 66) * fit.zoom + fit.x <= 768 + 1e-8);
    assert.ok((point.y - 22) * fit.zoom + fit.y >= 32 - 1e-8);
    assert.ok((point.y + 22) * fit.zoom + fit.y <= 568 + 1e-8);
  }
  for (let i = 0; i < 20; i++) key("-");
  assert.equal(transform().zoom, fit.zoom);
  for (let i = 0; i < 20; i++) key("+");
  assert.ok(!dialog.classes.has("is-overview"));
  centre.onclick();
  assert.deepEqual(transform(), fit);
});

test("edges require the parent's exact prerequisite level or higher", () => {
  const node = { after: [{ id: "smelter", level: 3 }] };
  for (const [level, expected] of [[0, false], [2, false], [3, true], [4, true]]) {
    assert.equal(edgeSatisfied(node, node.after[0], { smelter: { level } }), expected);
  }
  assert.equal(edgeSatisfied(node, node.after[0], {}), false);
});

test("missing count includes every failed condition, not just parents", () => {
  const checks = ["parent", "contract", "rank", "rp", "cost"].map((kind, i) => ({ kind, ok: i % 2 === 0 }));
  assert.equal(missingCount({ checks }), 2);
  assert.equal(missingCount({ checks: [] }), 0);
  assert.equal(missingCount({ checks: checks.map((check) => ({ ...check, ok: true })) }), 0);
});

test("path starts at the nearest owned ancestor across multiple branches", () => {
  const byId = {
    base: { state: "owned", after: [] },
    long: { state: "locked", after: [{ id: "base" }] },
    longer: { state: "locked", after: [{ id: "long" }] },
    near: { state: "owned", after: [{ id: "base" }] },
    target: { state: "locked", after: [{ id: "longer" }, { id: "near" }] },
    leaf: { state: "locked", after: [{ id: "target" }] },
  };
  assert.equal(JSON.stringify(pathToNearestOwned("leaf", byId)), '["near","target","leaf"]');
  assert.equal(JSON.stringify(pathToNearestOwned("near", byId)), "[]");
  assert.equal(JSON.stringify(pathToNearestOwned("base", byId)), "[]");
});

test("equally near ancestors choose a stable chain and unreachable targets have no path", () => {
  const byId = {
    a: { state: "owned", after: [] }, b: { state: "owned", after: [] },
    target: { state: "locked", after: [{ id: "b" }, { id: "a" }] },
    isolated: { state: "locked", after: [] },
  };
  assert.equal(JSON.stringify(pathToNearestOwned("target", byId)), '["a","target"]');
  assert.equal(JSON.stringify(pathToNearestOwned("isolated", byId)), "[]");
  assert.equal(JSON.stringify(pathToNearestOwned("missing", byId)), "[]");
});

test("name search matches both languages and case, but never reveals unknown nodes", () => {
  const node = { name: "مصهر الحديد", nameEn: "Iron Smelter", state: "locked" };
  for (const [query, expected] of [["الحديد", true], ["SMELTER", true], [" iron ", true], ["generator", false], ["", false], ["  ", false]]) {
    assert.equal(matchesSearch(node, query), expected);
  }
  for (const query of ["مصهر", "IRON", ""]) assert.equal(matchesSearch({ ...node, state: "unknown" }, query), false);
});

test("buy is enabled for available nodes and owned nodes whose every check passes", () => {
  for (const state of ["owned", "available", "near", "locked", "unknown"]) {
    for (const ok of [true, false]) {
      assert.equal(buyEnabled({ state, checks: [{ ok: true }, { ok }] }), state === "available" || state === "owned" && ok, `${state}, checks=${ok}`);
    }
  }
  assert.equal(buyEnabled({ state: "owned", checks: [] }), true);
});

test("nine branch hues are distinct, evenly spaced and stable under node ordering", () => {
  const branches = ["زراعة", "تعدين", "تصنيع", "طاقة", "أسطول", "بيانات", "خدمات", "بحث", "نقل"];
  const hues = branches.map((branch) => branchHue(branch, branches));
  assert.equal(new Set(hues).size, 9);
  assert.deepEqual([...hues].sort((a, b) => a - b).map((hue, i) => hue - i * 40), Array(9).fill(25));
  for (const branch of branches) {
    assert.equal(branchHue(branch, branches), branchHue(branch, [...branches].reverse()));
    assert.equal(branchHue(branch, branches), branchHue(branch, ["", ...branches, branch]));
  }
});


test("soon and max requirements stay disabled and render localized escaped labels without links", () => {
  for (const lang of ["ar", "en"]) {
    context.LANG = lang;
    for (const [kind, ar, en] of [["soon", "لسه بتتبني في تحديث جاي", "Coming in a later update"], ["max", "مكتمل", "Maxed"]]) {
      const check = { kind, ok: false, label: ar, labelEn: en };
      const node = { state: kind === "max" ? "owned" : "near", level: kind === "max" ? 18 : 0, checks: [check] };
      const markup = checkMarkup(check, 0);
      assert.equal(buyEnabled(node), false);
      assert.ok(markup.includes(lang === "en" ? en : ar));
      assert.match(markup, /is-missing/);
      assert.doesNotMatch(markup, /<button/);
      assert.equal(buyLabel(node), kind === "max" ? (lang === "en" ? "Maxed" : "مكتمل") : (lang === "en" ? "Build" : "ابني"));
      assert.match(checkMarkup({ ...check, label: "<unsafe>", labelEn: "<unsafe>" }, 0), /&lt;unsafe&gt;/);
    }
  }
});
