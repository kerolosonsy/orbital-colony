package main

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"orbital/colony/levels"
)

func TestSurveyStepsStayLocalizedInCatalogAndReview(t *testing.T) {
	l := levels.Get("6-1-survey")
	original, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	s, handler := profileServer(t)
	save := `{"progress":{"6-1-survey":{"stars":1},"6-2-top-three":{"stars":1},"6-3-census":{"stars":1},"6-4-stock-order":{"stars":1},"6-5-heat-grid":{"stars":1}}}`
	if response := profileRequest(t, handler, "PUT", "/api/profiles/default/progress", "default", save); response.Code != 200 {
		t.Fatal(response.Body)
	}
	for _, lang := range []string{"en", "ar", "en"} {
		want := l.Steps
		if lang == "en" {
			want = l.En.Steps
		}
		request := httptest.NewRequest("GET", "/api/levels", nil)
		request.Header.Set("X-Orbital-Lang", lang)
		recorder := httptest.NewRecorder()
		s.levels(recorder, request)
		var catalog struct{ Levels []*levels.Level }
		if err := json.Unmarshal(recorder.Body.Bytes(), &catalog); err != nil {
			t.Fatal(err)
		}
		for _, mission := range catalog.Levels {
			if mission.ID == l.ID && (!reflect.DeepEqual(mission.Steps, want) || mission.LessonRevision != l.LessonRevision) {
				t.Fatalf("%s catalog steps differ", lang)
			}
		}
		reviewRequest := httptest.NewRequest("POST", "http://localhost/api/profiles/default/review/start", strings.NewReader("{}"))
		reviewRequest.Header.Set("X-Orbital-Token", "test-token")
		reviewRequest.Header.Set("X-Orbital-Lang", lang)
		reviewResponse := httptest.NewRecorder()
		handler.ServeHTTP(reviewResponse, reviewRequest)
		var review levels.Level
		if err := json.Unmarshal(reviewResponse.Body.Bytes(), &review); err != nil {
			t.Fatal(err)
		}
		if reviewResponse.Code != 200 || review.ID != l.ID || review.Starter != "" || review.En != nil || review.LessonRevision != l.LessonRevision || !reflect.DeepEqual(review.Steps, want) {
			t.Fatalf("%s review: %s", lang, reviewResponse.Body)
		}
	}
	after, _ := json.Marshal(l)
	if string(after) != string(original) {
		t.Fatal("localization mutated shared Arabic/English lesson data")
	}
}
