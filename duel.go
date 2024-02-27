package bmgt

import (
	"math/rand"
	omwrand "github.com/sw965/omw/rand"
)

const (
	INIT_DRAW = 5
	MONSTER_ZONE_LENGTH = 5
	SPELL_TRAP_ZONE_LENGTH = 5
)

const (
	DRAW_PHASE_INDEX = iota
	STANDBY_PHASE_INDEX
	MAIN1_PHASE_INDEX
	BATTLE_PHASE_INDEX
	MAIN2_PHASE_INDEX
	END_PHASE_INDEX
)

type LifePoint int

const (
	INIT_LIFE_POINT LifePoint = 8000
)

type Phase int

const (
	DRAW_PHASE Phase = iota
	STANDBY_PHASE
	MAIN1_PHASE
	BATTLE_PHASE
	MAIN2_PHASE
	END_PHASE
)

type Phases []Phase

var PHASES = Phases{DRAW_PHASE, STANDBY_PHASE, MAIN1_PHASE, BATTLE_PHASE, MAIN2_PHASE, END_PHASE}

type OneSideDuelState struct {
	LifePoint LifePoint
	Hand Cards
	Deck Cards
	MonsterZone Cards
	SpellTrapZone Cards
	Graveyard Cards
}

func NewOneSideDuelState(deck Cards, r *rand.Rand) OneSideDuelState {
	deck = omwrand.Shuffled(deck, r)
	hand := deck[:INIT_DRAW]
	deck = deck[INIT_DRAW:]
	return OneSideDuelState{
		LifePoint:INIT_LIFE_POINT,
		Hand:hand, Deck:deck,
		MonsterZone:make(Cards, MONSTER_ZONE_LENGTH),
		SpellTrapZone:make(Cards, SPELL_TRAP_ZONE_LENGTH),
		Graveyard:make(Cards, 0, len(deck) + len(hand)),
	}
}

func (osds *OneSideDuelState) Draw(num int) {
	drawCards := osds.Deck[:num]
	deck = osds.Deck[num:]
	osds.Hand = append(osds.Hand, drawCards...)
}

func (osds *OneSideDuelState) Discard(idxs []int) {
	n := len(osds.Hand)
	newHand := make(Cards, n - len(idxs))
	for i := 0; i < n; i++ {
		if slices.Contains(idxs, i) {
			osds.Graveyard = append(osds.Graveyard)
		} else {
			newHand = append(newHand, osds.Hand[i])
		}
	}
	osds.Hand = newHand
}

type Duel struct {
	P1 OneSideDuelState
	P2 OneSideDuelState
	Phase Phase
}

func NewDuel(p1Deck, p2Deck Cards, r *rand.Rand) Duel {
	p1 := NewOneSideDuelState(p1Deck, r)
	p2 := NewOneSideDuelState(p2Deck, r)
	return Duel{P1:p1, P2:p2}
}

func (duel *Duel) Reverse() {
	p1 := duel.P1
	p2 := duel.P2
	duel.P1 = p2
	duel.P2 = p1
}

func(duel *Duel) DestructionByBattle(i int) {
	card := duel.P1.MonsterZone[i]
	duel.P1.Graveyard = append(duel.P1.Graveyard, card)
	duel.P1.MonsterZone[i] = Card{}
}

func (duel *Duel) Battle(action *Action) {
	attacker := duel.P1.MonsterZone[action.N1]
	defender := duel.P2.MonsterZone[action.N2]
	dbp := defender.BattlePosition()
	atkV := attacker.Atk
	defV := map[BattlePosition]int{
		FACE_UP_ATTACK_POSITION:defender.Atk,
		FACE_UP_DEFENSE_POSITION:defender.Def,
		FACE_DOWN_DEFENSE_POSITION:defender.Def,
	}[dbp]
	
	p1Dmg := 0
	p2Dmg := 0
	destroyP1 := false
	destroyP2 := false

	changeBattlePosition := func() {
		duel.P2.MonsterZone[action.N2].SetBattlePosition(FACE_DOWN_DEFENSE_POSITION)
	}

	destroyP2Func := func() {
		duel.Reverse()
		duel.DestructionByBattle(action.N2)
		duel.Reverse()
	}

	if atkV > defV {
		destroyP2 = true
		if dbp == FACE_UP_ATTACK_POSITION {
			p2Dmg = atkV - defV 
		}
	} else if atkV < defV {
		destroyP1 = true
		if dbp == FACE_UP_ATTACK_POSITION {
			p1Dmg = defV - atkV
		}
	} else {
		if dbp == FACE_UP_ATTACK_POSITION {
			destroyP1 = true
			destroyP2 = true
		}
	}

	if dbp == FACE_DOWN_DEFENSE_POSITION {
		changeBattlePosition()
	}

	duel.P1.LifePoint -= LifePoint(p1Dmg)
	duel.P2.LifePoint -= LifePoint(p2Dmg)

	if destroyP1 && destroyP2 {
		duel.DestructionByBattle(action.N1)
		destroyP2Func()
	} else if destroyP1 {
		duel.DestructionByBattle(action.N1)
	} else if destroyP2 {
		destroyP2Func()
	}
}