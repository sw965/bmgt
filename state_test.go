package bmgt

import (
	"testing"
	omwrand "github.com/sw965/omw/rand"
	"github.com/sw965/crow/game/sequential"
	omwjson "github.com/sw965/omw/json"
	"github.com/sw965/omw/fn"
	"fmt"
)

func TestPush(t *testing.T) {
	r := omwrand.NewMt19937()
	game := sequential.Game[State, Actions, Action]{
		LegalActions:NewLegalActions, Push:Push, IsEnd:IsGameEnd, Player:NewNotBadPlayer(r),
	}
	game.SetRandomActionPlayer(r)
	names := make(CardNames, 40)
	for i := 0; i < 40; i++ {
		names[i] = omwrand.Choice(CARD_NAMES[1:], r)
	}
	state := NewInitState(NewCards(names...), NewCards(names...), true, r)
	state, states, actions := game.PlayoutWithHistory(state, 128)

	d := fn.Map[[]EasyReadState](states, StateToEasyRead)
	fmt.Println(len(d))
	omwjson.Write(&d, "C:/Go/project/foo/js/history2.json")
	a := ActionsToEasyRead(actions)
	omwjson.Write(&a, "C:/Go/project/foo/js/a.json")
}