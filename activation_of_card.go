package bmgt

type ActivationOfCard struct {
	Activatable func(*State, *CardPlace) bool
	LegalCostActions func(*State, *CardPlace) Actions
	Cost func(*State, *Action)
	LegalEffectSelect func(*State) Actions
	LegalTargetActions func(*State) Actions
	LegalEffectActions func(*State, *CardPlace) Actions
	Effect func(*State, *Action) bool
}

func NewActivationOfSummonerMonk() ActivationOfCard {
	activatableF := func(state *State) bool {
		return omwslices.ContainsFunc(state.P1.Deck, IsLevel4MonsterCard)
	}

	costActionsF := func(state *State, place *CardPlace) Actions {
		cs := CategoriesOfCards(state.P1.Hand)
		idxs := omwslices.IndicesFunc(cs, IsSpellCategory)
		y := make(Actions, len(idxs))
		for i, idx := range idxs {
			action := Action{
				CardName:SUMMONER_MONK,
				HandIndices:[]int{idx},
				Type:COST_ACTION,
			}
			y[i] = action
		}
		return y
	}

	cost := func(state *State, action *Action) {
		state.P1.Discard(action.HandIndices)
	}

	effectActionsF := func(state *State) Actions {
		idxs := omwslices.IndicesFunc(state.P1.Deck, IsLevel4MonsterCard)
		for _, idx := range idxs {
			action := Action{
				DeckIndices:[]int{idx},
			}
		}
	}

	effect := func(state *State, action *Action) {
		state.P1.Recruit([]int{idx}, true, r)
	}

	return ActivationOfCard{
		ActivatableFunc:activatableF,
		LegalCostActionsFunc:costActionsF,
		Cost:cost,
		Effect:effect,
	}
}

func NewActivationOfOneDayOfPeace() ActivationOfCard {
	activatableF := func(state *State, place *CardPlace) bool {
		return len(state.P1.Deck) >= 1 && len(state.P2.Deck) >= 1
	}

	effect := func(state *State, action *Action) {
		state.P1.Draw(1)
		state.P2.Draw(1)
		state.P2.IsOneDayOfPeaceEndTrigger = true
	}
	return ActivationOfCard{
		ActivatableFunc:activatableF,
		Effect:effect,
	}
}

func NewActivationOfMagicalMallect() ActivationOfCard {
	effect := func(state *State, action *Action) {
		state.P1.HandToDeck(idxs, true, r)
		state.P1.Draw(len(idxs))
	}
	return ActivationOfCard{Effect:effect}
}

