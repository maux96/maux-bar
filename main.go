package main

import (
	"fmt"
	"log"

	sdlImg "github.com/veandco/go-sdl2/img"
	sdl "github.com/veandco/go-sdl2/sdl"
	// gi "my_graphics/google_icons"
)

const WIDTH, HEIGHT int32 = 800, 200

type Button struct {
	rect    sdl.Rect
	action  func()
	imgPath string
}

var BUTTONS []Button = []Button{
	{
		sdl.Rect{X: 0, Y: 0, W: 32, H: 32},
		func() {},
		"./google_icons/star.png",
	},
	{
		sdl.Rect{X: 0, Y: 0, W: 32, H: 32},
		func() {},
		"./google_icons/star.png",
	},
	{
		sdl.Rect{X: 0, Y: 0, W: 64, H: 64},
		func() {},
		"./google_icons/apps.png",
	},
}

func DrawButtons(surface *sdl.Surface) {
	for i, button := range BUTTONS {
		rw := sdl.RWFromFile(button.imgPath, "r")
		icon, err := sdlImg.LoadPNGRW(rw)

		if err != nil {
			log.Printf("Error loading the button Image (%s)\n", err.Error())
		}

		wW, wH := surface.W, surface.H
		totalButtons := int32(len(BUTTONS))

		icon.BlitScaled(nil, surface, &sdl.Rect{
			X: int32(int32(i+1)*(wW/(totalButtons+1)) - button.rect.W/2),
			Y: int32(wH/2 - button.rect.H/2),
			W: int32(button.rect.W),
			H: int32(button.rect.H),
		})
	}
}

func GetBackgroundRefreshFunction(surf *sdl.Surface, min uint8, max uint8) func() {
	bgColor := sdl.Color{R: 0, G: 0, B: (min + max) / 2, A: 255}
	scale, sign := 2, 1
	return func() {
		nextValue := int(bgColor.B) + sign*scale
		if nextValue > int(max) {
			sign = -1
		} else if nextValue <= int(min) {
			sign = 1
		}
		bgColor.B += uint8(sign * scale)
		fmt.Println(bgColor)
		surf.FillRect(nil, bgColor.Uint32())
	}
}

func main() {
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIDTH, HEIGHT,
		sdl.WINDOW_SHOWN|sdl.WINDOW_SKIP_TASKBAR,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surf, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	refreshBackground := GetBackgroundRefreshFunction(surf, 5, 80)
	for {
		refreshBackground()
		DrawButtons(surf)

		window.UpdateSurface()
		sdl.Delay(100)
	}
}
