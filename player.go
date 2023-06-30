package bmgt

type Player func(*State) Action

func EmptyPlayer(state *State) Action {
	return Action{}
}