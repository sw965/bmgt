package bmgt

type ActionType int

const (
	NORMAL_SUMMON_ACTION ActionType = iota
	TRIBUTE_SUMMON_ACTION
	FLIP_SUMMON_ACTION
	HAND_SPELL_ACTIVATION_ACTION
	HAND_SPELL_TRAP_SET_ACTION
	ZONE_SPELL_TRAP_ACTIVATION_ACTION
	COST_ACTION
	TARGET_ACTION
	EFFECT_SELECT_ACTION
)

type Action struct {
	CardName CardName
	CardIDs CardIDs
	HandIndices []int
	MonsterZoneIndices1 []int
	MonsterZoneIndices2 []int
	SpellTrapZoneIndices []int
	BattlePosition BattlePosition
	EffectSelectNumber int
	Type ActionType
}

type Actions []Action