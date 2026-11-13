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
    assert.doesNotMatch(available, /ct-repeat|ct-rp|↻/);
    assert.match(card({ status: "available" }, null, true), /disabled/);
    assert.match(card({ status: "done", run: { best: 615 } }), en ? /Best time 1:02/ : /أحسن وقت 1:02/);
    const active = card({ status: "active", frac: 0.25, left: 900, progress: "Moss 30/120" });
    assert.match(active, /width:25\.0%/);
    assert.match(active, en ? /Left 1:30/ : /باقي 1:30/);
    assert.match(active, /Moss 30\/120/);
    // live progress from the stream wins over the last summary
    assert.match(card({ status: "active", frac: 0.25, left: 900 }, { frac: 0.5, left: 300, progress: "Moss 60/120" }), /width:50\.0%[\s\S]*Moss 60\/120[\s\S]*0:30/);
    assert.doesNotMatch(card({ status: "active", left: -1 }), en ? /Left/ : /باقي/);

    // repeatable contracts show repeat mark and RP reward across states
    const rep = card({ status: "available", repeatable: true, rewardRP: 8 });
    assert.match(rep, en ? /↻ Repeatable/ : /↻ متكرر/);
    assert.match(rep, en ? /\+8<\/b> RP/ : /\+8<\/b> نقطة بحث/);
    assert.match(card({ status: "done", repeatable: true, rewardRP: 0 }), en ? /↻ Repeatable/ : /↻ متكرر/);
    assert.doesNotMatch(card({ status: "done", repeatable: true, rewardRP: 0 }), /ct-rp/);
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
const BUILDINGS = ["composter", "smelter", "kiln", "refinery", "assembler", "reactor-lab", "silo", "generator", "lab", "collector"];

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

function systemCtx(lang) {
  const { NAME, INV_COLOR } = appConstants(lang);
  const fullApp = vm.runInContext(
    appSource.slice(appSource.indexOf("const BURST"), appSource.indexOf("function setWorld"))
    + "\n;({ DRONE_COLORS });",
    vm.createContext({ LANG: lang }),
  );
  const ctx = vm.createContext({
    LANG: lang,
    T: (ar, en) => (lang === "en" ? en : ar),
    esc: (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"),
    NAME,
    INV_COLOR,
    DRONE_COLORS: fullApp.DRONE_COLORS,
  });
  const fmtCode = source.slice(source.indexOf("const fmt ="), source.indexOf("const actionMs"));
  const sysCode = source.slice(source.indexOf("// ---------- system readouts"), source.indexOf("function highlightRank("));
  return vm.runInContext(
    fmtCode + "\n" + sysCode + "\n;({ fmtNum, escAttr, rankGoalLabel, rankMarkup, statusExtraMarkup, hasCargoData, fleetVisible, fleetMarkup });",
    ctx,
  );
}

test("rankGoalLabel and rankMarkup render goal kinds, fractional bar, filling note, and hide when absent", () => {
  for (const lang of ["ar", "en"]) {
    const ctx = systemCtx(lang);
    const en = lang === "en";
    assert.equal(ctx.rankMarkup(null), "");
    assert.equal(ctx.rankMarkup(undefined), "");
    assert.equal(ctx.rankGoalLabel("output"), en ? "Output per minute" : "إنتاج في الدقيقة");
    assert.equal(ctx.rankGoalLabel("efficiency"), en ? "Value per command" : "قيمة لكل أمر");
    assert.equal(ctx.rankGoalLabel("<x>"), "<x>");

    const filling = ctx.rankMarkup({ level: 3, kind: "efficiency", target: 1.32, value: 1.1, filled: 0.8 });
    assert.match(filling, en ? /Rank 3/ : /المرتبة 3/);
    assert.match(filling, en ? /Value per command/ : /قيمة لكل أمر/);
    assert.match(filling, /1\.1 \/ 1\.32/);
    assert.match(filling, /width:83\.3%/);
    assert.match(filling, /class="rank-note"/);
    assert.match(filling, /80%/);

    const full = ctx.rankMarkup({ level: 4, kind: "output", target: 200, value: 250, filled: 1 });
    assert.match(full, en ? /Output per minute/ : /إنتاج في الدقيقة/);
    assert.match(full, /250 \/ 200/);
    assert.match(full, /width:100\.0%/);
    assert.doesNotMatch(full, /rank-note/);

    const escaped = ctx.rankMarkup({ level: 1, kind: "<script>", target: 10, value: 5, filled: 1 });
    assert.doesNotMatch(escaped, /<script>/);
    assert.match(escaped, /&lt;script&gt;/);
  }
});

test("statusExtraMarkup renders power, storage, rp, rent only when present and enforces warning thresholds", () => {
  for (const lang of ["ar", "en"]) {
    const ctx = systemCtx(lang);
    const en = lang === "en";
    assert.equal(ctx.statusExtraMarkup({}), "");
    assert.equal(ctx.statusExtraMarkup({ credits: 50, size: 4 }), "");

    // Power: exact match demand === supply is NOT warning; demand > supply IS warning
    const okPower = ctx.statusExtraMarkup({ power: { supply: 12.5, demand: 12.5, day: true } });
    assert.match(okPower, /12\.5 \/ 12\.5/);
    assert.match(okPower, en ? /Day/ : /نهار/);
    assert.doesNotMatch(okPower, /col-power warn/);

    const badPower = ctx.statusExtraMarkup({ power: { supply: 9, demand: 12.5, day: false } });
    assert.match(badPower, /9 \/ 12\.5/);
    assert.match(badPower, en ? /Night/ : /ليل/);
    assert.match(badPower, /col-power warn/);

    // Storage: exact 90% is NOT warning; > 90% IS warning
    const storageHtml = ctx.statusExtraMarkup({
      storage: {
        crops: { used: 180, cap: 200 },
        raw: { used: 91, cap: 100 },
        crafted: { used: 64, cap: 70 },
      },
    });
    assert.match(storageHtml, /col-bar-row" data-storage="crops"/);
    assert.match(storageHtml, /col-bar-row warn" data-storage="raw"/);
    assert.match(storageHtml, /col-bar-row warn" data-storage="crafted"/);
    assert.match(storageHtml, /180 \/ 200/);
    assert.match(storageHtml, /91 \/ 100/);
    assert.match(storageHtml, /64 \/ 70/);

    // RP and rent: 0 is valid and renders; absent hides
    const zeroHtml = ctx.statusExtraMarkup({ rp: 0, rent: 0 });
    assert.match(zeroHtml, /col-rp/);
    assert.match(zeroHtml, /col-rent/);
    assert.match(zeroHtml, en ? /0<\/b> RP/ : /0<\/b> نقطة/);
    assert.match(zeroHtml, /0 CR/);
  }
});

test("fleetVisible and fleetMarkup enforce visibility rules and render battery + cargo chips", () => {
  for (const lang of ["ar", "en"]) {
    const ctx = systemCtx(lang);
    const en = lang === "en";
    assert.equal(ctx.fleetVisible(undefined), false);
    assert.equal(ctx.fleetVisible(null), false);
    assert.equal(ctx.fleetVisible([]), false);
    assert.equal(ctx.fleetVisible([{ id: 0, battery: 80 }]), false);
    assert.equal(ctx.fleetVisible([{ id: 0, battery: 80, cargo: {}, cap: 0 }]), false);
    assert.equal(ctx.fleetMarkup([{ id: 0, battery: 80 }]), "");

    // Single drone WITH cargo system data shows
    assert.equal(ctx.fleetVisible([{ id: 0, battery: 80, cargo: {}, cap: 10 }]), true);
    const oneWithCargo = ctx.fleetMarkup([{ id: 0, battery: 80, cargo: { moss: 4 }, cap: 10 }]);
    assert.match(oneWithCargo, en ? /<h3>Fleet <small>1<\/small><\/h3>/ : /<h3>الأسطول <small>1<\/small><\/h3>/);
    assert.match(oneWithCargo, en ? /Drone 0/ : /درون 0/);
    assert.match(oneWithCargo, /80%/);
    assert.match(oneWithCargo, /4 \/ 10/);
    assert.match(oneWithCargo, /#93c45c/);

    // Multiple drones without cargo data still show; low battery (<20%) warns; unsafe keys escaped
    assert.equal(ctx.fleetVisible([{ id: 0, battery: 80 }, { id: 1, battery: 15 }]), true);
    const multi = ctx.fleetMarkup([
      { id: 1, battery: 15, cargo: { "<bad>": 2 }, cap: 5 },
      { id: 0, battery: 95 },
    ]);
    assert.match(multi, /flt-bat warn.*15%/);
    assert.match(multi, /&lt;bad&gt;/);
    assert.doesNotMatch(multi, /title="<bad>"/);
    assert.ok(multi.indexOf('data-drone="0"') < multi.indexOf('data-drone="1"'));
  }
});

test("Colony live events update rank, status extras, fleet, tree header and handle rank-up highlight", async () => {
  function makeEl(id = "") {
    const classes = new Set();
    return {
      id,
      innerHTML: "",
      textContent: "",
      hidden: false,
      dataset: {},
      classList: {
        add: (...cs) => cs.forEach((c) => classes.add(c)),
        remove: (...cs) => cs.forEach((c) => classes.delete(c)),
        toggle: (c, force) => {
          const next = force === undefined ? !classes.has(c) : !!force;
          if (next) classes.add(c);
          else classes.delete(c);
          return next;
        },
        contains: (c) => classes.has(c),
      },
      querySelector: () => makeEl(),
      querySelectorAll: () => [],
    };
  }
  const els = new Map();
  const getEl = (sel) => { if (!els.has(sel)) els.set(sel, makeEl(sel)); return els.get(sel); };
  const logs = [], sounds = [], treeUpdates = [], timers = [];
  let now = 1000;
  let refreshCount = 0;

  const ctx = vm.createContext({
    LANG: "en",
    T: (_ar, en) => en,
    esc: (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"),
    document: { body: { classList: makeEl().classList } },
    performance: { now: () => now },
    setTimeout: (fn, ms) => { timers.push({ fn, ms, cleared: false }); return timers.length; },
    clearTimeout: (id) => { if (id && timers[id - 1]) timers[id - 1].cleared = true; },
    $: getEl,
    S: { mode: "colony", levels: [], chapters: [] },
    W: { w: 4, h: 4, depth: 1, view: 0, follow: true, colony: true, live: true, t: 0, inv: {}, layers: [[]], drones: new Map() },
    TreeView: { isOpen: false, update: (snap) => treeUpdates.push(snap), close() {}, open() {} },
    Sound: { play: (name) => sounds.push(name) },
    log: (msg, cls) => logs.push({ msg, cls }),
    hud: () => {},
    setDigits: (el, v) => { el.textContent = String(v); },
    renderLegend: () => {},
    draw: () => {},
    api: async (path) => {
      if (path === "/api/colony") {
        refreshCount++;
        return {
          state: { tick: 20, size: 4, depth: 1, credits: 500, levels: {}, inv: {}, tiles: [], buildings: [], rank: 4, rp: 32 },
          tree: [{ id: "root" }],
        };
      }
      return {};
    },
  });
  vm.runInContext(
    appSource.slice(appSource.indexOf("const BURST"), appSource.indexOf("function setWorld"))
    + "\n;Object.assign(globalThis, { BURST, UNIT, ORES, ENT, NAME, INV_COLOR, DRONE_COLORS });",
    ctx,
  );
  const testColonySrc = source
    .slice(0, source.indexOf('$("#btn-colony").onclick'))
    .replace(
      "return { enter, exit, toggleRun, LEVEL, get running() { return running; } };",
      "return { enter, exit, toggleRun, onEvent, econ, renderAll, setData: (d) => { data = d; }, getData: () => data, treeSnapshot, LEVEL, get running() { return running; } };",
    );
  vm.runInContext(testColonySrc + "\n;globalThis.Colony = Colony;", ctx);

  // Older save: new elements stay hidden
  ctx.Colony.setData({
    state: { tick: 1, size: 4, depth: 1, credits: 100, levels: {}, inv: {}, tiles: [], buildings: [] },
    prices: {}, contracts: [], shop: [], levels: { total: 10 }, tree: [],
  });
  ctx.Colony.renderAll();
  assert.equal(getEl("#col-rank").classList.contains("hidden"), true);
  assert.equal(getEl("#col-sys").classList.contains("hidden"), true);
  assert.equal(getEl("#col-fleet").classList.contains("hidden"), true);

  // Live econ event with new systems populates all readouts
  ctx.Colony.onEvent({
    op: "econ",
    t: 10,
    state: {
      tick: 10, size: 4, depth: 1, credits: 250, levels: {}, inv: { moss: 4 }, tiles: [], buildings: [],
      power: { supply: 12.5, demand: 9, day: true },
      storage: { crops: { used: 40, cap: 200 }, raw: { used: 95, cap: 100 }, crafted: { used: 0, cap: 100 } },
      fleet: [{ id: 0, battery: 80, cargo: { moss: 4 }, cap: 10 }],
      rp: 20, rent: 5, rank: 3,
    },
    rank: { level: 3, kind: "efficiency", target: 1.32, value: 1.1, filled: 0.8 },
  });
  assert.equal(getEl("#col-rank").classList.contains("hidden"), false);
  assert.match(getEl("#col-rank").innerHTML, /Rank 3.*Value per command.*1\.1 \/ 1\.32/);
  assert.equal(getEl("#col-sys").classList.contains("hidden"), false);
  assert.match(getEl("#col-sys").innerHTML, /Day.*12\.5 \/ 9.*20.*RP.*5 CR/);
  assert.match(getEl("#col-sys").innerHTML, /col-bar-row warn" data-storage="raw"/);
  assert.equal(getEl("#col-fleet").classList.contains("hidden"), false);
  assert.match(getEl("#col-fleet").innerHTML, /Drone 0.*80%.*4 \/ 10/);

  // Tree snapshot exposes current rp and rank
  assert.equal(ctx.Colony.treeSnapshot().rp, 20);
  assert.equal(ctx.Colony.treeSnapshot().rank, 3);

  // Rank-up event logs message, plays goal sound, highlights bar, and preserves highlight across econ ticks until timeout
  ctx.TreeView.isOpen = true;
  ctx.Colony.onEvent({
    op: "rank",
    ok: true,
    msg: "Promoted to Rank 4!",
    rank: { level: 4, kind: "output", target: 200, value: 50, filled: 0.2 },
  });
  assert.deepEqual(logs.at(-1), { msg: "Promoted to Rank 4!", cls: "c-ok" });
  assert.equal(sounds.at(-1), "goal");
  assert.equal(getEl("#col-rank").classList.contains("rank-up"), true);
  assert.equal(ctx.Colony.getData().state.rank, 4);
  assert.equal(treeUpdates.at(-1).rank, 4);

  // Subsequent econ event 500ms later preserves rank-up highlight
  now += 500;
  ctx.Colony.onEvent({ op: "econ", t: 15, credits: 260, rank: { level: 4, kind: "output", target: 200, value: 80, filled: 0.4 } });
  assert.equal(getEl("#col-rank").classList.contains("rank-up"), true);

  // After highlight timer fires, class is removed
  const activeTimer = timers.findLast((t) => !t.cleared);
  activeTimer.fn();
  assert.equal(getEl("#col-rank").classList.contains("rank-up"), false);
  await Promise.resolve();
  assert.ok(refreshCount >= 1);
});
