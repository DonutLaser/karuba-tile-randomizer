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

	HoldStartTime   uint32
	ResetTime       uint32
	CurrentTime     uint32
	ResetInProgress bool
	DisableReset    bool

	RestartInstruction  RichText
	NextTileInstruction RichText
}

func NewApp(renderer *sdl.Renderer, windowWidth int32, windowHeight int32) (result *App) {
	result = &App{}

	result.WindowRect = sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight}

	allTiles := ReadDirectory("assets/images/tiles")
	result.Tiles = make([]string, len(allTiles))
	for index, entry := range allTiles {
		result.Tiles[index] = path.Join("assets/images/tiles", entry.Name())
	}
	result.RemainingTiles = int32(len(result.Tiles)) - 1
	result.RemainingText = fmt.Sprintf("Remaining: %d", result.RemainingTiles)

	result.Font = LoadFont("assets/fonts/consolab.ttf", 14)
	result.Renderer = renderer

	result.ResetTime = 2200 // 2.2 seconds

	result.RestartInstruction = *NewRichText(&result.Font)
	result.RestartInstruction.Add("Hold", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	result.RestartInstruction.Add("R", sdl.Color{R: 238, G: 204, B: 117, A: 255})
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

func (app *App) GetIcon() *sdl.Surface {
	return LoadIcon("assets/images/icon.png")
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
	if input.Space && app.RemainingTiles > 0 {
		app.ActiveTileIndex += 1
		app.RemainingTiles -= 1
		app.RemainingText = fmt.Sprintf("Remaining: %d", app.RemainingTiles)

		app.loadTile()

		return
	}

	if input.R == KeyStateWentDown {
		if !app.DisableReset {
			app.HoldStartTime = sdl.GetTicks()
			app.CurrentTime = app.HoldStartTime
			app.ResetInProgress = true
		}
	} else if input.R == KeyStateDown {
		if !app.DisableReset {
			app.CurrentTime = sdl.GetTicks()
			if app.CurrentTime >= app.HoldStartTime+app.ResetTime {
				app.reset()
				app.DisableReset = true
				app.ResetInProgress = false
			}
		}
	} else if input.R == KeyStateWentUp {
		app.DisableReset = false
		app.ResetInProgress = false
	}
}

func (app *App) Render() {
	app.Renderer.SetDrawColor(0, 0, 0, 255)
	app.Renderer.Clear()

	if app.ResetInProgress {
		progressRect := sdl.Rect{
			X: app.WindowRect.X,
			Y: app.WindowRect.Y,
			W: int32((float64(app.CurrentTime) - float64(app.HoldStartTime)) / float64(app.ResetTime) * float64(app.WindowRect.W)),
			H: app.WindowRect.H,
		}
		DrawRectTransparent(app.Renderer, &progressRect, sdl.Color{R: 238, G: 204, B: 117, A: 125})
	}

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
