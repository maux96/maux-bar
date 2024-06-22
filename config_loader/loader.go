package config_loader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ConfigData struct {
	Objects []struct {
		ObjType string
		W       int32
		H       int32
	}
}

func loadConfig(fileName string) (config *ConfigData, err error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var data ConfigData
	err = json.Unmarshal(content, &data)
	return &data, err
}

func PrepareBar(fileName string) (err error) {
	configData, err := loadConfig(fileName)
	if err != nil {
		return err
	}

	for i, val := range configData.Objects {
		fmt.Println(i, val)
	}
	return nil
}
