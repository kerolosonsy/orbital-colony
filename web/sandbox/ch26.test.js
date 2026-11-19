"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");
const { createSandboxHarness } = require("./harness.js");

const fixtures = JSON.parse(fs.readFileSync(path.join(__dirname, "ch26.fixtures.json"), "utf8"));

function panelFor(kind, lang) {
  const harness = createSandboxHarness({ chapter: 26, lang });
  const panel = harness.render(kind, fixtures[kind], { id: `c26-${kind}`, text: "fixture" });
  assert.ok(panel, `${kind} rendered in ${lang}`);
  for (const control of panel.querySelectorAll("button")) {
    assert.equal(control.tagName, "button");
    assert.equal(control.type, "button");
  }
  return { harness, panel };
}

function clickChanged(harness, panel, index) {
  const control = panel.querySelectorAll("button")[index];
  assert.ok(control, `button ${index} exists`);
  const before = panel.textContent;
  harness.click(control);
  assert.notEqual(panel.textContent, before, `button ${index} changes visible state`);
}

const languages = [
  ["ch26-handler-flow", "Next step", "الخطوة الجاية"],
  ["ch26-rest-atomic", "Records", "السجلات"],
  ["ch26-middleware-order", "Request enters outer.", "الطلب دخل outer."],
  ["ch26-error-envelope", "Error chain", "سلسلة الخطأ"],
  ["ch26-static-route", "Route", "المسار"],
  ["ch26-upload-stream", "Sample file", "الملف التوضيحي"],
];

test("chapter 26 fixtures cover every registered kind and both languages", () => {
  assert.deepEqual(Object.keys(fixtures).sort(), languages.map(([kind]) => kind).sort());
  for (const [kind, english, arabic] of languages) {
    const en = panelFor(kind, "en").panel.textContent;
    const ar = panelFor(kind, "ar").panel.textContent;
    assert.ok(en.includes(english), `${kind} has English text`);
    assert.ok(ar.includes(arabic), `${kind} has Arabic text`);
    assert.notEqual(en, ar, `${kind} changes language`);
  }
});

test("ch26-handler-flow sets metadata before body", () => {
  const { harness, panel } = panelFor("ch26-handler-flow", "en");
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 1);
  assert.match(panel.textContent, /207/);
  assert.match(panel.textContent, /sample/);
  clickChanged(harness, panel, 2);
});

test("ch26-rest-atomic allocates IDs with records in alternate orders", () => {
  const { harness, panel } = panelFor("ch26-rest-atomic", "en");
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 2);
});

test("ch26-middleware-order steps inward and unwinds outward", () => {
  const { harness, panel } = panelFor("ch26-middleware-order", "en");
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 2);
});

test("ch26-error-envelope follows a wrapped typed error", () => {
  const { harness, panel } = panelFor("ch26-error-envelope", "en");
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 1);
  assert.match(panel.textContent, /Found \*Fault/);
  clickChanged(harness, panel, 2);
});

test("ch26-static-route checks exact route entries", () => {
  const { harness, panel } = panelFor("ch26-static-route", "en");
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 1);
  assert.match(panel.textContent, /Denied/);
});

test("ch26-upload-stream copies bytes and reaches EOF", () => {
  const { harness, panel } = panelFor("ch26-upload-stream", "en");
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 0);
  assert.match(panel.textContent, /Reached EOF/);
  clickChanged(harness, panel, 1);
  assert.doesNotMatch(panel.textContent, /Reached EOF/);
});
