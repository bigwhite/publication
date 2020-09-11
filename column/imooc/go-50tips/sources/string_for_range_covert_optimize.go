package main

import (
	"fmt"
	"testing"
)

func convert() {
	s := "中国欢迎您，北京欢迎您"
	sl := []byte(s)
	for _, v := range sl {
		_ = v
	}
}
func convertWithOptimize() {
	s := "中国欢迎您，北京欢迎您"
	for _, v := range []byte(s) {
		_ = v
	}
}

func main() {
	fmt.Println(testing.AllocsPerRun(1, convert))
	fmt.Println(testing.AllocsPerRun(1, convertWithOptimize))
}
