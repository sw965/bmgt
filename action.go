package bmgt

type LowLevelMonsterNormalSummonAction  struct {
	HandIndex int
	MonsterZoneIndex int
}

type TributeSummonAction struct {
	HandIndex int
	TributeMonsterZoneIndices []int
	SummonMonsterZoneIndex int
}

type HandSpellTrapActivateAction struct {
	HandIndex int
	SpellTrapZone int
}

type HandSpellTrapSetAction struct {
	HandIndex int
	SpellTrapZone int
}

type Action struct {
	LowLevelMonsterNormalSummon LowLevelMonsterNormalSummonAction
	IsLowLevelMonsterNormalSummon bool

	TributeSummon TributeSummonAction
	IsTributeSummon bool

	HandSpellTrapActivate HandSpellTrapActivateAction
	HandSpellTrapSet HandSpellTrapSetAction
}

func NewEmptyAction() Action {
	return Action{
		PreMoveHandIndex:-1,
		PostMoveMonsterZoneIndex:-1,
		TributeIndicesOfTributeSummon:-1,
		PostMoveSpellTrapZoneIndex:-1,
		SpellTrapZoneIndex:-1,
	}
}

func NewHandSpellTrapActivateAction(handIdx, zoneIdx int) Action {
	result := NewEmptyAction()
	result.PreMoveHandIndex = handIdx
	result.PostMoveMonsterZoneIndex = zoneIdx
	result.IsSet = false
	return result
}

func NewSetSpellTrapActivateAction(zoneIdx int) Action {
	result := NewEmptyAction()
	result.SpellTrapZoneIndex = zoneIdx
	return result
}

func NewLowLevelMonsterNormalSummonAction(handIdx, zoneIdx int) Action {
	result := NewEmptyAction()
	result.PreMoveHandIndex = handIdx
	result.PostMoveMonsterZoneIndex = zoneIdx
	return result
}

func NewTributeActionOfTributeSummon(handIdx int, tributeIndices []int) Action {
	result := NewEmptyAction()
	result.TributeIndicesOfTributeSummon = tributeIndices
	return result
}

type Actions []Action