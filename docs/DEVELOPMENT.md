# 💻 开发指南

编码规范、项目流程和最佳实践。

## 开发环境设置

### 前置要求
- Go 1.20+
- MySQL 8.0+ 或 SQLite 3.0+
- Redis 6.0+
- git

### 初始化项目

```bash
# 1. 克隆项目
git clone https://github.com/Oyxiaoxi/GoHub-Service.git
cd GoHub-Service

# 2. 安装依赖
go mod download
go mod tidy

# 3. 复制配置文件
cp .env.example .env

# 4. 编辑配置文件
# 设置数据库、Redis、邮件等信息

# 5. 运行迁移
go run main.go migrate

# 6. 运行数据清理（可选）
go run main.go seed

# 7. 启动服务
go run main.go serve
```

## 项目命令

### 数据库迁移

```bash
# 运行所有迁移
go run main.go migrate

# 查看迁移状态
go run main.go migrate:status

# 回滚最后一次迁移
go run main.go migrate:rollback
```

### 数据填充

```bash
# 运行所有 seeder
go run main.go seed

# 运行特定 seeder
go run main.go seed --class=UsersSeeder
```

### CLI 命令生成

```bash
# 生成模型脚手架
go run main.go make:model Post

# 生成迁移文件
go run main.go make:migration create_posts_table

# 生成 controller
go run main.go make:controller PostController

# 生成 service
go run main.go make:service PostService

# 生成 repository
go run main.go make:repository PostRepository
```

## 代码结构与规范

### 1. 模型定义

文件位置: `app/models/{feature}/`

```go
package post

import "time"

type Post struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index"`
    Title     string    `gorm:"size:200;not null"`
    Content   string    `gorm:"type:text"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `gorm:"index"`
    
    // 关联关系
    User      User      `gorm:"foreignKey:UserID"`
    Comments  []Comment `gorm:"foreignKey:PostID"`
}

// 表名
func (Post) TableName() string {
    return "posts"
}
```

### 2. 数据库迁移

文件位置: `database/migrations/`

```go
package migrations

import (
    "database/sql"
    "time"
)

func CreatePostsTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTO_INCREMENT,
        user_id INTEGER NOT NULL,
        title VARCHAR(200) NOT NULL,
        content LONGTEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP NULL,
        INDEX idx_user_id (user_id),
        INDEX idx_deleted_at (deleted_at),
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `
    
    _, err := db.Exec(query)
    return err
}
```

### 3. Repository 数据访问

文件位置: `app/repositories/`

```go
package repositories

import (
    "GoHub-Service/app/models"
    "gorm.io/gorm"
)

type PostRepository struct {
    db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
    return &PostRepository{db: db}
}

// 创建
func (r *PostRepository) Create(post *models.Post) error {
    return r.db.Create(post).Error
}

// 查询
func (r *PostRepository) Find(id uint) (*models.Post, error) {
    var post models.Post
    err := r.db.First(&post, id).Error
    return &post, err
}

// 列表
func (r *PostRepository) FindAll(page, pageSize int) ([]models.Post, int64, error) {
    var posts []models.Post
    var total int64
    
    offset := (page - 1) * pageSize
    err := r.db.Offset(offset).Limit(pageSize).Find(&posts).Error
    r.db.Model(models.Post{}).Count(&total)
    
    return posts, total, err
}

// 更新
func (r *PostRepository) Update(post *models.Post) error {
    return r.db.Save(post).Error
}

// 删除
func (r *PostRepository) Delete(id uint) error {
    return r.db.Delete(&models.Post{}, id).Error
}
```

### 4. Service 业务逻辑

文件位置: `app/services/`

```go
package services

import (
    "GoHub-Service/app/models"
    "GoHub-Service/app/repositories"
    "GoHub-Service/app/cache"
    "GoHub-Service/pkg/logger"
    "errors"
)

type PostService struct {
    repo  *repositories.PostRepository
    cache *cache.PostCache
}

func NewPostService(repo *repositories.PostRepository, cache *cache.PostCache) *PostService {
    return &PostService{repo: repo, cache: cache}
}

// 创建话题
func (s *PostService) Create(post *models.Post) error {
    // 1. 验证业务规则
    if post.Title == "" {
        return errors.New("标题不能为空")
    }
    
    // 2. 调用 Repository
    if err := s.repo.Create(post); err != nil {
        logger.Error("创建文章失败", err)
        return err
    }
    
    // 3. 清除缓存
    s.cache.Clear()
    
    return nil
}

// 分页查询
func (s *PostService) GetPaginated(page, pageSize int) ([]models.Post, int64, error) {
    // 1. 先查缓存
    key := fmt.Sprintf("posts:page:%d:%d", page, pageSize)
    if cached := s.cache.Get(key); cached != nil {
        return cached.([]models.Post), 0, nil
    }
    
    // 2. 缓存未命中，查数据库
    posts, total, err := s.repo.FindAll(page, pageSize)
    if err != nil {
        return nil, 0, err
    }
    
    // 3. 写入缓存
    s.cache.Set(key, posts, 1*time.Hour)
    
    return posts, total, nil
}
```

### 5. Controller 请求处理

文件位置: `app/http/controllers/`

```go
package controllers

import (
    "GoHub-Service/app/models"
    "GoHub-Service/app/requests"
    "GoHub-Service/app/services"
    "GoHub-Service/pkg/response"
    "github.com/gin-gonic/gin"
    "net/http"
)

type PostController struct {
    service *services.PostService
}

func NewPostController(service *services.PostService) *PostController {
    return &PostController{service: service}
}

// 列表 GET /api/v1/posts
func (ctrl *PostController) Index(c *gin.Context) {
    // 1. 解析参数
    var req requests.PaginationRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    
    // 2. 调用 Service
    posts, total, err := ctrl.service.GetPaginated(req.Page, req.PageSize)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "获取列表失败")
        return
    }
    
    // 3. 返回响应
    response.Paginate(c, posts, total, req.Page, req.PageSize)
}

// 详情 GET /api/v1/posts/:id
func (ctrl *PostController) Show(c *gin.Context) {
    id := c.Param("id")
    post, err := ctrl.service.Get(id)
    if err != nil {
        response.Error(c, http.StatusNotFound, "文章不存在")
        return
    }
    
    response.Success(c, http.StatusOK, post)
}

// 创建 POST /api/v1/posts
func (ctrl *PostController) Store(c *gin.Context) {
    // 1. 解析并验证请求
    var req requests.CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    
    // 2. 创建模型
    post := &models.Post{
        Title:   req.Title,
        Content: req.Content,
    }
    
    // 3. 调用 Service
    if err := ctrl.service.Create(post); err != nil {
        response.Error(c, http.StatusInternalServerError, "创建失败")
        return
    }
    
    // 4. 返回响应
    response.Success(c, http.StatusCreated, post)
}
```

### 6. 请求验证

文件位置: `app/requests/`

```go
package requests

type CreatePostRequest struct {
    Title   string `json:"title" binding:"required,min=3,max=200"`
    Content string `json:"content" binding:"required,min=10"`
}

// 自定义验证
func (r *CreatePostRequest) Validate() error {
    if len(r.Title) < 3 {
        return errors.New("标题长度至少3个字符")
    }
    return nil
}
```

### 7. 路由定义

文件位置: `routes/`

```go
package routes

import (
    "GoHub-Service/app/http/controllers"
    "GoHub-Service/app/http/middlewares"
    "github.com/gin-gonic/gin"
)

func RegisterPostRoutes(r *gin.Engine, ctrl *controllers.PostController) {
    posts := r.Group("/api/v1/posts")
    {
        // 公开端点
        posts.GET("", ctrl.Index)
        posts.GET("/:id", ctrl.Show)
        
        // 需要认证的端点
        posts.POST("", 
            middlewares.Authenticate(),
            middlewares.RequirePermission("posts.create"),
            ctrl.Store)
            
        posts.PUT("/:id",
            middlewares.Authenticate(),
            middlewares.RequirePermission("posts.update"),
            ctrl.Update)
            
        posts.DELETE("/:id",
            middlewares.Authenticate(),
            middlewares.RequirePermission("posts.delete"),
            ctrl.Destroy)
    }
}
```

## 编码规范

### 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 包名 | 小写 | `repositories`, `services` |
| 常量 | 大写下划线 | `MAX_PAGE_SIZE`, `DEFAULT_TIMEOUT` |
| 函数 | 大驼峰 | `CreateUser()`, `GetPaginated()` |
| 变量 | 小驼峰 | `userID`, `pageSize` |
| 接口 | 大驼峰 + er 后缀 | `Reader`, `Writer` |
| 结构体 | 大驼峰 | `User`, `PostService` |

### 错误处理

```go
// ✅ 好的做法
if err != nil {
    logger.Error("操作失败", zap.Error(err))
    return fmt.Errorf("操作失败: %w", err)
}

// ❌ 不好的做法
if err != nil {
    panic(err)  // 不要使用 panic
}

if err != nil {
    // 不要忽略错误
}
```

### 注释规范

```go
// 导出函数必须有注释
// GetUser 通过 ID 查询用户
func (r *UserRepository) GetUser(id uint) (*User, error) {
    // ...
}

// 复杂逻辑添加注释
// 1. 先查缓存
// 2. 缓存未命中查数据库
// 3. 写入缓存
```

## 测试

### 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./app/services/...

# 显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 测试示例

```go
// app/services/user_service_test.go
package services

import (
    "testing"
)

func TestCreateUser(t *testing.T) {
    // Arrange
    service := NewUserService(mockRepo, mockCache)
    user := &User{Name: "John", Email: "john@example.com"}
    
    // Act
    err := service.Create(user)
    
    // Assert
    if err != nil {
        t.Fatalf("预期成功，实际出错: %v", err)
    }
}
```

## 提交 Git

### 提交消息规范

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型:
- `feat`: 新功能
- `fix`: 错误修复
- `docs`: 文档
- `style`: 格式变更
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试

示例:
```
feat(user): 添加用户注册功能

- 实现用户验证逻辑
- 集成邮箱验证
- 添加速率限制

Closes #123
```

## 常见问题

### Q: 如何处理关联关系？
A: 使用 GORM 的关联加载，避免 N+1 查询

```go
// 预加载关联
posts, _ := repo.FindAll()
db.Preload("User").Find(&posts)
```

### Q: 如何处理错误？
A: 返回有意义的错误，不要忽略或使用 panic

### Q: 缓存什么时候失效？
A: 数据更改时立即失效，避免缓存不一致

---

更多信息请查看 [ARCHITECTURE.md](./ARCHITECTURE.md)
