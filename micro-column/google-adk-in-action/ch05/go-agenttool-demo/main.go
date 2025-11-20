package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.0-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// --- 1. 定义三个“专家”子 Agent ---

	// Code Writer Agent: 接收用户请求，生成代码，并将结果存入 state["generated_code"]
	codeWriterAgent, err := llmagent.New(llmagent.Config{
		Name:  "CodeWriterAgent",
		Model: model,
		Instruction: `You are a Python Code Generator.
Based *only* on the user's request, write Python code that fulfills the requirement.
Output *only* the complete Python code block, enclosed in triple backticks ('''python ... ''').
Do not add any other text before or after the code block.`,
		Description: "Writes initial Python code based on a specification.",
		OutputKey:   "generated_code", // 将输出存入 state["generated_code"]
	})
	if err != nil {
		log.Fatalf("failed to create codeWriterAgent: %s", err)
	}

	// Code Reviewer Agent: 从 state 中读取代码，提供反馈，并将结果存入 state["review_comments"]
	codeReviewerAgent, err := llmagent.New(llmagent.Config{
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
	if err != nil {
		log.Fatalf("failed to create codeReviewerAgent: %s", err)
	}

	// Code Refactorer Agent: 从 state 中读取代码和反馈，进行重构
	codeRefactorerAgent, err := llmagent.New(llmagent.Config{
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
	if err != nil {
		log.Fatalf("failed to create codeRefactorerAgent: %s", err)
	}

	// --- 2. 定义我们的“代码流水线” Agent (作为 Sub-agent 集合) ---
	codePipelineAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name: "CodePipelineAgent",
			SubAgents: []agent.Agent{ // **定义执行顺序**
				codeWriterAgent,
				codeReviewerAgent,
				codeRefactorerAgent,
			},
			Description: "Executes a sequence of code writing, reviewing, and refactoring.",
		},
	})
	if err != nil {
		log.Fatalf("failed to create codePipelineAgent: %s", err)
	}

	// --- 2. 将“代码流水线” Agent 包装成一个 Tool ---
	pipelineTool := agenttool.New(codePipelineAgent, nil)

	// --- 3. 创建更高层次的“项目经理” Agent ---
	projectManagerAgent, err := llmagent.New(llmagent.Config{
		Name:  "ProjectManagerAgent",
		Model: model,
		Tools: []tool.Tool{pipelineTool},
		Instruction: `You are a senior project manager. When the user provides a software requirement,
                      delegate the entire implementation process to the 'CodePipelineAgent' tool.`,
	})
	if err != nil {
		log.Fatalf("failed to create projectManagerAgent: %s", err)
	}

	// --- 4. 启动 Launcher，并将“项目经理”设为根 Agent ---
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(projectManagerAgent),
	}
	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
