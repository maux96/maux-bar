package bar_items

import "github.com/veandco/go-sdl2/sdl"

type BarContext struct {
	/* TODO */
	Config         *BarConfigData
	HoveredItem    BarElement
	Window         *sdl.Window
	TooltipWindows *sdl.Window
	Elements       []BarElement
}

type BarDirection string
type BarItemsPlacement string

const (
	DIRECTION_HORIZONTAL BarDirection = "horizontal"
	DIRECTION_VERTICAL   BarDirection = "vertical"

	PLACE_ITEMS_CENTER        BarItemsPlacement = "center"
	PLACE_ITEMS_SPACE_BETWEEN BarItemsPlacement = "space-between"
)

type BarConfigData struct {
	W          int32
	H          int32
	PlaceItems BarItemsPlacement
	Direction  BarDirection
	Objects    []RawBarObjectData
	Background BackgroundConfig
}

type BackgroundConfig struct {
	Type   string
	Values map[string]string
}

func NewBarContext() *BarContext {
	return &BarContext{}
}
