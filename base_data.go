package bmgt

import (
	"fmt"
	"github.com/sw965/omw/json"
	osmw "github.com/sw965/omw/os"
	"strings"
)

type Level int
type Levels []Level

var LOW_LEVELS = Levels{1, 2, 3, 4}
var MEDIUM_LEVELS = Levels{5, 6}
const HIGH_LEVEL = 7

type Attribute int

const (
	DARK Attribute = iota
	LIGHT
	EARTH
	WATER
	FIRE
	WIND
)

var STRING_TO_ATTRIBUTE = map[string]Attribute{
	"闇":DARK,
	"光":LIGHT,
	"地":EARTH,
	"水":WATER,
	"炎":FIRE,
	"風":WIND,
}

type Attributes []Attribute

var ATTRIBUTES = Attributes{DARK, LIGHT, EARTH, WATER, FIRE, WIND}

type Type int

const (
	DRAGON Type = iota
	SPELLCASTER
	ZOMBLE
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
	WYRM
	CYBERSE
	DIVINE_BEAST
)

var STRING_TO_TYPE = map[string]Type{
	"ドラゴン":DRAGON,
	"魔法使い":SPELLCASTER,
	"ゾンビ":ZOMBLE,
	"戦士":WARRIOR,
	"獣戦士":BEAST_WARRIOR,
	"獣":BEAST,
	"鳥獣":WINGED_BEAST,
	"悪魔":FIEND,
	"天使":FAIRY,
	"恐竜":DINOSAUR,
	"爬虫類":REPTILE,
	"魚":FISH,
	"海竜":SEA_SERPENT,
	"機械":MACHINE,
	"雷":THUNDER,
	"水":AQUA,
	"炎":PYRO,
	"岩石":ROCK,
	"植物":PLANT,
	"サイキック":PSYCHIC,
	"幻神獣":DIVINE_BEAST,
}

type Types []Type

type Category int

const (
	NORMAL_MONSTER Category = iota
	EFFECT_MONSTER
	SPIRIT_MONSTER

	NORMAL_SPELL
	QUICK_PLAY_SPELL
	CONTINUOUS_SPELL

	NORMAL_TRAP
	CONTINUOUS_TRAP
	COUNTER_TRAP
)

func (c Category) IsMonster() bool {
	return c == NORMAL_MONSTER || c == EFFECT_MONSTER || c == SPIRIT_MONSTER
}

func (c Category) IsSpell() bool {
	return c == NORMAL_SPELL || c == QUICK_PLAY_SPELL || c == CONTINUOUS_SPELL
}

func (c Category) IsTrap() bool {
	return c == NORMAL_TRAP || c == CONTINUOUS_TRAP || c == COUNTER_TRAP
}

var STRING_TO_CATEGORY = map[string]Category{
	"通常モンスター":NORMAL_MONSTER,
	"効果モンスター":EFFECT_MONSTER,
	"スピリットモンスター":SPIRIT_MONSTER,

	"通常魔法":NORMAL_SPELL,
	"速攻魔法":QUICK_PLAY_SPELL,
	"永続魔法":CONTINUOUS_SPELL,

	"通常罠":NORMAL_TRAP,
	"永続罠":CONTINUOUS_TRAP,
	"カウンター罠":COUNTER_TRAP,
}

type cardBaseData struct {
	Level Level
	Atk int
	Def int

	Attribute string
	Type string
	Category string

	EffectNum int
	MaxSpellCounter int
}

type CardBaseData struct {
	Level Level
	Atk   int
	Def   int

	Attribute Attribute
	Type      Type
	Category Category

	EffectNum int
	MaxSpellCounter int
}

func LoadCardDataBase(path string) (CardBaseData, error) {
	data, err := json.Load[cardBaseData](path)
	if err != nil {
		return CardBaseData{}, err
	}

	category, ok := STRING_TO_CATEGORY[data.Category]
	if !ok {
		msg := fmt.Sprintf("不適な分類 %v", path)
		return CardBaseData{}, fmt.Errorf(msg)
	}

	isMonster := category.IsMonster()

	attribute, ok := STRING_TO_ATTRIBUTE[data.Attribute]
	if !ok && isMonster {
		return CardBaseData{}, fmt.Errorf("不適な属性")
	}

	t, ok := STRING_TO_TYPE[data.Type]
	if !ok && isMonster {
		return CardBaseData{}, fmt.Errorf("不適な種族")
	}

	y := CardBaseData{
		Level:data.Level,
		Atk:data.Atk,
		Def:data.Def,
		Attribute:attribute,
		Type:t,
		Category:category,
		EffectNum:data.EffectNum,
		MaxSpellCounter:data.MaxSpellCounter,
	}
	return y, nil
}

func (d CardBaseData) ToCard() Card {
	return Card{
		Category:d.Category,
		Attribute:d.Attribute,
		Level:d.Level,
		ThisTurnEffectActivationCounts: make([]int, d.EffectNum),
	}
}

type CardDatabase map[CardName]*CardBaseData

var CARD_DATA_BASE = func() CardDatabase {
	y := CardDatabase{}

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
			name := strings.TrimRight(dirName, ".json")
			data, err := LoadCardDataBase(path + dirName)
			if err != nil {
				panic(err)
			}
			y[STRING_TO_CARD_NAME[name]] = &data
		}
	}

	add(MONSTER_PATH)
	add(SPELL_PATH)
	add(TRAP_PATH)

	return y
}()

func init() {
	for name, data := range CARD_DATA_BASE {
		if data.Category != NORMAL_MONSTER && data.EffectNum == 0 {
			msg := fmt.Sprintf("通常モンスターではないのに、効果の数が0になっている。(%v)", CARD_NAME_TO_STRING[name])
			fmt.Println(msg)
		} 
	}
}