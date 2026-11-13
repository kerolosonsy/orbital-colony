package econ

import "testing"

func TestNeighborContractsCountTradesAndRequireReconnection(t *testing.T) {
	for _, tc := range []struct {
		id   string
		need int
	}{{"n31-trade-route", 10}, {"n31-bad-line", 5}} {
		t.Run(tc.id, func(t *testing.T) {
			st := New(31)
			st.Contracts["n31-trade-route"] = &ContractRun{Best: 1}
			if err := st.Accept(tc.id); err != nil {
				t.Fatal(err)
			}
			_, run := st.ActiveContract()
			if done, _, _ := st.Progress(0); done || st.Fraction() != 0 {
				t.Fatal("unstarted trades passed")
			}
			run.Neighbor.Trades = tc.need
			if tc.id == "n31-bad-line" {
				if done, _, _ := st.Progress(1); done {
					t.Fatal("resilient contract needs reconnection")
				}
				run.Neighbor.Reconnects = 4
				run.Neighbor.MaxLine = NeighborLineLimit + 1
				if done, _, _ := st.Progress(1); done {
					t.Fatal("unbounded lines passed")
				}
			}
			run.Neighbor.MaxLine = NeighborLineLimit
			if done, _, _ := st.Progress(1200); !done {
				t.Fatal("valid trades didn't pass")
			}
			if _, err := st.Settle(1200); err != nil {
				t.Fatal(err)
			}
			credits := st.Credits
			if _, err := st.Settle(1200); err == nil || st.Credits != credits {
				t.Fatal("reward paid twice")
			}
		})
	}
}
