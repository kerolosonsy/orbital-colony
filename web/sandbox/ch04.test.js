"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { createSandboxHarness } = require("./harness.js");
const fixtures = require("./ch04.fixtures.json");

const firstText = {
  "ch04-if-branch": { ar: "اختار حالة", en: "Choose a state" },
  "ch04-switch-case": { ar: "اختار قيمة", en: "Choose a value" },
  "ch04-ordered-branches": { ar: "ساعات التوصيل", en: "Delivery hours" },
};

test("chapter 04 playgrounds update in Arabic and English", () => {
  for (const [kind, localizedConfigs] of Object.entries(fixtures)) {
    for (const lang of ["ar", "en"]) {
      const config = localizedConfigs[lang];
      const harness = createSandboxHarness({ chapter: 4, lang });
      const panel = harness.render(kind, config, { id: `ch04-${kind}` });
      assert.ok(panel, `${kind} renders in ${lang}`);
      assert.ok(panel.textContent.includes(firstText[kind][lang]), `${kind} uses ${lang} config strings`);
      const controls = panel.querySelectorAll("button");
      assert.ok(controls.length > 0 && controls.length <= 4, `${kind} has a few real controls`);
      controls.forEach((control) => assert.equal(control.tagName, "button"));
      for (const control of controls) {
        const before = panel.textContent;
        harness.click(control);
        assert.notEqual(panel.textContent, before, `${kind} control changes visible state`);
      }
    }
  }
});
