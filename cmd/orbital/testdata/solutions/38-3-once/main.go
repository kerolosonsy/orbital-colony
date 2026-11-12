package main

import "sync"

type Calibration struct {
	once  sync.Once
	value int
}

func (c *Calibration) Value(load func() int) int {
	c.once.Do(func() { c.value = load() })
	return c.value
}
func LazyName(load func() string) func() string { return sync.OnceValue(load) }
func main()                                     {}
