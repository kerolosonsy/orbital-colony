package sim

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"

	"orbital/colony/econ"
)

func TestNeighborTransferAtomicAndExactlyOnce(t *testing.T) {
	st := econ.New(31)
	st.Inv["moss"] = 4
	var output bytes.Buffer
	w := NewColony(&output, st, 1)
	w.Realtime = false
	w.neighborCommand([]byte(`{"credit":20,"take":{"moss":3,"iron":1},"receipt":"missing"}`))
	if st.Credits != 0 || st.Inv["moss"] != 4 {
		t.Fatal("partial transfer applied")
	}
	var clients sync.WaitGroup
	for range 20 {
		clients.Add(1)
		go func() {
			defer clients.Done()
			w.neighborCommand([]byte(`{"credit":12.5,"take":{"moss":4},"receipt":"neighbor:1"}`))
		}()
	}
	clients.Wait()
	if st.Credits != 12.5 || st.Earned != 12.5 || st.Inv["moss"] != 0 {
		t.Fatalf("duplicate payment: %+v", st)
	}
	body, _ := json.Marshal(st)
	restored := new(econ.State)
	if err := json.Unmarshal(body, restored); err != nil {
		t.Fatal(err)
	}
	w.Econ = restored
	w.neighborCommand([]byte(`{"credit":12.5,"take":{"moss":4},"receipt":"neighbor:1"}`))
	if restored.Credits != 12.5 {
		t.Fatal("receipt lost across save")
	}
}
