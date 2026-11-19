"use strict";

const assert = require("node:assert/strict");
const { test } = require("node:test");
const fixtures = require("./ch14.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

const cases = [
  {
    kind: "ch14-generic-max",
    exercise(harness, panel, config, lang) {
      const find = (key) => button(panel, config.labels[key][lang]);
      clickChanged(harness, panel, find("next"));
      clickChanged(harness, panel, find("previous"));
      clickChanged(harness, panel, find("next"));
      clickChanged(harness, panel, find("reset"));
    },
    visible: "type",
  },
  {
    kind: "ch14-stack",
    exercise(harness, panel, config, lang) {
      const find = (key) => button(panel, config.labels[key][lang]);
      clickChanged(harness, panel, find("push"));
      clickChanged(harness, panel, find("peek"));
      clickChanged(harness, panel, find("pop"));
      clickChanged(harness, panel, find("reset"));
    },
    visible: "generic",
  },
  {
    kind: "ch14-type-constraints",
    exercise(harness, panel, config, lang) {
      const find = (key) => button(panel, config.labels[key][lang]);
      const next = find("next");
      assert.ok(panel.textContent.includes(config.labels.accepted[lang]));
      clickChanged(harness, panel, next);
      assert.ok(panel.textContent.includes(config.labels.accepted[lang]));
      clickChanged(harness, panel, next);
      assert.ok(panel.textContent.includes(config.labels.rejected[lang]));
      clickChanged(harness, panel, find("reset"));
    },
    visible: "constraint",
  },
  {
    kind: "ch14-queue",
    exercise(harness, panel, config, lang) {
      const find = (key) => button(panel, config.labels[key][lang]);
      clickChanged(harness, panel, find("add"));
      clickChanged(harness, panel, find("read"));
      clickChanged(harness, panel, find("depart"));
      clickChanged(harness, panel, find("reset"));
    },
    visible: "queue",
  },
  {
    kind: "ch14-recursion",
    exercise(harness, panel, config, lang) {
      const find = (key) => button(panel, config.labels[key][lang]);
      clickChanged(harness, panel, find("next"));
      clickChanged(harness, panel, find("previous"));
      clickChanged(harness, panel, find("next"));
      clickChanged(harness, panel, find("reset"));
    },
    visible: "calls",
  },
];

function button(panel, label) {
  const found = panel.querySelectorAll("button").find((candidate) => candidate.textContent === label);
  assert.ok(found, `missing button: ${label}`);
  return found;
}

function clickChanged(harness, panel, target) {
  assert.equal(target.tagName, "button", "controls must be real buttons");
  const before = panel.textContent;
  harness.click(target);
  assert.notEqual(panel.textContent, before, `button did not change visible state: ${target.textContent}`);
}

for (const lang of ["ar", "en"]) {
  for (const item of cases) {
    test(`${item.kind} renders and responds in ${lang}`, () => {
      const harness = createSandboxHarness({ chapter: 14, lang });
      const panel = harness.render(item.kind, fixtures[item.kind], { id: `${item.kind}-${lang}` });
      assert.ok(panel, `missing renderer for ${item.kind}`);
      assert.ok(panel.textContent.includes(fixtures[item.kind].labels[item.visible][lang]));
      assert.ok(panel.querySelectorAll("button").every((node) => node.tagName === "button"));
      item.exercise(harness, panel, fixtures[item.kind], lang);
    });
  }
}
