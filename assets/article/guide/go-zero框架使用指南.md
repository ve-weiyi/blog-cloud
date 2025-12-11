# Go-Zero 框架使用指南

Go-Zero 是一个集成了各种工程实践的 Web 和 RPC 框架，内建了 goctl 工具大幅提升开发效率。

## 官方资源

- 官网：https://go-zero.dev/
- GitHub：https://github.com/zeromicro/go-zero
- 文档：https://go-zero.dev/docs/concepts/overview

## 环境安装

### 安装 Goctl

```bash
# 安装 Goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 验证安装
goctl --version
```

### 安装 Protoc（推荐一键安装）

```bash
# 一键安装 protoc 及相关工具
goctl env check --install --verbose --force

# 验证安装
goctl env check --verbose
```

### 初始化模板

```bash
# 初始化代码模板
goctl template init
```

## 常用命令

### 创建项目

```bash
# 创建 API 服务
goctl api new blog --style go_zero

# 创建 RPC 服务
goctl rpc new blog --style go_zero
```

### 生成代码

```bash
# 生成 API 代码
goctl api go -api blog.api -dir ./ --style go_zero

# 生成 RPC 代码
goctl rpc protoc blog.proto --zrpc_out=. --style go_zero

# 生成 Model 代码
goctl model mysql datasource \
  -url="root:password@tcp(127.0.0.1:3306)/database" \
  -table="t_*" \
  -dir="./model" \
  --style go_zero
```

### 其他命令

```bash
# 格式化 API 文件
goctl api format --dir blog.api

# 验证 API 文件
goctl api validate -api blog.api

# 生成 TypeScript 代码
goctl api ts -api blog.api -dir ../ts

# 生成 Dockerfile
goctl docker --go service/rpc/blog/blog.go --exe blog
```

## API 文件语法

```go
syntax = "v1"

info(
title: "博客 API"
version: "v1.0"
)

type LoginReq {
Username string `json:"username"`
Password string `json:"password"`
}

type LoginResp {
Token string `json:"token"`
}

@server(
prefix: /api/v1
group: auth
)
service blog-api {
@doc "用户登录"
@handler Login
post /login (LoginReq) returns (LoginResp)
}
```

**中间件使用**

```go
@server(
prefix: /api/v1
group: user
middleware: Auth // 使用中间件
)
```

**JWT 认证**

```go
@server(
prefix: /api/v1
group: user
jwt: Auth // 使用 JWT
)
```

## Proto 文件语法

```protobuf
syntax = "proto3";

package blog;
option go_package = "./blog";

message LoginReq {
  string username = 1;
  string password = 2;
}

message LoginResp {
  string token = 1;
}

service Blog {
  rpc Login(LoginReq) returns(LoginResp);
}
```

## 项目结构

**API 服务**

```
blog-api/
├── etc/              # 配置文件
├── internal/
│   ├── config/      # 配置定义
│   ├── handler/     # 路由处理器
│   ├── logic/       # 业务逻辑
│   ├── svc/         # 服务上下文
│   └── types/       # 类型定义
├── blog.api         # API 定义
└── blog.go          # 入口文件
```

**RPC 服务**

```
blog-rpc/
├── etc/              # 配置文件
├── internal/
│   ├── config/      # 配置定义
│   ├── logic/       # 业务逻辑
│   ├── server/      # gRPC 服务器
│   └── svc/         # 服务上下文
├── blog/            # protobuf 生成
├── blog.proto       # Proto 定义
└── blog.go          # 入口文件
```

## 常见问题

**Goctl 命令找不到**

```bash
# 添加到 PATH
export PATH=$PATH:$GOPATH/bin
```

**Protoc 版本不兼容**

```bash
# 使用一键安装
goctl env check --install --verbose --force
```

**生成代码报错**

```bash
# 验证 API 文件
goctl api validate -api blog.api

# 重新初始化模板
goctl template clean
goctl template init
```

## 参考资料

- [Go-Zero 官方文档](https://go-zero.dev/)
- [Go-Zero GitHub](https://github.com/zeromicro/go-zero)
- [Goctl 工具文档](https://go-zero.dev/docs/tasks/installation/goctl)
