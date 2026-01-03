# 错误处理优化指南

## 概述

本项目实现了完善的错误处理机制，提供统一的错误码管理、错误链追踪和上下文信息。通过使用 Go 1.13+ 的 `errors.Is()` 和 `errors.As()` API，实现了更强大的错误判断和类型转换功能。

**版本**: v2.3  
**更新日期**: 2026-01-03

---

## 错误处理架构

### 错误类型（ErrorType）

系统定义了多种错误类型，用于分类不同场景的错误：

```go
const (
    ErrorTypeBusiness      ErrorType = "BUSINESS_ERROR"      // 业务错误
    ErrorTypeValidation    ErrorType = "VALIDATION_ERROR"    // 验证错误
    ErrorTypeAuthorization ErrorType = "AUTHORIZATION_ERROR" // 授权错误
    ErrorTypeNotFound      ErrorType = "NOT_FOUND_ERROR"     // 资源不存在
    ErrorTypeDatabase      ErrorType = "DATABASE_ERROR"      // 数据库错误
    ErrorTypeExternal      ErrorType = "EXTERNAL_ERROR"      // 外部服务错误
    ErrorTypeInternal      ErrorType = "INTERNAL_ERROR"      // 内部错误
    ErrorTypeNetwork       ErrorType = "NETWORK_ERROR"       // 网络错误
    ErrorTypeTimeout       ErrorType = "TIMEOUT_ERROR"       // 超时错误
    ErrorTypeConflict      ErrorType = "CONFLICT_ERROR"      // 冲突错误
)
```

### 错误码分配规则

| 区间 | 用途 | 说明 |
|------|------|------|
| 0 | 成功 | 操作成功 |
| 1000-1099 | 通用错误 | 参数错误、授权错误等 |
| 2000-2099 | 数据库错误 | 连接、查询、事务等 |
| 2100-2199 | 缓存错误 | Redis 相关错误 |
| 2200-2299 | 网络错误 | 超时、连接失败等 |
| 3000-3999 | 第三方服务 | 短信、邮件、支付等 |
| 4000-4999 | 业务错误 | 按模块分配 |

#### 业务模块错误码

- **4000-4099**: 用户模块
- **4100-4199**: 话题模块
- **4200-4299**: 评论模块
- **4300-4399**: 分类模块
- **4400-4499**: 权限模块

---

## AppError 结构

```go
type AppError struct {
    Type       ErrorType              // 错误类型
    Code       int                    // 业务错误码
    Message    string                 // 错误消息
    Details    map[string]interface{} // 错误详情
    Err        error                  // 原始错误（错误链）
    StackTrace string                 // 堆栈信息
    RequestID  string                 // 请求ID（用于追踪）
}
```

### 核心方法

```go
// Error 实现 error 接口
func (e *AppError) Error() string

// Unwrap 支持 errors.Unwrap，用于错误链解包
func (e *AppError) Unwrap() error

// Is 支持 errors.Is，用于错误比较
func (e *AppError) Is(target error) bool

// As 支持 errors.As，用于错误类型转换
func (e *AppError) As(target interface{}) bool

// WithRequestID 添加请求ID
func (e *AppError) WithRequestID(requestID string) *AppError

// WithDetails 添加错误详情
func (e *AppError) WithDetails(details map[string]interface{}) *AppError

// WithError 添加原始错误（用于错误包装）
func (e *AppError) WithError(err error) *AppError
```

---

## 使用指南

### 1. 创建错误

#### 通用错误

```go
// 资源不存在
err := apperrors.NotFoundError("用户")

// 带自定义错误码的资源不存在
err := apperrors.NotFoundErrorWithCode(apperrors.CodeUserNotFound, "用户")

// 参数验证错误
err := apperrors.ValidationError("参数错误", map[string]interface{}{
    "field": "email",
    "error": "邮箱格式不正确",
})

// 授权错误
err := apperrors.AuthorizationError("没有访问权限")

// 未授权
err := apperrors.UnauthorizedError("请先登录")

// 冲突错误（资源已存在）
err := apperrors.ConflictError("用户名")
```

#### 数据库错误

```go
// 通用数据库错误
err := apperrors.DatabaseError("查询", dbErr)

// 创建失败
err := apperrors.DatabaseCreateError("用户", dbErr)

// 更新失败
err := apperrors.DatabaseUpdateError("用户", dbErr)

// 删除失败
err := apperrors.DatabaseDeleteError("用户", dbErr)

// 重复记录
err := apperrors.DatabaseDuplicateError("用户名")
```

#### 其他错误

```go
// 超时错误
err := apperrors.TimeoutError("数据库查询")

// 网络错误
err := apperrors.NetworkError("连接失败", netErr)

// 缓存错误
err := apperrors.CacheError("设置", cacheErr)

// 外部服务错误
err := apperrors.ExternalError("短信服务", smsErr)

// 内部错误
err := apperrors.InternalError("系统异常", internalErr)
```

### 2. Repository 层错误处理

Repository 层提供了便捷的错误创建函数：

```go
// 资源不存在
err := repositories.NewNotFoundError("用户", userID)

// 创建失败
err := repositories.NewCreateError("用户", dbErr)

// 更新失败
err := repositories.NewUpdateError("用户", userID, dbErr)

// 删除失败
err := repositories.NewDeleteError("用户", userID, dbErr)

// 查询失败
err := repositories.NewQueryError("列表查询", dbErr)

// 重复记录
err := repositories.NewDuplicateError("用户", "email", "test@example.com")
```

#### 示例：CommentRepository

```go
func (r *commentRepository) Create(ctx context.Context, c *comment.Comment) error {
    c.Create()
    if c.ID == 0 {
        return repositories.NewCreateError("评论", nil)
    }
    return nil
}

func (r *commentRepository) Delete(ctx context.Context, id string) error {
    var commentModel comment.Comment
    if err := database.DB.WithContext(ctx).First(&commentModel, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return repositories.NewNotFoundError("评论", id)
        }
        return repositories.NewDeleteError("评论", id, err)
    }
    
    rowsAffected := commentModel.Delete()
    if rowsAffected == 0 {
        return repositories.NewDeleteError("评论", id, nil)
    }
    return nil
}
```

### 3. Service 层错误处理

Service 层使用 `apperrors` 包创建业务错误：

```go
func (s *RoleService) CreateRole(dto RoleCreateDTO) (*RoleResponseDTO, error) {
    // 检查是否存在
    existingRole, _ := s.repo.GetByName(dto.Name)
    if existingRole != nil {
        return nil, apperrors.ConflictError("角色")
    }

    newRole := &role.Role{
        Name:        dto.Name,
        DisplayName: dto.DisplayName,
        Description: dto.Description,
    }

    // 处理数据库错误
    if err := s.repo.Create(newRole); err != nil {
        return nil, apperrors.DatabaseCreateError("角色", err)
    }

    resp := toRoleResponseDTO(newRole)
    return &resp, nil
}

func (s *RoleService) GetRoleByID(id uint64) (*RoleResponseDTO, error) {
    role, err := s.repo.GetByID(id)
    if err != nil {
        // 使用特定错误码
        return nil, apperrors.NotFoundErrorWithCode(apperrors.CodeRoleNotFound, "角色")
    }

    resp := toRoleResponseDTO(role)
    return &resp, nil
}
```

### 4. 错误判断 - errors.Is()

使用 `errors.Is()` 判断错误类型：

```go
import (
    "errors"
    apperrors "GoHub-Service/pkg/errors"
)

// 判断是否为特定错误
if errors.Is(err, apperrors.ErrNotFound) {
    // 处理资源不存在
}

// 判断是否为数据库错误
if errors.Is(err, apperrors.ErrDatabaseError) {
    // 处理数据库错误
}

// 判断是否为特定错误码
if errors.Is(err, &apperrors.AppError{Code: apperrors.CodeUserNotFound}) {
    // 处理用户不存在
}
```

#### 预定义错误常量

```go
var (
    ErrNotFound          = &AppError{Type: ErrorTypeNotFound, Code: CodeNotFound}
    ErrUnauthorized      = &AppError{Type: ErrorTypeAuthorization, Code: CodeUnauthorized}
    ErrForbidden         = &AppError{Type: ErrorTypeAuthorization, Code: CodeForbidden}
    ErrInvalidParameter  = &AppError{Type: ErrorTypeValidation, Code: CodeInvalidParameter}
    ErrConflict          = &AppError{Type: ErrorTypeConflict, Code: CodeConflict}
    ErrDatabaseError     = &AppError{Type: ErrorTypeDatabase, Code: CodeDatabaseError}
    ErrInternalError     = &AppError{Type: ErrorTypeInternal, Code: CodeInternalError}
    ErrTimeout           = &AppError{Type: ErrorTypeTimeout, Code: CodeTimeout}
)
```

### 5. 错误类型转换 - errors.As()

使用 `errors.As()` 获取 AppError 详细信息：

```go
var appErr *apperrors.AppError
if errors.As(err, &appErr) {
    // 获取错误码
    code := appErr.GetCode()
    
    // 获取错误类型
    errorType := appErr.GetType()
    
    // 获取错误详情
    details := appErr.Details
    
    // 获取请求ID
    requestID := appErr.RequestID
    
    // 获取堆栈信息
    stackTrace := appErr.StackTrace
}
```

### 6. 错误包装

保留原始错误信息的同时添加上下文：

```go
// 方式1：使用 WrapError
err := apperrors.WrapError(dbErr, "获取用户列表失败")

// 方式2：使用 WithError
err := apperrors.NotFoundError("用户").WithError(originalErr)

// 添加详情
err := apperrors.DatabaseError("查询", dbErr).
    WithDetails(map[string]interface{}{
        "table": "users",
        "query": "SELECT * FROM users WHERE id = ?",
    }).
    WithRequestID(requestID)
```

### 7. Controller 层错误响应

```go
func (ctl *UsersController) Show(c *gin.Context) {
    ctx := ctx.FromGinContext(c)
    
    user, err := ctl.service.GetByID(ctx, id)
    if err != nil {
        // 提取 AppError
        var appErr *apperrors.AppError
        if errors.As(err, &appErr) {
            // 根据错误码返回不同状态码
            statusCode := getHTTPStatusCode(appErr.Code)
            c.JSON(statusCode, gin.H{
                "code":    appErr.Code,
                "message": appErr.Message,
                "details": appErr.Details,
            })
            return
        }
        
        // 未知错误
        c.JSON(500, gin.H{
            "code":    apperrors.CodeInternalError,
            "message": "服务器内部错误",
        })
        return
    }
    
    c.JSON(200, user)
}

// 错误码到HTTP状态码映射
func getHTTPStatusCode(code int) int {
    switch code {
    case apperrors.CodeNotFound:
        return 404
    case apperrors.CodeUnauthorized:
        return 401
    case apperrors.CodeForbidden:
        return 403
    case apperrors.CodeValidationError, apperrors.CodeInvalidParameter:
        return 400
    case apperrors.CodeConflict:
        return 409
    default:
        return 500
    }
}
```

---

## 最佳实践

### 1. 错误创建原则

✅ **推荐做法：**
```go
// 使用具体的错误创建函数
err := apperrors.NotFoundError("用户")
err := apperrors.DatabaseCreateError("用户", dbErr)
err := repositories.NewNotFoundError("用户", userID)
```

❌ **不推荐：**
```go
// 不要使用 fmt.Errorf
err := fmt.Errorf("用户不存在")

// 不要使用 errors.New
err := errors.New("创建失败")
```

### 2. 错误判断原则

✅ **推荐做法：**
```go
// 使用 errors.Is 判断
if errors.Is(err, apperrors.ErrNotFound) {
    // 处理
}

// 使用 errors.As 获取详情
var appErr *apperrors.AppError
if errors.As(err, &appErr) {
    code := appErr.GetCode()
}
```

❌ **不推荐：**
```go
// 不要使用字符串比较
if err.Error() == "用户不存在" {
    // 处理
}

// 不要使用类型断言（不安全）
if appErr, ok := err.(*apperrors.AppError); ok {
    // 处理
}
```

### 3. 错误信息原则

✅ **推荐做法：**
```go
// 提供清晰的错误消息
err := apperrors.NotFoundError("用户").
    WithDetails(map[string]interface{}{
        "user_id": userID,
        "action":  "login",
    })

// 包装原始错误
err := apperrors.DatabaseCreateError("用户", dbErr)
```

❌ **不推荐：**
```go
// 不要丢失原始错误信息
err := apperrors.NotFoundError("资源")  // 太笼统

// 不要暴露敏感信息
err := apperrors.ValidationError("密码错误", map[string]interface{}{
    "password": "123456",  // ❌ 不要记录密码
})
```

### 4. 错误传递原则

✅ **推荐做法：**
```go
// Repository → Service：直接返回
func (r *userRepository) GetByID(id string) (*user.User, error) {
    user := user.Get(id)
    if user.ID == 0 {
        return nil, repositories.NewNotFoundError("用户", id)
    }
    return &user, nil
}

// Service → Controller：添加业务上下文
func (s *UserService) GetByID(ctx context.Context, id string) (*UserDTO, error) {
    user, err := s.repo.GetByID(id)
    if err != nil {
        // 直接返回或添加上下文
        return nil, err
    }
    return toDTO(user), nil
}
```

❌ **不推荐：**
```go
// 不要重复包装
func (s *UserService) GetByID(id string) (*UserDTO, error) {
    user, err := s.repo.GetByID(id)
    if err != nil {
        // ❌ 不要重复创建错误
        return nil, apperrors.NotFoundError("用户")
    }
    return toDTO(user), nil
}
```

---

## 错误日志记录

### 日志级别建议

| 错误类型 | 日志级别 | 说明 |
|---------|---------|------|
| 资源不存在 | INFO/WARN | 可能是正常业务流程 |
| 参数验证失败 | INFO | 客户端错误 |
| 授权失败 | WARN | 可能的安全问题 |
| 数据库错误 | ERROR | 需要关注 |
| 外部服务错误 | ERROR | 需要关注 |
| 内部错误 | FATAL | 严重问题 |

### 日志记录示例

```go
import (
    "GoHub-Service/pkg/logger"
    apperrors "GoHub-Service/pkg/errors"
)

func (s *UserService) CreateUser(ctx context.Context, dto UserDTO) error {
    user, err := s.repo.Create(ctx, toModel(dto))
    if err != nil {
        var appErr *apperrors.AppError
        if errors.As(err, &appErr) {
            // 记录结构化日志
            logger.Error("创建用户失败", map[string]interface{}{
                "error_code":  appErr.Code,
                "error_type":  appErr.Type,
                "user_email":  dto.Email,
                "request_id":  appErr.RequestID,
                "stack_trace": appErr.StackTrace,
            })
        }
        return err
    }
    return nil
}
```

---

## 性能考虑

### 堆栈追踪

堆栈追踪对性能有一定影响，建议：

1. **生产环境**：可选择性关闭或采样
2. **开发环境**：始终开启
3. **关键错误**：始终保留堆栈

### 错误详情

避免在错误详情中存储大对象：

```go
// ✅ 推荐
err.WithDetails(map[string]interface{}{
    "user_id": userID,
    "action":  "login",
})

// ❌ 不推荐
err.WithDetails(map[string]interface{}{
    "user_object": entireUserObject,  // 可能很大
    "query_result": largeQueryResult,  // 可能很大
})
```

---

## 迁移指南

### 从旧错误系统迁移

**步骤1：更新导入**
```go
// 旧代码
import "errors"
import "fmt"

// 新代码
import (
    "errors"
    apperrors "GoHub-Service/pkg/errors"
)
```

**步骤2：替换错误创建**
```go
// 旧代码
return fmt.Errorf("用户不存在")

// 新代码
return apperrors.NotFoundError("用户")
```

**步骤3：更新错误判断**
```go
// 旧代码
if err.Error() == "resource not found" {
    // 处理
}

// 新代码
if errors.Is(err, apperrors.ErrNotFound) {
    // 处理
}
```

---

## 常见错误码速查

| 错误码 | 常量 | 说明 |
|-------|------|------|
| 0 | CodeSuccess | 成功 |
| 1001 | CodeInvalidParameter | 参数错误 |
| 1002 | CodeUnauthorized | 未授权 |
| 1003 | CodeForbidden | 禁止访问 |
| 1004 | CodeNotFound | 资源不存在 |
| 1009 | CodeConflict | 冲突 |
| 2001 | CodeDatabaseError | 数据库错误 |
| 2004 | CodeDatabaseCreate | 创建失败 |
| 2005 | CodeDatabaseUpdate | 更新失败 |
| 2006 | CodeDatabaseDelete | 删除失败 |
| 2007 | CodeDatabaseDuplicate | 重复记录 |
| 2101 | CodeCacheError | 缓存错误 |
| 2202 | CodeTimeout | 超时 |
| 4001 | CodeUserNotFound | 用户不存在 |
| 4101 | CodeTopicNotFound | 话题不存在 |
| 4201 | CodeCommentNotFound | 评论不存在 |
| 4301 | CodeCategoryNotFound | 分类不存在 |
| 4401 | CodeRoleNotFound | 角色不存在 |

---

## 总结

完善的错误处理系统提供了：

✅ **统一的错误码管理**  
✅ **错误链支持（Unwrap）**  
✅ **错误判断（Is）和类型转换（As）**  
✅ **详细的上下文信息**  
✅ **堆栈追踪**  
✅ **请求ID关联**  
✅ **结构化错误详情**  

通过遵循本指南的最佳实践，可以：

- 提高代码可维护性
- 更好地追踪和调试问题
- 提供更友好的错误信息
- 支持更精确的错误处理逻辑
