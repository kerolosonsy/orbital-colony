"use strict";

const assert = require("node:assert/strict");
const { test } = require("node:test");
const fixtures = require("./ch21.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

function render(kind, lang) {
  const harness = createSandboxHarness({ chapter: 21, lang });
  const panel = harness.render(kind, fixtures[kind], { id: `ch21-${kind}` });
  assert.ok(panel, `${kind} should render in ${lang}`);
  return { harness, panel };
}

function clickAndCheck(harness, panel, control) {
  const before = panel.textContent;
  harness.click(control);
  assert.notEqual(panel.textContent, before, "each control click should change visible state");
}

function exerciseLanguages(kind, key, exercise) {
  const config = fixtures[kind];
  for (const lang of ["ar", "en"]) {
    const { harness, panel } = render(kind, lang);
    assert.ok(panel.textContent.includes(config.labels[key][lang]), `${kind} should show its ${lang} label`);
    exercise(harness, panel);
  }
}

test("ch21-flags-env simulates flag and environment presence", () => {
  exerciseLanguages("ch21-flags-env", "toggle", (harness, panel) => {
    const controls = panel.querySelectorAll("button");
    assert.equal(controls.length, 2);
    clickAndCheck(harness, panel, controls[0]);
    clickAndCheck(harness, panel, controls[1]);
    clickAndCheck(harness, panel, controls[1]);
  });
});

test("ch21-reader-lines advances a scanner and writer", () => {
  exerciseLanguages("ch21-reader-lines", "next", (harness, panel) => {
    const [previous, next, reset] = panel.querySelectorAll("button");
    clickAndCheck(harness, panel, next);
    clickAndCheck(harness, panel, previous);
    clickAndCheck(harness, panel, next);
    clickAndCheck(harness, panel, reset);
  });
});

test("ch21-injected-draw compares repeatable and varied sources", () => {
  exerciseLanguages("ch21-injected-draw", "source", (harness, panel) => {
    const [source, draw] = panel.querySelectorAll("button");
    clickAndCheck(harness, panel, source);
    clickAndCheck(harness, panel, draw);
    clickAndCheck(harness, panel, draw);
  });
});

test("ch21-maneuver-values advances named enum examples", () => {
  exerciseLanguages("ch21-maneuver-values", "next", (harness, panel) => {
    const [next] = panel.querySelectorAll("button");
    clickAndCheck(harness, panel, next);
  });
});
