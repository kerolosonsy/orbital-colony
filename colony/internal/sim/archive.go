package sim

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"orbital/colony/econ"
)

const archiveLimit = 1 << 20
const ledgerHeader = "tick,good,amount,unit_price,credits\n"

type colonyArchive struct {
	dir     string
	rows    int
	valid   int
	credits map[string]float64
	days    map[string]map[string]float64
	pending map[string]float64
	day     int
	err     error
}

func (w *World) initArchive() {
	dir := os.Getenv("ORBITAL_COLONY_DIR")
	if !w.Econ.Has("sys.archive") || dir == "" {
		return
	}
	a := &colonyArchive{dir: dir, credits: map[string]float64{}, days: map[string]map[string]float64{}, pending: map[string]float64{}, day: w.Tick / econ.ArchiveDay}
	w.archive = a
	a.err = os.MkdirAll(dir, 0o700)
	if a.err == nil {
		a.err = a.append("ledger.csv", nil, ledgerHeader)
	}
	if a.err == nil {
		a.err = a.append("events.jsonl", nil, "")
	}
	if a.err == nil {
		a.loadLedger()
		a.loadDays()
	}
}

func archiveNumber(s string) (float64, bool) {
	n, err := strconv.ParseFloat(s, 64)
	return n, err == nil && !math.IsNaN(n) && !math.IsInf(n, 0) && n >= 0
}

func (a *colonyArchive) loadLedger() {
	f, err := os.Open(filepath.Join(a.dir, "ledger.csv"))
	if err != nil {
		a.err = err
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	_, _ = r.Read()
	for {
		row, err := r.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			var parseError *csv.ParseError
			if !errors.As(err, &parseError) {
				a.err = err
				return
			}
		}
		a.rows++
		if err != nil || len(row) != 5 {
			continue
		}
		_, tickOK := strconv.ParseInt(row[0], 10, 64)
		_, amountOK := archiveNumber(row[2])
		_, priceOK := archiveNumber(row[3])
		credits, creditsOK := archiveNumber(row[4])
		if tickOK == nil && amountOK && priceOK && creditsOK && econ.BasePrice[row[1]] > 0 {
			a.credits[row[1]] += credits
			a.valid++
		}
	}
}

type archiveHarvestEvent struct {
	Op    string             `json:"op"`
	T     int                `json:"t"`
	Day   int                `json:"day"`
	Goods map[string]float64 `json:"goods"`
}

func (a *colonyArchive) loadDays() {
	for _, suffix := range []string{".2", ".1", ""} {
		f, err := os.Open(filepath.Join(a.dir, "events.jsonl"+suffix))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			a.err = err
			return
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			var ev archiveHarvestEvent
			if json.Unmarshal(sc.Bytes(), &ev) == nil && ev.Op == "harvest" {
				for good, amount := range ev.Goods {
					a.addHarvest(ev.Day, good, amount)
				}
			}
		}
		a.err = sc.Err()
		f.Close()
		if a.err != nil {
			return
		}
	}
}

func (a *colonyArchive) append(name string, body []byte, header string) error {
	path := filepath.Join(a.dir, name)
	st, err := os.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if st != nil && st.Size()+int64(len(body)) > archiveLimit {
		if err := a.rotate(name); err != nil {
			return err
		}
		st = nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if st == nil || st.Size() == 0 {
		_, err = io.WriteString(f, header)
	}
	if err == nil {
		_, err = f.Write(body)
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func (a *colonyArchive) rotate(name string) error {
	path := filepath.Join(a.dir, name)
	if err := os.Rename(path+".1", path+".2"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.Rename(path, path+".1"); err != nil {
		return err
	}
	if name == "ledger.csv" {
		a.rows, a.valid = 0, 0
		clear(a.credits)
	}
	return nil
}

func (a *colonyArchive) event(ev any) {
	if a.err != nil {
		return
	}
	body, err := json.Marshal(ev)
	if err == nil {
		err = a.append("events.jsonl", append(body, '\n'), "")
	}
	a.err = err
}

func (w *World) archiveSale(good string, amount, credits float64) {
	a := w.archive
	if a == nil || a.err != nil || amount <= 0 {
		return
	}
	row := []string{strconv.Itoa(w.Tick), good, strconv.FormatFloat(amount, 'g', -1, 64), strconv.FormatFloat(credits/amount, 'g', -1, 64), strconv.FormatFloat(credits, 'g', -1, 64)}
	// Valid CSV quoting is a trap for strings.Split, while the unknown good makes that row invalid for the audit.
	variant := econ.EventHash(w.Econ.Seed, int64(w.Tick), uint64(a.rows)+1800) % 40
	switch variant {
	case 0:
		row = row[:4]
	case 1:
		row[2] = "bad-amount"
	case 2:
		row[1] += ",damaged"
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	_ = writer.Write(row)
	writer.Flush()
	a.err = a.append("ledger.csv", buf.Bytes(), ledgerHeader)
	if a.err != nil {
		return
	}
	a.rows++
	if variant > 2 {
		a.valid++
		a.credits[good] += credits
	}
	a.event(map[string]any{"op": "sale", "t": w.Tick, "good": good, "amount": amount, "unit_price": credits / amount, "credits": credits})
}

func (a *colonyArchive) addHarvest(day int, good string, amount float64) {
	key := strconv.Itoa(day)
	if a.days[key] == nil {
		a.days[key] = map[string]float64{}
	}
	a.days[key][good] += amount
}

func (w *World) archiveHarvest(good string, amount float64) {
	if a := w.archive; a != nil {
		a.addHarvest(w.Tick/econ.ArchiveDay, good, amount)
		a.pending[good] += amount
	}
}

func (w *World) flushArchive() {
	if a := w.archive; a != nil && len(a.pending) > 0 {
		a.event(archiveHarvestEvent{Op: "harvest", T: w.Tick, Day: a.day, Goods: a.pending})
		clear(a.pending)
	}
}

func (w *World) archiveTick(tick int) {
	a := w.archive
	if a == nil {
		return
	}
	if tick%120 == 0 {
		w.flushArchive()
	}
	a.day = tick / econ.ArchiveDay
	for day := range a.days {
		n, _ := strconv.Atoi(day)
		if n < a.day-3 {
			delete(a.days, day)
		}
	}
}

func (w *World) archiveEvent(ev Event) {
	// A player may return from main immediately after harvesting.
	w.flushArchive()
	if a := w.archive; a != nil && (ev.Op == "storm" || ev.Op == "fault") {
		a.event(colonyEventWire(ev))
	}
}

func (w *World) archiveAnswer(key string) (any, bool) {
	a := w.archive
	if a == nil || a.err != nil {
		return nil, false
	}
	switch key {
	case "audit":
		return a.credits, a.valid > 0
	case "shift-report":
		day := w.Tick / econ.ArchiveDay
		report := map[string]map[string]float64{}
		for n := day - 3; n < day; n++ {
			goods := a.days[strconv.Itoa(n)]
			if goods == nil {
				goods = map[string]float64{}
			}
			report[strconv.Itoa(n)] = goods
		}
		return report, day >= 3
	}
	return nil, false
}

func equalArchiveAnswer(want, got any) bool {
	wantJSON, err := json.Marshal(want)
	if err != nil {
		return false
	}
	gotJSON, err := json.Marshal(got)
	if err != nil {
		return false
	}
	var a, b any
	_ = json.Unmarshal(wantJSON, &a)
	_ = json.Unmarshal(gotJSON, &b)
	return equalArchiveJSON(a, b)
}

func equalArchiveJSON(a, b any) bool {
	switch x := a.(type) {
	case float64:
		y, ok := b.(float64)
		return ok && math.Abs(x-y) <= 1e-6*math.Max(1, math.Abs(x))
	case map[string]any:
		y, ok := b.(map[string]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for key, n := range x {
			if !equalArchiveJSON(n, y[key]) {
				return false
			}
		}
		return true
	}
	return false
}

func (w *World) Submit(key string, answer any) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil {
		return false
	}
	w.clock()
	w.advanceProduction(w.Tick)
	w.flushArchive()
	want, ready := w.archiveAnswer(key)
	if !ready {
		return false
	}
	correct := equalArchiveAnswer(want, answer)
	if !w.Econ.RecordSubmission(key, correct) {
		return false
	}
	w.checkContract()
	return correct
}
