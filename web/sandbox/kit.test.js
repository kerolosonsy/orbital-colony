"use strict";

const { test } = require("node:test");
const assert = require("node:assert/strict");
const { createSandboxHarness } = require("./harness.js");

function routeRenderer(block, { element, button, T, kit }) {
  const config = block.sandbox.config;
  const text = (pair) => T(pair.ar, pair.en);
  const label = (key) => text(config.labels[key]);
  const panel = element("section", "lesson-sandbox ch07-route-order");
  const stepper = kit.stepper({
    frames: config.routes,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(route, index) { return element("p", "route-frame", `${index + 1}. ${text(route)}`); },
  });
  panel.appendChild(stepper.node);
  return panel;
}

function config() {
  return {
    routes: [{ ar: "الحافة", en: "Rim" }, { ar: "الممر", en: "Relay" }, { ar: "الرصيف", en: "Dock" }],
    labels: {
      previous: { ar: "السابق", en: "Previous" },
      next: { ar: "التالي", en: "Next" },
      reset: { ar: "من الأول", en: "Reset" },
    },
  };
}

test("README route-order renderer runs through the chapter harness in Arabic and English", () => {
  for (const lang of ["ar", "en"]) {
    const harness = createSandboxHarness({ chapter: 7, lang });
    harness.register("ch07-route-order", routeRenderer);
    const box = harness.render("ch07-route-order", config(), { text: lang === "ar" ? "مسارات للتوضيح." : "Illustrative routes." });
    const buttons = box.querySelectorAll("button");
    const frame = () => box.querySelector(".route-frame").textContent;
    const counter = () => box.querySelector(".sandbox-kit-counter").textContent;
    assert.equal(buttons.length, 3);
    assert.equal(frame(), lang === "ar" ? "1. الحافة" : "1. Rim");
    assert.equal(counter(), "1 / 3");
    harness.click(buttons[1]);
    assert.equal(frame(), lang === "ar" ? "2. الممر" : "2. Relay");
    assert.equal(counter(), "2 / 3");
    harness.click(buttons[2]);
    assert.equal(frame(), lang === "ar" ? "1. الحافة" : "1. Rim");
  }
});
