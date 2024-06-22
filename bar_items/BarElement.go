package bar_items

import (
	sdl "github.com/veandco/go-sdl2/sdl"
)

type BarElement interface {
	Draweable
	Positionable
	Actionable
}

type Draweable interface {
	Draw(surf *sdl.Surface, ctx *BarContext)
}

type Positionable interface {
	SetPosition(X, Y int32)
	GetRect() sdl.Rect
}

type Actionable interface {
	Action(ctx *BarContext)
}
