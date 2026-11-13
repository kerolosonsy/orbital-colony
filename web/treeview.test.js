"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const escapeHTML = (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
const context = vm.createContext({ module: { exports: {} }, LANG: "ar", T: (ar, en) => context.LANG === "en" ? en : ar, esc: escapeHTML });
vm.runInContext(fs.readFileSync(__dirname + "/treeview.js", "utf8"), context);
const { edgeSatisfied, missingCount, pathToNearestOwned, matchesSearch, buyEnabled, buyLabel, checkMarkup, branchHue } = context.module.exports;

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
