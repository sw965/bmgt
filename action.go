package bmgt

import (
	"golang.org/x/exp/slices"
	"github.com/sw965/omw/fn"
)

type ActionType int

const (
	PHASE_TRANSITION_ACTION ActionType = iota
	NORMAL_SUMMON_ACTION
	BATTLE_POSITION_CHANGE_ACTION
	ATTACK_DECLARE_ACTION
)

func ActionTypeToString(t ActionType) string {
	switch t {
		case PHASE_TRANSITION_ACTION:
			return "フェイズ移行"
		case NORMAL_SUMMON_ACTION:
			return "通常召喚"
		case BATTLE_POSITION_CHANGE_ACTION:
			return "表示形式変更"
		case ATTACK_DECLARE_ACTION:
			return "攻撃宣言"
		default:
			return ""
	}
}

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

func GetTypeOfAction(action Action) ActionType {
	return action.Type
}

func IsBadAction(state *State) func(Action) bool {
	return func(action Action) bool {
		switch action.Type {
			case NORMAL_SUMMON_ACTION:
				idxs := action.Indices1()
				card := state.P1.Hand[idxs[0]]
				switch card.Name {
					case LUSTER_DRAGON:
						return action.BattlePosition != ATK_BATTLE_POSITION
					case GEMINI_ELF:
						return action.BattlePosition != ATK_BATTLE_POSITION
					case VORSE_RAIDER:
						return action.BattlePosition != ATK_BATTLE_POSITION
				}
			case ATTACK_DECLARE_ACTION:
				idxs1 := action.Indices1()
				idxs2 := action.Indices2()
				if len(idxs2) != 0 {
					p1Card := state.P1.MonsterZone[idxs1[0]]
					p2Card := state.P2.MonsterZone[idxs2[0]]
					if p2Card.BattlePosition == ATK_BATTLE_POSITION {
						if p1Card.Atk < p2Card.Atk {
							return true
						}
					}
				}
		}
		return false
	}
}

func IsNotBadAction(state *State) func(Action) bool {
	return func(action Action) bool {
		return !IsBadAction(state)(action)
	}
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
	y := make(Actions, 0, len(state.P1.Hand) * MONSTER_ZONE_LENGTH * len(poss))
	for i, hCard := range state.P1.Hand {
		for j, mCard := range state.P1.MonsterZone {
			isLow := slices.Contains(LOW_LEVELS, hCard.Level)
			isEmpty := IsEmptyCard(mCard) 
			if isLow && isEmpty {
				bs1 := BoolsOfAction{}
				bs1[i] = true
				bs2 := BoolsOfAction{}
				bs2[j] = true

				for _, pos := range poss {
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

func NewLegalBattlePositionChangeActions(state *State) Actions {
	if state.Phase != MAIN_PHASE {
		return Actions{}
	}

	y := make(Actions, 0, MONSTER_ZONE_LENGTH)
	for i, card := range state.P1.MonsterZone {
		if !IsEmptyCard(card) && card.IsBattlePositionChangeable {
			var pos BattlePosition
			switch card.BattlePosition {
				case ATK_BATTLE_POSITION:
					pos = FACE_UP_DEF_BATTLE_POSITION
				case FACE_UP_DEF_BATTLE_POSITION:
					pos = ATK_BATTLE_POSITION
				case FACE_DOWN_DEF_BATTLE_POSITION:
					pos = ATK_BATTLE_POSITION
			}

			bs1 := BoolsOfAction{}
			bs1[i] = true
			action := Action{
				Bools1:bs1,
				BattlePosition:pos,
				Type:BATTLE_POSITION_CHANGE_ACTION,
			}
			y = append(y, action)
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
			isNotEmpty := !IsEmptyCard(p1Card) && !IsEmptyCard(p2Card)
			canDeclare := !p1Card.IsAttackDeclared && (p1Card.BattlePosition == ATK_BATTLE_POSITION)
			if isNotEmpty && canDeclare {
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
		isNotEmpty := !IsEmptyCard(card)
		canDeclare := !card.IsAttackDeclared && card.BattlePosition == ATK_BATTLE_POSITION
		if isNotEmpty && canDeclare {
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
	isP2Turn := state.Turn%2 == 0
	if isP2Turn {
		stateV := state.Reverse()
		state = &stateV
	}

	phaseTransition := NewLegalPhaseTransitionActions(state)
	normalSummon := NewLegalNormalSummonActions(state)
	battlePositionChange := NewLegalBattlePositionChangeActions(state)
	attackDeclared := NewLegalAttackDeclareActions(state)
	directAttackDeclare := NewLegalDirectAttackDeclareActions(state)

	n := len(phaseTransition) +
		len(normalSummon) +
		len(battlePositionChange) +
		len(directAttackDeclare) +
		len(attackDeclared)

	y := make(Actions, 0, n)
	y = append(y, phaseTransition...)
	y = append(y, normalSummon...)
	y = append(y, battlePositionChange...)
	y = append(y, directAttackDeclare...)
	y = append(y, attackDeclared...)
	return y
}

func TypesOfActions(actions Actions) ActionTypes {
	return fn.Map[ActionTypes](actions, GetTypeOfAction)
}