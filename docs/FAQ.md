# ❓ 常见问题

常见的开发问题和解决方案。

## 设置和安装

### Q: 如何设置开发环境？

A: 参考 [QUICKSTART.md](./QUICKSTART.md)

关键步骤:
1. 克隆项目: `git clone ...`
2. 安装依赖: `go mod download`
3. 配置文件: 复制 `.env.example` 为 `.env`
4. 运行迁移: `go run main.go migrate`
5. 启动服务: `go run main.go serve`

### Q: 数据库连接失败怎么办？

A: 检查以下事项:

1. 验证数据库配置 (`.env` 文件)
```
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=gohub
DB_USERNAME=root
DB_PASSWORD=password
```

2. 验证数据库是否运行
```bash
# MySQL
mysql -h localhost -u root -p

# SQLite
ls -la database.db
```

3. 创建数据库 (如果使用 MySQL)
```sql
CREATE DATABASE gohub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. 运行迁移
```bash
go run main.go migrate
```

### Q: Redis 连接超时？

A: 检查以下事项:

1. 验证 Redis 是否运行
```bash
redis-cli ping  # 应返回 PONG
```

2. 验证 Redis 配置
```
REDIS_HOST=localhost
REDIS_PORT=6379
```

3. 检查网络连接
```bash
telnet localhost 6379
```

4. 查看 Redis 日志
```bash
redis-server --loglevel verbose
```

## API 使用

### Q: 如何进行认证？

A: 使用 JWT 令牌

1. **登录获取令牌**
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# 响应
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600
}
```

2. **使用令牌访问 API**
```bash
curl -X GET http://localhost:3000/api/v1/topics \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

3. **令牌过期**

当令牌过期时，使用刷新令牌获取新令牌：
```bash
curl -X POST http://localhost:3000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "..."
  }'
```

### Q: 如何获取 API 文档？

A: 使用 Swagger UI

访问: `http://localhost:3000/swagger/index.html`

或查看 JSON 格式: `http://localhost:3000/swagger.json`

### Q: API 返回 403 Forbidden，为什么？

A: 权限检查失败

1. **验证令牌有效性**
```bash
# 解码 JWT
echo "eyJhbGciOiJIUzI1NiIs..." | base64 -d
```

2. **检查用户角色**
```bash
curl -X GET http://localhost:3000/api/v1/user/roles \
  -H "Authorization: Bearer ..."
```

3. **检查权限列表**
```bash
curl -X GET http://localhost:3000/api/v1/user/permissions \
  -H "Authorization: Bearer ..."
```

### Q: 如何批量操作数据？

A: 使用批量 API

```bash
# 批量创建
curl -X POST http://localhost:3000/api/v1/topics/batch \
  -H "Authorization: Bearer ..." \
  -H "Content-Type: application/json" \
  -d '[
    {"title": "Topic 1", "content": "..."},
    {"title": "Topic 2", "content": "..."}
  ]'

# 批量删除
curl -X DELETE http://localhost:3000/api/v1/topics/batch \
  -H "Authorization: Bearer ..." \
  -H "Content-Type: application/json" \
  -d '{"ids": [1, 2, 3]}'
```

## 开发问题

### Q: 如何添加新的 API 端点？

A: 按照以下步骤：

1. **创建 Model**
```go
// app/models/feature/model.go
type Feature struct {
    ID    uint
    Name  string
    // ...
}
```

2. **创建 Repository**
```go
// app/repositories/feature_repository.go
type FeatureRepository struct { /* ... */ }

func (r *FeatureRepository) Create(feature *models.Feature) error { /* ... */ }
```

3. **创建 Service**
```go
// app/services/feature_service.go
type FeatureService struct { /* ... */ }

func (s *FeatureService) Create(feature *models.Feature) error { /* ... */ }
```

4. **创建 Controller**
```go
// app/http/controllers/feature_controller.go
type FeatureController struct { /* ... */ }

func (ctrl *FeatureController) Store(c *gin.Context) { /* ... */ }
```

5. **定义 Routes**
```go
// routes/feature.go
func RegisterFeatureRoutes(r *gin.Engine, ctrl *controllers.FeatureController) {
    r.POST("/api/v1/features", ctrl.Store)
}
```

6. **在主路由注册**
```go
// routes/api.go
RegisterFeatureRoutes(engine, featureCtrl)
```

### Q: 如何进行数据库迁移？

A: 使用迁移命令

```bash
# 创建迁移文件
go run main.go make:migration create_features_table

# 运行所有迁移
go run main.go migrate

# 查看迁移状态
go run main.go migrate:status

# 回滚上一次迁移
go run main.go migrate:rollback
```

### Q: 如何处理关联关系？

A: 使用 GORM 的关联功能

```go
// 一对多
type User struct {
    ID    uint
    Posts []Post `gorm:"foreignKey:UserID"`
}

// 多对多
type Post struct {
    ID       uint
    Tags     []Tag `gorm:"many2many:post_tags"`
}

// 预加载关联数据
db.Preload("Posts").Preload("Posts.Tags").Find(&users)
```

### Q: 如何处理并发请求？

A: 使用 Mutex 或 Channel

```go
// 方法 1: Mutex
var mu sync.Mutex

func UpdateCounter() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}

// 方法 2: Channel
updates := make(chan int)

go func() {
    for update := range updates {
        counter += update
    }
}()

updates <- 1
```

### Q: 如何缓存数据？

A: 使用 Redis 缓存

```go
// 获取缓存
val, err := redisClient.Get(ctx, "key").Result()

// 设置缓存
redisClient.Set(ctx, "key", "value", 1*time.Hour)

// 删除缓存
redisClient.Del(ctx, "key")

// 缓存失效
redisClient.FlushAll(ctx)
```

## 测试问题

### Q: 如何运行测试？

A: 使用 Go test 命令

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./app/services/...

# 显示详细输出
go test -v ./...

# 显示覆盖率
go test -cover ./...

# 生成覆盖率 HTML 报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Q: 如何进行单元测试？

A: 编写测试文件

```go
// app/services/user_service_test.go
package services

import (
    "testing"
)

func TestCreateUser(t *testing.T) {
    // Arrange
    service := NewUserService(mockRepo, mockCache)
    user := &User{Name: "John"}
    
    // Act
    err := service.Create(user)
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
```

### Q: 如何模拟数据库操作？

A: 使用 Mock

```go
import "github.com/stretchr/testify/mock"

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

// 使用
mockRepo := new(MockRepository)
mockRepo.On("Create", mock.Anything).Return(nil)

service := NewUserService(mockRepo, mockCache)
```

## 部署问题

### Q: 如何构建可执行文件？

A: 使用 Go build

```bash
# 构建
go build -o gohub main.go

# 交叉编译 (Linux)
GOOS=linux GOARCH=amd64 go build -o gohub main.go

# 交叉编译 (Windows)
GOOS=windows GOARCH=amd64 go build -o gohub.exe main.go

# 添加版本信息
go build -ldflags="-X main.Version=1.0.0" -o gohub main.go
```

### Q: 如何部署到生产环境？

A: 使用 systemd 服务

```ini
# /etc/systemd/system/gohub.service
[Unit]
Description=GoHub Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/gohub
ExecStart=/var/www/gohub/gohub serve
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl start gohub
sudo systemctl enable gohub
sudo systemctl status gohub
```

### Q: 如何监控应用性能？

A: 使用日志和指标

```bash
# 查看日志
tail -f storage/logs/gohub.log

# 监控资源使用
top -p $(pidof gohub)

# 检查端口占用
lsof -i :8080
```

## 常见错误

### Error: nil pointer dereference

**原因**: 访问了 nil 指针

**解决**:
```go
// ✅ 检查指针
if user != nil {
    fmt.Println(user.Name)
}

// ✅ 初始化指针
user := &User{}
```

### Error: database/sql: Scan error

**原因**: 类型不匹配

**解决**:
```go
// ✅ 确保类型匹配
var id int
err := row.Scan(&id)  // 确保 id 是 int

// ✅ 使用正确的类型
var timestamp sql.NullTime
row.Scan(&timestamp)
```

### Error: connection pool exhausted

**原因**: 连接池满了

**解决**:
```go
// 增加连接池大小
sqlDB.SetMaxOpenConns(50)

// 关闭连接
defer rows.Close()
defer db.Close()
```

### Error: context deadline exceeded

**原因**: 请求超时

**解决**:
```go
// ✅ 设置超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// ✅ 检查 context 错误
if err := ctx.Err(); err != nil {
    return err
}
```

## 获取帮助

- 📖 查看 [ARCHITECTURE.md](./ARCHITECTURE.md) 了解系统设计
- 🚀 查看 [QUICKSTART.md](./QUICKSTART.md) 快速开始
- 🔒 查看 [SECURITY.md](./SECURITY.md) 安全指南
- 💻 查看 [DEVELOPMENT.md](./DEVELOPMENT.md) 开发规范
- ⚡ 查看 [PERFORMANCE.md](./PERFORMANCE.md) 性能优化
- 🔐 查看 [RBAC.md](./RBAC.md) 权限系统

---

未找到答案？提交 Issue: https://github.com/Oyxiaoxi/GoHub-Service/issues
