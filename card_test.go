package bmgt_test

import (
	"github.com/sw965/bmgt"
	"golang.org/x/exp/slices"
	"testing"
)

func TestCards_Draw(t *testing.T) {
	cards := bmgt.NewCards(bmgt.EXODIA_PART_NAMES...)
	cards, drawCards, err := cards.Draw(2)
	if err != nil {
		panic(err)
	}

	expectedCards := bmgt.NewCards("封印されし者の右腕", "封印されし者の左足", "封印されし者の右足")
	expectedDrawCards := bmgt.NewCards("封印されしエクゾディア", "封印されし者の左腕")

	if !slices.Equal(cards, expectedCards) {
		t.Errorf("テスト失敗")
	}

	if !slices.Equal(drawCards, expectedDrawCards) {
		t.Errorf("テスト失敗")
	}
}