package bmgt

import (
	"golang.org/x/exp/slices"
	"github.com/sw965/omw/fn"
)

type ActionType int

const (
	PHASE_TRANSITION_ACTION ActionType = iota
	NORMAL_SUMMON_ACTION
	ATTACK_DECLARE_ACTION
)

type ActionTypes []ActionType

const LENGTH_OF_BOOLS_OF_ACTIONS = MIN_DECK_NUM
type BoolsOfAction [LENGTH_OF_BOOLS_OF_ACTIONS]bool

type Action struct {
	Bools1 BoolsOfAction
	Bools2 BoolsOfAction
	Phase Phase
	BattlePosition BattlePosition
	Type ActionType
}

func GetTypeOfAction(action Action) ActionType {
	return action.Type
}

func (action *Action) Indices1() []int {
	idxs := make([]int, 0, LENGTH_OF_BOOLS_OF_ACTIONS)
	for i, b := range action.Bools1 {
		if b {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func (action *Action) Indices2() []int {
	idxs := make([]int, 0, LENGTH_OF_BOOLS_OF_ACTIONS)
	for i, b := range action.Bools2 {
		if b {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

type Actions []Action

func NewLegalPhaseTransitionActions(state *State) Actions {
	switch state.Phase {
		case DRAW_PHASE:
			standby := Action{Phase:STANDBY_PHASE, Type:PHASE_TRANSITION_ACTION}
			return Actions{standby}
		case STANDBY_PHASE:
			main := Action{Phase:MAIN_PHASE, Type:PHASE_TRANSITION_ACTION}
			return Actions{main}
		case MAIN_PHASE:
			battle := Action{Phase:BATTLE_PHASE, Type:PHASE_TRANSITION_ACTION}
			end := Action{Phase:END_PHASE, Type:PHASE_TRANSITION_ACTION}
			if state.Turn == 1 {
				return Actions{end}
			} else {
				return Actions{battle, end}
			}
		case BATTLE_PHASE:
			end := Action{Phase:END_PHASE, Type:PHASE_TRANSITION_ACTION}
			return Actions{end}
	}
	return Actions{}
}

func NewLegalNormalSummonActions(state *State) Actions {
	if state.Phase != MAIN_PHASE || state.P1.ThisTurnNormalSummonCount != 0 {
		return Actions{}
	}
	poss := BattlePositions{ATK_BATTLE_POSITION, FACE_DOWN_DEF_BATTLE_POSITION}
	y := make(Actions, 0, len(state.P1.Hand) * MONSTER_ZONE_LENGTH)
	for i, hCard := range state.P1.Hand {
		for j, mCard := range state.P1.MonsterZone {
			if slices.Contains(LOW_LEVELS, hCard.Level) && IsEmptyCard(mCard) {
				for _, pos := range poss {
					bs1 := BoolsOfAction{}
					bs1[i] = true
					bs2 := BoolsOfAction{}
					bs2[j] = true
					action := Action{
						Bools1:bs1,
						Bools2:bs2,
						BattlePosition:pos,
						Type:NORMAL_SUMMON_ACTION,
					}
					y = append(y, action)
				}
			}
		}
	}
	return y
}

func NewLegalAttackDeclareActions(state *State) Actions {
	if state.Phase != BATTLE_PHASE {
		return Actions{}
	}

	y := make(Actions, 0, MONSTER_ZONE_LENGTH * MONSTER_ZONE_LENGTH)
	for i, p1Card := range state.P1.MonsterZone {
		for j, p2Card := range state.P2.MonsterZone {
			if !p1Card.IsAttackDeclared && !IsEmptyCard(p1Card) && !IsEmptyCard(p2Card) {
				bs1 := BoolsOfAction{}
				bs1[i] = true
				bs2 := BoolsOfAction{}
				bs2[j] = true
				action := Action{
					Bools1:bs1,
					Bools2:bs2,
					Type:ATTACK_DECLARE_ACTION,
				}
				y = append(y, action)
			}
		}
	}
	return y
}

func NewLegalDirectAttackDeclareActions(state *State) Actions {
	if state.Phase != BATTLE_PHASE || !fn.All(state.P2.MonsterZone, IsEmptyCard) {
		return Actions{}
	}
	y := make(Actions, 0, MONSTER_ZONE_LENGTH)
	for i, card := range state.P1.MonsterZone {
		if !IsEmptyCard(card) && !card.IsAttackDeclared {
			bs1 := BoolsOfAction{}
			bs1[i] = true
			action := Action{
				Bools1:bs1,
				Type:ATTACK_DECLARE_ACTION,
			}
			y = append(y, action)
		}
		return y
	}
	return Actions{}
}

func NewLegalActions(state *State) Actions {
	phaseTransition := NewLegalPhaseTransitionActions(state)
	normalSummon := NewLegalNormalSummonActions(state)
	attackDeclared := NewLegalAttackDeclareActions(state)
	directAttackDeclare := NewLegalDirectAttackDeclareActions(state)

	n := len(phaseTransition) +
		len(normalSummon) +
		len(directAttackDeclare) +
		len(attackDeclared)

	y := make(Actions, 0, n)
	y = append(y, phaseTransition...)
	y = append(y, normalSummon...)
	y = append(y, directAttackDeclare...)
	y = append(y, attackDeclared...)
	return y
}

func TypesOfActions(actions Actions) ActionTypes {
	return fn.Map[ActionTypes](actions, GetTypeOfAction)
}