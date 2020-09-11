package main

import "fmt"

func main() {
	var s []int
	s = append(s, 1)
	fmt.Println(len(s), cap(s))
	s = append(s, 2)
	fmt.Println(len(s), cap(s))
	s = append(s, 3)
	fmt.Println(len(s), cap(s))
	s = append(s, 4)
	fmt.Println(len(s), cap(s))
	s = append(s, 5)
	fmt.Println(len(s), cap(s))
}
