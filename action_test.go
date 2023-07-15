package bmgt

import (
	"fmt"
	omwrand "github.com/sw965/omw/rand"
	"golang.org/x/exp/slices"
	"testing"
)

func TestLegalActions(t *testing.T) {
	r := omwrand.NewMt19937()
	var state State
	var err error

	for i := 0; i < 128; i++ {
		state, err = NewInitState(OLD_LIBRARY_EXODIA_DECK, OLD_LIBRARY_EXODIA_DECK, r)
		if err != nil {
			panic(err)
		}
		if slices.ContainsFunc(state.P1.Hand, EqualNameCard("サンダー・ドラゴン")) {
			break
		}
	}
	state.Print()

	for _, action := range CostActionss.ThunderDragon(&state) {
		fmt.Println(action)
	}

	for _, action := range CostActionss.RoyalMagicalLibrary(&state) {
		fmt.Println(action)
	}

	for _, action := range CostActionss.SummonerMonk(&state) {
		fmt.Println(action)
	}
}
