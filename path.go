package bmgt

import (
	"os"
)

var (
	SW965_PATH = os.Getenv("GOPATH") + "sw965/"

	BMGT_PATH = SW965_PATH + "bmgt/"
	CARD_PATH = BMGT_PATH + "card/"
	MONSTER_PATH = CARD_PATH + "monster/"
	SPELL_PATH = CARD_PATH + "spell/"
	TRAP_PATH = CARD_PATH + "trap/"
)