package main

import (
	"fmt"
)

// RollingHash 封装滚动哈希逻辑
type RollingHash struct {
	base    int // 进制，通常用 256 或大素数
	mod     int // 取模防止溢出
	hash    int // 当前哈希值
	highPow int // 最高位的权重 base^(len-1)
	window  []byte
}

func NewRollingHash(windowLen int) *RollingHash {
	base := 256
	mod := 1000000007
	highPow := 1
	for i := 0; i < windowLen-1; i++ {
		highPow = (highPow * base) % mod
	}
	return &RollingHash{
		base:    base,
		mod:     mod,
		highPow: highPow,
		window:  make([]byte, 0, windowLen),
	}
}

// Append 添加字符，如果是第一个窗口则初始化
func (rh *RollingHash) Append(c byte) {
	rh.hash = (rh.hash*rh.base + int(c)) % rh.mod
	rh.window = append(rh.window, c)
}

// Roll 滑动窗口：移除 oldChar，添加 newChar
func (rh *RollingHash) Roll(oldChar, newChar byte) {
	// 1. 移除高位：hash - old * highPow
	// 注意负数处理：(a - b) % mod -> (a - b + mod) % mod
	rh.hash = (rh.hash - int(oldChar)*rh.highPow) % rh.mod
	if rh.hash < 0 {
		rh.hash += rh.mod
	}

	// 2. 左移并添加低位
	rh.hash = (rh.hash*rh.base + int(newChar)) % rh.mod
}

// SensitiveFilter 敏感词过滤器
type SensitiveFilter struct {
	wordLen int
	wordMap map[int]string // Hash -> Word
	rh      *RollingHash
}

func NewFilter(words []string) *SensitiveFilter {
	if len(words) == 0 {
		return nil
	}
	// 简化：假设所有敏感词长度相同
	wl := len(words[0])
	sf := &SensitiveFilter{
		wordLen: wl,
		wordMap: make(map[int]string),
		rh:      NewRollingHash(wl),
	}

	// 预计算敏感词哈希
	for _, w := range words {
		if len(w) != wl {
			panic("Simplicity: all words must have same length")
		}
		// 计算单个词的哈希（临时用 RollingHash 算一下）
		tempRh := NewRollingHash(wl)
		for i := 0; i < wl; i++ {
			tempRh.Append(w[i])
		}
		sf.wordMap[tempRh.hash] = w
	}
	return sf
}

func (sf *SensitiveFilter) Filter(text string) bool {
	n := len(text)
	if n < sf.wordLen {
		return false
	}

	// 1. 初始化第一个窗口
	sf.rh.hash = 0 // Reset
	for i := 0; i < sf.wordLen; i++ {
		sf.rh.Append(text[i])
	}

	// 检查第一个窗口
	if _, ok := sf.wordMap[sf.rh.hash]; ok {
		// Double check: 防止哈希冲突 (text[0:wl] == word?)
		if sf.doubleCheck(text[0:sf.wordLen]) {
			return true
		}
	}

	// 2. 开始滚动
	for i := 1; i <= n-sf.wordLen; i++ {
		oldChar := text[i-1]
		newChar := text[i+sf.wordLen-1]
		sf.rh.Roll(oldChar, newChar)

		if _, ok := sf.wordMap[sf.rh.hash]; ok {
			if sf.doubleCheck(text[i : i+sf.wordLen]) {
				return true
			}
		}
	}
	return false
}

func (sf *SensitiveFilter) doubleCheck(sub string) bool {
	// 简化：直接从 map 里找是否有这个 value
	// 实际工程中，Hash 冲突极低，这一步可根据需求优化
	for _, v := range sf.wordMap {
		if v == sub {
			return true
		}
	}
	return false
}

func main() {
	words := []string{"hack", "kill", "bomb"}
	filter := NewFilter(words)

	texts := []string{
		"this is a peaceful world",
		"i want to hack the system",
		"no bombs allowed",
	}

	for _, t := range texts {
		if filter.Filter(t) {
			fmt.Printf("Blocked: %s\n", t)
		} else {
			fmt.Printf("Passed:  %s\n", t)
		}
	}
}
