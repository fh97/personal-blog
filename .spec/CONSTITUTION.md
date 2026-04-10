# 项目原则（Constitution）

## 技术栈
- Go 1.24，标准库 `net/http` + `html/template`
- 外部依赖：`github.com/yuin/goldmark`（Markdown 渲染）
- 前端样式：Tailwind CSS（CDN，无构建步骤）
- 文章内容：Markdown 文件（`/content/posts/*.md`），带 YAML frontmatter

## 目录结构
```
├── main.go              # 入口，路由注册
├── model/post.go        # Post 数据结构
├── service/post.go      # 读取/解析文章业务逻辑
├── handler/blog.go      # HTTP Handler
├── templates/           # HTML 模板
│   ├── base.html
│   ├── index.html       # 主页
│   └── post.html        # 文章详情页
└── content/posts/       # Markdown 文章
```

## 代码规范
- 所有公开函数需有注释
- 错误必须向上返回，不在 service 层 panic
- Handler 层统一处理 500/404 响应

## 质量要求
- `go build ./...` 零报错
- service 层核心函数有单元测试
