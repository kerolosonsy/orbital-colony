"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const source = fs.readFileSync(__dirname + "/colony.js", "utf8");

function contracts(lang) {
  const context = vm.createContext({ T: (ar, en) => (lang === "en" ? en : ar), esc: (s) => String(s) });
  vm.runInContext(source.slice(source.indexOf("const clock"), source.indexOf("function renderContracts(")), context);
  return context;
}

test("contract cards show every state in both languages", () => {
  for (const lang of ["ar", "en"]) {
    const ctx = contracts(lang);
    const card = (c, now = null, running = false) => ctx.contractCard({ id: "c03", chapter: 3, title: "Full Sweep", brief: "Deliver", ...c }, now, running);
    const en = lang === "en";
    assert.match(card({ status: "locked" }), en ? /Finish chapter 3/ : /خلّص الفصل 3/);
    assert.match(card({ status: "locked", after: ["The Right Ground"] }), en ? /Finish the contract first: The Right Ground/ : /خلّص العقد الأول: The Right Ground/);
    const available = card({ status: "available" });
    assert.match(available, en ? />Accept</ : />اقبل</);
    assert.match(available, en ? /Lesson: chapter 3/ : /الدرس: فصل 3/);
    assert.match(card({ status: "available" }, null, true), /disabled/);
    assert.match(card({ status: "done", run: { best: 615 } }), en ? /Best time 1:02/ : /أحسن وقت 1:02/);
    const active = card({ status: "active", frac: 0.25, left: 900, progress: "Moss 30/120" });
    assert.match(active, /width:25\.0%/);
    assert.match(active, en ? /Left 1:30/ : /باقي 1:30/);
    assert.match(active, /Moss 30\/120/);
    // live progress from the stream wins over the last summary
    assert.match(card({ status: "active", frac: 0.25, left: 900 }, { frac: 0.5, left: 300, progress: "Moss 60/120" }), /width:50\.0%[\s\S]*Moss 60\/120[\s\S]*0:30/);
    assert.doesNotMatch(card({ status: "active", left: -1 }), en ? /Left/ : /باقي/);
  }
});
