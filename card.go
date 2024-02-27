package bmgt

import (
	"golang.org/x/exp/slices"
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

var CARD_NAME_TO_STRING()

func CardNameToString(name CardName) string {
	switch name {
		case NO_NAME:
			return ""
		case DARK_MAGICIAN_GIRL:
			return "ブラック・マジシャン・ガール"
		case EXODIA_THE_FORBIDDEN_ONE:
			return "封印されしエクゾディア"
		case LEFT_ARM_OF_THE_FORBIDDEN_ONE:
			return "封印されし者の左腕"
		case LEFT_LEG_OF_THE_FORBIDDEN_ONE:
			return "封印されし者の左足"
		case RIGHT_ARM_OF_THE_FORBIDDEN_ONE:
			return "封印されし者の右腕"
		case RIGHT_LEG_OF_THE_FORBIDDEN_ONE:
			return "封印されし者の右足"
		case POT_OF_GREED:
			return "強欲な壺"
		case MAGICAL_STONE_EXCAVATION:
			return "魔法石の採掘"
		default:
			return "存在しないカード名"
	}
}

func StringToCardName(s string) CardName {
	
}

type CardNames []CardName

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

type Type int

const (
	DRAGON Attribute = iota
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
	WYRM
	CYBERSE
	ILLUSION
	DIVINE_BEAST //幻神獣(三幻神)
	CREATOR_GOD //創造神(ホルアクティ)
)

type Face int

const (
	FACE_UP Face = iota
	FACE_DOWN
)

type Orientation int

const (
	VERTICAL Orientation = iota
	HORIZONTAL
)

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

type Card struct {
	Name CardName
	Level Level
	Attribute Attribute
	Type Type
	Atk int
	Def int
	Face Face
	Orientation Orientation
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

func (cards Cards) GetNames() CardNames {
	return fn.Map(cards, GetNameOfCard)
}

func (cards Cards) IsAllEmpty() bool {
	return fn.All(cards, func(card Card) bool { return card.Name == NO_NAME })
}