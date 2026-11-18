"use strict";
// Language and player profile, resolved before the rest of the app loads.
// Arabic is written in the HTML; English lives in data-en* attributes.
const PROFILE = localStorage.getItem("orbital.profile") || "default";
const LANG = localStorage.getItem("orbital.lang." + PROFILE) === "en" ? "en" : "ar";
const T = (ar, en) => (LANG === "en" ? en : ar);

(function applyLang() {
  const root = document.documentElement;
  root.lang = LANG;
  root.dir = LANG === "en" ? "ltr" : "rtl";
  if (LANG !== "en") return;
  document.title = "Orbital Colony — learn Go while you survey a planet";
  for (const el of document.querySelectorAll("[data-en]")) el.textContent = el.dataset.en;
  for (const el of document.querySelectorAll("[data-en-html]")) el.innerHTML = el.dataset.enHtml;
  for (const el of document.querySelectorAll("[data-en-title]")) el.title = el.dataset.enTitle;
  for (const el of document.querySelectorAll("[data-en-placeholder]")) el.placeholder = el.dataset.enPlaceholder;
  for (const el of document.querySelectorAll("[data-en-aria]")) el.setAttribute("aria-label", el.dataset.enAria);
})();

// Level and chapter text from the server carries an optional English copy.
function localizeLevel(l) {
  if (LANG !== "en" || !l.en) return l;
  const e = l.en;
  return { ...l, ...(e.steps?.length ? { steps: e.steps, lessonRevision: e.lessonRevision } : {}),
    title: e.title || l.title, goal: e.goal || l.goal, lesson: e.lesson || l.lesson,
    starter: e.starter || l.starter, hints: e.hints?.length ? e.hints : l.hints,
    files: e.files?.length ? e.files : l.files };
}
function localizeChapter(c) {
  return LANG === "en" ? { ...c, title: c.titleEn || c.title, concept: c.conceptEn || c.concept, track: c.trackEn || c.track } : c;
}
