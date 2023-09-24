package bmgt_test

import (
	"github.com/sw965/bmgt"
	"testing"
	"fmt"
)

func TestLegalNormalSummonActions(t *testing.T) {
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
	result := bmgt.ActionsF.ToStrings(bmgt.ActionsF.NewLegalNormalSummon(&state))
	for _, action := range result {
		fmt.Println(action)
	}
}

func TestLegalFlipSummon(t *testing.T) {
	monsterZone, err := bmgt.CardsF.New(
		bmgt.THUNDER_DRAGON,
		bmgt.NO_NAME,
		bmgt.SUMMONER_MONK,
		bmgt.RIGHT_ARM_OF_THE_FORBIDDEN_ONE,
		bmgt.EXODIA_THE_FORBIDDEN_ONE,
	)
	if err != nil {
		panic(err)
	}

	monsterZone[0].BattlePosition = bmgt.FACE_DOWN_DEFENSE_POSITION
	monsterZone[2].BattlePosition = bmgt.FACE_UP_DEFENSE_POSITION
	monsterZone[3].BattlePosition = bmgt.ATTACK_POSITION
	monsterZone[4].BattlePosition = bmgt.FACE_DOWN_DEFENSE_POSITION

	state := bmgt.State{}
	state.P1.MonsterZone = monsterZone
	result := bmgt.ActionsF.NewLegalFlipSummon(&state)
	for _, s := range bmgt.ActionsF.ToStrings(result) {
		fmt.Println(s)
	}
}