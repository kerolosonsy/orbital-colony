"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { createSandboxHarness } = require("./harness.js");
const fixtures = require("./ch33.fixtures.json");

const levels = [
  "33-1-retry-client",
  "33-2-reverse-proxy",
  "33-3-load-balancer",
  "33-4-gateway",
  "33-boss-survey-uplink",
];

for (const lang of ["ar", "en"]) {
  for (const id of levels) {
    test(`${id} sandbox changes state in ${lang}`, () => {
      const config = fixtures[id];
      const harness = createSandboxHarness({ chapter: 33, lang });
      const panel = harness.render("ch33-http-flow", config, {
        id: `ch33-${id}`,
        text: lang === "ar" ? "تابع الحالات." : "Step through states.",
      });
      const firstRun = config.runs[0];
      assert.ok(panel.textContent.includes(config.title[lang]));
      assert.ok(panel.textContent.includes(firstRun.label[lang]));
      assert.ok(panel.textContent.includes(firstRun.frames[0].caption[lang]));
      const controls = panel.querySelectorAll("button");
      assert.equal(controls.length, config.runs.length > 1 ? 3 : 2);

      let before = panel.textContent;
      harness.click(controls[0]);
      assert.notEqual(panel.textContent, before, "advance changes the visible frame");

      if (config.runs.length > 1) {
        before = panel.textContent;
        harness.click(controls[1]);
        assert.notEqual(panel.textContent, before, "another schedule changes the visible state");
        assert.ok(panel.textContent.includes(config.runs[1].label[lang]));
        harness.click(controls[0]);
      }

      before = panel.textContent;
      harness.click(controls[controls.length - 1]);
      assert.notEqual(panel.textContent, before, "reset changes the visible state");
      assert.ok(panel.textContent.includes(firstRun.frames[0].caption[lang]));
    });
  }
}
