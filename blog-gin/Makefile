.PHONY: help deps run build clean swagger swagger-fmt fmt

# 默认目标
help:
	@echo "可用命令："
	@echo "  make deps        - 安装依赖"
	@echo "  make run         - 启动服务"
	@echo "  make build       - 编译服务"
	@echo "  make swagger     - 生成 Swagger 文档"
	@echo "  make swagger-fmt - 格式化 Swagger 注解"
	@echo "  make fmt         - 格式化代码"
	@echo "  make clean       - 清理编译文件"

# 安装依赖
deps:
	go mod tidy
	go mod download

# 启动服务
run:
	go run main.go

# 编译服务
build:
	@echo "编译服务..."
	go build -o bin/blog-gin main.go
	@echo "编译完成！"

# 清理编译文件
clean:
	rm -rf bin/
	@echo "清理完成！"

# 生成 Swagger 文档
swagger:
	swag init

# 格式化 Swagger 注解
swagger-fmt:
	swag fmt

# 格式化代码
fmt:
	go fmt ./...
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .
