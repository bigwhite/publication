package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	// TrimSpace(string)
	fmt.Println(strings.TrimSpace("\t\n\f I love Go!! \n\r")) // I love Go!!
	fmt.Println(strings.TrimSpace("I love Go!! \f\v \n\r"))   // I love Go!!
	fmt.Println(strings.TrimSpace("I love Go!!"))             // I love Go!!

	// TrimSpace([]byte)
	fmt.Printf("%q\n", bytes.TrimSpace([]byte("\t\n\f I love Go!! \n\r"))) // "I love Go!!"
	fmt.Printf("%q\n", bytes.TrimSpace([]byte("I love Go!! \f\v \n\r")))   // "I love Go!!"
	fmt.Printf("%q\n", bytes.TrimSpace([]byte("I love Go!!")))             // "I love Go!!"

	// Trim、TrimLeft、TrimRight(string)
	fmt.Println(strings.Trim("\t\n fffI love Go!!\n \rfff", "\t\n\r f"))             // I love Go!!
	fmt.Printf("%q\n", strings.TrimLeft("\t\n fffI love Go!!\n \rfff", "\t\n\r f"))  // "I love Go!!\n \rfff"
	fmt.Printf("%q\n", strings.TrimRight("\t\n fffI love Go!!\n \rfff", "\t\n\r f")) // "\t\n fffI love Go!!"

	// Trim、TrimLeft、TrimRight([]byte)
	fmt.Printf("%q\n", bytes.Trim([]byte("\t\n fffI love Go!!\n \rfff"), "\t\n\r f"))      // I love Go!!
	fmt.Printf("%q\n", bytes.TrimLeft([]byte("\t\n fffI love Go!!\n \rfff"), "\t\n\r f"))  // "I love Go!!\n \rfff"
	fmt.Printf("%q\n", bytes.TrimRight([]byte("\t\n fffI love Go!!\n \rfff"), "\t\n\r f")) // "\t\n fffI love Go!!"

	// TrimPrefix、TrimSuffix(string)
	fmt.Printf("%q\n", strings.TrimLeft("prefix,prefix I love Go!!", "prefix,"))   // " I love Go!!"
	fmt.Printf("%q\n", strings.TrimPrefix("prefix,prefix I love Go!!", "prefix,")) // "prefix I love Go!!"
	fmt.Printf("%q\n", strings.TrimSuffix("I love Go!! suffix,suffix", ",suffix")) // "I love Go!! suffix"

	// TrimPrefix、TrimSuffix([]byte)
	fmt.Printf("%q\n", bytes.TrimPrefix([]byte("prefix,prefix I love Go!!"), []byte("prefix,"))) // "prefix I love Go!!"
	fmt.Printf("%q\n", bytes.TrimSuffix([]byte("I love Go!! suffix,suffix"), []byte(",suffix"))) // "I love Go!! suffix"

	// ToUpper、ToLower(string)
	fmt.Printf("%q\n", strings.ToUpper("i LoVe gOlaNg!!")) // "I LOVE GOLANG!!"
	fmt.Printf("%q\n", strings.ToLower("i LoVe gOlaNg!!")) // "i love golang!!"

	// ToUpper、ToLower([]byte)
	fmt.Printf("%q\n", bytes.ToUpper([]byte("i LoVe gOlaNg!!"))) // "I LOVE GOLANG!!"
	fmt.Printf("%q\n", bytes.ToLower([]byte("i LoVe gOlaNg!!"))) // "i love golang!!"

	// Map(string)
	trans := func(r rune) rune {
		switch {
		case r == 'p':
			return 'g'
		case r == 'y':
			return 'o'
		case r == 't':
			return 'l'
		case r == 'h':
			return 'a'
		case r == 'o':
			return 'n'
		case r == 'n':
			return 'g'
		}
		return r
	}
	fmt.Printf("%q\n", strings.Map(trans, "I like python!!")) // "I like golang!!"

	// Map([]byte)
	fmt.Printf("%q\n", bytes.Map(trans, []byte("I like python!!"))) // "I like golang!!"
}
