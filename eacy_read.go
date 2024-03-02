package bmgt

import (
	"github.com/sw965/omw/fn"
)

var CARD_NAME_TO_STRING = map[CardName]string{
	NO_NAME:"",
	DARK_MAGICIAN_GIRL:"ブラック・マジシャン・ガール",
	EXODIA_THE_FORBIDDEN_ONE:"封印されしエクゾディア",
	LEFT_ARM_OF_THE_FORBIDDEN_ONE:"封印されし者の左腕",
	LEFT_LEG_OF_THE_FORBIDDEN_ONE:"封印されし者の左足",
	RIGHT_ARM_OF_THE_FORBIDDEN_ONE:"封印されし者の右腕",
	RIGHT_LEG_OF_THE_FORBIDDEN_ONE:"封印されし者の左足",
	POT_OF_GREED:"強欲な壺",
	MAGICAL_STONE_EXCAVATION:"魔法石の採掘",
}

var ATTRIBUTE_TO_STRING = map[Attribute]string{
	DARK:"闇",
	LIGHT:"光",
	EARTH:"地",
	WATER:"水",
	FIRE:"炎",
	WIND:"風",
}

var TYPE_TO_STRING = map[Type]string{
	DRAGON:"ドラゴン",
	SPELLCASTET:"魔法使い",
	ZOMBIE:"アンデット",
	WARRIOR:"戦士",
	BEAST_WARRIOR:"獣戦士",
	BEAST:"獣",
	FIEND:"悪魔",
	FAIRY:"天使",
	INSECT:"昆虫",
	DINOSAUR:"恐竜",
	REPTILE:"爬虫類",
	FISH:"魚",
	SEA_SERPENT:"海竜",
	MACHINE:"機械",
	THUNDER:"雷",
	AQUA:"水",
	PYRO:"炎",
	ROCK:"岩石",
	PLANT:"植物",
	PSYCHIC:"サイキック",
	WYRM:"幻竜",
	CYBERSE:"サイバース",
	ILLUSION:"幻想魔",
	DIVINE_BEAST:"幻神獣",
	CREATOR_GOD:"創造神",
}

type EasyReadCard struct {
	Name string
	Level Level
	Attribute string
	Type string
	Atk int
	Def int
	Face string
	Orientation string
	ID CardID
}

func NewEasyReadCard(card Card) EasyReadCard {
	return EasyReadCard{
		Name:CARD_NAME_TO_STRING[card.Name],
		Level:card.Level,
		Attribute:ATTRIBUTE_TO_STRING[card.Attribute],
		Type:TYPE_TO_STRING[card.Type],
		Atk:card.Atk,
		Def:card.Def,
		Face:card.Face.ToString(),
		Orientation:card.Orientation.ToString(),
		ID:card.ID,
	}
}

type EasyReadCards []EasyReadCard

func NewEasyReadCards(cards Cards) EasyReadCards {
	return fn.Map[EasyReadCards](cards, NewEasyReadCard)
}

type EasyReadOneSide struct {
	LifePoint LifePoint
	Hand EasyReadCards
	Deck EasyReadCards
	MonsterZone EasyReadCards
	SpellTrapZone EasyReadCards
	Graveyard EasyReadCards
}

func NewEasyReadOneSide(o *OneSide) EasyReadOneSide {
	return EasyReadOneSide{
		LifePoint:o.LifePoint,
		Hand:NewEasyReadCards(o.Hand),
		Deck:NewEasyReadCards(o.Deck),
		MonsterZone:NewEasyReadCards(o.MonsterZone),
		SpellTrapZone:NewEasyReadCards(o.SpellTrapZone),
		Graveyard:NewEasyReadCards(o.Graveyard),
	}
}

type EasyReadDuel struct {
	P1 EasyReadOneSide
	P2 EasyReadOneSide
	Phase string
}

func NewEasyReadDuel(duel *Duel) EasyReadDuel {
	return EasyReadDuel{
		P1:NewEasyReadOneSide(&duel.P1),
		P2:NewEasyReadOneSide(&duel.P2),
		Phase:duel.Phase.ToString(),
	}
}