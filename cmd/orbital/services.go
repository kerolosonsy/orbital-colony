package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var serviceDefaults = map[string]string{
	"ORBITAL_PG_DSN":        "postgres://orbital:orbital@127.0.0.1:5432/orbital?sslmode=disable",
	"ORBITAL_MONGO_URI":     "mongodb://127.0.0.1:27017",
	"ORBITAL_KAFKA_BROKERS": "127.0.0.1:9092",
}

func serviceEnv(name string) string {
	if configured := os.Getenv(name); configured != "" {
		return configured
	}
	return serviceDefaults[name]
}

func serviceAddr(service string) (string, error) {
	var endpoint, port string
	switch service {
	case "postgres":
		endpoint, port = serviceEnv("ORBITAL_PG_DSN"), "5432"
	case "mongo":
		endpoint, port = serviceEnv("ORBITAL_MONGO_URI"), "27017"
	case "kafka":
		return strings.Split(serviceEnv("ORBITAL_KAFKA_BROKERS"), ",")[0], nil
	default:
		return "", fmt.Errorf("unknown service %q", service)
	}
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	if parsed.Hostname() == "" {
		return "", fmt.Errorf("service URI must include a host")
	}
	if parsed.Port() != "" {
		port = parsed.Port()
	}
	return net.JoinHostPort(parsed.Hostname(), port), nil
}

func checkService(ctx context.Context, service string) error {
	addr, err := serviceAddr(service)
	if err != nil {
		return err
	}
	conn, err := (&net.Dialer{Timeout: 500 * time.Millisecond}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	return conn.Close()
}

func (s *server) services(w http.ResponseWriter, r *http.Request) {
	statuses := map[string]bool{}
	for _, service := range []string{"postgres", "mongo", "kafka"} {
		statuses[service] = checkService(r.Context(), service) == nil
	}
	writeJSON(w, statuses)
}
