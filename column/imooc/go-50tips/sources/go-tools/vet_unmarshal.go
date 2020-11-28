package main

import (
	"encoding/json"
)

type Foo struct {
	Name string
	Age  int
}

func main() {
	var v Foo
	json.Unmarshal([]byte{}, v)
}
