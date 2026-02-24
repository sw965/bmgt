package bmgt

import (
	"fmt"
	"math/rand/v2"
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
	StandbyPhase
	Main1Phase
	BattlePhase
	Main2Phase
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
}

func NewInitOneSideState(deck Cards, rng *rand.Rand) OneSideState {
	state := OneSideState{
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

func (s *OneSideState) Draw() {
	if len(s.Deck) == 0 {
		return
	}
	card := s.Deck[0]
	s.Deck = s.Deck[1:]
	s.Hand = append(s.Hand, card)
}

type State struct {
	First      OneSideState
	Second     OneSideState
	Phase      Phase
	TurnPlayer TurnPlayer
	TurnCount  int
}

func NewInitState(deck1, deck2 Cards, rng *rand.Rand) *State {
	return &State{
		First:      NewInitOneSideState(deck1, rng),
		Second:     NewInitOneSideState(deck2, rng),
		Phase:      DrawPhase,
		TurnPlayer: First,
		TurnCount:  1,
	}
}

func (s *State) CurrentPlayerState() *OneSideState {
	if s.TurnPlayer == First {
		return &s.First
	}
	return &s.Second
}

func (s *State) NormalSummon(handIdx, zoneIdx int) error {
	if s.Phase != Main1Phase && s.Phase != Main2Phase {
		return fmt.Errorf("通常召喚はメインフェイズにのみ可能です")
	}

	cps := s.CurrentPlayerState()

	if handIdx < 0 || handIdx >= len(cps.Hand) {
		return fmt.Errorf("指定された手札のインデックスが無効です")
	}

	// モンスターゾーンに空きがあるかチェック
	if cps.MonsterZone[zoneIdx].Id == 0 {
		return fmt.Errorf("指定されたモンスターゾーンは空いていない")
	}

	// 手札からカードを取り出し、ゾーンに配置
	card := cps.Hand[handIdx]
	cps.MonsterZone[zoneIdx] = card

	// 手札からカードを削除
	cps.Hand = append(cps.Hand[:handIdx], cps.Hand[handIdx+1:]...)
	return nil
}

func (s *State) IsValidPhaseChange(target Phase) error {
	switch s.Phase {
	case DrawPhase:
		if target != StandbyPhase {
			return fmt.Errorf("DrawPhaseからはStandbyPhaseにのみ移行できます")
		}
	case StandbyPhase:
		if target != Main1Phase {
			return fmt.Errorf("StandbyPhaseからはMain1Phaseにのみ移行できます")
		}
	case Main1Phase:
		if target != BattlePhase && target != EndPhase {
			return fmt.Errorf("Main1PhaseからはBattlePhaseかEndPhaseにのみ移行できます")
		}
		// 先攻1ターン目はBattlePhaseに行けないルール
		if target == BattlePhase && s.TurnCount == 1 {
			return fmt.Errorf("先攻1ターン目はBattlePhaseを行えません")
		}
	case BattlePhase:
		if target != Main2Phase && target != EndPhase {
			return fmt.Errorf("BattlePhaseからはMain2PhaseかEndPhaseにのみ移行できます")
		}
	case Main2Phase:
		if target != EndPhase {
			return fmt.Errorf("Main2PhaseからはEndPhaseにのみ移行できます")
		}
	case EndPhase:
		return fmt.Errorf("EndPhaseからは直接フェイズ移行できません")
	}

	return nil
}

// ChangePhase は指定されたフェイズへ移行し、それに伴う状態更新を行います
func (s *State) ChangePhase(target Phase) error {
	// 1. まず合法性をチェック
	if err := s.IsValidPhaseChange(target); err != nil {
		return err
	}

	// 2. フェイズを更新
	s.Phase = target

	// 3. 移行先のフェイズ開始時の自動処理（ドローなど）
	if s.Phase == DrawPhase {
		// グローバルルールの設定を参照してドロー判定
		if !(s.TurnCount == 1 && s.TurnPlayer == First && !GlobalRule.IsFirstTurnDraw) {
			s.CurrentPlayerState().Draw()
		}
	}

	return nil
}

func (s *State) EndTurn() {
	s.TurnCount++
	s.TurnPlayer = s.TurnPlayer.Opposite()
	s.Phase = DrawPhase

	// 相手のドロー
	if !(s.TurnCount == 1 && s.TurnPlayer == First && !GlobalRule.IsFirstTurnDraw) {
		s.CurrentPlayerState().Draw()
	}
}

// LegalMoves は現在のフェイズと状態において、プレイヤーが選択可能な合法手のリストを返します。
func (s *State) LegalMoves() []Move {
	// キャパシティを少し多めに確保してアロケーションを減らす（MCTS向けの小技）
	moves := make([]Move, 0, 16)

	// 1. フェイズ移行（PhaseChange）の合法手
	// 移行先の候補に対して IsValidPhaseChange を呼び、エラーが nil なら追加する
	candidatePhases := []Phase{StandbyPhase, Main1Phase, BattlePhase, Main2Phase, EndPhase}
	for _, p := range candidatePhases {
		if s.IsValidPhaseChange(p) == nil {
			moves = append(moves, Move{
				Type:        PhaseChange,
				TargetPhase: p,
			})
		}
	}

	if s.Phase == Main1Phase || s.Phase == Main2Phase {
		cps := s.CurrentPlayerState()

		// モンスターゾーンに空きがあるか確認
		hasEmptyZone := false
		for i := range cps.MonsterZone {
			if cps.MonsterZone[i].Id == 0 {
				hasEmptyZone = true
				break
			}
		}

		// 空きがあれば、手札のカードを召喚するMoveを生成
		// （※本来はモンスターカードかどうかの判定が必要ですが、今はシンプルに手札全部を対象にします）
		if hasEmptyZone {
			for i := range cps.Hand {
				moves = append(moves, Move{
					Type:      NormalSummon,
					FromIndex: i,
				})
			}
		}
	}

	return moves
}

type MoveType int

const (
	NormalSummon MoveType = iota
	Attack
	PhaseChange
)

type Move struct {
	Type        MoveType
	FromIndex   int
	TargetIndex int
	TargetPhase Phase
}
