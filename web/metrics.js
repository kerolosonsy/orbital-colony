"use strict";

(function (global) {
  const failureKinds = ["compile", "runtime", "check"];
  let activeLevel = null;
  let activeSince = null;

  function metricCount(countValue) { return Number.isSafeInteger(countValue) && countValue > 0 ? countValue : 0; }

  function emptyLevel() {
    return { runs: 0, failedRuns: { compile: 0, runtime: 0, check: 0 }, hintCount: 0, hintsRevealed: [],
      lessonsOpened: 0, lessonStepsViewed: Object.create(null), interactions: 0, activeTimeMs: 0,
      solved: false, timeToFirstSolveMs: 0, bestTicks: 0 };
  }

  function normalizeLevel(savedMetrics = {}) {
    const entry = emptyLevel();
    entry.runs = metricCount(savedMetrics.runs);
    for (const kind of failureKinds) entry.failedRuns[kind] = metricCount(savedMetrics.failedRuns?.[kind]);
    entry.hintsRevealed = [...new Set((Array.isArray(savedMetrics.hintsRevealed) ? savedMetrics.hintsRevealed : []).filter((hint) => typeof hint === "string" && hint))];
    entry.hintCount = Math.max(metricCount(savedMetrics.hintCount), entry.hintsRevealed.length);
    entry.lessonsOpened = metricCount(savedMetrics.lessonsOpened);
    entry.interactions = metricCount(savedMetrics.interactions);
    entry.activeTimeMs = metricCount(savedMetrics.activeTimeMs);
    entry.solved = savedMetrics.solved === true;
    entry.timeToFirstSolveMs = metricCount(savedMetrics.timeToFirstSolveMs);
    entry.bestTicks = metricCount(savedMetrics.bestTicks);
    for (const [step, views] of Object.entries(savedMetrics.lessonStepsViewed || {})) entry.lessonStepsViewed[step] = metricCount(views);
    return entry;
  }

  function normalizeMetrics(raw) {
    const metrics = Object.create(null);
    if (!raw || typeof raw !== "object" || Array.isArray(raw)) return metrics;
    for (const [levelID, entry] of Object.entries(raw)) {
      if (entry && typeof entry === "object" && !Array.isArray(entry)) metrics[levelID] = normalizeLevel(entry);
    }
    return metrics;
  }

  function mergeMetrics(first, second) {
    const merged = normalizeMetrics(first);
    for (const [levelID, next] of Object.entries(normalizeMetrics(second))) {
      const current = merged[levelID] || emptyLevel();
      current.runs = Math.max(current.runs, next.runs);
      for (const kind of failureKinds) current.failedRuns[kind] = Math.max(current.failedRuns[kind], next.failedRuns[kind]);
      current.hintCount = Math.max(current.hintCount, next.hintCount);
      current.hintsRevealed = [...new Set([...current.hintsRevealed, ...next.hintsRevealed])];
      current.hintCount = Math.max(current.hintCount, current.hintsRevealed.length);
      current.lessonsOpened = Math.max(current.lessonsOpened, next.lessonsOpened);
      current.interactions = Math.max(current.interactions, next.interactions);
      current.activeTimeMs = Math.max(current.activeTimeMs, next.activeTimeMs);
      current.solved ||= next.solved;
      current.timeToFirstSolveMs = Math.max(current.timeToFirstSolveMs, next.timeToFirstSolveMs);
      if (!current.bestTicks || next.bestTicks && next.bestTicks < current.bestTicks) current.bestTicks = next.bestTicks;
      for (const [step, views] of Object.entries(next.lessonStepsViewed)) current.lessonStepsViewed[step] = Math.max(current.lessonStepsViewed[step] || 0, views);
      merged[levelID] = current;
    }
    return merged;
  }

  function levelMetrics(store, levelID) {
    store.blob.metrics = store.blob.metrics || Object.create(null);
    return store.blob.metrics[levelID] ||= emptyLevel();
  }

  function pageIsActive() {
    return document.visibilityState !== "hidden" && (typeof document.hasFocus !== "function" || document.hasFocus());
  }

  function accrueActiveTime(store) {
    if (activeLevel === null || activeSince === null) return false;
    const entry = levelMetrics(store, activeLevel);
    const elapsed = Math.max(0, Math.floor(global.performance.now() - activeSince));
    if (!entry.solved) entry.activeTimeMs += elapsed;
    activeSince = pageIsActive() && !entry.solved ? global.performance.now() : null;
    return elapsed > 0 && !entry.solved;
  }

  function saveAfterEvent(store, levelID) {
    if (activeLevel === levelID && pageIsActive() && !levelMetrics(store, levelID).solved) activeSince = global.performance.now();
    store.save();
  }

  function recordEvent(store, levelID, update) {
    accrueActiveTime(store);
    update(levelMetrics(store, levelID));
    saveAfterEvent(store, levelID);
  }

  function pauseActiveTime(store) {
    if (accrueActiveTime(store)) store.save();
    activeSince = null;
  }

  function resumeActiveTime(store) {
    if (activeLevel !== null && pageIsActive() && !levelMetrics(store, activeLevel).solved) activeSince = global.performance.now();
  }

  function bind(store) {
    document.addEventListener("visibilitychange", () => pageIsActive() ? resumeActiveTime(store) : pauseActiveTime(store));
    global.addEventListener("blur", () => pauseActiveTime(store));
    global.addEventListener("focus", () => resumeActiveTime(store));
    global.addEventListener("pagehide", () => {
      pauseActiveTime(store);
      store.push();
    });
  }

  function openLevel(store, levelID) {
    if (accrueActiveTime(store)) store.save();
    activeLevel = levelID;
    activeSince = pageIsActive() && !levelMetrics(store, levelID).solved ? global.performance.now() : null;
  }

  function recordRun(store, levelID, response) {
    accrueActiveTime(store);
    const entry = levelMetrics(store, levelID);
    entry.runs++;
    if (!response.result?.success) {
      const kind = failureKinds.includes(response.failure) ? response.failure : response.phase === "compile" ? "compile" : "check";
      entry.failedRuns[kind]++;
    } else {
      if (!entry.solved) entry.timeToFirstSolveMs = entry.activeTimeMs;
      entry.solved = true;
      if (!entry.bestTicks || response.result.ticks && response.result.ticks < entry.bestTicks) entry.bestTicks = metricCount(response.result.ticks);
      if (activeLevel === levelID) activeSince = null;
    }
    saveAfterEvent(store, levelID);
  }

  function recordHint(store, levelID, hintIndex) {
    recordEvent(store, levelID, (entry) => {
      entry.hintCount++;
      const hint = String(hintIndex);
      if (!entry.hintsRevealed.includes(hint)) entry.hintsRevealed.push(hint);
    });
  }

  function recordLessonOpened(store, levelID) {
    recordEvent(store, levelID, (entry) => { entry.lessonsOpened++; });
  }

  function recordStepViewed(store, levelID, stepID) {
    recordEvent(store, levelID, (entry) => { entry.lessonStepsViewed[stepID] = (entry.lessonStepsViewed[stepID] || 0) + 1; });
  }

  function recordInteraction(store, levelID) {
    recordEvent(store, levelID, (entry) => { entry.interactions++; });
  }

  function struggleScore(entry) {
    const failed = failureKinds.reduce((total, kind) => total + entry.failedRuns[kind], 0);
    return failed * 4 + entry.hintCount * 2 + Math.ceil(entry.activeTimeMs / 60000);
  }

  function rankStruggles(metrics, levels) {
    const order = new Map(levels.map((level, index) => [level.id, index]));
    return Object.entries(normalizeMetrics(metrics)).map(([id, entry]) => ({ id, entry, score: struggleScore(entry) }))
      .filter((row) => row.score > 0)
      .sort((a, b) => b.score - a.score || (order.get(a.id) ?? Infinity) - (order.get(b.id) ?? Infinity) || a.id.localeCompare(b.id));
  }

  function formatActive(ms) {
    const minutes = Math.ceil(ms / 60000);
    return T(`${minutes} د`, `${minutes} min`);
  }

  function renderStruggles(root, metrics, levels) {
    const byID = new Map(levels.map((level) => [level.id, level]));
    const rows = rankStruggles(metrics, levels).slice(0, 20);
    root.replaceChildren();
    if (!rows.length) {
      const empty = document.createElement("li");
      empty.textContent = T("لسه مفيش محاولات مسجلة.", "No learning activity recorded yet.");
      root.appendChild(empty);
      return;
    }
    rows.forEach(({ id, entry, score }) => {
      const row = document.createElement("li");
      const level = byID.get(id);
      const title = document.createElement("strong");
      title.textContent = `${id} · ${level?.title || id}`;
      const detail = document.createElement("span");
      const failures = failureKinds.reduce((total, kind) => total + entry.failedRuns[kind], 0);
      detail.textContent = T(`درجة ${score} · ${failures} فشل · ${entry.hintCount} تلميح · ${formatActive(entry.activeTimeMs)}`,
        `Score ${score} · ${failures} failed · ${entry.hintCount} hints · ${formatActive(entry.activeTimeMs)}`);
      row.append(title, detail);
      root.appendChild(row);
    });
  }

  function chapterTotals(chapter, metrics, levels) {
    const total = { levels: 0, solved: 0, runs: 0, failures: 0, hints: 0, activeTimeMs: 0 };
    for (const level of levels.filter((candidate) => candidate.chapter === chapter.n)) {
      const entry = metrics[level.id];
      if (!entry) continue;
      total.levels++;
      total.solved += entry.solved ? 1 : 0;
      total.runs += entry.runs;
      total.failures += failureKinds.reduce((sum, kind) => sum + entry.failedRuns[kind], 0);
      total.hints += entry.hintCount;
      total.activeTimeMs += entry.activeTimeMs;
    }
    return total;
  }

  function renderChapters(root, metrics, levels, chapters) {
    const normalized = normalizeMetrics(metrics);
    root.replaceChildren();
    const summaries = chapters.map((chapter) => ({ chapter, totals: chapterTotals(chapter, normalized, levels) }))
      .filter(({ totals }) => totals.levels > 0);
    if (!summaries.length) {
      const empty = document.createElement("p");
      empty.textContent = T("الملخص هيظهر مع اللعب.", "Chapter summaries will appear as you play.");
      root.appendChild(empty);
      return;
    }
    for (const { chapter, totals } of summaries) {
      const card = document.createElement("article");
      const title = document.createElement("strong");
      title.textContent = `${chapter.n} · ${chapter.title}`;
      const summary = document.createElement("p");
      summary.textContent = T(`${totals.solved} محلول · ${totals.runs} محاولة · ${totals.failures} فشل · ${totals.hints} تلميح · ${formatActive(totals.activeTimeMs)}`,
        `${totals.solved} solved · ${totals.runs} runs · ${totals.failures} failed · ${totals.hints} hints · ${formatActive(totals.activeTimeMs)}`);
      card.append(title, summary);
      root.appendChild(card);
    }
  }

  function renderStats(roots, metrics, catalog) {
    renderStruggles(roots.struggles, metrics, catalog.levels);
    renderChapters(roots.chapters, metrics, catalog.levels, catalog.chapters);
  }

  global.LearningMetrics = { normalizeMetrics, mergeMetrics, openLevel, bind, recordRun, recordHint,
    recordLessonOpened, recordStepViewed, recordInteraction, rankStruggles, renderStats };
})(globalThis);
