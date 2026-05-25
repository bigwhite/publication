package main

import (
	"context"
	"log"
	"os"

	ctxpkg "github.com/yourname/go-tiny-claw/internal/context"
	"github.com/yourname/go-tiny-claw/internal/engine"
	"github.com/yourname/go-tiny-claw/internal/observability"
	"github.com/yourname/go-tiny-claw/internal/provider"
	"github.com/yourname/go-tiny-claw/internal/schema"
	"github.com/yourname/go-tiny-claw/internal/tools"
)

func main() {
	if os.Getenv("ZHIPU_API_KEY") == "" {
		log.Fatal("请先导出 ZHIPU_API_KEY 环境变量")
	}

	workDir, _ := os.Getwd()
	modelName := "glm-4.5-air"

	// 1. 初始化真实的底层大脑
	realProvider := provider.NewZhipuOpenAIProvider(modelName)

	sessionID := "test_observability_001"
	sess := ctxpkg.GlobalSessionMgr.GetOrCreate(sessionID, workDir)

	// 2. 核心拼装：用 Tracker 将真实的大脑包裹起来
	trackedProvider := observability.NewCostTracker(realProvider, modelName, sess)

	registry := tools.NewRegistry()
	registry.Register(tools.NewBashTool(workDir))

	// 3. 将被包裹的 Provider 注入给 Engine (Engine 毫不知情)
	eng := engine.NewAgentEngine(trackedProvider, registry, false, false)
	reporter := engine.NewTerminalReporter()

	prompt := `请用 bash 帮我用 date 命令查一下现在的时间。`

	log.Println("\n>>> 🚀 启动带仪表盘的可观测性测试...")
	sess.Append(schema.Message{Role: schema.RoleUser, Content: prompt})

	err := eng.Run(context.Background(), sess, reporter)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}

	log.Printf("\n================ 财务报表 ================\n")
	log.Printf("会话 ID: %s\n", sess.ID)
	log.Printf("总消耗 Input Tokens: %d\n", sess.TotalPromptTokens)
	log.Printf("总消耗 Output Tokens: %d\n", sess.TotalCompletionTokens)
	log.Printf("总计费用 (CNY): ¥%.6f\n", sess.TotalCostCNY)
	log.Printf("==========================================\n")
}
