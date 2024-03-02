package bmgt

type Effect func(duel *Duel)
type Effects []Effect

func NewPotOfGreedEffects() Effects {
	var effect0 Effect
	effect0 = func(duel *Duel) {
		duel.P1.Draw(POT_OF_GREED_DRAW_NUM)
	}
	return Effects{effect0}
}