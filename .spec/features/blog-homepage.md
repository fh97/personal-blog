# Feature: 个人博客主页

## Goal
展示博主个人信息与文章列表，让访客快速了解博主并浏览文章。

## Acceptance Criteria
- [ ] 顶部 Hero 区：头像占位、姓名、一句话介绍、GitHub / Twitter 链接
- [ ] 文章列表按发布日期倒序排列
- [ ] 每篇文章卡片展示：标题、日期、标签（彩色 badge）、摘要（≤120字）
- [ ] 支持按标签点击筛选（URL query: `?tag=xxx`）
- [ ] 点击文章标题跳转到文章详情页（`/post/:slug`），渲染 Markdown 正文
- [ ] 响应式布局：移动端单列 / PC 端双列卡片
- [ ] 404 / 500 有友好错误页

## Frontmatter 格式（约定）
```markdown
---
title: 文章标题
date: 2026-04-10
tags: [Go, Blog]
summary: 这里是摘要，不超过120字。
---
正文内容...
```

## Constraints
- 纯静态文件驱动，不依赖数据库
- 文章 slug 即文件名（去掉 .md）
- 所有路由在 main.go 注册，handler 不感知路由

## Out of Scope（本期不做）
- 评论系统
- 全文搜索
- RSS Feed
- 暗色模式
