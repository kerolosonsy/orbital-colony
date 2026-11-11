package main

import (
	"errors"
	"sync"
)

type FuelAccount struct {
	mu             sync.Mutex
	capacity, fuel int64
}

func NewFuelAccount(capacity, initial int64) (*FuelAccount, error) {
	if capacity <= 0 || initial < 0 || initial > capacity {
		return nil, errors.New("invalid tank")
	}
	return &FuelAccount{capacity: capacity, fuel: initial}, nil
}
func (a *FuelAccount) Load(n int64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if n <= 0 || n > a.capacity-a.fuel {
		return errors.New("invalid load")
	}
	a.fuel += n
	return nil
}
func (a *FuelAccount) Burn(n int64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if n <= 0 || n > a.fuel {
		return errors.New("invalid burn")
	}
	a.fuel -= n
	return nil
}
func (a *FuelAccount) Remaining() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.fuel
}
func main() {}
