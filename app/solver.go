package main

import (
	"math"
	"sort"
)

type row []int

func (r row) Sum() int {
	sum := 0

	for _, n := range r {
		// Do not sum values which not initialized
		if n == math.MaxInt {
			return math.MaxInt - 1
		}

		sum += n
	}

	return sum
}

func (r row) Exists() bool {
	return r.Sum() != math.MaxInt-1
}

type state []row

func NewState(order, packsLen int) state {
	state := make(state, order+1)

	for i := range state {
		state[i] = make([]int, packsLen)

		for j := range packsLen {
			if i == 0 {
				state[i][j] = 0

				continue
			}

			state[i][j] = math.MaxInt
		}
	}

	return state
}

func (s state) Sum(i int) (sum int) {
	return s[i].Sum()
}

func (s state) Set(to, from, targetIdx, v int) {
	copy(s[to], s[from])

	// is value set for one pack index - set all other as 0 at this point
	for idx := range s[to] {
		if targetIdx == idx {
			s[to][idx] = v
		}
	}
}

func (s state) Exists(idx int) bool {
	return s[idx].Exists()
}

func Solve(order int, packs []int) map[int]int {
	// Sanity check.
	if len(packs) == 0 || order == 0 {
		return nil
	}

	// If there is only one pack size — we cannot do anything here and just send as is.
	if len(packs) == 1 {
		packsAmount := order / packs[0]
		if order%packs[0] > 0 {
			packsAmount += 1
		}

		return map[int]int{
			packs[0]: packsAmount,
		}
	}

	// Sorted packs reduce oversent.
	sort.Ints(packs)

	// First, we need to cut a few possible variants:
	//
	//   1. We cannot just try to guess the best packs by sending the biggest packs
	//   (optimize packs amount) and then slice remaining iteratively.
	//
	//   2. We cannot also optimize oversend by picking the lowest remaining from
	//   packs (optimize oversent) and then iteratively split to the biggest packs (to
	//   optimize packs).
	//
	// Both those variants stuck when we needed to reduce the initial pick (it can
	// be suboptimal especially if pack sized is not dividable to each other). So,
	// instead of doing some heuristic, we need to go with dynamic programming here
	// and try to step-by-step find optimal distribution until the desired amount is
	// reached.

	// We need to keep track of the lowest number of packs on each num. This
	// unfortunately could be large, but if needed, we can try to optimize it with a
	// window (we do not need to store more than max pack size number of values).
	stateInst := NewState(order+packs[len(packs)-1], len(packs))

	for i := range stateInst {
		for packIdx, packSize := range packs {
			// Skip values where we cannot look behind long enough
			if i < packSize {
				continue
			}

			// Possible value if we do nothing
			noAction := stateInst.Sum(i)
			// Possible value if we act at this point to pack of packSize
			act := stateInst.Sum(i-packSize) + 1

			// If we pick value as optimal at this point — set it to the result and save in
			// stateInst.
			if act < noAction {
				stateInst.Set(i, i-packSize, packIdx, stateInst[i-packSize][packIdx]+1)
			}
		}
	}

	result := map[int]int{}

	for _, r := range stateInst[order:] {
		if r.Exists() {
			for idx, num := range r {
				if num != 0 {
					result[packs[idx]] = num
				}
			}

			return result
		}
	}

	return result
}
