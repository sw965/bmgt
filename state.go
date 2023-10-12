package bmgt

import (
	"math/rand"
	omwrand "github.com/sw965/omw/rand"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
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
	ThisTurnNormalSummonOrSetLimit int
	ThisTurnNormalSummonOrSetCount int

	OncePerTurnLimitCardNames CardNames

	NormalSummonIndex int
	NormalSummonSuccessIndex int
	TributeSummonIndex int
	TributeSummonSuccessIndex int
	FlipSummonIndex int
	FlipSummonSuccessIndex int
	SpecialSummonIndex int
	SpecialSummonSuccessIndex int

	SpellActivationIndex int
	TrapActivationIndex int

	IsOneDayOfPeaceEndTrigger bool
	IsDoubleSummonApplied bool

	EffectProcessingNumber int
}

func NewOneSideState(deck Cards, r *rand.Rand, startID CardID) OneSideState {
	n := len(deck)

	deck = fn.MapIndex[Cards](deck, SetIDOfCard, startID)
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

func (oss *OneSideState) CanNormalSummon() bool {
	return oss.ThisTurnNormalSummonOrSetLimit > oss.ThisTurnNormalSummonOrSetCount
}

func (oss *OneSideState) Draw(num int) {
	draw := oss.Deck[:num]
	deck := oss.Deck[num:]
	oss.Hand = append(oss.Hand, draw...)
	oss.Deck = deck
}

func (oss *OneSideState) Discard(idxs []int) {
	hand, cards := omwslices.Delete(oss.Hand, idxs...)
	oss.Hand = hand
	oss.Graveyard = append(oss.Graveyard, cards...)
}

func (oss *OneSideState) DiscardBanish(idxs []int) {
	hand, cards := omwslices.Delete(oss.Hand, idxs...)
	oss.Hand = hand
	oss.Banish = append(oss.Banish, cards...)
}

func (oss *OneSideState) IDSalvage(ids CardIDs) {
	n := len(ids)
	hand := make(Cards, 0, len(oss.Hand) + n)
	for i, card := range oss.Hand {
		hand[i] = card
	}

	gy := make(Cards, 0, len(oss.Graveyard) - n)
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

func (oss *OneSideState) Search(idxs []int, shuffle bool, r *rand.Rand) {
	newDeck, addHand := omwslices.Delete(oss.Deck, idxs...)
	oss.Hand = append(oss.Hand, addHand...)
	oss.Deck = newDeck

	if shuffle {
		omwrand.Shuffle(oss.Deck, r)
	}
}

func (oss *OneSideState) HandToDeck(idxs []int, shuffle bool, r *rand.Rand) {
	newHand, addDeck := omwslices.Delete(oss.Hand, idxs...)
	oss.Hand = newHand
	oss.Deck = append(oss.Deck, addDeck...)
	if shuffle {
		omwrand.Shuffle(oss.Deck, r)
	}
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

func IsMainPhase(phase Phase) bool {
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
	return state
}

func NewLegalNormalSummonOrSetActions(state *State, pos BattlePosition) Actions {
	if !IsMainPhase(state.Phase) || !state.P1.CanNormalSummon() || len(state.Chain) != 0 {
		return Actions{}
	}

	y := make(Actions, 0, 128)
	for i := 0; i < len(state.P1.Hand); i++ {
		card := state.P1.Hand[i]
		if CanNormalSummonCard(card) {
			for j, mCard := range state.P1.MonsterZone {
				if IsEmptyCard(mCard) {
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

func NewLegalTributeSummonOrSetActions(state *State, pos BattlePosition) Actions {
	if !IsMainPhase(state.Phase) || !state.P1.CanNormalSummon() || len(state.Chain) != 0 { 
		return Actions{}
	}

	tributeCostIndices := omwslices.IndicesFunc(state.P1.MonsterZone, CanTributeSummonCostCard)
	monsterZoneEmptyIndices := omwslices.IndicesFunc(state.P1.MonsterZone, IsEmptyCard)
	costPossible := len(tributeCostIndices)
	if costPossible == 0 {
		return Actions{}
	}

	y := make(Actions, 0, 128)
	for i := 0; i < len(state.P1.Hand); i++ {
		card := state.P1.Hand[i]
		if CanTributeSummonCard(card) {
			cost := TributeSummonCostOfCard(card)
			if cost <= costPossible {
				tributeIdxss := omwslices.Combination[[][]int](tributeCostIndices, cost)
				for _, idxs := range tributeIdxss {
					summonIdxs := make([]int, 0, )
					summonIdxs = append(summonIdxs, idxs...)
					summonIdxs = append(summonIdxs, monsterZoneEmptyIndices...)
					action := Action{
						MonsterZoneIndices1:idxs,
						MonsterZoneIndices2:summonIdxs,
						BattlePosition:pos,
						Type:TRIBUTE_SUMMON_ACTION,
					}
					y = append(y, action)
				}
			}
		}
	}
	return y
}

func NewLegalHandSpellSpeed1SpellActions(state *State) Actions {
	if IsMainPhase(state.Phase) || len(state.Chain) != 0 {
		return Actions{}
	}

	spellTrapZoneEmptyIndices := omwslices.IndicesFunc(state.P1.SpellTrapZone, IsEmptyCard)
	y := make(Actions, 0, 128)

	for i := 0; i < len(state.P1.Hand); i++ {
		card := state.P1.Hand[i]
		if activatable, ok := HAND_SPELL_SPEED1_SPELL_ACTIVATABLE[card.Name]; ok {
			if activatable(state) {
				for _, j := range spellTrapZoneEmptyIndices {
					action := Action{
						HandIndices:[]int{i},
						SpellTrapZoneIndices:[]int{j},
						Type:HAND_SPELL_ACTIVATION_ACTION,
					}
					y = append(y, action)
				}
			}
		}
	}
	return y
}

func NewLegalHandQuickPlaySpellActions(state *State) Actions {
	end, err := omwslices.End(state.Chain)

	if err == nil && CARD_DATA_BASE[end.Card.Name].Category == Category(COUNTER_TRAP) {
		return Actions{}
	}

	y := make(Actions, 0, 128)
	spellTrapZoneEmptyIndices := omwslices.IndicesFunc(state.P1.MonsterZone, IsEmptyCard)

	for i := 0; i < len(state.P1.Hand); i++ {
		card := state.P1.Hand[i]
		if activatable, ok := HAND_QUICK_PLAY_SPELL_ACTIVATABLE[card.Name]; ok {
			if activatable(state) {
				for _, j := range spellTrapZoneEmptyIndices {
					action := Action{
						HandIndices:[]int{i},
						SpellTrapZoneIndices:[]int{j},
						Type:HAND_SPELL_ACTIVATION_ACTION,
					}
					y = append(y, action)
				}
			}
		}
	}
	return y
}

func NewLegalZoneNormalTrapActions(state *State) Actions {
	end, err := omwslices.End(state.Chain)

	if err == nil && CARD_DATA_BASE[end.Card.Name].Category == Category(COUNTER_TRAP) {
		return Actions{}
	}

	y := make(Actions, 0, 128)
	for i, card := range state.P1.SpellTrapZone {
		if activatable, ok := ZONE_NORMAL_TRAP_ACTIVATABLE[card.Name]; !card.IsSetTurn && ok {
			if activatable(state) {
				action := Action{
					SpellTrapZoneIndices:[]int{i},
					Type:ZONE_SPELL_TRAP_ACTIVATION_ACTION,
				}
				y = append(y, action)
			}
		}
	}
	return y
}

func NewLegalHandSpellTrapSetActions(state *State) Actions {
	if IsMainPhase(state.Phase) || len(state.Chain) != 0 { 
		return Actions{}
	}

	y := make(Actions, 0, 128)
	spellTrapZoneEmptyIndices := omwslices.IndicesFunc(state.P1.SpellTrapZone, IsEmptyCard)

	for i := 0; i < len(state.P1.Hand); i++ {
		card := state.P1.Hand[i]
		category := CARD_DATA_BASE[card.Name].Category
		if IsSpellCategory(category) || IsTrapCategory(category) {
			for _, j := range spellTrapZoneEmptyIndices {
				action := Action{
					HandIndices:[]int{i},
					SpellTrapZoneIndices:[]int{j},
					Type:HAND_SPELL_TRAP_SET_ACTION,
				}
				y = append(y, action)
			}
		}
	}
	return y
}

func NewLegalActions(state *State) Actions {
	//BattlePosではなく、裏表のいずれかの状態の型を作る魔法トラップも同様

	legalNormalSummonActions := NewLegalNormalSummonOrSetActions(state, ATTACK_POSITION)
	legalNormalSetActions := NewLegalNormalSummonOrSetActions(state, FACE_DOWN_DEFENSE_POSITION)
	legalTributeSummonActions := NewLegalTributeSummonOrSetActions(state, ATTACK_POSITION)
	legalTributeSetActions := NewLegalTributeSummonOrSetActions(state, FACE_DOWN_DEFENSE_POSITION)

	legalHandSpellSpeed1SpellActions := NewLegalHandSpellSpeed1SpellActions(state)
	legalHandQuickPlaySpellActions := NewLegalHandQuickPlaySpellActions(state)

	legalHandSpellTrapSetActions := NewLegalHandSpellTrapSetActions(state)

	n := len(legalNormalSummonActions) +
		len(legalNormalSetActions) +
		len(legalTributeSummonActions) +
		len(legalTributeSetActions) +
		len(legalHandSpellSpeed1SpellActions) +
		len(legalHandQuickPlaySpellActions) +
		len(legalHandSpellSpeed1SpellActions) +
		len(legalHandSpellTrapSetActions)

	y := make(Actions, 0, n)
	y = append(y, legalNormalSummonActions...)
	y = append(y, legalNormalSetActions...)
	y = append(y, legalTributeSummonActions...)
	y = append(y, legalTributeSetActions...)
	y = append(y, legalHandSpellSpeed1SpellActions...)
	y = append(y, legalHandQuickPlaySpellActions...)
	y = append(y, legalHandSpellTrapSetActions...)
	return y
}