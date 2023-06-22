package bmgt

import (
	osmw "github.com/sw965/omw/os"
	"github.com/sw965/omw/json"
	"strings"
)

type CardDatabase map[CardName]*Card

var CARD_DATA_BASE = func() CardDatabase {
	result := CardDatabase{}

	add := func (path string) {
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
			card, err := json.Load[Card](path + dirName)
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