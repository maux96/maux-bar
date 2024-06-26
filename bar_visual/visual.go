package bar_visual

import (
	"log"
	"maux_bar/bar_items"

	sdl "github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

var (
	cursorHandPointer   *sdl.Cursor
	cursorNormalPointer *sdl.Cursor
	fontUsed            *ttf.Font
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
			el := FindIntersectItem(&sdl.Rect{X: event.X, Y: event.Y, W: 1, H: 1}, bar.Elements)
			bar.HoveredItem = el
			/* change cursor */
			if el != nil {
				sdl.SetCursor(cursorHandPointer)
			} else {
				sdl.SetCursor(cursorNormalPointer)
			}

			/* tooltip */
			if butt, ok := el.(*bar_items.Button); el != nil && ok && len(butt.Label) > 0 {
				setTooltipPosition(event.X, event.Y, bar)
				setTooltipContent(butt.Label, bar)
				bar.TooltipWindows.Show()
			} else {
				bar.TooltipWindows.Hide()
			}
		case *sdl.WindowEvent:
			// nothing for now
		}
	}
}

func InitEverythingSDL() func() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	return func() {
		sdl.Quit()
		ttf.Quit()
	}
}

func StartBar(bar *bar_items.BarContext) {

	window, err := sdl.CreateWindow(
		"maux_bar",
		1, 1,
		int32(bar.Config.W), int32(bar.Config.H),
		sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_BORDERLESS|sdl.WINDOW_SKIP_TASKBAR,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	bar.Window = window

	err = initTooltipAndCursors(bar)
	if err != nil {
		panic(err)
	}

	surf, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	SetItemsPlacement(
		bar.Config.PlaceItems,
		sdl.Rect{X: 0, Y: 0, W: surf.W, H: surf.H},
		bar.Config.Direction,
		bar.Elements,
	)

	go handleEvents(bar)

	/* color selection in config file */
	refreshBackgroundColor := GetBackgroundRefreshFunction(
		surf,
		sdl.Color{R: 200, G: 30, B: 126, A: 255},
		sdl.Color{R: 30, G: 126, B: 200, A: 255},
		0.05,
	)
	for {
		refreshBackgroundColor()
		DrawItems(surf, bar.Elements, bar)

		updateErr := window.UpdateSurface()
		if updateErr != nil {
			log.Println(updateErr.Error())
		}
		sdl.Delay(100)
	}
}
