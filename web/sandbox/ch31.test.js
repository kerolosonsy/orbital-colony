"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixtures = require("./ch31.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

test("chapter 31 playgrounds render and respond in Arabic and English", () => {
  for (const [levelID, fixture] of Object.entries(fixtures)) {
    for (const lang of ["ar", "en"]) {
      const harness = createSandboxHarness({ chapter: 31, lang });
      const panel = harness.render(fixture.kind, fixture.config, { id: levelID });
      assert.ok(panel, `${levelID} should render in ${lang}`);
      assert.ok(panel.textContent.includes(fixture.config.labels.title[lang]), `${levelID} should use ${lang} text`);
      const buttons = panel.querySelectorAll("button");
      assert.ok(buttons.length > 0 && buttons.length <= 4, `${levelID} should have one to four controls`);
      for (const control of buttons) {
        assert.equal(control.tagName, "button", `${levelID} controls must be real buttons`);
        const before = panel.textContent;
        harness.click(control);
        assert.notEqual(panel.textContent, before, `${levelID} control should change visible state in ${lang}`);
      }
    }
  }
});
