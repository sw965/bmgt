package bmgt_test

import (
	"github.com/sw965/bmgt"
	"testing"
	"fmt"
)

func Test(t *testing.T) {
	hand, err := bmgt.CardsF.New(
		bmgt.SUMMONER_MONK,
		bmgt.JAR_OF_GREED,
		bmgt.THUNDER_DRAGON,
		bmgt.POT_OF_GREED,
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
	)
	if err != nil {
		panic(err)
	}

	monsterZone, err := bmgt.CardsF.New(
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.NO_NAME,
		bmgt.LEFT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.LEFT_LEG_OF_THE_FORBIDDEN_ONE,
		bmgt.NO_NAME,
	)
	if err != nil {
		panic(err)
	}

	state := bmgt.State{}
	state.P1.Hand = hand
	state.P1.MonsterZone = monsterZone
	result := bmgt.ActionsF.NewLegalNormalSummon(&state)
	for _, action := range result {
		fmt.Println(action)
	}
}