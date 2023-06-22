package bmgt

func PotOfGreedEffect(state State) (State, error) {
	deck, drawCards, err := state.P1.Deck.Draw(2)
	state.P1.Deck = deck
	state.P1.Hand = append(state.P1.Hand, drawCards...)
	return state, err
}

var EFFECT = map[CardName]func(State) (State, error) {
	"強欲な壺":PotOfGreedEffect,
}