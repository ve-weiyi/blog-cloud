
## 前言

### 什么是Mysql ？
Mysql是一个开源的关系型数据库管理系统，它支持SQL语言，可以用于存储、查询、更新和删除数据。

## 用途
Mysql用于存储用户数据、日志数据、配置数据等等，它是Web开发中最常用的数据库之一。

## 安装

使用docker安装Mysql

1. 复制以下内容到docker-compose.yml文件中
```yaml
version: "3.9"

# 持久化保存
volumes:
  mysql:

services:
  mysql:
    image: mysql:8.0.34       # 如果您是 arm64 架构：如 MacOS 的 M1，请修改镜像为 image: mysql/mysql-server:8.0.21
    container_name: mysql-server # 容器名
    restart: always
    ports:
      - "3306:3306"  # host物理直接映射端口
    environment:
      MYSQL_ROOT_PASSWORD: 'mysql7914' # root管理员用户密码
      MYSQL_DATABASE: 'blog-veweiyi' # 初始化启动时要创建的数据库的名称
      MYSQL_USER: 'veweiyi' # 初始数据库的访问用户的用户名
      MYSQL_PASSWORD: 'mysql7914' # 初始数据库的访问用户的密码
      TZ: 'Asia/Shanghai'
    volumes:
      - mysql:/var/lib/mysql
    command:
      --character-set-server=utf8mb4
      --collation-server=utf8mb4_general_ci #设置utf8字符集
```
2. 运行以下命令启动Mysql
```shell
docker-compose -f docker-compose.yaml up -d
```
