package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/users/", permissionsHandler)
	fmt.Println("Dummy User Service 启动于 :8090...")
	http.ListenAndServe(":8090", nil)
}

func permissionsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "users" || parts[2] != "permissions" {
		http.NotFound(w, r)
		return
	}
	userID := parts[1]

	// 这是一个非常简单的 dummy 逻辑
	// 在 E2E 测试中，我们只需要它能为特定的测试用户返回正确的结果
	canCreate := false
	if userID == "user-with-permission" {
		canCreate = true
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"can_create": canCreate})
}
