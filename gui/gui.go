package gui

import (
	"fmt"
	"image/color"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sw965/bmgt"
	"github.com/sw965/tendon"
)

// CalcHandPositions は手札の各カードの座標を計算して返します。
// screenWidth, screenHeight: 画面サイズ
// cardWidth, cardHeight: カードの描画サイズ (スケール適用後)
// numCards: 手札の枚数
// margin: カード同士の余白（負の値を指定するとカードが重なります）
// isPlayer: trueなら自分（画面下部）、falseなら相手（画面上部）
func CalcHandPositions(screenWidth, screenHeight, cardWidth, cardHeight float64, numCards int, margin float64, isPlayer bool) [][2]float64 {
	if numCards <= 0 {
		return nil
	}

	positions := make([][2]float64, numCards)

	// 手札全体の描画幅を計算
	// 1枚目の幅 + (2枚目以降の枚数 × (カード幅 + マージン))
	totalWidth := cardWidth + float64(numCards-1)*(cardWidth+margin)

	// 中央寄せするためのX座標の開始位置
	startX := (screenWidth - totalWidth) / 2

	// Y座標の決定
	var y float64
	const verticalPadding = 20.0 // 画面上下の端からの余白

	if isPlayer {
		// 自分（画面下部に配置）
		y = screenHeight - cardHeight - verticalPadding
	} else {
		// 相手（画面上部に配置）
		y = verticalPadding
	}

	// 各カードの座標を計算
	for i := 0; i < numCards; i++ {
		x := startX + float64(i)*(cardWidth+margin)
		positions[i] = [2]float64{x, y}
	}

	return positions
}

type Game struct {
	state *bmgt.State
	Hand  tendon.Elements
}

func NewGame(state *bmgt.State) (*Game, error) {
	hand := make(tendon.Elements, len(state.Second.Hand))
	for i, card := range state.Second.Hand {
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("img/%s.jpg", card.Name.ToString()))
		if err != nil {
			return nil, err
		}
		elem := tendon.NewElement()
		elem.Image = img
		elem.SetScale(0.1) // ※スケール適用後のサイズがWidth(), Height()で取得されます
		elem.Filter = ebiten.FilterLinear
		elem.Draggable = true
		hand[i] = elem
		hand[i].Z = i

		// ドロップ時の処理（とりあえず元の実装を残す）
		hand[i].OnLeftReleased = func(self *tendon.Element) {
			// 将来的には「フィールドの判定エリアに重なっていたら召喚」などのロジックが入ります
			// self.XRelativeToParent = 100
			// self.YRelativeToParent = 100
		}
	}

	screenWidth, screenHeight := 1280.0, 720.0
	margin := -50.0

	// 自分の手札の配置
	if len(hand) > 0 {
		cardW := hand[0].Width()
		cardH := hand[0].Height()
		
		// 座標リストを取得 (isPlayer = true)
		myHandPositions := CalcHandPositions(screenWidth, screenHeight, cardW, cardH, len(hand), margin, true)

		for i, elem := range hand {
			elem.XRelativeToParent = myHandPositions[i][0]
			elem.YRelativeToParent = myHandPositions[i][1]

			// ドラッグして離した時に元の位置に戻す処理も簡単に書けます
			elem.OnLeftReleased = func(self *tendon.Element) {
				self.XRelativeToParent = myHandPositions[i][0]
				self.YRelativeToParent = myHandPositions[i][1]
			}
		}
	}

	return &Game{
		state: state,
		Hand:  hand,
	}, nil
}

func (g *Game) Update() error {
	g.Hand.Update(0, 0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 100, 30, 255})
	g.Hand.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}