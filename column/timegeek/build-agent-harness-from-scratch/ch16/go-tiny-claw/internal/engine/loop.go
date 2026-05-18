package engine

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	ctxpkg "github.com/yourname/go-tiny-claw/internal/context"
	"github.com/yourname/go-tiny-claw/internal/provider"
	"github.com/yourname/go-tiny-claw/internal/schema"
	"github.com/yourname/go-tiny-claw/internal/tools"
)

type AgentEngine struct {
	provider       provider.LLMProvider
	registry       tools.Registry
	EnableThinking bool
	PlanMode       bool
	compactor      *ctxpkg.Compactor
	recovery       *ctxpkg.RecoveryManager
	injector       *ReminderInjector // 【新增】提醒注入器
}

func NewAgentEngine(p provider.LLMProvider, r tools.Registry, enableThinking bool, planMode bool) *AgentEngine {
	return &AgentEngine{
		provider:       p,
		registry:       r,
		EnableThinking: enableThinking,
		PlanMode:       planMode,
		compactor:      ctxpkg.NewCompactor(20000, 6),
		recovery:       ctxpkg.NewRecoveryManager(),
		injector:       NewReminderInjector(), // 【初始化注入器】
	}
}

func (e *AgentEngine) Run(ctx context.Context, session *ctxpkg.Session, reporter Reporter) error {
	log.Printf("[Engine] 唤醒会话 [%s]，锁定工作区: %s (PlanMode: %v)\n", session.ID, session.WorkDir, e.PlanMode)

	composer := ctxpkg.NewPromptComposer(session.WorkDir, e.PlanMode)
	systemMsg := composer.Build()

	for {
		availableTools := e.registry.GetAvailableTools()
		workingMemory := session.GetWorkingMemory(20)

		var contextHistory []schema.Message
		contextHistory = append(contextHistory, systemMsg)
		contextHistory = append(contextHistory, workingMemory...)
		compactedContext := e.compactor.Compact(contextHistory)

		var currentTurnThinkingContent string

		// Phase 1: Thinking
		if e.EnableThinking {
			if reporter != nil {
				reporter.OnThinking(ctx)
			}
			thinkResp, err := e.provider.Generate(ctx, compactedContext, nil)
			if err != nil {
				return fmt.Errorf("Thinking 阶段失败: %w", err)
			}
			if thinkResp.Content != "" {
				currentTurnThinkingContent = thinkResp.Content
				compactedContext = append(compactedContext, *thinkResp)
			}
		}

		// Phase 2: Action
		actionResp, err := e.provider.Generate(ctx, compactedContext, availableTools)
		if err != nil {
			return fmt.Errorf("Action 阶段失败: %w", err)
		}

		finalAssistantMsg := schema.Message{
			Role:      schema.RoleAssistant,
			Content:   strings.TrimSpace(currentTurnThinkingContent + "\n" + actionResp.Content),
			ToolCalls: actionResp.ToolCalls,
		}
		session.Append(finalAssistantMsg)

		if actionResp.Content != "" && reporter != nil {
			reporter.OnMessage(ctx, actionResp.Content)
		}

		if len(actionResp.ToolCalls) == 0 {
			break
		}

		observationMsgs := make([]schema.Message, len(actionResp.ToolCalls))
		var wg sync.WaitGroup

		// 用于收集本轮执行的最后一个工具供 Reminder 分析
		var lastToolCall schema.ToolCall
		var lastToolResult schema.ToolResult

		for i, toolCall := range actionResp.ToolCalls {
			wg.Add(1)

			go func(idx int, call schema.ToolCall) {
				defer wg.Done()

				if reporter != nil {
					reporter.OnToolCall(ctx, call.Name, string(call.Arguments))
				}

				result := e.registry.Execute(ctx, call)

				finalOutput := result.Output
				if result.IsError {
					finalOutput = e.recovery.AnalyzeAndInject(call.Name, result.Output)
				}

				if reporter != nil {
					displayOutput := finalOutput
					if len(displayOutput) > 200 {
						displayOutput = displayOutput[:200] + "... (已截断)"
					}
					reporter.OnToolResult(ctx, call.Name, displayOutput, result.IsError)
				}

				observationMsgs[idx] = schema.Message{
					Role:       schema.RoleUser,
					Content:    finalOutput,
					ToolCallID: call.ID,
				}

				if idx == 0 {
					lastToolCall = call
					lastToolResult = result
				}
			}(i, toolCall)
		}

		wg.Wait()

		session.Append(observationMsgs...)

		// 【核心防线】：在进入下一轮前，进行死循环探测与注入
		reminderMsg := e.injector.CheckAndInject(lastToolCall, lastToolResult)
		if reminderMsg != nil {
			session.Append(*reminderMsg)
		}
	}

	return nil
}
