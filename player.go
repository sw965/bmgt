package bmgt

import (
	"math/rand"
	omwrand "github.com/sw965/omw/rand"
)

type Player func(*Duel) Action

func NewRandomPlayer(r *rand.Rand) Player {
	return func(duel *Duel) Action {
		y := omwrand.Choice(NewLegalActions(duel), r)
		return y
	}
}