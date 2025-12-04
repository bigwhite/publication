# issue2md

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Build Status](https://img.shields.io/badge/Build-Passing-green.svg)](https://github.com/bigwhite/my-issue2md)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

一个高效的 GitHub Issue 到 Markdown 转换工具，支持 CLI 和 Web 服务两种使用模式。

## 🌟 核心特性 (Features)

- **多格式支持**: 支持 Markdown、HTML 和 JSON 格式输出
- **双模式运行**: 提供命令行工具 (CLI) 和 Web 服务两种使用方式
- **灵活配置**: 丰富的配置选项，满足不同使用场景
- **高安全性**: 非root用户运行，容器化部署
- **生产就绪**: 完整的测试覆盖、CI/CD支持和优雅关闭机制
- **开发者友好**: 遵循 Go 语言最佳实践，代码结构清晰

## 📦 安装指南 (Installation)

### 方式一：从源码构建

```bash
# 克隆仓库
git clone https://github.com/bigwhite/my-issue2md.git
cd my-issue2md

# 构建应用
make build

# 安装到 GOPATH/bin
make install
```

### 方式二：使用 Docker

```bash
# 构建 Docker 镜像
make docker-build

# 运行容器
make docker-run
```

### 方式三：下载预构建二进制文件

从 [Releases](https://github.com/bigwhite/my-issue2md/releases) 页面下载适合您系统的二进制文件。

## 🚀 使用方法 (Usage)

### CLI 工具使用

#### 基本语法

```bash
issue2md [owner/repo] [issue-number] [flags]
```

#### 命令行参数

| 参数 | 简写 | 描述 | 默认值 |
|------|------|------|--------|
| `--help` | `-h` | 显示帮助信息 | - |
| `--version` | `-v` | 显示版本信息 | - |
| `--output` | `-o` | 输出文件路径 | `"output.md"` |
| `--format` | `-f` | 输出格式：markdown, html, json | `"markdown"` |
| `--token` | `-t` | GitHub token (或设置 GITHUB_TOKEN 环境变量) | - |
| `--no-comments` | - | 排除评论内容 | `false` |
| `--no-metadata` | - | 排除元数据信息 | `false` |
| `--no-timestamps` | - | 排除时间戳信息 | `false` |
| `--overwrite` | - | 覆盖已存在的输出文件 | `false` |
| `--debug` | - | 启用调试日志 | `false` |

#### 使用示例

```bash
# 基本用法 - 转换 React 项目的第 12345 号 issue
issue2md facebook/react 12345

# 指定输出文件
issue2md facebook/react 12345 --output=issue.md

# 输出为 HTML 格式，不包含评论
issue2md facebook/react 12345 --format=html --no-comments

# 使用环境变量中的 GitHub token
export GITHUB_TOKEN=your_token_here
issue2md facebook/react 12345

# 调试模式
issue2md facebook/react 12345 --debug
```

### Web 服务使用

#### 启动 Web 服务

```bash
# 使用默认端口 8080 启动
./bin/issue2mdweb

# 或使用 Docker
docker run -p 8080:8080 -e GITHUB_TOKEN=your_token issue2md:latest
```

#### API 端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/` | GET | 服务首页，显示基本信息 |
| `/health` | GET | 健康检查端点 |
| `/api/v1/convert` | POST | Issue 转换 API (开发中) |

#### 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `PORT` | Web 服务端口 | `8080` |
| `GITHUB_TOKEN` | GitHub API 访问令牌 | - |
| `DEBUG` | 启用调试模式 | `false` |
| `NO_COLOR` | 禁用彩色输出 | `false` |

## 🔨 构建方法 (Building from Source)

### 前置要求

- Go 1.21 或更高版本
- Make 工具
- Git

### 开发环境设置

```bash
# 设置开发环境
make dev-setup

# 运行测试
make test

# 代码格式化
make format

# 静态分析
make lint
```

### 构建命令

```bash
# 构建所有应用 (CLI + Web)
make build

# 仅构建 CLI 工具
CGO_ENABLED=0 GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build -o bin/issue2md-cli ./cmd/issue2md

# 仅构建 Web 服务
CGO_ENABLED=0 GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build -o bin/issue2md-web ./cmd/issue2mdweb

# 构建 Docker 镜像
make docker-build

# 指定镜像标签
make docker-build DOCKER_TAG=v1.0.0
```

### 测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行基准测试
make test-benchmark
```

### 代码质量检查

```bash
# 格式化代码
make format

# 运行静态分析
make lint

# 运行完整验证流程 (格式化 + 静态分析 + 测试)
make verify
```

## 📁 项目结构

```
.
├── cmd/                    # 应用程序入口点
│   ├── issue2md/          # CLI 工具
│   └── issue2mdweb/       # Web 服务
├── internal/              # 内部包
│   ├── cli/              # CLI 框架
│   ├── config/           # 配置管理
│   ├── converter/        # 格式转换器
│   ├── github/           # GitHub API 客户端
│   └── parser/           # Markdown 解析器
├── specs/                 # 功能规格说明
├── .claude/              # Claude 配置
├── Makefile              # 构建脚本
├── Dockerfile            # Docker 镜像定义
├── go.mod                # Go 模块定义
└── constitution.md        # 项目开发宪法
```

## 🔧 配置选项

### 配置文件结构

应用支持通过配置文件和环境变量进行配置：

```json
{
  "github_token": "your_github_token",
  "output": {
    "format": "markdown",
    "filename": "output.md",
    "destination": "output",
    "overwrite": false
  },
  "parser": {
    "include_comments": true,
    "include_metadata": true,
    "include_timestamps": true,
    "include_user_links": true,
    "emojis_enabled": true,
    "preserve_line_breaks": true
  }
}
```

### 环境变量配置

| 环境变量 | 对应配置项 | 描述 |
|----------|------------|------|
| `GITHUB_TOKEN` | `github_token` | GitHub API 访问令牌 |
| `DEBUG` | 影响解析器配置 | 启用调试模式 |
| `NO_COLOR` | - | 禁用彩色输出 |

## 🐳 Docker 部署

### 构建镜像

```bash
# 默认标签 (latest)
make docker-build

# 指定标签
make docker-build DOCKER_TAG=v1.0.0
```

### 运行容器

```bash
# 基本运行
docker run -p 8080:8080 -e GITHUB_TOKEN=your_token issue2md:latest

# 挂载卷用于输出文件
docker run -p 8080:8080 -v $(pwd)/output:/app/output -e GITHUB_TOKEN=your_token issue2md:latest

# 作为 CLI 工具使用
docker run --rm -v $(pwd):/app -w /app issue2md:latest ./bin/issue2md facebook/react 12345
```

## 🧪 开发指南

### 代码规范

本项目严格遵循以下开发原则：

1. **简单性原则**: 遵循 Go 语言"少即是多"的哲学
2. **测试先行**: 所有新功能必须从失败的测试开始 (TDD)
3. **明确性原则**: 代码首要目标是让人类易于理解
4. **错误处理**: 所有错误都必须被显式处理

详细的开发规范请参考 [constitution.md](constitution.md)。

### 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

确保：
- 遵循 Conventional Commits 规范
- 所有测试通过 (`make test`)
- 代码通过静态分析 (`make lint`)
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证。详情请参考 [LICENSE](LICENSE) 文件。

## 🤝 支持

如果您遇到问题或有建议，请：

1. 查看文档
2. 搜索现有的 [Issues](https://github.com/bigwhite/my-issue2md/issues)
3. 创建新的 Issue 描述您的问题

## 📊 版本历史

- **v1.0.0** - 初始版本
  - 基础 CLI 功能
  - Web 服务框架
  - Docker 支持

---

**issue2md** - 让 GitHub Issue 转换变得简单高效！ 🚀