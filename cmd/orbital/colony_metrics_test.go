package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"

	"orbital/colony/econ"
)

func TestMetricsParsesTextAndRejectsAmbiguousSamples(t *testing.T) {
	for _, tc := range []struct {
		body string
		ok   bool
	}{
		{"# HELP colony_credits Current credits\n# TYPE colony_credits gauge\ncolony_credits 698.5\ncolony_goods_sold_total{good=\"moss\",site=\"a b\"} 5 1200\n", true},
		{`metric{label="a\"}b", path="c\\d", newline="a\nb"} 1`, true},
		{"metric 1\nmetric 2", false},
		{`metric{x="1",x="2"} 1`, false},
		{"metric{x=\"1\",y=\"2\"} 1\nmetric{y=\"2\",x=\"1\"} 2", false},
		{"metric NaN", false}, {"metric +Inf", false}, {"metric 1 2 3", false},
		{`metric{label="unterminated} 2`, false}, {`metric{x=1} 2`, false}, {`metric{x="\t"} 2`, false},
		{"metric 1 timestamp", false}, {"bad-name 1", false},
	} {
		t.Run(tc.body, func(t *testing.T) {
			metrics, err := parseColonyMetrics(tc.body)
			if (err == nil) != tc.ok {
				t.Fatalf("parsed %v error %v", metrics, err)
			}
		})
	}
}

func TestMetricsToleranceAndLifetimeLabelBound(t *testing.T) {
	for _, tc := range []struct {
		actual, expected float64
		ok               bool
	}{{102, 100, true}, {98, 100, true}, {102.01, 100, false}, {0, 0, true}, {0.001, 0, false}, {math.NaN(), 100, false}, {math.Inf(1), 100, false}} {
		if metricNear(tc.actual, tc.expected) != tc.ok {
			t.Fatalf("%g near %g", tc.actual, tc.expected)
		}
	}
	c := &colonyGame{}
	control := &colonyControl{labels: map[string]map[string]bool{}}
	for index := range 18 {
		metrics, err := parseColonyMetrics(fmt.Sprintf(`metric{drone="%d"} 1`, index))
		if err != nil {
			t.Fatal(err)
		}
		if c.acceptMetrics(control, metrics) != (index < metricsLabelLimit) {
			t.Fatalf("label bound at %d", index)
		}
	}
	if len(control.labels["metric"]) != metricsLabelLimit {
		t.Fatal("rejected labels accumulated")
	}
	for index := range 80 {
		c.acceptMetrics(control, map[string][]metricSample{fmt.Sprintf("m%d", index): {{value: 1}}})
	}
	if len(control.labels) > 64 {
		t.Fatal("rejected metric names accumulated")
	}
}

func TestMetricsSeriesBoundAndSnapshotIsolation(t *testing.T) {
	st := econ.New(1)
	st.Levels["sys.metrics"] = 1
	control := &colonyControl{metrics: true, series: map[string][][2]float64{}}
	c := &colonyGame{state: st, control: control}
	totals := map[string]float64{"colony_credits": 100, "colony_goods_sold_total": 1, "colony_drone_commands_total": 2, "colony_buildings": 3}
	for index := range 12 {
		totals[fmt.Sprintf("extra%d", index)] = float64(index)
	}
	for tick := range 75 {
		st.Tick = int64(tick)
		c.storeMetrics(control, totals)
	}
	if len(control.series) != 8 {
		t.Fatalf("series=%d", len(control.series))
	}
	for name, points := range control.series {
		if len(points) != 60 || points[0][0] != 15 || points[59][0] != 74 {
			t.Fatalf("%s: %v", name, points)
		}
	}
	c.mu.Lock()
	summary := c.controlSummary()
	c.mu.Unlock()
	series := summary["metrics"].(map[string]any)["series"].(map[string][][2]float64)
	series["colony_credits"][0][1] = -1
	if control.series["colony_credits"][0][1] != 100 {
		t.Fatal("summary shares live series")
	}
}

func TestMetricsScrapeRequiresMatchingAuthoritativeValues(t *testing.T) {
	for _, tc := range []struct {
		body string
		ok   bool
	}{
		{"colony_credits 100\ncolony_goods_sold_total 50\ncolony_drone_commands_total 200\ncolony_buildings 0", true},
		{"colony_credits 102\ncolony_goods_sold_total 50\ncolony_drone_commands_total 196\ncolony_buildings 0", true},
		{"colony_credits 103\ncolony_goods_sold_total 50\ncolony_drone_commands_total 200\ncolony_buildings 0", false},
		{"colony_credits 100\ncolony_goods_sold_total 50\ncolony_drone_commands_total 200", false},
		{"colony_credits 100\ncolony_goods_sold_total 50\ncolony_drone_commands_total 200\ncolony_buildings 0\nextra{x=\"a\"} 1e308\nextra{x=\"b\"} 1e308", false},
	} {
		st := econ.New(11)
		st.Credits = 100
		st.Systems.Sold = 50
		st.Systems.Commands = 200
		c := &colonyGame{state: st}
		control := &colonyControl{contract: "o36-dashboard", ctx: context.Background(), addr: "127.0.0.1:1", labels: map[string]map[string]bool{}, series: map[string][][2]float64{}}
		control.client = &http.Client{Transport: controlTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(tc.body)), Header: make(http.Header)}, nil
		})}
		if c.scrapeMetrics(control) != tc.ok {
			t.Fatalf("scrape judged incorrectly: %s", tc.body)
		}
		if !tc.ok && len(control.series) > 0 {
			t.Fatal("invalid scrape was stored")
		}
	}
}

func TestMetricsCollectsCustomSeriesOutsideDashboardContract(t *testing.T) {
	c := &colonyGame{state: econ.New(11)}
	control := &colonyControl{ctx: context.Background(), addr: "127.0.0.1:1", labels: map[string]map[string]bool{}, series: map[string][][2]float64{}}
	control.client = &http.Client{Transport: controlTransport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("colony_temperature 14")), Header: make(http.Header)}, nil
	})}
	if !c.scrapeMetrics(control) || len(control.series["colony_temperature"]) != 1 {
		t.Fatal("custom metric was not collected", control.series)
	}
}
