package bar_visual

import (
	"maux_bar/bar_items"

	sdl "github.com/veandco/go-sdl2/sdl"
)

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

func DrawItems(surface *sdl.Surface, items []bar_items.BarElement, barCtx *bar_items.BarContext) {
	for i := range items {
		var item bar_items.Draweable = items[i]
		item.Draw(surface, barCtx)
	}
}

func handleEvents(bar *bar_items.BarContext) {
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
}

func StartBar(bar *bar_items.BarContext) {
	window, err := sdl.CreateWindow(
		"maux_bar",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_UNDEFINED,
		int32(bar.Config.W), int32(bar.Config.H),
		sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN|sdl.WINDOW_SKIP_TASKBAR|sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_BORDERLESS,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surf, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	// SetItemsCentered(surf, bar.Config.Direction, bar.Elements)
	SetItemsSpaceBetween(surf, bar.Config.Direction, bar.Elements)

	go handleEvents(bar)

	/* color selection in config file */
	refreshBackgroundColor := GetBackgroundRefreshFunction(
		surf,
		sdl.Color{255, 30, 30, 126},
		sdl.Color{255, 126, 30, 30},
		0.01,
	)
	for {
		refreshBackgroundColor()
		DrawItems(surf, bar.Elements, bar)

		window.UpdateSurface()
		sdl.Delay(100)
	}
}
