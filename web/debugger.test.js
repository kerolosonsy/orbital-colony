"use strict";

const { test } = require("node:test");
const assert = require("node:assert/strict");
const Debugger = require("./debugger.js");

const frames = [
  { file: "main.go", line: 1, func: "main", goroutine: 1, depth: 0, event: "call" },
  { file: "main.go", line: 2, func: "main", goroutine: 1, depth: 0, event: "line" },
  { file: "main.go", line: 6, func: "helper", goroutine: 1, depth: 1, event: "call" },
  { file: "main.go", line: 7, func: "helper", goroutine: 1, depth: 1, event: "line" },
  { file: "main.go", line: 10, func: "nested", goroutine: 1, depth: 2, event: "call" },
  { file: "main.go", line: 11, func: "nested", goroutine: 1, depth: 2, event: "line" },
  { file: "main.go", line: 12, func: "nested", goroutine: 1, depth: 2, event: "return" },
  { file: "main.go", line: 8, func: "helper", goroutine: 1, depth: 1, event: "line" },
  { file: "main.go", line: 9, func: "helper", goroutine: 1, depth: 1, event: "return" },
  { file: "main.go", line: 3, func: "main", goroutine: 1, depth: 0, event: "line" },
  { file: "main.go", line: 4, func: "main", goroutine: 1, depth: 0, event: "return" },
];

test("debugger supports restart, back, into, over, and out", () => {
  const stepper = Debugger.createStepper(frames);
  assert.equal(stepper.restart().event, "call");
  assert.equal(stepper.set(1).line, 2);
  assert.equal(stepper.into().func, "helper");
  assert.equal(stepper.set(1).line, 2);
  assert.equal(stepper.over().line, 3);
  assert.equal(stepper.set(3).func, "helper");
  assert.equal(stepper.out().line, 3);
  assert.equal(stepper.back().event, "return");
});

test("continue stops on the next breakpoint or at the end", () => {
  const stepper = Debugger.createStepper(frames);
  assert.equal(stepper.set(1).line, 2);
  assert.equal(stepper.continueTo(["main.go:7"]).line, 7);
  assert.equal(stepper.continueTo([]).line, 4);
});

test("breakpoints toggle by file and line", () => {
  const one = Debugger.toggleBreakpointList([], "main.go", 8);
  assert.deepEqual(one, ["main.go:8"]);
  assert.deepEqual(Debugger.toggleBreakpointList(one, "main.go", 8), []);
});

test("inline values are readable, escaped, and shortened", () => {
  assert.equal(Debugger.inlineValueText([{ name: "i", value: "2" }, { name: "steps", value: "1" }]), "i = 2, steps = 1");
  assert.match(Debugger.renderInlineValues([{ name: "x", value: "<tag>" }]), /&lt;tag&gt;/);
  assert.ok(Debugger.inlineValueText([{ name: "text", value: "x".repeat(120) }], 20).endsWith("…"));
  const markup = Debugger.editorMarkup("x := 1", { active: true, line: 1, inlineValuesByLine: { 1: [{ name: "x", value: "1" }] } }, (line) => line);
  assert.match(markup, /is-current/);
  assert.match(markup, /<b>x<\/b> 1/);
});

test("inline rendering preserves multiline syntax highlighting", () => {
  const source = "/* first\ninside */\nx := 1";
  const markup = Debugger.editorMarkup(source, { active: true, line: 2 }, (code) => code
    .replace(/\/\*[\s\S]*?\*\//g, '<span class="tk-com">$&</span>'));
  assert.match(markup, /<span class="debug-line is-current" data-debug-line="2"><span class="tk-com">inside \*\/<\/span><\/span>/);
});

test("editing exits the active debug session", () => {
  const session = Debugger.createSession();
  session.begin({ available: true, frames: [{ line: 1 }] });
  assert.equal(session.active, true);
  session.edit();
  assert.equal(session.active, false);
  assert.equal(session.response, null);
});

test("gutter marks clickable breakpoint lines", () => {
  const markup = Debugger.gutterMarkup(2, new Map(), [2], 1);
  assert.match(markup, /data-debug-line="2"/);
  assert.match(markup, /has-breakpoint/);
  assert.match(markup, /debug-active/);
});

// ---- run playback: hooks fire in order and the end is reported once ----
function fakeDebuggerDOM() {
  const make = () => ({
    hidden: false, textContent: "", innerHTML: "", value: "0", max: "0", disabled: false, title: "", dataset: {},
    classList: { toggle() {}, add() {}, remove() {}, contains: () => false },
    setAttribute() {}, querySelectorAll: () => [], onclick: null, oninput: null,
  });
  const ids = ["debugger-toolbar", "debugger-controls", "debugger-status", "debugger-counter", "debugger-slider", "debugger-panel",
    "debugger-variables", "debugger-stack", "debugger-output", "debugger-order", "debugger-collapse", "debugger-where-text"];
  const els = Object.fromEntries(ids.map((id) => [`#${id}`, make()]));
  global.document = {
    documentElement: { lang: "en" },
    querySelector: (selector) => els[selector] || null,
    querySelectorAll: () => [],
    addEventListener() {},
  };
  global.matchMedia = () => ({ matches: false });
  return els;
}

function sessionWith(options) {
  fakeDebuggerDOM();
  const calls = { frames: [], end: 0, exit: 0, clear: 0 };
  Debugger.bind({ renderEditor() {}, getFile: () => "main.go", getBreakpoints: () => [] });
  Debugger.startWith("1-1", { available: true, frames: frames.map((f, i) => ({ ...f, events: i })), stdout: "" }, {
    ...options,
    onFrame: (frame, index) => calls.frames.push(index),
    onEnd: () => calls.end++,
    onExit: () => calls.exit++,
    onClear: () => calls.clear++,
  });
  return calls;
}

test("Run playback plays every step and reports the end once", (t) => {
  t.mock.timers.enable({ apis: ["setInterval"] });
  const calls = sessionWith({ autoplay: true });
  assert.equal(calls.frames[0], 0);
  assert.equal(Debugger.playing, true);
  t.mock.timers.tick(60000);
  assert.equal(calls.frames.at(-1), frames.length - 1);
  assert.equal(calls.end, 1);
  assert.equal(Debugger.playing, false);
  Debugger.clear();
  assert.equal(calls.clear, 1);
});

test("Step mode waits on the first step; skipping reveals once", () => {
  const calls = sessionWith({ autoplay: false });
  assert.deepEqual(calls.frames, [0]);
  assert.equal(Debugger.playing, false);
  assert.equal(calls.end, 0);
  Debugger.finishNow();
  assert.equal(calls.end, 1);
  Debugger.finishNow();
  assert.equal(calls.end, 1);
  Debugger.clear();
});

test("dropping a session from outside calls onClear exactly once", () => {
  const calls = sessionWith({ autoplay: false });
  Debugger.clear();
  Debugger.clear();
  assert.equal(calls.clear, 1);
  assert.equal(calls.end, 0);
});
