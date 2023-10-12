package bmgt

import (
	omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
)

func RoyalMagicalLibraryCost1(state *State, action *Action) {
	idx := action.MonsterZoneIndices1[0]
	state.P1.MonsterZone[idx].SpellCounter = 0
}

func ThunderDragonCost0(state *State, action *Action) {
	state.P1.Discard(action.HandIndices)
}

func ToonWorldCost0(state *State, action *Action) {
	state.P1.LifePoint -= 1000
}

func MagicalStoneExcavationCost0(state *State, action *Action) {
	state.P1.Discard(action.HandIndices)
}

func NewRoyalMagicalLibraryCost1Actions(state *State, idx int) Actions {
	action := Action{
		CardName:ROYAL_MAGICAL_LIBRARY,
		MonsterZoneIndices1:[]int{idx},
		Type:COST_ACTION,
	}
	return Actions{action}
}

func NewThunderDragonCost0Actions(state *State, idx int) Actions {
	action := Action{
		CardName:THUNDER_DRAGON,
		HandIndices:[]int{idx},
		Type:COST_ACTION,
	}
	return Actions{action}
}

func NewToonWorldCost0Actions(state *State) Actions {
	action := Action{
		CardName:TOON_WORLD,
		Type:COST_ACTION,
	}
	return Actions{action}
}

func NewMagicalStoneExcavationCost0Actions(state *State) Actions {
	c := omwmath.Combination{N:len(state.P1.Hand), R:2}
	idxss := c.Get()
	y := make(Actions, len(idxss))
	for i, idxs := range idxss {
		action := Action{
			CardName:MAGICAL_STONE_EXCAVATION,
			HandIndices:idxs,
			Type:COST_ACTION,
		}
		y[i] = action
	}
	return y
}