package main

import "fmt"

type FinishedHouse struct {
	style                  int    // 0: Chinese, 1: American, 2: European
	centralAirConditioning bool   // true or false
	floorMaterial          string // "ground-tile" or ”wood"
	wallMaterial           string // "latex" or "paper" or "diatom-mud"
}

func NewFinishedHouse(style int, centralAirConditioning bool,
	floorMaterial, wallMaterial string) *FinishedHouse {

	// here: you should do some check to the arguments passed

	h := &FinishedHouse{
		style:                  style,
		centralAirConditioning: centralAirConditioning,
		floorMaterial:          floorMaterial,
		wallMaterial:           wallMaterial,
	}

	return h
}

func main() {
	fmt.Printf("%+v\n", NewFinishedHouse(0, true, "wood", "paper"))
}
