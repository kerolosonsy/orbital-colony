"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixtures = require("./ch27.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

function clickAndCheck(harness, button, panel) {
  const before = panel.textContent;
  harness.click(button);
  assert.notEqual(panel.textContent, before, `"${button.textContent}" should change visible state`);
}

function checkCaptions(config) {
  function visit(value, key = "") {
    if (!value || typeof value !== "object") return;
    if (key === "caption") {
      for (const language of ["ar", "en"]) {
        assert.ok(value[language].trim().split(/\s+/).length <= 8, `${language} caption is too long`);
      }
    }
    for (const [childKey, child] of Object.entries(value)) visit(child, childKey);
  }
  visit(config);
}

function exerciseControls(harness, panel) {
  const buttons = panel.querySelectorAll("button");
  assert.ok(buttons.length > 0, "each playground has an interactive control");
  assert.ok(buttons.length <= 4, "playgrounds stay within the control limit");
  if (buttons.some((button) => button.disabled)) {
    const next = buttons.find((button) => !button.disabled);
    clickAndCheck(harness, next, panel);
  }
  for (const button of buttons) clickAndCheck(harness, button, panel);
}

test("chapter 27 playgrounds render and respond in Arabic and English", () => {
  assert.equal(Object.keys(fixtures).length, 5);
  for (const language of ["ar", "en"]) {
    const harness = createSandboxHarness({ chapter: 27, lang: language });
    for (const [kind, config] of Object.entries(fixtures)) {
      checkCaptions(config);
      const panel = harness.render(kind, config);
      assert.ok(panel, `${kind} should render`);
      const firstLabel = Object.values(config.labels)[0][language];
      assert.ok(panel.textContent.includes(firstLabel), `${kind} should render ${language} text`);
      exerciseControls(harness, panel);
    }
  }
});
