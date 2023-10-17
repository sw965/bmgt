package bmgt

import (
	"golang.org/x/exp/slices"
	"github.com/sw965/omw/fn"
)

type ActionType int

const (
	PHASE_TRANSITION_ACTION ActionType = iota
	NORMAL_DRAW_ACTION
	NORMAL_SUMMON_ACTION
	ATTACK_DECLARE_ACTION
)

type Action struct {
	HandIndices []int
	MonsterZoneIndices1 []int
	MonsterZoneIndices2 []int
	BattlePosition BattlePosition
	Phase Phase
	Type ActionType
}

func NewPhaseTransitionAction(phase Phase) Action {
	return Action{Phase:phase, Type:PHASE_TRANSITION_ACTION}
}

func IsDirectAttackDeclareAction(action Action) bool {
	return action.Type == ATTACK_DECLARE_ACTION && len(action.MonsterZoneIndices2) == 0
}

func IsBadAction(state *State) func(action Action) bool {
	return func(action Action) bool {
		//エクゾディアパーツを召喚する行為
		if action.Type == NORMAL_SUMMON_ACTION {
			idx := action.HandIndices[0]
			return slices.Contains(EXODIA_PARTS_NAMES, state.P1.Hand[idx].Name)
		}
		return false
	}
}

type Actions []Action

//責務として、適切なフェイズ・召喚権は考慮しない
func NewNormalSummonActions(state *State) Actions {
	y := make(Actions, 0, len(state.P1.Hand) * MONSTER_ZONE_LENGTH)
	for i, hCard := range state.P1.Hand {
		for j, mCard := range state.P1.MonsterZone {
			if slices.Contains(LOW_LEVELS, hCard.Level) && IsEmptyCard(mCard) {
				for _, pos := range NORMAL_SUMMON_BATTLE_POSITIONS {
					action := Action{
						HandIndices:[]int{i},
						MonsterZoneIndices1:[]int{j},
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

func NewDirectAttackDeclareActions(state *State) Actions {
	if fn.All(state.P2.MonsterZone, IsEmptyCard) {
		y := make(Actions, 0, MONSTER_ZONE_LENGTH)
		for i, card := range state.P1.MonsterZone {
			if !IsEmptyCard(card) {
				action := Action{
					MonsterZoneIndices1:[]int{i},
					Type:ATTACK_DECLARE_ACTION,
				}
				y = append(y, action)
			}
			return y
		}
	}
	return Actions{}
}

func NewAttackDeclareActions(state *State) Actions {
	y := make(Actions, 0, (MONSTER_ZONE_LENGTH * MONSTER_ZONE_LENGTH) + MONSTER_ZONE_LENGTH)
	y = append(y, NewDirectAttackDeclareActions(state)...)
	for i, p1Card := range state.P1.MonsterZone {
		for j, p2Card := range state.P2.MonsterZone {
			if !p1Card.IsAttackDeclared && !IsEmptyCard(p1Card) && !IsEmptyCard(p2Card) {
				action := Action{
					MonsterZoneIndices1:[]int{i},
					MonsterZoneIndices2:[]int{j},
					Type:ATTACK_DECLARE_ACTION,
				}
				y = append(y, action)
			}
		}
	}
	return y
}

func NewLegalActions(state *State) Actions {
	y := make(Actions, 0, 128)
	switch state.Phase {
		case DRAW_PHASE:
			if !state.P1.IsNormalDrawDone {
				y = append(y, Action{Type:NORMAL_DRAW_ACTION})
			} else {
				y = append(y, NewPhaseTransitionAction(STANDBY_PHASE))
			}
		case STANDBY_PHASE:
			y = append(y, NewPhaseTransitionAction(MAIN_PHASE))
		case MAIN_PHASE:
			if state.P1.ThisTurnNormalSummonCount == 0 {
				y = append(y, NewNormalSummonActions(state)...)
			}
			y = append(y, NewPhaseTransitionAction(BATTLE_PHASE))
			y = append(y, NewPhaseTransitionAction(END_PHASE))
		case BATTLE_PHASE:
			if state.P1.AttackDeclareIndex == -1 {
				y = append(y, NewPhaseTransitionAction(END_PHASE))
			} else {
				y = append(y, NewAttackDeclareActions(state)...)
			}
	}
	return y
}