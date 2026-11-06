#!/bin/bash

# 检查手写 Go 代码的格式与 import 分组；生成代码与 runtime 产物不参与。
# import 分组约定：标准库 / 第三方 / 本仓（blog-cloud、vkit、stompws）三组。
#
# 用法：
#   ./lint.sh         只检查，有问题则失败
#   ./lint.sh --fix   就地修复

set -e

GOIMPORTS=${GOIMPORTS:-goimports}
GCI=${GCI:-gci}
LOCAL_PREFIX="github.com/ve-weiyi"

case "${1:-}" in
    "") fix=0 ;;
    --fix) fix=1 ;;
    *) echo "用法：$0 [--fix]" >&2; exit 2 ;;
esac

need_tool() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "未找到 $1，请先安装：$2" >&2
        exit 1
    fi
}

need_tool "$GOIMPORTS" "go install golang.org/x/tools/cmd/goimports@latest"
need_tool "$GCI" "go install github.com/daixiang0/gci@latest"

# goctl 生成的 internal/server 在包名与路径末段不一致时会被 import 工具补别名，
# 重新生成即被抹掉，故与 pb 一样排除。
list_files() {
    find . -name "*.go" \
        -not -path "./.git/*" \
        -not -path "./runtime/*" \
        -not -path "*/internal/pb/*" \
        -not -path "*/internal/server/*" \
        -not -name "*.pb.go" \
        -print0
}

# gci 递归传目录（`.`）在本机会卡在 kqueue 上，必须传显式文件列表。
gci_args=(--skip-generated -s standard -s default -s "prefix($LOCAL_PREFIX)")

if [ "$fix" -eq 1 ]; then
    list_files | xargs -0 "$GOIMPORTS" -w
    # gci 把每个改动文件的 INFO 日志写到 stderr，过滤掉以免淹没错误。
    list_files | xargs -0 "$GCI" write "${gci_args[@]}" 2> >(grep -v 'INFO' >&2)
    echo "已完成修复"
    exit 0
fi

failed=0

unformatted=$(list_files | xargs -0 gofmt -l)
if [ -n "$unformatted" ]; then
    echo "以下文件未格式化，请运行 ./lint.sh --fix：" >&2
    echo "$unformatted" >&2
    failed=1
fi

# gci diff 有差异时退出码为 1，不能直接沿用 set -e。
ungrouped=$(list_files | xargs -0 "$GCI" diff "${gci_args[@]}" || true)
if [ -n "$ungrouped" ]; then
    echo "以下文件 import 分组不规范，请运行 ./lint.sh --fix：" >&2
    echo "$ungrouped" | sed -n 's/^--- //p' >&2
    failed=1
fi

if [ "$failed" -ne 0 ]; then
    exit 1
fi

echo "格式检查通过"
