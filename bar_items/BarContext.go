package bar_items

type BarContext struct {
	/* TODO */
	Config      *BarConfigData
	HoveredItem BarElement
	Elements    []BarElement
}

type BarDirection string

const (
	DIRECTION_HORIZONTAL BarDirection = "horizontal"
	DIRECTION_VERTICAL   BarDirection = "vertical"
)

type BarConfigData struct {
	W          int32
	H          int32
	PlaceItems string
	Direction  BarDirection
	Objects    []RawBarObjectData
}

func NewBarContext() *BarContext {
	return &BarContext{}
}
