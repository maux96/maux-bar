package bar_items

import (
	"io"
	"log"
	"maux_bar/utils"
	"os/exec"
	"strings"
	"time"

	sdl "github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Outputer struct {
	sdl.Rect
	output    <-chan string
	lastPrint string
	font      *ttf.Font
}

func NewOutputer(W, H int32, values map[string]any) *Outputer {
	but := Outputer{
		Rect: sdl.Rect{
			W: W,
			H: H,
		},
		lastPrint: "-",
	}

	if fontPath, ok := values["fontPath"]; ok {
		var fontSize float64 = 16
		if data, ok0 := values["fontSize"]; ok0 {
			fontSize = data.(float64)
		}
		fontUsed, err := ttf.OpenFont(fontPath.(string), int(fontSize))
		if err != nil {
			log.Fatalln(err.Error())
		}
		but.font = fontUsed
	} else {
		log.Println("font not found in outputer")
	}

	if action, ok := values["action"]; ok {
		isRepetitive := false
		repet, ok0 := values["repetitive"]
		if ok0 {
			isRepetitive = repet.(bool)
		}
		commandArgs, err := utils.ConvertSliceTo[string](action.([]any))
		if err != nil {
			log.Fatalln(err.Error())
		}
		ch, err := executeCommand(commandArgs, isRepetitive)
		if err != nil {
			log.Fatalln(err.Error())
		} else {
			but.output = ch
		}
	} else {
		log.Fatalln("no action found in outputer")
	}

	return &but
}

func (butt *Outputer) Draw(rend *sdl.Renderer, bar *BarContext) {
	var textToPrint string
	select {
	case x := <-butt.output:
		textToPrint = x
		butt.lastPrint = x
	default:
		textToPrint = butt.lastPrint
	}

	var usedFont *ttf.Font
	if butt.font != nil {
		usedFont = butt.font
	} else {
		usedFont = bar.Font
	}

	fontSurf, err := usedFont.RenderUTF8Solid(textToPrint, sdl.Color{R: 200, G: 200, B: 200, A: 200})
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer fontSurf.Free()

	// fontW, fontH := fontSurf.W, fontSurf.H

	fontTexture, err := rend.CreateTextureFromSurface(fontSurf)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer fontTexture.Destroy()

	rect := butt.GetRect()

	/* TODO use the size of the fontTexture */
	// W := min(fontW, rect.W)
	// H := min(fontH, rect.H)

	//rect_ := &sdl.Rect{X: rect.X, Y: rect.Y, W: W, H: H}
	rect_ := rect

	// rend.SetDrawColor(255, 255, 255, 255)
	// rend.DrawRect(rect_)

	rend.Copy(fontTexture, nil, &rect_)
}

func (butt *Outputer) SetPosition(x, y int32) {
	butt.X, butt.Y = x, y
}

func (butt *Outputer) GetRect() sdl.Rect {
	return butt.Rect
}

func (butt *Outputer) Action(ctx *BarContext) {
	// butt.action(ctx)
}

func executeCommand(commSplited []string, inRepetition bool) (out chan string, err error) {
	outChan := make(chan string)

	if inRepetition {
		go func() {
			for {
				command := exec.Command(commSplited[0], commSplited[1:]...)
				commandOuput, err := command.StdoutPipe()

				if err != nil {
					log.Println("error ejecutando el commando.", err.Error())
					<-time.After(time.Second * 5)
					continue
				}
				err = command.Start()
				if err != nil {
					log.Println("error ejecutando el commando.", err.Error())
					<-time.After(time.Second * 5)
					continue

				}
				content, err := io.ReadAll(commandOuput)
				if err != nil {
					log.Println("error ejecutando el commando.", err.Error())
					<-time.After(time.Second * 5)
					continue

				}

				outChan <- strings.Trim(string(content), " \n\t")
			}
		}()

	} else {

		command := exec.Command(commSplited[0], commSplited[1:]...)
		commandOuput, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		err = command.Start()
		if err != nil {
			log.Println("Problem executing the command:", err.Error())
			return nil, err
		} else {
			log.Println("Command Successfully executed.")

			go func() {
				buff := make([]byte, 64)
				for n, err := commandOuput.Read(buff); err == nil; n, err = commandOuput.Read(buff) {
					outChan <- strings.Trim(string(buff[:n]), " \n\t")
				}
			}()
		}
	}

	return outChan, nil
}
