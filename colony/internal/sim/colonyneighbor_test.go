package sim

import (
	"bytes"
	"encoding/json"
	"strings"
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

func TestColonyCommandsRouteSystemTransfersOnce(t *testing.T) {
	st := econ.New(31)
	st.Inv["moss"] = 10
	var output bytes.Buffer
	w := NewColony(&output, st, 1)
	w.Realtime = false
	w.colonyCommands(strings.NewReader(`{"receipt":"neighbor:1","credit":12,"take":{"moss":2}}
{"id":"market-1","credit":5,"take":{"moss":1}}
{"controlSnapshot":true}
`))
	if st.Credits != 17 || st.Inv["moss"] != 7 {
		t.Fatalf("misrouted transfers: credits=%v moss=%v", st.Credits, st.Inv["moss"])
	}
	if !bytes.Contains(output.Bytes(), []byte(`"op":"transfer"`)) || !bytes.Contains(output.Bytes(), []byte(`"op":"system-command"`)) {
		t.Fatal("missing sidecar receipts", output.String())
	}
}
