package  bmgt

import (
	"math/rand"
	omws "github.com/sw965/omw/slices"
)

type ChainLink struct {
	Card Card
	EffectNumber int
}

type Chain []ChainLink

func (c Chain) Resolution(state State, player Player, r *rand.Rand) (State, error) {
	c = omws.Reverse(c)
	var err error
	for _, l := range c {
		card := l.Card
		action := player(&state)
		state, err = EFFECT[card.Name](&action, &card, r)[l.EffectNumber](state)
		if err != nil {
			return state, err
		}
	}
	return state, nil
}
