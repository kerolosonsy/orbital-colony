package econ

import (
	"slices"

	"orbital/colony/i18n"
)

type Season struct {
	ID, Title, TitleEn, Desc, DescEn string
	Seed                             int64
	Ticks                            int64
	Kit                              map[string]int
	Credits                          float64
	Limits                           Limits
	Reference                        float64
}

var Seasons = []Season{
	{ID: "one-drone", Reference: 1169.45, Title: "درون واحد", TitleEn: "One Drone", Desc: "درون واحد بس لمدة 6000 tick.", DescEn: "Only one drone for 6000 ticks.", Seed: 6101, Ticks: 6000, Limits: Limits{MaxDrones: 1}},
	{ID: "blind-flight", Reference: 1277.83, Title: "طيران من غير حساسات", TitleEn: "Blind Flight", Desc: "الرادار والسكانر مقفولين.", DescEn: "Radar and scanner are disabled.", Seed: 6102, Ticks: 6000, Limits: Limits{Deny: []string{"api.radar", "api.scanner"}}},
	{ID: "volatile-market", Reference: 1909.44, Title: "سوق متقلب", TitleEn: "Volatile Market", Desc: "صدمة سوق كل 200 tick تقريباً.", DescEn: "A market shock roughly every 200 ticks.", Seed: 6103, Ticks: 6000, Limits: Limits{ShockPeriod: 200}},
	{ID: "storm-front", Reference: 963.75, Title: "جبهة عواصف", TitleEn: "Storm Front", Desc: "عاصفة كل 400 tick، بإنذار قبلها بـ60 tick.", DescEn: "Storms every 400 ticks, with 60 ticks of warning.", Seed: 6104, Ticks: 6000, Limits: Limits{StormPeriod: 400}},
	{ID: "fixed-kit", Reference: 1284.35, Title: "عدة ثابتة", TitleEn: "Fixed Kit", Desc: "ابدأ بعدة ثابتة؛ الشراء مقفول طول الموسم.", DescEn: "Start with a fixed kit; purchases are disabled throughout the season.", Seed: 6105, Ticks: 6000, Limits: Limits{NoBuy: true}},
	{ID: "small-field", Reference: 2540.24, Title: "أرض صغيرة", TitleEn: "Small Field", Desc: "طول وعرض الأرض نص الحجم المعتاد.", DescEn: "The field has half the usual width and height.", Seed: 6106, Ticks: 6000, Limits: Limits{FieldDivisor: 2}},
}

func init() {
	for i := range Seasons {
		Seasons[i].Credits = 200
		Seasons[i].Kit = map[string]int{"field": 4, "engine": 4, "yield.moss": 4, "grow.moss": 2, "api.try": 1, "api.scanner": 1, "api.radar": 1, "api.prices": 1, "api.spawn": 1, "hangar": 2}
		if Seasons[i].ID == "small-field" {
			Seasons[i].Kit["crop.glowcap"] = 1
			Seasons[i].Kit["yield.glowcap"] = 4
			Seasons[i].Kit["grow.glowcap"] = 2
		}
	}
}

func FindSeason(id string) *Season {
	for i := range Seasons {
		if Seasons[i].ID == id {
			return &Seasons[i]
		}
	}
	return nil
}

func (season *Season) New() *State {
	st := New(season.Seed)
	st.Season = season.ID
	st.Credits = season.Credits
	for id, level := range season.Kit {
		st.Levels[id] = level
		st.Blueprints[id] = true
	}
	st.Normalize()
	return st
}

func (s *State) SeasonLimits() Limits {
	if season := FindSeason(s.Season); season != nil {
		return season.Limits
	}
	return Limits{}
}

func (s *State) FieldSize() int {
	return max(1, FieldSize(s.Levels["field"])/max(1, s.SeasonLimits().FieldDivisor))
}

func (s *State) SeasonDenies(id string) bool { return slices.Contains(s.SeasonLimits().Deny, id) }

func (s *State) PurchaseReason() string {
	if s.SeasonLimits().NoBuy {
		return i18n.T("الشراء مقفول في الموسم ده", "Purchases are disabled in this season")
	}
	return ""
}

func (s *State) ShockPeriod() int64 {
	if period := s.SeasonLimits().ShockPeriod; period > 0 {
		return period
	}
	if s.Won("c02-price-watch") {
		return 1800
	}
	return 0
}

func (s *State) seasonDroneCap(n int) int {
	if cap := s.SeasonLimits().MaxDrones; cap > 0 {
		return min(n, cap)
	}
	return n
}
