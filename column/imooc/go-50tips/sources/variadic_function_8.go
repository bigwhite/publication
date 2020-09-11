package main

import "fmt"

type FinishedHouse struct {
	style                  int    // 0: Chinese, 1: American, 2: European
	centralAirConditioning bool   // true or false
	floorMaterial          string // "ground-tile" or ”wood"
	wallMaterial           string // "latex" or "paper" or "diatom-mud"
}

type Options struct {
	Style                  int    // 0: Chinese, 1: American, 2: European
	CentralAirConditioning bool   // true or false
	FloorMaterial          string // "ground-tile" or ”wood"
	WallMaterial           string // "latex" or "paper" or "diatom-mud"
}

func NewFinishedHouse(options *Options) *FinishedHouse {
	// use default style and materials if option is nil
	var style int = 0
	var centralAirConditioning = true
	var floorMaterial = "wood"
	var wallMaterial = "paper"

	if options != nil {
		// here: you should do some check to the options passed

		style = options.Style
		centralAirConditioning = options.CentralAirConditioning
		floorMaterial = options.FloorMaterial
		wallMaterial = options.WallMaterial
	}

	h := &FinishedHouse{
		style:                  style,
		centralAirConditioning: centralAirConditioning,
		floorMaterial:          floorMaterial,
		wallMaterial:           wallMaterial,
	}

	return h
}

func main() {
	fmt.Printf("%+v\n", NewFinishedHouse(nil)) // use default options
	fmt.Printf("%+v\n", NewFinishedHouse(&Options{
		Style:                  1,
		CentralAirConditioning: false,
		FloorMaterial:          "ground-tile",
		WallMaterial:           "paper",
	}))
}
