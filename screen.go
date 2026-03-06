package main

import "github.com/hajimehoshi/ebiten/v2"

type Screen interface {
	Update() error
	Draw(screen *ebiten.Image)
}
