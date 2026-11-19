"use strict";

// Complete: autocomplete for the learner package after "colony.".
// Pure helpers; the dropdown lives in app.js.
(function () {
  // context returns where the name being typed after "colony." starts, or null.
  function context(value, caret) {
    const before = value.slice(value.lastIndexOf("\n", caret - 1) + 1, caret);
    if (/\/\//.test(before.replace(/"(?:\\.|[^"\\])*"|`[^`]*`/g, '""'))) return null; // inside a comment
    const match = /(?:^|[^\w.])colony\.(\w*)$/.exec(before);
    if (!match) return null;
    return { start: caret - match[1].length, prefix: match[1] };
  }

  // rank keeps package-level names: prefix matches first, then other matches;
  // within each, functions before constants and types (stable otherwise).
  function rank(symbols, prefix, limit = 12) {
    const p = prefix.toLowerCase();
    const top = [], rest = [];
    for (const symbol of symbols || []) {
      if (symbol.kind === "method") continue;
      const name = symbol.name.toLowerCase();
      if (name.startsWith(p)) top.push(symbol);
      else if (p && name.includes(p)) rest.push(symbol);
    }
    // functions are what people call most, so they lead each group
    const order = { func: 0, const: 1, var: 1, type: 2 };
    const byKind = (a, b) => (order[a.kind] ?? 3) - (order[b.kind] ?? 3);
    return top.sort(byKind).concat(rest.sort(byKind)).slice(0, limit);
  }

  // insertion is the text for a symbol and where the caret goes in it:
  // functions get their parentheses, with the caret inside when they take arguments.
  function insertion(symbol) {
    if (symbol.kind !== "func") return { text: symbol.name, caret: symbol.name.length };
    const takesArgs = !/^\(\s*\)/.test(symbol.signature || "()");
    const text = symbol.name + "()";
    return { text, caret: takesArgs ? text.length - 1 : text.length };
  }

  const api = { context, rank, insertion };
  if (typeof window !== "undefined") window.Complete = api;
  if (typeof module !== "undefined") module.exports = api;
})();
