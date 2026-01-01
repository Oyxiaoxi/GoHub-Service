# 👨‍💻 开发指南与测试

**最后更新**: 2026年1月1日 | **版本**: v2.0

---

## 📖 目录

1. [开发环境配置](#开发环境配置)
2. [项目结构说明](#项目结构说明)
3. [编码规范](#编码规范)
4. [开发工作流](#开发工作流)
5. [单元测试](#单元测试)
6. [集成测试](#集成测试)
7. [性能测试](#性能测试)
8. [测试覆盖率](#测试覆盖率)
9. [常见问题](#常见问题)

---

## 🔧 开发环境配置

### 系统要求

```
✅ Go 1.21 或更高版本
✅ MySQL 8.0+
✅ Redis 7.0+
✅ Elasticsearch 8.5+
✅ Git 2.30+
✅ Docker & Docker Compose (可选，但推荐)
```

### IDE推荐配置

**VS Code**:
```json
{
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "go.lintArgs": ["--fast"],
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "golang.go"
}
```

**GoLand/IntelliJ**:
- 安装 Go Plugin
- 启用 Code Inspections
- 配置 Gofmt on Save

### 依赖安装

```bash
# 安装Go依赖
go mod download
go mod tidy

# 安装开发工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/cosmtrek/air@latest  # 热重载
```

---

## 📂 项目结构说明

```
GoHub-Service/
├── app/                    # 应用层代码
│   ├── cache/             # 缓存实现
│   │   ├── cache_tiers.go
│   │   ├── comment_cache.go
│   │   └── topic_cache.go
│   ├── cmd/               # CLI命令
│   │   ├── cmd.go         # 命令注册
│   │   ├── serve.go       # 启动服务
│   │   ├── migrate.go     # 数据库迁移
│   │   ├── seed.go        # 种子数据
│   │   └── elasticsearch.go
│   ├── http/              # HTTP处理层
│   │   ├── controllers/   # 控制器
│   │   │   ├── user_controller.go
│   │   │   ├── topic_controller.go
│   │   │   ├── comment_controller.go
│   │   │   └── search_controller.go
│   │   └── middlewares/   # 中间件
│   │       ├── auth.go
│   │       ├── cors.go
│   │       └── limiter.go
│   ├── models/            # 数据模型
│   │   ├── model.go       # 基础模型
│   │   ├── user/          # 用户模型
│   │   ├── topic/         # 话题模型
│   │   └── comment/       # 评论模型
│   ├── repositories/      # 数据访问层
│   │   ├── user_repository.go
│   │   ├── user_repository_test.go
│   │   ├── topic_repository.go
│   │   └── topic_repository_test.go
│   ├── requests/          # 请求验证
│   │   ├── user_request.go
│   │   ├── topic_request.go
│   │   └── comment_request.go
│   ├── services/          # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── topic_service.go
│   │   └── comment_service.go
│   └── policies/          # 权限策略
│       ├── topic_policy.go
│       └── comment_policy.go
│
├── pkg/                   # 公共包
│   ├── app/              # 应用实例
│   ├── auth/             # 认证逻辑
│   ├── cache/            # 缓存工具
│   ├── controller/       # 基础控制器
│   ├── database/         # 数据库工具
│   ├── elasticsearch/    # 搜索引擎
│   │   ├── client.go     # ES客户端
│   │   ├── index.go      # 索引管理
│   │   ├── search.go     # 搜索服务
│   │   └── sync.go       # 数据同步
│   ├── errors/           # 错误处理
│   ├── hash/             # 哈希工具
│   ├── helpers/          # 辅助函数
│   ├── jwt/              # JWT认证
│   ├── logger/           # 日志系统
│   ├── mail/             # 邮件服务
│   ├── paginator/        # 分页工具
│   ├── repository/       # 仓储基类
│   ├── response/         # 响应工具
│   ├── security/         # 安全防护
│   └── service/          # 服务基类
│
├── config/               # 配置管理
│   ├── app.go
│   ├── database.go
│   ├── redis.go
│   ├── jwt.go
│   └── elasticsearch.go
│
├── bootstrap/            # 启动初始化
│   ├── app.go
│   ├── database.go
│   ├── redis.go
│   └── elasticsearch.go
│
├── routes/               # 路由定义
│   ├── api.go
│   ├── admin.go
│   ├── topic.go
│   ├── comment.go
│   ├── user.go
│   └── search.go
│
├── database/             # 数据库资源
│   ├── migrations/       # 数据库迁移
│   │   ├── 2024_01_01_create_users_table.go
│   │   ├── 2024_01_02_create_topics_table.go
│   │   └── ...
│   └── seeders/          # 种子数据
│       ├── user_seeder.go
│       └── category_seeder.go
│
├── docs/                 # 📖 文档
│   ├── 00_INDEX.md       # 文档索引
│   ├── 01_QUICKSTART.md
│   ├── ...
│   └── 12_FAQ.md
│
├── scripts/              # 辅助脚本
│   ├── backup-database.sh
│   ├── pre-deploy-check.sh
│   └── run-tests.sh
│
├── main.go              # 应用入口
├── Makefile             # 快速命令
├── go.mod              # 模块定义
└── go.sum              # 依赖校验
```

### 分层架构说明

```
┌─────────────────────────────────────┐
│      HTTP Layer (Gin Framework)     │
├─────────────────────────────────────┤
│  Controllers (http/controllers/)    │
│  - 处理HTTP请求/响应                 │
│  - 请求验证                          │
│  - 调用Service层                     │
├─────────────────────────────────────┤
│  Services (app/services/)           │
│  - 业务逻辑处理                      │
│  - 事务管理                          │
│  - 数据组装                          │
├─────────────────────────────────────┤
│  Repositories (app/repositories/)   │
│  - 数据访问                          │
│  - SQL构建                           │
│  - 缓存操作                          │
├─────────────────────────────────────┤
│  Models (app/models/)               │
│  - 数据结构定义                      │
│  - 字段验证                          │
│  - 关系定义                          │
├─────────────────────────────────────┤
│  Infrastructure (pkg/)              │
│  - 数据库连接                        │
│  - 缓存连接                          │
│  - 日志系统                          │
│  - 认证系统                          │
└─────────────────────────────────────┘
```

---

## 📝 编码规范

### Go编码风格

遵循 [Effective Go](https://golang.org/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

#### 命名规范

```go
// ✅ 良好示例
type UserRepository struct { }
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) { }

// ❌ 避免
type user_repository struct { }
func GetUserByID(ctx context.Context, id int64) { }

// 常数
const (
    MaxLoginAttempts = 5
    DefaultPageSize  = 20
)

// 接口以-er结尾
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 私有变量小写开头
var (
    defaultClient *http.Client
    mu            sync.Mutex
)
```

#### 包组织

```go
package repositories

import (
    "context"
    "database/sql"
    
    "gohub/pkg/database"
    "gohub/app/models"
)

// 常量、类型定义
const TableName = "users"

type UserRepository struct {
    DB *gorm.DB
}

// 构造函数
func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{DB: db}
}

// 公开方法
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    // 实现
}

// 私有方法
func (r *UserRepository) formatQuery(q *gorm.DB) *gorm.DB {
    // 实现
}
```

#### 错误处理

```go
// ✅ 正确做法
if err != nil {
    logger.Errorf("failed to get user: %v", err)
    return nil, errors.Wrap(err, "get user failed")
}

// ❌ 避免
if err != nil {
    panic(err)
}

// ❌ 避免
if err != nil {
    return
}
```

#### 上下文使用

```go
// ✅ 所有IO操作都接受context
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    return r.DB.WithContext(ctx).First(&User{}, id).Error
}

// ❌ 不传递context
func (r *UserRepository) GetByID(id int64) (*User, error) {
    return r.DB.First(&User{}, id).Error
}
```

### 控制器规范

```go
package controllers

import (
    "gohub/pkg/controller"
    "gohub/pkg/response"
    "gohub/app/services"
)

type UserController struct {
    controller.BaseController
    userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}

// 控制器方法
func (uc *UserController) Show(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    user, err := uc.userService.GetByID(c.Request.Context(), userID)
    if err != nil {
        response.Error(c, "用户不存在")
        return
    }
    
    response.Success(c, user)
}
```

### 数据库操作规范

```go
// ✅ 使用事务处理复杂操作
func (s *TopicService) CreateWithTags(ctx context.Context, topic *Topic, tags []string) error {
    return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 创建话题
        if err := tx.Create(topic).Error; err != nil {
            return err
        }
        
        // 创建标签关联
        for _, tag := range tags {
            if err := tx.Create(&TopicTag{
                TopicID: topic.ID,
                Tag:     tag,
            }).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}

// ✅ 使用预加载优化查询
func (r *TopicRepository) GetWithRelations(ctx context.Context, id int64) (*Topic, error) {
    var topic Topic
    err := r.DB.WithContext(ctx).
        Preload("User").
        Preload("Category").
        Preload("Comments", func(db *gorm.DB) *gorm.DB {
            return db.Order("created_at DESC")
        }).
        First(&topic, id).Error
    return &topic, err
}
```

### 日志规范

```go
import "gohub/pkg/logger"

// ✅ 使用结构化日志
logger.Infof("user login", map[string]interface{}{
    "user_id": userID,
    "ip": c.ClientIP(),
    "duration_ms": duration,
})

// ✅ 错误日志包含堆栈跟踪
if err != nil {
    logger.Errorf("failed to create topic: %+v", err)
}

// ❌ 避免
log.Println("user login")
```

---

## 🔄 开发工作流

### 本地开发流程

```bash
# 1. 从main分支创建功能分支
git checkout main
git pull origin main
git checkout -b feature/user-authentication

# 2. 开发功能
# 编辑代码、添加测试

# 3. 运行测试确保通过
make test
make test-coverage

# 4. 代码提交
git add .
git commit -m "feat: add user authentication with JWT"

# 5. 推送到远程
git push origin feature/user-authentication

# 6. 创建Pull Request，请求代码审查

# 7. 代码审查完成，合并到main
```

### 分支命名规范

```
feature/功能名称          新功能
fix/问题描述              bug修复
refactor/重构内容         代码重构
docs/文档内容             文档更新
chore/杂务内容            维护工作

✅ 示例:
feature/user-registration
fix/elasticsearch-sync-timeout
refactor/topic-repository
docs/api-documentation-update
```

### Commit Message规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/)

```
<type>(<scope>): <subject>

<body>

<footer>

---

类型 (type):
- feat: 新功能
- fix: bug修复
- docs: 文档更新
- style: 代码风格调整
- refactor: 代码重构
- perf: 性能优化
- test: 测试相关
- chore: 构建/依赖相关

示例:
feat(auth): add JWT token refresh endpoint

- Add refresh_token endpoint
- Implement token refresh logic
- Add integration tests

Closes #123
```

---

## 🧪 单元测试

### 测试文件位置

```
app/repositories/
  ├── user_repository.go
  └── user_repository_test.go       ← 同一目录

app/services/
  ├── topic_service.go
  └── topic_service_test.go
```

### 基础单元测试

```go
package repositories

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "gorm.io/gorm"
)

func TestUserRepository_GetByID(t *testing.T) {
    // Arrange（准备）
    db := setupTestDB()
    repo := NewUserRepository(db)
    
    user := &User{Name: "John", Email: "john@example.com"}
    require.NoError(t, db.Create(user).Error)
    
    // Act（执行）
    result, err := repo.GetByID(context.Background(), user.ID)
    
    // Assert（断言）
    assert.NoError(t, err)
    assert.Equal(t, user.Name, result.Name)
    assert.Equal(t, user.Email, result.Email)
}

func TestUserRepository_Create(t *testing.T) {
    tests := []struct {
        name    string
        user    *User
        wantErr bool
    }{
        {
            name:    "valid user",
            user:    &User{Name: "John", Email: "john@example.com"},
            wantErr: false,
        },
        {
            name:    "missing email",
            user:    &User{Name: "John"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db := setupTestDB()
            repo := NewUserRepository(db)
            
            err := repo.Create(context.Background(), tt.user)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Mock测试

```go
import "github.com/golang/mock/gomock"

func TestUserService_Register(t *testing.T) {
    // 创建mock
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo)
    
    // 设置期望
    mockRepo.EXPECT().
        GetByEmail(gomock.Any(), "john@example.com").
        Return(nil, gorm.ErrRecordNotFound).
        Times(1)
    
    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil).
        Times(1)
    
    // 执行
    err := service.Register(context.Background(), "john@example.com", "password")
    
    // 验证
    assert.NoError(t, err)
}
```

### 运行单元测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./app/repositories -v

# 运行特定测试
go test -run TestUserRepository_GetByID ./app/repositories -v

# 并行运行（加速）
go test -parallel 4 ./...

# 显示覆盖率
go test -cover ./...
```

---

## 🔗 集成测试

### 数据库集成测试

```go
package repositories

import (
    "context"
    "testing"
    "github.com/stretchr/testify/require"
)

func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // 连接真实数据库
    db := setupTestDatabase()
    repo := NewUserRepository(db)
    
    // 测试创建和查询
    user := &User{Name: "John", Email: "john@example.com"}
    err := repo.Create(context.Background(), user)
    require.NoError(t, err)
    
    retrieved, err := repo.GetByID(context.Background(), user.ID)
    require.NoError(t, err)
    require.Equal(t, user.Name, retrieved.Name)
}
```

### HTTP集成测试

```go
package controllers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestUserController_Show(t *testing.T) {
    // 创建测试引擎
    router := gin.New()
    userService := setupMockUserService()
    controller := NewUserController(userService)
    
    router.GET("/users/:id", controller.Show)
    
    // 创建请求
    req, _ := http.NewRequest("GET", "/users/1", nil)
    w := httptest.NewRecorder()
    
    // 执行
    router.ServeHTTP(w, req)
    
    // 验证
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response struct {
        Code int
        Data json.RawMessage
    }
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, 200, response.Code)
}
```

### 运行集成测试

```bash
# 运行所有集成测试
make test-integration

# 运行包括单元和集成测试
go test -tags=integration ./...
```

---

## ⚡ 性能测试

### 基准测试

```go
func BenchmarkUserRepository_GetByID(b *testing.B) {
    db := setupBenchDB()
    repo := NewUserRepository(db)
    user := &User{Name: "John", Email: "john@example.com"}
    db.Create(user)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        repo.GetByID(context.Background(), user.ID)
    }
}

func BenchmarkTopicRepository_Search(b *testing.B) {
    db := setupBenchDB()
    repo := NewTopicRepository(db)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        repo.Search(context.Background(), "golang", 1, 20)
    }
}
```

### 运行性能测试

```bash
# 运行基准测试
go test -bench=. ./app/repositories

# 详细输出
go test -bench=. -benchmem ./app/repositories

# 比较两次运行
go test -bench=. -benchmem ./app/repositories | tee old.txt
# 修改代码...
go test -bench=. -benchmem ./app/repositories | tee new.txt
benchstat old.txt new.txt
```

---

## 📊 测试覆盖率

### 生成覆盖率报告

```bash
# 生成覆盖率文件
go test -coverprofile=coverage.out ./...

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html

# 查看覆盖率摘要
go tool cover -func=coverage.out | tail -1

# 获取特定包的覆盖率
go test -coverprofile=coverage.out -coverpkg=./app/services ./...
```

### 覆盖率目标

```
项目级别         目标覆盖率
────────────────────────
核心业务逻辑      > 80%
数据访问层        > 75%
HTTP控制器        > 70%
工具函数          > 60%
────────────────────────
整体目标          > 70%
```

### 测试总结（最新数据）

| 模块 | 覆盖率 | 单元测试 | 集成测试 |
|------|--------|---------|----------|
| **repositories** | 82% | 45 | 12 |
| **services** | 78% | 38 | 8 |
| **controllers** | 71% | 28 | 15 |
| **models** | 85% | 22 | 3 |
| **总计** | 79% | 133 | 38 |

---

## ❓ 常见问题

### 问题1：测试超时

```bash
# 增加超时时间
go test -timeout 5m ./...

# 调试特定测试
go test -run TestName -v -timeout 10m ./...
```

### 问题2：数据库状态污染

```go
// ✅ 使用事务隔离
func TestWithTransaction(t *testing.T) {
    db := setupTestDB()
    
    db.Exec("BEGIN")
    defer db.Exec("ROLLBACK")
    
    // 测试代码
}
```

### 问题3：并发测试失败

```go
// ✅ 使用同步原语
func TestConcurrent(t *testing.T) {
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // 测试代码
        }()
    }
    
    wg.Wait()
}
```

---

## 📚 相关文档

- [项目架构详解](02_ARCHITECTURE.md) - 深入理解系统设计
- [性能优化指南](07_PERFORMANCE.md) - 性能调优技巧

---
