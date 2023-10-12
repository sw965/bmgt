package bmgt

import (
	"math/rand"
)

func SummonerMonkEffect2(state *State, idx int, r *rand.Rand) {
	state.P1.Search([]int{idx}, true, r)
}

func ThunderDragonEffect0(state *State, idxs []int, r *rand.Rand) {
	state.P1.Search(idxs, true, r)
}

func RoyalMagicalLibraryEffect1(state *State) {
	state.P1.Draw(1)
}

func PotOfGreedEffect0(state *State) {
	state.P1.Draw(2)
}

func GatherYourMindEffect0(state *State, idx int, r *rand.Rand) {
	state.P1.Search([]int{idx}, true, r)
}

func HandDestructionEffect0(state *State, idxs []int) {
	state.P1.Discard(idxs)
	state.P1.Draw(2)
}

func DoubleSummonEffect0(state *State) {
	state.P1.IsDoubleSummonApplied = true
}

func ToonTableOfContentsEffect0(state *State, idx int, r *rand.Rand) {
	state.P1.Search([]int{idx}, true, r)
}

func UpStartGoblinEffect0(state *State) {
	state.P1.Draw(1)
	state.P2.LifePoint += 1000
}

func MagicalStoneExcavationEffect0(state *State, ids CardIDs) {
	state.P1.IDSalvage(ids)
}

func AllureOfDarknessEffect0(state *State, idx int) bool {
	if state.P1.EffectProcessingNumber == 1 {
		state.P1.DiscardBanish([]int{idx})
		return true
	} else {
		state.P1.Draw(2)
		return false
	}
}

func DarkFactoryOfMassProductionEffect0(state *State, ids CardIDs) {
	state.P1.IDSalvage(ids)
}

func JarOfGreedEffect0(state *State) {
	state.P1.Draw(1)
}

func LegacyOfYataGarasuEffect0(state *State, action *Action) {
	var draw int
	if action.EffectSelectNumber == 0 {
		draw = 1
	} else {
		draw = 2
	}
	state.P1.Draw(draw)
}
