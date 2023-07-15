package bmgt

import (
	"github.com/sw965/omw/fn"
	omws "github.com/sw965/omw/slices"
)

type effectActionss struct{}

var EffectActionss = effectActionss{}

// サンダー・ドラゴン
func (_ *effectActionss) ThunderDragon(state *State) Actionss {
	cardName := CardName("サンダー・ドラゴン")
	actions := NewDeckIndicesActions([]int{1, 2}, cardName, len(state.P1.Deck), 0, false)
	f := func(action Action) bool {
		for _, idx := range action.DeckIndices {
			if state.P1.Deck[idx].Name != cardName {
				return false
			}
		}
		return true
	}
	effect0 := fn.Filter(actions, f)
	return Actionss{effect0}
}

// 召喚僧サモンプリースト
func (_ *effectActionss) SummonerMonk(state *State) Actionss {
	effect0 := Actions{}
	effect1 := Actions{}

	cardName := CardName("召喚僧サモンプリースト")
	actions := NewDeckIndexActions(cardName, len(state.P1.Deck), 2, false)
	f := func(action Action) bool {
		return IsLevel4MonsterCard(state.P1.Deck[action.DeckIndices[0]])
	}
	effect2 := fn.Filter(actions, f)
	return Actionss{effect0, effect1, effect2}
}

// 打ち出の小槌
func (_ *effectActionss) MagicalMallet(state *State) Actionss {
	n := len(state.P1.Hand)
	rng := omws.IntegerRange[[]int, int]{Start: 0, End: n, Step: 1}
	selectCardNums := rng.Make()
	effect0 := NewHandIndicesActions(selectCardNums, "打ち出の小槌", n, 0, false)
	return Actionss{effect0}
}

// 手札断殺
func (_ *effectActionss) HandDestruction(state *State) Actionss {
	effect0 := NewHandIndicesActions([]int{2}, "手札断殺", len(state.P1.Hand), 0, false)
	return Actionss{effect0}
}

// トゥーンのもくじ
func (_ *effectActionss) ToonTableOfContents(state *State) Actionss {
	actions := NewDeckIndexActions("トゥーンのもくじ", len(state.P1.Deck), 0, false)
	f := func(action Action) bool {
		return IsToonCard(state.P1.Deck[action.DeckIndices[0]])
	}
	effect0 := fn.Filter(actions, f)
	return Actionss{effect0}
}

// 魔法石の採掘
func (_ *effectActionss) MagicalStoneExcavation(state *State) Actionss {
	actions := NewGraveyardIndexActions("魔法石の採掘", len(state.P1.Deck), 0, false)
	f := func(action Action) bool {
		return IsSpellCard(state.P1.Deck[action.GraveyardIndices[0]])
	}
	effect0 := fn.Filter(actions, f)
	return Actionss{effect0}
}
