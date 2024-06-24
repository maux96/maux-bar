package bar_visual

import (
	"maux_bar/bar_items"

	sdl "github.com/veandco/go-sdl2/sdl"
)

func SetItemsCentered(surface *sdl.Surface, direction bar_items.BarDirection, items []bar_items.BarElement) {
	const GAP int32 = 10

	totalSpace := GAP*int32(len(items)) - 1
	for _, item := range items {
		rect := item.GetRect()
		if direction == bar_items.DIRECTION_HORIZONTAL {
			totalSpace += rect.W
		} else {
			totalSpace += rect.H
		}
	}
	surfW, surfH := surface.W, surface.H

	var startPos int32
	if direction == bar_items.DIRECTION_HORIZONTAL {
		startPos = surfW/2 - totalSpace/2
	} else {
		startPos = surfH/2 - totalSpace/2
	}
	for _, item := range items {
		itemRect := item.GetRect()
		var X, Y int32
		if direction == bar_items.DIRECTION_HORIZONTAL {
			X = startPos
			Y = int32(surfH/2 - itemRect.H/2)
			startPos += itemRect.W + GAP
		} else {
			X = int32(surfW/2 - itemRect.W/2)
			Y = startPos
			startPos += itemRect.H + GAP
		}
		item.SetPosition(X, Y)
	}
}

func SetItemsSpaceBetween(surface *sdl.Surface, direction bar_items.BarDirection, items []bar_items.BarElement) {
	W, H := surface.W, surface.H
	totalItems := int32(len(items))

	for i := range items {
		var item bar_items.Positionable = items[i]
		var posX, posY int32
		itemRect := item.GetRect()
		if direction == bar_items.DIRECTION_HORIZONTAL {
			posX = int32(int32(i+1)*(W/(totalItems+1)) - itemRect.W/2)
			posY = int32(H/2 - itemRect.H/2)
		} else if direction == bar_items.DIRECTION_VERTICAL {
			posX = int32(W/2 - itemRect.W/2)
			posY = int32(int32(i+1)*(H/(totalItems+1)) - itemRect.H/2)
		}
		item.SetPosition(posX, posY)
	}

}
