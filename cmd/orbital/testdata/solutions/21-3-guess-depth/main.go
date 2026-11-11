package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func depthReply(text string, secret int, used int) (string, int, bool) {
	depth, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		return "invalid\n", used, false
	}
	if depth < 1 || depth > 64 {
		return "range\n", used, false
	}
	used++
	if depth == secret {
		return fmt.Sprintf("found=%d attempts=%d\n", secret, used), used, true
	}
	hint := "deeper\n"
	if depth > secret {
		hint = "shallower\n"
	}
	if used == 4 {
		hint += fmt.Sprintf("lost=%d\n", secret)
	}
	return hint, used, used == 4
}
func GuessDepth(in io.Reader, out io.Writer, draw func(int) int) error {
	if draw == nil {
		return fmt.Errorf("source required")
	}
	sample := draw(64)
	if sample < 0 || sample >= 64 {
		return fmt.Errorf("source outside range")
	}
	secret, used := sample+1, 0
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 4096), (128<<10)+2)
	for scanner.Scan() {
		if len(scanner.Text()) > 128<<10 {
			return fmt.Errorf("depth line too long")
		}
		reply, next, done := depthReply(scanner.Text(), secret, used)
		used = next
		if _, err := fmt.Fprint(out, reply); err != nil {
			return err
		}
		if done {
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	_, err := fmt.Fprintln(out, "aborted")
	return err
}
func main() {}
