package bmgt

import (
	"github.com/sw965/omw/fn"
	omwslices "github.com/sw965/omw/slices"
)

type Action struct {
	HandIndices []int
	MonsterZoneIndices []int
	SpellTrapZoneIndices []int

	IsNormalSummon bool
	BattlePosition BattlePosition
}

type Actions []Action

type actionF struct{}
var ActionF = actionF{}

func (f *actionF) NewNormalSummon(state *State, handI, zoneI int, pos BattlePosition) Action {
	return Action{
		HandIndices:[]int{handI},
		MonsterZoneIndices:[]int{zoneI},
		IsNormalSummon:true,
	}
}

func (f *actionF) GetHand(state *State, action *Action) Cards {
	return fn.Map[Cards](action.HandIndices, omwslices.IndexAccess[Cards](state.P1.Hand))
}

func (f *actionF) GetMonsterZone(state *State, action *Action) Cards {
	return fn.Map[Cards](action.MonsterZoneIndices, omwslices.IndexAccess[Cards](state.P1.MonsterZone))
}

func (f *actionF) IsHandNormalSummonPossible(state *State, action *Action) bool {
	cards := f.GetHand(state, action)
	return fn.All(cards, CardF.CanNormalSummon)
}

func (f *actionF) IsEmptyMonsterZoneIndices(state *State, action *Action) bool {
	return fn.All(f.GetMonsterZone(state, action), CardF.IsEmpty)
}

func (f *actionF) IsLegalNormalSummon(state *State) func(Action) bool {
	return func(action Action) bool {
		ok1 := f.IsHandNormalSummonPossible(state, &action)
		ok2 := f.IsEmptyMonsterZoneIndices(state, &action)
		return ok1 && ok2
	}
}

type actionsF struct{}
var ActionsF = actionsF{}

func (f *actionsF) NewNormalSummon(state *State) Actions {
	y := make(Actions, 0, len(state.P1.Hand) * MONSTER_ZONE_LENGTH * len(NORMAL_SUMMON_BATTLE_POSITIONS))
	for handI := range state.P1.Hand {
		for zoneI := range state.P1.MonsterZone {
			for _, pos := range NORMAL_SUMMON_BATTLE_POSITIONS {
				y = append(y, ActionF.NewNormalSummon(state, handI, zoneI, pos))
			}
		}
	}
	return y
}

func (f *actionsF) NewLegalNormalSummonActions(state *State) Actions {
	actions := f.NewNormalSummon(state)
	return fn.Filter(actions, ActionF.IsLegalNormalSummon(state))
}