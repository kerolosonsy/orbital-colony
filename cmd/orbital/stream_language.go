package main

import (
	"encoding/json"
	"sort"
	"strings"
)

// Player-written stdout and stderr stay verbatim. Only game event narration
// carries translation pairs, and those private fields are removed before SSE.
func localizeStreamLine(line, lang string) string {
	var event map[string]json.RawMessage
	if json.Unmarshal([]byte(line), &event) != nil {
		return line
	}
	var op string
	json.Unmarshal(event["op"], &op)
	if op == "stdout" || op == "stderr" {
		return line
	}
	var pairs [][2]string
	if blob := event["translations"]; blob != nil {
		if json.Unmarshal(blob, &pairs) != nil {
			return line
		}
	}
	delete(event, "translations")
	for _, field := range []string{"msg", "progress"} {
		if op == "report" && field == "msg" {
			continue
		}
		var text string
		if json.Unmarshal(event[field], &text) != nil {
			continue
		}
		text = streamText(lang, text, pairs)
		event[field], _ = json.Marshal(text)
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		return line
	}
	return string(encoded)
}

func streamText(lang, text string, pairs [][2]string) string {
	source, target := 0, 1
	if lang != "en" {
		source, target = 1, 0
	}
	sort.Slice(pairs, func(i, j int) bool { return len(pairs[i][source]) > len(pairs[j][source]) })
	replacements := make([]string, 0, len(pairs)*2)
	for _, pair := range pairs {
		if pair[source] != "" {
			replacements = append(replacements, pair[source], pair[target])
		}
	}
	text = strings.NewReplacer(replacements...).Replace(text)
	if text == "program exited" || text == "البرنامج خلص" {
		return tr(lang, "البرنامج خلص", "program exited")
	}
	return economyText(lang, text)
}
