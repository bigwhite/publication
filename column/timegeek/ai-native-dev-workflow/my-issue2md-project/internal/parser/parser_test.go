package parser

import (
	"strings"
	"testing"
)

// containsError 检查错误信息是否包含期望的内容
func containsError(actual, expected string) bool {
	// 如果期望的字符串在错误信息中，则返回true
	return strings.Contains(actual, expected)
}

func TestParse(t *testing.T) {
	// 创建解析器实例
	parser := NewURLParser()

	// 表格驱动测试用例
	tests := []struct {
		name          string
		rawURL        string
		wantType      string
		wantOwner     string
		wantRepo      string
		wantNumber    int
		wantError     bool
		errorContains string
	}{
		// 测试用例1: 合法的 Issue URL
		{
			name:       "Valid Issue URL",
			rawURL:     "https://github.com/facebook/react/issues/12345",
			wantType:   "issue",
			wantOwner:  "facebook",
			wantRepo:   "react",
			wantNumber: 12345,
			wantError:  false,
		},

		// 测试用例2: 合法的 Pull Request URL
		{
			name:       "Valid Pull Request URL",
			rawURL:     "https://github.com/microsoft/vscode/pull/9876",
			wantType:   "pull",
			wantOwner:  "microsoft",
			wantRepo:   "vscode",
			wantNumber: 9876,
			wantError:  false,
		},

		// 测试用例3: 合法的 Discussion URL
		{
			name:       "Valid Discussion URL",
			rawURL:     "https://github.com/github/roadmap/discussions/543",
			wantType:   "discussion",
			wantOwner:  "github",
			wantRepo:   "roadmap",
			wantNumber: 543,
			wantError:  false,
		},

		// 测试用例4: 无效的 URL（格式错误）
		{
			name:          "Invalid URL format",
			rawURL:        "https://example.com/not/github/url/123",
			wantError:     true,
			errorContains: "invalid URL",
		},

		// 测试用例5: 不支持的 URL 类型（仓库主页）
		{
			name:          "Unsupported URL type - repository homepage",
			rawURL:        "https://github.com/facebook/react",
			wantError:     true,
			errorContains: "unsupported URL type",
		},

		// 额外测试用例：空的URL
		{
			name:          "Empty URL",
			rawURL:        "",
			wantError:     true,
			errorContains: "empty URL",
		},

		// 额外测试用例：URL缺少number部分
		{
			name:          "URL missing number",
			rawURL:        "https://github.com/facebook/react/issues/",
			wantError:     true,
			errorContains: "invalid URL",
		},

		// 额外测试用例：URL包含非数字ID
		{
			name:          "URL with non-numeric ID",
			rawURL:        "https://github.com/facebook/react/issues/abc",
			wantError:     true,
			errorContains: "invalid URL",
		},

		// 额外测试用例：其他GitHub域名格式
		{
			name:          "Invalid GitHub domain",
			rawURL:        "https://gist.github.com/user/repo/123",
			wantError:     true,
			errorContains: "invalid URL",
		},

		// 额外测试用例：HTTPS协议缺失
		{
			name:          "Missing HTTPS protocol",
			rawURL:        "http://github.com/facebook/react/issues/123",
			wantError:     true,
			errorContains: "invalid URL",
		},
	}

	// 执行表格驱动测试
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用 Parse 函数
			got, err := parser.Parse(tt.rawURL)

			// 检查错误情况
			if tt.wantError {
				if err == nil {
					t.Errorf("Parse() expected error containing '%s', but got nil", tt.errorContains)
					return
				}
				if tt.errorContains != "" {
					errMsg := err.Error()
					// 检查错误信息是否包含期望的内容
					if !containsError(errMsg, tt.errorContains) {
						t.Errorf("Parse() error = %v, expected to contain '%s'", errMsg, tt.errorContains)
						return
					}
				}
				// 如果期望错误，不再检查返回值
				return
			}

			// 不期望错误但发生了错误
			if err != nil {
				t.Errorf("Parse() unexpected error = %v", err)
				return
			}

			// 检查返回的 ResourceURL
			if got == nil {
				t.Error("Parse() returned nil, expected ResourceURL")
				return
			}

			// 验证各个字段
			if got.Type != tt.wantType {
				t.Errorf("Parse().Type = %v, want %v", got.Type, tt.wantType)
			}
			if got.Owner != tt.wantOwner {
				t.Errorf("Parse().Owner = %v, want %v", got.Owner, tt.wantOwner)
			}
			if got.Repo != tt.wantRepo {
				t.Errorf("Parse().Repo = %v, want %v", got.Repo, tt.wantRepo)
			}
			if got.Number != tt.wantNumber {
				t.Errorf("Parse().Number = %v, want %v", got.Number, tt.wantNumber)
			}
			if got.URL != tt.rawURL {
				t.Errorf("Parse().URL = %v, want %v", got.URL, tt.rawURL)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	parser := NewURLParser()

	tests := []struct {
		name          string
		rawURL        string
		wantError     bool
		errorContains string
	}{
		{
			name:      "Valid Issue URL",
			rawURL:    "https://github.com/facebook/react/issues/12345",
			wantError: false,
		},
		{
			name:          "Invalid URL format",
			rawURL:        "https://example.com/not/github/url/123",
			wantError:     true,
			errorContains: "invalid URL",
		},
		{
			name:      "Valid Pull Request URL",
			rawURL:    "https://github.com/microsoft/vscode/pull/9876",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.Validate(tt.rawURL)

			if tt.wantError {
				if err == nil {
					t.Errorf("Validate() expected error containing '%s', but got nil", tt.errorContains)
					return
				}
				if tt.errorContains != "" {
					errMsg := err.Error()
					if !containsError(errMsg, tt.errorContains) {
						t.Errorf("Validate() error = %v, expected to contain '%s'", errMsg, tt.errorContains)
						return
					}
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestSupportedTypes(t *testing.T) {
	parser := NewURLParser()

	got := parser.SupportedTypes()

	// 验证返回支持的类型
	expected := []string{"issue", "pull", "discussion"}
	if len(got) != len(expected) {
		t.Errorf("SupportedTypes() length = %d, want %d", len(got), len(expected))
	}
	for i, v := range expected {
		if i >= len(got) || got[i] != v {
			t.Errorf("SupportedTypes() = %v, want %v", got, expected)
			break
		}
	}
}