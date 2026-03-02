package bmgt

import "github.com/sw965/crow/game"

func RankByAgentFunc(state *State) (game.RankByAgent[TurnPlayer], error) {
	rank := make(game.RankByAgent[TurnPlayer])

	firstLost := state.First.LifePoint <= 0 || state.First.IsDeckOut
	secondLost := state.Second.LifePoint <= 0 || state.Second.IsDeckOut

	if firstLost && secondLost {
		rank[First] = 1
		rank[Second] = 1
		return rank, nil
	} else if secondLost {
		rank[First] = 1
		rank[Second] = 2
		return rank, nil
	} else if firstLost {
		rank[First] = 2
		rank[Second] = 1
		return rank, nil
	}

	return rank, nil
}
