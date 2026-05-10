package main

import (
	"context"
	"log"
	"os"

	ctxpkg "github.com/yourname/go-tiny-claw/internal/context"
	"github.com/yourname/go-tiny-claw/internal/engine"
	"github.com/yourname/go-tiny-claw/internal/provider"
	"github.com/yourname/go-tiny-claw/internal/schema"
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
	registry.Register(tools.NewBashTool(workDir))

	eng := engine.NewAgentEngine(llmProvider, registry, false)
	reporter := engine.NewTerminalReporter()

	sessionID := "test_oom_protection_001"
	sess := ctxpkg.GlobalSessionMgr.GetOrCreate(sessionID, workDir)

	// 提示：你需要在终端先执行 yes "这是一段极其冗长的、无意义的服务器报错日志信息，用来模拟 OOM 场景" | head -n 2000 > mock_log.txt
	prompt := `
	请帮我执行以下三个步骤：
	1. 使用 bash 执行 echo "开始排查日志"
	2. 读取当前目录下的巨大文件 mock_log.txt
	3. 用 bash 执行 date 命令获取当前时间，并告诉我任务完成。
	`

	sess.Append(schema.Message{Role: schema.RoleUser, Content: prompt})

	err := eng.Run(context.Background(), sess, reporter)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
