package main

import "encoding/json"

func main() {

	var nilSlice []int
	var emptySlice = make([]int, 0, 5)

	println(nilSlice == nil)   // true
	println(emptySlice == nil) // false

	println(nilSlice, len(nilSlice), cap(nilSlice))
	println(emptySlice, len(emptySlice), cap(emptySlice))

	/*
		emptySlice = append(emptySlice, 1)
		println(emptySlice)
	*/

	m := map[string][]int{
		"nilSlice":   nilSlice,
		"emptySlice": emptySlice,
	}

	b, err := json.Marshal(m)
	if err != nil {
		println("json marshal error:", err)
		return
	}

	println(string(b))
}
