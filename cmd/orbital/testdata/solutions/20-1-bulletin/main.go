package main

import (
	"fmt"
	"html"
	"strings"
)

func RenderBulletin(source string) (string, error) {
	var out strings.Builder
	mode := ""
	for i := 0; i < len(source); {
		if strings.HasPrefix(source[i:], "[alert]") || strings.HasPrefix(source[i:], "[code]") {
			if mode != "" {
				return "", fmt.Errorf("nested tag")
			}
			if strings.HasPrefix(source[i:], "[alert]") {
				mode = "alert"
				out.WriteString("<strong>")
				i += 7
			} else {
				mode = "code"
				out.WriteString("<code>")
				i += 6
			}
			continue
		}
		if strings.HasPrefix(source[i:], "[/alert]") || strings.HasPrefix(source[i:], "[/code]") {
			end := "[/" + mode + "]"
			if mode == "" || !strings.HasPrefix(source[i:], end) {
				return "", fmt.Errorf("unexpected close")
			}
			if mode == "alert" {
				out.WriteString("</strong>")
			} else {
				out.WriteString("</code>")
			}
			i += len(end)
			mode = ""
			continue
		}
		if source[i] == '\n' && mode == "" {
			out.WriteString("<br>\n")
		} else {
			out.WriteString(html.EscapeString(source[i : i+1]))
		}
		i++
	}
	if mode != "" {
		return "", fmt.Errorf("unclosed tag")
	}
	return out.String(), nil
}
func main() {}
