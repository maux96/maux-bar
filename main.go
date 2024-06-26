package main

import (
	"log"
	"flag"

	"maux_bar/bar_items"
	"maux_bar/bar_visual"
	"maux_bar/config_loader"
)

func main() {
	flag.Parse()

	var configFilePath string 
	if flag.NArg() == 0 {
		log.Fatalln("JSON config file needed as first argument.")
	} else {
		configFilePath = flag.Arg(0)
	}

	bar := bar_items.NewBarContext()
	err := config_loader.PrepareBar(bar, configFilePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	QuitFunction := bar_visual.InitEverythingSDL()
	defer QuitFunction()

	config_loader.SetDefaultValues(bar)

	bar_visual.StartBar(bar)
}
