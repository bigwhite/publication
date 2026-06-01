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

// internal/feishu/approval.go (局部修正)

// IsDangerousCommand 简单的正则检查黑名单，判断该工具调用是否需要触发人类审批
func IsDangerousCommand(toolName string, args string) bool {
	// 白名单放行：对于纯读取工具，默认 YOLO 模式，全部放行
	if toolName == "read_file" {
		return false
	}

	// 【剧本设定】：在生产服务器的 AgentOps 场景下，修改任何文件都是高危操作！
	// 我们不允许 Agent 擅自使用 write_file 覆写文件，或使用 edit_file 篡改代码。
	if toolName == "write_file" || toolName == "edit_file" {
		return true
	}

	// 针对 bash 的高危模式匹配
	if toolName == "bash" {
		// 危险指令特征库 (模拟真实的运维黑名单)
		dangerousPatterns := []string{
			`rm\s+-r`,      // 级联删除
			`sudo\s+`,      // 提权操作
			`drop\s+`,      // 数据库危险命令
			`>.*\.go`,      // 恶意覆盖源代码
			`nginx\s+-s`,   // 【针对第 22 讲剧本】：拦截 Nginx 服务重启或停止
			`systemctl\s+`, // 拦截系统级服务管理
			`kill\s+`,      // 拦截杀进程操作
		}

		for _, p := range dangerousPatterns {
			if matched, _ := regexp.MatchString(p, args); matched {
				return true // 命中任何一条黑名单，必须挂起审批
			}
		}
	}

	// 如果没有命中高危特征，默认放行 (例如简单的 ls -la, tail -n 50 等探测命令)
	return false
}
