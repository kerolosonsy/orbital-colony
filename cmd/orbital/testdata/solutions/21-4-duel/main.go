package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Maneuver int

const (
	Pulse Maneuver = iota
	Shield
	Drill
)

type DuelScore struct{ Wins, Losses, Draws int }

func ParseManeuver(text string) (Maneuver, error) {
	switch strings.ToLower(strings.TrimSpace(text)) {
	case "pulse":
		return Pulse, nil
	case "shield":
		return Shield, nil
	case "drill":
		return Drill, nil
	}
	return Pulse, fmt.Errorf("invalid maneuver")
}
func ResolveDuel(player, rival Maneuver) (string, error) {
	if player < Pulse || player > Drill || rival < Pulse || rival > Drill {
		return "", fmt.Errorf("invalid maneuver")
	}
	if player == rival {
		return "draw", nil
	}
	if player == Pulse && rival == Drill || player == Drill && rival == Shield || player == Shield && rival == Pulse {
		return "win", nil
	}
	return "lose", nil
}
func duelLine(text string, draw func(int) int, score DuelScore) (string, string, error) {
	if strings.ToLower(strings.TrimSpace(text)) == "score" {
		return fmt.Sprintf("score=%d/%d/%d\n", score.Wins, score.Losses, score.Draws), "", nil
	}
	player, err := ParseManeuver(text)
	if err != nil {
		return "invalid\n", "", nil
	}
	sample := draw(3)
	if sample < 0 || sample >= 3 {
		return "", "", fmt.Errorf("source outside range")
	}
	rival := Maneuver(sample)
	verdict, err := ResolveDuel(player, rival)
	if err != nil {
		return "", "", err
	}
	names := []string{"pulse", "shield", "drill"}
	return fmt.Sprintf("player=%s rival=%s verdict=%s\n", names[player], names[rival], verdict), verdict, nil
}
func Duel(in io.Reader, out io.Writer, draw func(int) int) (DuelScore, error) {
	var score DuelScore
	if draw == nil {
		return score, fmt.Errorf("source required")
	}
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1024), 4098)
	for scanner.Scan() {
		if len(scanner.Text()) > 4096 {
			return score, fmt.Errorf("duel line too long")
		}
		reply, verdict, err := duelLine(scanner.Text(), draw, score)
		if err != nil {
			return score, err
		}
		if _, err := fmt.Fprint(out, reply); err != nil {
			return score, err
		}
		switch verdict {
		case "win":
			score.Wins++
		case "lose":
			score.Losses++
		case "draw":
			score.Draws++
		}
	}
	return score, scanner.Err()
}
func main() {}
