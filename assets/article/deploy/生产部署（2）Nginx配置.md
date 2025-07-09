# Nginx 反向代理配置

本文档介绍如何配置 Nginx 作为反向代理，实现域名访问和负载均衡。

## Nginx 简介

Nginx 是一个高性能的 HTTP 和反向代理服务器，具有以下优势：

- ✅ 高并发处理能力
- ✅ 低内存消耗
- ✅ 负载均衡
- ✅ 静态资源缓存
- ✅ SSL/TLS 支持

## 安装 Nginx

### Docker 方式（推荐）

```bash
docker run -d \
  --name nginx \
  --restart always \
  -p 80:80 \
  -p 443:443 \
  -v /etc/nginx/nginx.conf:/etc/nginx/nginx.conf:ro \
  -v /www/ssl:/www/ssl:ro \
  nginx:latest
```

### 系统安装

**CentOS/RHEL**

```bash
yum install -y nginx
systemctl start nginx
systemctl enable nginx
```

**Ubuntu/Debian**

```bash
apt update
apt install -y nginx
systemctl start nginx
systemctl enable nginx
```

## 配置文件

### 完整配置示例

```bash
# 复制项目配置文件
git clone https://github.com/ve-weiyi/ve-blog-golang.git
sudo cp ve-blog-golang/deploy/docker-compose/nginx.conf /etc/nginx/nginx.conf

# 重新加载配置
sudo nginx -t
sudo nginx -s reload
```

### 基础配置结构

```nginx
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # 日志格式
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    # 性能优化
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript 
               application/json application/javascript application/xml+rss;

    # WebSocket 支持
    map $http_upgrade $connection_upgrade {
        default upgrade;
        '' close;
    }

    # 引入站点配置
    include /etc/nginx/conf.d/*.conf;
}
```

## 站点配置

### 博客前台配置

创建 `/etc/nginx/conf.d/blog.conf`：

```nginx
server {
    listen 80;
    server_name blog.veweiyi.cn;

    # HTTP 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name blog.veweiyi.cn;

    # SSL 证书配置
    ssl_certificate /www/ssl/blog.veweiyi.cn.crt;
    ssl_certificate_key /www/ssl/blog.veweiyi.cn.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 前台 API 接口
    location ^~ /api {
        proxy_pass http://localhost:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket 支持
    location ^~ /ws {
        proxy_pass http://localhost:9090;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 前台网站
    location / {
        proxy_pass http://localhost:9420;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 管理后台配置

创建 `/etc/nginx/conf.d/admin.conf`：

```nginx
server {
    listen 80;
    server_name admin.veweiyi.cn;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name admin.veweiyi.cn;

    # SSL 证书配置
    ssl_certificate /www/ssl/admin.veweiyi.cn.crt;
    ssl_certificate_key /www/ssl/admin.veweiyi.cn.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 后台 API 接口
    location ^~ /admin_api {
        proxy_pass http://localhost:9091;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 后台网站
    location / {
        proxy_pass http://localhost:9421;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 性能优化

### 缓存配置

```nginx
# 静态资源缓存
location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
    expires 30d;
    add_header Cache-Control "public, immutable";
}

# 代理缓存
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m max_size=1g inactive=60m;

location / {
    proxy_cache my_cache;
    proxy_cache_valid 200 60m;
    proxy_cache_key $scheme$proxy_host$request_uri;
    add_header X-Cache-Status $upstream_cache_status;
}
```

### 限流配置

```nginx
# 限制请求频率
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

location /api {
    limit_req zone=api_limit burst=20 nodelay;
    proxy_pass http://localhost:9090;
}
```

### 负载均衡

```nginx
upstream blog_backend {
    least_conn;
    server localhost:9090 weight=1 max_fails=3 fail_timeout=30s;
    server localhost:9091 weight=1 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

location /api {
    proxy_pass http://blog_backend;
}
```

## 安全配置

### 基础安全

```nginx
# 隐藏 Nginx 版本
server_tokens off;

# 防止点击劫持
add_header X-Frame-Options "SAMEORIGIN" always;

# XSS 防护
add_header X-XSS-Protection "1; mode=block" always;

# 内容类型嗅探
add_header X-Content-Type-Options "nosniff" always;

# HSTS
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

### IP 访问限制

```nginx
# 白名单
location /admin {
    allow 192.168.1.0/24;
    deny all;
}

# 黑名单
deny 192.168.1.100;
```

## 常用命令

```bash
# 测试配置文件
nginx -t

# 重新加载配置
nginx -s reload

# 停止服务
nginx -s stop

# 优雅停止
nginx -s quit

# 查看版本
nginx -v

# 查看编译参数
nginx -V
```

## 日志管理

### 日志配置

```nginx
# 访问日志
access_log /var/log/nginx/blog_access.log main;

# 错误日志
error_log /var/log/nginx/blog_error.log warn;

# 关闭日志
access_log off;
```

### 日志分析

```bash
# 查看访问量
cat /var/log/nginx/access.log | wc -l

# 查看 IP 访问排名
awk '{print $1}' /var/log/nginx/access.log | sort | uniq -c | sort -rn | head -10

# 查看访问最多的页面
awk '{print $7}' /var/log/nginx/access.log | sort | uniq -c | sort -rn | head -10

# 查看状态码统计
awk '{print $9}' /var/log/nginx/access.log | sort | uniq -c | sort -rn
```

## 故障排查

### 常见问题

**配置文件错误**

```bash
# 检查配置语法
nginx -t

# 查看错误日志
tail -f /var/log/nginx/error.log
```

**502 Bad Gateway**

- 检查后端服务是否运行
- 检查端口是否正确
- 检查防火墙规则

**504 Gateway Timeout**

```nginx
# 增加超时时间
proxy_connect_timeout 60s;
proxy_send_timeout 60s;
proxy_read_timeout 60s;
```

**413 Request Entity Too Large**

```nginx
# 增加上传大小限制
client_max_body_size 100M;
```

## 参考资料

- [Nginx 官方文档](https://nginx.org/en/docs/)
- [Nginx 配置详解](https://cloud.tencent.com/developer/article/2412917)
- [Nginx 性能优化](https://www.nginx.com/blog/tuning-nginx/)
