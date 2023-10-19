package bmgt

import (
	"testing"
	omwrand "github.com/sw965/omw/rand"
	"fmt"
)

func TestPush(t *testing.T) {
	r := omwrand.NewMt19937()
	names := make(CardNames, 40)
	for i := 0; i < 40; i++ {
		names[i] = omwrand.Choice(EXODIA_PARTS_NAMES, r)
	}
	state := NewInitState(NewCards(names...), NewCards(names...), r)

	for {
		if IsGameEnd(&state) {
			fmt.Println("brak", state.P1.IsDeckDeath, state.P2.IsDeckDeath)
			break
		}
		legalActions := NewLegalActions(&state)
		fmt.Println("action = ", TypesOfActions(legalActions), len(legalActions))
		action := omwrand.Choice(legalActions, r)
		state = Push(state, action)
		p2HandNames := CardNamesToStrings(NamesOfCards(state.P2.Hand))
		p2MonsterZoneNames := CardNamesToStrings(NamesOfCards(state.P2.MonsterZone))
		p1MonsterZoneNames := CardNamesToStrings(NamesOfCards(state.P1.MonsterZone))
		p1HandNames := CardNamesToStrings(NamesOfCards(state.P1.Hand))
		fmt.Println(state.P1.LifePoint, state.P2.LifePoint, PhaseToString(state.Phase), state.Turn, state.P1.IsTurn, state.P2.IsTurn, len(state.P1.Deck), len(state.P2.Deck))
		if state.P2.IsTurn {
			state = state.Reverse()
		}
		fmt.Println(p2HandNames)
		fmt.Println(p2MonsterZoneNames)
		fmt.Println(p1MonsterZoneNames)
		fmt.Println(p1HandNames)
		fmt.Println("---")
	}
}