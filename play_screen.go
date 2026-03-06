package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
		leftPaddle:  NewPaddle(50, float64(GameAreaTop+(GameAreaBottom-GameAreaTop)/2-50), 15, 100),
		rightPaddle: NewPaddle(ScreenWidth-50-15, float64(GameAreaTop+(GameAreaBottom-GameAreaTop)/2-50), 15, 100),
		ball:        NewBall(ScreenWidth/2, float64(GameAreaTop+(GameAreaBottom-GameAreaTop)/2), 10),
		lastWinner:  PlayerNone,
	}
	p.serve()
	return p
}

func (p *PlayScreen) serve() {
	p.ball.Reset(ScreenWidth/2, float64(GameAreaTop+(GameAreaBottom-GameAreaTop)/2), int(p.lastWinner))
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

	// Retro colors
	bgHeader := color.RGBA{0, 0, 128, 255}    // Dark Blue
	bgFooter1 := color.RGBA{128, 0, 128, 255} // Purple
	bgFooter2 := color.RGBA{0, 100, 0, 255}   // Dark Green
	textHeader := color.RGBA{0, 255, 255, 255}
	textFooter1 := color.RGBA{255, 182, 193, 255} // Light Pink
	textFooter2 := color.RGBA{144, 238, 144, 255} // Light Green

	// Draw Backgrounds
	vector.FillRect(screen, 0, 0, float32(ScreenWidth), float32(GameAreaTop), bgHeader, true)
	vector.FillRect(screen, 0, float32(GameAreaBottom), float32(ScreenWidth), 40, bgFooter1, true)
	vector.FillRect(screen, 0, float32(GameAreaBottom)+40, float32(ScreenWidth), 40, bgFooter2, true)

	// Draw Header - Score (now using ButtonFontFace)
	scoreStr := fmt.Sprintf("%d - %d", p.scoreLeft, p.scoreRight)
	op := &text.DrawOptions{}
	wScore, _ := text.Measure(scoreStr, ButtonFontFace, ButtonFontFace.Size)
	op.GeoM.Translate((float64(ScreenWidth)-wScore)/2, 20)
	op.ColorScale.ScaleWithColor(textHeader)
	text.Draw(screen, scoreStr, ButtonFontFace, op)

	// Draw Header - Players
	headerStrLeft := "You"
	opHL := &text.DrawOptions{}
	opHL.GeoM.Translate(50, 20)
	opHL.ColorScale.ScaleWithColor(textHeader)
	text.Draw(screen, headerStrLeft, ButtonFontFace, opHL)

	headerStrRight := "CPU"
	opHR := &text.DrawOptions{}
	wHR, _ := text.Measure(headerStrRight, ButtonFontFace, ButtonFontFace.Size)
	opHR.GeoM.Translate(float64(ScreenWidth)-50-wHR, 20)
	opHR.ColorScale.ScaleWithColor(textHeader)
	text.Draw(screen, headerStrRight, ButtonFontFace, opHR)

	// Draw Footer: Row 1 (State)
	stateStr := "On"
	if p.state == StateServe {
		stateStr = "Serve"
	}
	opS := &text.DrawOptions{}
	wS, _ := text.Measure(stateStr, ButtonFontFace, ButtonFontFace.Size)
	opS.GeoM.Translate((float64(ScreenWidth)-wS)/2, float64(GameAreaBottom)+10)
	opS.ColorScale.ScaleWithColor(textFooter1)
	text.Draw(screen, stateStr, ButtonFontFace, opS)

	// Draw Footer: Row 2 (Controls)
	ctrlLeft := "Up/Down = move"
	ctrlCenter := "Space = serve"
	ctrlRight := "Q = quit"

	opCL := &text.DrawOptions{}
	opCL.GeoM.Translate(50, float64(GameAreaBottom)+50)
	opCL.ColorScale.ScaleWithColor(textFooter2)
	text.Draw(screen, ctrlLeft, ButtonFontFace, opCL)

	wC, _ := text.Measure(ctrlCenter, ButtonFontFace, ButtonFontFace.Size)
	opCC := &text.DrawOptions{}
	opCC.GeoM.Translate((float64(ScreenWidth)-wC)/2, float64(GameAreaBottom)+50)
	opCC.ColorScale.ScaleWithColor(textFooter2)
	text.Draw(screen, ctrlCenter, ButtonFontFace, opCC)

	wR, _ := text.Measure(ctrlRight, ButtonFontFace, ButtonFontFace.Size)
	opCR := &text.DrawOptions{}
	opCR.GeoM.Translate(float64(ScreenWidth)-50-wR, float64(GameAreaBottom)+50)
	opCR.ColorScale.ScaleWithColor(textFooter2)
	text.Draw(screen, ctrlRight, ButtonFontFace, opCR)
}
