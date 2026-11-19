"use strict";

const { test } = require("node:test");
const assert = require("node:assert/strict");
const fixtures = require("./ch17.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

test("chapter 17 playgrounds update state in Arabic and English", () => {
  for (const [kind, config] of Object.entries(fixtures)) {
    for (const lang of ["ar", "en"]) {
      const harness = createSandboxHarness({ chapter: 17, lang });
      const panel = harness.render(kind, config, { text: lang === "ar" ? "تجربة توضيحية." : "Illustrative try." });
      const state = panel.querySelector(".ch17-sandbox-state");
      const buttons = panel.querySelectorAll("button");
      assert.ok(state, `${kind} should render its state`);
      assert.equal(panel.getAttribute("dir"), lang === "ar" ? "rtl" : "ltr");
      assert.ok(buttons.length > 0 && buttons.length <= 4, `${kind} should use at most four buttons`);
      assert.equal(buttons.every((node) => node.tagName === "button"), true);
      assert.ok(buttons.some((node) => node.textContent === config.labels.action[lang]), `${kind} should localize its action`);

      for (const control of buttons) {
        const before = state.textContent;
        harness.click(control);
        assert.notEqual(state.textContent, before, `${kind} control should visibly change state`);
      }
      for (const caption of panel.querySelectorAll(".ch17-sandbox-caption")) {
        assert.ok(caption.textContent.trim().split(/\s+/).length <= 8, `${kind} caption should stay short`);
      }
    }
  }
});
