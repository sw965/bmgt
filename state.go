package bmgt

import (
	omwrand "github.com/sw965/omw/rand"
	omwslices "github.com/sw965/omw/slices"
	"math/rand"
)

type LifePoint int

const INIT_LIFE_POINT = 8000

type Phase int

const (
	DRAW_PHASE Phase = iota
	STANDBY_PHASE
	MAIN_PHASE
	BATTLE_PHASE
	END_PHASE
)

const (
	MONSTER_ZONE_LENGTH = 5
)

type OneSideState struct {
	LifePoint LifePoint
	Deck Cards
	Hand Cards
	MonsterZone Cards
	Graveyard Cards
	IsTurn bool
	IsNormalDrawDone bool
	ThisTurnNormalSummonCount int
	AttackDeclareIndex int
}

func NewInitOneSideState(deck Cards, r *rand.Rand) OneSideState {
	newDeck := omwrand.Shuffled(deck, r)
	newDeck = newDeck[5:]
	hand := newDeck[:5]

	y := OneSideState{}
	y.LifePoint = INIT_LIFE_POINT
	y.Deck = newDeck
	y.Hand = hand
	return y
}

func (oss *OneSideState) Draw(num int) {
	newDeck := oss.Deck[num:]
	oss.Hand = append(oss.Hand, oss.Deck[:num]...)
	oss.Deck = newDeck
}

type State struct {
	P1 OneSideState
	P2 OneSideState
	Phase Phase
}

func NewInitState(p1Deck, p2Deck Cards, r *rand.Rand) State {
	p1 := NewInitOneSideState(p1Deck, r)
	p2 := NewInitOneSideState(p2Deck, r)
	return State{P1:p1, P2:p2}
}

func Push(state State, action Action) {
	switch action.Type {
		case PHASE_TRANSITION_ACTION:
			var next Phase
			switch state.Phase {
				case DRAW_PHASE:
					next = STANDBY_PHASE
				case STANDBY_PHASE:
					next = MAIN_PHASE
				case BATTLE_PHASE:
					next = END_PHASE
				case END_PHASE:
					next = DRAW_PHASE
					state.P1.IsTurn = false
					state.P2.IsTurn = true
					state.P1.ThisTurnNormalSummonCount = 0
			}
			state.Phase = next
		case NORMAL_DRAW_ACTION:
			state.P1.Draw(1)
		case NORMAL_SUMMON_ACTION:
			newHand, summonCards := omwslices.Delete(state.P1.Hand, action.HandIndices[0])
			state.P1.Hand = newHand
			idx := action.MonsterZoneIndices1[0]
			state.P1.MonsterZone[idx] = summonCards[0]
			state.P1.MonsterZone[idx].BattlePosition = action.BattlePosition
			state.P1.ThisTurnNormalSummonCount += 1
		case ATTACK_DECLARE_ACTION:
			if IsDirectAttackDeclareAction(action) {
				state.P2.LifePoint -= action.MonsterZoneIndices1[0].Atk
			} else {
				atkIdx := action.MonsterZoneIndices1[0]
				defIdx := action.MonsterZoneIndices2[0]
				attackingCard := state.P1.MonsterZone[atkIdx]
				defendingCard := state.P2.MonsterZone[defIdx]
				p1Atk := attackingCard.Atk
				p2Atk := defendingCard.Def

				if defendingCard.Position == ATK_BATTLE_POSITION {
					if p1Atk > p2Atk {
						state.P2.MonsterZone[p2Idx] = Card{}
						state.P2.Graveyard = append(state.P2.Graveyard, defendingCard)
						state.P2.LifePoint -= p1Atk - p2Atk
					} else if p1Atk < p2Atk {
						state.P1.MonsterZone[p1Idx] = Card{}
						state.P1.Graveyard = append(state.P1.Graveyard, attackingCard)
						state.P1.LifePoint -= p2Atk - p1Atk
					} else {
						state.P1.MonsterZone[p1Idx] = Card{}
						state.P1.Graveyard = append(state.P1.Graveyard, attackingCard)
						state.P2.MonsterZone[p2Idx] = Card{}
						state.P2.Graveyard = append(state.P2.Graveyard, defendingCard)
					}
				} else if defendingCard.BattlePosition.IsDef() {
					state.P2.MonsterZone[p2Idx].BattlePosition = FACE_UP_ATK_POSITION
					if p1Atk > p2Deck {
						state.P1.MonsterZone[ak]
					}
				}
	}
}