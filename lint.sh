#!/bin/bash

# import 格式化: 先去除 import 块内空行，再运行 gofmt

set -e

lint_file() {
    local file="$1"

    awk '
    /^import \(/ { in_import=1 }
    /^\)$/ && in_import { in_import=0 }
    {
        if (in_import && $0 ~ /^[[:space:]]*$/) next
        print
    }
    ' "$file" > "$file.tmp" && mv "$file.tmp" "$file"

    gofmt -w "$file"
}

get_files() {
    find . -name "*.go" -not -path "./.git/*"
}

while IFS= read -r file; do
    lint_file "$file"
done < <(get_files)

echo "import 格式化完成"
