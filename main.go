package main

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"

	sdlImg "github.com/veandco/go-sdl2/img"
	sdl "github.com/veandco/go-sdl2/sdl"
	// gi "my_graphics/google_icons"
)

const WIDTH, HEIGHT uint32 = 800, 50

type Button struct {
	sdl.Rect
	action  func()
	imgPath string
}

type BarState struct {
	/* TODO */
	hoveredItem interface{}
}

func FindIntersectedButton(rect *sdl.Rect) *Button {
	for i := range BUTTONS {
		butt := &BUTTONS[i]
		if butt.HasIntersection(rect) {
			return butt
		}
	}
	return nil
}

var BUTTONS []Button = []Button{
	{
		Rect: sdl.Rect{X: 0, Y: 0, W: 32, H: 32},
		action: func() {
			fmt.Println("First Button Clicked!")
		},
		imgPath: "./google_icons/star.png",
	},
	{
		Rect: sdl.Rect{X: 0, Y: 0, W: 32, H: 32},
		action: func() {
			fmt.Println("Second Button Clicked!")
		},
		imgPath: "./google_icons/star.png",
	},
	{
		Rect: sdl.Rect{X: 0, Y: 0, W: 32, H: 32},
		action: func() {
			command := exec.Command("kitty", "-e", "ranger")
			/* decople child process */
			command.SysProcAttr = &syscall.SysProcAttr{
				Setsid: true,
			}
			err := command.Start()
			if err != nil {
				log.Println("Problem executing the command:", err.Error())
			} else {
				log.Println("Command Successfully executed.")
			}
		},
		imgPath: "./google_icons/apps.png",
	},
}

func SetButtonsCenteredPosition(surface *sdl.Surface) {
	for i := range BUTTONS[:] {
		button := &BUTTONS[i]
		wW, wH := surface.W, surface.H
		totalButtons := int32(len(BUTTONS))

		buttX := int32(int32(i+1)*(wW/(totalButtons+1)) - button.W/2)
		buttY := int32(wH/2 - button.H/2)
		button.X = buttX
		button.Y = buttY
	}
}

/* TODO Draweable interface */
func DrawButtons(surface *sdl.Surface) {
	for i := range BUTTONS {
		button := &BUTTONS[i]
		rw := sdl.RWFromFile(button.imgPath, "r")
		icon, err := sdlImg.LoadPNGRW(rw)
		if button == bar.hoveredItem {
			icon.SetAlphaMod(126)
		}

		if err != nil {
			log.Printf("Error loading the button Image (%s)\n", err.Error())
		}

		icon.BlitScaled(nil, surface, &sdl.Rect{
			X: button.X,
			Y: button.Y,
			W: int32(button.W),
			H: int32(button.H),
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
		surf.FillRect(nil, bgColor.Uint32())
	}
}

var bar *BarState

func main() {
	bar = &BarState{}

	window, err := sdl.CreateShapedWindow("test", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		WIDTH, HEIGHT,
		sdl.WINDOW_SHOWN|sdl.WINDOW_SKIP_TASKBAR|sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_BORDERLESS,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surf, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	SetButtonsCenteredPosition(surf)

	go func() {
		/* wating for mouse events */
		for {
			switch event := sdl.PollEvent().(type) {
			case *sdl.MouseButtonEvent:
				if event.Button != sdl.BUTTON_LEFT || event.Type != sdl.MOUSEBUTTONDOWN {
					continue
				}

				butt := FindIntersectedButton(&sdl.Rect{X: event.X, Y: event.Y, W: 1, H: 1})
				if butt != nil {
					go butt.action()
				}
			case *sdl.MouseMotionEvent:
				butt := FindIntersectedButton(&sdl.Rect{X: event.X, Y: event.Y, W: 1, H: 1})
				if butt != nil {
					bar.hoveredItem = butt
				} else {
					bar.hoveredItem = nil
				}
			}
		}
	}()

	refreshBackground := GetBackgroundRefreshFunction(surf, 5, 80)
	for {
		refreshBackground()
		DrawButtons(surf)

		window.UpdateSurface()
		sdl.Delay(100)
	}
}
