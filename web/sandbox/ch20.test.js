"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixtures = require("./ch20.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

const expectedLabels = {
  "ch20-bulletin-state": "source",
  "ch20-error-join": "causes",
  "ch20-lexer-offset": "source",
  "ch20-expression-tree": "expression",
};

function render(kind, config, lang) {
  const harness = createSandboxHarness({ chapter: 20, lang });
  const panel = harness.render(kind, config, { text: "sandbox" });
  assert.ok(panel, `${kind} should render in ${lang}`);
  assert.ok(panel.textContent.includes(config.labels[expectedLabels[kind]][lang]), `${kind} should show its ${lang} label`);
  return { harness, panel };
}

function englishConfig(value) {
  if (Array.isArray(value)) return value.map(englishConfig);
  if (value && typeof value === "object") {
    if (typeof value.ar === "string" && typeof value.en === "string") return { ar: value.en, en: value.en };
    return Object.fromEntries(Object.entries(value).map(([key, item]) => [key, englishConfig(item)]));
  }
  return value;
}

function clickAndCheck(harness, panel, control) {
  assert.equal(control.tagName, "button", "sandbox controls use real buttons");
  const before = panel.textContent;
  harness.click(control);
  assert.notEqual(panel.textContent, before, "clicking the control changes visible state");
}

function exerciseStepper(harness, panel) {
  const [previous, next, reset] = panel.querySelectorAll("button");
  assert.ok(previous && next && reset, "stepper renders previous, next, and reset buttons");
  clickAndCheck(harness, panel, next);
  clickAndCheck(harness, panel, previous);
  clickAndCheck(harness, panel, next);
  clickAndCheck(harness, panel, reset);
}

test("chapter 20 sandboxes render in Arabic and English and respond to every control", () => {
  for (const [kind, arabicConfig] of Object.entries(fixtures)) {
    for (const lang of ["ar", "en"]) {
      const config = lang === "ar" ? arabicConfig : englishConfig(arabicConfig);
      const { harness, panel } = render(kind, config, lang);
      if (kind === "ch20-bulletin-state" || kind === "ch20-lexer-offset" || kind === "ch20-expression-tree") {
        exerciseStepper(harness, panel);
      } else if (kind === "ch20-error-join") {
        const controls = panel.querySelectorAll("button");
        assert.equal(controls.length, config.fields.length, "each field has one toggle button");
        controls.forEach((control) => clickAndCheck(harness, panel, control));
      } else {
        assert.fail(`No interaction test for ${kind}`);
      }
    }
  }
});
