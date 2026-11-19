"use strict";

const { test } = require("node:test");
const assert = require("node:assert/strict");
const fixtures = require("./ch08.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

function text(pair, lang) {
  return pair[lang];
}

function buttons(box) {
  return box.querySelectorAll("button");
}

function clickAndCheck(harness, box, control, expected) {
  const before = box.textContent;
  harness.click(control);
  const after = box.textContent;
  assert.notEqual(after, before);
  if (expected !== undefined) assert.ok(after.includes(String(expected)), after);
}

test("chapter 8 playgrounds render and respond in Arabic and English", () => {
  const cases = [
    {
      kind: "ch08-copy-pointer",
      visible: (config) => config.labels.caller,
      controls: 3,
      exercise(harness, box) {
        const controls = buttons(box);
        clickAndCheck(harness, box, controls[0], "10");
        clickAndCheck(harness, box, controls[1], "10");
        clickAndCheck(harness, box, controls[2], "3");
      },
    },
    {
      kind: "ch08-receivers",
      visible: (config) => config.labels.caller,
      controls: 3,
      exercise(harness, box) {
        const controls = buttons(box);
        clickAndCheck(harness, box, controls[0], "5");
        clickAndCheck(harness, box, controls[1], "5");
        clickAndCheck(harness, box, controls[2], "2");
      },
    },
    {
      kind: "ch08-nil-guard",
      visible: (config) => config.labels.target,
      controls: 1,
      exercise(harness, box) {
        clickAndCheck(harness, box, buttons(box)[0], "6");
        assert.equal(box.querySelector(".sandbox-kit-value").textContent, "p = &6");
      },
    },
    {
      kind: "ch08-cart-chain",
      visible: (config) => config.labels.chain,
      controls: 2,
      exercise(harness, box) {
        const controls = buttons(box);
        clickAndCheck(harness, box, controls[0], "green");
        assert.equal(box.querySelectorAll(".ch08-node").filter((node) => node.classList.contains("current")).length, 1);
        clickAndCheck(harness, box, controls[1], "red");
        assert.equal(box.querySelector(".sandbox-kit-value").textContent, "red");
        assert.equal(box.querySelectorAll(".ch08-node").filter((node) => node.classList.contains("current")).length, 1);
      },
    },
  ];

  for (const lang of ["ar", "en"]) {
    const harness = createSandboxHarness({ chapter: 8, lang });
    for (const sandbox of cases) {
      const config = fixtures[sandbox.kind];
      const box = harness.render(sandbox.kind, config, {
        id: sandbox.kind,
        text: lang === "ar" ? "تجربة توضيحية." : "Illustrative practice.",
      });
      assert.ok(box, sandbox.kind + " did not render");
      assert.ok(box.textContent.includes(text(sandbox.visible(config), lang)), sandbox.kind + " did not render " + lang);
      assert.equal(buttons(box).length, sandbox.controls, sandbox.kind + " control count");
      sandbox.exercise(harness, box);
    }
  }
});
