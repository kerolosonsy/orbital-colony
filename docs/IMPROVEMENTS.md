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
23. **"Back to starter" is undoable** — the reset goes through the edit history (Ctrl/⌘+Z restores your code) and the console says so; before, it overwrote the textarea and the code was gone.
24. **Arrow keys step through the lesson** — ← / → move to the previous/next lesson step (mirrored for Arabic RTL) when focus is in the lesson and not on a control, sandbox or walkthrough.
25. **README covers the new features** — debugger, error explanations, playgrounds/stats, editor helpers, a shortcuts table and mission links.
26. **`make test` duration is honest** — help text said ~4 minutes; the full suite takes ~15–20.
27. **"Server is down" hint** — when the page can't reach the server (fetch fails), the console now says to run `make start` and refresh, instead of only "Failed to fetch".
28. **"Page is stale after restart" hint** — every restart makes a new API token, so an old tab got `403 missing token` forever. The console now says to refresh the page (Arabic and English messages both detected). Verified live by restarting the server under an open tab.
29. **Progress upload failures are no longer silent** — if saving progress to the server fails, the console says why once (and the hint is de-duplicated when a run fails at the same moment). The local copy keeps the progress meanwhile.
30. **Static files revalidate** — `/static/*` is served with `Cache-Control: no-cache` and a content-hash `ETag`, so a rebuilt server is always picked up on refresh while unchanged files come back as cheap `304`s. Test: `TestStaticFilesRevalidate`.
31. **Clean file paths in panics on macOS** — panic traces showed `/privatemain.go:8` because the binary records the resolved temp path (`/private/var/...`) while the server only stripped `/var/...`. New `stripDir` removes both spellings for build, test, lab, colony and panic output. Test: `TestStripDirRemovesResolvedTempPath`; verified with a real panic in the browser.
32. **Click a panic/test line to jump to the code** — console lines that name a mission file (`main.go:8 +0x1c`, `main_test.go:23: …`) are clickable and select that line in the editor, like compile errors already were.
33. **Move lines with Alt+↑/↓** — the current line (or every selected line) swaps with the line above/below, caret/selection following it.
34. **Duplicate lines with Shift+Alt+↑/↓** — copies the touched lines above or below. Both are in the shortcuts sheet and README, tested in `web/editkeys.test.js`.
35. **`make doctor` checks lab libraries and Docker** — it now says whether `labdeps` builds offline (else: run `make deps` once online) and whether Docker is there for the Postgres/MongoDB/Kafka missions.
36. **`make test-web`** — runs the browser-side tests (`node --test`, 278 tests), and `make check` now includes them; before, no make target ran the web tests at all.
37. **Caret lands on the first compile error** — after a failed build the editor opens the right file and selects the first error's line, so the fix starts there (clicking other errors still works).
38. **Bracket matching** — with the caret next to `(`/`[`/`{` or its closer, both brackets get a small box; brackets inside strings, runes and comments are ignored. Matcher `EditKeys.matchBracket` is tested; alignment checked in the browser with tabs.
39. **Favicon** — an inline SVG planet icon, so the tab is recognisable and the browser stops requesting a missing `/favicon.ico` on every load. (Bracket boxes also re-measure the character width once the code font finishes loading.)
40. **Enter goes to the next mission** — the result sheet focuses "Next mission" (it focused "Improve it" before), so finishing a mission and pressing Enter moves on.
41. **Welcome sheet mentions `?`** — new players learn about the shortcuts sheet on the first screen.
42. **Doc comments for the whole learner API** — 71 exported functions/methods and 9 exported types in `orbital/colony` (power, cargo, battery, buildings, melt/sift, dig, data sets, every `*Drone` method) had no doc comment, so hovering them in VS Code (gopls) showed nothing. Each now has a short, accurate comment checked against the simulator (e.g. `Charge` needs a refinery and the base tile and uses one fuel; building methods need the drone standing on the building; free calls are marked "no ticks").
43. **Doc comments for `labkit`** — the lab test helpers (`Clock`/`FakeClock`/`Ticker`, `Logs`/`CSV`/`Text` generators, `CA`/`CertOptions`/`TestCA`) are now documented, so labs show what each helper does on hover.
44. **⌘S / Ctrl+S saves** — instead of opening the browser's "Save page as…" dialog (a common reflex), it saves the code immediately, pushes progress to the server and shows "Saved". Listed in the shortcuts sheet.
45. **Smart Home key** — Home jumps to the first code character of the line, a second press to column 0 (Shift+Home selects), like code editors. Tested.
46. **Search in the mission index** — a search box (focused when the index opens) filters the 181 missions by sheet code, id, title or chapter name. Numbers are exact: `3-1` shows only sheet 3-1 (not 13-1/23-1) and `3` shows chapter 3; empty results say so. Verified in the browser.
47. **Enter in the index search opens the mission** — type `2-3` + Enter to jump straight to it (first unlocked match).
48. **Ctrl/⌘+] and Ctrl/⌘+[** indent / outdent the current line or selection without losing the caret (the VS Code keys; Tab/Shift+Tab still work). Tested; listed in the shortcuts sheet.
49. **Debugger keys that work on laptops** — F10/F11 need Fn on most Mac keyboards, so while debugging Alt+← steps back, Alt+→ steps forward and Alt+Space plays/pauses. The Step back button finally has a shortcut too. Verified live (1/12 → 3/12 → 2/12).
50. **Caret line number lights up** — the gutter number of the line you're on is brighter, matching the current-line band (cleared during debugging, where the debugger marks its own line).
51. **Console doesn't yank you to the bottom** — new output only auto-scrolls when you're already at the bottom, so you can read an earlier error while a long run keeps printing.
52. **Esc closes the lesson drawer** on narrow screens and returns focus to the code.
