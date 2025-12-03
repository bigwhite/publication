package tree

import "fmt"

// Node Huffman 树的节点
type Node struct {
	Char  byte  // 存储的字符（仅叶子节点有效）
	Freq  int   // 频率/权重
	Left  *Node // 左子节点
	Right *Node // 右子节点
}

// IsLeaf 判断是否为叶子节点
func (n *Node) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

func (n *Node) String() string {
	if n.IsLeaf() {
		return fmt.Sprintf("{Char: %c, Freq: %d}", n.Char, n.Freq)
	}
	return fmt.Sprintf("{Freq: %d}", n.Freq)
}
