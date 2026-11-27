package econ

import "testing"

// The deep-haul contract lends api.dig and needs drill level 2, so the drill
// must be buyable while the contract runs, before api.dig is owned.
func TestDeepHaulDrillBuyableWhileContractLendsDig(t *testing.T) {
	s := New(1)
	s.Credits, s.RP = 1000000, 1000
	s.Levels["field"] = 1
	if err := s.Accept("c17-deep-haul"); err != nil {
		t.Fatal(err)
	}
	if err := s.Buy("drill"); err != nil {
		t.Fatalf("drill not buyable during deep haul: %v", err)
	}
	s.Inv["iron"] = 10
	if err := s.Buy("drill"); err != nil || s.Level("drill") != 2 {
		t.Fatalf("drill level 2 not reachable during deep haul: %v", err)
	}
}
