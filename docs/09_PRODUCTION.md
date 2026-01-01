# 🚀 生产部署指南

**最后更新**: 2026年1月1日 | **版本**: v2.0

---

## 📖 目录

1. [预部署检查](#预部署检查)
2. [部署步骤](#部署步骤)
3. [配置管理](#配置管理)
4. [监控告警](#监控告警)
5. [故障处理](#故障处理)
6. [备份恢复](#备份恢复)

---

## ✅ 预部署检查

### 系统要求

- [ ] 服务器运行环境 (Linux/Unix)
- [ ] Go 1.21+ 已安装
- [ ] MySQL 8.0+ (主从+备份)
- [ ] Redis 7.0+ (集群模式)
- [ ] Elasticsearch 8.5+ (3节点+)
- [ ] Nginx (负载均衡)
- [ ] Docker & Docker Compose

### 安全检查

- [ ] 所有默认密码已修改
- [ ] SSL/TLS证书已部署
- [ ] 防火墙规则已配置
- [ ] 敏感数据已加密
- [ ] 日志审计已启用
- [ ] 备份策略已配置
- [ ] 应急响应计划已制定

---

## 🔧 部署步骤

### 1. 编译构建

```bash
# 编译
go build -o gohub main.go

# 生成可执行文件
make build

# 版本信息
./gohub version
```

### 2. 环境配置

```bash
# 复制环境配置
cp .env.production .env

# 修改关键配置
# - 数据库连接
# - Redis连接
# - Elasticsearch连接
# - JWT密钥
# - 邮件配置
```

### 3. 数据库初始化

```bash
# 创建数据库
mysql -h db-host -u root -p < database/create.sql

# 运行迁移
./gohub migrate

# 导入初始数据
./gohub seed
```

### 4. Elasticsearch初始化

```bash
# 创建索引
./gohub elasticsearch init

# 同步数据
./gohub elasticsearch sync

# 验证
curl http://localhost:9200/_cat/indices
```

### 5. 启动服务

```bash
# 方式1: 直接运行
./gohub serve

# 方式2: 使用systemd
sudo systemctl start gohub

# 方式3: 使用supervisor
supervisord -c /etc/supervisor/conf.d/gohub.conf

# 验证
curl http://localhost:8080/api/health
```

---

## 🔐 配置管理

### 环境变量

```bash
# 应用配置
APP_NAME=GoHub
APP_ENV=production
APP_DEBUG=false
APP_PORT=8080

# 数据库
DB_HOST=db.example.com
DB_PORT=3306
DB_USER=gohub_app
DB_PASSWORD=secure_password
DB_NAME=gohub

# Redis
REDIS_HOST=redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=redis_password

# Elasticsearch
ES_HOST=http://es.example.com:9200

# JWT
JWT_SECRET=your_super_secret_key_min_32_chars
JWT_EXPIRES_IN=24

# 邮件
MAIL_HOST=smtp.mailtrap.io
MAIL_FROM=noreply@example.com
```

### TLS证书配置

```bash
# 生成自签名证书（测试）
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365

# 配置Nginx反向代理HTTPS
server {
    listen 443 ssl;
    server_name example.com;
    
    ssl_certificate /etc/nginx/certs/cert.pem;
    ssl_certificate_key /etc/nginx/certs/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
    }
}
```

---

## 📊 监控告警

### 关键监控指标

```
应用层:
- HTTP请求延迟 (P50, P95, P99)
- 错误率 (5xx响应)
- 吞吐量 (QPS)

数据库:
- 查询延迟
- 连接数
- 慢查询日志

缓存:
- 命中率
- 大小
- 过期率

系统:
- CPU使用率
- 内存使用率
- 磁盘使用率
```

### Prometheus配置示例

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'gohub'
    static_configs:
      - targets: ['localhost:8080']
```

---

## 🚨 故障处理

### 常见问题解决

| 问题 | 症状 | 解决方案 |
|------|------|---------|
| **数据库连接失败** | 5xx错误 | 检查DB连接、重启MySQL |
| **Redis不可用** | 慢响应 | 检查Redis、清理内存 |
| **Elasticsearch宕机** | 搜索失败 | 检查ES集群、自动降级 |
| **磁盘满** | 写入失败 | 清理日志、扩展磁盘 |
| **内存泄漏** | OOM错误 | 重启应用、分析内存 |

---

## 💾 备份恢复

### 备份策略

```bash
# 数据库备份（每天）
mysqldump -h host -u user -p dbname > backup.sql

# Redis备份（每6小时）
redis-cli BGSAVE

# Elasticsearch备份（每周）
# 使用快照功能备份索引
```

### 恢复步骤

```bash
# 1. MySQL恢复
mysql -h host -u user -p dbname < backup.sql

# 2. Redis恢复
redis-cli SHUTDOWN
# 恢复dump.rdb
redis-server

# 3. Elasticsearch恢复
# 从快照恢复索引
```

---

**部署版本**: v2.0  
**可靠性**: 99.9% SLA  
**最后更新**: 2026年1月1日  
*由GoHub Operations Team维护* ✨
