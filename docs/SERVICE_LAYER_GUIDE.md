# Service 层架构指南

## 概述

Service 层是业务逻辑的核心层，位于 Controller 和 Model 之间。通过引入 Service 层，我们实现了：

- **业务逻辑集中化**：将复杂的业务逻辑从 Controller 中抽离
- **代码复用**：多个 Controller 可以共享同一个 Service
- **易于测试**：Service 层可以独立进行单元测试
- **清晰的职责划分**：Controller 负责 HTTP 处理，Service 负责业务逻辑

## 目录结构

```
app/
├── services/
│   ├── topic_service.go      # 话题业务逻辑
│   └── ...                    # 其他业务服务
├── http/
│   └── controllers/
│       └── api/v1/
│           ├── topics_controller.go  # 使用 TopicService
│           └── ...
└── models/
    └── topic/
        └── topic.go           # 数据模型
```

## 核心组件

### 1. 自定义错误类型 (pkg/errors/errors.go)

提供统一的错误处理机制：

```go
// 错误类型
const (
    ErrorTypeBusiness      = "BUSINESS_ERROR"       // 业务逻辑错误
    ErrorTypeValidation    = "VALIDATION_ERROR"     // 数据验证错误
    ErrorTypeAuthorization = "AUTHORIZATION_ERROR"  // 授权错误
    ErrorTypeNotFound      = "NOT_FOUND"           // 资源未找到
    ErrorTypeDatabase      = "DATABASE_ERROR"       // 数据库错误
    ErrorTypeExternal      = "EXTERNAL_ERROR"       // 外部服务错误
    ErrorTypeInternal      = "INTERNAL_ERROR"       // 内部系统错误
)

// 使用示例
func (s *TopicService) GetByID(id string) (*topic.Topic, error) {
    topicModel := topic.Get(id)
    if topicModel.ID == 0 {
        return nil, apperrors.NotFoundError("话题不存在", map[string]interface{}{
            "topic_id": id,
        })
    }
    return &topicModel, nil
}
```

### 2. RequestID 追踪 (app/http/middlewares/request_id.go)

为每个请求生成唯一标识，便于追踪和调试：

```go
// 中间件配置（在 routes/api.go）
v1.Use(middlewares.RequestID())

// 获取 RequestID
requestID := middlewares.GetRequestID(c)

// 错误日志自动包含 RequestID
logger.LogErrorWithContext(c, err, "操作失败")
```

### 3. 上下文日志 (pkg/logger/context.go)

提供带上下文的日志记录：

```go
// 记录错误日志（自动包含 RequestID 和错误详情）
logger.LogErrorWithContext(c, err, "创建话题失败",
    zap.String("title", request.Title),
    zap.String("user_id", auth.CurrentUID(c)),
)

// 普通日志（包含 RequestID）
logger.LogWithRequestID(c, "info", "处理请求",
    zap.String("action", "create_topic"),
)
```

## Service 层实现指南

### 创建 Service

```go
package services

import (
    "GoHub-Service/app/models/topic"
    apperrors "GoHub-Service/pkg/errors"
)

// TopicService 话题服务
type TopicService struct{}

// NewTopicService 创建服务实例
func NewTopicService() *TopicService {
    return &TopicService{}
}
```

### 定义 DTO (Data Transfer Object)

DTO 用于在不同层之间传递数据：

```go
// TopicCreateDTO 创建话题数据传输对象
type TopicCreateDTO struct {
    Title      string
    Body       string
    CategoryID string
    UserID     string
}

// TopicUpdateDTO 更新话题数据传输对象
type TopicUpdateDTO struct {
    Title      string
    Body       string
    CategoryID string
}
```

### 实现业务方法

```go
// GetByID 根据ID获取话题
func (s *TopicService) GetByID(id string) (*topic.Topic, error) {
    topicModel := topic.Get(id)
    if topicModel.ID == 0 {
        return nil, apperrors.NotFoundError("话题不存在", map[string]interface{}{
            "topic_id": id,
        })
    }
    return &topicModel, nil
}

// Create 创建话题
func (s *TopicService) Create(dto TopicCreateDTO) (*topic.Topic, error) {
    topicModel := topic.Topic{
        Title:      dto.Title,
        Body:       dto.Body,
        CategoryID: dto.CategoryID,
        UserID:     dto.UserID,
    }

    topicModel.Create()
    if topicModel.ID == 0 {
        return nil, apperrors.DatabaseError("创建话题失败", nil, nil)
    }

    return &topicModel, nil
}

// Update 更新话题
func (s *TopicService) Update(id string, dto TopicUpdateDTO) (*topic.Topic, error) {
    topicModel := topic.Get(id)
    if topicModel.ID == 0 {
        return nil, apperrors.NotFoundError("话题不存在", map[string]interface{}{
            "topic_id": id,
        })
    }

    topicModel.Title = dto.Title
    topicModel.Body = dto.Body
    topicModel.CategoryID = dto.CategoryID

    rowsAffected := topicModel.Save()
    if rowsAffected == 0 {
        return nil, apperrors.DatabaseError("更新话题失败", nil, nil)
    }

    return &topicModel, nil
}

// Delete 删除话题
func (s *TopicService) Delete(id string) error {
    topicModel := topic.Get(id)
    if topicModel.ID == 0 {
        return apperrors.NotFoundError("话题不存在", map[string]interface{}{
            "topic_id": id,
        })
    }

    rowsAffected := topicModel.Delete()
    if rowsAffected == 0 {
        return apperrors.DatabaseError("删除话题失败", nil, nil)
    }

    return nil
}

// CheckOwnership 检查用户是否拥有该话题
func (s *TopicService) CheckOwnership(topicID, userID string) (bool, error) {
    topicModel := topic.Get(topicID)
    if topicModel.ID == 0 {
        return false, apperrors.NotFoundError("话题不存在", map[string]interface{}{
            "topic_id": topicID,
        })
    }
    return topicModel.UserID == userID, nil
}
```

## Controller 层使用 Service

### 1. 初始化 Controller

```go
type TopicsController struct {
    BaseAPIController
    topicService *services.TopicService
}

// NewTopicsController 创建控制器实例
func NewTopicsController() *TopicsController {
    return &TopicsController{
        topicService: services.NewTopicService(),
    }
}
```

### 2. 在路由中使用

```go
// routes/api.go
tpc := controllers.NewTopicsController()
tpcGroup := v1.Group("/topics")
{
    tpcGroup.GET("", tpc.Index)
    tpcGroup.POST("", middlewares.AuthJWT(), tpc.Store)
    tpcGroup.PUT("/:id", middlewares.AuthJWT(), tpc.Update)
    tpcGroup.DELETE("/:id", middlewares.AuthJWT(), tpc.Delete)
    tpcGroup.GET("/:id", tpc.Show)
}
```

### 3. Controller 方法实现

```go
func (ctrl *TopicsController) Show(c *gin.Context) {
    // 调用 Service 获取数据
    topicModel, err := ctrl.topicService.GetByID(c.Param("id"))
    
    // 错误处理
    if err != nil {
        if apperrors.IsAppError(err) {
            appErr := apperrors.GetAppError(err)
            appErr.WithRequestID(middlewares.GetRequestID(c))
            response.Abort404(c)
            return
        }
        logger.LogErrorWithContext(c, err, "获取话题失败")
        response.Abort500(c)
        return
    }
    
    // 返回成功响应
    response.Data(c, topicModel)
}

func (ctrl *TopicsController) Store(c *gin.Context) {
    // 验证请求
    request := requests.TopicRequest{}
    if ok := requests.Validate(c, &request, requests.TopicSave); !ok {
        return
    }

    // 构建 DTO
    dto := services.TopicCreateDTO{
        Title:      request.Title,
        Body:       request.Body,
        CategoryID: request.CategoryID,
        UserID:     auth.CurrentUID(c),
    }

    // 调用 Service 创建
    topicModel, err := ctrl.topicService.Create(dto)
    if err != nil {
        logger.LogErrorWithContext(c, err, "创建话题失败",
            zap.String("title", request.Title),
            zap.String("user_id", auth.CurrentUID(c)),
        )
        response.Abort500(c, "创建失败，请稍后尝试~")
        return
    }

    response.Created(c, topicModel)
}

func (ctrl *TopicsController) Update(c *gin.Context) {
    // 验证请求
    request := requests.TopicRequest{}
    if ok := requests.Validate(c, &request, requests.TopicSave); !ok {
        return
    }

    // 检查所有权
    topicID := c.Param("id")
    currentUserID := auth.CurrentUID(c)
    
    isOwner, err := ctrl.topicService.CheckOwnership(topicID, currentUserID)
    if err != nil {
        if apperrors.IsAppError(err) {
            response.Abort404(c)
            return
        }
        logger.LogErrorWithContext(c, err, "检查话题所有权失败")
        response.Abort500(c)
        return
    }
    
    if !isOwner {
        response.Abort403(c, "无权限操作")
        return
    }

    // 构建 DTO 并更新
    dto := services.TopicUpdateDTO{
        Title:      request.Title,
        Body:       request.Body,
        CategoryID: request.CategoryID,
    }

    topicModel, err := ctrl.topicService.Update(topicID, dto)
    if err != nil {
        logger.LogErrorWithContext(c, err, "更新话题失败",
            zap.String("topic_id", topicID),
        )
        response.Abort500(c, "更新失败，请稍后尝试~")
        return
    }

    response.Data(c, topicModel)
}
```

## 错误处理流程

```
┌─────────────┐
│  Controller │
└──────┬──────┘
       │ 调用 Service
       ▼
┌─────────────┐
│   Service   │ ──┐ 业务逻辑错误
└──────┬──────┘   │
       │          ▼
       │    返回 AppError
       │          │
       │          │ 包含:
       │          │ - ErrorType
       │          │ - ErrorCode
       │          │ - Message
       │          │ - Details
       │          │ - StackTrace
       │          │
       ▼          ▼
┌─────────────────────────┐
│  Controller 错误处理    │
│  - 识别错误类型         │
│  - 记录上下文日志       │
│  - 返回适当的HTTP响应   │
└─────────────────────────┘
       │
       ▼
┌─────────────────────────┐
│  Logger 记录            │
│  - RequestID           │
│  - ErrorType           │
│  - ErrorCode           │
│  - StackTrace          │
│  - 业务上下文           │
└─────────────────────────┘
```

## 最佳实践

### 1. Service 职责

✅ **应该做**：
- 实现业务逻辑
- 数据验证（业务规则级别）
- 调用 Model 层进行数据操作
- 处理事务
- 返回明确的错误类型

❌ **不应该做**：
- 处理 HTTP 请求和响应
- 访问 gin.Context

## 事务与缓存示例

多步骤写操作请在 Service 层使用事务，同时在成功后刷新缓存：

```go
func (s *TopicService) CreateWithCache(dto TopicCreateDTO) (*topic.Topic, error) {
    var created topic.Topic
    err := database.DB.Transaction(func(tx *gorm.DB) error {
        t := topic.Topic{Title: dto.Title, Body: dto.Body, CategoryID: dto.CategoryID, UserID: dto.UserID}
        if err := tx.Create(&t).Error; err != nil {
            return err
        }
        created = t
        return nil
    })
    if err != nil {
        return nil, apperrors.DatabaseError("创建话题失败", nil, err)
    }

    // 写操作完成后刷新缓存
    cache.TopicStore(created.GetStringID(), created)
    return &created, nil
}
```

提示：
- 不要在 Service 中持有 gin.Context；必要上下文通过参数传入（如 userID、locale）。
- 将 DTO → Model 转换、缓存刷新、错误包装都集中在 Service，Controller 仅做绑定与响应。
- 直接返回 HTTP 状态码
- 处理表单验证（应在 Controller 层）

### 2. 错误处理

```go
// ✅ 好的做法
func (s *TopicService) GetByID(id string) (*topic.Topic, error) {
    topicModel := topic.Get(id)
    if topicModel.ID == 0 {
        return nil, apperrors.NotFoundError("话题不存在", map[string]interface{}{
            "topic_id": id,
        })
    }
    return &topicModel, nil
}

// ❌ 不好的做法
func (s *TopicService) GetByID(c *gin.Context, id string) {
    topicModel := topic.Get(id)
    if topicModel.ID == 0 {
        response.Abort404(c) // 不应在 Service 层处理 HTTP 响应
        return
    }
    response.Data(c, topicModel)
}
```

### 3. DTO 使用

```go
// ✅ 好的做法：使用 DTO 传递数据
dto := services.TopicCreateDTO{
    Title:      request.Title,
    Body:       request.Body,
    CategoryID: request.CategoryID,
    UserID:     auth.CurrentUID(c),
}
topicModel, err := ctrl.topicService.Create(dto)

// ❌ 不好的做法：直接传递 request 对象
topicModel, err := ctrl.topicService.Create(request) // Service 不应依赖 HTTP 请求对象
```

### 4. 日志记录

```go
// ✅ 好的做法：在 Controller 层记录日志，包含完整上下文
topicModel, err := ctrl.topicService.Create(dto)
if err != nil {
    logger.LogErrorWithContext(c, err, "创建话题失败",
        zap.String("title", request.Title),
        zap.String("user_id", auth.CurrentUID(c)),
    )
    response.Abort500(c, "创建失败，请稍后尝试~")
    return
}

// ❌ 不好的做法：在 Service 层访问 gin.Context
func (s *TopicService) Create(c *gin.Context, dto TopicCreateDTO) {
    // Service 不应该依赖 gin.Context
    logger.LogWithRequestID(c, "info", "创建话题")
}
```

## 迁移指南

### 从传统 Controller 迁移到 Service 层

**Before（传统方式）**：
```go
func (ctrl *TopicsController) Show(c *gin.Context) {
    topicModel := topic.Get(c.Param("id"))
    if topicModel.ID == 0 {
        response.Abort404(c)
        return
    }
    response.Data(c, topicModel)
}
```

**After（使用 Service 层）**：
```go
func (ctrl *TopicsController) Show(c *gin.Context) {
    topicModel, err := ctrl.topicService.GetByID(c.Param("id"))
    if err != nil {
        if apperrors.IsAppError(err) {
            appErr := apperrors.GetAppError(err)
            appErr.WithRequestID(middlewares.GetRequestID(c))
            response.Abort404(c)
            return
        }
        logger.LogErrorWithContext(c, err, "获取话题失败")
        response.Abort500(c)
        return
    }
    response.Data(c, topicModel)
}
```

## 总结

通过引入 Service 层架构，我们实现了：

1. **清晰的分层架构**：Controller → Service → Model
2. **统一的错误处理**：自定义错误类型 + 完整的错误追踪
3. **可追踪的请求链路**：RequestID 贯穿整个请求生命周期
4. **上下文感知的日志**：每条日志都包含 RequestID 和业务上下文
5. **易于测试和维护**：业务逻辑集中在 Service 层

这些改进大大提升了代码的可维护性、可测试性和可追溯性。
