"use strict";

// EditKeys: small code-editor behaviours for the plain <textarea>.
// handle() is pure: it looks at the text, the selection and the key and
// returns one edit {from, to, text, selStart, selEnd} (replace value[from:to]
// with text, then select selStart..selEnd) or null to let the browser act.
(function () {
  const PAIRS = { "(": ")", "[": "]", "{": "}", '"': '"', "`": "`" };
  const CLOSERS = new Set([")", "]", "}", '"', "`"]);

  function lineStart(value, i) { return value.lastIndexOf("\n", i - 1) + 1; }
  function lineEnd(value, i) {
    const n = value.indexOf("\n", i);
    return n < 0 ? value.length : n;
  }

  // The full lines touched by a selection (a selection ending at column 0
  // does not include that last line).
  function touchedLines(value, start, end) {
    const from = lineStart(value, start);
    let stop = end;
    if (end > start && value[end - 1] === "\n") stop = end - 1;
    return { from, to: lineEnd(value, stop) };
  }

  function mapLines(value, start, end, fn) {
    const { from, to } = touchedLines(value, start, end);
    const text = value.slice(from, to).split("\n").map(fn).join("\n");
    return { from, to, text, selStart: from, selEnd: from + text.length };
  }

  function outdent(line) {
    if (line.startsWith("\t")) return line.slice(1);
    return line.replace(/^ {1,4}/, "");
  }

  function toggleComment(value, start, end) {
    const { from, to } = touchedLines(value, start, end);
    const lines = value.slice(from, to).split("\n");
    const code = lines.filter((line) => line.trim() !== "");
    const commented = code.length > 0 && code.every((line) => /^\s*\/\//.test(line));
    const edit = mapLines(value, start, end, (line) => {
      if (line.trim() === "") return line;
      if (commented) return line.replace(/^(\s*)\/\/ ?/, "$1");
      return line.replace(/^(\s*)/, "$1// ");
    });
    if (start === end) {
      // keep the caret on the same code character
      const delta = edit.text.length - (to - from);
      const caret = Math.max(from, start + delta);
      edit.selStart = edit.selEnd = caret;
    }
    return edit;
  }

  // Alt+↑/↓ moves the touched lines; Shift+Alt+↑/↓ duplicates them.
  function moveOrCopyLines(value, start, end, down, copy) {
    const { from, to } = touchedLines(value, start, end);
    const block = value.slice(from, to);
    if (copy) {
      const text = block + "\n" + block;
      const shift = down ? block.length + 1 : 0;
      return { from, to, text, selStart: start + shift, selEnd: end + shift };
    }
    if (down) {
      if (to >= value.length) return null;
      const nextEnd = lineEnd(value, to + 1);
      const below = value.slice(to + 1, nextEnd);
      const shift = below.length + 1;
      return { from, to: nextEnd, text: below + "\n" + block, selStart: start + shift, selEnd: end + shift };
    }
    if (from === 0) return null;
    const prevStart = lineStart(value, from - 1);
    const above = value.slice(prevStart, from - 1);
    const shift = above.length + 1;
    return { from: prevStart, to, text: block + "\n" + above, selStart: start - shift, selEnd: end - shift };
  }

  function handle(value, start, end, e) {
    const mod = e.metaKey || e.ctrlKey;
    const collapsed = start === end;
    const prev = value[start - 1] || "";
    const next = value[end] || "";

    if (mod && (e.key === "/" || e.code === "Slash")) return toggleComment(value, start, end);
    // Ctrl/⌘+] and Ctrl/⌘+[ indent / outdent the touched lines (VS Code habit)
    if (mod && !e.altKey && (e.key === "]" || e.key === "[")) {
      const edit = mapLines(value, start, end, e.key === "]" ? (line) => (line === "" ? line : "\t" + line) : outdent);
      if (collapsed) {
        const caret = Math.max(edit.from, start + (edit.text.length - (edit.to - edit.from)));
        edit.selStart = edit.selEnd = caret;
      }
      return edit;
    }
    // Home jumps to the first code character, then to column 0
    if (e.key === "Home" && !mod && !e.altKey) {
      const ls = lineStart(value, end);
      const first = ls + value.slice(ls, lineEnd(value, ls)).match(/^\s*/)[0].length;
      const caret = end === first ? ls : first;
      return { from: end, to: end, text: "", selStart: e.shiftKey ? Math.min(start, caret) : caret, selEnd: e.shiftKey ? Math.max(start, caret) : caret };
    }
    if (e.altKey && !mod && (e.key === "ArrowUp" || e.key === "ArrowDown")) {
      return moveOrCopyLines(value, start, end, e.key === "ArrowDown", e.shiftKey);
    }
    if (mod || e.altKey || e.isComposing) return null;

    if (e.key === "Tab") {
      if (e.shiftKey) {
        const edit = mapLines(value, start, end, outdent);
        if (collapsed) {
          const caret = Math.max(edit.from, start - ((edit.to - edit.from) - edit.text.length));
          edit.selStart = edit.selEnd = caret;
        }
        return edit;
      }
      if (!collapsed && value.slice(start, end).includes("\n")) {
        return mapLines(value, start, end, (line) => (line === "" ? line : "\t" + line));
      }
      return { from: start, to: end, text: "\t", selStart: start + 1, selEnd: start + 1 };
    }

    if (e.key === "Enter") {
      const line = value.slice(lineStart(value, start), start);
      const base = line.match(/^\s*/)[0];
      let indent = base;
      if (/[{([]\s*$/.test(line) || /^\s*(case\b.*|default\s*):\s*$/.test(line)) indent += "\t";
      const opener = (line.trimEnd().slice(-1));
      if (collapsed && PAIRS[opener] && PAIRS[opener] === next && opener !== '"' && opener !== "`") {
        const text = "\n" + indent + "\n" + base;
        const caret = start + 1 + indent.length;
        return { from: start, to: end, text, selStart: caret, selEnd: caret };
      }
      const caret = start + 1 + indent.length;
      return { from: start, to: end, text: "\n" + indent, selStart: caret, selEnd: caret };
    }

    if (e.key === "Backspace" && collapsed && PAIRS[prev] && PAIRS[prev] === next) {
      return { from: start - 1, to: start + 1, text: "", selStart: start - 1, selEnd: start - 1 };
    }

    if (e.key === "}" && collapsed) {
      const ls = lineStart(value, start);
      if (/^\t+$/.test(value.slice(ls, start)) && next !== "}") {
        return { from: start - 1, to: start, text: "}", selStart: start, selEnd: start };
      }
    }

    // type over a closer that is already there
    if (CLOSERS.has(e.key) && collapsed && next === e.key) {
      return { from: start, to: start, text: "", selStart: start + 1, selEnd: start + 1 };
    }

    if (PAIRS[e.key]) {
      const close = PAIRS[e.key];
      if (!collapsed && !value.slice(start, end).includes("\n")) {
        const inner = value.slice(start, end);
        return { from: start, to: end, text: e.key + inner + close, selStart: start + 1, selEnd: end + 1 };
      }
      if (!collapsed) return null;
      const quote = e.key === '"' || e.key === "`";
      if (quote && (/[\w\\"`]/.test(prev))) return null;
      if (next && !/[\s)\]},;:]/.test(next)) return null;
      return { from: start, to: end, text: e.key + close, selStart: start + 1, selEnd: start + 1 };
    }
    return null;
  }

  // matchBracket finds the bracket next to the caret (before it first) and
  // its partner, skipping strings and comments. Returns [a, b] or null.
  function matchBracket(value, caret) {
    const OPEN = "([{", CLOSE = ")]}";
    let at = -1;
    for (const i of [caret - 1, caret]) {
      if (i >= 0 && i < value.length && (OPEN + CLOSE).includes(value[i])) { at = i; break; }
    }
    if (at < 0) return null;
    // mark code characters (outside strings, runes and comments)
    const code = new Uint8Array(value.length);
    for (let i = 0; i < value.length; i++) {
      const c = value[i], n = value[i + 1];
      if (c === "/" && n === "/") { while (i < value.length && value[i] !== "\n") i++; continue; }
      if (c === "/" && n === "*") { const e = value.indexOf("*/", i + 2); i = e < 0 ? value.length : e + 1; continue; }
      if (c === '"' || c === "'") { i++; while (i < value.length && value[i] !== c && value[i] !== "\n") { if (value[i] === "\\") i++; i++; } continue; }
      if (c === "`") { const e = value.indexOf("`", i + 1); i = e < 0 ? value.length : e; continue; }
      code[i] = 1;
    }
    if (!code[at]) return null;
    const ch = value[at];
    const forward = OPEN.includes(ch);
    const mate = forward ? CLOSE[OPEN.indexOf(ch)] : OPEN[CLOSE.indexOf(ch)];
    let depth = 0;
    for (let i = at; forward ? i < value.length : i >= 0; i += forward ? 1 : -1) {
      if (!code[i]) continue;
      if (value[i] === ch) depth++;
      else if (value[i] === mate && --depth === 0) return [at, i];
    }
    return null;
  }

  const api = { handle, matchBracket };
  if (typeof window !== "undefined") window.EditKeys = api;
  if (typeof module !== "undefined") module.exports = api;
})();
