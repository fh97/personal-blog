---
title: Git 核心知识完全指南
date: 2026-04-10
tags: [Git, 工程效能, 版本控制]
summary: 从基础命令到高级技巧，涵盖分支管理、冲突解决、rebase、cherry-pick、stash 等核心知识点，附带实际工作中最常用的场景示例。
---

## 一、Git 基础概念

Git 是分布式版本控制系统，理解它的三个区域是掌握 Git 的关键：

```
工作区（Working Directory）
    ↓ git add
暂存区（Staging Area / Index）
    ↓ git commit
本地仓库（Local Repository）
    ↓ git push
远程仓库（Remote Repository）
```

---

## 二、最常用的日常命令

### 查看状态

```bash
git status          # 查看工作区和暂存区状态
git log --oneline   # 简洁的提交历史
git log --graph --oneline --all  # 图形化分支历史
git diff            # 查看未暂存的变更
git diff --staged   # 查看已暂存的变更
```

### 提交变更

```bash
git add .                    # 暂存所有变更
git add src/foo.go           # 暂存指定文件
git commit -m "feat: xxx"    # 提交
git commit --amend           # 修改最近一次提交（未 push 时）
```

### 撤销操作

```bash
git restore foo.go           # 丢弃工作区修改
git restore --staged foo.go  # 从暂存区移出（保留工作区修改）
git reset HEAD~1             # 撤销最近一次 commit（保留改动到工作区）
git reset --hard HEAD~1      # 撤销最近一次 commit 并丢弃改动（危险！）
```

---

## 三、分支管理

### 基本操作

```bash
git branch feature/login      # 创建分支
git switch feature/login      # 切换分支（推荐，替代 checkout）
git switch -c feature/login   # 创建并切换
git branch -d feature/login   # 删除已合并的分支
git branch -D feature/login   # 强制删除分支
```

### 合并分支

```bash
# 方式一：merge（保留分支历史，产生 merge commit）
git switch main
git merge feature/login

# 方式二：rebase（线性历史，更干净）
git switch feature/login
git rebase main
```

**何时用 merge，何时用 rebase？**

| 场景 | 推荐方式 |
|------|---------|
| 合并 feature 到 main | merge（保留完整历史） |
| 同步 main 的最新变更到 feature | rebase（保持线性） |
| 已 push 的分支 | **不要 rebase**（会改写历史） |

---

## 四、rebase 详解

rebase 是 Git 中最强大也最容易踩坑的命令。

### 基础 rebase

```bash
# 将 feature 分支的基点移到 main 的最新提交
git switch feature/login
git rebase main
```

执行前：
```
main:    A - B - C
feature:     D - E
```

执行后：
```
main:    A - B - C
feature:         D' - E'
```

### 交互式 rebase（整理提交历史）

```bash
git rebase -i HEAD~3  # 对最近 3 个提交进行交互式编辑
```

常用操作：
- `pick`：保留提交
- `squash`（s）：合并到上一个提交
- `reword`（r）：修改提交信息
- `drop`（d）：删除这个提交

**实战场景**：提 PR 前把 "fix typo"、"wip" 等杂乱提交合并成一个干净的提交。

---

## 五、cherry-pick

从其他分支摘取特定提交应用到当前分支。

```bash
# 把 commit abc1234 应用到当前分支
git cherry-pick abc1234

# 摘取多个提交
git cherry-pick abc1234 def5678

# 摘取一个范围（不包含起点）
git cherry-pick abc1234..def5678
```

**典型场景**：hotfix 分支修复了一个 bug，需要同步到 release 分支，但不想合并整个 hotfix 分支。

---

## 六、stash 暂存工作现场

临时保存未完成的工作，切换到其他任务。

```bash
git stash                    # 暂存当前工作
git stash push -m "wip: 登录功能"  # 带描述的暂存
git stash list               # 查看所有暂存
git stash pop                # 恢复最近一次暂存（并删除）
git stash apply stash@{1}    # 恢复指定暂存（不删除）
git stash drop stash@{0}     # 删除指定暂存
git stash clear              # 清空所有暂存
```

---

## 七、解决冲突

冲突文件格式：

```
<<<<<<< HEAD（当前分支的内容）
func hello() string {
    return "hello"
}
=======（分隔线）
func hello() string {
    return "Hello, World!"
}
>>>>>>> feature/greeting（合入分支的内容）
```

解决步骤：
1. 手动编辑冲突文件，保留正确内容
2. `git add <冲突文件>`
3. `git commit` 或 `git rebase --continue`

**中途放弃**：
```bash
git merge --abort    # 放弃 merge
git rebase --abort   # 放弃 rebase
```

---

## 八、远程仓库操作

```bash
git remote -v                        # 查看远程仓库
git fetch origin                     # 拉取远程变更（不合并）
git pull origin main                 # 拉取并合并
git pull --rebase origin main        # 拉取并 rebase（推荐）
git push origin feature/login        # 推送分支
git push -f origin feature/login     # 强制推送（rebase 后需要）
git push origin :feature/login       # 删除远程分支
```

---

## 九、实用技巧

### 查找问题提交（二分法）

```bash
git bisect start
git bisect bad           # 标记当前版本有问题
git bisect good v1.0.0   # 标记某个版本正常
# Git 自动切换到中间提交，测试后继续标记
git bisect good / git bisect bad
git bisect reset         # 结束查找
```

### 找回误删的提交

```bash
git reflog               # 查看所有操作历史（包含已删除的）
git checkout abc1234     # 恢复到指定状态
```

### 忽略文件

```bash
# .gitignore 常用规则
*.log          # 忽略所有 .log 文件
/dist          # 忽略根目录的 dist 文件夹
!important.log # 不忽略这个文件
```

---

## 十、commit message 规范

推荐使用 **Conventional Commits** 格式：

```
<type>(<scope>): <subject>

[可选 body]

[可选 footer]
```

常用 type：

| type | 含义 |
|------|------|
| `feat` | 新功能 |
| `fix` | 修复 bug |
| `docs` | 文档变更 |
| `refactor` | 重构（不影响功能） |
| `test` | 测试相关 |
| `chore` | 构建/工具链变更 |

示例：
```
feat(auth): 添加 JWT 登录接口

支持用户名+密码登录，返回 access_token 和 refresh_token。
token 有效期分别为 2h 和 7d。

closes #DEV-1234
```
