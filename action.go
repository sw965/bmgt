package bmgt

import (
	"fmt"
)

type Action struct {
	CardName CardName

	HandIndices []int
	MonsterZoneIndices []int
	SpellTrapZoneIndices []int

	IsNormalSummon bool
	IsTributeSummonCost bool

	IsCardActivation bool
	IsCost bool
	IsSpellTrapSet bool
}

func (action *Action) Println() {
	fmt.Println("CardName =", CARD_NAME_TO_STRING[action.CardName])
	fmt.Println("HandIndices =", action.HandIndices)
	fmt.Println("MonsterZoneIndices =", action.MonsterZoneIndices)
	fmt.Println("SpellTrapZoneIndices =", action.SpellTrapZoneIndices)
	fmt.Println("IsNormalSummon =", action.IsNormalSummon)
	fmt.Println("IsTributeSummonCost =", action.IsTributeSummonCost)
	fmt.Println("IsCardActivation =", action.IsCardActivation)
	fmt.Println("IsCost =", action.IsCost)
	fmt.Println("IsSpellTrapSet =", action.IsSpellTrapSet)
}

type Actions []Action