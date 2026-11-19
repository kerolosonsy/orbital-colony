"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");
const { createSandboxHarness } = require("./harness");

const fixtures = JSON.parse(fs.readFileSync(path.join(__dirname, "ch23.fixtures.json"), "utf8"));

test("chapter 23 sandboxes render and step in Arabic and English", () => {
  for (const [kind, config] of Object.entries(fixtures)) {
    for (const language of ["ar", "en"]) {
      const harness = createSandboxHarness({ chapter: 23, lang: language });
      const panel = harness.render(kind, config, { id: kind, text: "playground" });
      assert.ok(panel, `${kind} should render in ${language}`);
      assert.equal(panel.getAttribute("dir"), language === "ar" ? "rtl" : "ltr");
      assert.ok(panel.textContent.includes(config.frames[0].caption[language]));
      assert.ok(panel.textContent.includes(config.action[language]));
      for (const row of config.frames[0].rows) assert.ok(panel.textContent.includes(row.label[language]));
      const controls = panel.querySelectorAll("button");
      assert.equal(controls.length, 1, `${kind} should have one button`);
      for (let index = 1; index < config.frames.length; index += 1) {
        const before = panel.textContent;
        harness.click(controls[0]);
        assert.notEqual(panel.textContent, before, `${kind} should change after a click`);
        assert.ok(panel.textContent.includes(config.frames[index].caption[language]));
        for (const row of config.frames[index].rows) assert.ok(panel.textContent.includes(row.label[language]));
      }
      for (const frame of config.frames) {
        for (const caption of Object.values(frame.caption)) assert.ok(caption.trim().split(/\s+/).length <= 8);
      }
    }
  }
});
