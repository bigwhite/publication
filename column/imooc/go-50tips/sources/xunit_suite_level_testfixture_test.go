package demo_test

import (
	"fmt"
	"testing"
)

func suiteSetUp(suiteName string) func() {
	fmt.Printf("\tsetUp fixture for suite %s\n", suiteName)
	return func() {
		fmt.Printf("\ttearDown fixture for suite %s\n", suiteName)
	}
}

func func1TestCase1(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func func1TestCase2(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func func1TestCase3(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func TestFunc1(t *testing.T) {
	t.Cleanup(suiteSetUp(t.Name()))
	t.Run("testcase1", func1TestCase1)
	t.Run("testcase2", func1TestCase2)
	t.Run("testcase3", func1TestCase3)
}

func func2TestCase1(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func func2TestCase2(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func func2TestCase3(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func TestFunc2(t *testing.T) {
	t.Cleanup(suiteSetUp(t.Name()))
	t.Run("testcase1", func2TestCase1)
	t.Run("testcase2", func2TestCase2)
	t.Run("testcase3", func2TestCase3)
}

func pkgSetUp(pkgName string) func() {
	fmt.Printf("package SetUp fixture for %s\n", pkgName)
	return func() {
		fmt.Printf("package TearDown fixture for %s\n", pkgName)
	}
}

func TestMain(m *testing.M) {
	defer pkgSetUp("package demo_test")()
	m.Run()
}
