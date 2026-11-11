package main

import (
	"fmt"
	"strconv"
)

type CommandToken struct {
	Kind, Text string
	Offset     int
}
type CommandError struct {
	Offset  int
	Message string
}

func (e *CommandError) Error() string { return fmt.Sprintf("byte %d: %s", e.Offset, e.Message) }
func flightLetter(b byte) bool        { return b >= 'a' && b <= 'z' }
func flightDigit(b byte) bool         { return b >= '0' && b <= '9' }
func LexFlight(source string) ([]CommandToken, error) {
	if len(source) > 4096 {
		return nil, &CommandError{4096, "program too long"}
	}
	var tokens []CommandToken
	for i := 0; i < len(source); {
		start := i
		switch b := source[i]; {
		case b == ' ' || b == '\t' || b == '\n' || b == '\r':
			i++
			continue
		case flightLetter(b):
			i++
			for i < len(source) && (flightLetter(source[i]) || flightDigit(source[i]) || source[i] == '_') {
				i++
			}
			tokens = append(tokens, CommandToken{"name", source[start:i], start})
		case flightDigit(b):
			i++
			for i < len(source) && flightDigit(source[i]) {
				i++
			}
			n, err := strconv.ParseInt(source[start:i], 10, 64)
			if err != nil || n > 1000 {
				return nil, &CommandError{start, "literal exceeds 1000"}
			}
			tokens = append(tokens, CommandToken{"number", source[start:i], start})
		case b == '+' || b == '-' || b == '*' || b == '/' || b == '(' || b == ')' || b == '=' || b == ';':
			tokens = append(tokens, CommandToken{string(b), string(b), i})
			i++
		default:
			return nil, &CommandError{i, "unknown byte"}
		}
	}
	return append(tokens, CommandToken{"eof", "", len(source)}), nil
}
func main() {}
