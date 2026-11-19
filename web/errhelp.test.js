"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");

class Element {
  constructor(tagName) { this.tagName = tagName; this.children = []; this.textContent = ""; }
  appendChild(child) { this.children.push(child); return child; }
}

test("error help renders a native collapsed disclosure with localized cause and check", () => {
  const window = {};
  const context = vm.createContext({
    window,
    document: { createElement: (name) => new Element(name) },
    T: (ar, en) => context.LANG === "en" ? en : ar,
    LANG: "en",
  });
  vm.runInContext(fs.readFileSync(__dirname + "/errhelp.js", "utf8"), context);

  const details = window.ErrorHelp.render({ meaning: "Index out of range.", cause: "The index is too large.", check: "Check the current length." });
  assert.equal(details.tagName, "details");
  assert.equal(details.className, "err-help");
  assert.equal(details.dir, "auto");
  assert.equal(details.open, undefined);
  assert.equal(details.children[0].tagName, "summary");
  assert.equal(details.children[0].textContent, "Index out of range.");
  const rows = details.children[1].children;
  assert.equal(rows[0].children[0].textContent, "Likely cause: ");
  assert.equal(rows[0].children[1].textContent, "The index is too large.");
  assert.equal(rows[1].children[0].textContent, "Check: ");
  assert.equal(rows[1].children[1].textContent, "Check the current length.");
  details.open = true;
  assert.equal(details.open, true);
});

test("error help handles Arabic and missing optional details", () => {
  const window = {};
  const context = vm.createContext({ window, document: { createElement: (name) => new Element(name) }, T: (ar) => ar, LANG: "ar" });
  vm.runInContext(fs.readFileSync(__dirname + "/errhelp.js", "utf8"), context);

  const details = window.ErrorHelp.render({ meaning: "فيه مشكلة.", cause: "السبب." });
  assert.equal(details.dir, "auto");
  assert.equal(details.children[0].textContent, "فيه مشكلة.");
  assert.equal(details.children[1].children.length, 1);
  assert.equal(details.children[1].children[0].children[0].textContent, "السبب غالباً: ");
  assert.equal(window.ErrorHelp.render({}), null);
});
