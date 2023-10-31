package bmgt

import (
	omwrand "github.com/sw965/omw/rand"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
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

func PhaseToString(phase Phase) string {
	switch phase {
		case DRAW_PHASE:
			return "ドローフェイズ"
		case STANDBY_PHASE:
			return "スタンバイフェイズ"
		case MAIN_PHASE:
			return "メインフェイズ"
		case BATTLE_PHASE:
			return "バトルフェイズ"
		case END_PHASE:
			return "エンドフェイズ"
		default:
			return ""
	}
}

const (
	INIT_HAND_NUM = 5
	MIN_DECK_NUM = 40
	MAX_DECK_NUM = 60
	MONSTER_ZONE_LENGTH = 5
	SPELL_TRAP_ZONE_LENGTH = 5
)

type OneSideState struct {
	LifePoint LifePoint
	Deck Cards
	Hand Cards
	MonsterZone Cards
	SpellTrapZone Cards
	Graveyard Cards
	IsNormalDrawDone bool
	ThisTurnNormalSummonCount int
	IsDeckDeath bool
}

func NewInitOneSideState(deck Cards, start CardID, r *rand.Rand) OneSideState {
	deck = fn.MapIndex[Cards](deck, SetIDOfCard, start)
	shuffledDeck := omwrand.Shuffled(deck, r)
	newDeck := shuffledDeck[INIT_HAND_NUM:]
	hand := shuffledDeck[:INIT_HAND_NUM]

	y := OneSideState{}
	y.LifePoint = INIT_LIFE_POINT
	y.Deck = newDeck
	y.Hand = hand
	y.MonsterZone = make(Cards, MONSTER_ZONE_LENGTH)
	y.SpellTrapZone = make(Cards, SPELL_TRAP_ZONE_LENGTH)
	y.Graveyard = make(Cards, 0, MAX_DECK_NUM)
	return y
}

func (oss *OneSideState) Draw(num int) {
	newDeck := oss.Deck[num:]
	oss.Hand = append(oss.Hand, oss.Deck[:num]...)
	oss.Deck = newDeck
}

func (atkOSS *OneSideState) Battle(defOSS *OneSideState, atkIdx, defIdx int) {
	destroy := func(oss *OneSideState, idx, dmg int) {
		card := oss.MonsterZone[idx]
		oss.MonsterZone[idx] = Card{}
		oss.Graveyard = append(oss.Graveyard, card)
		oss.LifePoint -= LifePoint(dmg)
	}

	attackDeclared := func() {
		atkOSS.MonsterZone[atkIdx].IsAttackDeclared = true
	}

	atkCard := atkOSS.MonsterZone[atkIdx]
	defCard := defOSS.MonsterZone[defIdx]

	if defCard.BattlePosition == ATK_BATTLE_POSITION {
		if atkCard.Atk > defCard.Atk {
			destroy(defOSS, defIdx, atkCard.Atk - defCard.Def)
			attackDeclared()
		} else if atkCard.Atk < defCard.Atk {
			destroy(atkOSS, atkIdx, defCard.Atk - atkCard.Atk)
			attackDeclared()
		} else {
			destroy(atkOSS, atkIdx, 0)
			destroy(defOSS, defIdx, 0)
		}
	} else {
		defOSS.MonsterZone[defIdx].BattlePosition = FACE_UP_DEF_BATTLE_POSITION
		if atkCard.Atk > defCard.Def {
			destroy(defOSS, defIdx, 0)
			attackDeclared()
		} else if atkCard.Atk < defCard.Def {
			atkOSS.LifePoint -= LifePoint(defCard.Def - atkCard.Atk)
			attackDeclared()
		}
	}
}

type State struct {
	P1 OneSideState
	P2 OneSideState
	Phase Phase
	Turn int
}

func NewInitState(p1Deck, p2Deck Cards, isStartDraw bool, r *rand.Rand) State {
	p1Deck = fn.Map[Cards](p1Deck, SetIsP1OfCard(true))
	p2Deck = fn.Map[Cards](p2Deck, SetIsP1OfCard(false))
	p1 := NewInitOneSideState(p1Deck, 1, r)
	p2 := NewInitOneSideState(p2Deck, CardID(len(p1Deck)+1), r)
	state := State{P1:p1, P2:p2, Phase:DRAW_PHASE, Turn:1}
	if isStartDraw {
		state.P1.Draw(1)
	}
	return state
}

func (state *State) Reverse() State {
	return State{P1:state.P2, P2:state.P1, Phase:state.Phase, Turn:state.Turn}
}

func (state *State) IsP1Turn() bool {
	return state.Turn%2 == 1
}

func (state *State) ChangeTurn() {
	state.P1.ThisTurnNormalSummonCount = 0
	for i := 0; i < MONSTER_ZONE_LENGTH; i++ {
		state.P1.MonsterZone[i].IsBattlePositionChangeable = true
	}

	for i := 0; i < MONSTER_ZONE_LENGTH; i++ {
		state.P1.MonsterZone[i].IsAttackDeclared = false
	}

	state.Turn += 1
	state.Phase = DRAW_PHASE

	if len(state.P2.Deck) == 0 {
		state.P2.IsDeckDeath = true
	} else {
		state.P2.Draw(1)
	}
}

func (state *State) ChangeBattlePosition(action *Action) {
	idx := action.Indices1()[0]
	state.P1.MonsterZone[idx].BattlePosition = action.BattlePosition
	state.P1.MonsterZone[idx].IsBattlePositionChangeable = false
}

func (state *State) NormalSummon(action *Action) {
	handIdx := action.Indices1()[0]
	mZoneIdx := action.Indices2()[0]

	newHand, summonCards := omwslices.Delete(state.P1.Hand, handIdx)
	state.P1.Hand = newHand
	state.P1.MonsterZone[mZoneIdx] = summonCards[0]
	state.P1.MonsterZone[mZoneIdx].BattlePosition = action.BattlePosition
	state.P1.MonsterZone[mZoneIdx].IsBattlePositionChangeable = false
	state.P1.ThisTurnNormalSummonCount += 1
}

func (state *State) AttackDeclare(action *Action) {
	p1MZoneIdx := action.Indices1()[0]
	p2MZoneIdxs := action.Indices2()
	isDirectAttackDeclare := len(p2MZoneIdxs) == 0

	if isDirectAttackDeclare {
		state.P2.LifePoint -= LifePoint(state.P1.MonsterZone[p1MZoneIdx].Atk)
	} else {
		p2MZoneIdx := p2MZoneIdxs[0]
		state.P1.Battle(&state.P2, p1MZoneIdx, p2MZoneIdx)
	}
}

func Push(state State, action Action) State {
	state.P1.Deck = CloneCards(state.P1.Deck)
	state.P1.Hand = CloneCards(state.P1.Hand)
	state.P1.MonsterZone = CloneCards(state.P1.MonsterZone)
	state.P1.Graveyard = CloneCards(state.P1.Graveyard)

	state.P2.Deck = CloneCards(state.P2.Deck)
	state.P2.Hand = CloneCards(state.P2.Hand)
	state.P2.MonsterZone = CloneCards(state.P2.MonsterZone)
	state.P2.Graveyard = CloneCards(state.P2.Graveyard)

	isP2Turn := !state.IsP1Turn()
	if isP2Turn {
		state = state.Reverse()
	}

	switch action.Type {
		case PHASE_TRANSITION_ACTION:
			state.Phase = action.Phase
			if state.Phase == END_PHASE {
				state.ChangeTurn()
			}
		case NORMAL_SUMMON_ACTION:
			state.NormalSummon(&action)
		case BATTLE_POSITION_CHANGE_ACTION:
			state.ChangeBattlePosition(&action)
		case ATTACK_DECLARE_ACTION:
			state.AttackDeclare(&action)
	}

	if isP2Turn {
		state = state.Reverse()
	}
	return state
}

func IsGameEnd(state *State) bool {
	return state.P1.LifePoint <= 0 || state.P2.LifePoint <= 0 || state.P1.IsDeckDeath || state.P2.IsDeckDeath
}