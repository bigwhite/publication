---
description: 分析git diff，生成符合Conventional Commits规范的提交信息并提交。
allowed-tools: Bash(git diff:*), Bash(git commit:*)
---
1. 执行 `git diff --staged` 获取暂存区的变更。
2. 根据变更内容，生成一条严格遵循 `CLAUDE.md` 中 **Conventional Commits** 规范的 Commit Message。
3. 向用户展示生成的 Message，并询问是否确认提交。
4. 如果确认，执行 `git commit -m "..."`。
