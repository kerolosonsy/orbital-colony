package main

import "orbital/colony"

type Score struct{ Correct, Attempted int }

func main() {
	score := Score{}
	for i, question := range colony.Exam() {
		if len(question.Choices) == 0 {
			continue
		}
		score.Attempted++
		ground := "soil"
		switch question.Prompt {
		case "lichen":
			ground = "ice"
		case "crystal":
			ground = "regolith"
		}
		for j, choice := range question.Choices {
			if choice == ground {
				if colony.Answer(i, j) {
					score.Correct++
				}
				break
			}
		}
	}
	percent := 0
	if score.Attempted > 0 {
		percent = score.Correct * 100 / score.Attempted
	}
	colony.Report(percent)
}
