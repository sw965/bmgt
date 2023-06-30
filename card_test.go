package bmgt_test

import (
	"github.com/sw965/bmgt"
	"testing"
)

func TestCards_Draw(t *testing.T) {
	cards := bmgt.NewCardsWithPanic(bmgt.EXODIA_PART_NAMES...)
	resultCards, drawCards, err := cards.Draw(2)
	if err != nil {
		panic(err)
	}

	expectedCards := bmgt.NewCardsWithPanic("封印されし者の右腕", "封印されし者の左足", "封印されし者の右足")
	expectedDrawCards := bmgt.NewCardsWithPanic("封印されしエクゾディア", "封印されし者の左腕")

	if !resultCards.Equal(expectedCards) {
		t.Errorf("テスト失敗")
	}

	if !drawCards.Equal(expectedDrawCards) {
		t.Errorf("テスト失敗")
	}
}