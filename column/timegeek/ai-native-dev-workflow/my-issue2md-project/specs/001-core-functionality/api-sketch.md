# API 设计草稿

本文档定义了 `internal/converter` 和 `internal/github` 包的主要接口设计，作为后续实现的参考。

---

## internal/github 包

### 1. 核心接口定义

```go
// GitHubClient 定义 GitHub API 客户端接口
type GitHubClient interface {
    // GetIssue 获取 GitHub Issue 完整信息
    GetIssue(ctx context.Context, owner, repo string, number int) (*Issue, error)

    // GetPullRequest 获取 GitHub PR 完整信息
    GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error)

    // GetDiscussion 获取 GitHub Discussion 完整信息
    GetDiscussion(ctx context.Context, owner, repo string, number int) (*Discussion, error)
}
```

### 2. 数据结构

```go
// Resource 表示 GitHub 资源的基础信息
type Resource struct {
    Title       string    `json:"title"`
    URL         string    `json:"url"`
    Author      User      `json:"author"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Status      string    `json:"status"` // "open", "closed", "merged"
    Body        string    `json:"body"`
    Reactions   Reactions `json:"reactions"`
    Comments    []Comment `json:"comments"`
}

// Issue 表示 GitHub Issue
type Issue struct {
    Resource
    Number int               `json:"number"`
    Labels []Label          `json:"labels"`
    Milestone *Milestone     `json:"milestone,omitempty"`
}

// PullRequest 表示 GitHub Pull Request
type PullRequest struct {
    Resource
    Number         int    `json:"number"`
    BaseBranch     string `json:"base_branch"`
    HeadBranch     string `json:"head_branch"`
    MergeCommitSHA string `json:"merge_commit_sha,omitempty"`
}

// Discussion 表示 GitHub Discussion
type Discussion struct {
    Resource
    Number      int         `json:"number"`
    Category    DiscussionCategory `json:"category"`
    Answer      *Comment    `json:"answer,omitempty"`
}

// User 表示 GitHub 用户信息
type User struct {
    Login     string `json:"login"`
    AvatarURL string `json:"avatar_url"`
    HTMLURL   string `json:"html_url"`
}

// Comment 表示评论
type Comment struct {
    ID        int64     `json:"id"`
    Author    User      `json:"author"`
    Body      string    `json:"body"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Reactions Reactions `json:"reactions"`
    IsAnswer  bool      `json:"is_answer"` // 仅用于 Discussion
}

// Reactions 表示 reactions 统计
type Reactions struct {
    ThumbsUp   int `json:"thumbs_up"`
    ThumbsDown int `json:"thumbs_down"`
    Laugh      int `json:"laugh"`
    Hooray     int `json:"hooray"`
    Confused   int `json:"confused"`
    Heart      int `json:"heart"`
    Rocket     int `json:"rocket"`
    Eyes       int `json:"eyes"`
}

// Label 表示 GitHub 标签
type Label struct {
    Name  string `json:"name"`
    Color string `json:"color"`
}

// Milestone 表示 GitHub Milestone
type Milestone struct {
    Title string `json:"title"`
    State string `json:"state"`
}

// DiscussionCategory 表示 Discussion 分类
type DiscussionCategory struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    Emoji string `json:"emoji"`
}
```

### 3. 客户端实现

```go
// NewGitHubClient 创建新的 GitHub API 客户端
func NewGitHubClient(token string) GitHubClient {
    // 返回实现了 GitHubClient 接口的实例
}

// NewGitHubClientWithHTTPClient 使用自定义 HTTP 客户端创建
func NewGitHubClientWithHTTPClient(token string, httpClient *http.Client) GitHubClient {
    // 用于测试时注入 mock HTTP 客户端
}
```

---

## internal/parser 包

### 1. URL 解析接口

```go
// ResourceURL 表示解析后的 GitHub 资源 URL
type ResourceURL struct {
    Type   string // "issue", "pull", "discussion"
    Owner  string
    Repo   string
    Number int
    URL    string // 原始 URL
}

// Parser 定义 URL 解析器接口
type Parser interface {
    // Parse 解析 GitHub URL，返回资源信息
    Parse(rawURL string) (*ResourceURL, error)

    // Validate 验证 URL 格式是否支持
    Validate(rawURL string) error
}
```

### 2. 实现函数

```go
// NewParser 创建新的解析器实例
func NewParser() Parser

// ParseGitHubURL 解析 GitHub URL (便捷函数)
func ParseGitHubURL(rawURL string) (*ResourceURL, error)
```

---

## internal/converter 包

### 1. Markdown 转换器接口

```go
// Converter 定义 Markdown 转换器接口
type Converter interface {
    // Convert 将 GitHub 资源转换为 Markdown
    Convert(ctx context.Context, resource Resource, options *ConvertOptions) ([]byte, error)
}

// ConvertOptions 转换选项
type ConvertOptions struct {
    EnableReactions  bool // 是否包含 reactions 统计
    EnableUserLinks  bool // 是否将用户名转换为链接
    OutputToStdout   bool // 是否输出到标准输出
    OutputFile       string // 输出文件路径
}
```

### 2. 资源到 Markdown 的转换函数

```go
// ConvertIssue 转换 Issue 为 Markdown
func ConvertIssue(ctx context.Context, issue *Issue, options *ConvertOptions) ([]byte, error)

// ConvertPullRequest 转换 PR 为 Markdown
func ConvertPullRequest(ctx context.Context, pr *PullRequest, options *ConvertOptions) ([]byte, error)

// ConvertDiscussion 转换 Discussion 为 Markdown
func ConvertDiscussion(ctx context.Context, discussion *Discussion, options *ConvertOptions) ([]byte, error)
```

### 3. 辅助函数

```go
// GenerateFrontmatter 生成 YAML frontmatter
func GenerateFrontmatter(resource Resource) (map[string]interface{}, error)

// FormatComment 格式化单个评论
func FormatComment(comment Comment, enableUserLinks bool) string

// FormatReactions 格式化 reactions 统计
func FormatReactions(reactions Reactions) string

// FormatUserLink 格式化用户链接
func FormatUserLink(user User) string
```

### 4. 模板相关

```go
// MarkdownTemplate 表示 Markdown 模板
type MarkdownTemplate struct {
    // 可以包含可配置的模板
}

// NewMarkdownTemplate 创建新的模板实例
func NewMarkdownTemplate() *MarkdownTemplate
```

---

## internal/cli 包

### 1. CLI 应用接口

```go
// CLIApp 表示命令行应用
type CLIApp interface {
    // Run 运行 CLI 应用
    Run(ctx context.Context, args []string) error
}
```

### 2. 命令行参数定义

```go
// CLIArgs 表示解析后的命令行参数
type CLIArgs struct {
    URL               string
    OutputFile        string
    EnableReactions   bool
    EnableUserLinks   bool
    ShowHelp          bool
    ShowVersion       bool
}
```

---

## internal/config 包

### 1. 配置接口

```go
// Config 表示应用配置
type Config interface {
    // GitHubToken 获取 GitHub Token
    GitHubToken() string

    // APITimeout 获取 API 超时时间
    APITimeout() time.Duration

    // UserAgent 获取 User-Agent
    UserAgent() string
}
```

### 2. 配置加载

```go
// NewConfig 创建新的配置实例
func NewConfig() Config

// NewConfigWithToken 使用指定 Token 创建配置
func NewConfigWithToken(token string) Config

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() Config
```

---

## 包间依赖关系

```
cmd/issue2md/
├── internal/cli/          # 命令行参数解析和执行
├── internal/config/       # 配置管理
└── internal/converter/    # Markdown 转换
    └── internal/github/   # GitHub API 调用
    └── internal/parser/   # URL 解析
```

**依赖流向:**
- `cmd/` → `cli/` → `config/`, `parser/`, `converter/`
- `converter/` → `github/`
- `cli/` → `parser/`, `converter/`
- `converter/` 不依赖其他内部包，只依赖 `github/`

---

## 下一步实现顺序

1. **internal/parser** - URL 解析功能 (最容易测试)
2. **internal/github** - GitHub API 客户端 (核心数据获取)
3. **internal/converter** - Markdown 转换功能
4. **internal/config** - 配置管理
5. **internal/cli** - 命令行接口
6. **cmd/issue2md** - 主程序入口

这个顺序遵循了自底向上、先易后难的原则，便于编写测试和逐步集成。