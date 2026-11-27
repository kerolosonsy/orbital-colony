"use strict";

// Grade: the code score out of 100 under the mission's stars. After a mission
// passes, /api/grade asks the helper's model to judge the craft of the code;
// the player can open why it scored that and how to reach 100, in words only.
(function () {
  const escape = (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");

  // codeHash names one version of the files, so the same code is graded once.
  function codeHash(files) {
    const text = Object.keys(files).sort().map((name) => `${name}\n${files[name]}`).join("\n\u0000");
    let h = 2166136261;
    for (let i = 0; i < text.length; i++) h = Math.imul(h ^ text.charCodeAt(i), 16777619);
    return (h >>> 0).toString(36) + ":" + text.length;
  }

  // point renders one reason or step: escaped, with `identifiers` as code.
  function point(text) {
    return escape(text).replace(/`([^`]+)`/g, '<code dir="ltr">$1</code>').replace(/\*\*([^*]+)\*\*/g, "<b>$1</b>");
  }

  function tone(score) { return score >= 90 ? "high" : score >= 70 ? "mid" : "low"; }

  // view returns the panel's HTML for a state:
  //   {kind: "hidden"} · {kind: "ask"} · {kind: "busy"} · {kind: "error", setup, limit}
  //   {kind: "done", grade: {score, summary, why, how}, stale, open: "why"|"how"|null}
  function view(state, T) {
    const label = `<span class="grade-label">${escape(T("تقييم الكود", "Code score"))}</span>`;
    if (state.kind === "ask") return `<div class="grade-head">${label}<button type="button" class="tool small" data-grade="again">${escape(T("قيّم كودي", "Score my code"))}</button></div>`;
    if (state.kind === "busy") return `<div class="grade-head">${label}<span class="grade-busy" role="status">${escape(T("بيقرا كودك…", "Reading your code…"))}</span></div>`;
    if (state.kind === "error") {
      const why = state.limit
        ? T("خلّصت تقييمات النهارده. بكرة تقدر تقيّم تاني.", "You've used today's code scores. You can score again tomorrow.")
        : state.setup
        ? T("التقييم محتاج Codex CLI متسطب ومتسجّل بحساب ChatGPT (brew install codex ثم codex login).", "Scoring needs the Codex CLI signed in with ChatGPT (brew install codex, then codex login).")
        : T("التقييم مجاش المرة دي.", "The score didn't come through this time.");
      return `<div class="grade-head">${label}<button type="button" class="tool small" data-grade="again">${escape(T("جرّب تاني", "Try again"))}</button></div><p class="grade-note">${escape(why)}</p>`;
    }
    if (state.kind !== "done") return "";
    const g = state.grade, open = state.open;
    const list = (name, items) => `<ul class="grade-list grade-${name}"${open === name ? "" : " hidden"}>${items.map((t) => `<li>${point(t)}</li>`).join("")}</ul>`;
    const toggle = (name, text) => `<button type="button" class="tool small" data-grade="${name}" aria-expanded="${open === name}">${escape(text)}</button>`;
    return `<div class="grade-head">${label}
        <b class="grade-score" data-tone="${tone(g.score)}">${g.score}<small>/100</small></b>
        <span class="grade-meter" aria-hidden="true"><i style="width:${g.score}%"></i></span></div>
      ${g.summary ? `<p class="grade-summary">${point(g.summary)}</p>` : ""}
      <div class="grade-actions">
        ${g.why?.length ? toggle("why", T(`ليه ${g.score}؟`, `Why ${g.score}?`)) : ""}
        ${g.how?.length ? toggle("how", T("إزاي أوصل لـ100؟", "How do I reach 100?")) : ""}
      </div>
      ${g.why?.length ? list("why", g.why) : ""}${g.how?.length ? list("how", g.how) : ""}
      ${state.stale ? `<p class="grade-note">${escape(T("الكود اتغيّر من آخر تقييم.", "The code changed since this score."))} <button type="button" class="tool small" data-grade="again">${escape(T("قيّمه تاني", "Score it again"))}</button></p>` : ""}`;
  }

  // ---- the panel (browser only) ----
  // host: { T, api(path, body), level() → current level, files(level) → {name: code},
  //         record(level) → progress record, save(), canGrade(level) }
  function mount(host) {
    const box = document.getElementById("lvl-grade");
    let open = null, busy = null, failed = null;

    function state() {
      const l = host.level();
      if (!l || !host.canGrade(l)) return { kind: "hidden" };
      const r = host.record(l);
      if (busy === l.id) return { kind: "busy" };
      if (failed?.id === l.id) return { kind: "error", setup: failed.setup, limit: failed.limit };
      if (!r.grade) return r.stars > 0 ? { kind: "ask" } : { kind: "hidden" };
      return { kind: "done", grade: r.grade, open, stale: r.grade.hash !== codeHash(host.files(l)) };
    }
    function draw() {
      const s = state();
      box.hidden = s.kind === "hidden";
      box.innerHTML = view(s, host.T);
    }

    // grade asks for a score unless this exact code already has one.
    async function grade(l, run = {}, force = false) {
      if (!l || !host.canGrade(l) || busy) return;
      const files = host.files(l), hash = codeHash(files), r = host.record(l);
      if (!force && r.grade?.hash === hash) return draw();
      busy = l.id; failed = null; open = null; draw();
      let result;
      try { result = await host.api("/api/grade", { level: l.id, files, stars: run.stars || r.stars, ticks: run.ticks || r.best || 0 }); }
      catch { result = { error: "failed" }; }
      busy = null;
      if (typeof result?.score === "number") {
        r.grade = { score: result.score, summary: result.summary || "", why: result.why || [], how: result.how || [], hash, at: Date.now() };
        host.save();
      } else failed = { id: l.id, setup: !!result?.setup, limit: result?.error === "limit" };
      if (host.level() === l) draw();
    }

    box.addEventListener("click", (e) => {
      const b = e.target.closest("button[data-grade]");
      if (!b) return;
      const action = b.dataset.grade;
      if (action === "again") return grade(host.level(), {}, true);
      open = open === action ? null : action;
      draw();
    });
    return { draw, grade, reset() { open = null; failed = null; draw(); } };
  }

  const api = { codeHash, point, tone, view, mount };
  if (typeof window !== "undefined") window.Grade = api;
  if (typeof module !== "undefined") module.exports = api;
})();
