package econ

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func TestRushPaysDoubleOnlyForRemainingUnits(t *testing.T) {
	for _, sold := range []float64{1.5, 3, 5} {
		s, calm := New(4), New(4)
		s.Inv["moss"], calm.Inv["moss"] = 10, 10
		s.Rush = &RushOrder{Good: "moss", Amount: 3, Left: 3, Until: 400, Bonus: 2}
		want := calm.Sell("moss", math.Min(3, sold)) * 2
		if sold > 3 {
			want += calm.Sell("moss", sold-3)
		}
		if got := s.Sell("moss", sold); math.Abs(got-want) > 1e-9 {
			t.Fatalf("sold %g: earned %g want %g", sold, got, want)
		}
		if sold >= 3 && s.Rush != nil {
			t.Fatal("filled order remained active")
		}
	}
	s := New(4)
	s.Inv["moss"] = 1
	s.Tick = 400
	s.Rush = &RushOrder{Good: "moss", Amount: 3, Left: 3, Until: 400, Bonus: 2}
	if s.Sell("moss", 1) != PriceV3["moss"] {
		t.Fatal("expired order paid bonus")
	}
}

func TestEventSaveRoundTripAndLegacyOwnership(t *testing.T) {
	s := New(17)
	s.Inv["moss"], s.Credits, s.RP = 2500, 1234, 12
	s.Levels["engine"], s.Blueprints["api.try"] = 3, true
	s.Contracts["c05-careful-harvest"] = &ContractRun{Best: 10}
	s.Weather = Weather{Storm: "warn", At: 120}
	s.Shock = &MarketShock{Good: "moss", Factor: 1.6, Until: 310}
	s.Rush = &RushOrder{Good: "moss", Amount: 30, Left: 10, Until: 410, Bonus: 2}
	s.Fleet[0].Fault, s.Fleet[0].Stunned = true, 70
	for _, legacy := range []bool{false, true} {
		blob, _ := json.Marshal(s)
		if legacy {
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(blob, &fields); err != nil {
				t.Fatal(err)
			}
			delete(fields, "weather")
			delete(fields, "rush")
			delete(fields, "shock")
			fields["version"] = json.RawMessage("2")
			blob, _ = json.Marshal(fields)
		}
		var loaded State
		if err := json.Unmarshal(blob, &loaded); err != nil {
			t.Fatal(err)
		}
		loaded.Normalize()
		if !reflect.DeepEqual(loaded.Inv, s.Inv) || loaded.Credits != s.Credits || loaded.RP != s.RP || !reflect.DeepEqual(loaded.Levels, s.Levels) || !loaded.Blueprints["api.try"] || !loaded.Won("c05-careful-harvest") {
			t.Fatal("save lost ownership")
		}
		if !legacy && (!reflect.DeepEqual(loaded.Weather, s.Weather) || !reflect.DeepEqual(loaded.Shock, s.Shock) || !reflect.DeepEqual(loaded.Rush, s.Rush) || !loaded.Fleet[0].Fault || loaded.Fleet[0].Stunned != 70) {
			t.Fatal("live event state lost on reload")
		}
		if legacy && loaded.Weather.Storm != "none" {
			t.Fatal("legacy weather default missing")
		}
	}
}
