package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// ChatMessage 定义了对话中单条消息的结构
type ChatMessage struct {
	Role    string `json:"role"`    // 角色：system, user, 或 assistant
	Content string `json:"content"` // 消息内容
}

// ChatCompletionRequest 定义了发送给聊天补全API的请求体结构
type ChatCompletionRequest struct {
	Model    string        `json:"model"`            // 使用的模型ID
	Messages []ChatMessage `json:"messages"`         // 对话消息列表
	Stream   bool          `json:"stream,omitempty"` // 是否流式响应，omitempty表示如果为false则不序列化此字段
	// 可以添加其他参数如 Temperature, MaxTokens 等
}

// ResponseChoice 定义了API响应中单个选择项的结构
type ResponseChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`       // 助手返回的消息
	FinishReason string      `json:"finish_reason"` // 结束原因，如 "stop", "length"
}

// ChatCompletionResponse 定义了聊天补全API响应体的主结构
type ChatCompletionResponse struct {
	ID      string           `json:"id"`      // 响应ID
	Object  string           `json:"object"`  // 对象类型，如 "chat.completion"
	Created int64            `json:"created"` // 创建时间戳
	Model   string           `json:"model"`   // 使用的模型ID
	Choices []ResponseChoice `json:"choices"` // 包含一个或多个回复选项的列表
	// Usage   UsageInfo        `json:"usage"`   // Token使用情况 (在非流式时通常有)
}

// UsageInfo 定义了Token使用统计的结构 (为简化，此处可以先不详细定义)
// type UsageInfo struct { ... }

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("错误：环境变量 OPENAI_API_KEY 未设置。")
	}

	// DeepSeek API端点 (OpenAI兼容)
	apiURL := "https://api.deepseek.com/v1/chat/completions" // 使用与OpenAI SDK默认更接近的路径
	modelID := "deepseek-chat"                               // DeepSeek提供的聊天模型

	// 构造请求体
	requestPayload := ChatCompletionRequest{
		Model: modelID,
		Messages: []ChatMessage{
			{Role: "system", Content: "你是一个乐于助人的AI助手。"},
			{Role: "user", Content: "你好AI，请问Go语言是什么时候发布的？"},
		},
		Stream: false, // 我们先尝试非流式
	}

	requestBodyBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Fatalf("序列化请求体失败: %v", err)
	}

	// 创建HTTP请求
	// 使用 context.WithTimeout 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // 确保在函数结束时取消上下文，释放资源

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatalf("创建HTTP请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	fmt.Println("正在发送请求到AI模型...")
	client := &http.Client{} // 可以配置client的超时等参数
	resp, err := client.Do(req)
	if err != nil {
		// 检查上下文是否已超时或被取消
		if errors.Is(err, context.DeadlineExceeded) {
			log.Fatalf("请求超时: %v", err)
		}
		log.Fatalf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取并处理响应
	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应体失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(responseBodyBytes))
	}

	var chatResponse ChatCompletionResponse
	if err := json.Unmarshal(responseBodyBytes, &chatResponse); err != nil {
		log.Fatalf("反序列化响应JSON失败: %v\n原始响应: %s", err, string(responseBodyBytes))
	}

	// 打印AI的回答
	if len(chatResponse.Choices) > 0 {
		assistantMessage := chatResponse.Choices[0].Message
		fmt.Printf("AI助手的回答 (模型: %s):\n%s\n", chatResponse.Model, assistantMessage.Content)
		fmt.Printf("(结束原因: %s)\n", chatResponse.Choices[0].FinishReason)
	} else {
		fmt.Println("AI未返回任何有效的回答选项。")
	}
}
