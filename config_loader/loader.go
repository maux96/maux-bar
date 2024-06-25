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
)

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
