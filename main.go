package main

import (
	"fmt"
	"log"
	"maux_bar/bar_context"
	"maux_bar/bar_items"
	"os/exec"
	"syscall"

	sdl "github.com/veandco/go-sdl2/sdl"
	// gi "my_graphics/google_icons"
)

const WIDTH, HEIGHT uint32 = 800, 50

func FindIntersectItem[T bar_items.Positionable](rect *sdl.Rect, items []T) *T {
	for i := range items {
		item := items[i]
		itemRect := item.GetRect()

		if rect.HasIntersection(&itemRect) {
			return &item
		}
	}
	return nil
}

var BUTTONS = []*bar_items.Button{
	bar_items.CreateButton(
		32, 32,
		"./google_icons/star.png",
		func() {
			fmt.Println("First bar_items.Button Clicked!")
		},
	),
	bar_items.CreateButton(
		32, 32,
		"./google_icons/star.png",
		func() {
			fmt.Println("Second bar_items.Button Clicked!")
		},
	),
	bar_items.CreateButton(
		32, 32,
		"./google_icons/apps.png",
		func() {
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
	),
}

func SetItemsCentered[T bar_items.Positionable](surface *sdl.Surface, items []T) {
	for i := range items {
		item := items[i]
		wW, wH := surface.W, surface.H
		totalItems := int32(len(items))

		itemRect := item.GetRect()
		buttX := int32(int32(i+1)*(wW/(totalItems+1)) - itemRect.W/2)
		buttY := int32(wH/2 - itemRect.H/2)
		item.SetPosition(buttX, buttY)
	}
}

func DrawItems[T bar_items.Draweable](surface *sdl.Surface, items []T, barCtx *bar_context.BarContext) {
	for i := range BUTTONS {
		items[i].Draw(surface, barCtx)
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

func main() {
	bar := bar_context.New()

	window, err := sdl.CreateShapedWindow(
		"test",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
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

	SetItemsCentered(surf, BUTTONS)

	go func() {
		/* wating for mouse events */
		for {
			switch event := sdl.PollEvent().(type) {
			case *sdl.MouseButtonEvent:
				if event.Button != sdl.BUTTON_LEFT || event.Type != sdl.MOUSEBUTTONDOWN {
					continue
				}
				butt := FindIntersectItem(&sdl.Rect{X: event.X, Y: event.Y, W: 1, H: 1}, BUTTONS)
				if butt != nil {
					go (*butt).Action()
				}
			case *sdl.MouseMotionEvent:
				butt := FindIntersectItem(&sdl.Rect{X: event.X, Y: event.Y, W: 1, H: 1}, BUTTONS)
				if butt != nil {
					bar.HoveredItem = *butt
				} else {
					bar.HoveredItem = nil
				}
			}
		}
	}()

	refreshBackground := GetBackgroundRefreshFunction(surf, 5, 80)
	for {
		refreshBackground()
		DrawItems(surf, BUTTONS, bar)

		window.UpdateSurface()
		sdl.Delay(100)
	}
}
