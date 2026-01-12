.PHONY: help deps docker-up docker-down run-rpc run-blog run-admin run-all build clean

# 默认目标
help:
	@echo "可用命令："
	@echo "  make deps        - 安装依赖"
	@echo "  make docker-up   - 启动依赖服务(MySQL/Redis/RabbitMQ)"
	@echo "  make docker-down - 停止依赖服务"
	@echo "  make run-rpc     - 启动 RPC 服务"
	@echo "  make run-blog    - 启动博客前台服务"
	@echo "  make run-admin   - 启动管理后台服务"
	@echo "  make run-all     - 启动所有服务"
	@echo "  make build       - 编译所有服务"
	@echo "  make clean       - 清理编译文件"

# 安装依赖
deps:
	go mod tidy
	go mod download
# 安装goctl工具
	go install github.com/zeromicro/go-zero/tools/goctl@latest
# 安装grpc工具(使用goctl)
	goctl env check --install --verbose --force

# 启动依赖服务
docker-up:
	cd ../deploy/docker-compose/data && docker-compose up -d

# 停止依赖服务
docker-down:
	cd ../deploy/docker-compose/data && docker-compose down

# 启动 RPC 服务（开发模式）
run-rpc:
	go run service/rpc/blog/blog.go -f service/rpc/blog/etc/blog.yaml

# 启动博客前台服务（开发模式）
run-blog:
	go run service/api/blog/blog.go -f service/api/blog/etc/blog.yaml

# 启动管理后台服务（开发模式）
run-admin:
	go run service/api/admin/admin.go -f service/api/admin/etc/admin.yaml

# 启动所有服务（需要在不同终端运行）
run-all:
	@echo "请在不同终端分别运行："
	@echo "  make run-rpc"
	@echo "  make run-blog"
	@echo "  make run-admin"

# 编译所有服务
build:
	@echo "编译 RPC 服务..."
	go build -o bin/blog-rpc service/rpc/blog/blog.go
	@echo "编译博客前台服务..."
	go build -o bin/blog-api service/api/blog/blog.go
	@echo "编译管理后台服务..."
	go build -o bin/admin-api service/api/admin/admin.go
	@echo "编译完成！"

# 清理编译文件
clean:
	rm -rf bin/
	@echo "清理完成！"
