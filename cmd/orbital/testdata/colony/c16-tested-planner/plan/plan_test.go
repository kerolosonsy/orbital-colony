package plan

import "testing"

func TestNext(t *testing.T) {
	for _, tc := range []struct {
		ready []int
		want  int
	}{
		{nil, -1}, {[]int{4}, 0}, {[]int{5, 2, 2, 8}, 1}, {[]int{3, 1, 0}, 2},
	} {
		if got := Next(tc.ready); got != tc.want {
			t.Errorf("Next(%v)=%d want %d", tc.ready, got, tc.want)
		}
	}
}
