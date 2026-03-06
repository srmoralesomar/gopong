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

type PlayScreen struct {
	game        *Game
	state       GameState
	leftPaddle  *Paddle
	rightPaddle *Paddle
	ball        *Ball
	scoreLeft   int
	scoreRight  int
	lastWinner  Player
}

func NewPlayScreen(g *Game) *PlayScreen {
	p := &PlayScreen{
		game:        g,
		state:       StateServe,
		leftPaddle:  NewPaddle(50, ScreenHeight/2-50, 15, 100),
		rightPaddle: NewPaddle(ScreenWidth-50-15, ScreenHeight/2-50, 15, 100),
		ball:        NewBall(ScreenWidth/2, ScreenHeight/2, 10),
		lastWinner:  PlayerNone,
	}
	p.serve()
	return p
}

func (p *PlayScreen) serve() {
	p.ball.Reset(ScreenWidth/2, ScreenHeight/2, int(p.lastWinner))
	p.state = StateServe
}

func (p *PlayScreen) startPlay() {
	p.state = StatePlay
}

func (p *PlayScreen) Update() error {
	p.leftPaddle.UpdateLeft()

	switch p.state {
	case StateServe:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			p.startPlay()
		}
	case StatePlay:
		p.rightPaddle.UpdateRightCPU(p.ball)
		p.ball.Update()

		// Check collisions
		ResolveCollisions(p.ball, p.leftPaddle)
		ResolveCollisions(p.ball, p.rightPaddle)

		// Check scoring boundaries
		if p.ball.X-p.ball.Radius < 0 {
			p.scoreRight++
			p.lastWinner = PlayerRight
			p.serve()
		} else if p.ball.X+p.ball.Radius > ScreenWidth {
			p.scoreLeft++
			p.lastWinner = PlayerLeft
			p.serve()
		}
	}

	return nil
}

func (p *PlayScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	p.leftPaddle.Draw(screen)
	p.rightPaddle.Draw(screen)
	p.ball.Draw(screen)

	scoreStr := fmt.Sprintf("%d - %d", p.scoreLeft, p.scoreRight)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(ScreenWidth/2-50), 20)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, scoreStr, ScoreFontFace, op)
}
