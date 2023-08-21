package bmgt

import (
	"math/rand"
	"golang.org/x/exp/slices"
	omwrand "github.com/sw965/omw/rand"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
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
	Actions Actions
}

func NewOneSideState(deck Cards, r *rand.Rand, startID CardID) OneSideState {
	n := len(deck)

	setID := func(i int, card Card) Card {
		card.ID = startID + CardID(i)
		return card
	}

	deck = fn.MapIndex[Cards](deck, setID)
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

func (oss *OneSideState) Draw(num int) {
	hand := make(Cards, 0, len(oss.Hand) + num)
	hand = append(hand, oss.Hand...)
	hand = append(hand, oss.Deck[:num]...)
	oss.Hand = hand
	oss.Deck = oss.Deck[num:]
}

//手札を墓地へ捨てる
func (oss *OneSideState) Discard(idxs []int) {
	hand := make(Cards, 0, len(oss.Hand) - len(idxs)) 
	for i, card := range oss.Hand {
		if slices.Contains(idxs, i) {
			oss.Graveyard = append(oss.Graveyard, card)
		} else {
			hand = append(hand, card)
		}
	}
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

func (state *State) LegalActions() Actions {
	if len(state.P1.Actions) != 0 {
		p1LastAction := state.P1.Actions[len(state.P1.Actions)-1]
		switch p1LastAction.CardName {
			case THUNDER_DRAGON:
				if p1LastAction.IsCardActivation {
					return Actions{Action{CardName:THUNDER_DRAGON, IsCost:true}}
				}
		}
	}

	y := make(Actions, 0, 128)

	monsterZoneEmptyIndices := omwslices.IndicesFunc(state.P1.MonsterZone, func(card Card) bool { return card.Name == NO_NAME })
	tributeSummonCostIndices := omwslices.IndicesFunc(
		state.P1.MonsterZone,
		func(card Card) bool {
			return card.Name != NO_NAME && card.Name != SUMMONER_MONK
		},
	)
	spellTrapZoneEmptyIndices := omwslices.IndicesFunc(state.P1.SpellTrapZone, func(card Card) bool { return card.Name == NO_NAME })

	isSpellSpeed1Activatable := state.Phase.IsMainPhase() && len(state.Chain) == 0
	isSpellSpeed2Activatable := len(state.Chain) == 0 || state.Chain[len(state.Chain)-1].Card.Category != COUNTER_TRAP

	if isSpellSpeed1Activatable {
		for handI, card := range state.P1.Hand {
			if card.Category.IsMonster() {
				if slices.Contains(LOW_LEVELS, card.Level) {
					for _, empI := range monsterZoneEmptyIndices {
						action := Action{
							CardName:card.Name,
							HandIndices:[]int{handI},
							MonsterZoneIndices:[]int{empI},
							IsNormalSummon:true,
						}
						y = append(y, action)
					}
				}
	
				var tributeCost int
				if slices.Contains(MEDIUM_LEVELS, card.Level) {
					tributeCost = 1
				} else {
					tributeCost = 2
				}
	
				if tributeCost >= len(tributeSummonCostIndices) {
					idxss := omwslices.Combination[[][]int, []int](tributeSummonCostIndices, tributeCost)
					for _, idxs := range idxss {
						action := Action{
							CardName:card.Name,
							MonsterZoneIndices:idxs,
							IsNormalSummon:true,
							IsTributeSummonCost:true,
						}
						y = append(y, action)
					}
				}
			}

			if f, ok := HAND_SPELL_SPEED1_MONSTER_ACTIVATABLE[card.Name]; ok {
				if omwslices.Any(f(state, &card)) {
					action := Action{
						CardName:card.Name,
						HandIndices:[]int{handI},
						IsCardActivation:true,
					}
					y = append(y, action)
				}
			}

			if f, ok := HAND_NORMAL_SPELL_ACTIVATABLE[card.Name]; ok {
				if omwslices.Any(f(state, &card)) {
					for _, zoneI := range spellTrapZoneEmptyIndices {
						action := Action{
							CardName:card.Name,
							HandIndices:[]int{handI},
							SpellTrapZoneIndices:[]int{zoneI},
							IsCardActivation:true,
						}
						y = append(y, action)
					}
				}
			}

			if card.Category.IsSpell() || card.Category.IsTrap() {
				for _, empI := range spellTrapZoneEmptyIndices {
					action := Action{
						HandIndices:[]int{handI},
						SpellTrapZoneIndices:[]int{empI},
						IsSpellTrapSet:true,
					}
					y = append(y, action)
				}
			}
		}

		for zoneI, card := range state.P1.SpellTrapZone {
			if f, ok := ZONE_NORMAL_SPELL_ACTIVATABLE[card.Name]; ok {
				if omwslices.Any(f(state, &card)) {
					action := Action{
						CardName:card.Name,
						SpellTrapZoneIndices:[]int{zoneI},
						IsCardActivation:true,
					}
					y = append(y, action)
				}
			}
		}
	}

	if isSpellSpeed2Activatable {
		if state.P1.IsTurn {
			for handI, card := range state.P1.Hand {
				if f, ok := HAND_QUICK_PLAY_SPELL_ACTIVATABLE[card.Name]; ok {
					if omwslices.Any(f(state, &card)) {
						for _, zoneI := range spellTrapZoneEmptyIndices {
							action := Action{
								CardName:card.Name,
								HandIndices:[]int{handI},
								SpellTrapZoneIndices:[]int{zoneI},
								IsCardActivation:true,
							}
							y = append(y, action)
						}
					}
				}
			}
		}

		for zoneI, card := range state.P1.SpellTrapZone {
			if f, ok := ZONE_QUICK_PLAY_SPELL_ACTIVATABLE[card.Name]; ok {
				if omwslices.Any(f(state, &card)) && !card.IsSetTurn {
					action := Action{
						CardName:card.Name,
						SpellTrapZoneIndices:[]int{zoneI},
						IsCardActivation:true,
					}
					y = append(y, action)
				}
			}

			if f, ok := ZONE_NORMAL_TRAP_ACTIVATABLE[card.Name]; ok {
				if omwslices.Any(f(state, &card)) && !card.IsSetTurn {
					action := Action{
						CardName:card.Name,
						SpellTrapZoneIndices:[]int{zoneI},
						IsCardActivation:true,
					}
					y = append(y, action)
				}
 			}
		}
	}
	return y
}

func (state State) Push(action *Action) State {
	if action.IsNormalSummon {
		if action.IsTributeSummonCost {
			state.P1.MonsterZoneToGraveyard(action.MonsterZoneIndices)
		}
	}
	return state
}