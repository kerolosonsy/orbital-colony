"use strict";

const TOKEN = document.querySelector('meta[name="orbital-token"]').content;
const $ = (s) => document.querySelector(s);

const HEADERS = { "X-Orbital-Token": TOKEN, "X-Orbital-Lang": LANG, "X-Orbital-Profile": PROFILE, "Content-Type": "application/json" };
async function api(path, body, method, signal) {
  const res = await fetch(path, {
    method: method || (body ? "POST" : "GET"),
    headers: HEADERS,
    body: body ? JSON.stringify(body) : undefined,
    signal,
  });
  if (!res.ok) throw new Error(`${res.status} ${await res.text()}`);
  return res.json();
}

// ---------- persistence ----------
// Each player profile keeps mission, quiz and lesson progress in separate fields.
// localStorage holds a copy so a failed upload is never lost.
const store = {
  blob: { progress: {}, prefs: {}, quizScores: {}, lessons: {}, metrics: {} }, timer: null,
  cacheKey: "orbital.p." + PROFILE,
  async load() {
    let b = null, online = true;
    try { b = await api(`/api/profiles/${PROFILE}/progress`); } catch { online = false; }
    // the local copy only stands in while the server can't be reached; an
    // empty answer from the server means this player really starts fresh
    if (!online) { try { b = JSON.parse(localStorage.getItem(this.cacheKey)); } catch { b = null; } }
    const legacy = !b?.progress && PROFILE === "default" && online && this.legacy();
    if (legacy) b = legacy;
    const lessons = b?.lessons && typeof b.lessons === "object" && !Array.isArray(b.lessons) ? b.lessons : {};
    this.blob = { progress: b?.progress || {}, prefs: b?.prefs || {}, quizScores: b?.quizScores || {}, review: b?.review || {}, lessons,
      metrics: LearningMetrics.normalizeMetrics(b?.metrics) };
    if (legacy) {
      await this.push();
      for (const k of ["orbital.progress.v1", "orbital.current", "orbital.welcomed", "orbital.linked"]) localStorage.removeItem(k);
    }
    return this.blob.progress;
  },
  // progress saved before profiles existed belongs to the first player
  legacy() {
    try {
      const progress = JSON.parse(localStorage.getItem("orbital.progress.v1"));
      if (!progress) return null;
      return { progress, prefs: { current: localStorage.getItem("orbital.current"), welcomed: !!localStorage.getItem("orbital.welcomed"),
        linked: JSON.parse(localStorage.getItem("orbital.linked") || "[]") } };
    } catch { return null; }
  },
  save() {
    // a full or blocked localStorage must not stop saving to the server
    try { localStorage.setItem(this.cacheKey, JSON.stringify(this.blob)); } catch {}
    clearTimeout(this.timer);
    this.timer = setTimeout(() => this.push(), 800);
  },
  push() {
    clearTimeout(this.timer); this.timer = null;
    return fetch(`/api/profiles/${PROFILE}/progress`, { method: "PUT", headers: HEADERS, body: JSON.stringify(this.blob), keepalive: true }).then(async (response) => {
      if (!response.ok) throw new Error(String(response.status));
      const saved = await response.json();
      this.blob.review = saved.review || {};
      this.blob.metrics = LearningMetrics.mergeMetrics(this.blob.metrics, saved.metrics);
      if ($("#missions").open) renderMissions();
      this.warned = false;
    }).catch((error) => {
      // progress stays in the local copy; say once why it isn't reaching the server
      if (this.warned || typeof logServerHint !== "function") return;
      this.warned = true;
      logServerHint(error.message === "403" ? new Error("403 token") : error);
    });
  },
  pref(k, v) {
    if (v === undefined) return this.blob.prefs[k];
    this.blob.prefs[k] = v; this.save();
  },
};
LearningMetrics.bind(store);
const S = { chapters: [], levels: [], cur: null, progress: {}, hintsShown: 0, running: false };
const rec = (id) => (S.progress[id] ||= { stars: 0, best: 0, code: null });

// ---------- text helpers ----------
const esc = (s) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");

const KW = new Set("break case chan const continue default defer else fallthrough for func go goto if import interface map package range return select struct switch type var".split(" "));
const TYPES = new Set("int int8 int16 int32 int64 uint uint8 uint16 uint32 uint64 float32 float64 string bool byte rune error any nil true false iota".split(" "));

function highlight(code) {
  const re = /(\/\/[^\n]*|\/\*[\s\S]*?\*\/)|("(?:\\.|[^"\\\n])*"|`[^`]*`|'(?:\\.|[^'\\\n])*')|\b(\d+(?:\.\d+)?)\b|\b([A-Za-z_]\w*)\b(\s*\()?/g;
  let out = "", last = 0, m;
  while ((m = re.exec(code))) {
    out += esc(code.slice(last, m.index));
    last = re.lastIndex;
    if (m[1]) out += `<span class="tk-com">${esc(m[1])}</span>`;
    else if (m[2]) out += `<span class="tk-str">${esc(m[2])}</span>`;
    else if (m[3]) out += `<span class="tk-num">${m[3]}</span>`;
    else {
      const w = m[4], paren = m[5] || "";
      const prev = code[m.index - 1];
      const next = code[m.index + w.length];
      let cls = "";
      if (KW.has(w)) cls = "tk-kw";
      else if (TYPES.has(w)) cls = "tk-type";
      else if (paren) cls = "tk-fn";
      else if (next === "." && prev !== ".") cls = "tk-pkg";
      out += cls ? `<span class="${cls}">${w}</span>${esc(paren)}` : w + esc(paren);
    }
  }
  return out + esc(code.slice(last));
}

function inline(s) {
  return esc(s)
    .replace(/~([^~]+)~/g, (_, c) => `<code>${c.replace(/\\\|/g, "|")}</code>`)
    .replace(/`([^`]+)`/g, "<code>$1</code>")
    .replace(/\*\*([^*]+)\*\*/g, "<b>$1</b>");
}

function md(src) {
  src = src.replace(/\s?\p{Extended_Pictographic}\uFE0F?/gu, "");
  const lines = src.split("\n");
  const out = [];
  for (let i = 0; i < lines.length; ) {
    const l = lines[i];
    if (l.startsWith("~~~")) {
      const buf = [];
      for (i++; i < lines.length && !lines[i].startsWith("~~~"); i++) buf.push(lines[i]);
      i++;
      out.push(`<pre><code>${highlight(buf.join("\n"))}</code></pre>`);
    } else if (l.startsWith("## ")) { out.push(`<h3>${inline(l.slice(3))}</h3>`); i++; }
    else if (l.startsWith("# ")) { out.push(`<h2>${inline(l.slice(2))}</h2>`); i++; }
    else if (/^\s*- /.test(l)) {
      const items = [];
      for (; i < lines.length && /^\s*- /.test(lines[i]); i++) items.push(`<li>${inline(lines[i].replace(/^\s*- /, ""))}</li>`);
      out.push(`<ul>${items.join("")}</ul>`);
    } else if (/^\d+\. /.test(l)) {
      const items = [];
      for (; i < lines.length && /^\d+\. /.test(lines[i]); i++) items.push(`<li>${inline(lines[i].replace(/^\d+\. /, ""))}</li>`);
      out.push(`<ol>${items.join("")}</ol>`);
    } else if (l.startsWith("|")) {
      const rows = [];
      for (; i < lines.length && lines[i].startsWith("|"); i++) {
        if (/^\|[\s|:-]+\|$/.test(lines[i])) continue;
        const cells = lines[i].slice(1, -1).split(/(?<!\\)\|/);
        rows.push(`<tr>${cells.map((c) => `<td>${inline(c.trim())}</td>`).join("")}</tr>`);
      }
      out.push(`<table>${rows.join("")}</table>`);
    } else if (!l.trim()) i++;
    else {
      const buf = [];
      for (; i < lines.length && lines[i].trim() && !/^(~~~|#|\s*- |\d+\. |\|)/.test(lines[i]); i++) buf.push(inline(lines[i]));
      out.push(`<p>${buf.join("<br>")}</p>`);
    }
  }
  return out.join("");
}

// ---------- files ----------
// Missions have one file (main.go); lab levels and colony have several tabs.
function levelFiles(l) {
  if (l.kind === "colony") {
    const r = rec(l.id);
    const extras = Object.keys(r.files || {}).filter((name) => name !== "main.go");
    return ["main.go", ...extras].map((name) => ({
      name,
      starter: name === "main.go" ? l.starter : "",
      readOnly: false,
    }));
  }
  return l.files?.length ? l.files.map((f) => ({ ...f, readOnly: f.readOnly || !f.name.endsWith(".go") })) : [{ name: "main.go", starter: l.starter, readOnly: false }];
}
function fileCode(l, f) {
  const r = S.review && S.cur === l ? S.reviewDraft : rec(l.id);
  if (l.kind === "colony") {
    if (f.name === "main.go") return r.files?.["main.go"] ?? r.code ?? f.starter;
    return r.files?.[f.name] ?? f.starter ?? "";
  }
  if (!l.files?.length) return r.code ?? f.starter;
  if (f.readOnly) return f.starter;
  const carried = l.carryFrom ? rec(l.carryFrom).files?.[f.name] : undefined;
  return r.files?.[f.name] ?? carried ?? f.starter;
}
function migrateColonyRecord(r, starter = "") {
  if (!r.files) {
    r.files = { "main.go": r.code ?? starter };
  } else if (r.files["main.go"] === undefined) {
    r.files["main.go"] = r.code ?? starter;
  }
  return r.files;
}
function saveFile(l, name, code) {
  const r = S.review && S.cur === l ? S.reviewDraft : rec(l.id);
  if (l.kind === "colony") {
    const files = migrateColonyRecord(r, l.starter);
    files[name] = code;
    if (name === "main.go") r.code = code;
    return;
  }
  if (!l.files?.length) r.code = code;
  else (r.files ||= {})[name] = code;
}
let lastErrors = [];
function showFile(name) {
  flushSave();
  const l = S.cur;
  const f = levelFiles(l).find((x) => x.name === name) || levelFiles(l)[0];
  S.file = f.name;
  ed.value = fileCode(l, f);
  ed.readOnly = !!f.readOnly || isLinked();
  errLines = new Map(lastErrors.filter((e) => (e.file || "main.go").replace(/^\.\//, "") === f.name).map((e) => [e.line, e.msg]));
  renderTabs();
  renderEditor();
}
function renderTabs() {
  const box = $("#file-tabs");
  box.innerHTML = "";
  const isColony = S.cur?.kind === "colony";
  const linkedNow = isLinked();
  for (const f of levelFiles(S.cur)) {
    if (!isColony || f.name === "main.go") {
      const b = document.createElement("button");
      b.type = "button";
      b.className = "tab" + (f.name === S.file ? " active" : "");
      b.innerHTML = esc(f.name) + (f.readOnly ? '<span class="ro">' + T("قراية بس", "read-only") + "</span>" : "");
      b.onclick = () => showFile(f.name);
      box.appendChild(b);
      continue;
    }
    const wrap = document.createElement("div");
    wrap.className = "tab" + (f.name === S.file ? " active" : "");
    const main = document.createElement("button");
    main.type = "button";
    main.className = "tab-main";
    main.textContent = f.name;
    main.onclick = () => showFile(f.name);
    const ren = document.createElement("button");
    ren.type = "button";
    ren.className = "tab-act";
    ren.textContent = "✎";
    ren.title = T("إعادة تسمية", "Rename");
    ren.setAttribute("aria-label", T(`غيّر اسم ${f.name}`, `Rename ${f.name}`));
    ren.disabled = linkedNow;
    ren.onclick = () => promptRenameColonyFile(f.name);
    const del = document.createElement("button");
    del.type = "button";
    del.className = "tab-act tab-del";
    del.textContent = "×";
    del.title = T("حذف", "Delete");
    del.setAttribute("aria-label", T(`امسح ${f.name}`, `Delete ${f.name}`));
    del.disabled = linkedNow;
    del.onclick = () => promptDeleteColonyFile(f.name);
    wrap.append(main, ren, del);
    box.appendChild(wrap);
  }
  if (isColony) {
    const add = document.createElement("button");
    add.type = "button";
    add.className = "tab tab-add";
    add.textContent = "+";
    add.title = T("ملف جديد", "New file");
    add.setAttribute("aria-label", T("ملف جديد", "New file"));
    add.disabled = linkedNow;
    add.onclick = () => promptCreateColonyFile();
    box.appendChild(add);
  }
}

const COLONY_EXTS = [".go", ".sql", ".json", ".yaml", ".txt", ".md", ".html"];
function validateColonyPath(raw, existing = []) {
  const path = String(raw ?? "").trim();
  if (!path) return { ok: false, error: T("اكتب مسار الملف.", "Enter a file path.") };
  if (path.startsWith("/") || path.includes("\\") || /^[a-zA-Z]:/.test(path)) {
    return { ok: false, error: T("لازم المسار يكون نسبي (من غير / في الأول أو \\).", "Path must be relative (no leading / or \\).") };
  }
  if (path.includes("..")) {
    return { ok: false, error: T("المسار مينفعش يحتوي على '..'.", "Path cannot contain '..'.") };
  }
  if (/[\x00-\x1f<>:"|?*]/.test(path)) {
    return { ok: false, error: T("المسار فيه رموز مش مسموحة.", "Path contains invalid characters.") };
  }
  const parts = path.split("/");
  if (parts.some((p) => !p || p === "." || p === ".." || p.startsWith(".") || p.trim() !== p)) {
    return { ok: false, error: T("مسار الملف مش سليم.", "Invalid file path.") };
  }
  const last = parts[parts.length - 1];
  const ext = COLONY_EXTS.find((e) => last.endsWith(e));
  if (!ext || last.length <= ext.length) {
    return { ok: false, error: T("الامتدادات المسموحة: " + COLONY_EXTS.join(" "), "Allowed extensions: " + COLONY_EXTS.join(" ")) };
  }
  if (existing.includes(path)) {
    return { ok: false, error: T("فيه ملف بالاسم ده بالفعل.", "A file with that path already exists.") };
  }
  if (existing.some((e) => e.startsWith(path + "/") || path.startsWith(e + "/"))) {
    return { ok: false, error: T("المسار ده بيتعارض مع مجلد أو ملف موجود.", "Path conflicts with an existing file or folder.") };
  }
  return { ok: true, path };
}
function colonyStarterFor(path) {
  if (!path.endsWith(".go")) return "";
  const parts = path.split("/");
  if (parts.length > 1 && /^[A-Za-z_][A-Za-z0-9_]*$/.test(parts[parts.length - 2])) {
    return `package ${parts[parts.length - 2]}\n\n`;
  }
  return "package main\n\n";
}
function colonyFilesPayload(l) {
  const files = {};
  for (const f of levelFiles(l)) files[f.name] = fileCode(l, f);
  return files;
}
function createColonyFile(l, rawPath) {
  const existing = levelFiles(l).map((f) => f.name);
  const res = validateColonyPath(rawPath, existing);
  if (!res.ok) return res;
  const files = migrateColonyRecord(rec(l.id), l.starter);
  files[res.path] = colonyStarterFor(res.path);
  return { ok: true, path: res.path };
}
function renameColonyFile(l, oldName, rawNewPath) {
  if (!oldName || oldName === "main.go") {
    return { ok: false, error: T("مينفعش تغيّر اسم main.go.", "Cannot rename main.go.") };
  }
  const trimmed = String(rawNewPath ?? "").trim();
  if (trimmed === oldName) return { ok: true, path: oldName, unchanged: true };
  const existing = levelFiles(l).map((f) => f.name).filter((n) => n !== oldName);
  const res = validateColonyPath(trimmed, existing);
  if (!res.ok) return res;
  const r = rec(l.id);
  const files = migrateColonyRecord(r, l.starter);
  if (files[oldName] === undefined) return { ok: false, error: T("الملف مش موجود.", "File not found.") };
  const next = {};
  for (const [k, v] of Object.entries(files)) next[k === oldName ? res.path : k] = v;
  r.files = next;
  return { ok: true, path: res.path };
}
function deleteColonyFile(l, name) {
  if (!name || name === "main.go") {
    return { ok: false, error: T("مينفعش تمسح main.go.", "Cannot delete main.go.") };
  }
  const r = rec(l.id);
  const files = migrateColonyRecord(r, l.starter);
  delete files[name];
  return { ok: true };
}
function firstErrorFile(errors, files, fallback = "main.go") {
  if (!Array.isArray(errors) || !errors.length) return fallback;
  const raw = (errors[0].file || "main.go").replace(/^\.\//, "");
  return files.some((f) => f.name === raw) ? raw : fallback;
}
function promptCreateColonyFile() {
  if (isLinked() || S.cur?.kind !== "colony") return;
  if (!ed.readOnly) saveFile(S.cur, S.file, ed.value);
  const raw = prompt(T("مسار الملف الجديد (مثال: plan/plan.go):", "New file path (e.g. plan/plan.go):"), "");
  if (raw === null) return;
  const res = createColonyFile(S.cur, raw);
  if (!res.ok) { alert(res.error); return; }
  store.save();
  showFile(res.path);
  $("#save-state").textContent = T("محفوظ", "Saved");
  ed.focus();
}
function promptRenameColonyFile(oldName) {
  if (isLinked() || S.cur?.kind !== "colony" || oldName === "main.go") return;
  if (!ed.readOnly) saveFile(S.cur, S.file, ed.value);
  const raw = prompt(T("المسار الجديد للملف:", "New file path:"), oldName);
  if (raw === null) return;
  const res = renameColonyFile(S.cur, oldName, raw);
  if (!res.ok) { alert(res.error); return; }
  if (res.unchanged) return;
  lastErrors = lastErrors.filter((e) => (e.file || "main.go").replace(/^\.\//, "") !== oldName);
  store.save();
  showFile(S.file === oldName ? res.path : S.file);
  $("#save-state").textContent = T("محفوظ", "Saved");
}
function promptDeleteColonyFile(name) {
  if (isLinked() || S.cur?.kind !== "colony" || name === "main.go") return;
  if (!confirm(T(`تمسح الملف ${name}؟`, `Delete ${name}?`))) return;
  const res = deleteColonyFile(S.cur, name);
  if (!res.ok) return;
  lastErrors = lastErrors.filter((e) => (e.file || "main.go").replace(/^\.\//, "") !== name);
  store.save();
  showFile(S.file === name ? "main.go" : S.file);
  $("#save-state").textContent = T("محفوظ", "Saved");
}

// ---------- VS Code link ----------
// A linked mission lives in workspace/<level>/ on disk. The player edits it
// in VS Code (gopls gives completion and docs); this page polls the files.
const linked = new Set();
let lastMtime = 0;
function isLinked() { return !S.review && !!S.cur && linked.has(S.cur.id); }
function saveLinked() { store.pref("linked", [...linked]); }
function linkUI() {
  const on = isLinked();
  const b = $("#btn-vscode");
  b.classList.toggle("linked", on);
  b.querySelector(".tool-label").textContent = on ? T("مربوط بـ VS Code", "Linked to VS Code") : "VS Code";
  b.title = on ? T("الكود بيتكتب في VS Code، واللعبة بتشوف كل save", "Your code lives in VS Code; the game sees every save") : T("اكتب في VS Code بالـ autocomplete، واللعبة هتشوف كل save", "Write in VS Code with autocomplete; the game sees every save");
  $("#link-box").classList.toggle("hidden", !on);
  $("#btn-reset").disabled = on;
  if (S.cur?.kind === "colony") renderTabs();
  if (on) $("#save-state").textContent = T("اكتب في VS Code · المحرر هنا للقراية", "Write in VS Code · this editor is read-only");
}

async function openVSCode() {
  flushSave();
  const files = {};
  for (const f of levelFiles(S.cur)) if (!f.readOnly) files[f.name] = fileCode(S.cur, f);
  try {
    const r = await api("/api/workspace/open", { level: S.cur.id, files });
    linked.add(S.cur.id); saveLinked();
    log(T(`الملف اتفتح في VS Code: ${r.path}`, `Opened in VS Code: ${r.path}`), "c-info");
    if (!r.launched) log(T(`مقدرتش أفتح VS Code لوحده. افتح الفولدر ده فيه: ${r.root}`, `Couldn't launch VS Code. Open this folder in it: ${r.root}`), "c-warn");
    log(T("اكتب هناك وخلّي بالك من الـ autocomplete (Ctrl+Space). كل ما تعمل Save (⌘S) اللعبة هتشوفه.", "Write there and use autocomplete (Ctrl+Space). Every save (⌘S) shows up here."), "c-tip");
    lastMtime = 0;
    await syncFromWorkspace(true);
    showFile(S.file);
    linkUI();
  } catch (e) {
    log(T("مقدرتش أجهّز الفولدر: ", "Couldn't prepare the folder: ") + e.message, "c-err");
  }
}

async function syncFromWorkspace(force = false) {
  if (!isLinked()) return false;
  const level = S.cur;
  let r;
  try {
    r = await api(`/api/workspace/files?level=${encodeURIComponent(level.id)}`);
  } catch {
    return false; // folder deleted or server restarting
  }
  if (S.cur !== level || (!force && r.mtime === lastMtime)) return false;
  const first = lastMtime === 0;
  lastMtime = r.mtime;
  let changed = false;
  for (const [name, code] of Object.entries(r.files)) {
    const f = levelFiles(level).find((x) => x.name === name);
    if (f && !f.readOnly && fileCode(level, f) !== code) {
      saveFile(level, name, code);
      changed = true;
    }
  }
  if (changed) {
    clearTraceRun();
    store.save();
    const pos = ed.scrollTop;
    showFile(S.file);
    ed.scrollTop = pos; syncScroll();
  }
  if (!force && !first) {
    log(T("اتحفظ في VS Code", "Saved in VS Code"), "c-info");
    if ($("#autorun").checked && !S.running) run();
  }
  return changed;
}
setInterval(() => syncFromWorkspace(), 1000);

// ---------- editor ----------
const ed = $("#code"), hl = $("#highlight"), gutter = $("#gutter");
let errLines = new Map();

function renderEditor() {
  const debugState = InlineDebugger.editorState(S.file);
  hl.innerHTML = InlineDebugger.editorMarkup(ed.value, debugState, highlight);
  const n = ed.value.split("\n").length;
  gutter.innerHTML = InlineDebugger.gutterMarkup(n, errLines, InlineDebugger.breakpointLines(S.cur?.id, S.file), debugState.active ? debugState.line : 0, debugState.visited || new Set());
  syncScroll();
  if (debugState.active) {
    const lineHeight = Number.parseFloat(getComputedStyle(ed).lineHeight) || 23;
    const lineTop = 14 + (debugState.line - 1) * lineHeight;
    if (lineTop < ed.scrollTop || lineTop + lineHeight > ed.scrollTop + ed.clientHeight) {
      ed.scrollTop = Math.max(0, lineTop - ed.clientHeight / 2);
      syncScroll();
    }
  }
}
function syncScroll() {
  hl.scrollTop = ed.scrollTop; hl.scrollLeft = ed.scrollLeft; gutter.scrollTop = ed.scrollTop;
  placeCaretLine();
}
const caretLine = $("#caret-line");
function placeCaretLine() {
  const show = !ed.readOnly && ed.selectionStart === ed.selectionEnd;
  caretLine.classList.toggle("on", show);
  // the caret line's number lights up too
  const caretRow = ed.readOnly ? 0 : ed.value.slice(0, ed.selectionStart).split("\n").length;
  for (const old of gutter.querySelectorAll(".gutter-line.caret-row")) if (Number(old.dataset.debugLine) !== caretRow) old.classList.remove("caret-row");
  gutter.querySelector(`[data-debug-line="${caretRow}"]`)?.classList.add("caret-row");
  if (!show) { bracketMarks.forEach((mark) => mark.classList.remove("on")); return; }
  const line = ed.value.slice(0, ed.selectionStart).split("\n").length;
  const lineHeight = Number.parseFloat(getComputedStyle(ed).lineHeight) || 23;
  caretLine.style.top = `${14 + (line - 1) * lineHeight - ed.scrollTop}px`;
  placeBrackets(lineHeight);
}
// The bracket at the caret and its partner get a small box.
const bracketMarks = [$("#bracket-a"), $("#bracket-b")];
let charWidth = 0;
document.fonts?.addEventListener?.("loadingdone", () => { charWidth = 0; }); // re-measure once the code font arrives
function placeBrackets(lineHeight) {
  const pair = ed.selectionStart === ed.selectionEnd && window.EditKeys?.matchBracket(ed.value, ed.selectionStart);
  bracketMarks.forEach((mark, i) => {
    mark.classList.toggle("on", !!pair);
    if (!pair) return;
    if (!charWidth) {
      const ctx = document.createElement("canvas").getContext("2d");
      ctx.font = getComputedStyle(ed).font;
      charWidth = ctx.measureText("0").width || 8.4;
    }
    const at = pair[i];
    const startOfLine = ed.value.lastIndexOf("\n", at - 1) + 1;
    const row = ed.value.slice(0, at).split("\n").length - 1;
    let col = 0;
    for (const c of ed.value.slice(startOfLine, at)) col = c === "\t" ? col + 4 - (col % 4) : col + 1;
    mark.style.top = `${14 + row * lineHeight - ed.scrollTop}px`;
    mark.style.left = `${16 + col * charWidth - ed.scrollLeft}px`;
    mark.style.width = `${charWidth}px`;
  });
}
document.addEventListener("selectionchange", () => { if (document.activeElement === ed) placeCaretLine(); });
let saveTimer;
function flushSave() {
  if (!saveTimer) return;
  clearTimeout(saveTimer);
  saveTimer = null;
  store.save();
  $("#save-state").textContent = T("محفوظ", "Saved");
}
function onEdit() {
  if (ed.readOnly) return;
  clearTraceRun();
  saveFile(S.cur, S.file, ed.value);
  errLines.clear();
  renderEditor();
  $("#save-state").textContent = "…";
  clearTimeout(saveTimer);
  saveTimer = setTimeout(flushSave, 400);
}
window.addEventListener("pagehide", () => { flushSave(); if (store.timer) store.push(); });
ed.addEventListener("input", onEdit);
ed.addEventListener("scroll", syncScroll);
gutter.addEventListener("click", (event) => {
  const line = Number(event.target.closest("[data-debug-line]")?.dataset.debugLine);
  if (!line || !S.cur) return;
  InlineDebugger.toggleBreakpoint(S.cur.id, S.file, line);
  renderEditor();
});
InlineDebugger.bind({
  requestTrace: (payload, signal) => api("/api/trace", payload, null, signal),
  getFile: () => S.file,
  showFile,
  renderEditor,
  getBreakpoints: (level) => store.pref("debugBreakpoints")?.[level] || [],
  saveBreakpoints: (level, breakpoints) => {
    const saved = { ...(store.pref("debugBreakpoints") || {}) };
    saved[level] = breakpoints;
    store.pref("debugBreakpoints", saved);
  },
});
// Apply one EditKeys edit through execCommand so Ctrl/⌘+Z can undo it.
function applyEdit(edit) {
  ed.setSelectionRange(edit.from, edit.to);
  if (edit.text) {
    if (!document.execCommand("insertText", false, edit.text)) ed.setRangeText(edit.text, edit.from, edit.to, "end");
  } else if (edit.to > edit.from) {
    if (!document.execCommand("delete")) ed.setRangeText("", edit.from, edit.to, "end");
  }
  ed.setSelectionRange(edit.selStart, edit.selEnd);
  if (edit.text || edit.to > edit.from) onEdit();
}
ed.addEventListener("keydown", (e) => {
  if ((e.metaKey || e.ctrlKey) && e.key === "Enter") { e.preventDefault(); run({ paused: e.shiftKey }); return; }
  if (e.key === "F9" && S.cur) {
    // toggle a breakpoint on the caret's line, like VS Code
    e.preventDefault();
    const line = ed.value.slice(0, ed.selectionStart).split("\n").length;
    InlineDebugger.toggleBreakpoint(S.cur.id, S.file, line);
    renderEditor();
    return;
  }
  if (ed.readOnly || !window.EditKeys) return;
  const edit = EditKeys.handle(ed.value, ed.selectionStart, ed.selectionEnd, e);
  if (!edit) return;
  e.preventDefault();
  applyEdit(edit);
});

// ---------- console ----------
const con = $("#console");
function appendConsoleNode(node) {
  con.appendChild(node);
  if (con.childElementCount > 400) con.firstChild.remove();
  con.scrollTop = con.scrollHeight;
}
// Copy the console as plain text (compile errors, output) to paste anywhere.
$("#btn-copy-console").onclick = async () => {
  const text = [...con.children].map((node) => node.innerText ?? node.textContent).join("\n").trim();
  if (!text) return;
  const button = $("#btn-copy-console");
  try {
    await navigator.clipboard.writeText(text);
    button.classList.add("copied");
    button.title = T("اتنسخ ✓", "Copied ✓");
  } catch {
    button.title = T("مقدرتش أنسخ", "Couldn't copy");
  }
  setTimeout(() => {
    button.classList.remove("copied");
    button.title = T("انسخ رسايل الكونسول", "Copy console output");
  }, 1500);
};
function log(text, cls = "c-out") {
  text = String(text).replace(/\p{Extended_Pictographic}\uFE0F?\s?/gu, "");
  const d = document.createElement("div");
  d.className = cls; d.textContent = text;
  appendConsoleNode(d);
}
// Two common local situations get a plain next step: the server is not
// running, or it was restarted (new token) while this tab stayed open.
function serverHint(error) {
  const message = String(error?.message || error);
  if (error instanceof TypeError || /Failed to fetch|NetworkError|Load failed/i.test(message)) {
    return T("السيرفر شكله مقفول. شغّله بـ make start وبعدين اعمل refresh للصفحة.", "The server seems to be down. Start it with make start, then refresh the page.");
  }
  if (/^403\b/.test(message) && /token|رمز الدخول/i.test(message)) {
    return T("السيرفر اتعمله restart والصفحة دي قديمة. اعمل refresh (⌘R / F5).", "The server was restarted and this page is stale. Refresh it (⌘R / F5).");
  }
  return "";
}
let lastServerHint = { text: "", at: 0 };
function logServerHint(error) {
  const hint = serverHint(error);
  // a run and a progress upload often fail together: say it once
  if (!hint || (hint === lastServerHint.text && Date.now() - lastServerHint.at < 10000)) return;
  lastServerHint = { text: hint, at: Date.now() };
  log(hint, "c-tip");
}
function logErrorHelp(meaning, cause, check) {
  const box = window.ErrorHelp?.render({ meaning, cause, check });
  if (box) appendConsoleNode(box);
}
function logCompileError(e) {
  const line = Math.max(1, Number(e.line) || 1);
  const file = firstErrorFile([e], levelFiles(S.cur), S.file);
  const button = document.createElement("button");
  button.type = "button";
  button.className = "c-error-line";
  button.textContent = `${e.file || "main.go"}:${e.line}:${e.col}  ${e.msg}`;
  button.onclick = () => jumpToLine(file, line);
  appendConsoleNode(button);
}
function jumpToLine(file, line) {
  showFile(file);
  const lines = ed.value.split("\n");
  const start = Math.min(ed.value.length, lines.slice(0, line - 1).reduce((n, value) => n + value.length + 1, 0));
  ed.focus();
  ed.setSelectionRange(start, Math.min(ed.value.length, start + (lines[line - 1] || "").length));
}
// Panic traces and failing test output name "main.go:12": make those lines
// clickable when the file is one of the mission's files.
function logOutputLine(text, cls) {
  const match = /(?:^|[\s/])([\w.-]+\.go):(\d+)\b/.exec(text);
  const file = match && levelFiles(S.cur).find((f) => f.name === match[1]);
  if (!file) return log(text, cls);
  const button = document.createElement("button");
  button.type = "button";
  button.className = `c-error-line ${cls}`;
  button.textContent = text;
  button.title = T(`روح لسطر ${match[2]} في ${match[1]}`, `Go to line ${match[2]} in ${match[1]}`);
  button.onclick = () => jumpToLine(file.name, Number(match[2]));
  appendConsoleNode(button);
}

// ---------- flip digits: changed characters roll in, like a counter ----------
function setDigits(el, text) {
  text = String(text);
  const prev = el.dataset.v ?? "";
  if (prev === text) return;
  el.dataset.v = text;
  el.classList.add("digits");
  el.innerHTML = [...text].map((ch, i) => {
    const changed = prev !== "" && prev[i] !== ch;
    return `<span${changed ? ' class="flip"' : ""}>${ch === " " ? "&nbsp;" : esc(ch)}</span>`;
  }).join("");
}

// ---------- sound: tiny synthesized cues, no audio files ----------
const Sound = {
  on: localStorage.getItem("orbital.sound") !== "off",
  ac: null, last: {},
  ensure() {
    if (!this.ac) this.ac = new (window.AudioContext || window.webkitAudioContext)();
    if (this.ac.state === "suspended") this.ac.resume();
    return this.ac;
  },
  tone(freq, dur, { type = "sine", gain = 0.05, at = 0, slide = 0 } = {}) {
    const ac = this.ensure(), t0 = ac.currentTime + at;
    const o = ac.createOscillator(), g = ac.createGain();
    o.type = type; o.frequency.setValueAtTime(freq, t0);
    if (slide) o.frequency.exponentialRampToValueAtTime(Math.max(40, freq + slide), t0 + dur);
    g.gain.setValueAtTime(0.0001, t0);
    g.gain.exponentialRampToValueAtTime(gain, t0 + 0.008);
    g.gain.exponentialRampToValueAtTime(0.0001, t0 + dur);
    o.connect(g).connect(ac.destination);
    o.start(t0); o.stop(t0 + dur + 0.02);
  },
  play(name) {
    if (!this.on) return;
    const now = performance.now(), gap = { step: 70, fail: 120, harvest: 45 }[name] || 0;
    if (gap && now - (this.last[name] || 0) < gap) return;
    this.last[name] = now;
    try {
      switch (name) {
        case "step": this.tone(240, 0.035, { gain: 0.012 }); break;
        case "harvest": this.tone(880, 0.07, { type: "triangle", gain: 0.045 }); this.tone(1320, 0.08, { type: "triangle", gain: 0.035, at: 0.045 }); break;
        case "plant": this.tone(520, 0.06, { type: "triangle", gain: 0.03, slide: 120 }); break;
        case "fail": this.tone(200, 0.16, { type: "square", gain: 0.018, slide: -70 }); break;
        case "report": this.tone(660, 0.07, { type: "triangle", gain: 0.04 }); this.tone(990, 0.1, { type: "triangle", gain: 0.035, at: 0.06 }); break;
        case "build": this.tone(330, 0.05, { type: "triangle", gain: 0.03 }); this.tone(495, 0.06, { type: "triangle", gain: 0.03, at: 0.05 }); break;
        case "error": this.tone(196, 0.14, { type: "sawtooth", gain: 0.025 }); this.tone(147, 0.22, { type: "sawtooth", gain: 0.025, at: 0.12 }); break;
        case "goal": [523, 659, 784, 1047].forEach((f, i) => this.tone(f, 0.16, { type: "triangle", gain: 0.05, at: i * 0.08 })); break;
        case "star": this.tone(1568, 0.18, { type: "sine", gain: 0.05 }); this.tone(2093, 0.22, { type: "sine", gain: 0.03, at: 0.04 }); break;
        case "hint": this.tone(740, 0.09, { type: "sine", gain: 0.03 }); break;
        case "storm": this.tone(110, 0.6, { type: "sawtooth", gain: 0.02, slide: -40 }); break;
      }
    } catch { /* audio is a nicety */ }
  },
};
function soundUI() {
  const b = $("#btn-sound");
  b.setAttribute("aria-pressed", String(Sound.on));
  b.querySelector("use").setAttribute("href", Sound.on ? "#i-sound" : "#i-mute");
  b.title = Sound.on ? T("اكتم الصوت", "Mute") : T("شغّل الصوت", "Unmute");
}
$("#btn-sound").onclick = () => {
  Sound.on = !Sound.on;
  localStorage.setItem("orbital.sound", Sound.on ? "on" : "off");
  soundUI();
  if (Sound.on) Sound.play("hint");
};
soundUI();

// ---------- world rendering: an isometric planet of textured blocks ----------
const cv = $("#world"), ctx = cv.getContext("2d");
let W = null;
const fx = [];
const particles = []; // {tx, ty: tile coords (y north), h: height px, vh, vx, vy, life, color, size}
const starfield = Array.from({ length: 160 }, () => ({ x: Math.random(), y: Math.random(), s: Math.random() < 0.15 ? 2 : 1, p: Math.random() * 6 }));
function burst(tx, ty, colors, n, { up = 1.6, spread = 1.2, size = 1 } = {}) {
  for (let i = 0; i < n; i++) {
    particles.push({ tx: tx + 0.5 + (Math.random() - 0.5) * 0.4, ty: ty + 0.5 + (Math.random() - 0.5) * 0.4,
      vx: (Math.random() - 0.5) * spread * 0.02, vy: (Math.random() - 0.5) * spread * 0.02,
      h: 2, vh: 0.6 + Math.random() * up * 0.5, life: 0.8,
      color: colors[Math.floor(Math.random() * colors.length)], size });
  }
}
const BURST = { moss: ["#74bd47", "#b2e57a", "#f6f39a"], lichen: ["#57c2b4", "#9be8d9"], crystal: ["#ff9ad5", "#ffd6ee", "#e45aa6"],
  glowcap: ["#9b6bff", "#8ff7ff"], reed: ["#ffcf5a", "#6fae4a"], iron: ["#d77a33", "#ffc07a", "#5c5566"], geode: ["#8f74ff", "#d6ccff"],
  copper: ["#d45842", "#ffa48e", "#5c5566"], sulfur: ["#e8e443", "#fffb9c", "#5c5566"],
  titanium: ["#8db3d9", "#d4e8fa", "#5c5566"], uranium: ["#40de7a", "#a6ffc4", "#5c5566"] };
const UNIT = { regolith: "#d39b55", soil: "#a45f3e", ice: "#9fd2cb", rock: "#77709a", stone: "#3a3542", tunnel: "#26202d", shaft: "#120f17" };
const ORES = new Set(["iron", "geode", "copper", "sulfur", "titanium", "uranium"]);
const ENT = { moss: "#93c45c", lichen: "#45b3a8", crystal: "#d75c9f", beacon: "#ef5a40", glowcap: "#9b6bff", reed: "#ffcf5a", iron: "#d77a33", geode: "#8f74ff",
  copper: "#d45842", sulfur: "#e8e443", titanium: "#8db3d9", uranium: "#40de7a" };
const NAME = LANG === "en"
  ? { regolith: "Regolith", soil: "Soil", ice: "Ice", rock: "Rock", stone: "Bedrock", tunnel: "Tunnel", shaft: "Shaft",
    moss: "Moss", lichen: "Lichen", crystal: "Crystal", beacon: "Beacon", glowcap: "Glowcap", reed: "Reed", iron: "Iron", geode: "Geode", cargo: "Cargo",
    water: "Water", silica: "Silica", copper: "Copper", sulfur: "Sulfur", titanium: "Titanium", helium: "Helium", uranium: "Uranium",
    biomass: "Biomass", plate: "Plate", wire: "Wire", glass: "Glass", fuel: "Fuel", fertilizer: "Fertilizer", circuit: "Circuit",
    alloy: "Alloy", battery: "Battery", coolant: "Coolant", processor: "Processor", frame: "Frame", cell: "Cell", core: "Core",
    composter: "Composter", smelter: "Smelter", kiln: "Kiln", refinery: "Refinery", assembler: "Assembler",
    "reactor-lab": "Reactor Lab", silo: "Silo", generator: "Generator", lab: "Lab", collector: "Gas Collector" }
  : { regolith: "ريجوليث", soil: "تربة", ice: "جليد", rock: "صخر", stone: "صخر صلب", tunnel: "نفق", shaft: "بير لتحت",
    moss: "طحالب", lichen: "أشنة", crystal: "كريستال", beacon: "منارة", glowcap: "فطر مضيء", reed: "بوص", iron: "حديد", geode: "جيود", cargo: "شحنة",
    water: "المياه", silica: "السيليكا", copper: "النحاس", sulfur: "الكبريت", titanium: "التيتانيوم", helium: "الهيليوم", uranium: "اليورانيوم",
    biomass: "الكتلة الحيوية", plate: "الصفيحة", wire: "السلك", glass: "الزجاج", fuel: "الوقود", fertilizer: "السماد", circuit: "الدائرة الإلكترونية",
    alloy: "السبيكة", battery: "البطارية", coolant: "سائل التبريد", processor: "المعالج", frame: "الهيكل", cell: "خلية الطاقة", core: "النواة",
    composter: "السمّادة", smelter: "المصهر", kiln: "الفرن", refinery: "المصفاة", assembler: "المجمّع",
    "reactor-lab": "معمل المفاعل", silo: "الصومعة", generator: "المولّد", lab: "المعمل", collector: "جامع الغاز" };
const INV_COLOR = {
  moss: ENT.moss, lichen: ENT.lichen, crystal: ENT.crystal, glowcap: ENT.glowcap, reed: ENT.reed, iron: ENT.iron, geode: ENT.geode,
  water: "#3a95e8", silica: "#d9cbb0", copper: ENT.copper, sulfur: ENT.sulfur, titanium: ENT.titanium, helium: "#f7bad2", uranium: ENT.uranium,
  biomass: "#58993e", plate: "#c6d0de", wire: "#f2a07e", glass: "#b0f2f5", fuel: "#e04665", fertilizer: "#9e7d47", circuit: "#2ab572",
  alloy: "#6e7ba3", battery: "#4e80f0", coolant: "#5ecdf5", processor: "#d166eb", frame: "#969689", cell: "#baf54e", core: "#6bf2c5",
  cargo: "#efe7d6",
};
const DRONE_COLORS = ["#ff7a52", "#f2d16b", "#8fd3ff", "#c7a2ff", "#efe7d6", "#93c45c"];
const INK = "#17142a", CONTACT = "#0d0a1c";
const COND = '"Barlow Condensed", "Readex Pro", sans-serif';
const noEmoji = (t) => String(t).replace(/\p{Extended_Pictographic}️?/gu, "").trim();

function setWorld(init) {
  const depth = init.layers || 1;
  const layers = Array.from({ length: depth }, () => Array.from({ length: init.w }, () => []));
  for (const t of init.tiles) layers[t.z || 0][t.x][t.y] = t;
  for (let z = 0; z < depth; z++) for (let x = 0; x < init.w; x++) for (let y = 0; y < init.h; y++) {
    layers[z][x][y] ||= { x, y, z, g: z ? "stone" : "regolith" };
  }
  if (init.gold) parText(init.par, init.gold, true);
  if (typeof Buildings !== "undefined") Buildings.syncState(init.buildings || []);
  W = { w: init.w, h: init.h, layers, tiles: layers[0], depth, view: 0, follow: true, t: 0, maxTicks: init.maxTicks, inv: init.inv || {}, storm: 0,
        drones: new Map([[0, { x: init.x, y: init.y, z: init.z || 0, px: init.x, py: init.y, bump: 0 }]]) };
  lastInv = { ...W.inv };
  // let narrow layouts size the map frame to the isometric sheet
  $(".map-frame").style.setProperty("--map-ar", `${(init.w + init.h) * 2} / ${init.w + init.h + 4}`);
  renderLayerBar();
  renderLegend();
  hud(init.progress);
}

// layer switcher: shown only when the world has more than one layer
function renderLayerBar() {
  let bar = $("#layer-bar");
  if (!bar) {
    bar = document.createElement("div");
    bar.id = "layer-bar"; bar.className = "layer-bar";
    $(".map-frame").appendChild(bar);
  }
  if (!W || W.depth < 2) { bar.innerHTML = ""; bar.hidden = true; return; }
  bar.hidden = false;
  bar.innerHTML = Array.from({ length: W.depth }, (_, z) =>
    `<button class="layer${z === W.view ? " on" : ""}" data-z="${z}" title="${z ? T(`طبقة ${z} تحت السطح`, `Layer ${z} below the surface`) : T("السطح", "Surface")}">${z ? `−${z}` : T("سطح", "Surface")}</button>`).join("")
    + `<label class="follow" title="${T("الكاميرا تتبع الدرون", "Camera follows the drone")}"><input type="checkbox" ${W.follow ? "checked" : ""}> ${T("تتبّع", "Follow")}</label>`;
  bar.querySelectorAll("button").forEach((b) => (b.onclick = () => { W.view = +b.dataset.z; W.follow = false; renderLayerBar(); renderLegend(); }));
  bar.querySelector("input").onchange = (e) => { W.follow = e.target.checked; renderLayerBar(); };
}

function renderLegend() {
  const g = new Set(), e = new Set(), b = new Set();
  for (const col of W.layers[W.view]) for (const t of col) { g.add(t.g); if (t.e) e.add(t.e); if (t.b) b.add(t.b); }
  const bKinds = typeof Buildings !== "undefined" ? Buildings.KINDS.filter((k) => b.has(k)) : [];
  const items = [
    ...["regolith", "soil", "ice", "rock", "stone", "tunnel", "shaft"].filter((k) => g.has(k)),
    ...["moss", "lichen", "glowcap", "reed", "crystal", "iron", "geode", "copper", "sulfur", "titanium", "uranium", "beacon"].filter((k) => e.has(k)),
    ...bKinds,
  ];
  const thumb = (k) => UNIT[k] ? Sprites.ground(k, 1)
    : ORES.has(k) ? Sprites.ore(k, 1)
    : typeof Buildings !== "undefined" && b.has(k) ? Buildings.sprite(k, 0)
    : Sprites.entity({ e: k, p: 0, r: 0, s: 6 }, 0, 1);
  $("#legend").innerHTML = items.map((k) => `<span><i style="background:${UNIT[k] ? "" : "#2a2447 "}url(${Sprites.dataURL(thumb(k))})"></i>${NAME[k] || k}</span>`).join("")
    + `<span><i style="background:#2a2447 url(${Sprites.dataURL(Sprites.drone(DRONE_COLORS[0], 0))})"></i>${T("الدرون", "Drone")}</span>`;
}

let lastInv = {};
function hud(progress) {
  if (!W) return;
  const tickEl = $("#hud-tick"), fill = $("#energy-fill");
  if (W.colony) {
    fill.style.transform = "scaleX(1)"; fill.classList.remove("low");
    setDigits(tickEl, `T ${W.t}`);
  } else if (!W.maxTicks) {
    fill.style.transform = "scaleX(1)"; fill.classList.remove("low");
    setDigits(tickEl, "LAB");
  } else {
    const left = Math.max(0, W.maxTicks - W.t);
    fill.style.transform = `scaleX(${left / W.maxTicks})`;
    fill.classList.toggle("low", left / W.maxTicks < 0.2);
    setDigits(tickEl, `T ${W.t} / ${W.maxTicks}`);
  }
  const inv = $("#hud-inv");
  const keys = Object.keys(W.inv);
  if ([...inv.children].map((c) => c.dataset.k).join() !== keys.join()) {
    inv.innerHTML = keys.map((k) => `<span class="chip" data-k="${k}" title="${NAME[k] || k}"><i style="background:${INV_COLOR[k] || "#ccc"}"></i><b></b></span>`).join("");
  }
  for (const c of inv.children) setDigits(c.querySelector("b"), W.inv[c.dataset.k]);
  lastInv = { ...W.inv };
  if (progress !== undefined) $("#hud-progress").textContent = progress;
}

function hash(x, y) { const s = Math.sin(x * 127.1 + y * 311.7) * 43758.5453; return s - Math.floor(s); }

// Block height, as a fraction of the half-width, for each ground.
const BLOCK = { regolith: 0.34, soil: 0.28, ice: 0.3, rock: 0.78, stone: 0.5, tunnel: 0.08, shaft: 0.08 };

// Isometric geometry: tile (x, y) with y pointing north; "row" runs south.
// Screen: X = ox + (x - row) * hw, Y = oy + (x + row) * hh, hh = hw / 2.
function isoGeometry(cw, ch) {
  const n = W.w + W.h;
  let hw = Math.floor(Math.min((cw - 40) / n, (ch - 70) / (n / 2 + 1.6)));
  hw = Math.max(8, hw);
  if (hw >= 16) hw = Math.floor(hw / 8) * 8; // keep 16px textures on whole pixels
  const hh = hw / 2;
  const width = n * hw, height = n * hh;
  const ox = Math.round((cw - width) / 2 + W.h * hw);
  const oy = Math.round((ch - height) / 2 + hw * 0.55);
  const toScreen = (x, row, lift = 0) => [ox + (x - row) * hw, oy + (x + row) * hh - lift];
  return { hw, hh, ox, oy, toScreen };
}

// Compass rose drawn along the real iso axes: each arm points the way
// colony.North/East/South/West actually moves the drone on screen.
function drawCompass(G, cw) {
  const L = Math.hypot(G.hw, G.hh), R = 30, CX = cw - 52, CY = 52;
  const arms = [["N", G.hw, -G.hh, "#f2b84b"], ["E", G.hw, G.hh, "#efe7d6"], ["S", -G.hw, G.hh, "#efe7d6"], ["W", -G.hw, -G.hh, "#efe7d6"]];
  ctx.save();
  ctx.fillStyle = "rgba(20,16,32,0.55)"; ctx.strokeStyle = "rgba(239,231,214,0.25)"; ctx.lineWidth = 1;
  ctx.beginPath(); ctx.arc(CX, CY, R + 14, 0, Math.PI * 2); ctx.fill(); ctx.stroke();
  ctx.font = `700 12px ${COND}`; ctx.textAlign = "center"; ctx.textBaseline = "middle";
  for (const [k, dx, dy, col] of arms) {
    const ux = dx / L, uy = dy / L, px = -uy, py = ux;
    const tx = CX + ux * R, ty = CY + uy * R;
    ctx.fillStyle = col;
    ctx.beginPath(); ctx.moveTo(tx, ty); ctx.lineTo(CX + px * 4, CY + py * 4); ctx.lineTo(CX - px * 4, CY - py * 4); ctx.closePath();
    ctx.globalAlpha = k === "N" ? 1 : 0.7; ctx.fill(); ctx.globalAlpha = 1;
    ctx.fillText(k, CX + ux * (R + 8), CY + uy * (R + 8));
  }
  ctx.restore();
}

// centre of a tile's top face, in fractional tile coords (y north)
function tileTop(G, tx, ty, lift) {
  const row = W.h - ty; // ty is a centre like y + 0.5
  return G.toScreen(tx, row, lift);
}

let hover = null;
cv.addEventListener("pointermove", (e) => {
  if (!W) return;
  const r = cv.getBoundingClientRect();
  const G = isoGeometry(r.width, r.height);
  const X = e.clientX - r.left - G.ox, Y = e.clientY - r.top - G.oy + G.hw * 0.3;
  const u = (X / G.hw + Y / G.hh) / 2, v = (Y / G.hh - X / G.hw) / 2;
  const x = Math.floor(u), row = Math.floor(v), y = W.h - 1 - row;
  hover = x >= 0 && y >= 0 && x < W.w && y < W.h ? { x, y } : null;
});
cv.addEventListener("pointerleave", () => (hover = null));

// map a 16px texture onto a parallelogram: o = origin, a = image x axis, b = image y axis
function texQuad(img, o, a, b) {
  ctx.save();
  ctx.setTransform(1, 0, 0, 1, 0, 0);
  const dpr = window.devicePixelRatio || 1;
  ctx.setTransform(dpr * a[0] / 16, dpr * a[1] / 16, dpr * b[0] / 16, dpr * b[1] / 16, dpr * o[0], dpr * o[1]);
  ctx.drawImage(img, 0, 0);
  ctx.restore();
}

function poly(points, fill) {
  ctx.beginPath(); ctx.moveTo(points[0][0], points[0][1]);
  for (let i = 1; i < points.length; i++) ctx.lineTo(points[i][0], points[i][1]);
  ctx.closePath(); ctx.fillStyle = fill; ctx.fill();
}

function drawBlock(G, t, x, y, now) {
  const row = W.h - 1 - y;
  const h = Math.round(BLOCK[t.g] * G.hw);
  const [X, Y] = G.toScreen(x, row);
  const T = [X, Y - h], R = [X + G.hw, Y + G.hh - h], B = [X, Y + 2 * G.hh - h], L = [X - G.hw, Y + G.hh - h];
  const v = Math.floor(hash(x * 3 + 1, y * 7 + 2) * 8);
  const glint = t.g === "ice" && (Math.floor(now / 800) + x * 3 + y * 5) % 9 === 0 ? 1 : 0;
  const tex = Sprites.ground(t.g, v, glint);
  // side faces: the same material, cut and shaded
  if (h > 0) {
    texQuad(tex, L, [B[0] - L[0], B[1] - L[1]], [0, h]);
    poly([L, B, [B[0], B[1] + h], [L[0], L[1] + h]], "rgba(10,6,20,0.30)");
    texQuad(tex, B, [R[0] - B[0], R[1] - B[1]], [0, h]);
    poly([B, R, [R[0], R[1] + h], [B[0], B[1] + h]], "rgba(10,6,20,0.50)");
  }
  // top face
  texQuad(tex, T, [R[0] - T[0], R[1] - T[1]], [L[0] - T[0], L[1] - T[1]]);
  if (t.g === "shaft") { // a hole that drops into the dark
    const k = 0.2, c = [X, Y + G.hh - h];
    const p = (P) => [c[0] + (P[0] - c[0]) * (1 - k), c[1] + (P[1] - c[1]) * (1 - k)];
    poly([p(T), p(R), p(B), p(L)], "#07050c");
    poly([p(T), p(R), [p(R)[0], p(R)[1] + G.hh * 0.5], [p(T)[0], p(T)[1] + G.hh * 0.5]], "rgba(70,60,90,0.55)");
    poly([p(L), p(T), [p(T)[0], p(T)[1] + G.hh * 0.5], [p(L)[0], p(L)[1] + G.hh * 0.5]], "rgba(110,98,130,0.45)");
  }
  if (ORES.has(t.e) && t.g === "stone") { // ore showing in the rock
    texQuad(Sprites.ore(t.e, v), T, [R[0] - T[0], R[1] - T[1]], [L[0] - T[0], L[1] - T[1]]);
    texQuad(Sprites.ore(t.e, v + 1), B, [R[0] - B[0], R[1] - B[1]], [0, h]);
  }
  // crisp top edges: lit back edges, dark front edges
  ctx.lineWidth = 1;
  ctx.strokeStyle = "rgba(255,245,225,0.18)";
  ctx.beginPath(); ctx.moveTo(L[0], L[1]); ctx.lineTo(T[0], T[1]); ctx.lineTo(R[0], R[1]); ctx.stroke();
  ctx.strokeStyle = "rgba(8,5,18,0.55)";
  ctx.beginPath(); ctx.moveTo(L[0], L[1]); ctx.lineTo(B[0], B[1]); ctx.lineTo(R[0], R[1]);
  if (h > 0) { ctx.moveTo(B[0], B[1]); ctx.lineTo(B[0], B[1] + h); }
  ctx.stroke();
  return h;
}

let darkCanvas = null; // reused every frame for the underground lighting
function darkLayer() { return (darkCanvas ||= document.createElement("canvas")); }

function draw(now) {
  const dpr = window.devicePixelRatio || 1;
  const cw = cv.clientWidth, ch = cv.clientHeight;
  if (cv.width !== Math.round(cw * dpr) || cv.height !== Math.round(ch * dpr)) { cv.width = Math.round(cw * dpr); cv.height = Math.round(ch * dpr); }
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  ctx.imageSmoothingEnabled = false;
  ctx.fillStyle = INK; ctx.fillRect(0, 0, cw, ch);
  for (const st of starfield) {
    ctx.globalAlpha = 0.25 + 0.3 * (0.5 + 0.5 * Math.sin(now / 1100 + st.p));
    ctx.fillStyle = "#efe7d6"; ctx.fillRect(Math.floor(st.x * cw), Math.floor(st.y * ch), st.s, st.s);
  }
  ctx.globalAlpha = 1;
  if (!W) return;

  const G = isoGeometry(cw, ch);
  const grid = W.layers[W.view];
  const under = W.view > 0;

  // the slab the sheet sits on: a deep cut through the planet
  {
    const depthPx = G.hw * (under ? 1.6 : 1.1);
    const P = (x, row) => G.toScreen(x, row);
    const L = P(0, W.h), B = P(W.w, W.h), R = P(W.w, 0);
    const strata = under ? ["#1b1722", "#231d2c", "#16121b"] : ["#5c3a22", "#3e2a1e", "#2a2131"];
    poly([L, B, [B[0], B[1] + depthPx], [L[0], L[1] + depthPx]], strata[0]);
    poly([B, R, [R[0], R[1] + depthPx], [B[0], B[1] + depthPx]], strata[1]);
    for (let k = 1; k < 4; k++) { // strata lines
      const d = (depthPx * k) / 4;
      ctx.strokeStyle = "rgba(0,0,0,0.35)"; ctx.lineWidth = 1;
      ctx.beginPath(); ctx.moveTo(L[0], L[1] + d); ctx.lineTo(B[0], B[1] + d); ctx.lineTo(R[0], R[1] + d); ctx.stroke();
    }
  }

  // back to front: diagonals of x + row
  const order = [];
  for (let x = 0; x < W.w; x++) for (let y = 0; y < W.h; y++) order.push([x, y, x + (W.h - 1 - y)]);
  order.sort((a, b) => a[2] - b[2] || a[0] - b[0]);
  const lifts = {};
  for (const [x, y] of order) {
    const t = grid[x][y];
    const h = drawBlock(G, t, x, y, now);
    lifts[x + "," + y] = h;
    if (t.e && !(t.g === "stone" && ORES.has(t.e))) drawEntity(G, t, x, y, h, now);
    if (t.b && typeof Buildings !== "undefined") {
      const [bx, by] = G.toScreen(x + 0.5, W.h - 1 - y + 0.5, h);
      Buildings.draw(ctx, bx, by, G.hw, t.b, W.view === 0 ? Buildings.at(x, y) : null, now);
    }
  }

  // underground: darkness, with the drones' lamps lighting the rock around them
  if (under) {
    ctx.save();
    const dark = darkLayer();
    if (dark.width !== cv.width || dark.height !== cv.height) { dark.width = cv.width; dark.height = cv.height; }
    const dc = dark.getContext("2d");
    dc.setTransform(1, 0, 0, 1, 0, 0); dc.globalCompositeOperation = "source-over"; dc.clearRect(0, 0, dark.width, dark.height);
    dc.setTransform(dpr, 0, 0, dpr, 0, 0);
    dc.fillStyle = "rgba(6,4,12,0.62)"; dc.fillRect(0, 0, cw, ch);
    dc.globalCompositeOperation = "destination-out";
    for (const [, d] of W.drones) {
      if ((d.z || 0) !== W.view) continue;
      const [X, Y] = tileTop(G, d.x + 0.5, d.y + 0.5, 0);
      const r = G.hw * 4.5, gr = dc.createRadialGradient(X, Y, G.hw * 0.5, X, Y, r);
      gr.addColorStop(0, "rgba(0,0,0,1)"); gr.addColorStop(1, "rgba(0,0,0,0)");
      dc.fillStyle = gr; dc.beginPath(); dc.ellipse(X, Y, r, r * 0.7, 0, 0, 7); dc.fill();
    }
    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.drawImage(dark, 0, 0);
    ctx.restore();
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  }

  // hovered tile
  if (hover && !playback) {
    const t = grid[hover.x][hover.y], h = lifts[hover.x + "," + hover.y] || 0;
    const [X, Y] = G.toScreen(hover.x, W.h - 1 - hover.y, h);
    ctx.strokeStyle = "#efe7d6"; ctx.lineWidth = 2; ctx.setLineDash([4, 3]);
    ctx.beginPath(); ctx.moveTo(X, Y); ctx.lineTo(X + G.hw, Y + G.hh); ctx.lineTo(X, Y + 2 * G.hh); ctx.lineTo(X - G.hw, Y + G.hh); ctx.closePath(); ctx.stroke();
    ctx.setLineDash([]);
    const label = `(${hover.x}, ${hover.y}${W.depth > 1 ? `, z ${W.view}` : ""}) ${t.g}${t.e ? " · " + t.e + (t.s ? " " + t.s : "") : ""}${t.b ? " · " + t.b : ""}`;
    ctx.font = `600 12.5px ${COND}`;
    const tw = ctx.measureText(label).width + 14;
    const lx = Math.min(Math.max(X - tw / 2, 4), cw - tw - 4), ly = Math.max(Y - 30, 4);
    ctx.fillStyle = "#efe7d6"; ctx.beginPath(); ctx.roundRect(lx, ly, tw, 20, 4); ctx.fill();
    ctx.fillStyle = INK; ctx.textAlign = "left"; ctx.textBaseline = "middle"; ctx.fillText(label, lx + 7, ly + 10.5);
  }

  // drones on this layer (interpolated)
  const k = playback ? Math.min(1, playback.acc / playback.step) : 1;
  const e = 1 - Math.pow(1 - k, 3);
  const drones = [...W.drones].filter(([, d]) => (d.z || 0) === W.view);
  drones.sort((a, b) => (a[1].x - a[1].y) - (b[1].x - b[1].y));
  for (const [id, d] of drones) {
    let ee = e;
    if (W.live) { const kk = Math.min(1, (now - (d.t0 || 0)) / (W.liveStep || 260)); ee = 1 - Math.pow(1 - kk, 3); }
    const dx = d.px + (d.x - d.px) * ee, dy = d.py + (d.y - d.py) * ee;
    let shake = 0;
    if (d.bump > 0) { shake = Math.sin(now / 18) * G.hw * 0.08 * d.bump; d.bump = Math.max(0, d.bump - 0.05); }
    const h = lifts[Math.round(dx) + "," + Math.round(dy)] || 0;
    const [X, Y] = tileTop(G, dx + 0.5, dy + 0.5, h);
    drawDrone(G, X + shake, Y, now, id);
  }
  // other layers' drones: a small marker so they are not lost
  const elsewhere = [...W.drones].filter(([, d]) => (d.z || 0) !== W.view).length;
  if (elsewhere) {
    ctx.font = `600 12px ${COND}`; ctx.fillStyle = "#a69dc2"; ctx.textAlign = "right"; ctx.textBaseline = "top";
    ctx.fillText(`${elsewhere} DRONE${elsewhere > 1 ? "S" : ""} ON OTHER LAYERS`, cw - 12, 12);
  }

  // particles
  const u = Math.max(2, Math.round(G.hw / 8));
  for (let i = particles.length - 1; i >= 0; i--) {
    const p = particles[i];
    p.life -= 0.03; p.vh -= 0.1; p.h += p.vh; p.tx += p.vx; p.ty += p.vy;
    if (p.h < 0) { p.h = 0; p.vh *= -0.3; }
    if (p.life <= 0) { particles.splice(i, 1); continue; }
    const base = lifts[Math.floor(p.tx) + "," + Math.floor(p.ty)] || 0;
    if (particles.length > 260) particles.splice(0, particles.length - 260);
    const [X, Y] = tileTop(G, p.tx, p.ty, base + p.h * (G.hw / 24));
    ctx.globalAlpha = Math.min(1, p.life * 1.6);
    ctx.fillStyle = p.color; ctx.fillRect(Math.round(X), Math.round(Y), u * p.size, u * p.size);
  }
  ctx.globalAlpha = 1;

  // dust storm over the whole sheet
  if (W.storm > 0) {
    ctx.fillStyle = W.storm === 2 ? "rgba(211,155,85,0.30)" : "rgba(211,155,85,0.09)";
    ctx.fillRect(0, 0, cw, ch);
    ctx.strokeStyle = W.storm === 2 ? "rgba(239,231,214,0.32)" : "rgba(239,231,214,0.1)"; ctx.lineWidth = 2;
    const off = (now / (W.storm === 2 ? 9 : 30)) % 18;
    for (let i = -ch; i < cw + ch; i += 18) { ctx.beginPath(); ctx.moveTo(i + off, 0); ctx.lineTo(i + off - ch, ch); ctx.stroke(); }
  }

  // coordinates along the two front edges, and the sheet's collar
  ctx.font = `600 ${Math.max(10, Math.min(13, G.hw * 0.42))}px ${COND}`; ctx.fillStyle = "#a69dc2";
  ctx.textAlign = "center"; ctx.textBaseline = "middle";
  const slab = G.hw * (under ? 1.6 : 1.1);
  for (let x = 0; x < W.w; x++) { const [X, Y] = G.toScreen(x + 0.5, W.h); ctx.fillText(String(x), X - 10, Y + slab * 0.5 + 10); }
  for (let y = 0; y < W.h; y++) { const [X, Y] = G.toScreen(W.w, W.h - 1 - y + 0.5); ctx.fillText(String(y), X + 10, Y + slab * 0.5 + 10); }
  ctx.textAlign = "left"; ctx.textBaseline = "alphabetic"; ctx.fillStyle = "#cfc6b2"; ctx.font = `700 12px ${COND}`;
  const code = S.cur ? sheetCode(S.cur) : "";
  ctx.fillText(`KEPLER-G · SHEET ${code} · ${W.w}×${W.h}${W.depth > 1 ? ` · LAYER ${W.view ? "−" + W.view : "SURFACE"}` : ""}`, 12, ch - 12);
  drawCompass(G, cw);

  // floating labels
  for (let i = fx.length - 1; i >= 0; i--) {
    const f = fx[i];
    f.age += 16;
    if (f.age > 1100) { fx.splice(i, 1); continue; }
    const h = lifts[f.x + "," + f.y] || 0;
    const [X, Y] = tileTop(G, f.x + 0.5, f.y + 0.5, h + G.hw * 0.9 + f.age / 30);
    ctx.globalAlpha = 1 - f.age / 1100;
    ctx.font = `700 ${Math.max(13, G.hw * 0.5)}px ${COND}`;
    ctx.textAlign = "center"; ctx.textBaseline = "middle";
    ctx.lineWidth = 4; ctx.strokeStyle = INK; ctx.strokeText(f.text, X, Y);
    ctx.fillStyle = f.color; ctx.fillText(f.text, X, Y);
    ctx.globalAlpha = 1;
  }
}

// upright sprites (plants, crystals, beacons) standing on a block's top face
function drawEntity(G, t, x, y, h, now) {
  const [X, Y0] = G.toScreen(x + 0.5, W.h - 1 - y + 0.5, h);
  const Y = Y0;
  const size = Math.round(G.hw * 1.15 / 8) * 8 || G.hw;
  const ripe = t.e !== "crystal" && t.e !== "beacon" && W.t >= (t.r || 0);
  ctx.fillStyle = "rgba(10,6,20,0.35)"; // contact shadow on the top face
  ctx.beginPath(); ctx.ellipse(X, Y, G.hw * 0.42, G.hh * 0.42, 0, 0, 7); ctx.fill();
  if (t.e === "beacon") {
    const p = (now / 1400) % 1;
    ctx.strokeStyle = `rgba(239,90,64,${0.7 * (1 - p)})`; ctx.lineWidth = 2;
    ctx.beginPath(); ctx.ellipse(X, Y, G.hw * 0.4 * (1 + p), G.hh * 0.4 * (1 + p), 0, 0, 7); ctx.stroke();
  }
  const img = Sprites.entity(t, now, W.t);
  const plant = t.e === "moss" || t.e === "reed" || t.e === "glowcap" || t.e === "lichen";
  const phase = now / 650 + x * 0.9 + y * 1.7;
  ctx.save();
  ctx.translate(Math.round(X), Math.round(Y + size * 0.2));
  if (plant) {
    const bob = ripe ? -Math.abs(Math.sin(phase * 1.3)) * Math.max(1, G.hw / 24) : 0;
    ctx.translate(0, bob);
    ctx.transform(1, 0, Math.sin(phase) * (t.e === "reed" ? 0.09 : 0.035), 1, 0, 0);
  }
  ctx.drawImage(img, -size / 2, -size, size, size);
  ctx.restore();
  const glint = (Math.floor(now / 160) + x * 13 + y * 7) % (t.e === "crystal" ? 23 : 37);
  if ((t.e === "crystal" || ripe) && glint < 3) {
    const u = Math.max(2, Math.round(G.hw / 12));
    const gx = X + size * (hash(x, y) - 0.5) * 0.5, gy = Y - size * (0.45 + 0.3 * hash(y, x));
    ctx.fillStyle = glint === 1 ? "#ffffff" : "rgba(255,255,255,0.6)";
    ctx.fillRect(gx, gy - u, u, u * 3); ctx.fillRect(gx - u, gy, u * 3, u);
  }
  if (t.e === "crystal" && t.s) {
    ctx.font = `700 ${Math.max(10, G.hw * 0.42)}px ${COND}`; ctx.textAlign = "center"; ctx.textBaseline = "alphabetic";
    ctx.lineWidth = 3; ctx.strokeStyle = CONTACT; ctx.strokeText(t.s, X + size * 0.42, Y - size * 0.55);
    ctx.fillStyle = "#ffd6ee"; ctx.fillText(t.s, X + size * 0.42, Y - size * 0.55);
  }
}

// a hovering drone with its shadow on the block under it
function drawDrone(G, X, Y, now, id = 0) {
  const col = DRONE_COLORS[id % DRONE_COLORS.length];
  const size = Math.round(G.hw * 1.5 / 10) * 10 || G.hw * 1.5;
  const bob = Math.sin(now / 280 + id) * Math.max(1, G.hw * 0.06);
  ctx.save();
  ctx.fillStyle = "rgba(10,6,20,0.42)";
  ctx.beginPath(); ctx.ellipse(X, Y - G.hh * 0.05, G.hw * 0.55 - bob * 0.3, G.hh * 0.55 - bob * 0.15, 0, 0, 7); ctx.fill();
  ctx.imageSmoothingEnabled = false;
  const lift = G.hw * 0.75;
  ctx.drawImage(Sprites.drone(col, Math.floor(now / 55) % 8), Math.round(X - size / 2), Math.round(Y - lift - size / 2 + bob), size, size);
  if (id > 0) {
    ctx.font = `700 ${Math.max(10, G.hw * 0.4)}px ${COND}`; ctx.textAlign = "center"; ctx.textBaseline = "middle";
    ctx.lineWidth = 3; ctx.strokeStyle = CONTACT; ctx.strokeText(String(id), X, Y - lift - size * 0.62);
    ctx.fillStyle = "#efe7d6"; ctx.fillText(String(id), X, Y - lift - size * 0.62);
  }
  ctx.restore();
}

// a little pixel world for the welcome screen
function drawSurveyArt(canvas) {
  const c = canvas.getContext("2d"), dpr = window.devicePixelRatio || 1;
  const w = canvas.clientWidth, h = canvas.clientHeight;
  canvas.width = w * dpr; canvas.height = h * dpr; c.setTransform(dpr, 0, 0, dpr, 0, 0);
  c.imageSmoothingEnabled = false;
  c.fillStyle = INK; c.fillRect(0, 0, w, h);
  const cell = 32, cols = Math.ceil(w / cell), rows = Math.ceil(h / cell);
  const seeds = Array.from({ length: 8 }, (_, i) => ({ x: Math.random() * cols, y: Math.random() * rows, u: ["soil", "regolith", "ice", "rock", "soil", "regolith", "ice", "soil"][i] }));
  const unitAt = (x, y) => { let b = seeds[0], bd = 1e9; for (const s of seeds) { const d = (s.x - x) ** 2 + (s.y - y) ** 2; if (d < bd) { bd = d; b = s; } } return b.u; };
  for (let x = 0; x < cols; x++) for (let y = 0; y < rows; y++) {
    const g = unitAt(x, y);
    c.drawImage(Sprites.ground(g, Math.floor(hash(x, y) * 8)), x * cell, y * cell, cell, cell);
    const r = hash(y * 3, x * 5);
    let e = null;
    if (g === "soil" && r < 0.55) e = { e: r < 0.2 ? "glowcap" : "moss", p: 0, r: r < 0.4 ? 0 : 99 };
    if (g === "ice" && r < 0.35) e = { e: "lichen", p: 0, r: 0 };
    if (g === "rock" && r < 0.3) e = { e: "crystal", s: 3 + Math.floor(r * 20) };
    if (e) c.drawImage(Sprites.entity(e, 0, 50), x * cell, y * cell, cell, cell);
  }
  const dx = Math.floor(cols / 2) * cell, dy = Math.floor(rows / 2) * cell;
  c.fillStyle = "rgba(10,6,20,0.4)"; c.beginPath(); c.ellipse(dx + 20, dy + 26, 15, 6, 0, 0, 7); c.fill();
  c.drawImage(Sprites.drone(DRONE_COLORS[0], 2), dx - 4, dy - 12, 40, 40);
  c.strokeStyle = "#5a4f8c"; c.lineWidth = 2; c.strokeRect(8, 8, w - 16, h - 16);
  c.fillStyle = "#efe7d6"; c.font = `700 13px ${COND}`; c.textAlign = "left"; c.textBaseline = "alphabetic";
  c.lineWidth = 4; c.strokeStyle = INK; c.strokeText("KEPLER-G · SURVEY SHEET 1", 18, h - 18); c.fillText("KEPLER-G · SURVEY SHEET 1", 18, h - 18);
}

// ---------- playback ----------
let playback = null;
const speedEl = $("#speed");

function startPlayback(events, done) {
  const init = events.find((e) => e.op === "init");
  if (init) setWorld(init);
  // realtime missions interleave drones; play them in tick order
  const q = events.filter((e) => e.op !== "init").map((e, i) => [e, i]);
  q.sort((a, b) => a[0].t - b[0].t || a[1] - b[1]);
  playback = { q: q.map((p) => p[0]), i: 0, acc: 0, step: 1, done };
  $("#btn-skip").disabled = false;
}

const DRONE_OPS = new Set(["move", "harvest", "plant", "till", "scan", "wait", "unload", "export", "swap", "dig", "digdown", "descend", "ascend"]);
function apply(ev) {
  W.t = Math.max(W.t, ev.t);
  const id = ev.d || 0;
  if (ev.op === "spawn") {
    W.drones.set(id, { x: ev.x, y: ev.y, px: ev.x, py: ev.y, bump: 0 });
    log(`t=${ev.t} ${noEmoji(ev.msg)}`, "c-info");
    renderLegend();
  } else if (DRONE_OPS.has(ev.op)) {
    const d = W.drones.get(id) || { bump: 0 };
    if (W.live && (d.x !== ev.x || d.y !== ev.y)) { d.px = d.x ?? ev.x; d.py = d.y ?? ev.y; d.t0 = performance.now(); }
    d.x = ev.x; d.y = ev.y; d.z = ev.z || 0;
    if (ev.op === "descend" || ev.op === "ascend") { d.px = d.x; d.py = d.y; }
    W.drones.set(id, d);
    if (id === 0 && W.follow && d.z !== W.view) { W.view = d.z; renderLayerBar(); renderLegend(); }
  }
  if (ev.op === "building") {
    if (typeof Buildings !== "undefined") {
      const prev = W.layers?.[ev.z || 0]?.[ev.x]?.[ev.y]?.b;
      Buildings.applyEvent(ev, W);
      if (prev !== ev.building?.kind) renderLegend();
    }
    hud(ev.progress);
    return;
  }
  if (ev.op === "storm-warning") { W.storm = 1; log(`t=${ev.t} ${noEmoji(ev.msg)}`, "c-warn"); Sound.play("fail"); }
  else if (ev.op === "storm") { W.storm = 2; log(`t=${ev.t} ${noEmoji(ev.msg)}`, "c-err"); Sound.play("storm"); }
  else if (ev.op === "clear") { W.storm = 0; log(`t=${ev.t} ${noEmoji(ev.msg)}`, "c-ok"); Sound.play("report"); }
  if (ev.tile) {
    const z = ev.tile.z || 0;
    const prev = (W.layers[z] || W.layers[0])[ev.tile.x]?.[ev.tile.y];
    (W.layers[z] || W.layers[0])[ev.tile.x][ev.tile.y] = ev.tile;
    if (z === 0 && !ev.tile.b && typeof Buildings !== "undefined") Buildings.remove(ev.tile.x, ev.tile.y);
    if (prev?.b !== ev.tile.b) renderLegend();
  }
  if (ev.ok && (ev.op === "dig" || ev.op === "digdown") && ev.tile) burst(ev.tile.x, ev.tile.y, ["#5c5566", "#3a3542", "#736b7d"], 10, { up: 1.2, spread: 1.6 });
  if (ev.inv) W.inv = ev.inv;
  if (ev.op === "report") {
    log(`Report: ${ev.msg}`, "c-report");
    fx.push({ x: ev.x, y: ev.y, text: ev.msg, color: "#9fd2cb", age: 0 });
    Sound.play("report");
  } else if (ev.op === "recall") {
    log(`t=${ev.t} ${noEmoji(ev.msg)}`, "c-warn"); W.storm = 1; Sound.play("fail");
  } else if (ev.op === "scan" && ev.msg) {
    fx.push({ x: ev.x, y: ev.y, text: noEmoji(ev.msg), color: ev.ok ? "#efe7d6" : "#f0a14a", age: 0 });
    if (!ev.ok) log(`t=${ev.t}${id ? T(` [درون ${id}]`, ` [drone ${id}]`) : ""} ${noEmoji(ev.msg)}`, "c-warn");
  } else if (ev.op === "goal") {
    log(ev.msg, "c-ok");
  } else if (ev.op === "fail") {
    log(ev.msg, "c-err");
    S.failLogged = true;
    Sound.play("error");
  } else if (ev.ok && ev.op === "unload") {
    fx.push({ x: ev.x, y: ev.y, text: noEmoji(ev.msg).replace(" crystal", ""), color: "#f2b33d", age: 0 });
    Sound.play("report");
  } else if (ev.ok && ev.op === "move") {
    if (ev.msg) { fx.push({ x: ev.x, y: ev.y, text: noEmoji(ev.msg), color: "#efe7d6", age: 0 }); Sound.play("report"); }
    else Sound.play("step");
    burst(ev.x, ev.y, ["#cfc6b2", "#a69dc2"], 3, { up: 0.35, spread: 0.9 });
  } else if (ev.ok && ev.op === "harvest" && ev.msg) {
    const m = ev.msg.match(/^\+(\d+(?:\.\d+)?) (\w+)/);
    burst(ev.x, ev.y, BURST[m && m[2]] || BURST.moss, 9, { up: 2, spread: 2 });
    fx.push({ x: ev.x, y: ev.y, text: m ? `+${m[1]}` : noEmoji(ev.msg), color: (m && INV_COLOR[m[2]]) || "#93c45c", age: 0 });
    Sound.play("harvest");
  } else if (ev.ok && ev.op === "plant") {
    burst(ev.x, ev.y, ["#87613f", "#6e4b2e", "#a8692f"], 7, { up: 0.9, spread: 1.4 });
    Sound.play("plant");
  } else if (!ev.ok && !(ev.op === "tile" && ev.ok === undefined)) {
    if (ev.op === "move") (W.drones.get(id) || {}).bump = 1;
    log(`t=${ev.t}${id ? T(` [درون ${id}]`, ` [drone ${id}]`) : ""} ${ev.op}: ${noEmoji(ev.msg || T("فشل", "failed"))}`, "c-warn");
    if (ev.tile) fx.push({ x: ev.x, y: ev.y, text: "×", color: "#ff6b6b", age: 0 });
    Sound.play("fail");
  }
  hud(ev.progress);
}

function tick(dt) {
  if (!playback) return;
  const speed = +speedEl.value;
  playback.step = 380 / speed;
  playback.acc += dt;
  while (playback.acc >= playback.step && playback.i < playback.q.length) {
    playback.acc -= playback.step;
    for (const d of W.drones.values()) { d.px = d.x; d.py = d.y; }
    // everything that happened on the same tick plays in one step
    const t = playback.q[playback.i].t;
    while (playback.i < playback.q.length && playback.q[playback.i].t === t) apply(playback.q[playback.i++]);
  }
  if (playback.i >= playback.q.length && playback.acc >= playback.step) finishPlayback();
}

function finishPlayback() {
  if (!playback) return;
  try {
    while (playback.i < playback.q.length) apply(playback.q[playback.i++]);
  } catch (e) { console.error(e); }
  for (const d of W.drones.values()) { d.px = d.x; d.py = d.y; }
  const done = playback.done;
  playback = null;
  $("#btn-skip").disabled = true;
  done && done();
}
$("#btn-skip").onclick = () => {
  if (dbgRun && !dbgRun.revealed) InlineDebugger.finishNow();
  else finishPlayback();
};

let lastT = performance.now();
function frame(now) {
  requestAnimationFrame(frame);
  const dt = now - lastT;
  lastT = now;
  try { tick(dt); } catch (e) { console.error(e); finishPlayback(); }
  try { draw(now); } catch (e) { console.error(e); }
}
requestAnimationFrame(frame);

// ---------- levels ----------
function parText(par, gold, dynamic) {
  $("#goal-strip-par").innerHTML = par > 0 ? `<span class="stars">${starHTML(3)}</span> ${gold} ${T("tick", "ticks")}` : "";
  $("#lvl-par").innerHTML = par > 0
    ? `<span class="stars">${starHTML(2)}</span> ${par} ${T("tick أو أقل", "ticks or fewer")} · <span class="stars">${starHTML(3)}</span> ${gold} ${T("tick أو أقل", "ticks or fewer")}${dynamic ? T(" (الخريطة دي)", " (this map)") : ""}`
    : "";
}
function idx(id) { return S.levels.findIndex((l) => l.id === id); }
function unlocked(i) {
  if (i < 0 || i >= S.levels.length) return false;
  const level = S.levels[i];
  if (level.chapter === 38) {
    if (S.levels.slice(i).some((l) => l.chapter === 38 && rec(l.id).stars > 0)) return true;
    const previous = S.levels[i - 1];
    if (previous?.chapter === 38) return rec(previous.id).stars > 0;
    const chapter17 = S.levels.filter((l) => l.chapter === 17);
    return chapter17.length > 0 && rec(chapter17[chapter17.length - 1].id).stars > 0;
  }
  if (i === 0 || S.levels.slice(i).some((l) => l.chapter !== 38 && rec(l.id).stars > 0)) return true;
  const previous = S.levels[i - 1];
  return !!previous.service || rec(previous.id).stars > 0;
}
function starHTML(n, total = 3) {
  return Array.from({ length: total }, (_, i) => `<svg class="star${i < n ? " on" : ""}" viewBox="0 0 24 24" aria-hidden="true"><use href="#i-star"/></svg>`).join("");
}
function sheetCode(l) { return l.id.split("-").slice(0, 2).join("-"); }
function updateOverall() {
  const done = S.levels.filter((l) => rec(l.id).stars > 0).length;
  setDigits($("#overall-count"), done);
  $("#overall-total").textContent = S.levels.length;
  $("#overall-fill").style.transform = `scaleX(${done / Math.max(1, S.levels.length)})`;
}

async function refreshService(level) {
  let badge = $("#service-badge");
  if (!badge) {
    badge = document.createElement("small");
    badge.id = "service-badge";
    badge.className = "service-badge";
    badge.setAttribute("role", "status");
    $("#lvl-goal").after(badge);
  }
  badge.hidden = !level.service;
  if (!level.service) return;
  const label = T("الخدمة المطلوبة: ", "Required service: ") + level.service;
  badge.textContent = label + T(" · بنفحص الحالة…", " · Checking status…");
  try {
    const statuses = await api("/api/services");
    if (S.cur !== level) return;
    badge.textContent = label + (statuses[level.service]
      ? T(" · متاحة", " · Available")
      : T(" · مش متاحة — شغّل make services-up", " · Unavailable — run make services-up"));
  } catch (e) {
    if (S.cur === level) badge.textContent = label + T(" · تعذّر فحص الحالة", " · Status check failed");
  }
}

let previewVersion = 0;
async function loadLevel(id, review = null) {
  if (S.running) stopRun(true);
  if (S.mode === "colony") Colony.exit();
  const version = ++previewVersion;
  flushSave();
  if (playback) finishPlayback();
  const l = review || S.levels[idx(id)];
  S.review = !!review; S.reviewDraft = {};
  S.cur = l; S.hintsShown = 0;
  clearTraceRun();
  LearningMetrics.openLevel(store, id);
  store.pref("current", id);
  // the tab shows the mission, and the URL hash links straight to it
  document.title = `${sheetCode(l)} · ${l.title} — Orbital Colony`;
  if (!review && location.hash !== "#" + id) history.replaceState(null, "", "#" + id);
  const ch = S.chapters.find((c) => c.n === l.chapter);
  $("#lvl-chapter").textContent = `${sheetCode(l)} · ${ch ? ch.title : ""}`;
  $("#lvl-stars").setAttribute("aria-label", T(`${rec(id).stars} من 3 نجوم`, `${rec(id).stars} of 3 stars`));
  $("#lvl-title").textContent = (S.review ? T("مراجعة · ", "Review · ") : "") + l.title;
  $("#lvl-goal").innerHTML = inline(l.goal);
  $("#goal-strip-text").innerHTML = inline(l.goal);
  refreshService(l);
  parText(l.par, l.gold, false);
  Lessons.open(l, (stepID) => LearningMetrics.recordStepViewed(store, id, stepID));
  LearningMetrics.recordLessonOpened(store, id);
  $("#hint-list").innerHTML = "";
  $("#btn-hint").disabled = !l.hints?.length;
  $("#btn-hint").querySelector("span").textContent = T(`تلميح (${l.hints?.length || 0})`, `Hint (${l.hints?.length || 0})`);
  const r = rec(id);
  $("#lvl-stars").innerHTML = starHTML(r.stars);
  lastErrors = [];
  lastMtime = 0;
  showFile(levelFiles(l)[0].name);
  linkUI();
  $("#btn-vscode").disabled = S.review;
  if (isLinked()) syncFromWorkspace(true);
  $("#save-state").textContent = "";
  $("#review-notice").classList.toggle("hidden", !S.review);
  $("#review-notice").textContent = l.kind === "lab"
    ? T("مراجعة · المختبر بيبدأ بكود البداية الأصلي، من غير حلك المحفوظ.", "Review · This lab starts with its original starter files, without your saved solution.")
    : T("مراجعة · كود البداية فاضي عن قصد. اكتب الحل من ذاكرتك.", "Review · The starter is intentionally empty. Write the solution from memory.");
  $("#banner").classList.add("hidden");
  con.innerHTML = "";
  log(T(`ورقة ${sheetCode(l)}: ${l.title}. اكتب الكود ودوس Run.`, `Sheet ${sheetCode(l)}: ${l.title}. Write the code and press Run.`), "c-info");
  document.body.classList.remove("lesson-open");
  $(".lesson-panel").scrollTop = 0;
  try {
    const preview = await api(`/api/preview?level=${encodeURIComponent(id)}`);
    if (S.cur === l && version === previewVersion) setWorld(preview);
  } catch (e) { if (S.cur === l) log(String(e), "c-err"); }
}

$("#btn-hint").onclick = () => {
  const l = S.cur;
  if (S.hintsShown >= l.hints.length) return;
  const d = document.createElement("div");
  d.className = "hint";
  const hintIndex = S.hintsShown;
  const hintText = inline(l.hints[S.hintsShown++].replace(/\p{Extended_Pictographic}\uFE0F?\s?/gu, ""));
  LearningMetrics.recordHint(store, l.id, hintIndex);
  const hintLabel = l.steps?.length ? [T("الفكرة", "Concept"), T("التنظيم", "Structure"), T("ركّز هنا", "Specific")][S.hintsShown - 1] : "";
  d.innerHTML = hintLabel ? `<strong>${hintLabel}</strong><div>${hintText}</div>` : hintText;
  Sound.play("hint");
  $("#hint-list").appendChild(d);
  $("#btn-hint").querySelector("span").textContent = S.hintsShown < l.hints.length ? T(`تلميح كمان (${l.hints.length - S.hintsShown})`, `Another hint (${l.hints.length - S.hintsShown})`) : T("خلصت التلميحات", "No more hints");
  $("#btn-hint").disabled = S.hintsShown >= l.hints.length;
};

$("#btn-reset").onclick = () => {
  if (!confirm(T("ترجع لكود البداية؟ الكود الحالي هيتمسح.", "Go back to the starter code? Your current code will be erased."))) return;
  const f = levelFiles(S.cur).find((x) => x.name === S.file);
  if (!f || f.readOnly) return;
  // replace through the edit history so Ctrl/⌘+Z brings the old code back
  ed.focus();
  applyEdit({ from: 0, to: ed.value.length, text: f.starter, selStart: 0, selEnd: 0 });
  if (ed.value !== f.starter) { ed.value = f.starter; onEdit(); }
  log(T("رجعت لكود البداية. Ctrl/⌘+Z يرجّع كودك.", "Back to the starter code. Ctrl/⌘+Z brings your code back."), "c-info");
};

let quizAttempt = null;
function openQuiz(chapter) {
  quizAttempt = { chapter, answers: new Map() };
  const panel = $("#quiz-panel");
  panel.classList.remove("hidden");
  $("#quiz-title").textContent = `${chapter.title} · ${T("اختبر نفسك", "Check yourself")}`;
  $("#quiz-score").textContent = T("٥ أسئلة فهم، من غير كتابة كود. الاختبار مش بيقفل أي مهمة.", "5 understanding questions, no code to write. This quiz never locks missions.");
  const questions = $("#quiz-questions");
  questions.replaceChildren();
  chapter.quiz.forEach((question, index) => questions.appendChild(quizQuestion(question, index)));
  $("#quiz-retry").onclick = () => openQuiz(chapter);
  panel.scrollIntoView({ block: "start" });
  $("#quiz-title").focus();
}
function quizQuestion(question, index) {
  const text = LANG === "en" ? question.en : question;
  const field = document.createElement("fieldset");
  const legend = document.createElement("legend");
  legend.textContent = `${index + 1}. ${text.question}`;
  field.appendChild(legend);
  const feedback = document.createElement("p");
  feedback.className = "quiz-feedback";
  feedback.setAttribute("role", "status");
  text.options.forEach((option, choice) => {
    const button = document.createElement("button");
    button.className = "tool quiz-option";
    button.textContent = option;
    button.onclick = () => answerQuiz(field, question, index, choice);
    field.appendChild(button);
  });
  field.appendChild(feedback);
  return field;
}
function answerQuiz(field, question, index, choice) {
  if (quizAttempt.answers.has(index)) return;
  const correct = choice === question.answer;
  quizAttempt.answers.set(index, correct);
  field.querySelectorAll("button").forEach((button, i) => {
    button.disabled = true;
    if (i === question.answer) button.classList.add("quiz-correct");
    else if (i === choice) button.classList.add("quiz-wrong");
  });
  const text = LANG === "en" ? question.en : question;
  const feedback = field.querySelector(".quiz-feedback");
  feedback.textContent = `${correct ? T("صح!", "Correct!") : T("مش صح.", "Incorrect.")} ${text.explanation}`;
  if (quizAttempt.answers.size === quizAttempt.chapter.quiz.length) finishQuiz();
}
function finishQuiz() {
  const score = [...quizAttempt.answers.values()].filter(Boolean).length;
  const chapter = quizAttempt.chapter.n;
  store.blob.quizScores[chapter] = Math.max(store.blob.quizScores[chapter] ?? 0, score);
  store.save();
  $("#quiz-score").textContent = `${T("نتيجتك", "Your score")}: ${score}/5 · ${T("أفضل نتيجة", "Best score")}: ${store.blob.quizScores[chapter]}/5`;
  $(`[data-quiz-chapter="${chapter}"]`).innerHTML = quizEntryText(chapter);
}
function quizEntryText(chapter) {
  const best = store.blob.quizScores[chapter];
  return `<span class="cell-title">${T("اختبر نفسك", "Check yourself")}</span>${best === undefined ? "" : `<span class="quiz-badge">${T("أفضل نتيجة", "Best score")} <bdi>${esc(String(best))}/5</bdi></span>`}`;
}

const CHAPTER_UNITS = ["#d39b55", "#9fd2cb", "#93c45c", "#d75c9f", "#a45f3e", "#45b3a8", "#f2b33d", "#c7a2ff"];
function renderMissions() {
  const list = $("#chapter-list");
  list.innerHTML = "";
  if (store.blob.review?.due) {
    const card = document.createElement("section");
    card.className = "review-card";
    const mission = S.levels.find((l) => l.id === store.blob.review.level);
    card.innerHTML = `<h3>${T("مراجعة", "Review")}</h3><p>${esc(mission?.title || "")} · ${T("جرّب من ذاكرتك ببيانات جديدة. كود البداية فاضي إلا في المختبرات. النجوم والتقدّم زي ما هم.", "Try from memory with a fresh seed. Starters are empty except in labs. Stars and progress stay the same.")}</p>`;
    for (const [action, label] of [["start", T("ابدأ", "Start")], ["skip", T("تخطّي", "Skip")]]) {
      const button = document.createElement("button");
      button.className = "tool"; button.textContent = label; button.disabled = S.running;
      button.onclick = async () => {
        button.disabled = true;
        try {
          await store.push();
          const response = await api(`/api/profiles/${PROFILE}/review/${action}`, {});
          if (action === "start") { closeOverlay($("#missions")); await loadLevel(response.id, response); }
          else { store.blob.review = response; renderMissions(); }
        } catch (e) { log(String(e), "c-err"); }
        finally { button.disabled = false; }
      };
      card.appendChild(button);
    }
    list.appendChild(card);
  }
  const groups = new Map();
  for (const ch of S.chapters) {
    let group = groups.get(ch.track);
    if (!group) {
      group = document.createElement("details");
      group.className = "track-group";
      group.dataset.track = ch.track;
      groups.set(ch.track, group);
      group.open = true;
      const summary = document.createElement("summary");
      summary.textContent = ch.track;
      group.appendChild(summary);
      list.appendChild(group);
    }
    const row = document.createElement("div");
    const levels = S.levels.map((l, i) => [l, i]).filter(([l]) => l.chapter === ch.n);
    const anyOpen = levels.some(([, i]) => unlocked(i));
    row.className = "quad-row" + (anyOpen ? "" : " locked");
    const unit = CHAPTER_UNITS[(ch.n - 1) % CHAPTER_UNITS.length];
    row.style.setProperty("--unit", unit);
    row.style.setProperty("--fill", `radial-gradient(#0000002e 1px, transparent 1.5px) 0 0 / 6px 6px, ${unit}`);
    row.innerHTML = `<div class="quad-name"><span class="quad-num">${ch.n}</span><div><b>${esc(ch.title)}</b><small>${esc(ch.concept)}</small></div></div>`;
    const cells = document.createElement("div");
    cells.className = "quad-cells";
    for (const [l, i] of levels) {
      const b = document.createElement("button");
      const r = rec(l.id), open = unlocked(i);
      b.className = "cell" + (r.stars > 0 ? " done" : "") + (S.cur && l.id === S.cur.id ? " current" : "");
      b.disabled = !open || S.running;
      b.innerHTML = `<span class="cell-top"><span class="cell-code">${sheetCode(l)}</span>${open ? `<span class="stars">${starHTML(r.stars)}</span>` : ""}</span><span class="cell-title">${esc(l.title)}</span>${l.boss ? `<span class="boss-badge">${T("ختام المسار", "Track Boss")}</span>` : ""}`;
      b.title = open ? (l.kind === "lab" ? T("مختبر", "Lab") : T("مهمة", "Mission")) : T("مقفولة: خلّص اللي قبلها", "Locked: finish the one before it");
      b.onclick = () => { closeOverlay($("#missions")); loadLevel(l.id); };
      b.dataset.search = `${l.id} ${sheetCode(l)} ${l.title}`.toLowerCase();
      b.dataset.code = sheetCode(l);
      cells.appendChild(b);
    }
    if (ch.quiz?.length) {
      const quiz = document.createElement("button");
      quiz.className = "cell quiz-entry";
      quiz.dataset.quizChapter = ch.n;
      quiz.innerHTML = quizEntryText(ch.n);
      quiz.onclick = () => openQuiz(ch);
      cells.appendChild(quiz);
    }
    row.appendChild(cells);
    group.appendChild(row);
  }
  if (typeof applyMissionFilter === "function") applyMissionFilter(); // keep the search while re-rendering
  requestAnimationFrame(() => list.querySelector(store.blob.review?.due ? ".review-card" : ".cell.current")?.scrollIntoView({ block: "center" }));
}
function showOverlay(selector) {
  const overlay = $(selector);
  if (overlay.open) return;
  overlay.classList.remove("hidden");
  overlay.showModal();
}
function closeOverlay(overlay) {
  overlay.close();
  overlay.classList.add("hidden");
}
// Filter the index by mission id, sheet code, title or chapter name.
function applyMissionFilter() {
  const query = $("#mission-search").value.trim().toLowerCase();
  // "3" means chapter 3 and "3-1" means that sheet, not 13-1 or 23-1
  const numeric = /^\d+(-\d+)?$/.test(query);
  let shown = 0;
  for (const row of document.querySelectorAll("#chapter-list .quad-row")) {
    const chapterHit = !query || (!numeric && row.querySelector(".quad-name").textContent.toLowerCase().includes(query));
    let any = false;
    for (const cell of row.querySelectorAll(".cell")) {
      const code = cell.dataset.code || "";
      const hit = chapterHit || (numeric ? code === query || code.startsWith(query + "-") : (cell.dataset.search || "").includes(query));
      cell.hidden = !hit;
      any ||= hit;
    }
    row.hidden = !any;
    if (any) shown++;
  }
  for (const group of document.querySelectorAll("#chapter-list .track-group")) {
    group.hidden = !group.querySelector(".quad-row:not([hidden])");
    if (query) group.open = true;
  }
  $("#mission-search-empty").hidden = shown > 0;
}
$("#mission-search").addEventListener("input", applyMissionFilter);
// Enter opens the first open mission left after filtering
$("#mission-search").addEventListener("keydown", (e) => {
  if (e.key !== "Enter" || !e.target.value.trim()) return;
  const first = [...document.querySelectorAll("#chapter-list .cell[data-code]")].find((cell) => !cell.hidden && !cell.closest("[hidden]") && !cell.disabled);
  if (first) { e.preventDefault(); first.click(); }
});

function cycleOverlayFocus(e) {
  if (e.key !== "Tab") return;
  const controls = [...e.currentTarget.querySelectorAll("button, a[href], input, textarea, select, [tabindex]")]
    .filter((el) => !el.disabled && el.tabIndex >= 0 && el.getClientRects().length);
  const first = controls[0], last = controls.at(-1);
  const outsideControls = !controls.includes(document.activeElement);
  const target = e.shiftKey && (outsideControls || document.activeElement === first) ? last
    : !e.shiftKey && (outsideControls || document.activeElement === last) ? first : null;
  if (target) { e.preventDefault(); target.focus(); }
}
document.querySelectorAll("dialog.overlay").forEach((overlay) => {
  overlay.addEventListener("keydown", cycleOverlayFocus);
  overlay.addEventListener("close", () => {
    if (!overlay.open) overlay.classList.add("hidden");
  });
});
async function refreshLearningStats() {
  let saved = null;
  try {
    saved = await api(`/api/profiles/${PROFILE}/metrics`);
  } catch (error) {
    if (!(error instanceof Error)) throw error;
    // The profile's local save remains usable while the localhost API is offline.
  }
  if (saved) {
    store.blob.metrics = LearningMetrics.mergeMetrics(store.blob.metrics, saved);
    store.save();
  }
  LearningMetrics.renderStats({ struggles: $("#metrics-struggles"), chapters: $("#metrics-chapters") }, store.blob.metrics,
    { levels: S.levels, chapters: S.chapters });
}

$("#btn-missions").onclick = () => { renderMissions(); showOverlay("#missions"); };
$("#btn-stats").onclick = async () => {
  closeOverlay($("#missions"));
  showOverlay("#metrics");
  LearningMetrics.renderStats({ struggles: $("#metrics-struggles"), chapters: $("#metrics-chapters") }, store.blob.metrics,
    { levels: S.levels, chapters: S.chapters });
  await refreshLearningStats();
};
$("#btn-metrics-export").onclick = () => {
  const blob = new Blob([JSON.stringify(LearningMetrics.normalizeMetrics(store.blob.metrics), null, 2)], { type: "application/json" });
  const link = document.createElement("a");
  link.href = URL.createObjectURL(blob);
  link.download = "orbital-learning-metrics.json";
  document.body.appendChild(link);
  link.click();
  link.remove();
  const url = link.href;
  setTimeout(() => URL.revokeObjectURL(url), 0);
};
function countLessonInteraction(event) {
  const target = event.target.closest("button, select, input");
  if (!target || !target.closest(".lesson-walkthrough, .lesson-sandbox")) return;
  if (event.type === "click" && target.tagName !== "BUTTON") return;
  if (event.type === "change" && target.tagName === "BUTTON") return;
  if (S.cur) LearningMetrics.recordInteraction(store, S.cur.id);
}
$("#lesson").addEventListener("click", countLessonInteraction);
$("#lesson").addEventListener("change", countLessonInteraction);
document.querySelectorAll("[data-close]").forEach((b) => (b.onclick = () => closeOverlay(b.closest("dialog"))));

// ---------- run ----------
// One run at a time. Stop (or switching missions) cancels it right away:
// the request is aborted, which makes the server kill the build or program.
let runId = 0, runCtl = null;
const RUN_TIMEOUT = 150000; // longer than the server's own build + run limits
function clearTraceRun() {
  InlineDebugger.clear();
}
function showTraceRun(level, response, payload, sources) {
  if (!response.traceable) {
    if (level.kind === "lab" && response.phase === "lab") InlineDebugger.unavailable(T("المختبر ده مفيهوش main function للتصحيح.", "This lab has no main function to debug."));
    return;
  }
  const files = new Map(levelFiles(level).map((file) => [file.name, file]));
  if (!Object.entries(sources).every(([name, code]) => files.has(name) && fileCode(level, files.get(name)) === code)) return;
  const tracePayload = { level: payload.level };
  if (level.kind === "lab") tracePayload.files = { ...payload.files };
  else tracePayload.code = payload.code;
  InlineDebugger.start(level.id, tracePayload);
}
// ---------- debug playback ----------
// For missions that return a single-run trace, Run plays the recorded steps:
// the debugger walks the code while the map shows the world at each step.
// The result appears when the last step is reached (or on skip/exit).
let dbgRun = null;
function startDebugRun(level, r, paused) {
  const list = r.events || [];
  const initAt = list.findIndex((e) => e.op === "init");
  if (initAt >= 0) setWorld(list[initAt]);
  dbgRun = { level, r, list, base: initAt + 1, applied: initAt + 1, revealed: false };
  W.live = true; W.liveStep = 220;
  $("#btn-skip").disabled = false;
  InlineDebugger.startWith(level.id, r.trace, {
    autoplay: !paused,
    onFrame: (frame, index, total) => { seekMap(frame.events); mapBadge(index, total); },
    onEnd: () => revealDebugRun(false),
    onExit: () => revealDebugRun(false),
    onClear: dropDebugRun,
  });
}
function quietly(fn) {
  const realLog = log, realPlay = Sound.play, realBurst = burst, fxCount = fx.length;
  log = () => {}; Sound.play = () => {}; burst = () => {};
  try { fn(); } finally { log = realLog; Sound.play = realPlay; burst = realBurst; fx.length = fxCount; }
}
function seekMap(target) {
  const run = dbgRun;
  if (!run || !W) return;
  const n = Math.max(run.base, Math.min(Number(target) || 0, run.list.length));
  if (n < run.applied) {
    // going back: rebuild the world up to the step without sound or logs
    if (run.base > 0) setWorld(run.list[run.base - 1]);
    W.live = true; W.liveStep = 220;
    quietly(() => { for (let i = run.base; i < n; i++) apply(run.list[i]); });
    for (const d of W.drones.values()) { d.px = d.x; d.py = d.y; }
  } else {
    for (let i = run.applied; i < n; i++) apply(run.list[i]);
  }
  run.applied = n;
}
function mapBadge(index, total) {
  const badge = $("#debug-map-badge");
  if (!badge) return;
  if (index == null) { badge.hidden = true; return; }
  badge.hidden = false;
  const icon = InlineDebugger.playing && index < total - 1 ? "▶" : "⏸";
  badge.textContent = T(`${icon} خطوة ${index + 1} / ${total}`, `${icon} step ${index + 1} / ${total}`);
}
function revealDebugRun(silent) {
  const run = dbgRun;
  if (!run || run.revealed) return;
  run.revealed = true;
  $("#btn-skip").disabled = true;
  if (silent) return finish(run.level, run.r.result, run.r.missing || [], true);
  seekMap(run.list.length);
  finish(run.level, run.r.result, run.r.missing || []);
}
function dropDebugRun() {
  const run = dbgRun;
  if (!run) return;
  if (!run.revealed) revealDebugRun(true);
  // leave the map at the run's end state
  if (W && S.cur === run.level) quietly(() => { for (let i = run.applied; i < run.list.length; i++) apply(run.list[i]); });
  if (W) { for (const d of W.drones.values()) { d.px = d.x; d.py = d.y; } W.live = false; }
  dbgRun = null;
  mapBadge(null);
  $("#btn-skip").disabled = true;
}

function runButton(running) {
  const btn = $("#btn-run");
  btn.classList.toggle("stop", running);
  btn.querySelector("use").setAttribute("href", running ? "#i-stop" : "#i-play");
  btn.querySelector("span").textContent = running ? "Stop" : "Run";
  btn.title = running ? T("وقّف الكود (Esc)", "Stop the code (Esc)") : "Ctrl/⌘ + Enter";
}
function stopRun(quiet) {
  if (!S.running) return;
  runId++;
  runCtl?.abort();
  runCtl = null;
  if (playback) { playback = null; $("#btn-skip").disabled = true; }
  for (const d of W.drones.values()) { d.px = d.x; d.py = d.y; }
  S.running = false;
  runButton(false);
  if (!quiet) { log(T("اتوقف الكود.", "Code stopped."), "c-warn"); Sound.play("fail"); }
}
document.addEventListener("keydown", (e) => {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "s") {
    // the game saves as you type; ⌘S/Ctrl+S just confirms it instead of the browser's "Save page"
    e.preventDefault();
    if (!ed.readOnly && S.cur) {
      saveFile(S.cur, S.file, ed.value);
      clearTimeout(saveTimer); saveTimer = null;
      store.save(); store.push();
      $("#save-state").textContent = T("محفوظ", "Saved");
    }
    return;
  }
  const typing = e.target.closest?.("input, textarea, select, [contenteditable]");
  if (e.key === "?" && !typing && !e.metaKey && !e.ctrlKey && !document.querySelector("dialog[open]")) {
    e.preventDefault();
    showOverlay("#shortcuts");
    return;
  }
  if (e.key === "Escape" && S.running && S.mode !== "colony" && !document.querySelector("dialog[open]")) stopRun();
});

async function run(opts = {}) {
  if (S.mode === "colony") { clearTraceRun(); return Colony.toggleRun(); }
  if (S.running) return stopRun();
  clearTraceRun();
  previewVersion++;
  if (playback) finishPlayback();
  const me = ++runId;
  const ctl = runCtl = new AbortController();
  const timer = setTimeout(() => ctl.abort("timeout"), S.cur.kind === "lab" ? 30 * 60 * 1000 : RUN_TIMEOUT);
  S.running = true;
  S.failLogged = false;
  runButton(true);
  Sound.play("build");
  $("#banner").classList.add("hidden");
  con.innerHTML = "";
  const level = S.cur;
  try {
    if (isLinked()) await syncFromWorkspace(true);
    else if (!ed.readOnly) saveFile(S.cur, S.file, ed.value);
    store.save();
    if (me !== runId) return;
    const isLab = level.kind === "lab";
    const payload = { level: level.id, review: !!S.review };
    if (isLab) {
      payload.files = {};
      for (const f of levelFiles(level)) if (!f.readOnly) payload.files[f.name] = fileCode(level, f);
    } else payload.code = fileCode(level, levelFiles(level)[0]);
    const traceSources = Object.fromEntries(
      levelFiles(level).filter((f) => f.name.endsWith(".go")).map((f) => [f.name, fileCode(level, f)]),
    );
    const r = await api("/api/run", payload, null, ctl.signal);
    if (me !== runId || S.cur !== level) return;
    LearningMetrics.recordRun(store, level.id, r);
    if (r.phase === "require") {
      finish(level, r.result, r.missing || []);
      return;
    }
    if (r.phase === "compile") {
      lastErrors = r.compileErrors || [];
      log(r.result.message, "c-err");
      Sound.play("error");
      if (lastErrors.length) {
        for (const e of lastErrors) {
          logCompileError(e);
          logErrorHelp(e.tip, e.tipCause, e.tipCheck);
        }
        // land on the first error so the fix starts right there
        const line = Number(lastErrors[0].line);
        const file = firstErrorFile(lastErrors, levelFiles(level), S.file);
        if (line > 0 && !ed.readOnly) jumpToLine(file, line);
        else showFile(file);
      } else log(r.compileOutput, "c-err");
      if (r.runtimeTip) logErrorHelp(r.runtimeTip, r.runtimeTipCause, r.runtimeTipCheck);
      banner(r.result.message, false);
      return;
    }
    lastErrors = [];
    showTraceRun(level, r, payload, traceSources);
    if (r.phase === "lab") {
      refreshService(level);
      for (const st of r.steps || []) {
        log(`${st.ok ? T("نجح:", "Passed:") : T("فشل:", "Failed:")} ${st.title}`, st.ok ? "c-ok" : "c-err");
        if (!st.ok && st.output) st.output.split("\n").slice(0, 40).forEach((l) => logOutputLine("   " + l, "c-err"));
        if (!st.ok) logErrorHelp(st.tip, st.tipCause, st.tipCheck);
      }
      finish(level, r.result, r.missing || []);
      return;
    }
    if (r.stdout) { log("fmt output:", "c-info"); r.stdout.trimEnd().split("\n").forEach((l) => log(l)); }
    if (r.stderr) r.stderr.trimEnd().split("\n").forEach((l) => logOutputLine(l, "c-err"));
    if (r.runtimeTip) logErrorHelp(r.runtimeTip, r.runtimeTipCause, r.runtimeTipCheck);
    log(T("الدرون شغال…", "Drone running…"), "c-info");
    if (r.trace?.available && r.trace.synced && r.trace.frames?.length) {
      startDebugRun(level, r, !!opts.paused);
      return;
    }
    await new Promise((res) => startPlayback(r.events, res));
    if (me !== runId) return;
    finish(level, r.result, r.missing || []);
  } catch (e) {
    if (me !== runId) return; // stopped by the player
    if (ctl.signal.aborted) {
      log(T("السيرفر اتأخر جداً في الرد، فوقّفت الكود. جرّب تاني.", "The server took too long to answer, so the run was stopped. Try again."), "c-err");
      logErrorHelp(T("التشغيل ماخلصش قبل المهلة.", "The run did not finish before its time limit."), T("ممكن يكون فيه انتظار مفتوح أو عملية بطيئة.", "The program may be waiting indefinitely or running slowly."), T("راجع آخر output ونقاط الانتظار.", "Check the last output and the places where the program waits."));
    }
    else { log(T("خطأ في الاتصال بالسيرفر: ", "Server connection error: ") + e.message, "c-err"); logServerHint(e); }
  } finally {
    clearTimeout(timer);
    if (me === runId) {
      S.running = false;
      runCtl = null;
      runButton(false);
    }
  }
}
$("#btn-run").onclick = () => run();
$("#btn-step").onclick = () => run({ paused: true });
$("#btn-vscode").onclick = openVSCode;
$("#btn-unlink").onclick = () => {
  linked.delete(S.cur.id); saveLinked();
  showFile(S.file); linkUI();
  $("#save-state").textContent = "";
  log(T("اتفك الربط. الكود اللي اتكتب في VS Code متسجل هنا، وتقدر تكمّل في المحرر ده.", "Unlinked. The code you wrote in VS Code is saved here; keep going in this editor."), "c-info");
};

function banner(text, ok) {
  const b = $("#banner");
  b.textContent = text; b.className = "banner" + (ok ? " ok" : "");
}

// silent: record the result without any UI (the learner left the level or
// edited the code while a debug playback was still showing the run).
function finish(level, res, missing = [], silent = false) {
  if (!silent && S.cur !== level) return;
  if (silent && (!res.success || S.review)) return;
  if (!res.success) {
    banner(res.message + (res.progress ? ` (${res.progress})` : "") + (missing.length ? " " + missing.join(" ") : ""), false);
    missing.forEach((m) => log("   • " + m, "c-tip"));
    if (!S.failLogged) { log(res.message, "c-err"); Sound.play("error"); }
    if (res.progress) log("   " + res.progress, "c-tip");
    return;
  }
  if (S.review) {
    store.blob.review = { ...store.blob.review, due: false };
    banner(T("المراجعة اكتملت! النجوم والتقدّم زي ما هم.", "Review complete! Stars and progress stay the same."), true);
    Sound.play("goal");
    return;
  }
  const r = rec(level.id);
  const improved = res.stars > r.stars || (r.best === 0 || res.ticks < r.best);
  r.stars = Math.max(r.stars, res.stars);
  r.best = r.best ? Math.min(r.best, res.ticks) : res.ticks;
  store.save();
  if (silent) { updateOverall(); return; }
  $("#lvl-stars").innerHTML = starHTML(r.stars);
  $("#lvl-stars").setAttribute("aria-label", T(`${r.stars} من 3 نجوم`, `${r.stars} of 3 stars`));
  banner(res.message, true);

  const i = idx(level.id);
  const next = S.levels[i + 1];
  $("#result-stars").innerHTML = starHTML(res.stars);
  $("#stamp-code").textContent = `SHEET ${sheetCode(level)}`;
  updateOverall();
  Sound.play("goal");
  for (let k = 0; k < res.stars; k++) setTimeout(() => Sound.play("star"), 420 + k * 180);
  let text = level.kind === "lab" ? res.message : T(`خلصت في <b>${res.ticks}</b> tick.`, `Done in <b>${res.ticks}</b> ticks.`);
  if (res.par > 0) {
    if (res.stars < 3) text += T(`<br>النجوم التلاتة محتاجة <b>${res.gold}</b> tick أو أقل. تقدر توصلها؟`, `<br>Three stars need <b>${res.gold}</b> ticks or fewer. Can you get there?`);
    else text += T("<br>حل مثالي. مفيش أحسن من كده.", "<br>Perfect solution. It doesn't get better than this.");
  }
  if (improved && r.best === res.ticks && res.ticks > 0) text += T(`<br>أحسن رقم ليك: <b>${r.best}</b>`, `<br>Your best: <b>${r.best}</b>`);
  if (level.boss) text += `<br><span class="bridge-badge">${T("المسار اكتمل!", "Track complete!")}</span>`;
  for (const project of level.bridge || []) {
    const number = String(project).padStart(3, "0");
    text += `<br><span class="bridge-badge">${esc(T(`جاهز لمشروع ${number}`, `Ready for project ${number}`))}</span>`;
  }
  $("#result-text").innerHTML = text;
  const nb = $("#btn-next");
  nb.classList.toggle("hidden", !next);
  nb.onclick = () => { closeOverlay($("#result")); loadLevel(next.id); };
  showOverlay("#result");
}

// ---------- boot ----------
(async function boot() {
  try {
    const [data] = await Promise.all([api("/api/levels"), store.load().then((p) => { S.progress = p; })]);
    S.chapters = data.chapters.map(localizeChapter); S.levels = data.levels.map(localizeLevel);
    for (const id of store.pref("linked") || []) linked.add(id);
    Profiles.init();
  } catch (e) {
    log(T("مقدرتش أوصل للسيرفر: ", "Can't reach the server: ") + e.message, "c-err");
    logServerHint(e);
    return;
  }
  const linkedID = hashLevelID();
  let id = idx(linkedID) >= 0 && unlocked(idx(linkedID)) ? linkedID : store.pref("current");
  if (idx(id) < 0 || !unlocked(idx(id))) {
    const firstOpen = S.levels.findIndex((l, i) => unlocked(i) && rec(l.id).stars === 0);
    id = S.levels[firstOpen >= 0 ? firstOpen : 0].id;
  }
  restoreSpeeds();
  updateOverall();
  loadLevel(id);
  if (!store.pref("welcomed")) showWelcome();
})();

// The map speed and the debugger speed are remembered per player.
function restoreSpeeds() {
  const map = Number(store.pref("mapSpeed"));
  if (map >= +speedEl.min && map <= +speedEl.max) speedEl.value = map;
  speedEl.addEventListener("change", () => store.pref("mapSpeed", +speedEl.value));
  const select = $("#debugger-speed");
  const saved = String(store.pref("debugSpeed") || "");
  if (select && [...select.options].some((o) => o.value === saved)) {
    select.value = saved;
    select.dispatchEvent(new Event("change"));
  }
  select?.addEventListener("change", () => store.pref("debugSpeed", select.value));
}

// Opening a #mission link in a running tab jumps to it (if unlocked).
// "#1-2" (sheet code) or "#1-2-some-title" (full id) both name a mission.
function hashLevelID() {
  const raw = decodeURIComponent(location.hash.slice(1));
  return S.levels.find((l) => l.id === raw)?.id || S.levels.find((l) => sheetCode(l) === raw)?.id || raw;
}
window.addEventListener("hashchange", () => {
  const id = hashLevelID();
  if (id && id !== S.cur?.id && unlocked(idx(id))) loadLevel(id);
});

function showWelcome() {
  showOverlay("#welcome");
  requestAnimationFrame(() => drawSurveyArt($("#welcome-map")));
}
$("#btn-start").onclick = () => {
  store.pref("welcomed", true);
  closeOverlay($("#welcome"));
  Sound.play("goal");
  ed.focus();
};
$("#brand").onclick = (e) => { e.preventDefault(); showWelcome(); };
$("#btn-lesson").onclick = () => {
  const open = document.body.classList.toggle("lesson-open");
  $("#btn-lesson").setAttribute("aria-expanded", String(open));
};
