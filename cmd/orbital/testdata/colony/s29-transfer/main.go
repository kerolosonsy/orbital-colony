package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"orbital/colony"
)

func main() {
	stop := colony.Shutdown()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	var ready atomic.Bool
	ready.Store(true)
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if !ready.Load() {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(200)
	})
	setup(mux)
	var orders sync.Mutex
	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		var order struct {
			Good   string
			Amount int
		}
		if json.NewDecoder(r.Body).Decode(&order) != nil || order.Good != "moss" || order.Amount <= 0 {
			w.WriteHeader(400)
			return
		}
		if delay, _ := strconv.Atoi(r.Header.Get("X-Orbital-Delay-Ms")); delay > 0 && delay <= 1000 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
		orders.Lock()
		defer orders.Unlock()
		if colony.Count(colony.Moss) < order.Amount {
			w.WriteHeader(409)
			return
		}
		colony.Export(colony.Moss, order.Amount)
		w.WriteHeader(202)
	})
	srv := &http.Server{Addr: os.Getenv("ORBITAL_LISTEN_ADDR"), Handler: mux, ReadHeaderTimeout: time.Second}
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Print(err)
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				colony.Wait()
			}
		}
	}()
	<-stop
	ready.Store(false)
	time.Sleep(150 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Print(err)
	}
	// Keep fd 4 alive until the game observes the last export and drain response.
	time.Sleep(300 * time.Millisecond)
}
