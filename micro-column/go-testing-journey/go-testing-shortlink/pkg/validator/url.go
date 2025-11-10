package validator

import (
	"net/url"
)

// IsValidURL 检查一个字符串是否是有效的 URL
func IsValidURL(rawURL string) bool {
	// 使用 net/url 包来解析
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	// 必须包含 http 或 https scheme
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// 必须包含 host
	if u.Host == "" {
		return false
	}

	return true
}
