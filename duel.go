package bmgt

import (
	"math/rand"
	omwrand "github.com/sw965/omw/rand"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
)

const (
	INIT_DRAW = 5
	MONSTER_ZONE_LENGTH = 5
	SPELL_TRAP_ZONE_LENGTH = 5
)

type LifePoint int

const (
	INIT_LIFE_POINT LifePoint = 8000
)

type Phase int

const (
	DRAW_PHASE Phase = iota
	STANDBY_PHASE
	MAIN1_PHASE
	BATTLE_PHASE
	MAIN2_PHASE
	END_PHASE
)

type Phases []Phases

var PHASES = Phases{DRAW_PHASE, STANDBY_PHASE, MAIN1_PHASE, BATTLE_PHASE, MAIN2_PHASE, END_PHASE}

type OneSideDuelState struct {
	LifePoint LifePoint
	Hand Cards
	Deck Cards
	MonsterZone Cards
	SpellTrapZone Cards
	Graveyard Cards
}

func NewOneSideDuelState(deck Cards, r *rand.Rand) OneSideDuelState {
	deck = omwrand.Shuffled(deck, r)
	hand := deck[:INIT_DRAW]
	deck = deck[INIT_DRAW:]
	return OneSideDuelState{
		LifePoint:INIT_LIFE_POINT,
		Hand:hand, Deck:deck,
		MonsterZone:make(Cards, MONSTER_ZONE_LENGTH),
		SpellTrapZone:make(Cards, SPELL_TRAP_ZONE_LENGTH),
		Graveyard:make(Cards, len(deck) + len(hand)),
	}
}

type Duel struct {
	P1 OneSideDuelState
	P2 OneSideDuelState
	Phase Phase
}

func NewDuel(p1Deck, p2Deck Cards, r *rand.Rand) Duel {
	p1 := NewOneSideDuelState(p1Deck, r)
	p2 := NewOneSideDuelState(p2Deck, r)
	return Duel{P1:p1, P2:p2}
}

type ActionType int

const (
	PHASE_TRANSITION_ACTION ActionType = iota
	NORMAL_SUMMON_ACTION
)

type Action struct {
	N1 int
	N2 int
	Type ActionType
}

type Actions []*Action

func NewPhaseTransitionActions() Actions {
	
}

func NewNormalSummonActions(duel *Duel) Actions {
	is := omwslices.NewSequentialInteger[[]int](0, len(duel.P1.Hand))
	js := omwslices.NewSequentialInteger[[]int](0, MONSTER_ZONE_LENGTH)
	f := func(i, j int) *Action {
		return &Action{N1:i, N2:j, Type:NORMAL_SUMMON_ACTION}
	}
	return omwslices.Product2[[]int, []int, Actions](is, js, f)
}

func IsNormalAction(duel *Duel) func(*Action) bool {
	return func(action *Action) bool {	
		if !CanNormalSummonCard(duel.P1.Hand[action.N1]) {
			return false
		}

		if duel.P1.MonsterZone[action.N2].Name == NO_NAME {
			return false
		}
		return  true
	}
}

func NewNormalSummonLegalActions(duel *Duel) Actions {
	as := NewNormalSummonActions(duel)
	return fn.Filter(as, IsNormalAction(duel))
}