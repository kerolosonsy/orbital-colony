package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var dockBays = []string{"north", "south", "east", "west", "other"}

type DockMetrics struct {
	reg      *prometheus.Registry
	landings *prometheus.CounterVec
	seconds  *prometheus.HistogramVec
	queue    prometheus.Gauge
	inflight *prometheus.GaugeVec
}

func NewDockMetrics(reg *prometheus.Registry) (*DockMetrics, error) {
	m := &DockMetrics{
		reg: reg,
		landings: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "dock_landings_total", Help: "Landings by bay and outcome.",
		}, []string{"bay", "outcome"}),
		seconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "dock_landing_seconds", Help: "Landing duration.", Buckets: []float64{1, 5, 15, 60},
		}, []string{"bay"}),
		queue: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "dock_queue_depth", Help: "Ships waiting to land.",
		}),
		inflight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dock_inflight", Help: "Landings in progress.",
		}, []string{"bay"}),
	}
	for _, c := range []prometheus.Collector{m.landings, m.seconds, m.queue, m.inflight} {
		if err := reg.Register(c); err != nil {
			return nil, err
		}
	}
	for _, bay := range dockBays {
		m.landings.WithLabelValues(bay, "ok")
		m.landings.WithLabelValues(bay, "aborted")
	}
	return m, nil
}

func dockBay(bay string) string {
	switch bay {
	case "north", "south", "east", "west":
		return bay
	}
	return "other"
}

func (m *DockMetrics) RecordLanding(bay string, ok bool, d time.Duration) {
	bay = dockBay(bay)
	outcome := "aborted"
	if ok {
		outcome = "ok"
	}
	m.landings.WithLabelValues(bay, outcome).Inc()
	m.seconds.WithLabelValues(bay).Observe(max(d, 0).Seconds())
}

func (m *DockMetrics) SetQueue(n int) {
	m.queue.Set(float64(max(n, 0)))
}

func (m *DockMetrics) Inflight(bay string) func() {
	g := m.inflight.WithLabelValues(dockBay(bay))
	g.Inc()
	var once sync.Once
	return func() { once.Do(g.Dec) }
}

func (m *DockMetrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})
}

func main() {}
