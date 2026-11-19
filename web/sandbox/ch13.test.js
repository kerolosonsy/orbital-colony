"use strict";

const assert = require("node:assert/strict");
const { test } = require("node:test");
const fixtures = require("./ch13.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

const cases = [
  {
    kind: "ch13-recall",
    exercise(harness, panel, config, lang) {
      const find = (key) => button(panel, config.labels[key][lang]);
      clickChanged(harness, panel, find("recall"));
      const advance = find("advance");
      for (let i = 0; i < config.workers.length; i++) clickChanged(harness, panel, advance);
      clickChanged(harness, panel, find("wait"));
      clickChanged(harness, panel, find("shuffleReset"));
    },
    visible: "signal",
  },
  {
    kind: "ch13-timeout",
    exercise(harness, panel, config, lang) {
      const find = (key) => button(panel, config.labels[key][lang]);
      const tick = find("tick");
      for (let i = 0; i < config.deadlineTicks; i++) clickChanged(harness, panel, tick);
      clickChanged(harness, panel, find("reset"));
      clickChanged(harness, panel, find("cancel"));
    },
    visible: "deadline",
  },
  {
    kind: "ch13-first-wins",
    exercise(harness, panel, config, lang) {
      const find = (key) => button(panel, config.labels[key][lang]);
      const race = find("raceStep");
      clickChanged(harness, panel, race);
      const winnerLabel = config.labels.found[lang];
      const firstWinner = panel.querySelectorAll(".ch13-worker").find((row) => row.textContent.includes(winnerLabel)).textContent;
      clickChanged(harness, panel, find("respond"));
      clickChanged(harness, panel, find("shuffleAgain"));
      assert.notEqual(panel.querySelector(".ch13-order").textContent, `${config.labels.order[lang]}: ${config.orders[0].map((index) => config.workers[index]).join(" → ")}`);
      clickChanged(harness, panel, find("raceStep"));
      const nextWinner = panel.querySelectorAll(".ch13-worker").find((row) => row.textContent.includes(winnerLabel)).textContent;
      assert.notEqual(nextWinner, firstWinner);
      clickChanged(harness, panel, find("respond"));
    },
    visible: "order",
  },
  {
    kind: "ch13-context-chain",
    exercise(harness, panel, config, lang) {
      const find = (key) => button(panel, config.labels[key][lang]);
      clickChanged(harness, panel, find("next"));
      clickChanged(harness, panel, find("previous"));
      const reason = panel.querySelectorAll("button").at(-1);
      clickChanged(harness, panel, reason);
      const next = find("next");
      clickChanged(harness, panel, next);
      clickChanged(harness, panel, next);
      clickChanged(harness, panel, find("reset"));
    },
    visible: "same",
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
      const harness = createSandboxHarness({ chapter: 13, lang });
      const panel = harness.render(item.kind, fixtures[item.kind], { id: `${item.kind}-${lang}` });
      assert.ok(panel, `missing renderer for ${item.kind}`);
      assert.ok(panel.textContent.includes(fixtures[item.kind].labels[item.visible]?.[lang] || (lang === "ar" ? "السياق" : "Context")));
      assert.ok(panel.querySelectorAll("button").every((node) => node.tagName === "button"));
      item.exercise(harness, panel, fixtures[item.kind], lang);
    });
  }
}
