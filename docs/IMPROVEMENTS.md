# Improvement sweep (branch `improve-sweep`)

Each entry: what changed and why. Numbered in the order they landed.

1. **Removed dead error-tip helpers** — `compileTips`, `runTips`, `legacyTipList`, `tipFor` in `cmd/orbital/errtips.go` were left over from the old tip list and never called (staticcheck U1000). Also removed unused `ch01Variable` in `colony/levels/ch01_lessons.go`.
2. **Raw-string regexes in lesson tests** — ch12/21/26/27 lesson tests used double-escaped regex strings; switched to raw strings so they read like the pattern they match (staticcheck S1007).
3. **Auto-closing brackets and quotes in the editor** — typing `(`, `[`, `{`, `"` or `` ` `` adds its closer; typing the closer when it is already next just steps over it. Quotes don't pair after a word character, and nothing pairs in the middle of a word. Logic lives in the new pure module `web/editkeys.js` (tested in `web/editkeys.test.js`).
4. **Enter between a pair opens a block** — `{|}` + Enter gives the `{`, an indented empty line with the caret, and `}` on its own line (also for `(` and `[`).
5. **Indent after `case x:` / `default:`** — Enter after a switch case label indents the body like gofmt.
6. **Backspace deletes an empty pair** — `(|)` + Backspace removes both characters.
7. **Block indent / outdent** — Tab with a multi-line selection indents every selected line; Shift+Tab outdents the current line or selection (tab or up to 4 spaces). Caret stays on the same code.
8. **Ctrl/⌘+/ toggles `//` comments** on the current line or all selected lines (blank lines untouched).
9. **Selection wrapping** — with a single-line selection, typing `(`/`[`/`{`/`"` wraps it instead of replacing it.
10. **Editor edits are undoable** — every edit above goes through `execCommand` with the exact range, so Ctrl/⌘+Z undoes it like normal typing; the read-only editor (debugger) ignores them.
11. **Security headers on every response** — `X-Frame-Options: DENY` + `frame-ancestors 'none'` (another site can't frame the page that holds the API token and click Run for you), `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`. Test: `TestGuardSetsSecurityHeaders`.
12. **Map speed is remembered** — the playback speed slider is saved per player and restored on the next visit.
13. **Debugger speed is remembered** — the 0.5×/1×/2×/4× choice is saved per player too.
14. **Local save can't break editing** — `store.save()` wraps `localStorage.setItem` in try/catch, so a full or blocked storage no longer throws on every keystroke; the server save still runs.
15. **F9 toggles a breakpoint on the caret line** (VS Code habit), without reaching for the mouse.
16. **Debugger buttons show their shortcut** — tooltips now say F10 / F11 / Shift+F11 / F5 / Shift+F5.
17. **Keyboard shortcuts sheet** — press `?` (outside the editor) to see every shortcut in Arabic or English.
18. **Browser tab shows the mission** — `document.title` becomes `1-2 · <title> — Orbital Colony`, so several tabs/bookmarks are tellable apart.
19. **Mission links** — the URL hash follows the open mission (`#1-2-…`); opening `#1-2` (sheet code) or the full id jumps straight to it when unlocked, also in an already-open tab.
20. **Colony tab title fixes** — the mission title is re-captured on every colony entry (it used to restore the title from the very first entry), and the open colony shows `Colony — Orbital Colony`.
21. **Copy console button** — one click copies the whole console (errors, output) as plain text, with a short "Copied ✓" confirmation.
22. **Current-line highlight** — a faint band marks the caret's line while editing (hidden while selecting or when the editor is read-only during debugging).
