package econ

import (
	"fmt"
	"math"
)

const NeighborLineLimit = 1024

type NeighborProgress struct {
	Trades     int `json:"trades"`
	Reconnects int `json:"reconnects"`
	MaxLine    int `json:"maxLine"`
}

func addNeighborTree() {
	n := Node{ID: "sys.neighbor", Branch: "الشبكات", Name: "المستعمرة الجارة", NameEn: "Neighbour Colony", Desc: "تجارة مع جار عبر TCP؛ الاتصال محتاج حدود سطور ومواعيد نهائية وإعادة اتصال.", DescEn: "Trade with a neighbour over TCP using bounded lines, deadlines and reconnection.", After: []Prereq{{"api.prices", 1}}, Chapter: 31}
	addResearchNode(n)
	Find(n.ID).Max = 1
	Find(n.ID).GroupEn = "Networks"
}

func init() {
	Contracts = append(Contracts,
		&Contract{ID: "n31-trade-route", Chapter: 31, Requires: []string{"sel:net.DialTimeout", "sel:.SetDeadline"}, Window: 3000,
			Title: "طريق التجارة", TitleEn: "Trade Route", Brief: "كمّل عشر صفقات TCP مع الجار. وافق وبعدين اشحن الطحالب. استخدم deadlines وقارئ محدود وإعادة اتصال. محتاج sys.neighbor. الدرس: فصل 31.", BriefEn: "Complete ten TCP trades with the neighbour. Accept, then ship moss. Use deadlines, a bounded reader and reconnection. Requires sys.neighbor. Lesson: chapter 31.", Goal: Goal{Submit: "neighbor.trades", MaxLost: NoLimit}, Reward: Reward{Credits: 300}},
		&Contract{ID: "n31-bad-line", Chapter: 31, After: []string{"n31-trade-route"}, Requires: []string{"sel:net.DialTimeout", "sel:.SetDeadline", "sel:bufio.NewReaderSize", "sel:.ReadSlice"}, Window: 3000,
			Title: "سطور بايظة", TitleEn: "Bad Lines", Brief: "كمّل خمس صفقات مع أعلى معدل أعطال. الجار بيبعت سطور أكبر من 1024 بايت وgarbage، بيقف في نص السطر وبيقطع الاتصال. اقفل الاتصال عند تجاوز الحد واستخدم deadlines وإعادة اتصال. محتاج sys.neighbor. الدرس: فصل 31.", BriefEn: "Complete five trades at the highest fault rate. The neighbour sends lines over 1024 bytes, garbage, stalls mid-line and disconnects. Close on oversized lines; use deadlines and reconnect. Requires sys.neighbor. Lesson: chapter 31.", Goal: Goal{Submit: "neighbor.resilient", MaxLost: NoLimit}, Reward: Reward{Credits: 500}},
	)
}

func (s *State) IsNeighborContract() bool {
	return s.Active == "n31-trade-route" || s.Active == "n31-bad-line"
}

func (s *State) NeighborProgressPair(tick int64) (bool, bool, [2]string) {
	c, r := s.ActiveContract()
	need := 10
	if c.ID == "n31-bad-line" {
		need = 5
	}
	progress := r.Neighbor
	done := progress.Trades >= need
	if c.ID == "n31-bad-line" {
		done = done && progress.Reconnects >= 4 && progress.MaxLine <= NeighborLineLimit
	}
	why := [2]string{fmt.Sprintf("صفقات %d/%d، إعادة اتصال %d", progress.Trades, need, progress.Reconnects), fmt.Sprintf("trades %d/%d, reconnects %d", progress.Trades, need, progress.Reconnects)}
	return done, !done && tick-r.Start > c.Window, why
}

func (s *State) NeighborFraction() float64 {
	_, r := s.ActiveContract()
	if s.Active == "n31-bad-line" {
		return math.Min(1, math.Min(float64(r.Neighbor.Trades)/5, float64(r.Neighbor.Reconnects)/4))
	}
	return math.Min(1, float64(r.Neighbor.Trades)/10)
}
