package bar_items

import (
	"log"

	sdlImg "github.com/veandco/go-sdl2/img"
	sdl "github.com/veandco/go-sdl2/sdl"
)

type Button struct {
	sdl.Rect
	Label        string
	imgAsSurface *sdl.Surface
	action       func(ctx *BarContext)
}

func NewButton(W, H int32, values map[string]string, action func(ctx *BarContext)) *Button {
	but := Button{
		Rect: sdl.Rect{
			W: W,
			H: H,
		},
		imgAsSurface: nil,
		action:       action,
	}

	if imgPath, ok := values["imgPath"]; ok {
		rw := sdl.RWFromFile(imgPath, "r")
		img, err := sdlImg.LoadPNGRW(rw)
		if err != nil {
			log.Printf("Error loading the button Image (%s)\n", err.Error())
		}
		but.imgAsSurface = img
	} else {
		panic("no imagePath in button")
	}

	if label, ok := values["label"]; ok {
		but.Label = label
	}

	return &but
}

func (butt *Button) Draw(surf *sdl.Surface, bar *BarContext) {
	if butt == bar.HoveredItem {
		butt.imgAsSurface.SetAlphaMod(126)
	} else {
		butt.imgAsSurface.SetAlphaMod(255)
	}

	butt.imgAsSurface.BlitScaled(nil, surf, &sdl.Rect{
		X: butt.X,
		Y: butt.Y,
		W: int32(butt.W),
		H: int32(butt.H),
	})
}

func (butt *Button) SetPosition(x, y int32) {
	butt.X, butt.Y = x, y
}

func (butt *Button) GetRect() sdl.Rect {
	return butt.Rect
}

func (butt *Button) Action(ctx *BarContext) {
	butt.action(ctx)
}
