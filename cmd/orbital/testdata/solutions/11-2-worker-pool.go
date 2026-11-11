package main

import (
	"orbital/colony"
	"sync"
)

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
}

func worker(d *colony.Drone, jobs <-chan colony.Point, results chan<- int) {
	for p := range jobs {
		goTo(d, p)
		n, _ := d.TryHarvest()
		results <- n
	}
}

func main() {
	jobs := make(chan colony.Point)
	results := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		d := colony.Spawn()
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(d, jobs, results)
		}()
	}
	go func() {
		for _, t := range colony.Targets() {
			jobs <- t
		}
		close(jobs)
	}()
	go func() {
		wg.Wait()
		close(results)
	}()
	total := 0
	for n := range results {
		total += n
	}
	colony.Report(total)
}
