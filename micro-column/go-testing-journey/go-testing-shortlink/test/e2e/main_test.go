// test/e2e/main_test.go
//go:build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("启动 E2E 测试环境...")
	// 使用 --build 标志确保每次都使用最新的代码构建镜像
	cmdUp := exec.Command("docker-compose", "-f", "../../docker-compose.e2e.yml", "up", "-d", "--build")
	if output, err := cmdUp.CombinedOutput(); err != nil {
		fmt.Printf("无法启动 docker-compose 环境: %v\nOutput: %s\n", err, string(output))
		os.Exit(1)
	}

	// 更好的等待方式：轮询健康检查端点
	if err := waitForServices(); err != nil {
		fmt.Println("服务启动失败:", err)
		// 确保在退出前也尝试清理环境
		exec.Command("docker-compose", "-f", "../../docker-compose.e2e.yml", "down").Run()
		os.Exit(1)
	}
	fmt.Println("E2E 测试环境已就绪。")

	exitCode := m.Run()

	fmt.Println("销毁 E2E 测试环境...")
	cmdDown := exec.Command("docker-compose", "-f", "../../docker-compose.e2e.yml", "down")
	if err := cmdDown.Run(); err != nil {
		fmt.Printf("无法销毁 docker-compose 环境: %v\n", err)
	}

	os.Exit(exitCode)
}

// waitForServices 轮询 shortlink-app 的 /healthz 端点，直到它返回 200 OK
func waitForServices() error {
	const (
		maxRetries = 200
		retryDelay = 3 * time.Second
		appURL     = "http://localhost:8080/healthz"
	)

	fmt.Println("等待服务启动...")
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(appURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Println("服务已健康！")
			return nil // 成功
		}
		if err != nil {
			fmt.Printf("尝试 #%d: 连接失败: %v\n", i+1, err)
		} else {
			fmt.Printf("尝试 #%d: 服务未就绪，状态码: %d\n", i+1, resp.StatusCode)
		}
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("服务在 %d 次尝试后仍未就绪", maxRetries)
}

// TestShortlink_E2E_HappyPath 验证一个完整的用户旅程
func TestShortlink_E2E_HappyPath(t *testing.T) {
	appURL := "http://localhost:8080"
	originalURL := "https://www.very-long-url-for-e2e-testing.com"
	var shortCode string

	// 1. 场景一：用户成功创建一个短链接
	t.Run("Create a short link", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{"url": originalURL})
		// 注意：为了让 dummy user-service 通过，我们需要在请求中模拟一个有权限的用户
		req, _ := http.NewRequest("POST", appURL+"/api/links", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "user-with-permission") // 假设 handler 会读取这个头

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("创建请求失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("期望状态码 201, 得到了 %d", resp.StatusCode)
		}

		var createResp map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
			t.Fatalf("解析创建响应失败: %v", err)
		}
		
		shortCode = createResp["short_code"]
		if shortCode == "" { t.Fatal("响应中没有找到 short_code") }
		t.Logf("成功创建短链接，短码为: %s", shortCode)
	})

	// 2. 场景二：用户通过短链接访问，被成功重定向
	t.Run("Redirect via short link", func(t *testing.T) {
		if shortCode == "" { t.Skip("由于创建失败，跳过重定向测试") }
		
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		redirectURL := fmt.Sprintf("%s/%s", appURL, shortCode)
		resp, err := client.Get(redirectURL)
		if err != nil { t.Fatalf("重定向请求失败: %v", err) }
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("期望状态码 301 或 302, 得到了 %d", resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		if location != originalURL {
			t.Fatalf("期望重定向到 %s, 但实际为 %s", originalURL, location)
		}
		t.Logf("成功验证重定向，Location: %s", location)
	})

	// 3. 场景三：用户查询访问统计
	t.Run("Get link statistics", func(t *testing.T) {
		if shortCode == "" { t.Skip("由于创建失败，跳过统计测试") }
		
		// 等待异步的 Redis 计数器生效
		time.Sleep(1 * time.Second) 
		
		statsURL := fmt.Sprintf("%s/api/links/%s/stats", appURL, shortCode)
		resp, err := http.Get(statsURL)
		if err != nil { t.Fatalf("统计请求失败: %v", err) }
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("期望状态码 200, 得到了 %d", resp.StatusCode)
		}

		var statsResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&statsResp); err != nil {
			t.Fatalf("解析统计响应失败: %v", err)
		}
		
		visits, ok := statsResp["visits"].(float64)
		if !ok || visits < 1 { t.Fatalf("期望访问次数至少为 1, 得到了 %v", statsResp["visits"]) }
		t.Logf("成功获取统计信息，访问次数为: %.0f", visits)
	})
}
