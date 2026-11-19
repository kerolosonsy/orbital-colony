"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixtures = require("./ch10.fixtures.json");
const { createSandboxHarness } = require("./harness");

test("chapter 10 playgrounds render in Arabic and English and vary event order", () => {
  for (const [kind, config] of Object.entries(fixtures)) {
    for (const lang of ["ar", "en"]) {
      const harness = createSandboxHarness({ chapter: 10, lang });
      const panel = harness.render(kind, config, { id: `${kind}-test`, text: "" });
      assert.ok(panel, `${kind} should render in ${lang}`);
      assert.ok(panel.textContent.includes(config.intro[lang]), `${kind} intro should use ${lang}`);

      const controls = panel.querySelectorAll("button");
      assert.equal(controls.length, 1, `${kind} should have one control`);
      assert.equal(controls[0].textContent, config.button[lang]);
      const before = panel.textContent;
      harness.click(controls[0]);
      assert.notEqual(panel.textContent, before, `${kind} control should change visible state`);
      assert.ok(panel.textContent.includes(config.runs[0].events[0][lang]));

      for (let i = 1; i < config.runs[0].events.length; i++) harness.click(controls[0]);
      const firstRun = panel.textContent;
      harness.click(controls[0]);
      assert.notEqual(panel.textContent, firstRun, `${kind} should show another valid order`);
      assert.ok(panel.textContent.includes(config.runs[1].events[0][lang]));
    }
  }
});
