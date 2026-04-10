package model

import "time"

// Post 表示一篇博客文章，包含元数据和可选的渲染后正文内容。
type Post struct {
	// Title 是文章标题，来自 frontmatter 的 title 字段。
	Title string

	// Date 是文章发布日期，来自 frontmatter 的 date 字段（格式 2006-01-02）。
	Date time.Time

	// Tags 是文章标签列表，来自 frontmatter 的 tags 字段。
	Tags []string

	// Summary 是文章摘要，来自 frontmatter 的 summary 字段。
	Summary string

	// Slug 是文章的唯一标识符，由文件名去掉 .md 后缀得到。
	Slug string

	// Content 是渲染后的 HTML 正文内容。
	// 在主页列表场景下为空字符串，仅在详情页场景下填充。
	Content string
}
