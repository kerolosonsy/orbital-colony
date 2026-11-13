package main

import (
	"bufio"
	"fmt"
	"math"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

const metricsLabelLimit = 16

var metricName = regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)
var metricLabelName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

type metricSample struct {
	labels map[string]string
	value  float64
}

func parseMetricLabels(text string) (map[string]string, error) {
	labels := map[string]string{}
	for strings.TrimSpace(text) != "" {
		name, rest, ok := strings.Cut(strings.TrimSpace(text), "=")
		name, rest = strings.TrimSpace(name), strings.TrimSpace(rest)
		if !ok || !metricLabelName.MatchString(name) || len(rest) == 0 || rest[0] != '"' || len(labels) >= 8 {
			return nil, fmt.Errorf("invalid metric labels")
		}
		end := 1
		for end < len(rest) {
			if rest[end] == '\\' {
				if end+1 >= len(rest) || !strings.ContainsRune("\\\"n", rune(rest[end+1])) {
					return nil, fmt.Errorf("invalid metric escape")
				}
				end += 2
				continue
			}
			if rest[end] == '"' {
				break
			}
			end++
		}
		if end >= len(rest) {
			return nil, fmt.Errorf("unterminated metric label")
		}
		value, err := strconv.Unquote(rest[:end+1])
		if _, duplicate := labels[name]; duplicate || err != nil || len(value) > 128 {
			return nil, fmt.Errorf("invalid metric label value")
		}
		labels[name] = value
		text = strings.TrimSpace(rest[end+1:])
		if text == "" {
			break
		}
		if text[0] != ',' {
			return nil, fmt.Errorf("invalid metric label separator")
		}
		text = text[1:]
	}
	return labels, nil
}

func parseMetricLine(line string) (string, metricSample, error) {
	end := strings.IndexAny(line, "{ \t")
	if end < 0 {
		return "", metricSample{}, fmt.Errorf("missing metric value")
	}
	name, rest := line[:end], strings.TrimSpace(line[end:])
	if !metricName.MatchString(name) {
		return "", metricSample{}, fmt.Errorf("invalid metric name")
	}
	sample := metricSample{labels: map[string]string{}}
	if strings.HasPrefix(rest, "{") {
		closing, quoted, escaped := -1, false, false
		for i, char := range rest {
			if escaped {
				escaped = false
				continue
			}
			if char == '\\' && quoted {
				escaped = true
				continue
			}
			if char == '"' {
				quoted = !quoted
			}
			if char == '}' && !quoted {
				closing = i
				break
			}
		}
		if closing < 0 {
			return "", sample, fmt.Errorf("unterminated metric labels")
		}
		labels, err := parseMetricLabels(rest[1:closing])
		if err != nil {
			return "", sample, err
		}
		sample.labels, rest = labels, strings.TrimSpace(rest[closing+1:])
	}
	fields := strings.Fields(rest)
	if len(fields) < 1 || len(fields) > 2 {
		return "", sample, fmt.Errorf("invalid metric sample")
	}
	value, err := strconv.ParseFloat(fields[0], 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return "", sample, fmt.Errorf("non-finite metric value")
	}
	if len(fields) == 2 {
		if _, err := strconv.ParseInt(fields[1], 10, 64); err != nil {
			return "", sample, err
		}
	}
	sample.value = value
	return name, sample, nil
}

func parseColonyMetrics(body string) (map[string][]metricSample, error) {
	metrics := map[string][]metricSample{}
	seen := map[string]bool{}
	scanner := bufio.NewScanner(strings.NewReader(body))
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, sample, err := parseMetricLine(line)
		if err != nil {
			return nil, err
		}
		keys := make([]string, 0, len(sample.labels))
		for label := range sample.labels {
			keys = append(keys, label)
		}
		sort.Strings(keys)
		identity := name
		for _, label := range keys {
			identity += "\x00" + label + "=" + strconv.Quote(sample.labels[label])
		}
		if seen[identity] || count >= 256 || len(metrics) >= 64 && metrics[name] == nil {
			return nil, fmt.Errorf("too many or duplicate metric samples")
		}
		seen[identity], count = true, count+1
		metrics[name] = append(metrics[name], sample)
	}
	return metrics, scanner.Err()
}

func metricNear(actual, expected float64) bool {
	return !math.IsNaN(actual) && !math.IsInf(actual, 0) && math.Abs(actual-expected) <= math.Max(1e-9, math.Abs(expected)*0.02)
}

func (c *colonyGame) metricsExpected() map[string]float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return map[string]float64{"colony_credits": c.state.Credits, "colony_goods_sold_total": c.state.Systems.Sold, "colony_drone_commands_total": float64(c.state.Systems.Commands), "colony_buildings": float64(len(c.state.Buildings))}
}

func (c *colonyGame) acceptMetrics(control *colonyControl, metrics map[string][]metricSample) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for name, samples := range metrics {
		if control.labels[name] == nil {
			if len(control.labels) >= 64 {
				return false
			}
			control.labels[name] = map[string]bool{}
		}
		for _, sample := range samples {
			for label, value := range sample.labels {
				identity := label + "\x00" + value
				if !control.labels[name][identity] && len(control.labels[name]) >= metricsLabelLimit {
					return false
				}
				control.labels[name][identity] = true
			}
		}
	}
	return true
}

func metricsTotals(metrics map[string][]metricSample) map[string]float64 {
	totals := map[string]float64{}
	for name, samples := range metrics {
		for _, sample := range samples {
			totals[name] += sample.value
		}
	}
	return totals
}

func (c *colonyGame) scrapeMetrics(control *colonyControl) bool {
	before := c.metricsExpected()
	status, body, elapsed, err := control.request(control.ctx, "/metrics", nil)
	metrics, parseErr := parseColonyMetrics(string(body))
	if err != nil || parseErr != nil || status != 200 || elapsed > controlLatency || !c.acceptMetrics(control, metrics) {
		return false
	}
	after, totals := c.metricsExpected(), metricsTotals(metrics)
	for _, total := range totals {
		if math.IsNaN(total) || math.IsInf(total, 0) {
			return false
		}
	}
	if control.contract == "o36-dashboard" {
		for name, expected := range before {
			actual, exists := totals[name]
			if !exists || !metricNear(actual, expected) && !metricNear(actual, after[name]) {
				return false
			}
		}
	}
	c.storeMetrics(control, totals)
	return true
}

func (c *colonyGame) storeMetrics(control *colonyControl, totals map[string]float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	names := []string{"colony_credits", "colony_goods_sold_total", "colony_drone_commands_total", "colony_buildings"}
	var extra []string
	for name := range totals {
		if !slices.Contains(names, name) {
			extra = append(extra, name)
		}
	}
	sort.Strings(extra)
	names = append(names, extra...)
	for _, name := range names {
		value, exists := totals[name]
		if !exists || len(control.series) >= 8 && control.series[name] == nil {
			continue
		}
		points := append(control.series[name], [2]float64{float64(c.state.Tick), value})
		if len(points) > 60 {
			points = points[len(points)-60:]
		}
		control.series[name] = points
	}
}

func (c *colonyGame) runMetrics(control *colonyControl) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	good := 0
	for {
		select {
		case <-control.ctx.Done():
			return
		case <-ticker.C:
			ok := c.scrapeMetrics(control)
			c.broadcastControl()
			if !ok && control.contract == "o36-dashboard" {
				c.controlVerdict(control, false)
			} else {
				good++
			}
			if good == 3 && control.contract == "o36-dashboard" {
				c.controlVerdict(control, c.probeDrain(control))
			}
		}
	}
}
