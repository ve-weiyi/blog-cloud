#!/bin/bash

PROTO_DIR="."                 # Proto 文件所在目录
OUTPUT_DIR="../internal/pb"   # 生成代码的输出目录
ETC_DIR="../etc"              # YAML 配置文件目录

shopt -s nullglob             # 当没有匹配的文件时，* 不再保留

for file in "$PROTO_DIR"/*.proto; do
  if [ -f "$file" ]; then
    # 生成代码
    goctl rpc protoc "$file" --go_out="$OUTPUT_DIR" --go-grpc_out="$OUTPUT_DIR" --zrpc_out=../ -m --style go_zero

    # 删除 YAML 配置文件
    rm -f "$ETC_DIR/$(basename "$file" .proto).yaml"

    # 删除生成的 Go 文件
    rm -f "../$(basename "$file" .proto).go"
  fi
done

#goctl rpc protoc ./article.proto --go_out=../internal/pb --go-grpc_out=../internal/pb --zrpc_out=../ -m --style go_zero
#rm -f ../etc/article.yaml
#rm ../article.go
