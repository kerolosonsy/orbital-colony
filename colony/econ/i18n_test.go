package econ

import (
	"orbital/colony/i18n"
	"testing"
)

func TestEnglishEconomyFeedback(t *testing.T) {
	previous := i18n.Lang()
	i18n.SetLang("en")
	t.Cleanup(func() { i18n.SetLang(previous) })
	cases := []struct {
		name     string
		feedback func(*State) string
		want     string
	}{
		{"unknown upgrade", func(s *State) string { return s.Buy("missing").Error() }, `No upgrade named "missing"`},
		{"maximum level", func(s *State) string {
			u := Find("field")
			s.Levels[u.ID] = u.Max
			_, _, why := s.Status(u)
			return why
		}, "All levels completed"},
		{"missing prerequisite", func(s *State) string { _, _, why := s.Status(Find("hangar")); return why }, "Requires first: Hangar License"},
		{"missing credits", func(s *State) string { _, _, why := s.Status(Find("field")); return why }, "You need 30 more credits"},
		{"missing goods", func(s *State) string {
			s.Credits = 1000
			s.Levels["crop.lichen"] = 1
			s.Blueprints["crop.glowcap"] = true
			_, _, why := s.Status(Find("crop.glowcap"))
			return why
		}, "You need 40 more Lichen"},
		{"missing blueprint", func(s *State) string { _, _, why := s.Status(Find("api.try")); return why }, "Finish the contract first: Careful Harvest"},
		{"claimed reward", func(s *State) string { s.Claimed["chapter-1"] = true; _, err := s.Claim(1); return err.Error() }, "This reward has already been claimed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.feedback(New(1)); got != tc.want {
				t.Errorf("feedback = %q, want %q", got, tc.want)
			}
		})
	}
	for good, want := range map[string]string{"moss": "Moss", "lichen": "Lichen", "glowcap": "Glowcap", "reed": "Pulse Reed", "crystal": "Crystal", "iron": "Iron", "geode": "Geode"} {
		if got := CropName(good); got != want {
			t.Errorf("CropName(%q) = %q, want %q", good, got, want)
		}
	}
}
