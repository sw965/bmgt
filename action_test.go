package bmgt

import (
	"testing"
	omwrand "github.com/sw965/omw/rand"
	"fmt"
	"golang.org/x/exp/slices"
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

	for _, action := range NewThunderDragonCostLegalActionss(&state) {
		fmt.Println(action)
	}

	for _, action := range NewRoyalMagicalLibraryCostLegalActionss(&state) {
		fmt.Println(action)
	}

	for _, action := range NewSummonerMonkCostLegalActionss(&state) {
		fmt.Println(action)
	}

	actions := NewThunderDragonEffectLegalActionss(&state)[0]
	state1, err := EFFECT[actions[0].CardName](&actions[0], &Card{}, r)[actions[0].EffectNumber](state)
	if err != nil {
		panic(err)
	}
	state1.Print()
}