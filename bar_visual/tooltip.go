package bar_visual

import (
	"log"
	"maux_bar/bar_items"

	"github.com/veandco/go-sdl2/sdl"
)

func setTooltipContent(content string, bar *bar_items.BarContext) {
	const PADDING = 5

	rend, err := sdl.CreateRenderer(bar.TooltipWindows, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer rend.Destroy()

	fontSurf, err := bar.Font.RenderUTF8Solid(content, sdl.Color{R: 200, G: 200, B: 200, A: 200})
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
	defer fontTexture.Destroy()
	rend.SetDrawColor(0, 0, 0, 255)
	rend.Clear()

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
