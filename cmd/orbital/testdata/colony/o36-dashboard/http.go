package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"orbital/colony"
)

func setup(mux *http.ServeMux) {
	mux.HandleFunc("GET /inventory", func(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(colony.Inventory()) })
	registry := prometheus.NewRegistry()
	for name, read := range map[string]func() float64{
		"colony_credits":              func() float64 { return colony.SystemStats().Credits },
		"colony_goods_sold_total":     func() float64 { return colony.SystemStats().Sold },
		"colony_drone_commands_total": func() float64 { return float64(colony.SystemStats().Commands) },
		"colony_buildings":            func() float64 { return float64(colony.SystemStats().Buildings) },
	} {
		registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: name, Help: name}, read))
	}
	mux.Handle("GET /metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
}
