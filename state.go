package bmgt

import (
	"fmt"
	"github.com/sw965/omw/fn"
	omwrand "github.com/sw965/omw/rand"
	omws "github.com/sw965/omw/slices"
	"math/rand"
)

type LifePoint int

const (
	INIT_LIFE_POINT = 8000
)

const (
	MONSTER_ZONE_LENGTH    = 5
	SPELL_TRAP_ZONE_LENGTH = 5
)

type OneSideState struct {
	LifePoint     LifePoint
	Hand          Cards
	Deck          Cards
	MonsterZone   Cards
	SpellTrapZone Cards
	Graveyard     Cards

	IsPriorityWaiver                  bool
	CurrentTurnNormalSummonUpperLimit int
	CurrentTurnNormalSummonNum        int
	IsDeclareAnAttack                 bool
	OncePerTurn                       CardNames
}

func NewOneSideState(deck Cards, r *rand.Rand, startID CardID) (OneSideState, error) {
	var hand Cards
	var err error

	deck = fn.Map[Cards, Cards](deck, CloneCard)
	for i := range deck {
		deck[i].ID = startID + CardID(i)
	}

	deck = omwrand.Shuffled(deck, r)
	deck, hand, err = deck.Draw(5)

	result := OneSideState{}
	result.LifePoint = INIT_LIFE_POINT
	result.Hand = hand
	result.Deck = deck
	result.MonsterZone = make(Cards, MONSTER_ZONE_LENGTH)
	result.SpellTrapZone = make(Cards, SPELL_TRAP_ZONE_LENGTH)
	result.Graveyard = make(Cards, 0, len(deck))
	return result, err
}

func (oss OneSideState) Draw(num int) (OneSideState, error) {
	deck, drawCards, err := oss.Deck.Draw(num)
	oss.Deck = deck
	oss.Hand = append(oss.Hand, drawCards...)
	return oss, err
}

// 手札を捨てる
func (oss OneSideState) Discard(idxs []int) OneSideState {
	hand, gy := omws.Pops(oss.Hand, idxs)
	oss.Hand = hand
	oss.Graveyard = append(oss.Graveyard, gy...)
	return oss
}

// デッキから手札に加える
func (oss OneSideState) Search(idxs []int, r *rand.Rand) OneSideState {
	var cards Cards
	oss.Deck, cards = omws.Pops(oss.Deck, idxs)
	oss.Hand = append(oss.Hand, cards...)
	oss.Deck = omwrand.Shuffled(oss.Deck, r)
	return oss
}

// 墓地から手札に加える
func (oss OneSideState) Salvage(idxs []int) OneSideState {
	var cards Cards
	oss.Hand, cards = omws.Pops(oss.Graveyard, idxs)
	oss.Hand = append(oss.Hand, cards...)
	return oss
}

func (oss *OneSideState) CanNormalSummon() bool {
	return oss.CurrentTurnNormalSummonNum < oss.CurrentTurnNormalSummonUpperLimit
}

type Phase string

const (
	DRAW_PHASE    = Phase("ドロー")
	STANDBY_PHASE = Phase("スタンバイ")
	MAIN1_PHASE   = Phase("メイン1")
	BATTLE_PHASE  = Phase("バトル")
	MAIN2_PHASE   = Phase("メイン2")
	END_PHASE     = Phase("エンド")
)

type State struct {
	P1           OneSideState
	P2           OneSideState
	IsP1Turn     bool
	IsP1Priority bool
	Phase        Phase
	Chain        Chain

	SelectCards              Cards
	EffectProcessingCardName CardName
	EffectProcessingNumber   int

	//一時休戦
	OneDayOfPeace bool
}

func NewInitState(p1Deck, p2Deck Cards, r *rand.Rand) (State, error) {
	p1, err := NewOneSideState(p1Deck, r, 0)
	if err != nil {
		return State{}, err
	}
	p2, err := NewOneSideState(p2Deck, r, CardID(len(p1Deck)))
	p1Deck = omwrand.Shuffled(p1Deck, r)
	p2Deck = omwrand.Shuffled(p2Deck, r)

	state := State{P1: p1, P2: p2}
	state.IsP1Turn = true
	state.IsP1Priority = true
	state.Phase = DRAW_PHASE
	return state, err
}

func (state *State) IsMainPhase() bool {
	return state.Phase == MAIN1_PHASE || state.Phase == MAIN2_PHASE
}

func (state *State) CanSpellSpeed1Activation() bool {
	return state.IsMainPhase() && len(state.Chain) == 0
}

func (state State) Print() {
	fmt.Println(state.P2.Hand.Names())
	fmt.Println(state.P2.SpellTrapZone.Names())
	fmt.Println(state.P2.MonsterZone.Names())
	fmt.Println(state.P1.MonsterZone.Names())
	fmt.Println(state.P1.SpellTrapZone.Names())
	fmt.Println(state.P1.Hand.Names())
}

type StateTransition func(State) (State, error)
