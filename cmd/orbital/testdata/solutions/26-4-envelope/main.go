package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type RelayStatus struct {
	Name   string `json:"name"`
	Charge int    `json:"charge"`
}
type RelayFault struct {
	Kind  string
	Cause error
}

func (e *RelayFault) Error() string { return fmt.Sprintf("%s: %v", e.Kind, e.Cause) }
func (e *RelayFault) Unwrap() error { return e.Cause }

type relayProblem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type relayEnvelope struct {
	Data  *RelayStatus  `json:"data"`
	Error *relayProblem `json:"error"`
}

func relayFailure(err error) (int, *relayProblem) {
	var fault *RelayFault
	if errors.As(err, &fault) {
		switch fault.Kind {
		case "invalid":
			return 400, &relayProblem{"invalid", "invalid relay"}
		case "missing":
			return 404, &relayProblem{"missing", "relay missing"}
		case "blocked":
			return 409, &relayProblem{"blocked", "relay blocked"}
		}
	}
	return 500, &relayProblem{"internal", "relay unavailable"}
}
func relayResponse(w http.ResponseWriter, r *http.Request, status int, envelope relayEnvelope) {
	encoded, _ := json.Marshal(envelope)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if r.Method != "HEAD" {
		_, _ = w.Write(encoded)
	}
}
func NewRelayEnvelope(fetch func(context.Context, string) (RelayStatus, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/relays/")
		if !strings.HasPrefix(r.URL.Path, "/relays/") || name == "" || strings.Contains(name, "/") {
			relayResponse(w, r, 404, relayEnvelope{Error: &relayProblem{"missing", "relay missing"}})
			return
		}
		if r.Method != "GET" {
			w.Header().Set("Allow", "GET")
			relayResponse(w, r, 405, relayEnvelope{Error: &relayProblem{"method", "use GET"}})
			return
		}
		relay, err := fetch(r.Context(), name)
		if err != nil {
			status, problem := relayFailure(err)
			relayResponse(w, r, status, relayEnvelope{Error: problem})
			return
		}
		relayResponse(w, r, 200, relayEnvelope{Data: &relay})
	})
}
func main() {}
