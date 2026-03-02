package bmgt

import "fmt"

type battleManager struct {
	state              *State
	fromIdx            int
	targetIdx          int
	attacker           *Card
	defender           *Card
	attackPlayerState  *OneSideState
	defensePlayerState *OneSideState
	hook               UIHook
}

func (bm *battleManager) handleVictory() {
	// 相手モンスターを破壊し、墓地へ送る
	bm.defensePlayerState.Graveyard = append(bm.defensePlayerState.Graveyard, *bm.defender)
	bm.defensePlayerState.MonsterZone[bm.targetIdx] = Card{}
	bm.state.triggerUI(bm.hook, UICommand{
		Type:                          BattleDestructionUICommand,
		TurnPlayerMonsterZoneIndex:    -1,
		NonTurnPlayerMonsterZoneIndex: bm.targetIdx,
	})

	dmg := LifePoint(bm.attacker.Atk - bm.defender.Atk)
	bm.defensePlayerState.LifePoint -= dmg
	bm.state.triggerUI(bm.hook, UICommand{
		Type:                        LifePointChangeUICommand,
		NonTurnPlayerLifePointDelta: -dmg,
	})
}

func (bm *battleManager) handleDefeat() {
	bm.attackPlayerState.Graveyard = append(bm.attackPlayerState.Graveyard, *bm.attacker)
	bm.attackPlayerState.MonsterZone[bm.fromIdx] = Card{}
	bm.state.triggerUI(bm.hook, UICommand{
		Type:                       BattleDestructionUICommand,
		TurnPlayerMonsterZoneIndex: bm.fromIdx,
	})

	dmg := LifePoint(bm.defender.Atk - bm.attacker.Atk)
	bm.attackPlayerState.LifePoint -= dmg
	bm.state.triggerUI(bm.hook, UICommand{
		Type:                     LifePointChangeUICommand,
		TurnPlayerLifePointDelta: -dmg,
	})
}

func (bm *battleManager) handleDraw() {
	bm.attackPlayerState.Graveyard = append(bm.attackPlayerState.Graveyard, *bm.attacker)
	bm.attackPlayerState.MonsterZone[bm.fromIdx] = Card{}
	bm.defensePlayerState.Graveyard = append(bm.defensePlayerState.Graveyard, *bm.defender)
	bm.defensePlayerState.MonsterZone[bm.targetIdx] = Card{}
	bm.state.triggerUI(bm.hook, UICommand{
		Type:                          BattleDestructionUICommand,
		TurnPlayerMonsterZoneIndex:    bm.fromIdx,
		NonTurnPlayerMonsterZoneIndex: bm.targetIdx,
	})
}

func (s *State) Battle(fromIdx, targetIdx int, hook UIHook) error {
	tps := s.TurnPlayerState()
	ntps := s.NonTurnPlayerState()

	if tps.MonsterZone[fromIdx].Id == 0 {
		return fmt.Errorf("空のモンスターゾーンで攻撃しようとした")
	}
	if ntps.MonsterZone[targetIdx].Id == 0 {
		return fmt.Errorf("空のモンスターゾーンに対して攻撃しようとした")
	}

	// 攻撃宣言をUIに通知
	s.triggerUI(hook, UICommand{
		Type:                          DeclareAnAttackUICommand,
		TurnPlayerMonsterZoneIndex:    fromIdx,
		NonTurnPlayerMonsterZoneIndex: targetIdx,
	})

	bm := &battleManager{
		state:              s,
		fromIdx:            fromIdx,
		targetIdx:          targetIdx,
		attacker:           &tps.MonsterZone[fromIdx],
		defender:           &ntps.MonsterZone[targetIdx],
		attackPlayerState:  tps,
		defensePlayerState: ntps,
		hook:               hook,
	}

	diff := bm.attacker.Atk - bm.defender.Atk
	switch {
	case diff > 0:
		bm.handleVictory()
	case diff < 0:
		bm.handleDefeat()
	default:
		bm.handleDraw()
	}

	// 自分のモンスターがフィールドに残っていれば攻撃済みフラグを立てる
	if tps.MonsterZone[fromIdx].Id != 0 {
		tps.MonsterZone[fromIdx].IsAttacked = true
	}
	return nil
}

func (s *State) DirectAttack(fromIdx int, hook UIHook) error {
	tps := s.TurnPlayerState()
	if tps.MonsterZone[fromIdx].Id == 0 {
		return fmt.Errorf("空のモンスターゾーンでダイレクトアタックしようとした")
	}

	attacker := tps.MonsterZone[fromIdx]
	ntps := s.NonTurnPlayerState()

	s.triggerUI(hook, UICommand{
		Type:                       DirectAttackUICommand,
		TurnPlayerMonsterZoneIndex: fromIdx,
	})

	// 直接攻撃のダメージ適用
	dmg := LifePoint(attacker.Atk)
	ntps.LifePoint -= dmg
	s.triggerUI(hook, UICommand{
		Type:                        LifePointChangeUICommand,
		NonTurnPlayerLifePointDelta: -dmg,
	})

	tps.MonsterZone[fromIdx].IsAttacked = true
	return nil
}
