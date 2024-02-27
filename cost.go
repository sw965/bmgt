package bmgt

func NewMagicalStoneExcavationActionss(duel *Duel) Actionss {
	cost0 := make(Actions, 0)
	n := len(duel.P1.Hand)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			cost0 := Action{N1 = 1<<(1) + 1<<(j)}
		}
	}
	return Actionss{cost0}
}

type Cost func(*Duel)
type Costs []Cost

func NewMagicalStoneExcavationCosts() Costs {
	cost0 := func(duel *Duel, action *Action) {
		idxs := omwslices.Indices(omwslices.Binary(action.N1), 1)
		duel.P1.Discard(idxs)
	}
	return Costs{cost0}
}