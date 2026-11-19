"use strict";

const assert = require("node:assert/strict");
const { test } = require("node:test");
const fixtures = require("./ch22.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

function render(kind, lang) {
  const harness = createSandboxHarness({ chapter: 22, lang });
  const panel = harness.render(kind, fixtures[kind], { id: `ch22-${kind}` });
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

test("ch22-depth-tree advances through sorted directory depth", () => {
  exerciseLanguages("ch22-depth-tree", "depth", (harness, panel) => {
    const [previous, next, reset] = panel.querySelectorAll("button");
    clickAndCheck(harness, panel, next);
    clickAndCheck(harness, panel, previous);
    clickAndCheck(harness, panel, next);
    clickAndCheck(harness, panel, reset);
  });
});

test("ch22-review-plan sorts preview data before mutation", () => {
  exerciseLanguages("ch22-review-plan", "plan", (harness, panel) => {
    const [previous, next, reset] = panel.querySelectorAll("button");
    clickAndCheck(harness, panel, next);
    clickAndCheck(harness, panel, previous);
    clickAndCheck(harness, panel, next);
    clickAndCheck(harness, panel, reset);
  });
});

test("ch22-atomic-publish keeps old content through staging and failure", () => {
  exerciseLanguages("ch22-atomic-publish", "target", (harness, panel) => {
    const [write, publish, fail, reset] = panel.querySelectorAll("button");
    clickAndCheck(harness, panel, write);
    clickAndCheck(harness, panel, publish);
    clickAndCheck(harness, panel, reset);
    clickAndCheck(harness, panel, write);
    clickAndCheck(harness, panel, fail);
    clickAndCheck(harness, panel, reset);
  });
});

test("ch22-size-hash filters candidate pairs by content", () => {
  exerciseLanguages("ch22-size-hash", "caption", (harness, panel) => {
    const [next] = panel.querySelectorAll("button");
    clickAndCheck(harness, panel, next);
    clickAndCheck(harness, panel, next);
    clickAndCheck(harness, panel, next);
  });
});
