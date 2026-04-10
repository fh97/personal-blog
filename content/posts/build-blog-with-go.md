---
title: 用 Go 从零搭建个人博客
date: 2026-04-10
tags: [Go, Blog, Web]
summary: 记录一次用纯 Go 标准库 + Tailwind CSS 搭建个人博客的过程，从路由、模板到 Markdown 渲染，全程不依赖框架。
---

## 为什么选择 Go？

Go 的标准库 `net/http` 足够强大，对于一个个人博客来说完全够用。不需要引入 Gin、Echo 等框架，保持依赖最小化。

## 项目结构

```
├── main.go
├── model/
├── service/
├── handler/
├── templates/
└── content/posts/
```

## Markdown 解析

使用 `goldmark` 将 Markdown 渲染为 HTML：

```go
var buf bytes.Buffer
if err := goldmark.Convert([]byte(body), &buf); err != nil {
    return "", err
}
return buf.String(), nil
```

## 总结

整个博客不到 500 行 Go 代码，模板用 Tailwind CSS，部署一个二进制文件搞定。
