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
	workDir += "/workspace"

	llmProvider := provider.NewZhipuOpenAIProvider("glm-4.5-air")
	reporter := engine.NewTerminalReporter()

	// 【防御沙箱】为子智能体准备受限的只读注册表
	readOnlyRegistry := tools.NewRegistry()
	readOnlyRegistry.Register(tools.NewReadFileTool(workDir))
	readOnlyRegistry.Register(tools.NewBashTool(workDir)) // 允许简单的 grep 等搜索操作

	// 为主智能体准备全功能注册表
	mainRegistry := tools.NewRegistry()
	mainRegistry.Register(tools.NewReadFileTool(workDir))
	mainRegistry.Register(tools.NewWriteFileTool(workDir))
	mainRegistry.Register(tools.NewBashTool(workDir))
	mainRegistry.Register(tools.NewEditFileTool(workDir))

	// 初始化主引擎
	eng := engine.NewAgentEngine(llmProvider, mainRegistry, false, false)

	// 【核心装配】：将带有 Engine 引用和只读 Registry 的 Subagent 工具注册进主线
	mainRegistry.Register(tools.NewSubagentTool(eng, readOnlyRegistry, reporter))

	sessionID := "test_subagent_001"
	sess := ctxpkg.GlobalSessionMgr.GetOrCreate(sessionID, workDir)

	prompt := `
	我需要你在这个遗留项目里，找到那个“核心密码”。
	为了防止污染主上下文，请你务必派出子智能体（spawn_subagent）去执行探索任务。
	你可以让子智能体使用 bash 去查找当前目录（及其所有子目录）下名为 config.txt 的文件。
	子智能体拿到密码向你汇报后，请你亲自使用 write_file 工具，将密码写在根目录的 answer.txt 里。
	`

	log.Println("\n>>> 🚀 启动多智能体协同测试...")
	sess.Append(schema.Message{Role: schema.RoleUser, Content: prompt})

	err := eng.Run(context.Background(), sess, reporter)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
