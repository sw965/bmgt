package bmgt

import (
	"math/rand"
	"golang.org/x/exp/slices"
	omwrand "github.com/sw965/omw/rand"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
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

func (phase Phase) ToString() string {
	switch phase {
		case DRAW_PHASE:
			return "ドローフェイズ"
		case STANDBY_PHASE:
			return "スタンバイフェイズ"
		case MAIN1_PHASE:
			return "メイン1フェイズ"
		case BATTLE_PHASE:
			return "バトルフェイズ"
		case MAIN2_PHASE:
			return "メイン2フェイズ"
		case END_PHASE:
			return "エンドフェイズ"
		default:
			return ""
	}
}

type Phases []Phase

var PHASES = Phases{DRAW_PHASE, STANDBY_PHASE, MAIN1_PHASE, BATTLE_PHASE, MAIN2_PHASE, END_PHASE}

type OneSide struct {
	LifePoint LifePoint
	Hand Cards
	Deck Cards
	MonsterZone Cards
	SpellTrapZone Cards
	Graveyard Cards
	IsDeckOut bool
}

func NewOneSide(deck Cards, r *rand.Rand) OneSide {
	deck = omwrand.Shuffled(deck, r)
	hand := deck[:INIT_DRAW]
	deck = deck[INIT_DRAW:]
	return OneSide{
		LifePoint:INIT_LIFE_POINT,
		Hand:hand, Deck:deck,
		MonsterZone:make(Cards, MONSTER_ZONE_LENGTH),
		SpellTrapZone:make(Cards, SPELL_TRAP_ZONE_LENGTH),
		Graveyard:make(Cards, 0, len(deck) + len(hand)),
	}
}

func (o *OneSide) Draw(num int) {
	drawCards := o.Deck[:num]
	drawCards = fn.Map[Cards](drawCards, func(card Card) Card {
		card.Face = FACE_UP
		return card
	})
	o.Deck = o.Deck[num:]
	o.Hand = append(o.Hand, drawCards...)
}

func (o *OneSide) Discard(idxs []int) {
	n := len(o.Hand)
	newHand := make(Cards, n - len(idxs))
	for i := 0; i < n; i++ {
		if slices.Contains(idxs, i) {
			o.Graveyard = append(o.Graveyard)
		} else {
			newHand = append(newHand, o.Hand[i])
		}
	}
	o.Hand = newHand
}

type Duel struct {
	P1 OneSide
	P2 OneSide
	Phase Phase
	Turn int
}

func NewDuel(p1Deck, p2Deck Cards, r *rand.Rand) Duel {
	p1 := NewOneSide(p1Deck, r)
	p2 := NewOneSide(p2Deck, r)
	return Duel{P1:p1, P2:p2}
}

func (duel *Duel) Reverse() {
	p1 := duel.P1
	p2 := duel.P2
	duel.P1 = p2
	duel.P2 = p1
}

func (duel *Duel) DrawPhase(p1, p2 Player) {
	if len(duel.P1.Deck) == 0 {
		duel.P1.IsDeckOut = true
		return
	}

	duel.P1.Draw(1)

	for {
		action := p1(duel)
		if action.Type == PASS_PRIORITY_ACTION {
			duel.Reverse()
			continue
		}
	}
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

func Push(duel Duel, action Action) Duel {
	if action.Type == PHASE_TRANSITION_ACTION {
		duel.Phase = PHASES[action.N1]
		if duel.Phase == END_PHASE {
			duel.Reverse()
			duel.Turn += 1
			duel.Phase = DRAW_PHASE
			duel.P1.Draw(1)
			duel.Phase = STANDBY_PHASE
			duel.Phase = MAIN1_PHASE
		}
	}
	return duel
}

func IsDuelEnd(duel *Duel) bool {
	if omwslices.IsSubset(duel.P1.Hand.Names() , EXODIA_PARTS_NAMES) {
		return true
	}
	return duel.P1.LifePoint <= 0 || duel.P2.LifePoint <= 0
}

type Duels []Duel