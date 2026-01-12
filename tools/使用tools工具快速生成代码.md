# Tools - 代码生成工具

ve-blog-golang 项目的代码生成工具集，用于从 API 定义文件自动生成后端、前端和数据库模型代码，提高开发效率。

## 目录结构

```
tools/
├── cmd/                    # 命令行工具实现
│   ├── api/               # API 代码生成命令
│   │   └── gin/          # Gin 框架生成
│   ├── model/             # Model 代码生成命令
│   │   └── mysql/        # MySQL 模型生成
│   └── web/               # Web 代码生成命令
│       └── ts/           # TypeScript 生成
├── parserx/               # 解析器
│   ├── apiparser/         # API 文件解析器
│   ├── dbparser/          # 数据库解析器
│   └── swagparser/        # Swagger 解析器
├── template/              # 代码模板
│   ├── api/              # API 模板
│   │   ├── gin/         # Gin 框架模板
│   │   └── go-zero/     # Go-zero 框架模板
│   ├── model/            # Model 模板
│   └── web/              # Web 前端模板
│       └── ts/          # TypeScript 模板
├── testdata/              # 测试数据
├── Makefile               # 构建脚本
├── main.go                # 入口文件
└── go.mod                 # 依赖管理
```

## 功能概述

| 功能         | 说明                                  | 生成内容                       |
|------------|-------------------------------------|----------------------------|
| API 代码生成   | 从 `.api` 或 `swagger.json` 生成 Gin 代码 | types/logic/handler/router |
| Model 代码生成 | 从数据库或 SQL 文件生成 GORM 模型              | 数据库模型文件                    |
| Web 代码生成   | 从 `.api` 或 `swagger.json` 生成 TS 代码  | API 接口调用代码和类型定义            |


## 命令结构

```
tools
├── api                    # API 代码生成
│   └── gin               # Gin 框架
│       ├── api           # 从 .api 文件生成
│       └── swagger       # 从 swagger.json 生成
├── model                  # Model 代码生成
│   └── mysql             # MySQL 数据库
│       ├── ddl           # 从 SQL 文件生成
│       └── dsn           # 从数据库连接生成
└── web                    # Web 代码生成
    └── ts                # TypeScript
        ├── api           # 从 .api 文件生成
        └── swagger       # 从 swagger.json 生成
```

## 命令参数说明

| 参数   | 说明       | 示例                                                  |
|------|----------|-----------------------------------------------------|
| `-f` | 输入文件路径   | `../blog-gozero/service/api/blog/proto/blog.api`    |
| `-t` | 模板目录路径   | `./template/api/gin`                                |
| `-o` | 输出目录路径   | `../blog-gin/api/blog`                              |
| `-n` | 输出文件名格式  | `%s.go` / `%v_model.go`                             |
| `-c` | 上下文包路径   | `github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx` |
| `-s` | SQL 文件路径 | `../blog-veweiyi-init.sql`                          |
| `-u` | 数据库连接字符串 | `root:password@(host:3306)/database`                |

## 快速开始

### 使用 Makefile（推荐）

```bash
# 查看所有可用命令
make help

# 安装依赖
make deps

# 生成 Gin 框架代码（从 .api 文件）
make gen-api-gin-api

# 生成 Gin 框架代码（从 swagger.json 文件）
make gen-api-gin-swagger

# 生成数据库模型代码（从 SQL 文件）
make gen-model-ddl

# 生成数据库模型代码（从数据库连接）
make gen-model-dsn

# 生成 TypeScript 代码（从 .api 文件）
make gen-web-ts-api

# 生成 TypeScript 代码（从 swagger.json 文件）
make gen-web-ts-swagger

# 清理生成的代码
make clean
```

### 使用命令行

**生成 Gin API 代码**

```bash
# 从 .api 文件生成
go run main.go api gin api \
  -f ../blog-gozero/service/api/blog/proto/blog.api \
  -t ./template/api/gin \
  -o ../blog-gin/api/blog \
  -c github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx \
  -n '%s.go'

# 从 swagger.json 文件生成
go run main.go api gin swagger \
  -f ../blog-gozero/service/api/blog/docs/blog.json \
  -t ./template/api/gin \
  -o ../blog-gin/api/blog \
  -c github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx \
  -n '%s.go'
```

**生成 Model 代码**

```bash
# 从 SQL 文件生成
go run main.go model mysql ddl \
  -s ../blog-veweiyi-init.sql \
  -t ./template/model/model.tpl \
  -o ./runtime/model \
  -n '%v_model.go'

# 从数据库连接生成
go run main.go model mysql dsn \
  -u 'root:mysql7914@(127.0.0.1:3306)/blog-veweiyi?charset=utf8mb4&parseTime=True&loc=Local' \
  -t ./template/model/model.tpl \
  -o ./runtime/model \
  -n '%v_model.go'
```

**生成 TypeScript 代码**

```bash
# 从 .api 文件生成
go run main.go web ts api \
  -f ../blog-gozero/service/api/blog/proto/blog.api \
  -t ./template/web/ts \
  -o ./runtime/web/ts/blog/api \
  -n '%v.ts'

# 从 swagger.json 文件生成
go run main.go web ts swagger \
  -f ../blog-gozero/service/api/blog/docs/blog.json \
  -t ./template/web/ts \
  -o ./runtime/web/ts/blog/api \
  -n '%v.ts'
```

## 工作原理

### 解析器

**SpecParser** - 解析 `.api` 文件

- 使用 go-zero 的 API 解析器
- 提取接口路径、方法、参数、返回值
- 生成统一的 ApiSpec 结构

**SwaggerParser** - 解析 `swagger.json` 文件

- 使用 go-openapi 解析 Swagger 文档
- 转换为统一的 ApiSpec 结构
- 支持与 .api 文件相同的生成流程

**DBParser** - 解析数据库结构

- 支持从 SQL DDL 文件解析
- 支持从数据库 DSN 连接解析
- 提取表结构、字段、索引等信息

### 代码生成流程

1. 解析输入文件（`.api`、`swagger.json` 或数据库）
2. 转换为统一的数据结构
3. 根据模板生成目标代码
4. 输出到指定目录

## 技术特点

- ✅ 基于 Cobra 构建命令行工具
- ✅ 支持自定义模板，灵活扩展
- ✅ 支持多种代码生成场景
- ✅ 自动化生成，减少重复代码编写
- ✅ 统一代码风格，提高代码质量

## 使用场景

1. **新增接口** - 修改 `.api` 文件后，运行 `make gen-api-gin-api` 自动生成 Gin 代码
2. **数据库变更** - 修改 SQL 文件后，运行 `make gen-model-ddl` 更新模型
3. **前后端协作** - 运行 `make gen-web-ts-api` 生成前端 TypeScript 代码
4. **已有 Swagger** - 从现有 Swagger 文档生成代码，运行 `make gen-api-gin-swagger`

## 注意事项

1. 生成代码前建议先备份或提交现有代码
2. 生成的代码可能需要手动调整部分逻辑
3. 模板文件位于 `template/` 目录，可根据需求自定义
4. 生成的文件会覆盖同名文件，请谨慎操作
