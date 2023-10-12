package bmgt

import (
	"golang.org/x/exp/slices"
)

func NewLegacyOfYataGarasuEffectSelect0Actions(state *State) Actions {
	y := make(Actions, 0, 2)
	action := Action{
		CardName:LEGACY_OF_YATA_GARASU,
		EffectSelectNumber:0,
		Type:EFFECT_SELECT_ACTION,
	}
	y = append(y, action)

	cs := CategoriesOfCards(state.P1.MonsterZone)
	if slices.ContainsFunc(cs, IsSpiritMonsterCategory) {
		action := Action{
			CardName:LEGACY_OF_YATA_GARASU,
			EffectSelectNumber:1,
			Type:EFFECT_SELECT_ACTION,
		}
		y = append(y, action)
	}
	return y
}