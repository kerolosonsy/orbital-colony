"use strict";
// Pixel-art sprites for Kepler-G, painted procedurally at 16x16 and cached.
// Every material gets its own palette and texture; nothing here is borrowed
// from another game's art.
const Sprites = (() => {
  const S = 16;
  const cache = new Map();

  function rng(seed) {
    let a = seed >>> 0;
    return () => {
      a = (a + 0x6d2b79f5) | 0;
      let t = Math.imul(a ^ (a >>> 15), 1 | a);
      t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
      return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
    };
  }

  function sprite(key, w, h, paint) {
    let c = cache.get(key);
    if (c) return c;
    c = document.createElement("canvas");
    c.width = w; c.height = h;
    const g = c.getContext("2d");
    const px = (x, y, col) => {
      if (x < 0 || y < 0 || x >= w || y >= h || !col) return;
      g.fillStyle = col; g.fillRect(x | 0, y | 0, 1, 1);
    };
    paint(px, g);
    cache.set(key, c);
    return c;
  }

  // smooth value noise on a coarse grid, for natural colour patches
  function noiseField(r, cells = 4) {
    const grid = Array.from({ length: cells + 1 }, () => Array.from({ length: cells + 1 }, r));
    return (x, y) => {
      const fx = (x / S) * cells, fy = (y / S) * cells;
      const ix = Math.floor(fx), iy = Math.floor(fy), tx = fx - ix, ty = fy - iy;
      const a = grid[iy][ix], b = grid[iy][ix + 1], c = grid[iy + 1][ix], d = grid[iy + 1][ix + 1];
      return a + (b - a) * tx + (c - a) * ty + (a - b - c + d) * tx * ty;
    };
  }

  // raised-tile bevel: lit top-left edge, shaded bottom-right edge
  function bevel(px, pal, read) {
    for (let i = 0; i < S; i++) {
      px(i, 0, pal[Math.min(pal.length - 1, read(i, 0) + 1)]);
      px(0, i, pal[Math.min(pal.length - 1, read(0, i) + 1)]);
      px(i, S - 1, pal[Math.max(0, read(i, S - 1) - 1)]);
      px(S - 1, i, pal[Math.max(0, read(S - 1, i) - 1)]);
    }
  }

  function grid() { return Array.from({ length: S }, () => new Array(S).fill(0)); }
  function flush(px, pal, m) { for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) px(x, y, pal[m[y][x]]); }

  // ---------- ground ----------
  const PAL = {
    regolith: ["#8a5426", "#a8692f", "#bd7c3c", "#cf914c", "#dfa862", "#f0c88c"],
    soil: ["#24140b", "#33200f", "#452b18", "#5a3a22", "#6e4b2e", "#87613f"],
    ice: ["#3f7fa6", "#5b9fc4", "#7dbedb", "#a3d8ec", "#c7ecf7", "#f2fcff"],
    rock: ["#25212f", "#363145", "#4b455e", "#625b79", "#7d7697", "#a39dbc"],
    stone: ["#2b2733", "#3d3847", "#524b5e", "#675f74", "#7d7489", "#958ca1"],
    tunnel: ["#2a1d15", "#3b2a1e", "#4f3826", "#634630", "#7a573b", "#94704f"],
  };

  // underground stone: layered strata with a few fractures
  function stone(v) {
    return sprite(`g:stone:${v}`, S, S, (px) => {
      const r = rng(5000 + v), m = grid(), pal = PAL.stone;
      for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) {
        const band = Math.floor((y + Math.sin((x + v) * 0.7) * 1.2) / 4) % 2;
        m[y][x] = 2 + band + (r() < 0.2 ? (r() < 0.5 ? -1 : 1) : 0);
      }
      for (let y = 3; y < S; y += 4) for (let x = 0; x < S; x++) if (r() < 0.7) m[y][x] = 1;
      for (let c = 0; c < 2; c++) { let x = Math.floor(r() * S), y = Math.floor(r() * S); for (let k = 0; k < 5; k++) { if (m[y]) m[y][x] = 0; x += r() < 0.5 ? 1 : 0; y += 1; if (y >= S) break; } }
      flush(px, pal, m.map((row) => row.map((k) => Math.max(0, Math.min(5, k)))));
      bevel(px, pal, (x, y) => Math.max(0, Math.min(5, m[y][x])));
    });
  }

  // dug-out tunnel floor: dark dust and gravel
  function tunnel(v) {
    return sprite(`g:tunnel:${v}`, S, S, (px) => {
      const r = rng(6000 + v), n = noiseField(r, 3), m = grid(), pal = PAL.tunnel;
      for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) m[y][x] = n(x, y) < 0.5 ? 2 : 3;
      for (let p = 0; p < 7; p++) { const x = Math.floor(r() * 15), y = Math.floor(r() * 15); m[y][x] = 5; m[y][x + 1] = 1; }
      flush(px, pal, m);
    });
  }

  // ore seen in the face of the rock (drawn over stone)
  function ore(kind, v) {
    return sprite(`o:${kind}:${v}`, S, S, (px) => {
      const r = rng(7000 + v + (kind === "geode" ? 50 : 0));
      const pal = kind === "geode" ? ["#2a1a55", "#5b3cc4", "#8f74ff", "#d6ccff"] : ["#4a2410", "#9a4d1e", "#d77a33", "#ffc07a"];
      const n = kind === "geode" ? 2 : 4;
      for (let i = 0; i < n; i++) {
        const cx = 3 + Math.floor(r() * 10), cy = 3 + Math.floor(r() * 10), rad = kind === "geode" ? 3 : 1.6 + r();
        for (let y = -4; y <= 4; y++) for (let x = -4; x <= 4; x++) {
          const d = Math.hypot(x, y);
          if (d > rad) continue;
          px(cx + x, cy + y, d > rad - 0.8 ? pal[0] : (x + y < -1 ? pal[3] : x + y < 1 ? pal[2] : pal[1]));
        }
        if (kind === "geode") { px(cx, cy, "#ffffff"); px(cx - 1, cy + 1, "#d6ccff"); }
      }
    });
  }

  function regolith(v) {
    return sprite(`g:regolith:${v}`, S, S, (px) => {
      const r = rng(1000 + v), n = noiseField(r), m = grid(), pal = PAL.regolith;
      for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) {
        const k = n(x, y) * 2.2 + r() * 0.9;
        m[y][x] = k < 0.9 ? 1 : k < 1.6 ? 2 : k < 2.3 ? 3 : 4;
      }
      if (v % 4 === 0) { // a small crater
        const cx = 5 + Math.floor(r() * 6), cy = 5 + Math.floor(r() * 6);
        for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) {
          const d = Math.hypot(x - cx, y - cy);
          if (d < 2.2) m[y][x] = 1;
          else if (d < 3.2) m[y][x] = (x - cx) + (y - cy) > 0 ? 5 : 0;
        }
      }
      for (let p = 0; p < 3 + (v % 3); p++) { // pebbles with a lit top-left pixel
        const x = 1 + Math.floor(r() * 13), y = 1 + Math.floor(r() * 13);
        m[y][x] = 0; m[y][x + 1] = 0; if (r() < 0.5) m[y + 1][x] = 1;
        if (y > 0) m[y - 1][x] = 5;
      }
      flush(px, pal, m);
      bevel(px, pal, (x, y) => m[y][x]);
    });
  }

  function soil(v) {
    return sprite(`g:soil:${v}`, S, S, (px) => {
      const r = rng(2000 + v), m = grid(), pal = PAL.soil;
      for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) {
        const row = (y + (v % 2) * 2) % 4; // furrow, its shadow, the ridge, its lit crest
        let k = [1, 2, 3, 4][row];
        if (r() < 0.18) k += r() < 0.5 ? -1 : 1;
        m[y][x] = Math.max(0, Math.min(5, k));
      }
      for (let p = 0; p < 6; p++) { // clods and damp specks
        const x = Math.floor(r() * 15), y = Math.floor(r() * 15);
        if (r() < 0.5) { m[y][x] = 5; m[y][x + 1] = 4; } else m[y][x] = 0;
      }
      flush(px, pal, m);
      bevel(px, pal, (x, y) => m[y][x]);
    });
  }

  function ice(v, frame) {
    return sprite(`g:ice:${v}:${frame}`, S, S, (px) => {
      const r = rng(3000 + v), n = noiseField(r, 3), m = grid(), pal = PAL.ice;
      for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) {
        let k = n(x, y) < 0.45 ? 2 : 3;
        const band = (x + y + v * 3) % 12;
        if (band === 0 || band === 1) k = 4;
        if (band === 2) k = 5;
        m[y][x] = k;
      }
      for (let c = 0; c < 2; c++) { // hairline cracks with a bright edge
        let x = Math.floor(r() * S), y = Math.floor(r() * S);
        for (let s = 0; s < 6 + Math.floor(r() * 5); s++) {
          if (x < 0 || y < 0 || x >= S || y >= S) break;
          m[y][x] = 0;
          if (y + 1 < S && m[y + 1][x] !== 0) m[y + 1][x] = 5;
          x += r() < 0.5 ? 1 : -1; y += r() < 0.6 ? 1 : 0;
        }
      }
      flush(px, pal, m);
      bevel(px, pal, (x, y) => m[y][x]);
      if (frame) { // a glint that comes and goes
        const gx = 3 + (v * 5) % 10, gy = 3 + (v * 7) % 10;
        px(gx, gy, "#ffffff"); px(gx - 1, gy, "#e6fbff"); px(gx + 1, gy, "#e6fbff"); px(gx, gy - 1, "#e6fbff"); px(gx, gy + 1, "#e6fbff");
      }
    });
  }

  function rock(v) {
    return sprite(`g:rock:${v}`, S, S, (px) => {
      const r = rng(4000 + v), m = grid(), pal = PAL.rock;
      for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) m[y][x] = r() < 0.2 ? 0 : 1;
      const stones = [[4 + r() * 2, 5 + r() * 2, 3.6 + r()], [11 + r() * 2, 6 + r() * 3, 3 + r()], [7 + r() * 3, 12 + r(), 2.6 + r()]];
      for (const [cx, cy, rad] of stones) {
        for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) {
          const dx = x - cx, dy = (y - cy) * 1.15, d = Math.hypot(dx, dy);
          if (d > rad) continue;
          if (d > rad - 1) { m[y][x] = 0; continue; }
          const lit = (-dx - dy) / rad;
          m[y][x] = lit > 0.55 ? 5 : lit > 0.15 ? 4 : lit > -0.35 ? 3 : 2;
        }
      }
      const [cx, cy] = stones[0]; // a crack across the biggest stone
      for (let i = -1; i < 2; i++) { const x = Math.round(cx + i), y = Math.round(cy + i * 0.5); if (m[y]?.[x] > 1) m[y][x] = 1; }
      flush(px, pal, m);
      bevel(px, pal, (x, y) => m[y][x]);
    });
  }

  function ground(kind, v, frame = 0) {
    switch (kind) {
      case "stone": return stone(v % 6);
      case "tunnel": case "shaft": return tunnel(v % 6);
      case "soil": return soil(v % 6);
      case "ice": return ice(v % 6, frame);
      case "rock": return rock(v % 6);
      default: return regolith(v % 8);
    }
  }

  // ---------- life and minerals ----------
  function blob(px, cx, cy, rad, pal, r) {
    // a lumpy dome lit from the top-left; pal = [outline, dark, mid, light, highlight]
    const lumps = Array.from({ length: 8 }, () => 0.82 + r() * 0.3);
    for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) {
      const dx = x + 0.5 - cx, dy = y + 0.5 - cy, a = Math.atan2(dy, dx);
      const rr = rad * lumps[Math.floor(((a + Math.PI) / (Math.PI * 2)) * 8) % 8];
      const d = Math.hypot(dx, dy);
      if (d > rr) continue;
      if (d > rr - 1) { px(x, y, pal[0]); continue; }
      const lit = (-dx - dy * 1.2) / rr;
      px(x, y, lit > 0.5 ? pal[4] : lit > 0.1 ? pal[3] : lit > -0.4 ? pal[2] : pal[1]);
    }
  }

  function moss(stage, ripe) {
    // moss is a cushion of many tiny tufts, not one ball
    return sprite(`e:moss:${stage}:${ripe}`, S, S, (px) => {
      const r = rng(50 + stage), pal = ["#173a17", "#2d5e24", "#4a8f32", "#74bd47", "#b2e57a"];
      if (stage === 0) {
        for (const [x, h] of [[6, 3], [9, 2], [11, 2]]) {
          for (let i = 0; i < h; i++) px(x, 13 - i, pal[2]);
          px(x - 1, 13 - h, pal[3]); px(x + 1, 13 - h + 1, pal[3]);
        }
        return;
      }
      const spread = [0, 2.2, 3.4, 4.6, 5.6][stage];
      const count = [0, 4, 7, 10, 13][stage];
      const tufts = [];
      for (let i = 0; i < count; i++) {
        const a = r() * Math.PI * 2, d = Math.sqrt(r()) * spread;
        tufts.push([8 + Math.cos(a) * d, 9.5 + Math.sin(a) * d * 0.8, 1.5 + r() * 1.1]);
      }
      tufts.sort((a, b) => a[1] - b[1]); // back rows first
      for (const [cx, cy, rad] of tufts) {
        for (let y = Math.floor(cy - rad - 1); y <= cy + rad + 1; y++) for (let x = Math.floor(cx - rad - 1); x <= cx + rad + 1; x++) {
          const dx = x + 0.5 - cx, dy = y + 0.5 - cy, d = Math.hypot(dx, dy);
          if (d > rad + 0.4) continue;
          if (d > rad - 0.55) { px(x, y, pal[0]); continue; }
          const lit = (-dx - dy) / rad;
          px(x, y, lit > 0.45 ? pal[4] : lit > -0.1 ? pal[3] : pal[2]);
        }
      }
      if (ripe) for (let i = 0; i < 5; i++) {
        const [cx, cy] = tufts[Math.floor(r() * tufts.length)];
        px(Math.round(cx), Math.round(cy - 2), "#f6f39a"); px(Math.round(cx), Math.round(cy - 3), "#d9cf5a");
      }
    });
  }

  function lichen(stage, ripe) {
    return sprite(`e:lichen:${stage}:${ripe}`, S, S, (px) => {
      const r = rng(70 + stage), pal = ["#0f3f3c", "#1f6f6a", "#2f9a90", "#57c2b4", "#9be8d9"];
      const rad = [1.6, 2.8, 4, 5.2, 6.4][stage];
      const scallop = Array.from({ length: 10 }, () => 0.8 + r() * 0.3);
      for (let y = 0; y < S; y++) for (let x = 0; x < S; x++) {
        const dx = x + 0.5 - 8, dy = y + 0.5 - 8, d = Math.hypot(dx, dy), a = Math.atan2(dy, dx);
        const rr = rad * scallop[Math.floor(((a + Math.PI) / (Math.PI * 2)) * 10) % 10];
        if (d > rr) continue;
        if (d > rr - 1) { px(x, y, pal[0]); continue; }
        px(x, y, pal[1 + (Math.floor(d * 1.3) % 3)]);
      }
      px(8, 8, ripe ? "#e8fff9" : pal[4]);
      if (ripe) for (const [x, y] of [[5, 6], [11, 9], [7, 12]]) px(x, y, "#e8fff9");
    });
  }

  function glowcap(stage, ripe) {
    return sprite(`e:glowcap:${stage}:${ripe}`, S, S, (px) => {
      const spots = [[8, 10, 3], [4, 12, 2], [12, 11, 2], [6, 7, 2], [11, 6, 2]];
      const n = Math.min(spots.length, 1 + stage);
      const cap = ripe ? ["#3a1d73", "#6a3ccf", "#9b6bff", "#c7a8ff"] : ["#2b1a4d", "#4b2e8a", "#6a4bb3", "#8e74cf"];
      for (let i = 0; i < n; i++) {
        const [cx, base, w] = spots[i];
        const h = 1 + Math.min(3, stage);
        for (let s = 0; s < h; s++) { px(cx, base - s, "#e8dcc8"); px(cx + 1, base - s, "#b9ab96"); }
        const top = base - h;
        for (let x = cx - w; x <= cx + w + 1; x++) {
          px(x, top, cap[1]);
          if (x > cx - w && x < cx + w + 1) px(x, top - 1, cap[2]);
        }
        px(cx - w + 1, top - 1, cap[3]);
        px(cx, top - 2, cap[2]); px(cx + 1, top - 2, cap[3]);
        if (ripe) { px(cx - 1, top, "#8ff7ff"); px(cx + 2, top - 1, "#8ff7ff"); }
      }
    });
  }

  function reed(stage, ripe, rotten) {
    return sprite(`e:reed:${stage}:${ripe}:${rotten}`, S, S, (px) => {
      const xs = [4, 7, 10, 12, 6];
      const n = Math.min(xs.length, 2 + stage);
      for (let i = 0; i < n; i++) {
        const x = xs[i], h = 3 + stage * 2 - (i % 2);
        for (let s = 0; s < h; s++) {
          const lean = rotten ? Math.floor(s / 3) : 0;
          px(x + lean, 14 - s, rotten ? "#6b4a2a" : s < 2 ? "#3f6f2e" : "#6fae4a");
        }
        const head = rotten ? "#4e3420" : ripe ? "#ffcf5a" : "#c9913a";
        const tx = x + (rotten ? Math.floor(h / 3) : 0), ty = 14 - h;
        px(tx, ty, head); px(tx, ty - 1, head); if (ripe && !rotten) px(tx, ty - 2, "#fff2b8");
      }
    });
  }

  function crystal(size) {
    const sz = Math.max(1, Math.min(10, size || 1));
    return sprite(`e:crystal:${sz}`, S, S, (px) => {
      const r = rng(90 + sz);
      const shards = Math.min(5, 1 + Math.floor(sz / 2));
      const spots = [[8, 13], [5, 14], [11, 14], [7, 14], [10, 13]];
      for (let i = 0; i < shards; i++) {
        const [cx, base] = spots[i];
        const h = Math.min(12, 4 + Math.round(sz * 0.6) - (i ? 1 + Math.floor(r() * 2) : 0));
        for (let s = 0; s < h; s++) {
          const y = base - s, half = s > h - 3 ? 0 : 1;
          px(cx - half - 1, y, "#5c1a3e");
          px(cx - half, y, "#ff9ad5");
          px(cx, y, s % 3 === 1 ? "#ffd6ee" : "#e45aa6");
          px(cx + 1, y, "#a1336f");
          px(cx + 2, y, half ? "#5c1a3e" : null);
        }
        px(cx, base - h, "#5c1a3e");
      }
      for (let x = 4; x < 13; x++) px(x, 15, "#3a1028");
    });
  }

  function beacon(frame) {
    return sprite(`e:beacon:${frame}`, S, S, (px) => {
      for (let x = 4; x < 12; x++) { px(x, 14, "#3b3a4d"); px(x, 13, "#6b6a85"); }
      for (let y = 4; y < 13; y++) {
        const w = Math.floor((13 - y) / 3);
        px(8 - w, y, "#9c9bb6"); px(8 + w, y, "#6b6a85");
        if (y % 3 === 0) for (let x = 8 - w; x <= 8 + w; x++) px(x, y, "#8a89a3");
      }
      px(8, 3, frame ? "#ff4d3a" : "#7a2a22"); px(8, 2, frame ? "#ffb3a8" : "#5a1f1a");
      if (frame) { px(7, 3, "#ff4d3a"); px(9, 3, "#ff4d3a"); }
    });
  }

  function entity(t, now, tick) {
    switch (t.e) {
      case "beacon": return beacon(Math.floor(now / 500) % 2);
      case "crystal": return crystal(t.s);
      default: {
        const span = Math.max(1, (t.r || 0) - (t.p || 0));
        const grow = Math.max(0, Math.min(1, (tick - (t.p || 0)) / span));
        const ripe = tick >= (t.r || 0);
        const stage = ripe ? 4 : Math.min(3, Math.floor(grow * 4));
        if (t.e === "lichen") return lichen(stage, ripe);
        if (t.e === "glowcap") return glowcap(stage, ripe);
        if (t.e === "reed") return reed(stage, ripe, ripe && tick > (t.r || 0) + 90);
        return moss(stage, ripe);
      }
    }
  }

  // ---------- the drone (20x20, top-down) ----------
  function shade(hex, f) {
    const n = parseInt(hex.slice(1), 16);
    const c = [n >> 16, (n >> 8) & 255, n & 255].map((v) => Math.max(0, Math.min(255, Math.round(f > 0 ? v + (255 - v) * f : v * (1 + f)))));
    return `#${c.map((v) => v.toString(16).padStart(2, "0")).join("")}`;
  }

  function drone(color, frame) {
    const D = 20;
    return sprite(`d:${color}:${frame}`, D, D, (px) => {
      const arm = "#2a2638", armHi = "#4a4560";
      for (let i = 4; i < 16; i++) { // diagonal arms
        px(i, i, arm); px(i + 1, i, armHi);
        px(D - 1 - i, i, arm); px(D - 2 - i, i, armHi);
      }
      for (const [cx, cy] of [[4, 4], [15, 4], [4, 15], [15, 15]]) {
        for (let y = -3; y <= 3; y++) for (let x = -3; x <= 3; x++) { // spinning disc
          const d = Math.hypot(x, y);
          if (d <= 3.4 && d > 2.4) px(cx + x, cy + y, "#d8d4ee55");
        }
        const a = (frame / 4) * Math.PI + (cx + cy) * 0.3; // blades
        for (let s = -3; s <= 3; s++) px(cx + Math.round(Math.cos(a) * s), cy + Math.round(Math.sin(a) * s), "#efeaff");
        px(cx, cy, "#1d1a2c");
      }
      const body = [shade(color, -0.55), shade(color, -0.25), color, shade(color, 0.35), shade(color, 0.65)];
      for (let y = 6; y < 14; y++) for (let x = 6; x < 14; x++) {
        const corner = (x === 6 || x === 13) && (y === 6 || y === 13);
        if (corner) continue;
        const edge = x === 6 || x === 13 || y === 6 || y === 13;
        px(x, y, edge ? "#1d1a2c" : y === 7 || x === 7 ? body[3] : y === 12 || x === 12 ? body[1] : body[2]);
      }
      px(8, 7, body[4]); px(9, 7, body[4]);
      // glass dome
      px(9, 9, "#7fe3ff"); px(10, 9, "#4fb8e0"); px(9, 10, "#4fb8e0"); px(10, 10, "#2f86ad"); px(9, 9, "#dffaff");
      // status lights
      px(7, 12, frame % 2 ? "#6dff9b" : "#2f6b44");
      px(12, 12, frame % 2 ? "#ff6b5a" : "#6b2a24");
    });
  }

  function dataURL(canvas) { return canvas.toDataURL(); }

  return { S, ground, entity, drone, dataURL, ore, shade, moss, crystal, glowcap, lichen, reed };
})();
