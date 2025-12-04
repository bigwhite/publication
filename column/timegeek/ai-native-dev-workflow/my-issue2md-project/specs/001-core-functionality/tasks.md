# issue2md 实现任务列表

**版本**: 1.0
**基于**: specs/001-core-functionality/plan.md
**原则**: TDD（测试先行），原子化任务，依赖关系明确

---

## Phase 1: Foundation (数据结构定义)

### 1.1 项目基础设置

#### 1.1.1 依赖管理
1. **[P]** 创建 `go.mod` 依赖声明
   - 添加 `github.com/google/go-github` 依赖
   - 添加 `golang.org/x/oauth2` 依赖

### 1.2 核心数据结构 (internal/github/types.go)

#### 1.2.1 基础类型定义测试
2. 创建 `internal/github/types_test.go` - User结构体测试
3. 创建 `internal/github/types_test.go` - Comment结构体测试
4. 创建 `internal/github/types_test.go` - Reactions结构体测试
5. **[P]** 创建 `internal/github/types_test.go` - Label结构体测试
6. **[P]** 创建 `internal/github/types_test.go` - Milestone结构体测试
7. **[P]** 创建 `internal/github/types_test.go` - DiscussionCategory结构体测试

#### 1.2.2 基础类型实现
8. 创建 `internal/github/types.go` - User结构体
9. 创建 `internal/github/types.go` - Comment结构体
10. 创建 `internal/github/types.go` - Reactions结构体
11. **[P]** 创建 `internal/github/types.go` - Label结构体
12. **[P]** 创建 `internal/github/types.go` - Milestone结构体
13. **[P]** 创建 `internal/github/types.go` - DiscussionCategory结构体

#### 1.2.3 资源类型测试
14. 创建 `internal/github/types_test.go` - Resource基础结构体测试
15. 创建 `internal/github/types_test.go` - Issue结构体测试
16. 创建 `internal/github/types_test.go` - PullRequest结构体测试
17. 创建 `internal/github/types_test.go` - Discussion结构体测试

#### 1.2.4 资源类型实现
18. 创建 `internal/github/types.go` - Resource基础结构体
19. 创建 `internal/github/types.go` - Issue结构体
20. 创建 `internal/github/types.go` - PullRequest结构体
21. 创建 `internal/github/types.go` - Discussion结构体

### 1.3 URL解析数据结构 (internal/parser/types.go)

#### 1.3.1 URL类型测试
22. 创建 `internal/parser/types_test.go` - ResourceURL结构体测试
23. 创建 `internal/parser/types_test.go` - URLPattern结构体测试
24. **[P]** 创建 `internal/parser/types_test.go` - SupportedPatterns常量测试

#### 1.3.2 URL类型实现
25. 创建 `internal/parser/types.go` - ResourceURL结构体
26. 创建 `internal/parser/types.go` - URLPattern结构体
27. **[P]** 创建 `internal/parser/types.go` - SupportedPatterns常量定义

### 1.4 配置数据结构 (internal/config/config.go)

#### 1.4.1 配置类型测试
28. 创建 `internal/config/config_test.go` - Config结构体测试
29. **[P]** 创建 `internal/config/config_test.go` - ConvertOptions结构体测试
30. **[P]** 创建 `internal/config/config_test.go` - CLIArgs结构体测试

#### 1.4.2 配置类型实现
31. 创建 `internal/config/config.go` - Config结构体
32. **[P]** 创建 `internal/config/config.go` - ConvertOptions结构体
33. **[P]** 创建 `internal/config/config.go` - CLIArgs结构体

### 1.5 统一错误定义

#### 1.5.1 错误测试
34. **[P]** 创建 `internal/github/errors_test.go` - 统一错误类型测试

#### 1.5.2 错误实现
35. **[P]** 创建 `internal/github/errors.go` - 统一错误类型定义

---

## Phase 2: URL Parser (URL解析逻辑，TDD)

### 2.1 URL解析核心功能 (internal/parser/parser.go)

#### 2.1.1 解析器测试
36. 创建 `internal/parser/parser_test.go` - Parse方法基础测试
37. 创建 `internal/parser/parser_test.go` - Parse方法Issue URL测试
38. 创建 `internal/parser/parser_test.go` - Parse方法PR URL测试
39. 创建 `internal/parser/parser_test.go` - Parse方法Discussion URL测试
40. 创建 `internal/parser/parser_test.go` - Parse方法无效URL测试
41. 创建 `internal/parser/parser_test.go` - Parse方法不支持URL测试

#### 2.1.2 解析器实现
42. 创建 `internal/parser/parser.go` - Parser接口定义
43. 创建 `internal/parser/parser.go` - Parse方法实现
44. 创建 `internal/parser/parser.go` - NewParser构造函数

### 2.2 URL验证功能 (internal/parser/validation.go)

#### 2.2.1 验证功能测试
45. 创建 `internal/parser/validation_test.go` - Validate方法测试
46. 创建 `internal/parser/validation_test.go` - SupportedTypes方法测试
47. **[P]** 创建 `internal/parser/validation_test.go` - 边界情况测试

#### 2.2.2 验证功能实现
48. 创建 `internal/parser/validation.go` - Validate方法实现
49. 创建 `internal/parser/validation.go` - SupportedTypes方法实现

### 2.3 便捷函数 (internal/parser/parser.go)

#### 2.3.1 便捷函数测试
50. 创建 `internal/parser/parser_test.go` - ParseGitHubURL函数测试

#### 2.3.2 便捷函数实现
51. 在 `internal/parser/parser.go` 中添加 ParseGitHubURL函数

---

## Phase 3: GitHub Fetcher (API交互逻辑，TDD)

### 3.1 GitHub认证处理 (internal/github/auth.go)

#### 3.1.1 认证测试
52. 创建 `internal/github/auth_test.go` - Token认证测试
53. 创建 `internal/github/auth_test.go` - 匿名访问测试
54. **[P]** 创建 `internal/github/auth_test.go` - 无效Token测试

#### 3.1.2 认证实现
55. 创建 `internal/github/auth.go` - Token认证实现
56. 创建 `internal/github/auth.go` - HTTP客户端构建

### 3.2 GraphQL查询构建 (internal/github/queries.go)

#### 3.2.1 查询构建测试
57. 创建 `internal/github/queries_test.go` - Issue查询构建测试
58. 创建 `internal/github/queries_test.go` - PR查询构建测试
59. 创建 `internal/github/queries_test.go` - Discussion查询构建测试
60. **[P]** 创建 `internal/github/queries_test.go` - 查询参数转义测试

#### 3.2.2 查询构建实现
61. 创建 `internal/github/queries.go` - Issue GraphQL查询
62. 创建 `internal/github/queries.go` - PR GraphQL查询
63. 创建 `internal/github/queries.go` - Discussion GraphQL查询

### 3.3 GitHub客户端实现 (internal/github/client.go)

#### 3.3.1 客户端接口测试
64. 创建 `internal/github/client_test.go` - GitHubClient接口测试
65. 创建 `internal/github/client_test.go` - GetIssue方法测试
66. 创建 `internal/github/client_test.go` - GetPullRequest方法测试
67. 创建 `internal/github/client_test.go` - GetDiscussion方法测试

#### 3.3.2 客户端实现
68. 创建 `internal/github/client.go` - GitHubClient接口定义
69. 创建 `internal/github/client.go` - GetIssue方法实现
70. 创建 `internal/github/client.go` - GetPullRequest方法实现
71. 创建 `internal/github/client.go` - GetDiscussion方法实现

### 3.4 客户端构建器 (internal/github/client.go)

#### 3.4.1 构建器测试
72. 创建 `internal/github/client_test.go` - NewGitHubClient测试
73. 创建 `internal/github/client_test.go` - GitHubClientBuilder测试
74. **[P]** 创建 `internal/github/client_test.go` - 构建器链式调用测试

#### 3.4.2 构建器实现
75. 在 `internal/github/client.go` 中添加 NewGitHubClient函数
76. 在 `internal/github/client.go` 中添加 GitHubClientBuilder接口
77. 在 `internal/github/client.go` 中添加构建器实现

### 3.5 错误处理和限流 (internal/github/client.go)

#### 3.5.1 错误处理测试
78. 创建 `internal/github/client_test.go` - API错误处理测试
79. 创建 `internal/github/client_test.go` - 网络超时测试
80. **[P]** 创建 `internal/github/client_test.go` - 限流重试测试

#### 3.5.2 错误处理实现
81. 在 `internal/github/client.go` 中添加错误处理逻辑
82. 在 `internal/github/client.go` 中添加限流重试逻辑

---

## Phase 4: Markdown Converter (转换逻辑，TDD)

### 4.1 YAML Frontmatter生成 (internal/converter/frontmatter.go)

#### 4.1.1 Frontmatter测试
83. 创建 `internal/converter/frontmatter_test.go` - Generate方法测试
84. 创建 `internal/converter/frontmatter_test.go` - ToYAML方法测试
85. 创建 `internal/converter/frontmatter_test.go` - Issue frontmatter测试
86. 创建 `internal/converter/frontmatter_test.go` - PR frontmatter测试
87. 创建 `internal/converter/frontmatter_test.go` - Discussion frontmatter测试

#### 4.1.2 Frontmatter实现
88. 创建 `internal/converter/frontmatter.go` - FrontmatterGenerator接口
89. 创建 `internal/converter/frontmatter.go` - Generate方法实现
90. 创建 `internal/converter/frontmatter.go` - ToYAML方法实现

### 4.2 内容格式化 (internal/converter/formatter.go)

#### 4.2.1 格式化测试
91. 创建 `internal/converter/formatter_test.go` - FormatComment方法测试
92. 创建 `internal/converter/formatter_test.go` - FormatReactions方法测试
93. 创建 `internal/converter/formatter_test.go` - FormatUserLink方法测试
94. **[P]** 创建 `internal/converter/formatter_test.go` - 时间格式化测试

#### 4.2.2 格式化实现
95. 创建 `internal/converter/formatter.go` - FormatComment方法
96. 创建 `internal/converter/formatter.go` - FormatReactions方法
97. 创建 `internal/converter/formatter.go` - FormatUserLink方法
98. **[P]** 在 `internal/converter/formatter.go` 中添加时间格式化函数

### 4.3 Markdown模板 (internal/converter/templates.go)

#### 4.3.1 模板测试
99. 创建 `internal/converter/templates_test.go` - Issue模板测试
100. 创建 `internal/converter/templates_test.go` - PR模板测试
101. 创建 `internal/converter/templates_test.go` - Discussion模板测试
102. **[P]** 创建 `internal/converter/templates_test.go` - 模板选项测试

#### 4.3.2 模板实现
103. 创建 `internal/converter/templates.go` - Issue模板
104. 创建 `internal/converter/templates.go` - PR模板
105. 创建 `internal/converter/templates.go` - Discussion模板
106. **[P]** 创建 `internal/converter/templates.go` - 模板选项处理

### 4.4 转换器核心 (internal/converter/converter.go)

#### 4.4.1 转换器测试
107. 创建 `internal/converter/converter_test.go` - Converter接口测试
108. 创建 `internal/converter/converter_test.go` - Convert方法测试
109. 创建 `internal/converter/converter_test.go` - ConvertIssue方法测试
110. 创建 `internal/converter/converter_test.go` - ConvertPullRequest方法测试
111. 创建 `internal/converter/converter_test.go` - ConvertDiscussion方法测试

#### 4.4.2 转换器实现
112. 创建 `internal/converter/converter.go` - Converter接口定义
113. 创建 `internal/converter/converter.go` - Convert方法实现
114. 创建 `internal/converter/converter.go` - ConvertIssue方法实现
115. 创建 `internal/converter/converter.go` - ConvertPullRequest方法实现
116. 创建 `internal/converter/converter.go` - ConvertDiscussion方法实现

### 4.5 转换器构建器 (internal/converter/converter.go)

#### 4.5.1 构建器测试
117. 创建 `internal/converter/converter_test.go` - NewConverter测试

#### 4.5.2 构建器实现
118. 在 `internal/converter/converter.go` 中添加 NewConverter函数

---

## Phase 5: Configuration Management (配置管理，TDD)

### 5.1 环境变量处理 (internal/config/env.go)

#### 5.1.1 环境变量测试
119. 创建 `internal/config/env_test.go` - LoadFromEnv测试
120. 创建 `internal/config/env_test.go` - GITHUB_TOKEN环境变量测试
121. **[P]** 创建 `internal/config/env_test.go` - 默认值测试

#### 5.1.2 环境变量实现
122. 创建 `internal/config/env.go` - LoadFromEnv函数
123. 创建 `internal/config/env.go` - 环境变量读取逻辑

### 5.2 配置实现 (internal/config/config.go)

#### 5.2.1 配置测试
124. 创建 `internal/config/config_test.go` - Config接口测试
125. 创建 `internal/config/config_test.go` - GitHubToken方法测试
126. 创建 `internal/config/config_test.go` - UserAgent方法测试
127. 创建 `internal/config/config_test.go` - APITimeout方法测试
128. **[P]** 创建 `internal/config/config_test.go` - HTTPClient方法测试

#### 5.2.2 配置实现
129. 在 `internal/config/config.go` 中添加 Config接口
130. 在 `internal/config/config.go` 中添加 Config结构体实现
131. 在 `internal/config/config.go` 中添加 NewConfig函数
132. **[P]** 在 `internal/config/config.go` 中添加 HTTPClient方法

### 5.3 配置加载器 (internal/config/config.go)

#### 5.3.1 加载器测试
133. 创建 `internal/config/config_test.go` - ConfigLoader接口测试
134. 创建 `internal/config/config_test.go` - LoadWithToken测试
135. **[P]** 创建 `internal/config/config_test.go` - Validate测试

#### 5.3.2 加载器实现
136. 在 `internal/config/config.go` 中添加 ConfigLoader接口
137. 在 `internal/config/config.go` 中添加 LoadWithToken函数
138. **[P]** 在 `internal/config/config.go` 中添加 Validate函数

---

## Phase 6: CLI Assembly (命令行入口集成)

### 6.1 命令行参数解析 (internal/cli/args.go)

#### 6.1.1 参数解析测试
139. 创建 `internal/cli/args_test.go` - ArgParser接口测试
140. 创建 `internal/cli/args_test.go` - Parse方法测试
141. 创建 `internal/cli/args_test.go` - Validate方法测试
142. 创建 `internal/cli/args_test.go` - ShowUsage测试
143. 创建 `internal/cli/args_test.go` - ShowVersion测试
144. **[P]** 创建 `internal/cli/args_test.go` - 参数边界测试

#### 6.1.2 参数解析实现
145. 创建 `internal/cli/args.go` - ArgParser接口
146. 创建 `internal/cli/args.go` - Parse方法实现
147. 创建 `internal/cli/args.go` - Validate方法实现
148. 创建 `internal/cli/args.go` - ShowUsage方法实现
149. 创建 `internal/cli/args.go` - ShowVersion方法实现
150. **[P]** 创建 `internal/cli/args.go` - 参数边界处理

### 6.2 版本管理 (internal/cli/version.go)

#### 6.2.1 版本测试
151. **[P]** 创建 `internal/cli/version_test.go` - 版本信息测试

#### 6.2.2 版本实现
152. **[P]** 创建 `internal/cli/version.go` - 版本信息定义

### 6.3 CLI应用主逻辑 (internal/cli/app.go)

#### 6.3.1 CLI应用测试
153. 创建 `internal/cli/app_test.go` - CLIApp接口测试
154. 创建 `internal/cli/app_test.go` - Run方法测试
155. 创建 `internal/cli/app_test.go` - Execute方法测试
156. **[P]** 创建 `internal/cli/app_test.go` - 错误处理测试

#### 6.3.2 CLI应用实现
157. 创建 `internal/cli/app.go` - CLIApp接口
158. 创建 `internal/cli/app.go` - CLIApp结构体
159. 创建 `internal/cli/app.go` - Run方法实现
160. 创建 `internal/cli/app.go` - Execute方法实现
161. **[P]** 创建 `internal/cli/app.go` - 错误处理逻辑

### 6.4 文件输出处理 (internal/cli/app.go)

#### 6.4.1 输出处理测试
162. 创建 `internal/cli/app_test.go` - OutputHandler接口测试
163. 创建 `internal/cli/app_test.go` - WriteToFile测试
164. 创建 `internal/cli/app_test.go` - WriteToStdout测试
165. **[P]** 创建 `internal/cli/app_test.go` - EnsureDirectory测试

#### 6.4.2 输出处理实现
166. 在 `internal/cli/app.go` 中添加 OutputHandler接口
167. 在 `internal/cli/app.go` 中添加 WriteToFile方法
168. 在 `internal/cli/app.go` 中添加 WriteToStdout方法
169. **[P]** 在 `internal/cli/app.go` 中添加 EnsureDirectory方法

### 6.5 应用构建器 (internal/cli/app.go)

#### 6.5.1 构建器测试
170. 创建 `internal/cli/app_test.go` - NewCLIApp测试

#### 6.5.2 构建器实现
171. 在 `internal/cli/app.go` 中添加 NewCLIApp函数

---

## Phase 7: Main Program (主程序入口)

### 7.1 应用入口点 (cmd/issue2md/main.go)

#### 7.1.1 主程序测试
172. 创建 `cmd/issue2md/main_test.go` - main函数测试
173. 创建 `cmd/issue2md/main_test.go` - 集成测试准备

#### 7.1.2 主程序实现
174. 创建 `cmd/issue2md/main.go` - main函数实现

### 7.2 构建脚本 (Makefile)

#### 7.2.1 构建脚本实现
175. 创建 `Makefile` - 构建规则
176. 创建 `Makefile` - 测试规则
177. 创建 `Makefile` - 清理规则
178. **[P]** 创建 `Makefile` - 安装规则

---

## Phase 8: Integration & Testing (集成与测试)

### 8.1 端到端测试

#### 8.1.1 E2E测试创建
179. 创建 `internal/integration_test.go` - Issue转换E2E测试
180. 创建 `internal/integration_test.go` - PR转换E2E测试
181. 创建 `internal/integration_test.go` - Discussion转换E2E测试
182. **[P]** 创建 `internal/integration_test.go` - 错误场景E2E测试

### 8.2 性能测试

#### 8.2.1 性能测试创建
183. **[P]** 创建 `internal/benchmark_test.go` - API调用性能测试
184. **[P]** 创建 `internal/benchmark_test.go` - 转换性能测试

### 8.3 测试数据准备

#### 8.3.1 测试数据
185. 创建 `testdata/` 目录结构
186. 准备GitHub测试仓库数据

---

## 任务执行指南

### 并行执行说明
- **[P]** 标记的任务可以并行执行，无依赖关系
- 同一阶段内的[P]任务可以同时开发
- 跨阶段的任务必须按依赖顺序执行

### TDD执行流程
每个功能模块遵循以下流程：
1. **先创建测试文件**（任务编号为奇数）
2. **运行测试确保失败**（Red阶段）
3. **实现对应功能**（任务编号为偶数）
4. **运行测试确保通过**（Green阶段）
5. **重构优化代码**（Refactor阶段）

### 错误处理原则
- 所有错误必须显式处理
- 使用 `fmt.Errorf("context: %w", err)` 进行错误包装
- 提供清晰的错误信息

### 代码质量标准
- 所有函数必须有文档注释
- 遵循Go官方代码规范
- 使用 `golangci-lint` 进行静态检查
- 单元测试覆盖率 >= 90%

### 验收标准
- 所有测试通过
- 代码质量检查通过
- 性能指标达标（API响应 < 5秒，内存 < 50MB）
- 完整的CLI功能可用

---

**下一步**: 从Phase 1开始，按照TDD流程逐个执行任务。每完成一个任务后，运行测试确保功能正确性。