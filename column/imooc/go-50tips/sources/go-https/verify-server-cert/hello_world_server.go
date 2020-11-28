package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!\n")
	})
	fmt.Println(http.ListenAndServeTLS("localhost:8081",
		"../server-signed-by-ca.crt",
		"../server.key", nil))
}
