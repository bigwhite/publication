package main

import "fmt"

// Permission 定义权限位掩码
type Permission uint8

const (
	PermRead   Permission = 1 << iota // 0001
	PermWrite                         // 0010
	PermDelete                        // 0100
	PermAdmin                         // 1000
)

// String 方便打印
func (p Permission) String() string {
	s := ""
	if p&PermRead != 0 {
		s += "READ "
	}
	if p&PermWrite != 0 {
		s += "WRITE "
	}
	if p&PermDelete != 0 {
		s += "DELETE "
	}
	if p&PermAdmin != 0 {
		s += "ADMIN "
	}
	if s == "" {
		s = "NONE"
	}
	return fmt.Sprintf("[%04b] %s", p, s)
}

// Has 检查是否拥有某权限
func (p Permission) Has(target Permission) bool {
	// 核心：AND 运算结果不为 0
	return p&target == target
}

// Add 添加权限
func (p *Permission) Add(target Permission) {
	// 核心：OR 运算
	*p |= target
}

// Remove 移除权限
func (p *Permission) Remove(target Permission) {
	// 核心：AND NOT (Go 特有 &^)
	// 相当于 p & (^target)
	*p &^= target
}

// Toggle 切换权限（有则删，无则增）
func (p *Permission) Toggle(target Permission) {
	// 核心：XOR 运算
	*p ^= target
}

func main() {
	// 初始：只有读权限
	var userPerm Permission = PermRead
	fmt.Println("Initial:", userPerm)

	// 1. 添加写权限
	userPerm.Add(PermWrite)
	fmt.Println("Add Write:", userPerm)

	// 2. 检查权限
	fmt.Printf("Can Write? %v\n", userPerm.Has(PermWrite))
	fmt.Printf("Can Delete? %v\n", userPerm.Has(PermDelete))

	// 3. 添加管理员（顺便带上了删除，如果逻辑如此设计）
	// 这里演示单独添加删除
	userPerm.Add(PermDelete)
	fmt.Println("Add Delete:", userPerm)

	// 4. 移除写权限
	userPerm.Remove(PermWrite)
	fmt.Println("Remove Write:", userPerm)

	// 5. 切换读取权限（原本有，变无）
	userPerm.Toggle(PermRead)
	fmt.Println("Toggle Read:", userPerm)
}
