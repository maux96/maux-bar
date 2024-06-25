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

func setTooltipContent(content string, bar *bar_items.BarContext) {
	const PADDING = 5

	rend, err := sdl.CreateRenderer(bar.TooltipWindows, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer rend.Destroy()

	fontSurf, err := fontUsed.RenderUTF8Solid(content, sdl.Color{R: 255, G: 0, B: 0, A: 255})
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer fontSurf.Free()

	W, H := fontSurf.W, fontSurf.H
	bar.TooltipWindows.SetSize(PADDING*2+W, PADDING*2+H)

	fontTexture, err := rend.CreateTextureFromSurface(fontSurf)
	if err != nil {
		log.Println(err.Error())
		return
	}

	rend.Copy(fontTexture, nil, &sdl.Rect{X: PADDING, Y: PADDING, W: W, H: H})
	rend.Present()
}
func setTooltipPosition(x, y int32, bar *bar_items.BarContext) {
	/* TODO show the tooltip in specific direction based in window position in screen */
	X, Y := bar.Window.GetPosition()
	tW, tH := bar.TooltipWindows.GetSize()
	if bar.Config.Direction == bar_items.DIRECTION_VERTICAL {
		bar.TooltipWindows.SetPosition((10)+X+x, Y+y-(tH/2))
	} else {
		bar.TooltipWindows.SetPosition(X+x-(tW/2), Y+y-(tH+5))
	}
}

func initTooltipAndCursors(bar *bar_items.BarContext) (err error) {

	cursorHandPointer = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND)
	cursorNormalPointer = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW)

	fontUsed, err = ttf.OpenFont("./google_icons/fonts/Freedom-10eM.ttf", 16)
	if err != nil {
		return err
	}

	window, err := sdl.CreateWindow("maux_bar_tooltip", 1, 1, 70, 30, sdl.WINDOW_TOOLTIP|sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_HIDDEN)
	if err != nil {
		return err
	}

	bar.TooltipWindows = window
	return nil
}

func StartBar(bar *bar_items.BarContext) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		panic(err)
	}
	defer ttf.Quit()

	window, err := sdl.CreateWindow(
		"maux_bar",
		100, 1,
		int32(bar.Config.W), int32(bar.Config.H),
		sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_BORDERLESS|sdl.WINDOW_SKIP_TASKBAR,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	bar.Window = window

	surf, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	err = initTooltipAndCursors(bar)
	if err != nil {
		panic(err)
	}

	// SetItemsCentered(surf, bar.Config.Direction, bar.Elements)
	SetItemsSpaceBetween(surf, bar.Config.Direction, bar.Elements)

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
			log.Printf(updateErr.Error())
		}
		sdl.Delay(100)
	}
}
