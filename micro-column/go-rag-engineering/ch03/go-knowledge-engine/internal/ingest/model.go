package ingest

import "time"

// ParsedDocument 表示解析器输出的文档。
type ParsedDocument struct {
	Title    string
	Headings []string
	Content  string
}

// KnowledgeItem 表示知识系统中的统一知识条目。
type KnowledgeItem struct {
	ID          string
	SourcePath  string
	Format      string
	Title       string
	Headings    []string
	Content     string
	Fingerprint string
	UpdatedAt   time.Time
}

// Parser 定义不同格式解析器的统一接口。
type Parser interface {
	Parse(path string, raw []byte) (ParsedDocument, error)
}

// Config 定义知识摄取流水线的配置。
type Config struct {
	SourceRoot    string
	KnowledgeRoot string
	ManifestPath  string
}

// Stats 记录一次摄取任务的统计信息。
type Stats struct {
	Scanned int
	Added   int
	Updated int
	Skipped int
	Deleted int
	Failed  int
}

// manifestEntry 记录单个来源文件上一次的摄取结果。
type manifestEntry struct {
	ID          string `json:"id"`
	Fingerprint string `json:"fingerprint"`
	Output      string `json:"output"`
}

// manifest 保存来源路径与摄取结果之间的映射。
type manifest map[string]manifestEntry
