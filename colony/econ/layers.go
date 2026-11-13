package econ

import "math"

// AtZ preserves At's surface meaning while exposing saved underground tiles.
func (s *State) AtZ(x, y, z int) *Tile {
	if z == 0 {
		return s.At(x, y)
	}
	return &s.Under[z-1][y*s.Size+x]
}

func (s *State) normalizeLayers() {
	// Five drill upgrades need five underground layers, including in old saves.
	s.Depth = max(6, s.Depth, len(s.Under)+1)
	for len(s.Under) < s.Depth-1 {
		s.Under = append(s.Under, nil)
	}
	s.resizeUnder(s.Size, s.Size)
}

func (s *State) resizeUnder(n, oldSize int) {
	for z, old := range s.Under {
		layer := make([]Tile, n*n)
		for y := 0; y < n; y++ {
			for x := 0; x < n; x++ {
				if x < oldSize && y < oldSize && y*oldSize+x < len(old) && old[y*oldSize+x].G != "" {
					layer[y*n+x] = old[y*oldSize+x]
				} else {
					layer[y*n+x] = undergroundAt(s.Seed, x, y, z+1)
				}
			}
		}
		s.Under[z] = layer
	}
}

func undergroundAt(seed int64, x, y, z int) Tile {
	// Coordinate hashing keeps newly revealed ore independent of expansion order.
	h := uint64(seed) ^ uint64(x)*0x9E3779B97F4A7C15 ^ uint64(y)*0xC2B2AE3D27D4EB4F ^ uint64(z)*0x94D049BB133111EB
	h ^= h >> 31
	h *= 0xBF58476D1CE4E5B9
	h ^= h >> 29
	t := Tile{G: "stone"}
	switch r := h % 100; {
	case z >= 2 && r < 8:
		t.E = "geode"
	case r < 30:
		t.E = "iron"
	case z >= 2 && r < 42:
		t.E = "copper"
	case z >= 3 && r >= 42 && r < 52:
		t.E = "sulfur"
	case z >= 4 && r >= 52 && r < 60:
		t.E = "titanium"
	case z == 5 && r >= 60 && r < 64:
		t.E = "uranium"
	}
	return t
}

func addDepthUpgrades() {
	add(&Upgrade{ID: "drill", Group: "الأعماق", Name: "الحفار", Desc: "كل مستوى يسمح بالحفر لطبقة أعمق.", Max: 5, cost: func(l int) Cost {
		c := Cost{Credits: math.Round(800 * math.Pow(2.6, float64(l)))}
		if l >= 1 {
			c.Goods = map[string]float64{"iron": 10 * math.Pow(2, float64(l-1))}
		}
		return c
	}})
	add(&Upgrade{ID: "api.dig", Group: "الأعماق", Name: "رخصة الحفر", Desc: "Dig و DigDown و Descend و Ascend.", Max: 1, Lesson: 17, cost: credits(600, 1)})
}

func (s *State) migrateV3() {
	for z, layer := range s.Under {
		for i := range layer {
			// Revealed ore and excavated tiles belong to the player already.
			if layer[i] == (Tile{G: "stone"}) {
				layer[i] = undergroundAt(s.Seed, i%s.Size, i/s.Size, z+1)
			}
		}
	}
	s.Version = 3
}
