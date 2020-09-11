// +build gofuzz

package parser

func Fuzz(data []byte) int {
	ParseComplex(data)
	return 0
}

