package parser

import (
	"github.com/bigwhite/my-issue2md/internal/github"
)

// Parser 定义Markdown解析器接口
type Parser interface {
	Parse(issue *github.Issue, comments []*github.Comment) (*MarkdownDocument, error)
}

// MarkdownParser Markdown解析器实现
type MarkdownParser struct {
	options *Options
}

// Options 解析器配置选项
type Options struct {
	IncludeComments    bool `json:"include_comments"`
	IncludeMetadata    bool `json:"include_metadata"`
	IncludeTimestamps  bool `json:"include_timestamps"`
	IncludeUserLinks   bool `json:"include_user_links"`
	EmojisEnabled      bool `json:"emojis_enabled"`
	PreserveLineBreaks bool `json:"preserve_line_breaks"`
}

// DefaultOptions 返回默认解析器选项
func DefaultOptions() *Options {
	return &Options{
		IncludeComments:    true,
		IncludeMetadata:    true,
		IncludeTimestamps:  true,
		IncludeUserLinks:   true,
		EmojisEnabled:      true,
		PreserveLineBreaks: true,
	}
}

// NewParser 创建新的Markdown解析器
func NewParser(opts *Options) *MarkdownParser {
	if opts == nil {
		opts = DefaultOptions()
	}
	return &MarkdownParser{
		options: opts,
	}
}

// MarkdownDocument 表示解析后的Markdown文档
type MarkdownDocument struct {
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

// IssueMetadata Issue元数据
type IssueMetadata struct {
	Number     int    `json:"number"`
	State      string `json:"state"`
	Author     string `json:"author"`
	Assignees  string `json:"assignees"`
	Labels     string `json:"labels"`
	Milestone  string `json:"milestone"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	ClosedAt   string `json:"closed_at"`
	URL        string `json:"url"`
	Repository string `json:"repository"`
}

// CommentMetadata 评论元数据
type CommentMetadata struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	URL       string `json:"url"`
}

// ProcessingError 处理错误
type ProcessingError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// Error 实现error接口
func (e *ProcessingError) Error() string {
	return e.Message
}

// NewProcessingError 创建处理错误
func NewProcessingError(message, code, details string) *ProcessingError {
	return &ProcessingError{
		Message: message,
		Code:    code,
		Details: details,
	}
}