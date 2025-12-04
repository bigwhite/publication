# issue2md 技术实现方案

**版本**: 1.0
**状态**: 准备实施
**基于**: specs/001-core-functionality/spec.md v1.0

---

## 1. 技术上下文总结

### 1.1 核心技术栈

- **语言**: Go (版本 >= 1.21.0，当前项目配置为 1.25.0)
- **Web框架**: 仅使用标准库 `net/http`，严格遵循"简单性原则"
- **GitHub API客户端**: 使用 `google/go-github` 库，结合 GraphQL API v4
- **Markdown处理**: 不使用第三方库，直接使用 `fmt` 和字符串处理
- **构建工具**: 使用 `Makefile` 进行标准化操作
- **测试方法**: 表格驱动测试（Table-Driven Tests）优先

### 1.2 项目外部依赖

```go
// go.mod 将添加以下依赖
require (
    github.com/google/go-github v45.1.0
    golang.org/x/oauth2 v0.12.0  // GitHub认证
)
```

### 1.3 架构设计原则

- **简单性**: 避免过度抽象，使用Go标准库
- **内聚性**: 每个内部包职责单一明确
- **可测试性**: 所有组件都支持单元测试和集成测试
- **错误处理**: 严格遵循错误包装原则

---

## 2. "合宪性"审查

### 2.1 符合性检查表

| 宪法条款 | 符合性 | 具体体现 |
|---------|--------|----------|
| **第一条：简单性原则** | ✅ 完全符合 | 仅使用标准库`net/http`，最小化外部依赖，避免过度抽象 |
| **第二条：测试先行铁律** | ✅ 完全符合 | 采用TDD方法，优先表格驱动测试，避免mocks |
| **第三条：明确性原则** | ✅ 完全符合 | 显式错误处理，无全局变量，依赖注入 |

### 2.2 具体合规措施

#### 2.2.1 简单性原则实施
- **YAGNI**: 仅实现spec.md要求的功能，不预置扩展性
- **标准库优先**: Web服务使用`net/http`，JSON处理使用`encoding/json`
- **反过度工程**: 使用具体函数而非复杂接口体系

#### 2.2.2 测试先行实施
- **TDD循环**: 每个功能从失败测试开始
- **表格驱动**: 所有单元测试使用`[]struct`模式
- **集成测试优先**: 使用真实GitHub API进行测试

#### 2.2.3 明确性原则实施
- **错误处理**: 严格使用`fmt.Errorf("context: %w", err)`
- **依赖注入**: 通过构造函数传递所有依赖
- **配置管理**: 通过环境变量显式配置，无隐藏全局状态

---

## 3. 项目结构细化

### 3.1 完整目录结构

```
issue2md/
├── cmd/
│   └── issue2md/
│       └── main.go                 # 应用入口点
├── internal/
│   ├── cli/                        # 命令行接口处理
│   │   ├── app.go                  # CLI应用主逻辑
│   │   ├── args.go                 # 命令行参数解析
│   │   └── version.go              # 版本信息管理
│   ├── config/                     # 配置管理
│   │   ├── config.go               # 配置接口和实现
│   │   └── env.go                  # 环境变量处理
│   ├── converter/                  # Markdown转换核心
│   │   ├── converter.go            # 转换器接口和实现
│   │   ├── frontmatter.go          # YAML frontmatter生成
│   │   ├── formatter.go            # 内容格式化
│   │   └── templates.go            # Markdown模板定义
│   ├── github/                     # GitHub API交互
│   │   ├── client.go               # GitHub客户端实现
│   │   ├── types.go                # GitHub数据结构定义
│   │   ├── queries.go              # GraphQL查询构建
│   │   └── auth.go                 # 认证处理
│   └── parser/                     # URL解析与识别
│       ├── parser.go               # URL解析器实现
│       ├── types.go                # 解析相关类型定义
│       └── validation.go           # URL验证逻辑
├── pkg/                            # (可选) 公共包
├── go.mod
├── go.sum
├── Makefile                        # 构建脚本
├── README.md
├── LICENSE
└── specs/
    └── 001-core-functionality/
        ├── spec.md                 # 功能规格
        └── plan.md                 # 本技术方案
```

### 3.2 包职责与依赖关系

#### 3.2.1 包职责矩阵

| 包名 | 主要职责 | 输入 | 输出 |
|------|----------|------|------|
| **parser** | GitHub URL解析和类型识别 | Raw URL字符串 | ResourceURL结构 |
| **github** | GitHub API数据获取 | ResourceURL, 认证Token | Issue/PR/Discussion结构 |
| **converter** | 数据到Markdown转换 | GitHub资源数据, 转换选项 | Markdown字节流 |
| **config** | 应用配置管理 | 环境变量, 启动参数 | Config接口实现 |
| **cli** | 命令行接口和应用流程 | 命令行参数, stdin | CLI执行结果 |

#### 3.2.2 依赖关系图

```
cmd/issue2md/main.go
    ↓
internal/cli/
    ↓ ↓ ↓
internal/config/  internal/parser/  internal/converter/
                                    ↓
                              internal/github/
```

**依赖规则**:
- `cmd/` 仅依赖 `internal/cli`
- `cli/` 可依赖 `config/`, `parser/`, `converter/`
- `converter/` 仅依赖 `github/`
- `parser/` 和 `config/` 独立，不依赖其他内部包
- 所有包都可以依赖Go标准库

### 3.3 模块间通信接口

```go
// 统一错误类型
var (
    ErrInvalidURL     = errors.New("invalid GitHub URL")
    ErrUnsupportedURL = errors.New("unsupported URL type")
    ErrAPIError       = errors.New("GitHub API error")
    ErrAuthRequired   = errors.New("authentication required")
)

// 核心数据传递格式
type ResourceData struct {
    Type string // "issue", "pr", "discussion"
    Data interface{} // *Issue, *PullRequest, 或 *Discussion
}
```

---

## 4. 核心数据结构

### 4.1 统一资源表示

```go
// Resource 表示所有GitHub资源的通用字段
type Resource struct {
    Title       string    `json:"title"`
    URL         string    `json:"url"`
    Number      int       `json:"number"`
    Author      User      `json:"author"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Status      string    `json:"status"`      // "open", "closed", "merged"
    Body        string    `json:"body"`
    Reactions   Reactions `json:"reactions"`
    Comments    []Comment `json:"comments"`
}

// User 表示GitHub用户信息
type User struct {
    Login     string `json:"login"`
    AvatarURL string `json:"avatar_url"`
    HTMLURL   string `json:"html_url"`
}

// Comment 表示评论内容
type Comment struct {
    ID        int64     `json:"id"`
    Author    User      `json:"author"`
    Body      string    `json:"body"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Reactions Reactions `json:"reactions"`
    IsAnswer  bool      `json:"is_answer"` // 仅用于Discussion
}

// Reactions 表示reaction统计
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
```

### 4.2 特定资源类型

```go
// Issue 表示GitHub Issue
type Issue struct {
    Resource
    Labels    []Label    `json:"labels"`
    Milestone *Milestone `json:"milestone,omitempty"`
}

// PullRequest 表示GitHub Pull Request
type PullRequest struct {
    Resource
    BaseBranch     string `json:"base_branch"`
    HeadBranch     string `json:"head_branch"`
    MergeCommitSHA string `json:"merge_commit_sha,omitempty"`
}

// Discussion 表示GitHub Discussion
type Discussion struct {
    Resource
    Category    DiscussionCategory `json:"category"`
    Answer      *Comment           `json:"answer,omitempty"`
}

// Label 表示GitHub标签
type Label struct {
    Name  string `json:"name"`
    Color string `json:"color"`
}

// Milestone 表示GitHub里程碑
type Milestone struct {
    Title string `json:"title"`
    State string `json:"state"`
}

// DiscussionCategory 表示Discussion分类
type DiscussionCategory struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    Emoji string `json:"emoji"`
}
```

### 4.3 配置和选项结构

```go
// Config 表示应用配置
type Config struct {
    GitHubToken string
    UserAgent   string
    APITimeout  time.Duration
}

// ConvertOptions 表示Markdown转换选项
type ConvertOptions struct {
    EnableReactions bool // 是否包含reactions统计
    EnableUserLinks bool // 是否将@username转换为链接
}

// CLIArgs 表示命令行参数
type CLIArgs struct {
    URL             string
    OutputFile      string
    EnableReactions bool
    EnableUserLinks bool
    ShowHelp        bool
    ShowVersion     bool
}
```

### 4.4 URL解析结构

```go
// ResourceURL 表示解析后的GitHub资源URL
type ResourceURL struct {
    Type   string // "issue", "pull", "discussion"
    Owner  string
    Repo   string
    Number int
    URL    string // 原始URL
}

// URLPattern 支持的URL模式定义
type URLPattern struct {
    Pattern string
    Type    string
}

var SupportedPatterns = []URLPattern{
    {`^https://github\.com/([^/]+)/([^/]+)/issues/(\d+)$`, "issue"},
    {`^https://github\.com/([^/]+)/([^/]+)/pull/(\d+)$`, "pull"},
    {`^https://github\.com/([^/]+)/([^/]+)/discussions/(\d+)$`, "discussion"},
}
```

---

## 5. 接口设计

### 5.1 GitHub API客户端接口

```go
// GitHubClient 定义GitHub API交互接口
type GitHubClient interface {
    // GetIssue 获取Issue完整信息
    GetIssue(ctx context.Context, owner, repo string, number int) (*Issue, error)

    // GetPullRequest 获取PR完整信息
    GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error)

    // GetDiscussion 获取Discussion完整信息
    GetDiscussion(ctx context.Context, owner, repo string, number int) (*Discussion, error)
}

// GitHubClientBuilder 客户端构建器
type GitHubClientBuilder interface {
    WithToken(token string) GitHubClientBuilder
    WithHTTPClient(client *http.Client) GitHubClientBuilder
    WithTimeout(timeout time.Duration) GitHubClientBuilder
    Build() (GitHubClient, error)
}
```

### 5.2 URL解析器接口

```go
// Parser 定义URL解析接口
type Parser interface {
    // Parse 解析GitHub URL
    Parse(rawURL string) (*ResourceURL, error)

    // Validate 验证URL格式
    Validate(rawURL string) error

    // SupportedTypes 返回支持的资源类型
    SupportedTypes() []string
}
```

### 5.3 Markdown转换器接口

```go
// Converter 定义Markdown转换接口
type Converter interface {
    // Convert 将GitHub资源转换为Markdown
    Convert(ctx context.Context, resource interface{}, options *ConvertOptions) ([]byte, error)

    // ConvertIssue 转换Issue
    ConvertIssue(ctx context.Context, issue *Issue, options *ConvertOptions) ([]byte, error)

    // ConvertPullRequest 转换PR
    ConvertPullRequest(ctx context.Context, pr *PullRequest, options *ConvertOptions) ([]byte, error)

    // ConvertDiscussion 转换Discussion
    ConvertDiscussion(ctx context.Context, discussion *Discussion, options *ConvertOptions) ([]byte, error)
}

// FrontmatterGenerator YAML frontmatter生成器接口
type FrontmatterGenerator interface {
    Generate(resource interface{}) (map[string]interface{}, error)
    ToYAML(data map[string]interface{}) ([]byte, error)
}
```

### 5.4 配置管理接口

```go
// Config 定义配置接口
type Config interface {
    GitHubToken() string
    UserAgent() string
    APITimeout() time.Duration
    HTTPClient() *http.Client
}

// ConfigLoader 配置加载器接口
type ConfigLoader interface {
    LoadFromEnv() Config
    LoadWithToken(token string) Config
    Validate() error
}
```

### 5.5 CLI应用接口

```go
// CLIApp 定义CLI应用接口
type CLIApp interface {
    Run(ctx context.Context, args []string) error
    Execute(ctx context.Context, parsedArgs *CLIArgs) error
}

// ArgParser 命令行参数解析器接口
type ArgParser interface {
    Parse(args []string) (*CLIArgs, error)
    Validate(args *CLIArgs) error
    ShowUsage() error
    ShowVersion() error
}
```

### 5.6 文件输出接口

```go
// OutputHandler 输出处理器接口
type OutputHandler interface {
    WriteToFile(data []byte, filename string) error
    WriteToStdout(data []byte) error
    EnsureDirectory(filename string) error
}
```

---

## 6. 实施计划与优先级

### 6.1 开发阶段划分

#### 阶段1：基础设施 (第1-2周)
1. **parser包** - URL解析和验证 (无外部依赖，易测试)
2. **config包** - 配置管理 (环境变量处理)
3. **核心数据结构** - 定义所有struct和接口

#### 阶段2：数据获取 (第3-4周)
1. **github包** - GitHub API客户端实现
2. **认证处理** - Token管理和请求认证
3. **错误处理** - 统一错误类型和处理

#### 阶段3：数据处理 (第5-6周)
1. **converter包** - Markdown转换核心逻辑
2. **模板实现** - YAML frontmatter和内容格式化
3. **选项处理** - reactions和用户链接功能

#### 阶段4：用户接口 (第7周)
1. **cli包** - 命令行参数解析和应用流程
2. **输出处理** - 文件和标准输出
3. **帮助和版本信息**

#### 阶段5：集成与优化 (第8周)
1. **cmd/issue2md** - 主程序入口点
2. **Makefile** - 构建和测试脚本
3. **集成测试** - 端到端测试
4. **性能优化** - API限流和缓存

### 6.2 TDD实施策略

#### 6.2.1 测试优先级
1. **单元测试** - 每个包独立测试 (90%+覆盖率)
2. **集成测试** - 真实GitHub API调用测试
3. **端到端测试** - 完整CLI命令测试

#### 6.2.2 测试数据管理
- 使用GitHub的公开测试仓库
- 创建标准测试数据集
- 支持Mock API用于CI/CD

### 6.3 质量保证措施

#### 6.3.1 代码质量
- 使用`golangci-lint`进行静态检查
- 严格遵循Go官方代码规范
- 100%类型安全的接口设计

#### 6.3.2 性能要求
- API请求响应时间 < 5秒
- 内存使用峰值 < 50MB
- 支持GitHub API限流自动重试

---

## 7. 关键技术决策与理由

### 7.1 使用GraphQL而非REST API

**决策**: 使用GitHub GraphQL API v4获取数据

**理由**:
- 单次请求获取所有需要的数据，减少API调用次数
- 更好的性能和更低的延迟
- 精确获取所需字段，减少数据传输量
- 更好的类型安全性

### 7.2 标准库优先策略

**决策**: 严格使用Go标准库，避免引入额外Web框架

**理由**:
- 符合项目宪法"简单性原则"
- 减少外部依赖和潜在的安全风险
- 提高代码可维护性和理解性
- 降低学习成本

### 7.3 表格驱动测试优先

**决策**: 所有单元测试采用表格驱动测试模式

**理由**:
- 提高测试用例的可读性和维护性
- 便于添加新的测试场景
- 减少重复代码，提高测试质量
- 符合Go社区最佳实践

---

## 8. 风险评估与缓解策略

### 8.1 技术风险

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| GitHub API限流 | 中 | 中 | 实现指数退避重试，支持Token认证 |
| GraphQL查询复杂度 | 低 | 中 | 从简单查询开始，逐步优化 |
| 私有仓库访问 | 中 | 低 | 明确文档说明，提供Token配置指导 |

### 8.2 项目风险

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| 功能范围蔓延 | 中 | 高 | 严格遵循spec.md，避免镀金 |
| 测试覆盖不足 | 低 | 中 | TDD开发，自动化覆盖率检查 |
| 性能不达标 | 低 | 中 | 早期性能测试，渐进优化 |

---

## 9. 成功标准

### 9.1 功能完整性
- ✅ 支持Issue、PR、Discussion三种类型URL解析
- ✅ 完整实现spec.md中定义的所有功能
- ✅ 命令行接口完全符合设计规范

### 9.2 质量标准
- ✅ 单元测试覆盖率 >= 90%
- ✅ 所有静态检查通过
- ✅ 零已知安全漏洞

### 9.3 性能标准
- ✅ 单个资源处理时间 < 5秒
- ✅ 内存使用峰值 < 50MB
- ✅ 成功处理所有验收测试用例

### 9.4 可维护性
- ✅ 代码结构清晰，包职责明确
- ✅ 完整的API文档和使用示例
- ✅ 符合Go语言最佳实践

---

**下一步**: 开始实施阶段1，从`internal/parser`包开始，遵循TDD方法逐步实现各个模块。