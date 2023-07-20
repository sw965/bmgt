package bmgt_test

import (
	"github.com/sw965/bmgt"
	"testing"
	"fmt"
	omwrand "github.com/sw965/omw/rand"
)

func TestLegalActions(t *testing.T) {
	r := omwrand.NewMt19937()
	state, err := bmgt.NewInitState(bmgt.OLD_LIBRARY_EXODIA_DECK, bmgt.OLD_LIBRARY_EXODIA_DECK, r)
	if err != nil {
		panic(err)
	}
	zone, err := bmgt.NewCards("封印されしエクゾディア", "", "", "封印されし者の右腕", "")
	if err != nil {
		panic(err)
	}
	state.P1.MonsterZone = zone
	hand, err := bmgt.NewCards("成金ゴブリン", "サンダー・ドラゴン", "手札断殺", "八汰烏の骸", "召喚僧サモンプリースト")
	if err != nil {
		panic(err)
	}

	state.P1.Hand = hand

	state.Phase = bmgt.MAIN1_PHASE
	state.Print()
	for _, action := range state.LegalActions() {
		fmt.Println(action)
	}
} 