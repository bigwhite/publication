package main

import "fmt"

// PluginName 是一个导出的变量
var PluginName = "MyNativeGoPlugin"

// Greet 是一个导出的函数
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s, from %s!", name, PluginName)
}

// Version 是另一个导出变量示例
var Version = "1.0.0"

// 为了能被编译为plugin，必须有一个main函数，即使它是空的
func main() {
	// 通常插件的main函数不执行任何操作，因为其代码是被主程序加载和调用的
}
