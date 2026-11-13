package sim

import "orbital/colony/econ"

// InitLayers preserves the surface and creates solid underground layers. Setup only.
func (w *World) InitLayers(depth int) {
	w.Depth = max(1, depth)
	w.under = make([][]Tile, w.Depth-1)
	for z := range w.under {
		w.under[z] = make([]Tile, w.W*w.H)
		for i := range w.under[z] {
			w.under[z][i].Ground = Stone
		}
	}
}

// AtZ is the setup and goal-check accessor for a specific depth.
func (w *World) AtZ(x, y, z int) *Tile {
	if z == 0 {
		return w.At(x, y)
	}
	return &w.under[z-1][y*w.W+x]
}

func (w *World) here() *Tile { return w.AtZ(w.X, w.Y, w.Z) }

// Underground seeds mission ore using the mission's random source. Setup only.
func (w *World) Underground(depth int, ironChance, geodeChance float64) {
	w.InitLayers(depth)
	for z := 1; z < w.Depth; z++ {
		for i := range w.under[z-1] {
			r := w.Rand.Float64()
			if z >= 2 && r < geodeChance {
				w.under[z-1][i].Entity = Geode
			} else if r < ironChance || z >= 2 && r < geodeChance+ironChance {
				w.under[z-1][i].Entity = Iron
			}
		}
	}
}

// layerCount preserves the single-layer event format used by existing missions.
func (w *World) layerCount() int {
	if w.Depth > 1 {
		return w.Depth
	}
	return 0
}

func (w *World) layer(id int) int {
	if id == 0 {
		return w.Z
	}
	return w.Drones[id].Z
}

func (w *World) Layer(id int) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.layer(id)
}

func (w *World) LayerDepth() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Depth
}

func walkable(g Ground) bool { return g == Tunnel || g == Shaft }

func directionStep(dir Direction) (dx, dy int, valid bool) {
	switch dir {
	case North:
		return 0, 1, true
	case East:
		return 1, 0, true
	case South:
		return 0, -1, true
	case West:
		return -1, 0, true
	}
	return 0, 0, false
}

func (w *World) collectOre(t *Tile) error {
	if t.Entity != Iron && t.Entity != Geode && t.Entity != Copper && t.Entity != Sulfur && t.Entity != Titanium && t.Entity != Uranium {
		return nil
	}
	n := 1
	if t.Entity == Geode {
		n = 3
	}
	if w.Econ != nil {
		name := t.Entity.String()
		if err := w.collectGood(t.Entity, float64(n)*w.Econ.YieldMult(name)); err != nil {
			return err
		}
		w.archiveHarvest(name, float64(n)*w.Econ.YieldMult(name))
		w.syncInv()
	} else {
		w.Inv[t.Entity] += n
	}
	t.Entity, t.Planted, t.Ready, t.Size = Empty, 0, 0, 0
	return nil
}

func (w *World) tileChanged(id, x, y, z int) {
	v := w.viewZ(x, y, z)
	w.emit(Event{Op: "tile", T: w.Tick, D: id, OK: true, X: w.X, Y: w.Y, Z: w.Z, Tile: &v})
}

func (w *World) Dig(id int, dir Direction) bool { return w.TryDig(id, dir) == nil }

func (w *World) TryDig(id int, dir Direction) (err error) {
	_, blocked := w.actErr(id, "dig", func() (bool, string, bool) {
		if m := w.gate("api.dig"); m != "" {
			err = ErrLocked
			return false, m, false
		}
		if w.Z == 0 {
			err = econ.ErrNotHere
			return false, say("على السطح استخدم DigDown", "Use DigDown on the surface"), false
		}
		dx, dy, valid := directionStep(dir)
		x, y := w.X+dx, w.Y+dy
		if !valid || x < 0 || y < 0 || x >= w.W || y >= w.H {
			err = econ.ErrNotHere
			return false, say("مفيش خانة في الاتجاه ده", "No tile in that direction"), false
		}
		t := w.AtZ(x, y, w.Z)
		if t.Ground != Stone {
			err = econ.ErrTileNotFree
			return false, say("الحفر محتاج صخر صلب", "Digging requires solid rock"), false
		}
		if err = w.collectOre(t); err != nil {
			return false, buildingMessage(err), false
		}
		t.Ground = Tunnel
		w.tileChanged(id, x, y, w.Z)
		return true, "", false
	})
	if blocked != nil {
		err = blocked
	}
	return err
}

func (w *World) DigDown(id int) bool { return w.TryDigDown(id) == nil }

func (w *World) TryDigDown(id int) (err error) {
	_, blocked := w.actErr(id, "digdown", func() (bool, string, bool) {
		if m := w.gate("api.dig"); m != "" {
			err = ErrLocked
			return false, m, false
		}
		if w.Z+1 >= w.Depth {
			err = econ.ErrNotHere
			return false, say("وصلت لأعمق طبقة", "You reached the deepest layer"), false
		}
		if w.Econ != nil && w.Econ.Level("drill") < w.Z+1 {
			err = ErrLocked
			return false, say("🔒 محتاج ترقية الحفار للعمق ده", "🔒 Upgrade the drill to reach this depth"), false
		}
		g := w.here().Ground
		if !(w.Z == 0 && (g == Regolith || g == Soil) || w.Z > 0 && g == Tunnel) {
			err = econ.ErrNotHere
			return false, say("مينفعش تحفر لتحت من الخانة دي", "You cannot dig down from this tile"), false
		}
		if w.Z == 0 && w.buildingAt(w.X, w.Y) != nil {
			err = econ.ErrTileNotFree
			return false, buildingMessage(econ.ErrTileNotFree), false
		}
		t := w.AtZ(w.X, w.Y, w.Z+1)
		if err = w.collectOre(t); err != nil {
			return false, buildingMessage(err), false
		}
		w.here().Ground = Shaft
		t.Ground = Tunnel
		w.tileChanged(id, w.X, w.Y, w.Z+1)
		return true, "", true
	})
	if blocked != nil {
		err = blocked
	}
	return err
}

func (w *World) Descend(id int) bool {
	return w.act(id, "descend", func() (bool, string, bool) {
		if m := w.gate("api.dig"); m != "" {
			return false, m, false
		}
		if w.Z+1 >= w.Depth || w.here().Ground != Shaft {
			return false, say("مفيش بئر للنزول هنا", "No shaft here to descend"), false
		}
		w.Z++
		return true, "", false
	})
}

func (w *World) Ascend(id int) bool {
	return w.act(id, "ascend", func() (bool, string, bool) {
		if m := w.gate("api.dig"); m != "" {
			return false, m, false
		}
		if w.Z == 0 || w.AtZ(w.X, w.Y, w.Z-1).Ground != Shaft {
			return false, say("مفيش بئر للصعود هنا", "No shaft here to ascend"), false
		}
		w.Z--
		return true, "", false
	})
}

// LayerGrid returns grid[x][y] of the ground on layer z (any layer, free).
func (w *World) LayerGrid(z int) [][]Ground {
	w.mu.Lock()
	defer w.mu.Unlock()
	if z < 0 || z >= max(1, w.Depth) {
		return nil
	}
	g := make([][]Ground, w.W)
	for x := range g {
		g[x] = make([]Ground, w.H)
		for y := range g[x] {
			g[x][y] = w.AtZ(x, y, z).Ground
		}
	}
	return g
}
