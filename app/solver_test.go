package main

import (
	"maps"
	"testing"
)

func TestSolve(t *testing.T) {
	table := []struct {
		order int
		packs []int
		want  map[int]int
	}{
		// Just a simple cases
		{order: 100, packs: []int{10, 20, 30}, want: map[int]int{30: 3, 10: 1}},
		{order: 50, packs: []int{5, 10, 20}, want: map[int]int{10: 1, 20: 2}},

		// Some edge cases
		{order: 46, packs: []int{3, 23}, want: map[int]int{23: 2}},
		{order: 500_000, packs: []int{23, 31, 53}, want: map[int]int{23: 2, 31: 7, 53: 9429}},
	}

	for _, tt := range table {
		got := Solve(tt.order, tt.packs)
		if !maps.Equal(got, tt.want) {
			t.Errorf("Solve(%d, %v) = %v, want %v", tt.order, tt.packs, got, tt.want)
		}
	}
}
