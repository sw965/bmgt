package bmgt

func IsPotOfGreedActivationPossible(duel *Duel) []bool {
	effect0 := len(duel.P1.Deck) >= 2
	return []bool{effect0}
}