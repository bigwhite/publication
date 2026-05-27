package main

import (
	"fmt"
	"strconv"
	"strings"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Codec 封装序列化器
type Codec struct{}

func Constructor() Codec {
	return Codec{}
}

// Serializes a tree to a single string.
// 格式: "1,2,#,#,3,4,#,#,5,#,#"
func (this *Codec) serialize(root *TreeNode) string {
	var sb strings.Builder
	var dfs func(*TreeNode)

	dfs = func(node *TreeNode) {
		if node == nil {
			sb.WriteString("#,")
			return
		}
		sb.WriteString(strconv.Itoa(node.Val) + ",")
		dfs(node.Left)
		dfs(node.Right)
	}

	dfs(root)
	return sb.String()
}

// Deserializes your encoded data to tree.
func (this *Codec) deserialize(data string) *TreeNode {
	vals := strings.Split(data, ",")
	// 也就是一个 Queue，反序列化时不断消费队头
	var build func() *TreeNode

	build = func() *TreeNode {
		if len(vals) == 0 {
			return nil
		}
		// Pop 队头
		valStr := vals[0]
		vals = vals[1:]

		if valStr == "#" || valStr == "" {
			return nil
		}

		val, _ := strconv.Atoi(valStr)
		node := &TreeNode{Val: val}
		node.Left = build()  // 递归构建左子树
		node.Right = build() // 递归构建右子树
		return node
	}

	return build()
}

func main() {
	//     1
	//    / \
	//   2   3
	//      / \
	//     4   5
	root := &TreeNode{1,
		&TreeNode{2, nil, nil},
		&TreeNode{3, &TreeNode{4, nil, nil}, &TreeNode{5, nil, nil}},
	}

	ser := Constructor()
	deser := Constructor()

	data := ser.serialize(root)
	fmt.Printf("Serialized: %s\n", data)

	newRoot := deser.deserialize(data)
	fmt.Printf("Deserialized Root Val: %d\n", newRoot.Val)
}
