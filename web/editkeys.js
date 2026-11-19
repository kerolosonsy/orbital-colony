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

  function handle(value, start, end, e) {
    const mod = e.metaKey || e.ctrlKey;
    const collapsed = start === end;
    const prev = value[start - 1] || "";
    const next = value[end] || "";

    if (mod && (e.key === "/" || e.code === "Slash")) return toggleComment(value, start, end);
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

  const api = { handle };
  if (typeof window !== "undefined") window.EditKeys = api;
  if (typeof module !== "undefined") module.exports = api;
})();
