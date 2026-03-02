package bmgt

type UIHook func(UICommand)

type UICommandType int

const (
	NormalSummonUICommand UICommandType = iota
	DeclareAnAttackUICommand
	DirectAttackUICommand
	BattleDestructionUICommand
	LifePointChangeUICommand
)

type UICommand struct {
	Type                          UICommandType
	TurnPlayerHandIndex           int
	NonTurnPlayerHandIndex        int
	TurnPlayerMonsterZoneIndex    int
	NonTurnPlayerMonsterZoneIndex int
	TurnPlayerLifePointDelta      LifePoint
	NonTurnPlayerLifePointDelta   LifePoint
	TurnPlayer                    TurnPlayer
}
