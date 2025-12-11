# Docker 环境搭建指南

## Docker 简介

Docker 是一个开源的容器化平台，使开发者能够创建、部署和运行应用程序在隔离的环境中。

### 核心优势

- ✅ **轻量高效**：容器只包含运行应用所需的代码和依赖
- ✅ **跨平台兼容**：实现"构建一次，随处运行"
- ✅ **版本控制**：支持镜像版本管理
- ✅ **隔离安全**：容器间相互隔离，防止资源干扰

### 通俗理解

**传统方式**：

- 安装 MySQL、Redis、Java 环境
- 配置各种环境变量
- 处理版本兼容问题

**使用 Docker**：
- 只需安装 Docker
- 拉取所需服务的镜像
- 运行容器即可

💡 **比喻**：Docker 就像一个公司，各种服务（MySQL、Redis等）就是不同的团队。需要时"招聘"团队，不需要时"解散"团队，不会对系统产生任何影响。

## 核心概念

| 概念             | 说明                  | 类比     |
|----------------|---------------------|--------|
| 镜像（Image）      | 只读的模板，包含应用运行所需的所有内容 | 软件安装包  |
| 容器（Container）  | 镜像的运行实例，独立的进程环境     | 运行中的程序 |
| 仓库（Repository） | 存储和分发镜像的地方          | 应用商店   |

## 安装 Docker

### Linux (CentOS/RHEL)

```bash
# 一键安装
curl -fsSL https://get.docker.com | bash -s docker

# 启动 Docker
sudo systemctl start docker
sudo systemctl enable docker

# 验证安装
docker --version
docker-compose --version
```

### macOS / Windows

下载并安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)

- 提供可视化界面
- 简化容器管理
- 自动配置环境

## 常用命令

### 镜像管理

```bash
# 搜索镜像
docker search mysql

# 拉取镜像
docker pull mysql:8.0

# 查看镜像
docker images

# 删除镜像
docker rmi mysql:8.0
```

### 容器管理

```bash
# 运行容器
docker run -d --name mysql -p 3306:3306 mysql:8.0

# 查看运行中的容器
docker ps

# 查看所有容器
docker ps -a

# 停止容器
docker stop mysql

# 启动容器
docker start mysql

# 删除容器
docker rm mysql

# 查看容器日志
docker logs mysql

# 进入容器
docker exec -it mysql bash
```

### 系统管理

```bash
# 查看 Docker 信息
docker info

# 清理无用资源
docker system prune -a

# 查看资源使用
docker stats
```

## 服务安装示例

### MySQL

```bash
docker run -d \
  --name mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=your_password \
  -v /data/mysql:/var/lib/mysql \
  --restart=always \
  mysql:8.0
```

### Redis

```bash
docker run -d \
  --name redis \
  -p 6379:6379 \
  -v /data/redis:/data \
  --restart=always \
  redis:latest redis-server --appendonly yes
```

### RabbitMQ

```bash
docker run -d \
  --name rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin \
  -e RABBITMQ_DEFAULT_PASS=admin \
  --restart=always \
  rabbitmq:management
```

## Docker Compose

### 安装

```bash
# Linux
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证
docker-compose --version
```

### 使用示例

```yaml
version: '3'
services:
  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: password
    volumes:
      - /data/mysql:/var/lib/mysql
    restart: always

  redis:
    image: redis:latest
    ports:
      - "6379:6379"
    volumes:
      - /data/redis:/data
    restart: always
```

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 查看日志
docker-compose logs -f
```

## 最佳实践

### 1. 数据持久化

- ✅ 使用数据卷挂载（`-v`）
- ✅ 避免容器内存储重要数据
- ✅ 定期备份数据

### 2. 资源管理

```bash
# 限制内存和 CPU
docker run -d \
  --memory="512m" \
  --cpus="1.0" \
  mysql:8.0
```

### 3. 安全建议

- ✅ 使用官方镜像
- ✅ 定期更新镜像
- ✅ 限制容器权限
- ✅ 使用非 root 用户
- ✅ 配置防火墙规则

### 4. 网络配置

```bash
# 创建自定义网络
docker network create blog-network

# 容器加入网络
docker run -d --network blog-network --name mysql mysql:8.0
```

## 常见问题

### 容器无法启动

```bash
# 查看日志
docker logs <container_name>

# 检查端口占用
netstat -tunlp | grep <port>

# 检查资源使用
docker stats
```

### 数据备份

```bash
# 备份容器数据
docker exec mysql mysqldump -u root -p database > backup.sql

# 备份数据卷
docker run --rm -v /data/mysql:/backup -v $(pwd):/backup-dest busybox tar czf /backup-dest/mysql-backup.tar.gz /backup
```

### 清理空间

```bash
# 清理无用镜像
docker image prune -a

# 清理无用容器
docker container prune

# 清理无用数据卷
docker volume prune

# 清理所有无用资源
docker system prune -a --volumes
```

## 参考资料

- [Docker 官方文档](https://docs.docker.com/)
- [Docker Hub](https://hub.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
