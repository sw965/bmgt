package bmgt

import (
	osmw "github.com/sw965/omw/os"
	"github.com/sw965/omw/json"
	"strings"
	"fmt"
	"golang.org/x/exp/slices"
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
	add(TRAP_PATH)
	
	return result
}()

func init() {
	for name, card := range CARD_DATA_BASE {
		if name != card.Name {
			fmt.Println(name, card.Name)
		}

		isMonster := IsMonsterCard(*card)
		isSpell := IsSpellCard(*card)
		isTrap := IsTrapCard(*card)

		if !isMonster && !isSpell && !isTrap {
			fmt.Println(name, "モンスター/魔法/罠 のどれでもない")
		}

		if !slices.Contains(ATTRIBUTES, card.Attribute) && isMonster {
			fmt.Println(name, card.Attribute)
		}

		if !slices.Contains(TYPES, card.Type) && isMonster {
			fmt.Println(name, card.Type)
		}
	}
}