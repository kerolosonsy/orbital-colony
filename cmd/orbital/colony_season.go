package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"orbital/colony/econ"
)

type colonySeason struct {
	main *econ.State
	path string
}

type seasonView struct {
	ID    string  `json:"id"`
	Left  int64   `json:"left"`
	Score float64 `json:"score"`
}

func (c *colonyGame) seasonView() *seasonView {
	season := econ.FindSeason(c.state.Season)
	if season == nil {
		return nil
	}
	return &seasonView{ID: season.ID, Left: max(0, season.Ticks-c.state.Tick), Score: c.state.Earned}
}

func (c *colonyGame) seasonSummary(summary map[string]any) map[string]any {
	if view := c.seasonView(); view != nil {
		summary["season"] = view
	}
	return summary
}

func (c *colonyGame) seasonBests() (map[string]float64, error) {
	bests := map[string]float64{}
	body, err := os.ReadFile(filepath.Join(filepath.Dir(c.path), "seasons.json"))
	if errors.Is(err, os.ErrNotExist) {
		return bests, nil
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &bests)
	if bests == nil {
		bests = map[string]float64{}
	}
	return bests, err
}

func seasonBestKey(season *econ.Season) string { return fmt.Sprintf("%s:%d", season.ID, season.Seed) }

func (s *server) registerColonySeasons(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/colony/seasons", s.colonySeasons)
	mux.HandleFunc("POST /api/colony/season", s.colonySeason)
}

func (s *server) colonySeasons(w http.ResponseWriter, r *http.Request) {
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	bests, err := c.seasonBests()
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	seasons := make([]map[string]any, 0, len(econ.Seasons))
	lang := requestLang(r)
	for _, season := range econ.Seasons {
		seasons = append(seasons, map[string]any{"id": season.ID, "title": tr(lang, season.Title, season.TitleEn), "desc": tr(lang, season.Desc, season.DescEn), "ticks": season.Ticks, "best": bests[seasonBestKey(&season)], "reference": season.Reference})
	}
	writeJSON(w, map[string]any{"seasons": seasons, "active": c.seasonView()})
}

func (s *server) colonySeason(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	if err := decoder.Decode(&req); err != nil {
		requestError(w, r, "bad request", 400)
		return
	}
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.ops.Lock()
	defer c.ops.Unlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active != nil {
		requestError(w, r, "Stop the colony first", http.StatusConflict)
		return
	}
	if req.ID != "" && econ.FindSeason(req.ID) == nil {
		requestError(w, r, tr(requestLang(r), "موسم مش معروف", "Unknown season"), 400)
		return
	}
	if err := c.selectSeason(req.ID); err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"ok": true, "active": c.seasonView()})
}

func (c *colonyGame) selectSeason(id string) error {
	if c.season != nil {
		if err := c.persist(); err != nil {
			return err
		}
	}
	if id == "" {
		if c.season != nil {
			c.state, c.path = c.season.main, c.season.path
			c.season = nil
		}
		return nil
	}
	season := econ.FindSeason(id)
	path := filepath.Join(filepath.Dir(c.path), "season-"+id+".json")
	st := season.New()
	body, err := os.ReadFile(path)
	if err == nil {
		var saved econ.State
		if err := json.Unmarshal(body, &saved); err != nil {
			return err
		}
		if saved.Season == id && saved.Seed == season.Seed && saved.Tick < season.Ticks {
			saved.Normalize()
			st = &saved
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	encoded, err := json.Marshal(st)
	if err != nil {
		return err
	}
	if err := atomicJSON(path, encoded); err != nil {
		return err
	}
	if c.season == nil {
		c.season = &colonySeason{main: c.state, path: c.path}
	}
	c.state, c.path = st, path
	return nil
}

// Called after the final snapshot, with mu held; completion is broadcast after unlock.
func (c *colonyGame) finishSeason(lang string) string {
	season := econ.FindSeason(c.state.Season)
	if season == nil || c.state.Tick < season.Ticks {
		return ""
	}
	bests, err := c.seasonBests()
	if err == nil {
		key := seasonBestKey(season)
		bests[key] = max(bests[key], c.state.Earned)
		var encoded []byte
		encoded, err = json.Marshal(bests)
		if err == nil {
			err = atomicJSON(filepath.Join(filepath.Dir(c.path), "seasons.json"), encoded)
		}
	}
	ar, en := "الموسم خلص واتحفظت النتيجة", "Season finished and score saved"
	if err != nil {
		ar, en = "تعذر حفظ أحسن نتيجة: "+err.Error(), "Could not save best score: "+err.Error()
	}
	line, _ := json.Marshal(map[string]any{"op": "season", "ok": err == nil, "score": c.state.Earned, "best": bests[seasonBestKey(season)], "msg": tr(lang, ar, en), "translations": [][2]string{{ar, en}}})
	return string(line)
}

func (c *colonyGame) publishSeason(line string) {
	if line != "" {
		c.broadcast(line)
	}
}
