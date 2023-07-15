package bmgt

import (
	"github.com/sw965/omw/fn"
)

type Cost func(*Action) []StateTransition

// 王立魔法図書館
func RoyalMagicalLibraryCost(action *Action) []StateTransition {
	effect0 := fn.IdentityWithNilError[State]
	effect1 := func(state State) (State, error) {
		monsterZone := state.P1.MonsterZone.Clone()
		monsterZone[action.MonsterZoneIndices[0]].SpellCounter = 0
		state.P1.MonsterZone = monsterZone
		return state, nil
	}
	return []StateTransition{effect0, effect1}
}

// サンダードラゴン
func ThunderDragonCost(action *Action) []StateTransition {
	effect0 := func(state State) (State, error) {
		state.P1 = state.P1.Discard(action.HandIndices)
		return state, nil
	}
	return []StateTransition{effect0}
}

// 召喚僧サモンプリースト
func SummonerMonkCost(action *Action) []StateTransition {
	effect0 := fn.IdentityWithNilError[State]
	effect1 := fn.IdentityWithNilError[State]
	effect2 := func(state State) (State, error) {
		state.P1 = state.P1.Discard(action.HandIndices)
		monsterZone := state.P1.MonsterZone.Clone()
		monsterZone[action.MonsterZoneIndices[0]].ThisTurnEffectActivationCounts[2] += 1
		state.P1.MonsterZone = monsterZone
		return state, nil
	}
	return []StateTransition{effect0, effect1, effect2}
}

// トゥーン・ワールド
func ToonWorldCost(action *Action) []StateTransition {
	effect0 := func(state State) (State, error) {
		state.P1.LifePoint -= 1000
		return state, nil
	}
	return []StateTransition{effect0}
}

// 魔法石の採掘
func MagicalStoneExcavationCost(action *Action) []StateTransition {
	effect0 := func(state State) (State, error) {
		state.P1 = state.P1.Discard(action.HandIndices)
		return state, nil
	}
	return []StateTransition{effect0}
}

// //八汰烏の骸
func LegacyOfYataGarasuCost(action *Action) []StateTransition {
	effect0 := func(state State) (State, error) {
		spellTrapZone := state.P1.SpellTrapZone.Clone()
		spellTrapZone[action.SpellTrapZoneIndices[0]].SelectEffectNumber = action.SelectEffectNumber
		state.P1.SpellTrapZone = spellTrapZone
		return state, nil
	}
	return []StateTransition{effect0}
}

var COST = map[CardName]Cost{
	"王立魔法図書館":     RoyalMagicalLibraryCost,
	"サンダー・ドラゴン":   ThunderDragonCost,
	"召喚僧サモンプリースト": SummonerMonkCost,
	"トゥーン・ワールド":   ToonWorldCost,
	"魔法石の採掘":      MagicalStoneExcavationCost,
	"八汰烏の骸":       LegacyOfYataGarasuCost,
}
