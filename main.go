package main

import (
	"fmt"
	"maux_bar/bar_items"
	"maux_bar/config_loader"

	sdl "github.com/veandco/go-sdl2/sdl"
)

const WIDTH, HEIGHT uint32 = 800, 50

func FindIntersectItem(rect *sdl.Rect, items []bar_items.BarElement) bar_items.BarElement {
	for i := range items {
		var item bar_items.Positionable = items[i]
		itemRect := item.GetRect()

		if rect.HasIntersection(&itemRect) {
			return items[i]
		}
	}
	return nil
}

func SetItemsCentered(surface *sdl.Surface, items []bar_items.BarElement) {
	const GAP int32 = 10

	totalSpaceX := GAP*int32(len(items)) - 1
	for _, item := range items {
		rect := item.GetRect()
		W, _ := rect.W, rect.H
		totalSpaceX += W
	}

	surfW, surfH := surface.W, surface.H
	startPosX := surfW/2 - totalSpaceX/2
	for _, item := range items {
		itemRect := item.GetRect()
		buttX := startPosX
		buttY := int32(surfH/2 - itemRect.H/2)

		startPosX += itemRect.W + GAP
		item.SetPosition(buttX, buttY)
	}
}

func SetItemsSpaceBetween(surface *sdl.Surface, items []bar_items.BarElement) {
	for i := range items {
		var item bar_items.Positionable = items[i]
		wW, wH := surface.W, surface.H
		totalItems := int32(len(items))

		itemRect := item.GetRect()
		buttX := int32(int32(i+1)*(wW/(totalItems+1)) - itemRect.W/2)
		buttY := int32(wH/2 - itemRect.H/2)
		item.SetPosition(buttX, buttY)
	}
}

func DrawItems(surface *sdl.Surface, items []bar_items.BarElement, barCtx *bar_items.BarContext) {
	for i := range items {
		var item bar_items.Draweable = items[i]
		item.Draw(surface, barCtx)
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
	bar := bar_items.NewBarContext()
	config_loader.PrepareBar(bar, "./testConfig.json")
	fmt.Println(bar)

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

	SetItemsCentered(surf, bar.Elements)

	go func() {
		/* wating for mouse events */
		for {
			switch event := sdl.PollEvent().(type) {
			case *sdl.MouseButtonEvent:
				if event.Button != sdl.BUTTON_LEFT || event.Type != sdl.MOUSEBUTTONDOWN {
					continue
				}
				butt := FindIntersectItem(&sdl.Rect{X: event.X, Y: event.Y, W: 1, H: 1}, bar.Elements)
				if butt != nil {
					go butt.Action(bar)
				}
			case *sdl.MouseMotionEvent:
				butt := FindIntersectItem(&sdl.Rect{X: event.X, Y: event.Y, W: 1, H: 1}, bar.Elements)
				if butt != nil {
					bar.HoveredItem = butt
				} else {
					bar.HoveredItem = nil
				}
			}
		}
	}()

	refreshBackground := GetBackgroundRefreshFunction(surf, 5, 230)
	for {
		refreshBackground()
		DrawItems(surf, bar.Elements, bar)

		window.UpdateSurface()
		sdl.Delay(100)
	}
}
