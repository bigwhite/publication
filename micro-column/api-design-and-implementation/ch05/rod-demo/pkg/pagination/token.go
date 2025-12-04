package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// PageToken 是我们在 Token 字符串中隐藏的结构
// 实际生产中，你可以加入 Salt 签名或加密，防止客户端伪造
type PageToken struct {
	Offset string `json:"o"` // 这里的 Offset 不是 SQL offset，而是"偏移的锚点值"(Cursor)
	// 可以扩展：Time int64 `json:"t"` // 如果按时间排序
}

// Encode 生成 next_page_token
func Encode(cursor string) string {
	if cursor == "" {
		return ""
	}
	t := PageToken{Offset: cursor}
	b, _ := json.Marshal(t)
	// 使用 URL 安全的 Base64 编码
	return base64.URLEncoding.EncodeToString(b)
}

// Decode 解析 page_token
func Decode(tokenStr string) (string, error) {
	if tokenStr == "" {
		return "", nil
	}
	b, err := base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		return "", errors.New("invalid page_token format")
	}
	var t PageToken
	if err := json.Unmarshal(b, &t); err != nil {
		return "", errors.New("invalid page_token payload")
	}
	return t.Offset, nil
}
