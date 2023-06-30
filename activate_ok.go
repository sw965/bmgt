package bmgt

import (
	"golang.org/x/exp/slices"
)

//https://yugioh-wiki.net/index.php?%C8%AF%C6%B0

type ActivateOK func(*State) bool

//召喚僧サモンプリースト
func SummonerMonkActivateOK(state *State) bool {
	effect2 := slices.ContainsFunc(state.P1.Hand, IsSpellCard) && slices.ContainsFunc(state.P1.Deck, IsLevel4MonsterCard)
	return effect2
}

//強欲な壺
func PotOfGreedActivateOK(state *State) bool {
	return len(state.P1.Deck) >= 2
}

//強欲な瓶
func JarOfGreedActivateOK(state *State) bool {
	return len(state.P1.Deck) >= 1
}

//八汰烏の骸
func LegacyOfYataGarasuActivateOK(state *State) bool {
	select0 := len(state.P1.Deck) >= 1
	select1 := len(state.P1.Deck) >= 2 && slices.ContainsFunc(state.P1.MonsterZone, IsSpiritMonsterCard)
	return select0 || select1
}

//トゥーンのもくじ
func ToonTableOfContentsActivateOK(state *State) bool {
	return slices.ContainsFunc(state.P1.Deck, IsToonCard)
}

//手札断殺
func HandDestructionActivateOK(state *State) bool {
	return len(state.P1.Hand) >= 3 && len(state.P2.Hand) >= 2 && len(state.P1.Deck) >= 2 && len(state.P2.Deck) >= 2
}

var ACTIVATE_OK = map[CardName]ActivateOK {
	"強欲な壺":PotOfGreedActivateOK,
	"強欲な瓶":JarOfGreedActivateOK,
	"八汰烏の骸":LegacyOfYataGarasuActivateOK,
	TOON+"のもくじ":ToonTableOfContentsActivateOK,
}