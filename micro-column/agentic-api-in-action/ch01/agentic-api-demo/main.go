package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// AgenticRequest 代表了 AI 智能体发来的标准任务请求
type AgenticRequest struct {
	// 明确的意图动作
	Action string `json:"action"`
	// 动作的目标上下文 (例如文档ID)
	ContextID string `json:"context_id"`
	// 动作需要的特定参数
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// AgenticResponse 代表了返回给 AI 的标准结构化响应
type AgenticResponse struct {
	Status  string      `json:"status"` // SUCCESS, FAILED, REQUIRE_CONFIRM
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
}

func main() {
	// 定义一个面向动作的路由前缀
	http.HandleFunc("/api/v1/actions", actionHandler)

	fmt.Println("Agentic API Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// actionHandler 充当了“任务调度中心”
func actionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed for actions", http.StatusMethodNotAllowed)
		return
	}

	var req AgenticRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, http.StatusBadRequest, "FAILED", nil, "Invalid JSON payload")
		return
	}

	// 核心：基于 Action (动词) 进行路由分发，而不是基于资源名词
	switch strings.ToUpper(req.Action) {
	case "SUMMARIZE":
		handleSummarize(w, req)
	case "TRANSLATE":
		// handleTranslate(w, req)
		sendResponse(w, http.StatusNotImplemented, "FAILED", nil, "Action TRANSLATE not implemented yet")
	default:
		sendResponse(w, http.StatusBadRequest, "FAILED", nil, fmt.Sprintf("Unknown action: %s", req.Action))
	}
}

// handleSummarize 处理具体的总结任务
func handleSummarize(w http.ResponseWriter, req AgenticRequest) {
	docID := req.ContextID
	if docID == "" {
		sendResponse(w, http.StatusBadRequest, "FAILED", nil, "context_id (Document ID) is required")
		return
	}

	// 解析可选参数 (Agentic API 应该允许 AI 传入控制参数)
	maxLength := 100 // 默认值
	if ml, ok := req.Parameters["max_length"].(float64); ok {
		maxLength = int(ml)
	}

	// 模拟从数据库获取文档并进行总结的复杂逻辑
	log.Printf("Executing SUMMARIZE for doc: %s, max length: %d\n", docID, maxLength)

	// 模拟生成的摘要
	mockSummary := fmt.Sprintf("这是关于文档 %s 的核心总结，长度被限制在 %d 字以内：Agentic API 是未来的趋势。", docID, maxLength)

	// 返回标准化响应
	sendResponse(w, http.StatusOK, "SUCCESS", mockSummary, "Document summarized successfully")
}

// 统一的响应封装助手
func sendResponse(w http.ResponseWriter, statusCode int, status string, result interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := AgenticResponse{
		Status:  status,
		Result:  result,
		Message: message,
	}
	json.NewEncoder(w).Encode(resp)
}
