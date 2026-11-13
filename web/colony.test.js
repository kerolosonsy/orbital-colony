"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const source = fs.readFileSync(__dirname + "/colony.js", "utf8");

function contracts(lang) {
  const context = vm.createContext({ T: (ar, en) => (lang === "en" ? en : ar), esc: (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;") });
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

const visibilityContext = vm.createContext({});
vm.runInContext(source.slice(source.indexOf("function visibleContracts("), source.indexOf("function treeSnapshot(")), visibilityContext);

test("contract links reveal a hidden target without expanding every locked contract", () => {
  const list = [
    { id: "late", chapter: 20, status: "locked" },
    { id: "next", chapter: 4, status: "locked" },
    { id: "done", chapter: 1, status: "done" },
    { id: "active", chapter: 3, status: "active" },
    { id: "other", chapter: 10, status: "locked" },
    { id: "second", chapter: 5, status: "locked" },
    { id: "open", chapter: 2, status: "available" },
  ];
  const ids = (target) => Array.from(visibilityContext.visibleContracts(list, target), (c) => c.id);
  assert.deepEqual(ids(""), ["active", "open", "done", "next", "second"]);
  assert.deepEqual(ids("late"), ["active", "open", "done", "next", "second", "late"]);
  assert.deepEqual(ids("next"), ids(""));
  assert.deepEqual(ids("missing"), ids(""));
  assert.deepEqual(list.map((c) => c.id), ["late", "next", "done", "active", "other", "second", "open"]);
  assert.equal(visibilityContext.visibleContracts([], "late").length, 0);
  assert.deepEqual(Array.from(visibilityContext.visibleContracts(list.slice(2, 4), ""), (c) => c.id), ["active", "done"]);
});


test("contract cards escape contract links and visible text", () => {
  const ctx = contracts("en");
  for (const status of ["active", "available", "done", "locked"]) {
    const markup = ctx.contractCard({ id: 'c" onclick="bad', chapter: 3, title: "<script>", brief: "<b>brief</b>", status }, null, false);
    assert.match(markup, /data-id="c&quot; onclick=&quot;bad"/);
    assert.match(markup, /&lt;script&gt;/);
    assert.doesNotMatch(markup, /<script>|<b>brief<\/b>|data-id="c" onclick=/);
  }
});

const appSource = fs.readFileSync(__dirname + "/app.js", "utf8");
const RAW_OLD = ["moss", "lichen", "glowcap", "reed", "crystal", "iron", "geode"];
const RAW_NEW = ["water", "silica", "copper", "sulfur", "titanium", "helium", "uranium"];
const CRAFTED = ["biomass", "plate", "wire", "glass", "fuel", "fertilizer", "circuit", "alloy", "battery", "coolant", "processor", "frame", "cell", "core"];
const ALL_GOODS = [...RAW_OLD, ...RAW_NEW, ...CRAFTED];
const BUILDINGS = ["composter", "smelter", "kiln", "refinery", "assembler", "reactor-lab", "silo", "generator", "lab"];

function appConstants(lang) {
  const ctx = vm.createContext({ LANG: lang });
  return vm.runInContext(
    appSource.slice(appSource.indexOf("const BURST"), appSource.indexOf("const DRONE_COLORS"))
    + "\n;({ BURST, UNIT, ORES, ENT, NAME, INV_COLOR });",
    ctx,
  );
}

function marketCtx(lang) {
  const { NAME, INV_COLOR } = appConstants(lang);
  const marketEl = { innerHTML: "" };
  const ctx = vm.createContext({
    LANG: lang,
    T: (ar, en) => (lang === "en" ? en : ar),
    esc: (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"),
    NAME,
    INV_COLOR,
    $: (sel) => (sel === "#col-market" ? marketEl : null),
  });
  const head = source.slice(source.indexOf("const OLD_GOODS"), source.indexOf("let data ="));
  const fmtCode = source.slice(source.indexOf("const fmt ="), source.indexOf("const actionMs"));
  const marketCode = source.slice(source.indexOf("function marketRowVisible("), source.indexOf("function renderBuildings("));
  const exposed = vm.runInContext(
    head + "\n" + fmtCode + "\n" + marketCode + "\n;({ BASE, marketMarkup, marketTable, marketRowVisible });",
    ctx,
  );
  return { ctx: exposed, marketEl };
}

test("names and colours cover every good and building in both languages with unique good colours", () => {
  const ar = appConstants("ar"), en = appConstants("en");
  for (const id of [...ALL_GOODS, ...BUILDINGS]) {
    assert.ok(typeof ar.NAME[id] === "string" && ar.NAME[id].length > 0, `missing ar NAME for ${id}`);
    assert.ok(typeof en.NAME[id] === "string" && en.NAME[id].length > 0, `missing en NAME for ${id}`);
    assert.notEqual(ar.NAME[id], en.NAME[id], `ar and en NAME should differ for ${id}`);
  }
  const seen = new Set();
  for (const good of ALL_GOODS) {
    const col = en.INV_COLOR[good];
    assert.match(col, /^#[0-9a-f]{6}$/i, `missing hex INV_COLOR for ${good}`);
    const norm = col.toLowerCase();
    assert.ok(!seen.has(norm), `duplicate colour ${norm} for ${good}`);
    seen.add(norm);
  }
  assert.ok(!seen.has(en.INV_COLOR.cargo.toLowerCase()), "cargo colour must stay distinct from goods");
  for (const ore of ["iron", "geode", "copper", "sulfur", "titanium", "uranium"]) {
    assert.match(en.ENT[ore], /^#[0-9a-f]{6}$/i, `missing ENT for ore ${ore}`);
    assert.ok(Array.isArray(en.BURST[ore]) && en.BURST[ore].length >= 2, `missing BURST for ore ${ore}`);
  }
});

test("market lists new goods only when held and shows group headings only when non-empty", () => {
  const expectedBase = {
    moss: 1, lichen: 3, glowcap: 2, reed: 14, crystal: 6, iron: 9, geode: 40,
    water: 2, silica: 5, copper: 12, sulfur: 18, titanium: 30, helium: 45, uranium: 90,
    biomass: 6, plate: 29, wire: 38, glass: 24, fuel: 22, fertilizer: 48, circuit: 170,
    alloy: 141, battery: 136, coolant: 78, processor: 880, frame: 941, cell: 619, core: 6818,
  };
  for (const lang of ["ar", "en"]) {
    const { ctx } = marketCtx(lang);
    const { NAME } = appConstants(lang);
    const rawHeading = lang === "en" ? "Raw" : "خامات";
    const craftedHeading = lang === "en" ? "Crafted" : "منتجات مصنّعة";
    for (const g of ALL_GOODS) assert.equal(ctx.BASE[g], expectedBase[g]);

    // Empty inventory: moss shows under Raw; no new good shows; Crafted heading stays hidden
    const emptyHtml = ctx.marketMarkup({ inv: {}, levels: {} }, {});
    assert.ok(emptyHtml.includes(NAME.moss));
    assert.ok(emptyHtml.includes(rawHeading));
    assert.ok(!emptyHtml.includes(craftedHeading));
    for (const g of [...RAW_NEW, ...CRAFTED, "lichen", "glowcap", "reed", "crystal", "iron", "geode"]) {
      assert.ok(!emptyHtml.includes(NAME[g]), `${g} should be hidden when not held`);
    }

    // Old rules preserved, fractional < 1 new goods stay hidden
    const partialHtml = ctx.marketMarkup({ inv: { water: 0.5, biomass: 0.9 }, levels: { "crop.lichen": 1, "api.dig": 1 } }, {});
    assert.ok(partialHtml.includes(NAME.lichen));
    assert.ok(partialHtml.includes(NAME.iron));
    assert.ok(partialHtml.includes(NAME.geode));
    assert.ok(!partialHtml.includes(NAME.water));
    assert.ok(!partialHtml.includes(NAME.biomass));
    assert.ok(!partialHtml.includes(craftedHeading));

    // Holding >= 1 of a new raw and crafted good shows them under their respective headings
    const heldHtml = ctx.marketMarkup({ inv: { water: 1, copper: 3, plate: 2, core: 1 }, levels: {} }, {});
    assert.ok(heldHtml.includes(rawHeading));
    assert.ok(heldHtml.includes(craftedHeading));
    assert.ok(heldHtml.includes(NAME.water));
    assert.ok(heldHtml.includes(NAME.copper));
    assert.ok(heldHtml.includes(NAME.plate));
    assert.ok(heldHtml.includes(NAME.core));
    assert.ok(heldHtml.indexOf(rawHeading) < heldHtml.indexOf(NAME.water));
    assert.ok(heldHtml.indexOf(NAME.water) < heldHtml.indexOf(craftedHeading));
    assert.ok(heldHtml.indexOf(craftedHeading) < heldHtml.indexOf(NAME.plate));

    // Group heading helper omits Raw when raw rows are empty and omits both when both are empty
    const onlyCrafted = ctx.marketTable([], ["plate"], { plate: 2 }, {});
    assert.ok(!onlyCrafted.includes(rawHeading));
    assert.ok(onlyCrafted.includes(craftedHeading));
    const none = ctx.marketTable([], [], {}, {});
    assert.ok(!none.includes(rawHeading));
    assert.ok(!none.includes(craftedHeading));
  }
});
