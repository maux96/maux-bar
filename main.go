package main

import (
	"maux_bar/bar_items"
	"maux_bar/bar_visual"
	"maux_bar/config_loader"
)

func main() {
	bar := bar_items.NewBarContext()
	config_loader.PrepareBar(bar, "./testConfig.json")

	bar_visual.StartBar(bar)
}
