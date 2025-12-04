package converter

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/bigwhite/my-issue2md/internal/parser"
)

// Converter 定义转换器接口
type Converter interface {
	Convert(doc *parser.MarkdownDocument) ([]byte, error)
}

// MarkdownConverter Markdown转换器
type MarkdownConverter struct {
	options *ConverterOptions
}

// HTMLConverter HTML转换器
type HTMLConverter struct {
	options *ConverterOptions
}

// JSONConverter JSON转换器
type JSONConverter struct {
	options *ConverterOptions
}

// ConverterOptions 转换器选项
type ConverterOptions struct {
	EnableSyntaxHighlighting bool     `json:"enable_syntax_highlighting"`
	EnableTableOfContents    bool     `json:"enable_table_of_contents"`
	EnableEmojis            bool     `json:"enable_emojis"`
	EnableLinks             bool     `json:"enable_links"`
	EnableImages            bool     `json:"enable_images"`
	EnableCodeBlocks        bool     `json:"enable_code_blocks"`
	EnableTables            bool     `json:"enable_tables"`
	EnableStrikethrough     bool     `json:"enable_strikethrough"`
	EnableTaskLists         bool     `json:"enable_task_lists"`
	CustomCSS               string   `json:"custom_css"`
	Template                string   `json:"template"`
}

// DefaultConverterOptions 返回默认转换器选项
func DefaultConverterOptions() *ConverterOptions {
	return &ConverterOptions{
		EnableSyntaxHighlighting: true,
		EnableTableOfContents:    false,
		EnableEmojis:            true,
		EnableLinks:             true,
		EnableImages:            true,
		EnableCodeBlocks:        true,
		EnableTables:            true,
		EnableStrikethrough:     true,
		EnableTaskLists:         true,
	}
}

// NewMarkdownConverter 创建Markdown转换器
func NewMarkdownConverter(opts *ConverterOptions) *MarkdownConverter {
	if opts == nil {
		opts = DefaultConverterOptions()
	}
	return &MarkdownConverter{
		options: opts,
	}
}

// NewHTMLConverter 创建HTML转换器
func NewHTMLConverter(opts *ConverterOptions) *HTMLConverter {
	if opts == nil {
		opts = DefaultConverterOptions()
	}
	return &HTMLConverter{
		options: opts,
	}
}

// NewJSONConverter 创建JSON转换器
func NewJSONConverter(opts *ConverterOptions) *JSONConverter {
	if opts == nil {
		opts = DefaultConverterOptions()
	}
	return &JSONConverter{
		options: opts,
	}
}

// OutputFormat 输出格式
type OutputFormat string

const (
	FormatMarkdown OutputFormat = "markdown"
	FormatHTML     OutputFormat = "html"
	FormatJSON     OutputFormat = "json"
)

// Writer 输出写入器
type Writer interface {
	Write(data []byte) error
	WriteString(s string) error
	Flush() error
	Close() error
}

// FileWriter 文件写入器
type FileWriter struct {
	path     string
	file     *os.File
	buffer   *bytes.Buffer
	overwrite bool
}

// NewFileWriter 创建文件写入器
func NewFileWriter(path string, overwrite bool) *FileWriter {
	return &FileWriter{
		path:      path,
		buffer:    bytes.NewBuffer(nil),
		overwrite: overwrite,
	}
}

// StdoutWriter 标准输出写入器
type StdoutWriter struct {
	writer io.Writer
}

// NewStdoutWriter 创建标准输出写入器
func NewStdoutWriter(writer io.Writer) *StdoutWriter {
	return &StdoutWriter{
		writer: writer,
	}
}

// ConversionError 转换错误
type ConversionError struct {
	Message   string
	Format    OutputFormat
	SourceErr error
}

func (e *ConversionError) Error() string {
	if e.SourceErr != nil {
		return e.Message + ": " + e.SourceErr.Error()
	}
	return e.Message
}

// Unwrap 返回原始错误
func (e *ConversionError) Unwrap() error {
	return e.SourceErr
}

// Convert 将Markdown文档转换为字节数组
func (mc *MarkdownConverter) Convert(doc *parser.MarkdownDocument) ([]byte, error) {
	if doc == nil {
		return nil, &ConversionError{
			Message: "document is nil",
			Format:  FormatMarkdown,
		}
	}

	var result string
	if mc.options.EnableTableOfContents && doc.Title != "" {
		result += "# " + doc.Title + "\n\n"
	}
	result += doc.Content
	return []byte(result), nil
}

// Convert 将Markdown文档转换为HTML
func (hc *HTMLConverter) Convert(doc *parser.MarkdownDocument) ([]byte, error) {
	if doc == nil {
		return nil, &ConversionError{
			Message: "document is nil",
			Format:  FormatHTML,
		}
	}

	// 简单的HTML转换，完整的实现将在Phase 2中完成
	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	if doc.Title != "" {
		html += "<title>" + doc.Title + "</title>\n"
	}
	html += "</head>\n<body>\n"
	if doc.Title != "" {
		html += "<h1>" + doc.Title + "</h1>\n"
	}
	html += "<div>" + doc.Content + "</div>\n"
	html += "</body>\n</html>"

	return []byte(html), nil
}

// Convert 将Markdown文档转换为JSON
func (jc *JSONConverter) Convert(doc *parser.MarkdownDocument) ([]byte, error) {
	if doc == nil {
		return nil, &ConversionError{
			Message: "document is nil",
			Format:  FormatJSON,
		}
	}

	// 使用标准库进行安全的JSON序列化
	jsonDoc := map[string]interface{}{
		"title":   doc.Title,
		"content": doc.Content,
	}

	if len(doc.Metadata) > 0 {
		jsonDoc["metadata"] = doc.Metadata
	}

	result, err := json.MarshalIndent(jsonDoc, "", "  ")
	if err != nil {
		return nil, &ConversionError{
			Message:   "failed to marshal JSON",
			Format:    FormatJSON,
			SourceErr: err,
		}
	}

	return result, nil
}