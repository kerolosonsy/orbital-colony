"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { createSandboxHarness } = require("./harness.js");
const fixtures = require("./ch05.fixtures.json");

function render(kind, lang) {
  const config = fixtures[kind];
  const harness = createSandboxHarness({ chapter: 5, lang });
  const panel = harness.render(kind, config, { id: `c05-${kind}`, text: config.source[lang] });
  assert.ok(panel, `${kind} did not render in ${lang}`);
  assert.ok(panel.textContent.includes(config.source[lang]), `${kind} did not localize its source label in ${lang}`);
  const controls = panel.querySelectorAll("button");
  assert.ok(controls.length <= 4, `${kind} has more than four controls`);
  return { config, harness, panel, controls };
}

function assertChanged(panel, action, kind, controlName) {
  const before = panel.textContent;
  action();
  assert.notEqual(panel.textContent, before, `${kind}: ${controlName} did not change visible state`);
}

function exerciseStepper(kind, lang) {
  const { config, harness, panel, controls } = render(kind, lang);
  assert.equal(controls.length, 3, `${kind} should have previous, next, and reset`);
  const [previous, next, reset] = controls;
  assert.ok(previous.disabled, `${kind} previous should be disabled at the first frame`);
  assertChanged(panel, () => harness.click(next), kind, "next");
  assert.equal(previous.disabled, false);
  assertChanged(panel, () => harness.click(previous), kind, "previous");
  assertChanged(panel, () => harness.click(next), kind, "next again");
  assertChanged(panel, () => harness.click(reset), kind, "reset");
  assert.ok(panel.textContent.includes(config.frames[0].caption[lang]));
}

function exerciseToggle(kind, lang) {
  const { harness, panel, controls } = render(kind, lang);
  assert.equal(controls.length, 1, `${kind} should have one toggle`);
  assertChanged(panel, () => harness.click(controls[0]), kind, "toggle");
}

for (const lang of ["ar", "en"]) {
  test(`ch05 parameter playground responds in ${lang}`, () => exerciseStepper("ch05-params", lang));
  test(`ch05 error value playground responds in ${lang}`, () => exerciseToggle("ch05-error-value", lang));
  test(`ch05 error wrapping playground responds in ${lang}`, () => exerciseToggle("ch05-wrap", lang));
  test(`ch05 defer playground responds in ${lang}`, () => exerciseStepper("ch05-defer", lang));

  test(`ch05 closure playground responds in ${lang}`, () => {
    const { config, harness, panel, controls } = render("ch05-closure", lang);
    assert.equal(controls.length, 3, "closure should have store, call, and reset");
    const [store, call, reset] = controls;
    assert.ok(call.disabled, "call should wait for a stored function value");
    assertChanged(panel, () => harness.click(store), "ch05-closure", "store function");
    assert.equal(call.disabled, false);
    assertChanged(panel, () => harness.click(call), "ch05-closure", "call function");
    assert.ok(panel.textContent.includes(String(config.initialCount + 1)));
    assertChanged(panel, () => harness.click(reset), "ch05-closure", "reset");
    assert.ok(panel.textContent.includes(config.states.idle[lang]));
  });
}
