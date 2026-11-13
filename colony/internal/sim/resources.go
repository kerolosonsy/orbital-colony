package sim

func (w *World) Melt(id int) bool {
	return w.act(id, "melt", func() (bool, string, bool) {
		if w.Econ == nil || w.Z != 0 || w.here().Ground != Ice {
			return false, say("التسييح محتاج جليد على سطح المستعمرة", "Melting requires surface ice in the colony"), false
		}
		t := w.here()
		t.Ground, t.MeltAt = Regolith, w.Tick+w.Econ.MeltRegrowthTicks()
		w.Econ.Inv["water"] += w.Econ.YieldMult("water")
		w.syncInv()
		return true, "", true
	})
}

func (w *World) Sift(id int) bool {
	return w.act(id, "sift", func() (bool, string, bool) {
		if w.Econ == nil || w.Z != 0 || w.here().Ground != Regolith {
			return false, say("الغربلة محتاجة ريجوليث على سطح المستعمرة", "Sifting requires surface regolith in the colony"), false
		}
		t := w.here()
		refilled := (w.Tick - t.SiftAt) / w.Econ.SiftRefillTicks()
		if t.SiftUsed > 0 && refilled > 0 {
			t.SiftUsed = max(0, t.SiftUsed-refilled)
			t.SiftAt += refilled * w.Econ.SiftRefillTicks()
		}
		if t.SiftUsed == 0 {
			t.SiftAt = w.Tick
		}
		if t.SiftUsed >= w.Econ.SiftStock() {
			return false, say("مخزون السيليكا خلص؛ استنى يتجدد", "Silica stock is depleted; wait for it to refill"), false
		}
		t.SiftUsed++
		w.Econ.Inv["silica"] += w.Econ.YieldMult("silica")
		w.syncInv()
		return true, "", false
	})
}

func (w *World) restoreIce(tick int) {
	w.Each(func(x, y int, t *Tile) {
		if t.MeltAt != 0 && tick >= t.MeltAt {
			t.Ground, t.MeltAt = Ice, 0
			v := w.view(x, y)
			w.emit(Event{Op: "tile", T: tick, OK: true, X: x, Y: y, Tile: &v})
		}
	})
}
