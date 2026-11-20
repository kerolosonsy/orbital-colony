"use strict";

// QuickFix: one-click fixes for the two most common import errors.
// fixFor(message) describes the fix; apply(code, fix) returns the new code
// (or null when there is nothing to change).
(function () {
  // standard packages a learner reaches for; others stay manual
  const STD = new Set(["fmt", "strings", "strconv", "math", "sort", "errors", "time", "sync", "os", "slices", "maps", "unicode", "bytes", "bufio", "io", "context", "regexp", "json"]);
  const PATHS = { json: "encoding/json" };

  function fixFor(message) {
    let m = /^undefined: (\w+)$/.exec(String(message).trim());
    if (m && STD.has(m[1])) return { kind: "add", path: PATHS[m[1]] || m[1] };
    m = /^"([\w./-]+)" imported and not used/.exec(String(message).trim());
    if (m) return { kind: "remove", path: m[1] };
    return null;
  }

  // imports finds the import declaration (without its trailing newline)
  function imports(code) {
    const block = /^import \(\n([\s\S]*?)\n\)$/m.exec(code);
    if (block) return { match: block, paths: block[1].split("\n").map((l) => l.trim()).filter(Boolean) };
    const single = /^import ("[^"]+")$/m.exec(code);
    if (single) return { match: single, paths: [single[1]] };
    return null;
  }

  function render(paths) {
    if (paths.length === 1) return `import ${paths[0]}`;
    return `import (\n${paths.map((p) => "\t" + p).join("\n")}\n)`;
  }

  function apply(code, fix) {
    const quoted = `"${fix.path}"`;
    const found = imports(code);
    if (fix.kind === "add") {
      if (found?.paths.includes(quoted)) return null;
      if (!found) {
        const pkg = /^package \w+\n/m.exec(code);
        if (!pkg) return null;
        const at = pkg.index + pkg[0].length;
        return code.slice(0, at) + `\nimport ${quoted}\n` + code.slice(at);
      }
      const end = found.match.index + found.match[0].length;
      return code.slice(0, found.match.index) + render([...found.paths, quoted].sort()) + code.slice(end);
    }
    if (!found?.paths.includes(quoted)) return null;
    const paths = found.paths.filter((p) => p !== quoted);
    const end = found.match.index + found.match[0].length;
    if (paths.length) return code.slice(0, found.match.index) + render(paths) + code.slice(end);
    // no imports left: drop the declaration and the blank line after it
    return code.slice(0, found.match.index) + code.slice(end).replace(/^\n+/, "");
  }

  const api = { fixFor, apply };
  if (typeof window !== "undefined") window.QuickFix = api;
  if (typeof module !== "undefined") module.exports = api;
})();
