package main

import "fmt"

// Node 双向链表节点
type Node struct {
	Key, Val   int
	Prev, Next *Node
}

// LRUCache LRU 缓存结构体
type LRUCache struct {
	capacity int
	cache    map[int]*Node // Key -> Node 指针

	// 使用 dummy head 和 dummy tail 简化边界操作
	// 链表方向：head(最新) <-> ... <-> tail(最旧)
	head, tail *Node
}

func Constructor(capacity int) LRUCache {
	l := LRUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		head:     &Node{}, // Dummy Head
		tail:     &Node{}, // Dummy Tail
	}
	// 初始化双向链表： head <-> tail
	l.head.Next = l.tail
	l.tail.Prev = l.head
	return l
}

// removeNode 从链表中摘除一个节点
func (this *LRUCache) removeNode(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

// addToHead 将节点插入到头部（Dummy Head 之后）
func (this *LRUCache) addToHead(node *Node) {
	node.Prev = this.head
	node.Next = this.head.Next

	this.head.Next.Prev = node
	this.head.Next = node
}

// moveToHead 将一个已存在的节点移动到头部
func (this *LRUCache) moveToHead(node *Node) {
	this.removeNode(node)
	this.addToHead(node)
}

// removeTail 删除尾部节点（Dummy Tail 之前），并返回它（以便从 map 删除）
func (this *LRUCache) removeTail() *Node {
	node := this.tail.Prev
	this.removeNode(node)
	return node
}

func (this *LRUCache) Get(key int) int {
	if node, ok := this.cache[key]; ok {
		// 命中缓存，将其标记为最近使用（移到头部）
		this.moveToHead(node)
		return node.Val
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if node, ok := this.cache[key]; ok {
		// key 存在，更新值，并移到头部
		node.Val = value
		this.moveToHead(node)
	} else {
		// key 不存在，创建新节点
		newNode := &Node{Key: key, Val: value}
		this.cache[key] = newNode
		this.addToHead(newNode) // 放入头部

		// 检查容量
		if len(this.cache) > this.capacity {
			// 淘汰最久未使用的（尾部）
			removed := this.removeTail()
			delete(this.cache, removed.Key) // 从 map 中删除
		}
	}
}

func main() {
	lru := Constructor(2) // 容量为 2

	lru.Put(1, 1)                      // cache: {1=1}
	lru.Put(2, 2)                      // cache: {1=1, 2=2}
	fmt.Println("Get(1):", lru.Get(1)) // 返回 1, cache: {2=2, 1=1} (1 变为最新)

	lru.Put(3, 3)                      // 容量满，淘汰 2。cache: {1=1, 3=3}
	fmt.Println("Get(2):", lru.Get(2)) // 返回 -1 (未找到)

	lru.Put(4, 4)                      // 容量满，淘汰 1。cache: {3=3, 4=4}
	fmt.Println("Get(1):", lru.Get(1)) // 返回 -1 (未找到)
	fmt.Println("Get(3):", lru.Get(3)) // 返回 3
	fmt.Println("Get(4):", lru.Get(4)) // 返回 4
}
