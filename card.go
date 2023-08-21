package bmgt

import (
	"fmt"
	omws "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
)

type CardName int

const (
	NO_NAME CardName = iota
	ONE_DAY_OF_PEACE
	MAGICAL_MALLET
	ROYAL_MAGICAL_LIBRARY
	SOLEMN_JUDGMENT
	JAR_OF_GREED
	POT_OF_GREED
	SUMMONER_MONK
	THUNDER_DRAGON
	GATHER_YOUR_MIND
	HAND_DESTRUCTION
	DOUBLE_SUMMON
	TOON_TABLE_OF_CONTENTS
	TOON_WORLD
	UPSTART_GOBLIN
	EXODIA_THE_FORBIDDEN_ONE
	LEFT_LEG_OF_THE_FORBIDDEN_ONE
	LEFT_ARM_OF_THE_FORBIDDEN_ONE
	RIGHT_LEG_OF_THE_FORBIDDEN_ONE
	RIGHT_ARM_OF_THE_FORBIDDEN_ONE
	MAGICAL_STONE_EXCAVATION
	ALLURE_OF_DARKNESS
	DARK_FACTORY_OF_MASS_PRODUCTION
	LEGACY_OF_YATA_GARASU
)

var STRING_TO_CARD_NAME = map[string]CardName{
	"一時休戦":ONE_DAY_OF_PEACE,
	"打ち出の小槌":MAGICAL_MALLET,
	"王立魔法図書館":ROYAL_MAGICAL_LIBRARY,
	"強欲な瓶":JAR_OF_GREED,
	"召喚僧サモンプリースト":SUMMONER_MONK,
	"サンダー・ドラゴン":THUNDER_DRAGON,
	"精神統一":GATHER_YOUR_MIND,
	"手札断殺":HAND_DESTRUCTION,
	"二重召喚":DOUBLE_SUMMON,
	"トゥーンのもくじ":TOON_TABLE_OF_CONTENTS,
	"トゥーン・ワールド":TOON_WORLD,
	"成金ゴブリン":UPSTART_GOBLIN,
	"封印されしエクゾディア":EXODIA_THE_FORBIDDEN_ONE,
	"封印されし者の左足":LEFT_LEG_OF_THE_FORBIDDEN_ONE,
	"封印されし者の左腕":LEFT_ARM_OF_THE_FORBIDDEN_ONE,
	"封印されし者の右足":RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
	"封印されし者の右腕":RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
	"魔法石の採掘":MAGICAL_STONE_EXCAVATION,
	"闇の誘惑":ALLURE_OF_DARKNESS,
	"闇の量産工場":DARK_FACTORY_OF_MASS_PRODUCTION,
	"八汰烏の骸":LEGACY_OF_YATA_GARASU,
}

var CARD_NAME_TO_STRING = func() map[CardName]string {
	y := map[CardName]string{}
	for k, v := range STRING_TO_CARD_NAME {
		y[v] = k
	}
	return y
}()

type CardNames []CardName

var NORMAL_MONSTER_NAMES = func() CardNames {
	y := make(CardNames, 0, 128)
	for name, data := range CARD_DATA_BASE {
		if data.Category == NORMAL_MONSTER {
			y = append(y, name)
		}
	}
	return y
}()

var SPELL_CARD_NAMES = func() CardNames {
	y := make(CardNames, 0, 128)
	for name, data := range CARD_DATA_BASE {
		if data.Category.IsSpell() {
			y = append(y, name)			
		}
	}
	return y
}()

var EXODIA_PART_NAMES = CardNames{
	EXODIA_THE_FORBIDDEN_ONE,
	LEFT_ARM_OF_THE_FORBIDDEN_ONE,
	RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
	LEFT_LEG_OF_THE_FORBIDDEN_ONE,
	RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
}

var TOON_CARD_NAMES = CardNames{
	TOON_TABLE_OF_CONTENTS,
	TOON_WORLD,
}

func (names CardNames) ToStrings() ([]string, error) {
	f := func(name CardName) (string, error) {
		if y, ok := CARD_NAME_TO_STRING[name]; !ok {
			return "", fmt.Errorf("不適な名前")
		} else {
			return y, nil
		}
	}
	return fn.MapError[[]string](names, f)
}

type BattlePosition int

const (
	ATTACK_POSITION BattlePosition = iota
	FACE_UP_DEFENSE_POSITION
	FACE_DOWN_DEFENSE_POSITION
)

type CardID int
type CardIDs []CardID

type Card struct {
	Name           CardName
	Category Category

	Level Level
	BattlePosition BattlePosition
	Attribute Attribute

	IsSet     bool
	IsSetTurn bool

	ThisTurnEffectActivationCounts []int
	SelectEffectNumber             int
	SpellCounter                   int

	ID        CardID
}

var EMPTY_CARD = Card{}

func IsDarkMonster(card Card) bool {
	return card.Attribute == DARK
}

func (card Card) Clone() Card {
	counts := make([]int, len(card.ThisTurnEffectActivationCounts))
	for i, c := range card.ThisTurnEffectActivationCounts {
		counts[i] = c
	}
	card.ThisTurnEffectActivationCounts = counts
	return card
}

type Cards []Card

var OLD_LIBRARY_EXODIA_DECK = func() Cards {
	y, err := NewCards(
		EXODIA_THE_FORBIDDEN_ONE,
		LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		RIGHT_LEG_OF_THE_FORBIDDEN_ONE,

		ROYAL_MAGICAL_LIBRARY,
		ROYAL_MAGICAL_LIBRARY,
		ROYAL_MAGICAL_LIBRARY,

		SUMMONER_MONK,
		SUMMONER_MONK,

		THUNDER_DRAGON,
		THUNDER_DRAGON,
		THUNDER_DRAGON,

		ONE_DAY_OF_PEACE,

		UPSTART_GOBLIN,
		UPSTART_GOBLIN,
		UPSTART_GOBLIN,

		TOON_TABLE_OF_CONTENTS,
		TOON_TABLE_OF_CONTENTS,
		TOON_TABLE_OF_CONTENTS,
		TOON_WORLD,

		GATHER_YOUR_MIND,
		GATHER_YOUR_MIND,
		GATHER_YOUR_MIND,

		HAND_DESTRUCTION,
		HAND_DESTRUCTION,
		HAND_DESTRUCTION,

		MAGICAL_MALLET,
		MAGICAL_MALLET,
		MAGICAL_MALLET,

		ALLURE_OF_DARKNESS,
		DOUBLE_SUMMON,
		MAGICAL_STONE_EXCAVATION,
		DARK_FACTORY_OF_MASS_PRODUCTION,

		JAR_OF_GREED,
		JAR_OF_GREED,
		JAR_OF_GREED,

		LEGACY_OF_YATA_GARASU,
		LEGACY_OF_YATA_GARASU,
		LEGACY_OF_YATA_GARASU,
	)
	if err != nil {
		panic(err)
	}
	return y
}()

func NewCards(names ...CardName) (Cards, error) {
	result := make(Cards, len(names))
	for i, name := range names {
		var card Card
		if name == NO_NAME {
			card = EMPTY_CARD.Clone()
		} else {
			data, ok := CARD_DATA_BASE[name]
			if !ok {
				msg := fmt.Sprintf("データベースに存在しないカード名が入力された %v", name)
				return Cards{}, fmt.Errorf(msg)
			}
			card = Card{Name: name, Category:data.Category, Attribute:data.Attribute, Level:data.Level, ThisTurnEffectActivationCounts: make([]int, data.EffectNum)}
		}
		result[i] = card
	}
	return result, nil
}

func (cards Cards) Names() CardNames {
	result := make(CardNames, len(cards))
	for i, card := range cards {
		result[i] = card.Name
 	}
	return result
}

func (cards Cards) IDs() CardIDs {
	result := make(CardIDs, len(cards))
	for i, card := range cards {
		result[i] = card.ID
	}
	return result
}

func (cards Cards) IDSorted() Cards {
	return omws.SortedFunc(cards, func(c1, c2 Card) bool { return c1.ID < c2.ID })
}