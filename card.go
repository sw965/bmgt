package bmgt

import (
	"fmt"
	"github.com/sw965/omw/fn"
	omathw "github.com/sw965/omw/math"
	omws "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
	"strings"
)

type BattlePosition string

const (
	ATTACK_POSITION            = BattlePosition("攻撃表示")
	FACE_UP_DEFENSE_POSITION   = BattlePosition("表側守備表示")
	FACE_DOWN_DEFENSE_POSITION = BattlePosition("裏側守備表示")
)

type CardID int

type Card struct {
	Name                         CardName
	BattlePosition               BattlePosition
	IsSetTurn                    bool
	ThisTurnEffectActivateCounts []int
	SelectEffectNumber           int
	ID                           CardID
}

var EMPTY_CARD = Card{}

func IsEmptyCard(card Card) bool {
	return card.Name == ""
}

func IsNotEmptyCard(card Card) bool {
	return card.Name != ""
}

func CloneCard(card Card) Card {
	card.ThisTurnEffectActivateCounts = slices.Clone(card.ThisTurnEffectActivateCounts)
	return card
}

func EqualNameCard(name CardName) func(Card) bool {
	return func(card Card) bool {
		return card.Name == name
	}
}

func EqualIDCard(id CardID) func(Card) bool {
	return func(card Card) bool {
		return card.ID == id
	}
}

func IsSpellSpeed2Card(card Card) bool {
	data := CARD_DATA_BASE[card.Name]
	return data.IsQuickPlaySpell || data.IsTrap()
}

func IsMonsterCard(card Card) bool {
	data := CARD_DATA_BASE[card.Name]
	return data.IsMonster()
}

func IsLowLevelMonsterCard(card Card) bool {
	data := CARD_DATA_BASE[card.Name]
	return slices.Contains(LOW_LEVELS, data.Level)
}

func IsLevel4MonsterCard(card Card) bool {
	data := CARD_DATA_BASE[card.Name]
	return data.Level == 4
}

func IsMediumLevelMonsterCard(card Card) bool {
	data := CARD_DATA_BASE[card.Name]
	return slices.Contains(MEDIUM_LEVELS, data.Level)
}

func IsHighLevelMonsterCard(card Card) bool {
	data := CARD_DATA_BASE[card.Name]
	return data.Level > omathw.Max(MEDIUM_LEVELS...)
}

func IsSpiritMonsterCard(card Card) bool {
	data := CARD_DATA_BASE[card.Name]
	return data.IsSpiritMonster
}

func IsSpellCard(card Card) bool {
	data := CARD_DATA_BASE[card.Name]
	return data.IsTrap()
}

func IsTrapCard(card Card) bool {
	data := CARD_DATA_BASE[card.Name]
	return data.IsTrap()
}

func IsToonCard(card Card) bool {
	return strings.Contains(string(card.Name), string(TOON))
}

type Cards []Card

var OLD_LIBRARY_EXODIA_DECK = func() Cards {
	result, err := NewCards(
		"封印されしエクゾディア",
		"封印されし者の左腕",
		"封印されし者の右腕",
		"封印されし者の左足",
		"封印されし者の右足",
		"王立魔法図書館",
		"王立魔法図書館",
		"王立魔法図書館",
		"召喚僧サモンプリースト",
		"召喚僧サモンプリースト",
		"サンダー・ドラゴン",
		"サンダー・ドラゴン",
		"サンダー・ドラゴン",

		"一時休戦",
		"成金ゴブリン",
		"成金ゴブリン",
		"成金ゴブリン",
		"トゥーンのもくじ",
		"トゥーンのもくじ",
		"トゥーンのもくじ",
		"トゥーン・ワールド",
		"精神統一",
		"精神統一",
		"精神統一",
		"手札断殺",
		"手札断殺",
		"手札断殺",
		"打ち出の小槌",
		"打ち出の小槌",
		"打ち出の小槌",
		"闇の誘惑",
		"二重召喚",
		"魔法石の採掘",
		"闇の量産工場",

		"強欲な瓶",
		"強欲な瓶",
		"強欲な瓶",
		"八汰烏の骸",
		"八汰烏の骸",
		"八汰烏の骸",
	)
	if err != nil {
		panic(err)
	}
	return result
}()

func NewCards(names ...CardName) (Cards, error) {
	result := make(Cards, len(names))
	for i, name := range names {
		var card Card
		if name == "" {
			cloneCard := CloneCard(EMPTY_CARD)
			card = cloneCard
		} else {
			data, ok := CARD_DATA_BASE[name]
			if !ok {
				msg := fmt.Sprintf("データベースに存在しないカード名が入力された。入力されたカード名 = %v", name)
				return Cards{}, fmt.Errorf(msg)
			}
			card = Card{Name: name, ThisTurnEffectActivateCounts: make([]int, len(data.EffectTypes))}
		}
		result[i] = card
	}
	return result, nil
}

func (cards Cards) Draw(num int) (Cards, Cards, error) {
	drawCards := make(Cards, num)
	for i := 0; i < num; i++ {
		if len(cards) == 0 {
			return cards, drawCards, fmt.Errorf("ドローしようとしたが、カードがなかった")
		}
		var drawCard Card
		cards, drawCard = omws.Pop(cards, 0)
		drawCards[i] = drawCard
	}
	return cards, drawCards, nil
}

func (cards Cards) Clone() Cards {
	return fn.Map[Cards, Cards](cards, CloneCard)
}

func (cards Cards) NameIndices(name CardName) []int {
	return omws.IndicesFunc(cards, EqualNameCard(name))
}

func (cards Cards) IDIndex(id CardID) int {
	return slices.IndexFunc(cards, EqualIDCard(id))
}
