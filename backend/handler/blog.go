// Copyright 2026 Baidu Inc. All rights reserved.
// Use of this source code is governed by a xxx
// license that can be found in the LICENSE file.

// Package handler provides HTTP handlers for the personal blog.
package handler

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"icode.baidu.com/baidu/personal-code/fullStack/model"
	"icode.baidu.com/baidu/personal-code/fullStack/service"
)

// Profile holds the blog owner's personal information displayed on the homepage.
type Profile struct {
	Name    string
	Bio     string
	GitHub  string
	Twitter string
}

// IndexData is the template data passed to templates/index.html.
type IndexData struct {
	Profile   Profile
	Posts     []*model.Post
	Tags      []string // de-duplicated, sorted list of all tags
	ActiveTag string   // currently filtered tag; empty means "all"
}

// PostData is the template data passed to templates/post.html.
// It embeds *model.Post so all Post fields are accessible directly in the template.
type PostData struct {
	*model.Post
	Content template.HTML // safe-typed HTML for {{.Content}}
}

// ErrorData is the template data passed to templates/error.html.
type ErrorData struct {
	Code    int
	Message string
}

// defaultProfile is the hard-coded blog owner profile.
var defaultProfile = Profile{
	Name:    "fh97",
	Bio:     "开发小菜鸡，喜欢编程和折腾工具",
	GitHub:  "https://github.com/fh97",
	Twitter: "",
}

// renderError renders the error template and writes the given HTTP status code.
// If template rendering itself fails, it falls back to a plain-text response.
func renderError(w http.ResponseWriter, templatesDir string, code int, message string) {
	tmpl, err := template.ParseFiles(
		filepath.Join(templatesDir, "base.html"),
		filepath.Join(templatesDir, "error.html"),
	)
	if err != nil {
		http.Error(w, message, code)
		return
	}
	w.WriteHeader(code)
	if err := tmpl.ExecuteTemplate(w, "base.html", ErrorData{Code: code, Message: message}); err != nil {
		log.Printf("renderError: execute base.html: %v", err)
	}
}

// collectTags returns a de-duplicated, order-stable list of all tags across posts.
func collectTags(posts []*model.Post) []string {
	seen := make(map[string]struct{})
	var tags []string
	for _, p := range posts {
		for _, t := range p.Tags {
			if _, ok := seen[t]; !ok {
				seen[t] = struct{}{}
				tags = append(tags, t)
			}
		}
	}
	return tags
}

// isPostNotFound reports whether the error from service.GetPostBySlug indicates
// the post does not exist (as opposed to a filesystem or parse error).
func isPostNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}

// postsDir resolves the posts directory path. It falls back to the working
// directory's "posts" sub-folder when postsDir is empty.
func resolvePostsDir(postsDir string) string {
	if postsDir != "" {
		return postsDir
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, "posts")
}

// IndexHandler returns an http.HandlerFunc that handles GET /.
//
// It reads an optional ?tag= query parameter to filter posts, collects all
// unique tags for the tag-filter bar, and renders templates/index.html.
//
// Parameters:
//   - postsDir: directory containing .md post files.
//   - templatesDir: directory containing *.html template files.
func IndexHandler(postsDir, templatesDir string) http.HandlerFunc {
	postsDir = resolvePostsDir(postsDir)

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(
			filepath.Join(templatesDir, "base.html"),
			filepath.Join(templatesDir, "index.html"),
		)
		if err != nil {
			log.Printf("IndexHandler: parse templates: %v", err)
			renderError(w, templatesDir, http.StatusInternalServerError, "模板加载失败")
			return
		}

		activeTag := r.URL.Query().Get("tag")

		// Fetch all posts (needed for tag collection regardless of filter).
		allPosts, err := service.GetAllPosts(postsDir)
		if err != nil {
			log.Printf("IndexHandler: GetAllPosts: %v", err)
			renderError(w, templatesDir, http.StatusInternalServerError, "文章加载失败")
			return
		}

		tags := collectTags(allPosts)

		// Apply tag filter if requested.
		var posts []*model.Post
		if activeTag != "" {
			posts, err = service.GetPostsByTag(postsDir, activeTag)
			if err != nil {
				log.Printf("IndexHandler: GetPostsByTag(%q): %v", activeTag, err)
				renderError(w, templatesDir, http.StatusInternalServerError, "文章筛选失败")
				return
			}
		} else {
			posts = allPosts
		}

		data := IndexData{
			Profile:   defaultProfile,
			Posts:     posts,
			Tags:      tags,
			ActiveTag: activeTag,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
			log.Printf("IndexHandler: execute base.html: %v", err)
		}
	}
}

// PostHandler returns an http.HandlerFunc that handles GET /post/{slug}.
//
// It extracts the slug from the URL path, loads and renders the corresponding
// Markdown file, then renders templates/post.html.
//
// Parameters:
//   - postsDir: directory containing .md post files.
//   - templatesDir: directory containing *.html template files.
func PostHandler(postsDir, templatesDir string) http.HandlerFunc {
	postsDir = resolvePostsDir(postsDir)

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse templates on every request (hot-reload for development).
		tmpl, err := template.ParseFiles(
			filepath.Join(templatesDir, "base.html"),
			filepath.Join(templatesDir, "post.html"),
		)
		if err != nil {
			log.Printf("PostHandler: parse templates: %v", err)
			renderError(w, templatesDir, http.StatusInternalServerError, "模板加载失败")
			return
		}

		// Extract slug from path: /post/<slug>
		// Support both net/http ServeMux and custom router patterns.
		slug := strings.TrimPrefix(r.URL.Path, "/post/")
		slug = strings.TrimSuffix(slug, "/") // tolerate trailing slash
		slug = strings.TrimSpace(slug)

		if slug == "" {
			renderError(w, templatesDir, http.StatusNotFound, "请求的文章不存在")
			return
		}

		post, err := service.GetPostBySlug(postsDir, slug)
		if err != nil {
			if isPostNotFound(err) {
				log.Printf("PostHandler: post not found: slug=%q", slug)
				renderError(w, templatesDir, http.StatusNotFound, "文章「"+slug+"」不存在")
				return
			}
			log.Printf("PostHandler: GetPostBySlug(%q): %v", slug, err)
			renderError(w, templatesDir, http.StatusInternalServerError, "文章加载失败")
			return
		}

		// Wrap the raw HTML string in template.HTML so the template engine
		// does not escape it (equivalent to the |safe filter in other engines).
		data := PostData{
			Post:    post,
			Content: template.HTML(post.Content), //nolint:gosec // content is rendered from trusted Markdown files
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
			log.Printf("PostHandler: execute base.html: %v", err)
		}
	}
}
