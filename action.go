package bmgt

import (
	"fmt"
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
	omwslices "github.com/sw965/omw/slices"
	omwmath "github.com/sw965/omw/math"
)

type Action struct {
	CardNames CardNames
	CardIDs CardIDs

	HandIndices []int
	MonsterZoneIndices []int
	SpellTrapZoneIndices []int

	IsNormalSummon bool
	IsTributeSummonCost bool
	IsTributeSummon bool
	IsFlipSummon bool
	IsCardActivation bool
	IsSpellTrapSet bool

	IsWaiverOfPriority bool

	BattlePosition BattlePosition
}

type Actions []Action

type actionF struct{}
var ActionF = actionF{}

func (f *actionF) ToString(action Action) string {
	msg := ""
	msg += fmt.Sprintf("Name:%v, ", CardNamesToStrings(action.CardNames))
	msg += fmt.Sprintf("CardID:%v, ", action.CardIDs)
	msg += fmt.Sprintf("HandIndices:%v, ", action.HandIndices)
	msg += fmt.Sprintf("MonsterZoneIndices:%v, ", action.MonsterZoneIndices)
	msg += fmt.Sprintf("SpellTrapZoneIndices:%v, ", action.SpellTrapZoneIndices)
	msg += fmt.Sprintf("IsNormalSummon:%v, ", action.IsNormalSummon)
	msg += fmt.Sprintf("IsFlipSummon:%v, ", action.IsFlipSummon)
	msg += fmt.Sprintf("BattlePosition:%v", action.BattlePosition.ToString())
	return msg
}

func (f *actionF) MonsterZoneIndicesAccess(state *State, action Action) Cards {
	return omwslices.IndicesAccess(state.P1.MonsterZone)(action.MonsterZoneIndices)
}

func (f *actionF) IsLegalNormalSummon(state *State) func(Action) bool {
	return func(action Action) bool {
		hIdx := action.HandIndices[0]
		mIdx := action.MonsterZoneIndices[0]

		ok1 := CardF.CanNormalSummon(state.P1.Hand[hIdx])
		ok2 := CardF.IsEmpty(state.P1.MonsterZone[mIdx])
		ok3 := action.IsNormalSummon
		return ok1 && ok2 && ok3
	}
}

func (f *actionF) IsLegalTributeSummonCost(state *State, cost int) func(Action) bool {
	return func(action Action) bool {
		costs := CardsF.TributeSUmmonCosts(state.P1.Hand)
		mZone := ActionF.MonsterZoneIndicesAccess(state, action)

		ok1 := slices.Contains(costs, cost)
		ok2 := fn.All(mZone, CardF.CanTributeSummonCost)
		ok3 := action.IsTributeSummonCost
		return ok1 && ok2 && ok3
	}
}

func (f *actionF) IsLegalTributeSummon(state *State) func(Action) bool {
	return func(action Action) bool {
		hIdx := action.HandIndices[0]
		mIdx := action.MonsterZoneIndices[0]

		ok1 := CardF.TributeSummonCost(state.P1.Hand[hIdx]) == len(state.P1.TributeSummonCostIndices)
		ok2 := CardF.IsEmpty(state.P1.MonsterZone[mIdx])
		ok3 := action.IsTributeSummon
		return ok1 && ok2 && ok3
	}
}

func (f *actionF) IsLegalFlipSummon(state *State) func(Action) bool {
	return func(action Action) bool {
		idx := action.MonsterZoneIndices[0]

		ok1 := CardF.IsNotEmpty(state.P1.MonsterZone[idx])
		ok2 := CardF.CanFlipSummon(state.P1.MonsterZone[idx])
		ok3 := action.IsFlipSummon
		return ok1 && ok2 && ok3
	}
}

func (f *actionF) IsLegalHandSpellCardActivation(state *State, category SpellCategory) func(Action) bool {
	return func(action Action) bool {
		hIdx := action.HandIndices[0]
		spIdx := action.SpellTrapZoneIndices[0]

		ok1 := state.P1.Hand[hIdx].Category == Category(category)
		ok2 := CardF.IsEmpty(state.P1.SpellTrapZone[spIdx])
		ok3 := action.IsCardActivation
		return ok1 && ok2 && ok3
	}	
}

func (f *actionF) IsLegalHandNormalSpellCardActivation(state *State) func(Action) bool {
	return f.IsLegalHandSpellCardActivation(state, NORMAL_SPELL)
}

func (f *actionF) IsLegalHandQuickPlaySpellCardActivation(state *State) func(Action) bool {
	return f.IsLegalHandSpellCardActivation(state, QUICK_PLAY_SPELL)
}

func (f *actionF) IsLegalHandContinuousSpellCardActivation(state *State) func(Action) bool {
	return f.IsLegalHandSpellCardActivation(state, CONTINUOUS_SPELL)
}

func (f *actionF) IsLegalHandTrapSet(state *State) func(Action) bool {
	return func(action Action) bool {
		hIdx := action.HandIndices[0]
		spIdx := action.SpellTrapZoneIndices[0]
		ok1 := state.P1.Hand[hIdx].Category.IsTrap()
		ok2 := CardF.IsEmpty(state.P1.SpellTrapZone[spIdx])
		ok3 := action.IsSpellTrapSet
		return ok1 && ok2 && ok3
	}
}

type actionsF struct{}
var ActionsF = actionsF{}

func (f *actionsF) ToStrings(actions Actions) []string {
	return fn.Map[[]string](actions, ActionF.ToString)
}

func (f *actionsF) NewNormalSummon(state *State) Actions {
	y := make(Actions, 0, len(state.P1.Hand) * MONSTER_ZONE_LENGTH * len(NORMAL_SUMMON_BATTLE_POSITIONS))
	for handI := range state.P1.Hand {
		for zoneI := range state.P1.MonsterZone {
			for _, pos := range NORMAL_SUMMON_BATTLE_POSITIONS {
				card := state.P1.Hand[handI]
				action := Action{
					CardNames:CardNames{card.Name},
					CardIDs:CardIDs{card.ID},
					HandIndices:[]int{handI},
					MonsterZoneIndices:[]int{zoneI},
					IsNormalSummon:true,
					BattlePosition:pos,
				}
				y = append(y, action)
			}
		}
	}
	return y
}

func (f *actionsF) NewLegalNormalSummon(state *State) Actions {
	actions := f.NewNormalSummon(state)
	return fn.Filter(actions, ActionF.IsLegalNormalSummon(state))
}

func (f *actionsF) NewTributeSummonCost(state *State, cost int) Actions {
	c := omwmath.Combination{N:MONSTER_ZONE_LENGTH, R:cost}
	y := make(Actions, c.TotalNum())
	for i, idxs := range c.Get() {
		cards := omwslices.IndicesAccess(state.P1.MonsterZone)(idxs)
		y[i] = Action{
			CardNames:CardsF.Names(cards),
			CardIDs:CardsF.IDs(cards),
			MonsterZoneIndices:idxs,
			IsTributeSummonCost:true,
		}
	}
	return y
}

func (f *actionsF) NewLegalTributeSummonCost(state *State, cost int) Actions {
	actions := f.NewTributeSummonCost(state, cost)
	return fn.Filter(actions, ActionF.IsLegalTributeSummonCost(state, cost))
}

func (f *actionsF) NewTributeSummon(state *State) Actions {
	y := make(Actions, 0, MONSTER_ZONE_LENGTH * len(state.P1.Hand))
	for handI, hCard := range state.P1.Hand {
		for zoneI := 0; zoneI < MONSTER_ZONE_LENGTH; zoneI++ {
			action := Action{
				CardNames:CardNames{hCard.Name},
				CardIDs:CardIDs{hCard.ID},
				HandIndices:[]int{handI},
				MonsterZoneIndices:[]int{zoneI},
				IsTributeSummon:true,
			}
			y = append(y, action)
		}
	}
	return y
}

func (f *actionsF) NewLegalTributeSummon(state *State) Actions {
	actions := f.NewTributeSummon(state)
	return fn.Filter(actions, ActionF.IsLegalTributeSummon(state))
}

func (f *actionsF) NewFlipSummon(state *State) Actions {
	y := make(Actions, 0, MONSTER_ZONE_LENGTH)
	for i, card := range state.P1.MonsterZone {
		action := Action{
			CardNames:CardNames{card.Name},
			CardIDs:CardIDs{card.ID},
			MonsterZoneIndices:[]int{i},
			BattlePosition:ATTACK_POSITION,
			IsFlipSummon:true,
		}
		y = append(y, action)
	}
	return y
}

func (f *actionsF) NewLegalFlipSummon(state *State) Actions {
	actions := f.NewFlipSummon(state)
	return fn.Filter(actions, ActionF.IsLegalFlipSummon(state))
}

func (f *actionsF) NewHandSpellCardActivation(state *State) Actions {
	y := make(Actions, 0, len(state.P1.Hand) * SPELL_TRAP_ZONE_LENGTH)
	for handI, hCard  := range state.P1.Hand {
		for zoneI := 0; zoneI < SPELL_TRAP_ZONE_LENGTH; zoneI++ {
			action := Action{
				CardNames:CardNames{hCard.Name},
				CardIDs:CardIDs{hCard.ID},
				HandIndices:[]int{handI},
				SpellTrapZoneIndices:[]int{zoneI},
				IsCardActivation:true,
			}
			y = append(y, action)
		}
	}
	return y
}

func (f *actionsF) NewLegalHandNormalSpell(state *State) Actions {
	actions := f.NewHandSpellCardActivation(state)
	return fn.Filter(actions, ActionF.IsLegalHandNormalSpellCardActivation(state))
}

func (f *actionsF) NewLegalHandQuickPlaySpellCardActivation(state *State) Actions {
	actions := f.NewHandSpellCardActivation(state)
	return fn.Filter(actions, ActionF.IsLegalHandQuickPlaySpellCardActivation(state))
}

func (f *actionsF) NewLegalHandContinuousSpellCardActivation(state *State) Actions {
	actions := f.NewHandSpellCardActivation(state)
	return fn.Filter(actions, ActionF.IsLegalHandContinuousSpellCardActivation(state))
}

func (f *actionsF) NewHandTrapSet(state *State) Actions {
	y := make(Actions, 0, len(state.P1.Hand) * SPELL_TRAP_ZONE_LENGTH)
	for handI, hCard := range state.P1.Hand {
		for zoneI := 0; zoneI < SPELL_TRAP_ZONE_LENGTH; zoneI++ {
			action := Action{
				CardNames:CardNames{hCard.Name},
				CardIDs:CardIDs{hCard.ID},
				HandIndices:[]int{handI},
				SpellTrapZoneIndices:[]int{zoneI},
				IsSpellTrapSet:true,
			}
			y = append(y, action)
		}
	}
	return y
}

func (f *actionsF) NewLegalHandTrapSet(state *State) Actions {
	actions := f.NewHandTrapSet(state)
	return fn.Filter(actions, ActionF.IsLegalHandTrapSet(state))
}