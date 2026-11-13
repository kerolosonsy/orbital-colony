package sim

import (
	"math"
	"time"

	"orbital/colony/econ"
)

func (w *World) seasonActionClock(speed int) bool {
	if !w.seasonClock() {
		return false
	}
	// Single-drone seasons use action time so scheduler delays cannot affect scores.
	w.Tick += max(1, int(math.Round(float64(w.ActionDur*time.Duration(speed))/float64(w.Unit))))
	return true
}

func (w *World) seasonClock() bool { return w.Econ != nil && w.Econ.SeasonLimits().MaxDrones == 1 }

func (w *World) finishSeason() bool {
	if w.Econ == nil {
		return false
	}
	season := econ.FindSeason(w.Econ.Season)
	if season == nil || int64(w.Tick) < season.Ticks {
		return false
	}
	w.Tick = int(season.Ticks)
	w.econSnapshot("season-end", say("الموسم خلص", "Season finished"), true, true)
	w.done = true
	w.exit(0)
	return true
}

func (w *World) seasonStorm(tick int64) bool {
	period := w.Econ.SeasonLimits().StormPeriod
	if period == 0 {
		return false
	}
	if w.Econ.Weather.At == 0 {
		hit := (tick/period + 1) * period
		if hit-econ.StormWarningTicks <= tick {
			hit += period
		}
		w.Econ.Weather = econ.Weather{Storm: "none", At: hit - econ.StormWarningTicks}
	}
	return true
}

func (w *World) clampSeasonTick() {
	if w.Econ != nil {
		if season := econ.FindSeason(w.Econ.Season); season != nil {
			w.Tick = min(w.Tick, int(season.Ticks))
		}
	}
}

func (w *World) recoverSeasonMarket() {
	if w.seasonClock() {
		w.Econ.Recover(1)
	}
}
