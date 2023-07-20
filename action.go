package bmgt

import (
	omws "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
)

type Action struct {
	CardName CardName
	HandIndices []int
	MonsterZoneIndices1 []int
	MonsterZoneIndices2 []int
	SpellTrapZoneIndices []int

	IsFaceUpNormalSummon bool
	IsFaceDownNormalSummon bool
	IsFaceUpTributeSummon bool
	IsFaceDownTributeSummon bool
	IsActivationOfCard bool
	IsSetSpellTrap bool

	IsCost bool
	EffectNumber int
}

type Actions []Action

func NewNormalSummonActions(name CardName, state *State, handIdx int) Actions {
	monsterZoneEmptyIndices := omws.IndicesFunc(state.P1.MonsterZone, IsEmptyCard)
	result := make(Actions, 0, len(monsterZoneEmptyIndices))
	for _, zoneI := range monsterZoneEmptyIndices {
		action := Action{
			CardName:name,
			HandIndices:[]int{handIdx},
			MonsterZoneIndices1:[]int{zoneI},
			IsFaceUpNormalSummon:true,
		}
		result = append(result, action)
	}

	f := func(action Action) Action {
		action.IsFaceUpNormalSummon = false
		action.IsFaceDownNormalSummon = true
		return action
	}
	result = append(result, fn.Map[Actions](result, f)...)
	return result
}

func NewTributeSummonActions(name CardName, state *State, handIdx, costNum int) Actions {
	monsterZoneEmptyIndices := omws.IndicesFunc(state.P1.MonsterZone, IsEmptyCard)
	monsterZoneEmptyCount := len(monsterZoneEmptyIndices)
	monsterZoneNotEmptyIndices := omws.IndicesFunc(state.P1.MonsterZone, IsNotEmptyCard)
	monsterZoneNotEmptyCount := len(monsterZoneNotEmptyIndices)
	costIdxss := omws.Combination[[][]int, []int](monsterZoneNotEmptyIndices, costNum)
	result := make(Actions, 0, len(costIdxss) * (costNum + monsterZoneNotEmptyCount) * 2)

	for _, costIdxs := range costIdxss {
		canNotCost := func(costIdx int) bool {
			card := state.P1.MonsterZone[costIdx]
			if card.Name == "召喚僧サモンプリースト" {
				if card.BattlePosition.IsFaceUp() {
					return true
				}
			}
			return false
		}

		if fn.All(costIdxs, canNotCost) {
			continue
		}

		emptyIdxs := make([]int, 0, len(costIdxs) + monsterZoneEmptyCount)
		emptyIdxs = append(emptyIdxs, costIdxs...)
		emptyIdxs = append(emptyIdxs, monsterZoneEmptyIndices...)
		for _, emptyI := range emptyIdxs {
			action := Action{
				CardName:name,
				HandIndices:[]int{handIdx},
				MonsterZoneIndices1:costIdxs,
				MonsterZoneIndices2:[]int{emptyI},
				IsFaceUpTributeSummon:true,
			}
			result = append(result, action)
		}
	}

	f := func(action Action) Action {
		action.IsFaceUpTributeSummon = false
		action.IsFaceDownTributeSummon = true
		return action
	}
	result = append(result, fn.Map[Actions](result, f)...)
	return result
}

func NewHandCardActivationActions(name CardName, state *State, handIdx int) Actions {
	spellTrapEmptyIndices := omws.IndicesFunc(state.P1.SpellTrapZone, IsEmptyCard)
	actions := make(Actions, 0, len(spellTrapEmptyIndices))
	for _, zoneI := range spellTrapEmptyIndices {
		action := Action{
			CardName:name,
			HandIndices:[]int{handIdx},
			SpellTrapZoneIndices:[]int{zoneI},
			IsActivationOfCard:true,
		}
		actions = append(actions, action)
	}
	return actions
}

func NewHandSpellTrapSetActions(name CardName, state *State, handIdx int) Actions {
	spellTrapEmptyIndices := omws.IndicesFunc(state.P1.SpellTrapZone, IsEmptyCard)
	n := len(spellTrapEmptyIndices)
	result := make(Actions, 0, n)
	for _, zoneI := range spellTrapEmptyIndices {
		action := Action{
			CardName:name,
			HandIndices:[]int{handIdx},
			SpellTrapZoneIndices:[]int{zoneI},
			IsSetSpellTrap:true,
		}
		result = append(result, action)
	}
	return result
}