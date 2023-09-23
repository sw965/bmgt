package bmgt

import (
	"github.com/sw965/omw"
)

var (
	BMGT_PATH    = omw.SW965_PATH + "bmgt/"
	DATA_PATH    = BMGT_PATH + "data/"
	MONSTER_PATH = DATA_PATH + "monster/"
	SPELL_PATH   = DATA_PATH + "spell/"
	TRAP_PATH    = DATA_PATH + "trap/"
)

const TEMPLATE_JSON_NAME = "テンプレート.json"

func IsTemplateJsonName(name string) bool {
	return name == TEMPLATE_JSON_NAME
}

func IsNotTemplateJsonName(name string) bool {
	return !IsTemplateJsonName(name)
}
