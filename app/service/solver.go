package service

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

var (
	ErrCannotSolve = errors.New("cannot solve state")

	ErrParams = errors.New("invalid parameters")
)

// inf value to initialize values in the state
//
// -1 here needed because this allows us to do unconditional +1 which then will
// be replaced because any meaningful value will be lower than this. Without -1
// we should add an explicit check if the value is set.
const inf = math.MaxInt - 1

// Solve will solve packing for a specified order and packs configuration
//
// ErrCannotSolve error can be returned from this function only when we have a
// bug. With any package configuration we should be able to solve a packing task
// at least with massive oversend or vast amount of pack utilization. So treat
// this error as a strong signal to re-check logic and tests, not just a plain
// error.
func Solve(order int, packs []int) (map[int]int, error) {
	// Sanity check.
	if len(packs) == 0 || order == 0 {
		return nil, nil
	}

	if order < 0 {
		return nil, fmt.Errorf("invalid order: %d: %w", order, ErrParams)
	}

	sort.Ints(packs)

	if packs[0] <= 0 {
		return nil, fmt.Errorf("invalid pack: %d: %w", packs[0], ErrParams)
	}

	// If there is only one pack size, we cannot do anything here and just send as is.
	if len(packs) == 1 {
		packsAmount := order / packs[0]
		if order%packs[0] > 0 {
			packsAmount += 1
		}

		return map[int]int{
			packs[0]: packsAmount,
		}, nil
	}

	// How long we should go in the worst case
	iterLength := order + packs[len(packs)-1]

	// We need to keep track of the lowest number of packs on each num. This
	// unfortunately could be large, but if needed, we can try to optimize it with a
	// window (we do not need to store more than max pack size number of values).
	stateInst := NewState(iterLength, len(packs))

	result := map[int]int{}

	for i := range iterLength {
		for packIdx, packSize := range packs {
			// Skip values where we cannot look behind long enough
			if i < packSize {
				continue
			}

			// Possible value if we do nothing
			noAction := stateInst.Get(i).Sum()
			// Possible value if we act at this point to pack of packSize
			act := stateInst.Get(i-packSize).Sum() + 1

			// If we pick value as optimal at this point — set it to the result and save in
			// stateInst.
			if act < noAction {
				stateInst.Modify(i, i-packSize, packIdx, stateInst.Get(i - packSize)[packIdx]+1)
			}
		}

		// If we calculated enough state — return first exists distribution as a result.
		if i == order {
			if stateInst.Get(i).Exists() {
				for idx, num := range stateInst[i] {
					if num != 0 {
						result[packs[idx]] = num
					}
				}

				return result, nil
			}
		}
	}

	// If we somehow go here - there is a bug in the solver func.
	//
	// Might be worth panic here, but for now choose a safer approach and write
	// notice about proper monitoring.
	return nil, ErrCannotSolve
}

// row type represents the state in the moment of time
type row []int

// Sum returns sum packs amount if the row is initialized. Otherwise, it returns
// inf const.
func (r row) Sum() int {
	sum := 0

	for _, n := range r {
		// Do not sum values which not initialized
		if n == inf {
			return inf
		}

		sum += n
	}

	return sum
}

// Exists method checks if the row is initialized.
func (r row) Exists() bool {
	return r.Sum() != inf
}

// state provides an interface to access previous calculations. It may be a
// little bit overthinking, but since in naive implementation we should store
// history size of the order, we need to think about possible optimization. With
// this type we can minimize code changes if we decide to have a sliding window
// state (for example, implemented via ring-buffer) in the future.
type state []row

// NewState creates a state for DB calculation that allows us to look back.
func NewState(order, packsLen int) state {
	state := make(state, order+1)

	for i := range state {
		state[i] = make([]int, packsLen)

		for j := range packsLen {
			if i == 0 {
				state[i][j] = 0

				continue
			}

			state[i][j] = inf
		}
	}

	return state
}

// Get - lookup in the state at point i
func (s state) Get(i int) row {
	return s[i]
}

// Modify method modifies state and inherits previous values.
func (s state) Modify(to, from, targetIdx, v int) {
	// Explicitly copy to avoid ref-copy slice
	copy(s[to], s[from])

	s[to][targetIdx] = v
}
