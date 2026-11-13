package econ

import "math"

const StormWarningTicks = 60
const StormStunTicks = 150
const ShockTicks = 300
const RushTicks = 400
const RushBonus = 2

const CleanStopContract = "c13b-clean-stop"

type Weather struct {
	Storm string `json:"storm"`
	At    int64  `json:"at"`
}

type RushOrder struct {
	Good   string  `json:"good"`
	Amount int     `json:"amount"`
	Left   float64 `json:"left"`
	Until  int64   `json:"until"`
	Bonus  float64 `json:"bonus"`
}

type MarketShock struct {
	Good   string  `json:"good"`
	Factor float64 `json:"factor"`
	Until  int64   `json:"until"`
}

// EventHash keeps event choices independent of the simulation's random stream.
func EventHash(seed, tick int64, kind uint64) uint64 {
	h := uint64(seed) ^ uint64(tick)*0x9E3779B97F4A7C15 ^ kind*0xC2B2AE3D27D4EB4F
	h ^= h >> 31
	h *= 0xBF58476D1CE4E5B9
	h ^= h >> 29
	return h
}

// EventAt chooses one ten-tick slot per period, away from the period's edges.
func EventAt(seed, epoch, period int64, kind uint64) int64 {
	return epoch*period + period/4 + int64(EventHash(seed, epoch, kind)%uint64(period/20))*10
}

func CleanStopAt(seed, start int64) int64 {
	return start + 200 + int64(EventHash(seed, start, 5)%101)
}

func (s *State) normalizeEvents() {
	if s.Weather.Storm == "" {
		s.Weather.Storm = "none"
	}
}

func (s *State) rushPayment(good string, amount, price float64) float64 {
	r := s.Rush
	if r == nil || r.Good != good || s.Tick >= r.Until {
		return amount * price
	}
	bonus := math.Min(amount, r.Left)
	r.Left -= bonus
	if r.Left <= 1e-9 {
		s.Rush = nil
	}
	return (amount + bonus*(r.Bonus-1)) * price
}

func (s *State) ExpireEvents(tick int64) {
	if s.Shock != nil && tick >= s.Shock.Until {
		s.Shock = nil
	}
	if s.Rush != nil && tick >= s.Rush.Until {
		s.Rush = nil
	}
}

func init() {
	c := &Contract{ID: CleanStopContract, Chapter: 13, After: []string{"c13-reed-clock"}, Requires: []string{"select"}, Window: 800,
		Title: "الإغلاق المرتب", TitleEn: "Clean Stop",
		Brief:   "هيجيلك SIGTERM في وقت مفاجئ. استخدم select مع colony.Shutdown، رجّع الدرونات للقاعدة وفرّغ الشحنة قبل الخروج. مفيش حمولة تضيع.",
		BriefEn: "SIGTERM arrives unexpectedly. Use select with colony.Shutdown, return drones to base and unload before exiting. Lose no cargo.",
		Goal:    Goal{MaxLost: NoLimit}, Reward: Reward{Credits: 300}}
	for i, prev := range Contracts {
		if prev.ID == "c13-reed-clock" {
			Contracts = append(Contracts[:i+1], append([]*Contract{c}, Contracts[i+1:]...)...)
			break
		}
	}
}

// DiscardCargo is called only when a player process ends; stored goods remain owned.
func (s *State) DiscardCargo() {
	for i := range s.Fleet {
		if s.Active == CleanStopContract {
			_, r := s.ActiveContract()
			for _, n := range s.Fleet[i].Cargo {
				r.CargoLost += n
			}
		}
		clear(s.Fleet[i].Cargo)
	}
}

// FinishCleanStop also handles programs returning from main before the grace timer.
func (s *State) FinishCleanStop() (done, settled bool) {
	s.DiscardCargo()
	if s.Active != CleanStopContract {
		return false, false
	}
	_, r := s.ActiveContract()
	r.Stopped = r.Stopping
	done, failed, _ := s.Progress(s.Tick)
	if done {
		_, _ = s.Settle(s.Tick)
	} else if failed {
		s.Fail()
	}
	return done, done || failed
}
