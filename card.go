package bmgt

import (
	"fmt"
	omws "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
	"strings"
)

type CardName string

const (
	TOON = "トゥーン"
)

type CardNames []CardName

var EXODIA_PART_NAMES = CardNames{
	"封印されしエクゾディア",
	"封印されし者の左腕",
	"封印されし者の右腕",
	"封印されし者の左足",
	"封印されし者の右足",
}

type Level int
type Levels []Level

var LOW_LEVELS = Levels{1, 2, 3, 4}
var MEDIUM_LEVELS = Levels{5, 6}

type Attribute string

const (
	DARK = Attribute("闇")
	LIGHT = Attribute("光")
	EARTH = Attribute("地")
	WATER = Attribute("水")
	FIRE = Attribute("炎")
	WIND = Attribute("風")
)

type Attributes []Attribute

var ATTRIBUTES = Attributes{DARK, LIGHT, EARTH, WATER, FIRE, WIND}

type Type string

const (
	DRAGON = "ドラゴン"
	SPELLCASTER = "魔法使い"
	ZOMBLE = "ゾンビ"
	WARRIOR = "戦士"
	BEAST_WARRIOR = "獣戦士"
	BEAST = "獣"
	WINGED_BEAST = "鳥獣"
	FIEND = "悪魔"
	FAIRY = "天使"
	INSECT = "昆虫"
	DINOSAUR = "恐竜"
	REPTILE = "爬虫類"
	FISH = "魚"
	SEA_SERPENT = "海竜"
	MACHINE = "機械"
	THUNDER = "雷"
	AQUA = "水"
	PYRO = "炎"
	ROCK = "岩石"
	PLANT = "植物"
	PSYCHIC = "サイキック"
	WYRM = "幻竜"
	CYBERSE = "サイバース"
	DIVINE_BEAST = "幻神獣"
)

type Types []Type

var TYPES = Types{
	DRAGON, SPELLCASTER, ZOMBLE, WARRIOR, BEAST_WARRIOR,
	BEAST, WINGED_BEAST, FIEND, FAIRY, INSECT,
	DINOSAUR, REPTILE, FISH, SEA_SERPENT, MACHINE,
	THUNDER, AQUA, PYRO, ROCK, PLANT,
	PSYCHIC, WYRM, CYBERSE, DIVINE_BEAST,
}

type Card struct {
	Name CardName
	Level Level
	Atk int
	Def int

	Attribute Attribute
	Type Type

	IsNormalMonster bool
	IsEffectMonster bool
	IsSpiritMonster bool
	CanNormalSummon bool

    IsNormalSpell bool
    IsQuickPlaySpell bool
    IsContinuousSpell bool

	IsNormalTrap bool
	IsContinuousTrap bool
	IsCounterTrap bool

	IsAttackPosition bool
	IsFaceUpDefensePosition bool
	IsFaceDownDefensePosition bool

	IsSetTurn bool
	IsNormalSummoned bool
}

var EMPTY_CARD = Card{}

func (card Card) Clone() Card {
	return card
}

func IsSpellSpeed2Card(card Card) bool {
	return card.IsQuickPlaySpell || IsTrapCard(card)
}

func IsEmptyCard(card Card) bool {
	return card == EMPTY_CARD
}

func IsNotEmptyCard(card Card) bool {
	return card != EMPTY_CARD
}

func IsMonsterCard(card Card) bool {
	return card.IsNormalMonster || card.IsEffectMonster
}

func IsLowLevelMonsterCard(card Card) bool {
	return slices.Contains(LOW_LEVELS, card.Level)
}

func IsLevel4MonsterCard(card Card) bool {
	return card.Level == 4
}

func IsMediumLevelMonsterCard(card Card) bool {
	return slices.Contains(MEDIUM_LEVELS, card.Level)
}

func IsHighLevelMonsterCard(card Card) bool {
	return card.Level > omathw.Max(MEDIUM_LEVELS...)
}

func IsSpiritMonsterCard(card Card) bool {
	return card.IsSpiritMonster
}

func IsSpellCard(card Card) bool {
	return card.IsNormalSpell || card.IsQuickPlaySpell || card.IsContinuousSpell
}

func IsTrapCard(card Card) bool {
	return card.IsNormalTrap || card.IsContinuousTrap
}

func IsToonCard(card Card) bool {
	return strings.Contains(string(card.Name), string(TOON))
}

type Cards []Card

var OLD_LIBRARY_EXODIA_DECK = func() Cards {
	result, err := NewCards(
		"封印されしエクゾディア",
		"封印されし者の左腕",
		"封印されし者の右腕",
		"封印されし者の左足",
		"封印されし者の右足",
		"王立魔法図書館",
		"王立魔法図書館",
		"王立魔法図書館",
		"召喚僧サモンプリースト",
		"召喚僧サモンプリースト",
		"サンダー・ドラゴン",
		"サンダー・ドラゴン",
		"サンダー・ドラゴン",

		"一時休戦",
		"成金ゴブリン",
		"成金ゴブリン",
		"成金ゴブリン",
		"トゥーンのもくじ",
		"トゥーンのもくじ",
		"トゥーンのもくじ",
		"トゥーン・ワールド",
		"精神統一",
		"精神統一",
		"精神統一",
		"手札断殺",
		"手札断殺",
		"手札断殺",
		"打ち出の小槌",
		"打ち出の小槌",
		"打ち出の小槌",
		"闇の誘惑",
		"二重召喚",
		"魔法石の採掘",
		"闇の量産工場",

		"強欲な瓶",
		"強欲な瓶",
		"強欲な瓶",
		"八汰烏の骸",
		"八汰烏の骸",
		"八汰烏の骸",
	)
	if err != nil {
		panic(err)
	}
	return result
}()

func NewCards(names ...CardName) (Cards, error) {
	result := make(Cards, len(names))
	for i, name := range names {
		var card *Card
		if name == "" {
			copyCard := EMPTY_CARD.Clone()
			card = &copyCard
		} else {
			var ok bool
			card, ok = CARD_DATA_BASE[name]
			copyCard := card.Clone()
			card = &copyCard
			if !ok {
				msg := fmt.Sprintf("データベースに存在しないカード名が入力された。入力されたカード名 = %v", name)
				return Cards{}, fmt.Errorf(msg)
			}
		}
		result[i] = *card
	}
	return result, nil
}

func NewCardsWithPanic(names ...CardName) Cards {
	cards, err := NewCards(names...)
	if err != nil {
		panic(err)
	}
	return cards
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