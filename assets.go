package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	ScoreFontFace *text.GoTextFace
)

func initAssets() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal("failed to load font: ", err)
	}
	ScoreFontFace = &text.GoTextFace{
		Source: s,
		Size:   48,
	}
}
