# API 设计优化 v4.0 - 快速使用指南

## 🎯 优化概览

API 设计优化已完成，提供以下核心功能：
1. ✅ **API 版本管理** - 完整的版本生命周期管理
2. ✅ **统一响应格式** - 标准化的 API 响应结构
3. ✅ **OpenAPI 文档** - 交互式 Swagger 文档

---

## 1. API 版本管理

### 访问版本信息
```bash
# 获取所有支持的 API 版本
curl http://localhost:3000/api/versions
```

响应示例：
```json
{
  "current_version": "v1",
  "versions": {
    "v1": {
      "version": "v1",
      "status": "active",
      "release_date": "2024-01-01",
      "features": ["用户管理", "话题管理", "评论管理"]
    }
  },
  "api_docs": "/swagger/index.html"
}
```

### 使用特定版本

**方式1：URL 路径（推荐）**
```bash
GET /api/v1/users
```

**方式2：请求头**
```bash
GET /api/users
X-API-Version: v1
```

### 版本状态
- **active** - 活跃版本，完全支持
- **deprecated** - 已废弃（响应头会包含警告）
- **sunset** - 已停用（返回 410 Gone）
- **planned** - 计划中

**文档**: [docs/23_API_VERSIONING.md](./23_API_VERSIONING.md)

---

## 2. 统一响应格式

### 新增标准响应函数

在 `pkg/response/standard.go` 中提供：

#### 成功响应
```go
// 简单成功
response.StandardSuccess(c, userData)

// 带自定义消息
response.StandardSuccessWithMessage(c, "创建成功", userData)

// 带分页信息
response.StandardSuccessWithMeta(c, userList, &response.MetaInfo{
    CurrentPage: 1,
    PerPage:     20,
    Total:       100,
    TotalPages:  5,
})
```

#### 错误响应
```go
// 通用错误
response.StandardError(c, http.StatusNotFound, response.CodeUserNotFound, "用户不存在")

// 验证错误
response.StandardValidationError(c, map[string]string{
    "email": "邮箱格式不正确",
    "phone": "手机号已存在",
})
```

### 响应格式示例

**成功响应**：
```json
{
  "success": true,
  "code": 20000,
  "message": "success",
  "data": {...},
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  },
  "timestamp": 1704067200,
  "request_id": "abc123",
  "version": "v1"
}
```

**错误响应**：
```json
{
  "success": false,
  "code": 40400,
  "message": "User not found",
  "error": {
    "type": "not_found",
    "details": "User with ID 999 not found",
    "fields": {...}
  },
  "timestamp": 1704067200,
  "request_id": "abc123"
}
```

### 向后兼容
- 旧的响应函数（`response.Success`, `response.Data` 等）**仍然可用**
- 新代码推荐使用 `response.Standard*` 系列函数
- 两种响应格式可以并存，逐步迁移

---

## 3. OpenAPI 文档（Swagger）

### 访问文档
```
http://localhost:3000/swagger/index.html
```

### 生成文档

#### 方式1：使用 Makefile
```bash
make -f Makefile.swagger swagger
```

#### 方式2：直接使用 swag
```bash
# 先安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init --parseDependency --parseInternal
```

### 添加 Swagger 注解

在控制器方法上添加注解：

```go
// CurrentUser 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 返回当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /user [get]
func (ctrl *UsersController) CurrentUser(c *gin.Context) {
    // 方法实现
}
```

### 常用注解标签

| 标签 | 说明 | 示例 |
|------|------|------|
| @Summary | 简短描述 | 获取用户信息 |
| @Description | 详细描述 | 根据用户ID获取详细信息 |
| @Tags | 分组标签 | 用户管理 |
| @Param | 参数定义 | id path string true "用户ID" |
| @Success | 成功响应 | 200 {object} User |
| @Failure | 失败响应 | 404 {object} Error |
| @Security | 安全认证 | Bearer |

**完整指南**: [docs/24_OPENAPI_GUIDE.md](./24_OPENAPI_GUIDE.md)

---

## 4. 快速开始

### Step 1: 启动应用
```bash
go run main.go serve
```

### Step 2: 查看 API 版本
```bash
curl http://localhost:3000/api/versions
```

### Step 3: 访问 Swagger 文档
打开浏览器访问：
```
http://localhost:3000/swagger/index.html
```

### Step 4: 在控制器中使用新响应格式

**更新现有控制器**：
```go
// 旧代码
func (ctrl *UsersController) Show(c *gin.Context) {
    user := ...
    response.Data(c, user)  // 旧格式
}

// 新代码（推荐）
func (ctrl *UsersController) Show(c *gin.Context) {
    user := ...
    response.StandardSuccess(c, user)  // 新格式，包含更多元数据
}
```

---

## 5. 实际应用示例

### 示例1：用户列表（带分页）

```go
func (ctrl *UsersController) Index(c *gin.Context) {
    // 获取分页参数
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
    
    // 查询数据
    users, total, err := ctrl.userService.List(c, page, perPage)
    if err != nil {
        response.StandardError(c, http.StatusInternalServerError, 
            response.CodeInternalError, err.Error())
        return
    }
    
    // 计算总页数
    totalPages := int(math.Ceil(float64(total) / float64(perPage)))
    
    // 返回标准响应（带分页）
    response.StandardSuccessWithMeta(c, users, &response.MetaInfo{
        CurrentPage: page,
        PerPage:     perPage,
        Total:       total,
        TotalPages:  totalPages,
    })
}
```

### 示例2：创建用户（带验证）

```go
// @Summary 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body requests.UserRequest true "用户信息"
// @Success 201 {object} response.StandardResponse{data=user.User}
// @Failure 422 {object} response.StandardResponse
// @Router /users [post]
func (ctrl *UsersController) Store(c *gin.Context) {
    // 验证请求
    var req requests.UserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.StandardValidationError(c, map[string]string{
            "error": err.Error(),
        })
        return
    }
    
    // 创建用户
    user, err := ctrl.userService.Create(c, &req)
    if err != nil {
        response.StandardError(c, http.StatusInternalServerError,
            response.CodeInternalError, err.Error())
        return
    }
    
    // 返回创建结果
    c.JSON(http.StatusCreated, response.NewStandardResponse(
        c, true, response.CodeSuccess, "用户创建成功", user,
    ))
}
```

---

## 6. 迁移指南

### 控制器迁移步骤

1. **添加 Swagger 注解**
   ```go
   // @Summary ...
   // @Description ...
   // @Tags ...
   ```

2. **更新响应格式**
   ```go
   // 旧格式
   response.Data(c, data)
   
   // 新格式
   response.StandardSuccess(c, data)
   ```

3. **重新生成文档**
   ```bash
   make -f Makefile.swagger swagger
   ```

### 建议迁移顺序
1. ✅ 核心控制器（已完成示例：users_controller.go）
2. ⏳ 业务控制器（topics, comments, categories）
3. ⏳ 认证控制器（auth）
4. ⏳ 其他控制器

---

## 7. 常见问题

### Q: 旧的响应函数还能用吗？
A: 能。`response.Success`, `response.Data` 等函数**完全向后兼容**，不影响现有代码。

### Q: 必须迁移到新格式吗？
A: 不强制，但**推荐新代码使用新格式**，获得更丰富的响应元数据（timestamp, request_id, version等）。

### Q: Swagger 文档如何更新？
A: 修改注解后，运行 `make -f Makefile.swagger swagger` 重新生成。

### Q: 如何在 Swagger UI 中测试需要认证的接口？
A: 
1. 点击页面右上角的 "Authorize" 按钮
2. 输入：`Bearer {your_jwt_token}`
3. 点击 "Authorize"

### Q: 版本废弃如何通知客户端？
A: 服务器会在响应头中添加 `X-API-Warn` 警告信息。

---

## 8. 相关文档

| 文档 | 说明 |
|------|------|
| [docs/23_API_VERSIONING.md](./23_API_VERSIONING.md) | API 版本管理详细策略 |
| [docs/24_OPENAPI_GUIDE.md](./24_OPENAPI_GUIDE.md) | OpenAPI/Swagger 完整指南 |
| [Makefile.swagger](../Makefile.swagger) | Swagger 文档生成命令 |

---

## 9. 优化效果

### ✅ 已完成
- API 版本管理系统（版本信息、废弃警告、停用控制）
- 统一响应格式（Standard* 系列函数）
- OpenAPI/Swagger 文档支持
- 核心控制器注解示例
- 完整使用文档

### ⏳ 待完善
- 其他控制器 Swagger 注解补充（可逐步进行）
- 生产环境 Swagger 访问控制（可选）

### 📊 代码统计
- 新增文件：7 个
- 新增代码：1000+ 行
- 新增文档：500+ 行

---

## 📝 总结

v4.0 API 设计优化完成了：
1. ✅ API 版本管理策略
2. ✅ 响应格式统一（向后兼容）
3. ✅ OpenAPI 文档支持

**下一步**: 根据需要继续补充其他控制器的 Swagger 注解，或开始最后的「15. 安全加固」优化。
