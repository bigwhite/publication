package main

import "fmt"

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

type Trie struct {
	root *TrieNode
}

func Constructor() Trie {
	return Trie{root: NewTrieNode()}
}

// Insert 插入单词
func (this *Trie) Insert(word string) {
	node := this.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// Search 查找单词是否存在
func (this *Trie) Search(word string) bool {
	node := this.searchNode(word)
	return node != nil && node.isEnd
}

// StartsWith 查找是否有前缀
func (this *Trie) StartsWith(prefix string) bool {
	return this.searchNode(prefix) != nil
}

// 辅助函数
func (this *Trie) searchNode(s string) *TrieNode {
	node := this.root
	for _, char := range s {
		if next, ok := node.children[char]; ok {
			node = next
		} else {
			return nil
		}
	}
	return node
}

func main() {
	trie := Constructor()
	trie.Insert("apple")
	fmt.Println("Search 'apple':", trie.Search("apple"))     // true
	fmt.Println("Search 'app':", trie.Search("app"))         // false
	fmt.Println("StartsWith 'app':", trie.StartsWith("app")) // true
	trie.Insert("app")
	fmt.Println("Search 'app':", trie.Search("app")) // true
}
