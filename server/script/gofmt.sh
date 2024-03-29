#!/bin/bash

format_directory() {
    for file in "$1"/*; do
        if [ -d "$file" ]; then
            format_directory "$file"
        elif [ "${file##*.}" = "go" ]; then
            go fmt "$file"
        fi
    done
}

format_directory .

# 递归格式化指定目录下的所有 Go 文件
format_go_files() {
    local dir="$1"

    # 遍历目录中的所有文件和子目录
    for file in "$dir"/*; do
        if [ -d "$file" ]; then
            # 如果是子目录，则递归处理
            format_go_files "$file"
        elif [ "${file##*.}" = "go" ]; then
            # 如果是 Go 文件，则使用 goimports 进行格式化
            goimports -w "$file"
        fi
    done
}

# 指定要格式化的目录
target_dir="."

# 调用函数来格式化目录下的所有 Go 文件
format_go_files "$target_dir"
