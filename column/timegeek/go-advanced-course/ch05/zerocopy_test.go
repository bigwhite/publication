package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

const (
	testStringSmall = "hello, world"
	testStringLarge = "Go is an open source programming language that makes it easy to build simple, reliable, and efficient sof  tware. Go语言是一种开源编程语言，它能让我们能够轻松地构建简单、可靠且高效的软件。"
)

// 测试工具函数，检查内存分配情况
func testAllocs(t *testing.T, name string, f func()) {
	t.Helper()
	n := testing.AllocsPerRun(100, f)
	t.Logf("%-40s %f allocs/op", name, n)
	if n > 0 {
		t.Logf("⚠️ %s: 有内存分配发生", name)
	} else {
		t.Logf("✅ %s: 零内存分配", name)
	}
}

// 1. 使用range遍历string转换为[]byte的场景
func TestRangeOverConvertedBytes(t *testing.T) {
	s := testStringLarge
	testAllocs(t, "Range over []byte(string)", func() {
		sum := 0
		for _, v := range []byte(s) {
			sum += int(v)
		}
		_ = sum
	})
}

// 2. 使用[]byte转换为string作为map键的场景
func TestMapKeyConversion(t *testing.T) {
	m := make(map[string]int)
	m[testStringLarge] = 42
	b := []byte(testStringLarge)

	testAllocs(t, "Map lookup with string([]byte) key", func() {
		v := m[string(b)]
		_ = v
	})
}

// 3. append([]byte, string...)操作
func TestAppendStringToBytes(t *testing.T) {
	dst := make([]byte, 0, 100)
	s := testStringSmall

	testAllocs(t, "append([]byte, string...)", func() {
		result := append(dst[:0], s...)
		_ = result
	})
}

// 4. copy([]byte, string)操作
func TestCopyStringToBytes(t *testing.T) {
	s := testStringLarge
	dst := make([]byte, len(s))

	testAllocs(t, "copy([]byte, string)", func() {
		n := copy(dst, s)
		_ = n
	})
}

// 5. 字符串比较操作
func TestCompareStringWithBytes(t *testing.T) {
	s := testStringLarge
	b := []byte(s)

	testAllocs(t, "Compare: string([]byte) == string", func() {
		equal := string(b) == s
		_ = equal
	})

	testAllocs(t, "Compare: string([]byte) != string", func() {
		notEqual := string(b) != s
		_ = notEqual
	})

	b1 := []byte(testStringLarge)
	b2 := []byte(testStringLarge)

	testAllocs(t, "Compare: string([]byte) == string([]byte)", func() {
		equal := string(b1) == string(b2)
		_ = equal
	})
}

// 6. bytes包函数的string参数
func TestBytesPackageWithString(t *testing.T) {
	s := testStringSmall

	testAllocs(t, "bytes.Contains([]byte, []byte(string))", func() {
		c := bytes.Contains([]byte("hello world"), []byte(s))
		_ = c
	})
}

// 7. for循环遍历string转换为[]byte的场景
func TestForLoopOverConvertedBytes(t *testing.T) {
	s := testStringLarge

	testAllocs(t, "for loop over []byte(string)", func() {
		bs := []byte(s)
		sum := 0
		for i := 0; i < len(bs); i++ {
			sum += int(bs[i])
		}
		_ = sum
	})
}

// 8. switch语句中使用string([]byte)
func TestSwitchWithConvertedBytes(t *testing.T) {
	b := []byte(testStringSmall)

	testAllocs(t, "switch with string([]byte)", func() {
		switch string(b) {
		case "hello":
			// 不执行
		case "world":
			// 不执行
		default:
			// 执行
		}
	})
}

func byteSliceToStringUnsafe(b []byte) string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&b[0])),
		Len:  len(b),
	}))
}

func stringToByteSliceUnsafe(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data,
		Len:  len(s),
		Cap:  len(s),
	}))
}

// 使用unsafe.String将字节切片转换为字符串
// 这是零拷贝操作，但必须确保字节切片在字符串使用期间不被修改
func byteSliceToStringUnsafeGo121(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// 使用unsafe.Slice将字符串转换为字节切片
// 这是零拷贝操作，但必须确保不修改返回的切片内容，否则会破坏字符串的不可变性
func stringToByteSliceUnsafeGo121(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func TestUnsafeConversions(t *testing.T) {
	b := []byte(testStringLarge)
	s := testStringLarge

	testAllocs(t, "Unsafe []byte to string", func() {
		s2 := byteSliceToStringUnsafe(b)
		_ = s2
	})

	testAllocs(t, "Unsafe string to []byte", func() {
		b2 := stringToByteSliceUnsafe(s)
		_ = b2
	})
}

func TestUnsafeConversionsGo121(t *testing.T) {
	b := []byte(testStringLarge)
	s := testStringLarge

	testAllocs(t, "Unsafe []byte to string(Go1.21)", func() {
		s2 := byteSliceToStringUnsafeGo121(b)
		_ = s2
	})

	testAllocs(t, "Unsafe string to []byte(Go1.21)", func() {
		b2 := stringToByteSliceUnsafeGo121(s)
		_ = b2
	})
}

// 验证优化是否存在的主测试入口

// 验证优化是否存在的主测试入口
func TestMain(m *testing.M) {
	fmt.Println("开始测试Go编译器对string和[]byte转换的零拷贝优化...")
	m.Run()
	fmt.Println("测试完成。")
}
