package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// 设计哲学 1: 不要无限信任网络。为你的每次调用设定一个期望的时间边界。
	client := http.Client{
		Timeout: 2 * time.Second, // 设定2秒超时
	}

	// 设计哲学 2: 失败是常态，必须有预案。
	var body []byte
	var err error

	// 尝试最多3次
	for i := 0; i < 3; i++ {
		resp, err := client.Get("http://inventory-service/check/product-123")
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ = io.ReadAll(resp.Body)
				log.Println("Success on attempt", i+1)
				break // 成功，跳出循环
			}
			err = fmt.Errorf("bad status: %s", resp.Status)
		}

		log.Printf("Attempt %d failed: %v. Retrying in 1 second...", i+1, err)
		time.Sleep(1 * time.Second) // 简单的重试间隔
	}

	if body == nil {
		log.Fatalf("Failed to check inventory after multiple retries: %v", err)
		// 在这里，我们可以触发降级逻辑，比如返回一个默认的库存状态
		return
	}

	fmt.Println("Inventory check result:", string(body))
}
