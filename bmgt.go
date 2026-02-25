package bmgt

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

type Rule struct {
	IsFirstTurnDraw bool
}

var GlobalRule = Rule{
	IsFirstTurnDraw: false,
}

type LifePoint int

const (
	InitLifePoint = 8000
)

type TurnPlayer int

const (
	First TurnPlayer = iota
	Second
)

func (t TurnPlayer) Opposite() TurnPlayer {
	if t == First {
		return Second
	}
	return First
}

type Phase int

const (
	DrawPhase Phase = iota
	// StandbyPhase // 現時点で目指す実装ではスタンバイフェイズは不要
	Main1Phase
	BattlePhase
	// Main2Phase // 現時点で目指す実装では不要
	EndPhase
)

type OneSideState struct {
	Deck             Cards
	Hand             Cards
	MonsterZone      [5]Card
	SpellAndTrapZone [5]Card
	Graveyard        Cards
	Banish           Cards
	LifePoint        LifePoint
	IsDeckOut        bool
}

func NewInitOneSideState(deck Cards, rng *rand.Rand) *OneSideState {
	state := &OneSideState{
		Deck:      deck,
		Hand:      make(Cards, 0, 40),
		LifePoint: InitLifePoint,
	}

	rng.Shuffle(len(state.Deck), func(i, j int) {
		state.Deck[i], state.Deck[j] = state.Deck[j], state.Deck[i]
	})

	for i := 0; i < 5; i++ {
		state.Draw()
	}
	return state
}

func (s *OneSideState) Clone() *OneSideState {
	cloned := *s
	cloned.Deck = append(Cards(nil), s.Deck...)
	cloned.Hand = append(Cards(nil), s.Hand...)
	cloned.Graveyard = append(Cards(nil), s.Graveyard...)
	cloned.Banish = append(Cards(nil), s.Banish...)
	return &cloned
}

func (s *OneSideState) Equal(other *OneSideState) bool {
	if s == nil || other == nil {
		return s == other
	}

	if s.LifePoint != other.LifePoint || s.IsDeckOut != other.IsDeckOut {
		return false
	}

	if s.MonsterZone != other.MonsterZone || s.SpellAndTrapZone != other.SpellAndTrapZone {
		return false
	}

	return slices.Equal(s.Deck, other.Deck) &&
		slices.Equal(s.Hand, other.Hand) &&
		slices.Equal(s.Graveyard, other.Graveyard) &&
		slices.Equal(s.Banish, other.Banish)
}

func (s *OneSideState) Draw() {
	if len(s.Deck) == 0 {
		s.IsDeckOut = true
		return
	}
	card := s.Deck[0]
	s.Deck = s.Deck[1:]
	s.Hand = append(s.Hand, card)
}

type State struct {
	First      *OneSideState
	Second     *OneSideState
	Phase      Phase
	TurnPlayer TurnPlayer
	TurnCount  int
}

func NewInitState(deck1, deck2 Cards, rng *rand.Rand) *State {
	return &State{
		First:      NewInitOneSideState(deck1, rng),
		Second:     NewInitOneSideState(deck2, rng),
		Phase:      Main1Phase,
		TurnPlayer: First,
		TurnCount:  1,
	}
}

func (s *State) Clone() *State {
	cloned := *s
	cloned.First = s.First.Clone()
	cloned.Second = s.Second.Clone()
	return &cloned
}

func (s *State) Equal(other *State) bool {
	if s == nil || other == nil {
		return s == other
	}
	if s.Phase != other.Phase || s.TurnPlayer != other.TurnPlayer || s.TurnCount != other.TurnCount {
		return false
	}
	return s.First.Equal(other.First) && s.Second.Equal(other.Second)
}

func (s *State) TurnPlayerState() *OneSideState {
	if s.TurnPlayer == First {
		return s.First
	}
	return s.Second
}

func (s *State) NonTurnPlayerState() *OneSideState {
	if s.TurnPlayer == First {
		return s.Second
	}
	return s.First
}

func (s *State) NormalSummon(handIdx, zoneIdx int) error {
	if s.Phase != Main1Phase { // && s.Phase != Main2Phase {
		return fmt.Errorf("通常召喚はメインフェイズにのみ可能です")
	}

	tps := s.TurnPlayerState()

	if handIdx < 0 || handIdx >= len(tps.Hand) {
		return fmt.Errorf("指定された手札のインデックスが無効です")
	}

	// モンスターゾーンに空きがあるかチェック
	if tps.MonsterZone[zoneIdx].Id != 0 {
		return fmt.Errorf("指定されたモンスターゾーンは空いていない")
	}

	// 手札からカードを取り出し、ゾーンに配置
	card := tps.Hand[handIdx]
	tps.MonsterZone[zoneIdx] = card

	// 手札からカードを削除
	tps.Hand = append(tps.Hand[:handIdx], tps.Hand[handIdx+1:]...)
	return nil
}

func (s *State) Battle(fromIdx, targetIdx int) error {
	tps := s.TurnPlayerState()
	// MonsterCard -> MonsterZone に修正
	if tps.MonsterZone[fromIdx].Id == 0 {
		return fmt.Errorf("空のモンスターゾーンで攻撃しようとした")
	}

	ntps := s.NonTurnPlayerState()
	if ntps.MonsterZone[targetIdx].Id == 0 {
		return fmt.Errorf("空のモンスターゾーンに対して攻撃しようとした")
	}

	attackCard := tps.MonsterZone[fromIdx]
	defenseCard := tps.MonsterZone[targetIdx]

	diff := attackCard.Atk - defenseCard.Atk
	if diff > 0 {
		// 1. 相手モンスターを破壊し、墓地へ送る
		ntps.Graveyard = append(ntps.Graveyard, defenseCard)
		ntps.MonsterZone[targetIdx] = Card{} // ゾーンを空（初期値）にする

		// 相手のライフポイントを減らす
		ntps.LifePoint -= LifePoint(diff)
	} else if diff == 0 {
		// 2. 相打ち: 両方のモンスターを破壊し、墓地へ送る
		tps.Graveyard = append(tps.Graveyard, attackCard)
		tps.MonsterZone[fromIdx] = Card{}

		ntps.Graveyard = append(ntps.Graveyard, defenseCard)
		ntps.MonsterZone[targetIdx] = Card{}
	} else {
		// 3. 攻撃の失敗 (diff < 0): 自分のモンスターが破壊される
		tps.Graveyard = append(tps.Graveyard, attackCard)
		tps.MonsterZone[fromIdx] = Card{}

		// 自分のライフポイントを減らす（diffがマイナスなので -diff にする）
		tps.LifePoint -= LifePoint(-diff)
	}

	// モンスターが相打ち・自爆特攻で破壊されていなければ、攻撃済みフラグを立てる
	if tps.MonsterZone[fromIdx].Id != 0 {
		tps.MonsterZone[fromIdx].IsAttacked = true
	}
	return nil
}

func (s *State) DirectAttack(fromIdx int) error {
	tps := s.TurnPlayerState()

	if tps.MonsterZone[fromIdx].Id == 0 {
		return fmt.Errorf("空のモンスターゾーンでダイレクトアタックしようとした")
	}

	attackCard := tps.MonsterZone[fromIdx]
	ntps := s.NonTurnPlayerState()
	ntps.LifePoint -= LifePoint(attackCard.Atk)
	tps.MonsterZone[fromIdx].IsAttacked = true
	return nil
}

func (s *State) LegalMoves() []Move {
	moves := make([]Move, 0, 16)
	if s.Phase == Main1Phase {
		moves = append(moves, Move{
			Type:        PhaseChange,
			TargetPhase: BattlePhase,
		})

		moves = append(moves, Move{
			Type:        PhaseChange,
			TargetPhase: EndPhase,
		})
	}

	if s.Phase == BattlePhase {
		moves = append(moves, Move{
			Type:        PhaseChange,
			TargetPhase: EndPhase,
		})
	}

	if s.Phase == Main1Phase { // || s.Phase == Main2Phase {
		tps := s.TurnPlayerState()
		for fromIdx := range tps.Hand {
			if tps.Hand[fromIdx].Id == 0 {
				continue
			}
			for targetIdx := range tps.MonsterZone {
				if tps.MonsterZone[targetIdx].Id != 0 {
					continue
				}
				moves = append(moves, Move{
					Type:        NormalSummon,
					FromIndex:   fromIdx,
					TargetIndex: targetIdx,
				})
			}
		}
	}

	if s.Phase == BattlePhase {
		tps := s.TurnPlayerState()
		ntps := s.NonTurnPlayerState()

		// 相手のモンスターゾーンが空かどうかを判定する
		isOpponentFieldEmpty := true
		for _, card := range ntps.MonsterZone {
			if card.Id != 0 {
				isOpponentFieldEmpty = false
				break
			}
		}

		for fromIdx := range tps.MonsterZone {
			if tps.MonsterZone[fromIdx].Id == 0 {
				continue
			}

			if tps.MonsterZone[fromIdx].IsAttacked {
				continue
			}

			if isOpponentFieldEmpty {
				// 相手の場が空ならダイレクトアタックを追加
				moves = append(moves, Move{
					Type:      DirectAttack,
					FromIndex: fromIdx,
				})
			} else {
				// 相手の場にモンスターがいるなら、各モンスターへの攻撃を追加
				for targetIdx := range ntps.MonsterZone {
					if ntps.MonsterZone[targetIdx].Id == 0 {
						continue
					}

					moves = append(moves, Move{
						Type:        Attack,
						FromIndex:   fromIdx,
						TargetIndex: targetIdx,
					})
				}
			}
		}
	}
	return moves
}

func (s *State) Move(move Move) error {
	switch move.Type {
	case NormalSummon:
		return s.NormalSummon(move.FromIndex, move.TargetIndex)
	case Attack:
		return s.Battle(move.FromIndex, move.TargetIndex)
	case DirectAttack:
		return s.DirectAttack(move.FromIndex)
	case PhaseChange:
		s.Phase = move.TargetPhase
		if s.Phase == EndPhase {
			tps := s.TurnPlayerState()
			for i := range tps.MonsterZone {
				tps.MonsterZone[i].IsAttacked = false
			}

			s.TurnCount++
			s.TurnPlayer = s.TurnPlayer.Opposite()

			nextTps := s.TurnPlayerState()
			nextTps.Draw()
			s.Phase = Main1Phase
		}
	}
	return nil
}

type MoveType int

const (
	NormalSummon MoveType = iota
	Attack
	DirectAttack
	PhaseChange
)

type Move struct {
	Type        MoveType
	FromIndex   int
	TargetIndex int
	TargetPhase Phase
}
