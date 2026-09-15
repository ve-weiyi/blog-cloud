.PHONY: help deps docker-deps docker-app docker-up docker-deps-down docker-down run-app-rpc run-app-api run-admin-api build clean k8s-up k8s-down

# 默认目标
help:
	@echo "可用命令："
	@echo "  make deps        - 安装依赖"
	@echo "  make docker-deps    - 启动依赖服务(MySQL→Redis→RabbitMQ)"
	@echo "  make docker-app     - 启动 App 容器服务"
	@echo "  make docker-up      - 启动所有容器服务(deps+app)"
	@echo "  make docker-deps-down - 停止依赖服务"
	@echo "  make docker-down    - 停止所有容器服务"
	@echo "  make run-app-rpc   - 启动 RPC 服务"
	@echo "  make run-app-api   - 启动博客前台服务"
	@echo "  make run-admin-api - 启动管理后台服务"
	@echo "  make build         - 编译所有服务"
	@echo "  make clean         - 清理编译文件"
	@echo "  make k8s-up        - 部署到 K8s"
	@echo "  make k8s-down      - 移除 K8s 资源"

# 安装依赖
deps:
	go mod tidy
	go mod download
# 安装goctl工具
	go install github.com/zeromicro/go-zero/tools/goctl@latest
# 安装grpc工具(使用goctl)
	goctl env check --install --verbose --force

# 启动依赖服务
docker-deps:
	@echo "启动 MySQL..."
	docker compose -f ../deploy/docker-compose/mysql/mysql.yaml up -d
	@echo "启动 Redis..."
	docker compose -f ../deploy/docker-compose/redis/redis.yaml up -d
	@echo "启动 RabbitMQ..."
	docker compose -f ../deploy/docker-compose/rabbitmq/rabbitmq.yaml up -d
	@echo "依赖服务启动完成！"

# 启动 App 容器服务
docker-app:
	@echo "启动 App 服务..."
	docker compose -f ../deploy/docker/docker-compose.yml up -d
	@echo "App 服务启动完成！"

# 启动所有容器服务
docker-up:
	@$(MAKE) docker-deps
	@$(MAKE) docker-app

# 停止依赖服务
docker-deps-down:
	@echo "停止 RabbitMQ..."
	docker compose -f ../deploy/docker-compose/rabbitmq/rabbitmq.yaml down
	@echo "停止 Redis..."
	docker compose -f ../deploy/docker-compose/redis/redis.yaml down
	@echo "停止 MySQL..."
	docker compose -f ../deploy/docker-compose/mysql/mysql.yaml down
	@echo "依赖服务已停止！"

# 停止所有容器服务
docker-down:
	@echo "停止 App 服务..."
	docker compose -f ../deploy/docker/docker-compose.yml down
	@echo "停止 RabbitMQ..."
	docker compose -f ../deploy/docker-compose/rabbitmq/rabbitmq.yaml down
	@echo "停止 Redis..."
	docker compose -f ../deploy/docker-compose/redis/redis.yaml down
	@echo "停止 MySQL..."
	docker compose -f ../deploy/docker-compose/mysql/mysql.yaml down
	@echo "所有服务已停止！"

# 部署到 K8s
k8s-up:
	@echo "创建 Namespace..."
	kubectl apply -f ../deploy/k8s/namespace.yaml
	@echo "部署 App 服务..."
	kubectl apply -f ../deploy/k8s/app/deployment.yaml
	kubectl apply -f ../deploy/k8s/app/ingress.yaml
	kubectl apply -f ../deploy/k8s/app/admin-web.yaml
	kubectl apply -f ../deploy/k8s/app/app-web.yaml
	@echo "K8s 部署完成！"

# 移除 K8s 资源
k8s-down:
	@echo "移除 App 服务..."
	kubectl delete -f ../deploy/k8s/app/app-web.yaml --ignore-not-found
	kubectl delete -f ../deploy/k8s/app/admin-web.yaml --ignore-not-found
	kubectl delete -f ../deploy/k8s/app/ingress.yaml --ignore-not-found
	kubectl delete -f ../deploy/k8s/app/deployment.yaml --ignore-not-found
	@echo "K8s 资源已移除！"

# 启动 RPC 服务（开发模式）
run-app-rpc:
	go run rpc/blog/blog.go -f rpc/blog/etc/app-rpc.yaml

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
