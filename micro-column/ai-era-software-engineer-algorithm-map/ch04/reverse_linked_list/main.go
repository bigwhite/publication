package main

import "fmt"

// ListNode 定义链表节点
type ListNode struct {
	Val  int
	Next *ListNode
}

// ReverseList 反转链表
func ReverseList(head *ListNode) *ListNode {
	var prev *ListNode = nil
	curr := head

	for curr != nil {
		// 1. 记录下一步要去哪里，防止断链
		nextTemp := curr.Next

		// 2. 斩断过去，回首掏（指向前驱）
		curr.Next = prev

		// 3. 整体向后移动一步
		prev = curr
		curr = nextTemp
	}

	// 最后 curr 是 nil，prev 是新的头节点
	return prev
}

// 辅助函数：打印链表
func printList(head *ListNode) {
	for head != nil {
		fmt.Printf("%d -> ", head.Val)
		head = head.Next
	}
	fmt.Println("nil")
}

func main() {
	// 构建链表 1->2->3->4->5
	head := &ListNode{1, &ListNode{2, &ListNode{3, &ListNode{4, &ListNode{5, nil}}}}}

	fmt.Print("Original: ")
	printList(head)

	newHead := ReverseList(head)

	fmt.Print("Reversed: ")
	printList(newHead)
}
