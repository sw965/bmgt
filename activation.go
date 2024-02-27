package bmgt

func IsPotOfGreedActivationPossible(duel *Duel) []bool {
	effect0 := duel.Deck >= len(duel.P1.Deck)
	return []bool{effect0}
}