package bmgt

import (
	"github.com/sw965/omw/fn"
	"math/rand"
	omws "github.com/sw965/omw/slices"
	omwrand "github.com/sw965/omw/rand"
	"golang.org/x/exp/slices"
)

type Effect func(*Action, *Card, *rand.Rand) []StateTransition

//王立魔法図書館
func RoyalMagicalLibraryEffect(action *Action, card *Card, r *rand.Rand) []StateTransition {
	effect0 := fn.IdentityWithNilError[State]
	effect1 := func(state State) (State, error) {
		var err error
		state.P1, err = state.P1.Draw(1)
		return state, err
	}
	return []StateTransition{effect0, effect1}
}

//サンダー・ドラゴン
func ThunderDragonEffect(action *Action, card *Card, r *rand.Rand) []StateTransition {
	effect0 := func(state State) (State, error) {
		state.P1 = state.P1.Search(action.DeckIndices, r)
		return state, nil
	}
	return []StateTransition{effect0}
}

//召喚僧サモンプリースト
func SummonerMonkEffect(action *Action, card *Card, r *rand.Rand) []StateTransition {
	effect0 := func(state State) (State, error) {
		idx := slices.IndexFunc(state.P1.MonsterZone, EqualIDCard(card.ID))
		monsterZone := state.P1.MonsterZone.Clone()
		monsterZone[idx].BattlePosition = FACE_UP_DEFENSE_POSITION
		state.P1.MonsterZone = monsterZone
		return state, nil
	}

	effect1 := fn.IdentityWithNilError[State]

	effect2 := func(state State) (State, error) {
		state.P1 = state.P1.Search(action.DeckIndices, r)
		return state, nil
	}
	return []StateTransition{effect0, effect1, effect2}
}

//一時休戦
func OneDayOfPeaceEffect(action *Action, card *Card, r *rand.Rand) []StateTransition {
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
func MagicalMalletEffect(action *Action, card *Card, r *rand.Rand) []StateTransition {
	effect0 := func(state State) (State, error) {
		var cards Cards
		idxs := action.HandIndices
		state.P1.Hand, cards = omws.Pops(state.P1.Hand, action.HandIndices)
		state.P1.Deck = append(state.P1.Deck, cards...)
		state.P1.Deck = omwrand.Shuffled(state.P1.Deck, r)
		var err error
		state.P1, err = state.P1.Draw(len(idxs))
		return state, err
	}
	return []StateTransition{effect0}
}

var EFFECT = map[CardName]Effect{
	"王立魔法図書館":RoyalMagicalLibraryEffect,
	"サンダー・ドラゴン":   ThunderDragonEffect,
	"召喚僧サモンプリースト": SummonerMonkEffect,
	"一時休戦":        OneDayOfPeaceEffect,
	"打ち出の小槌":      MagicalMalletEffect,
}