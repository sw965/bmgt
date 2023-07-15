package bmgt

import (
	omathw "github.com/sw965/omw/math"
)

type Action struct {
	CardName CardName

	HandIndices          []int
	MonsterZoneIndices   []int
	SpellTrapZoneIndices []int
	DeckIndices          []int
	GraveyardIndices     []int

	IsActivationOfCard bool
	EffectNumber       int
	SelectEffectNumber int
	IsCost             bool
}

type Actions []Action

func NewHandIndexActions(name CardName, handLen, effectNum int, isCost bool) Actions {
	result := make(Actions, handLen)
	for i := 0; i < handLen; i++ {
		result[i] = Action{
			CardName:     name,
			HandIndices:  []int{i},
			EffectNumber: effectNum,
			IsCost:       isCost,
		}
	}
	return result
}

func NewHandIndicesActions(selectCardNums []int, name CardName, handLen, effectNum int, isCost bool) Actions {
	cs := make([][][]int, len(selectCardNums))
	yn := 0
	for i, r := range selectCardNums {
		c := omathw.Combination{N: handLen, R: r}
		yn += c.TotalNum()
		cs[i] = c.Get()
	}

	result := make(Actions, 0, yn)
	for _, c := range cs {
		for _, idxs := range c {
			action := Action{
				CardName:     name,
				HandIndices:  idxs,
				EffectNumber: effectNum,
				IsCost:       isCost,
			}
			result = append(result, action)
		}
	}
	return result
}

func NewMonsterZoneIndexActions(name CardName, effectNum int, isCost bool) Actions {
	result := make(Actions, MONSTER_ZONE_LENGTH)
	for i := 0; i < MONSTER_ZONE_LENGTH; i++ {
		result[i] = Action{
			CardName:           name,
			MonsterZoneIndices: []int{i},
			EffectNumber:       effectNum,
			IsCost:             isCost,
		}
	}
	return result
}

func NewDeckIndexActions(name CardName, deckLen, effectNum int, isCost bool) Actions {
	result := make(Actions, deckLen)
	for i := 0; i < deckLen; i++ {
		result[i] = Action{
			CardName:     name,
			DeckIndices:  []int{i},
			EffectNumber: effectNum,
			IsCost:       isCost,
		}
	}
	return result
}

func NewDeckIndicesActions(selectCardNums []int, name CardName, deckLen, effectNum int, isCost bool) Actions {
	cs := make([][][]int, len(selectCardNums))
	yn := 0
	for i, r := range selectCardNums {
		c := omathw.Combination{N: deckLen, R: r}
		yn += c.TotalNum()
		cs[i] = c.Get()
	}

	result := make(Actions, 0, yn)
	for _, c := range cs {
		for _, indices := range c {
			action := Action{
				CardName:    name,
				DeckIndices: indices,
			}
			result = append(result, action)
		}
	}
	return result
}

func NewGraveyardIndexActions(name CardName, graveyardLen, effectNum int, isCost bool) Actions {
	result := make(Actions, graveyardLen)
	for i := 0; i < graveyardLen; i++ {
		result[i] = Action{
			CardName:         name,
			GraveyardIndices: []int{i},
			EffectNumber:     effectNum,
			IsCost:           isCost,
		}
	}
	return result
}

func NewGraveyardIndicesActions(selectCardNums []int, name CardName, graveyardLen, effectNum int, isCost bool) Actions {
	cs := make([][][]int, len(selectCardNums))
	yn := 0
	for i, r := range selectCardNums {
		c := omathw.Combination{N: graveyardLen, R: r}
		yn += c.TotalNum()
		cs[i] = c.Get()
	}

	result := make(Actions, 0, yn)
	for _, c := range cs {
		for _, idxs := range c {
			action := Action{
				CardName:         name,
				GraveyardIndices: idxs,
				EffectNumber:     effectNum,
				IsCost:           isCost,
			}
			result = append(result, action)
		}
	}
	return result
}

func NewHandIndexAndMonsterZoneIndexActions(name CardName, handLen, effectNum int, isCost bool) Actions {
	result := make(Actions, 0, MONSTER_ZONE_LENGTH*handLen)
	for handI := 0; handI < handLen; handI++ {
		for zoneI := 0; zoneI < MONSTER_ZONE_LENGTH; zoneI++ {
			action := Action{
				CardName:           name,
				HandIndices:        []int{handI},
				MonsterZoneIndices: []int{zoneI},
				EffectNumber:       effectNum,
				IsCost:             isCost,
			}
			result = append(result, action)
		}
	}
	return result
}

func NewHandNormalSpellCardActivationLegalActions(state *State) Actions {
	if !state.CanSpellSpeed1Activation() {
		return Actions{}
	}
	result := make(Actions, 0, len(state.P1.Hand))
	for handI, handCard := range state.P1.Hand {
		for zoneI, zoneCard := range state.P1.SpellTrapZone {
			data := CARD_DATA_BASE[handCard.Name]
			if data.IsNormalSpell && IsEmptyCard(zoneCard) {
				action := Action{
					CardName:             handCard.Name,
					HandIndices:          []int{handI},
					SpellTrapZoneIndices: []int{zoneI},
					IsActivationOfCard:   true,
				}
				result = append(result, action)
			}
		}
	}
	return result
}

func NewSetNormalSpellCardActivationLegalActions(state *State) Actions {
	if !state.CanSpellSpeed1Activation() {
		return Actions{}
	}
	result := make(Actions, 0, SPELL_TRAP_ZONE_LENGTH)
	for i, card := range state.P1.SpellTrapZone {
		if card.IsSet {
			data := CARD_DATA_BASE[card.Name]
			if data.IsNormalSpell {
				action := Action{
					CardName:             card.Name,
					SpellTrapZoneIndices: []int{i},
					IsActivationOfCard:   true,
				}
				result = append(result, action)
			}
		}
	}
	return result
}

func NewHandQuickPlaySpellCardActivationLegalActions(state *State) Actions {
	result := make(Actions, 0, len(state.P1.Hand))
	for handI, handCard := range state.P1.Hand {
		for zoneI, zoneCard := range state.P1.SpellTrapZone {
			data := CARD_DATA_BASE[handCard.Name]
			if data.IsQuickPlaySpell && IsEmptyCard(zoneCard) {
				action := Action{
					CardName:             handCard.Name,
					HandIndices:          []int{handI},
					SpellTrapZoneIndices: []int{zoneI},
					IsActivationOfCard:   true,
				}
				result = append(result, action)
			}
		}
	}
	return result
}

func NewSetNormalTrapCardActivationLegalActions(state *State) Actions {
	result := make(Actions, 0, SPELL_TRAP_ZONE_LENGTH)
	for i, card := range state.P1.SpellTrapZone {
		data := CARD_DATA_BASE[card.Name]
		if data.IsNormalTrap && card.IsSet && !card.IsSetTurn {
			action := Action{
				CardName:             card.Name,
				SpellTrapZoneIndices: []int{i},
				IsActivationOfCard:   true,
			}
			result = append(result, action)
		}
	}
	return result
}

type Actionss []Actions
