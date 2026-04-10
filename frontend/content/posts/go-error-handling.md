---
title: 聊聊 Go 的错误处理哲学
date: 2026-03-20
tags: [Go, 最佳实践]
summary: Go 的错误处理是显式的、啰嗦的，但也是清晰的。这篇文章聊聊我对 Go 错误处理的一些理解和实践心得。
---

## 为什么 Go 不用异常？

Go 的设计者认为异常会让控制流变得难以追踪。显式的错误返回值，虽然啰嗦，但每一处错误处理都是可见的。

## 几个常用模式

### Sentinel Error

```go
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {
    // handle
}
```

### Error Wrapping

```go
return nil, fmt.Errorf("getPost: %w", err)
```

### 自定义错误类型

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s - %s", e.Field, e.Message)
}
```

## 总结

不要怕写 `if err != nil`，清晰比简洁更重要。
