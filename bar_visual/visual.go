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

func DrawItems(rend *sdl.Renderer, items []bar_items.BarElement, barCtx *bar_items.BarContext) {
	for i := range items {
		var item bar_items.Draweable = items[i]
		item.Draw(rend, barCtx)
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

	err = initTooltipCursorsFont(bar)
	if err != nil {
		panic(err)
	}

	rend, err := sdl.CreateRenderer(bar.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer rend.Destroy()

	SetItemsPlacement(
		bar.Config.PlaceItems,
		sdl.Rect{X: 0, Y: 0, W: bar.Config.W, H: bar.Config.H},
		bar.Config.Direction,
		bar.Elements,
	)

	go handleEvents(bar)

	/* color selection in config file */
	refreshBackgroundColor := GetBackgroundRefreshFunction(
		rend,
		bar.Config.Background,
	)
	for {
		refreshBackgroundColor()
		DrawItems(rend, bar.Elements, bar)

		// updateErr := window.UpdateSurface()
		// if updateErr != nil {
		// 	log.Println(updateErr.Error())
		// }

		rend.Present()
		sdl.Delay(100)
	}
}

func initTooltipCursorsFont(bar *bar_items.BarContext) (err error) {

	cursorHandPointer = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND)
	cursorNormalPointer = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW)

	fontUsed, err := ttf.OpenFont(bar.Config.Font.FontPath, bar.Config.Font.Size)
	if err != nil {
		return err
	}
	bar.Font = fontUsed

	window, err := sdl.CreateWindow("maux_bar_tooltip", 1, 1, 70, 30, sdl.WINDOW_TOOLTIP|sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_HIDDEN)
	if err != nil {
		return err
	}

	bar.TooltipWindows = window
	return nil
}
