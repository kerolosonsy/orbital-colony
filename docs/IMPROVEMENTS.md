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
