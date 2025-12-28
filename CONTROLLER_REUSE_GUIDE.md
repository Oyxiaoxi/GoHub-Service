# Controller代码复用指南

> 更新时间：2025年12月28日  
> 版本：v1.0

本文档说明如何使用新的CRUD助手和授权中间件来减少Controller中的重复代码。

---

## 📋 目录

1. [CRUD助手使用](#crud助手使用)
2. [授权中间件使用](#授权中间件使用)
3. [完整示例](#完整示例)
4. [迁移指南](#迁移指南)

---

## CRUD助手使用

### 1. 基本用法

#### 1.1 创建CRUD助手实例

```go
import "GoHub-Service/pkg/controller"

// 在Controller中创建助手
crudHelper := controller.NewCRUDHelper("话题")
```

#### 1.2 使用HandleShow

```go
func (ctrl *TopicsController) Show(c *gin.Context) {
    topicModel := topic.Get(c.Param("id"))
    ctrl.crudHelper.HandleShow(c, &topicModel)
}
```

**替代原有代码**：
```go
// 旧代码
if topicModel.ID == 0 {
    response.Abort404(c)
    return
}
response.Data(c, topicModel)
```

#### 1.3 使用HandleStore

```go
func (ctrl *TopicsController) Store(c *gin.Context) {
    request := requests.TopicRequest{}
    if ok := requests.Validate(c, &request, requests.TopicSave); !ok {
        return
    }

    topicModel := topic.Topic{
        Title:      request.Title,
        Body:       request.Body,
        CategoryID: request.CategoryID,
        UserID:     auth.CurrentUID(c),
    }
    
    ctrl.crudHelper.HandleStore(c, &topicModel)
}
```

**替代原有代码**：
```go
// 旧代码
topicModel.Create()
if topicModel.ID > 0 {
    response.Created(c, topicModel)
} else {
    response.Abort500(c, "创建失败，请稍后尝试~")
}
```

#### 1.4 使用HandleUpdate

```go
func (ctrl *TopicsController) Update(c *gin.Context) {
    topicModel := topic.Get(c.Param("id"))
    
    request := requests.TopicRequest{}
    if ok := requests.Validate(c, &request, requests.TopicSave); !ok {
        return
    }
    
    // 检查权限
    if ok := middlewares.CheckModelOwnership(c, &topicModel); !ok {
        return
    }
    
    // 更新字段
    topicModel.Title = request.Title
    topicModel.Body = request.Body
    topicModel.CategoryID = request.CategoryID
    
    ctrl.crudHelper.HandleUpdate(c, &topicModel)
}
```

#### 1.5 使用HandleDelete

```go
func (ctrl *TopicsController) Delete(c *gin.Context) {
    topicModel := topic.Get(c.Param("id"))
    
    // 检查权限
    if ok := middlewares.CheckModelOwnership(c, &topicModel); !ok {
        return
    }
    
    ctrl.crudHelper.HandleDelete(c, &topicModel)
}
```

#### 1.6 使用HandleList

```go
func (ctrl *TopicsController) Index(c *gin.Context) {
    request := requests.PaginationRequest{}
    if ok := requests.Validate(c, &request, requests.Pagination); !ok {
        return
    }

    data, pager := topic.Paginate(c, 10)
    ctrl.crudHelper.HandleList(c, data, pager)
}
```

---

## 授权中间件使用

### 1. 在Controller中使用CheckModelOwnership

#### 示例：Update方法

```go
func (ctrl *TopicsController) Update(c *gin.Context) {
    topicModel := topic.Get(c.Param("id"))
    
    // 使用通用的所有权检查
    if ok := middlewares.CheckModelOwnership(c, &topicModel); !ok {
        return // 已自动返回403错误
    }
    
    // ... 后续更新逻辑
}
```

**替代原有代码**：
```go
// 旧代码
if ok := policies.CanModifyTopic(c, topicModel); !ok {
    response.Abort403(c)
    return
}
```

### 2. 模型实现OwnershipChecker接口

要使用`CheckModelOwnership`，模型需要实现`OwnershipChecker`接口：

```go
// 在模型中添加此方法
func (topic *Topic) GetOwnerID() string {
    return topic.UserID
}
```

### 3. 在路由中使用中间件（可选）

如果想在路由级别检查权限，可以使用中间件：

```go
// routes/api.go

// 方式1：在特定路由使用
topicsGroup.PUT("/:id", middlewares.CheckOwnership(func(c *gin.Context) string {
    topicModel := topic.Get(c.Param("id"))
    return topicModel.UserID
}), topicsController.Update)

// 方式2：使用策略检查
topicsGroup.PUT("/:id", middlewares.CheckPolicy(
    func(c *gin.Context, model interface{}) bool {
        return policies.CanModifyTopic(c, model.(topic.Topic))
    },
    func(c *gin.Context) interface{} {
        return topic.Get(c.Param("id"))
    },
), topicsController.Update)
```

---

## 完整示例

### 优化前的TopicsController

```go
package v1

import (
    "GoHub-Service/app/models/topic"
    "GoHub-Service/app/policies"
    "GoHub-Service/app/requests"
    "GoHub-Service/pkg/auth"
    "GoHub-Service/pkg/response"
    "github.com/gin-gonic/gin"
)

type TopicsController struct {
    BaseAPIController
}

func (ctrl *TopicsController) Show(c *gin.Context) {
    topicModel := topic.Get(c.Param("id"))
    if topicModel.ID == 0 {
        response.Abort404(c)
        return
    }
    response.Data(c, topicModel)
}

func (ctrl *TopicsController) Update(c *gin.Context) {
    topicModel := topic.Get(c.Param("id"))
    if topicModel.ID == 0 {
        response.Abort404(c)
        return
    }

    request := requests.TopicRequest{}
    if ok := requests.Validate(c, &request, requests.TopicSave); !ok {
        return
    }
    
    if ok := policies.CanModifyTopic(c, topicModel); !ok {
        response.Abort403(c)
        return
    }

    topicModel.Title = request.Title
    topicModel.Body = request.Body
    topicModel.CategoryID = request.CategoryID
    rowsAffected := topicModel.Save()
    if rowsAffected > 0 {
        response.Data(c, topicModel)
    } else {
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}
```

### 优化后的TopicsController

```go
package v1

import (
    "GoHub-Service/app/http/middlewares"
    "GoHub-Service/app/models/topic"
    "GoHub-Service/app/requests"
    "GoHub-Service/pkg/auth"
    "GoHub-Service/pkg/controller"
    "github.com/gin-gonic/gin"
)

type TopicsController struct {
    BaseAPIController
    crudHelper *controller.CRUDHelper
}

func NewTopicsController() *TopicsController {
    return &TopicsController{
        crudHelper: controller.NewCRUDHelper("话题"),
    }
}

func (ctrl *TopicsController) Show(c *gin.Context) {
    topicModel := topic.Get(c.Param("id"))
    ctrl.crudHelper.HandleShow(c, &topicModel)
}

func (ctrl *TopicsController) Update(c *gin.Context) {
    topicModel := topic.Get(c.Param("id"))
    
    request := requests.TopicRequest{}
    if ok := requests.Validate(c, &request, requests.TopicSave); !ok {
        return
    }
    
    // 统一的所有权检查
    if ok := middlewares.CheckModelOwnership(c, &topicModel); !ok {
        return
    }

    // 更新字段
    topicModel.Title = request.Title
    topicModel.Body = request.Body
    topicModel.CategoryID = request.CategoryID
    
    // 统一的更新处理
    ctrl.crudHelper.HandleUpdate(c, &topicModel)
}
```

**代码减少**：约30-40%的重复代码

---

## 迁移指南

### 步骤1：为模型添加接口实现

```go
// 1. 实现Model接口（如果使用CRUD助手）
func (model *YourModel) GetID() uint64 {
    return model.ID
}

// 2. 实现OwnershipChecker接口（如果需要权限检查）
func (model *YourModel) GetOwnerID() string {
    return model.UserID
}
```

### 步骤2：更新Controller

```go
// 1. 添加crudHelper字段
type YourController struct {
    BaseAPIController
    crudHelper *controller.CRUDHelper
}

// 2. 创建构造函数
func NewYourController() *YourController {
    return &YourController{
        crudHelper: controller.NewCRUDHelper("资源名称"),
    }
}

// 3. 逐个替换方法
```

### 步骤3：替换权限检查

```go
// 旧代码
if ok := policies.CanModifyXxx(c, model); !ok {
    response.Abort403(c)
    return
}

// 新代码
if ok := middlewares.CheckModelOwnership(c, &model); !ok {
    return
}
```

---

## 优点总结

### 1. 代码复用
- ✅ 减少30-40%的重复代码
- ✅ 统一错误处理逻辑
- ✅ 统一响应格式

### 2. 可维护性
- ✅ 修改一处，全局生效
- ✅ 代码更简洁易读
- ✅ 减少出错可能

### 3. 一致性
- ✅ 所有CRUD操作行为一致
- ✅ 所有权限检查逻辑一致
- ✅ 错误消息格式一致

### 4. 扩展性
- ✅ 易于添加新功能
- ✅ 支持自定义策略
- ✅ 灵活的中间件组合

---

## 注意事项

1. **渐进式迁移**：不需要一次性修改所有Controller，可以逐步迁移
2. **向后兼容**：旧的写法仍然可用，不影响现有代码
3. **灵活使用**：可以根据具体需求选择性使用助手方法
4. **接口实现**：确保模型实现了必要的接口

---

**相关文档**：
- [代码规范文档](CODING_STANDARDS.md)
- [优化计划](OPTIMIZATION_PLAN.md)
- [Service层架构指南](docs/SERVICE_LAYER_GUIDE.md)
