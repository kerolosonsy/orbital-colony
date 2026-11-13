package sim

import (
	"bytes"
	"encoding/json"
	"testing"

	"orbital/colony/econ"
)

func TestControlResultOnlySettlesMatchingActiveContract(t *testing.T) {
	var output bytes.Buffer
	st := econ.New(11)
	if err := st.Accept("h26-control-api"); err != nil {
		t.Fatal(err)
	}
	w := NewColony(&output, st, 11)
	w.colonyCommands(bytes.NewBufferString("{\"controlResult\":{\"id\":\"s29-transfer\",\"ok\":true}}\n"))
	if st.Won("h26-control-api") || st.Active != "h26-control-api" {
		t.Fatal("unrelated result settled contract")
	}
	w.colonyCommands(bytes.NewBufferString("{\"controlResult\":{\"id\":\"h26-control-api\",\"ok\":true}}\n"))
	if !st.Won("h26-control-api") {
		t.Fatal(output.String())
	}
	if err := st.Accept("h26-control-api"); err != nil {
		t.Fatal(err)
	}
	w.colonyCommands(bytes.NewBufferString("{\"controlResult\":{\"id\":\"h26-control-api\",\"ok\":false}}\n"))
	if !st.Won("h26-control-api") || st.Active != "" {
		t.Fatal("failed replay removed old win")
	}
}

func TestSystemSnapshotCountersPersistAndInventoryIsDetached(t *testing.T) {
	st := econ.New(11)
	st.Inv["moss"] = 12.5
	w := NewColony(&bytes.Buffer{}, st, 11)
	w.Realtime = false
	w.Export(0, Moss, 5)
	snapshot := w.SystemStats()
	if snapshot.Sold != 5 || snapshot.Commands != 1 || snapshot.Inventory["moss"] != 7.5 {
		t.Fatalf("%+v", snapshot)
	}
	snapshot.Inventory["moss"] = 0
	if w.SystemStats().Inventory["moss"] != 7.5 {
		t.Fatal("snapshot changed inventory")
	}
	body, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	var loaded econ.State
	if err := json.Unmarshal(body, &loaded); err != nil {
		t.Fatal(err)
	}
	restarted := NewColony(&bytes.Buffer{}, &loaded, 11)
	if got := restarted.SystemStats(); got.Sold != 5 || got.Commands != 1 || got.Inventory["moss"] != 7.5 {
		t.Fatalf("restart: %+v", got)
	}
}
