#!/bin/bash
# =========================================
# 同步当前仓库到 GitHub，修改提交人信息
# v3 版本 - 基于 v2 小改
# 使用 git filter-branch 方法（稳定可靠）
# 
# 改进：
# - 添加 --dry-run 预览模式
# - 添加帮助信息
# - 显示提交总数
# =========================================

set -e

# --- 可自定义参数 ---
GITHUB_BRANCH="main"
TMP_BRANCH="github-sync"
NEW_NAME="EdgeNext Team"
NEW_EMAIL="team@edgenext.com"

# 排除的文件夹（不会同步到 GitHub）
EXCLUDE_DIRS="openspec .cline .clinerules bin"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 解析参数
DRY_RUN=false
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            echo "用法: $0 [选项]"
            echo ""
            echo "选项:"
            echo "  -d, --dry-run    仅预览，不执行实际操作"
            echo "  -h, --help       显示此帮助信息"
            echo ""
            echo "示例:"
            echo "  $0               # 执行同步"
            echo "  $0 --dry-run     # 仅预览"
            exit 0
            ;;
        *)
            echo -e "${RED}未知选项: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🔍 检查当前目录...${NC}"
if [ ! -d ".git" ]; then
  echo -e "${RED}❌ 当前目录不是 Git 仓库，请在项目根目录执行。${NC}"
  exit 1
fi

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo -e "${BLUE}📦 当前分支：${CURRENT_BRANCH}${NC}"

# 自动检测当前分支的上游远程
UPSTREAM=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null || echo "")
if [ -n "$UPSTREAM" ]; then
  UPSTREAM_REMOTE=$(echo "$UPSTREAM" | cut -d'/' -f1)
  echo -e "${BLUE}🌐 检测到上游远程：${UPSTREAM_REMOTE}${NC}"
  
  if [ "$UPSTREAM_REMOTE" = "origin" ]; then
    GITHUB_REMOTE="github"
    echo -e "${BLUE}🔄 上游是 origin，将使用 github 远程${NC}"
  else
    GITHUB_REMOTE="$UPSTREAM_REMOTE"
  fi
else
  GITHUB_REMOTE="github"
  echo -e "${YELLOW}⚠️  当前分支未设置上游，使用默认远程：${GITHUB_REMOTE}${NC}"
fi

# 检查远程是否存在
if ! git remote | grep -q "^${GITHUB_REMOTE}$"; then
  echo -e "${YELLOW}⚠️ 未检测到远程 '${GITHUB_REMOTE}'。${NC}"
  
  if [ "$DRY_RUN" = true ]; then
    echo -e "${BLUE}[DRY-RUN] 将添加 GitHub 远程仓库...${NC}"
  else
    echo -e "${BLUE}👉 添加 GitHub 远程仓库...${NC}"
    git remote add "${GITHUB_REMOTE}" "https://github.com/edgenextapisdk/terraform-provider-edgenext.git"
    echo -e "${GREEN}✅ 已添加远程：${GITHUB_REMOTE}${NC}"
  fi
else
  REPO_URL=$(git remote get-url ${GITHUB_REMOTE})
  echo -e "${BLUE}🌐 已检测到远程：${REPO_URL}${NC}"
fi

echo ""
echo -e "${YELLOW}即将执行以下操作：${NC}"
echo "-------------------------------------"
echo "当前分支：    $CURRENT_BRANCH"
echo "临时分支：    $TMP_BRANCH"
echo "同步目标：    ${GITHUB_REMOTE}/${GITHUB_BRANCH}"
echo "新作者名：    ${NEW_NAME}"
echo "新作者邮箱：  ${NEW_EMAIL}"
if [ "$DRY_RUN" = true ]; then
  echo "模式：        DRY-RUN（仅预览）"
fi
echo "-------------------------------------"
echo ""
echo -e "${RED}⚠️  警告：这将重写 Git 历史，修改所有提交的作者信息！${NC}"
echo -e "${YELLOW}💡 本地 ${CURRENT_BRANCH} 分支不会被修改，仅在临时分支 ${TMP_BRANCH} 上操作。${NC}"
echo ""

# DRY-RUN 模式确认
if [ "$DRY_RUN" = true ]; then
  echo -e "${BLUE}🔍 DRY-RUN 模式：仅显示将要执行的操作，不进行实际修改${NC}"
  echo ""
fi

read -p "❓是否确认继续？(y/N): " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
  echo -e "${RED}🚫 已取消操作。${NC}"
  exit 0
fi

echo ""
echo -e "${BLUE}🔍 检查远端状态...${NC}"

# 获取远端最新信息
if [ "$DRY_RUN" != true ]; then
  git fetch "${GITHUB_REMOTE}" "${GITHUB_BRANCH}" 2>/dev/null || true
fi

# 检查远端分支是否存在
REMOTE_REF="${GITHUB_REMOTE}/${GITHUB_BRANCH}"
if git rev-parse --verify "${REMOTE_REF}" >/dev/null 2>&1; then
  REMOTE_HEAD=$(git rev-parse "${REMOTE_REF}")
  LOCAL_HEAD=$(git rev-parse HEAD)
  
  if [ "$REMOTE_HEAD" != "$LOCAL_HEAD" ]; then
    echo -e "${YELLOW}⚠️  警告：远端分支 ${REMOTE_REF} 与本地不同！${NC}"
    echo -e "${YELLOW}   强制推送将覆盖远端历史。${NC}"
    echo ""
    read -p "❓ 是否继续（将覆盖远端）？(y/N): " continue_confirm
    if [[ ! "$continue_confirm" =~ ^[Yy]$ ]]; then
      echo -e "${RED}🚫 已取消操作。${NC}"
      exit 0
    fi
  fi
else
  echo -e "${BLUE}ℹ️  远端分支 ${REMOTE_REF} 不存在，将创建新分支。${NC}"
fi

# DRY-RUN 模式结束
if [ "$DRY_RUN" = true ]; then
  echo ""
  echo -e "${GREEN}✅ DRY-RUN 完成！以上是将要执行的操作。${NC}"
  echo -e "${BLUE}💡 移除 --dry-run 参数以执行实际操作。${NC}"
  exit 0
fi

echo ""
echo -e "${BLUE}🚀 开始重写作者信息...${NC}"

# 删除旧的临时分支（如果存在）
git branch -D "$TMP_BRANCH" 2>/dev/null || true

# 创建临时分支
echo -e "${BLUE}📦 创建临时分支：${TMP_BRANCH}${NC}"
git checkout -b "$TMP_BRANCH"

# 显示提交总数
TOTAL_COMMITS=$(git rev-list --count HEAD)
echo -e "${BLUE}📊 共 ${TOTAL_COMMITS} 个提交需要处理${NC}"

# 使用 git filter-branch 重写作者信息
echo -e "${BLUE}🔄 正在重写提交历史...${NC}"
echo -e "${YELLOW}   将所有提交的作者改为：${NEW_NAME} <${NEW_EMAIL}>${NC}"

# 设置环境变量禁用警告
export FILTER_BRANCH_SQUELCH_WARNING=1

# 构建 index-filter 命令移除排除的文件夹
INDEX_FILTER_CMD=""
for dir in $EXCLUDE_DIRS; do
  INDEX_FILTER_CMD="${INDEX_FILTER_CMD}git rm -rf --cached --ignore-unmatch ${dir} 2>/dev/null || true; "
done

# 执行 filter-branch（同时重写作者信息和移除排除文件夹）
git filter-branch -f --env-filter "
  export GIT_AUTHOR_NAME='${NEW_NAME}'
  export GIT_AUTHOR_EMAIL='${NEW_EMAIL}'
  export GIT_COMMITTER_NAME='${NEW_NAME}'
  export GIT_COMMITTER_EMAIL='${NEW_EMAIL}'
" --index-filter "${INDEX_FILTER_CMD}" HEAD 2>&1 | grep -v "WARNING:" | grep -v "git-filter-branch" || true

# 清理 filter-branch 的备份
rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now --quiet

echo ""
echo -e "${GREEN}✅ 重写完成！${NC}"
echo -e "${BLUE}   已移除排除文件夹并统一作者信息${NC}"

echo ""
echo -e "${BLUE}📋 查看修改后的提交：${NC}"
git log --oneline --format="%h %an <%ae> %s" -5

echo ""
echo -e "${BLUE}🔍 对比本地和临时分支：${NC}"
echo ""
echo -e "${YELLOW}=== 本地 ${CURRENT_BRANCH} 分支（保留真实信息）===${NC}"
git log ${CURRENT_BRANCH} --oneline --format="%h %an <%ae> %s" -3
echo ""
echo -e "${YELLOW}=== 临时 ${TMP_BRANCH} 分支（团队信息）===${NC}"
git log ${TMP_BRANCH} --oneline --format="%h %an <%ae> %s" -3

echo ""
read -p "⚠️ 是否推送到 GitHub (${GITHUB_REMOTE}/${GITHUB_BRANCH})？(y/N): " push_confirm
if [[ ! "$push_confirm" =~ ^[Yy]$ ]]; then
  echo -e "${RED}🚫 已取消推送。${NC}"
  echo -e "${BLUE}💡 临时分支 ${TMP_BRANCH} 已保留，如需删除：git branch -D ${TMP_BRANCH}${NC}"
  git checkout "$CURRENT_BRANCH"
  exit 0
fi

echo -e "${BLUE}📤 正在推送到 ${GITHUB_REMOTE}/${GITHUB_BRANCH}...${NC}"
git push -f "${GITHUB_REMOTE}" "${TMP_BRANCH}:${GITHUB_BRANCH}"

# 切回原分支
git checkout "$CURRENT_BRANCH"

echo ""
echo -e "${GREEN}✅ 推送完成！${NC}"
echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}📦 分支：${GITHUB_REMOTE}/${GITHUB_BRANCH}${NC}"
echo -e "${GREEN}👤 作者：${NEW_NAME} <${NEW_EMAIL}>${NC}"
echo -e "${GREEN}🌐 仓库：$(git remote get-url ${GITHUB_REMOTE})${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}💡 本地分支 ${CURRENT_BRANCH} 保持不变（真实作者信息）${NC}"
echo -e "${BLUE}💡 临时分支 ${TMP_BRANCH} 已保留，如需删除：git branch -D ${TMP_BRANCH}${NC}"
echo ""
echo -e "${YELLOW}🔍 验证方式：${NC}"
echo -e "   1. 查看本地原始分支：git log ${CURRENT_BRANCH} --oneline --format=\"%h %an <%ae> %s\" -5"
echo -e "   2. 查看 GitHub 分支：git log ${GITHUB_REMOTE}/${GITHUB_BRANCH} --oneline --format=\"%h %an <%ae> %s\" -5"
echo -e "   3. 访问 GitHub 网页：$(git remote get-url ${GITHUB_REMOTE} | sed 's/\.git$//')"
echo ""