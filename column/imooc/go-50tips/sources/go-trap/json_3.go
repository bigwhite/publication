package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	m := map[string]interface{}{
		"byteSlice": []byte("hello, go"),
		"string":    "hello, go",
	}

	b, _ := json.Marshal(m)
	fmt.Println(string(b))
}
