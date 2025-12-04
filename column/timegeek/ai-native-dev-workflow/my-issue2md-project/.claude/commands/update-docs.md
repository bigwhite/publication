---
description: 扫描整个项目，并更新或重新生成README.md文档。
model: opus
---

你是一位优秀的技术文档工程师。请全面扫描当前目录（`@.`）下的所有文件，特别是`go.mod`, `Makefile`, 以及`cmd/`目录下的`main.go`，以理解项目的最新状态。

然后，请为`issue2md`项目，重新生成一份完整、清晰、与当前代码完全同步的`README.md`文件内容。

这份README必须包含以下部分：
1.  **项目简介 (Overview)**
2.  **核心特性 (Features)**
3.  **安装指南 (Installation)**
4.  **使用方法 (Usage):** 必须包含所有命令行参数的详细说明和示例。
5.  **构建方法 (Building from Source):** 必须引用`Makefile`中的命令。
