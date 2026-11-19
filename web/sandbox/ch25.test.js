"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");
const { createSandboxHarness } = require("./harness.js");

const fixtures = JSON.parse(fs.readFileSync(path.join(__dirname, "ch25.fixtures.json"), "utf8"));

function panelFor(kind, lang) {
  const harness = createSandboxHarness({ chapter: 25, lang });
  const panel = harness.render(kind, fixtures[kind], { id: `c25-${kind}`, text: "fixture" });
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
  ["ch25-locked-account", "Balance", "الرصيد"],
  ["ch25-race-report", "Counter", "العداد"],
  ["ch25-token-clock", "Elapsed seconds", "الوقت اللي عدى"],
  ["ch25-nonblocking-bus", "Buffer", "الـbuffer"],
  ["ch25-ttl-cache", "Advance clock", "قدّم الساعة"],
  ["ch25-atomic-counter", "Current Load", "Load الحالي"],
  ["ch25-closed-drain", "Job channel", "قناة الشغل"],
];

test("chapter 25 fixtures cover every registered kind and both languages", () => {
  assert.deepEqual(Object.keys(fixtures).sort(), languages.map(([kind]) => kind).sort());
  for (const [kind, english, arabic] of languages) {
    const en = panelFor(kind, "en").panel.textContent;
    const ar = panelFor(kind, "ar").panel.textContent;
    assert.ok(en.includes(english), `${kind} has English text`);
    assert.ok(ar.includes(arabic), `${kind} has Arabic text`);
    assert.notEqual(en, ar, `${kind} changes language`);
  }
});

test("ch25-locked-account steps through locked transactions and alternate order", () => {
  const { harness, panel } = panelFor("ch25-locked-account", "en");
  clickChanged(harness, panel, 0);
  assert.match(panel.textContent, /Balance3/);
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 0);
  assert.match(panel.textContent, /Worker B spends 5; left 2/);
  clickChanged(harness, panel, 2);
  assert.match(panel.textContent, /Balance7/);
});

test("ch25-race-report demonstrates two possible interleavings", () => {
  const { harness, panel } = panelFor("ch25-race-report", "en");
  for (let step = 0; step < 4; step += 1) harness.click(panel.querySelectorAll("button")[0]);
  assert.match(panel.textContent, /Counter15/);
  clickChanged(harness, panel, 1);
  for (let step = 0; step < 4; step += 1) harness.click(panel.querySelectorAll("button")[0]);
  assert.match(panel.textContent, /Counter13/);
  clickChanged(harness, panel, 2);
  assert.match(panel.textContent, /Counter10/);
});

test("ch25-token-clock recalculates complete intervals and permits", () => {
  for (const lang of ["en", "ar"]) {
    const { panel } = panelFor("ch25-token-clock", lang);
    const slider = panel.querySelector("input");
    assert.equal(slider.type, "range");
    const before = panel.textContent;
    slider.value = "16";
    slider.oninput();
    assert.notEqual(panel.textContent, before);
    const completeLabel = lang === "en" ? "Complete intervals" : "فترات كاملة";
    const permitsLabel = lang === "en" ? "Permits" : "التصاريح";
    assert.ok(panel.textContent.includes(`${completeLabel}3`));
    assert.ok(panel.textContent.includes(`${permitsLabel}4 / 6`));
  }
});

test("ch25-nonblocking-bus publishes, receives, and resets", () => {
  const { harness, panel } = panelFor("ch25-nonblocking-bus", "en");
  clickChanged(harness, panel, 0);
  assert.match(panel.textContent, /full; skipped/);
  clickChanged(harness, panel, 1);
  assert.match(panel.textContent, /Buffer: 0\/1/);
  clickChanged(harness, panel, 2);
  assert.match(panel.textContent, /ready/);
});

test("ch25-ttl-cache advances time, reads expiry, and resets", () => {
  const { harness, panel } = panelFor("ch25-ttl-cache", "en");
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 1);
  for (let step = 0; step < 4; step += 1) clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 1);
  assert.match(panel.textContent, /Expired/);
  clickChanged(harness, panel, 2);
});

test("ch25-atomic-counter applies adds in alternate orders", () => {
  const { harness, panel } = panelFor("ch25-atomic-counter", "en");
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 0);
  assert.match(panel.textContent, /Current Load6/);
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 0);
  assert.match(panel.textContent, /Current Load3/);
  clickChanged(harness, panel, 2);
  assert.match(panel.textContent, /Current Load0/);
});

test("ch25-closed-drain preserves queued values after close", () => {
  const { harness, panel } = panelFor("ch25-closed-drain", "en");
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 0);
  clickChanged(harness, panel, 1);
  clickChanged(harness, panel, 2);
  assert.match(panel.textContent, /Open/);
});
