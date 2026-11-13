"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const escapeHTML = (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
const context = vm.createContext({ module: { exports: {} }, LANG: "ar", T: (ar, en) => context.LANG === "en" ? en : ar, esc: escapeHTML });
vm.runInContext(fs.readFileSync(__dirname + "/treeview.js", "utf8"), context);
const { fitView, edgeSatisfied, missingCount, pathToNearestOwned, matchesSearch, buyEnabled, buyLabel, stateLabel, checkLabel, checkMarkup, costMarkup, headerStatsMarkup, branchHue } = context.module.exports;

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

test("rp and rank checks render clear localized labels without clickable links", () => {
  for (const lang of ["ar", "en"]) {
    context.LANG = lang;
    const rpWithLabels = checkMarkup({ kind: "rp", ok: true, label: "٤٠ نقطة بحث", labelEn: "40 research points" }, 0);
    assert.match(rpWithLabels, /is-ok/);
    assert.doesNotMatch(rpWithLabels, /<button/);
    assert.ok(rpWithLabels.includes(lang === "en" ? "40 research points" : "٤٠ نقطة بحث"));

    const rpNumericFallback = checkMarkup({ kind: "rp", ok: false }, 1, { rp: 40 });
    assert.match(rpNumericFallback, /is-missing/);
    assert.doesNotMatch(rpNumericFallback, /<button/);
    assert.ok(rpNumericFallback.includes(lang === "en" ? "Research points: 40" : "نقاط البحث: 40"));

    const rpBareFallback = checkMarkup({ kind: "rp", ok: false }, 2);
    assert.ok(rpBareFallback.includes(lang === "en" ? "Enough research points" : "نقاط بحث كافية"));

    const rankFromRef = checkMarkup({ kind: "rank", ok: false, ref: "3" }, 3);
    assert.match(rankFromRef, /is-missing/);
    assert.doesNotMatch(rankFromRef, /<button/);
    assert.ok(rankFromRef.includes(lang === "en" ? "Rank 3 or higher" : "المرتبة 3 أو أعلى"));

    const rankBareFallback = checkMarkup({ kind: "rank", ok: true }, 4);
    assert.match(rankBareFallback, /is-ok/);
    assert.ok(rankBareFallback.includes(lang === "en" ? "Colony rank" : "مرتبة المستعمرة"));
  }
});

test("costMarkup places RP cost next to credit cost only when rp > 0 and reflects affordability", () => {
  for (const lang of ["ar", "en"]) {
    context.LANG = lang;
    const nodeNoRp = { cost: { credits: 100, goods: {} } };
    const noRpHtml = costMarkup(nodeNoRp, { credits: 120, rp: 50, inv: {} });
    assert.doesNotMatch(noRpHtml, lang === "en" ? />RP</ : />نقاط البحث</);

    const nodeZeroRp = { rp: 0, cost: { credits: 100, goods: {} } };
    assert.doesNotMatch(costMarkup(nodeZeroRp, { credits: 120, rp: 50, inv: {} }), lang === "en" ? />RP</ : />نقاط البحث</);

    const nodeWithRp = { rp: 40, cost: { credits: 100, goods: { moss: 5 } } };
    const missingRpHtml = costMarkup(nodeWithRp, { credits: 120, rp: 25, inv: { moss: 5 } });
    assert.match(missingRpHtml, lang === "en" ? /Credits.*120 \/ 100.*RP.*25 \/ 40/ : /كريدت.*120 \/ 100.*نقاط البحث.*25 \/ 40/);
    assert.match(missingRpHtml, /class="is-missing"><span>(RP|نقاط البحث)<\/span>/);

    const metRpHtml = costMarkup(nodeWithRp, { credits: 120, rp: 40, inv: { moss: 5 } });
    assert.match(metRpHtml, /class="is-ok"><span>(RP|نقاط البحث)<\/span>/);
  }
});

test("tree header stats render current RP and rank when present and omit absent fields", () => {
  for (const lang of ["ar", "en"]) {
    context.LANG = lang;
    assert.equal(headerStatsMarkup({}), "");
    assert.equal(headerStatsMarkup({ rp: 0 }), lang === "en" ? '<span class="tv-stat tv-stat-rp"><span>RP</span><b dir="ltr">0</b></span>' : '<span class="tv-stat tv-stat-rp"><span>نقاط البحث</span><b dir="ltr">0</b></span>');
    const both = headerStatsMarkup({ rank: 3, rp: 45 });
    assert.match(both, lang === "en" ? /Rank.*3.*RP.*45/ : /المرتبة.*3.*نقاط البحث.*45/);
  }
});

test("tree overlay updates header stats and detail card live from TreeView.update", () => {
  const element = () => ({
    attributes: new Map(), events: {}, classes: new Set(), innerHTML: "", hidden: false,
    setAttribute(name, value) { this.attributes.set(name, value); },
    addEventListener(name, listener) { this.events[name] = listener; },
    closest() { return null; },
  });
  const dialog = element(), canvas = element(), world = element(), stats = element();
  const centre = element(), input = element(), detail = element(), nav = element(), close = element();
  const selectors = { ".tv-canvas": canvas, ".tv-world": world, ".tv-detail": detail, ".tv-branches": nav, ".tv-stats": stats, input, '[data-action="centre"]': centre, '[data-action="close"]': close };
  dialog.querySelector = (selector) => selectors[selector];
  dialog.querySelectorAll = () => [];
  dialog.showModal = () => { dialog.open = true; };
  dialog.classList = { toggle(name, active) { if (active) dialog.classes.add(name); else dialog.classes.delete(name); } };
  canvas.getBoundingClientRect = () => ({ width: 800, height: 600 });
  const ui = vm.createContext({
    document: { activeElement: null, createElement: () => dialog, body: { appendChild() {} } },
    ResizeObserver: class { observe() {} },
    LANG: "en", T: (_ar, en) => en, esc: escapeHTML,
    NAME: { moss: "Moss" },
  });
  vm.runInContext(fs.readFileSync(__dirname + "/tree.js", "utf8"), ui);
  vm.runInContext(fs.readFileSync(__dirname + "/treeview.js", "utf8"), ui);
  const view = vm.runInContext("TreeView", ui);
  const nodes = [
    { id: "root", branch: "", after: [], state: "owned", level: 1, name: "Base", desc: "D", effect: "E", checks: [], cost: { credits: 0, goods: {} } },
    { id: "tech", branch: "research", after: [{ id: "root", level: 1 }], state: "near", level: 0, name: "Lab Tech", desc: "D", effect: "E", rp: 40, checks: [{ kind: "rank", ok: false, ref: "3" }, { kind: "rp", ok: false }], cost: { credits: 200, goods: {} } },
  ];
  view.open({ nodes, credits: 500, inv: {} });
  assert.equal(stats.hidden, true);
  assert.equal(stats.innerHTML, "");

  view.update({ nodes, credits: 500, inv: {}, rp: 25, rank: 2 });
  assert.equal(stats.hidden, false);
  assert.match(stats.innerHTML, /Rank.*2.*RP.*25/);

  // Selecting node via click renders rank + rp checks and RP cost in detail card
  canvas.events.click({ detail: 0, target: { closest: () => ({ dataset: { node: "tech" } }) } });
  assert.match(detail.innerHTML, /Rank 3 or higher/);
  assert.match(detail.innerHTML, /Research points: 40/);
  assert.match(detail.innerHTML, /25 \/ 40/);
});

test("max 0 endless tree nodes are never maxed and render plain numeric levels", () => {
  for (const lang of ["ar", "en"]) {
    context.LANG = lang;
    const en = lang === "en";
    const endlessZero = {
      id: "endless", max: 0, level: 0, state: "available",
      checks: [{ kind: "max", ok: false }, { kind: "credits", ok: true }],
    };
    assert.equal(buyEnabled(endlessZero), true);
    assert.equal(buyLabel(endlessZero), en ? "Build" : "ابني");
    assert.equal(stateLabel(endlessZero), en ? "Available" : "متاحة");
    assert.equal(missingCount(endlessZero), 0);

    const endlessHigh = {
      id: "endless", max: 0, level: 7, state: "owned",
      checks: [{ kind: "max", ok: false }, { kind: "credits", ok: true }],
    };
    assert.equal(buyEnabled(endlessHigh), true);
    assert.equal(buyLabel(endlessHigh), en ? "Upgrade" : "رقّي");
    assert.equal(stateLabel(endlessHigh), en ? "Level 7" : "مستوى 7");

    const endlessBroke = {
      id: "endless", max: 0, level: 7, state: "owned",
      checks: [{ kind: "credits", ok: false }],
    };
    assert.equal(buyEnabled(endlessBroke), false);
    assert.equal(buyLabel(endlessBroke), en ? "Upgrade" : "رقّي");
  }
});

test("costMarkup renders colour chips for crafted and raw goods in costs", () => {
  context.LANG = "en";
  context.NAME = { plate: "Plate", moss: "Moss" };
  context.INV_COLOR = { plate: "#c6d0de", moss: "#93c45c" };
  const html = costMarkup(
    { cost: { credits: 100, goods: { plate: 5, moss: 10 } } },
    { credits: 200, rp: 0, inv: { plate: 5, moss: 2 } },
  );
  assert.match(html, /tv-good-chip[^>]*background:#c6d0de/);
  assert.match(html, /tv-good-chip[^>]*background:#93c45c/);
  assert.match(html, /class="is-ok"><span class="tv-cost-good"><i class="tv-good-chip"[^>]*><\/i>Plate<\/span><b dir="ltr">5 \/ 5<\/b>/);
  assert.match(html, /class="is-missing"><span class="tv-cost-good"><i class="tv-good-chip"[^>]*><\/i>Moss<\/span><b dir="ltr">2 \/ 10<\/b>/);
});
