"use strict";

const { test } = require("node:test");
const assert = require("node:assert/strict");
const fixtures = require("./ch07.fixtures.json");
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

test("chapter 7 playgrounds render and respond in Arabic and English", () => {
  const cases = [
    {
      kind: "ch07-point-fields",
      visible: (config) => config.labels.object,
      controls: 2,
      exercise(harness, box) {
        const controls = buttons(box);
        clickAndCheck(harness, box, controls[0], "3");
        clickAndCheck(harness, box, controls[1], "1");
      },
    },
    {
      kind: "ch07-nav-pointer",
      visible: (config) => config.labels.position,
      controls: 2,
      exercise(harness, box) {
        const controls = buttons(box);
        clickAndCheck(harness, box, controls[0], "3");
        clickAndCheck(harness, box, controls[1], "2");
      },
    },
    {
      kind: "ch07-manhattan",
      visible: (config) => config.labels.total,
      controls: 3,
      exercise(harness, box) {
        const controls = buttons(box);
        clickAndCheck(harness, box, controls[0]);
        assert.deepEqual(box.querySelectorAll(".sandbox-kit-value").map((node) => node.textContent), ["2", "2", "4"]);
        clickAndCheck(harness, box, controls[1]);
        assert.deepEqual(box.querySelectorAll(".sandbox-kit-value").map((node) => node.textContent), ["2", "3", "5"]);
        clickAndCheck(harness, box, controls[2]);
        assert.deepEqual(box.querySelectorAll(".sandbox-kit-value").map((node) => node.textContent), ["1", "2", "3"]);
      },
    },
    {
      kind: "ch07-embedded-pack",
      visible: (config) => config.frames[0].value,
      controls: 3,
      exercise(harness, box) {
        const controls = buttons(box);
        clickAndCheck(harness, box, controls[1]);
        clickAndCheck(harness, box, controls[0]);
        clickAndCheck(harness, box, controls[1]);
        clickAndCheck(harness, box, controls[2]);
      },
    },
    {
      kind: "ch07-range-questions",
      visible: (config) => config.questions[0].prompt,
      controls: 3,
      exercise(harness, box) {
        const controls = buttons(box);
        clickAndCheck(harness, box, controls[1]);
        clickAndCheck(harness, box, controls[0]);
        clickAndCheck(harness, box, controls[1]);
        clickAndCheck(harness, box, controls[2]);
      },
    },
  ];

  for (const lang of ["ar", "en"]) {
    const harness = createSandboxHarness({ chapter: 7, lang });
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
