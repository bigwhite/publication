package main

import "fmt"

func findRepeatedDnaSequences(s string) []string {
	if len(s) < 10 {
		return []string{}
	}

	// 映射表：将字符映射为 0, 1, 2, 3 (4进制)
	toInt := map[byte]int{'A': 0, 'C': 1, 'G': 2, 'T': 3}

	// 窗口长度 L = 10, 进制 K = 4
	// 最高位的权重 modulus = 4^9
	L, K := 10, 4
	modulus := 1
	for i := 0; i < L-1; i++ {
		modulus *= K
	}

	// 计算第一个窗口的哈希值
	hash := 0
	for i := 0; i < L; i++ {
		hash = hash*K + toInt[s[i]]
	}

	seen := map[int]int{hash: 1}
	var res []string

	// 开始滑动
	for i := 1; i <= len(s)-L; i++ {
		// 1. 移除最左边字符 (Remove Leading)
		// hash - s[i-1] * (4^9)
		hash = hash - toInt[s[i-1]]*modulus

		// 2. 左移一位 (Shift)
		hash = hash * K

		// 3. 加入新字符 (Add Trailing)
		hash = hash + toInt[s[i+L-1]]

		// 4. 记录与检查
		if count, ok := seen[hash]; ok {
			if count == 1 {
				res = append(res, s[i:i+L])
			}
			seen[hash]++
		} else {
			seen[hash] = 1
		}
	}

	return res
}

func main() {
	s := "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"
	fmt.Printf("Repeated: %v\n", findRepeatedDnaSequences(s))
}
