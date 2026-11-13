package sim

import (
	"encoding/json"
	"math"

	"orbital/colony/econ"
)

type neighborTransfer struct {
	Credit   float64                `json:"credit"`
	Take     map[string]float64     `json:"take"`
	Receipt  string                 `json:"receipt"`
	Neighbor *econ.NeighborProgress `json:"neighbor,omitempty"`
}

func (w *World) neighborCommand(line []byte) bool {
	var command neighborTransfer
	if json.Unmarshal(line, &command) != nil || (command.Take == nil && command.Credit == 0 && command.Receipt == "") {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		return true
	}
	w.clock()
	w.advanceProduction(w.Tick)
	ok, msg := w.transferNeighbor(command)
	if ok {
		w.checkContract()
	}
	w.syncInv()
	w.econSnapshot("state", "", true, true)
	encoded, _ := json.Marshal(map[string]any{"op": "transfer", "id": command.Receipt, "ok": ok, "msg": msg, "credit": command.Credit})
	w.out.Write(append(encoded, '\n'))
	return true
}

func (w *World) transferNeighbor(command neighborTransfer) (bool, string) {
	st := w.Econ
	if command.Receipt != "" {
		if paid, ok := st.NeighborPayments[command.Receipt]; ok {
			return paid == command.Credit, ""
		}
	}
	if command.Credit < 0 || math.IsNaN(command.Credit) || math.IsInf(command.Credit, 0) {
		return false, say("مبلغ مش صالح", "Invalid credit amount")
	}
	for good, amount := range command.Take {
		if _, known := econ.PriceV3[good]; !known || amount <= 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
			return false, say("شحنة مش صالحة", "Invalid shipment")
		}
		if st.Inv[good] < amount {
			return false, say("المخزون مش كفاية", "Insufficient goods")
		}
	}
	for good, amount := range command.Take {
		st.Inv[good] -= amount
	}
	st.Credits += command.Credit
	st.Earned += command.Credit
	if command.Receipt != "" {
		if st.NeighborPayments == nil {
			st.NeighborPayments = map[string]float64{}
		}
		st.NeighborPayments[command.Receipt] = command.Credit
	}
	if command.Neighbor != nil && st.IsNeighborContract() {
		_, r := st.ActiveContract()
		r.Neighbor.Trades++
		r.Neighbor.Reconnects = max(r.Neighbor.Reconnects, command.Neighbor.Reconnects)
		r.Neighbor.MaxLine = max(r.Neighbor.MaxLine, command.Neighbor.MaxLine)
	}
	return true, ""
}
