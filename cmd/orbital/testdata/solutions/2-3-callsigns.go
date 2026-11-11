package main

import (
	"orbital/colony"
	"unicode/utf8"
)

func main() {
	firstName := colony.Callsign(0)
	first, _ := utf8.DecodeRuneInString(firstName)
	last, _ := utf8.DecodeLastRuneInString(firstName)
	colony.Report(len(firstName))
	colony.Report(utf8.RuneCountInString(firstName))
	colony.Report(string(first))
	colony.Report(string(last))
	secondName := colony.Callsign(1)
	secondFirst, _ := utf8.DecodeRuneInString(secondName)
	secondLast, _ := utf8.DecodeLastRuneInString(secondName)
	colony.Report(len(secondName))
	colony.Report(utf8.RuneCountInString(secondName))
	colony.Report(string(secondFirst))
	colony.Report(string(secondLast))
}
