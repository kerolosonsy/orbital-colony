"use strict";
const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");

const spritesSrc = fs.readFileSync(__dirname + "/sprites.js", "utf8");
const buildingsSrc = fs.readFileSync(__dirname + "/buildings.js", "utf8");
const appSrc = fs.readFileSync(__dirname + "/app.js", "utf8");
const colonySrc = fs.readFileSync(__dirname + "/colony.js", "utf8");

const BUILDING_KINDS = [
  "composter", "smelter", "kiln", "refinery", "assembler",
  "reactor-lab", "silo", "generator", "lab",
];

function makeFakeDocument() {
  return {
    createElement(tag) {
      assert.equal(tag, "canvas");
      const ops = [];
      const pixels = new Map();
      const ctx = {
        fillStyle: "#000000",
        strokeStyle: "#000000",
        lineWidth: 1,
        imageSmoothingEnabled: true,
        ops,
        pixels,
        save() { ops.push(["save"]); },
        restore() { ops.push(["restore"]); },
        beginPath() { ops.push(["beginPath"]); },
        closePath() { ops.push(["closePath"]); },
        moveTo(x, y) { ops.push(["moveTo", x, y]); },
        lineTo(x, y) { ops.push(["lineTo", x, y]); },
        ellipse(...args) { ops.push(["ellipse", ...args]); },
        fill() { ops.push(["fill", this.fillStyle]); },
        stroke() { ops.push(["stroke", this.strokeStyle, this.lineWidth]); },
        drawImage(img, x, y, w, h) { ops.push(["drawImage", img, x, y, w, h]); },
        fillRect(x, y, w, h) {
          ops.push(["fillRect", x, y, w, h, this.fillStyle]);
          if (w === 1 && h === 1) pixels.set(`${x},${y}`, this.fillStyle);
        },
      };
      return {
        width: 0,
        height: 0,
        _ctx: ctx,
        getContext() { return ctx; },
      };
    },
  };
}

function loadBuildingEnv(lang = "en") {
  const doc = makeFakeDocument();
  const context = vm.createContext({
    LANG: lang,
    T: (ar, en) => (lang === "en" ? en : ar),
    esc: (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"),
    document: doc,
    performance: { now: () => 1000 },
  });
  vm.runInContext(
    appSrc.slice(appSrc.indexOf("const BURST"), appSrc.indexOf("function setWorld"))
    + "\n;Object.assign(globalThis, { BURST, UNIT, ORES, ENT, NAME, INV_COLOR });",
    context,
  );
  vm.runInContext(spritesSrc + "\n;globalThis.Sprites = Sprites;", context);
  vm.runInContext(buildingsSrc + "\n;globalThis.Buildings = Buildings;", context);
  return context;
}

function makeFakeElement(id = "") {
  const classes = new Set();
  return {
    id,
    innerHTML: "",
    textContent: "",
    hidden: false,
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
    querySelector: () => makeFakeElement(),
    querySelectorAll: () => [],
  };
}

function loadLiveColonyEnv(lang = "en") {
  const elements = new Map();
  const getEl = (sel) => {
    if (!elements.has(sel)) elements.set(sel, makeFakeElement(sel));
    return elements.get(sel);
  };
  const doc = {
    ...makeFakeDocument(),
    body: { classList: makeFakeElement().classList },
  };
  const logs = [];
  const sounds = [];
  const fx = [];
  let legendCalls = 0;
  let drawCalls = 0;
  let now = 1200;

  const context = vm.createContext({
    LANG: lang,
    T: (ar, en) => (lang === "en" ? en : ar),
    esc: (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"),
    document: doc,
    performance: { now: () => now },
    setNow: (v) => { now = v; },
    $: getEl,
    $$: () => [],
    S: { mode: "colony", levels: [], chapters: [], failLogged: false },
    W: null,
    fx,
    TreeView: { isOpen: false, update() {}, close() {}, open() {} },
    Sound: { play: (name) => sounds.push(name) },
    log: (msg, cls) => logs.push({ msg, cls }),
    hud: () => {},
    burst: () => {},
    setDigits: (el, v) => { el.textContent = String(v); },
    renderLegend: () => { legendCalls += 1; },
    renderLayerBar: () => {},
    draw: () => { drawCalls += 1; },
    setWorld: () => {},
    parText: () => {},
    lastInv: {},
  });

  vm.runInContext(
    appSrc.slice(appSrc.indexOf("const BURST"), appSrc.indexOf("function setWorld"))
    + "\n;Object.assign(globalThis, { BURST, UNIT, ORES, ENT, NAME, INV_COLOR, noEmoji });",
    context,
  );
  vm.runInContext(buildingsSrc + "\n;globalThis.Buildings = Buildings;", context);
  vm.runInContext(
    appSrc.slice(appSrc.indexOf("const DRONE_OPS"), appSrc.indexOf("function tick(")),
    context,
  );

  // Expose internal Colony handlers cleanly in-memory for SSE integration testing.
  const testColonySrc = colonySrc
    .slice(0, colonySrc.indexOf('$("#btn-colony").onclick'))
    .replace(
      "return { enter, exit, toggleRun, LEVEL, get running() { return running; } };",
      "return { enter, exit, toggleRun, onEvent, econ, renderBuildings, renderMarket, setData: (d) => { data = d; }, getData: () => data, LEVEL, get running() { return running; } };",
    );
  vm.runInContext(testColonySrc + "\n;globalThis.Colony = Colony;", context);

  return {
    context,
    elements,
    logs,
    sounds,
    get legendCalls() { return legendCalls; },
    get drawCalls() { return drawCalls; },
  };
}

test("Buildings status map tracks syncState, live building events, removal, and tile cleanup", () => {
  const ctx = loadBuildingEnv("en");
  const { Buildings } = ctx;
  assert.deepEqual(Array.from(Buildings.KINDS), BUILDING_KINDS);

  const world = {
    layers: [[
      [{ x: 0, y: 0, g: "regolith" }, { x: 0, y: 1, g: "regolith" }, { x: 0, y: 2, g: "regolith" }],
      [{ x: 1, y: 0, g: "regolith" }, { x: 1, y: 1, g: "regolith" }, { x: 1, y: 2, g: "regolith" }],
      [{ x: 2, y: 0, g: "regolith" }, { x: 2, y: 1, g: "regolith" }, { x: 2, y: 2, g: "regolith" }],
    ]],
  };

  Buildings.syncState([
    { kind: "smelter", x: 1, y: 2, level: 2, recipe: "plate", in: { iron: 2 }, out: { plate: 1 }, progress: 0.75, wear: 81 },
    { kind: "silo", x: 0, y: 0, level: 1, recipe: "", in: {}, out: {}, progress: 0, wear: 12 },
  ]);
  assert.equal(Buildings.list().length, 2);
  // Sorted by y, then x.
  assert.deepEqual(Array.from(Buildings.list(), (b) => [b.x, b.y, b.kind]), [[0, 0, "silo"], [1, 2, "smelter"]]);
  assert.equal(Buildings.at(1, 2).progress, 0.75);

  Buildings.applyEvent({
    op: "building",
    t: 123,
    x: 2,
    y: 1,
    building: { kind: "composter", recipe: "biomass", in: { moss: 4 }, out: { biomass: 1 }, progress: 0.4, wear: 3, level: 1 },
  }, world);

  assert.equal(Buildings.list().length, 3);
  assert.equal(world.layers[0][2][1].b, "composter");
  assert.deepEqual(JSON.parse(JSON.stringify(Buildings.at(2, 1))), {
    kind: "composter",
    x: 2,
    y: 1,
    level: 1,
    recipe: "biomass",
    in: { moss: 4 },
    out: { biomass: 1 },
    progress: 0.4,
    wear: 3,
  });

  // Clearing a building event clears both the status map and tile.b
  Buildings.applyEvent({ op: "building", x: 2, y: 1, building: null }, world);
  assert.equal(Buildings.at(2, 1), null);
  assert.equal("b" in world.layers[0][2][1], false);

  Buildings.clear();
  assert.deepEqual(Array.from(Buildings.list()), []);
});

test("Buildings.listMarkup renders localized names, positions, I/O chips, wear > 80 warning, and escapes HTML", () => {
  for (const lang of ["ar", "en"]) {
    const ctx = loadBuildingEnv(lang);
    const { Buildings, NAME, INV_COLOR } = ctx;
    assert.equal(Buildings.listMarkup([]), "");

    const html = Buildings.listMarkup([
      { kind: "composter", x: 2, y: 1, level: 1, recipe: "biomass", in: { moss: 4 }, out: { biomass: 1 }, progress: 0.4, wear: 80 },
      { kind: "reactor-lab", x: 0, y: 3, level: 2, recipe: "core", in: { uranium: 2, circuit: 1 }, out: {}, progress: 0.2, wear: 81 },
      { kind: "<img src=x onerror=1>", x: 1, y: 1, level: 1, recipe: "<script>", in: { "x\"y": 3 }, out: { plate: 2 }, progress: 0, wear: 5 },
    ]);

    assert.match(html, lang === "en" ? /<h3>Buildings <small>3<\/small><\/h3>/ : /<h3>المباني <small>3<\/small><\/h3>/);
    assert.match(html, new RegExp(NAME.composter));
    assert.match(html, new RegExp(NAME["reactor-lab"]));
    assert.match(html, new RegExp(NAME.biomass));
    assert.match(html, new RegExp(NAME.core));
    assert.match(html, /\(2, 1\)/);
    assert.match(html, /\(0, 3\)/);
    assert.match(html, new RegExp(INV_COLOR.moss));
    assert.match(html, new RegExp(INV_COLOR.biomass));
    assert.match(html, new RegExp(INV_COLOR.circuit));
    assert.match(html, new RegExp(INV_COLOR.uranium));
    assert.match(html, /<span class="bld-none">—<\/span>/);

    // Wear warning applies strictly for wear > 80 (80 -> normal, 81 -> worn).
    const rows = html.split('<div class="item bld-row').slice(1);
    assert.equal(rows.length, 3);
    assert.equal(rows[0].startsWith(' bld-worn"'), false);
    assert.doesNotMatch(rows[0], /class="bld-wear warn"/);
    assert.equal(rows[1].startsWith(' bld-worn"'), true);
    assert.match(rows[1], /class="bld-wear warn"/);
    assert.match(rows[1], /81%/);

    // Unsafe strings are escaped.
    assert.doesNotMatch(html, /<img /);
    assert.doesNotMatch(html, /<script>/);
    assert.match(html, /&lt;img src=x onerror=1&gt;/);
    assert.match(html, /&lt;script&gt;/);
  }
});

test("Sprites.ore and Buildings.sprite produce distinct procedural pixel art and draw handles bars & wear", () => {
  const ctx = loadBuildingEnv("en");
  const { Sprites, Buildings } = ctx;

  // All 6 underground ores draw non-empty pixels and distinct palettes.
  const orePalettes = new Map();
  for (const ore of ["iron", "geode", "copper", "sulfur", "titanium", "uranium"]) {
    const c = Sprites.ore(ore, 1);
    const cols = new Set(c._ctx.pixels.values());
    assert.ok(c._ctx.pixels.size > 0, `ore ${ore} should draw pixels`);
    assert.ok(cols.size >= 3, `ore ${ore} should use shaded tones`);
    orePalettes.set(ore, [...cols].sort().join(","));
  }
  assert.equal(new Set(orePalettes.values()).size, 6, "each ore must have its own colour palette");

  // All 9 buildings draw non-empty 16x16 grids with unique silhouettes.
  const signatures = new Set();
  for (const kind of BUILDING_KINDS) {
    const idle = Buildings.sprite(kind, 0);
    const active = Buildings.sprite(kind, 1);
    assert.equal(idle.width, 16);
    assert.equal(idle.height, 16);
    assert.ok(idle._ctx.pixels.size > 30, `building ${kind} should have substantial silhouette`);
    const coords = [...idle._ctx.pixels.keys()].sort().join(";");
    signatures.add(coords);
    // Active pulse frame alters at least one pixel on producing buildings.
    const diff = [...active._ctx.pixels.entries()].some(([k, v]) => idle._ctx.pixels.get(k) !== v);
    assert.ok(diff, `building ${kind} active frame should pulse visibly`);
  }
  assert.equal(signatures.size, BUILDING_KINDS.length, "each building kind should have a distinct silhouette");

  // Progress bar appears iff progress > 0; wear warning mark appears iff wear > 80.
  const makeDrawCtx = () => {
    const ops = [];
    return {
      ops,
      fillStyle: "",
      strokeStyle: "",
      lineWidth: 1,
      imageSmoothingEnabled: true,
      save() {},
      restore() {},
      beginPath() { ops.push(["beginPath"]); },
      closePath() { ops.push(["closePath"]); },
      moveTo(x, y) { ops.push(["moveTo", x, y]); },
      lineTo(x, y) { ops.push(["lineTo", x, y]); },
      ellipse(...a) { ops.push(["ellipse", ...a]); },
      fill() { ops.push(["fill", this.fillStyle]); },
      stroke() { ops.push(["stroke", this.strokeStyle]); },
      drawImage(img, x, y, w, h) { ops.push(["drawImage", img, x, y, w, h]); },
      fillRect(x, y, w, h) { ops.push(["fillRect", x, y, w, h, this.fillStyle]); },
    };
  };

  const idleCtx = makeDrawCtx();
  Buildings.draw(idleCtx, 100, 100, 24, "composter", { progress: 0, wear: 80 }, 0);
  assert.equal(idleCtx.ops.some((op) => op[0] === "fillRect" && op[5] === "#d39b55"), false);
  assert.equal(idleCtx.ops.some((op) => op[0] === "fill" && (op[1] === "#f0a14a" || op[1] === "#ff6b6b")), false);

  const activeCtx = makeDrawCtx();
  Buildings.draw(activeCtx, 100, 100, 24, "composter", { progress: 0.4, wear: 81 }, 500);
  assert.equal(activeCtx.ops.some((op) => op[0] === "fillRect" && op[5] === "#d39b55"), true);
  assert.equal(activeCtx.ops.some((op) => op[0] === "fill" && op[1] === "#f0a14a"), true);

  const brokenCtx = makeDrawCtx();
  Buildings.draw(brokenCtx, 100, 100, 24, "smelter", { progress: 0.9, wear: 100 }, 500);
  assert.equal(brokenCtx.ops.some((op) => op[0] === "fill" && op[1] === "#ff6b6b"), true);
});

test("Colony live flow updates buildings list, world tiles, legend and redraws without errors", () => {
  const env = loadLiveColonyEnv("en");
  const { context, elements, logs, sounds } = env;
  const { Colony, Buildings, apply } = context;

  const makeGrid = (n) => Array.from({ length: n }, (_, x) =>
    Array.from({ length: n }, (_, y) => ({ x, y, z: 0, g: "regolith" })));

  context.W = {
    w: 4,
    h: 4,
    depth: 1,
    view: 0,
    follow: true,
    colony: true,
    live: true,
    t: 10,
    inv: {},
    layers: [makeGrid(4)],
    drones: new Map(),
  };

  Colony.setData({
    state: { tick: 10, size: 4, depth: 1, credits: 100, levels: {}, inv: {}, tiles: [], buildings: [] },
    prices: {},
    contracts: [],
    shop: [],
    levels: { total: 10 },
  });

  // 1) Initial empty state keeps #col-buildings hidden.
  Colony.renderBuildings();
  assert.equal(elements.get("#col-buildings").classList.contains("hidden"), true);
  assert.equal(elements.get("#col-buildings").innerHTML, "");

  // 2) Live building event updates status map, W tile.b, #col-buildings, and redraws without fail sound/log.
  const drawsBefore = env.drawCalls;
  const legendBefore = env.legendCalls;
  Colony.onEvent({
    op: "building",
    t: 123,
    x: 2,
    y: 1,
    building: { kind: "composter", recipe: "biomass", in: { moss: 4 }, out: { biomass: 1 }, progress: 0.4, wear: 3, level: 1 },
  });
  assert.equal(context.W.layers[0][2][1].b, "composter");
  assert.equal(Buildings.at(2, 1).recipe, "biomass");
  assert.equal(elements.get("#col-buildings").classList.contains("hidden"), false);
  assert.match(elements.get("#col-buildings").innerHTML, /Composter/);
  assert.match(elements.get("#col-buildings").innerHTML, /\(2, 1\)/);
  assert.equal(env.drawCalls, drawsBefore + 1);
  assert.equal(env.legendCalls, legendBefore + 1);
  assert.equal(sounds.includes("fail"), false);
  assert.equal(logs.length, 0);

  // 3) Full state event syncs tile b + buildings array and updates panel without reloading.
  const nextTiles = Array.from({ length: 16 }, () => ({ g: "regolith" }));
  nextTiles[1 * 4 + 3] = { g: "regolith", b: "kiln" };
  Colony.onEvent({
    op: "state",
    t: 130,
    state: {
      tick: 130,
      size: 4,
      depth: 1,
      credits: 200,
      levels: {},
      inv: { glass: 2 },
      tiles: nextTiles,
      buildings: [
        { kind: "kiln", x: 3, y: 1, level: 1, recipe: "glass", in: { silica: 2 }, out: { glass: 1 }, progress: 0.5, wear: 85 },
      ],
    },
  });
  assert.equal(context.W.layers[0][2][1].b, undefined);
  assert.equal(context.W.layers[0][3][1].b, "kiln");
  assert.equal(Buildings.at(2, 1), null);
  assert.equal(Buildings.at(3, 1).wear, 85);
  assert.match(elements.get("#col-buildings").innerHTML, /Kiln/);
  assert.match(elements.get("#col-buildings").innerHTML, /bld-worn/);
  assert.match(elements.get("#col-market").innerHTML, /Glass/);

  // 4) Tile event setting and clearing tile.b works through apply without false failure logs.
  apply({ op: "tile", t: 135, tile: { x: 0, y: 2, z: 0, g: "regolith", b: "silo" } });
  assert.equal(context.W.layers[0][0][2].b, "silo");
  apply({ op: "tile", t: 136, tile: { x: 3, y: 1, z: 0, g: "regolith" } });
  assert.equal(Buildings.at(3, 1), null);
  assert.equal(sounds.includes("fail"), false);
  assert.equal(logs.length, 0);

  // 5) Empty buildings array hides #col-buildings again.
  Colony.onEvent({
    op: "state",
    t: 140,
    state: { tick: 140, size: 4, depth: 1, credits: 200, levels: {}, inv: {}, tiles: Array.from({ length: 16 }, () => ({ g: "regolith" })), buildings: [] },
  });
  assert.equal(elements.get("#col-buildings").classList.contains("hidden"), true);
  assert.equal(elements.get("#col-buildings").innerHTML, "");
});
