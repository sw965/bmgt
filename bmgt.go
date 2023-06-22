package bmgt

import (
	"math/rand"
	omwrand "github.com/sw965/omw/rand"
)

const (
	MONSTER_ZONE_SIZE = 5
	SPELL_AND_TRAP_ZONE_SIZE = 5
)

type OneSideState struct {
	LifePoint int
	Hand Cards
	Deck Cards
	MonsterZone [MONSTER_ZONE_SIZE]*Card
	SpellAndTrapZone [SPELL_AND_TRAP_ZONE_SIZE]*Card
	Graveyard Cards
}

type Phase string

const (
	DRAW_PHASE = Phase("ドロー")
	STANDBY_PHASE = Phase("スタンバイ")
	MAIN1_PHASE = Phase("メイン1")
	BATTLE_PHASE = Phase("バトル")
	MAIN2_PHASE = Phase("メイン2")
	END_PHASE = Phase("エンド")
)

type State struct {
	P1 OneSideState
	P2 OneSideState
	Turn int
	Priority int
	Phase Phase
}

func NewInitState(p1Deck, p2Deck Cards, r *rand.Rand) (State, error) {
	p1Deck = omwrand.Shuffled(p1Deck, r)
	p2Deck = omwrand.Shuffled(p2Deck, r)

	p1Hand, p1Deck, err := p1Deck.Draw(5)
	if err != nil {
		return State{}, err
	}

	p2Hand, p2Deck, err := p2Deck.Draw(5)
	if err != nil {
		return State{}, err
	}

	p1 := OneSideState{}
	p1.LifePoint = 8000
	p1.Hand = p1Hand
	p1.Deck = p1Deck

	p2 := OneSideState{}
	p2.LifePoint = 8000
	p2.Hand = p2Hand
	p2.Deck = p2Deck

	state := State{P1:p1, P2:p2}
	state.Turn = 0
	state.Priority = 0
	state.Phase = MAIN1_PHASE
	return state, nil
}