package main

func sum(args ...int) int {
	var total int

	for _, v := range args {
		total += v
	}

	return total
}

func main() {
	a, b, c := 1, 2, 3
	println(sum(a, b, c))
	nums := []int{4, 5, 6}
	println(sum(nums...))
}
