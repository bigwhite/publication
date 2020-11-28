package main

import (
	"fmt"
	"time"
)

func subTwoTimeHasMonotonic() {
	t1 := time.Now()
	time.Sleep(time.Second * 5)
	t2 := time.Now()
	diff := t2.Sub(t1)
	fmt.Printf("[hasMonotonic = 1] t2 - t1 = %v\n", diff)

}

func subTwoTimeNoMonotonic() {
	t1 := time.Date(2020, 6, 18, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2020, 6, 18, 12, 0, 0, 0, time.UTC)
	diff := t2.Sub(t1)
	fmt.Printf("[hasMonotonic = 0] t2 - t1 = %v\n", diff)
}

func main() {
	subTwoTimeHasMonotonic()
	subTwoTimeNoMonotonic()
}
