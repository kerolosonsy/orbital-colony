package main

import (
	"encoding/json"
	"go/format"
	"net/http"
	"strings"
)

// formatCode: POST /api/format {"code": "..."} → {"code": formatted} or
// {"error": "..."} when the code doesn't parse. It uses go/format (gofmt's
// own printer) in-process, so nothing is executed.
func (s *server) formatCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10)).Decode(&req); err != nil {
		requestError(w, r, "bad request", http.StatusBadRequest)
		return
	}
	out, err := format.Source([]byte(req.Code))
	if err != nil {
		writeJSON(w, map[string]string{"error": strings.TrimPrefix(err.Error(), "<standard input>:")})
		return
	}
	writeJSON(w, map[string]string{"code": string(out)})
}
