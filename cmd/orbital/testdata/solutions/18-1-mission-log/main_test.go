package main

import (
	"reflect"
	"testing"
)

func TestMissionLog(t *testing.T) {
	for _, tc := range []struct {
		name  string
		erase int
	}{
		{"erase head", 1}, {"erase middle", 2}, {"erase tail", 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var log MissionLog
			if len(log.Entries()) != 0 {
				t.Fatal("fresh log")
			}
			for _, bad := range []string{"", " \t"} {
				if id, err := log.Log(bad); err == nil || id != 0 {
					t.Fatal("invalid sector")
				}
			}
			for i, sector := range []string{" ice ", "ore", "relay"} {
				if id, err := log.Log(sector); err != nil || id != i+1 {
					t.Fatal("issued ID", id, err)
				}
			}
			if log.Entries()[0].Sector != "ice" {
				t.Fatal("trim sector")
			}
			for n := 0; n < 2; n++ {
				if err := log.Seal(2); err != nil {
					t.Fatal(err)
				}
			}
			if !log.Entries()[1].Sealed {
				t.Fatal("seal did not persist")
			}
			before := log.Entries()
			for _, id := range []int{-1, 0, 99} {
				if log.Seal(id) == nil || log.Erase(id) == nil {
					t.Fatal("missing ID accepted")
				}
			}
			if !reflect.DeepEqual(log.Entries(), before) {
				t.Fatal("error mutated log")
			}
			if err := log.Erase(tc.erase); err != nil {
				t.Fatal(err)
			}
			var want []MissionEntry
			for _, entry := range before {
				if entry.ID != tc.erase {
					want = append(want, entry)
				}
			}
			if !reflect.DeepEqual(log.Entries(), want) {
				t.Fatal("erase order")
			}
			if id, err := log.Log("new"); err != nil || id != 4 {
				t.Fatal("reused ID", id, err)
			}
			snapshot := log.Entries()
			snapshot[0].Sector = "changed"
			if log.Entries()[0].Sector == "changed" {
				t.Fatal("snapshot aliases log")
			}
		})
	}
}
