package main

import (
	"fmt"
	"sort"
)

// Suggestion 建议词及其热度
type Suggestion struct {
	Word  string
	Score int
}

type TrieNode struct {
	children map[rune]*TrieNode
	// 缓存该节点（作为前缀）下热度最高的 Top 3 单词
	top3 []Suggestion
}

func NewNode() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

type AutocompleteSystem struct {
	root *TrieNode
}

func NewAutocompleteSystem() *AutocompleteSystem {
	return &AutocompleteSystem{root: NewNode()}
}

// Insert 插入单词及热度
// 实际工程中，热度可能来自日志分析
func (as *AutocompleteSystem) Insert(word string, score int) {
	item := Suggestion{Word: word, Score: score}
	node := as.root

	// 1. 遍历路径，将 item 加入沿途所有节点的 top3 列表
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = NewNode()
		}
		node = node.children[char]
		as.updateTop3(node, item)
	}
}

// updateTop3 维护节点的 Top 3 列表
func (as *AutocompleteSystem) updateTop3(node *TrieNode, item Suggestion) {
	// 简单实现：先看看是否已存在（更新），不存在则追加，然后排序截断
	// 生产环境可以用最小堆优化
	found := false
	for i, s := range node.top3 {
		if s.Word == item.Word {
			node.top3[i].Score = item.Score // 更新分数
			found = true
			break
		}
	}
	if !found {
		node.top3 = append(node.top3, item)
	}

	// 排序：分数降序，字典序升序
	sort.Slice(node.top3, func(i, j int) bool {
		if node.top3[i].Score == node.top3[j].Score {
			return node.top3[i].Word < node.top3[j].Word
		}
		return node.top3[i].Score > node.top3[j].Score
	})

	// 截断 Top 3
	if len(node.top3) > 3 {
		node.top3 = node.top3[:3]
	}
}

// Search 根据前缀返回推荐列表
func (as *AutocompleteSystem) Search(prefix string) []string {
	node := as.root
	for _, char := range prefix {
		if next, ok := node.children[char]; ok {
			node = next
		} else {
			return []string{}
		}
	}

	// 直接返回缓存好的 Top 3，速度极快 O(1)
	res := make([]string, len(node.top3))
	for i, s := range node.top3 {
		res[i] = s.Word
	}
	return res
}

func main() {
	sys := NewAutocompleteSystem()

	// 模拟写入历史数据
	sys.Insert("i love you", 5)
	sys.Insert("island", 3)
	sys.Insert("ironman", 2)
	sys.Insert("i love leetcode", 2)

	fmt.Println("Query 'i':", sys.Search("i"))     // 预期包含 "i love you", "island", "ironman"
	fmt.Println("Query 'i l':", sys.Search("i l")) // 预期 "i love you", "i love leetcode"
	fmt.Println("Query 'ir':", sys.Search("ir"))   // 预期 "ironman"

	// 动态更新热度
	sys.Insert("ironman", 10)                            // 钢铁侠热度飙升
	fmt.Println("Query 'i' (Updated):", sys.Search("i")) // ironman 应该排到第一
}
