package bar_items

import sdl "github.com/veandco/go-sdl2/sdl"

type Positionable interface {
	SetPosition(X, Y int32)
	GetRect() sdl.Rect
}
