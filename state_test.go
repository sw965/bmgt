package bmgt_test

import (
	"github.com/sw965/bmgt"
	"testing"
	"golang.org/x/exp/slices"
)

func TestState_LegalZeroTributeNormalSummonAndSetActions(t *testing.T) {
	p1State := bmgt.OneSideState{}
	p1State.Hand = bmgt.NewCardsWithPanic("封印されしエクゾディア", "強欲な瓶", "サンダー・ドラゴン", "トゥーンのもくじ")
	p1State.MonsterZone = bmgt.NewCardsWithPanic("", "", "封印されし者の右腕", "", "封印されし者の左足").ToMonsterZoneWithPanic()
	state := bmgt.State{P1:p1State}

	result := state.LegalZeroTributeNormalSummonAndSetActions()
	expected := bmgt.Actions{
		bmgt.Action{PreMoveHandIndex:0, PostMoveMonsterZoneIndex:0, IsNormalSummon:true},
		bmgt.Action{PreMoveHandIndex:0, PostMoveMonsterZoneIndex:1, IsNormalSummon:true},
		bmgt.Action{PreMoveHandIndex:0, PostMoveMonsterZoneIndex:3, IsNormalSummon:true},
	}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestState_Push(t *testing.T) {
	var err error
	p1State := bmgt.OneSideState{}
	p1State.Hand = bmgt.NewCardsWithPanic("サンダー・ドラゴン", "八汰烏の骸", "精神統一", "封印されし者の左腕", "封印されし者の右腕")
	p1State.MonsterZone = bmgt.NewCardsWithPanic("", "封印されしエクゾディア", "", "", "").ToMonsterZoneWithPanic()
	p1State.MonsterZone[1].IsAttackPosition = true

	state := bmgt.State{P1:p1State}
	state, err = state.Push(&bmgt.Action{PreMoveHandIndex:3, PostMoveMonsterZoneIndex:3, IsNormalSummon:true})
	if err != nil {
		panic(err)
	}

	expectedHand := bmgt.NewCardsWithPanic("サンダー・ドラゴン", "八汰烏の骸", "精神統一", "封印されし者の右腕")
	if !state.P1.Hand.Equal(expectedHand) {
		t.Errorf("テスト失敗")
	}

	expectedMonsterZone := bmgt.NewCardsWithPanic("", "封印されしエクゾディア", "", "封印されし者の左腕", "").ToMonsterZoneWithPanic()
	expectedMonsterZone[1].IsAttackPosition = true
	expectedMonsterZone[3].IsAttackPosition = true

	if !state.P1.MonsterZone.Equal(&expectedMonsterZone) {
		t.Errorf("テスト失敗")
	}
}