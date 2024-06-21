package bar_items

import (
	"log"
	"maux_bar/bar_context"

	sdlImg "github.com/veandco/go-sdl2/img"
	sdl "github.com/veandco/go-sdl2/sdl"
)

type Button struct {
	sdl.Rect
	imgPath string
	action  func()
}

func CreateButton(W, H int32, imgPath string, action func()) *Button {
	return &Button{
		Rect: sdl.Rect{
			W: W,
			H: H,
		},
		imgPath: imgPath,
		action:  action,
	}
}

func (butt *Button) Draw(surf *sdl.Surface, bar *bar_context.BarContext) {
	rw := sdl.RWFromFile(butt.imgPath, "r")
	icon, err := sdlImg.LoadPNGRW(rw)
	if butt == bar.HoveredItem {
		icon.SetAlphaMod(126)
	}

	if err != nil {
		log.Printf("Error loading the button Image (%s)\n", err.Error())
	}

	icon.BlitScaled(nil, surf, &sdl.Rect{
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

func (butt *Button) Action() {
	butt.action()
}
