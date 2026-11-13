package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

func setup(mux *http.ServeMux) {
	db, err := sql.Open("sqlite", os.Getenv("ORBITAL_DB_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA journal_mode=WAL;
CREATE TABLE IF NOT EXISTS batches(id TEXT PRIMARY KEY);
CREATE TABLE IF NOT EXISTS goods(good TEXT PRIMARY KEY, amount REAL NOT NULL CHECK(amount>=0));
CREATE TABLE IF NOT EXISTS stores(name TEXT PRIMARY KEY, amount INTEGER NOT NULL CHECK(amount>=0));
INSERT INTO stores(name,amount) VALUES ('a',100),('b',100) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS transfers(id TEXT PRIMARY KEY, source TEXT, destination TEXT, amount INTEGER);`); err != nil {
		log.Fatal(err)
	}
	mux.HandleFunc("GET /inventory", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT good,amount FROM goods")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer rows.Close()
		inventory := map[string]float64{}
		for rows.Next() {
			var good string
			var amount float64
			if rows.Scan(&good, &amount) != nil {
				w.WriteHeader(500)
				return
			}
			inventory[good] = amount
		}
		if rows.Err() != nil {
			w.WriteHeader(500)
			return
		}
		json.NewEncoder(w).Encode(inventory)
	})
	mux.HandleFunc("POST /inventory", func(w http.ResponseWriter, r *http.Request) {
		var batch struct {
			ID    string
			Goods map[string]float64
		}
		if json.NewDecoder(r.Body).Decode(&batch) != nil || batch.ID == "" {
			w.WriteHeader(400)
			return
		}
		for _, amount := range batch.Goods {
			if amount < 0 {
				w.WriteHeader(400)
				return
			}
		}
		tx, err := db.Begin()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer tx.Rollback()
		inserted, err := tx.Exec("INSERT INTO batches(id) VALUES (?) ON CONFLICT DO NOTHING", batch.ID)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		count, err := inserted.RowsAffected()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if count == 0 {
			w.WriteHeader(204)
			return
		}
		for good, amount := range batch.Goods {
			if _, err := tx.Exec("INSERT INTO goods(good,amount) VALUES (?,?) ON CONFLICT(good) DO UPDATE SET amount=amount+excluded.amount", good, amount); err != nil {
				w.WriteHeader(500)
				return
			}
		}
		if tx.Commit() != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(204)
	})
	mux.HandleFunc("GET /stores", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT name,amount FROM stores")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer rows.Close()
		stores := map[string]int{}
		for rows.Next() {
			var name string
			var amount int
			if rows.Scan(&name, &amount) != nil {
				w.WriteHeader(500)
				return
			}
			stores[name] = amount
		}
		if rows.Err() != nil {
			w.WriteHeader(500)
			return
		}
		json.NewEncoder(w).Encode(stores)
	})
	mux.HandleFunc("POST /transfer", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			ID, From, To string
			Amount       int
		}
		if json.NewDecoder(r.Body).Decode(&request) != nil || request.ID == "" || request.Amount <= 0 || request.From == request.To {
			w.WriteHeader(400)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer tx.Rollback()
		var source, destination string
		var amount int
		err = tx.QueryRow("SELECT source,destination,amount FROM transfers WHERE id=?", request.ID).Scan(&source, &destination, &amount)
		if err == nil {
			if source != request.From || destination != request.To || amount != request.Amount {
				w.WriteHeader(409)
				return
			}
			w.WriteHeader(200)
			return
		}
		if err != sql.ErrNoRows {
			w.WriteHeader(500)
			return
		}
		var balance int
		if tx.QueryRow("SELECT amount FROM stores WHERE name=?", request.From).Scan(&balance) != nil || balance < request.Amount {
			w.WriteHeader(409)
			return
		}
		if _, err := tx.Exec("UPDATE stores SET amount=amount-? WHERE name=?", request.Amount, request.From); err != nil {
			w.WriteHeader(500)
			return
		}
		changed, err := tx.Exec("UPDATE stores SET amount=amount+? WHERE name=?", request.Amount, request.To)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if count, err := changed.RowsAffected(); err != nil || count != 1 {
			w.WriteHeader(409)
			return
		}
		if _, err := tx.Exec("INSERT INTO transfers(id,source,destination,amount) VALUES (?,?,?,?)", request.ID, request.From, request.To, request.Amount); err != nil {
			w.WriteHeader(500)
			return
		}
		if tx.Commit() != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	})
}
