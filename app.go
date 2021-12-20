package main

import (
	"fmt"
	"math/rand"
	"path"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type App struct {
	WindowRect sdl.Rect

	Tiles           []string
	ActiveTileIndex int32
	ActiveTile      Image
	RemainingTiles  int32
	RemainingText   string

	Font     Font
	Renderer *sdl.Renderer

	RestartInstruction  RichText
	NextTileInstruction RichText
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
	result.RemainingText = fmt.Sprintf("Remaining: %d", result.RemainingTiles)

	result.Font = LoadFont("assets/fonts/consolab.ttf", 14)
	result.Renderer = renderer

	result.RestartInstruction = *NewRichText(&result.Font)
	result.RestartInstruction.Add("Press", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	result.RestartInstruction.Add("Ctrl + Space", sdl.Color{R: 238, G: 204, B: 117, A: 255})
	result.RestartInstruction.Add("to restart", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	result.NextTileInstruction = *NewRichText(&result.Font)
	result.NextTileInstruction.Add("Press", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	result.NextTileInstruction.Add("Space", sdl.Color{R: 238, G: 204, B: 117, A: 255})
	result.NextTileInstruction.Add("to get the next tile", sdl.Color{R: 255, G: 255, B: 255, A: 255})

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
	app.RemainingText = fmt.Sprintf("Remaining: %d", app.RemainingTiles)

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
			app.RemainingText = fmt.Sprintf("Remaining: %d", app.RemainingTiles)

			app.loadTile()
		}
	}
}

func (app *App) Render() {
	app.Renderer.SetDrawColor(0, 0, 0, 255)
	app.Renderer.Clear()

	tilePosition := sdl.Point{X: app.WindowRect.W/2 - app.ActiveTile.Width/2, Y: app.WindowRect.H/2 - app.ActiveTile.Height/2}
	app.ActiveTile.Render(app.Renderer, tilePosition, sdl.Color{R: 255, G: 255, B: 255, A: 255})

	tileRect := sdl.Rect{
		X: tilePosition.X - 5,
		Y: tilePosition.Y - 5,
		W: app.ActiveTile.Width + 10,
		H: app.ActiveTile.Height + 10,
	}

	if strings.Contains(app.Tiles[app.ActiveTileIndex], "silver") {
		DrawRectOutline(app.Renderer, &tileRect, sdl.Color{R: 216, G: 232, B: 232, A: 255}, 5)
	} else if strings.Contains(app.Tiles[app.ActiveTileIndex], "gold") {
		DrawRectOutline(app.Renderer, &tileRect, sdl.Color{R: 237, G: 239, B: 93, A: 255}, 5)
	}

	if app.RemainingTiles > 0 {
		remainingWidth := app.Font.GetStringWidth(app.RemainingText)
		remainingRect := sdl.Rect{
			X: app.WindowRect.X + (app.WindowRect.W/2 - remainingWidth/2),
			Y: tilePosition.Y + app.ActiveTile.Height + 20,
			W: remainingWidth,
			H: app.Font.Size,
		}
		DrawText(app.Renderer, &app.Font, app.RemainingText, &remainingRect, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	}

	app.NextTileInstruction.Render(app.Renderer, sdl.Point{X: app.WindowRect.X + 10, Y: app.WindowRect.Y + 10})
	app.RestartInstruction.Render(app.Renderer, sdl.Point{X: app.WindowRect.X + 10, Y: app.WindowRect.Y + 10 + app.Font.Size + 10})

	app.Renderer.Present()
}
