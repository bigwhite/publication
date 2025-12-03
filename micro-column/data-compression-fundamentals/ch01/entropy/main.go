package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

// CalculateEntropy 读取 reader 中的数据并计算香农熵
func CalculateEntropy(r io.Reader) (float64, int64, error) {
	// 1. 统计频率
	// 我们假设处理的是字节流，所以符号集大小为 256 (0x00-0xFF)
	counts := make([]int64, 256)
	var totalBytes int64

	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := r.Read(buf)
		if n > 0 {
			totalBytes += int64(n)
			for _, b := range buf[:n] {
				counts[b]++
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, err
		}
	}

	if totalBytes == 0 {
		return 0, 0, nil
	}

	// 2. 计算熵
	// H(X) = - Σ p(x) * log2(p(x))
	var entropy float64
	for _, count := range counts {
		if count == 0 {
			continue
		}
		// p 是该字节出现的概率
		p := float64(count) / float64(totalBytes)
		entropy -= p * math.Log2(p)
	}

	return entropy, totalBytes, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: entropy <file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	entropy, size, err := CalculateEntropy(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Size: %d bytes\n", size)
	fmt.Printf("Shannon Entropy: %.4f bits/byte\n", entropy)

	// 计算理论最小体积
	// 理论体积 = (熵 * 总字节数) / 8 (bits转bytes)
	minSize := (entropy * float64(size)) / 8

	// 计算压缩后的理论体积 / 原体积 * 100
	rate := (minSize / float64(size)) * 100

	fmt.Printf("Theoretical Min Size: %.0f bytes (%.2f%% of original)\n",
		minSize, rate)
}
