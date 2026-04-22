package main

import "fmt"

func main() {
    // 启动服务器
    fmt.Println("Server is starting on port 8080...")
    
    // TODO: 增加鉴权逻辑
    if true {
        fmt.Println("No auth, everyone can access.")
    }
}
