---
id: 8ff830f9858167dcf6cce268b33ea10cf15e09991098562332cf1fb598dade6f
source: operations/deployment.md
format: md
updated_at: 2026-11-20T06:49:21Z
---

# 部署与监控规范

## Document Headings

- 部署与监控规范

# 部署与监控规范

服务部署后，需要持续监控请求延迟、错误率和资源使用情况。

对于 HTTP 服务，需要关注：

- P50、P95 和 P99 请求延迟。
- HTTP 5xx 错误比例。
- 活跃连接数和连接池使用率。
- 请求超时数量。
- 重试请求数量。

当 P99 延迟持续升高时，应结合上下游服务指标定位问题，
不能只根据单个服务的 CPU 使用率判断故障原因。
