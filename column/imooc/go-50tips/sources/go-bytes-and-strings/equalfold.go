package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.EqualFold("GoLang", "golang"))               // true
	fmt.Println(bytes.Equal([]byte("GoLang"), []byte("Golang")))     // false
	fmt.Println(bytes.EqualFold([]byte("GoLang"), []byte("Golang"))) // true
}
