# Git 仓库同步工具

同步内部 Git 仓库到 GitHub，修改提交人信息为团队统一身份。

## 脚本说明

### sync_to_github_v3.sh

将当前仓库同步到 GitHub，自动修改所有提交的作者信息为团队统一身份。

**主要功能：**
- 重写 Git 历史，修改所有提交的作者信息
- 自动排除指定文件夹（不同步到 GitHub）
- 支持 dry-run 预览模式
- 本地分支保持不变，仅在临时分支上操作

## 使用方法

```bash
# 进入项目根目录
cd /path/to/terraform-provider-edgenext

# 先预览将要执行的操作（推荐）
./bin/sync_to_github_v3.sh --dry-run

# 执行同步
./bin/sync_to_github_v3.sh

# 查看帮助
./bin/sync_to_github_v3.sh --help
```

## 命令行选项

| 选项 | 说明 |
|------|------|
| `-d, --dry-run` | 仅预览，不执行实际操作 |
| `-h, --help` | 显示帮助信息 |

## 工作流程

1. **检查环境** - 验证 Git 仓库和远程配置
2. **创建临时分支** - 在临时分支上操作，不影响本地分支
3. **重写作者信息** - 使用 `git filter-branch` 修改所有提交的作者
4. **排除文件夹** - 从 GitHub 分支移除指定的排除文件夹
5. **更新 .gitignore** - 将排除目录添加到 .gitignore
6. **推送到 GitHub** - 强制推送到 GitHub 远程仓库
7. **切回原分支** - 本地分支保持不变

## 配置说明

可在脚本开头修改以下配置：

```bash
# 目标分支
GITHUB_BRANCH="main"

# 临时分支名称
TMP_BRANCH="github-sync"

# 团队作者信息
NEW_NAME="EdgeNext Team"
NEW_EMAIL="team@edgenext.com"

# 排除的文件夹（不会同步到 GitHub）
EXCLUDE_DIRS="openspec .cline .clinerules bin"
```

## 排除文件夹

以下文件夹不会同步到 GitHub：
- `openspec/` - OpenSpec 配置和变更记录
- `.cline/` - Cline 技能配置
- `.clinerules/` - Cline 工作流规则
- `bin/` - 内部脚本工具

## 注意事项

⚠️ **重要提醒：**

1. **本地分支安全** - 脚本仅在临时分支上操作，本地分支不会被修改
2. **强制推送** - 推送到 GitHub 时会使用 `--force`，覆盖远端历史
3. **团队协调** - 执行前请确保团队成员知晓
4. **首次运行** - 建议先使用 `--dry-run` 预览

## 前置条件

- 当前目录为 Git 仓库根目录
- 已配置 GitHub 远程仓库（或脚本会自动添加）
- 有推送到 GitHub 的权限

## 示例输出

```
🔍 检查当前目录...
📦 当前分支：master
🌐 检测到上游远程：origin
🔄 上游是 origin，将使用 github 远程
🌐 已检测到远程：https://github.com/edgenextapisdk/terraform-provider-edgenext.git

即将执行以下操作：
-------------------------------------
当前分支：    master
临时分支：    github-sync
同步目标：    github/main
新作者名：    EdgeNext Team
新作者邮箱：  team@edgenext.com
-------------------------------------

⚠️  警告：这将重写 Git 历史，修改所有提交的作者信息！
💡 本地 master 分支不会被修改，仅在临时分支 github-sync 上操作。
```

## 故障排查

### 问题：提示"当前目录不是 Git 仓库"

**解决方案：** 确保在项目根目录执行脚本

```bash
cd /path/to/terraform-provider-edgenext
./bin/sync_to_github_v3.sh
```

### 问题：推送失败

**解决方案：** 检查网络连接和 GitHub 权限

```bash
# 测试连接
git fetch github

# 检查远程配置
git remote -v
```

### 问题：临时分支残留

**解决方案：** 手动删除临时分支

```bash
git branch -D github-sync
```

## 版本历史

- **v3** - 添加 dry-run 模式、帮助信息、排除文件夹功能
- **v2** - 使用 git filter-branch 方法，稳定版本
- **v1** - 初始版本，使用 git filter-repo