package bar_items

import (
	bar_context "maux_bar/bar_context"

	sdl "github.com/veandco/go-sdl2/sdl"
)

type Draweable interface {
	Draw(surf *sdl.Surface, ctx *bar_context.BarContext)
}
