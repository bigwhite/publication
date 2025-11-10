// pkg/parser/parser_fuzz_test.go
//go:build fuzz

package parser

import "testing"

func FuzzUnsafeAndPanicReverse(f *testing.F) {
    // 1. 添加“种子语料”
    f.Add("hello, world")
    f.Add("!12345")
    // 添加一个超过10个字符的种子，帮助 fuzzing 引擎更快地探索
    f.Add("a long string for fuzzing") 

    // 2. 启动 Fuzzing 引擎
    f.Fuzz(func(t *testing.T, original string) {
        // 这是我们的 Fuzzing 目标函数
        // 它会以 Fuzzing 引擎生成的各种 original 值为参数，被反复执行
        _, err := UnsafeAndPanicReverse(original)
        if err != nil {
            // 如果函数正常返回错误，是可接受的
            // 我们可以选择性地 t.Skip() 来告诉引擎，这类输入是“无趣的”
            t.Skip()
        }
    })
}
