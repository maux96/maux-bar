package bar_visual

import (
	"github.com/veandco/go-sdl2/sdl"
)

/*
The returned function changes the background by interpolating the two selected colors with a step size of delta.
*/
func GetBackgroundRefreshFunction(surf *sdl.Surface, col1 sdl.Color, col2 sdl.Color, delta float64) func() {
	sign := 1.0
	step := 0.0 // go from 0 to 1
	return func() {
		currentColor := sdl.Color{
			R: uint8(float64(col1.R)*(1-step) + float64(col2.R)*step),
			G: uint8(float64(col1.G)*(1-step) + float64(col2.G)*step),
			B: uint8(float64(col1.B)*(1-step) + float64(col2.B)*step),
			A: uint8(float64(col1.A)*(1-step) + float64(col2.A)*step),
		}
		nextStepValue := step + sign*delta
		if nextStepValue <= 0 || nextStepValue >= 1.0 {
			sign *= -1.0
			step += sign * delta
		} else {
			step = nextStepValue
		}
		surf.FillRect(nil, sdl.MapRGBA(surf.Format, currentColor.R, currentColor.G, currentColor.B, currentColor.A))
	}
}
