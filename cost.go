package bmgt

import (
	omwslices "github.com/sw965/omw/slices"
	omwmath "github.com/sw965/omw/math"
)

func NewMagicalStoneExcavationActionss(duel *Duel) Actionss {
	n := len(duel.P1.Hand)
	if n < MAGICAL_STONE_EXCAVATION_COST {
		return Actionss{}
	}
	c := omwmath.Combination{N:n, R:MAGICAL_STONE_EXCAVATION_COST}
	cost0 := make(Actions, 0, c.TotalNum())

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			cost0 = append(cost0, Action{N1:1<<(1) + 1<<(j)})
		}
	}
	return Actionss{cost0}
}

type Cost func(*Duel, *Action)
type Costs []Cost

func NewMagicalStoneExcavationCosts() Costs {
	cost0 := func(duel *Duel, action *Action) {
		idxs := omwslices.Indices(omwslices.Binary(action.N1), 1)
		duel.P1.Discard(idxs)
	}
	return Costs{cost0}
}