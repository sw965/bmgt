package bmgt

import (
	omws "github.com/sw965/omw/slices"
)

type ActivateOK func(*State, *Card) []bool

//サンダー・ドラゴン
func ThunderDragonActivateOK(state *State, card *Card) []bool {
	hc := omws.CountFunc(state.P1.Hand, EqualIDCard(card.ID))
	dc := omws.CountFunc(state.P1.Deck, EqualNameCard("サンダー・ドラゴン"))
	effect0 := hc >= 1 && dc >= 1
	return []bool{effect0}
}

//召喚僧サモンプリースト
func SummonerMonkActivateOK(state *State, card *Card) []bool {
	ac := card.ThisTurnEffectActivateCounts[2]
	hc := omws.CountFunc(state.P1.Hand, IsSpellCard)
	dc := omws.CountFunc(state.P1.Deck, IsLevel4MonsterCard)
	effect2 := ac == 1 && hc >= 1 && dc >= 1
	return []bool{false, false, effect2}
}

//一時休戦
func OneDayOfPeaceActivateOK(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 1 && len(state.P2.Deck) >= 1
	return []bool{effect0}
}

//打ち出の小槌
func MagicalMalletActivateOK(state *State, card *Card) []bool {
	effect0 := len(state.P1.Hand) >= 1
	return []bool{effect0}
}

//強欲な壺
func PotOfGreedActivateOK(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 2
	return []bool{effect0}
}

//手札断殺
func HandDestructionActivateOK(state *State, card *Card) []bool {
	effect0 := len(state.P1.Hand) >= 3 && len(state.P2.Hand) >= 2 && len(state.P1.Deck) >= 2 && len(state.P2.Deck) >= 2
	return []bool{effect0}
}

//トゥーン・ワールド
func ToonWorldActivateOK(state *State, card *Card) []bool {
	effect0 := state.P1.LifePoint > 1000
	return []bool{effect0}
}

//強欲な瓶
func JarofGreedActivateOK(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 1
	return []bool{effect0}
}

//八汰烏の骸
func LegacyOfYataGarasuActivateOK(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 1
	return []bool{effect0}
}

var ACTIVATE_OK = map[CardName]ActivateOK{
	"サンダー・ドラゴン":   ThunderDragonActivateOK,
	"召喚僧サモンプリースト": SummonerMonkActivateOK,
	"一時休戦":        OneDayOfPeaceActivateOK,
	"打ち出の小槌":      MagicalMalletActivateOK,
	"強欲な壺":        PotOfGreedActivateOK,
	"手札断殺":        HandDestructionActivateOK,
	"トゥーン・ワールド":   ThunderDragonActivateOK,
	"八汰烏の骸":       LegacyOfYataGarasuActivateOK,
}
