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
  let runs = [], raced = false, raceEnabled = false, checking = false;

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
    syncWeatherToWorld(data.state.weather);
    if (keepView) { W.view = Math.min(keepView.view, W.depth - 1); W.follow = keepView.follow; renderLayerBar(); renderLegend(); }
    hud("");
  }

  function syncWeatherToWorld(weather) {
    if (!W || !W.colony) return;
    W.storm = weather && weather.storm === "active" ? 2 : weather && weather.storm === "warn" ? 1 : 0;
  }

  function getRacePref() {
    return typeof store !== "undefined" && typeof store.pref === "function" ? !!store.pref("colonyRace") : false;
  }
  function setRacePref(on) {
    if (typeof store !== "undefined" && typeof store.pref === "function") store.pref("colonyRace", !!on);
  }
  function getFilesForRun(level) {
    if (typeof colonyFilesPayload === "function") return colonyFilesPayload(level);
    const list = typeof levelFiles === "function" ? levelFiles(level) : ["main.go"];
    const first = Array.isArray(list) && list.length ? list[0] : "main.go";
    const name = typeof first === "string" ? first : (first && first.name) || "main.go";
    return { [name]: typeof fileCode === "function" ? fileCode(level, first) : "" };
  }
  function resolveErrorFile(errs, files, fallback) {
    if (typeof firstErrorFile === "function") return firstErrorFile(errs, files, fallback);
    return (errs && errs[0] && errs[0].file) || fallback || "main.go";
  }
  function buildRunPayload(level, race) {
    return { files: getFilesForRun(level), race: !!race };
  }
  function raceUI() {
    const b = $("#btn-race");
    if (b) b.setAttribute("aria-pressed", String(!!raceEnabled));
  }
  function toggleRace() {
    raceEnabled = !raceEnabled;
    setRacePref(raceEnabled);
    raceUI();
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
    raced = !!data.raced;
    raceEnabled = getRacePref();
    runs = [];
    samples = []; live = null;
    clearTimeout(rankHighlightTimer); rankHighlightUntil = 0;
    openStream();
    renderAll();
    runUI();
    fetchRuns();
  }

  function exit() {
    TreeView.close();
    clearTimeout(contractHighlightTimer);
    contractTarget = ""; contractHighlightUntil = 0;
    clearTimeout(rankHighlightTimer);
    rankHighlightUntil = 0;
    running = false; raced = false; runs = [];
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
      case "race":
        log(ev.msg || T("اتكشف تعارض سباق (Data race).", "Data race detected."), "c-err");
        if (!raced) {
          raced = true;
          Sound.play("error");
          runUI();
        }
        return;
      case "storm": {
        const phase = ev.phase === "hit" || ev.phase === "clear" ? ev.phase : "warn";
        const cls = phase === "hit" ? "c-err" : phase === "clear" ? "c-ok" : "c-warn";
        const snd = phase === "hit" ? "storm" : phase === "clear" ? "report" : "fail";
        const fallback = phase === "hit"
          ? T("عاصفة ضربت المستعمرة!", "Storm hit the colony!")
          : phase === "clear"
            ? T("العاصفة هديت.", "Storm cleared.")
            : T("تحذير: عاصفة جاية!", "Warning: incoming storm!");
        log(ev.msg || fallback, cls);
        Sound.play(snd);
        if (ev.state) {
          econ(ev);
        } else if (data && data.state) {
          if (typeof ev.t === "number") data.state.tick = Math.max(data.state.tick || 0, ev.t);
          if (ev.weather) data.state.weather = ev.weather;
          else {
            const nextStorm = phase === "hit" ? "active" : phase === "clear" ? "none" : "warn";
            data.state.weather = { ...(data.state.weather || {}), storm: nextStorm, ...(typeof ev.at === "number" ? { at: ev.at } : {}) };
          }
          syncWeatherToWorld(data.state.weather);
          renderWeather();
          if (typeof draw === "function") draw(performance.now());
        }
        return;
      }
      case "fault": {
        const id = typeof ev.d === "number" ? ev.d : 0;
        log(ev.msg || T(`عطل في الدرون ${id}.`, `Drone ${id} faulted.`), "c-err");
        Sound.play("error");
        if (ev.state) {
          econ(ev);
        } else if (data && data.state && Array.isArray(data.state.fleet) && typeof ev.d === "number") {
          const drone = data.state.fleet.find((x) => Number(x.id) === ev.d);
          if (drone) {
            drone.fault = true;
            if (typeof ev.stunned === "number") drone.stunned = ev.stunned;
          }
          renderFleet();
        }
        return;
      }
      case "shock":
        log(ev.msg || T("صدمة في السوق!", "Market shock!"), "c-warn");
        Sound.play("fail");
        if (ev.state || ev.prices) {
          econ(ev);
        } else if (data && data.state) {
          if (typeof ev.t === "number") data.state.tick = Math.max(data.state.tick || 0, ev.t);
          if ("shock" in ev) data.state.shock = ev.shock;
          else if (ev.good) data.state.shock = { good: ev.good, factor: ev.factor, until: ev.until };
          renderMarket();
        }
        return;
      case "rush": {
        const ok = ev.ok !== false;
        log(ev.msg || T("طلب عاجل في السوق!", "Rush order in the market!"), ok ? "c-ok" : "c-err");
        Sound.play(ok ? "report" : "fail");
        if (ev.state) {
          econ(ev);
        } else if (data && data.state) {
          if (typeof ev.t === "number") data.state.tick = Math.max(data.state.tick || 0, ev.t);
          if ("rush" in ev) data.state.rush = ev.rush;
          else if (ev.good) data.state.rush = { good: ev.good, amount: ev.amount, left: ev.left ?? ev.amount, until: ev.until, bonus: ev.bonus };
          renderRush();
        }
        return;
      }
      case "stdout": log(ev.msg); return;
      case "stderr": log(ev.msg, "c-err"); return;
      case "stopped":
        running = false; live = null; runUI();
        log(T(`المستعمرة وقفت. ${ev.msg || ""}`, `Colony stopped. ${ev.msg || ""}`), "c-info");
        refresh();
        fetchRuns();
        return;
    }
    if (W && W.colony) {
      apply(ev);
      if (data && data.state && typeof ev.t === "number" && ev.t > (data.state.tick || 0)) {
        const prevSec = Math.floor((data.state.tick || 0) / 10);
        data.state.tick = ev.t;
        if (Math.floor(ev.t / 10) !== prevSec) {
          if (data.state.weather && data.state.weather.storm && data.state.weather.storm !== "none") renderWeather();
          if (data.state.rush) renderRush();
        }
      }
    }
  }

  function syncStateToWorld(st) {
    if (typeof Buildings !== "undefined" && Array.isArray(st.buildings)) Buildings.syncState(st.buildings);
    syncWeatherToWorld(st.weather);
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
    if (ev.stats) data.state.stats = ev.stats;
    if (ev.weather) data.state.weather = ev.weather;
    if ("rush" in ev) data.state.rush = ev.rush;
    if ("shock" in ev) data.state.shock = ev.shock;
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
    if (typeof ev.t === "number") {
      data.state.tick = Math.max(data.state.tick || 0, ev.t);
      if (W) W.t = Math.max(W.t, ev.t);
    }
    syncWeatherToWorld(data.state.weather);
    if (ev.contract && ev.op !== "contract") { live = { id: ev.contract.id, frac: ev.contract.frac, left: ev.contract.left, progress: ev.progress || "" }; renderContracts(); }
    const now = performance.now();
    samples.push([now, data.state.credits]);
    while (samples.length && now - samples[0][0] > 90000) samples.shift();
    renderWeather(); renderRank(); renderStatus(); renderLiveStats(); renderRush(); renderMarket(); renderBuildings(); renderFleet();
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

  async function fetchRuns() {
    try {
      const res = await api("/api/colony/runs");
      if (S.mode !== "colony") return;
      runs = Array.isArray(res && res.runs) ? res.runs : [];
      renderLastRun();
    } catch {
      // Older or offline backends may not expose run history yet.
    }
  }

  function applyCheckResult(kind, r) {
    const ok = !!(r && r.ok);
    const cls = ok ? "c-ok" : "c-err";
    const head = kind === "test"
      ? (ok ? T("الاختبارات نجحت.", "Tests passed.") : T("الاختبارات فشلت.", "Tests failed."))
      : (ok ? T("فحص vet تمام.", "Vet passed.") : T("فحص vet لقى ملاحظات.", "Vet failed."));
    log(head, cls);
    if (r && r.output) {
      const out = String(r.output).trimEnd();
      if (out) out.split("\n").forEach((line) => log(line, cls));
    }
    const errs = Array.isArray(r && r.compileErrors) ? r.compileErrors : [];
    lastErrors = errs;
    const active = (typeof S !== "undefined" && S.file) || "main.go";
    if (errs.length) {
      for (const e of errs) {
        log(`${e.file || "main.go"}:${e.line}:${e.col}  ${e.msg}`, "c-err");
        if (e.tip) log(`   ${e.tip}`, "c-tip");
      }
      showFile(resolveErrorFile(errs, levelFiles(LEVEL), active));
    } else {
      showFile(active);
    }
    Sound.play(ok ? "goal" : "error");
  }

  async function runCheck(kind) {
    if (checking) return;
    const active = (typeof S !== "undefined" && S.file) || "main.go";
    if (isLinked()) await syncFromWorkspace(true);
    else if (!ed.readOnly) saveFile(LEVEL, active, ed.value);
    store.save();
    checking = true;
    runUI();
    con.innerHTML = "";
    Sound.play("build");
    log(kind === "test" ? T("بشغّل الاختبارات…", "Running tests…") : T("بفحص الكود بـ go vet…", "Running go vet…"), "c-info");
    try {
      const r = await api(kind === "test" ? "/api/colony/test" : "/api/colony/vet", { files: getFilesForRun(LEVEL) });
      applyCheckResult(kind, r);
    } catch (e) {
      log(T("خطأ في الاتصال بالسيرفر: ", "Server connection error: ") + e.message, "c-err");
      Sound.play("error");
    } finally {
      checking = false;
      runUI();
    }
  }

  // ---------- running ----------
  async function toggleRun() {
    const btn = $("#btn-run");
    const active = (typeof S !== "undefined" && S.file) || "main.go";
    if (running) {
      btn.disabled = true;
      try { await api("/api/colony/stop", {}); } catch (e) { log(T("مقدرتش أوقفها: ", "Couldn't stop it: ") + e.message, "c-err"); }
      btn.disabled = false;
      running = false; runUI(); refresh(); fetchRuns();
      return;
    }
    if (isLinked()) await syncFromWorkspace(true);
    else if (!ed.readOnly) saveFile(LEVEL, active, ed.value);
    store.save();
    raced = false;
    runUI();
    btn.classList.add("busy"); btn.querySelector("span").textContent = "Build…";
    con.innerHTML = "";
    Sound.play("build");
    try {
      const r = await api("/api/colony/run", buildRunPayload(LEVEL, raceEnabled));
      if (r.phase === "requires") {
        const c = (data.contracts || []).find((x) => x.id === r.contract);
        log(T(`الكود لسه مش جاهز للعقد «${c ? c.title : r.contract}»:`, `Your code isn't ready for the contract “${c ? c.title : r.contract}”:`), "c-err");
        for (const m of r.missing || []) log("   " + m, "c-tip");
        Sound.play("error");
        return;
      }
      if (r.phase === "require") {
        const msg = (r.result && r.result.message) || r.message || T("الاستيراد ده مش مسموح.", "That import is not allowed.");
        log(msg, "c-err");
        for (const m of r.missing || []) log("   • " + m, "c-tip");
        if (r.result && r.result.progress) log("   " + r.result.progress, "c-tip");
        Sound.play("error");
        return;
      }
      if (r.phase === "compile") {
        lastErrors = r.compileErrors || [];
        const msg = (r.result && r.result.message) || T("الكود مااتبناش. صلّح الأخطاء دي:", "The code didn't build. Fix these errors:");
        log(msg, "c-err");
        for (const e of lastErrors) {
          log(`${e.file || "main.go"}:${e.line}:${e.col}  ${e.msg}`, "c-err");
          if (e.tip) log(`   ${e.tip}`, "c-tip");
        }
        if (!lastErrors.length && r.compileOutput) log(r.compileOutput, "c-err");
        showFile(resolveErrorFile(lastErrors, levelFiles(LEVEL), active));
        Sound.play("error");
        return;
      }
      lastErrors = []; showFile(active);
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
    if (pill) {
      const base = running ? T("شغالة", "Running") : T("واقفة", "Stopped");
      pill.textContent = raced ? `${base} · ${T("سباق", "Raced")}` : base;
      pill.classList.toggle("on", running);
      pill.classList.toggle("raced", !!raced);
    }
    const testBtn = $("#btn-test");
    const vetBtn = $("#btn-vet");
    if (testBtn) testBtn.disabled = checking;
    if (vetBtn) vetBtn.disabled = checking;
    raceUI();
    if (colony) {
      renderLastRun();
      renderLiveStats();
      if (data) renderContracts();
    }
  }

  // ---------- panel ----------
  function renderAll() { renderWeather(); renderRank(); renderStatus(); renderLastRun(); renderLiveStats(); renderRush(); renderMarket(); renderBuildings(); renderFleet(); renderRewards(); renderContracts(); shopDirty = true; renderShop(); updateTree(); }

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
    $("#col-levels").textContent = data.levels
      ? (data.levels.total > 0 ? `${owned} / ${data.levels.total} ${T("ترقية", "upgrades")}` : `${owned} ${T("ترقية", "upgrades")}`)
      : "";
    $("#col-size").textContent = `${data.state.size}×${data.state.size}${data.state.depth > 1 ? ` · ${data.state.depth} ${T("طبقات", "layers")}` : ""}`;
    const sys = $("#col-sys");
    if (sys) {
      const extra = statusExtraMarkup(data.state);
      sys.classList.toggle("hidden", !extra);
      sys.innerHTML = extra;
    }
  }

  function marketRowVisible(g, inv, levels, shock) {
    if (shock && shock.good === g) return true;
    if ((inv[g] || 0) >= 1) return true;
    if (!OLD_GOODS.has(g)) return false;
    return g === "moss" || !!levels["crop." + g] || ((g === "iron" || g === "geode") && !!levels["api.dig"]);
  }

  function marketRows(list, inv, prices, shock, tick) {
    return list.map((g) => {
      const p = prices[g] ?? BASE[g], low = p < BASE[g] * 0.97;
      const shocked = !!(shock && shock.good === g);
      const factorText = shocked && typeof shock.factor === "number" ? ` ×${Number(shock.factor.toFixed(2))}` : "";
      const secLeft = shocked && typeof shock.until === "number" && typeof tick === "number" && shock.until > tick
        ? Math.max(0, Math.round((shock.until - tick) / 10))
        : null;
      const shockTitle = secLeft !== null
        ? ` title="${esc(T("باقي", "Left") + " " + Math.floor(secLeft / 60) + ":" + String(secLeft % 60).padStart(2, "0")).replace(/"/g, "&quot;")}"`
        : "";
      const shockMark = shocked
        ? ` <span class="shock-badge"${shockTitle}>${esc(T("صدمة", "Shock") + factorText)}</span>`
        : "";
      return `<tr${shocked ? ' class="market-shock"' : ""}><td><i style="background:${INV_COLOR[g] || "#ccc"}"></i>${esc(NAME[g] || g)}${shockMark}</td><td class="num">${fmt(inv[g] || 0)}</td><td class="num${low ? " low" : ""}">${p.toFixed(2)}</td></tr>`;
    }).join("");
  }

  function marketTable(raw, crafted, inv, prices, shock, tick) {
    const group = (label, list) => (list.length ? `<tr class="market-group"><td colspan="3">${esc(label)}</td></tr>` + marketRows(list, inv, prices || {}, shock, tick) : "");
    return `<tr><td></td><td>${T("مخزون", "Stock")}</td><td>${T("السعر", "Price")}</td></tr>`
      + group(T("خامات", "Raw"), raw)
      + group(T("منتجات مصنّعة", "Crafted"), crafted);
  }

  function marketMarkup(state, prices) {
    const inv = state.inv || {}, levels = state.levels || {}, shock = state.shock || null;
    const tick = typeof state.tick === "number" ? state.tick : undefined;
    return marketTable(
      RAW_GOODS.filter((g) => marketRowVisible(g, inv, levels, shock)),
      CRAFTED_GOODS.filter((g) => marketRowVisible(g, inv, levels, shock)),
      inv,
      prices,
      shock,
      tick,
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

  function hasFleetAlert(drone) {
    return !!(drone && (drone.fault || Number(drone.stunned) > 0));
  }

  function fleetVisible(fleet) {
    if (!Array.isArray(fleet) || !fleet.length) return false;
    if (fleet.length === 1 && !hasCargoData(fleet[0]) && !hasFleetAlert(fleet[0])) return false;
    return true;
  }

  const clockFmt = (ticks) => {
    const s = Math.max(0, Math.round((Number(ticks) || 0) / 10));
    return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, "0")}`;
  };

  function runComparison(last, prev) {
    if (!last || !prev) return null;
    const prevRate = Number(prev.perMin);
    const lastRate = Number(last.perMin);
    if (!Number.isFinite(prevRate) || prevRate <= 0 || !Number.isFinite(lastRate)) return null;
    const raw = Math.round(((lastRate - prevRate) / prevRate) * 100);
    const pct = Object.is(raw, -0) ? 0 : raw;
    return {
      pct,
      text: `${pct > 0 ? "+" : ""}${pct}%`,
      cls: pct > 0 ? "up" : pct < 0 ? "down" : "flat",
    };
  }

  function droneShares(d) {
    const raw = [
      Math.max(0, Number(d && d.work) || 0),
      Math.max(0, Number(d && d.move) || 0),
      Math.max(0, Number(d && d.wait) || 0),
      Math.max(0, Number(d && d.failed) || 0),
    ];
    const total = raw[0] + raw[1] + raw[2] + raw[3];
    if (!(total > 0)) return { work: 0, move: 0, wait: 0, failed: 0, total: 0 };
    const exact = raw.map((v) => (v / total) * 100);
    const ints = exact.map((x) => Math.floor(x));
    let rem = 100 - (ints[0] + ints[1] + ints[2] + ints[3]);
    const order = [0, 1, 2, 3]
      .filter((i) => raw[i] > 0)
      .sort((a, b) => (exact[b] - ints[b]) - (exact[a] - ints[a]) || a - b);
    for (let k = 0; rem > 0 && k < order.length; k++) {
      ints[order[k]] += 1;
      rem -= 1;
    }
    return { work: ints[0], move: ints[1], wait: ints[2], failed: ints[3], total };
  }

  function droneBarsMarkup(stats, deps = {}) {
    if (!stats || !Array.isArray(stats.drones) || !stats.drones.length) return "";
    const tr = deps.T || T;
    const escapeFn = deps.esc || esc;
    const droneColors = deps.DRONE_COLORS || (typeof DRONE_COLORS !== "undefined" ? DRONE_COLORS : ["#ff7a52"]);
    const compact = !!deps.compact;
    const keys = [
      ["work", tr("عمل", "Work")],
      ["move", tr("حركة", "Move")],
      ["wait", tr("انتظار", "Wait")],
      ["failed", tr("فشل", "Failed")],
    ];
    const sorted = [...stats.drones].sort((a, b) => (Number(a.id) || 0) - (Number(b.id) || 0));
    const rows = sorted.map((d) => {
      const id = Number(d.id) || 0;
      const col = escAttr(droneColors[Math.abs(id) % droneColors.length], escapeFn);
      const s = droneShares(d);
      const segs = keys
        .filter(([k]) => s[k] > 0)
        .map(([k, label]) => `<i class="lr-seg lr-seg-${k}" style="width:${s[k]}%" data-share="${s[k]}" title="${escAttr(`${label}: ${s[k]}%`, escapeFn)}"></i>`)
        .join("");
      const ariaParts = keys.map(([k, label]) => `${label} ${s[k]}%`).join(", ");
      const pctSummary = compact
        ? `${s.work}%`
        : `${s.work}% / ${s.move}% / ${s.wait}% / ${s.failed}%`;
      return `<div class="lr-drone" data-drone="${id}"><span class="lr-drone-id"><i class="flt-dot" style="background:${col}"></i>${escapeFn(tr(`د${id}`, `D${id}`))}</span><div class="lr-bar" role="img" aria-label="${escAttr(ariaParts, escapeFn)}">${segs}</div><span class="lr-drone-pct" dir="ltr">${pctSummary}</span></div>`;
    }).join("");
    const legend = `<div class="lr-legend">${keys.map(([k, label]) => `<span class="lr-leg-item"><i class="lr-seg-${k}"></i>${escapeFn(label)}</span>`).join("")}</div>`;
    return `<div class="lr-drones${compact ? " compact" : ""}">${rows}</div>${legend}`;
  }

  function lastRunMarkup(runsList, deps = {}) {
    if (!Array.isArray(runsList) || !runsList.length) return "";
    const tr = deps.T || T;
    const escapeFn = deps.esc || esc;
    const last = runsList[0];
    if (!last || typeof last !== "object") return "";
    const prev = runsList.length > 1 ? runsList[1] : null;
    const cmp = runComparison(last, prev);
    const cmpHtml = cmp ? `<span class="lr-delta ${cmp.cls}" dir="ltr">${escapeFn(cmp.text)}</span>` : "";
    const raceBadge = last.race ? `<span class="lr-race-tag">-race</span>` : "";
    const stats = last.stats && typeof last.stats === "object" ? last.stats : null;
    const barsHtml = stats ? droneBarsMarkup(stats, { ...deps, compact: false }) : "";
    const totalIdle = stats
      ? (typeof stats.idleMoves === "number"
        ? stats.idleMoves
        : (Array.isArray(stats.drones) ? stats.drones.reduce((sum, d) => sum + (Number(d.idleMoves) || 0), 0) : 0))
      : 0;
    const topFail = stats && stats.topFail ? String(stats.topFail) : "";
    const footHtml = stats
      ? `<div class="lr-meta"><span>${escapeFn(tr("أكتر أمر فشل", "Top failed"))}: <code dir="ltr">${escapeFn(topFail || tr("مفيش", "None"))}</code></span><span>${escapeFn(tr("حركات فاضية", "Idle moves"))}: <b dir="ltr">${fmtNum(totalIdle)}</b></span></div>`
      : "";
    return `<h3>${escapeFn(tr("آخر تشغيل", "Last run"))}${raceBadge}</h3>`
      + `<div class="lr-metrics">`
      + `<div class="lr-metric"><span class="lr-k">${escapeFn(tr("كريدت / دقيقة", "CR / min"))}</span><span class="lr-v"><b dir="ltr">${fmtNum(last.perMin || 0)}</b>${cmpHtml}</span></div>`
      + `<div class="lr-metric"><span class="lr-k">${escapeFn(tr("التكات", "Ticks"))}</span><span class="lr-v"><b dir="ltr">${fmtNum(last.ticks || 0)}</b></span></div>`
      + `<div class="lr-metric"><span class="lr-k">${escapeFn(tr("المكسب", "Earned"))}</span><span class="lr-v"><b dir="ltr">${fmtNum(last.earned || 0)}</b></span></div>`
      + `</div>`
      + barsHtml
      + footHtml;
  }

  function liveStatsMarkup(stats, deps = {}) {
    if (!stats || !Array.isArray(stats.drones) || !stats.drones.length) return "";
    const tr = deps.T || T;
    const escapeFn = deps.esc || esc;
    return `<h3>${escapeFn(tr("إحصائيات الدرونز", "Drone stats"))}</h3>${droneBarsMarkup(stats, { ...deps, compact: true })}`;
  }

  function weatherMarkup(weather, tick, deps = {}) {
    if (!weather || !weather.storm || weather.storm === "none") return "";
    const tr = deps.T || T;
    const escapeFn = deps.esc || esc;
    const phase = weather.storm === "active" || weather.storm === "hit" ? "active" : "warn";
    const hasAt = typeof weather.at === "number";
    const rem = hasAt ? Math.max(0, weather.at - (Number(tick) || 0)) : 0;
    const text = phase === "active"
      ? (hasAt ? tr(`عاصفة شغالة — تنتهي خلال ${clockFmt(rem)}`, `Storm active — clears in ${clockFmt(rem)}`) : tr("عاصفة شغالة", "Storm active"))
      : (hasAt ? tr(`تحذير عاصفة خلال ${clockFmt(rem)}`, `Storm incoming in ${clockFmt(rem)}`) : tr("تحذير: عاصفة جاية", "Storm warning"));
    return `<div class="col-storm is-${phase}"><span class="col-storm-dot" aria-hidden="true"></span><span>${escapeFn(text)}</span></div>`;
  }

  function rushMarkup(rush, tick, deps = {}) {
    if (!rush || !rush.good) return "";
    const tr = deps.T || T;
    const escapeFn = deps.esc || esc;
    const names = deps.NAME || NAME;
    const colors = deps.INV_COLOR || INV_COLOR;
    const goodName = names[rush.good] || rush.good;
    const chipColor = escAttr(colors[rush.good] || "#ccc", escapeFn);
    const total = Number(rush.amount) || 0;
    const left = Math.max(0, Number(rush.left ?? total) || 0);
    const done = total > 0 ? Math.max(0, total - left) : 0;
    const pct = total > 0 ? Math.max(0, Math.min(100, (done / total) * 100)) : 0;
    const bonus = Number(rush.bonus) || 0;
    const rem = typeof rush.until === "number" ? Math.max(0, rush.until - (Number(tick) || 0)) : null;
    const bar = total > 0
      ? `<div class="ct-bar" role="progressbar" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${Math.round(pct)}"><i style="width:${pct.toFixed(1)}%"></i></div>`
      : "";
    return `<div class="item rush-card">`
      + `<div class="item-top"><b>${escapeFn(tr("طلب عاجل", "Rush order"))}</b><span class="rush-bonus" dir="ltr">${escapeFn(tr("مكافأة", "Bonus"))}: <b>${fmtNum(bonus)}</b></span></div>`
      + `<div class="rush-row"><span class="rush-good"><i style="background:${chipColor}"></i>${escapeFn(goodName)}</span><span class="rush-amt" dir="ltr">${escapeFn(tr("متبقي", "Left"))} <b>${fmtNum(left)}</b>${total > 0 ? ` / ${fmtNum(total)}` : ""}</span></div>`
      + bar
      + (rem !== null ? `<div class="item-foot"><span class="rush-time">${escapeFn(tr("الوقت المتبقي", "Time left"))}</span><b dir="ltr">${clockFmt(rem)}</b></div>` : "")
      + `</div>`;
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
      const isFault = !!d.fault;
      const stunnedTicks = Number(d.stunned) || 0;
      const isStunned = stunnedTicks > 0;
      const tags = [
        isFault ? `<span class="flt-tag fault">${escapeFn(tr("عطل", "Fault"))}</span>` : "",
        isStunned ? `<span class="flt-tag stunned">${escapeFn(tr("مصاب", "Stunned"))} <b dir="ltr">${clockFmt(stunnedTicks)}</b></span>` : "",
      ].filter(Boolean).join("");
      const tagsHtml = tags ? `<span class="flt-tags">${tags}</span>` : "";
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
      return `<div class="item flt-row${isFault ? " flt-fault" : ""}${isStunned ? " flt-stunned" : ""}" data-drone="${id}"><div class="item-top"><b><i class="flt-dot" style="background:${col}"></i>${escapeFn(tr(`درون ${id}`, `Drone ${id}`))}</b>${tagsHtml}</div>${batMarkup}${cargoMarkup}</div>`;
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

  const currentTick = () => Math.max(Number(data && data.state && data.state.tick) || 0, Number(typeof W !== "undefined" && W && W.colony ? W.t : 0) || 0);

  function renderWeather() {
    const box = $("#col-weather");
    if (!box) return;
    const html = weatherMarkup(data && data.state && data.state.weather, currentTick());
    box.classList.toggle("hidden", !html);
    box.innerHTML = html;
  }

  function renderLastRun() {
    const box = $("#col-lastrun");
    if (!box) return;
    const html = !running && runs.length ? lastRunMarkup(runs) : "";
    box.classList.toggle("hidden", !html);
    box.innerHTML = html;
  }

  function renderLiveStats() {
    const box = $("#col-live-stats");
    if (!box) return;
    const html = running && data && data.state ? liveStatsMarkup(data.state.stats) : "";
    box.classList.toggle("hidden", !html);
    box.innerHTML = html;
  }

  function renderRush() {
    const box = $("#col-rush");
    if (!box) return;
    const html = data && data.state ? rushMarkup(data.state.rush, currentTick()) : "";
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
    if (u.max > 0 && lv >= u.max) return { maxed: true };
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
      const hasEndless = items.some((u) => u.max === 0);
      const owned = items.reduce((a, u) => a + (u.max === 0 ? (u.level || 0) : Math.min(u.level, u.max)), 0);
      const total = items.reduce((a, u) => a + (u.max || 0), 0);
      const groupCount = hasEndless ? `${owned}` : `${owned}/${total}`;
      return `<h3>${esc(g)} <small>${groupCount}</small></h3>` + items.map((u) => {
        const s = status(u);
        const goods = Object.entries(u.cost.goods || {}).map(([k, n]) => `<span class="cg" title="${esc(NAME[k] || k).replace(/"/g, "&quot;")}"><i style="background:${INV_COLOR[k] || "#ccc"}"></i>${fmt(n)}</span>`).join("");
        const cost = s.maxed ? `<span class="maxed">${T("مكتمل", "Maxed")}</span>` : `<span class="cc">${fmt(u.cost.credits)}</span>${goods}`;
        const lesson = u.lesson && !s.maxed && !s.can ? `<span class="lesson-link">${T(`بتتعلمها في الفصل ${u.lesson}`, `Taught in chapter ${u.lesson}`)}</span>` : "";
        const lvText = u.max === 0 ? String(u.level || 0) : (u.max > 1 ? `${u.level}/${u.max}` : u.level ? "✓" : "");
        return `<div class="item${s.can ? " can" : ""}${s.maxed ? " done" : ""}">
          <div class="item-top"><b>${esc(uName(u))}</b><span class="lv">${lvText}</span></div>
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
  const raceBtn = $("#btn-race"); if (raceBtn) raceBtn.onclick = toggleRace;
  const testBtn = $("#btn-test"); if (testBtn) testBtn.onclick = () => runCheck("test");
  const vetBtn = $("#btn-vet"); if (vetBtn) vetBtn.onclick = () => runCheck("vet");

  return { enter, exit, toggleRun, LEVEL, get running() { return running; } };
})();

$("#btn-colony").onclick = () => (S.mode === "colony" ? loadLevel(store.pref("current") || S.levels[0].id) : Colony.enter());
