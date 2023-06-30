package bmgt_test

import (
	"testing"
	"github.com/sw965/bmgt"
	"fmt"
)

func TestPrintCARD_DATA_BASE(t *testing.T) {
	for cardName, card := range bmgt.CARD_DATA_BASE {
		fmt.Println(cardName, card)
	}
}