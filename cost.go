package bmgt

import (
	"github.com/sw965/omw/fn"
)

type Cost func(Player, CardID) []StateTransition

//サンダードラゴン
func ThunderDragonCost(player Player, id CardID) []StateTransition {
	effect0 := func(state State) (State, error) {
		idx := state.P1.Hand.IDIndex(id)
		state.P1 = state.P1.Discard([]int{idx})
		return state, nil
	}
	return []StateTransition{effect0}
}

//召喚僧サモンプリースト
func SummonerMonkCost(player Player, id CardID) []StateTransition {
	effect0 := fn.IdentityWithNilError[State]
	effect1 := fn.IdentityWithNilError[State]
	effect2 := func(state State) (State, error) {
		action := player(&state)
		state.P1 = state.P1.Discard(action.DiscardIndices)
		monsterZone := state.P1.MonsterZone.Clone()
		idx := monsterZone.IDIndex(id)
		monsterZone[idx].ThisTurnEffectActivateCounts[2] += 1
		state.P1.MonsterZone = monsterZone
		return state, nil
	}
	return []StateTransition{effect0, effect1, effect2}
}

//トゥーン・ワールド
func ToonWorldCost(player Player, id CardID) []StateTransition {
	effect0 := func(state State) (State, error) {
		state.P1.LifePoint -= 1000
		return state, nil
	}
	return []StateTransition{effect0}
}

//魔法石の採掘
func MagicalStoneExcavationCost(player Player, id CardID) []StateTransition {
	effect0 := func(state State) (State, error) {
		action := player(&state)
		state.P1 = state.P1.Discard(action.DiscardIndices)
		return state, nil
	}
	return []StateTransition{effect0}
}

//八汰烏の骸
func LegacyOfYataGarasuCost(player Player, id CardID) []StateTransition {
	effect0 := func(state State) (State, error) {
		action := player(&state)
		spellTrapZone := state.P1.SpellTrapZone.Clone()
		idx := spellTrapZone.IDIndex(id)
		spellTrapZone[idx].SelectEffectNumber = action.SelectEffectNumber
		state.P1.SpellTrapZone = spellTrapZone
		return state, nil
	}
	return []StateTransition{effect0}
}

var COSTS = map[CardName]Cost{
	"サンダー・ドラゴン":   ThunderDragonCost,
	"召喚僧サモンプリースト": SummonerMonkCost,
	"トゥーン・ワールド":   ToonWorldCost,
	"魔法石の採掘":      MagicalStoneExcavationCost,
	"八汰烏の骸":       LegacyOfYataGarasuCost,
}
