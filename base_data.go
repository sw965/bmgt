package bmgt

var STRING_CARD_NAMES = func() CardNames {
	entries := omwos.NewDirEntries(MONSTER_DATA_PATH)
	monsterNames := fn.Map(entries.Names(),  omwstrings.Replace(".json", "", 1))
	y := make(CardNames, 0, len(monsterNames))
	y = append(y, monsterNames...)
	return y
}

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
	Atk int
}

func NewCardBaseData(path string) CardBaseData {
	entries := os.NewDirEntries(path)
	for _, fileName := range entries.Names() {
		jsonData :=
	}
}

type CardDatabase map[CardName]CardBaseData

func NewCardDatabase() CardDatabase {
	y := CardDatabase{}
	string
}