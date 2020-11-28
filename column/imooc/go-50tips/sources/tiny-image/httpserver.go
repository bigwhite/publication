package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("http daemon start")
	fmt.Println("  -> listen on port:8080")
	http.ListenAndServe(":8080", nil)
}
