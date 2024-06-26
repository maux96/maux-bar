package bar_visual

import (
	"errors"
	"log"
	"maux_bar/bar_items"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

/* format: R G B A? */
func stringToColor(colAsString string) (col *sdl.Color, err error) {
	stringList := strings.Split(colAsString, " ")
	if len(stringList) < 3 {
		return nil, errors.New("incorrect color format")
	}

	R, err0 := strconv.ParseUint(stringList[0], 10, 8)
	G, err1 := strconv.ParseUint(stringList[1], 10, 8)
	B, err2 := strconv.ParseUint(stringList[2], 10, 8)
	err = errors.Join(err0, err1, err2)
	if err != nil {
		return nil, err
	}

	col = &sdl.Color{
		R: uint8(R),
		G: uint8(G),
		B: uint8(B),
		A: 255,
	}
	if len(stringList) == 4 {
		A, err := strconv.ParseUint(stringList[2], 10, 8)
		if err != nil {
			return nil, err
		}
		col.A = uint8(A)
	}
	return col, nil
}

/*
The returned function changes the background by interpolating the two selected colors with a step size of delta.
*/
func GetBackgroundRefreshFunction(surf *sdl.Surface, config bar_items.BackgroundConfig) func() {
	var (
		col1  *sdl.Color
		col2  *sdl.Color
		delta float64

		sign = 1.0
		step = 0.0 // go from 0 to 1
	)
	var errCol1, errCol2, errDelta error
	if colAsString, ok := config.Values["col1"]; ok {
		col1, errCol1 = stringToColor(colAsString)
	} else {
		col1 = &sdl.Color{R: 200, G: 30, B: 126, A: 255}
	}
	if colAsString, ok := config.Values["col2"]; ok {
		col2, errCol2 = stringToColor(colAsString)
	} else {
		col2 = &sdl.Color{R: 30, G: 126, B: 200, A: 255}
	}
	if deltaAsString, ok := config.Values["delta"]; ok {
		delta, errDelta = strconv.ParseFloat(deltaAsString, 64)
	} else {
		delta = 0.05
	}
	err := errors.Join(errCol1, errCol2, errDelta)
	if err != nil {
		log.Fatalln(err.Error())
	}

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
