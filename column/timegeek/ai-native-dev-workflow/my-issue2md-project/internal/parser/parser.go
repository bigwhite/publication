package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// ResourceURL 表示解析后的GitHub资源URL
type ResourceURL struct {
	Type   string // "issue", "pull", "discussion"
	Owner  string
	Repo   string
	Number int
	URL    string // 原始URL
}

// URLParser 定义URL解析接口
type URLParser interface {
	// Parse 解析GitHub URL
	Parse(rawURL string) (*ResourceURL, error)

	// Validate 验证URL格式
	Validate(rawURL string) error

	// SupportedTypes 返回支持的资源类型
	SupportedTypes() []string
}

// GitHubURLParser GitHub URL解析器实现
type GitHubURLParser struct{}

// NewURLParser 创建新的URL解析器
func NewURLParser() *GitHubURLParser {
	return &GitHubURLParser{}
}

// Parse 解析GitHub URL
func (p *GitHubURLParser) Parse(rawURL string) (*ResourceURL, error) {
	// 基本URL格式检查
	if rawURL == "" {
		return nil, fmt.Errorf("empty URL")
	}

	// 检查是否为有效的GitHub URL
	if !strings.HasPrefix(rawURL, "https://github.com/") {
		return nil, fmt.Errorf("invalid URL: %w", fmt.Errorf("not a GitHub URL"))
	}

	// 分割和验证路径
	owner, repo, resourceType, numberStr, err := p.splitAndValidatePath(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// 解析资源类型和编号
	resultType, err := p.parseResourceType(resourceType)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// 验证并解析数字
	number, err := p.parseNumber(numberStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// 构造返回结果
	return &ResourceURL{
		Type:   resultType,
		Owner:  owner,
		Repo:   repo,
		Number: number,
		URL:    rawURL,
	}, nil
}

// Validate 验证URL格式
func (p *GitHubURLParser) Validate(rawURL string) error {
	_, err := p.Parse(rawURL)
	return err
}

// splitAndValidatePath 分割URL路径并进行基础验证
func (p *GitHubURLParser) splitAndValidatePath(rawURL string) (owner, repo, resourceType, numberStr string, err error) {
	// 移除协议前缀进行路径解析
	path := strings.TrimPrefix(rawURL, "https://github.com/")

	// 分割路径部分
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 3 {
		// 检查是否是仓库主页（只有 owner/repo）
		if len(pathParts) == 2 && pathParts[0] != "" && pathParts[1] != "" {
			return "", "", "", "", fmt.Errorf("unsupported URL type")
		}
		return "", "", "", "", fmt.Errorf("insufficient path components")
	}

	owner = pathParts[0]
	repo = pathParts[1]

	if owner == "" || repo == "" {
		return "", "", "", "", fmt.Errorf("missing owner or repository")
	}

	if len(pathParts) < 4 {
		return "", "", "", "", fmt.Errorf("missing number component")
	}

	resourceType = pathParts[2]
	numberStr = pathParts[3]

	return owner, repo, resourceType, numberStr, nil
}

// parseResourceType 解析资源类型
func (p *GitHubURLParser) parseResourceType(resourceType string) (string, error) {
	switch resourceType {
	case "issues":
		return "issue", nil
	case "pull":
		return "pull", nil
	case "discussions":
		return "discussion", nil
	default:
		return "", fmt.Errorf("unsupported URL type: %s", resourceType)
	}
}

// parseNumber 解析并验证数字
func (p *GitHubURLParser) parseNumber(numberStr string) (int, error) {
	if numberStr == "" {
		return 0, fmt.Errorf("missing number component")
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number format '%s'", numberStr)
	}

	if number <= 0 {
		return 0, fmt.Errorf("number must be positive, got %d", number)
	}

	return number, nil
}

// SupportedTypes 返回支持的资源类型
func (p *GitHubURLParser) SupportedTypes() []string {
	return []string{"issue", "pull", "discussion"}
}