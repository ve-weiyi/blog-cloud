.PHONY: help deps init-config run-blog-rpc run-app-api run-admin-api build clean

# 本 Makefile 只管后端程序的编译与本地启动。部署编排（Docker / K8s）见 deploy/ 下的文档。

# 默认目标
help:
	@echo "可用命令："
	@echo "  make deps        - 安装依赖"
	@echo "  make init-config - 从 *.example.yaml 生成本地配置（已存在则跳过）"
	@echo "  make run-blog-rpc  - 启动 RPC 服务"
	@echo "  make run-app-api   - 启动博客前台服务"
	@echo "  make run-admin-api - 启动管理后台服务"
	@echo "  make build         - 编译所有服务"
	@echo "  make clean         - 清理编译文件"

# 安装依赖
deps:
	go mod tidy
	go mod download
# 安装goctl工具
	go install github.com/zeromicro/go-zero/tools/goctl@latest
# 安装grpc工具(使用goctl)
	goctl env check --install --verbose --force

# 从 *.example.yaml 生成本地运行时配置。
# 这些 *.yaml 在 .gitignore 中，不入库；各人按需改端口、库名与凭证。
init-config:
	@for base in rpc/blog/etc/blog-rpc api/app/etc/app-api api/admin/etc/admin-api; do \
		src="$$base.example.yaml"; dst="$$base.yaml"; \
		if [ -f "$$dst" ]; then \
			echo "  已存在，跳过: $$dst"; \
		else \
			cp "$$src" "$$dst"; \
			echo "  已生成: $$dst"; \
		fi; \
	done

# 启动 RPC 服务（开发模式）
run-blog-rpc:
	go run rpc/blog/blog.go -f rpc/blog/etc/blog-rpc.yaml

# 启动博客前台服务（开发模式）
run-app-api:
	go run api/app/app.go -f api/app/etc/app-api.yaml

# 启动管理后台服务（开发模式）
run-admin-api:
	go run api/admin/admin.go -f api/admin/etc/admin-api.yaml

# 编译所有服务
build:
	@echo "编译 RPC 服务..."
	go build -o bin/blog-rpc rpc/blog/blog.go
	@echo "编译博客前台服务..."
	go build -o bin/blog-api api/app/app.go
	@echo "编译管理后台服务..."
	go build -o bin/admin-api api/admin/admin.go
	@echo "编译完成！"

# 清理编译文件
clean:
	rm -rf bin/
	@echo "清理完成！"
