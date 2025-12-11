# Cobra 命令行工具使用指南

Cobra 是一个强大的 Go 语言命令行框架，用于创建现代化的 CLI 应用程序。

## 官方资源

- GitHub: https://github.com/spf13/cobra
- 文档: https://cobra.dev/

## 快速开始

### 1. 安装 Cobra CLI

```bash
go install github.com/spf13/cobra-cli@latest
```

验证安装：

```bash
cobra-cli --version
```

### 2. 初始化项目

```bash
# 基础初始化
cobra-cli init

# 指定作者信息
cobra-cli init --author "your-email@example.com"

# 指定开源协议
cobra-cli init --license apache
cobra-cli init --license mit
```

### 3. 添加命令

```bash
# 添加版本命令
cobra-cli add version

# 添加数据库迁移命令
cobra-cli add migrate

# 添加 API 服务命令
cobra-cli add api
```

### 4. 运行命令

```bash
# 查看帮助信息
go run main.go -h
go run main.go migrate --help

# 运行命令
go run main.go version
go run main.go migrate --action=reset
go run main.go api -c file -f ./config.yaml
```

## 项目结构

```
project/
├── cmd/
│   ├── root.go      # 根命令
│   ├── version.go   # 版本命令
│   ├── migrate.go   # 迁移命令
│   └── api.go       # API 命令
├── main.go          # 入口文件
└── go.mod
```

## 常用功能

### 添加命令参数

```go
// 持久化参数（全局参数）
rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

// 局部参数
migrateCmd.Flags().StringP("action", "a", "up", "migration action")
```

### 添加子命令

```go
// 在 cmd/migrate.go 中
var migrateUpCmd = &cobra.Command{
Use:   "up",
Short: "Run migrations",
Run: func (cmd *cobra.Command, args []string) {
// 执行迁移
},
}

func init() {
migrateCmd.AddCommand(migrateUpCmd)
}
```

## 最佳实践

1. **命令命名**：使用小写字母和连字符，如 `db-migrate`
2. **参数设计**：提供简写和全写两种形式，如 `-c` 和 `--config`
3. **帮助信息**：为每个命令添加清晰的 `Short` 和 `Long` 描述
4. **错误处理**：使用 `cobra.CheckErr()` 统一处理错误
5. **配置文件**：使用 Viper 配合 Cobra 管理配置


