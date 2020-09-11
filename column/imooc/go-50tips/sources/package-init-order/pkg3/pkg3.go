package pkg3

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
		fmt.Println("pkg3: const c init")
	}
	return ""
}

func variableInit(name string) string {
	fmt.Printf("pkg3: var %s init\n", name)
	return name
}

func init() {
	fmt.Println("pkg3: init")
}
