package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	InitialBallSpeed = 4.0
	MaxBallSpeed     = 12.0
	SpeedIncrement   = 0.5
)

type Ball struct {
	X      float64
	Y      float64
	Radius float64
	Vx     float64
	Vy     float64
	Speed  float64
}

func NewBall(x, y, radius float64) *Ball {
	b := &Ball{
		X:      x,
		Y:      y,
		Radius: radius,
	}
	b.Reset(x, y)
	return b
}

func (b *Ball) Reset(x, y float64) {
	b.X = x
	b.Y = y
	b.Speed = InitialBallSpeed

	// Start at a random angle either to the left or right, bounded safely
	// We want roughly an angle between -45 deg (-pi/4) and +45 deg (pi/4)
	angle := (rand.Float64() * math.Pi / 2) - (math.Pi / 4)
	
	// Randomly flip to the other side (leftwards)
	if rand.Intn(2) == 0 {
		angle += math.Pi
	}

	b.Vx = math.Cos(angle)
	b.Vy = math.Sin(angle)
}

func (b *Ball) Update() {
	b.X += b.Vx * b.Speed
	b.Y += b.Vy * b.Speed

	// Top and bottom screen bounds collisions
	if b.Y-b.Radius < 0 {
		b.Y = b.Radius // Positional correction
		b.Vy = math.Abs(b.Vy) // Ensure positive Y velocity
	} else if b.Y+b.Radius > ScreenHeight {
		b.Y = ScreenHeight - b.Radius // Positional correction
		b.Vy = -math.Abs(b.Vy) // Ensure negative Y velocity
	}
}

func (b *Ball) Draw(screen *ebiten.Image) {
	// Draw a white circle
	vector.FillCircle(screen, float32(b.X), float32(b.Y), float32(b.Radius), color.White, true)
}

func (b *Ball) IncreaseSpeed() {
	b.Speed += SpeedIncrement
	if b.Speed > MaxBallSpeed {
		b.Speed = MaxBallSpeed
	}
}
