"use strict";
// The open colony: an endless, persistent run. The player's program streams
// events live; this module owns the side panel (credits, market, shop,
// lesson rewards) and switches the shared editor and map into colony mode.
const Colony = (() => {
  const LEVEL = { id: "colony", kind: "colony", title: T("المستعمرة", "Colony"), goal: T("كبّر المستعمرة: ازرع، احصد، بيع، واشتري ترقيات", "Grow the colony: plant, harvest, sell and buy upgrades"), lesson: "", hints: [], starter: "", par: 0, gold: 0 };
  const GOODS = ["moss", "lichen", "glowcap", "reed", "crystal", "iron", "geode"];
  const BASE = { moss: 1, lichen: 3, glowcap: 2, reed: 14, crystal: 6, iron: 9, geode: 40 };
  let data = null, es = null, running = false, samples = [], shopDirty = true;

  const fmt = (n) => Math.floor(n).toLocaleString("en-US");
  const actionMs = (st) => 260 * Math.pow(0.92, (st.levels && st.levels.engine) || 0);

  function stateToInit(st) {
    const tiles = [];
    const n = st.size;
    st.tiles.forEach((t, i) => tiles.push({ ...t, x: i % n, y: Math.floor(i / n), z: 0 }));
    (st.under || []).forEach((layer, z) => layer.forEach((t, i) => tiles.push({ ...t, x: i % n, y: Math.floor(i / n), z: z + 1 })));
    const inv = {};
    for (const [k, v] of Object.entries(st.inv || {})) if (v >= 1) inv[k] = Math.floor(v);
    return { w: n, h: n, tiles, layers: 1 + (st.under || []).length, x: 0, y: 0, z: 0, maxTicks: 0, inv };
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
    $("#goal-strip-text").textContent = LEVEL.goal;
    $("#goal-strip-par").textContent = "";
    lastErrors = []; lastMtime = 0;
    showFile("main.go"); linkUI();
    if (isLinked()) syncFromWorkspace(true);
    con.innerHTML = "";
    log(T("المستعمرة: الكود هنا بيفضل شغال لحد ما تدوس Stop. كل اللي بتجمعه بيتحفظ.", "Colony: your code keeps running until you press Stop. Everything you gather is saved."), "c-info");
    log(T("بيع بـ colony.Export(colony.Moss, 0)، واشتري ترقيات من اللوحة. الأوامر المقفولة بتقولك تشتري إيه.", "Sell with colony.Export(colony.Moss, 0) and buy upgrades from the panel. Locked commands tell you what to buy."), "c-tip");
    liveWorld(stateToInit(data.state));
    running = data.running;
    samples = [];
    openStream();
    renderAll();
    runUI();
  }

  function exit() {
    S.mode = "mission";
    document.body.classList.remove("colony");
    if (es) { es.close(); es = null; }
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
        liveWorld(ev);
        W.colony = true;
        return;
      case "resize": {
        const drones = W.drones;
        liveWorld({ ...ev, x: 0, y: 0, maxTicks: 0, inv: W.inv });
        W.drones = drones;
        log(T(`الأرض كبرت: ${ev.w}×${ev.h}`, `Field expanded: ${ev.w}×${ev.h}`), "c-ok");
        return;
      }
      case "econ":
      case "state":
      case "bought":
        econ(ev);
        if (ev.op === "bought") {
          if (ev.ok) { log(T(`اتشترت: ${shopName(ev.msg)}`, `Bought: ${shopName(ev.msg)}`), "c-ok"); Sound.play("goal"); refresh(); }
          else log(T(`الشراء فشل: ${ev.msg}`, `Purchase failed: ${ev.msg}`), "c-err");
        }
        return;
      case "stdout": log(ev.msg); return;
      case "stderr": log(ev.msg, "c-err"); return;
      case "stopped":
        running = false; runUI();
        log(T(`المستعمرة وقفت. ${ev.msg || ""}`, `Colony stopped. ${ev.msg || ""}`), "c-info");
        refresh();
        return;
    }
    if (W && W.colony) apply(ev);
  }

  function econ(ev) {
    if (ev.state) { data.state = ev.state; W && (W.liveStep = actionMs(data.state)); shopDirty = true; }
    if (typeof ev.credits === "number") data.state.credits = ev.credits;
    if (ev.prices) data.prices = ev.prices;
    if (ev.inv) {
      data.state.inv = ev.inv;
      if (W) { W.inv = {}; for (const [k, v] of Object.entries(ev.inv)) if (v >= 1) W.inv[k] = Math.floor(v); hud(); }
    }
    if (typeof ev.t === "number" && W) W.t = Math.max(W.t, ev.t);
    const now = performance.now();
    samples.push([now, data.state.credits]);
    while (samples.length && now - samples[0][0] > 90000) samples.shift();
    renderStatus(); renderMarket();
    renderShop();
  }

  async function refresh() {
    try { data = { ...data, ...(await api("/api/colony")) }; } catch { return; }
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
  }

  // ---------- panel ----------
  function renderAll() { renderStatus(); renderMarket(); renderRewards(); shopDirty = true; renderShop(); }

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
  }

  function renderMarket() {
    const inv = data.state.inv || {}, prices = data.prices || {};
    const rows = GOODS.filter((g) => (inv[g] || 0) >= 1 || g === "moss" || (data.state.levels || {})["crop." + g] || (g === "iron" || g === "geode") && (data.state.levels || {})["api.dig"]);
    $("#col-market").innerHTML = `<tr><td></td><td>${T("مخزون", "Stock")}</td><td>${T("السعر", "Price")}</td></tr>` + rows.map((g) => {
      const p = prices[g] ?? BASE[g], low = p < BASE[g] * 0.97;
      return `<tr><td><i style="background:${INV_COLOR[g] || "#ccc"}"></i>${NAME[g] || g}</td><td class="num">${fmt(inv[g] || 0)}</td><td class="num${low ? " low" : ""}">${p.toFixed(2)}</td></tr>`;
    }).join("");
  }

  function shopName(id) { const u = (data.shop || []).find((x) => x.id === id); return u ? uName(u) : id; }
  const uName = (u) => (LANG === "en" && u.nameEn) || u.name;
  const uDesc = (u) => (LANG === "en" && u.descEn) || u.desc;
  const uGroup = (u) => (LANG === "en" && u.groupEn) || u.group;

  // affordability is checked live against current credits and goods
  function status(u) {
    const lv = (data.state.levels || {})[u.id] || 0;
    if (lv >= u.max) return { maxed: true };
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
      data.state = r.state; data.shop = r.shop; shopDirty = true;
      liveWorld(stateToInit(data.state));
      renderAll();
    } catch (e) {
      log(T("الشراء فشل: ", "Purchase failed: ") + e.message, "c-err");
      shopDirty = true; renderShop();
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

  return { enter, exit, toggleRun, LEVEL, get running() { return running; } };
})();

$("#btn-colony").onclick = () => (S.mode === "colony" ? loadLevel(store.pref("current") || S.levels[0].id) : Colony.enter());
