package econ

func (b *Building) SetMode(mode string) error {
	if b.Kind != "generator" || mode != "solar" && mode != "burner" && mode != "reactor" {
		return ErrBadRecipe
	}
	b.Mode = mode
	return nil
}

func (b *Building) generatorSupply(tick int64) float64 {
	if b.Kind != "generator" || !b.Enabled || b.Wear >= 100 {
		return 0
	}
	good, n := GeneratorFuel(b.Mode, 1)
	if n > 0 && b.In[good]+1e-9 < n {
		return 0
	}
	return GeneratorOutput(b.Mode, tick)
}

func (s *State) PowerSupply(tick int64) float64 {
	supply := 0.0
	for _, b := range s.Buildings {
		supply += b.generatorSupply(tick)
	}
	return max(BaselinePower, supply)
}

func (s *State) ConsumePower(tick int64) float64 {
	supply := s.PowerSupply(tick)
	for _, b := range s.Buildings {
		if b.generatorSupply(tick) == 0 {
			continue
		}
		good, n := GeneratorFuel(b.Mode, 1)
		if n > 0 {
			b.In[good] = max(0, b.In[good]-n)
		}
	}
	return supply
}

func (s *State) RefreshPower(tick int64) {
	s.Power = PowerStatus{Supply: s.PowerSupply(tick), Day: Daylight(tick)}
	powered, demand := Allocate(s.Power.Supply, s.Buildings)
	s.Power.Demand = demand
	for i, b := range s.Buildings {
		b.Powered = powered[i] && (s.RentRate(len(s.Fleet)) == 0 || s.Credits > 0)
	}
}

func collectorNode() Node {
	return Node{ID: "collector", Branch: "الطاقة", Name: "جامع الغاز", NameEn: "Gas Collector", Desc: "يجمع الهيليوم بالليل بس.", DescEn: "Collects helium only at night.", After: []Prereq{{"generator", 1}}, Chapter: 19, Cost: nodeCost(500, 1.6)}
}

func addInfrastructureUpgrades() {
	add(&Upgrade{ID: "cargo", Group: "الهانجر", GroupEn: "Hangar", Name: "حمولة أكبر", NameEn: "Larger Cargo", Desc: "كل مستوى يزوّد حمولة الدرون 5 وحدات.", DescEn: "Each level adds 5 units of drone cargo.", Max: 8, Requires: []string{"composter"}, Lesson: 10, cost: mixed(100, 1.6, "moss", 10, 1.15)})
}
