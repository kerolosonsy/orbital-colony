"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const fixtures = require("./ch02.fixtures.json");
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
  if (kind === "ch02-inventory-flow") {
    assert.equal(buttons.length, 4);
    const start = node.textContent;
    harness.click(buttons[1]);
    const firstRead = node.textContent;
    harness.click(buttons[0]);
    assert.equal(node.textContent, start);
    harness.click(buttons[1]);
    harness.click(buttons[1]);
    harness.click(buttons[1]);
    assert.match(node.textContent, /9/);
    harness.click(buttons[2]);
    assert.equal(node.textContent, start);
    const beforeToggle = node.textContent;
    harness.click(buttons[3]);
    assert.notEqual(node.textContent, beforeToggle);
    assert.notEqual(firstRead, start);
    return;
  }
  if (kind === "ch02-return-pair") {
    assert.equal(buttons.length, 0);
    const inputs = node.querySelectorAll("input");
    assert.equal(inputs.length, 2);
    for (const input of inputs) {
      const before = node.textContent;
      const prior = Number(input.value);
      input.value = String(prior + 1);
      input.oninput();
      assert.notEqual(node.textContent, before);
      assert.notEqual(input.value, String(prior));
    }
    return;
  }
  if (kind === "ch02-utf8-span") {
    assert.equal(buttons.length, 2);
    const start = node.textContent;
    harness.click(buttons[1]);
    const secondRune = node.textContent;
    assert.notEqual(secondRune, start);
    assert.equal(node.querySelectorAll(".ch02-utf8-group").filter((group) => group.classList.contains("selected")).length, 1);
    harness.click(buttons[0]);
    assert.equal(node.textContent, start);
    return;
  }
  assert.equal(kind, "ch02-unit-align");
  assert.equal(buttons.length, 3);
  let previous = node.textContent;
  for (const choice of buttons) {
    harness.click(choice);
    assert.notEqual(node.textContent, previous);
    previous = node.textContent;
  }
  assert.equal(node.querySelectorAll('button[aria-pressed="true"]').length, 1);
}

for (const [levelID, fixture] of Object.entries(fixtures)) {
  test(`${fixture.kind} renders and responds in Arabic and English`, () => {
    for (const lang of ["ar", "en"]) {
      const harness = createSandboxHarness({ chapter: 2, lang });
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
