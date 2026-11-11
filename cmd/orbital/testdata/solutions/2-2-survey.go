package main

import "orbital/colony"

func main() {
	w, h := colony.WorldSize()
	colony.Report(w * h)
	colony.Report(2 * (w + h))
}
