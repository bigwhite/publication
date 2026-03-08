package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func MergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	// 技巧：虚拟头节点，避免处理头节点为空的边界情况
	dummy := &ListNode{Val: -1}
	tail := dummy // tail 始终指向新链表的末尾

	// 只要两个链表都有数据，就比较
	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			tail.Next = list1
			list1 = list1.Next
		} else {
			tail.Next = list2
			list2 = list2.Next
		}
		// tail 向后移动
		tail = tail.Next
	}

	// 扫尾：把剩余的链表直接接在后面
	if list1 != nil {
		tail.Next = list1
	} else if list2 != nil {
		tail.Next = list2
	}

	// 返回 dummy 的下一个节点
	return dummy.Next
}

func printList(head *ListNode) {
	for head != nil {
		fmt.Printf("%d -> ", head.Val)
		head = head.Next
	}
	fmt.Println("nil")
}

func main() {
	l1 := &ListNode{1, &ListNode{2, &ListNode{4, nil}}}
	l2 := &ListNode{1, &ListNode{3, &ListNode{4, nil}}}

	fmt.Print("L1: ")
	printList(l1)
	fmt.Print("L2: ")
	printList(l2)

	merged := MergeTwoLists(l1, l2)
	fmt.Print("Merged: ")
	printList(merged)
}
