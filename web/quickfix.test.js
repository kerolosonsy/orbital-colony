"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const { fixFor, apply } = require("./quickfix.js");

test("fixFor recognises missing and unused standard imports", () => {
  assert.deepEqual(fixFor("undefined: fmt"), { kind: "add", path: "fmt" });
  assert.deepEqual(fixFor("undefined: json"), { kind: "add", path: "encoding/json" });
  assert.equal(fixFor("undefined: totalCount"), null);
  assert.deepEqual(fixFor('"strings" imported and not used'), { kind: "remove", path: "strings" });
  assert.equal(fixFor("declared and not used: x"), null);
});

test("adding an import grows a single import into a sorted block", () => {
  const code = 'package main\n\nimport "orbital/colony"\n\nfunc main() {}\n';
  assert.equal(apply(code, { kind: "add", path: "fmt" }), 'package main\n\nimport (\n\t"fmt"\n\t"orbital/colony"\n)\n\nfunc main() {}\n');
  const block = 'package main\n\nimport (\n\t"fmt"\n\t"orbital/colony"\n)\n\nfunc main() {}\n';
  assert.equal(apply(block, { kind: "add", path: "strings" }), 'package main\n\nimport (\n\t"fmt"\n\t"orbital/colony"\n\t"strings"\n)\n\nfunc main() {}\n');
  assert.equal(apply(block, { kind: "add", path: "fmt" }), null);
  assert.equal(apply("package main\n\nfunc main() {}\n", { kind: "add", path: "fmt" }), 'package main\n\nimport "fmt"\n\nfunc main() {}\n');
});

test("removing an import shrinks the block or drops it", () => {
  const block = 'package main\n\nimport (\n\t"fmt"\n\t"orbital/colony"\n)\n\nfunc main() {}\n';
  assert.equal(apply(block, { kind: "remove", path: "fmt" }), 'package main\n\nimport "orbital/colony"\n\nfunc main() {}\n');
  const single = 'package main\n\nimport "fmt"\n\nfunc main() {}\n';
  assert.equal(apply(single, { kind: "remove", path: "fmt" }), "package main\n\nfunc main() {}\n");
  assert.equal(apply(single, { kind: "remove", path: "os" }), null);
});
