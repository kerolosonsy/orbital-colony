"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const { render, suggestions, errorText } = require("./chat.js");
const T = (ar, en) => en;

test("render escapes HTML and keeps the small Markdown a reply uses", () => {
  assert.equal(render("<b>hi</b>"), "<p>&lt;b&gt;hi&lt;/b&gt;</p>");
  assert.equal(render("use `Measure()` **first**"), '<p>use <code dir="ltr">Measure()</code> <b>first</b></p>');
  assert.equal(render("- one\n- two\nafter"), "<ul><li>one</li><li>two</li></ul><p>after</p>");
  assert.equal(render("1. a\n2. b"), "<ol><li>a</li><li>b</li></ol>");
  assert.equal(render("see:\n```go\nx := 1 < 2\n```\ndone"), '<p>see:</p><pre dir="ltr"><code>x := 1 &lt; 2</code></pre><p>done</p>');
});

test("suggestions lead with the error when the last run failed", () => {
  assert.equal(suggestions(T, true)[0], "Why does my code fail?");
  assert.equal(suggestions(T, false)[0], "I don't get the mission");
  assert.equal(suggestions(T, false).length, 3);
});

test("errorText tells how to set up Codex when it is missing", () => {
  assert.match(errorText(T, { setup: true }), /codex login/);
  assert.match(errorText(T, { error: "failed" }), /Try again/);
});
