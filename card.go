package bmgt

import (
	"golang.org/x/exp/slices"
	omwmaps "github.com/sw965/omw/maps"
	"github.com/sw965/omw/fn"
)

type CardName int

const (
	NO_NAME CardName = iota
	DARK_MAGICIAN_GIRL

	EXODIA_THE_FORBIDDEN_ONE
	LEFT_ARM_OF_THE_FORBIDDEN_ONE
	LEFT_LEG_OF_THE_FORBIDDEN_ONE
	RIGHT_ARM_OF_THE_FORBIDDEN_ONE
	RIGHT_LEG_OF_THE_FORBIDDEN_ONE

	POT_OF_GREED
	MAGICAL_STONE_EXCAVATION
)

var STRING_TO_CARD_NAME = omwmaps.Reverse[map[string]CardName](CARD_NAME_TO_STRING)

type CardNames []CardName

var EXODIA_PARTS_NAMES = CardNames{
	EXODIA_THE_FORBIDDEN_ONE,
	LEFT_ARM_OF_THE_FORBIDDEN_ONE,
	LEFT_LEG_OF_THE_FORBIDDEN_ONE,
	RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
	RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
}

type Level int

type Levels []Level

var LOW_LEVELS = Levels{1, 2, 3, 4}

type Attribute int

const (
	DARK Attribute = iota
	LIGHT
	EARTH
	WATER
	FIRE
	WIND
)

var STRING_TO_ATTRIBUTE = omwmaps.Reverse[map[string]Attribute](ATTRIBUTE_TO_STRING)

type Type int

const (
	DRAGON Type = iota
	SPELLCASTET // 魔法使い
	ZOMBIE //アンデット
	WARRIOR //戦士
	BEAST_WARRIOR
	BEAST
	WINGED_BEAST //鳥獣
	FIEND //悪魔
	FAIRY
	INSECT
	DINOSAUR //恐竜
	REPTILE //爬虫類
	FISH
	SEA_SERPENT //海竜
	MACHINE
	THUNDER
	AQUA
	PYRO
	ROCK
	PLANT
	PSYCHIC
	WYRM //幻竜
	CYBERSE //サイバース
	ILLUSION //幻想魔
	DIVINE_BEAST //幻神獣(三幻神)
	CREATOR_GOD //創造神(ホルアクティ)
)

var STRING_TO_TYPE = omwmaps.Reverse[map[string]Type](TYPE_TO_STRING)

type Face int

const (
	FACE_UP Face = iota
	FACE_DOWN
)

func (face Face) ToString() string {
	switch face {
		case FACE_UP:
			return "表"
		case FACE_DOWN:
			return "裏"
		default:
			return ""
	}
}

type Orientation int

const (
	VERTICAL Orientation = iota
	HORIZONTAL
)

func (o Orientation) ToString() string {
	switch o {
		case VERTICAL:
			return "縦"
		case HORIZONTAL:
			return "裏"
		default:
			return ""
	}
}

type BattlePosition int

const (
	FACE_UP_ATTACK_POSITION BattlePosition = iota
	FACE_UP_DEFENSE_POSITION
	FACE_DOWN_DEFENSE_POSITION
)

func NewBattlePosition(face Face, o Orientation) BattlePosition {
	return map[Face]map[Orientation]BattlePosition{
		FACE_UP:map[Orientation]BattlePosition{VERTICAL:FACE_UP_ATTACK_POSITION, HORIZONTAL:FACE_UP_DEFENSE_POSITION},
		FACE_DOWN:map[Orientation]BattlePosition{HORIZONTAL:FACE_DOWN_DEFENSE_POSITION},
	}[face][o]
}

type CardID int

type Card struct {
	Name CardName
	Level Level
	Attribute Attribute
	Type Type
	Atk int
	Def int
	Face Face
	Orientation Orientation
	ID CardID
}

func NewCard(name CardName) Card {
	if name == NO_NAME {
		return Card{}
	} else {
		y := Card{}
		data := CARD_DATA_BASE[name]
		y.Name = name
		y.Level = data.Level
		y.Attribute = data.Attribute
		y.Type = data.Type
		y.Atk = data.Atk
		y.Def = data.Def
		y.Face = FACE_DOWN
		y.Orientation = VERTICAL
		return y
	}
}

func (card *Card) BattlePosition() BattlePosition {
	return NewBattlePosition(card.Face, card.Orientation)
}

func (card *Card) SetBattlePosition(pos BattlePosition) {
	upAtk := func() {
		card.Face = FACE_UP
		card.Orientation = VERTICAL
	}

	upDef := func() {
		card.Face = FACE_UP
		card.Orientation = HORIZONTAL
	}

	downDef := func() {
		card.Face = FACE_DOWN
		card.Orientation = HORIZONTAL
	}

	map[BattlePosition]func(){
		FACE_UP_ATTACK_POSITION:upAtk,
		FACE_UP_DEFENSE_POSITION:upDef,
		FACE_DOWN_DEFENSE_POSITION:downDef,
	}[pos]()
}

func GetNameOfCard(card Card) CardName {
	return card.Name
}

func CanNormalSummonCard(card Card) bool {
	return slices.Contains(LOW_LEVELS, card.Level)
}

type Cards []Card

func NewCards(names ...CardName) Cards {
	return fn.Map[Cards](names, NewCard)
}

func (cards Cards) Names() CardNames {
	return fn.Map[CardNames](cards, GetNameOfCard)
}

func (cards Cards) IsAllEmpty() bool {
	return fn.All(cards, func(card Card) bool { return card.Name == NO_NAME })
}