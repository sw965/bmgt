package bmgt

import (
	omwslices "github.com/sw965/omw/slices"
)

func NewMagicalStoneExcavationTarget0Actions(state *State) Actions {
	cs := CategoriesOfCards(state.P1.Graveyard)
	idxs := omwslices.IndicesFunc(cs, IsSpellCategory)
	y := make(Actions, len(idxs))
	for i, idx := range idxs {
		action := Action{
			CardName:MAGICAL_STONE_EXCAVATION,
			CardIDs:CardIDs{state.P1.Graveyard[idx].ID},
			Type:TARGET_ACTION,
		}
		y[i] = action
	}
	return y
}

func NewDarkFactoryOfMassProductionTarget0Actions(state *State) Actions {
	darkIdxs := omwslices.IndicesFunc(state.P1.Graveyard, IsDarkMonsterCard)
	idxss := omwslices.Combination[[][]int](darkIdxs, 2)
	y := make(Actions, len(idxss))
	for i, idxs := range idxss {
		cards := omwslices.IndicesAccess(state.P1.Graveyard)(idxs)
		ids := GetIDsOfCards(cards)
		action := Action{
			CardName:DARK_FACTORY_OF_MASS_PRODUCTION,
			CardIDs:ids,
			Type:TARGET_ACTION,
		}
		y[i] = action
	}
	return y
}