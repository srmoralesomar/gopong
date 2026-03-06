package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type EndScreen struct {
	game         *Game
	scoreLeft    int
	scoreRight   int
	isWin        bool
	difficulty   Difficulty
	hoveredIndex int
}

func NewEndScreen(g *Game, scoreLeft, scoreRight int, isWin bool, diff Difficulty) *EndScreen {
	return &EndScreen{
		game:         g,
		scoreLeft:    scoreLeft,
		scoreRight:   scoreRight,
		isWin:        isWin,
		difficulty:   diff,
		hoveredIndex: -1,
	}
}

func (s *EndScreen) Update() error {
	buttonWidth := float32(200)
	buttonHeight := float32(80)
	buttonGap := float32(20)

	labels := []string{"Play again", "Home"}
	numButtons := len(labels)
	totalButtonsHeight := float32(numButtons)*buttonHeight + float32(numButtons-1)*buttonGap

	// Determine vertical offset logic to match Draw
	totalContentHeight := float32(100) + totalButtonsHeight // Roughly 100px for title and score
	startY := float32(ScreenHeight/2) - totalContentHeight/2 + 100 // Buttons start below text

	startX := float32(ScreenWidth/2) - buttonWidth/2

	s.hoveredIndex = -1

	x, y := ebiten.CursorPosition()
	fx, fy := float32(x), float32(y)

	for i := 0; i < numButtons; i++ {
		by := startY + float32(i)*(buttonHeight+buttonGap)
		if fx >= startX && fx <= startX+buttonWidth && fy >= by && fy <= by+buttonHeight {
			s.hoveredIndex = i
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		switch s.hoveredIndex {
		case 0: // Play again
			s.game.SwitchScreen(NewPlayScreen(s.game, s.difficulty))
		case 1: // Home
			s.game.SwitchScreen(NewStartScreen(s.game))
		}
	}

	touchIDs := inpututil.AppendJustPressedTouchIDs(nil)
	for _, id := range touchIDs {
		tx, ty := ebiten.TouchPosition(id)
		tfx, tfy := float32(tx), float32(ty)
		for i := 0; i < numButtons; i++ {
			by := startY + float32(i)*(buttonHeight+buttonGap)
			if tfx >= startX && tfx <= startX+buttonWidth && tfy >= by && tfy <= by+buttonHeight {
				switch i {
				case 0: // Play again
					s.game.SwitchScreen(NewPlayScreen(s.game, s.difficulty))
				case 1: // Home
					s.game.SwitchScreen(NewStartScreen(s.game))
				}
			}
		}
	}

	return nil
}

func (s *EndScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Draw Title ("You won!" or "You lost")
	title := "You lost"
	titleColor := color.RGBA{255, 0, 0, 255} // Red
	if s.isWin {
		title = "You won!"
		titleColor = color.RGBA{0, 255, 0, 255} // Green
	}

	opTitle := &text.DrawOptions{}
	wTitle, hTitle := text.Measure(title, ButtonFontFace, ButtonFontFace.Size)
	
	totalContentHeight := float64(100) + float64(2*80+20) // approx height for text + 2 buttons + gap
	baseY := (float64(ScreenHeight) - totalContentHeight) / 2

	opTitle.GeoM.Translate((float64(ScreenWidth)-wTitle)/2, baseY)
	opTitle.ColorScale.ScaleWithColor(titleColor)
	text.Draw(screen, title, ButtonFontFace, opTitle)

	// Draw Score
	scoreStr := fmt.Sprintf("Final Score: %d - %d", s.scoreLeft, s.scoreRight)
	opScore := &text.DrawOptions{}
	wScore, _ := text.Measure(scoreStr, ButtonFontFace, ButtonFontFace.Size)
	opScore.GeoM.Translate((float64(ScreenWidth)-wScore)/2, baseY+hTitle+20) // 20px gap below title
	opScore.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, scoreStr, ButtonFontFace, opScore)

	// Draw Buttons
	buttonWidth := float32(200)
	buttonHeight := float32(80)
	buttonGap := float32(20)

	labels := []string{"Play again", "Home"}

	startX := float32(ScreenWidth/2) - buttonWidth/2
	startY := float32(baseY) + 100 // Same offset as used in Update

	for i, label := range labels {
		by := startY + float32(i)*(buttonHeight+buttonGap)
		isHovered := s.hoveredIndex == i

		var bgColor, textColor color.Color
		if isHovered {
			bgColor = color.White
			textColor = color.Black
		} else {
			bgColor = color.Black
			textColor = color.White
		}

		if isHovered {
			vector.FillRect(screen, startX, by, buttonWidth, buttonHeight, bgColor, true)
		} else {
			vector.StrokeRect(screen, startX, by, buttonWidth, buttonHeight, 2, color.White, true)
		}

		w, h := text.Measure(label, ButtonFontFace, ButtonFontFace.Size)

		textX := float64(startX) + (float64(buttonWidth)-w)/2
		textY := float64(by) + (float64(buttonHeight)-h)/2

		op := &text.DrawOptions{}
		op.GeoM.Translate(textX, textY)
		op.ColorScale.ScaleWithColor(textColor)

		text.Draw(screen, label, ButtonFontFace, op)
	}
}
