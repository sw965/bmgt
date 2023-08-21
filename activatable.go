package bmgt

import (
	"golang.org/x/exp/slices"
	omwslices "github.com/sw965/omw/slices"
	"fmt"
)

type ActivatableFunc func(*State, *Card) []bool
type Activatable map[CardName]ActivatableFunc

type handSpellSpeed1MonsterActivatable struct{}
var HandSpellSpeed1MonsterActivatable = handSpellSpeed1MonsterActivatable{}

func (h *handSpellSpeed1MonsterActivatable) ThunderDragon(state *State, card *Card) []bool {
	effect0 := slices.Contains(state.P1.Deck.Names(), THUNDER_DRAGON)
	return []bool{effect0}
}

var HAND_SPELL_SPEED1_MONSTER_ACTIVATABLE = Activatable{
	THUNDER_DRAGON:HandSpellSpeed1MonsterActivatable.ThunderDragon,
}

type handNormalSpellActivatable struct{}
var HandNormalSpellActivatable = &handNormalSpellActivatable{}

func (h *handNormalSpellActivatable) OneDayOfPeace(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 1 && len(state.P2.Deck) >= 1
	return []bool{effect0}
}

func (h *handNormalSpellActivatable) MagicalMallet(state *State, card *Card) []bool {
	effect0 := len(state.P1.Hand) >= 2
	return []bool{effect0}
}

func (h *handNormalSpellActivatable) PotOfGreed(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 2
	return []bool{effect0}
}

func (h *handNormalSpellActivatable) GatherYourMind(state *State, card *Card) []bool {
	effect0 := slices.Contains(state.P1.OncePerTurnLimitCardNames, GATHER_YOUR_MIND) && slices.Contains(state.P1.OncePerTurnLimitCardNames, GATHER_YOUR_MIND)
	return []bool{effect0}
}

func (h *handNormalSpellActivatable) DoubleSummon(state *State, card *Card) []bool {
	effect0 := slices.Contains(state.P1.OncePerTurnLimitCardNames, DOUBLE_SUMMON)
	return []bool{effect0}
}

func (h *handNormalSpellActivatable) ToonTableOfContents(state *State, card *Card) []bool {
	effect0 := slices.ContainsFunc(state.P1.Deck.Names(), func(name CardName) bool { return slices.Contains(TOON_CARD_NAMES, name)})
	return []bool{effect0}
}

func (h *handNormalSpellActivatable) UpstartGoblin(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 1
	return []bool{effect0}
}

func (h *handNormalSpellActivatable) MagicalStoneExcavation(state *State, card *Card) []bool {
	effect0 := len(state.P1.Hand) >= 3 && slices.ContainsFunc(state.P1.Graveyard, func(card Card) bool { return card.Category.IsSpell() })
	return []bool{effect0}
}

func (h *handNormalSpellActivatable) AllureOfDarkness(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 2
	return []bool{effect0}
}

func (h *handNormalSpellActivatable) DarkFactoryOfMassProduction(state *State, card *Card) []bool {
	effect0 := omwslices.CountFunc(state.P1.Graveyard.Names(), func(name CardName) bool { return slices.Contains(NORMAL_MONSTER_NAMES, name) }) >= 2
	return []bool{effect0}
}

var HAND_NORMAL_SPELL_ACTIVATABLE = Activatable{
	ONE_DAY_OF_PEACE:HandNormalSpellActivatable.OneDayOfPeace,
	MAGICAL_MALLET:HandNormalSpellActivatable.MagicalMallet,
	GATHER_YOUR_MIND:HandNormalSpellActivatable.GatherYourMind,
	DOUBLE_SUMMON:HandNormalSpellActivatable.DoubleSummon,
	TOON_TABLE_OF_CONTENTS:HandNormalSpellActivatable.ToonTableOfContents,
	UPSTART_GOBLIN:HandNormalSpellActivatable.UpstartGoblin,
	MAGICAL_STONE_EXCAVATION:HandNormalSpellActivatable.MagicalStoneExcavation,
	ALLURE_OF_DARKNESS:HandNormalSpellActivatable.AllureOfDarkness,
	DARK_FACTORY_OF_MASS_PRODUCTION:HandNormalSpellActivatable.DarkFactoryOfMassProduction,
}

type handQuickPlaySpellActivatable struct{}
var HandQuickPlaySpellActivatable = handQuickPlaySpellActivatable{}

func (h *handQuickPlaySpellActivatable) HandDestruction(state *State, card *Card) []bool {
	effect0 := len(state.P1.Hand) >= 3 && len(state.P1.Deck) >= 2 && len(state.P2.Hand) >= 2 && len(state.P2.Deck) >= 2 
	return []bool{effect0}
}

var HAND_QUICK_PLAY_SPELL_ACTIVATABLE = Activatable{
	HAND_DESTRUCTION:HandQuickPlaySpellActivatable.HandDestruction,
}

type handContinuousSpellActivatable struct{}
var HandContinuousSpellActivatable = handContinuousSpellActivatable{}

func(h *handContinuousSpellActivatable) ToonWorld(state *State, card *Card) []bool {
	effect0 := state.P1.LifePoint > 1000
	return []bool{effect0}
}

var HAND_CONTINUOUS_SPELL_ACTIVATABLE = Activatable{
	TOON_WORLD:HandContinuousSpellActivatable.ToonWorld,
}

type zoneSpellSpeed1MonsterActivatable struct{}
var ZoneSpellSpeed1MonsterActivatable = zoneSpellSpeed1MonsterActivatable{}

func (z zoneSpellSpeed1MonsterActivatable) RoyalMagicalLibrary(state *State, card *Card) []bool {
	effect0 := false
	effect1 := card.SpellCounter == 3
	return []bool{effect0, effect1} 
}

func (z zoneSpellSpeed1MonsterActivatable) SummonerMonk(state *State, card *Card) []bool {
	effect0 := false
	effect1 := false
	effect2 := slices.ContainsFunc(state.P1.Hand, func(card Card) bool { return card.Category.IsSpell() }) && slices.ContainsFunc(state.P1.Deck, func(card Card) bool { return card.Level == 4 })
	return []bool{effect0, effect1, effect2}
}

type zoneNormalSpellActivatable struct{}
var ZoneNormalSpellActivatable = zoneNormalSpellActivatable{}

func(z *zoneNormalSpellActivatable) MagicalMallet(state *State, card *Card) []bool {
	effect0 := len(state.P1.Hand) >= 1
	return []bool{effect0}
}

func (z *zoneNormalSpellActivatable) MagicalStoneExcavation(state *State, card *Card) []bool {
	effect0 := len(state.P1.Hand) >= 2
	return []bool{effect0}
}

var ZONE_NORMAL_SPELL_ACTIVATABLE = func() Activatable{
	y := Activatable{
		MAGICAL_MALLET:ZoneNormalSpellActivatable.MagicalMallet,
		MAGICAL_STONE_EXCAVATION:ZoneNormalSpellActivatable.MagicalStoneExcavation,
	}
	for name, f := range HAND_NORMAL_SPELL_ACTIVATABLE {
		if _, ok := y[name]; !ok {
			y[name] = f
		}
	}
	return y
}()


type zoneQuickPlaySpellActivatable struct{}
var ZoneQuickPlaySpellActivatable = zoneQuickPlaySpellActivatable{}

func(z *zoneQuickPlaySpellActivatable) HandDestruction(state *State, card *Card) []bool {
	effect0 := len(state.P1.Hand) >= 2 && len(state.P1.Deck) >= 2 && len(state.P2.Hand) >= 2 && len(state.P2.Deck) >= 2
	return []bool{effect0}
}

var ZONE_QUICK_PLAY_SPELL_ACTIVATABLE = func() Activatable {
	y := Activatable{
		HAND_DESTRUCTION:ZoneQuickPlaySpellActivatable.HandDestruction,
	}
	for name, f := range HAND_QUICK_PLAY_SPELL_ACTIVATABLE {
		if _, ok := y[name]; !ok {
			y[name] = f
		}
	}
	return y
}()

type zoneContinuousSpellActivatable struct{}
var ZoneContinuousSpellActivatable = zoneContinuousSpellActivatable{}

var ZONE_CONTINUOUS_SPELL_ACTIVATABLE = func() Activatable {
	y := Activatable{}
	for name, f := range HAND_CONTINUOUS_SPELL_ACTIVATABLE {
		if _, ok := y[name]; !ok {
			y[name] = f
		}
	}
	return y
}()

type zoneNormalTrapActivatable struct{}
var ZoneNormalTrapActivatable = zoneNormalTrapActivatable{}

func(z *zoneNormalTrapActivatable) JarOfGreed(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 1
	return []bool{effect0}
}

func (z *zoneNormalTrapActivatable) LegacyOfYataGarasu(state *State, card *Card) []bool {
	effect0 := len(state.P1.Deck) >= 1
	return []bool{effect0}
}

var ZONE_NORMAL_TRAP_ACTIVATABLE = Activatable{
	JAR_OF_GREED:ZoneNormalTrapActivatable.JarOfGreed,
	LEGACY_OF_YATA_GARASU:ZoneNormalTrapActivatable.LegacyOfYataGarasu,
}

type zoneCounterTrapActivatable struct{}
var ZoneCounterTrapActivatable = zoneCounterTrapActivatable{}

func(z *zoneCounterTrapActivatable) SolemnJudgment(state *State, card *Card) []bool {
	return []bool{false}
}

var ZONE_COUNTER_TRAP_ACTIVATABLE = Activatable{
	SOLEMN_JUDGMENT:ZoneCounterTrapActivatable.SolemnJudgment,
}

func init() {
	f := func(a Activatable, category Category, s string) {
		for name, _ := range a {
			data := CARD_DATA_BASE[name]
			if data.Category != category {
				msg := fmt.Sprintf(s + " に %v が 入っている", name)
				panic(fmt.Errorf(msg))
			}
		}

		for name, data := range CARD_DATA_BASE {
			if data.Category == category {
				if _, ok := a[name]; !ok {
					fmt.Println(fmt.Sprintf(s + "に %v が 入っていない", CARD_NAME_TO_STRING[name]))
				}
			}
		}
	}
	f(HAND_SPELL_SPEED1_MONSTER_ACTIVATABLE, EFFECT_MONSTER, "HAND_SPELL_SPEED1_MONSTER_ACTIVATABLE")
	f(HAND_NORMAL_SPELL_ACTIVATABLE, NORMAL_SPELL, "HAND_NORMAL_SPELL_ACTIVATABLE")
	f(HAND_QUICK_PLAY_SPELL_ACTIVATABLE, QUICK_PLAY_SPELL, "HAND_QUICK_PLAY_SPELL_ACTIVATABLE")
	f(HAND_CONTINUOUS_SPELL_ACTIVATABLE, CONTINUOUS_SPELL, "HAND_CONTINUOUS_SPELL_ACTIVATABLE")

	f(ZONE_NORMAL_SPELL_ACTIVATABLE, NORMAL_SPELL, "ZONE_NORMAL_SPELL_ACTIVATABLE")
	f(ZONE_QUICK_PLAY_SPELL_ACTIVATABLE, QUICK_PLAY_SPELL, "ZONE_QUICK_PLAY_SPELL_ACTIVATABLE")
	f(ZONE_CONTINUOUS_SPELL_ACTIVATABLE, CONTINUOUS_SPELL, "ZONE_CONTINUOUS_SPELL_ACTIVATABLE")
	f(ZONE_NORMAL_TRAP_ACTIVATABLE, NORMAL_TRAP, "ZONE_NORMAL_TRAP_ACTIVATABLE")
	f(ZONE_CONTINUOUS_SPELL_ACTIVATABLE, CONTINUOUS_SPELL, "ZONE_CONTINUOUS_SPELL_ACTIVATABLE")
}