# MySQL 环境搭建指南

## MySQL 简介

MySQL 是一个开源的关系型数据库管理系统（RDBMS），具有高性能、高可靠性、跨平台等特点。

**主要用途**：
- 用户数据存储
- 业务数据管理
- 日志数据存储
- 配置信息管理

## Docker 安装

### docker-compose.yml

```yaml
version: "3.9"

volumes:
  mysql_data:

services:
  mysql:
    image: mysql:8.0
    container_name: mysql
    restart: always
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: 'your_password'
      MYSQL_DATABASE: 'blog_veweiyi'
      TZ: 'Asia/Shanghai'
    volumes:
      - mysql_data:/var/lib/mysql
    command:
      --character-set-server=utf8mb4
      --collation-server=utf8mb4_unicode_ci
      --max_connections=1000
```

### 启动服务

```bash
# 启动
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f mysql
```

## 基本使用

### 连接数据库

```bash
# 进入容器
docker exec -it mysql mysql -u root -p

# 直接执行 SQL
docker exec -it mysql mysql -u root -p -e "SHOW DATABASES;"
```

### 常用命令

```sql
-- 查看数据库
SHOW DATABASES;

-- 创建数据库
CREATE
DATABASE blog_veweiyi CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE
blog_veweiyi;

-- 查看表
SHOW TABLES;

-- 创建用户
CREATE
USER 'blog'@'%' IDENTIFIED BY 'password';

-- 授予权限
GRANT ALL PRIVILEGES ON blog_veweiyi.* TO
'blog'@'%';
FLUSH
PRIVILEGES;

-- 查看用户权限
SHOW
GRANTS FOR 'blog'@'%';
```

## 数据备份与恢复

### 备份数据

```bash
# 备份整个数据库
docker exec mysql mysqldump -u root -p blog_veweiyi > backup.sql

# 备份所有数据库
docker exec mysql mysqldump -u root -p --all-databases > all_backup.sql

# 备份特定表
docker exec mysql mysqldump -u root -p blog_veweiyi t_user t_article > tables_backup.sql
```

### 恢复数据

```bash
# 恢复数据库
docker exec -i mysql mysql -u root -p blog_veweiyi < backup.sql

# 恢复所有数据库
docker exec -i mysql mysql -u root -p < all_backup.sql
```

## 性能优化

### 配置优化

```yaml
services:
  mysql:
    command:
      --character-set-server=utf8mb4
      --collation-server=utf8mb4_unicode_ci
      --max_connections=1000              # 最大连接数
      --innodb_buffer_pool_size=1G        # InnoDB 缓冲池大小
      --innodb_log_file_size=256M         # 日志文件大小
      --slow_query_log=1                  # 开启慢查询日志
      --long_query_time=2                 # 慢查询时间阈值（秒）
```

### 监控命令

```sql
-- 查看连接数
SHOW
STATUS LIKE 'Threads_connected';

-- 查看最大连接数
SHOW
VARIABLES LIKE 'max_connections';

-- 查看慢查询
SHOW
VARIABLES LIKE 'slow_query%';

-- 查看数据库大小
SELECT table_schema AS 'Database', ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.tables
GROUP BY table_schema;
```

## 安全配置

### 1. 修改 root 密码

```sql
ALTER
USER 'root'@'localhost' IDENTIFIED BY 'new_password';
FLUSH
PRIVILEGES;
```

### 2. 限制远程访问

```sql
-- 只允许特定 IP 访问
CREATE
USER 'blog'@'192.168.1.100' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON blog_veweiyi.* TO
'blog'@'192.168.1.100';
```

### 3. 删除匿名用户

```sql
DELETE
FROM mysql.user
WHERE User = '';
FLUSH
PRIVILEGES;
```

### 4. 安全建议

- ✅ 使用强密码
- ✅ 限制远程访问
- ✅ 定期备份数据
- ✅ 及时更新版本
- ✅ 监控异常访问

## 常见问题

### 连接被拒绝

```bash
# 检查容器状态
docker ps | grep mysql

# 检查端口占用
netstat -tunlp | grep 3306

# 查看错误日志
docker logs mysql
```

### 忘记 root 密码

```bash
# 停止容器
docker stop mysql

# 以跳过权限检查方式启动
docker run -d --name mysql-temp \
  -e MYSQL_ROOT_PASSWORD=new_password \
  mysql:8.0

# 重置密码后重启
docker restart mysql
```

### 数据库损坏

```bash
# 检查并修复表
docker exec -it mysql mysqlcheck -u root -p --auto-repair --all-databases
```

## 参考资料

- [MySQL 官方文档](https://dev.mysql.com/doc/)
- [Docker MySQL 镜像](https://hub.docker.com/_/mysql)
- [MySQL 性能优化](https://dev.mysql.com/doc/refman/8.0/en/optimization.html)
