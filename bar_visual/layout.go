package bar_visual

import (
	"log"
	"maux_bar/bar_items"

	sdl "github.com/veandco/go-sdl2/sdl"
)

func SetItemsPlacement(placementMode bar_items.BarItemsPlacement, containerRect sdl.Rect, direction bar_items.BarDirection, items []bar_items.BarElement) {
	switch placementMode {
	case bar_items.PLACE_ITEMS_CENTER:
		SetItemsCentered(containerRect, direction, items)
	case bar_items.PLACE_ITEMS_SPACE_BETWEEN:
		SetItemsSpaceBetween(containerRect, direction, items)
	default:
		log.Println("unknown placement. using center")
		SetItemsCentered(containerRect, direction, items)
	}
}

func SetItemsCentered(containerRect sdl.Rect, direction bar_items.BarDirection, items []bar_items.BarElement) {
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
	containerW, containerH := containerRect.W, containerRect.H

	var startPos int32
	if direction == bar_items.DIRECTION_HORIZONTAL {
		startPos = containerRect.X + containerW/2 - totalSpace/2
	} else {
		startPos = containerRect.Y + containerH/2 - totalSpace/2
	}
	for _, item := range items {
		itemRect := item.GetRect()
		var X, Y int32
		if direction == bar_items.DIRECTION_HORIZONTAL {
			X = startPos
			Y = containerRect.Y + int32(containerH/2-itemRect.H/2)
			startPos += itemRect.W + GAP
		} else {
			X = containerRect.X + int32(containerW/2-itemRect.W/2)
			Y = startPos
			startPos += itemRect.H + GAP
		}
		item.SetPosition(X, Y)
	}
}

func SetItemsSpaceBetween(containerRect sdl.Rect, direction bar_items.BarDirection, items []bar_items.BarElement) {
	W, H := containerRect.W, containerRect.H
	totalItems := int32(len(items))

	for i := range items {
		var item bar_items.Positionable = items[i]
		var posX, posY int32
		itemRect := item.GetRect()
		if direction == bar_items.DIRECTION_HORIZONTAL {
			posX = containerRect.X + int32(int32(i+1)*(W/(totalItems+1))-itemRect.W/2)
			posY = containerRect.Y + int32(H/2-itemRect.H/2)
		} else if direction == bar_items.DIRECTION_VERTICAL {
			posX = containerRect.X + int32(W/2-itemRect.W/2)
			posY = containerRect.Y + int32(int32(i+1)*(H/(totalItems+1))-itemRect.H/2)
		}
		item.SetPosition(posX, posY)
	}

}
