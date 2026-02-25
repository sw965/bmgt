package main

import (
	"fmt"
	"math/rand/v2"
	"runtime"

	"github.com/sw965/bmgt"
	"github.com/sw965/crow/game"
	"github.com/sw965/crow/game/sequential"
	"github.com/sw965/crow/pucb"
	"github.com/sw965/crow/mcts/puct"
	"github.com/sw965/omw/mathx/randx"
)

func main() {
	// --- 1. 基本設定 ---
	numGames := 100
	mctsSimulations := 16000
	numWorkers := runtime.NumCPU() // 並列実行数
	masterRng := rand.New(rand.NewPCG(42, 42))

	// --- 2. 共通ロジックの定義 ---
	logic := sequential.Logic[*bmgt.State, bmgt.Move, bmgt.TurnPlayer]{
		LegalMovesFunc: func(s *bmgt.State) []bmgt.Move { return s.LegalMoves() },
		MoveFunc: func(s *bmgt.State, m bmgt.Move) (*bmgt.State, error) {
			next := s.Clone()
			if err := next.Move(m); err != nil {
				return nil, err
			}
			return next, nil
		},
		EqualFunc: func(s1, s2 *bmgt.State) bool { return s1.Equal(s2) },
		CurrentAgentFunc: func(s *bmgt.State) bmgt.TurnPlayer { return s.TurnPlayer },
	}

	engine := sequential.Engine[*bmgt.State, bmgt.Move, bmgt.TurnPlayer]{
		Logic:           logic,
		RankByAgentFunc: bmgt.RankByAgentFunc,
		Agents:          []bmgt.TurnPlayer{bmgt.First, bmgt.Second},
	}
	engine.SetStandardResultScoreByAgentFunc()

	// --- 3. アクターの作成 ---
	// A. ランダムアクター
	randomActor := sequential.NewRandomActorCritic[*bmgt.State, bmgt.Move, bmgt.TurnPlayer]()
	randomActor.Name = "Random"

	// B. MCTSアクターの設定
	puctEng := puct.Engine[*bmgt.State, bmgt.Move, bmgt.TurnPlayer]{
		Game:         engine,
		PUCBFunc:     pucb.NewAlphaGoFunc(1.414), // C=sqrt(2)
		NextNodesCap: 10,
		VirtualValue: 0.0, // 仮想損失（並列探索用）
	}
	puctEng.SetUniformPolicyFunc() // 方策は一様分布
	
	// リーフ評価はランダムプレイアウトで行う
	playoutRng := rand.New(rand.NewPCG(1, 1))
	puctEng.SetPlayout(randomActor, playoutRng)

	// 探索用のワーカー乱数生成
	workerRngs := randx.NewPCGs(1)
	
	mctsActor := sequential.ActorCritic[*bmgt.State, bmgt.Move, bmgt.TurnPlayer]{
		Name: "MCTS-16000",
		// 探索結果（訪問回数）に基づく方策と、根ノードの評価値を返す
		PolicyValueFunc: puctEng.NewPolicyValueFunc(mctsSimulations, workerRngs),
		// 探索後、最も訪問回数が多い手を選択（Greedy）
		SelectFunc: game.MaxSelectFunc[bmgt.Move, bmgt.TurnPlayer],
	}

	// --- 4. 対局の実行 (Cross Playout) ---
	// 100戦行うために、初期状態を100個用意（すべてBMG 40枚デッキ）
	inits := make([]*bmgt.State, numGames)
	for i := 0; i < numGames; i++ {
		deck1 := createBMGDeck()
		deck2 := createBMGDeck()
		inits[i] = bmgt.NewInitState(deck1, deck2, masterRng)
	}

	// 2つのアクター（Random vs MCTS）を登録
	recorder, _ := engine.NewCrossPlayoutRecorder(inits, []sequential.ActorCritic[*bmgt.State, bmgt.Move, bmgt.TurnPlayer]{
		randomActor,
		mctsActor,
	}, numWorkers)

	fmt.Printf("対局開始: %s vs %s (%d試合)\n", randomActor.Name, mctsActor.Name, numGames*2) // 先後入れ替え含むため2倍

	// 全対局を回す
	_, err := recorder.Collect()
	if err != nil {
		panic(err)
	}

	// --- 5. 結果出力 ---
	avgScores, _ := recorder.AverageScoreByActorCriticName()
	numWins := recorder.NumGamesByActorCriticName()

	fmt.Println("\n=== 最終集計結果 ===")
	for name, score := range avgScores {
		fmt.Printf("アクター: %-12s | 平均スコア: %.4f | 対局数: %d\n", name, score, numWins[name])
	}
}

// BMGデッキ作成用 (Id=0を避ける)
func createBMGDeck() bmgt.Cards {
	deck := make(bmgt.Cards, 40)
	for i := 0; i < 40; i++ {
		deck[i] = bmgt.Card{Name: 1, Atk: 2000, Def: 1700, Id: i + 1}
	}
	return deck
}