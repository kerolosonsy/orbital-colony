package main

import (
	"strings"

	"orbital/colony/econ"
	"orbital/colony/levels"
)

func levelForLang(l *levels.Level, lang string) *levels.Level {
	if l == nil || lang != "en" || l.En == nil {
		return l
	}
	localized := *l
	localized.Title, localized.Goal = l.En.Title, l.En.Goal
	localized.Lesson, localized.Starter = l.En.Lesson, l.En.Starter
	localized.Hints = l.En.Hints
	if l.En.Files != nil {
		localized.Files = l.En.Files
	}
	return &localized
}

func shopForLang(st *econ.State, lang string) []econ.ShopItem {
	shop := st.Shop()
	for i := range shop {
		upgrade := *shop[i].Upgrade
		if lang == "en" {
			upgrade.Name, upgrade.Desc, upgrade.Group = upgrade.NameEn, upgrade.DescEn, upgrade.GroupEn
		}
		shop[i].Upgrade = &upgrade
		shop[i].Reason = economyText(lang, shop[i].Reason)
	}
	return shop
}

// Economy methods use the process language. Translate their diagnostic text
// without changing that shared language or duplicating purchase rules.
func economyText(lang, message string) string {
	for _, pair := range [][2]string{
		{"خلصت كل المستويات", "All levels completed"},
		{"المكافأة دي اتصرفت قبل كده", "This reward has already been claimed"},
		{"مفيش عقد نشط", "No active contract"},
		{"عندك عقد نشط بالفعل. خلّصه أو ألغيه الأول", "You already have an active contract. Finish or cancel it first"},
		{"العقد لسه مخلصش", "The contract is not finished yet"},
		{"محتاج مخطط من عقد", "Needs a contract blueprint"},
	} {
		if message == pair[0] || message == pair[1] {
			return tr(lang, pair[0], pair[1])
		}
	}
	for _, pair := range [][2]string{
		{"محتاج الأول: ", "Requires first: "},
		{"مفيش ترقية اسمها ", "No upgrade named "},
		{"مفيش عقد اسمه ", "No contract named "},
	} {
		for _, prefix := range pair {
			if tail, ok := strings.CutPrefix(message, prefix); ok {
				return tr(lang, pair[0], pair[1]) + economyName(lang, tail)
			}
		}
	}
	for _, prefix := range []string{"خلّص العقد الأول: ", "Finish the contract first: "} {
		if tail, ok := strings.CutPrefix(message, prefix); ok {
			for _, c := range econ.Contracts {
				if tail == c.Title || tail == c.TitleEn {
					tail = tr(lang, c.Title, c.TitleEn)
				}
			}
			return tr(lang, "خلّص العقد الأول: ", "Finish the contract first: ") + tail
		}
	}
	return economyShortfall(lang, message)
}

// contractView is a contract as the panel shows it, in one language.
type contractView struct {
	ID       string            `json:"id"`
	Chapter  int               `json:"chapter"`
	Title    string            `json:"title"`
	Brief    string            `json:"brief"`
	Window   int64             `json:"window"`
	Status   string            `json:"status"` // locked, available, active, done
	After    []string          `json:"after,omitempty"`
	Run      *econ.ContractRun `json:"run,omitempty"`
	Frac     float64           `json:"frac,omitempty"`
	Progress string            `json:"progress,omitempty"`
	Left     int64             `json:"left,omitempty"`
}

func contractsForLang(st *econ.State, lang string, done map[int]bool) []contractView {
	open := map[string]bool{}
	for _, c := range st.Available(done) {
		open[c.ID] = true
	}
	out := make([]contractView, 0, len(econ.Contracts))
	for _, c := range econ.Contracts {
		v := contractView{ID: c.ID, Chapter: c.Chapter, Title: tr(lang, c.Title, c.TitleEn), Brief: tr(lang, c.Brief, c.BriefEn), Window: c.Window, Run: st.Contracts[c.ID]}
		for _, a := range c.After {
			if !st.Won(a) {
				p := econ.FindContract(a)
				v.After = append(v.After, tr(lang, p.Title, p.TitleEn))
			}
		}
		switch {
		case st.Active == c.ID:
			v.Status = "active"
			_, _, why := st.ProgressPair(st.Tick)
			v.Frac, v.Progress, v.Left = st.Fraction(), tr(lang, why[0], why[1]), st.Left(st.Tick)
		case st.Won(c.ID):
			v.Status = "done"
		case open[c.ID] || st.Active != "" && done[c.Chapter] && len(v.After) == 0:
			v.Status = "available"
		default:
			v.Status = "locked"
		}
		out = append(out, v)
	}
	return out
}

func economyName(lang, name string) string {
	for _, resources := range [][]econ.Resource{econ.RawV3, econ.Crafted} {
		for _, g := range resources {
			if name == g.Name || name == g.NameEn {
				return tr(lang, g.Name, g.NameEn)
			}
		}
	}
	for _, u := range econ.Catalog {
		if name == u.Name || name == u.NameEn {
			return tr(lang, u.Name, u.NameEn)
		}
		if strings.HasPrefix(u.ID, "yield.") {
			ar, en := strings.TrimPrefix(u.Name, "حصاد "), strings.TrimPrefix(u.NameEn, "Harvest ")
			if name == ar || name == en {
				return tr(lang, ar, en)
			}
		}
	}
	return name
}

func economyShortfall(lang, message string) string {
	amount, good := "", ""
	if tail, ok := strings.CutPrefix(message, "ناقصك "); ok {
		amount, good, _ = strings.Cut(tail, " ")
	} else if tail, ok := strings.CutPrefix(message, "You need "); ok {
		amount, good, _ = strings.Cut(tail, " more ")
	} else {
		return message
	}
	if good == "كريدت" || good == "credits" {
		good = tr(lang, "كريدت", "credits")
	} else {
		good = economyName(lang, good)
	}
	return tr(lang, "ناقصك "+amount+" "+good, "You need "+amount+" more "+good)
}

func mutantDescription(lang, description string) string {
	for _, pair := range [][2]string{
		{"خرق عقد المسار", "Route contract violation"},
		{"إطار بترتيب بايتات غلط", "Wrong byte order"},
		{"نوع مجهول مقبول", "Unknown kind accepted"},
		{"محتوى بيشارك المدخل", "Payload aliases input"},
		{"تجاهل المسح السالب", "Skip negative scans"},
		{"نسيان آخر مسح", "Omit final scan"},
		{"مسح متحسب مرتين", "Double counted scan"},
		{"حد الـ 50 كجم: جرّبت 50 بالظبط؟", "The 50 kg boundary: did you try exactly 50?"},
		{"المسافة صفر مع شحنة", "Zero distance with cargo"},
		{"تكلفة كل خانة مسافة", "Fuel cost per tile of distance"},
		{"الكسور في الشحنة (6 كجم مثلاً)", "Cargo fractions (6 kg, for example)"},
		{"رسوم الشحنة التقيلة", "Heavy cargo surcharge"},
		{"شحنة تقيلة جداً (أكتر من 100)", "Very heavy cargo (over 100)"},
	} {
		if description == pair[0] || description == pair[1] {
			return tr(lang, pair[0], pair[1])
		}
	}
	return description
}
