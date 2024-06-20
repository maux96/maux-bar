package main

import (
	"log"

	sdlImg "github.com/veandco/go-sdl2/img"
	sdl "github.com/veandco/go-sdl2/sdl"
	// gi "my_graphics/google_icons"
)

const WIDTH, HEIGHT int32 = 800, 600

type Button struct {
	rect    sdl.Rect
	action  func()
	imgPath string
}

var BUTTONS []Button = []Button{
	{
		sdl.Rect{X: 0, Y: 0, W: 64, H: 64},
		func() {},
		"./google_icons/star.png",
	},
	{
		sdl.Rect{X: 0, Y: 0, W: 64, H: 64},
		func() {},
		"./google_icons/apps.png",
	},
	{
		sdl.Rect{X: 0, Y: 0, W: 126, H: 126},
		func() {},
		"./google_icons/star.png",
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
			X: int32(int32(i)*(wW/totalButtons) - button.rect.W/2),
			Y: int32(wH/2 - button.rect.H/2),
			W: int32(button.rect.W),
			H: int32(button.rect.H),
		})
	}
}

func main() {
	// err := gi.GetPngIcon("grade")

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

	var colorValue uint32 = 0xFF00000F
	for {
		surf.FillRect(nil, colorValue)
		colorValue++

		DrawButtons(surf)

		window.UpdateSurface()
		sdl.Delay(100)
	}
}
