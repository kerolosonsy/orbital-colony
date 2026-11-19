"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");
const { createSandboxHarness } = require("./harness.js");

const fixtures = JSON.parse(fs.readFileSync(path.join(__dirname, "ch12.fixtures.json"), "utf8"));
const buttonLabel = {
  "ch12-shared-counter": "next",
  "ch12-registry-lock": "next",
  "ch12-deadlock-order": "next",
};

function assertShortCaptions(value) {
  if (Array.isArray(value)) {
    value.forEach(assertShortCaptions);
  } else if (value && typeof value === "object") {
    if (value.caption?.ar && value.caption?.en) {
      for (const language of ["ar", "en"]) {
        assert.ok(value.caption[language].trim().split(/\s+/).length <= 8, `${language} caption is too long`);
      }
    }
    Object.values(value).forEach(assertShortCaptions);
  }
}

for (const [kind, config] of Object.entries(fixtures)) {
  assertShortCaptions(config);
  for (const lang of ["ar", "en"]) {
    test(`${kind} renders and responds in ${lang}`, () => {
      const harness = createSandboxHarness({ chapter: 12, lang });
      const panel = harness.render(kind, config, { id: `12-${kind}`, text: config.labels[buttonLabel[kind]][lang] });
      assert.ok(panel, "registered kind should render");
      assert.ok(panel.textContent.includes(config.labels[buttonLabel[kind]][lang]), "renders localized Config text");
      const controls = panel.querySelectorAll("button");
      assert.ok(controls.length > 0);
      for (const control of controls) {
        assert.equal(control.tagName, "button");
        const before = panel.textContent;
        harness.click(control);
        assert.notEqual(panel.textContent, before, `control ${control.textContent} should change visible state`);
      }
    });
  }
}
