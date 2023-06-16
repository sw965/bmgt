package bmgt

import (
	"github.com/sw965/omw"
	"strings"
)

type CardDatabase map[CardName]*Card

var CARD_DATA_BASE = func() CardDatabase {
	result := CardDatabase{}

	add := func (path string) {
		dirNames, err := omw.DirNames(path)
		if err != nil {
			panic(err)
		}
		for _, dirName := range dirNames {
			if dirName == "テンプレート.json" {
				continue
			}
			cardName := CardName(strings.TrimRight(dirName, ".json"))
			card, err := omw.LoadJson[Card](path + dirName)
			if err != nil {
				panic(err)
			}
			result[cardName] = &card
		}
	}

	add(MONSTER_PATH)
	add(SPELL_PATH)
	
	return result
}()