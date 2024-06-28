package bar_items

import (
	sdl "github.com/veandco/go-sdl2/sdl"
)

type RawBarObjectData struct {
	ObjType string
	W       int32
	H       int32
	Values  map[string]any
}

type BarElement interface {
	Draweable
	Positionable
	Actionable
}

type Draweable interface {
	Draw(rend *sdl.Renderer, ctx *BarContext)
}

type Positionable interface {
	SetPosition(X, Y int32)
	GetRect() sdl.Rect
}

type Actionable interface {
	Action(ctx *BarContext)
}
