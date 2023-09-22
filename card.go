package bmgt

import (
	omwos "github.com/sw965/omw/os"
	omwstrings "github.com/sw965/omw/strings"
	"github.com/sw965/omw/fn"
	"fmt"
	"golang.org/x/exp/slices"
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
	MACRO_COSMOS
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

func StringToCardName(s string) CardName {
	return STRING_TO_CARD_NAME[s]
}

var CARD_NAME_TO_STRING = func() map[CardName]string {
	y := map[CardName]string{}
	for k, v := range STRING_TO_CARD_NAME {
		y[v] = k
	}
	return y
}()

type cardNameF struct{}
var CardNameF = cardNameF{}

func (f *cardNameF) IsMonster(name CardName) bool {
	return CARD_DATA_BASE[name].Category.IsMonster()
}

func (f *cardNameF) IsNormalMonster(name CardName) bool {
	return CARD_DATA_BASE[name].Category == NORMAL_MONSTER
}

type CardNames []CardName

var MONSTER_NAMES = func() CardNames {
	entries, err := omwos.NewDirEntries(MONSTER_PATH)
	if err != nil {
		panic(err)
	}
	dirs := fn.Filter(entries.Names(), IsNotTemplateJsonName)
	dirs = fn.Map[[]string](dirs, omwstrings.Replace(omwos.JSON_EXTENSION, "", 1))
	return fn.Map[CardNames](dirs, StringToCardName)
}()

var NORMAL_MONSTER_NAMES = func() CardNames {
	return fn.Filter(MONSTER_NAMES, CardNameF.IsNormalMonster)
}()

var SPELL_CARD_NAMES = func() CardNames {
	entries, err := omwos.NewDirEntries(SPELL_PATH)
	if err != nil {
		panic(err)
	}
	dirs := fn.Filter(entries.Names(), IsNotTemplateJsonName)
	dirs = fn.Map[[]string](dirs, omwstrings.Replace(omwos.JSON_EXTENSION, "", 1))
	return fn.Map[CardNames](dirs, StringToCardName)
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

type BattlePosition int

const (
	ATTACK_POSITION BattlePosition = iota
	FACE_UP_DEFENSE_POSITION
	FACE_DOWN_DEFENSE_POSITION
)

type BattlePositions []BattlePosition

var NORMAL_SUMMON_BATTLE_POSITIONS = BattlePositions{ATTACK_POSITION, FACE_DOWN_DEFENSE_POSITION}

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
	SpellCounter                   int

	ID        CardID
}

func NewCard(name CardName) (Card, error) {
	if name == NO_NAME {
		return Card{}, nil
	} else {
		data, ok := CARD_DATA_BASE[name]
		if !ok {
			msg := fmt.Sprintf("データベースに存在しないカード名が入力された %v", name)
			return Card{}, fmt.Errorf(msg)
		} else {
			card := data.ToCard()
			card.Name = name
			return card, nil
		}
	}
}

type cardF struct{}
var CardF = cardF{}

func (f *cardF) IsEmpty(card Card) bool {
	return card.Name == NO_NAME
}

func (f *cardF) IsNotEmpty(card Card) bool {
	return !f.IsEmpty(card)
}

func (f *cardF) GetName(card Card) CardName {
	return card.Name
}

func (f *cardF) SetID(id CardID, card Card) Card {
	card.ID = id
	return card
}

func (f *cardF) Clone(card Card) Card {
	card.ThisTurnEffectActivationCounts = slices.Clone(card.ThisTurnEffectActivationCounts)
	return card
}

func (f *cardF) IsMonster(card Card) bool {
	return slices.Contains(MONSTER_NAMES, card.Name)
}

func (f *cardF) IsLowLevelMonster(card Card) bool {
	return slices.Contains(LOW_LEVELS, card.Level)
}

func (f *cardF) IsDarkMonster(card Card) bool {
	return card.Attribute == DARK
}

func (f *cardF) CanNormalSummon(card Card) bool {
	return f.IsLowLevelMonster(card)
}

func (f *cardF) CanTributeSummonCost(card Card) bool {
	return f.IsMonster(card) && card.Name != SUMMONER_MONK
}

type Cards []Card

var OLD_LIBRARY_EXODIA_DECK = func() Cards {
	names := CardNames{
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
	}
	y, err := CardsF.New(names)
	if err != nil {
		panic(err)
	}
	return y
}()

type cardsF struct{}
var CardsF = cardsF{}

func (f *cardsF) New(names CardNames) (Cards, error) {
	return fn.MapError[Cards](names, NewCard)
}

func (f *cardsF) Names(cards Cards) CardNames {
	return fn.Map[CardNames](cards, CardF.GetName)
}

func (f *cardsF) Clone(cards Cards) Cards {
	return fn.Map[Cards](cards, CardF.Clone)
}