package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	checkError(err)
	defer sdl.Quit()

	err = ttf.Init()
	checkError(err)
	defer ttf.Quit()

	window, err := sdl.CreateWindow("Karuba Tile Randomizer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_RESIZABLE)
	checkError(err)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	checkError(err)
	defer renderer.Destroy()

	windowWidth, windowHeight := window.GetSize()

	app := NewApp(renderer, windowWidth, windowHeight)
	input := Input{}

	running := true
	for running {
		input.Clear()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				keycode := t.Keysym.Sym

				switch keycode {
				case sdl.K_LCTRL:
					fallthrough
				case sdl.K_RCTRL:
					input.Ctrl = t.Type == sdl.KEYDOWN
				case sdl.K_SPACE:
					// @TODO (!important) should not set true on key repeat
					if t.State != sdl.RELEASED {
						input.Space = true
					}
				}
			case *sdl.WindowEvent:
				if t.Event == sdl.WINDOWEVENT_RESIZED {
					app.Resize(t.Data1, t.Data2)
				}
			}
		}

		app.Tick(&input)
		app.Render()
	}

	app.Close()
}
