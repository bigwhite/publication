package feishu

import (
	"fmt"
	"log"
	"regexp"
	"sync"
)

type ApprovalResult struct {
	Allowed bool
	Reason  string
}

type ApprovalManager struct {
	mu           sync.RWMutex
	pendingTasks map[string]chan ApprovalResult
}

var GlobalApprovalMgr = &ApprovalManager{
	pendingTasks: make(map[string]chan ApprovalResult),
}

func (m *ApprovalManager) WaitForApproval(taskID string, toolName string, args string, reporter *FeishuReporter) (bool, string) {
	ch := make(chan ApprovalResult, 1)

	m.mu.Lock()
	m.pendingTasks[taskID] = ch
	m.mu.Unlock()

	noticeMsg := fmt.Sprintf(`⚠️ **高危操作审批请求**
Agent 试图执行以下动作:
- 工具: %s
- 参数: %s

任务 ID: **%s**

👉 请回复 "approve %s" 或 "reject %s" 决定是否放行。`, toolName, args, taskID, taskID, taskID)

	if reporter != nil {
		reporter.sendMsg(noticeMsg)
	} else {
		fmt.Printf("\n\033[31m[需要审批 TaskID: %s]\033[0m %s\n", taskID, noticeMsg)
	}

	log.Printf("[Approval] 发送审批请求 (TaskID: %s)，协程挂起等待...\n", taskID)

	// 阻塞当前 Goroutine
	result := <-ch

	m.mu.Lock()
	delete(m.pendingTasks, taskID)
	m.mu.Unlock()

	return result.Allowed, result.Reason
}

func (m *ApprovalManager) ResolveApproval(taskID string, allowed bool, reason string) {
	m.mu.RLock()
	ch, exists := m.pendingTasks[taskID]
	m.mu.RUnlock()

	if exists {
		log.Printf("[Approval] 收到飞书审批结果 (TaskID: %s, Allowed: %v)\n", taskID, allowed)
		ch <- ApprovalResult{Allowed: allowed, Reason: reason}
	}
}

func IsDangerousCommand(toolName string, args string) bool {
	if toolName != "bash" && toolName != "write_file" && toolName != "edit_file" {
		return false
	}

	if toolName == "bash" {
		dangerousPatterns := []string{`rm\s+-r`, `sudo\s+`, `drop\s+`, `>.*\.go`}
		for _, p := range dangerousPatterns {
			if matched, _ := regexp.MatchString(p, args); matched {
				return true
			}
		}
	}
	return false
}
