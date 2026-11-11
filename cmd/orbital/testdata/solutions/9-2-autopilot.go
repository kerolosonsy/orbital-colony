package main

import "orbital/colony"

type Farmer struct{}

func (Farmer) Plan(t colony.Tile) []colony.Action {
	switch t.Ground {
	case colony.Ice:
		return []colony.Action{colony.PlantLichen}
	case colony.Regolith:
		return []colony.Action{colony.DoTill, colony.PlantMoss}
	case colony.Soil:
		if t.Entity == colony.Moss && t.Ripe {
			return []colony.Action{colony.DoHarvest, colony.PlantMoss}
		}
		return []colony.Action{colony.PlantMoss}
	}
	return nil
}

func main() {
	colony.Autopilot(Farmer{})
}
