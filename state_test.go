package bmgt_test

import (
	"github.com/sw965/bmgt"
	"testing"
	"fmt"
	"golang.org/x/exp/slices"
	omwrand "github.com/sw965/omw/rand"
)

func TestDraw(t *testing.T) {
	r := omwrand.NewMt19937()
	deck := bmgt.OLD_LIBRARY_EXODIA_DECK
	state := bmgt.NewInitState(deck, deck, r)
	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Deck.Names().ToStrings())
	state.P1.Draw(2)
	fmt.Println("")
	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Deck.Names().ToStrings())
}

func TestDiscard(t *testing.T) {
	r := omwrand.NewMt19937()
	deck := bmgt.OLD_LIBRARY_EXODIA_DECK
	state := bmgt.NewInitState(deck, deck, r)
	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Graveyard.Names().ToStrings())
	state.P1.Discard([]int{1, 3})
	fmt.Println("")
	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Graveyard.Names().ToStrings())
}

func TestDiscardBanish(t *testing.T) {
	r := omwrand.NewMt19937()
	deck := bmgt.OLD_LIBRARY_EXODIA_DECK
	state := bmgt.NewInitState(deck, deck, r)
	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Banish.Names().ToStrings())
	state.P1.DiscardBanish([]int{0, 4})
	fmt.Println("")
	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Banish.Names().ToStrings())
}

func TestSearch(t *testing.T) {
	r := omwrand.NewMt19937()
	deck := bmgt.OLD_LIBRARY_EXODIA_DECK
	state := bmgt.NewInitState(deck, deck, r)

	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Deck.Names().ToStrings())
	state.P1.Search([]int{0, 34}, r)
	fmt.Println("")
	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Deck.Names().ToStrings())
}

func TestHandToDeck(t *testing.T) {
	r := omwrand.NewMt19937()
	deck := bmgt.OLD_LIBRARY_EXODIA_DECK
	state := bmgt.NewInitState(deck, deck, r)
	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Hand.IDs())
	fmt.Println(state.P1.Deck.IDs())
	state.P1.HandToDeck([]int{2}, r)
	fmt.Println("")
	fmt.Println(state.P1.Hand.Names().ToStrings())
	fmt.Println(state.P1.Hand.IDs())
	fmt.Println(state.P1.Deck.IDs())
}

func TestLegalActions(t *testing.T) {
	r := omwrand.NewMt19937()
	deck := bmgt.OLD_LIBRARY_EXODIA_DECK
	
	var state bmgt.State
	for {
		state = bmgt.NewInitState(deck, deck, r)
		if slices.ContainsFunc(state.P1.Hand, func(card bmgt.Card) bool { return card.Name == bmgt.HAND_DESTRUCTION }) {
			break
		}
	}
	//state.Phase = bmgt.MAIN1_PHASE
	fmt.Println(state.P1.IsTurn)
	fmt.Println(state.P1.Hand.Names().ToStrings())
	for _, action := range state.LegalActions() {
		action.Println()
		fmt.Println("")
	}
}