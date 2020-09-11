package pkg1

import (
	"fmt"

	_ "github.com/bigwhite/package-init-order/pkg2"
)

var (
	_ = constInitCheck()
	v = variableInit("v")
)

const (
	c = "c"
)

func constInitCheck() string {
	if c != "" {
		fmt.Println("pkg1: const c init")
	}
	return ""
}

func variableInit(name string) string {
	fmt.Printf("pkg1: var %s init\n", name)
	return name
}

func init() {
	fmt.Println("pkg1: init")
}
