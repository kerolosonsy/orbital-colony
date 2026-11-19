"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { createSandboxHarness } = require("./harness.js");
const fixtures = require("./ch34.fixtures.json");

const levels = [
  "34-1-kv-store",
  "34-2-wal",
  "34-3-compaction",
  "34-4-chunk-store",
];

for (const lang of ["ar", "en"]) {
  for (const id of levels) {
    test(`${id} sandbox changes state in ${lang}`, () => {
      const config = fixtures[id];
      const harness = createSandboxHarness({ chapter: 34, lang });
      const panel = harness.render("ch34-storage-flow", config, {
        id: `ch34-${id}`,
        text: lang === "ar" ? "تابع الحالات." : "Step through states.",
      });
      assert.ok(panel.textContent.includes(config.title[lang]));
      assert.ok(panel.textContent.includes(config.frames[0].caption[lang]));
      const controls = panel.querySelectorAll("button");
      assert.equal(controls.length, 2);

      const beforeAdvance = panel.textContent;
      harness.click(controls[0]);
      assert.notEqual(panel.textContent, beforeAdvance, "advance changes the visible frame");
      assert.ok(panel.textContent.includes(config.frames[1].caption[lang]));

      const beforeReset = panel.textContent;
      harness.click(controls[1]);
      assert.notEqual(panel.textContent, beforeReset, "reset changes the visible frame");
      assert.ok(panel.textContent.includes(config.frames[0].caption[lang]));
    });
  }
}
