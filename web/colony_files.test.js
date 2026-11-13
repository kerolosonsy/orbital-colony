"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");

const appSource = fs.readFileSync(__dirname + "/app.js", "utf8");
const filesSlice = appSource.slice(
  appSource.indexOf("function levelFiles("),
  appSource.indexOf("// ---------- VS Code link"),
);

function createFilesHarness(lang = "en", initialRecord = {}) {
  const progress = { colony: structuredClone(initialRecord) };
  const tabBox = { innerHTML: "", children: [], append(...args) { this.children.push(...args); }, appendChild(x) { this.children.push(x); } };
  const saveState = { textContent: "" };
  const prompts = [], alerts = [], confirms = [];
  let savedCount = 0;

  const makeDomNode = (tag) => ({
    tagName: tag.toUpperCase(),
    type: "",
    className: "",
    innerHTML: "",
    textContent: "",
    title: "",
    disabled: false,
    children: [],
    attrs: new Map(),
    setAttribute(k, v) { this.attrs.set(k, String(v)); },
    getAttribute(k) { return this.attrs.get(k); },
    append(...args) { this.children.push(...args); },
    appendChild(x) { this.children.push(x); },
    focus() {},
  });

  const ctx = vm.createContext({
    LANG: lang,
    T: (ar, en) => (lang === "en" ? en : ar),
    esc: (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"),
    S: { cur: null, file: "main.go", review: false },
    rec: (id) => (progress[id] ||= {}),
    ed: { value: "", readOnly: false, focus() {} },
    errLines: new Map(),
    flushSave: () => {},
    renderEditor: () => {},
    isLinked: () => false,
    store: { save() { savedCount++; } },
    $: (sel) => (sel === "#file-tabs" ? tabBox : sel === "#save-state" ? saveState : makeDomNode("div")),
    document: { createElement: (tag) => makeDomNode(tag) },
    prompt: (msg, def) => (prompts.length ? prompts.shift() : def ?? null),
    alert: (msg) => alerts.push(msg),
    confirm: (msg) => (confirms.length ? confirms.shift() : true),
    Map, Object, Array, String, Set, RegExp,
  });

  vm.runInContext(filesSlice, ctx);
  return { ctx, progress, tabBox, prompts, alerts, confirms, getSaved: () => savedCount };
}

test("validateColonyPath accepts allowed relative paths and rejects unsafe or unknown paths", () => {
  const { ctx } = createFilesHarness("en");
  for (const okPath of ["helper.go", "plan/plan.go", "db/schema.sql", "cfg/settings.json", "cfg/app.yaml", "notes.txt", "README.md", "ui/panel.html"]) {
    const res = ctx.validateColonyPath(okPath, ["main.go"]);
    assert.equal(res.ok, true, `${okPath} should be valid: ${res.error}`);
    assert.equal(res.path, okPath);
  }

  for (const badPath of [
    "", "   ",
    "/main.go", "\\plan\\a.go", "C:a.go",
    "../secret.go", "plan/../main.go", "a..b/c.go",
    "./helper.go", ".env.go", "plan//a.go", "plan/", ".go",
    "script.py", "run.sh", "index.ts", "noext",
    "main.go", "main.go/sub.go",
  ]) {
    const res = ctx.validateColonyPath(badPath, ["main.go", "dir/sub.go"]);
    assert.equal(res.ok, false, `${badPath} should be rejected`);
  }
  assert.equal(ctx.validateColonyPath("dir", ["dir/sub.go"]).ok, false);
  assert.equal(ctx.validateColonyPath("dir.go/nested.go", ["dir.go"]).ok, false);
});

test("legacy single-file colony saves appear as editable main.go, migrate cleanly, and build files payload", () => {
  const level = { id: "colony", kind: "colony", starter: "package main\n// starter\n" };
  const { ctx, progress } = createFilesHarness("en", { code: "package main\n// legacy\n" });

  // Before migration: shows main.go editable with legacy code and builds payload
  const filesBefore = ctx.levelFiles(level);
  assert.deepEqual(JSON.parse(JSON.stringify(filesBefore)), [{ name: "main.go", starter: level.starter, readOnly: false }]);
  assert.equal(ctx.fileCode(level, filesBefore[0]), "package main\n// legacy\n");
  assert.deepEqual(JSON.parse(JSON.stringify(ctx.colonyFilesPayload(level))), { "main.go": "package main\n// legacy\n" });

  // Saving edits keeps both r.files["main.go"] and r.code updated
  ctx.saveFile(level, "main.go", "package main\n// updated\n");
  assert.equal(progress.colony.code, "package main\n// updated\n");
  assert.equal(progress.colony.files["main.go"], "package main\n// updated\n");
});

test("colony project file create, starter inference, non-Go editability, rename, and delete work cleanly", () => {
  const level = { id: "colony", kind: "colony", starter: "package main\n" };
  const { ctx, progress } = createFilesHarness("en", { code: "package main\n" });

  // Create Go file in subfolder -> infers package name; create non-Go file -> empty & editable
  const c1 = ctx.createColonyFile(level, "plan/plan.go");
  assert.deepEqual(JSON.parse(JSON.stringify(c1)), { ok: true, path: "plan/plan.go" });
  assert.equal(progress.colony.files["plan/plan.go"], "package plan\n\n");

  const c2 = ctx.createColonyFile(level, "notes.md");
  assert.deepEqual(JSON.parse(JSON.stringify(c2)), { ok: true, path: "notes.md" });

  const list = ctx.levelFiles(level);
  assert.deepEqual(Array.from(list, (f) => f.name), ["main.go", "plan/plan.go", "notes.md"]);
  assert.equal(list.every((f) => f.readOnly === false), true);

  ctx.saveFile(level, "notes.md", "# Plan\n");
  const payload = ctx.colonyFilesPayload(level);
  assert.deepEqual(JSON.parse(JSON.stringify(payload)), {
    "main.go": "package main\n",
    "plan/plan.go": "package plan\n\n",
    "notes.md": "# Plan\n",
  });

  // Cannot rename or delete main.go
  assert.equal(ctx.renameColonyFile(level, "main.go", "other.go").ok, false);
  assert.equal(ctx.deleteColonyFile(level, "main.go").ok, false);

  // Rename preserves content; delete removes file
  const ren = ctx.renameColonyFile(level, "notes.md", "docs/readme.md");
  assert.deepEqual(JSON.parse(JSON.stringify(ren)), { ok: true, path: "docs/readme.md" });
  assert.equal(progress.colony.files["notes.md"], undefined);
  assert.equal(progress.colony.files["docs/readme.md"], "# Plan\n");

  assert.equal(ctx.deleteColonyFile(level, "docs/readme.md").ok, true);
  assert.deepEqual(Array.from(ctx.levelFiles(level), (f) => f.name), ["main.go", "plan/plan.go"]);
});

test("firstErrorFile routes to matching file, strips ./ prefix, and falls back safely", () => {
  const { ctx } = createFilesHarness("en");
  const files = [{ name: "main.go" }, { name: "plan/plan.go" }, { name: "util.go" }];
  assert.equal(ctx.firstErrorFile([{ file: "./plan/plan.go", line: 3 }], files, "main.go"), "plan/plan.go");
  assert.equal(ctx.firstErrorFile([{ file: "util.go", line: 8 }], files, "main.go"), "util.go");
  assert.equal(ctx.firstErrorFile([{ file: "unknown.go", line: 1 }], files, "plan/plan.go"), "plan/plan.go");
  assert.equal(ctx.firstErrorFile([], files, "util.go"), "util.go");
});

test("renderTabs renders colony add/rename/delete controls on non-main files only and switches with firstErrorFile", () => {
  const level = { id: "colony", kind: "colony", starter: "package main\n" };
  const { ctx, tabBox, prompts, alerts } = createFilesHarness("en", { files: { "main.go": "package main\n", "plan/plan.go": "package plan\n" } });
  ctx.S.cur = level;
  ctx.S.file = "main.go";
  ctx.renderTabs();

  // First child is main.go plain button (no rename/delete), second is wrap div with tab-main + ren + del, last is '+' add button
  assert.equal(tabBox.children.length, 3);
  assert.equal(tabBox.children[0].className, "tab active");
  assert.equal(tabBox.children[0].children.length, 0);
  assert.equal(tabBox.children[1].className, "tab");
  assert.equal(tabBox.children[1].children.map((c) => c.className).join(","), "tab-main,tab-act,tab-act tab-del");
  assert.equal(tabBox.children[2].className, "tab tab-add");

  // Triggering new file prompt switches active file on valid input and alerts on bad extension
  prompts.push("../bad.go");
  tabBox.children[2].onclick();
  assert.equal(alerts.length, 1);
  assert.equal(ctx.S.file, "main.go");

  prompts.push("cfg/items.json");
  tabBox.children[2].onclick();
  assert.equal(ctx.S.file, "cfg/items.json");
});
