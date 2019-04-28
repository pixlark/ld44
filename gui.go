package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var _ = fmt.Println

func button(renderer *sdl.Renderer, font *ttf.Font, rect sdl.Rect, text string) bool {
	// Render bg
	renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
	renderer.FillRect(&rect)

	// Render font
	texture := fontRender(renderer, font, text, sdl.Color{0, 0, 0, 0xff})
	defer texture.Destroy()
	
	_, _, fW, fH, _ := texture.Query()
	fontRect := centerRectInRect(sdl.Rect{0, 0, fW, fH}, rect)
	
	renderer.Copy(texture, nil, &fontRect)

	// Check for click
	mx, my, _ := sdl.GetMouseState()
	if (globalState.leftClick) {
		if mx > rect.X && mx < rect.X+rect.W &&
			my > rect.Y && my < rect.Y+rect.H {
			return true
		}
	}
	return false
}
