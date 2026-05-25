# go-tiny-claw: 极简智能体驾驭引擎

这是一个基于“驾驭工程 (Harness Engineering)”理念，由 Go 语言从零实现的微型 AI Agent 操作系统。

## 核心设计哲学
- **Harness over Framework**: 真正的壁垒不在于调用大模型 API，而在于如何调度工具、管理上下文和安全拦截。
- **极简即是正义**: 我们拒绝臃肿的插件，仅向大模型提供 Read、Write、Edit 和 Bash 四大图灵完备的原语。
- **状态外部化**: 抛弃内存状态机，将记忆与计划持久化在 PLAN.md 与 TODO.md 中。

## 当前系统组件
1. **心脏 (Main Loop)**: 纯手写的 ReAct 循环。
2. **大脑 (Provider)**: 适配了官方 Claude SDK 与智谱 GLM API。
3. **手脚 (Tool Registry)**: 支持并发执行的动态工具集。
4. **神经元 (Reporter)**: 已成功接入飞书群聊事件流。

---
> "在大模型时代，每一位工程师都应该拥有属于自己的 Agent 驱动引擎。"

*本仓库由极客时间专栏《从零构建 Agent Harness》实战产出。*
*作者：Tony Bai*
