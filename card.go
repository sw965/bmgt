package bmgt

import (
	"github.com/sw965/omw/fn"
	omwmaps "github.com/sw965/omw/maps"
)

type CardName int

const (
	NO_NAME CardName = iota
	LUSTER_DRAGON
	GEMINI_ELF
	VORSE_RAIDER
	EXODIA_THE_FORBIDDEN_ONE
	RIGHT_LEG_OF_THE_FORBIDDEN_ONE
	RIGHT_ARM_OF_THE_FORBIDDEN_ONE
	LEFT_LEG_OF_THE_FORBIDDEN_ONE
	LEFT_ARM_OF_THE_FORBIDDEN_ONE
)

func CardNameToString(cardName CardName) string {
	switch cardName {
		case NO_NAME:
			return ""
		case LUSTER_DRAGON:
			return "サファイアドラゴン"
		case GEMINI_ELF:
			return "ヂェミナイ・エルフ"
		case VORSE_RAIDER:
			return "ブラッド・ヴォルス"
		case EXODIA_THE_FORBIDDEN_ONE:
			return "封印されしエクゾディア"
		case RIGHT_LEG_OF_THE_FORBIDDEN_ONE:
			return "封印されし者の右足"
		case RIGHT_ARM_OF_THE_FORBIDDEN_ONE:
			return "封印されし者の右腕"
		case LEFT_LEG_OF_THE_FORBIDDEN_ONE:
			return "封印されし者の左足"
		case LEFT_ARM_OF_THE_FORBIDDEN_ONE:
			return "封印されし者の左腕"
		default:
			return ""
	}
}

var CARD_NAME_TO_STRING = fn.Memo[map[CardName]string](CARD_NAMES, CardNameToString)
var STRING_TO_CARD_NAME = omwmaps.Reverse[map[string]CardName](CARD_NAME_TO_STRING)

type CardNames []CardName

var CARD_NAMES = CardNames{
	NO_NAME,
	LUSTER_DRAGON,
	GEMINI_ELF,
	VORSE_RAIDER,
	EXODIA_THE_FORBIDDEN_ONE,
	RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
	RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
	LEFT_LEG_OF_THE_FORBIDDEN_ONE,
	LEFT_ARM_OF_THE_FORBIDDEN_ONE,
}

var EXODIA_PARTS_NAMES = CardNames{
	EXODIA_THE_FORBIDDEN_ONE,
	RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
	RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
	LEFT_LEG_OF_THE_FORBIDDEN_ONE,
	LEFT_ARM_OF_THE_FORBIDDEN_ONE,
}

func CardNamesToStrings(names CardNames) []string {
	return fn.Map[[]string](names, CardNameToString)
}

type Level int
type Levels []Level

var LOW_LEVELS = Levels{1, 2, 3, 4}

type Card struct {
	Name CardName
	Attribute Attribute
	Level Level
	Type Type
	Atk int
	Def int
	BattlePosition BattlePosition
	IsBattlePositionChangeable bool
	IsAttackDeclared bool
}

func NewCard(name CardName) Card {
	data := CARD_DATABASE[name]
	return Card{
		Name:name,
		Attribute:data.Attribute,
		Level:data.Level,
		Type:data.Type,
		Atk:data.Atk,
		Def:data.Def,
	}
}

func GetNameOfCard(card Card) CardName {
	return card.Name
}

func CloneCard(card Card) Card {
	return card
}

func IsEmptyCard(card Card) bool {
	return card.Name == NO_NAME
}

type Cards []Card

func NewCards(names ...CardName) Cards {
	return fn.Map[Cards](names, NewCard)
}

func NamesOfCards(cards Cards) CardNames {
	return fn.Map[CardNames](cards, GetNameOfCard)
}

func CloneCards(cards Cards) Cards {
	return fn.Map[Cards](cards, CloneCard)
}

type BattlePosition int

const (
	ATK_BATTLE_POSITION BattlePosition = iota
	FACE_UP_DEF_BATTLE_POSITION
	FACE_DOWN_DEF_BATTLE_POSITION
)

func BattlePositionToString(pos BattlePosition) string {
	switch pos {
		case ATK_BATTLE_POSITION:
			return "攻撃表示"
		case FACE_UP_DEF_BATTLE_POSITION:
			return "表側守備表示"
		case FACE_DOWN_DEF_BATTLE_POSITION:
			return "裏側守備表示"
		default:
			return ""
	}
}

type BattlePositions []BattlePosition

var BATTLE_POSITIONS = BattlePositions{ATK_BATTLE_POSITION, FACE_UP_DEF_BATTLE_POSITION, FACE_DOWN_DEF_BATTLE_POSITION}