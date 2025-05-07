package bmgt

import (
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
)

type ActionType int

const (
	PHASE_TRANSITION_ACTION ActionType = iota
	NORMAL_SUMMON_ACTION
	DIRECT_ATTACK_DECLARATION_ACTION
	ATTACK_DECLARATION_ACTION
	PASS_PRIORITY_ACTION
)

type Action struct {
	N1 int
	N2 int
	Type ActionType
}

type Actions []Action

func NewPhaseTransitionActions(duel *Duel) Actions {
	switch duel.Phase {
		case MAIN1_PHASE:
			return Actions{Action{N1:BATTLE_PHASE_INDEX}, Action{N1:END_PHASE_INDEX}}
		case BATTLE_PHASE:
			return Actions{Action{N1:MAIN2_PHASE_INDEX}, Action{N1:END_PHASE_INDEX}}
		case MAIN2_PHASE:
			return Actions{Action{N1:END_PHASE_INDEX}}
		default:
			return Actions{}
	}
}

func newNormalSummonActions(duel *Duel) Actions {
	is := omwslices.NewSequentialInteger[[]int](0, len(duel.P1.Hand))
	js := omwslices.NewSequentialInteger[[]int](0, MONSTER_ZONE_LENGTH)
	f := func(i, j int) Action {
		return Action{N1:i, N2:j, Type:NORMAL_SUMMON_ACTION}
	}
	return omwslices.Product2[[]int, []int, Actions](is, js, f)
}

func IsLegalNormalSummonAction(duel *Duel) func(Action) bool {
	return func(action Action) bool {	
		if !CanNormalSummonCard(duel.P1.Hand[action.N1]) {
			return false
		}

		if duel.P1.MonsterZone[action.N2].Name == NO_NAME {
			return false
		}
		return  true
	}
}

func NewLegalNormalSummonActions(duel *Duel) Actions {
	as := newNormalSummonActions(duel)
	return fn.Filter(as, IsLegalNormalSummonAction(duel))
}

func newDirectAttackDeclarationActions(duel *Duel) Actions {
	is := omwslices.NewSequentialInteger[[]int](0, MONSTER_ZONE_LENGTH)
	f := func(i int) Action { return Action{N1:i, Type:DIRECT_ATTACK_DECLARATION_ACTION} }
	return fn.Map[Actions](is, f)
}

func IsLegalDirectAttackDeclarationActions(duel *Duel) func(Action) bool {	
	return func(action Action) bool {
		i := action.N1
		if duel.P1.MonsterZone[i].Name == NO_NAME {
			return false
		}

		if duel.P1.MonsterZone[i].BattlePosition() != FACE_UP_ATTACK_POSITION {
			return false
		}

		return true
	}
}

func NewLegalDirectAttackDeclarationActions(duel *Duel) Actions {
	if !duel.P2.MonsterZone.IsAllEmpty() {
		return Actions{}
	}
	as := newDirectAttackDeclarationActions(duel)
	return fn.Filter(as, IsLegalDirectAttackDeclarationActions(duel))
}

func NewLegalActions(duel *Duel) Actions {
	return NewPhaseTransitionActions(duel)
}

type Actionss []Actions