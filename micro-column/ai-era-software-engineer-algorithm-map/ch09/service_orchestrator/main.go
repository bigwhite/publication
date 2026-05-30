package main

import (
	"fmt"
	"strings"
	"time"
)

// ServiceConfig 定义服务及其依赖
type ServiceConfig struct {
	Name string
	Deps []string
}

// Orchestrator 编排器
type Orchestrator struct {
	services []ServiceConfig
}

func NewOrchestrator(services []ServiceConfig) *Orchestrator {
	return &Orchestrator{services: services}
}

// ResolveStartupOrder 计算启动顺序
// 返回一个二维切片，外层是批次，内层是该批次并行启动的服务名
func (o *Orchestrator) ResolveStartupOrder() ([][]string, error) {
	// 1. 建图
	graph := make(map[string][]string) // dep -> [service...]
	inDegree := make(map[string]int)
	allServices := make(map[string]bool) // 记录所有出现过的服务名

	// 初始化入度表（防止漏掉没有依赖的服务）
	for _, svc := range o.services {
		allServices[svc.Name] = true
		if _, ok := inDegree[svc.Name]; !ok {
			inDegree[svc.Name] = 0
		}
		for _, dep := range svc.Deps {
			allServices[dep] = true
			graph[dep] = append(graph[dep], svc.Name)
			inDegree[svc.Name]++
		}
	}

	// 2. 找入口
	queue := []string{}
	for name := range allServices {
		// 注意：配置文件可能没写基础服务(如MySQL)的配置项，默认它们入度为0
		if inDegree[name] == 0 {
			queue = append(queue, name)
		}
	}

	var startupPlan [][]string
	processedCount := 0

	// 3. 拓扑排序 (分层 BFS)
	for len(queue) > 0 {
		levelSize := len(queue)
		var currentBatch []string

		// 这一层的服务都可以并行启动
		for i := 0; i < levelSize; i++ {
			svc := queue[0]
			queue = queue[1:]
			currentBatch = append(currentBatch, svc)
			processedCount++

			// 解锁下游
			for _, nextSvc := range graph[svc] {
				inDegree[nextSvc]--
				if inDegree[nextSvc] == 0 {
					queue = append(queue, nextSvc)
				}
			}
		}
		startupPlan = append(startupPlan, currentBatch)
	}

	// 4. 环检测
	if processedCount != len(allServices) {
		return nil, fmt.Errorf("cyclic dependency detected! processed %d/%d services", processedCount, len(allServices))
	}

	return startupPlan, nil
}

func main() {
	// 定义依赖关系
	config := []ServiceConfig{
		{"Web-Backend", []string{"DB-Proxy", "Redis"}},
		{"DB-Proxy", []string{"MySQL-Master", "MySQL-Slave"}},
		{"Recommend-Sys", []string{"Web-Backend", "AI-Model"}},
		// MySQL, Redis, AI-Model 是基础服务，没有依赖
	}

	orc := NewOrchestrator(config)
	batches, err := orc.ResolveStartupOrder()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("--- Service Startup Plan ---")
	for i, batch := range batches {
		fmt.Printf("Batch %d: Starting [%s]...\n", i+1, strings.Join(batch, ", "))
		// 模拟并行启动耗时
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("        -> %d services started.\n", len(batch))
	}
	fmt.Println("--- All Systems Operational ---")
}
