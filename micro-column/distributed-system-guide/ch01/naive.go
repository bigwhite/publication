package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// 假设这是我们的库存服务
	resp, err := http.Get("http://inventory-service/check/product-123")
	if err != nil {
		// 在单体世界，这通常意味着一个致命的、需要程序员介入的Bug。
		// 在分布式世界，这可能只是网络的一次偶然抖动。
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Inventory check result:", string(body))
}
