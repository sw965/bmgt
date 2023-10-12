package bmgt

import (
	"golang.org/x/exp/slices"
	omwslices "github.com/sw965/omw/slices"
)

func IsThunderDragonActivatable(state *State) bool {
	return slices.ContainsFunc(state.P1.Deck, EqualNameOfCard(THUNDER_DRAGON))
}

func IsHandMagicalMalletActivatable(state *State) bool {
	return len(state.P1.Hand) >= 2
}

func IsZoneMagicalMalletActivatable(state *State) bool {
	return len(state.P1.Hand) >= 1
}

func IsPotOfGreedActivatable(state *State) bool {
	return len(state.P1.Deck) >= 2
}

func IsGatherYourMindActivatable(state *State) bool {
	return slices.Contains(state.P1.OncePerTurnLimitCardNames, GATHER_YOUR_MIND)
}

func IsHandHandDestructionActivatable(state *State) bool {
	return len(state.P1.Hand) >= 3 && len(state.P2.Hand) >= 2 && len(state.P1.Deck) >= 2 && len(state.P2.Deck) >= 2
}

func IsDoubleSummonActivatable(state *State) bool {
	return !state.P1.IsDoubleSummonApplied
}

func IsToonTableOfContentsActivatable(state *State) bool {
	names := GetNamesOfCards(state.P1.Deck)
	for _, name := range TOON_CARD_NAMES {
		if slices.Contains(names, name) {
			return true
		}
	}
	return false
}

func IsToonWorldActivatable(state *State) bool {
	return state.P1.LifePoint > 1000
}

func IsUpstartGoblinActivatable(state *State) bool {
	return len(state.P1.Deck) >= 1
}

func IsHandMagicalStoneExcavationActivatable(state *State) bool {
	return len(state.P1.Hand) >= 3 && slices.ContainsFunc(CategoriesOfCards(state.P1.Graveyard), IsSpellCategory)
}

func IsZoneMagicalStoneExcavationActivatable(state *State) bool {
	return len(state.P1.Hand) >= 2
}

func IsAllureOfDarknessActivatable(state *State) bool {
	return len(state.P1.Deck) >= 2
}

func IsDarkFactoryOfMassProductionActivatable(state *State) bool {
	return omwslices.CountFunc(state.P1.Graveyard, IsDarkMonsterCard) >= 2
}

func IsSolemnJudgmentActivatable(state *State) bool {
	if state.P1.NormalSummonIndex != -1 {
		return true
	}

	if state.P1.TributeSummonIndex != -1 {
		return true
	}

	if state.P1.FlipSummonIndex != -1 {
		return true
	}

	if state.P1.SpellActivationIndex != -1 {
		return true
	}

	if state.P2.TrapActivationIndex != -1 {
		return true
	}

	if state.P2.NormalSummonIndex != -1 {
		return true
	}

	if state.P2.TributeSummonIndex != -1 {
		return true
	}

	if state.P2.FlipSummonIndex != -1 {
		return true
	}

	if state.P2.SpellActivationIndex != -1 {
		return true
	}

	if state.P2.TrapActivationIndex != -1 {
		return true
	}
	return false
}

func IsJarOfGreedActivatable(state *State) bool {
	return len(state.P1.Deck) >= 1
}

func IsLegacyOfYataGarasuActivatable(state *State) bool {
	return len(state.P1.Deck) >= 1
}

var HAND_SPELL_SPEED1_SPELL_ACTIVATABLE = map[CardName]func(*State)bool{
	ONE_DAY_OF_PEACE:IsOneDayOfPeaceActivatable,
	MAGICAL_MALLET:IsHandMagicalMalletActivatable,
	POT_OF_GREED:IsPotOfGreedActivatable,
	GATHER_YOUR_MIND:IsGatherYourMindActivatable,
	DOUBLE_SUMMON:IsDoubleSummonActivatable,
	TOON_TABLE_OF_CONTENTS:IsToonTableOfContentsActivatable,
	TOON_WORLD:IsToonWorldActivatable,
	UPSTART_GOBLIN:IsUpstartGoblinActivatable,
	MAGICAL_STONE_EXCAVATION:IsHandMagicalStoneExcavationActivatable,
	ALLURE_OF_DARKNESS:IsAllureOfDarknessActivatable,
	DARK_FACTORY_OF_MASS_PRODUCTION:IsDarkFactoryOfMassProductionActivatable,
}

var HAND_QUICK_PLAY_SPELL_ACTIVATABLE = map[CardName]func(*State)bool{
	HAND_DESTRUCTION:IsHandHandDestructionActivatable,
}

var ZONE_SPELL_SPEED1_SPELL_ACTIVATABLE = func() map[CardName]func(*State)bool{
	y := map[CardName]func(*State)bool{
		MAGICAL_MALLET:IsZoneMagicalMalletActivatable,
		MAGICAL_STONE_EXCAVATION:IsZoneMagicalStoneExcavationActivatable,
	}
	for name, f := range HAND_SPELL_SPEED1_SPELL_ACTIVATABLE {
		if _, ok := y[name]; !ok {
			y[name] = f
		}
	}
	return y
}()

var ZONE_QUICK_PLAY_SPELL_ACTIVATABLE = func() map[CardName]func(*State)bool{
	y := map[CardName]func(*State)bool{}
	for name, f := range HAND_QUICK_PLAY_SPELL_ACTIVATABLE {
		if _, ok := y[name]; !ok {
			y[name] = f
		}
	}
	return y
}

var ZONE_NORMAL_TRAP_ACTIVATABLE = map[CardName]func(*State)bool {
	JAR_OF_GREED:IsJarOfGreedActivatable,
	LEGACY_OF_YATA_GARASU:IsLegacyOfYataGarasuActivatable,
}

var ZONE_COUNTER_TRAP_ACTIVATABLE = map[CardName]func(*State)bool {
	SOLEMN_JUDGMENT:IsSolemnJudgmentActivatable,
}