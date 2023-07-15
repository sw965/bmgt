package bmgt

import (
	"fmt"
	"github.com/sw965/omw/json"
	osmw "github.com/sw965/omw/os"
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
	DARK  = Attribute("闇")
	LIGHT = Attribute("光")
	EARTH = Attribute("地")
	WATER = Attribute("水")
	FIRE  = Attribute("炎")
	WIND  = Attribute("風")
)

type Attributes []Attribute

var ATTRIBUTES = Attributes{DARK, LIGHT, EARTH, WATER, FIRE, WIND}

type Type string

const (
	DRAGON        = "ドラゴン"
	SPELLCASTER   = "魔法使い"
	ZOMBLE        = "ゾンビ"
	WARRIOR       = "戦士"
	BEAST_WARRIOR = "獣戦士"
	BEAST         = "獣"
	WINGED_BEAST  = "鳥獣"
	FIEND         = "悪魔"
	FAIRY         = "天使"
	INSECT        = "昆虫"
	DINOSAUR      = "恐竜"
	REPTILE       = "爬虫類"
	FISH          = "魚"
	SEA_SERPENT   = "海竜"
	MACHINE       = "機械"
	THUNDER       = "雷"
	AQUA          = "水"
	PYRO          = "炎"
	ROCK          = "岩石"
	PLANT         = "植物"
	PSYCHIC       = "サイキック"
	WYRM          = "幻竜"
	CYBERSE       = "サイバース"
	DIVINE_BEAST  = "幻神獣"
)

type Types []Type

var TYPES = Types{
	DRAGON, SPELLCASTER, ZOMBLE, WARRIOR, BEAST_WARRIOR,
	BEAST, WINGED_BEAST, FIEND, FAIRY, INSECT,
	DINOSAUR, REPTILE, FISH, SEA_SERPENT, MACHINE,
	THUNDER, AQUA, PYRO, ROCK, PLANT,
	PSYCHIC, WYRM, CYBERSE, DIVINE_BEAST,
}

type EffectType string

const (
	IGNITION_EFFECT   = "起動効果"
	TRIGGER_EFFECT    = "誘発効果"
	CONTINUOUS_EFFECT = "永続効果"
)

type EffectTypes []EffectType

type CardBaseData struct {
	Level Level
	Atk   int
	Def   int

	Attribute Attribute
	Type      Type

	IsNormalMonster bool
	IsEffectMonster bool
	IsSpiritMonster bool
	CanNormalSummon bool

	IsNormalSpell     bool
	IsQuickPlaySpell  bool
	IsContinuousSpell bool

	IsNormalTrap     bool
	IsContinuousTrap bool
	IsCounterTrap    bool

	MaxSpellCounter int
	EffectTypes     EffectTypes
}

func (data *CardBaseData) IsMonster() bool {
	return data.IsNormalMonster || data.IsEffectMonster || data.IsSpiritMonster
}

func (data *CardBaseData) IsSpell() bool {
	return data.IsNormalSpell || data.IsQuickPlaySpell || data.IsContinuousSpell
}

func (data *CardBaseData) IsTrap() bool {
	return data.IsNormalTrap || data.IsContinuousTrap || data.IsCounterTrap
}

type CardDatabase map[CardName]*CardBaseData

var CARD_DATA_BASE = func() CardDatabase {
	result := CardDatabase{}

	add := func(path string) {
		dirEntries, err := osmw.NewDirEntries(path)
		if err != nil {
			panic(err)
		}

		dirNames := dirEntries.Names()
		for _, dirName := range dirNames {
			if dirName == "テンプレート.json" {
				continue
			}
			cardName := CardName(strings.TrimRight(dirName, ".json"))
			data, err := json.Load[CardBaseData](path + dirName)
			if err != nil {
				panic(err)
			}
			result[cardName] = &data
		}
	}

	add(MONSTER_PATH)
	add(SPELL_PATH)
	add(TRAP_PATH)

	return result
}()

func init() {
	for name, data := range CARD_DATA_BASE {
		isMonster := data.IsMonster()
		isSpell := data.IsSpell()
		isTrap := data.IsTrap()

		if data.IsSpiritMonster && !data.IsEffectMonster {
			fmt.Println(name, "スピリットモンスターであるのに、効果モンスターではない (スピリットモンスターは、同時に効果モンスターでなければならない)")
		}

		if !isMonster && !isSpell && !isTrap {
			fmt.Println(name, "モンスター/魔法/罠 のどれでもない")
		}

		if !slices.Contains(ATTRIBUTES, data.Attribute) && isMonster {
			fmt.Println(name, data.Attribute)
		}

		if !slices.Contains(TYPES, data.Type) && isMonster {
			fmt.Println(name, data.Type)
		}
	}
}

const (
	SAME_CARD_NAME_LIMIT = 3
)
