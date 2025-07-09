# Quickstart

Quickstart 是一个基于 Go AST（抽象语法树）的代码生成和转换工具包，提供了简单易用的 API 来操作和生成 Go
代码。该工具包主要包含三个核心模块：AST 操作、代码注入和代码生成。

## 目录结构

```
quickstart/
├── astx/           # AST 操作核心包
│   ├── inject_ast.go    # AST 注入工具
│   ├── inject_func.go   # 函数注入工具
│   ├── inject_import.go # 导入注入工具
│   ├── inject_struct.go # 结构体注入工具
│   ├── utils.go         # 通用工具函数
│   └── easy_test.go     # 测试用例
├── inject/         # 代码注入包(使用dst库)
│   ├── inject_ast.go    # AST 注入实现
│   ├── inject_decl.go   # 声明注入实现
│   ├── inject_func.go   # 函数注入实现
│   ├── inject_import.go # 导入注入实现
│   ├── inject_struct.go # 结构体注入实现
│   ├── utils.go         # 工具函数
│   └── easy_test.go     # 测试用例
└── invent/         # 代码生成包
    ├── template.go      # 模板生成工具
    ├── file_utils.go    # 文件操作工具
    └── func_utils.go    # 函数生成工具
```

## 核心特性

### 1. AST 操作 (astx)

- 基于标准库 `go/ast` 的代码解析和转换
- 支持完整的 AST 节点操作
- 提供友好的 API 接口
- 适用于代码分析和转换场景

### 2. 代码注入 (inject)

- 基于 `dst` 库的代码注入实现
- 支持导入声明注入
- 支持结构体定义注入
- 支持函数声明注入
- 支持通用声明注入
- 适用于代码增强和修改场景

### 3. 代码生成 (invent)

- 基于模板的代码生成
- 提供文件操作工具
- 提供函数生成工具
- 适用于代码自动生成场景

## 快速开始

### 安装

```bash
go get github.com/ve-weiyi/ve-blog-golang/kit/quickstart
```

### 使用示例

#### 1. AST 操作示例

```go
import (
"github.com/ve-weiyi/ve-blog-golang/kit/quickstart/astx"
"go/parser"
"go/token"
)

// 创建文件集
fset := token.NewFileSet()

// 解析 Go 文件
f, err := parser.ParseFile(fset, "example.go", nil, parser.ParseComments)
if err != nil {
log.Fatal(err)
}

// 使用 AST 工具
astTool := astx.NewASTTool(fset, f)
```

#### 2. 代码注入示例

```go
import "github.com/ve-weiyi/ve-blog-golang/kit/quickstart/inject"

// 创建注入器
injector := inject.NewInjector()

// 注入导入
injector.AddImport("fmt")
injector.AddImport("strings")

// 注入结构体
injector.AddStruct("User", []inject.Field{
{Name: "Name", Type: "string", Tag: `json:"name"`},
{Name: "Age", Type: "int", Tag: `json:"age"`},
})

// 注入函数
injector.AddFunction("GetUser", []inject.Param{
{Name: "id", Type: "int"},
}, "User", []string{
"return User{Name: \"test\", Age: 18}",
})
```

#### 3. 代码生成示例

```go
import "github.com/ve-weiyi/ve-blog-golang/kit/quickstart/invent"

// 创建生成器
generator := invent.NewGenerator()

// 生成代码
err := generator.Generate("template.tmpl", "output.go", map[string]interface{}{
"PackageName": "example",
"StructName":  "User",
"Fields": []map[string]string{
{"Name": "Name", "Type": "string"},
{"Name": "Age", "Type": "int"},
},
})
```

## 最佳实践

1. **AST 操作**
    - 使用 `astx` 包进行代码分析和转换
    - 适用于需要精确控制代码结构的场景
    - 建议用于代码质量检查和重构

2. **代码注入**
    - 使用 `inject` 包进行代码注入
    - 适用于需要动态修改代码的场景
    - 建议用于代码增强和功能扩展

3. **代码生成**
    - 使用 `invent` 包进行代码生成
    - 适用于需要批量生成代码的场景
    - 建议用于项目脚手架和模板生成

## 注意事项

1. 使用前请确保已安装 Go 1.16 或更高版本
2. 代码生成和注入操作可能会修改源文件，请做好备份
3. 建议在开发环境中使用，生产环境请谨慎使用
4. 使用 `inject` 包时，需要安装 `dst` 库：`go get github.com/dave/dst`
