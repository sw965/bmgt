package bmgt_test

import (
	"testing"
	"github.com/sw965/bmgt"
	omwrand "github.com/sw965/omw/rand"
	"fmt"
	"github.com/sw965/crow/game/sequential"
)

func Test(t *testing.T) {
	deck1 := bmgt.NewCards(
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
		bmgt.LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
		bmgt.LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
		bmgt.LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
		bmgt.LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
		bmgt.LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
		bmgt.LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
		bmgt.LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
		bmgt.LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.RIGHT_LEG_OF_THE_FORBIDDEN_ONE,
	)
	r := omwrand.NewMt19937()
	duel := bmgt.NewDuel(deck1, deck1, r)
	game := sequential.Game()

	fmt.Println(duel)
}