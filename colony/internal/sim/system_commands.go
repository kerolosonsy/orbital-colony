package sim

import (
	"encoding/json"
	"math"

	"orbital/colony/econ"
)

type systemReceipt struct {
	ID      string  `json:"id"`
	OK      bool    `json:"ok"`
	Credits float64 `json:"credits"`
	Error   string  `json:"error,omitempty"`
}

type economyCommand struct {
	ID     string             `json:"id"`
	Credit *float64           `json:"credit"`
	Take   map[string]float64 `json:"take"`
	Market string             `json:"market"`
}

func finitePositive(n float64) bool { return n > 0 && !math.IsNaN(n) && !math.IsInf(n, 0) }

func (w *World) commandRefusal(cmd economyCommand) string {
	if cmd.Credit != nil && (*cmd.Credit < 0 || math.IsNaN(*cmd.Credit) || math.IsInf(*cmd.Credit, 0)) {
		return say("مبلغ غير صالح", "Invalid credit amount")
	}
	if cmd.Market != "" && (!w.Econ.Has("sys.market") || len(cmd.Take) != 1) {
		return say("السوق مقفول أو الطلب غير صالح", "Market locked or invalid order")
	}
	for good, amount := range cmd.Take {
		if econ.BasePrice[good] <= 0 || !finitePositive(amount) {
			return say("سلعة أو كمية غير صالحة", "Invalid good or amount")
		}
		if w.Econ.Inv[good] < amount {
			return say("المخزون مش كفاية", "Insufficient stock")
		}
	}
	return ""
}

func (w *World) applySystemCommand(cmd economyCommand) systemReceipt {
	receipt := systemReceipt{ID: cmd.ID}
	if cmd.Market != "" {
		if previous, ok := w.marketKeys[cmd.Market]; ok {
			previous.ID = cmd.ID
			return previous
		}
	}
	if receipt.Error = w.commandRefusal(cmd); receipt.Error != "" {
		return receipt
	}
	if cmd.Credit != nil {
		receipt.Credits = *cmd.Credit
	}
	for good, amount := range cmd.Take {
		if cmd.Market != "" {
			receipt.Credits = 1.2 * w.Econ.Price(good) * amount
			if !finitePositive(receipt.Credits) {
				receipt.Error = say("مبلغ غير صالح", "Invalid credit amount")
				return receipt
			}
		}
	}
	if math.IsInf(w.Econ.Credits+receipt.Credits, 0) || math.IsInf(w.Econ.Earned+receipt.Credits, 0) {
		receipt.Error = say("مبلغ غير صالح", "Invalid credit amount")
		return receipt
	}
	for good, amount := range cmd.Take {
		w.Econ.Inv[good] -= amount
		if cmd.Market != "" {
			w.archiveSale(good, amount, receipt.Credits)
			w.Econ.RecordSale(-1, good, amount, receipt.Credits, receipt.Credits/amount)
			w.recordRankSale(good, amount)
		}
	}
	w.Econ.Credits += receipt.Credits
	w.Econ.Earned += receipt.Credits
	receipt.OK = true
	if cmd.Market != "" {
		if w.marketKeys == nil {
			w.marketKeys = map[string]systemReceipt{}
		}
		w.marketKeys[cmd.Market] = receipt
		w.completeMarketContract()
	}
	w.syncInv()
	return receipt
}

func (w *World) completeMarketContract() {
	c, _ := w.Econ.ActiveContract()
	if c != nil && len(w.marketKeys) >= 20 && (c.ID == "h24-bulk-seller" || c.ID == "h33-stubborn-client") {
		w.Econ.RecordSubmission(c.Goal.Submit, true)
		w.checkContract()
	}
}

func (w *World) systemCommand(line []byte) bool {
	var cmd economyCommand
	if json.Unmarshal(line, &cmd) != nil || cmd.Credit == nil && cmd.Take == nil && cmd.Market == "" {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		return true
	}
	w.clock()
	w.advanceProduction(w.Tick)
	receipt := w.applySystemCommand(cmd)
	w.syncEcon()
	w.emitRaw(econEvent{Op: "system-command", T: w.Tick, OK: receipt.OK, Credits: w.Econ.Credits, State: w.Econ, Receipt: &receipt})
	return true
}

func (w *World) emitArchiveStatus() {
	if a := w.archive; a != nil {
		body, _ := json.Marshal(map[string]any{"op": "systems", "systems": map[string]any{"archive": map[string]any{"on": a.err == nil, "rows": a.rows}}})
		w.out.Write(append(body, '\n'))
	}
}
