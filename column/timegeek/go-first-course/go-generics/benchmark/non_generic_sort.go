package main

func siftDownOrdered(data []int, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && (data[first+child] < data[first+child+1]) {
			child++
		}
		if !(data[first+root] < data[first+child]) {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func heapSortOrdered(data []int, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownOrdered(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDownOrdered(data, lo, i, first)
	}
}

func medianOfThreeOrdered(data []int, m1, m0, m2 int) {
	// sort 3 elements
	if data[m1] < data[m0] {
		data[m1], data[m0] = data[m0], data[m1]
	}
	// data[m0] <= data[m1]
	if data[m2] < data[m1] {
		data[m2], data[m1] = data[m1], data[m2]
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if data[m1] < data[m0] {
			data[m1], data[m0] = data[m0], data[m1]
		}
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func doPivotOrdered(data []int, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's "Ninther" median of three medians of three.
		s := (hi - lo) / 8
		medianOfThreeOrdered(data, lo, lo+s, lo+2*s)
		medianOfThreeOrdered(data, m, m-s, m+s)
		medianOfThreeOrdered(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeOrdered(data, lo, m, hi-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && (data[a] < data[pivot]); a++ {
	}
	b := a
	for {
		for ; b < c && !(data[pivot] < data[b]); b++ { // data[b] <= pivot
		}
		for ; b < c && (data[pivot] < data[c-1]); c-- { // data[c-1] > pivot
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		data[b], data[c-1] = data[c-1], data[b]
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let's be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if !(data[pivot] < data[hi-1]) { // data[hi-1] = pivot
			data[c], data[hi-1] = data[hi-1], data[c]
			c++
			dups++
		}
		if !(data[b-1] < data[pivot]) { // data[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if !(data[m] < data[pivot]) { // data[m] = pivot
			data[m], data[b-1] = data[b-1], data[m]
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			for ; a < b && !(data[b-1] < data[pivot]); b-- { // data[b] == pivot
			}
			for ; a < b && (data[a] < data[pivot]); a++ { // data[a] < pivot
			}
			if a >= b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			data[a], data[b-1] = data[b-1], data[a]
			a++
			b--
		}
	}
	// Swap pivot into middle
	data[pivot], data[b-1] = data[b-1], data[pivot]
	return b - 1, c
}

func insertionSortOrdered(data []int, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && (data[j] < data[j-1]); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func quickSortOrdered(data []int, a, b, maxDepth int) {
	for b-a > 12 { // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 {
			heapSortOrdered(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotOrdered(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSortOrdered(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSortOrdered(data, mhi, b)
		} else {
			quickSortOrdered(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSortOrdered(data, a, mlo)
		}
	}
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ {
			if data[i] < data[i-6] {
				data[i], data[i-6] = data[i-6], data[i]
			}
		}
		insertionSortOrdered(data, a, b)
	}
}

func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func Sort(x []int) {
	n := len(x)
	quickSortOrdered(x, 0, n, maxDepth(n))
}
