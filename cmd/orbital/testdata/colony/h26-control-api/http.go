package main

import (
	"encoding/json"
	"net/http"
	"orbital/colony"
)

func setup(mux *http.ServeMux) {
	mux.HandleFunc("GET /inventory", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(colony.Inventory())
	})
}
