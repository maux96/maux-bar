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

func NewButton(W, H int32, values map[string]any, action func(ctx *BarContext)) *Button {
	but := Button{
		Rect: sdl.Rect{
			W: W,
			H: H,
		},
		imgAsSurface: nil,
		action:       action,
	}

	if data, ok := values["imgPath"]; ok {
		if imgPath, isString := data.(string); isString {
			rw := sdl.RWFromFile(imgPath, "r")
			img, err := sdlImg.LoadPNGRW(rw)
			if err != nil {
				log.Printf("Error loading the button Image (%s)\n", err.Error())
			}
			but.imgAsSurface = img
		}
	} else {
		panic("no imagePath in button")
	}

	if data, ok := values["label"]; ok {
		if label, isString := data.(string); isString {
			but.Label = label
		}
	}

	return &but
}

func (butt *Button) Draw(rend *sdl.Renderer, bar *BarContext) {
	if butt == bar.HoveredItem {
		butt.imgAsSurface.SetAlphaMod(126)
	} else {
		butt.imgAsSurface.SetAlphaMod(255)
	}

	imgTexture, err := rend.CreateTextureFromSurface(butt.imgAsSurface)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer imgTexture.Destroy()

	rend.Copy(imgTexture, nil, &sdl.Rect{X: butt.X, Y: butt.Y, W: butt.W, H: butt.H})
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
