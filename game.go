package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameState int

const (
	StateServe GameState = iota
	StatePlay
)

type Player int

const (
	PlayerNone Player = iota
	PlayerLeft
	PlayerRight
)

type Game struct {
	state       GameState
	leftPaddle  *Paddle
	rightPaddle *Paddle
	ball        *Ball
	scoreLeft   int
	scoreRight  int
	lastWinner  Player // Tracks who won the last point
}

func NewGame() *Game {
	g := &Game{
		state:       StateServe,
		leftPaddle:  NewPaddle(50, ScreenHeight/2-50, 15, 100),
		rightPaddle: NewPaddle(ScreenWidth-50-15, ScreenHeight/2-50, 15, 100),
		ball:        NewBall(ScreenWidth/2, ScreenHeight/2, 10),
		lastWinner:  PlayerNone,
	}
	g.serve()
	return g
}

func (g *Game) serve() {
	g.ball.Reset(ScreenWidth/2, ScreenHeight/2, int(g.lastWinner))
	g.state = StateServe
}

func (g *Game) startPlay() {
	g.state = StatePlay
}

func (g *Game) Update() error {
	g.leftPaddle.UpdateLeft()

	switch g.state {
	case StateServe:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.startPlay()
		}
	case StatePlay:
		g.rightPaddle.UpdateRightCPU(g.ball)
		g.ball.Update()

		// Check collisions
		ResolveCollisions(g.ball, g.leftPaddle)
		ResolveCollisions(g.ball, g.rightPaddle)

		// Check scoring boundaries
		if g.ball.X-g.ball.Radius < 0 {
			g.scoreRight++
			g.lastWinner = PlayerRight
			g.serve()
		} else if g.ball.X+g.ball.Radius > ScreenWidth {
			g.scoreLeft++
			g.lastWinner = PlayerLeft
			g.serve()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	g.leftPaddle.Draw(screen)
	g.rightPaddle.Draw(screen)
	g.ball.Draw(screen)

	scoreStr := fmt.Sprintf("%d - %d", g.scoreLeft, g.scoreRight)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(ScreenWidth/2-50), 20)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, scoreStr, ScoreFontFace, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
