# issue2md Core Functionality Specification

**Version**: 1.0
**Status**: Ready for Implementation
**Target**: MVP CLI Tool

---

## 1. 用户故事 (User Stories)

### 1.1 CLI 核心用户故事 (MVP)

**作为一名开发者，我希望能够通过一个简单的命令将 GitHub Issue/PR/Discussion 转换为 Markdown 文件，以便于本地文档整理和知识归档。**

**验收标准:**
- 能够自动识别 GitHub URL 类型并正确处理
- 生成的 Markdown 文件包含完整的讨论内容
- 支持命令行参数控制输出格式和内容
- 错误处理清晰，用户能够快速定位问题

### 1.2 Web 界面用户故事 (Future)

**作为一名团队协作者，我希望通过 Web 界面上传 GitHub URL 链接列表，批量转换为 Markdown 并打包下载，以便于团队知识库建设。**

*注：此功能作为未来迭代，当前版本不实现*

---

## 2. 功能性需求 (Functional Requirements)

### 2.1 URL 解析与识别

**需求**: 工具必须能够自动识别并处理以下类型的 GitHub URL:

| 类型 | URL 模式 | 示例 |
|------|----------|------|
| Issue | `https://github.com/{owner}/{repo}/issues/{number}` | `https://github.com/bigwhite/issue2md/issues/1` |
| Pull Request | `https://github.com/{owner}/{repo}/pull/{number}` | `https://github.com/bigwhite/issue2md/pull/42` |
| Discussion | `https://github.com/{owner}/{repo}/discussions/{number}` | `https://github.com/bigwhite/issue2md/discussions/123` |

**技术要求:**
- URL 验证：必须检查 URL 格式的有效性
- 类型识别：解析 URL 路径确定资源类型
- 错误处理：无效 URL 或不支持的类型应返回明确的错误信息

### 2.2 命令行接口设计

**基本语法:**
```bash
issue2md [flags] <url> [output_file]
```

**参数说明:**
- `<url>`: GitHub Issue/PR/Discussion URL (必需)
- `[output_file]`: 输出文件路径 (可选，默认输出到 stdout)

**Flags:**
```bash
-enable-reactions    # 包含 reactions 统计信息 (默认: false)
-enable-user-links   # 将用户名渲染为 GitHub 主页链接 (默认: false)
-h, -help           # 显示帮助信息
-v, -version        # 显示版本信息
```

**环境变量:**
```bash
GITHUB_TOKEN        # GitHub Personal Access Token (可选，用于提高 API 限制)
```

### 2.3 数据获取与转换

#### 2.3.1 核心内容获取

**必须包含的信息:**
- 标题 (Title)
- 作者信息 (Author: username, avatar URL)
- 创建时间 (Created At)
- 最后更新时间 (Last Updated At)
- 状态 (Status: Open/Closed/Merged)
- 主体描述/内容 (Body/Description)
- 所有评论 (Comments)

#### 2.3.2 可选内容

**Reactions 统计 (通过 -enable-reactions 控制):**
- 主楼 reactions (👍👎😄🎉😕❤️🚀👀)
- 每条评论的 reactions

**用户链接 (通过 -enable-user-links 控制):**
- `@username` 转换为 `[[@username](https://github.com/username)]`

#### 2.3.3 特殊处理规则

**Pull Request:**
- 仅包含 PR 描述和评论，**不包含**代码 diff
- Review comments 与普通评论按时间统一排序
- 如果 PR 已合并，状态显示为 "Merged"

**Discussion:**
- 被标记为 Answer 的评论需要特殊标识
- 支持 Discussion 的特殊状态 (Open/Closed/Answered)

### 2.4 Markdown 输出格式

#### 2.4.1 YAML Frontmatter

每个输出的 Markdown 文件必须包含 YAML frontmatter:

```yaml
---
title: "GitHub Issue Title"
url: "https://github.com/owner/repo/issues/123"
author: username
author_url: "https://github.com/username"
created_at: "2024-01-01T10:00:00Z"
updated_at: "2024-01-02T15:30:00Z"
status: "open" # open/closed/merged
type: "issue" # issue/pr/discussion
reaction_counts:
  thumbs_up: 5
  thumbs_down: 0
  laugh: 2
  hooray: 1
  confused: 0
  heart: 3
  rocket: 0
  eyes: 1
total_comments: 15
---
```

#### 2.4.2 Markdown 正文结构

```markdown
# [Issue Title] - Open/Closed

**作者:** @username
**创建时间:** 2024-01-01 10:00:00 UTC
**最后更新:** 2024-01-02 15:30:00 UTC
**状态:** Open
**评论数:** 15

## Description
[原始内容，保持原有 Markdown 格式...]

## Comments (15)

### @user1 - 2024-01-01 11:00:00 UTC
[评论内容...]

### @user2 - 2024-01-01 12:00:00 UTC
[评论内容...]

### ✅ @author - 2024-01-01 13:00:00 UTC [Accepted Answer]
[被标记为答案的评论内容...]
```

---

## 3. 非功能性需求 (Non-Functional Requirements)

### 3.1 架构设计

**核心原则:**
- 遵循 Go 语言"少即是多"哲学
- 模块化设计，便于测试和维护
- 使用标准库优先，最小化外部依赖

**模块结构:**
```
internal/
├── github/      # GitHub API 交互
├── parser/      # URL 解析与类型识别
├── converter/   # 数据转换为 Markdown
├── cli/         # 命令行接口
└── config/      # 配置管理
```

### 3.2 错误处理

**错误分类:**
1. **URL 错误**: 无效格式、不支持的类型
2. **网络错误**: API 请求失败、超时
3. **API 错误**: 资源不存在、权限不足、限流
4. **文件错误**: 输出文件无法写入

**错误处理原则:**
- 所有错误必须被显式处理
- 使用 `fmt.Errorf("...: %w", err)` 进行错误包装
- 友好的错误信息输出到 stderr
- 适当的退出码 (exit code)

### 3.3 性能要求

- 单个 Issue/PR/Discussion 处理时间 < 5 秒
- 内存使用 < 50MB
- 支持 GitHub API 限流处理

### 3.4 兼容性

- Go 版本: >= 1.21
- 操作系统: Linux, macOS, Windows
- GitHub API: v4 (GraphQL)

---

## 4. 验收标准 (Acceptance Criteria)

### 4.1 基本功能测试

**测试用例 1: Issue 转换**
```bash
# Given: 一个有效的 GitHub Issue URL
# When: 执行 issue2md 命令
# Then: 输出包含完整 Issue 信息的 Markdown

issue2md "https://github.com/golang/go/issues/12345"
```

**预期结果:**
- 输出有效的 Markdown 格式
- 包含 YAML frontmatter
- 包含 Issue 标题、描述、所有评论
- 时间戳格式正确

**测试用例 2: PR 转换**
```bash
# Given: 一个有效的 GitHub PR URL
# When: 执行 issue2md 命令
# Then: 输出包含 PR 描述和评论的 Markdown

issue2md "https://github.com/golang/go/pull/12345"
```

**测试用例 3: Discussion 转换**
```bash
# Given: 一个有效的 GitHub Discussion URL
# When: 执行 issue2md 命令
# Then: 输出包含 Discussion 内容和答案标识

issue2md "https://github.com/golang/go/discussions/12345"
```

### 4.2 功能标志测试

**测试用例 4: Reactions 支持**
```bash
issue2md -enable-reactions "https://github.com/golang/go/issues/12345"
```
**预期结果:** Markdown 中包含 reactions 统计信息

**测试用例 5: 用户链接支持**
```bash
issue2md -enable-user-links "https://github.com/golang/go/issues/12345"
```
**预期结果:** @username 被转换为 GitHub 主页链接

### 4.3 错误处理测试

**测试用例 6: 无效 URL**
```bash
issue2md "invalid-url"
```
**预期结果:** 返回明确的错误信息，非零退出码

**测试用例 7: 不存在的资源**
```bash
issue2md "https://github.com/golang/go/issues/99999"
```
**预期结果:** 返回资源不存在的错误信息

### 4.4 文件输出测试

**测试用例 8: 输出到文件**
```bash
issue2md "https://github.com/golang/go/issues/12345" output.md
```
**预期结果:** 内容正确写入指定文件

### 4.5 集成测试

**测试用例 9: 带 Token 的私有仓库访问**
```bash
export GITHUB_TOKEN=ghp_xxx
issue2md "https://github.com/private/repo/issues/1"
```
**预期结果:** 能够成功访问私有仓库 (如果 token 有效)

---

## 5. 输出格式示例

### 5.1 Issue 转换示例

**输入 URL:** `https://github.com/bigwhite/issue2md/issues/1`

**输出 Markdown:**
```markdown
---
title: "Add support for GitHub Discussions"
url: "https://github.com/bigwhite/issue2md/issues/1"
author: johndoe
author_url: "https://github.com/johndoe"
created_at: "2024-01-01T10:00:00Z"
updated_at: "2024-01-02T15:30:00Z"
status: "open"
type: "issue"
reaction_counts:
  thumbs_up: 8
  thumbs_down: 0
  laugh: 1
  hooray: 3
  confused: 0
  heart: 5
  rocket: 2
  eyes: 1
total_comments: 12
---

# Add support for GitHub Discussions - Open

**作者:** @johndoe
**创建时间:** 2024-01-01 10:00:00 UTC
**最后更新:** 2024-01-02 15:30:00 UTC
**状态:** Open
**评论数:** 12

## Description
Currently, issue2md only supports Issues and Pull Requests. It would be great to also support GitHub Discussions.

### Requirements
- Parse Discussion URLs
- Handle Answer marking
- Support Discussion reactions

## Comments (12)

### @alice - 2024-01-01 11:00:00 UTC
Great idea! Discussions are becoming more important for community engagement.

### @bob - 2024-01-01 12:30:00 UTC
I agree. This would be very useful for documenting community decisions.

### ✅ @johndoe - 2024-01-02 15:30:00 UTC [Accepted Answer]
Thanks for the feedback! I'll start working on this feature. The main challenge will be handling the different data structure for Discussions vs Issues.
```

### 5.2 PR 转换示例

**输入 URL:** `https://github.com/bigwhite/issue2md/pull/42`

**输出 Markdown:**
```markdown
---
title: "feat: add GitHub API client"
url: "https://github.com/bigwhite/issue2md/pull/42"
author: contributor
author_url: "https://github.com/contributor"
created_at: "2024-01-05T09:00:00Z"
updated_at: "2024-01-06T14:00:00Z"
status: "merged"
type: "pr"
reaction_counts:
  thumbs_up: 15
  thumbs_down: 0
  laugh: 0
  hooray: 8
  confused: 0
  heart: 12
  rocket: 6
  eyes: 2
total_comments: 8
---

# feat: add GitHub API client - Merged

**作者:** @contributor
**创建时间:** 2024-01-05 09:00:00 UTC
**最后更新:** 2024-01-06 14:00:00 UTC
**状态:** Merged
**评论数:** 8

## Description
This PR adds a GitHub API client to interact with the GitHub GraphQL API for fetching issue, PR, and discussion data.

### Changes Made
- Added GitHub GraphQL client
- Implemented basic query builders
- Added authentication support via environment variables

## Comments (8)

### @maintainer1 - 2024-01-05 10:00:00 UTC
Looks good! I have a few suggestions on the GraphQL query structure...

### @reviewer1 - 2024-01-05 11:30:00 UTC
The authentication approach looks solid. Have you considered rate limiting?
```

---

## 6. 实现注意事项

### 6.1 GitHub API 使用

- 使用 GitHub GraphQL API v4
- 实现 Basic Rate Limiting 处理
- 支持匿名访问 (公开仓库) 和 Token 认证访问

### 6.2 安全考虑

- Token 不通过命令行参数传递 (避免 shell 历史泄露)
- 不记录敏感信息到日志文件
- 适当的输入验证和清理

### 6.3 测试策略

- 单元测试：每个模块独立测试
- 集成测试：完整流程测试
- Mock GitHub API 进行测试
- 表格驱动测试优先

---

**下一步:** 根据此规格文档开始实现，优先完成 CLI 核心功能。