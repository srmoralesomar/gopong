package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type StartScreen struct {
	game      *Game
	isHovered bool
}

func NewStartScreen(g *Game) *StartScreen {
	return &StartScreen{game: g}
}

func (s *StartScreen) Update() error {
	buttonWidth := float32(200)
	buttonHeight := float32(80)
	buttonX := float32(ScreenWidth/2) - buttonWidth/2
	buttonY := float32(ScreenHeight/2) - buttonHeight/2

	// Check hover state
	x, y := ebiten.CursorPosition()
	fx, fy := float32(x), float32(y)
	s.isHovered = fx >= buttonX && fx <= buttonX+buttonWidth && fy >= buttonY && fy <= buttonY+buttonHeight

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && s.isHovered {
		s.game.SwitchScreen(NewPlayScreen(s.game))
	}

	touchIDs := inpututil.AppendJustPressedTouchIDs(nil)
	for _, id := range touchIDs {
		tx, ty := ebiten.TouchPosition(id)
		tfx, tfy := float32(tx), float32(ty)
		if tfx >= buttonX && tfx <= buttonX+buttonWidth && tfy >= buttonY && tfy <= buttonY+buttonHeight {
			s.game.SwitchScreen(NewPlayScreen(s.game))
		}
	}

	return nil
}

func (s *StartScreen) Draw(screen *ebiten.Image) {
	// Black background
	screen.Fill(color.RGBA{0, 0, 0, 255})

	buttonWidth := float32(200)
	buttonHeight := float32(80)
	buttonX := float32(ScreenWidth/2) - buttonWidth/2
	buttonY := float32(ScreenHeight/2) - buttonHeight/2

	// Button colors based on hover
	var bgColor, textColor color.Color
	if s.isHovered {
		bgColor = color.White
		textColor = color.Black
	} else {
		bgColor = color.Black
		textColor = color.White
	}

	// Draw Button background (and border if not hovered)
	if s.isHovered {
		vector.FillRect(screen, buttonX, buttonY, buttonWidth, buttonHeight, bgColor, true)
	} else {
		// Draw white border
		vector.StrokeRect(screen, buttonX, buttonY, buttonWidth, buttonHeight, 2, color.White, true)
	}

	msg := "START"
	w, h := text.Measure(msg, ButtonFontFace, ButtonFontFace.Size) // line spacing doesn't matter for single line

	// Calculate center offset
	textX := float64(buttonX) + (float64(buttonWidth) - w) / 2
	textY := float64(buttonY) + (float64(buttonHeight) - h) / 2

	op := &text.DrawOptions{}
	op.GeoM.Translate(textX, textY)
	op.ColorScale.ScaleWithColor(textColor)
	
	text.Draw(screen, msg, ButtonFontFace, op)
}
