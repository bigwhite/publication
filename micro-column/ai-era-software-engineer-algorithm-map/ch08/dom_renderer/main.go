package main

import (
	"fmt"
	"strings"
)

// ElementType 节点类型
type ElementType int

const (
	NodeElement ElementType = iota
	NodeText
)

// Node DOM 节点结构
type Node struct {
	Type     ElementType
	TagName  string            // for Element
	Props    map[string]string // 属性, e.g. class="container"
	Children []*Node           // 子节点
	Text     string            // for Text Node
}

// NewElement 创建元素节点
func NewElement(tag string, props map[string]string, children ...*Node) *Node {
	return &Node{
		Type:     NodeElement,
		TagName:  tag,
		Props:    props,
		Children: children,
	}
}

// NewText 创建文本节点
func NewText(text string) *Node {
	return &Node{
		Type: NodeText,
		Text: text,
	}
}

// RenderDOM 渲染函数：DFS 遍历 DOM 树
func RenderDOM(node *Node, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)

	if node.Type == NodeText {
		fmt.Printf("%s%s\n", indent, node.Text)
		return
	}

	// 1. 打印开始标签
	propsStr := ""
	for k, v := range node.Props {
		propsStr += fmt.Sprintf(" %s=\"%s\"", k, v)
	}
	fmt.Printf("%s<%s%s>\n", indent, node.TagName, propsStr)

	// 2. 递归处理子节点 (DFS)
	for _, child := range node.Children {
		RenderDOM(child, indentLevel+1)
	}

	// 3. 打印结束标签
	fmt.Printf("%s</%s>\n", indent, node.TagName)
}

func main() {
	// 构建 DOM 树
	// <div id="app">
	//   <h1>Hello World</h1>
	//   <ul class="list">
	//     <li>Item 1</li>
	//     <li>Item 2</li>
	//   </ul>
	// </div>

	root := NewElement("div", map[string]string{"id": "app"},
		NewElement("h1", nil, NewText("Hello World")),
		NewElement("ul", map[string]string{"class": "list"},
			NewElement("li", nil, NewText("Item 1")),
			NewElement("li", nil, NewText("Item 2")),
		),
	)

	fmt.Println("--- Rendered HTML ---")
	RenderDOM(root, 0)
}
