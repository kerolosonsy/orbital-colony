package main

import "orbital/colony"

type CargoKind int
type Cells int

const (
	Cell CargoKind = iota
	Pack
	Crate
)
const (
	CellFactor  Cells = 1
	PackFactor  Cells = 6
	CrateFactor Cells = 24
)

func main() {
	cells := Cells(colony.CargoAmount(int(Cell)))
	packs := Cells(colony.CargoAmount(int(Pack)))
	crates := Cells(colony.CargoAmount(int(Crate)))
	total := cells*CellFactor + packs*PackFactor + crates*CrateFactor
	colony.Report(total)
}
