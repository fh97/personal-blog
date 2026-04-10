---
title: Spec-Driven Development 实践指南
date: 2026-04-08
tags: [工程效能, AI, Spec]
summary: 先写清楚规格，再让 AI 按规格生成代码。这不只是一种工具用法，更是一种思维方式的转变——让"做什么"先于"怎么做"。
---

## 什么是 Spec-Driven Development？

传统开发里，代码是核心，规格文档只是脚手架——写完就扔。

**Spec-Driven Development 把规格变成可执行的**，直接驱动生成可工作的实现。

## 三阶段工作流

### 1. Spec（定义约束）

```markdown
## Acceptance Criteria
- [ ] 支持 PDF、DOCX 格式，单文件 ≤ 50MB
- [ ] 上传成功后 5s 内返回文档 ID
```

### 2. Plan（拆分任务）

识别依赖关系，确定哪些 Task 可以并行执行。

### 3. Subagent Dev（分发执行）

每个 Subagent 只拿它需要的上下文，专注完成一个 Task。

## 核心原则

> **Spec 是合同，代码是履约。**
