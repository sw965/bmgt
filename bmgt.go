package bmgt

import (
	"fmt"
	"golang.org/x/exp/slices"
	"math/rand"
)

type CardName string
type Attribute string
type Type string

type Card struct {
	Name CardName
	Level int
	Atk int
	Def int

	Attribute Attribute
	Type Type

	IsNormalMonster bool
	IsEffectMonster bool

    IsNormalSpell bool
    IsQuickPlaySpell bool
    IsContinuousSpell bool

	IsNormalTrap bool
	IsContinuousTrap bool
	IsCounterTrap bool
}

type Cards []Card

var EXODIA_DECK = func() Cards {
	result := Cards{
		*CARD_DATA_BASE["封印されしエクゾディア"],
		*CARD_DATA_BASE["封印されしエクゾディア"],
		*CARD_DATA_BASE["封印されしエクゾディア"],
		*CARD_DATA_BASE["封印されしエクゾディア"],
		*CARD_DATA_BASE["封印されしエクゾディア"],
		*CARD_DATA_BASE["封印されしエクゾディア"],

		*CARD_DATA_BASE["封印されし者の左腕"],
		*CARD_DATA_BASE["封印されし者の左腕"],
		*CARD_DATA_BASE["封印されし者の左腕"],
		*CARD_DATA_BASE["封印されし者の左腕"],
		*CARD_DATA_BASE["封印されし者の左腕"],
		*CARD_DATA_BASE["封印されし者の左腕"],

		*CARD_DATA_BASE["封印されし者の右腕"],
		*CARD_DATA_BASE["封印されし者の右腕"],
		*CARD_DATA_BASE["封印されし者の右腕"],
		*CARD_DATA_BASE["封印されし者の右腕"],
		*CARD_DATA_BASE["封印されし者の右腕"],
		*CARD_DATA_BASE["封印されし者の右腕"],

		*CARD_DATA_BASE["封印されし者の左足"],
		*CARD_DATA_BASE["封印されし者の左足"],
		*CARD_DATA_BASE["封印されし者の左足"],
		*CARD_DATA_BASE["封印されし者の左足"],
		*CARD_DATA_BASE["封印されし者の左足"],
		*CARD_DATA_BASE["封印されし者の左足"],

		*CARD_DATA_BASE["封印されし者の右足"],
		*CARD_DATA_BASE["封印されし者の右足"],
		*CARD_DATA_BASE["封印されし者の右足"],
		*CARD_DATA_BASE["封印されし者の右足"],
		*CARD_DATA_BASE["封印されし者の右足"],
		*CARD_DATA_BASE["封印されし者の右足"],

		*CARD_DATA_BASE["強欲な壺"],
		*CARD_DATA_BASE["強欲な壺"],
		*CARD_DATA_BASE["強欲な壺"],
		*CARD_DATA_BASE["強欲な壺"],
		*CARD_DATA_BASE["強欲な壺"],
		*CARD_DATA_BASE["強欲な壺"],
		*CARD_DATA_BASE["強欲な壺"],
		*CARD_DATA_BASE["強欲な壺"],
		*CARD_DATA_BASE["強欲な壺"],
		*CARD_DATA_BASE["強欲な壺"],
	}
	return result
}()

func (cards Cards) Shuffle(r *rand.Rand) Cards {
	cards = slices.Clone(cards)
	r.Shuffle(len(cards), func(i, j int) {cards[i], cards[j] = cards[j], cards[i]})
	return cards
}

func (cards Cards) Draw(num int) (Cards, Cards, error) {
	n := len(cards)
	if n < num {
		return Cards{}, Cards{}, fmt.Errorf("ドローしようとした枚数 > 残りの枚数")
	}
	draw := make(Cards, num)
	result := make(Cards, n-num)

	for i, card := range cards {
		if i > num {
			draw[i] = card
		} else {
			result = append(result, card)
		}
	}

	return draw, result, nil
}

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
	Self OneSideState
	Opponent OneSideState
	Turn int
	Priority int
	Phase Phase
}

func NewInitState(selfDeck, opponentDeck Cards, r *rand.Rand) (State, error) {
	selfDeck = selfDeck.Shuffle(r)
	opponentDeck = opponentDeck.Shuffle(r)

	selfHand, selfDeck, err := selfDeck.Draw(5)
	if err != nil {
		return State{}, err
	}

	opponentHand, opponentDeck, err := opponentDeck.Draw(5)
	if err != nil {
		return State{}, err
	}

	self := OneSideState{}
	self.LifePoint = 8000
	self.Hand = selfHand
	self.Deck = selfDeck

	opponent := OneSideState{}
	opponent.LifePoint = 8000
	opponent.Hand = opponentHand
	opponent.Deck = opponentDeck

	state := State{Self:self, Opponent:opponent}
	state.Turn = 0
	state.Priority = 0
	state.Phase = MAIN1_PHASE
	return state
}