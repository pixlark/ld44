package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var _ = fmt.Println

const (
	pathLeftPad   int32 = 40
	pathTopPad    int32 = 100
	pathTop       int32 = pathTopPad
	pathThickness int32 = 5
	pathVertSpace int32 = screenH / 5
	pathLeft      int32 = pathLeftPad
	pathRight     int32 = screenW - pathLeftPad

	orbSize int32 = 40

	stopperSize int32 = 25

	secondsPerStep float32 = 1.0
)

func loadTexture(renderer *sdl.Renderer, path string) *sdl.Texture {
	texture, err := img.LoadTexture(renderer, "orb.png")
	if err != nil {
		fatal("Could not open texture")
	}
	return texture
}

func loadFont(path string, size int) *ttf.Font {
	font, err := ttf.OpenFont(path, size)
	if err != nil {
		fatal("Could not open font")
	}
	return font
}

type GameState struct {
	level     Level
	assets    map[string]*sdl.Texture
	font      *ttf.Font
	going     bool
	stepTimer float32

	transientResetRow    int
	transientResetColumn int
	transientTool        Tool
}

func (this *GameState) init(renderer *sdl.Renderer) {
	this.assets = make(map[string]*sdl.Texture)
	this.assets["orb"] = loadTexture(renderer, "orb.png")

	this.font = loadFont("DejaVuSans.ttf", 15)
	this.going = false

	this.level.init()
}

func (this *GameState) update(events []sdl.Event) Response {
	for _, event := range events {
		switch event := event.(type) {
		case *sdl.KeyboardEvent:
			switch event.Type {
			case sdl.KEYDOWN:
				break
			}
		}
	}

	// Update step timer if we're going
	if this.going {
		this.stepTimer -= globalState.deltaTime
		if this.stepTimer <= 0.0 {
			this.level.step()
			this.stepTimer = secondsPerStep
		}
	}

	mx, my, _ := sdl.GetMouseState()
	if !this.going && this.transientTool == nil {
		// Check for tool dragging if we're not going and not already dragging
		dragged, row, column := this.level.canDragTool()
		if dragged != nil {
			dragged.removeFromLevel(&this.level, row)
			this.transientTool = dragged
			this.transientResetRow = row
			this.transientResetColumn = column
		}
	} else if this.transientTool != nil {
		// Deal with resetting the transient tool
		reset := false
		if this.going {
			reset = true
		} else if globalState.leftClick {
			row, col, ok := this.level.inRangeOfToolSpot(mx, my)
			if ok {
				this.transientTool.addToLevel(&this.level, row, col)
				this.transientTool = nil
			} else {
				reset = true
			}
		}
		if reset {
			this.transientTool.addToLevel(&this.level, this.transientResetRow, this.transientResetColumn)
			this.transientTool = nil
		}
	}

	return Response{RESPONSE_OK, nil}
}

func fontRender(renderer *sdl.Renderer, font *ttf.Font, text string) *sdl.Texture {
	surface, _ := font.RenderUTF8Solid(text, sdl.Color{0, 0, 0, 0xff})
	defer surface.Free()
	texture, _ := renderer.CreateTextureFromSurface(surface)
	return texture
}

func (this *GameState) render(renderer *sdl.Renderer) Response {
	renderer.SetDrawColor(0, 0, 0, 0xff)
	renderer.Clear()

	// Draw path lines
	renderer.SetDrawColor(0xaa, 0xaa, 0xaa, 0xff)
	for i := range this.level.paths {
		rect := this.level.pathRect(i)
		renderer.FillRect(&rect)
	}

	// Draw stoppers
	for r, path := range this.level.paths {
		for _, stopper := range path.stoppers {
			if stopper.active {
				renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
			} else {
				renderer.SetDrawColor(0xcc, 0xcc, 0xcc, 0xff)
			}
			rect := this.level.stopperRect(r, stopper.position)
			renderer.FillRect(&rect)
		}
	}

	// Draw orbs
	for i, path := range this.level.paths {
		// Orb texture
		rect := this.level.baseRect(i, path.orbPosition)
		// Move to position along path
		// Offset to center
		rect.X -= orbSize / 2
		rect.Y -= orbSize / 2
		rect.W = orbSize
		rect.H = orbSize
		renderer.Copy(this.assets["orb"], nil, &rect)

		// Orb index
		// TODO(pixlark): Make this look less shite
		fontTexture := fontRender(renderer, this.font, fmt.Sprintf("%d", path.orbIndex))
		defer fontTexture.Destroy()
		renderer.Copy(fontTexture, nil, &rect)
	}

	// Go/Reset button
	var buttonText string
	if this.going {
		buttonText = "Reset"
	} else {
		buttonText = "Go"
	}
	pressed := button(renderer, this.font, sdl.Rect{0, 0, 100, 50}, buttonText)

	// Transient tool
	mx, my, _ := sdl.GetMouseState()
	if this.transientTool != nil {
		switch this.transientTool.(type) {
		case Stopper:
			renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
			rect := sdl.Rect{
				mx - (stopperSize / 2), my - (stopperSize / 2),
				stopperSize, stopperSize}
			renderer.FillRect(&rect)
		}
	}

	if pressed {
		this.going = !this.going
		this.stepTimer = secondsPerStep
		if !this.going {
			this.level.reset()
		}
	}

	return Response{RESPONSE_OK, nil}
}

func (this *GameState) exit() {

}
