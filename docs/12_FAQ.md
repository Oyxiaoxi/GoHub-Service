# ❓ 常见问题解答 (FAQ)

**最后更新**: 2026年1月1日 | **版本**: v2.0

---

## 🚀 快速启动问题

### Q1: 如何快速启动项目？

**A:** 使用Docker Compose一键启动：

```bash
docker-compose -f docker-compose.elasticsearch.yml up -d
make init
make serve
```

详见 [01_QUICKSTART.md](01_QUICKSTART.md)

### Q2: 启动时出现MySQL连接失败？

**A:** 检查以下项：

```bash
# 1. 检查MySQL状态
docker-compose ps | grep mysql

# 2. 查看MySQL日志
docker-compose logs mysql

# 3. 重启MySQL
docker-compose restart mysql

# 4. 检查连接字符串
echo $DB_HOST $DB_PORT $DB_USER
```

### Q3: go mod tidy报错？

**A:** 尝试以下步骤：

```bash
# 1. 清理缓存
go clean -modcache

# 2. 重新下载
go mod download

# 3. 整理模块
go mod tidy

# 4. 验证
go mod verify
```

---

## 🗄️ 数据库问题

### Q4: 如何重置数据库？

**A:**

```bash
# 1. 删除所有表
./gohub migrate:reset

# 2. 重新运行迁移
./gohub migrate

# 3. 导入种子数据
./gohub seed
```

### Q5: 如何查看慢查询？

**A:**

```bash
# 1. 查看慢查询日志
tail -f /var/log/mysql/slow.log

# 2. 分析查询
mysql> EXPLAIN SELECT ...;

# 3. 添加索引
ALTER TABLE topics ADD INDEX idx_name (column);
```

### Q6: 数据库占用空间太大？

**A:**

```bash
# 1. 检查表大小
SELECT table_name, ROUND(data_length/1024/1024) as MB 
FROM information_schema.tables 
WHERE table_schema='gohub';

# 2. 清理日志表
DELETE FROM logs WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

# 3. 优化表
OPTIMIZE TABLE topics, comments, users;
```

---

## 🔍 搜索问题

### Q7: Elasticsearch搜索无结果？

**A:**

```bash
# 1. 检查ES状态
curl http://localhost:9200/_cluster/health

# 2. 检查索引
curl http://localhost:9200/_cat/indices

# 3. 同步数据
./gohub elasticsearch sync

# 4. 验证数据
curl "http://localhost:9200/topics/_search?pretty"
```

### Q8: 搜索响应慢？

**A:**

```bash
# 1. 检查索引大小
curl "http://localhost:9200/_cat/indices?pretty"

# 2. 查看集群状态
curl http://localhost:9200/_cluster/stats

# 3. 重建索引
./gohub elasticsearch reindex
```

---

## 🔐 权限和认证问题

### Q9: 用户无法登录？

**A:**

```bash
# 1. 检查用户是否存在
SELECT * FROM users WHERE email='user@example.com';

# 2. 检查JWT配置
echo $JWT_SECRET

# 3. 查看认证日志
grep "auth" logs/error.log

# 4. 重置密码
# 在数据库中更新或使用API重置密码
```

### Q10: Token过期？

**A:**

```bash
# 1. 检查Token过期时间
# 默认24小时，可在config/jwt.go修改

# 2. 调用刷新接口
POST /api/auth/refresh
Authorization: Bearer {current_token}

# 3. 重新登录获取新Token
POST /api/auth/login
```

### Q11: 某个用户无法执行某项操作？

**A:**

```bash
# 1. 检查用户权限
SELECT p.name FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN roles r ON rp.role_id = r.id
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = {user_id};

# 2. 为用户分配角色
INSERT INTO user_roles (user_id, role_id) VALUES ({user_id}, {role_id});

# 3. 为角色分配权限
INSERT INTO role_permissions (role_id, permission_id) VALUES ({role_id}, {perm_id});
```

---

## ⚡ 性能问题

### Q12: API响应很慢？

**A:**

```bash
# 1. 检查数据库查询
echo $ENABLE_SLOW_LOG  # 应为true

# 2. 分析慢查询
tail -f /var/log/mysql/slow.log

# 3. 检查缓存命中率
redis-cli INFO stats | grep hits

# 4. 使用pprof分析
curl http://localhost:8080/debug/pprof/profile > cpu.prof
go tool pprof cpu.prof
```

### Q13: 内存占用过高？

**A:**

```bash
# 1. 查看内存使用
docker stats gohub

# 2. 分析内存泄漏
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# 3. 重启应用
docker-compose restart gohub
```

### Q14: Redis占用空间过大？

**A:**

```bash
# 1. 查看Redis大小
redis-cli INFO memory

# 2. 查看键的大小
redis-cli --bigkeys

# 3. 清理过期键
redis-cli FLUSHDB  # 谨慎使用！

# 4. 配置过期策略
redis-cli CONFIG SET maxmemory-policy allkeys-lru
```

---

## 📊 监控和日志

### Q15: 如何查看应用日志？

**A:**

```bash
# 1. Docker容器日志
docker-compose logs -f gohub

# 2. 文件日志
tail -f logs/error.log
tail -f logs/info.log

# 3. 过滤特定日志
grep "ERROR" logs/error.log | tail -20
grep "user_id=123" logs/info.log
```

### Q16: 如何设置日志级别？

**A:**

编辑 `.env`：

```bash
APP_DEBUG=true       # 调试模式，输出详细日志
LOG_LEVEL=debug      # debug, info, warning, error
```

---

## 🐛 调试技巧

### Q17: 如何在本地调试？

**A:**

```bash
# 1. 使用Delve调试器
dlv debug main.go

# 2. IDE调试 (VS Code/GoLand)
# 设置断点并F5开始调试

# 3. 添加日志
logger.Debugf("variable value: %v", var)

# 4. 打印HTTP请求
curl -v http://localhost:8080/api/users
```

### Q18: 如何追踪请求？

**A:**

```bash
# 1. 添加请求ID
# 框架自动在每个请求中生成唯一ID

# 2. 查看完整日志链路
grep "request_id=abc123" logs/*.log

# 3. 使用分布式追踪
# Jaeger / Zipkin 集成
```

---

## 🔄 版本更新

### Q19: 如何升级到新版本？

**A:**

```bash
# 1. 备份数据
./scripts/backup-database.sh

# 2. 拉取新代码
git pull origin main

# 3. 更新依赖
go mod download && go mod tidy

# 4. 运行迁移
./gohub migrate

# 5. 重启服务
docker-compose restart gohub
```

### Q20: 如何回滚版本？

**A:**

```bash
# 1. 查看提交历史
git log --oneline

# 2. 回到之前版本
git checkout <commit_hash>

# 3. 恢复数据库备份
mysql -u root -p gohub < backup.sql

# 4. 重启应用
docker-compose restart gohub
```

---

## 📞 获取更多帮助

### 相关文档

- [快速开始](01_QUICKSTART.md) - 初步了解
- [架构设计](02_ARCHITECTURE.md) - 深入理解
- [开发指南](05_DEVELOPMENT.md) - 开发相关
- [性能优化](07_PERFORMANCE.md) - 性能调优
- [部署指南](09_PRODUCTION.md) - 生产部署
