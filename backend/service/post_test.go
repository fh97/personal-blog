package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testPostsDir = "/tmp/test-posts"

// testPost1 是第一篇测试文章的内容（较新的日期）。
var testPost1 = `---
title: Hello Go
date: 2024-03-15
tags: [Go, Blog]
summary: Introduction to Go programming
---

# Hello Go

This is the **first** post about Go.

- Item 1
- Item 2
`

// testPost2 是第二篇测试文章的内容（较旧的日期）。
var testPost2 = `---
title: My Second Post
date: 2023-11-20
tags: [Blog, Tips]
summary: Some useful tips
---

## Second Post

Here is some _markdown_ content.
`

func setupTestPosts(t *testing.T) {
	t.Helper()

	if err := os.MkdirAll(testPostsDir, 0755); err != nil {
		t.Fatalf("failed to create test posts dir: %v", err)
	}

	files := map[string]string{
		"hello-go.md":     testPost1,
		"second-post.md":  testPost2,
	}

	for name, content := range files {
		path := filepath.Join(testPostsDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file %s: %v", name, err)
		}
	}
}

func cleanupTestPosts(t *testing.T) {
	t.Helper()
	if err := os.RemoveAll(testPostsDir); err != nil {
		t.Logf("warning: failed to cleanup test posts dir: %v", err)
	}
}

func TestGetAllPosts(t *testing.T) {
	setupTestPosts(t)
	defer cleanupTestPosts(t)

	posts, err := GetAllPosts(testPostsDir)
	if err != nil {
		t.Fatalf("GetAllPosts returned error: %v", err)
	}

	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}

	// 验证按日期倒序排列：第一篇应为 2024-03-15
	if posts[0].Slug != "hello-go" {
		t.Errorf("expected first post slug to be 'hello-go', got %q", posts[0].Slug)
	}
	if posts[1].Slug != "second-post" {
		t.Errorf("expected second post slug to be 'second-post', got %q", posts[1].Slug)
	}

	// 验证第一篇文章元数据
	first := posts[0]
	if first.Title != "Hello Go" {
		t.Errorf("expected title 'Hello Go', got %q", first.Title)
	}

	expectedDate := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if !first.Date.Equal(expectedDate) {
		t.Errorf("expected date %v, got %v", expectedDate, first.Date)
	}

	if len(first.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d: %v", len(first.Tags), first.Tags)
	} else {
		if first.Tags[0] != "Go" {
			t.Errorf("expected first tag 'Go', got %q", first.Tags[0])
		}
		if first.Tags[1] != "Blog" {
			t.Errorf("expected second tag 'Blog', got %q", first.Tags[1])
		}
	}

	if first.Summary != "Introduction to Go programming" {
		t.Errorf("expected summary 'Introduction to Go programming', got %q", first.Summary)
	}

	// 主页列表场景：Content 应为空
	if first.Content != "" {
		t.Errorf("expected Content to be empty in list mode, got %q", first.Content)
	}

	// 验证第二篇文章元数据
	second := posts[1]
	if second.Title != "My Second Post" {
		t.Errorf("expected title 'My Second Post', got %q", second.Title)
	}

	expectedDate2 := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)
	if !second.Date.Equal(expectedDate2) {
		t.Errorf("expected date %v, got %v", expectedDate2, second.Date)
	}
}

func TestGetAllPosts_EmptyDir(t *testing.T) {
	emptyDir := "/tmp/test-posts-empty"
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("failed to create empty dir: %v", err)
	}
	defer os.RemoveAll(emptyDir)

	posts, err := GetAllPosts(emptyDir)
	if err != nil {
		t.Fatalf("GetAllPosts on empty dir returned error: %v", err)
	}

	if len(posts) != 0 {
		t.Errorf("expected 0 posts from empty dir, got %d", len(posts))
	}
}

func TestGetAllPosts_NonExistentDir(t *testing.T) {
	_, err := GetAllPosts("/tmp/non-existent-posts-dir-xyz")
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}

func TestGetPostBySlug(t *testing.T) {
	setupTestPosts(t)
	defer cleanupTestPosts(t)

	post, err := GetPostBySlug(testPostsDir, "hello-go")
	if err != nil {
		t.Fatalf("GetPostBySlug returned error: %v", err)
	}

	if post.Slug != "hello-go" {
		t.Errorf("expected slug 'hello-go', got %q", post.Slug)
	}

	if post.Title != "Hello Go" {
		t.Errorf("expected title 'Hello Go', got %q", post.Title)
	}

	// 详情页：Content 应包含渲染后的 HTML
	if post.Content == "" {
		t.Error("expected non-empty Content for detail page, got empty string")
	}

	// 验证渲染后的 HTML 包含预期的标签
	if !containsSubstring(post.Content, "<h1>") && !containsSubstring(post.Content, "<strong>") {
		t.Errorf("expected rendered HTML to contain heading or bold elements, got: %q", post.Content)
	}
}

func TestGetPostBySlug_NotFound(t *testing.T) {
	setupTestPosts(t)
	defer cleanupTestPosts(t)

	_, err := GetPostBySlug(testPostsDir, "non-existent-slug")
	if err == nil {
		t.Error("expected error for non-existent slug, got nil")
	}
}

func TestGetPostBySlug_SecondPost(t *testing.T) {
	setupTestPosts(t)
	defer cleanupTestPosts(t)

	post, err := GetPostBySlug(testPostsDir, "second-post")
	if err != nil {
		t.Fatalf("GetPostBySlug returned error: %v", err)
	}

	if post.Title != "My Second Post" {
		t.Errorf("expected title 'My Second Post', got %q", post.Title)
	}

	if len(post.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(post.Tags))
	}

	if post.Content == "" {
		t.Error("expected non-empty Content, got empty string")
	}

	// 验证渲染后的 HTML 包含 em 标签（markdown 中的 _markdown_）
	if !containsSubstring(post.Content, "<em>") {
		t.Errorf("expected rendered HTML to contain <em> tag, got: %q", post.Content)
	}
}

func TestGetPostsByTag(t *testing.T) {
	setupTestPosts(t)
	defer cleanupTestPosts(t)

	// 查找含 "Blog" tag 的文章（两篇都有）
	posts, err := GetPostsByTag(testPostsDir, "Blog")
	if err != nil {
		t.Fatalf("GetPostsByTag returned error: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("expected 2 posts with tag 'Blog', got %d", len(posts))
	}

	// 查找含 "Go" tag 的文章（只有第一篇）
	posts, err = GetPostsByTag(testPostsDir, "Go")
	if err != nil {
		t.Fatalf("GetPostsByTag returned error: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("expected 1 post with tag 'Go', got %d", len(posts))
	}

	if len(posts) > 0 && posts[0].Slug != "hello-go" {
		t.Errorf("expected slug 'hello-go', got %q", posts[0].Slug)
	}

	// 查找不存在的 tag
	posts, err = GetPostsByTag(testPostsDir, "NonExistentTag")
	if err != nil {
		t.Fatalf("GetPostsByTag returned error: %v", err)
	}

	if len(posts) != 0 {
		t.Errorf("expected 0 posts for non-existent tag, got %d", len(posts))
	}
}

// containsSubstring 是测试辅助函数，检查字符串 s 是否包含子串 sub。
func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
