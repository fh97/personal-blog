// Copyright 2026 Baidu Inc. All rights reserved.
// Use of this source code is governed by a xxx
// license that can be found in the LICENSE file.

// Package main is the entry point for the personal blog server.
package main

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"icode.baidu.com/baidu/personal-code/fullStack/handler"
)

// main starts the blog HTTP server on :8080.
func main() {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)

	postsDir := filepath.Join(root, "content", "posts")
	templatesDir := filepath.Join(root, "templates")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.IndexHandler(postsDir, templatesDir))
	mux.HandleFunc("/post/", handler.PostHandler(postsDir, templatesDir))

	addr := ":8080"
	log.Printf("博客启动，访问 http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
