package bmgt

import (
	"fmt"
	omws "github.com/sw965/omw/slices"
)

type CardName string
type CardNames []CardName

var EXODIA_PART_NAMES = CardNames{
	"封印されしエクゾディア",
	"封印されし者の左腕",
	"封印されし者の右腕",
	"封印されし者の左足",
	"封印されし者の右足",
}

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

func NewCards(names ...CardName) Cards {
	result := make(Cards, len(names))
	for i, name := range names {
		result[i] = *CARD_DATA_BASE[name]
	}
	return result
}

func (cards Cards) Names() CardNames {
	result := make(CardNames, len(cards))
	for i, card := range cards {
		result[i] = card.Name
	}
	return result
}

func (cards Cards) Draw(num int) (Cards, Cards, error) {
	drawCards := make(Cards, num)
	for i := 0; i < num; i++ {
		if len(cards) == 0 {
			return cards, drawCards, fmt.Errorf("ドローしようとしたが、カードがなかった")
		}
		var drawCard Card
		cards, drawCard = omws.Pop(cards, 0)
		drawCards[i] = drawCard
	}
	return cards, drawCards, nil
}