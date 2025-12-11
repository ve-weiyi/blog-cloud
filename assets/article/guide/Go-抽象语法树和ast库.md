# Go 抽象语法树（AST）详解

## 前言

作为一名代码开发程序员，我常常思考一个程序文件，比如 `hello.go` 是如何被编译器理解的。平常在编写程序时，IDE
又是如何提供代码提示的。在这奥妙无穷的背后，AST（Abstract Syntax Tree）抽象语法树功不可没，它站在每一行程序的身后，默默无闻地工作，为繁荣的互联网世界立下了汗马功劳。

## 什么是 AST

AST（Abstract Syntax Tree，抽象语法树）是源代码语法结构的一种抽象表示形式。它是编译器理解和处理代码的基础，也是 IDE
提供代码提示、重构等功能的核心。

AST 通常是一个树形结构，用于表示程序的语法层次关系。每个节点代表源代码中的一个语法结构，例如表达式、语句、函数声明等。

## 编译过程

编译过程是将高级语言源代码转换为目标机器代码或其他中间表示的过程。它通常包括以下几个主要阶段：

```
源代码 → 词法分析 → 语法分析 → 语义分析 → 中间代码生成 → 优化 → 目标代码生成 → 链接 → 可执行文件
```

### 1. 词法分析（Lexical Analysis）

- **输入**：源代码
- **输出**：词法单元（tokens）
- **任务**：将源代码分解为基本的词法单元，如关键字、标识符、运算符等

### 2. 语法分析（Syntax Analysis）

- **输入**：词法单元序列
- **输出**：语法树（AST）或其他中间表示
- **任务**：根据语法规则检查词法单元序列的结构，构建语法树表示源代码的语法结构

### 3. 语义分析（Semantic Analysis）

- **输入**：语法树或中间表示
- **输出**：语义信息
- **任务**：检查源代码的语义正确性，捕捉并处理语法上合法但语义上错误的结构，生成符号表等

### 4. 中间代码生成（Intermediate Code Generation）

- **输入**：语法树或中间表示
- **输出**：中间代码
- **任务**：将语法树转换为一种中间表示，以便进行后续优化和目标代码生成

### 5. 优化（Optimization）

- **输入**：中间代码
- **输出**：优化后的中间代码
- **任务**：对中间代码进行各种优化，提高程序性能、减小代码体积等

### 6. 目标代码生成（Code Generation）

- **输入**：优化后的中间代码
- **输出**：目标机器代码或汇编代码
- **任务**：将中间代码转换为目标机器代码，根据目标平台生成可执行文件

### 7. 链接（Linking）

- **输入**：目标机器代码或汇编代码，可能包括库文件
- **输出**：可执行文件
- **任务**：将多个目标文件链接在一起，解决符号引用、地址重定位等问题，生成最终的可执行文件

## Token（词法单元）

Token 是源代码中的最小语法单元，是编译器或解释器在进行词法分析时识别和处理的基本单位。Token
是源代码经过词法分析器处理后得到的标识符，它代表了源代码中的不可分割的词法结构。

### Token 类型

| 类型               | 说明          | 示例                              |
|------------------|-------------|---------------------------------|
| 关键字（Keywords）    | 保留字，具有特殊含义  | `if`、`else`、`for`、`func`        |
| 标识符（Identifiers） | 变量、函数、类等的名称 | `userName`、`getUser`、`MyStruct` |
| 常量（Literals）     | 固定值的词法单元    | `42`、`3.14`、`"hello"`           |
| 运算符（Operators）   | 执行操作的符号     | `+`、`-`、`*`、`/`、`==`            |
| 分隔符（Delimiters）  | 源代码结构的分隔符   | `(`、`)`、`{`、`}`、`;`             |
| 注释（Comments）     | 注释内容        | `// comment`、`/* comment */`    |

### Token 示例

```go
// 源代码
x = 10 + 5

// 对应的 Token
Token(IDENT, "x")  // 标识符
Token(ASSIGN, "=") // 赋值运算符
Token(INT, "10") // 整数常量
Token(PLUS, "+") // 加法运算符
Token(INT, "5")  // 整数常量
```

Token 的产生是由词法分析器（Lexer 或 Scanner）负责的，它扫描源代码，将字符序列组合成有意义的 Token。Token
提供给语法分析器（Parser）使用，用于构建 AST（抽象语法树）等进一步的语法分析和语义分析。

## Go AST 包使用

### FileSet

FileSet 是用于跟踪源代码文件和位置信息的结构。它维护源代码文件的集合，每个文件都关联有唯一的标识符，并记录了源代码中各个位置的行号、列号等信息。

### 解析文件

```go
import (
"go/ast"
"go/parser"
"go/token"
"log"
)

// 解析 Go 代码文件
fset := token.NewFileSet()
f, err := parser.ParseFile(fset, "hello.go", nil, parser.ParseComments)
if err != nil {
log.Fatal(err)
}

// 打印 AST
ast.Print(fset, f)
```

### ParseFile 函数详解

```go
func ParseFile(fset *token.FileSet, filename string, src any, mode Mode) (f *ast.File, err error)
```

**参数说明**

| 参数         | 类型               | 说明                         |
|------------|------------------|----------------------------|
| `fset`     | `*token.FileSet` | 用于记录位置信息的文件集合              |
| `filename` | `string`         | 源文件的文件名                    |
| `src`      | `any`            | 源代码，可以是字符串、字节数组或 io.Reader |
| `mode`     | `Mode`           | 控制解析器的模式，例如是否跳过对象解析阶段      |

**返回值**

| 返回值   | 类型          | 说明                      |
|-------|-------------|-------------------------|
| `f`   | `*ast.File` | 解析后的 AST 树，表示整个源文件的语法结构 |
| `err` | `error`     | 解析过程中的错误，如果有的话          |

**解析模式**

```go
const (
PackageClauseOnly Mode = 1 << iota // 只解析 package 子句
ImportsOnly                        // 只解析 import 声明
ParseComments     // 解析注释
Trace             // 打印跟踪信息
DeclarationErrors // 报告声明错误
SpuriousErrors    // 报告所有错误
)
```

### AST 节点类型

在 Go 语言中，AST 主要由三个大类别的节点构成：

#### 1. 基本节点（Basic Nodes）

| 节点类型               | 说明    | 示例         |
|--------------------|-------|------------|
| `ast.Ident`        | 标识符   | 变量名、函数名    |
| `ast.BasicLit`     | 基本字面量 | 整数、浮点数、字符串 |
| `ast.CompositeLit` | 复合字面量 | 数组、切片、字典   |

#### 2. 表达式节点（Expression Nodes）

| 节点类型               | 说明               |
|--------------------|------------------|
| `ast.BinaryExpr`   | 二元表达式（如 `a + b`） |
| `ast.UnaryExpr`    | 一元表达式（如 `-a`）    |
| `ast.CallExpr`     | 函数调用表达式          |
| `ast.IndexExpr`    | 索引表达式（如 `a[i]`）  |
| `ast.SelectorExpr` | 选择器表达式（如 `a.b`）  |
| `ast.StarExpr`     | 指针表达式（如 `*p`）    |

#### 3. 语句节点（Statement Nodes）

| 节点类型             | 说明           |
|------------------|--------------|
| `ast.ExprStmt`   | 表达式语句        |
| `ast.AssignStmt` | 赋值语句         |
| `ast.IfStmt`     | 条件语句         |
| `ast.ForStmt`    | 循环语句         |
| `ast.SwitchStmt` | 开关语句         |
| `ast.RangeStmt`  | for range 循环 |
| `ast.ReturnStmt` | 返回语句         |
| `ast.DeferStmt`  | defer 语句     |
| `ast.GoStmt`     | go 语句        |

#### 4. 声明节点（Declaration Nodes）

| 节点类型             | 说明                          |
|------------------|-----------------------------|
| `ast.GenDecl`    | 通用声明（import、const、var、type） |
| `ast.FuncDecl`   | 函数声明                        |
| `ast.TypeSpec`   | 类型声明                        |
| `ast.ValueSpec`  | 值声明（const、var）              |
| `ast.ImportSpec` | 导入声明                        |

### 节点接口

所有节点都实现了 `ast.Node` 接口：

```go
// 所有节点都实现 Node 接口
type Node interface {
Pos() token.Pos // 起始位置
End() token.Pos // 结束位置
}

// 所有表达式节点实现 Expr 接口
type Expr interface {
Node
exprNode()
}

// 所有语句节点实现 Stmt 接口
type Stmt interface {
Node
stmtNode()
}

// 所有声明节点实现 Decl 接口
type Decl interface {
Node
declNode()
}
```

## 遍历 AST

### 使用 ast.Inspect

`ast.Inspect` 是最简单的遍历方式，它会深度优先遍历 AST 树：

```go
ast.Inspect(f, func (n ast.Node) bool {
// 查找所有函数声明
if fn, ok := n.(*ast.FuncDecl); ok {
fmt.Printf("Function: %s\n", fn.Name.Name)

// 打印函数参数
for _, param := range fn.Type.Params.List {
for _, name := range param.Names {
fmt.Printf("  Param: %s\n", name.Name)
}
}
}

// 返回 true 继续遍历子节点，返回 false 跳过子节点
return true
})
```

### 使用 ast.Walk

`ast.Walk` 提供了更灵活的遍历方式，需要实现 `ast.Visitor` 接口：

```go
type visitor struct{}

func (v visitor) Visit(n ast.Node) ast.Visitor {
if n == nil {
return nil
}

// 根据节点类型进行不同处理
switch x := n.(type) {
case *ast.FuncDecl:
fmt.Printf("Function: %s\n", x.Name.Name)
case *ast.AssignStmt:
fmt.Println("Assignment statement")
case *ast.CallExpr:
fmt.Println("Function call")
}

return v
}

ast.Walk(visitor{}, f)
```

## 实际应用示例

### 1. 查找所有函数

```go
func findFunctions(fset *token.FileSet, f *ast.File) {
ast.Inspect(f, func (n ast.Node) bool {
if fn, ok := n.(*ast.FuncDecl); ok {
pos := fset.Position(fn.Pos())
fmt.Printf("Function %s at %s\n", fn.Name.Name, pos)
}
return true
})
}
```

### 2. 查找所有导入

```go
func findImports(f *ast.File) {
for _, imp := range f.Imports {
path := imp.Path.Value
if imp.Name != nil {
fmt.Printf("import %s %s\n", imp.Name.Name, path)
} else {
fmt.Printf("import %s\n", path)
}
}
}
```

### 3. 统计代码行数

```go
func countLines(fset *token.FileSet, f *ast.File) int {
return fset.Position(f.End()).Line
}
```

### 4. 查找未使用的变量

```go
func findUnusedVars(fset *token.FileSet, f *ast.File) {
declared := make(map[string]bool)
used := make(map[string]bool)

// 查找声明的变量
ast.Inspect(f, func (n ast.Node) bool {
if assign, ok := n.(*ast.AssignStmt); ok {
for _, lhs := range assign.Lhs {
if ident, ok := lhs.(*ast.Ident); ok {
declared[ident.Name] = true
}
}
}
return true
})

// 查找使用的变量
ast.Inspect(f, func (n ast.Node) bool {
if ident, ok := n.(*ast.Ident); ok {
used[ident.Name] = true
}
return true
})

// 找出未使用的变量
for name := range declared {
if !used[name] {
fmt.Printf("Unused variable: %s\n", name)
}
}
}
```

## 实际应用场景

### 1. 代码分析工具

- **golint**：代码风格检查
- **go vet**：静态分析工具
- **staticcheck**：高级静态分析

### 2. 代码生成工具

- **go generate**：代码生成
- **mockgen**：生成 mock 代码
- **stringer**：生成 String() 方法

### 3. 代码重构工具

- **gofmt**：代码格式化
- **goimports**：自动管理导入
- **gorename**：重命名标识符

### 4. IDE 功能

- 代码补全
- 跳转到定义
- 查找引用
- 代码重构

## 参考资料

- [Go AST 官方文档](https://pkg.go.dev/go/ast)
- [Go Parser 官方文档](https://pkg.go.dev/go/parser)
- [Go Token 官方文档](https://pkg.go.dev/go/token)
- [Go AST Viewer](https://yuroyoro.github.io/goast-viewer/)
