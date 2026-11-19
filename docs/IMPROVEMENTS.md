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
26. **`make test` duration is honest** — help text said ~4 minutes; the full suite took ~15–20 (now ~8, see 58).
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
53. **Opened hints stay open** — hints you revealed for a mission are remembered per player and shown again (silently, no sound or second hint metric) when you come back or reload; before, they vanished and had to be clicked again. Review runs keep starting without hints. Verified by reload in the browser; test extended.
54. **File tabs keep your place** — switching between files of a lab (or colony) restores each file's caret and scroll position instead of jumping to the top.
55. **Labs reopen on the file you were editing** — the last open file per multi-file mission is remembered per player (review runs start fresh).
56. **Caret survives a reload** — the caret and scroll of the open file are saved when the page closes and restored on the next load of the same mission (verified: position 40 → reload → 40).
57. **`/api/levels` is encoded once per language** — the 3.6 MB level payload was rebuilt and re-encoded on every page load; it is now cached per language with an `ETag` (`no-cache` + `304` when unchanged). Measured: ~10–20 ms → ~2 ms per request, and 0 bytes on revalidation. Test: `TestLevelsPayloadIsCachedWithETag`.
58. **Level tests 3.4× faster** — `TestLevelsSolvable` runs independent levels in parallel (half the CPUs, `ORBITAL_TEST_PARALLEL` overrides); realtime chapters 10–13 stay serial so CPU load can't change their wall-clock results. Measured on this machine: 882 s → 259 s, all levels passing.
59. **10 more error explanations** — common beginner compiler errors now get the Arabic/English "meaning / likely cause / check" box: `break`/`continue` outside a loop, unused label, missing function body (`{` on its own line), constant overflow, constant index out of bounds, `append` on a non-slice, recursive struct type, division by constant zero, using a no-value call as a value, too many return values. Samples are the exact Go 1.27 messages; priority tests make sure the specific tip wins over the generic `cannot use`/`invalid operation` one.
60. **3 more runtime error explanations** — `makeslice: len out of range` (negative computed length), `sync: unlock of unlocked mutex` (double Unlock), and `strconv…: invalid syntax` (converting non-numeric text, check the error first). Messages taken from real Go 1.27 runs.
61. **Speed slider says what it means** — hovering (and screen readers via `aria-valuetext`) shows "about N ticks per second", updated as you drag.
62. **Clear "port in use" message** — starting the server on a busy port printed a raw `listen tcp … bind: address already in use`; it now says the game may already be running (`make status`) and how to pick another port. Verified against the running preview server.
63. **`make services-up` waits until services are ready** — Postgres (`pg_isready`), MongoDB (`mongosh ping`) and Redpanda (`rpk cluster health`) got healthchecks and the target uses `up -d --wait`, so a service mission started right after doesn't fail on a database that is still booting. Validated with `docker compose config`; not run live (Docker wasn't running here).
64. **Feedback while building** — pressing Run used to leave an empty console until the compiler finished; it now says "Building…" (or "Building and running go test…" for labs) right away.
65. **Colony gets the same server hints** — a failed colony run, test or vet now also explains "server down → make start" or "server restarted → refresh".
66. **Colony says when its live connection dies** — after a server restart the colony's event stream is rejected and closes for good, so the colony froze silently. It now prints "live connection lost, refresh the page" once. Verified by restarting the server with the colony open.
67. **Format with gofmt (Shift+Alt+F)** — new `POST /api/format` runs Go's own `go/format` in-process (nothing executed) and the editor replaces the file through the undo history, keeping the caret on the same code character (verified: caret stays before `Move(` after reformatting). Code that doesn't parse yet gets a short "can't format yet: <reason>". Test: `TestFormatCode`; shortcut in the sheet and README.
68. **README testing section** — `make check` / `make test` / `make test-web`, and how to tune `ORBITAL_TEST_PARALLEL`.
69. **`{ }` format button** — a small button next to the file tabs formats the file with gofmt, for people who don't know the shortcut (tooltip shows Shift+Alt+F).
70. **Mission intro mentions the keys** — the first console line of every mission now says Run is Ctrl/⌘+Enter and `?` lists the shortcuts.
71. **VS Code workspace suggests the Go extension** — the `workspace/` module now gets `.vscode/extensions.json` (recommends `golang.go`, which provides the gopls autocomplete/hover the README promises) and `.vscode/settings.json` (gofmt on save for Go files). Files the player already has are never overwritten. Test: `TestWorkspaceSuggestsGoExtensionWithoutOverwriting`.
72. **Whole server suite ~2× faster end to end** — with the parallel level tests, `go test ./cmd/orbital/` now takes 492 s (8.2 min) instead of ~15–21 min; `make test` help and README updated.
73. **Fixed broken code formatting in two lessons** — 20-1 (Arabic and English) wrote `~<br>⏎~` with a real line break inside the code span, so the rule "newlines become `<br>\n`" rendered as stray tildes; it now shows `<br>\n`. 20-4 (English) left `~(9-3)-1` unclosed. 16-3 (Arabic) used a bare `~50` for "about 50", now "حوالي 50".
74. **Regression test for lesson markup** — `TestInlineCodeSpansCloseOnTheirLine` checks every goal, lesson, hint and step block in both languages: each `~code~` span must close on its own line (the client renders line by line), allowing the standalone generics `~` and table rows.
