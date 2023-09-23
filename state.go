package bmgt

import (
	"math"
	"math/rand"
	"golang.org/x/exp/slices"
	omwrand "github.com/sw965/omw/rand"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
)

const (
	MEDIUM_LEVEL_TRIBUTE_SUMMON_COST = 1
	HIGH_LEVEL_TRIBUTE_SUMMON_COST = 2
)

const (
	INIT_DRAW = 5
	MAX_DECK_NUM = 60
	MAX_EX_DECK = 15
	GRAVEYARD_CAP = MAX_DECK_NUM + MAX_EX_DECK
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
	Banish Cards

	IsTurn bool

	OncePerTurnLimitCardNames CardNames

	NormalSummonIndex int
	NormalSummonSuccessIndex int
	FlipSummonIndex int
	FlipSummonSuccessIndex int
	SpecialSummonIndex int
	SpecialSummonSuccessIndex int

	SpellCardActivationIndex int
	TrapCardActivationIndex int

	AfterNormalDraw bool
	AfterDraw bool

	IsOneDayOfPeaceEndTrigger bool
}

func NewOneSideState(deck Cards, r *rand.Rand, startID CardID) OneSideState {
	n := len(deck)

	deck = fn.MapIndex[Cards](deck, CardF.SetID, startID)
	deck = omwrand.Shuffled(deck, r)
	hand := deck[:INIT_DRAW]
	deck = deck[INIT_DRAW:]

	result := OneSideState{}
	result.LifePoint = INIT_LIFE_POINT
	result.Hand = hand
	result.Deck = deck
	result.MonsterZone = make(Cards, MONSTER_ZONE_LENGTH)
	result.SpellTrapZone = make(Cards, SPELL_TRAP_ZONE_LENGTH)
	result.Graveyard = make(Cards, 0, n)
	result.Banish = make(Cards, 0, n)
	return result
}

func (oss *OneSideState) NormalSummon(handI, zoneI int, pos BattlePosition) {
	hand, summonCards := omwslices.Delete(oss.Hand, handI)
	oss.MonsterZone[zoneI] = summonCards[0]
	oss.NormalSummonIndex = zoneI
	oss.Hand = hand
}

//手札を除外ゾーンへ捨てる
func (oss *OneSideState) DiscardBanish(idxs []int) {
	hand := make(Cards, 0, len(oss.Hand) - len(idxs))
	for i, card := range oss.Hand {
		if slices.Contains(idxs, i) {
			oss.Banish = append(oss.Banish, card)
		} else {
			hand = append(hand, card)
		}
	}
	oss.Hand = hand
}

//デッキから手札へ
func (oss *OneSideState) Search(idxs []int, r *rand.Rand) {
	n := len(idxs)
	hand := make(Cards, 0, len(oss.Hand) + n)
	hand = append(hand, oss.Hand...)
	deck := make(Cards, 0, len(oss.Deck) - n)

	for i, card := range oss.Deck {
		if slices.Contains(idxs, i) {
			hand = append(hand, card)
		} else {
			deck = append(deck, card)
		}
	}

	r.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	oss.Hand = hand
	oss.Deck = deck
}

//墓地から手札へ
func (oss *OneSideState) IDSalvage(ids CardIDs) {
	hand := make(Cards, 0, len(oss.Hand) + 1)
	hand = append(hand, oss.Hand...)
	gy := make(Cards, 0, GRAVEYARD_CAP - 1)
	for _, card := range oss.Graveyard {
		if slices.Contains(ids, card.ID) {
			hand = append(hand, card)
		} else {
			gy = append(gy, card)
		}
	}
	oss.Hand = hand
	oss.Graveyard = gy
}

func (oss *OneSideState) HandToDeck(idxs []int, r *rand.Rand) {
	n := len(idxs)
	hand := make(Cards, 0, len(oss.Hand)-n)
	deck := make(Cards, 0, len(oss.Deck)+n)
	deck = append(deck, oss.Deck...)
	for i, card := range oss.Hand {
		if slices.Contains(idxs, i) {
			deck = append(deck, card)
		} else {
			hand = append(hand, card)
		}
	}

	r.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	oss.Hand = hand
	oss.Deck = deck
}

func (oss *OneSideState) MonsterZoneToGraveyard(zoneIdxs []int) {
	for i := 0; i < MONSTER_ZONE_LENGTH; i++ {
		if slices.Contains(zoneIdxs, i) {
			card := oss.MonsterZone[i]
			oss.Graveyard = append(oss.Graveyard, card)
			oss.MonsterZone[i] = Card{}
		}
	}
}

func (oss *OneSideState) HandToSpellTrapZone(handI, zoneI int) {
	hand := make(Cards, 0, len(oss.Hand)-1)
	for i, card := range oss.Hand {
		if i == handI {
			oss.SpellTrapZone[zoneI] = card
		} else {
			hand = append(hand, card)
		}
	}
	oss.Hand = hand
}

func (oss *OneSideState) HandToMonsterZone(handIdx, zoneIdx int) {
	hand := make(Cards, 0, len(oss.Hand)-1)
	for handI, card := range oss.Hand {
		if handI == handIdx {
			oss.MonsterZone[zoneIdx] = oss.Hand[handIdx]
		} else {
			hand = append(hand, card)
		}
	}
	oss.Hand = hand
}

type Phase int

const (
	DRAW_PHASE Phase = iota
	STANDBY_PHASE
	MAIN1_PHASE
	BATTLE_PHASE
	MAIN2_PHASE
	END_PHASE
)

func (phase Phase) IsMainPhase() bool {
	return phase == MAIN1_PHASE || phase == MAIN2_PHASE
}

type State struct {
	P1           OneSideState
	P2           OneSideState
	Phase        Phase
	Chain Chain
}

func NewInitState(p1Deck, p2Deck Cards, r *rand.Rand) State {
	p1 := NewOneSideState(p1Deck, r, 0)
	p2 := NewOneSideState(p2Deck, r, CardID(len(p1Deck)))
	p1.IsTurn = true
	p2.IsTurn = false
	state := State{P1: p1, P2: p2}
	state.Phase = DRAW_PHASE
	state.Chain = make(Chain, 0, 16)
	return state
}

type stateF struct{}
var StateF = stateF{}

type StateChanger func(*State)

type stateChangerF struct{}
var StateChangerF = stateChangerF{}

func (f *stateChangerF) ReversePlayer1AndPlayer2(state *State) {
	p1 := state.P1
	p2 := state.P2
	state.P1 = p2
	state.P2 = p1
}

func (f *stateChangerF) PayHalfLifePoint(state *State) {
	state.P1.LifePoint -= LifePoint(math.Round(float64(state.P1.LifePoint) / 2.0))
}

func (f *stateChangerF) HandToDeck(idxs []int, r *rand.Rand) StateChanger {
	return func(state *State) {
		var deck Cards
		state.P1.Hand, deck = omwslices.Delete(state.P1.Hand, idxs...)
		state.P1.Deck = append(state.P1.Deck, deck...)
		f.Shuffle(r)(state)
	}
}

func (f *stateChangerF) Draw(num int) StateChanger {
	return func(state *State) {
		hand := make(Cards, 0, len(state.P1.Hand) + num)
		hand = append(hand, state.P1.Hand...)
		hand = append(hand, state.P1.Deck[:num]...)
		state.P1.Hand = hand
		state.P1.Deck = state.P1.Deck[num:]
	}
}

func (f *stateChangerF) Shuffle(r *rand.Rand) StateChanger {
	return func(state *State) {
		omwrand.Shuffle(state.P1.Deck, r)
	}
}

func (f *stateChangerF) Discard(idxs []int) StateChanger {
	return func(state *State) {
		hand, gy := omwslices.Delete(state.P1.Hand, idxs...)
		state.P1.Hand = hand
		state.P1.Graveyard = append(state.P1.Graveyard, gy...)
	}
}

func (f *stateChangerF) NegateNormalSummon(destroy bool) StateChanger {
	return func(state *State) {
		p1Idx := state.P1.NormalSummonIndex
		p2Idx := state.P2.NormalSummonIndex

		negate := func(zone Cards, idx int) {
			zone[idx].NegatedNormalSummon = true
			zone[idx].Destroyed = destroy
		}
	
		if p1Idx != -1 {
			negate(state.P1.MonsterZone, p1Idx)
			return
		}
	
		if p2Idx != -1 {
			negate(state.P2.MonsterZone, p2Idx)
			return
		}
	}
}

func (f *stateChangerF) NegateFlipSummon(destroy bool) StateChanger {
	return func(state *State) {
		p1Idx := state.P1.FlipSummonIndex
		p2Idx := state.P2.FlipSummonIndex

		negate := func(zone Cards, idx int) {
			zone[idx].NegatedFlipSummon = true
			zone[idx].Destroyed = destroy
		}
	
		if p1Idx != -1 {
			negate(state.P1.MonsterZone, p1Idx)
			return
		}
	
		if p2Idx != -1 {
			negate(state.P2.MonsterZone, p2Idx)
			return
		}
	}
}

func (f *stateChangerF) NegateSpecialSummon(destroy bool) StateChanger {
	return func(state *State) {
		p1Idx := state.P1.SpecialSummonIndex
		p2Idx := state.P2.SpecialSummonIndex

		negate := func(zone Cards, idx int) {
			zone[idx].NegatedSpecialSummon = true
			zone[idx].Destroyed = destroy
		}
	
		if p1Idx != -1 {
			negate(state.P1.MonsterZone, p1Idx)
			return
		}
	
		if p2Idx != -1 {
			negate(state.P2.MonsterZone, p2Idx)
			return
		}
	}
}

func (f *stateChangerF) MonsterZoneSpellCounterRemoval(idx, removal int) StateChanger {
	return func(state *State) {
		state.P1.MonsterZone[idx].SpellCounter -= removal
	}
}

func (f *stateChangerF) NegateSpellCardActivation(destroy bool) StateChanger {
	return func(state *State) {
		p1Idx := state.P1.SpellCardActivationIndex
		p2Idx := state.P2.SpellCardActivationIndex

		negate := func(zone Cards, idx int) {
			zone[idx].NegatedCardActivation = true
			zone[idx].Destroyed = destroy
		}
	
		if p1Idx != -1 {
			negate(state.P1.SpellTrapZone, p1Idx)
			return
		}
	
		if p2Idx != -1 {
			negate(state.P2.SpellTrapZone, p2Idx)
			return
		}
	}
}

func (f *stateChangerF) NegateTrapCardActivation(destroy bool) StateChanger {
	return func(state *State) {
		p1Idx := state.P1.TrapCardActivationIndex
		p2Idx := state.P2.TrapCardActivationIndex

		negate := func(zone Cards, idx int) {
			zone[idx].NegatedCardActivation = true
			zone[idx].Destroyed = destroy
		}
	
		if p1Idx != -1 {
			negate(state.P1.SpellTrapZone, p1Idx)
			return
		}
	
		if p2Idx != -1 {
			negate(state.P2.SpellTrapZone, p2Idx)
			return
		}
	}
}

func (f *stateChangerF) OneDayOfPeace(state *State) {
	state.P2.IsOneDayOfPeaceEndTrigger = true
}

func (f *stateChangerF) SolemnJudgment(state *State) {
	f.NegateNormalSummon(true)
	f.NegateFlipSummon(true)
	f.NegateSpecialSummon(true)
	f.NegateSpellCardActivation(true)
	f.NegateTrapCardActivation(true)
}

type StateChangers []StateChanger
type StateChangerss []StateChangers