package main

import "fmt"

func DailyTemperatures(temperatures []int) []int {
	n := len(temperatures)
	result := make([]int, n)
	// stack 存储索引（index），而不是温度值
	// 因为我们需要计算距离（index 差值）
	stack := []int{}

	for i, currentTemp := range temperatures {
		// 当栈不为空，且当前温度 > 栈顶那天的温度
		// 说明栈顶那天找到了“下一个更高温度”
		for len(stack) > 0 && currentTemp > temperatures[stack[len(stack)-1]] {
			// 弹出栈顶索引
			prevIndex := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// 计算距离，并记录结果
			result[prevIndex] = i - prevIndex
		}
		// 当前天入栈，等待它的“更高温度”
		stack = append(stack, i)
	}

	return result
}

func main() {
	temps := []int{73, 74, 75, 71, 69, 72, 76, 73}
	fmt.Printf("Temps:  %v\n", temps)
	fmt.Printf("Result: %v\n", DailyTemperatures(temps))
}
