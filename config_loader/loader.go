package config_loader

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"maux_bar/bar_items"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/veandco/go-sdl2/sdl"
)

func SetDefaultValues(bar *bar_items.BarContext) {
	if bar.Config == nil {
		bar.Config = &bar_items.BarConfigData{}
	}
	if bar.Config.Direction == "" {
		bar.Config.Direction = bar_items.DIRECTION_HORIZONTAL
	}

	mode, err := sdl.GetDesktopDisplayMode(0)
	if err != nil {
		log.Println(err.Error())
	} else {
		if bar.Config.W == 0 && bar.Config.Direction == bar_items.DIRECTION_HORIZONTAL {
			bar.Config.W = mode.W
			if bar.Config.H == 0 {
				bar.Config.H = 48
			}
		}
		if bar.Config.H == 0 && bar.Config.Direction == bar_items.DIRECTION_VERTICAL {
			bar.Config.H = mode.H
			if bar.Config.W == 0 {
				bar.Config.W = 48
			}
		}
	}

	if bar.Config.PlaceItems == "" {
		bar.Config.PlaceItems = bar_items.PLACE_ITEMS_CENTER
	}
}

func loadConfig(fileName string) (config *bar_items.BarConfigData, err error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var data bar_items.BarConfigData
	err = json.Unmarshal(content, &data)
	return &data, err
}

func PrepareBar(bar *bar_items.BarContext, fileName string) (err error) {
	configData, err := loadConfig(fileName)
	if err != nil {
		return err
	}
	bar.Config = configData
	for i := range configData.Objects {
		obj, err := objectCreator(&configData.Objects[i])
		if err != nil {
			return err
		}
		bar.Elements = append(bar.Elements, obj)
	}
	return nil
}

func objectCreator(objectData *bar_items.RawBarObjectData) (bar_items.BarElement, error) {
	switch objectData.ObjType {
	case "button":
		actionFunc := createExecuter(strings.Split(objectData.Values["action"], " "))
		return bar_items.NewButton(objectData.W, objectData.H, objectData.Values, actionFunc), nil
	}
	return nil, errors.New("object type not found")
}

func createExecuter(command []string) func(*bar_items.BarContext) {
	return func(_ *bar_items.BarContext) {
		command := exec.Command(command[0], command[1:]...)
		command.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}
		err := command.Start()
		if err != nil {
			log.Println("Problem executing the command:", err.Error())
		} else {
			log.Println("Command Successfully executed.")
		}
	}
}
