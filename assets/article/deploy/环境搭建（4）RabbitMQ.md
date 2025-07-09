# RabbitMQ 环境搭建指南

## RabbitMQ 简介

RabbitMQ 是一个开源的消息代理软件，支持多种消息协议，具有高可用性、消息持久化、集群支持等特点。

**主要用途**：
- 异步消息处理
- 应用解耦
- 流量削峰
- 任务队列
- 日志处理

## Docker 安装

### docker-compose.yml

```yaml
version: "3.9"

volumes:
  rabbitmq_data:

services:
  rabbitmq:
    image: rabbitmq:3.12-management
    container_name: rabbitmq
    restart: always
    ports:
      - "5672:5672"    # AMQP 协议端口
      - "15672:15672"  # Web 管理界面
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: your_password
      TZ: 'Asia/Shanghai'
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
```

### 启动服务

```bash
# 启动
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f rabbitmq
```

## 管理界面

### 访问地址

- URL: http://localhost:15672
- 用户名: admin
- 密码: your_password

### 主要功能

- 查看队列状态
- 监控消息流量
- 管理用户权限
- 配置交换机和队列
- 查看连接信息

## 核心概念

| 概念          | 说明    | 作用         |
|-------------|-------|------------|
| Exchange    | 消息交换机 | 接收消息并路由到队列 |
| Queue       | 消息队列  | 存储消息       |
| Binding     | 绑定关系  | 连接交换机和队列   |
| Routing Key | 路由键   | 决定消息路由规则   |

### Exchange 类型

| 类型      | 说明    | 使用场景      |
|---------|-------|-----------|
| Direct  | 直连交换机 | 精确匹配路由键   |
| Fanout  | 扇出交换机 | 广播消息到所有队列 |
| Topic   | 主题交换机 | 模式匹配路由键   |
| Headers | 头交换机  | 根据消息头路由   |

## 消息模式

### 1. 简单队列

```
Producer → Queue → Consumer
```

一个生产者，一个消费者，一个队列。

### 2. 工作队列

```
Producer → Queue → Consumer1
                 → Consumer2
```

一个生产者，多个消费者，轮询分发消息。

### 3. 发布/订阅

```
Producer → Exchange → Queue1 → Consumer1
                   → Queue2 → Consumer2
```

一个生产者，多个消费者，广播消息。

### 4. 路由模式

```
Producer → Exchange → Queue1 (routing_key: error) → Consumer1
                   → Queue2 (routing_key: info)  → Consumer2
```

根据路由键分发消息。

### 5. 主题模式

```
Producer → Exchange → Queue1 (pattern: *.error) → Consumer1
                   → Queue2 (pattern: log.*)   → Consumer2
```

根据模式匹配分发消息。

## 基本操作

### 创建队列

```bash
# 进入容器
docker exec -it rabbitmq bash

# 创建队列
rabbitmqadmin declare queue name=my_queue durable=true

# 查看队列
rabbitmqctl list_queues
```

### 创建交换机

```bash
# 创建交换机
rabbitmqadmin declare exchange name=my_exchange type=direct

# 查看交换机
rabbitmqctl list_exchanges
```

### 创建绑定

```bash
# 绑定队列到交换机
rabbitmqadmin declare binding source=my_exchange destination=my_queue routing_key=my_key
```

### 用户管理

```bash
# 创建用户
rabbitmqctl add_user username password

# 设置用户角色
rabbitmqctl set_user_tags username administrator

# 设置权限
rabbitmqctl set_permissions -p / username ".*" ".*" ".*"

# 查看用户
rabbitmqctl list_users
```

## 性能优化

### 配置优化

```yaml
services:
  rabbitmq:
    environment:
      # 内存限制
      RABBITMQ_VM_MEMORY_HIGH_WATERMARK: 0.8
      # 磁盘限制
      RABBITMQ_DISK_FREE_LIMIT: 2GB
      # 心跳超时
      RABBITMQ_HEARTBEAT: 60
```

### 消息持久化

```bash
# 声明持久化队列
rabbitmqadmin declare queue name=my_queue durable=true

# 发送持久化消息
# delivery_mode=2 表示持久化
```

### 预取数量

```bash
# 设置预取数量（消费者一次获取的消息数）
# 在消费者代码中设置 prefetch_count
```

## 监控与维护

### 监控命令

```bash
# 查看节点状态
rabbitmqctl status

# 查看队列信息
rabbitmqctl list_queues name messages consumers

# 查看交换机
rabbitmqctl list_exchanges

# 查看绑定
rabbitmqctl list_bindings

# 查看连接
rabbitmqctl list_connections

# 查看通道
rabbitmqctl list_channels
```

### 清理操作

```bash
# 清空队列
rabbitmqctl purge_queue my_queue

# 删除队列
rabbitmqctl delete_queue my_queue

# 删除交换机
rabbitmqctl delete_exchange my_exchange
```

## 安全配置

### 1. 修改默认密码

```bash
# 修改密码
rabbitmqctl change_password admin new_password
```

### 2. 限制访问

```yaml
services:
  rabbitmq:
    ports:
      - "127.0.0.1:5672:5672"    # 只允许本地访问
      - "127.0.0.1:15672:15672"
```

### 3. 启用 SSL

```yaml
services:
  rabbitmq:
    environment:
      RABBITMQ_SSL_CERTFILE: /path/to/cert.pem
      RABBITMQ_SSL_KEYFILE: /path/to/key.pem
      RABBITMQ_SSL_CACERTFILE: /path/to/ca.pem
```

### 4. 安全建议

- ✅ 修改默认密码
- ✅ 限制网络访问
- ✅ 启用 SSL/TLS
- ✅ 配置访问控制
- ✅ 定期备份数据
- ✅ 监控异常访问

## 常见问题

### 消息丢失

**原因**：

- 消息未持久化
- 队列未持久化
- 消费者未确认

**解决方案**：

```bash
# 1. 启用队列持久化
rabbitmqadmin declare queue name=my_queue durable=true

# 2. 发送持久化消息（在代码中设置 delivery_mode=2）

# 3. 启用消费者确认（在代码中手动确认）
```

### 消息积压

**原因**：

- 消费者处理慢
- 消费者数量不足
- 网络问题

**解决方案**：

```bash
# 1. 增加消费者数量

# 2. 优化消费者处理逻辑

# 3. 检查队列状态
rabbitmqctl list_queues name messages consumers
```

### 连接失败

```bash
# 检查容器状态
docker ps | grep rabbitmq

# 检查端口
netstat -tunlp | grep 5672

# 查看日志
docker logs rabbitmq
```

## 参考资料

- [RabbitMQ 官方文档](https://www.rabbitmq.com/documentation.html)
- [RabbitMQ 中文文档](https://rabbitmq.mr-ping.com/)
- [Docker RabbitMQ 镜像](https://hub.docker.com/_/rabbitmq)
