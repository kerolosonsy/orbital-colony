"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixtures = require("./ch35.fixtures.json");
const { createSandboxHarness } = require("./harness");

const kind = "ch35-event-trace";

for (const fixture of fixtures.levels) {
  test(`${fixture.id} renders its bilingual event trace`, () => {
    assert.equal(fixture.kind, kind);
    for (const lang of ["en", "ar"]) {
      const harness = createSandboxHarness({ chapter: 35, lang });
      const config = fixture.configs[lang];
      const panel = harness.render(kind, config, { id: `${fixture.id}-playground` });
      assert.ok(panel, `${fixture.id} has a registered renderer`);
      const buttons = panel.querySelectorAll("button");
      assert.equal(buttons.length, 3, "the playground stays within its control limit");
      assert.deepEqual(buttons.map((node) => node.textContent), [
        config.labels.advance[lang],
        config.labels.again[lang],
        config.labels.reset[lang],
      ]);
      assert.ok(panel.textContent.includes(config.variants[0][0].caption[lang]));
      for (const variant of config.variants) {
        for (const frame of variant) {
          assert.ok(frame.caption[lang].trim().split(/\s+/).length <= 8, "captions stay short");
        }
      }

      const initial = panel.textContent;
      harness.click(buttons[0]);
      const advanced = panel.textContent;
      assert.notEqual(advanced, initial, "advance changes the visible state");
      harness.click(buttons[1]);
      const alternate = panel.textContent;
      assert.notEqual(alternate, advanced, "alternate run changes the visible state");
      harness.click(buttons[2]);
      assert.notEqual(panel.textContent, alternate, "reset returns to the first state");
    }
  });
}
