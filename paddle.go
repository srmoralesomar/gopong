package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const PaddleSpeed = 6.0

type Paddle struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func NewPaddle(x, y, w, h float64) *Paddle {
	return &Paddle{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
	}
}

func (p *Paddle) UpdateLeft() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.Y -= PaddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.Y += PaddleSpeed
	}

	// Stay within screen bounds
	if p.Y < 0 {
		p.Y = 0
	}
	if p.Y+p.Height > ScreenHeight {
		p.Y = ScreenHeight - p.Height
	}
}

func (p *Paddle) UpdateRightCPU(ball *Ball) {
	// Simple AI: tries to center the paddle on the ball's Y position
	targetY := ball.Y - p.Height/2
	
	// Limit CPU speed slightly so you can beat it at high speeds
	cpuSpeedY := PaddleSpeed * 0.85
	
	if p.Y < targetY-cpuSpeedY {
		p.Y += cpuSpeedY
	} else if p.Y > targetY+cpuSpeedY {
		p.Y -= cpuSpeedY
	} else {
		p.Y = targetY // Snap to target if very close
	}

	if p.Y < 0 {
		p.Y = 0
	}
	if p.Y+p.Height > ScreenHeight {
		p.Y = ScreenHeight - p.Height
	}
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	// Draw a white rectangle
	vector.FillRect(screen, float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height), color.White, true)
}
