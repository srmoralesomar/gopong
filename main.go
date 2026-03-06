package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

func main() {
	initAssets()
	game := NewGame()

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Pong")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
