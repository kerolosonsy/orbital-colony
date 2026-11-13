"use strict";
// Building sprites, live status map, and colony panel markup.
const Buildings = (() => {
  const S = 16;
  const KINDS = ["composter", "smelter", "kiln", "refinery", "assembler", "reactor-lab", "silo", "generator", "lab"];
  const cache = new Map();
  const statusMap = new Map();

  const key = (x, y) => `${x},${y}`;

  function clear() { statusMap.clear(); }

  function normalize(b, fx, fy) {
    return {
      kind: b.kind,
      x: Number(b.x !== undefined ? b.x : fx) || 0,
      y: Number(b.y !== undefined ? b.y : fy) || 0,
      level: b.level ?? 1,
      recipe: b.recipe || "",
      in: b.in || {},
      out: b.out || {},
      progress: Number(b.progress) || 0,
      wear: Number(b.wear) || 0,
    };
  }

  function syncState(buildings) {
    statusMap.clear();
    for (const b of buildings || []) {
      if (!b || !b.kind) continue;
      const entry = normalize(b);
      statusMap.set(key(entry.x, entry.y), entry);
    }
  }

  function applyEvent(ev, world) {
    const b = ev && ev.building;
    const x = b && b.x !== undefined ? b.x : ev && ev.x;
    const y = b && b.y !== undefined ? b.y : ev && ev.y;
    if (x === undefined || y === undefined) return null;
    const tile = world && world.layers && world.layers[ev.z || 0] && world.layers[ev.z || 0][x] && world.layers[ev.z || 0][x][y];
    if (!b || !b.kind) {
      statusMap.delete(key(x, y));
      if (tile) delete tile.b;
      return null;
    }
    const entry = normalize(b, x, y);
    statusMap.set(key(entry.x, entry.y), entry);
    if (tile) tile.b = entry.kind;
    return entry;
  }

  function remove(x, y) { statusMap.delete(key(x, y)); }
  function at(x, y) { return statusMap.get(key(x, y)) || null; }
  function list() {
    return [...statusMap.values()].sort((a, b) => a.y - b.y || a.x - b.x || (a.kind < b.kind ? -1 : a.kind > b.kind ? 1 : 0));
  }

  // ---------- procedural 16x16 building sprites ----------
  function makeSprite(k, paint) {
    let c = cache.get(k);
    if (c) return c;
    c = document.createElement("canvas");
    c.width = S; c.height = S;
    const g = c.getContext("2d");
    const px = (x, y, col) => {
      if (x < 0 || y < 0 || x >= S || y >= S || !col) return;
      g.fillStyle = col; g.fillRect(x | 0, y | 0, 1, 1);
    };
    paint(px, g);
    cache.set(k, c);
    return c;
  }

  function paintGrid(px, rows, pal) {
    for (let y = 0; y < rows.length; y++) {
      const row = rows[y];
      for (let x = 0; x < row.length; x++) {
        const col = pal[row[x]];
        if (col) px(x, y, col);
      }
    }
  }

  const ART = {
    composter: {
      pal: { "#": "#0d0a1c", o: "#b88a58", b: "#7a5330", d: "#4f3219", g: "#58993e", G: "#93c45c" },
      rows: [
        "................",
        "................",
        "................",
        "................",
        "..############..",
        ".#ggggGGGGgggg#.",
        ".#GGGGdGGdGGGG#.",
        ".#oooooooooooo#.",
        ".#bbbbbbbbbbbb#.",
        ".##b#bbbbbb#b##.",
        ".#dddddddddddd#.",
        ".#b#bbbbbbbb#b#.",
        ".#dddddddddddd#.",
        "..############..",
        "................",
        "................",
      ],
      pulse: (px) => { px(5, 4, "#f6f39a"); px(10, 4, "#b2e57a"); px(7, 5, "#f6f39a"); },
    },
    smelter: {
      pal: { "#": "#0d0a1c", S: "#8da4be", s: "#5a6b82", b: "#4b455e", d: "#312c40", f: "#d77a33", F: "#ef5a40", Y: "#ffcf5a" },
      rows: [
        "..####..........",
        "..#ss#..........",
        "..#Ss#..........",
        "..#Ss#..........",
        "..#Ss#######....",
        "..#Ssbbbbbbb#...",
        ".##Ssbbbbbbbb#..",
        ".#bbbbbbbbbbb#..",
        ".#bb######bbb#..",
        ".#bb#ffff#bbb#..",
        ".#bb#fFFf#bbb#..",
        ".#dd#fYYf#ddd#..",
        ".#dd#ffff#ddd#..",
        ".#ddddddddddd#..",
        ".#############..",
        "................",
      ],
      pulse: (px) => { px(3, 0, "#a69dc2"); px(4, 0, "#ffcf5a"); px(7, 11, "#fff2b8"); px(8, 11, "#fff2b8"); },
    },
    kiln: {
      pal: { "#": "#0d0a1c", c: "#77709a", T: "#d39b55", t: "#a45f3e", d: "#6e3b23", g: "#5ecdf5", G: "#b0f2f5", W: "#ffffff" },
      rows: [
        "................",
        "................",
        ".....######.....",
        "....#cccccc#....",
        "...#TTTTTTTT#...",
        "..#TTTttttttt#..",
        "..#TTtttttttt#..",
        ".#TTtt#####ttt#.",
        ".#Ttt#ggggg#tt#.",
        ".#ttt#gGGGg#dd#.",
        ".#ddd#gWWWg#dd#.",
        ".#ddd#ggggg#dd#.",
        ".#dddddddddddd#.",
        "..############..",
        "................",
        "................",
      ],
      pulse: (px) => { px(7, 1, "#d9cbb0"); px(8, 1, "#efe7d6"); px(7, 9, "#ffffff"); px(8, 9, "#ffffff"); },
    },
    refinery: {
      pal: { "#": "#0d0a1c", l: "#b8c1cf", M: "#7b96b5", d: "#4b5973", a: "#e04665" },
      rows: [
        "...###..........",
        "..#lll#.........",
        "..#lMl#.........",
        "..#lMl######....",
        "..#ddd#....#....",
        "..#lMl#..#####..",
        "..#aaa#.#lllll#.",
        "..#lMl#.#lMMMl#.",
        "..#ddd#.#MMMMd#.",
        "..#lMl#.#MddMd#.",
        "..#lMl#.#dddd##.",
        "..#ddd#..##dd#..",
        "..#ddd#...#dd#..",
        "..############..",
        "................",
        "................",
      ],
      pulse: (px) => { px(4, 0, "#ffcf5a"); px(4, 6, "#5ecdf5"); px(11, 7, "#5ecdf5"); px(11, 8, "#b0f2f5"); },
    },
    assembler: {
      pal: { "#": "#0d0a1c", L: "#c6d0de", l: "#7a89a6", d: "#465066", y: "#d39b55", Y: "#ffcf5a", c: "#2ab572" },
      rows: [
        "................",
        ".##############.",
        ".#LLLLLLLLLLLL#.",
        ".#ll##yyyy##ll#.",
        ".#ll.#yYy#.#ll#.",
        ".#ll..#yy#.#ll#.",
        ".#ll...#y#.#ll#.",
        ".#ll.......#ll#.",
        ".#ll.#####.#ll#.",
        ".#dd.#ccc#.#dd#.",
        ".#dd.#ddd#.#dd#.",
        ".#dd.......#dd#.",
        ".#dd.......#dd#.",
        ".####.....#####.",
        "................",
        "................",
      ],
      pulse: (px) => { px(8, 7, "#8ff7ff"); px(7, 8, "#ffcf5a"); px(7, 9, "#baf54e"); },
    },
    "reactor-lab": {
      pal: { "#": "#0d0a1c", l: "#969689", M: "#6e7ba3", d: "#3f465e", u: "#1e8745", U: "#40de7a", G: "#a6ffc4" },
      rows: [
        "......####......",
        ".....#llll#.....",
        ".##..#lMMl#..##.",
        ".#l#.#uuuu#.#d#.",
        ".#l#.#uUUu#.#d#.",
        ".#l###uGGu###d#.",
        ".#ll##uUUu##dd#.",
        ".#dd##uuuu##dd#.",
        ".#dd##lMMl##dd#.",
        ".#dd.#dddd#.#d#.",
        ".#dd.#dddd#.#d#.",
        ".#dd########dd#.",
        ".#dddddddddddd#.",
        ".##############.",
        "................",
        "................",
      ],
      pulse: (px) => { px(7, 4, "#a6ffc4"); px(8, 4, "#ffffff"); px(7, 5, "#ffffff"); px(8, 5, "#a6ffc4"); },
    },
    silo: {
      pal: { "#": "#0d0a1c", L: "#efe7d6", l: "#b8c1cf", d: "#68738a", g: "#3a95e8", G: "#5ecdf5" },
      rows: [
        "....########....",
        "...#LLLLLLLL#...",
        "..#LLLlllllll#..",
        "..#llllllllll#..",
        "..#dddddddddd#..",
        "..#llggggllll#..",
        "..#llgGGgllld#..",
        "..#llgggglldd#..",
        "..#dddddddddd#..",
        "..#llllllllld#..",
        "..#lllllllldd#..",
        "..#ddddddddd##..",
        "..#dddddddddd#..",
        "..############..",
        "................",
        "................",
      ],
      pulse: (px) => { px(6, 5, "#baf54e"); px(7, 5, "#baf54e"); px(6, 6, "#ffffff"); },
    },
    generator: {
      pal: { "#": "#0d0a1c", l: "#8db3d9", d: "#465a73", c: "#d45842", C: "#f2a07e", y: "#ffcf5a", r: "#4e80f0", R: "#b0f2f5" },
      rows: [
        "...##......##...",
        "...#y#....#y#...",
        "...#y#....#y#...",
        "..############..",
        ".#llllllllllll#.",
        ".#llcccccc##ll#.",
        ".#llcCCCCc#rRl#.",
        ".#llcccccc#Rrr#.",
        ".#llcCCCCc#rRl#.",
        ".#ddcccccc##dd#.",
        ".#dddddddddddd#.",
        ".#dddddddddddd#.",
        "..############..",
        "................",
        "................",
        "................",
      ],
      pulse: (px) => {
        px(6, 1, "#fff2b8"); px(8, 1, "#fff2b8");
        px(11, 6, "#4e80f0"); px(12, 6, "#b0f2f5"); px(11, 7, "#b0f2f5"); px(12, 7, "#4e80f0");
      },
    },
    lab: {
      pal: { "#": "#0d0a1c", W: "#efe7d6", w: "#cfc6b2", d: "#77709a", s: "#463f66", i: "#45b3a8", I: "#9fd2cb", p: "#8db3d9", P: "#d4e8fa" },
      rows: [
        ".........###....",
        "........#pPp#...",
        ".........#P#....",
        "......####p#....",
        "....##WWWWWW##..",
        "...#WWWwwwwwww#.",
        "..#WWwwwwwwwww#.",
        "..#ww#######ww#.",
        "..#ww#iiiii#dd#.",
        "..#ww#iIIIi#dd#.",
        "..#ddddddddddd#.",
        "..#sssssssssss#.",
        "..#############.",
        "................",
        "................",
        "................",
      ],
      pulse: (px) => { px(10, 1, "#8f74ff"); px(11, 0, "#d6ccff"); px(7, 9, "#ffffff"); },
    },
  };

  function sprite(kind, frame = 0) {
    const f = frame ? 1 : 0;
    return makeSprite(`b:${kind}:${f}`, (px) => {
      const def = ART[kind];
      if (!def) {
        paintGrid(px, ART.silo.rows, { ...ART.silo.pal, L: "#a69dc2", l: "#77709a" });
        return;
      }
      paintGrid(px, def.rows, def.pal);
      if (f && def.pulse) def.pulse(px);
    });
  }

  // Draw the building sprite on its tile top face, plus a progress bar when
  // producing and a warning mark when wear > 80.
  function draw(ctx, X, Y, hw, kind, status, now) {
    if (!kind) return;
    const size = Math.round((hw * 1.3) / 8) * 8 || hw;
    const producing = !!(status && status.progress > 0);
    const frame = producing ? Math.floor((now || 0) / 450) % 2 : 0;
    const img = sprite(kind, frame);

    ctx.save();
    ctx.fillStyle = "rgba(10,6,20,0.35)";
    ctx.beginPath();
    ctx.ellipse(X, Y, hw * 0.48, hw * 0.24, 0, 0, Math.PI * 2);
    ctx.fill();
    ctx.imageSmoothingEnabled = false;
    const topY = Math.round(Y + size * 0.22 - size);
    ctx.drawImage(img, Math.round(X - size / 2), topY, size, size);

    if (producing) {
      const p = Math.max(0, Math.min(1, status.progress));
      const barW = Math.max(14, Math.round(hw * 1.1));
      const barH = Math.max(3, Math.round(hw * 0.14));
      const bx = Math.round(X - barW / 2);
      const by = Math.round(topY - barH - 3);
      ctx.fillStyle = "#0d0a1c";
      ctx.fillRect(bx - 1, by - 1, barW + 2, barH + 2);
      ctx.fillStyle = "#17142a";
      ctx.fillRect(bx, by, barW, barH);
      ctx.fillStyle = "#d39b55";
      ctx.fillRect(bx, by, Math.max(1, Math.round(barW * p)), barH);
    }

    if (status && status.wear > 80) {
      const r = Math.max(4, Math.round(hw * 0.24));
      const wx = Math.round(X + size * 0.36);
      const wy = Math.round(topY + r + 1);
      ctx.beginPath();
      ctx.moveTo(wx, wy - r);
      ctx.lineTo(wx + r, wy + r);
      ctx.lineTo(wx - r, wy + r);
      ctx.closePath();
      ctx.fillStyle = status.wear >= 100 ? "#ff6b6b" : "#f0a14a";
      ctx.fill();
      ctx.lineWidth = 1.5;
      ctx.strokeStyle = "#0d0a1c";
      ctx.stroke();
      ctx.fillStyle = "#17142a";
      ctx.fillRect(wx - 1, wy - Math.max(1, Math.round(r * 0.35)), 2, Math.max(2, Math.round(r * 0.6)));
      ctx.fillRect(wx - 1, wy + Math.max(1, Math.round(r * 0.5)), 2, 1.5);
    }
    ctx.restore();
  }

  // ---------- colony panel markup ----------
  const escAttr = (s, escapeFn) => escapeFn(String(s)).replace(/"/g, "&quot;").replace(/'/g, "&#39;");
  const fmtAmt = (n) => (Number.isInteger(n) ? n.toLocaleString("en-US") : Number(n.toFixed(1)).toLocaleString("en-US"));

  function chipsMarkup(map, names, colors, escapeFn) {
    const entries = Object.entries(map || {}).filter(([, n]) => n > 0).sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0));
    if (!entries.length) return '<span class="bld-none">—</span>';
    return entries.map(([g, n]) => {
      const label = escAttr(names[g] || g, escapeFn);
      const col = escAttr(colors[g] || "#ccc", escapeFn);
      return `<span class="cg bld-chip" title="${label}"><i style="background:${col}"></i>${fmtAmt(n)}</span>`;
    }).join("");
  }

  function listMarkup(buildings, deps = {}) {
    if (!buildings || !buildings.length) return "";
    const tr = deps.T || T;
    const escapeFn = deps.esc || esc;
    const names = deps.NAME || NAME;
    const colors = deps.INV_COLOR || INV_COLOR;
    const rows = buildings.map((b) => {
      const name = escapeFn(String(names[b.kind] || b.kind));
      const recipeName = b.recipe ? escapeFn(String(names[b.recipe] || b.recipe)) : escapeFn(tr("بلا وصفة", "No recipe"));
      const recipeChip = b.recipe && colors[b.recipe] ? `<i style="background:${escAttr(colors[b.recipe], escapeFn)}"></i>` : "";
      const wear = Math.round(Number(b.wear) || 0);
      const worn = wear > 80;
      const x = Number(b.x) || 0, y = Number(b.y) || 0;
      return `<div class="item bld-row${worn ? " bld-worn" : ""}" data-x="${x}" data-y="${y}">`
        + `<div class="item-top"><b>${name}</b><span class="bld-pos" dir="ltr">(${x}, ${y})</span></div>`
        + `<div class="bld-meta"><span class="bld-recipe">${recipeChip}${recipeName}</span><span class="bld-wear${worn ? " warn" : ""}">${escapeFn(tr("تآكل", "Wear"))} <b dir="ltr">${wear}%</b></span></div>`
        + `<div class="bld-io"><span class="bld-dir"><span class="bld-lbl">${escapeFn(tr("داخل", "In"))}</span>${chipsMarkup(b.in, names, colors, escapeFn)}</span><span class="bld-dir"><span class="bld-lbl">${escapeFn(tr("خارج", "Out"))}</span>${chipsMarkup(b.out, names, colors, escapeFn)}</span></div>`
        + `</div>`;
    }).join("");
    return `<h3>${escapeFn(tr("المباني", "Buildings"))} <small>${buildings.length}</small></h3>` + rows;
  }

  return { KINDS, clear, syncState, applyEvent, remove, at, list, sprite, draw, listMarkup };
})();
