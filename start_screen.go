package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type StartScreen struct {
	game         *Game
	hoveredIndex int
}

func NewStartScreen(g *Game) *StartScreen {
	return &StartScreen{game: g, hoveredIndex: -1}
}

func (s *StartScreen) Update() error {
	buttonWidth := float32(200)
	buttonHeight := float32(80)
	buttonGap := float32(20)

	labels := []string{"Easy", "Medium", "Hard"}
	numButtons := len(labels)
	totalHeight := float32(numButtons)*buttonHeight + float32(numButtons-1)*buttonGap

	startX := float32(ScreenWidth/2) - buttonWidth/2
	startY := float32(ScreenHeight/2) - totalHeight/2

	s.hoveredIndex = -1

	x, y := ebiten.CursorPosition()
	fx, fy := float32(x), float32(y)

	for i := 0; i < numButtons; i++ {
		by := startY + float32(i)*(buttonHeight+buttonGap)
		if fx >= startX && fx <= startX+buttonWidth && fy >= by && fy <= by+buttonHeight {
			s.hoveredIndex = i
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && s.hoveredIndex != -1 {
		s.game.SwitchScreen(NewPlayScreen(s.game))
	}

	touchIDs := inpututil.AppendJustPressedTouchIDs(nil)
	for _, id := range touchIDs {
		tx, ty := ebiten.TouchPosition(id)
		tfx, tfy := float32(tx), float32(ty)
		for i := 0; i < numButtons; i++ {
			by := startY + float32(i)*(buttonHeight+buttonGap)
			if tfx >= startX && tfx <= startX+buttonWidth && tfy >= by && tfy <= by+buttonHeight {
				s.game.SwitchScreen(NewPlayScreen(s.game))
			}
		}
	}

	return nil
}

func (s *StartScreen) Draw(screen *ebiten.Image) {
	// Black background
	screen.Fill(color.RGBA{0, 0, 0, 255})

	buttonWidth := float32(200)
	buttonHeight := float32(80)
	buttonGap := float32(20)

	labels := []string{"Easy", "Medium", "Hard"}
	numButtons := len(labels)
	totalHeight := float32(numButtons)*buttonHeight + float32(numButtons-1)*buttonGap

	startX := float32(ScreenWidth/2) - buttonWidth/2
	startY := float32(ScreenHeight/2) - totalHeight/2

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

		// Draw Button background (and border if not hovered)
		if isHovered {
			vector.FillRect(screen, startX, by, buttonWidth, buttonHeight, bgColor, true)
		} else {
			// Draw white border
			vector.StrokeRect(screen, startX, by, buttonWidth, buttonHeight, 2, color.White, true)
		}

		w, h := text.Measure(label, ButtonFontFace, ButtonFontFace.Size) // line spacing doesn't matter for single line

		// Calculate center offset
		textX := float64(startX) + (float64(buttonWidth)-w)/2
		textY := float64(by) + (float64(buttonHeight)-h)/2

		op := &text.DrawOptions{}
		op.GeoM.Translate(textX, textY)
		op.ColorScale.ScaleWithColor(textColor)

		text.Draw(screen, label, ButtonFontFace, op)
	}
}
