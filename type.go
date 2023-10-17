package bmgt

import (
	"github.com/sw965/omw/fn"
	omwmaps "github.com/sw965/omw/maps"
)

type Type int

const (
	DRAGON Type = iota
	SPELLCASTER
	ZOMBIE
	WARRIOR
	BEAST_WARRIOR
	BEAST
	WINGED_BEAST
	FIEND
	FAIRY
	INSECT
	DINOSAUR
	REPTILE
	FISH
	SEA_SERPENT
	MACHINE
	THUNDER
	AQUA
	PYRO
	ROCK
	PLANT
	PSYCHIC
	CYBERSE
)

type Types []Type

var TYPES = Types{
	DRAGON,
	SPELLCASTER,
	ZOMBIE,
	WARRIOR,
	BEAST_WARRIOR,
	BEAST,
	WINGED_BEAST,
	FIEND,
	FAIRY,
	INSECT,
	DINOSAUR,
	REPTILE,
	FISH,
	SEA_SERPENT,
	MACHINE,
	THUNDER,
	AQUA,
	PYRO,
	ROCK,
	PLANT,
	PSYCHIC,
	CYBERSE,
}

func TypeToString(t Type) string {
	switch t {
		case DRAGON:
			return "ドラゴン"
		case SPELLCASTER:
			return "魔法使い"
		case ZOMBIE:
			return "アンデット"
		case WARRIOR:
			return "戦士"
		case BEAST_WARRIOR:
			return "獣戦士"
		case BEAST:
			return "獣"
		case WINGED_BEAST:
			return "鳥獣"
		case FIEND:
			return "悪魔"
		case FAIRY:
			return "天使"
		case INSECT:
			return "昆虫"
		case DINOSAUR:
			return "恐竜"
		case REPTILE:
			return "爬虫類"
		case FISH:
			return "魚"
		case SEA_SERPENT:
			return "海竜"
		case MACHINE:
			return "機械"
		case THUNDER:
			return "雷"
		case AQUA:
			return "水"
		case PYRO:
			return "炎"
		case ROCK:
			return "岩石"
		case PLANT:
			return "植物"
		case PSYCHIC:
			return "サイキック"
		case CYBERSE:
			return "サイバース"
		default:
			return ""
	}
}

var TYPE_TO_STRING = fn.Memo[map[Type]string](TYPES, TypeToString)
var STRING_TO_TYPE = omwmaps.Reverse[map[string]Type](TYPE_TO_STRING)