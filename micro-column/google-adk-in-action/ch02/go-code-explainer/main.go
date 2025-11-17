package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// 1. 初始化模型
	model, err := gemini.NewModel(ctx, "gemini-2.0-flash-lite", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// 2. 创建 LLMAgent
	codeExplainerAgent, err := llmagent.New(llmagent.Config{
		Name:        "go_code_explainer",
		Description: "An agent that explains snippets of Go code in a clear and concise way.",
		Model:       model,
		Instruction: "You are an expert Go programmer and a patient teacher. When you are given a snippet of Go code, your sole purpose is to explain what it does, line by line, in plain English. Your explanation should be easy for a beginner to understand. Do not get sidetracked. Do not answer questions that are not about the provided code snippet.",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// 3. 配置并启动 Launcher
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(codeExplainerAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
