"use strict";

const { test } = require("node:test");
const assert = require("node:assert/strict");
const TracePanel = require("./trace.js");

test("trace panel stepper moves through frames and clamps both ends", () => {
  const frames = [{ line: 2 }, { line: 5 }, { line: 8 }];
  const updates = [];
  const stepper = TracePanel.createStepper(frames, (frame, index) => updates.push([frame?.line || null, index]));
  assert.equal(stepper.frame.line, 2);
  assert.equal(stepper.next().line, 5);
  assert.equal(stepper.set(50).line, 8);
  assert.equal(stepper.next().line, 8);
  assert.equal(stepper.previous().line, 5);
  assert.deepEqual(updates.map(([line]) => line), [5, 8, 8, 5]);
});

test("trace code panel highlights the current source line and escapes learner code", () => {
  const markup = TracePanel.codeMarkup("first\ncall(x < y)\nlast", 2);
  assert.match(markup, /data-trace-line="2"><span class="trace-line-number">2<\/span><code>call\(x &lt; y\)<\/code>/);
  assert.equal((markup.match(/is-current/g) || []).length, 1);
  assert.match(markup, /data-trace-line="1"/);
  assert.match(markup, /data-trace-line="3"/);
});
