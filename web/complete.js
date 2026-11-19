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
    if (symbol.kind !== "func" && symbol.kind !== "method") return { text: symbol.name, caret: symbol.name.length };
    const takesArgs = !/^\(\s*\)/.test(symbol.signature || "()");
    const text = symbol.name + "()";
    return { text, caret: takesArgs ? text.length - 1 : text.length };
  }

  // callContext names the colony function whose arguments the caret is in:
  // "colony.Move(colony.Ea|" → "Move". Only simple, same-line calls.
  function callContext(value, caret) {
    const before = value.slice(value.lastIndexOf("\n", caret - 1) + 1, caret);
    if (/\/\//.test(before.replace(/"(?:\\.|[^"\\])*"|`[^`]*`/g, '""'))) return null;
    let depth = 0;
    for (let i = before.length - 1; i >= 0; i--) {
      if (before[i] === ")") depth++;
      else if (before[i] === "(") {
        if (depth === 0) {
          const match = /(?:^|[^\w.])colony\.(\w+)$/.exec(before.slice(0, i));
          return match ? match[1] : null;
        }
        depth--;
      }
    }
    return null;
  }

  // memberContext handles "d." where d came from the colony package:
  // d := colony.Spawn()  →  Spawn returns *Drone  →  list Drone's methods.
  // Also "var d *colony.Drone" and parameters "d *colony.Drone".
  function memberContext(value, caret, symbols) {
    const before = value.slice(value.lastIndexOf("\n", caret - 1) + 1, caret);
    if (/\/\//.test(before.replace(/"(?:\\.|[^"\\])*"|`[^`]*`/g, '""'))) return null;
    const match = /(?:^|[^\w.])([A-Za-z_]\w*)\.(\w*)$/.exec(before);
    if (!match || match[1] === "colony") return null;
    const name = match[1].replace(/[$]/g, "");
    const code = value.slice(0, caret);
    let recv = null;
    const assigned = [...code.matchAll(new RegExp(`\\b${name}\\s*(?:,\\s*\\w+\\s*)?:?=\\s*colony\\.(\\w+)\\(`, "g"))].pop();
    if (assigned) {
      const fn = (symbols || []).find((s) => s.kind === "func" && s.name === assigned[1]);
      recv = /\)\s*\(?\s*\*?(\w+)/.exec(fn?.signature || "")?.[1] || null;
    }
    const typed = [...code.matchAll(new RegExp(`\\b${name}\\s+\\*?colony\\.(\\w+)`, "g"))].pop();
    if (typed && (!assigned || typed.index > assigned.index)) recv = typed[1];
    if (!recv || !(symbols || []).some((s) => s.kind === "method" && s.recv === recv)) return null;
    return { start: caret - match[2].length, prefix: match[2], recv };
  }

  // rankMembers lists one type's methods, prefix matches first.
  function rankMembers(symbols, recv, prefix, limit = 12) {
    const p = prefix.toLowerCase();
    const methods = (symbols || []).filter((s) => s.kind === "method" && s.recv === recv);
    const top = methods.filter((s) => s.name.toLowerCase().startsWith(p));
    const rest = p ? methods.filter((s) => !s.name.toLowerCase().startsWith(p) && s.name.toLowerCase().includes(p)) : [];
    return top.concat(rest).slice(0, limit);
  }

  const api = { context, rank, insertion, callContext, memberContext, rankMembers };
  if (typeof window !== "undefined") window.Complete = api;
  if (typeof module !== "undefined") module.exports = api;
})();
