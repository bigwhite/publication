package main

import (
	"fmt"
	"os"
)

func main() {
	writer := os.Stdout
	data := "Hello, World!"
	n, err := fmt.Fprint(writer, data, "\n")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Fprint: written %d bytes\n", n)

	n, err = fmt.Fprintf(writer, "%s\n", data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Fprintf: written %d bytes\n", n)

	n, err = fmt.Fprintln(writer, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Fprintln: written %d bytes\n", n)
}
