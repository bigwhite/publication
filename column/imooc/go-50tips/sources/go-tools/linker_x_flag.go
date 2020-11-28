// build cmd: go build -ldflags "-X main.version=v0.7.0" linker_x_flag.go
package main

import (
	"fmt"
	"os"
)

var (
	version string
)

func main() {
	if os.Args[1] == "version" {
		fmt.Println("version:", version)
		return
	}
}
