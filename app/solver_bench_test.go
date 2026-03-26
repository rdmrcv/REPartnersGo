package main

import (
	"maps"
	"testing"
)

func BenchmarkSolve(b *testing.B) {
	b.Run(
		"small pack", func(b *testing.B) {
			var v map[int]int

			for b.Loop() {
				v, _ = Solve(100, []int{10, 20, 30})
			}

			if !maps.Equal(v, map[int]int{30: 3, 10: 1}) {
				b.FailNow()
			}
		},
	)

	b.Run(
		"large pack", func(b *testing.B) {
			var v map[int]int

			for b.Loop() {
				v, _ = Solve(500_000, []int{23, 31, 53})
			}

			if !maps.Equal(v, map[int]int{23: 2, 31: 7, 53: 9429}) {
				b.FailNow()
			}
		},
	)
}
