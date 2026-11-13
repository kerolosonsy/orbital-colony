package sim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"orbital/colony/econ"
)

// Purchases used to race action scheduling, and resizing raced path queries.
func TestColonyConcurrentDronesTickerAndPurchases(t *testing.T) {
	st := econ.New(11)
	st.Credits = 1e12
	st.Rank = 100
	for good := range econ.PriceV3 {
		st.Inv[good] = 1e12
	}
	st.Levels["api.spawn"], st.Levels["hangar"] = 1, 3
	st.Levels["lab"], st.Levels["refinery"] = 1, 1
	st.Buildings = []*econ.Building{
		{Kind: "lab", X: 1, Y: 1, Level: 1, Enabled: true, In: map[string]float64{"biomass": 10}, Out: map[string]float64{}},
		{Kind: "refinery", X: 2, Y: 1, Level: 1, Enabled: true, In: map[string]float64{}, Out: map[string]float64{}},
	}
	st.Weather = econ.Weather{Storm: "none", At: 10}
	var trace bytes.Buffer
	w := NewColony(&trace, st, 11)
	w.Unit = time.Millisecond
	w.RockBlocks = true
	w.applyEcon()
	w.Start()
	w.Launch()
	ids := []int{0, w.Spawn(), w.Spawn(), w.Spawn()}
	for _, id := range ids {
		if id < 0 {
			t.Fatal("could not launch fleet")
		}
	}
	tickerDone := make(chan struct{})
	go func() { defer close(tickerDone); w.colonyTicker() }()
	reader, writer := io.Pipe()
	commandsDone := make(chan struct{})
	go func() { defer close(commandsDone); w.colonyCommands(reader) }()
	t.Cleanup(func() {
		writer.Close()
		<-commandsDone
		reader.Close()
		w.mu.Lock()
		w.done = true
		w.mu.Unlock()
		<-tickerDone
	})
	stopQueries, queriesDone := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(queriesDone)
		for {
			select {
			case <-stopQueries:
				return
			default:
				w.Neighbors(Point{})
				w.StepCost(Point{})
				w.BFS(Point{})
				w.Cheapest(Point{})
			}
		}
	}()
	defer func() { close(stopQueries); <-queriesDone }()

	var drones sync.WaitGroup
	for _, id := range ids {
		drones.Add(1)
		go func() {
			defer drones.Done()
			for i := 0; i < 100; i++ {
				w.Move(id, East)
				w.Move(id, West)
				w.Wait(id)
				w.Carry(id)
				w.CargoFree(id)
				w.Battery(id)
				w.Storage()
				w.Power()
				w.RankGoal()
				w.Research()
				w.StormClear()
				w.Fault(id)
			}
		}()
	}
	for i := 0; i < 30; i++ {
		upgrade := "engine"
		if i%10 == 0 {
			upgrade = "field"
		}
		if _, err := fmt.Fprintf(writer, "{\"buy\":%q}\n", upgrade); err != nil {
			t.Fatal(err)
		}
		w.Wait(0)
	}
	writer.Close()
	<-commandsDone
	dronesDone := make(chan struct{})
	go func() { drones.Wait(); close(dronesDone) }()
	select {
	case <-dronesDone:
	case <-time.After(5 * time.Second):
		t.Fatal("concurrent colony actions did not finish")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if st.Level("engine") != 27 || st.Level("field") != 3 || st.RP == 0 {
		t.Fatalf("purchases or research stalled: engine=%d field=%d rp=%g", st.Level("engine"), st.Level("field"), st.RP)
	}
	bought, snapshots := 0, 0
	for _, line := range bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'}) {
		var ev struct {
			Op    string
			OK    bool
			Stats *runStats
		}
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Op == "bought" && ev.OK {
			bought++
		}
		if ev.Op == "econ" && ev.Stats != nil {
			snapshots++
			if len(ev.Stats.Drones) != len(ids) {
				t.Fatalf("snapshot lost drones: %+v", ev.Stats)
			}
		}
	}
	if bought != 30 || snapshots == 0 {
		t.Fatalf("missing purchases or ticker snapshots: bought=%d snapshots=%d", bought, snapshots)
	}
}
