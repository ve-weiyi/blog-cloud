#!/bin/bash

PROTO_DIR="."                 # Proto 文件所在目录
OUTPUT_DIR="../internal/pb"   # Proto 生成代码的输出目录 (设置在internal内部是为了防止api使用，api应该使用client下的pb)
ZRPC_OUT="../"                # zrpc_out 输出目录
ETC_DIR="../etc"              # YAML 配置文件目录

shopt -s nullglob             # 当没有匹配的文件时，* 不再保留

for file in "$PROTO_DIR"/*.proto; do
  if [ -f "$file" ]; then
    # 生成代码
    goctl rpc protoc "$file" \
      --go_out="$OUTPUT_DIR" \
      --go-grpc_out="$OUTPUT_DIR" \
      --zrpc_out="$ZRPC_OUT" \
      --style go_zero \
      -m

    # 删除 YAML 配置文件
    rm -f "$ETC_DIR/$(basename "$file" .proto).yaml"

    # 删除生成的 Go 文件
    rm -f "${ZRPC_OUT}$(basename "$file" .proto).go"
  fi
done

#goctl rpc protoc ./message.proto --go_out=../pb --go-grpc_out=../pb --zrpc_out=../ -m --style go_zero
#rm -f ../etc/message.yaml
#rm ../message.go
