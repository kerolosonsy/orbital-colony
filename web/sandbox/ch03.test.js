"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { createSandboxHarness } = require("./harness.js");
const fixtures = require("./ch03.fixtures.json");

const firstText = {
  "ch03-repeat-counter": { ar: "العداد i", en: "Counter i" },
  "ch03-nested-grid": { ar: "الصف والعمود", en: "Row and column" },
  "ch03-condition-loop": { ar: "الشرط", en: "Condition" },
};

test("chapter 03 playgrounds update in Arabic and English", () => {
  for (const [kind, localizedConfigs] of Object.entries(fixtures)) {
    for (const lang of ["ar", "en"]) {
      const config = localizedConfigs[lang];
      const harness = createSandboxHarness({ chapter: 3, lang });
      const panel = harness.render(kind, config, { id: `ch03-${kind}` });
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
