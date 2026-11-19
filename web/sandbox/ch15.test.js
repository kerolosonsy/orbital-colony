"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixture = require("./ch15.fixtures.json");
const { createSandboxHarness } = require("./harness");

test("chapter 15 algorithm trace renders and responds in Arabic and English", () => {
  for (const lang of ["ar", "en"]) {
    const harness = createSandboxHarness({ chapter: 15, lang });
    const panel = harness.render(fixture.kind, fixture.config, { id: fixture.id });
    assert.ok(panel);
    assert.ok(panel.textContent.includes(fixture.config.labels.next[lang]));
    assert.ok(panel.textContent.includes(fixture.config.labels.reset[lang]));
    assert.ok(panel.textContent.includes(fixture.config.frames[0].caption[lang]));
    assert.ok(panel.textContent.includes(fixture.config.frames[0].rows[0].label[lang]));

    const controls = panel.querySelectorAll("button");
    assert.equal(controls.length, 3);
    const [previous, next, reset] = controls;
    assert.equal(previous.disabled, true);
    const initial = panel.textContent;
    harness.click(next);
    const advanced = panel.textContent;
    assert.notEqual(advanced, initial);
    assert.equal(previous.disabled, false);
    harness.click(previous);
    assert.notEqual(panel.textContent, advanced);
    harness.click(next);
    const secondFrame = panel.textContent;
    harness.click(reset);
    assert.notEqual(panel.textContent, secondFrame);
    assert.equal(panel.textContent, initial);
  }
});
