# Redis 环境搭建指南

## Redis 简介

Redis 是一个开源的内存数据结构存储系统，具有高性能、支持多种数据结构、持久化存储等特点。

**主要用途**：
- 缓存系统
- 会话管理
- 消息队列
- 排行榜系统
- 实时计数器

## Docker 安装

### docker-compose.yml

```yaml
version: "3.9"

volumes:
  redis_data:

services:
  redis:
    image: redis:7.0
    container_name: redis
    restart: always
    ports:
      - "6379:6379"
    environment:
      TZ: 'Asia/Shanghai'
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes --requirepass your_password
```

### 启动服务

```bash
# 启动
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f redis
```

## 基本使用

### 连接 Redis

```bash
# 进入容器
docker exec -it redis redis-cli

# 使用密码连接
docker exec -it redis redis-cli -a your_password

# 验证密码
AUTH your_password
```

### 常用命令

**键值操作**

```bash
# 设置键值
SET key value
SET key value EX 3600  # 设置过期时间（秒）

# 获取值
GET key

# 删除键
DEL key

# 检查键是否存在
EXISTS key

# 设置过期时间
EXPIRE key 3600

# 查看剩余时间
TTL key
```

**字符串操作**

```bash
# 递增
INCR counter
INCRBY counter 10

# 递减
DECR counter
DECRBY counter 10

# 追加
APPEND key value
```

**列表操作**

```bash
# 左侧插入
LPUSH list value1 value2

# 右侧插入
RPUSH list value1 value2

# 获取列表
LRANGE list 0 -1

# 弹出元素
LPOP list
RPOP list
```

**集合操作**

```bash
# 添加成员
SADD set member1 member2

# 获取所有成员
SMEMBERS set

# 删除成员
SREM set member1

# 检查成员
SISMEMBER set member1
```

**哈希操作**

```bash
# 设置字段
HSET hash field1 value1 field2 value2

# 获取字段
HGET hash field1

# 获取所有字段
HGETALL hash

# 删除字段
HDEL hash field1
```

**有序集合操作**

```bash
# 添加成员
ZADD zset 1 member1 2 member2

# 获取成员（按分数排序）
ZRANGE zset 0 -1 WITHSCORES

# 获取成员分数
ZSCORE zset member1
```

**说明**：

- `appendonly yes` - 开启 AOF
- `appendfsync everysec` - 每秒同步一次

### 监控命令

```bash
# 查看信息
INFO
INFO memory
INFO stats

# 查看慢日志
SLOWLOG GET 10

# 查看客户端连接
CLIENT LIST

# 查看键空间
DBSIZE

# 监控命令
MONITOR
```

## 安全配置

### 1. 设置密码

```bash
# 临时设置
CONFIG SET requirepass your_password

# 永久设置（在 docker-compose.yml 中）
command: redis-server --requirepass your_password
```

### 2. 禁用危险命令

```yaml
services:
  redis:
    command: >
      redis-server
      --rename-command FLUSHDB ""
      --rename-command FLUSHALL ""
      --rename-command CONFIG ""
      --requirepass your_password
```

### 3. 安全建议

- ✅ 设置强密码
- ✅ 限制网络访问
- ✅ 禁用危险命令
- ✅ 启用持久化
- ✅ 定期备份数据
- ✅ 监控异常访问

## 数据备份与恢复

### 备份数据

```bash
# 手动触发 RDB 备份
docker exec redis redis-cli -a your_password BGSAVE

# 复制 RDB 文件
docker cp redis:/data/dump.rdb ./backup/

# 复制 AOF 文件
docker cp redis:/data/appendonly.aof ./backup/
```

### 恢复数据

```bash
# 停止 Redis
docker stop redis

# 复制备份文件到数据卷
docker cp ./backup/dump.rdb redis:/data/

# 启动 Redis
docker start redis
```

## 常见问题

### 连接超时

```bash
# 检查容器状态
docker ps | grep redis

# 检查端口
netstat -tunlp | grep 6379

# 测试连接
redis-cli -h localhost -p 6379 -a your_password ping
```

### 内存不足

```bash
# 查看内存使用
docker exec redis redis-cli -a your_password INFO memory

# 清理过期键
docker exec redis redis-cli -a your_password --scan --pattern "*" | xargs redis-cli -a your_password DEL

# 手动触发淘汰
docker exec redis redis-cli -a your_password MEMORY PURGE
```

### 性能问题

```bash
# 查看慢日志
docker exec redis redis-cli -a your_password SLOWLOG GET 10

# 查看统计信息
docker exec redis redis-cli -a your_password INFO stats

# 实时监控
docker exec redis redis-cli -a your_password --stat
```

## 参考资料

- [Redis 官方文档](https://redis.io/documentation)
- [Redis 中文网](http://www.redis.cn/)
- [Docker Redis 镜像](https://hub.docker.com/_/redis)
