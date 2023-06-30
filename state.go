package bmgt

import (
	"math/rand"
	omwrand "github.com/sw965/omw/rand"
	omws "github.com/sw965/omw/slices"
	"fmt"
)

type LifePoint int

const (
	INIT_LIFE_POINT = 8000
)

const (
	MONSTER_ZONE_SIZE = 5
	SPELL_TRAP_ZONE_SIZE = 5
)

type OneSideState struct {
	LifePoint LifePoint
	Hand Cards
	Deck Cards
	MonsterZone Cards
	SpellTrapZone Cards
	Graveyard Cards

	IsPriorityWaiver bool
	CurrentTurnNormalSummonUpperLimit int
	CurrentTurnNormalSummonNum int
	IsDeclareAnAttack bool
}

func NewOneSideState(deck Cards, r *rand.Rand) (OneSideState, error) {
	var hand Cards
	var err error
	deck = omwrand.Shuffled(deck, r)
	deck, hand, err = deck.Draw(5)

	result := OneSideState{}
	result.LifePoint = INIT_LIFE_POINT
	result.Hand = hand
	result.Deck = deck
	result.Graveyard = make(Cards, 0, len(deck))
	return result, err
}

func (oss OneSideState) Draw(num int) (OneSideState, error) {
	deck, drawCards, err := oss.Deck.Draw(num)
	oss.Deck = deck
	oss.Hand = append(oss.Hand, drawCards...)
	return oss, err
}

func (oss *OneSideState) CanNormalSummon() bool {
	return oss.CurrentTurnNormalSummonNum < oss.CurrentTurnNormalSummonUpperLimit
}

type Phase string

const (
	DRAW_PHASE = Phase("ドロー")
	STANDBY_PHASE = Phase("スタンバイ")
	MAIN1_PHASE = Phase("メイン1")
	BATTLE_PHASE = Phase("バトル")
	MAIN2_PHASE = Phase("メイン2")
	END_PHASE = Phase("エンド")
)

type State struct {
	P1 OneSideState
	P2 OneSideState
	IsP1Turn bool
	IsP1Priority bool
	Phase Phase
	Chain Cards
}

func NewInitState(p1Deck, p2Deck Cards, r *rand.Rand) (State, error) {
	p1Deck = omwrand.Shuffled(p1Deck, r)
	p2Deck = omwrand.Shuffled(p2Deck, r)

	p1Hand, p1Deck, err := p1Deck.Draw(5)
	if err != nil {
		return State{}, err
	}

	p2Hand, p2Deck, err := p2Deck.Draw(5)
	if err != nil {
		return State{}, err
	}

	p1 := OneSideState{}
	p1.LifePoint = INIT_LIFE_POINT
	p1.Hand = p1Hand
	p1.Deck = p1Deck

	p2 := OneSideState{}
	p2.LifePoint = INIT_LIFE_POINT
	p2.Hand = p2Hand
	p2.Deck = p2Deck

	state := State{P1:p1, P2:p2}
	state.IsP1Turn = true
	state.IsP1Priority = true
	state.Phase = DRAW_PHASE
	return state, nil
}

func (state State) LegalActions() Actions {
	eqNameF := func(name CardName) func(card Card) bool {
		return func(card Card) bool {
			return card.Name == name
		}
	}
	p1MonsterZoneEmptyIndices := omws.IndicesFunc(state.P1.MonsterZone, IsEmptyCard)
	p1MonsterZoneNotEmptyIndices := omws.IndicesFunc(state.P1.MonsterZone, IsNotEmptyCard)
	p1SpellTrapZoneEmptyIndices := omws.IndicesFunc(state.P1.SpellTrapZone, IsEmptyCard)

	result := make(Actions, 0, 64)
	if state.Phase == DRAW_PHASE || state.Phase == STANDBY_PHASE {
		//手札にある速攻魔法の合法手
		for _, card := range state.P1.Hand {
			if !card.IsQuickPlaySpell {
				continue
			}

			if ACTIVATE_OK[card.Name](&state) {
				for _, handIdx := range omws.IndicesFunc(state.P1.Hand, eqNameF(card.Name)) {
					for _, zoneIdx := range p1SpellTrapZoneEmptyIndices {
						result = append(result, Action{ActivateHandSpellTrap:MoveIndex{Pre:handIdx, Post:zoneIdx}})
					}
				}
			}
		}

		//セットされている速攻魔法と罠の合法手
		for _, card := range state.P1.SpellTrapZone {
			if !IsSpellSpeed2Card(card) {
				continue
			}

			if ACTIVATE_OK[card.Name](&state) && !card.IsSetTurn {
				for _, zoneIdx := range omws.IndicesFunc(state.P1.SpellTrapZone, eqNameF(card.Name)) {
					result = append(result, Action{SetHandSpellTrap:zoneIdx})
				}
			}
		}
	}

	if state.Phase == MAIN1_PHASE {
		for _, card := range state.P1.Hand {
			handNameIndices := omws.IndicesFunc(state.P1.Hand, eqNameF(card.Name))

			//生贄なしの通常召喚
			if IsLowLevelMonsterCard(card) && card.CanNormalSummon && state.P1.CanNormalSummon() {
				for _, handIdx := range handNameIndices {
					for _, zoneIdx := range p1MonsterZoneEmptyIndices {
						action := LowLevelMonsterNormalSummonAction{HandIndex:handIdx, MonsterZoneIndex:zoneIdx}
						result = append(result, Action{LowLevelMonsterNormalSummon:action, IsLowLevelMonsterNormalSummon:true})
					}
				}
			//一体生贄の通常召喚
			} else if IsMediumLevelMonsterCard(card) && card.CanNormalSummon && state.P1.CanNormalSummon() {
				for _, handIdx := range handNameIndices {
					for _, tributeIdx := range p1MonsterZoneNotEmptyIndices {
						action := 
						result = append(result, NewTributeActionOfTributeSummon([]int{tributeIdx}) )
					}
				}
			//手札にある魔法罠
			} else if card.IsSpellTrap() && ACTIVATE_OK[card.Name] {
				for _, handIdx := range handNameIndices {
					for _, zoneIdx := range p1SpellTrapZoneEmptyIndices {
						//魔法カードの時のみセットせずに発動可能(手札から罠が発動出来る場合もあるが、そのようなカードは未実装な為)
						if card.IsSpell() {
							result = append(result, NewHandSpellTrapAction(handIdx, zoneIdx, false)...)
						}
						result = append(result, NewHandSpellTrapAction(handIdx, zoneIdx, true)...)
					}
				}
			}
		}

		//起動効果
		for _, card := range state.P1.MonsterZone {
			
		}

		//セットされている魔法・罠
		for _, card := range state.P1.SpellTrapZone {
			var activateOK bool
			if card.IsQuickPlaySpell || card.IsTrap() {
				activateOK = ACTIVATE_OK[card.Name] && !card.IsSetTurn
			} else {
				activateOK = ACTIVATE_OK[card.Name]
			}

			if activateOK {
				for _, zoneIdx := range omw.IndicesFunc(state.P1.SpellTrapZone, eqNameF(card.Name)) {
					result = append(result, NewSetSpellTrapActivateAction(zoneIdx))
				}
			}
		}
	}
	return result
}