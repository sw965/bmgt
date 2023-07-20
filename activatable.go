package bmgt

import (
	"golang.org/x/exp/slices"
	omws "github.com/sw965/omw/slices"
)

type handActivatable struct{}
var HandActivatable = &handActivatable{}

//一時休戦
func (_ *handActivatable) OneDayOfPeace(state *State) bool {
	return len(state.P1.Deck) >= 1 && len(state.P2.Deck) >= 1
}

//打ち出の小槌
func (_ *handActivatable) MagicalMallet(state *State) bool {
	return len(state.P1.Hand) >= 2
}

//強欲な壺
func (_ *handActivatable) PotOfGreed(state *State) bool {
	return len(state.P1.Deck) >= 2
}

//精神統一
func (_ *handActivatable) GatherYourMind(state *State) bool {
	return slices.ContainsFunc(state.P1.Deck, EqualNameCard("精神統一"))  && slices.Contains(state.P1.OncePerTurnLimitCardNames, "精神統一")
}

//手札断殺
func (_ *handActivatable) HandDestruction(state *State) bool {
	return len(state.P1.Deck) >= 2 && len(state.P2.Deck) >= 2 && len(state.P1.Hand) >= 3 && len(state.P2.Hand) >= 2
}

//トゥーンのもくじ
func (_ *handActivatable) ToonTableOfContents(state *State) bool {
	return slices.ContainsFunc(state.P1.Deck, IsToonCard)
}

//トゥーン・ワールド
func (_ *handActivatable) ToonWorld(state *State) bool {
	return state.P1.LifePoint >= 1000
}

//成金ゴブリン
func (_ *handActivatable) UpstartGoblin(state *State) bool {
	return len(state.P1.Deck) >= 1
}

//魔法石の採掘
func (_ *handActivatable) MagicalStoneExcavation(state *State) bool {
	return len(state.P1.Hand) >= 3 && slices.ContainsFunc(state.P1.Graveyard, IsSpellCard)
}

//闇の誘惑
func (_ *handActivatable) AllureOfDarkness(state *State) bool {
	return len(state.P1.Deck) >= 2
}

//闇の量産工場
func (_ *handActivatable) DarkFactoryOfMassProduction(state *State) bool {
	return omws.CountFunc(state.P1.Graveyard, IsNormalMonsterCard) >= 2
}

type CardActivatableData map[CardName]func(*State) bool

var HAND_CARD_ACTIVATABLE = CardActivatableData {
	"一時休戦":HandActivatable.OneDayOfPeace,
	"打ち出の小槌":HandActivatable.MagicalMallet,
	"強欲な壺":HandActivatable.PotOfGreed,
	"精神統一":HandActivatable.GatherYourMind,
	"手札断殺":HandActivatable.HandDestruction,
	"トゥーンのもくじ":HandActivatable.ToonTableOfContents,
	"トゥーン・ワールド":HandActivatable.ToonWorld,
	"成金ゴブリン":HandActivatable.UpstartGoblin,
	"魔法石の採掘":HandActivatable.MagicalStoneExcavation,
	"闇の誘惑":HandActivatable.AllureOfDarkness,
	"闇の量産工場":HandActivatable.DarkFactoryOfMassProduction,
}

var HAND_QUICK_PLAY_SPELL_ACTIVATABLE = func() CardActivatableData {
	result := CardActivatableData{}
	for name, f := range HAND_CARD_ACTIVATABLE {
		data := CARD_DATA_BASE[name]
		if data.IsQuickPlaySpell {
			result[name] = f
		}
	}
	return result
}()