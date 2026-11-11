package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	moss := colony.Count(colony.Moss)
	crystals := colony.Count(colony.Crystal)
	fmt.Println(moss, crystals)
	colony.Report(moss + crystals)
}
