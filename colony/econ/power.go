package econ

import "sort"

func Daylight(tick int64) bool { return (tick%1200+1200)%1200 < 600 }

// GeneratorOutput assumes the caller supplies the fuel returned by GeneratorFuel.
func GeneratorOutput(kind string, tick int64) float64 {
	switch kind {
	case "solar":
		if Daylight(tick) {
			return 10
		}
	case "burner":
		return 25
	case "reactor":
		return 100
	}
	return 0
}

func GeneratorFuel(kind string, ticks int) (good string, amount float64) {
	switch kind {
	case "burner":
		return "fuel", float64(max(0, ticks)) / 100
	case "reactor":
		return "cell", float64(max(0, ticks)) / 1000
	}
	return "", 0
}

func Draw(kind string) float64 {
	switch kind {
	case "collector":
		return 2
	case "composter":
		return 1
	case "smelter", "kiln":
		return 4
	case "refinery":
		return 6
	case "assembler":
		return 8
	case "reactor-lab":
		return 20
	case "lab":
		return 10
	}
	return 0
}

func Allocate(supply float64, buildings []*Building) (powered map[int]bool, demand float64) {
	powered = make(map[int]bool)
	order := make([]int, 0, len(buildings))
	for i, b := range buildings {
		if b.Enabled {
			order = append(order, i)
			demand += Draw(b.Kind)
		}
	}
	sort.SliceStable(order, func(i, j int) bool {
		return buildings[order[i]].Priority > buildings[order[j]].Priority
	})
	for _, i := range order {
		draw := Draw(buildings[i].Kind)
		if draw <= supply {
			powered[i] = true
			supply -= draw
		}
	}
	return powered, demand
}
