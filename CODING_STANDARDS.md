# GoHub-Service 代码规范文档

> 更新时间：2025年12月28日  
> 版本：v1.0

---

## 📋 目录

1. [项目结构规范](#项目结构规范)
2. [代码风格规范](#代码风格规范)
3. [命名规范](#命名规范)
4. [注释规范](#注释规范)
5. [错误处理规范](#错误处理规范)
6. [API响应规范](#api响应规范)
7. [数据库规范](#数据库规范)
8. [测试规范](#测试规范)

---

## 项目结构规范

### 目录结构

```
GoHub-Service/
├── app/                    # 应用核心代码
│   ├── cmd/               # 命令行工具
│   ├── http/              # HTTP相关
│   │   ├── controllers/   # 控制器层
│   │   └── middlewares/   # 中间件
│   ├── models/            # 数据模型层
│   ├── policies/          # 授权策略
│   └── requests/          # 请求验证
├── bootstrap/             # 初始化引导
├── config/                # 配置文件
├── database/              # 数据库相关
│   ├── factories/         # 数据工厂
│   ├── migrations/        # 数据迁移
│   └── seeders/           # 数据填充
├── pkg/                   # 公共包
│   ├── auth/              # 认证
│   ├── cache/             # 缓存
│   ├── controller/        # 控制器辅助
│   ├── repository/        # 数据仓库模式
│   ├── response/          # 响应处理
│   └── service/           # 服务层
├── routes/                # 路由定义
└── main.go               # 入口文件
```

### 分层架构

```
请求 → 中间件 → 路由 → Controller → Service → Repository → Model → 数据库
                                 ↓
                              响应返回
```

---

## 代码风格规范

### 1. 使用gofmt格式化代码

```bash
# 格式化单个文件
gofmt -w file.go

# 格式化整个项目
gofmt -w .
```

### 2. 遵循Go语言官方规范

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### 3. 代码缩进

- 使用Tab缩进（Go标准）
- 每行代码不超过120个字符

### 4. 导入包规范

```go
import (
    // 标准库
    "fmt"
    "time"
    
    // 第三方库
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    // 项目内部包
    "GoHub-Service/app/models/user"
    "GoHub-Service/pkg/response"
)
```

---

## 命名规范

### 1. 包命名

- 全小写，不使用下划线或驼峰
- 简短且有意义
- 单数形式

```go
// ✅ 正确
package user
package config

// ❌ 错误
package userController
package user_service
```

### 2. 文件命名

- 全小写，使用下划线分隔
- 与包内容相关

```go
// ✅ 正确
user_controller.go
user_model.go
user_util.go

// ❌ 错误
UserController.go
userController.go
```

### 3. 变量命名

- 驼峰命名
- 简短且有意义
- 缩写词保持大写（如ID、HTTP、URL）

```go
// ✅ 正确
var userID uint64
var httpClient *http.Client
var apiURL string

// ❌ 错误
var userId uint64
var httpclient *http.Client
var apiUrl string
```

### 4. 函数命名

- 驼峰命名
- 导出函数首字母大写
- 动词开头

```go
// ✅ 正确
func CreateUser()
func GetUserByID()
func validateEmail()

// ❌ 错误
func create_user()
func getuserbyid()
func EmailValidate()
```

### 5. 常量命名

- 驼峰命名或全大写+下划线
- 首字母大写表示导出

```go
// ✅ 正确
const MaxRetries = 3
const DefaultTimeout = 30
const CODE_SUCCESS = 0

// ❌ 错误
const max_retries = 3
const DEFAULTTIMEOUT = 30
```

---

## 注释规范

### 1. 包注释

```go
// Package user 用户相关业务逻辑
//
// 本包提供用户的创建、查询、更新和删除功能
// 以及用户认证、授权等相关操作
package user
```

### 2. 函数注释

```go
// CreateUser 创建新用户
//
// 参数：
//   - name: 用户名称
//   - email: 用户邮箱
//
// 返回：
//   - *User: 创建的用户对象
//   - error: 错误信息，成功时为nil
func CreateUser(name, email string) (*User, error) {
    // 实现代码
}
```

### 3. 结构体注释

```go
// User 用户模型
// 存储用户的基本信息和认证信息
type User struct {
    models.BaseModel
    
    // Name 用户名称，长度3-255字符
    Name string `gorm:"column:name;type:varchar(255)" json:"name"`
    
    // Email 用户邮箱，必须唯一
    Email string `gorm:"column:email;type:varchar(255);unique" json:"email"`
    
    models.CommonTimestampsField
}
```

### 4. 复杂逻辑注释

```go
func ProcessPayment(amount float64) error {
    // 1. 验证金额是否有效
    if amount <= 0 {
        return errors.New("无效的金额")
    }
    
    // 2. 检查用户余额
    // 这里需要加锁以防止并发问题
    balance := getUserBalance()
    if balance < amount {
        return errors.New("余额不足")
    }
    
    // 3. 执行支付逻辑
    return executePayment(amount)
}
```

---

## 错误处理规范

### 1. 使用统一错误码

```go
// ✅ 正确：使用定义好的错误码
response.ApiErrorWithCode(c, response.CodeUserNotFound)

// ❌ 错误：硬编码错误信息
c.JSON(404, gin.H{"error": "user not found"})
```

### 2. 错误检查

```go
// ✅ 正确：立即检查错误
user, err := userService.GetUser(id)
if err != nil {
    return err
}

// ❌ 错误：忽略错误
user, _ := userService.GetUser(id)
```

### 3. 错误包装

```go
// ✅ 正确：包装错误，提供上下文
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

### 4. 错误日志

```go
// ✅ 正确：记录重要错误
if err != nil {
    logger.Error("数据库操作失败", 
        zap.Error(err),
        zap.Uint64("user_id", userID),
    )
    return err
}
```

---

## API响应规范

### 1. 统一响应格式

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "id": 1,
        "name": "张三"
    }
}
```

### 2. 成功响应

```go
// 返回数据
response.ApiSuccess(c, userData)

// 仅返回成功消息
response.ApiSuccessWithMessage(c, "操作成功")
```

### 3. 错误响应

```go
// 使用错误码
response.ApiErrorWithCode(c, response.CodeUserNotFound)

// 自定义错误消息
response.ApiError(c, http.StatusBadRequest, response.CodeInvalidParams, "参数格式错误")
```

### 4. 分页响应

```go
response.JSON(c, gin.H{
    "data": items,
    "pager": gin.H{
        "page": 1,
        "per_page": 10,
        "total": 100,
    },
})
```

---

## 数据库规范

### 1. 表命名

- 全小写，使用下划线分隔
- 复数形式

```sql
-- ✅ 正确
users
user_profiles
topic_categories

-- ❌ 错误
User
userProfile
topic_category
```

### 2. 字段命名

- 全小写，使用下划线分隔
- 语义明确

```sql
-- ✅ 正确
user_id
created_at
email_verified_at

-- ❌ 错误
userId
createdAt
emailVerifiedAt
```

### 3. 索引

```go
// 单列索引
`gorm:"column:email;index"`

// 唯一索引
`gorm:"column:email;uniqueIndex"`

// 复合索引
`gorm:"index:idx_user_email"`
```

### 4. 外键

```go
type Topic struct {
    models.BaseModel
    
    UserID     uint64 `gorm:"column:user_id;index" json:"user_id"`
    User       user.User `gorm:"foreignKey:UserID" json:"user"`
    
    CategoryID uint64 `gorm:"column:category_id;index" json:"category_id"`
    Category   category.Category `gorm:"foreignKey:CategoryID" json:"category"`
}
```

---

## 测试规范

### 1. 测试文件命名

```
user.go       → user_test.go
controller.go → controller_test.go
```

### 2. 测试函数命名

```go
func TestCreateUser(t *testing.T) {}
func TestGetUserByID(t *testing.T) {}
func TestValidateEmail(t *testing.T) {}
```

### 3. 测试用例结构

```go
func TestCreateUser(t *testing.T) {
    // 1. Setup - 准备测试数据
    user := &User{
        Name:  "测试用户",
        Email: "test@example.com",
    }
    
    // 2. Execute - 执行测试
    err := CreateUser(user)
    
    // 3. Assert - 验证结果
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

### 4. 表格驱动测试

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct{
        name    string
        email   string
        wantErr bool
    }{
        {"有效邮箱", "test@example.com", false},
        {"无效邮箱", "invalid-email", true},
        {"空邮箱", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error = %v, want error = %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## 最佳实践

### 1. DRY原则（Don't Repeat Yourself）

- 提取公共逻辑到函数
- 使用Service层复用业务逻辑
- 使用Repository层复用数据访问逻辑

### 2. 单一职责原则

- Controller只负责HTTP请求处理
- Service负责业务逻辑
- Repository负责数据访问
- Model只定义数据结构

### 3. 依赖注入

```go
// ✅ 正确：通过参数传递依赖
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// ❌ 错误：直接在内部创建依赖
type UserService struct {}

func (s *UserService) GetUser(id uint64) {
    db := database.DB  // 直接使用全局变量
}
```

### 4. 接口设计

```go
// ✅ 正确：小而精的接口
type UserReader interface {
    GetUser(id uint64) (*User, error)
}

type UserWriter interface {
    CreateUser(user *User) error
    UpdateUser(user *User) error
}

// ❌ 错误：大而全的接口
type UserService interface {
    CreateUser()
    UpdateUser()
    DeleteUser()
    GetUser()
    ListUsers()
    // ... 更多方法
}
```

---

## 代码审查清单

- [ ] 代码格式化（gofmt）
- [ ] 命名规范正确
- [ ] 注释完整清晰
- [ ] 错误处理得当
- [ ] 没有硬编码
- [ ] 没有重复代码
- [ ] 单元测试覆盖
- [ ] 性能考虑合理
- [ ] 安全性检查
- [ ] 日志记录适当

---

**注意**：本规范是活的文档，会随着项目发展持续更新和完善。
