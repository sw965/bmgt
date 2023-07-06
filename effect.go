package bmgt

import (
	"github.com/sw965/omw/fn"
	"math/rand"
)

type Effect func(Player, CardID, *rand.Rand) []StateTransition

//サンダー・ドラゴン
func ThunderDragonEffect(player Player, id CardID, r *rand.Rand) []StateTransition {
	effect0 := func(state State) (State, error) {
		action := player(&state)
		state.P1 = state.P1.Search(action.SearchIndices, r)
		return state, nil
	}
	return []StateTransition{effect0}
}

//召喚僧サモンプリースト
func SummonerMonkEffect(player Player, id CardID, r *rand.Rand) []StateTransition {
	effect0 := func(state State) (State, error) {
		monsterZone := state.P1.MonsterZone.Clone()
		idx := monsterZone.IDIndex(id)
		monsterZone[idx].BattlePosition = FACE_UP_DEFENSE_POSITION
		state.P1.MonsterZone = monsterZone
		return state, nil
	}

	effect1 := fn.IdentityWithNilError[State]

	effect2 := func(state State) (State, error) {
		action := player(&state)
		state.P1 = state.P1.Search(action.SearchIndices, r)
		return state, nil
	}
	return []StateTransition{effect0, effect1, effect2}
}

//一時休戦
func OneDayOfPeaceEffect(player Player, id CardID, r *rand.Rand) []StateTransition {
	effect0 := func(state State) (State, error) {
		var err error
		state.P1, err = state.P1.Draw(1)
		if err != nil {
			return state, err
		}
		state.P2, err = state.P2.Draw(1)
		state.OneDayOfPeace = true
		return state, err
	}
	return []StateTransition{effect0}
}

//打ち出の小槌
func MagicalMalletEffect(player Player, id CardID, r *rand.Rand) []StateTransition {
	effect0 := func(state State) (State, error) {
		var cards Cards
		action := player(&state)
		idxs := action.HandIndices
		state.P1.Hand, cards = omws.Pop(state.P1.Hand, action.HandIndices)
		state.P1.Deck = append(state.P1.Deck, cards...)
		state.P1.Deck = omwrand.Shuffled(state.P1.Deck, r)
		state.P1, err = state.P1.Draw(len(idxs))
		return state, err
	}

	return []StateTransition{effect0}
}

var EFFECT = map[CardName]Effect{
	"サンダー・ドラゴン":   ThunderDragonEffect,
	"召喚僧サモンプリースト": SummonerMonkEffect,
	"一時休戦":        OneDayOfPeaceEffect,
	"打ち出の小槌":      MagicalMalletEffect,
}
