package bmgt

import (
	"math/rand"
	"github.com/sw965/omw/fn"
	omwrand "github.com/sw965/omw/rand"
)

type Player func(*State) Action

func NewGoodPlayer(r *rand.Rand) Player {
	return func(state *State) Action {
		legalActions := NewLegalActions(state)
		if len(legalActions) == 1 {
			return legalActions[0]
		}
	
		notBadActions := fn.Filter(legalActions, IsBadAction(state))
		return omwrand.Choice(notBadActions, r)
	}
}