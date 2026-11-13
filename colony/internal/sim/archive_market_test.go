package sim

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"orbital/colony/econ"
)

func archiveWorld(t *testing.T, seed int64) (*World, *bytes.Buffer) {
	t.Helper()
	t.Setenv("ORBITAL_COLONY_DIR", t.TempDir())
	st := econ.New(seed)
	st.Levels["sys.archive"] = 1
	var trace bytes.Buffer
	w := NewColony(&trace, st, seed)
	w.Realtime = false
	return w, &trace
}

func readArchive(t *testing.T, w *World, name string) []byte {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(w.archive.dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func TestArchiveLedgerDeterministicMalformedRowsAndReload(t *testing.T) {
	generate := func(seed int64) (*World, []byte) {
		w, _ := archiveWorld(t, seed)
		for i := 0; i < 400; i++ {
			w.Tick = i
			w.archiveSale("moss", 2, 3.5)
		}
		return w, readArchive(t, w, "ledger.csv")
	}
	first, ledger := generate(71)
	_, again := generate(71)
	if !bytes.Equal(ledger, again) {
		t.Fatal("same seed produced different malformed rows")
	}
	_, other := generate(72)
	if bytes.Equal(ledger, other) {
		t.Fatal("malformed row pattern did not depend on seed")
	}
	r := csv.NewReader(bytes.NewReader(ledger))
	r.FieldsPerRecord = -1
	header, _ := r.Read()
	if strings.Join(header, ",") != strings.TrimSpace(ledgerHeader) {
		t.Fatal(header)
	}
	missing, nonnumeric, quoted, valid := 0, 0, 0, 0
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		switch {
		case len(row) != 5:
			missing++
		case row[2] == "bad-amount":
			nonnumeric++
		case row[1] == "moss,damaged":
			quoted++
		default:
			valid++
		}
	}
	if missing == 0 || nonnumeric == 0 || quoted == 0 || valid < 300 {
		t.Fatalf("malformed variants missing: %d %d %d valid=%d", missing, nonnumeric, quoted, valid)
	}
	if first.archive.rows != 400 || first.archive.credits["moss"] != float64(valid)*3.5 {
		t.Fatalf("truth includes malformed rows: %+v", first.archive)
	}
	t.Setenv("ORBITAL_COLONY_DIR", first.archive.dir)
	reloaded := NewColony(io.Discard, first.Econ, 71)
	if !reflect.DeepEqual(first.archive.credits, reloaded.archive.credits) || reloaded.archive.rows != 400 {
		t.Fatal("archive truth did not survive a new run")
	}
}

func TestArchiveHarvestDaysStormsFaultsAndSubmit(t *testing.T) {
	for _, correct := range []bool{false, true} {
		t.Run(fmt.Sprint(correct), func(t *testing.T) {
			w, trace := archiveWorld(t, 31)
			if err := w.Econ.Accept("d19-shift-report"); err != nil {
				t.Fatal(err)
			}
			want := map[string]map[string]float64{}
			for day := 0; day < 3; day++ {
				w.Tick = day * econ.ArchiveDay
				w.archiveTick(w.Tick)
				w.archiveHarvest("moss", float64(day+2))
				w.archiveHarvest("water", 1.5)
				w.Tick += 120
				w.archiveTick(w.Tick)
				want[fmt.Sprint(day)] = map[string]float64{"moss": float64(day + 2), "water": 1.5}
			}
			w.emit(Event{Op: "storm", Phase: "hit", T: w.Tick})
			w.emit(Event{Op: "fault", D: 0, T: w.Tick})
			w.Tick = 3600
			w.archiveTick(w.Tick)
			if w.Submit("wrong-key", want) {
				t.Fatal("accepted wrong key")
			}
			if !correct {
				want["0"]["moss"]++
			}
			if got := w.Submit("shift-report", want); got != correct {
				t.Fatalf("submit=%v", got)
			}
			run := w.Econ.Contracts["d19-shift-report"]
			if w.Econ.Active != "" || w.Econ.Won("d19-shift-report") != correct {
				t.Fatalf("did not settle: %+v", run)
			}
			body := readArchive(t, w, "events.jsonl")
			if !bytes.Contains(body, []byte(`"op":"storm"`)) || !bytes.Contains(body, []byte(`"op":"fault"`)) {
				t.Fatal(string(body))
			}
			if !bytes.Contains(body, []byte(`"d":0`)) {
				t.Fatal("fault lost drone zero")
			}
			if !bytes.Contains(trace.Bytes(), []byte(`"op":"contract"`)) {
				t.Fatal("missing contract event")
			}
		})
	}
}

func TestArchiveAuditSubmitRejectsWrongAndNonJSONAnswers(t *testing.T) {
	for _, answer := range []any{map[string]float64{"moss": 17}, map[string]int{"moss": 18}, make(chan int)} {
		w, _ := archiveWorld(t, 9)
		if err := w.Econ.Accept("d18-audit"); err != nil {
			t.Fatal(err)
		}
		w.archive.credits["moss"], w.archive.valid = 17, 1
		correct := reflect.DeepEqual(answer, map[string]float64{"moss": 17})
		if w.Submit("audit", answer) != correct || w.Econ.Won("d18-audit") != correct || w.Econ.Active != "" {
			t.Fatal("incorrect audit settlement")
		}
	}
}

func TestArchiveRotationCapsFilesAndRetainsHeaders(t *testing.T) {
	w, _ := archiveWorld(t, 6)
	for i := 0; i < 4; i++ {
		w.archiveSale("moss", 1, 2)
		w.archive.event(map[string]string{"op": "sale", "padding": strings.Repeat("x", archiveLimit-100)})
		if err := os.WriteFile(filepath.Join(w.archive.dir, "ledger.csv"), []byte(ledgerHeader+strings.Repeat("x", archiveLimit-len(ledgerHeader)-1)), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	w.archiveSale("moss", 1, 2)
	if w.archive.err != nil {
		t.Fatal(w.archive.err)
	}
	for _, name := range []string{"ledger.csv", "events.jsonl"} {
		for _, suffix := range []string{"", ".1", ".2"} {
			st, err := os.Stat(filepath.Join(w.archive.dir, name+suffix))
			if err != nil || st.Size() > archiveLimit {
				t.Fatalf("%s%s: %v %v", name, suffix, st, err)
			}
		}
		if _, err := os.Stat(filepath.Join(w.archive.dir, name+".3")); !os.IsNotExist(err) {
			t.Fatal("unbounded rotations")
		}
	}
	if !bytes.HasPrefix(readArchive(t, w, "ledger.csv"), []byte(ledgerHeader)) || w.archive.rows != 1 {
		t.Fatal("rotation lost header or current row count")
	}
}

func TestMarketCommandsConcurrentDuplicatesMoveGoodsAndCreditsOnce(t *testing.T) {
	w, trace := archiveWorld(t, 5)
	w.Econ.Levels["sys.market"] = 1
	w.Econ.Inv["moss"] = 100
	before, earned := w.Econ.Credits, w.Econ.Earned
	price := w.Econ.Price("moss")
	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			w.systemCommand([]byte(fmt.Sprintf(`{"id":"%d","market":"same","take":{"moss":20}}`, i)))
		}(i)
	}
	wg.Wait()
	if w.Econ.Inv["moss"] != 80 || w.Econ.Credits != before+24*price || w.Econ.Earned != earned+24*price {
		t.Fatalf("duplicate sale: stock=%v credits=%v", w.Econ.Inv["moss"], w.Econ.Credits-before)
	}
	for _, line := range bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'}) {
		var ev struct{ Receipt *systemReceipt }
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Receipt != nil && (!ev.Receipt.OK || ev.Receipt.Credits != 24*price) {
			t.Fatalf("different duplicate receipt: %+v", ev.Receipt)
		}
	}
}

func TestSystemCommandsRefuseWholeTakeAndLockedMarket(t *testing.T) {
	w, _ := archiveWorld(t, 8)
	w.Econ.Inv["moss"] = 10
	before := w.Econ.Credits
	for _, command := range []string{
		`{"credit":12.5,"take":{"moss":4,"lichen":1}}`,
		`{"credit":-1,"take":{"moss":1}}`,
		`{"credit":12,"take":{"unknown":1}}`,
		`{"market":"locked","take":{"moss":1}}`,
		`{"credit":1,"take":{"moss":-2}}`,
	} {
		w.systemCommand([]byte(command))
		if w.Econ.Inv["moss"] != 10 || w.Econ.Credits != before {
			t.Fatal("partially applied refusal", command)
		}
	}
	w.systemCommand([]byte(`{"credit":12.5,"take":{"moss":4}}`))
	if w.Econ.Inv["moss"] != 6 || w.Econ.Credits != before+12.5 {
		t.Fatal("valid atomic command failed")
	}
}

func TestArchiveAbsentNodeDoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ORBITAL_COLONY_DIR", dir)
	w := NewColony(io.Discard, econ.New(1), 1)
	w.Realtime = false
	if w.archive != nil || w.Submit("audit", map[string]int{}) {
		t.Fatal("archive enabled without node")
	}
	files, err := os.ReadDir(dir)
	if err != nil || len(files) != 0 {
		t.Fatal(files, err)
	}
}

func TestArchiveRealHarvestAndExportProduceParseableFiles(t *testing.T) {
	w, _ := archiveWorld(t, 12)
	w.ActionDur = 0
	w.Start()
	w.At(0, 0).Ground, w.At(0, 0).Entity = Soil, Moss
	if !w.Harvest(0) {
		t.Fatal("harvest failed")
	}
	if !bytes.Contains(readArchive(t, w, "events.jsonl"), []byte(`"op":"harvest"`)) {
		t.Fatal("harvest would be lost if main returned immediately")
	}
	earned := w.Export(0, Moss, 1)
	if earned <= 0 {
		t.Fatal("sale failed")
	}
	w.flushArchive()
	if body := readArchive(t, w, "events.jsonl"); !bytes.Contains(body, []byte(`"op":"sale"`)) || !bytes.Contains(body, []byte(`"op":"harvest"`)) {
		t.Fatal(string(body))
	}
	if w.archive.rows != 1 {
		t.Fatal("real sale not in ledger")
	}
}

func TestMarketContractsCountOnlyDistinctFilledOrders(t *testing.T) {
	for _, id := range []string{"h24-bulk-seller", "h33-stubborn-client"} {
		w, _ := archiveWorld(t, 55)
		w.Econ.Levels["sys.market"] = 1
		w.Econ.Inv["moss"] = 100
		before, price := w.Econ.Credits, w.Econ.Price("moss")
		if err := w.Econ.Accept(id); err != nil {
			t.Fatal(err)
		}
		for n := 0; n < 10; n++ {
			w.systemCommand([]byte(`{"market":"order-0","take":{"moss":1}}`))
		}
		w.systemCommand([]byte(`{"market":"refused","take":{"moss":1000}}`))
		for n := 1; n < 19; n++ {
			w.systemCommand([]byte(fmt.Sprintf(`{"market":"order-%d","take":{"moss":1}}`, n)))
		}
		if w.Econ.Won(id) || w.Econ.Active != id {
			t.Fatal("duplicate or refused order counted toward contract")
		}
		w.systemCommand([]byte(`{"market":"order-19","take":{"moss":1}}`))
		if !w.Econ.Won(id) || w.Econ.Active != "" || w.Econ.Inv["moss"] != 80 {
			t.Fatal("twentieth distinct fill did not settle")
		}
		if delta := w.Econ.Credits - before - econ.FindContract(id).Reward.Credits - 24*price; delta > 1e-8 || delta < -1e-8 {
			t.Fatal("incorrect bulk credits or duplicated reward", delta)
		}
		paid := w.Econ.Credits
		w.systemCommand([]byte(`{"market":"order-19","take":{"moss":1}}`))
		if w.Econ.Credits != paid || w.Econ.Inv["moss"] != 80 {
			t.Fatal("settled contract paid duplicate order")
		}
	}
}
