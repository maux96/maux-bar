package bar_items

type BarContext struct {
	/* TODO */
	HoveredItem BarElement
	Elements    []BarElement
}

func NewBarContext() *BarContext {
	return &BarContext{}
}
