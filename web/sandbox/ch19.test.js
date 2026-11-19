"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixtures = require("./ch19.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

const expectedLabels = {
  "ch19-duration-add": "start",
  "ch19-zone-view": "instant",
  "ch19-select-race": "schedule",
  "ch19-worker-order": "finishOrder",
};

function render(kind, config, lang) {
  const harness = createSandboxHarness({ chapter: 19, lang });
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

function changeSlider(panel, value) {
  const slider = panel.querySelector(".sandbox-kit-range");
  assert.ok(slider, "sandbox renders a slider");
  const before = panel.textContent;
  slider.value = String(value);
  slider.oninput();
  assert.notEqual(panel.textContent, before, "changing the slider changes visible state");
}

function exerciseStepper(harness, panel) {
  const [previous, next, reset] = panel.querySelectorAll("button");
  assert.ok(previous && next && reset, "stepper renders previous, next, and reset buttons");
  clickAndCheck(harness, panel, next);
  clickAndCheck(harness, panel, previous);
  clickAndCheck(harness, panel, next);
  clickAndCheck(harness, panel, reset);
}

test("chapter 19 sandboxes render in Arabic and English and respond to every control", () => {
  for (const [kind, arabicConfig] of Object.entries(fixtures)) {
    for (const lang of ["ar", "en"]) {
      const config = lang === "ar" ? arabicConfig : englishConfig(arabicConfig);
      const { harness, panel } = render(kind, config, lang);
      if (kind === "ch19-duration-add") changeSlider(panel, 45);
      else if (kind === "ch19-zone-view") changeSlider(panel, 2);
      else if (kind === "ch19-select-race") {
        const [runAgain, reset] = panel.querySelectorAll("button");
        clickAndCheck(harness, panel, runAgain);
        clickAndCheck(harness, panel, reset);
      } else if (kind === "ch19-worker-order") {
        const [runAgain] = panel.querySelectorAll("button");
        clickAndCheck(harness, panel, runAgain);
      } else {
        assert.fail(`No interaction test for ${kind}`);
      }
    }
  }
});
