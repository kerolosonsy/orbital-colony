package main

import "orbital/colony"

type Ledger struct{ TotalCents, People int }
type CrewUnit struct{ Ledger }

func (u CrewUnit) Shares() []int {
	shares := make([]int, u.People)
	for i := range shares {
		shares[i] = u.TotalCents / u.People
		if i < u.TotalCents%u.People {
			shares[i]++
		}
	}
	return shares
}
func main() {
	total, people := colony.Invoice()
	unit := CrewUnit{Ledger{total, people}}
	colony.Report(unit.Shares())
}
