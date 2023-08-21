package bmgt

type Player func(*State) Action

// func NewRandomActionPlayer(r *rand.Rand) Player {
// 	return func(state *State) Action {
// 		actions := state.Leg
// 	}
// }