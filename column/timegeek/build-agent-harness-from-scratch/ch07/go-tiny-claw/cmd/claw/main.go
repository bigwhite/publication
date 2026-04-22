package main

import (
	"context"
	"log"
	"os"

	"github.com/yourname/go-tiny-claw/internal/engine"
	"github.com/yourname/go-tiny-claw/internal/provider"
	"github.com/yourname/go-tiny-claw/internal/tools"
)

func main() {
	if os.Getenv("ZHIPU_API_KEY") == "" {
		log.Fatal("请先导出 ZHIPU_API_KEY 环境变量")
	}

	workDir, _ := os.Getwd()

	llmProvider := provider.NewZhipuOpenAIProvider("glm-4.5-air")
	registry := tools.NewRegistry()

	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir)) // 挂载 Edit 工具

	// 开启慢思考模式
	eng := engine.NewAgentEngine(llmProvider, registry, workDir, false)

	prompt := `
	我当前目录下有一个 server.go 文件。
	请帮我把里面 "TODO: 增加鉴权逻辑" 下面的那个 if 语句，整个替换为：
	if user == nil {
		fmt.Println("Forbidden!")
		return
	}
	`

	err := eng.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
