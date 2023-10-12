package bmgt

type CardPlace struct {
	HandIndex int
}

type ActivationOfCard struct {
	Activatable func(*State) bool
	LegalCostActions func(*State, *CardPlace) Actions
	Cost func(*State, *Action)
	Pledge func(*State)
	LegalEffectSelect func(*State) Actions
	LegalTargetActions func(*State) Actions
	LegalEffectActions func(*State, *CardPlace) Actions
	Effect func(*State, *Action) bool
}

func NewActivationOfSummonerMonk() ActivationOfCard {
	activatable := func(state *State) bool {
		return omwslices.ContainsFunc(state.P1.Deck, IsLevel4MonsterCard)
	}

	costActions := func(state *State, place *CardPlace) Actions {
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

	effectActions := func(state *State) Actions {
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
		Activatable:activatable,
		LegalCostActions:costActions,
		Cost:cost,
		Effect:effect,
	}
}

func NewActivationOfOneDayOfPeace() ActivationOfCard {
	activatable := func(state *State, place *CardPlace) bool {
		return len(state.P1.Deck) >= 1 && len(state.P2.Deck) >= 1
	}

	effect := func(state *State, action *Action) {
		state.P1.Draw(1)
		state.P2.Draw(1)
		state.P2.IsOneDayOfPeaceEndTrigger = true
	}
	return ActivationOfCard{
		Activatable:activatable,
		Effect:effect,
	}
}

func NewActivationOfMagicalMallect() ActivationOfCard {
	effectActions := func(state *State) Actions {
		y := make(Actions, 0, 128)
		n := len(state.P1.Hand)
		for i := 0; i < len(state.P1.Hand); i++ {
			r := i + 1
			c := omwmath.Combination{N:n, R:r}
			idxss := c.Get()
			for _, idxs := range idxss {
				action := Action{
					HandIndices:idxs,
				}
				y = append(y, action)
			}
		}
		return y
	}

	effect := func(state *State, action *Action) {
		state.P1.HandToDeck(idxs, true, r)
		state.P1.Draw(len(idxs))
	}
	return ActivationOfCard{effectActions:effectActions, Effect:effect}
}

func NewActivationOfPotOfGreed(state *State) ActivationOfCard {
	activatable := func(state *State) bool {
		return len(state.P1.Deck) >= 2
	}

	effect := func(state *State, action *Action) {
		state.P1.Draw(2)
	}
	return ActivationOfCard{Activatable:activatable, Effect:effect}
}

func NewActivatableOfGatherYourMind(state *State) ActivationOfCard {
	activatable := func(state *State) bool {
		return !omwslices.Contains(state.P1.OncePerTurnLimitCardNames, GATHER_YOUR_MIND)
	}

	pledge := func(state *State) {
		state.P1.OncePerTurnLimitCardNames = append(state.P1.OncePerTurnLimitCardNames, GATHER_YOUR_MIND)
	}

	effectActions := func(state *State) Actions {
		idxs := omwslices.IndicesFunc(state.P1.Deck, EqualNameOfCard(GATHER_YOUR_MIND))
		y := make(Actions, len(idxs))
		for i, idx := range idxs {
			action := Action{
				DeckIndices:[]int{idx}
			}
			y[i] = action
		}
		return y
	}

	effect := func(state *State, action *Action) bool {
		state.P1.Search(action.DeckIndices)
		return true
	}

	return ActivationOfCard{
		Activatable:activatable,
		Pledge:pledge,
		LegalEffectActions:effectActions,
		Effect:effect,
	}
}

func newActivationOfHandDestruction() ActivationOfCard {
	effectActions := func(state *State) Actions {
		c := omwmath.Combination{N:len(state.P1.Hand), R:2}
		idxss := c.Get()
		y := make(Actions, len(idxss))
		for i, idxs := range idxss {
			y[i] = Action{
				HandIndices:idxs,
			}
		}
		return y
	}

	effect := func(state *State, _ *Action) {
		state.P1.Draw(2)
	}

	return ActivationOfCard{
		LegalEffectActions:effectActions,
		Effect:effect,
	}
}

func NewActivationOfHandHandDestruction(state *State) ActivationOfCard {
	activatable := func(state *State) bool {
		return len(state.P1.Hand) >= 3 && len(state.P2.Hand) >= 2 && len(state.P1.Deck) >= 2 && len(state.P2.Deck) >= 2
	}
	y := newActivationOfHandDestruction()
	y.Activatable = activatable
	return y
}

func NewActivationOfZoneHandDestruction(state *State) ActivationOfCard {
	activatable := func(state *State) bool {
		return len(state.P1.Hand) >= 2 && len(state.P2.Hand) >= 2 && len(state.P1.Deck) >= 2 %% len(state.P2.Deck) >= 2
	}
	y := newActivationOfHandDestruction()
	y.Activatable = activatable
	return y
}

func NewActivationOfDoubleSummon(state *State) ActivationOfCard {
	activatable := func(state *State) bool {
		return !state.P1.IsDoubleSummonApplied
	}

	effect := func(state *State, _ *Action) {
		state.P1.IsDoubleSummonApplied = true
	}

	return ActivationOfCard{
		Activatable:activatable,
		Effect:effect,
	}
}

func NewActivatableOfToonTableOfContents(state *State) ActivationOfCard {
	activatable := func(state *State) bool {
		return slices.ContainsFunc(state.P1.Deck, IsToonCard)
	}

	effectActions := func(state *State) Actions {
		idxs := omwslices.IndicesFunc(state.P1.Deck, IsToonCard)
		y := make(Actions, len(idxs))
		for i, idx := range idxs {
			y[i] = Action{
				DeckIndices:[]int{idx},
			}
		}
		return y
	}

	effect := func(state *State, action *Action) {
		state.P1.Deck(action.DeckIndices)
	}

	return ActivationOfCard{
		Activatable:activatable,
		LegalEffectActions:effectActions,
		Effect:effect,
	}
}

func NewActivationOfToonWolrd(state *State) ActivationOfCard {
	activatable := func(state *State) bool {
		return state.P1.LifePoint > 1000
	}

	cost := func(state *State) {
		state.P1.LifePoint -= 1000
	}

	return ActivationOfCard{
		Activatable:activatable,
		Cost:cost,
	}
}

func NewActivationOfUpStartGoblin(state *State) ActivationOfCard {
	activatable := func(state *State) bool [
		return len(state.P1.Deck) >= 1
	]

	return ActivationOfCard{
		Activatable:activatable,
	}
}

func newActivationOfMagicalStoneExcavation(state *State) ActivationOfCard {
	targetActions := func(state *State) Actions {
		cs := CategoriesOfCards(state.P1.Deck)
		idxs := omwslices.IndicesFunc(cs, IsSpellCategory)
		y := make(Actions, len(idxs))
		for i, idx := range idxs {
			id := state.P1.Graveyard[idx]
			action := Action{
				CardIDs:CardIDs{id},
			}
			y[i] = action
		}
		return y
	}

	effect := func(state *State, _ *Action) {
		state.P1.IDSalvage(action.CardIDs)
	}

	return ActivationOfCard{
		LegalEffectActions:effectActions,
		Effect:effect,
	}
}