"use strict";
// The open colony: an endless, persistent run. The player's program streams
// events live; this module owns the side panel (credits, market, shop,
// lesson rewards) and switches the shared editor and map into colony mode.
const Colony = (() => {
  const LEVEL = { id: "colony", kind: "colony", title: T("المستعمرة", "Colony"), goal: T("كبّر المستعمرة: ازرع، احصد، بيع، واشتري ترقيات", "Grow the colony: plant, harvest, sell and buy upgrades"), lesson: "", hints: [], starter: "", par: 0, gold: 0 };
  const OLD_GOODS = new Set(["moss", "lichen", "glowcap", "reed", "crystal", "iron", "geode"]);
  const RAW_GOODS = ["moss", "lichen", "glowcap", "reed", "crystal", "iron", "geode", "water", "silica", "copper", "sulfur", "titanium", "helium", "uranium"];
  const CRAFTED_GOODS = ["biomass", "plate", "wire", "glass", "fuel", "fertilizer", "circuit", "alloy", "battery", "coolant", "processor", "frame", "cell", "core"];
  const GOODS = [...RAW_GOODS, ...CRAFTED_GOODS];
  const BASE = {
    moss: 1, lichen: 3, glowcap: 2, reed: 14, crystal: 6, iron: 9, geode: 40,
    water: 2, silica: 5, copper: 12, sulfur: 18, titanium: 30, helium: 45, uranium: 90,
    biomass: 6, plate: 29, wire: 38, glass: 24, fuel: 22, fertilizer: 48, circuit: 170,
    alloy: 141, battery: 136, coolant: 78, processor: 880, frame: 941, cell: 619, core: 6818,
  };
  let data = null, es = null, running = false, samples = [], shopDirty = true, live = null;

  let contractTarget = "", contractHighlightUntil = 0, contractHighlightTimer = null;
  let rankHighlightUntil = 0, rankHighlightTimer = null;

  const fmt = (n) => Math.floor(n).toLocaleString("en-US");
  const actionMs = (st) => 260 * Math.pow(0.92, (st.levels && st.levels.engine) || 0);

  function stateToInit(st) {
    const tiles = [];
    const n = st.size;
    st.tiles.forEach((t, i) => tiles.push({ ...t, x: i % n, y: Math.floor(i / n), z: 0 }));
    (st.under || []).forEach((layer, z) => layer.forEach((t, i) => tiles.push({ ...t, x: i % n, y: Math.floor(i / n), z: z + 1 })));
    const inv = {};
    for (const [k, v] of Object.entries(st.inv || {})) if (v >= 1) inv[k] = Math.floor(v);
    return { w: n, h: n, tiles, layers: 1 + (st.under || []).length, x: 0, y: 0, z: 0, maxTicks: 0, inv, buildings: st.buildings || [] };
  }

  function liveWorld(init) {
    const keepView = W && W.colony ? { view: W.view, follow: W.follow } : null;
    setWorld(init);
    W.colony = true; W.live = true; W.liveStep = actionMs(data.state); W.t = data.state.tick || 0;
    if (keepView) { W.view = Math.min(keepView.view, W.depth - 1); W.follow = keepView.follow; renderLayerBar(); renderLegend(); }
    hud("");
  }

  // ---------- entering and leaving ----------
  async function enter() {
    if (playback) finishPlayback();
    S.mode = "colony";
    document.body.classList.add("colony");
    try {
      data = await api("/api/colony");
    } catch (e) {
      log(T("مقدرتش أفتح المستعمرة: ", "Couldn't open the colony: ") + e.message, "c-err");
      return exit();
    }
    S.review = false;
    $("#review-notice").classList.add("hidden");
    $("#btn-vscode").disabled = false;
    LEVEL.starter = data.starter;
    S.cur = LEVEL;
    $("#lvl-title").textContent = LEVEL.title;
    $("#lvl-chapter").textContent = "OPEN COLONY · KEPLER-G";
    $("#lvl-stars").innerHTML = "";
    $("#goal-strip-par").textContent = "";
    lastErrors = []; lastMtime = 0;
    showFile("main.go"); linkUI();
    if (isLinked()) syncFromWorkspace(true);
    con.innerHTML = "";
    log(T("المستعمرة: الكود هنا بيفضل شغال لحد ما تدوس Stop. كل اللي بتجمعه بيتحفظ.", "Colony: your code keeps running until you press Stop. Everything you gather is saved."), "c-info");
    log(T("بيع بـ colony.Export(colony.Moss, 0)، واشتري ترقيات من اللوحة. الأوامر المقفولة بتقولك تشتري إيه.", "Sell with colony.Export(colony.Moss, 0) and buy upgrades from the panel. Locked commands tell you what to buy."), "c-tip");
    liveWorld(stateToInit(data.state));
    running = data.running;
    samples = []; live = null;
    clearTimeout(rankHighlightTimer); rankHighlightUntil = 0;
    openStream();
    renderAll();
    runUI();
  }

  function exit() {
    TreeView.close();
    clearTimeout(contractHighlightTimer);
    contractTarget = ""; contractHighlightUntil = 0;
    clearTimeout(rankHighlightTimer);
    rankHighlightUntil = 0;
    S.mode = "mission";
    document.body.classList.remove("colony");
    if (es) { es.close(); es = null; }
    if (typeof Buildings !== "undefined") Buildings.clear();
    if (W) { W.live = false; W.colony = false; }
    runUI();
  }

  function openStream() {
    if (es) es.close();
    es = new EventSource(`/api/colony/stream?token=${encodeURIComponent(TOKEN)}&profile=${encodeURIComponent(PROFILE)}&lang=${LANG}`);
    es.onmessage = (m) => {
      let ev;
      try { ev = JSON.parse(m.data); } catch { return; }
      onEvent(ev);
    };
  }

  // ---------- live events ----------
  function onEvent(ev) {
    if (S.mode !== "colony") return;
    switch (ev.op) {
      case "init":
        data.state.tick = ev.t || data.state.tick;
        liveWorld({ ...ev, buildings: ev.buildings || data.state.buildings || [] });
        W.colony = true;
        renderBuildings();
        return;
      case "resize": {
        const drones = W.drones;
        liveWorld({ ...ev, x: 0, y: 0, maxTicks: 0, inv: W.inv, buildings: ev.buildings || data.state.buildings || [] });
        W.drones = drones;
        renderBuildings();
        log(T(`الأرض كبرت: ${ev.w}×${ev.h}`, `Field expanded: ${ev.w}×${ev.h}`), "c-ok");
        return;
      }
      case "building":
        if (W && W.colony) apply(ev);
        else if (typeof Buildings !== "undefined") Buildings.applyEvent(ev, W);
        if (data && data.state && typeof Buildings !== "undefined") data.state.buildings = Buildings.list();
        renderBuildings();
        if (typeof draw === "function") draw(performance.now());
        return;
      case "econ":
      case "state":
      case "bought":
        econ(ev);
        if (ev.op === "bought") {
          if (ev.ok) { log(T(`اتشترت: ${shopName(ev.msg)}`, `Bought: ${shopName(ev.msg)}`), "c-ok"); Sound.play("goal"); refresh(); }
          else log(T(`الشراء فشل: ${ev.msg}`, `Purchase failed: ${ev.msg}`), "c-err");
        }
        return;
      case "contract":
        econ(ev);
        live = null;
        log(ev.msg, ev.ok ? "c-ok" : "c-err");
        Sound.play(ev.ok ? "goal" : "error");
        refresh();
        return;
      case "rank":
        econ(ev);
        if (ev.msg) log(ev.msg, ev.ok === false ? "c-err" : "c-ok");
        if (ev.ok !== false) { Sound.play("goal"); highlightRank(); }
        if (!ev.tree && !TreeView.isOpen) refresh();
        return;
      case "stdout": log(ev.msg); return;
      case "stderr": log(ev.msg, "c-err"); return;
      case "stopped":
        running = false; live = null; runUI();
        log(T(`المستعمرة وقفت. ${ev.msg || ""}`, `Colony stopped. ${ev.msg || ""}`), "c-info");
        refresh();
        return;
    }
    if (W && W.colony) apply(ev);
  }

  function syncStateToWorld(st) {
    if (typeof Buildings !== "undefined" && Array.isArray(st.buildings)) Buildings.syncState(st.buildings);
    if (!W || !W.colony) return;
    if (st.size && st.size !== W.w) { liveWorld(stateToInit(st)); return; }
    let legendDirty = false;
    if (Array.isArray(st.tiles) && W.layers && W.layers[0]) {
      const n = st.size || W.w;
      st.tiles.forEach((t, i) => {
        const tile = W.layers[0][i % n] && W.layers[0][i % n][Math.floor(i / n)];
        if (!tile) return;
        if (t.b) { if (tile.b !== t.b) legendDirty = true; tile.b = t.b; }
        else if (tile.b) { delete tile.b; legendDirty = true; }
      });
    }
    for (const b of st.buildings || []) {
      const tile = W.layers && W.layers[0] && W.layers[0][b.x] && W.layers[0][b.x][b.y];
      if (tile && b.kind && tile.b !== b.kind) { tile.b = b.kind; legendDirty = true; }
    }
    if (legendDirty && typeof renderLegend === "function") renderLegend();
  }

  function econ(ev) {
    if (!data || !data.state) return;
    if (ev.state) {
      data.state = ev.state;
      W && (W.liveStep = actionMs(data.state));
      shopDirty = true;
      syncStateToWorld(ev.state);
    }
    if (ev.rank && typeof ev.rank === "object") {
      data.rank = ev.rank;
      if (typeof ev.rank.level === "number") data.state.rank = ev.rank.level;
    }
    if (typeof ev.rp === "number") data.state.rp = ev.rp;
    if (typeof ev.credits === "number") data.state.credits = ev.credits;
    if (ev.prices) data.prices = ev.prices;
    if (ev.inv) {
      data.state.inv = ev.inv;
      if (W) { W.inv = {}; for (const [k, v] of Object.entries(ev.inv)) if (v >= 1) W.inv[k] = Math.floor(v); hud(); }
    }
    if (typeof ev.t === "number" && W) W.t = Math.max(W.t, ev.t);
    if (ev.contract && ev.op !== "contract") { live = { id: ev.contract.id, frac: ev.contract.frac, left: ev.contract.left, progress: ev.progress || "" }; renderContracts(); }
    const now = performance.now();
    samples.push([now, data.state.credits]);
    while (samples.length && now - samples[0][0] > 90000) samples.shift();
    renderRank(); renderStatus(); renderMarket(); renderBuildings(); renderFleet();
    renderShop();
    if (ev.state && typeof draw === "function") draw(now);
    if (ev.tree) { data.tree = ev.tree; updateTree(); }
    else if (TreeView.isOpen && (ev.state || typeof ev.credits === "number" || ev.inv || ev.rank || typeof ev.rp === "number")) { updateTree(); refresh(); }
  }

  async function refresh() {
    try {
      const fresh = await api("/api/colony");
      if (S.mode !== "colony") return;
      data = { ...data, ...fresh };
    } catch { return; }
    if (data && data.state) syncStateToWorld(data.state);
    shopDirty = true;
    renderAll();
  }

  // ---------- running ----------
  async function toggleRun() {
    const btn = $("#btn-run");
    if (running) {
      btn.disabled = true;
      try { await api("/api/colony/stop", {}); } catch (e) { log(T("مقدرتش أوقفها: ", "Couldn't stop it: ") + e.message, "c-err"); }
      btn.disabled = false;
      running = false; runUI(); refresh();
      return;
    }
    if (isLinked()) await syncFromWorkspace(true);
    else if (!ed.readOnly) saveFile(LEVEL, "main.go", ed.value);
    store.save();
    btn.classList.add("busy"); btn.querySelector("span").textContent = "Build…";
    con.innerHTML = "";
    Sound.play("build");
    try {
      const r = await api("/api/colony/run", { code: fileCode(LEVEL, levelFiles(LEVEL)[0]) });
      if (r.phase === "requires") {
        const c = (data.contracts || []).find((x) => x.id === r.contract);
        log(T(`الكود لسه مش جاهز للعقد «${c ? c.title : r.contract}»:`, `Your code isn't ready for the contract “${c ? c.title : r.contract}”:`), "c-err");
        for (const m of r.missing || []) log("   " + m, "c-tip");
        Sound.play("error");
        return;
      }
      if (r.phase === "compile") {
        lastErrors = r.compileErrors || [];
        log(T("الكود مااتبناش. صلّح الأخطاء دي:", "The code didn't build. Fix these errors:"), "c-err");
        for (const e of lastErrors) {
          log(`${e.file || "main.go"}:${e.line}:${e.col}  ${e.msg}`, "c-err");
          if (e.tip) log(`   ${e.tip}`, "c-tip");
        }
        if (!lastErrors.length) log(r.compileOutput, "c-err");
        showFile("main.go");
        Sound.play("error");
        return;
      }
      lastErrors = []; showFile("main.go");
      running = true;
      log(T("المستعمرة شغالة.", "Colony running."), "c-ok");
    } catch (e) {
      log(T("خطأ في الاتصال بالسيرفر: ", "Server connection error: ") + e.message, "c-err");
    } finally {
      btn.classList.remove("busy");
      runUI();
    }
  }

  function runUI() {
    const btn = $("#btn-run");
    const colony = S.mode === "colony";
    btn.classList.toggle("stop", colony && running);
    btn.querySelector("span").textContent = colony && running ? "Stop" : "Run";
    btn.querySelector("use").setAttribute("href", colony && running ? "#i-stop" : "#i-play");
    $("#btn-colony").classList.toggle("linked", colony);
    const pill = $("#col-state");
    if (pill) { pill.textContent = running ? T("شغالة", "Running") : T("واقفة", "Stopped"); pill.classList.toggle("on", running); }
    if (colony && data) renderContracts();
  }

  // ---------- panel ----------
  function renderAll() { renderRank(); renderStatus(); renderMarket(); renderBuildings(); renderFleet(); renderRewards(); renderContracts(); shopDirty = true; renderShop(); updateTree(); }

  // ---------- contracts: code challenges that unlock blueprints ----------
  // ticks are 100ms of colony time
  const clock = (ticks) => { const s = Math.max(0, Math.round(ticks / 10)); return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, "0")}`; };

  function contractCard(c, now, running) {
    const run = c.run || {};
    const id = esc(c.id).replace(/"/g, "&quot;");
    const best = run.best ? `<span class="ct-best">${T("أحسن وقت", "Best time")} ${clock(run.best)}</span>` : "";
    const head = `<div class="item-top"><b>${esc(c.title)}</b><span class="lv">${T(`فصل ${c.chapter}`, `Ch. ${c.chapter}`)}</span></div>`;
    const badges = [
      c.repeatable ? `<span class="ct-repeat" title="${esc(T("عقد متكرر", "Repeatable")).replace(/"/g, "&quot;")}">↻ ${esc(T("متكرر", "Repeatable"))}</span>` : "",
      c.rewardRP > 0 ? `<span class="ct-rp"><b dir="ltr">+${Math.floor(c.rewardRP).toLocaleString("en-US")}</b> ${esc(T("نقطة بحث", "RP"))}</span>` : "",
    ].filter(Boolean).join("");
    const meta = badges ? `<div class="ct-badges">${badges}</div>` : "";
    const lesson = `<button class="lesson-link ct-lesson" data-ch="${c.chapter}">${T(`الدرس: فصل ${c.chapter}`, `Lesson: chapter ${c.chapter}`)}</button>`;
    const stopped = running ? `<span class="why">${T("وقّف المستعمرة الأول", "Stop the colony first")}</span>` : "";
    switch (c.status) {
      case "active": {
        const frac = now ? now.frac : c.frac || 0, left = now ? now.left : c.left, text = now ? now.progress : c.progress || "";
        return `<div class="item contract active" data-id="${id}">${head}${meta}<p>${esc(c.brief)}</p>
          <div class="ct-bar" role="progressbar" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${Math.round(frac * 100)}"><i style="width:${(frac * 100).toFixed(1)}%"></i></div>
          <div class="item-foot"><span class="ct-progress">${esc(text)}</span>${left >= 0 ? `<span class="ct-left">${T("باقي", "Left")} ${clock(left)}</span>` : ""}</div>
          <div class="item-foot">${lesson}<button class="ct-act ct-cancel" ${running ? "disabled" : ""}>${T("ألغي", "Cancel")}</button></div>${stopped}</div>`;
      }
      case "available":
      case "done":
        return `<div class="item contract${c.status === "done" ? " done" : ""}" data-id="${id}">${head}${meta}<p>${esc(c.brief)}</p>
          <div class="item-foot">${c.status === "done" ? best : lesson}<button class="ct-act ct-accept" data-id="${id}" ${running ? "disabled" : ""}>${c.status === "done" ? T("العبه تاني", "Replay") : T("اقبل", "Accept")}</button></div>${stopped}</div>`;
      default:
        return `<div class="item contract locked" data-id="${id}">${head}${meta}<span class="why">${(c.after || []).length ? esc(T("خلّص العقد الأول: ", "Finish the contract first: ") + c.after.join(T("، ", ", "))) : T(`خلّص الفصل ${c.chapter}`, `Finish chapter ${c.chapter}`)}</span></div>`;
    }
  }

  function renderContracts() {
    const list = data.contracts || [];
    const box = $("#col-contracts");
    const active = list.find((c) => c.status === "active");
    const now = active && live && live.id === active.id ? live : null;
    $("#goal-strip-text").textContent = active ? `${T("عقد", "Contract")}: ${active.title} — ${(now ? now.progress : active.progress) || ""}` : LEVEL.goal;
    if (!list.length) { box.innerHTML = ""; return; }
    const visible = visibleContracts(list, contractTarget);
    const done = list.filter((c) => c.status === "done").length;
    box.innerHTML = `<h3>${T("العقود", "Contracts")} <small>${done}/${list.length}</small></h3>` + visible.map((c) => contractCard(c, c.status === "active" ? now : null, running)).join("");
    if (performance.now() < contractHighlightUntil) {
      [...box.querySelectorAll(".contract")].find((card) => card.dataset.id === contractTarget)?.classList.add("tree-contract-highlight");
    }
    box.querySelectorAll(".ct-accept").forEach((b) => (b.onclick = () => contract(b.dataset.id)));
    box.querySelectorAll(".ct-cancel").forEach((b) => (b.onclick = () => contract("")));
    box.querySelectorAll(".ct-lesson").forEach((b) => (b.onclick = () => {
      const first = S.levels.find((l) => l.chapter === +b.dataset.ch);
      if (first) loadLevel(first.id);
    }));
  }

  function visibleContracts(list, targetId) {
    const order = { active: 0, available: 1, done: 2, locked: 3 };
    const shown = [...list].sort((a, b) => order[a.status] - order[b.status] || a.chapter - b.chapter);
    const lockedAt = shown.findIndex((c) => c.status === "locked");
    return lockedAt < 0 ? shown : shown.filter((c, i) => i < lockedAt + 2 || c.id === targetId);
  }

  function treeSnapshot() { return { nodes: data.tree, inv: data.state.inv, credits: data.state.credits, rp: data.state.rp, rank: data.state.rank }; }

  function updateTree() {
    if (TreeView.isOpen) TreeView.update(treeSnapshot());
  }

  async function openTree() {
    await refresh();
    if (S.mode !== "colony") return;
    TreeView.open(treeSnapshot(), {
      onBuy: (id) => buy(id, document.querySelector(".treeview .tv-buy")),
      onOpenContract: openTreeContract,
    });
  }

  function openTreeContract(ref) {
    TreeView.close();
    contractTarget = ref; contractHighlightUntil = performance.now() + 1800;
    clearTimeout(contractHighlightTimer);
    renderContracts();
    $("#col-contracts").querySelector(".tree-contract-highlight")?.scrollIntoView({ block: "center", behavior: "smooth" });
    contractHighlightTimer = setTimeout(() => {
      $("#col-contracts").querySelector(".tree-contract-highlight")?.classList.remove("tree-contract-highlight");
    }, 1800);
  }

  async function contract(id) {
    try {
      const r = await api("/api/colony/contract", { id });
      if (r.ok) {
        const c = (r.contracts || []).find((x) => x.id === id);
        log(id ? T(`قبلت العقد: ${c ? c.title : id}. دوس Run عشان يبدأ.`, `Contract accepted: ${c ? c.title : id}. Press Run to start.`) : T("العقد اتلغى.", "Contract cancelled."), "c-info");
      } else log(r.error, "c-err");
      data.state = r.state; data.contracts = r.contracts; live = null;
      await refresh();
    } catch (e) {
      log(e.message.includes("409") ? T("وقّف المستعمرة الأول.", "Stop the colony first.") : T("مقدرتش أغيّر العقد: ", "Couldn't change the contract: ") + e.message, "c-err");
    }
  }

  function renderStatus() {
    setDigits($("#col-credits"), fmt(data.state.credits));
    let rate = "—";
    if (samples.length > 2) {
      const [t0, c0] = samples[0], [t1, c1] = samples[samples.length - 1];
      if (t1 - t0 > 8000) rate = `${c1 - c0 >= 0 ? "+" : ""}${fmt(((c1 - c0) / (t1 - t0)) * 60000)} ${T("/ دقيقة", "/ min")}`;
    }
    $("#col-rate").textContent = rate;
    const owned = Object.values(data.state.levels || {}).reduce((a, b) => a + b, 0);
    $("#col-levels").textContent = data.levels ? `${owned} / ${data.levels.total} ${T("ترقية", "upgrades")}` : "";
    $("#col-size").textContent = `${data.state.size}×${data.state.size}${data.state.depth > 1 ? ` · ${data.state.depth} ${T("طبقات", "layers")}` : ""}`;
    const sys = $("#col-sys");
    if (sys) {
      const extra = statusExtraMarkup(data.state);
      sys.classList.toggle("hidden", !extra);
      sys.innerHTML = extra;
    }
  }

  function marketRowVisible(g, inv, levels) {
    if ((inv[g] || 0) >= 1) return true;
    if (!OLD_GOODS.has(g)) return false;
    return g === "moss" || !!levels["crop." + g] || ((g === "iron" || g === "geode") && !!levels["api.dig"]);
  }

  function marketRows(list, inv, prices) {
    return list.map((g) => {
      const p = prices[g] ?? BASE[g], low = p < BASE[g] * 0.97;
      return `<tr><td><i style="background:${INV_COLOR[g] || "#ccc"}"></i>${esc(NAME[g] || g)}</td><td class="num">${fmt(inv[g] || 0)}</td><td class="num${low ? " low" : ""}">${p.toFixed(2)}</td></tr>`;
    }).join("");
  }

  function marketTable(raw, crafted, inv, prices) {
    const group = (label, list) => (list.length ? `<tr class="market-group"><td colspan="3">${esc(label)}</td></tr>` + marketRows(list, inv, prices || {}) : "");
    return `<tr><td></td><td>${T("مخزون", "Stock")}</td><td>${T("السعر", "Price")}</td></tr>`
      + group(T("خامات", "Raw"), raw)
      + group(T("منتجات مصنّعة", "Crafted"), crafted);
  }

  function marketMarkup(state, prices) {
    const inv = state.inv || {}, levels = state.levels || {};
    return marketTable(
      RAW_GOODS.filter((g) => marketRowVisible(g, inv, levels)),
      CRAFTED_GOODS.filter((g) => marketRowVisible(g, inv, levels)),
      inv,
      prices,
    );
  }

  function renderMarket() {
    $("#col-market").innerHTML = marketMarkup(data.state, data.prices);
  }

  function renderBuildings() {
    const box = $("#col-buildings");
    if (!box) return;
    const list = typeof Buildings !== "undefined" ? Buildings.list() : [];
    box.classList.toggle("hidden", !list.length);
    box.innerHTML = list.length ? Buildings.listMarkup(list) : "";
  }

  // ---------- system readouts ----------
  const fmtNum = (n) => {
    const v = Number(n) || 0;
    return Number.isInteger(v) ? v.toLocaleString("en-US") : Number(v.toFixed(2)).toLocaleString("en-US");
  };
  const escAttr = (s, escapeFn = esc) => escapeFn(String(s)).replace(/"/g, "&quot;").replace(/'/g, "&#39;");

  function rankGoalLabel(kind, tr = T) {
    if (kind === "output") return tr("إنتاج في الدقيقة", "Output per minute");
    if (kind === "efficiency") return tr("قيمة لكل أمر", "Value per command");
    return kind ? String(kind) : "";
  }

  function rankMarkup(rank, deps = {}) {
    if (!rank || typeof rank !== "object") return "";
    const tr = deps.T || T;
    const escapeFn = deps.esc || esc;
    const level = Number(rank.level) || 0;
    const value = Number(rank.value) || 0;
    const target = Number(rank.target) || 0;
    const pct = target > 0 ? Math.max(0, Math.min(100, (value / target) * 100)) : (value > 0 ? 100 : 0);
    const kindText = rankGoalLabel(rank.kind, tr);
    const filling = typeof rank.filled === "number" && rank.filled < 1;
    const fillPct = filling ? Math.max(0, Math.round(rank.filled * 100)) : 100;
    const note = filling
      ? `<div class="rank-note">${escapeFn(tr(`نافذة القياس بتتملي (${fillPct}%)`, `Measuring window filling (${fillPct}%)`))}</div>`
      : "";
    return `<div class="rank-top"><span class="rank-lv">${escapeFn(tr(`المرتبة ${fmt(level)}`, `Rank ${fmt(level)}`))}</span><span class="rank-kind">${escapeFn(kindText)}</span><span class="rank-val" dir="ltr">${fmtNum(value)} / ${fmtNum(target)}</span></div>`
      + `<div class="rank-bar" role="progressbar" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${Math.round(pct)}"><i style="width:${pct.toFixed(1)}%"></i></div>`
      + note;
  }

  function statusExtraMarkup(state, deps = {}) {
    if (!state || typeof state !== "object") return "";
    const tr = deps.T || T;
    const escapeFn = deps.esc || esc;
    const parts = [];

    if (state.power && typeof state.power === "object" && (typeof state.power.supply === "number" || typeof state.power.demand === "number")) {
      const supply = Number(state.power.supply) || 0;
      const demand = Number(state.power.demand) || 0;
      const warn = demand > supply;
      const dayMark = typeof state.power.day === "boolean"
        ? `<span class="col-day ${state.power.day ? "is-day" : "is-night"}">${escapeFn(state.power.day ? tr("نهار", "Day") : tr("ليل", "Night"))}</span>`
        : "";
      parts.push(`<div class="col-stat col-power${warn ? " warn" : ""}" title="${escAttr(tr("إمداد / طلب", "Supply / Demand"), escapeFn)}"><span class="col-stat-lbl">${escapeFn(tr("الطاقة", "Power"))}${dayMark}</span><span class="col-stat-val" dir="ltr">${fmtNum(supply)} / ${fmtNum(demand)}</span></div>`);
    }

    if (typeof state.rp === "number" && Number.isFinite(state.rp)) {
      parts.push(`<div class="col-stat col-rp"><span class="col-stat-lbl">${escapeFn(tr("نقاط البحث", "Research"))}</span><span class="col-stat-val"><b dir="ltr">${fmt(state.rp)}</b> ${escapeFn(tr("نقطة", "RP"))}</span></div>`);
    }

    if (typeof state.rent === "number" && Number.isFinite(state.rent)) {
      parts.push(`<div class="col-stat col-rent"><span class="col-stat-lbl">${escapeFn(tr("الإيجار", "Rent"))}</span><span class="col-stat-val"><b dir="ltr">${fmtNum(state.rent)} CR</b> ${escapeFn(tr("/ دقيقة", "/ min"))}</span></div>`);
    }

    if (state.storage && typeof state.storage === "object") {
      const cats = [
        ["crops", tr("محاصيل", "Crops")],
        ["raw", tr("خامات", "Raw")],
        ["crafted", tr("مصنّع", "Crafted")],
      ];
      const bars = cats.filter(([k]) => state.storage[k] && typeof state.storage[k] === "object").map(([k, label]) => {
        const used = Number(state.storage[k].used) || 0;
        const cap = Number(state.storage[k].cap) || 0;
        const pct = cap > 0 ? Math.max(0, Math.min(100, (used / cap) * 100)) : (used > 0 ? 100 : 0);
        const warn = cap > 0 ? used * 10 > cap * 9 : used > 0;
        return `<div class="col-bar-row${warn ? " warn" : ""}" data-storage="${k}"><div class="col-bar-head"><span>${escapeFn(label)}</span><span class="col-bar-num" dir="ltr">${fmtNum(used)} / ${fmtNum(cap)}</span></div><div class="col-bar" role="progressbar" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${Math.round(pct)}"><i style="width:${pct.toFixed(1)}%"></i></div></div>`;
      }).join("");
      if (bars) parts.push(`<div class="col-storage">${bars}</div>`);
    }

    return parts.join("");
  }

  function hasCargoData(drone) {
    if (!drone || typeof drone !== "object") return false;
    if (Number(drone.cap) > 0) return true;
    return Object.values(drone.cargo || {}).some((n) => Number(n) > 0);
  }

  function fleetVisible(fleet) {
    if (!Array.isArray(fleet) || !fleet.length) return false;
    if (fleet.length === 1 && !hasCargoData(fleet[0])) return false;
    return true;
  }

  function fleetMarkup(fleet, deps = {}) {
    if (!fleetVisible(fleet)) return "";
    const tr = deps.T || T;
    const escapeFn = deps.esc || esc;
    const names = deps.NAME || NAME;
    const colors = deps.INV_COLOR || INV_COLOR;
    const droneColors = deps.DRONE_COLORS || (typeof DRONE_COLORS !== "undefined" ? DRONE_COLORS : ["#ff7a52"]);
    const sorted = [...fleet].sort((a, b) => (Number(a.id) || 0) - (Number(b.id) || 0));
    const rows = sorted.map((d) => {
      const id = Number(d.id) || 0;
      const col = escAttr(droneColors[Math.abs(id) % droneColors.length], escapeFn);
      const bat = typeof d.battery === "number" ? Math.max(0, Math.min(100, d.battery)) : null;
      const batLow = bat !== null && bat < 20;
      const batMarkup = bat !== null
        ? `<div class="flt-bat${batLow ? " warn" : ""}"><span class="flt-lbl">${escapeFn(tr("البطارية", "Battery"))}</span><div class="flt-bar" role="progressbar" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${Math.round(bat)}"><i style="width:${bat.toFixed(1)}%"></i></div><span class="flt-num" dir="ltr">${Math.round(bat)}%</span></div>`
        : "";
      const entries = Object.entries(d.cargo || {}).filter(([, n]) => Number(n) > 0).sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0));
      const used = entries.reduce((sum, [, n]) => sum + Number(n), 0);
      const cap = Number(d.cap) || 0;
      const chips = entries.length
        ? entries.map(([g, n]) => `<span class="cg bld-chip" title="${escAttr(names[g] || g, escapeFn)}"><i style="background:${escAttr(colors[g] || "#ccc", escapeFn)}"></i>${fmtNum(n)}</span>`).join("")
        : '<span class="bld-none">—</span>';
      const cargoMarkup = cap > 0 || used > 0
        ? `<div class="flt-cargo"><div class="flt-cargo-head"><span class="flt-lbl">${escapeFn(tr("الحمولة", "Cargo"))}</span><span class="flt-num" dir="ltr">${fmtNum(used)} / ${fmtNum(cap)}</span></div><div class="flt-chips">${chips}</div></div>`
        : "";
      return `<div class="item flt-row" data-drone="${id}"><div class="item-top"><b><i class="flt-dot" style="background:${col}"></i>${escapeFn(tr(`درون ${id}`, `Drone ${id}`))}</b></div>${batMarkup}${cargoMarkup}</div>`;
    }).join("");
    return `<h3>${escapeFn(tr("الأسطول", "Fleet"))} <small>${sorted.length}</small></h3>` + rows;
  }

  function highlightRank() {
    rankHighlightUntil = performance.now() + 1800;
    clearTimeout(rankHighlightTimer);
    const box = $("#col-rank");
    if (box) box.classList.add("rank-up");
    rankHighlightTimer = setTimeout(() => {
      const el = $("#col-rank");
      if (el) el.classList.remove("rank-up");
    }, 1800);
  }

  function renderRank() {
    const box = $("#col-rank");
    if (!box) return;
    const html = data && data.rank ? rankMarkup(data.rank) : "";
    box.classList.toggle("hidden", !html);
    box.classList.toggle("rank-up", !!html && performance.now() < rankHighlightUntil);
    box.innerHTML = html;
  }

  function renderFleet() {
    const box = $("#col-fleet");
    if (!box) return;
    const html = data && data.state ? fleetMarkup(data.state.fleet) : "";
    box.classList.toggle("hidden", !html);
    box.innerHTML = html;
  }

  function shopName(id) { const u = (data.shop || []).find((x) => x.id === id); return u ? uName(u) : id; }
  const uName = (u) => (LANG === "en" && u.nameEn) || u.name;
  const uDesc = (u) => (LANG === "en" && u.descEn) || u.desc;
  const uGroup = (u) => (LANG === "en" && u.groupEn) || u.group;

  // affordability is checked live against current credits and goods
  function status(u) {
    const lv = (data.state.levels || {})[u.id] || 0;
    if (lv >= u.max) return { maxed: true };
    if (u.blueprint) return { why: u.reason };
    for (const r of u.requires || []) if (!((data.state.levels || {})[r] > 0)) return { why: T("محتاج الأول: ", "Needs first: ") + shopName(r) };
    if (data.state.credits < u.cost.credits) return { why: T(`ناقصك ${fmt(u.cost.credits - data.state.credits)} كريدت`, `Need ${fmt(u.cost.credits - data.state.credits)} more credits`), money: true };
    for (const [g, n] of Object.entries(u.cost.goods || {})) if ((data.state.inv[g] || 0) < n) return { why: T(`ناقصك ${fmt(n - (data.state.inv[g] || 0))} ${NAME[g] || g}`, `Need ${fmt(n - (data.state.inv[g] || 0))} more ${NAME[g] || g}`), money: true };
    return { can: true };
  }

  let lastShopKey = "";
  function renderShop() {
    if (!data.shop) return;
    const key = data.shop.map((u) => { const s = status(u); return u.id + (u.level) + (s.can ? 1 : s.maxed ? 2 : 0); }).join("|");
    if (!shopDirty && key === lastShopKey) return;
    lastShopKey = key; shopDirty = false;
    const groups = new Map();
    for (const u of data.shop) { const gk = uGroup(u); if (!groups.has(gk)) groups.set(gk, []); groups.get(gk).push(u); }
    const box = $("#col-shop");
    box.innerHTML = [...groups].map(([g, items]) => {
      const owned = items.reduce((a, u) => a + Math.min(u.level, u.max), 0), total = items.reduce((a, u) => a + u.max, 0);
      return `<h3>${esc(g)} <small>${owned}/${total}</small></h3>` + items.map((u) => {
        const s = status(u);
        const goods = Object.entries(u.cost.goods || {}).map(([k, n]) => `<span class="cg"><i style="background:${INV_COLOR[k] || "#ccc"}"></i>${fmt(n)}</span>`).join("");
        const cost = s.maxed ? `<span class="maxed">${T("مكتمل", "Maxed")}</span>` : `<span class="cc">${fmt(u.cost.credits)}</span>${goods}`;
        const lesson = u.lesson && !s.maxed && !s.can ? `<span class="lesson-link">${T(`بتتعلمها في الفصل ${u.lesson}`, `Taught in chapter ${u.lesson}`)}</span>` : "";
        return `<div class="item${s.can ? " can" : ""}${s.maxed ? " done" : ""}">
          <div class="item-top"><b>${esc(uName(u))}</b><span class="lv">${u.max > 1 ? `${u.level}/${u.max}` : u.level ? "✓" : ""}</span></div>
          <p>${esc(uDesc(u))}</p>
          <div class="item-foot"><span class="cost">${cost}</span>${s.maxed ? "" : `<button class="buy" data-id="${u.id}" ${s.can ? "" : "disabled"}>${T("اشتري", "Buy")}</button>`}</div>
          ${!s.can && !s.maxed && !s.money ? `<span class="why">${esc(s.why)}</span>` : ""}${lesson}
        </div>`;
      }).join("");
    }).join("");
    box.querySelectorAll(".buy").forEach((b) => (b.onclick = () => buy(b.dataset.id, b)));
  }

  async function buy(id, btn) {
    btn.disabled = true; btn.textContent = "…";
    try {
      const r = await api("/api/colony/buy", { id });
      if (r.pending) return; // the running program applies it and reports "bought"
      if (r.ok) { log(T(`اتشترت: ${shopName(id)}`, `Bought: ${shopName(id)}`), "c-ok"); Sound.play("goal"); }
      else log(T(`الشراء فشل: ${r.error}`, `Purchase failed: ${r.error}`), "c-err");
      data.state = r.state; data.shop = r.shop; data.tree = r.tree; shopDirty = true;
      liveWorld(stateToInit(data.state));
      renderAll();
    } catch (e) {
      log(T("الشراء فشل: ", "Purchase failed: ") + e.message, "c-err");
      shopDirty = true; renderShop(); updateTree();
    }
  }

  // finishing every mission of a chapter pays into the colony, once
  function renderRewards() {
    const done = (S.chapters || []).filter((ch) => {
      const ls = S.levels.filter((l) => l.chapter === ch.n);
      return ls.length && ls.every((l) => rec(l.id).stars > 0) && !(data.state.claimed || {})["chapter-" + ch.n];
    });
    $("#col-rewards").innerHTML = done.length
      ? `<h3>${T("مكافآت الدروس", "Lesson rewards")}</h3>` + done.map((ch) => `<button class="claim" data-n="${ch.n}">${T(`استلم مكافأة الفصل ${ch.n}`, `Claim chapter ${ch.n} reward`)}: ${esc(ch.title)}</button>`).join("")
      : "";
    $("#col-rewards").querySelectorAll(".claim").forEach((b) => (b.onclick = () => claim(+b.dataset.n)));
  }

  async function claim(n) {
    try {
      const r = await api("/api/colony/claim", { chapter: n });
      if (r.ok) {
        log(T(`مكافأة الفصل ${n}: +${fmt(r.reward.credits)} كريدت`, `Chapter ${n} reward: +${fmt(r.reward.credits)} credits`) + ((r.reward.grants || []).length ? " + " + r.reward.grants.map(shopName).join(T("، ", ", ")) : ""), "c-ok");
        Sound.play("goal");
      } else log(r.error, "c-err");
      data.state = r.state; await refresh();
    } catch (e) {
      log(e.message.includes("409") ? T("وقّف المستعمرة الأول، وبعدين استلم المكافأة.", "Stop the colony first, then claim the reward.") : T("مقدرتش أستلم: ", "Couldn't claim: ") + e.message, "c-err");
    }
  }

  $("#btn-tech-tree").onclick = openTree;

  return { enter, exit, toggleRun, LEVEL, get running() { return running; } };
})();

$("#btn-colony").onclick = () => (S.mode === "colony" ? loadLevel(store.pref("current") || S.levels[0].id) : Colony.enter());
