package bmgt

import (
	omwmaps "github.com/sw965/omw/maps"
	omwjson "github.com/sw965/omw/json"
	omwstrings "github.com/sw965/omw/strings"
	omwos "github.com/sw965/omw/os"
)

type Category int

const (
	MONSTER Category = iota
	NORMAL_SPELL
	NORMAL_TRAP
)

var CATEGORY_TO_STRING = map[Category]string{
	MONSTER:"モンスター", NORMAL_SPELL:"通常魔法", NORMAL_TRAP:"通常罠",
}

var STRING_TO_CATEGORY = omwmaps.Reverse[map[string]Category, map[Category]string](CATEGORY_TO_STRING)

type cardBaseData struct {
	Attribute string
	Level Level
	Type string
	Atk int
	Def int	
	Category string
}

func (d *cardBaseData) ToCardBaseData() CardBaseData {
	y := CardBaseData{}
	y.Attribute = STRING_TO_ATTRIBUTE[d.Attribute]
	y.Level = d.Level
	y.Type = STRING_TO_TYPE[d.Type]
	y.Atk = d.Atk
	y.Def = d.Def
	y.Category = STRING_TO_CATEGORY[d.Category]
	return y
}

type CardBaseData struct {
	Attribute Attribute
	Level Level
	Type Type
	Atk int
	Def int
	Category Category
}

func NewCardBaseData(fileName string) CardBaseData {
	d, err := omwjson.Load[cardBaseData](JSON_PATH + fileName)
	if err != nil {
		panic(err)
	}
	return d.ToCardBaseData()
}

type CardDatabase map[CardName]CardBaseData

var CARD_DATA_BASE = func() CardDatabase {
	entries, err := omwos.NewDirEntries(JSON_PATH)
	if err != nil {
		panic(err)
	}
	y := CardDatabase{}
	for _, fileName := range entries.Names() {
		cardName := STRING_TO_CARD_NAME[omwstrings.Replace(omwjson.EXTENSION, "", 1)(fileName)]
		y[cardName] = NewCardBaseData(fileName)
	}
	return y
}()