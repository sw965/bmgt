package bmgt

import (
	"math/rand"
	omwrand "github.com/sw965/omw/rand"
	"github.com/sw965/omw/fn"
)

func NewNotBadPlayer(r *rand.Rand) func(state *State) Action {
	return func(state *State) Action {
		legalActions := NewLegalActions(state)
		notBadActions := fn.Filter(legalActions, IsNotBadAction(state))
		if len(notBadActions) == 0 {
			return omwrand.Choice(legalActions, r)
		} else {
			return omwrand.Choice(notBadActions, r)
		}
	}
}