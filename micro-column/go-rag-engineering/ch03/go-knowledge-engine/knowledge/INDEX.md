# 团队知识库

本知识库采用 File-first 组织方式。

根目录下带有 ID 前缀的 Markdown 文件是由知识摄取流水线生成的规范化知识文件。

各主题目录下的 Markdown 文件是经过整理的主题知识，用于帮助 Agent 理解、导航和关联不同来源的知识。

## 知识主题

### 工程设计

- [工程知识索引](engineering/INDEX.md)
- [系统架构与服务边界](engineering/architecture.md)
- [服务调用超时与重试规范](engineering/timeout.md)

### 历史故障

- [故障知识索引](incidents/INDEX.md)
- [订单服务超时故障复盘](incidents/order-timeout.md)

### 运维管理

- [运维知识索引](operations/INDEX.md)
- [部署与监控规范](operations/deployment.md)

### 项目知识

- [项目知识索引](projects/INDEX.md)
- [订单服务项目概览](projects/order-service.md)

## 知识组织约定

1. 原始来源由摄取流水线维护，不应随意修改。
2. 主题知识用于归纳、解释和关联原始来源。
3. 每份主题知识应尽可能保留来源路径。
4. 新增主题知识后，应同步更新对应目录的 INDEX.md。
5. 主题知识之间可以通过相对路径建立 Markdown 链接。
