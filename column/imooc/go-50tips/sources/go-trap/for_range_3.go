package main

import (
	"fmt"
	"sort"
)

func main() {
	var indexes []int
	heros := map[int]string{
		1: "superman",
		2: "batman",
		3: "spiderman",
		4: "the flash",
	}
	for k, v := range heros {
		indexes = append(indexes, k)
		fmt.Println(k, v)
	}

	sort.Ints(indexes)
	for _, idx := range indexes {
		fmt.Println(heros[idx])
	}
}
