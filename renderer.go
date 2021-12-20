package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func DrawText(renderer *sdl.Renderer, font *Font, text string, rect *sdl.Rect, color sdl.Color) {
	surface, err := font.Data.RenderUTF8Blended(text, color)
	checkError(err)
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	checkError(err)
	defer texture.Destroy()

	renderer.Copy(texture, nil, rect)
}

func DrawRect(renderer *sdl.Renderer, rect *sdl.Rect, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.FillRect(rect)
}

func DrawRectOutline(renderer *sdl.Renderer, rect *sdl.Rect, color sdl.Color, width int32) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)

	top := sdl.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: width}
	right := sdl.Rect{X: rect.X + rect.W - width, Y: rect.Y, W: width, H: rect.H}
	bottom := sdl.Rect{X: rect.X, Y: rect.Y + rect.H - width, W: rect.W, H: width}
	left := sdl.Rect{X: rect.X, Y: rect.Y, W: width, H: rect.H}

	renderer.FillRect(&top)
	renderer.FillRect(&right)
	renderer.FillRect(&bottom)
	renderer.FillRect(&left)
}

func DrawRectTransparent(renderer *sdl.Renderer, rect *sdl.Rect, color sdl.Color) {
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	DrawRect(renderer, rect, color)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_NONE)
}

func DrawImage(renderer *sdl.Renderer, texture *sdl.Texture, rect sdl.Rect, color sdl.Color) {
	texture.SetColorMod(color.R, color.G, color.B)
	renderer.Copy(texture, nil, &rect)
}
