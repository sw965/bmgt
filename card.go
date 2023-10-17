package bmgt

import (
	"github.com/sw965/omw/fn"
	omwmaps "github.com/sw965/omw/maps"
)

type CardName int

const (
	EXODIA_THE_FORBIDDEN_ONE CardName = iota
	RIGHT_LEG_OF_THE_FORBIDDEN_ONE
	RIGHT_ARM_OF_THE_FORBIDDEN_ONE
	LEFT_LEG_OF_THE_FORBIDDEN_ONE
	LEFT_ARM_OF_THE_FORBIDDEN_ONE
)

func CardNameToString(cardName CardName) string {
	switch cardName {
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

func IsEmptyCard(card Card) bool {
	return card == Card{}
}

type Cards []Card

type BattlePosition int

const (
	ATK_POSITION BattlePosition = iota
	FACE_UP_DEF_POSITION
	FACE_DOWN_DEF_POSITION
)

type BattlePositions []BattlePosition

var BATTLE_POSITIONS = BattlePositions{ATK_POSITION, FACE_UP_DEF_POSITION, FACE_DOWN_DEF_POSITION}
var NORMAL_SUMMON_BATTLE_POSITIONS = BattlePositions{ATK_POSITION, FACE_DOWN_DEF_POSITION}