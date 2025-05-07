package bmgt

import (
	"github.com/sw965/omw/fn"
	omwjson "github.com/sw965/omw/json"
)

type IDOneSide struct {
	LifePoint LifePoint
	Hand CardIDs
	Deck CardIDs
	MonsterZone CardIDs
	SpellTrapZone CardIDs
	Graveyard CardIDs
}

func NewIDOneSide(o *OneSide) IDOneSide {
	return IDOneSide{
		LifePoint:o.LifePoint,
		Hand:o.Hand.IDs(),
		Deck:o.Deck.IDs(),
		MonsterZone:o.MonsterZone.IDs(),
		SpellTrapZone:o.SpellTrapZone.IDs(),
		Graveyard:o.Graveyard.IDs(),
	}
}

type IDuel struct {
	P1 IDOneSide
	P2 IDOneSide
}

func NewIDuel(duel Duel) IDuel {
	return IDuel{
		P1:NewIDOneSide(&duel.P1),
		P2:NewIDOneSide(&duel.P2),
	}
}

func (i *IDuel) WriteJSON(path string) error {
	return omwjson.Write(i, path)
}

type IDuels []IDuel

func NewIDuels(ds Duels) IDuels {
	return fn.Map[IDuels](ds, NewIDuel)
}

type Replay struct {
	CardNames CardNames
	CardIDs CardIDs
	IDuels IDuels
}

func (r *Replay) WriteJSON(path string) error {
	return omwjson.Write(r, path)
}