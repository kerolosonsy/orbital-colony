package sim

import "orbital/colony/econ"

func (w *World) advanceRank() {
	w.Econ.RankWindow.Advance(int64(w.Tick - w.lastRankTick))
	w.lastRankTick = w.Tick
}

func (w *World) recordRankSale(good string, amount float64) {
	w.advanceRank()
	w.Econ.RecordRankSale(good, amount)
}

func (w *World) recordRankCommand() {
	w.advanceRank()
	w.Econ.RankWindow.RecordCommand()
}

func (w *World) checkRank() {
	w.advanceRank()
	if w.Econ.PromoteRank() {
		w.syncEcon()
		msg := sayf("⬆ المستعمرة وصلت للرتبة %d", "⬆ Colony reached rank %d", w.Econ.Rank)
		w.emitRaw(econEvent{Op: "rank", T: w.Tick, OK: true, Msg: msg, Credits: w.Econ.Credits, State: w.Econ})
	}
}

func (w *World) Rank() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil {
		return 0
	}
	return w.Econ.Rank
}

func (w *World) RankGoal() (kind string, target, current float64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil {
		goal := econ.GoalFor(0)
		return goal.Kind, goal.Target, 0
	}
	w.clock()
	w.advanceRank()
	progress := w.Econ.RankProgress()
	return progress.Kind, progress.Target, progress.Value
}
