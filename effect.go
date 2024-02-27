package bmgt

type Effect func(duel *Duel)
type Effects []Effect

func NewPotOfGreedEffects() Effects {
	var effect0 Effect
	effect0 = func(duel *Duel) {
		duel.Draw(2)
	}
	return Effects{effect0}
}