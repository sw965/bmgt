package bmgt

import (
	"golang.org/x/exp/slices"
)

type SelectOK func(*State, Player) [][]bool

//八汰烏の骸
func LegacyOfYataGarasuSelectOK(state *State, player Player) [][]bool {
	effect00 := len(state.P1.Deck) >= 1
	effect01 := len(state.P1.Deck) >= 2 && slices.ContainsFunc(state.P2.MonsterZone, IsSpiritMonsterCard)
	return [][]bool{[]bool{effect00, effect01}}
}
