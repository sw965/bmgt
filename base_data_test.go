package bmgt_test

import (
	"fmt"
	"github.com/sw965/bmgt"
	"testing"
)

func TestPrintCARD_DATA_BASE(t *testing.T) {
	for cardName, card := range bmgt.CARD_DATA_BASE {
		fmt.Println(cardName, card)
	}
}
