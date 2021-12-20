package main

import (
	"math/rand"
	"path"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type App struct {
	WindowRect sdl.Rect

	Tiles           []string
	ActiveTileIndex int32
	ActiveTile      Image
	RemainingTiles  int32

	Renderer *sdl.Renderer
}

func NewApp(renderer *sdl.Renderer, windowWidth int32, windowHeight int32) (result *App) {
	result = &App{}

	result.WindowRect = sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight}

	allTiles := ReadDirectory("assets/images")
	result.Tiles = make([]string, len(allTiles))
	for index, entry := range allTiles {
		result.Tiles[index] = path.Join("assets/images", entry.Name())
	}
	result.RemainingTiles = int32(len(result.Tiles)) - 1

	result.Renderer = renderer

	rand.Seed(time.Now().UnixNano())
	result.shuffleTiles()

	result.loadTile()

	return
}

func (app *App) Close() {
}

func (app *App) shuffleTiles() {
	currentIndex := len(app.Tiles)

	for currentIndex != 0 {
		randomIndex := rand.Intn(currentIndex)
		currentIndex -= 1

		app.Tiles[currentIndex], app.Tiles[randomIndex] = app.Tiles[randomIndex], app.Tiles[currentIndex]
	}
}

func (app *App) loadTile() {
	// @TODO maybe cache tiles?
	app.ActiveTile = LoadImage(app.Tiles[app.ActiveTileIndex], app.Renderer)
}

func (app *App) reset() {
	app.ActiveTileIndex = 0
	app.RemainingTiles = int32(len(app.Tiles)) - 1

	app.shuffleTiles()
	app.loadTile()
}

func (app *App) Resize(windowWidth int32, windowHeight int32) {
	app.WindowRect.W = windowWidth
	app.WindowRect.H = windowHeight
}

func (app *App) Tick(input *Input) {
	if input.Space {
		if input.Ctrl {
			// @TODO (!important) hold for some amount of time to reset
			app.reset()
		} else if app.RemainingTiles > 0 {
			app.ActiveTileIndex += 1
			app.RemainingTiles -= 1

			app.loadTile()
		}
	}
}

func (app *App) Render() {
	app.Renderer.SetDrawColor(0, 0, 0, 255)
	app.Renderer.Clear()

	tilePosition := sdl.Point{X: app.WindowRect.W/2 - app.ActiveTile.Width/2, Y: app.WindowRect.H/2 - app.ActiveTile.Height/2}
	app.ActiveTile.Render(app.Renderer, tilePosition, sdl.Color{R: 255, G: 255, B: 255, A: 255})

	// @TODO (!important) show remaining tiles
	// @TODO (!important) show silver border for tiles containing silver
	// @TODO (!important) show gold border for tiles containing gold
	// @TODO (!important) show instructions

	app.Renderer.Present()
}
