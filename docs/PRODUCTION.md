# 生产环境配置指南

## 目录
- [环境变量配置](#环境变量配置)
- [数据库配置](#数据库配置)
- [Redis 配置](#redis-配置)
- [安全配置](#安全配置)
- [性能优化](#性能优化)
- [日志配置](#日志配置)
- [监控配置](#监控配置)
- [部署检查清单](#部署检查清单)

---

## 环境变量配置

### 基础配置

```env
# 应用配置
APP_NAME=GoHub-Service
APP_ENV=production                    # ⚠️ 必须设置为 production
APP_KEY=your-32-character-secret-key  # ⚠️ 必须使用强随机字符串
APP_DEBUG=false                       # ⚠️ 生产环境必须关闭调试
APP_URL=https://api.yourdomain.com
APP_PORT=3000
TIMEZONE=Asia/Shanghai
```

**注意事项：**
- `APP_KEY` 必须是 32 位以上的强随机字符串，用于加密会话和 JWT
- `APP_DEBUG` 在生产环境必须设置为 `false`，否则会暴露敏感信息
- `APP_ENV` 设置为 `production` 会启用性能优化和错误处理优化

---

## 数据库配置

### MySQL 生产配置

```env
DB_CONNECTION=mysql
DB_HOST=your-db-host.com              # 使用内网地址
DB_PORT=3306
DB_DATABASE=gohub_production
DB_USERNAME=gohub_user                # ⚠️ 不要使用 root
DB_PASSWORD=strong-password-here      # ⚠️ 使用强密码
DB_DEBUG=0                            # ⚠️ 生产环境关闭 SQL 日志
```

### 数据库优化建议

**1. 连接池配置**
```go
// 在 config/database.go 中设置
MaxIdleConns: 10     // 最大空闲连接数
MaxOpenConns: 100    // 最大打开连接数
ConnMaxLifetime: 1h  // 连接最大生命周期
```

**2. 索引优化**
确保以下字段已创建索引：
- `users` 表：`email`, `phone`, `name`
- `topics` 表：`category_id`, `user_id`, `created_at`
- `comments` 表：`topic_id`, `user_id`, `parent_id`
- `follows` 表：`user_id`, `follow_id`（联合索引）
- `likes` 表：`user_id`, `target_type`, `target_id`（联合索引）

**3. 查询优化**
```sql
-- 检查慢查询
SHOW VARIABLES LIKE 'slow_query_log';
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;

-- 分析表
ANALYZE TABLE users, topics, comments;
```

---

## Redis 配置

### 生产配置

```env
REDIS_HOST=your-redis-host.com        # 使用内网地址
REDIS_PORT=6379
REDIS_PASSWORD=strong-redis-password  # ⚠️ 必须设置密码
REDIS_CACHE_DB=0
REDIS_MAIN_DB=1
```

### Redis 优化建议

**1. 内存配置**
```redis
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
```

**2. 持久化配置**
```redis
# 生产环境建议同时启用 RDB 和 AOF
save 900 1
save 300 10
save 60 10000

appendonly yes
appendfsync everysec
```

**3. 连接池配置**
```go
PoolSize: 10
MinIdleConns: 5
MaxRetries: 3
```

---

## 安全配置

### 1. CORS 配置

```env
# 生产环境只允许特定域名
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

**配置文件：** `config/cors.go`
```go
"allowed_origins": config.Env("CORS_ALLOWED_ORIGINS", ""),
"allowed_methods": "GET,POST,PUT,PATCH,DELETE,OPTIONS",
"allowed_headers": "Authorization,Content-Type,Accept,Origin",
"exposed_headers": "Content-Length",
"max_age": 12 * 3600,
```

### 2. JWT 配置

```env
JWT_EXPIRE_TIME=7200              # 2小时（秒）
JWT_MAX_REFRESH_TIME=604800       # 7天（秒）
JWT_KEY=your-jwt-secret-key       # ⚠️ 必须是强随机字符串
```

### 3. 限流配置

```env
LIMITER_DEFAULT_MESSAGE=请求过于频繁，请稍后再试

# IP 限流（每小时）
LIMIT_GLOBAL_IP_RATE=1000-H       # 全局限流
LIMIT_AUTH_IP_RATE=5000-H         # 认证用户限流

# 接口限流（每小时）
LIMIT_SIGNUP_EXIST_RATE=60-H      # 注册检查
LIMIT_VERIFY_PHONE_RATE=10-H      # 手机验证码
LIMIT_VERIFY_EMAIL_RATE=10-H      # 邮箱验证码
LIMIT_VERIFY_CAPTCHA_RATE=20-H    # 图形验证码
```

### 4. 内容安全配置

```env
CONTENT_CHECK_ENABLED=true
SENSITIVE_WORD_FILTER_ENABLED=true
SENSITIVE_WORD_REPLACEMENT=***
XSS_PROTECTION_ENABLED=true
MAX_CONTENT_LENGTH=10000
MAX_TITLE_LENGTH=200
ALLOW_HTML_TAGS=false
IMAGE_CHECK_ENABLED=true
MAX_IMAGE_SIZE=5242880            # 5MB
AUDIT_LOG_ENABLED=true
```

### 5. HTTPS/TLS 配置

**使用 Nginx 反向代理（推荐）：**

```nginx
server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL 证书
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # 反向代理
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 限制请求体大小
    client_max_body_size 10M;
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}
```

---

## 性能优化

### 1. 文件上传配置

```env
STORAGE_BASE_PATH=/var/www/uploads    # 使用绝对路径
STORAGE_PUBLIC_PREFIX=/uploads
STORAGE_MAX_SIZE_MB=5
STORAGE_MAX_IMAGE_WIDTH=2048
STORAGE_JPEG_QUALITY=85
STORAGE_ALLOWED_EXT=jpg,jpeg,png,gif,webp
```

### 2. 缓存策略

**缓存配置建议：**
```go
// 热点数据缓存时间
UserCache: 1小时
TopicCache: 30分钟
CategoryCache: 24小时
CommentCache: 15分钟
```

### 3. 数据库查询优化

- 使用预加载（Preload）减少 N+1 查询
- 使用分页避免大量数据加载
- 使用索引加速查询
- 避免 SELECT * 查询

### 4. 静态资源优化

- 使用 CDN 托管静态资源
- 启用 Gzip 压缩
- 设置合理的缓存头

---

## 日志配置

### 生产环境日志配置

```env
LOG_TYPE=daily                        # 每日日志文件
LOG_LEVEL=info                        # ⚠️ 不要使用 debug
LOG_MAX_SIZE=100                      # 日志文件最大 100MB
LOG_MAX_BACKUPS=30                    # 保留 30 个备份
LOG_MAX_AGE=30                        # 保留 30 天
LOG_COMPRESS=true                     # 压缩旧日志
```

### 日志级别说明

- `debug`: 详细的调试信息（仅开发环境）
- `info`: 一般信息（生产环境推荐）
- `warn`: 警告信息
- `error`: 错误信息
- `fatal`: 致命错误

### 日志轮转配置

```go
// pkg/logger/logger.go
&lumberjack.Logger{
    Filename:   logFilename,
    MaxSize:    config.GetInt("log.max_size"),      // 100MB
    MaxBackups: config.GetInt("log.max_backups"),   // 30个
    MaxAge:     config.GetInt("log.max_age"),       // 30天
    Compress:   config.GetBool("log.compress"),     // 压缩
}
```

---

## 监控配置

### 1. 健康检查端点

在 `routes/api.go` 中添加：

```go
// 健康检查
api.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "healthy",
        "timestamp": time.Now().Unix(),
    })
})

// 就绪检查
api.GET("/ready", func(c *gin.Context) {
    // 检查数据库连接
    if err := database.DB.Ping(); err != nil {
        c.JSON(503, gin.H{"status": "not ready", "reason": "database"})
        return
    }
    
    // 检查 Redis 连接
    if err := redis.Redis.Ping(c).Err(); err != nil {
        c.JSON(503, gin.H{"status": "not ready", "reason": "redis"})
        return
    }
    
    c.JSON(200, gin.H{"status": "ready"})
})
```

### 2. Prometheus 指标

项目已集成 Prometheus 指标，访问：`/metrics`

**关键指标：**
- HTTP 请求总数和延迟
- 数据库连接池状态
- Redis 连接池状态
- Goroutine 数量
- 内存使用情况

### 3. 错误监控

建议集成以下服务之一：
- Sentry（错误追踪）
- New Relic（应用性能监控）
- DataDog（全栈监控）

---

## 部署检查清单

### 部署前检查

- [ ] **环境变量配置完整**
  - [ ] `APP_ENV=production`
  - [ ] `APP_DEBUG=false`
  - [ ] `APP_KEY` 已设置为强随机值
  - [ ] 数据库配置正确
  - [ ] Redis 配置正确
  - [ ] JWT 密钥已配置

- [ ] **安全配置**
  - [ ] CORS 只允许特定域名
  - [ ] Redis 设置了密码
  - [ ] 数据库不使用 root 用户
  - [ ] 所有密码为强密码
  - [ ] 限流配置合理
  - [ ] 内容安全检查已启用

- [ ] **数据库**
  - [ ] 运行了所有迁移：`./main migrate up`
  - [ ] 创建了必要的索引
  - [ ] 配置了备份策略
  - [ ] 设置了连接池参数

- [ ] **文件权限**
  - [ ] storage 目录可写：`chmod -R 755 storage/`
  - [ ] public/uploads 目录可写：`chmod -R 755 public/uploads/`
  - [ ] 日志目录可写：`chmod -R 755 storage/logs/`

- [ ] **依赖和构建**
  - [ ] 运行了 `go mod tidy`
  - [ ] 编译成功：`go build -o main`
  - [ ] 测试通过：`go test ./...`

- [ ] **反向代理**
  - [ ] Nginx/Caddy 配置正确
  - [ ] SSL 证书已安装
  - [ ] HTTPS 重定向已配置
  - [ ] 安全头已添加

- [ ] **监控和日志**
  - [ ] 日志级别设置为 `info`
  - [ ] 日志轮转已配置
  - [ ] 健康检查端点可访问
  - [ ] Prometheus 指标可访问

### 部署后验证

- [ ] **功能测试**
  - [ ] 用户注册/登录正常
  - [ ] JWT Token 正常
  - [ ] 文件上传正常
  - [ ] 数据库读写正常
  - [ ] Redis 缓存正常

- [ ] **性能测试**
  - [ ] 响应时间在可接受范围
  - [ ] 并发请求处理正常
  - [ ] 内存使用正常
  - [ ] CPU 使用正常

- [ ] **安全测试**
  - [ ] XSS 防护生效
  - [ ] 敏感词过滤生效
  - [ ] 限流功能正常
  - [ ] CORS 配置生效
  - [ ] HTTPS 正常工作

- [ ] **监控告警**
  - [ ] 健康检查正常
  - [ ] 错误日志正常记录
  - [ ] 监控指标正常采集
  - [ ] 告警规则已配置

---

## 服务管理

### 使用 Systemd 管理服务

创建服务文件 `/etc/systemd/system/gohub.service`：

```ini
[Unit]
Description=GoHub Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/gohub
Environment="GIN_MODE=release"
ExecStart=/var/www/gohub/main
Restart=always
RestartSec=10

# 安全配置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ReadWritePaths=/var/www/gohub/storage /var/www/gohub/public/uploads

[Install]
WantedBy=multi-user.target
```

**启动服务：**
```bash
sudo systemctl daemon-reload
sudo systemctl enable gohub
sudo systemctl start gohub
sudo systemctl status gohub
```

---

## Docker 部署

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/public ./public
COPY --from=builder /app/storage ./storage

ENV TZ=Asia/Shanghai
EXPOSE 3000

CMD ["./main"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - APP_ENV=production
      - DB_HOST=mysql
      - REDIS_HOST=redis
    env_file:
      - .env
    depends_on:
      - mysql
      - redis
    volumes:
      - ./storage:/root/storage
      - ./public/uploads:/root/public/uploads
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_DATABASE: ${DB_DATABASE}
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - app
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
```

---

## 备份策略

### 数据库备份

```bash
#!/bin/bash
# backup-db.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/mysql"
DB_NAME="gohub_production"

mkdir -p $BACKUP_DIR

mysqldump -u$DB_USER -p$DB_PASSWORD \
  --single-transaction \
  --routines \
  --triggers \
  $DB_NAME | gzip > $BACKUP_DIR/backup_$DATE.sql.gz

# 删除 30 天前的备份
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete
```

### 文件备份

```bash
#!/bin/bash
# backup-files.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/files"

mkdir -p $BACKUP_DIR

tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz \
  /var/www/gohub/public/uploads

# 删除 30 天前的备份
find $BACKUP_DIR -name "uploads_*.tar.gz" -mtime +30 -delete
```

### 自动备份（Crontab）

```cron
# 每天凌晨 2 点备份数据库
0 2 * * * /path/to/backup-db.sh

# 每天凌晨 3 点备份文件
0 3 * * * /path/to/backup-files.sh
```

---

## 故障排查

### 常见问题

**1. 连接数据库失败**
```bash
# 检查数据库连接
mysql -h$DB_HOST -u$DB_USERNAME -p$DB_PASSWORD $DB_DATABASE

# 检查防火墙
sudo ufw status
```

**2. Redis 连接失败**
```bash
# 检查 Redis 连接
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD ping

# 检查 Redis 状态
sudo systemctl status redis
```

**3. 服务无法启动**
```bash
# 查看日志
sudo journalctl -u gohub -n 50 -f

# 检查端口占用
sudo lsof -i :3000
```

**4. 性能问题**
```bash
# 查看进程状态
top
htop

# 查看慢查询
mysql -e "SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;"

# 查看 Redis 慢查询
redis-cli slowlog get 10
```

---

## 安全最佳实践

1. **定期更新依赖**
   ```bash
   go get -u ./...
   go mod tidy
   ```

2. **定期更新系统**
   ```bash
   sudo apt update && sudo apt upgrade
   ```

3. **使用防火墙**
   ```bash
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw allow 22/tcp
   sudo ufw enable
   ```

4. **限制 SSH 访问**
   - 禁用密码登录，只使用 SSH 密钥
   - 修改 SSH 默认端口
   - 使用 fail2ban 防止暴力破解

5. **定期审计日志**
   ```bash
   # 查看最近的错误日志
   tail -n 100 storage/logs/error.log
   
   # 查看访问日志
   tail -n 100 /var/log/nginx/access.log
   ```

6. **定期检查安全漏洞**
   ```bash
   # Go 依赖安全检查
   go list -json -m all | nancy sleuth
   ```

---

## 性能调优建议

### Go 应用优化

```bash
# 编译优化
go build -ldflags="-s -w" -o main

# 启用 PGO（Profile-Guided Optimization）
go build -pgo=auto -o main
```

### 系统优化

```bash
# 增加文件描述符限制
echo "* soft nofile 65536" >> /etc/security/limits.conf
echo "* hard nofile 65536" >> /etc/security/limits.conf

# TCP 优化
echo "net.core.somaxconn = 1024" >> /etc/sysctl.conf
echo "net.ipv4.tcp_max_syn_backlog = 2048" >> /etc/sysctl.conf
sysctl -p
```

---

## 总结

生产环境部署需要注意：
1. ✅ 关闭调试模式
2. ✅ 使用强密码和密钥
3. ✅ 配置 HTTPS 和安全头
4. ✅ 启用限流和内容安全检查
5. ✅ 配置监控和告警
6. ✅ 设置日志轮转
7. ✅ 配置数据库和文件备份
8. ✅ 使用反向代理和负载均衡
9. ✅ 定期更新和安全审计
10. ✅ 准备故障恢复计划

完整的配置示例请参考 `.env.production.example` 文件。
