package bmgt

import (
	"golang.org/x/exp/slices"
	"github.com/sw965/omw/fn"
	omathw "github.com/sw965/omw/math"
	omws "github.com/sw965/omw/slices"
)

type Action struct {
	CardName CardName

	HandIndices []int
	MonsterZoneIndices []int
	SpellTrapZoneIndices []int
	DeckIndices []int

	IsActivationOfCard bool
	EffectNumber int
	IsCost bool
}

type Actions []Action

func NewHandIndexActions(name CardName, handLength, effectNum int, isCost bool) Actions {
	result := make(Actions, handLength)
	for i := 0; i < handLength; i++ {
		result[i] = Action{
			CardName:name,
			HandIndices:[]int{i},
			EffectNumber:effectNum,
			IsCost:isCost,
		}
	}
	return result
}

func NewHandIndicesActions(selectCardNums []int, name CardName, handLength, effectNum int, isCost bool) Actions {
	cs := make([][][]int, len(selectCardNums))
	yn := 0
	for i, r := range selectCardNums {
		c := omathw.Combination{N:handLength, R:r}
		yn += c.TotalNum()
		cs[i] = c.Get()
	}

	result := make(Actions, 0, yn)
	for _, c := range cs {
		for _, idxs := range c {
			action := Action{
				CardName:name,
				HandIndices:idxs,
				EffectNumber:effectNum,
				IsCost:isCost,
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
			CardName:name,
			MonsterZoneIndices:[]int{i},
			EffectNumber:effectNum,
			IsCost:isCost,
		}
	}
	return result
}

func NewDeckIndexActions(name CardName, deckLength, effectNum int, isCost bool) Actions {
	result := make(Actions, deckLength)
	for i := 0; i < deckLength; i++ {
		result[i] = Action{
			CardName:name,
			DeckIndices:[]int{i},
			EffectNumber:effectNum,
			IsCost:isCost,
		}
	}
	return result
}

func NewDeckIndicesActions(selectCardNums []int, name CardName, deckLength, effectNum int, isCost bool) Actions {
	cs := make([][][]int, len(selectCardNums))
	yn := 0
	for i, r := range selectCardNums {
		c := omathw.Combination{N:deckLength, R:r}
		yn += c.TotalNum()
		cs[i] = c.Get()
	}

	result := make(Actions, 0, yn)
	for _, c := range cs {
		for _, indices := range c {
			action := Action{
				CardName:name,
				DeckIndices:indices,
			}
			result = append(result, action)
		}
	}
	return result
}

func NewHandIndexAndMonsterZoneIndexActions(name CardName, handLength, effectNum int, isCost bool) Actions {
	result := make(Actions, 0, MONSTER_ZONE_LENGTH * handLength)
	for handI := 0; handI < handLength; handI++ {
		for zoneI := 0; zoneI < MONSTER_ZONE_LENGTH; zoneI++ {
			action := Action{
				CardName:name,
				HandIndices:[]int{handI},
				MonsterZoneIndices:[]int{zoneI},
				EffectNumber:effectNum,
				IsCost:isCost,
			}
			result = append(result, action)
		}
	}
	return result
}

func NewNormalSummonLegalActions(state *State) Actions {
	return Actions{}
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
					CardName:handCard.Name,
					HandIndices:[]int{handI},
					SpellTrapZoneIndices:[]int{zoneI},
					IsActivationOfCard:true,
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
					CardName:card.Name,
					SpellTrapZoneIndices:[]int{i},
					IsActivationOfCard:true,
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
					CardName:handCard.Name,
					HandIndices:[]int{handI},
					SpellTrapZoneIndices:[]int{zoneI},
					IsActivationOfCard:true,
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
				CardName:card.Name,
				SpellTrapZoneIndices:[]int{i},
				IsActivationOfCard:true,
			}
			result = append(result, action)
		}
	}
	return result
}

type Actionss []Actions

//王立魔法図書館
func NewRoyalMagicalLibraryCostLegalActionss(state *State) Actionss {
	cardName := CardName("王立魔法図書館")
	var effect0 Actions
	if state.CanSpellSpeed1Activation() {
		actions := NewMonsterZoneIndexActions(cardName, 0, true)
		f := func(action Action) bool {
			card := state.P1.MonsterZone[action.MonsterZoneIndices[0]]
			return card.Name == cardName && card.SpellCounter == 3
		}
		effect0 = fn.Filter(actions, f)
	}
	return Actionss{effect0}
}

//サンダー・ドラゴン
func NewThunderDragonCostLegalActionss(state *State) Actionss {
	cardName := CardName("サンダー・ドラゴン")
	var effect0 Actions
	if state.CanSpellSpeed1Activation() && slices.ContainsFunc(state.P1.Deck, EqualNameCard(cardName)) {
		actions := NewHandIndexActions(cardName, len(state.P1.Hand), 0, true)
		f := func(action Action) bool {
			card := state.P1.Hand[action.HandIndices[0]]
			return card.Name == cardName
		}
		effect0 = fn.Filter(actions, f)
	}
	return Actionss{effect0}
}

//サンダー・ドラゴン
func NewThunderDragonEffectLegalActionss(state *State) Actionss {
	cardName := CardName("サンダー・ドラゴン")
	actions := NewDeckIndicesActions([]int{1, 2}, cardName, len(state.P1.Deck), 0, false)
	f := func(action Action) bool {
		for _, idx := range action.DeckIndices {
			if state.P1.Deck[idx].Name != cardName {
				return false
			}
		}
		return true
	}
	effect0 := fn.Filter(actions, f)
	return Actionss{effect0}
}

//召喚僧サモンプリースト
func NewSummonerMonkCostLegalActionss(state *State) Actionss {
	effect0 := Actions{}
	effect1 := Actions{}
	effect2 := Actions{}
	if state.CanSpellSpeed1Activation() && slices.ContainsFunc(state.P1.Deck, IsLevel4MonsterCard){
		cardName := CardName("召喚僧サモンプリースト")
		actions := NewHandIndexAndMonsterZoneIndexActions(cardName, len(state.P1.Hand), 2, true)
		f := func(action Action) bool {
			zoneCard := state.P1.MonsterZone[action.MonsterZoneIndices[0]]
			handCard := state.P1.Hand[action.HandIndices[0]]
			return zoneCard.Name == "召喚僧サモンプリースト" && zoneCard.ThisTurnEffectActivationCounts[2] == 0 && IsSpellCard(handCard)
		}
		effect2 = fn.Filter(actions, f)
	}
	return Actionss{effect0, effect1, effect2}
}

//召喚僧サモンプリースト
func NewSummonerMonkEffectLegalActionss(state *State) Actionss {
	effect0 := Actions{}
	effect1 := Actions{}

	cardName := CardName("召喚僧サモンプリースト")
	actions := NewDeckIndexActions(cardName, len(state.P1.Deck), 2, false)
	f := func(action Action) bool {
		return IsLevel4MonsterCard(state.P1.Deck[action.DeckIndices[0]])
	}
	effect2 := fn.Filter(actions, f)
	return Actionss{effect0, effect1, effect2}
}

//手札断殺
func NewHandDestructionEffectLegalActionss(state *State) Actionss {
	effect0 := NewHandIndicesActions([]int{2}, "手札断殺", len(state.P1.Hand), 0, false)
	return Actionss{effect0}
}

//打ち出の小槌
func NewMagicalMalletEffectLegalActionss(state *State) Actionss {
	n := len(state.P1.Hand)
	rng := omws.IntegerRange[[]int, int]{Start:0, End:n, Step:1}
	selectCardNums := rng.Make()
	effect0 := NewHandIndicesActions(selectCardNums, "打ち出の小槌", n, 0, false)
	return Actionss{effect0}
}