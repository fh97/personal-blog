package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"icode.baidu.com/baidu/personal-code/fullStack/model"
)

// parseFrontmatter 解析 Markdown 文件内容，提取 frontmatter 元数据和正文。
// frontmatter 以 "---\n" 开头和结尾，中间为 key: value 格式的元数据。
// 返回解析后的 Post（不含 Content）和正文原始内容。
func parseFrontmatter(filename string, data []byte) (*model.Post, []byte, error) {
	const delimiter = "---\n"

	content := string(data)

	if !strings.HasPrefix(content, delimiter) {
		return nil, data, fmt.Errorf("file %s: missing frontmatter opening delimiter", filename)
	}

	// 跳过第一个 "---\n"
	rest := content[len(delimiter):]

	// 找第二个 "---\n"
	endIdx := strings.Index(rest, delimiter)
	if endIdx == -1 {
		return nil, data, fmt.Errorf("file %s: missing frontmatter closing delimiter", filename)
	}

	frontmatter := rest[:endIdx]
	body := rest[endIdx+len(delimiter):]

	post := &model.Post{}

	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])

		switch key {
		case "title":
			post.Title = value
		case "date":
			t, err := time.Parse("2006-01-02", value)
			if err != nil {
				return nil, nil, fmt.Errorf("file %s: invalid date format %q: %w", filename, value, err)
			}
			post.Date = t
		case "tags":
			// 解析格式 [Go, Blog] 或 [Go,Blog]
			value = strings.TrimPrefix(value, "[")
			value = strings.TrimSuffix(value, "]")
			if value != "" {
				parts := strings.Split(value, ",")
				for _, p := range parts {
					tag := strings.TrimSpace(p)
					if tag != "" {
						post.Tags = append(post.Tags, tag)
					}
				}
			}
		case "summary":
			post.Summary = value
		}
	}

	// slug 为文件名去掉 .md 后缀
	base := filepath.Base(filename)
	post.Slug = strings.TrimSuffix(base, ".md")

	return post, []byte(body), nil
}

// renderMarkdown 将 Markdown 内容渲染为 HTML 字符串。
func renderMarkdown(src []byte) (string, error) {
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return "", fmt.Errorf("markdown render error: %w", err)
	}
	return buf.String(), nil
}

// GetAllPosts 读取 postsDir 下所有 .md 文件，解析 frontmatter，按 Date 倒序返回。
// Content 字段留空，适用于主页列表场景，不加载正文内容。
func GetAllPosts(postsDir string) ([]*model.Post, error) {
	entries, err := os.ReadDir(postsDir)
	if err != nil {
		return nil, fmt.Errorf("read posts dir %q: %w", postsDir, err)
	}

	var posts []*model.Post

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}

		fullPath := filepath.Join(postsDir, name)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("read file %q: %w", fullPath, err)
		}

		post, _, err := parseFrontmatter(name, data)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	// 按 Date 倒序排列（最新的在前面）
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

// GetPostBySlug 根据 slug 找到对应的 .md 文件，解析 frontmatter 并渲染 Markdown 正文为 HTML。
// slug 对应文件名去掉 .md 后缀的部分。
func GetPostBySlug(postsDir, slug string) (*model.Post, error) {
	filename := slug + ".md"
	fullPath := filepath.Join(postsDir, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("post with slug %q not found", slug)
		}
		return nil, fmt.Errorf("read file %q: %w", fullPath, err)
	}

	post, body, err := parseFrontmatter(filename, data)
	if err != nil {
		return nil, err
	}

	html, err := renderMarkdown(body)
	if err != nil {
		return nil, fmt.Errorf("render markdown for slug %q: %w", slug, err)
	}

	post.Content = html
	return post, nil
}

// GetPostsByTag 从 postsDir 下的所有文章中筛选包含指定 tag 的文章，按 Date 倒序返回。
func GetPostsByTag(postsDir, tag string) ([]*model.Post, error) {
	all, err := GetAllPosts(postsDir)
	if err != nil {
		return nil, err
	}

	var result []*model.Post
	for _, post := range all {
		for _, t := range post.Tags {
			if t == tag {
				result = append(result, post)
				break
			}
		}
	}

	return result, nil
}
