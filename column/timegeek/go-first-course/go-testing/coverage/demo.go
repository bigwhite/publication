package main

// Add 函数将两个整数相加并返回结果
func Add(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	return a + b
}

// IsPositive 函数检查给定的整数是否为正数
func IsPositive(num int) bool {
	if num > 0 {
		return true
	}
	return false
}
