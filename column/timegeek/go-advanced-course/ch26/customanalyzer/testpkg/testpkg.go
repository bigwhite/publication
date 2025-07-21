package testpkg // A library package (not main)

import "fmt"

// This function is correctly exported.
func ExportedFunction() {
	fmt.Println("This is an exported function.")
}

// thisFunctionIsUnexported violates our hypothetical rule if we want all top-level funcs exported.
// Or, it's just a note if we only want to list unexported top-level functions.
// For this analyzer, we assume the rule is "top-level functions should be exported".
func thisFunctionIsUnexported() { // Analyzer should flag this
	fmt.Println("This function is not exported.")
}

type MyStruct struct{}

// This is a method, should be ignored by our simple check (Recv != nil)
func (s *MyStruct) ExportedMethod() {
	fmt.Println("This is an exported method.")
}
func (s *MyStruct) unexportedMethod() {
	fmt.Println("This is an unexported method.")
}

// TestHelperFunction is a test helper, should be ignored by our check
func TestHelperFunction() {}

func init() {
	fmt.Println("testpkg init")
}
