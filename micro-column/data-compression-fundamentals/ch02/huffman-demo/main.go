package main

import (
	"container/heap"
	"fmt"
	"huffman-demo/tree"
)

// 1. 统计频率
func countFreq(text string) map[byte]int {
	freq := make(map[byte]int)
	for i := 0; i < len(text); i++ {
		freq[text[i]]++
	}
	return freq
}

// 2. 构建 Huffman 树
func buildHuffmanTree(freqMap map[byte]int) *tree.Node {
	// 将所有字符转换为叶子节点
	nodes := make([]*tree.Node, 0, len(freqMap))
	for char, freq := range freqMap {
		nodes = append(nodes, &tree.Node{Char: char, Freq: freq})
	}

	// 初始化优先队列
	pq := tree.NewQueue(nodes)

	// 循环合并，直到只剩一个根节点
	for pq.Len() > 1 {
		// 取出频率最小的两个节点
		left := heap.Pop(pq).(*tree.Node)
		right := heap.Pop(pq).(*tree.Node)

		// 合并为父节点 (频率相加)
		parent := &tree.Node{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}

		// 放回队列
		heap.Push(pq, parent)
	}

	return heap.Pop(pq).(*tree.Node)
}

// 3. 生成编码表 (递归遍历树)
func generateCodes(node *tree.Node, currentCode string, codes map[byte]string) {
	if node == nil {
		return
	}

	// 如果是叶子节点，保存编码
	if node.IsLeaf() {
		codes[node.Char] = currentCode
		return
	}

	// 往左添 '0'，往右添 '1'
	generateCodes(node.Left, currentCode+"0", codes)
	generateCodes(node.Right, currentCode+"1", codes)
}

func main() {
	// 测试文本
	text := "this is an example for huffman encoding"
	fmt.Printf("Original Text: %q\n", text)

	// Step 1: 统计
	freqs := countFreq(text)
	fmt.Println("\n[Frequency Table]")
	for c, f := range freqs {
		fmt.Printf("'%c': %d\n", c, f)
	}

	// Step 2: 建树
	root := buildHuffmanTree(freqs)

	// Step 3: 生成编码
	codes := make(map[byte]string)
	generateCodes(root, "", codes)

	// Step 4: 打印结果与压缩率估算
	fmt.Println("\n[Huffman Codes]")
	var totalBits int
	for char, code := range codes {
		fmt.Printf("'%c': %s\n", char, code)
		// 计算压缩后的总位数: 频率 * 编码长度
		totalBits += freqs[char] * len(code)
	}

	originalBits := len(text) * 8
	fmt.Printf("\nOriginal Size: %d bits\n", originalBits)
	fmt.Printf("Compressed Size: %d bits\n", totalBits)

	// 压缩率 = (1 - 压缩后/压缩前) * 100%
	fmt.Printf("Compression Rate: %.2f%% of original\n",
		float64(totalBits)/float64(originalBits)*100)
}
