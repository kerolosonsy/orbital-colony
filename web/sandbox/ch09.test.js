"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixtures = require("./ch09.fixtures.json");
const { createSandboxHarness } = require("./harness");

test("chapter 09 playgrounds render in Arabic and English and respond to their controls", () => {
  const firstControlLabels = {
    "ch09-stringer": "showString",
    "ch09-plan-callback": "button",
    "ch09-type-switch": "button",
    "ch09-trim-space": "trim",
    "ch09-map-updates": "button",
  };
  const nextControlLabels = {
    "ch09-stringer": "showDefault",
    "ch09-trim-space": "restore",
  };
  for (const [kind, config] of Object.entries(fixtures)) {
    for (const lang of ["ar", "en"]) {
      const harness = createSandboxHarness({ chapter: 9, lang });
      const panel = harness.render(kind, config, { id: `${kind}-test`, text: "" });
      assert.ok(panel, `${kind} should render in ${lang}`);
      assert.ok(panel.textContent.includes(config.intro[lang]), `${kind} intro should use ${lang}`);

      const controls = panel.querySelectorAll("button");
      assert.equal(controls.length, 1, `${kind} should have one control`);
      assert.equal(controls[0].textContent, config[firstControlLabels[kind]][lang]);
      const before = panel.textContent;
      harness.click(controls[0]);
      assert.notEqual(panel.textContent, before, `${kind} control should change visible state`);
      if (nextControlLabels[kind]) assert.equal(controls[0].textContent, config[nextControlLabels[kind]][lang]);
    }
  }
});
