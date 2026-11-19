"use strict";

const { test } = require("node:test");
const assert = require("node:assert/strict");
const fixtures = require("./ch18.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

test("chapter 18 playgrounds update state in Arabic and English", () => {
  for (const [kind, config] of Object.entries(fixtures)) {
    for (const lang of ["ar", "en"]) {
      const harness = createSandboxHarness({ chapter: 18, lang });
      const panel = harness.render(kind, config, { text: lang === "ar" ? "تجربة توضيحية." : "Illustrative try." });
      const state = panel.querySelector(".ch18-sandbox-state");
      const buttons = panel.querySelectorAll("button");
      const actionLabels = config.mode === "unicode"
        ? [config.labels.split[lang], config.labels.lower[lang]]
        : [config.labels.action[lang]];
      assert.ok(state, `${kind} should render its state`);
      assert.equal(panel.getAttribute("dir"), lang === "ar" ? "rtl" : "ltr");
      assert.ok(buttons.length > 0 && buttons.length <= 4, `${kind} should use at most four buttons`);
      assert.equal(buttons.every((node) => node.tagName === "button"), true);
      for (const actionLabel of actionLabels) {
        assert.ok(buttons.some((node) => node.textContent === actionLabel), `${kind} should localize its controls`);
      }

      for (const control of buttons) {
        const before = state.textContent;
        harness.click(control);
        assert.notEqual(state.textContent, before, `${kind} control should visibly change state`);
      }
      for (const caption of panel.querySelectorAll(".ch18-sandbox-caption")) {
        assert.ok(caption.textContent.trim().split(/\s+/).length <= 8, `${kind} caption should stay short`);
      }
    }
  }
});
