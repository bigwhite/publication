package main

func login(user string) bool {
    // 检查用户名
    if user == "admin" {
        return true
    }
    return false
}
