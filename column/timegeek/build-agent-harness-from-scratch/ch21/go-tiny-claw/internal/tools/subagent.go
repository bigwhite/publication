package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/yourname/go-tiny-claw/internal/schema"
)

// AgentRunner 定义了引擎向外部工具暴露的特定执行能力接口
type AgentRunner interface {
	RunSub(ctx context.Context, taskPrompt string, readOnlyRegistry Registry, reporter interface{}) (string, error)
}

type SubagentTool struct {
	runner           AgentRunner
	readOnlyRegistry Registry
	reporter         interface{} // 暂时用 interface 规避包依赖，底层通过反射或断言使用
}

func NewSubagentTool(runner AgentRunner, readOnlyRegistry Registry, reporter interface{}) *SubagentTool {
	return &SubagentTool{
		runner:           runner,
		readOnlyRegistry: readOnlyRegistry,
		reporter:         reporter,
	}
}

func (t *SubagentTool) Name() string {
	return "spawn_subagent"
}

func (t *SubagentTool) Definition() schema.ToolDefinition {
	return schema.ToolDefinition{
		Name:        t.Name(),
		Description: "派出一个专门用于深度探索（Exploration）的子智能体。当你需要阅读大量代码、跨文件查找逻辑时请调用此工具。它在探索完毕后，会给你返回一份极度精炼的摘要报告。",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"task_prompt": map[string]interface{}{
					"type":        "string",
					"description": "给子智能体下达的明确探索指令。",
				},
			},
			"required": []string{"task_prompt"},
		},
	}
}

type subagentArgs struct {
	TaskPrompt string `json:"task_prompt"`
}

func (t *SubagentTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var input subagentArgs
	if err := json.Unmarshal(args, &input); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	log.Printf("[Subagent] 🚀 主 Agent 发起委派！正在拉起探路者: [%s]...\n", input.TaskPrompt)

	// 【修改】：在接口调用中，将工具持有的 reporter 透传下去
	// (在 loop.go 的 RunSub 实现中，可以通过断言恢复 Reporter 接口)
	summary, err := t.runner.RunSub(ctx, input.TaskPrompt, t.readOnlyRegistry, t.reporter)

	if err != nil {
		return fmt.Errorf("子智能体执行失败: %v", err).Error(), nil
	}

	log.Printf("[Subagent] ✅ 子智能体任务结束。报告返回给主干...")

	return fmt.Sprintf("【子智能体探索报告】:\n%s", summary), nil
}
