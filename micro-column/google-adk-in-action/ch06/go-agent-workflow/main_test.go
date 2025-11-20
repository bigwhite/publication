package main

import (
	"context"
	"os"
	"strings"
	"testing"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

// createTestPipelineAgent 是一个辅助函数，用于在测试中创建我们的 Agent
// 注意：为了让测试代码能访问到它，你需要将这个函数也放在 main 包中（即 main_test.go）
func createTestPipelineAgent(t *testing.T) agent.Agent {
	t.Helper() // 标记这是一个测试辅助函数

	ctx := context.Background()
	model, err := gemini.NewModel(ctx, "gemini-2.0-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		t.Fatalf("Failed to create model for test: %v", err)
	}

	// Code Writer Agent: 接收用户请求，生成代码，并将结果存入 state["generated_code"]
	codeWriterAgent, _ := llmagent.New(llmagent.Config{
		Name:  "CodeWriterAgent",
		Model: model,
		Instruction: `You are a Python Code Generator.
Based *only* on the user's request, write Python code that fulfills the requirement.
Output *only* the complete Python code block, enclosed in triple backticks ('''python ... ''').
Do not add any other text before or after the code block.`,
		Description: "Writes initial Python code based on a specification.",
		OutputKey:   "generated_code", // 将输出存入 state["generated_code"]
	})

	// Code Reviewer Agent: 从 state 中读取代码，提供反馈，并将结果存入 state["review_comments"]
	codeReviewerAgent, _ := llmagent.New(llmagent.Config{
		Name:  "CodeReviewerAgent",
		Model: model,
		Instruction: `You are an expert Python Code Reviewer.
Your task is to provide constructive feedback on the provided code.

**Code to Review:**
'''python
{generated_code}
'''

Provide your feedback as a concise, bulleted list.
If the code is excellent and requires no changes, simply state: "No major issues found."`,
		Description: "Reviews code and provides feedback.",
		OutputKey:   "review_comments", // 将输出存入 state["review_comments"]
	})

	// Code Refactorer Agent: 从 state 中读取代码和反馈，进行重构
	codeRefactorerAgent, _ := llmagent.New(llmagent.Config{
		Name:  "CodeRefactorerAgent",
		Model: model,
		Instruction: `You are a Python Code Refactoring AI.
Your goal is to improve the given Python code based on the provided review comments.

**Original Code:**
'''python
{generated_code}
'''

**Review Comments:**
{review_comments}

**Task:**
Apply the suggestions from the review comments to refactor the original code.
Output *only* the final, refactored Python code block.`,
		Description: "Refactors code based on review comments.",
		OutputKey:   "refactored_code", // 将输出存入 state["refactored_code"]
	})

	// Sequential Agent for the pipeline
	codePipelineAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name: "CodePipelineAgent",
			SubAgents: []agent.Agent{
				codeWriterAgent,
				codeReviewerAgent,
				codeRefactorerAgent,
			},
			Description: "Executes a sequence of code writing, reviewing, and refactoring.",
		},
	})
	if err != nil {
		t.Fatalf("failed to create codePipelineAgent for test: %s", err)
	}
	return codePipelineAgent
}

func TestCodePipelineAgent_Evaluation(t *testing.T) {
	// 如果没有设置 API Key，则跳过此集成测试
	if os.Getenv("GOOGLE_API_KEY") == "" {
		t.Skip("Skipping evaluation test: GOOGLE_API_KEY is not set")
	}

	// 1. Inputs: 定义 Agent 和测试用例
	ctx := context.Background()
	agentToTest := createTestPipelineAgent(t)

	testQuery := "Write a Python function that takes a list of integers and returns a new list containing only the even numbers."

	// 2. Run Agent: 执行 Agent
	sessionService := session.InMemoryService()
	r, err := runner.New(runner.Config{
		AppName:        "test-app",
		Agent:          agentToTest,
		SessionService: sessionService,
	})
	if err != nil {
		t.Fatalf("Failed to create runner: %v", err)
	}

	sessionID := "eval-session-1"
	_, err = sessionService.Create(ctx, &session.CreateRequest{AppName: "test-app", UserID: "eval-user", SessionID: sessionID})
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	stream := r.Run(ctx, "eval-user", sessionID, genai.NewContentFromText(testQuery, genai.RoleUser), agent.RunConfig{})

	// 3. Actual Results: 收集实际结果
	var actualEvents []*session.Event
	for event, err := range stream {
		if err != nil {
			t.Fatalf("Agent run failed: %v", err)
		}
		actualEvents = append(actualEvents, event)
	}

	// 4. Compare: 对比实际结果与预期

	// 4.1 评估轨迹 (Trajectory)
	t.Run("Evaluate Trajectory", func(t *testing.T) {
		expectedTrajectory := []string{"CodeWriterAgent", "CodeReviewerAgent", "CodeRefactorerAgent"}
		var actualTrajectory []string
		for _, event := range actualEvents {
			if event.Author != "user" && event.Author != "" {
				// 去重，只记录每个 Agent 第一次出现
				isNew := true
				for _, name := range actualTrajectory {
					if name == event.Author {
						isNew = false
						break
					}
				}
				if isNew {
					actualTrajectory = append(actualTrajectory, event.Author)
				}
			}
		}

		expected := strings.Join(expectedTrajectory, " -> ")
		actual := strings.Join(actualTrajectory, " -> ")

		if actual != expected {
			t.Errorf("Trajectory mismatch!\n  Expected: %s\n  Actual:   %s", expected, actual)
		} else {
			t.Logf("✅ Trajectory validation passed: %s", actual)
		}
	})

	// 4.2 评估最终结果 (Final Response)
	t.Run("Evaluate Final Response", func(t *testing.T) {
		if len(actualEvents) == 0 {
			t.Fatal("No events were generated.")
		}
		// 最后一个非 partial 事件通常包含最终结果
		var finalCode string
		for i := len(actualEvents) - 1; i >= 0; i-- {
			event := actualEvents[i]
			if !event.Partial && event.LLMResponse.Content != nil && len(event.LLMResponse.Content.Parts) > 0 {
				finalCode = event.LLMResponse.Content.Parts[0].Text
				break
			}
		}

		if finalCode == "" {
			t.Fatal("Could not determine final response from events.")
		}
		
		// 我们期望最终的代码使用了列表推导式
		expectedCodeSnippet := "[number for number in numbers if number % 2 == 0]"
		if !strings.Contains(finalCode, expectedCodeSnippet) {
			t.Errorf("Final response validation failed!\n  Expected code to contain: %s\n  But got:\n%s", expectedCodeSnippet, finalCode)
		} else {
			t.Logf("✅ Final response validation passed.")
		}
	})

	// 5. Evaluation Report: Go 测试框架本身就是我们的评估报告
	t.Log("Evaluation finished.")
}
