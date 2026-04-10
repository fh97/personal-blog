# Plan: 个人博客主页实现

## 依赖关系
```
Task1 (数据层) ──┐
                 ├──► Task3 (Handler) ──► Task4 (组装验证)
Task2 (模板)  ──┘
```
Task1、Task2 可并行；Task3 依赖 Task1+Task2；Task4 最后执行。

---

## Task 1：数据层 [model + service]
**文件：`model/post.go`、`service/post.go`**

- `Post` struct：`Title`、`Date`、`Tags`、`Summary`、`Slug`、`Content（HTML）`
- `service.GetAllPosts(dir)`：读取 content/posts/*.md，解析 frontmatter，按 date 倒序
- `service.GetPostBySlug(dir, slug)`：读取单篇，渲染 Markdown → HTML
- `service.GetPostsByTag(dir, tag)`：按标签过滤
- frontmatter 用纯 Go 解析（`---` 分隔，strings.Split），无外部依赖
- Markdown 渲染用 `github.com/yuin/goldmark`

**交付：** `model/post.go` + `service/post.go` + `service/post_test.go`

---

## Task 2：HTML 模板 [templates]
**文件：`templates/base.html`、`templates/index.html`、`templates/post.html`、`templates/error.html`**

- `base.html`：引入 Tailwind CSS CDN，定义 `{{block "content"}}`
- `index.html`：
  - Hero Section（头像、名字、简介、社交链接）
  - TagFilter（标签按钮，当前激活标签高亮）
  - PostList（响应式双列卡片：标题、日期、tags badge、摘要）
- `post.html`：文章详情，渲染 HTML 正文
- `error.html`：404 / 500 友好页面

**交付：** 4 个 HTML 模板文件

---

## Task 3：HTTP Handler + 路由 [handler]
**文件：`handler/blog.go`**

- `IndexHandler`：处理 `GET /`，读 `?tag=` 参数，调 service，渲染 index.html
- `PostHandler`：处理 `GET /post/{slug}`，调 service，渲染 post.html
- 统一错误处理：404 → error.html，500 → error.html
- Handler 不感知路由，参数从 `r.URL` 提取

**交付：** `handler/blog.go`

---

## Task 4：组装 + 示例内容 + 验证
**文件：`main.go`（更新）、`content/posts/*.md`（2-3 篇示例文章）**

- 更新 `main.go`：注册路由、静态文件服务、启动端口 8080
- 添加 `go get github.com/yuin/goldmark`
- 写 2-3 篇示例 Markdown 文章
- 执行 `go build ./...` 验证零报错
