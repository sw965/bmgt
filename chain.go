package bmgt

type ChainLink struct {
	Card Card
}

type Chain []ChainLink

func (chain Chain) GetLast() ChainLink {
	return chain[len(chain) - 1]
}