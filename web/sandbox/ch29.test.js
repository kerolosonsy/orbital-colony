"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");
const { createSandboxHarness } = require("./harness.js");

const fixtures = JSON.parse(fs.readFileSync(path.join(__dirname, "ch29.fixtures.json"), "utf8"));
const controlOrder = {
  "ch29-migration-boundary": [0, 1, 2],
  "ch29-pool-state": [0, 2, 1],
};

for (const [levelID, fixture] of Object.entries(fixtures)) {
  for (const lang of ["ar", "en"]) {
    test(`${levelID} ${lang} sandbox controls`, () => {
      const harness = createSandboxHarness({ chapter: 29, lang });
      const sandbox = harness.render(fixture.kind, fixture.config, {
        id: fixture.id,
        text: fixture.config.intro[lang],
      });
      assert.ok(sandbox, `${fixture.kind} rendered`);
      assert.ok(harness.article.textContent.includes(fixture.config.intro[lang]), "localized intro is visible");
      const localizedLabel = Object.values(fixture.config.labels)[0][lang];
      assert.ok(sandbox.textContent.includes(localizedLabel), "localized Config label is visible");

      const controls = sandbox.querySelectorAll("button");
      assert.ok(controls.length > 0 && controls.length <= 4, "uses one to four real buttons");
      controls.forEach((control) => assert.equal(control.tagName, "button"));
      const order = controlOrder[fixture.kind] || controls.map((_, index) => index);
      assert.equal(order.length, controls.length, "every control is clicked once");
      for (const index of order) {
        const before = sandbox.textContent;
        harness.click(controls[index]);
        assert.notEqual(sandbox.textContent, before, `${fixture.kind} control ${index} changes visible state`);
      }
    });
  }
}
