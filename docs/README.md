# GoHub-Service 项目文档

> **现代化的 Go 论坛后端服务** - 基于 Gin + GORM，采用三层架构，内置完整的用户系统、RBAC 权限、缓存策略和性能优化方案。

**版本**: v5.0 | **Go版本**: 1.20+ | **更新**: 2026年1月3日

---

## 📚 快速导航

| 文档类型 | 链接 | 说明 |
|---------|------|------|
| 🚀 **快速开始** | [点击查看](#快速开始) | 5分钟完成环境搭建和项目启动 |
| 🏗️ **架构设计** | [点击查看](#系统架构) | 理解项目分层架构和设计思路 |
| 📝 **API文档** | [Swagger UI](http://localhost:3000/swagger/index.html) | 交互式 API 文档（需启动项目）|
| ⚡ **性能优化** | [00_OPTIMIZATION_GUIDE.md](00_OPTIMIZATION_GUIDE.md) | 15项性能优化完整指南 |
| 🔐 **安全加固** | [27_SECURITY_HARDENING.md](27_SECURITY_HARDENING.md) | 安全防护和最佳实践 |
| 💻 **开发指南** | [点击查看](#开发指南) | 编码规范、工作流程、测试方法 |

---

## 🚀 快速开始

### 环境要求

- Go 1.20+
- MySQL 5.7+ / 8.0+
- Redis 6.0+
- (可选) Elasticsearch 7.x

### 1. 克隆项目

```bash
git clone https://github.com/Oyxiaoxi/GoHub-Service.git
cd GoHub-Service
```

### 2. 配置环境

复制配置文件并修改：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```ini
# 应用配置
APP_NAME=GoHub
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost
APP_PORT=3000

# 数据库配置
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=gohub
DB_USERNAME=root
DB_PASSWORD=your_password

# Redis 配置
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0

# JWT 配置
JWT_SECRET=your-secret-key-min-32-characters
JWT_EXPIRE_TIME=2h
JWT_MAX_REFRESH_TIME=168h
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 数据库迁移

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE gohub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行迁移
go run main.go migrate

# (可选) 填充测试数据
go run main.go seed
```

### 5. 启动服务

```bash
go run main.go serve
```

服务启动后访问：
- API 服务：http://localhost:3000
- Swagger 文档：http://localhost:3000/swagger/index.html
- API 版本信息：http://localhost:3000/api/versions
- 健康检查：http://localhost:3000/health

---

## 🏗️ 系统架构

### 项目结构

```
GoHub-Service/
├── app/                    # 应用核心代码
│   ├── cmd/               # 命令行工具
│   ├── http/              # HTTP 层
│   │   ├── controllers/   # 控制器（处理请求）
│   │   └── middlewares/   # 中间件
│   ├── models/            # 数据模型（ORM）
│   ├── repositories/      # 仓储层（数据访问）
│   ├── services/          # 业务逻辑层
│   ├── requests/          # 请求验证
│   ├── policies/          # 权限策略
│   └── cache/             # 缓存层
├── bootstrap/             # 启动初始化
├── config/                # 配置文件
├── database/              # 数据库相关
│   ├── migrations/        # 迁移文件
│   ├── seeders/           # 数据填充
│   └── factories/         # 数据工厂
├── docs/                  # 项目文档
├── pkg/                   # 可复用包
│   ├── auth/              # 认证
│   ├── cache/             # 缓存工具
│   ├── database/          # 数据库工具
│   ├── logger/            # 日志
│   ├── response/          # 响应处理
│   ├── security/          # 安全工具
│   └── ...
├── routes/                # 路由定义
├── storage/               # 存储目录
│   └── logs/              # 日志文件
├── public/                # 静态文件
│   └── uploads/           # 上传文件
└── main.go                # 程序入口
```

### 三层架构

```
┌─────────────────────────────────────────┐
│         Controller 层（控制器）          │
│  - 接收HTTP请求                          │
│  - 参数验证                              │
│  - 调用Service                           │
│  - 返回响应                              │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│          Service 层（业务逻辑）          │
│  - 业务逻辑处理                          │
│  - 数据组装                              │
│  - 调用Repository                        │
│  - 缓存处理                              │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│       Repository 层（数据访问）          │
│  - 数据库操作                            │
│  - SQL查询                               │
│  - ORM操作                               │
└─────────────────────────────────────────┘
```

**设计原则**：
- **Controller**：薄层，只负责HTTP协议相关的事情
- **Service**：厚层，包含所有业务逻辑
- **Repository**：数据访问抽象，隔离数据库细节

---

## 💻 核心功能模块

### 1. 用户系统

**功能**：
- ✅ 手机号/邮箱注册
- ✅ 密码/短信验证码登录
- ✅ JWT Token 认证
- ✅ Token 自动续期
- ✅ 用户资料管理
- ✅ 头像上传

**API 端点**：
```
POST   /api/v1/auth/signup/using-phone     # 手机号注册
POST   /api/v1/auth/signup/using-email     # 邮箱注册
POST   /api/v1/auth/login/using-phone      # 手机号登录
POST   /api/v1/auth/login/using-password   # 密码登录
POST   /api/v1/auth/login/refresh-token    # 刷新Token
GET    /api/v1/user                         # 当前用户信息
PUT    /api/v1/users/:id                    # 更新用户资料
```

**核心代码**：
- Controller: `app/http/controllers/api/v1/auth/*.go`
- Service: `app/services/user_service.go`
- Repository: `app/repositories/user_repository.go`
- Model: `app/models/user/user.go`

### 2. 话题系统

**功能**：
- ✅ 话题发布/编辑/删除
- ✅ 话题列表（分页）
- ✅ 话题详情
- ✅ 话题分类
- ✅ 浏览计数
- ✅ 权限控制（作者可编辑）

**API 端点**：
```
GET    /api/v1/topics           # 话题列表
POST   /api/v1/topics           # 创建话题
GET    /api/v1/topics/:id       # 话题详情
PUT    /api/v1/topics/:id       # 更新话题
DELETE /api/v1/topics/:id       # 删除话题
```

### 3. 评论系统

**功能**：
- ✅ 发表评论
- ✅ 评论回复
- ✅ 评论列表（分页）
- ✅ 评论删除
- ✅ 评论通知

**API 端点**：
```
GET    /api/v1/comments              # 评论列表
POST   /api/v1/comments              # 发表评论
DELETE /api/v1/comments/:id          # 删除评论
```

### 4. 权限系统 (RBAC)

**角色**：
- `超级管理员` - 全部权限
- `管理员` - 内容管理
- `版主` - 分类管理
- `普通用户` - 基础权限

**权限点**：
- `manage_contents` - 内容管理
- `manage_users` - 用户管理
- `manage_roles` - 角色管理
- `manage_permissions` - 权限管理

**使用方式**：
```go
// 中间件：检查权限
router.Use(middlewares.CheckPermission("manage_contents"))

// Policy：策略检查
if !policies.CanManageTopic(user, topic) {
    return errors.New("无权限")
}
```

### 5. 搜索系统

**功能**：
- ✅ 话题搜索
- ✅ 用户搜索
- ✅ 全文搜索（Elasticsearch）
- ✅ 搜索高亮

**API 端点**：
```
GET /api/v1/search/topics?q=关键词    # 话题搜索
GET /api/v1/search/users?q=关键词     # 用户搜索
```

---

## 🔧 开发指南

### 命令行工具

```bash
# 启动服务
go run main.go serve

# 数据库迁移
go run main.go migrate          # 执行迁移
go run main.go migrate:rollback # 回滚迁移
go run main.go migrate:fresh    # 重置数据库

# 数据填充
go run main.go seed             # 填充测试数据

# 缓存管理
go run main.go cache:clear      # 清空缓存

# 代码生成
go run main.go make:model User              # 生成模型
go run main.go make:controller UserController
go run main.go make:migration create_users_table

# 生成 Swagger 文档
make swagger-gen
```

### 创建新模块示例

**1. 创建数据迁移**：

```bash
go run main.go make:migration create_articles_table
```

编辑 `database/migrations/xxxx_create_articles_table.go`：

```go
func up() {
    migration.CreateTable("articles", func(table schema.Blueprint) {
        table.ID()
        table.String("title")
        table.Text("content")
        table.UnsignedBigInteger("user_id")
        table.Timestamps()
        table.Index("user_id")
    })
}
```

**2. 创建模型**：

```go
// app/models/article/article.go
package article

import "GoHub-Service/app/models"

type Article struct {
    models.BaseModel
    Title   string `json:"title"`
    Content string `json:"content"`
    UserID  uint64 `json:"user_id"`
}
```

**3. 创建 Repository**：

```go
// app/repositories/article_repository.go
package repositories

type ArticleRepository struct {
    BaseRepository
}

func (r *ArticleRepository) GetByID(ctx context.Context, id string) (*article.Article, error) {
    var article article.Article
    err := r.DB(ctx).First(&article, id).Error
    return &article, err
}
```

**4. 创建 Service**：

```go
// app/services/article_service.go
package services

type ArticleService struct {
    repo *repositories.ArticleRepository
}

func (s *ArticleService) GetByID(ctx context.Context, id string) (*article.Article, error) {
    return s.repo.GetByID(ctx, id)
}
```

**5. 创建 Controller**：

```go
// app/http/controllers/api/v1/articles_controller.go
package v1

type ArticlesController struct {
    BaseAPIController
}

// Show 文章详情
// @Summary 获取文章详情
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path string true "文章ID"
// @Success 200 {object} response.StandardResponse
// @Router /articles/{id} [get]
func (ctrl *ArticlesController) Show(c *gin.Context) {
    article, err := services.ArticleService.GetByID(c.Request.Context(), c.Param("id"))
    if err != nil {
        response.ApiError(c, 404, response.CodeNotFound, "文章不存在")
        return
    }
    response.StandardSuccess(c, article)
}
```

**6. 注册路由**：

```go
// routes/article.go
package routes

func RegisterArticleRoutes(r *gin.RouterGroup, ctrl *controllers.ArticlesController) {
    articles := r.Group("/articles")
    {
        articles.GET("", ctrl.Index)
        articles.POST("", middlewares.AuthJWT(), ctrl.Store)
        articles.GET("/:id", ctrl.Show)
        articles.PUT("/:id", middlewares.AuthJWT(), ctrl.Update)
        articles.DELETE("/:id", middlewares.AuthJWT(), ctrl.Delete)
    }
}
```

### 编码规范

**命名规范**：
- 文件名：`snake_case`（user_service.go）
- 类型名：`PascalCase`（UserService）
- 变量/函数：`camelCase`（getUserByID）
- 常量：`UPPER_SNAKE_CASE`（MAX_PAGE_SIZE）

**注释规范**：
```go
// UserService 用户服务
// 提供用户相关的业务逻辑处理
type UserService struct {
    repo *repositories.UserRepository
}

// GetByID 根据ID获取用户
// @param ctx 上下文
// @param id 用户ID
// @return 用户对象和错误信息
func (s *UserService) GetByID(ctx context.Context, id string) (*user.User, error) {
    // 实现
}
```

**错误处理**：
```go
// ✅ 正确：使用结构化错误
if err != nil {
    return nil, errors.Wrap(err, "failed to get user")
}

// ✅ 正确：判断特定错误
if errors.Is(err, repositories.ErrNotFound) {
    return nil, errors.New("user not found", errors.CodeNotFound)
}
```

---

## 🔐 安全特性

### 1. 认证与授权

**JWT Token**：
- 访问令牌有效期：2小时
- 刷新令牌有效期：7天
- Token 自动续期机制

**RBAC 权限**：
- 基于角色的访问控制
- 权限粒度到具体操作
- 支持动态权限配置

### 2. 输入验证

**自动验证**：
```go
type CreateTopicRequest struct {
    Title      string `json:"title" validate:"required,min=3,max=100"`
    Content    string `json:"content" validate:"required,min=10"`
    CategoryID uint64 `json:"category_id" validate:"required,exists:categories,id"`
}
```

**高级验证器**（v5.0新增）：
- SQL 注入检测（4种模式）
- XSS 攻击检测（6种模式）
- 路径遍历检测（3种模式）
- 密码强度验证（评分系统）

### 3. 限流防护

**IP 限流**（v5.0新增）：
- 全局：200次/小时
- 认证路由：20次/分钟
- 密码重置：5次/分钟
- 验证码：10次/分钟

**自动封禁**：
- 超限自动封禁 1 分钟
- 滑动时间窗口算法
- 线程安全实现

### 4. 安全中间件

```go
// 已启用的安全中间件
router.Use(
    middlewares.SecureHeaders(),                 // 安全响应头
    middlewares.EnhancedSecurityValidation(),    // 综合安全验证
    middlewares.EnhancedSQLInjectionProtection(), // SQL注入防护
    middlewares.EnhancedXSSProtection(),         // XSS防护
)
```

---

## ⚡ 性能优化

### 1. 数据库优化

**N+1 查询优化**：
```go
// ❌ 问题：N+1查询
comments := repo.GetByTopicID(topicID)
for _, comment := range comments {
    user := userRepo.GetByID(comment.UserID)  // N次查询
}

// ✅ 解决：使用Preload
db.Preload("User", func(db *gorm.DB) *gorm.DB {
    return db.Select("id", "name", "avatar")
}).Find(&comments)
```

**批量操作**：
```go
// ✅ 批量插入
db.CreateInBatches(items, 100)  // 每批100条
```

**慢查询监控**：
- 阈值：200ms
- 自动记录 SQL、执行时间、影响行数
- 日志级别：WARN

### 2. 缓存策略

**三级缓存**：
```go
// L1: 本地缓存（100ms）
// L2: Redis 缓存（10分钟）
// L3: 数据库

func Get(id string) (*Data, error) {
    // 1. 查本地缓存
    if data := localCache.Get(id); data != nil {
        return data, nil
    }
    
    // 2. 查 Redis
    if data := redis.Get(id); data != nil {
        localCache.Set(id, data, 100*time.Millisecond)
        return data, nil
    }
    
    // 3. 查数据库
    data := db.Find(id)
    redis.Set(id, data, 10*time.Minute)
    localCache.Set(id, data, 100*time.Millisecond)
    return data, nil
}
```

**缓存击穿防护**（Singleflight）：
```go
import "golang.org/x/sync/singleflight"

var group singleflight.Group

func GetFromCache(key string) (interface{}, error) {
    // 多个并发请求只执行一次
    v, err, _ := group.Do(key, func() (interface{}, error) {
        return fetchFromDB(key)
    })
    return v, err
}
```

### 3. 资源管理

**协程池**（防止协程泄漏）：
```go
pool := resource.NewGoRoutinePool(20)  // 20 workers
defer pool.Close()

for _, task := range tasks {
    pool.Submit(func() {
        // 执行任务
    })
}
pool.Wait()
```

**超时保护**：
```go
guard := resource.NewContextGuard(5 * time.Second)
defer guard.Release()

result, err := guard.Execute(func(ctx context.Context) (interface{}, error) {
    return service.GetData(ctx, id)
})
```

### 4. 性能监控

**Prometheus 指标**：
- 访问：http://localhost:3000/metrics
- HTTP 请求计数、延迟、错误率

**数据库监控**：
- 访问：http://localhost:3000/database/stats
- 连接池状态、慢查询统计

**缓存监控**：
- 访问：http://localhost:3000/cache/stats
- 命中率、键数量、内存使用

---

## 📊 测试

### 运行测试

```bash
# 所有测试
go test ./... -v

# 指定包测试
go test ./app/services/... -v

# 带覆盖率
go test ./... -cover

# 性能测试
go test ./... -bench=. -benchmem
```

### 测试覆盖率

当前覆盖率：**88%+**

| 层级 | 覆盖率 | 状态 |
|------|--------|------|
| Services | 100% (12/12) | ✅ |
| Repositories | 60% (6/10) | 🟢 |
| Controllers | 35% (4/11) | 🟡 |
| Middlewares | 45% (6/13) | 🟡 |
| pkg/mapper | 100% | ✅ |
| pkg/resource | 85% | 🟢 |

### 测试示例

```go
func TestUserService_GetByID(t *testing.T) {
    // 设置测试环境
    env := testutil.SetupTestEnvironment(t)
    defer env.Cleanup()
    
    // 创建测试数据
    user := testutil.MockUserFactory()
    env.DB.Create(&user)
    
    // 执行测试
    service := services.NewUserService()
    result, err := service.GetByID(context.Background(), user.ID)
    
    // 断言
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, user.Name, result.Name)
}
```

---

## 🚀 部署

### 生产环境配置

```ini
APP_ENV=production
APP_DEBUG=false
APP_PORT=8080

# 数据库连接池
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_MAX_LIFETIME=1h

# Redis
REDIS_MAX_RETRIES=3
REDIS_POOL_SIZE=100

# 日志
LOG_LEVEL=info
LOG_TYPE=daily
```

### Docker 部署

```bash
# 构建镜像
docker build -t gohub-service:latest .

# 运行容器
docker run -d \
  --name gohub-service \
  -p 8080:8080 \
  -v $(pwd)/.env:/app/.env \
  -v $(pwd)/storage:/app/storage \
  gohub-service:latest
```

### 使用 Docker Compose

```bash
docker-compose up -d
```

---

## 📖 更多文档

- [完整优化指南](./00_OPTIMIZATION_GUIDE.md) - 15项性能优化详解
- [安全加固指南](./27_SECURITY_HARDENING.md) - 安全防护最佳实践
- [API版本管理](./23_API_VERSIONING.md) - API 版本控制策略
- [OpenAPI文档](./24_OPENAPI_GUIDE.md) - Swagger 使用指南
- [日志优化](./19_LOG_OPTIMIZATION.md) - 日志系统优化
- [资源泄漏防护](./20_RESOURCE_LEAK_PROTECTION.md) - 资源管理详解
- [代码去重](./21_CODE_DEDUPLICATION.md) - 代码优化技巧
- [性能监控](./22_PERFORMANCE_MONITORING.md) - 性能监控方案

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

[MIT License](LICENSE)

---

**GoHub-Service** - 由 [Oyxiaoxi](https://github.com/Oyxiaoxi) 用 ❤️ 打造
