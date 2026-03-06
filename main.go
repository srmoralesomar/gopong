package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth    = 800
	ScreenHeight   = 600
	HeaderHeight   = 60
	FooterHeight   = 80
	GameAreaTop    = HeaderHeight
	GameAreaBottom = ScreenHeight - FooterHeight
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
