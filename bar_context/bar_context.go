package bar_context

type BarContext struct {
	/* TODO */
	HoveredItem interface{}
}

func New() *BarContext {
	return &BarContext{}
}
