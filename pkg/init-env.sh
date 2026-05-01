#!/bin/bash
set -euo pipefail

# ====================== 模板配置：支持 2 种写法 ======================
# 写法1：远程URL  →  自动下载覆盖
# 写法2：本地文件路径  →  自动复制覆盖
# 可同时配置多个，支持混用
TEMPLATE_MAP=(
  # 远程URL示例
  # 本地文件示例（把本地文件覆盖到goctl模板目录）
   "api/handler.tpl=./pkg/goctl/template/api/handler.tpl"
  # "model/model.tpl=./local-tpl/model/model.tpl"
)

# -------------------------- 函数1：检查基础环境 --------------------------
check_environment() {
  echo "==================================================="
  echo " 正在检查基础环境..."
  echo "==================================================="

  if [ ! -f "go.mod" ]; then
    echo "错误：当前目录未找到 go.mod 文件"
    exit 1
  fi

  if ! command -v go &> /dev/null; then
    echo "错误：未安装 Go 环境"
    exit 1
  fi

  GO_ZERO_VERSION=$(grep 'github.com/zeromicro/go-zero' go.mod | awk '{print $2}' | sed 's/^v//')
  if [ -z "$GO_ZERO_VERSION" ]; then
    echo "错误：未在 go.mod 中找到 go-zero 依赖"
    exit 1
  fi

  GOCTL_TEMPLATE_DIR="$HOME/.goctl/$GO_ZERO_VERSION"

  echo "✅ 环境检查通过"
  echo "Go-zero 版本：v$GO_ZERO_VERSION"
  echo "goctl 模板目录：$GOCTL_TEMPLATE_DIR"
  echo
}

# -------------------------- 函数2：初始化goctl环境 --------------------------
init_goctl_environment() {
  echo "==================================================="
  echo " 正在初始化 go-zero 开发环境 v$GO_ZERO_VERSION"
  echo "==================================================="

  if ! command -v goctl &> /dev/null; then
    echo "正在安装 goctl..."
    go install github.com/zeromicro/go-zero/tools/goctl@v$GO_ZERO_VERSION
  else
    echo "goctl 已安装"
  fi

  echo "检查 goctl 环境并安装依赖..."
  if ! goctl env check --install --verbose --force; then
    echo "❌ goctl 环境检查失败"
    exit 1
  fi

  echo "初始化官方模板..."
  goctl template init

  echo "✅ goctl 环境初始化完成"
  echo
}

# -------------------------- 函数3：覆盖自定义模板（URL+本地双支持） --------------------------
override_custom_templates() {
  echo -e "\n==================================================="
  echo " 开始覆盖自定义模板（支持远程URL / 本地文件）"
  echo "==================================================="

  for item in "${TEMPLATE_MAP[@]}"; do
    IFS='=' read -r target_file source <<< "$item"
    dest="$GOCTL_TEMPLATE_DIR/$target_file"
    mkdir -p "$(dirname "$dest")"

    # ============== 核心逻辑：自动判断是URL还是本地文件 ==============
    if [[ "$source" =~ ^https?:// ]]; then
      # 处理远程URL
      echo "🌐 下载远程模板: $source"
      echo "└─→ 覆盖到: $dest"
      if ! curl -fsSL "$source" -o "$dest"; then
        echo "❌ 下载失败: $source"
        exit 1
      fi
    else
      # 处理本地文件
      if [ -f "$source" ]; then
        echo "📁 使用本地模板: $source"
        echo "└─→ 覆盖到: $dest"
        cp -f "$source" "$dest"
      else
        echo "❌ 错误：本地文件不存在 → $source"
        exit 1
      fi
    fi

    echo "✅ 成功处理: $target_file"
    echo
  done
}

# ========================== 主执行流程 ==========================
main() {
  check_environment
  init_goctl_environment
  override_custom_templates

  echo "==================================================="
  echo " ✅ 环境初始化全部完成！"
  echo "==================================================="
}

main