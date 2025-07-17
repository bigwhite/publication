package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
)

func main() {
	// 确定插件路径，这里假设.so文件与可执行文件在同一目录或特定子目录
	// 在实际部署中，路径管理需要更健壮
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	pluginDir := filepath.Dir(exePath)
	// 如果从项目根目录运行 go run ./plugins/native/main.go,
	// .so文件可能在 ./plugins/native/myplugin/myplugin.so
	// 为简单起见，我们假设 .so 文件已拷贝到与 main.go 同级或易于访问的位置
	// 或者直接指定绝对/相对路径
	// 对于本示例，假设myplugin.so在./myplugin/myplugin.so (相对于main.go的父目录)
	// 为了能在 `go run` 时找到，我们调整一下路径期望
	// 或者在运行前将 myplugin.so 拷贝到 main.go 旁边

	// 尝试更可靠的路径定位（假设.so在与main.go同级的myplugin子目录编译生成）
	// 或者，你可以在运行前手动将 myplugin.so 放到与 main.go 同级
	// e.g., cp ./myplugin/myplugin.so .
	// pluginPath := "./myplugin.so"

	// 为了让 `go run ch23/plugins/native/main.go` 能工作，
	// 编译后的 myplugin.so 需要放在 ch23/plugins/native/ 目录下
	// 通常编译插件会在其自己的目录中，然后主程序加载时指定正确的路径。
	// 假设我们已经将 myplugin.so 放在了 plugins/native/ 目录下：
	// cd ch23/plugins/native/myplugin && go build -buildmode=plugin -o ../myplugin.so plugin.go
	// cd ../../.. (回到ch23根目录)
	// go run ./plugins/native/main.go

	pluginPath := pluginDir + "/myplugin.so" // 假设.so文件在当前工作目录下，或相对于main.go
	// 如果从模块根目录运行 go run ./plugins/native/main.go,
	// 则期望 myplugin.so 在 ./plugins/native/myplugin.so
	// 我们调整一下，假设运行 `go run main.go` 时，`myplugin.so` 位于 `./`
	// 这意味着你需要先 `cd ch23/plugins/native` 然后 `go run main.go`
	// 并且 `myplugin.so` 也拷贝到了 `ch23/plugins/native` 目录

	// 为了让示例可直接运行，假设编译后的.so与main.go在同一查找路径
	// 在实践中，主程序会有一个配置项来指定插件的搜索路径或具体文件
	// 这里我们手动构建一个预期的路径，假设你在ch23/plugins/native目录下运行
	// 或者在编译插件后，将其拷贝到 main.go 的旁边。

	// 最简单的运行方式：
	// 1. cd ch23/plugins/native/myplugin
	// 2. go build -buildmode=plugin -o ../myplugin.so plugin.go  (将.so输出到上一级目录)
	// 3. cd .. (进入 ch23/plugins/native 目录)
	// 4. go run main.go

	log.Printf("Attempting to load plugin from: %s\n", pluginPath)
	p, err := plugin.Open(pluginPath)
	if err != nil {
		log.Fatalf("Failed to open plugin '%s': %v\nIf running with 'go run', ensure '%s' is in the current directory or adjust path.", pluginPath, err, filepath.Base(pluginPath))
	}
	log.Println("Plugin loaded successfully.")

	// 查找导出的变量 PluginName
	pluginNameSymbol, err := p.Lookup("PluginName")
	if err != nil {
		log.Fatalf("Failed to lookup PluginName: %v", err)
	}
	pluginName, ok := pluginNameSymbol.(*string) // 需要类型断言，因为Lookup返回plugin.Symbol (interface{})
	if !ok {
		log.Fatalf("PluginName is not a *string, actual type: %T", pluginNameSymbol)
	}
	fmt.Printf("Plugin's registered name: %s\n", *pluginName)

	// 查找导出的变量 Version
	versionSymbol, err := p.Lookup("Version")
	if err != nil {
		log.Fatalf("Failed to lookup Version: %v", err)
	}
	version, ok := versionSymbol.(*string)
	if !ok {
		log.Fatalf("Version is not a *string, actual type: %T", versionSymbol)
	}
	fmt.Printf("Plugin's version: %s\n", *version)

	// 查找导出的函数 Greet
	greetSymbol, err := p.Lookup("Greet")
	if err != nil {
		log.Fatalf("Failed to lookup Greet: %v", err)
	}
	greetFunc, ok := greetSymbol.(func(string) string) // 类型断言为函数类型
	if !ok {
		log.Fatalf("Greet is not a func(string) string, actual type: %T", greetSymbol)
	}

	// 调用插件中的函数
	message := greetFunc("Go Developer")
	fmt.Println(message)

	// 尝试修改插件中的可导出变量 (如果插件设计允许)
	log.Printf("Original PluginName in plugin: %s\n", *pluginName)
	*pluginName = "MyUpdatedNativePlugin" // 修改的是主程序持有的指针指向的值

	// 再次调用函数，看看插件内部是否感知到变化（取决于插件如何使用该变量）
	// 如果Greet函数直接使用全局的PluginName，它会看到变化
	// 但如果Greet函数在调用时捕获了PluginName的副本，则可能看不到
	// 在我们的简单示例中，Greet函数每次都会读取全局的PluginName
	messageAfterChange := greetFunc("Gopher")
	fmt.Println(messageAfterChange)
}
