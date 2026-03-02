package main

import (
	"log"
	"math/rand/v2" //

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/bmgt"
	"github.com/sw965/bmgt/gui"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Yu-Gi-Oh! Simulator")

	// ダミーのデッキと初期状態を作成
	deck1 := make(bmgt.Cards, 40)
	deck2 := make(bmgt.Cards, 40)
	for i := 0; i < 40; i++ {
		deck1[i] = bmgt.Card{Id: i + 1, Atk: 2000, Name: bmgt.DarkMagicianGirl}
		deck2[i] = bmgt.Card{Id: i + 41, Atk: 2000, Name: bmgt.DarkMagicianGirl}
	}
	rng := rand.New(rand.NewPCG(42, 42))
	state := bmgt.NewInitState(deck1, deck2, rng)

	game, err := gui.NewGame(state)
	if err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
