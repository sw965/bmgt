package bmgt

import (
	"github.com/sw965/omw/fn"
)

type EasyReadCard struct {
	Name string
	Attribute string
	Level Level
	Type string
	Atk int
	Def int
	BattlePosition string
	IsAttackDeclared bool
}

func CardToEasyRead(card Card) EasyReadCard {
	return EasyReadCard{
		Name:CARD_NAME_TO_STRING[card.Name],
		Attribute:ATTRIBUTE_TO_STRING[card.Attribute],
		Level:card.Level,
		Type:TYPE_TO_STRING[card.Type],
		Atk:card.Atk,
		Def:card.Def,
		BattlePosition:BattlePositionToString(card.BattlePosition),
		IsAttackDeclared:card.IsAttackDeclared,
	}
}

type EasyReadCards []EasyReadCard

func CardsToEasyRead(cards Cards) EasyReadCards {
	return fn.Map[EasyReadCards](cards, CardToEasyRead)
}

type EasyReadAction struct {
	Indices1 []int
	Indices2 []int
	Phase string
	BattlePosition string
	Type string
}

func ActionToEasyRead(action Action) EasyReadAction {
	return EasyReadAction{
		Indices1:action.Indices1(),
		Indices2:action.Indices2(),
		Phase:PhaseToString(action.Phase),
		BattlePosition:BattlePositionToString(action.BattlePosition),
		Type:ActionTypeToString(action.Type),
	}
}

type EasyReadActions []EasyReadAction

func ActionsToEasyRead(actions Actions) EasyReadActions {
	return fn.Map[EasyReadActions](actions, ActionToEasyRead)
}

type EasyReadOneSideState struct {
	LifePoint LifePoint
	Deck EasyReadCards
	Hand EasyReadCards
	MonsterZone EasyReadCards
	SpellTrapZone EasyReadCards
	Graveyard EasyReadCards
	IsTurn bool
	IsNormalDrawDone bool
	ThisTurnNormalSummonCount int
	IsDeckDeath bool
}

func OneSideStateToEasyRead(oss OneSideState) EasyReadOneSideState {
	return EasyReadOneSideState{
		LifePoint:oss.LifePoint,
		Deck:CardsToEasyRead(oss.Deck),
		Hand:CardsToEasyRead(oss.Hand),
		MonsterZone:CardsToEasyRead(oss.MonsterZone),
		SpellTrapZone:CardsToEasyRead(oss.SpellTrapZone),
		Graveyard:CardsToEasyRead(oss.Graveyard),
		IsNormalDrawDone:oss.IsNormalDrawDone,
		ThisTurnNormalSummonCount:oss.ThisTurnNormalSummonCount,
		IsDeckDeath:oss.IsDeckDeath,
	}
}

type EasyReadState struct {
	P1 EasyReadOneSideState
	P2 EasyReadOneSideState
	Phase string
	Turn int
}

func StateToEasyRead(state State) EasyReadState {
	return EasyReadState{
		P1:OneSideStateToEasyRead(state.P1),
		P2:OneSideStateToEasyRead(state.P2),
		Phase:PhaseToString(state.Phase),
		Turn:state.Turn,
	}
}