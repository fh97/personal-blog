// Copyright 2026 Baidu Inc. All rights reserved.
// Use of this source code is governed by a xxx
// license that can be found in the LICENSE file.

// Package main is the entry point for the personal blog server.
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"icode.baidu.com/baidu/personal-code/fullStack/handler"
)

// main starts the blog HTTP server on :8300.
// It expects to be run from the repository root, where frontend/ lives.
func main() {
	root, _ := os.Getwd()

	postsDir := filepath.Join(root, "frontend", "content", "posts")
	templatesDir := filepath.Join(root, "frontend", "templates")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.IndexHandler(postsDir, templatesDir))
	mux.HandleFunc("/post/", handler.PostHandler(postsDir, templatesDir))

	addr := ":8300"
	log.Printf("博客启动，访问 http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
