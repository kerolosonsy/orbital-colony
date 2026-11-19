"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixtures = require("./ch01.fixtures.json");
const { createSandboxHarness } = require("./harness.js");

function localePairs(value, result = []) {
  if (!value || typeof value !== "object") return result;
  if (typeof value.ar === "string" && typeof value.en === "string" && value.ar !== value.en) result.push(value);
  for (const child of Object.values(value)) {
    localePairs(child, result);
  }
  return result;
}

function captionPairs(value, key = "", result = []) {
  if (!value || typeof value !== "object") return result;
  if (key === "caption" && typeof value.ar === "string" && typeof value.en === "string") result.push(value);
  for (const [childKey, child] of Object.entries(value)) captionPairs(child, childKey, result);
  return result;
}

function exercise(kind, harness, node) {
  const buttons = node.querySelectorAll("button");
  assert.ok(buttons.length <= 4, `${kind} has at most four controls`);
  if (kind === "ch01-route-stepper" || kind === "ch01-planting-state") {
    assert.equal(buttons.length, 3);
    const start = node.textContent;
    harness.click(buttons[1]);
    assert.notEqual(node.textContent, start);
    const advanced = node.textContent;
    harness.click(buttons[0]);
    assert.equal(node.textContent, start);
    harness.click(buttons[1]);
    harness.click(buttons[2]);
    assert.equal(node.textContent, start);
    assert.notEqual(advanced, start);
    return;
  }
  assert.equal(kind, "ch01-harvest-call");
  assert.equal(buttons.length, 3);
  const start = node.textContent;
  harness.click(buttons[1]);
  const moved = node.textContent;
  harness.click(buttons[1]);
  assert.notEqual(node.textContent, moved);
  assert.match(node.textContent, /1/);
  harness.click(buttons[0]);
  assert.equal(node.textContent, moved);
  harness.click(buttons[2]);
  assert.equal(node.textContent, start);
}

for (const [levelID, fixture] of Object.entries(fixtures)) {
  test(`${fixture.kind} renders and responds in Arabic and English`, () => {
    for (const lang of ["ar", "en"]) {
      const harness = createSandboxHarness({ chapter: 1, lang });
      const node = harness.render(fixture.kind, fixture.config, { id: levelID });
      assert.ok(node);
      assert.ok(node.className.includes(fixture.kind));
      const pairs = localePairs(fixture.config);
      assert.ok(pairs.length > 0);
      assert.ok(pairs.some((pair) => node.textContent.includes(pair[lang])), `${fixture.kind} renders ${lang} labels`);
      for (const caption of captionPairs(fixture.config)) {
        assert.ok(caption[lang].trim().split(/\s+/).length <= 8, `${fixture.kind} ${lang} caption is concise`);
      }
      exercise(fixture.kind, harness, node);
    }
  });
}
