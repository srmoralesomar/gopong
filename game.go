package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	currentScreen Screen
}

func NewGame() *Game {
	g := &Game{}
	g.currentScreen = NewStartScreen(g)
	return g
}

func (g *Game) SwitchScreen(s Screen) {
	g.currentScreen = s
}

func (g *Game) Update() error {
	if g.currentScreen != nil {
		return g.currentScreen.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScreen != nil {
		g.currentScreen.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
