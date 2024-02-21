package bmgt

import (
	"golang.org/x/exp/slices"
)

type CardName int

const (
	NO_NAME CardName = iota
	DARK_MAGICIAN_GIRL
)

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

type Card struct {
	Name CardName
	Level Level
	Attribute Attribute
	Type Type
	Atk int
	Def int
}

func CanNormalSummonCard(card Card) bool {
	return slices.Contains(LOW_LEVELS, card.Level)
}

type Cards []Card