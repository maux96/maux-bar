package bar_items

import (
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	//sdlImg "github.com/veandco/go-sdl2/img"
	sdl "github.com/veandco/go-sdl2/sdl"
)

type Outputer struct {
	sdl.Rect
	output    <-chan string
	lastPrint string
}

func NewOutputer(W, H int32, values map[string]string) *Outputer {
	but := Outputer{
		Rect: sdl.Rect{
			W: W,
			H: H,
		},
		lastPrint: "-",
	}

	if action, ok := values["action"]; ok {
		repet, ok0 := values["repetitive"]
		isRepetitive := false
		/* TODO change this  */
		if ok0 && repet == "true" {
			isRepetitive = true
		}
		ch, err := executeCommand(action, isRepetitive)
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

	fontSurf, err := bar.Font.RenderUTF8Solid(textToPrint, sdl.Color{R: 200, G: 200, B: 200, A: 200})
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

func executeCommand(action string, inRepetition bool) (out chan string, err error) {
	commSplited := strings.Split(action, " ")
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
