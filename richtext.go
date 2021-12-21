package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type RichTextPiece struct {
	Text  string
	Color sdl.Color
	Width int32
}

type RichText struct {
	Pieces []RichTextPiece

	Font *Font
}

func NewRichText(font *Font) *RichText {
	return &RichText{
		Pieces: []RichTextPiece{},
		Font:   font,
	}
}

func (rt *RichText) Add(text string, color sdl.Color) {
	rt.Pieces = append(rt.Pieces, RichTextPiece{
		Text:  text,
		Color: color,
		Width: rt.Font.GetStringWidth(text),
	})
}

func (rt *RichText) Render(renderer *sdl.Renderer, position sdl.Point) {
	cursor := int32(0)
	for _, piece := range rt.Pieces {
		rect := sdl.Rect{
			X: position.X + cursor,
			Y: position.Y,
			W: piece.Width,
			H: rt.Font.Size,
		}
		DrawText(renderer, rt.Font, piece.Text, &rect, piece.Color)

		cursor += piece.Width + int32(rt.Font.CharacterWidth)
	}
}
