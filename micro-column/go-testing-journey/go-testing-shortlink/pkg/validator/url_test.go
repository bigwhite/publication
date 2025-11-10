//go:build unit

package validator

import "testing"

func TestIsValidURL(t *testing.T) {
	// 使用“表驱动测试”法，这是 Go 社区的最佳实践
	testCases := []struct {
		name     string // 子测试的名称
		inputURL string // 输入
		want     bool   // 期望的输出
	}{
		{
			name:     "合法的 HTTP URL",
			inputURL: "http://example.com",
			want:     true,
		},
		{
			name:     "合法的 HTTPS URL",
			inputURL: "https://example.com/path?query=1",
			want:     true,
		},
		{
			name:     "缺少 Scheme",
			inputURL: "example.com",
			want:     false,
		},
		{
			name:     "不支持的 Scheme (ftp)",
			inputURL: "ftp://example.com",
			want:     false,
		},
		{
			name:     "缺少 Host",
			inputURL: "http://",
			want:     false,
		},
		{
			name:     "空字符串",
			inputURL: "",
			want:     false,
		},
	}

	for _, tc := range testCases {
		// t.Run 让每个测试用例成为一个独立的子测试
		t.Run(tc.name, func(t *testing.T) {
			got := IsValidURL(tc.inputURL)
			if got != tc.want {
				t.Errorf("IsValidURL(%q) = %v; want %v", tc.inputURL, got, tc.want)
			}
		})
	}
}
