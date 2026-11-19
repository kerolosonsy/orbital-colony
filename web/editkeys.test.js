"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const { handle } = require("./editkeys.js");

// type applies the key at the "|" caret (or "[...]" selection) and returns
// the new text with the caret/selection marked the same way.
function type(marked, key, mods = {}) {
  let value = marked, start, end;
  if (marked.includes("|")) { start = end = marked.indexOf("|"); value = marked.replace("|", ""); }
  else { start = marked.indexOf("["); value = marked.replace("[", ""); end = value.indexOf("]"); value = value.replace("]", ""); }
  const edit = handle(value, start, end, { key, ...mods });
  if (!edit) return null;
  const out = value.slice(0, edit.from) + edit.text + value.slice(edit.to);
  if (edit.selStart === edit.selEnd) return out.slice(0, edit.selStart) + "|" + out.slice(edit.selStart);
  return out.slice(0, edit.selStart) + "[" + out.slice(edit.selStart, edit.selEnd) + "]" + out.slice(edit.selEnd);
}

test("Enter keeps indentation and indents after an opener", () => {
  assert.equal(type("\tx := 1|", "Enter"), "\tx := 1\n\t|");
  assert.equal(type("\tif ok {|", "Enter"), "\tif ok {\n\t\t|");
  assert.equal(type("\tcase 1:|", "Enter"), "\tcase 1:\n\t\t|");
  assert.equal(type("\tdefault:|", "Enter"), "\tdefault:\n\t\t|");
});

test("Enter between braces puts the closer on its own line", () => {
  assert.equal(type("func main() {|}", "Enter"), "func main() {\n\t|\n}");
  assert.equal(type("\tf(|)", "Enter"), "\tf(\n\t\t|\n\t)");
});

test("brackets and quotes close themselves and are typed over", () => {
  assert.equal(type("fmt.Println|", "("), "fmt.Println(|)");
  assert.equal(type("x := |", '"'), 'x := "|"');
  assert.equal(type("f(x|)", ")"), "f(x)|");
  assert.equal(type('"hi|"', '"'), '"hi"|');
  // no pairing in the middle of a word or after a word character for quotes
  assert.equal(type("|abc", "("), null);
  assert.equal(type("abc|", '"'), null);
});

test("a selection is wrapped by the bracket", () => {
  assert.equal(type("[x + 1]", "("), "([x + 1])");
});

test("Backspace removes an empty pair", () => {
  assert.equal(type("f(|)", "Backspace"), "f|");
  assert.equal(type("f(x|)", "Backspace"), null);
});

test("Tab indents selected lines and Shift+Tab outdents them", () => {
  assert.equal(type("[a\nb]", "Tab"), "[\ta\n\tb]");
  assert.equal(type("[\ta\n\tb]", "Tab", { shiftKey: true }), "[a\nb]");
  assert.equal(type("\t\tx|", "Tab", { shiftKey: true }), "\tx|");
  assert.equal(type("x|", "Tab"), "x\t|");
});

test("} on an indent-only line dedents one level", () => {
  assert.equal(type("\t\t|", "}"), "\t}|");
});

test("Ctrl/⌘+/ toggles a line comment", () => {
  assert.equal(type("\tx := 1|", "/", { ctrlKey: true }), "\t// x := 1|");
  assert.equal(type("\t// x := 1|", "/", { metaKey: true }), "\tx := 1|");
  assert.equal(type("[a\n\nb]", "/", { ctrlKey: true }), "[// a\n\n// b]");
});

test("other keys and modified keys are left to the browser", () => {
  assert.equal(type("x|", "a"), null);
  assert.equal(type("x|", "(", { altKey: true }), null);
  assert.equal(type("x|", "Enter", { isComposing: true }), null);
});

test("Alt+arrows move lines and Shift+Alt+arrows duplicate them", () => {
  assert.equal(type("a\nb|\nc", "ArrowUp", { altKey: true }), "b|\na\nc");
  assert.equal(type("a|\nb\nc", "ArrowDown", { altKey: true }), "b\na|\nc");
  assert.equal(type("a|\nb", "ArrowUp", { altKey: true }), null);
  assert.equal(type("a\nb|", "ArrowDown", { altKey: true }), null);
  assert.equal(type("[a\nb]\nc", "ArrowDown", { altKey: true }), "c\n[a\nb]");
  assert.equal(type("x|", "ArrowDown", { altKey: true, shiftKey: true }), "x\nx|");
  assert.equal(type("x|", "ArrowUp", { altKey: true, shiftKey: true }), "x|\nx");
});
