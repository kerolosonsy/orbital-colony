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
    setTimeout: () => 1,
    clearTimeout: () => {},
    NAME,
    INV_COLOR,
    DRONE_COLORS: fullApp.DRONE_COLORS,
  });
  const fmtCode = source.slice(source.indexOf("const fmt ="), source.indexOf("const actionMs"));
  const sysCode = source.slice(source.indexOf("// ---------- system readouts"), source.indexOf("function highlightRank("));
  return vm.runInContext(
    fmtCode + "\n" + sysCode + "\n;({ fmtNum, escAttr, clockFmt, runComparison, droneShares, droneBarsMarkup, lastRunMarkup, liveStatsMarkup, weatherMarkup, rushMarkup, rankGoalLabel, rankMarkup, statusExtraMarkup, hasCargoData, hasFleetAlert, fleetVisible, fleetMarkup, SYSTEM_ORDER, systemsVisible, addrChipMarkup, normalizeSeriesPoints, sparklinePoints, sparklinePath, metricCardMarkup, metricsGridMarkup, systemRowMarkup, systemsMarkup, copyText, seasonComparison, seasonCountdownText, seasonBannerMarkup, seasonsListMarkup, seasonResultInfo, seasonResultMarkup });",
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

test("runComparison handles increase, decrease, flat, zero previous run, and missing previous run", () => {
  const ctx = systemCtx("en");
  assert.equal(ctx.runComparison({ perMin: 30 }, null), null);
  assert.equal(ctx.runComparison({ perMin: 30 }, undefined), null);
  assert.equal(ctx.runComparison({ perMin: 30 }, { perMin: 0 }), null);
  assert.equal(ctx.runComparison({ perMin: 30 }, { perMin: -5 }), null);
  assert.equal(ctx.runComparison(null, { perMin: 20 }), null);

  assert.equal(JSON.stringify(ctx.runComparison({ perMin: 30 }, { perMin: 24 })), JSON.stringify({ pct: 25, text: "+25%", cls: "up" }));
  assert.equal(JSON.stringify(ctx.runComparison({ perMin: 16 }, { perMin: 20 })), JSON.stringify({ pct: -20, text: "-20%", cls: "down" }));
  assert.equal(JSON.stringify(ctx.runComparison({ perMin: 19.999 }, { perMin: 20 })), JSON.stringify({ pct: 0, text: "0%", cls: "flat" }));
});

test("droneShares sums to 100 for any positive distribution and keeps zero shares at zero", () => {
  const ctx = systemCtx("en");
  for (const d of [
    { work: 120, move: 40, wait: 300, failed: 3 },
    { work: 1, move: 1, wait: 1, failed: 0 },
    { work: 1, move: 2, wait: 3, failed: 4 },
    { work: 7, move: 7, wait: 7, failed: 0 },
    { work: 500, move: 0, wait: 0, failed: 1 },
  ]) {
    const s = ctx.droneShares(d);
    assert.equal(s.work + s.move + s.wait + s.failed, 100);
    if (!d.failed) assert.equal(s.failed, 0);
    if (!d.move) assert.equal(s.move, 0);
    if (!d.wait) assert.equal(s.wait, 0);
  }
  const zero = ctx.droneShares({ work: 0, move: 0, wait: 0, failed: 0 });
  assert.equal(JSON.stringify({ work: zero.work, move: zero.move, wait: zero.wait, failed: zero.failed, total: zero.total }), JSON.stringify({ work: 0, move: 0, wait: 0, failed: 0, total: 0 }));
});

test("lastRunMarkup and liveStatsMarkup render comparison, bars summing to 100%, legend, topFail, and idleMoves", () => {
  for (const lang of ["ar", "en"]) {
    const ctx = systemCtx(lang);
    const en = lang === "en";
    assert.equal(ctx.lastRunMarkup([]), "");
    assert.equal(ctx.liveStatsMarkup(null), "");

    const runs = [
      {
        ticks: 1234, earned: 560.5, perMin: 27.2, race: true,
        stats: {
          drones: [
            { id: 0, work: 120, move: 40, wait: 300, failed: 3, idleMoves: 5 },
            { id: 1, work: 10, move: 10, wait: 10, failed: 0, idleMoves: 2 },
          ],
          topFail: "<harvest>",
        },
      },
      { ticks: 1000, earned: 400, perMin: 20 },
    ];
    const html = ctx.lastRunMarkup(runs);
    assert.match(html, en ? /Last run/ : /آخر تشغيل/);
    assert.match(html, /27\.2/);
    assert.match(html, /\+36%/);
    assert.match(html, /1,234/);
    assert.match(html, /560\.5/);
    assert.match(html, /-race/);
    assert.match(html, /&lt;harvest&gt;/);
    assert.doesNotMatch(html, /<harvest>/);
    assert.match(html, /7/); // summed idleMoves 5 + 2
    assert.match(html, en ? /Work.*Move.*Wait.*Failed/ : /عمل.*حركة.*انتظار.*فشل/);

    // Segments in each stacked bar sum to 100%
    for (const bar of html.match(/class="lr-bar"[\s\S]*?<\/div>/g) || []) {
      const shares = [...bar.matchAll(/data-share="(\d+)"/g)].map((m) => Number(m[1]));
      assert.equal(shares.reduce((a, b) => a + b, 0), 100);
    }

    const liveHtml = ctx.liveStatsMarkup(runs[0].stats);
    assert.match(liveHtml, en ? /Drone stats/ : /إحصائيات الدرونز/);
    assert.match(liveHtml, /lr-drones compact/);
  }
});

test("weatherMarkup renders storm warning and active countdowns and hides when none or absent", () => {
  for (const lang of ["ar", "en"]) {
    const ctx = systemCtx(lang);
    const en = lang === "en";
    assert.equal(ctx.weatherMarkup(null, 100), "");
    assert.equal(ctx.weatherMarkup({ storm: "none", at: 1234 }, 100), "");

    const warn = ctx.weatherMarkup({ storm: "warn", at: 1234 }, 1000);
    assert.match(warn, /col-storm is-warn/);
    assert.match(warn, en ? /Storm incoming in 0:23/ : /تحذير عاصفة خلال 0:23/);

    const active = ctx.weatherMarkup({ storm: "active", at: 1500 }, 1200);
    assert.match(active, /col-storm is-active/);
    assert.match(active, en ? /Storm active — clears in 0:30/ : /عاصفة شغالة — تنتهي خلال 0:30/);
  }
});

test("rushMarkup, marketMarkup shock, and fleet faults/stuns render when present and stay clean when absent", () => {
  for (const lang of ["ar", "en"]) {
    const sys = systemCtx(lang);
    const { ctx: mkt } = marketCtx(lang);
    const { NAME, INV_COLOR } = appConstants(lang);
    const en = lang === "en";

    // Rush order rendering and absence
    assert.equal(sys.rushMarkup(null, 100), "");
    const rush = sys.rushMarkup({ good: "moss", amount: 40, left: 25, until: 1600, bonus: 2 }, 1300);
    assert.match(rush, /rush-card/);
    assert.ok(rush.includes(NAME.moss));
    assert.match(rush, new RegExp(INV_COLOR.moss, "i"));
    assert.match(rush, /<b>25<\/b> \/ 40/);
    assert.match(rush, /0:30/);
    assert.match(rush, en ? /Bonus.*2/ : /مكافأة.*2/);

    // Shock in market table rendering and absence
    const plainMkt = mkt.marketMarkup({ inv: { crystal: 5 }, levels: { "crop.crystal": 1 } }, { crystal: 6 });
    assert.doesNotMatch(plainMkt, /market-shock|shock-badge/);
    const shockMkt = mkt.marketMarkup({ inv: {}, levels: {}, tick: 1200, shock: { good: "crystal", factor: 0.4, until: 1500 } }, { crystal: 2.4 });
    assert.match(shockMkt, /class="market-shock"/);
    assert.ok(shockMkt.includes(NAME.crystal));
    assert.match(shockMkt, /×0\.4/);
    assert.match(shockMkt, /0:30/);

    // Single drone with fault or stunned stays visible and displays tags
    assert.equal(sys.fleetVisible([{ id: 0, cap: 0, cargo: {} }]), false);
    assert.equal(sys.fleetVisible([{ id: 0, cap: 0, cargo: {}, fault: true }]), true);
    assert.equal(sys.fleetVisible([{ id: 0, cap: 0, cargo: {}, stunned: 150 }]), true);
    const flt = sys.fleetMarkup([{ id: 1, battery: 75, fault: true, stunned: 150, cap: 10, cargo: {} }]);
    assert.match(flt, /flt-row flt-fault flt-stunned/);
    assert.match(flt, en ? /Fault/ : /عطل/);
    assert.match(flt, en ? /Stunned.*0:15/ : /مصاب.*0:15/);
  }
});

test("Colony handles race toggle persistence, test/vet endpoints, compile-error routing, runs history, world events, and max:0 shop items", async () => {
  const els = new Map();
  const makeEl = () => ({
    innerHTML: "", textContent: "", disabled: false, value: "", readOnly: false,
    attrs: new Map(),
    setAttribute(k, v) { this.attrs.set(k, String(v)); },
    getAttribute(k) { return this.attrs.get(k); },
    querySelector(sel) {
      if (sel === "span") return (this._span ||= { textContent: "" });
      if (sel === "use") return (this._use ||= { attrs: new Map(), setAttribute(k, v) { this.attrs.set(k, String(v)); } });
      return null;
    },
    querySelectorAll() { return []; },
    classList: {
      set: new Set(["hidden"]),
      add(...c) { c.forEach((x) => this.set.add(x)); },
      remove(...c) { c.forEach((x) => this.set.delete(x)); },
      toggle(c, f) { const on = f === undefined ? !this.set.has(c) : !!f; if (on) this.set.add(c); else this.set.delete(c); return on; },
      contains(c) { return this.set.has(c); },
    },
  });
  const getEl = (sel) => { if (!els.has(sel)) els.set(sel, makeEl()); return els.get(sel); };
  const logs = [], sounds = [], prefs = { colonyRace: true }, apiCalls = [];
  let shownFiles = [], lastErrors = [], runResponse = { ok: true }, runsResponses = [[
    { ticks: 900, earned: 300, perMin: 30, race: false, stats: { drones: [{ id: 0, work: 50, move: 25, wait: 25, failed: 0, idleMoves: 2 }], topFail: "" } },
    { ticks: 800, earned: 200, perMin: 20, race: false },
  ]];
  let testResponse = { ok: true, output: "ok colony\n" }, vetResponse = { ok: false, output: "vet failed\n", compileErrors: [{ file: "plan/plan.go", line: 4, col: 2, msg: "unreachable", tip: "remove dead code" }] };

  const colonyFiles = { "main.go": "package main\n", "plan/plan.go": "package plan\n" };
  const ctx = vm.createContext({
    LANG: "en", TOKEN: "tok", PROFILE: "default",
    T: (_ar, en) => en,
    esc: (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"),
    setTimeout: () => 1, clearTimeout: () => {},
    performance: { now: () => 1000 },
    Math, Object, Array, Number, String, Boolean, JSON, Map, Set,
    EventSource: class { close() {} },
    document: { body: makeEl(), querySelector: () => null },
    $: getEl,
    ed: { value: colonyFiles["plan/plan.go"], readOnly: false },
    con: getEl("#console"),
    S: { mode: "mission", file: "plan/plan.go", levels: [], chapters: [] },
    playback: false,
    lastErrors: [],
    lastMtime: 0,
    finishPlayback() {},
    renderLayerBar() {},
    W: { w: 4, h: 4, depth: 1, view: 0, follow: true, colony: true, live: true, t: 0, inv: {}, layers: [[]], drones: new Map() },
    TreeView: { isOpen: false, update() {}, close() {}, open() {} },
    Sound: { play: (name) => sounds.push(name) },
    log: (msg, cls) => logs.push({ msg, cls }),
    hud() {}, setWorld() {}, renderLegend() {}, draw() {}, linkUI() {},
    isLinked: () => false,
    syncFromWorkspace: async () => {},
    levelFiles: () => [{ name: "main.go" }, { name: "plan/plan.go" }],
    fileCode: (_l, f) => colonyFiles[typeof f === "string" ? f : f.name] || "",
    saveFile: (_l, name, code) => { colonyFiles[name] = code; },
    colonyFilesPayload: () => ({ ...colonyFiles }),
    firstErrorFile: (errs, files, fb) => {
      const raw = errs && errs[0] ? (errs[0].file || "main.go").replace(/^\.\//, "") : "";
      return files.some((f) => f.name === raw) ? raw : fb;
    },
    showFile: (f) => { shownFiles.push(f); ctx.S.file = f; },
    rec: () => ({}),
    setDigits: (el, v) => { el.textContent = String(v); },
    store: {
      pref: (k, v) => (v === undefined ? prefs[k] : (prefs[k] = v)),
      save() {},
    },
    api: async (path, body) => {
      apiCalls.push({ path, body });
      if (path === "/api/colony") {
        return {
          running: false, starter: "package main\n", prices: {}, contracts: [],
          levels: { total: 0 },
          shop: [
            {
              id: "endless.drone", name: "Drone+", nameEn: "Drone+", desc: "More", descEn: "More",
              group: "Fleet", groupEn: "Fleet", max: 0, level: 4,
              cost: { credits: 100, goods: { plate: 2 } },
            },
          ],
          state: {
            tick: 100, size: 4, depth: 1, credits: 500, levels: { "endless.drone": 4 },
            inv: { plate: 5 }, tiles: [], buildings: [],
          },
        };
      }
      if (path === "/api/colony/runs") return { runs: runsResponses[0] || [] };
      if (path === "/api/colony/run") return runResponse;
      if (path === "/api/colony/test") return testResponse;
      if (path === "/api/colony/vet") return vetResponse;
      if (path === "/api/colony/stop") return { ok: true };
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
      "return { enter, exit, toggleRun, toggleRace, runCheck, onEvent, status, renderShop, getData: () => data, LEVEL, get running() { return running; }, get raced() { return raced; }, get raceEnabled() { return raceEnabled; } };",
    );
  vm.runInContext(testColonySrc + "\n;globalThis.Colony = Colony;", ctx);

  await ctx.Colony.enter();
  await Promise.resolve();
  // Enter restored persisted raceEnabled and fetched /api/colony/runs -> showing last run card
  assert.equal(ctx.Colony.raceEnabled, true);
  assert.equal(getEl("#btn-race").getAttribute("aria-pressed"), "true");
  assert.equal(getEl("#col-lastrun").classList.contains("hidden"), false);
  assert.match(getEl("#col-lastrun").innerHTML, /\+50%/);

  // Endless shop item (max: 0) is never maxed, shows plain level, omits /total, keeps buy button + colour chip
  assert.equal(ctx.Colony.status(ctx.Colony.getData().shop[0]).can, true);
  const shopHtml = getEl("#col-shop").innerHTML;
  assert.match(shopHtml, /<small>4<\/small>/);
  assert.match(shopHtml, /class="lv">4</);
  assert.match(shopHtml, /<button class="buy" data-id="endless\.drone" >/);
  assert.match(shopHtml, new RegExp(appConstants("en").INV_COLOR.plate, "i"));
  assert.doesNotMatch(shopHtml, /Maxed/);

  // Clicking Run sends { files, race: true } and clears any prior raced mark
  await ctx.Colony.toggleRun();
  const runCall = apiCalls.findLast((c) => c.path === "/api/colony/run");
  assert.deepEqual(Object.keys(runCall.body.files).sort(), ["main.go", "plan/plan.go"]);
  assert.equal(runCall.body.race, true);
  assert.equal(ctx.Colony.running, true);
  assert.equal(getEl("#col-lastrun").classList.contains("hidden"), true);

  // Race live events play error sound ONCE per run and mark state pill raced
  sounds.length = 0;
  ctx.Colony.onEvent({ op: "race", msg: "WARNING: DATA RACE" });
  ctx.Colony.onEvent({ op: "race", msg: "WARNING: DATA RACE 2" });
  assert.equal(sounds.filter((s) => s === "error").length, 1);
  assert.equal(ctx.Colony.raced, true);
  assert.match(getEl("#col-state").textContent, /Raced/);

  // Storm, fault, shock, rush events log and update UI
  ctx.Colony.onEvent({ op: "storm", phase: "warn", msg: "Storm coming", t: 110, at: 310 });
  assert.deepEqual(logs.at(-1), { msg: "Storm coming", cls: "c-warn" });
  assert.equal(ctx.W.storm, 1);
  assert.match(getEl("#col-weather").innerHTML, /0:20/);

  ctx.Colony.onEvent({ op: "fault", d: 0, stunned: 120, msg: "Drone 0 faulted", state: { ...ctx.Colony.getData().state, fleet: [{ id: 0, battery: 50, fault: true, stunned: 120 }] } });
  assert.deepEqual(logs.at(-1), { msg: "Drone 0 faulted", cls: "c-err" });
  assert.match(getEl("#col-fleet").innerHTML, /Fault.*Stunned.*0:12/);

  ctx.Colony.onEvent({ op: "shock", msg: "Crystal market shock!", shock: { good: "crystal", factor: 0.4, until: 400 }, prices: { crystal: 2.4 } });
  assert.deepEqual(logs.at(-1), { msg: "Crystal market shock!", cls: "c-warn" });
  assert.equal(sounds.at(-1), "fail");
  assert.match(getEl("#col-market").innerHTML, /Shock ×0\.4/);

  ctx.Colony.onEvent({ op: "rush", msg: "Urgent moss order!", rush: { good: "moss", amount: 40, left: 25, until: 500, bonus: 2 } });
  assert.deepEqual(logs.at(-1), { msg: "Urgent moss order!", cls: "c-ok" });
  assert.equal(sounds.at(-1), "report");
  assert.equal(getEl("#col-rush").classList.contains("hidden"), false);
  assert.match(getEl("#col-rush").innerHTML, /Rush order[\s\S]*<b>25<\/b> \/ 40/);

  // Stopped event updates runs history and unhides last-run card
  runsResponses[0] = [{ ticks: 500, earned: 400, perMin: 40, race: true, stats: { drones: [{ id: 0, work: 80, move: 20, wait: 0, failed: 0, idleMoves: 0 }], topFail: "harvest" } }];
  ctx.Colony.onEvent({ op: "stopped", msg: "done" });
  await Promise.resolve();
  assert.equal(ctx.Colony.running, false);
  assert.equal(getEl("#col-lastrun").classList.contains("hidden"), false);

  // Test & Vet endpoints route compile errors to the failing file
  await ctx.Colony.runCheck("test");
  assert.equal(logs.at(-1).cls, "c-ok");
  await ctx.Colony.runCheck("vet");
  assert.equal(shownFiles.at(-1), "plan/plan.go");
  assert.equal(logs.some((l) => l.msg.includes("plan/plan.go:4:2") && l.cls === "c-err"), true);
  assert.equal(logs.some((l) => l.msg.includes("remove dead code") && l.cls === "c-tip"), true);

  // Toggling race off persists preference and updates payload
  ctx.Colony.toggleRace();
  assert.equal(prefs.colonyRace, false);
  assert.equal(getEl("#btn-race").getAttribute("aria-pressed"), "false");
  runResponse = { phase: "compile", result: { message: "Import disallowed" }, compileErrors: [{ file: "plan/plan.go", line: 2, col: 1, msg: "bad import" }] };
  shownFiles.length = 0;
  await ctx.Colony.toggleRun();
  assert.equal(apiCalls.findLast((c) => c.path === "/api/colony/run").body.race, false);
  assert.equal(ctx.Colony.raced, false);
  assert.equal(shownFiles.at(-1), "plan/plan.go");
  assert.equal(logs.some((l) => l.msg === "Import disallowed" && l.cls === "c-err"), true);
});

test("systemsVisible and systemsMarkup stay hidden when off and render every system cleanly in both languages", async () => {
  for (const lang of ["ar", "en"]) {
    const ctx = systemCtx(lang);
    const en = lang === "en";
    assert.equal(ctx.systemsVisible(null), false);
    assert.equal(ctx.systemsVisible(undefined), false);
    assert.equal(ctx.systemsVisible({}), false);
    assert.equal(ctx.systemsVisible({ archive: { on: false, rows: 10 }, market: { on: false } }), false);
    assert.equal(ctx.systemsMarkup({ archive: { on: false } }), "");

    const html = ctx.systemsMarkup({
      archive: { on: true, rows: 120 },
      market: { on: true, url: "http://127.0.0.1:41233?a=1&b=<2>", orders: 12, rejected: 3 },
      control: { on: true, addr: "127.0.0.1:41240", healthy: true, lastMs: 12, probes: 40, failed: 1 },
      warehouse: { on: true, restarts: 2 },
      metrics: { on: true, series: { colony_credits: [[1200, 698.5], [1250, 701]] } },
      neighbor: { on: true, addr: "127.0.0.1:41250", trades: 4, drops: 2 },
    });

    assert.match(html, en ? /Systems/ : /الأنظمة/);
    assert.match(html, en ? /Archive/ : /الأرشيف/);
    assert.match(html, en ? /Rows/ : /الصفوف/);
    assert.match(html, /120/);

    assert.match(html, en ? /Market/ : /السوق/);
    assert.match(html, en ? /Orders/ : /الطلبات/);
    assert.match(html, /12/);
    assert.match(html, /3/);
    assert.ok(html.includes("http://127.0.0.1:41233?a=1&amp;b=&lt;2&gt;"));
    assert.match(html, /class="sys-copy"/);

    assert.match(html, en ? /Control Centre/ : /مركز القيادة/);
    assert.match(html, en ? /Healthy/ : /سليم/);
    assert.match(html, /12 ms/);
    assert.match(html, /40/);
    assert.match(html, /1/);
    assert.ok(html.includes("127.0.0.1:41240"));

    assert.match(html, en ? /Warehouse/ : /المخزن/);
    assert.match(html, /2/);

    assert.match(html, en ? /Metrics/ : /المراقبة/);
    assert.match(html, /colony_credits/);
    assert.match(html, /701/);

    assert.match(html, en ? /Neighbour/ : /الجيران/);
    assert.match(html, /127\.0\.0\.1:41250/);
    assert.match(html, /4/);

    const badCtrl = ctx.systemRowMarkup("control", { on: true, healthy: false, probes: 5, failed: 5 });
    assert.match(badCtrl, /sys-health is-bad/);
    assert.match(badCtrl, en ? /Unhealthy/ : /عطلان/);

    const unkCtrl = ctx.systemRowMarkup("control", { on: true });
    assert.match(unkCtrl, /sys-health is-unknown/);
    assert.match(unkCtrl, en ? /Unknown/ : /غير معروف/);
  }

  // copyText handles Clipboard API success and missing/rejected fallback silently
  const ctx = systemCtx("en");
  let copied = "";
  const btn = { textContent: "Copy", dataset: {}, classList: { add() {}, remove() {} } };
  assert.equal(await ctx.copyText("127.0.0.1:41240", btn, { clipboard: { writeText: async (v) => { copied = v; } } }), true);
  assert.equal(copied, "127.0.0.1:41240");
  assert.equal(btn.textContent, "Copied");
  assert.equal(await ctx.copyText("127.0.0.1:41240", null, { clipboard: null }), false);
  assert.equal(await ctx.copyText("127.0.0.1:41240", null, { clipboard: { writeText: () => Promise.reject(new Error("denied")) } }), false);
});

test("sparklinePath, point/series limits, and under-two-points metrics rendering", () => {
  const ctx = systemCtx("en");
  assert.equal(ctx.sparklinePath([]), "");
  assert.equal(ctx.sparklinePath([[10, 50]]), "");
  assert.equal(ctx.sparklinePath([[0, 0], [10, 100]], 120, 28, 2), "M2 26 L118 2");
  assert.equal(ctx.sparklinePath([[0, 50], [5, 100], [10, 0]], 120, 28, 2), "M2 14 L60 2 L118 26");
  // Flat series stays along vertical midline; duplicate ticks spread evenly in X
  assert.equal(ctx.sparklinePath([[10, 42], [20, 42]], 120, 28, 2), "M2 14 L118 14");
  assert.equal(ctx.sparklinePath([[10, 0], [10, 100]], 120, 28, 2), "M2 26 L118 2");

  // Fewer than two points shows value (or dash) without an inline SVG
  const single = ctx.metricCardMarkup("solar_w", [[1200, 48.5]]);
  assert.ok(!single.includes("<svg"));
  assert.match(single, /solar_w/);
  assert.match(single, /48\.5/);
  const empty = ctx.metricCardMarkup("empty_metric", []);
  assert.ok(!empty.includes("<svg"));
  assert.match(empty, /—/);

  // Latest 60 points limit and max 8 series limit are enforced
  const manyPts = Array.from({ length: 75 }, (_, i) => [i * 10, i + 1]);
  const norm = ctx.normalizeSeriesPoints(manyPts);
  assert.equal(norm.length, 60);
  assert.deepEqual(Array.from(norm[0]), [150, 16]);
  assert.deepEqual(Array.from(norm[59]), [740, 75]);

  const tenSeries = {};
  for (let i = 0; i < 10; i++) tenSeries[`m_${i}<tag>`] = [[10, i], [20, i + 1]];
  const grid = ctx.metricsGridMarkup(tenSeries);
  assert.equal((grid.match(/data-metric="/g) || []).length, 8);
  assert.ok(!grid.includes("m_8"));
  assert.ok(grid.includes("m_0&lt;tag&gt;"));
});

test("seasons markup covers list comparison, missing best, countdown, running disable, and result card", () => {
  for (const lang of ["ar", "en"]) {
    const ctx = systemCtx(lang);
    const en = lang === "en";

    assert.equal(ctx.seasonCountdownText(6000), "10:00");
    assert.equal(ctx.seasonCountdownText(4200), "7:00");
    assert.equal(ctx.seasonCountdownText(-5), "0:00");

    const listStopped = ctx.seasonsListMarkup({
      seasons: [
        { id: "one-drone", title: "One Drone <1>", desc: "Run & grow", ticks: 6000, best: 1234, reference: 2000 },
        { id: "night-rush", title: "Night Rush", desc: "Survive storms", ticks: 3600, best: null, reference: 1500 },
      ],
      active: null,
    }, false);

    assert.ok(listStopped.includes("One Drone &lt;1&gt;"));
    assert.ok(listStopped.includes("Run &amp; grow"));
    assert.match(listStopped, /10:00/);
    assert.match(listStopped, /1,234/);
    assert.match(listStopped, /2,000/);
    assert.match(listStopped, /width:61\.7%/);
    assert.match(listStopped, />—</);
    assert.match(listStopped, /width:0\.0%/);
    assert.ok(!/data-season-start="one-drone"[^>]*disabled/.test(listStopped));

    const listRunning = ctx.seasonsListMarkup({
      seasons: [{ id: "one-drone", title: "One Drone", desc: "Desc", ticks: 6000, best: 2100, reference: 2000 }],
      active: null,
    }, true);
    assert.match(listRunning, /season-run-warn/);
    assert.match(listRunning, /season-why/);
    assert.match(listRunning, /is-met/);
    assert.match(listRunning, /width:100\.0%/);
    assert.match(listRunning, /data-season-start="one-drone"\s*disabled/);

    const banner = ctx.seasonBannerMarkup({ id: "one-drone", left: 4200, score: 310 }, { title: "One Drone" }, false);
    assert.match(banner, /One Drone/);
    assert.match(banner, /7:00/);
    assert.match(banner, /310/);
    assert.match(banner, /btn-leave-season/);

    const newBestInfo = ctx.seasonResultInfo({ op: "season", ok: true, score: 1234, best: 1234, msg: "Nice!" }, { title: "One Drone", best: 900, reference: 2000 });
    assert.equal(newBestInfo.newBest, true);
    const newBestHtml = ctx.seasonResultMarkup({ op: "season", ok: true, score: 1234, best: 1234, msg: "Nice <1>!" }, { title: "One Drone", best: 900, reference: 2000 }, {}, newBestInfo);
    assert.match(newBestHtml, en ? /New best!/ : /رقم قياسي جديد!/);
    assert.match(newBestHtml, /1,234 \/ 2,000 \(62%\)/);
    assert.ok(newBestHtml.includes("Nice &lt;1&gt;!"));

    const notBestInfo = ctx.seasonResultInfo({ op: "season", ok: true, score: 900, best: 1500, msg: "Done" }, { title: "One Drone", best: 1500, reference: 2000 });
    assert.equal(notBestInfo.newBest, false);
    const notBestHtml = ctx.seasonResultMarkup({ op: "season", ok: true, score: 900, best: 1500, msg: "Done" }, { title: "One Drone", best: 1500, reference: 2000 }, {}, notBestInfo);
    assert.ok(!notBestHtml.includes("season-new-best"));
    assert.match(notBestHtml, /1,500/);
  }
});

test("Colony seasons and systems end-to-end: title distinction, running disable, countdown, leave confirm, and season result sound", async () => {
  const elements = new Map();
  const makeEl = (id = "") => {
    const classes = new Set();
    const attrs = new Map();
    return {
      id,
      innerHTML: "",
      textContent: "",
      value: "",
      open: false,
      disabled: false,
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
      setAttribute: (k, v) => attrs.set(k, String(v)),
      getAttribute: (k) => (attrs.has(k) ? attrs.get(k) : null),
      querySelector: () => makeEl(),
      querySelectorAll: () => [],
      showModal() { this.open = true; classes.delete("hidden"); },
      close() { this.open = false; classes.add("hidden"); },
    };
  };
  const getEl = (sel) => {
    if (!elements.has(sel)) elements.set(sel, makeEl(sel));
    return elements.get(sel);
  };

  const apiCalls = [];
  const logs = [];
  const sounds = [];
  const confirmCalls = [];
  let confirmReply = false;
  let colonyState = {
    running: false,
    starter: "package main\n",
    prices: {},
    contracts: [],
    shop: [],
    systems: { archive: { on: true, rows: 55 } },
    season: null,
    state: { tick: 100, size: 4, depth: 1, credits: 500, levels: {}, inv: {}, tiles: [], buildings: [] },
  };
  const seasonsData = {
    seasons: [
      { id: "one-drone", title: "One Drone Challenge", desc: "Build with one drone", ticks: 6000, best: 1000, reference: 2000 },
    ],
    active: null,
  };

  const doc = { title: "Orbital Colony", body: makeEl("body"), querySelector: () => null };
  const ctx = vm.createContext({
    LANG: "en", TOKEN: "tok", PROFILE: "default",
    T: (_ar, en) => en,
    esc: (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"),
    setTimeout: () => 1,
    clearTimeout: () => {},
    performance: { now: () => 1000 },
    Math, Object, Array, Number, String, Boolean, JSON, Map, Set,
    EventSource: class { close() {} },
    document: doc,
    confirm: (msg) => { confirmCalls.push(msg); return confirmReply; },
    $: getEl,
    ed: { value: "", readOnly: false },
    con: getEl("#console"),
    S: { mode: "mission", file: "main.go", levels: [], chapters: [] },
    playback: false,
    lastErrors: [],
    lastMtime: 0,
    finishPlayback() {},
    renderLayerBar() {},
    apply() {},
    showOverlay: (sel) => getEl(sel).showModal(),
    closeOverlay: (el) => el && el.close(),
    W: { w: 4, h: 4, depth: 1, view: 0, follow: true, colony: true, live: true, t: 0, inv: {}, layers: [[]], drones: new Map() },
    TreeView: { isOpen: false, update() {}, close() {}, open() {} },
    Sound: { play: (name) => sounds.push(name) },
    log: (msg, cls) => logs.push({ msg, cls }),
    hud() {}, setWorld() {}, renderLegend() {}, draw() {}, linkUI() {},
    isLinked: () => false,
    levelFiles: () => [{ name: "main.go" }],
    fileCode: () => "package main\n",
    saveFile() {},
    showFile() {},
    rec: () => ({}),
    setDigits: (el, v) => { el.textContent = String(v); },
    store: { pref: () => false, save() {} },
    api: async (path, body) => {
      apiCalls.push({ path, body });
      if (path === "/api/colony") return JSON.parse(JSON.stringify(colonyState));
      if (path === "/api/colony/runs") return { runs: [] };
      if (path === "/api/colony/seasons") return JSON.parse(JSON.stringify(seasonsData));
      if (path === "/api/colony/season") {
        if (colonyState.running) throw new Error("409 colony running");
        const nextId = body && body.id ? body.id : "";
        seasonsData.active = nextId || null;
        colonyState = {
          ...colonyState,
          season: nextId ? { id: nextId, left: 4200, score: 310 } : null,
          state: { ...colonyState.state, tick: 200 },
        };
        return { ok: true };
      }
      if (path === "/api/colony/run") {
        colonyState.running = true;
        return { started: true };
      }
      if (path === "/api/colony/stop") {
        colonyState.running = false;
        return { ok: true };
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
      "return { enter, exit, toggleRun, onEvent, openSeasons, switchSeason, leaveSeason, LEVEL, get running() { return running; } };",
    );
  vm.runInContext(testColonySrc + "\n;globalThis.Colony = Colony;", ctx);

  await ctx.Colony.enter();
  assert.equal(getEl("#col-systems").classList.contains("hidden"), false);
  assert.match(getEl("#col-systems").innerHTML, /Archive[\s\S]*55/);
  assert.equal(getEl("#col-season").classList.contains("hidden"), true);
  assert.equal(doc.body.classList.contains("colony-season"), false);
  assert.equal(getEl("#lvl-title").textContent, "Colony");

  // Live systems event turning all systems off hides #col-systems
  ctx.Colony.onEvent({ op: "systems", systems: { archive: { on: false } } });
  assert.equal(getEl("#col-systems").classList.contains("hidden"), true);

  // Starting a season while running is rejected cleanly and seasons dialog marks buttons disabled
  await ctx.Colony.toggleRun();
  assert.equal(ctx.Colony.running, true);
  await ctx.Colony.openSeasons();
  assert.match(getEl("#seasons-list").innerHTML, /data-season-start="one-drone"\s*disabled/);
  await ctx.Colony.switchSeason("one-drone");
  assert.equal(apiCalls.some((c) => c.path === "/api/colony/season"), false);

  // Stop, start season -> banner appears and header title area distinguishes season mode
  await ctx.Colony.toggleRun();
  await ctx.Colony.switchSeason("one-drone");
  assert.deepEqual(JSON.parse(JSON.stringify(apiCalls.findLast((c) => c.path === "/api/colony/season").body)), { id: "one-drone" });
  assert.equal(doc.body.classList.contains("colony-season"), true);
  assert.match(getEl("#lvl-title").textContent, /Season: One Drone Challenge/);
  assert.match(getEl("#lvl-chapter").textContent, /ONE-DRONE/);
  assert.match(doc.title, /Season: One Drone Challenge/);
  assert.equal(getEl("#col-season").classList.contains("hidden"), false);
  assert.match(getEl("#col-season").innerHTML, /One Drone Challenge/);
  assert.match(getEl("#col-season").innerHTML, /7:00/);
  assert.match(getEl("#col-season").innerHTML, /310/);

  // Running live tick decrements countdown in #col-season
  await ctx.Colony.toggleRun();
  ctx.Colony.onEvent({ op: "tick", t: 250 });
  assert.match(getEl("#col-season").innerHTML, /6:55/);

  // Cancelling leave confirmation aborts without API call; confirming leaves season when stopped
  confirmReply = false;
  assert.equal(await ctx.Colony.leaveSeason(), false);
  assert.equal(confirmCalls.length, 1);
  assert.equal(apiCalls.filter((c) => c.path === "/api/colony/season").length, 1);

  await ctx.Colony.toggleRun();
  confirmReply = true;
  assert.equal(await ctx.Colony.leaveSeason(), true);
  assert.deepEqual(JSON.parse(JSON.stringify(apiCalls.findLast((c) => c.path === "/api/colony/season").body)), { id: "" });
  assert.equal(doc.body.classList.contains("colony-season"), false);
  assert.equal(doc.title, "Orbital Colony");
  assert.equal(getEl("#lvl-title").textContent, "Colony");
  assert.equal(getEl("#col-season").classList.contains("hidden"), true);

  // Season-end event: new best plays goal sound and renders new-best result card; non-new-best plays report
  seasonsData.seasons[0].best = 1000;
  await ctx.Colony.openSeasons();
  sounds.length = 0;
  ctx.Colony.onEvent({ op: "season", id: "one-drone", ok: true, score: 1500, best: 1500, msg: "Season finished!" });
  await Promise.resolve();
  assert.equal(sounds.includes("goal"), true);
  assert.match(getEl("#season-result-body").innerHTML, /New best!/);
  assert.match(getEl("#season-result-body").innerHTML, /1,500 \/ 2,000 \(75%\)/);

  sounds.length = 0;
  ctx.Colony.onEvent({ op: "season", id: "one-drone", ok: true, score: 900, best: 1500, msg: "Finished" });
  await Promise.resolve();
  assert.equal(sounds.includes("goal"), false);
  assert.equal(sounds.includes("report"), true);
  assert.ok(!getEl("#season-result-body").innerHTML.includes("season-new-best"));
});
