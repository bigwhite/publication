package main

import (
	"bytes"
	"fmt"
	"rle-demo/protocol"
)

func main() {
	// 构造测试数据
	// Case 1: 大量重复 -> 应该变小
	// Case 2: 包含转义符本身 -> 应该被正确处理
	// Case 3: 随机数据 -> 可能变大（膨胀）

	original := []byte{
		'A', 'A', 'A', 'A', 'A', 'A', // 6个A -> 应该被压成 3字节: FF 06 41
		'B', 'C', 'D', // 不重复 -> 原样: 42 43 44
		0xFF,          // 转义符本身 -> 应该被编码: FF 01 FF
		'X', 'X', 'X', // 3个X (未达阈值) -> 原样: 58 58 58
	}

	fmt.Printf("Original (%d bytes): %v\n", len(original), original)

	// 1. 压缩
	var compressedBuf bytes.Buffer
	encoder := protocol.NewEncoder(&compressedBuf)
	if err := encoder.Encode(original); err != nil {
		panic(err)
	}

	compressedData := compressedBuf.Bytes()
	fmt.Printf("Compressed (%d bytes): %v\n", len(compressedData), compressedData)

	// 2. 解压
	decoder := protocol.NewDecoder(&compressedBuf)
	decodedData, err := decoder.Decode()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded (%d bytes):  %v\n", len(decodedData), decodedData)

	// 3. 验证
	if bytes.Equal(original, decodedData) {
		fmt.Println("✅ Success: Data matched!")
	} else {
		fmt.Println("❌ Failure: Data mismatch!")
	}

	// 4. 计算压缩率
	ratio := float64(len(compressedData)) / float64(len(original)) * 100
	fmt.Printf("Compression Ratio: %.2f%%\n", ratio)
}
