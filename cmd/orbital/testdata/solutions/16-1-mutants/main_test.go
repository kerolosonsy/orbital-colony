package main

import "testing"

func TestFuel(t *testing.T) {
	cases := []struct{ d, c, want int }{
		{3, 12, 8},
		{0, 40, 0},
		{0, 0, 0},
		{2, 6, 5},
		{1, 50, 12},
		{1, 51, 22},
		{1, 200, 52},
	}
	for _, tc := range cases {
		if got := Fuel(tc.d, tc.c); got != tc.want {
			t.Errorf("Fuel(%d, %d) = %d, want %d", tc.d, tc.c, got, tc.want)
		}
	}
}
