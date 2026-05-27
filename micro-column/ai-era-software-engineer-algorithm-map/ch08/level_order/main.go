package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func LevelOrder(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}

	var res [][]int
	// 队列：初始化加入根节点
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levelSize := len(queue) // 关键：锁定当前层的节点数
		var currentLevel []int

		for i := 0; i < levelSize; i++ {
			// 1. 出队
			node := queue[0]
			queue = queue[1:]

			// 2. 记录值
			currentLevel = append(currentLevel, node.Val)

			// 3. 子节点入队（下一层）
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		res = append(res, currentLevel)
	}

	return res
}

func main() {
	// 构建树: [3,9,20,null,null,15,7]
	root := &TreeNode{3,
		&TreeNode{9, nil, nil},
		&TreeNode{20, &TreeNode{15, nil, nil}, &TreeNode{7, nil, nil}},
	}

	fmt.Printf("Level Order: %v\n", LevelOrder(root))
}
