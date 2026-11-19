# Improvement sweep (branch `improve-sweep`)

Each entry: what changed and why. Numbered in the order they landed.

1. **Removed dead error-tip helpers** — `compileTips`, `runTips`, `legacyTipList`, `tipFor` in `cmd/orbital/errtips.go` were left over from the old tip list and never called (staticcheck U1000). Also removed unused `ch01Variable` in `colony/levels/ch01_lessons.go`.
2. **Raw-string regexes in lesson tests** — ch12/21/26/27 lesson tests used double-escaped regex strings; switched to raw strings so they read like the pattern they match (staticcheck S1007).
