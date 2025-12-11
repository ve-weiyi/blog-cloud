# HTTPS访问配置

本文档介绍如何为博客系统配置 HTTPS 证书，实现安全的加密传输。

## HTTPS 简介

HTTPS（HyperText Transfer Protocol Secure）是 HTTP 的安全版本，通过 SSL/TLS 协议加密数据传输。

### 为什么需要 HTTPS

- ✅ 数据加密传输，防止中间人攻击
- ✅ 验证网站身份，防止钓鱼网站
- ✅ 提升 SEO 排名，搜索引擎优先收录
- ✅ 浏览器不显示"不安全"警告
- ✅ 支持 HTTP/2 协议，提升性能

## SSL 证书申请与配置

SSL 证书用于启用 HTTPS 加密传输，保护网站数据安全，避免浏览器显示"不安全"警告。

### 免费证书

| 类型            | 有效期 | 适用场景    | 获取方式   |
|---------------|-----|---------|--------|
| Let's Encrypt | 90天 | 个人网站    | 自动续期   |
| 阿里云免费证书       | 1年  | 个人/小型网站 | 云服务商申请 |
| 腾讯云免费证书       | 1年  | 个人/小型网站 | 云服务商申请 |

### 付费证书

| 类型    | 价格 | 适用场景   | 特点              |
|-------|----|--------|-----------------|
| DV 证书 | 低  | 个人网站   | 域名验证            |
| OV 证书 | 中  | 企业网站   | 组织验证            |
| EV 证书 | 高  | 金融/电商  | 扩展验证，绿色地址栏      |
| 通配符证书 | 较高 | 多子域名网站 | 支持 *.domain.com |

## 申请证书

### 方式一：腾讯云免费证书（推荐）

#### 1. 进入 SSL 证书控制台

访问 [腾讯云 SSL 证书控制台](https://console.cloud.tencent.com/ssl)

![SSL证书控制台](https://static.veweiyi.cn/article/deploy/img_3.png)

#### 2. 选择证书类型

根据需求选择免费或付费证书

![选择证书类型](https://static.veweiyi.cn/article/deploy/img_4.png)

#### 3. 填写证书信息

- 域名：填写你的域名（如 blog.veweiyi.cn）
- 验证方式：选择 DNS 验证或文件验证
- 申请邮箱：填写有效邮箱

![填写证书信息](https://static.veweiyi.cn/article/deploy/img_5.png)

#### 4. 域名验证

**DNS 验证步骤**
```
1. 登录域名解析控制台
2. 添加 TXT 记录：
   主机记录：_dnsauth
   记录类型：TXT
   记录值：（证书系统提供的验证值）
   TTL：600
3. 等待验证通过
```

#### 5. 下载证书

验证通过后，下载证书文件

![下载证书](https://static.veweiyi.cn/article/deploy/img_7.png)

### 方式二：Let's Encrypt 免费证书

#### 使用 Certbot 自动申请

```bash
# 安装 Certbot
yum install -y certbot python3-certbot-nginx  # CentOS
apt install -y certbot python3-certbot-nginx  # Ubuntu

# 申请证书（自动配置 Nginx）
certbot --nginx -d blog.veweiyi.cn -d www.blog.veweiyi.cn

# 仅申请证书（手动配置）
certbot certonly --webroot -w /var/www/html -d blog.veweiyi.cn

# 测试自动续期
certbot renew --dry-run
```

## 安装证书

### 1. 上传证书文件

```bash
# 创建证书目录
mkdir -p /www/ssl

# 上传证书文件
scp blog.veweiyi.cn_bundle.crt root@your_server:/www/ssl/
scp blog.veweiyi.cn.key root@your_server:/www/ssl/

# 设置权限
chmod 600 /www/ssl/*.key
chmod 644 /www/ssl/*.crt
```

证书文件目录结构：

![证书文件](https://static.veweiyi.cn/article/deploy/img_8.png)

### 2. 配置 Nginx

创建或编辑 Nginx 配置文件：

```nginx
# HTTP 自动跳转 HTTPS
server {
    listen 80;
    server_name blog.veweiyi.cn;
    return 301 https://$server_name$request_uri;
}

# HTTPS 配置
server {
    listen 443 ssl http2;
    server_name blog.veweiyi.cn;

    # SSL 证书配置
    ssl_certificate /www/ssl/blog.veweiyi.cn_bundle.crt;
    ssl_certificate_key /www/ssl/blog.veweiyi.cn.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 安全头
    add_header Strict-Transport-Security "max-age=31536000" always;

    # 博客前台接口
    location ^~ /api {
        proxy_pass http://localhost:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 博客前台网站
    location / {
        proxy_pass http://localhost:9420;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. 测试并重启 Nginx

```bash
# 测试配置
nginx -t

# 重新加载配置
nginx -s reload

# 或重启 Nginx
systemctl restart nginx
```

## 证书验证

### 浏览器测试

1. 访问 https://blog.veweiyi.cn
2. 检查地址栏是否显示锁图标
3. 点击锁图标查看证书信息

### 命令行测试

```bash
# 测试 SSL 连接
openssl s_client -connect blog.veweiyi.cn:443 -servername blog.veweiyi.cn

# 查看证书信息
openssl s_client -connect blog.veweiyi.cn:443 -servername blog.veweiyi.cn | openssl x509 -noout -text

# 测试证书有效期
echo | openssl s_client -connect blog.veweiyi.cn:443 2>/dev/null | openssl x509 -noout -dates

# 测试 HTTPS 访问
curl -I https://blog.veweiyi.cn
```

### 在线测试工具

- [SSL Labs](https://www.ssllabs.com/ssltest/) - 权威 SSL 测试
- [MySSL](https://myssl.com/) - 国内 SSL 检测

## 证书续期

### Let's Encrypt 自动续期

```bash
# 查看定时任务
systemctl list-timers | grep certbot

# 手动续期
certbot renew

# 强制续期
certbot renew --force-renewal
```

### 云服务商证书续期

```
1. 证书到期前 30 天收到提醒
2. 重新申请新证书
3. 下载新证书文件
4. 替换服务器上的证书
5. 重新加载 Nginx
```

## 常见问题

### 1. 证书不受信任

**原因**：

- 证书链不完整
- 证书已过期
- 域名不匹配

**解决方案**：

```bash
# 检查证书链
openssl s_client -connect blog.veweiyi.cn:443 -showcerts

# 使用完整证书链
ssl_certificate /www/ssl/blog.veweiyi.cn_bundle.crt;

# 检查证书有效期
openssl x509 -in /www/ssl/blog.veweiyi.cn.crt -noout -dates
```

### 2. 混合内容警告

**原因**：HTTPS 页面加载了 HTTP 资源

**解决方案**：

```nginx
# 添加 CSP 头，自动升级不安全请求
add_header Content-Security-Policy "upgrade-insecure-requests" always;
```

### 3. 性能优化

```nginx
# 启用 SSL 会话复用
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;

# 启用 HTTP/2
listen 443 ssl http2;
```

## 安全建议

- ✅ 定期更新证书
- ✅ 使用强加密算法
- ✅ 启用 HSTS
- ✅ 禁用不安全的协议（SSLv3, TLSv1.0, TLSv1.1）
- ✅ 配置证书过期监控
- ✅ 使用 CAA 记录保护域名
- ✅ 定期检查证书配置

## 参考资料

- [Let's Encrypt 官网](https://letsencrypt.org/)
- [SSL Labs 测试](https://www.ssllabs.com/ssltest/)
- [Mozilla SSL 配置生成器](https://ssl-config.mozilla.org/)
- [腾讯云 SSL 证书](https://cloud.tencent.com/product/ssl)
