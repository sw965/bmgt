package bmgt

import (
	"os"
)

var (
	ROOT_PATH = os.Getenv("GOPATH") + "sw965/" + "bmgt/"
	JSON_PATH = ROOT_PATH + "json/"
)