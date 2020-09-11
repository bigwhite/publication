package pkg2

import "fmt"

var (
	_ = constInitCheck()
	v = variableInit("v")
)

const (
	c = "c"
)

func constInitCheck() string {
	if c != "" {
		fmt.Println("pkg2: const c init")
	}
	return ""
}

func variableInit(name string) string {
	fmt.Printf("pkg2: var %s init\n", name)
	return name
}

func init() {
	fmt.Println("pkg2: init")
}
