package econ

import "math"

type PowerStatus struct {
	Supply float64 `json:"supply"`
	Demand float64 `json:"demand"`
	Day    bool    `json:"day"`
}

type StorageStatus struct {
	Used int `json:"used"`
	Cap  int `json:"cap"`
}

type DroneState struct {
	Fault   bool               `json:"fault,omitempty"`
	Stunned int                `json:"stunned,omitempty"`
	ID      int                `json:"id"`
	Battery int                `json:"battery"`
	Cargo   map[string]float64 `json:"cargo"`
	Cap     int                `json:"cap"`
}

const BaselinePower = 7
const BatteryCapacity = 100
const BuildingRent = 0.5
const DroneRent = 0.5

func StorageGroup(good string) string {
	switch good {
	case "moss", "lichen", "glowcap", "reed":
		return "crops"
	}
	if TierOf(good) > 0 {
		return "crafted"
	}
	return "raw"
}

func (s *State) StorageUsage() map[string]int {
	amounts := map[string]float64{"crops": 0, "raw": 0, "crafted": 0}
	for _, good := range sortedKeys(s.Inv) {
		amounts[StorageGroup(good)] += s.Inv[good]
	}
	used := map[string]int{}
	for group, n := range amounts {
		used[group] = int(math.Ceil(n - 1e-9))
	}
	return used
}

func (s *State) RefreshStorage() {
	s.Storage = map[string]StorageStatus{}
	for group, used := range s.StorageUsage() {
		capacity := s.StorageBase[group]
		for _, b := range s.Buildings {
			if b.Kind == "silo" {
				capacity += siloCapacity(group)
			}
		}
		s.Storage[group] = StorageStatus{Used: used, Cap: capacity}
	}
}

func siloCapacity(group string) int {
	if group == "crops" {
		return 200
	}
	return 100
}

func (s *State) Store(goods map[string]float64) error {
	amounts := map[string]float64{}
	for _, good := range sortedKeys(s.Inv) {
		amounts[StorageGroup(good)] += s.Inv[good]
	}
	for _, good := range sortedKeys(goods) {
		amounts[StorageGroup(good)] += goods[good]
	}
	s.RefreshStorage()
	for group, n := range amounts {
		if n > float64(s.Storage[group].Cap)+1e-9 {
			return ErrFull
		}
	}
	for good, n := range goods {
		s.Inv[good] += n
	}
	s.RefreshStorage()
	return nil
}

func (s *State) CargoCapacity() int { return 10 + int(linearMult(5, s.Level("cargo"), 8)-1) }

func (s *State) RentRate(drones int) float64 {
	return float64(len(s.Buildings))*BuildingRent + float64(max(0, drones-1))*DroneRent
}

func (s *State) PayRent(drones int) bool {
	s.Rent = s.RentRate(drones)
	cost := s.Rent / 600
	if s.Credits+1e-9 < cost {
		s.Credits = 0
		return false
	}
	s.Credits = math.Max(0, s.Credits-cost)
	return true
}

func (s *State) normalizeInfrastructure() {
	if s.StorageBase == nil {
		s.StorageBase = map[string]int{"crops": 2000, "raw": 1000, "crafted": 1000}
		for group, used := range s.StorageUsage() {
			s.StorageBase[group] = max(s.StorageBase[group], used)
		}
	}
	for _, b := range s.Buildings {
		if b.Kind == "generator" && b.Mode == "" {
			b.Mode = "solar"
		}
		if b.Kind == "collector" && b.Recipe == "" {
			b.Recipe = "helium"
		}
	}
	if len(s.Fleet) == 0 {
		s.Fleet = []DroneState{{Battery: BatteryCapacity, Cargo: map[string]float64{}, Cap: s.CargoCapacity()}}
	}
	for i := range s.Fleet {
		s.Fleet[i].ID = i
		s.Fleet[i].Battery = min(BatteryCapacity, max(0, s.Fleet[i].Battery))
		s.Fleet[i].Cap = s.CargoCapacity()
		if s.Fleet[i].Cargo == nil {
			s.Fleet[i].Cargo = map[string]float64{}
		}
	}
	s.RefreshStorage()
	s.RefreshPower(s.Tick)
	s.Rent = s.RentRate(len(s.Fleet))
}
