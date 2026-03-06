package main

import (
	"image/color"
	"math"

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

	// Stay within bounds
	if p.Y < float64(GameAreaTop) {
		p.Y = float64(GameAreaTop)
	}
	if p.Y+p.Height > float64(GameAreaBottom) {
		p.Y = float64(GameAreaBottom) - p.Height
	}
}

func (p *Paddle) UpdateRightCPU(ball *Ball, diff Difficulty) {
	var targetY float64
	var cpuSpeedY float64

	switch diff {
	case DifficultyEasy:
		targetY = ball.Y - p.Height/2
		cpuSpeedY = PaddleSpeed * 0.5
	case DifficultyMedium:
		targetY = ball.Y - p.Height/2
		cpuSpeedY = PaddleSpeed * 0.85
	case DifficultyHard:
		cpuSpeedY = PaddleSpeed * 1.0
		if ball.Vx <= 0 {
			// Center paddle if ball is moving away
			targetY = float64(GameAreaTop+(GameAreaBottom-GameAreaTop)/2) - p.Height/2
		} else {
			// Predict intersection
			vX := ball.Vx * ball.Speed
			vY := ball.Vy * ball.Speed
			
			// Time until collision with paddle front edge
			timeToIntercept := (p.X - ball.Radius - ball.X) / vX
			if timeToIntercept < 0 {
				targetY = ball.Y - p.Height/2
			} else {
				predictedY := ball.Y + vY*timeToIntercept
				
				top := float64(GameAreaTop) + ball.Radius
				bottom := float64(GameAreaBottom) - ball.Radius
				rangeY := bottom - top
				
				if rangeY > 0 {
					relativeY := predictedY - top
					bounces := int(math.Floor(relativeY / rangeY))
					rem := relativeY - float64(bounces)*rangeY
					
					if bounces%2 == 0 {
						predictedY = top + rem
					} else {
						predictedY = bottom - rem
					}
				}
				targetY = predictedY - p.Height/2
			}
		}
	}

	if p.Y < targetY-cpuSpeedY {
		p.Y += cpuSpeedY
	} else if p.Y > targetY+cpuSpeedY {
		p.Y -= cpuSpeedY
	} else {
		p.Y = targetY // Snap to target if very close
	}

	if p.Y < float64(GameAreaTop) {
		p.Y = float64(GameAreaTop)
	}
	if p.Y+p.Height > float64(GameAreaBottom) {
		p.Y = float64(GameAreaBottom) - p.Height
	}
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	// Draw a white rectangle
	vector.FillRect(screen, float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height), color.White, true)
}
