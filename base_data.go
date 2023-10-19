package bmgt

import (
	omwjson "github.com/sw965/omw/json"
	omwos "github.com/sw965/omw/os"
	"strings"
)

type cardBaseData struct {
	Attribute string
	Level Level
	Type string
	Atk int
	Def int
}

type CardBaseData struct {
	Attribute Attribute
	Level Level
	Type Type
	Atk int
	Def int
}

func LoadCardBaseData(path string) CardBaseData {
	old, err := omwjson.Load[cardBaseData](path)
	if err != nil {
		panic(err)
	}
	new := CardBaseData{}
	new.Attribute = STRING_TO_ATTRIBUTE[old.Attribute]
	new.Level = old.Level
	new.Type = STRING_TO_TYPE[old.Type]
	new.Atk = old.Atk
	new.Def = old.Def
	return new
}

type CardDatabase map[CardName]*CardBaseData

var CARD_DATABASE = func() CardDatabase {
	y := CardDatabase{}
	entries, err := omwos.NewDirEntries(MONSTER_JSON_PATH)
	if err != nil {
		panic(err)
	}
	for _, dirName := range entries.Names() {
		path := MONSTER_JSON_PATH + dirName
		baseData := LoadCardBaseData(path)
		cardName := STRING_TO_CARD_NAME[strings.Replace(dirName, ".json", "", 1)]
		y[cardName] = &baseData
	}
	return y
}()