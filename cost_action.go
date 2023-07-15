package bmgt

import (
	"github.com/sw965/omw/fn"
	omws "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
)

type costActionss struct{}

var CostActionss = costActionss{}

// 王立魔法図書館
func (_ *costActionss) RoyalMagicalLibrary(state *State) Actionss {
	cardName := CardName("王立魔法図書館")
	var effect0 Actions
	if state.CanSpellSpeed1Activation() {
		actions := NewMonsterZoneIndexActions(cardName, 0, true)
		f := func(action Action) bool {
			card := state.P1.MonsterZone[action.MonsterZoneIndices[0]]
			return card.Name == cardName && card.SpellCounter == 3
		}
		effect0 = fn.Filter(actions, f)
	}
	return Actionss{effect0}
}

// サンダー・ドラゴン
func (_ *costActionss) ThunderDragon(state *State) Actionss {
	cardName := CardName("サンダー・ドラゴン")
	var effect0 Actions
	if state.CanSpellSpeed1Activation() && slices.ContainsFunc(state.P1.Deck, EqualNameCard(cardName)) {
		actions := NewHandIndexActions(cardName, len(state.P1.Hand), 0, true)
		f := func(action Action) bool {
			card := state.P1.Hand[action.HandIndices[0]]
			return card.Name == cardName
		}
		effect0 = fn.Filter(actions, f)
	}
	return Actionss{effect0}
}

// 召喚僧サモンプリースト
func (_ *costActionss) SummonerMonk(state *State) Actionss {
	effect0 := Actions{}
	effect1 := Actions{}
	effect2 := Actions{}
	if state.CanSpellSpeed1Activation() && slices.ContainsFunc(state.P1.Deck, IsLevel4MonsterCard) {
		cardName := CardName("召喚僧サモンプリースト")
		actions := NewHandIndexAndMonsterZoneIndexActions(cardName, len(state.P1.Hand), 2, true)
		f := func(action Action) bool {
			zoneCard := state.P1.MonsterZone[action.MonsterZoneIndices[0]]
			handCard := state.P1.Hand[action.HandIndices[0]]
			return zoneCard.Name == "召喚僧サモンプリースト" && zoneCard.ThisTurnEffectActivationCounts[2] == 0 && IsSpellCard(handCard)
		}
		effect2 = fn.Filter(actions, f)
	}
	return Actionss{effect0, effect1, effect2}
}

// 闇の量産工場
func (_ *costActionss) DarkFactoryOfMassProduction(state *State) Actionss {
	var effect0 Actions
	f := func(card Card) bool {
		data := CARD_DATA_BASE[card.Name]
		return data.IsNormalMonster
	}
	if omws.CountFunc(state.P1.Graveyard, f) >= 2 {
		actions := NewGraveyardIndicesActions([]int{2}, "闇の量産工場", len(state.P1.Graveyard), 0, false)
		g := func(action Action) bool {
			card1 := state.P1.Graveyard[action.GraveyardIndices[0]]
			card2 := state.P2.Graveyard[action.GraveyardIndices[1]]
			data1 := CARD_DATA_BASE[card1.Name]
			data2 := CARD_DATA_BASE[card2.Name]
			return data1.IsNormalMonster && data2.IsNormalMonster
		}
		effect0 = fn.Filter(actions, g)
	}
	return Actionss{effect0}
}

func (_ *costActionss) LegacyOfYataGarasu(state *State) Actionss {
	cardName := CardName("八汰烏の骸")
	actions := Actions{
		Action{CardName: cardName, SelectEffectNumber: 0},
		Action{CardName: cardName, SelectEffectNumber: 1},
	}
	f := func(action Action) bool {
		if action.SelectEffectNumber == 0 {
			return len(state.P1.Deck) >= 1
		} else {
			return slices.ContainsFunc(state.P2.MonsterZone, IsSpiritMonsterCard) && len(state.P1.Deck) >= 2
		}
	}
	effect0 := fn.Filter(actions, f)
	return Actionss{effect0}
}
