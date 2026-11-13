package sim

import "orbital/colony/econ"

const (
	siftStock       = 3
	siftRefillTicks = 100
)

func (w *World) Melt(id int) bool { n, _ := w.TryMelt(id); return n > 0 }

func (w *World) TryMelt(id int) (n int, err error) {
	w.act(id, "melt", func() (bool, string, bool) {
		if w.Econ == nil || w.Z != 0 || w.here().Ground != Ice {
			err = econ.ErrNotHere
			return false, say("التسييح محتاج جليد على سطح المستعمرة", "Melting requires surface ice in the colony"), false
		}
		err = w.collectGood(Water, 1)
		if err != nil {
			return false, buildingMessage(err), false
		}
		n = 1
		t := w.here()
		t.Ground, t.MeltAt = Regolith, w.Tick+600
		w.syncInv()
		return true, "", true
	})
	return n, err
}

func (w *World) Sift(id int) bool { n, _ := w.TrySift(id); return n > 0 }

func (w *World) TrySift(id int) (n int, err error) {
	w.act(id, "sift", func() (bool, string, bool) {
		if w.Econ == nil || w.Z != 0 || w.here().Ground != Regolith {
			err = econ.ErrNotHere
			return false, say("الغربلة محتاجة ريجوليث على سطح المستعمرة", "Sifting requires surface regolith in the colony"), false
		}
		t := w.here()
		refilled := (w.Tick - t.SiftAt) / siftRefillTicks
		if t.SiftUsed > 0 && refilled > 0 {
			t.SiftUsed = max(0, t.SiftUsed-refilled)
			t.SiftAt += refilled * siftRefillTicks
		}
		if t.SiftUsed == 0 {
			t.SiftAt = w.Tick
		}
		if t.SiftUsed >= siftStock {
			err = econ.ErrNoInput
			return false, say("مخزون السيليكا خلص؛ استنى يتجدد", "Silica stock is depleted; wait for it to refill"), false
		}
		err = w.collectGood(Silica, 1)
		if err != nil {
			return false, buildingMessage(err), false
		}
		n = 1
		t.SiftUsed++
		w.syncInv()
		return true, "", false
	})
	return n, err
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
