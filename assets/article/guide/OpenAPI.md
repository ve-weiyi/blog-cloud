# OpenAPI 规范指南

## 什么是 OpenAPI

OpenAPI 规范（原 Swagger 规范）是一种用于描述 RESTful API 的标准格式。它使用 JSON 或 YAML 格式定义 API 的结构、端点、参数、响应等信息。

## OpenAPI vs RESTful API

| 对比项 | OpenAPI          | RESTful API    |
|-----|------------------|----------------|
| 定义  | API 描述规范         | API 设计风格       |
| 作用  | 定义 API 文档格式      | 定义 API 设计原则    |
| 关系  | 可以描述 RESTful API | 可以用 OpenAPI 描述 |

## OpenAPI 版本

### OpenAPI 2.0（Swagger 2.0）

- 发布时间：2014 年
- 格式：JSON/YAML
- 特点：广泛使用，工具支持完善

### OpenAPI 3.0

- 发布时间：2017 年
- 格式：JSON/YAML
- 改进：更灵活的结构、更好的复用性、支持多服务器

**主要区别**

| 特性    | OpenAPI 2.0           | OpenAPI 3.0                  |
|-------|-----------------------|------------------------------|
| 服务器定义 | `host` + `basePath`   | `servers` 数组                 |
| 请求体   | `parameters`          | `requestBody`                |
| 响应定义  | 简单                    | 更详细，支持多媒体类型                  |
| 组件复用  | `definitions`         | `components`                 |
| 安全定义  | `securityDefinitions` | `components/securitySchemes` |

## OpenAPI 3.0 基本结构

```yaml
openapi: 3.0.0
info:
  title: 博客 API
  version: 1.0.0
  description: 博客系统接口文档

servers:
  - url: https://api.example.com/v1
    description: 生产环境
  - url: https://test-api.example.com/v1
    description: 测试环境

paths:
  /users:
    get:
      summary: 获取用户列表
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/User'

components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
```

## 核心概念

### 1. Info 对象

定义 API 的基本信息

```yaml
info:
  title: API 标题
  version: 1.0.0
  description: API 描述
  contact:
    name: 联系人
    email: contact@example.com
```

### 2. Servers 对象

定义 API 服务器地址

```yaml
servers:
  - url: https://api.example.com/v1
    description: 生产环境
  - url: http://localhost:8080
    description: 本地开发
```

### 3. Paths 对象

定义 API 端点

```yaml
paths:
  /users/{id}:
    get:
      summary: 获取用户信息
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: 成功
```

### 4. Components 对象

定义可复用的组件

```yaml
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string

  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

## 常用工具

### 1. Swagger Editor

在线编辑和验证 OpenAPI 文档

- 地址：https://editor.swagger.io/
- 功能：实时预览、语法检查、代码生成

### 2. Swagger UI

生成交互式 API 文档

```bash
# Docker 运行
docker run -p 80:8080 -e SWAGGER_JSON=/foo/openapi.yaml -v /path/to/openapi.yaml:/foo/openapi.yaml swaggerapi/swagger-ui
```

### 3. Swagger Codegen

根据 OpenAPI 文档生成客户端和服务端代码

```bash
# 生成 Go 服务端代码
swagger-codegen generate -i openapi.yaml -l go-server -o ./server
```

## 最佳实践

1. **版本管理**：使用语义化版本号
2. **描述清晰**：为每个端点添加详细的 `summary` 和 `description`
3. **使用组件**：通过 `$ref` 引用复用组件
4. **错误处理**：定义统一的错误响应格式
5. **安全定义**：明确定义认证和授权方式
6. **示例数据**：为请求和响应添加示例

## 参考资料

- [OpenAPI 官方规范](https://spec.openapis.org/oas/latest.html)
- [OpenAPI 中文文档](https://openapi.apifox.cn/)
- [Swagger 官方网站](https://swagger.io/)
- [OpenAPI 2.0 vs 3.0 对比](https://blog.readme.com/an-example-filled-guide-to-swagger-3-2/)
