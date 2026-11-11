package main

import "fmt"

func heatShape(grid [][]float64) (int, int, error) {
	if len(grid) == 0 {
		return 0, 0, fmt.Errorf("empty heat grid")
	}
	cols := len(grid[0])
	for i, row := range grid {
		if len(row) != cols {
			return 0, 0, fmt.Errorf("ragged row %d: got %d columns, want %d", i, len(row), cols)
		}
	}
	if cols == 0 {
		return 0, 0, fmt.Errorf("empty heat grid")
	}
	return len(grid), cols, nil
}
func heatSpace(rows, cols int) [][]float64 {
	grid := make([][]float64, rows)
	for i := range grid {
		grid[i] = make([]float64, cols)
	}
	return grid
}
func OverlayHeat(a, b [][]float64) ([][]float64, error) {
	rows, cols, err := heatShape(a)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	br, bc, err := heatShape(b)
	if err != nil {
		return nil, fmt.Errorf("overlay: %w", err)
	}
	if rows != br || cols != bc {
		return nil, fmt.Errorf("overlay shapes %dx%d and %dx%d differ", rows, cols, br, bc)
	}
	out := heatSpace(rows, cols)
	for i := range out {
		for j := range out[i] {
			out[i][j] = a[i][j] + b[i][j]
		}
	}
	return out, nil
}
func ThroughHeat(a, b [][]float64) ([][]float64, error) {
	rows, cols, err := heatShape(a)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	br, bc, err := heatShape(b)
	if err != nil {
		return nil, fmt.Errorf("coupler: %w", err)
	}
	if cols != br {
		return nil, fmt.Errorf("source columns %d differ from coupler rows %d", cols, br)
	}
	out := heatSpace(rows, bc)
	for i := range out {
		for j := range out[i] {
			for k := 0; k < cols; k++ {
				out[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return out, nil
}
func FlipHeat(a [][]float64) ([][]float64, error) {
	rows, cols, err := heatShape(a)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	out := heatSpace(cols, rows)
	for i := range a {
		for j, n := range a[i] {
			out[j][i] = n
		}
	}
	return out, nil
}
func main() {}
