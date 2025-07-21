package vetexamples

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Example 1: Printf format error
func PrintfError(name string, age int) {
	fmt.Printf("Name: %s, Age: %d years, Height: %.2f\n", name, age) // Missing argument for %.2f
}

// Example 2: Loop closure
func LoopClosureProblem() {
	var wg sync.WaitGroup
	s := []string{"a", "b", "c"}
	for _, v := range s { // v is reused in each iteration
		wg.Add(1)
		go func() { // This goroutine captures the loop variable v by reference
			defer wg.Done()
			// All goroutines will likely print 'c' because v will be 'c' when they run
			fmt.Printf("Loop var (problem): %s\n", v)
		}()
	}
	wg.Wait()
}
func LoopClosureFixed() {
	var wg sync.WaitGroup
	s := []string{"a", "b", "c"}
	for _, v := range s {
		v := v // Create a new variable v shadowing the loop variable
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("Loop var (fixed): %s\n", v)
		}()
	}
	wg.Wait()
}

// Example 3: Lost cancel
func ProcessWithContext(parentCtx context.Context) context.Context {
	newCtx, _ := context.WithTimeout(parentCtx, 5*time.Second)
	go func() { // Simulate work that respects cancellation
		<-newCtx.Done()
		fmt.Println("ProcessWithContext: context done (e.g. timeout or manual cancel)")
	}()

	// For this example, let's make a clear lost cancel case for go vet to find:
	if time.Now().Year() > 2000 { // Dummy condition
		_, cancelFuncThatWillBeLost := context.WithCancel(parentCtx)
		_ = cancelFuncThatWillBeLost // Suppress unused variable, but vet checks if called.
	}
	return newCtx // Returning the context, but what about cancel from line 60?
}

// Dummy main for package to be vet-able
func main() {
	PrintfError("Alice", 30)
	LoopClosureProblem()
	LoopClosureFixed()

	ctx := context.Background()
	derivedCtx := ProcessWithContext(ctx)
	_ = derivedCtx
}
