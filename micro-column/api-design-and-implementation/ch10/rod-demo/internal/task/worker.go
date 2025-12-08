package task

import (
	"fmt"
	"rod-demo/internal/domain"
	"sync"
	"time"
)

// OpStore 模拟全局存储 (生产环境应为 Redis)
var (
	OpStore = make(map[string]*domain.Operation)
	mu      sync.RWMutex
)

// GetOperation 安全地获取操作状态
func GetOperation(id string) *domain.Operation {
	mu.RLock()
	defer mu.RUnlock()
	if op, ok := OpStore[id]; ok {
		// 返回副本以避免并发读写冲突
		val := *op
		return &val
	}
	return nil
}

// StartImageGeneration 启动异步任务
func StartImageGeneration(opID string, prompt string) {
	// 1. 初始化任务状态
	mu.Lock()
	OpStore[opID] = &domain.Operation{
		ID:   opID,
		Done: false,
		Metadata: domain.OperationMetadata{
			Status:   "Queued",
			Progress: 0,
		},
	}
	mu.Unlock()

	// 2. 启动协程模拟耗时操作 (AI 推理)
	go func() {
		// 模拟 5 秒处理，每秒更新进度
		for i := 1; i <= 5; i++ {
			time.Sleep(1 * time.Second)
			mu.Lock()
			if op, ok := OpStore[opID]; ok {
				op.Metadata = domain.OperationMetadata{
					Status:   "Processing",
					Progress: i * 20,
				}
			}
			mu.Unlock()
		}

		// 3. 任务完成
		mu.Lock()
		if op, ok := OpStore[opID]; ok {
			op.Done = true
			op.Metadata = domain.OperationMetadata{Progress: 100, Status: "Done"}
			// 模拟生成了一张图片 URL
			op.Response = map[string]string{
				"image_url": fmt.Sprintf("https://cdn.example.com/images/%s.png", opID),
			}
		}
		mu.Unlock()
	}()
}
