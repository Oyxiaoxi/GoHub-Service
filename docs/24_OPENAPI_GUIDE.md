# OpenAPI 文档使用指南

## 概述

GoHub-Service 使用 Swagger/OpenAPI 3.0 标准生成 API 文档，提供交互式 API 测试界面。

## 访问文档

### 本地开发环境
```
http://localhost:3000/swagger/index.html
```

### 生产环境
```
https://api.gohub.com/swagger/index.html
```

## 安装依赖

```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 安装 gin-swagger
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

## 生成文档

### 命令行生成
```bash
# 在项目根目录执行
swag init

# 指定输出目录
swag init -o docs

# 指定 API 目录
swag init -d ./,./app/http/controllers
```

### Makefile 集成
```bash
# 添加到 Makefile
swagger:
	swag init --parseDependency --parseInternal

# 使用
make swagger
```

## 注解语法

### 主应用注解（main.go）

```go
// @title GoHub-Service API
// @version 1.0
// @description GoHub 社区论坛 API 文档
// @termsOfService https://github.com/Oyxiaoxi/GoHub-Service

// @contact.name API Support
// @contact.url https://github.com/Oyxiaoxi/GoHub-Service/issues
// @contact.email support@gohub.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main
```

### 控制器方法注解

#### 基础注解
```go
// CurrentUser 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 返回当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.StandardResponse{data=user.User}
// @Failure 401 {object} response.StandardResponse
// @Router /user [get]
func (ctrl *UsersController) CurrentUser(c *gin.Context) {
    // 实现代码
}
```

#### 带参数的注解
```go
// Show 获取用户详情
// @Summary 获取指定用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} response.StandardResponse{data=user.User}
// @Failure 404 {object} response.StandardResponse
// @Router /users/{id} [get]
func (ctrl *UsersController) Show(c *gin.Context) {
    // 实现代码
}
```

#### 带请求体的注解
```go
// Store 创建用户
// @Summary 创建新用户
// @Description 创建一个新的用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body requests.UserRequest true "用户信息"
// @Success 201 {object} response.StandardResponse{data=user.User}
// @Failure 422 {object} response.StandardResponse
// @Router /users [post]
func (ctrl *UsersController) Store(c *gin.Context) {
    // 实现代码
}
```

#### 带查询参数的注解
```go
// Index 用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(20)
// @Param order query string false "排序方式" default(id)
// @Success 200 {object} response.StandardResponse{data=[]user.User,meta=response.MetaInfo}
// @Router /users [get]
func (ctrl *UsersController) Index(c *gin.Context) {
    // 实现代码
}
```

## 注解参数说明

### @Param 参数类型

| 参数位置 | 说明 | 示例 |
|---------|------|------|
| path | URL 路径参数 | `/users/{id}` |
| query | URL 查询参数 | `/users?page=1` |
| header | HTTP 头部 | `X-Request-ID` |
| body | 请求体 | JSON/XML 数据 |
| formData | 表单数据 | 文件上传 |

### @Success/@Failure 格式

```
@Success {httpCode} {responseType} {dataType} {description}
@Failure {httpCode} {responseType} {dataType} {description}
```

示例：
```go
// @Success 200 {object} response.StandardResponse{data=user.User}
// @Failure 404 {object} response.StandardResponse{error=response.ErrorInfo}
```

### 常用注解标签

| 标签 | 说明 | 示例 |
|------|------|------|
| @Summary | 简短描述 | 获取用户信息 |
| @Description | 详细描述 | 根据用户ID获取详细信息 |
| @Tags | 分组标签 | 用户管理 |
| @Accept | 接受的MIME类型 | json, xml, multipart/form-data |
| @Produce | 返回的MIME类型 | json, xml |
| @Security | 安全认证 | Bearer |
| @Deprecated | 标记废弃 | - |

## 数据模型注解

### 结构体注解
```go
// User 用户模型
type User struct {
    // 用户ID
    ID uint64 `json:"id" example:"1"`
    
    // 用户名
    // required: true
    // minLength: 3
    // maxLength: 20
    Name string `json:"name" example:"张三" validate:"required,min=3,max=20"`
    
    // 邮箱
    // format: email
    Email string `json:"email" example:"zhangsan@example.com" validate:"email"`
    
    // 手机号
    // pattern: ^1[3-9]\d{9}$
    Phone string `json:"phone" example:"13800138000"`
    
    // 创建时间
    // format: date-time
    CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
} // @name User
```

### 枚举类型
```go
// UserStatus 用户状态
type UserStatus string

const (
    // 活跃
    StatusActive UserStatus = "active"
    // 禁用
    StatusDisabled UserStatus = "disabled"
    // 删除
    StatusDeleted UserStatus = "deleted"
) // @name UserStatus
```

## 响应格式示例

### 成功响应
```json
{
  "success": true,
  "code": 20000,
  "message": "success",
  "data": {
    "id": 1,
    "name": "张三",
    "email": "zhangsan@example.com"
  },
  "timestamp": 1704067200,
  "request_id": "abc123",
  "version": "v1"
}
```

### 分页响应
```json
{
  "success": true,
  "code": 20000,
  "message": "success",
  "data": [...],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  },
  "timestamp": 1704067200,
  "request_id": "abc123"
}
```

### 错误响应
```json
{
  "success": false,
  "code": 40400,
  "message": "Resource not found",
  "error": {
    "type": "not_found",
    "details": "User with ID 999 not found"
  },
  "timestamp": 1704067200,
  "request_id": "abc123"
}
```

## 认证配置

### JWT Token 认证
```go
// 在 Swagger UI 中点击 "Authorize" 按钮
// 输入格式：Bearer {your_token}
// 示例：Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 自定义配置

### swag.yml 配置文件
```yaml
# 项目根目录创建 swag.yml
version: "1.0"
swagger: "2.0"
info:
  title: "GoHub-Service API"
  version: "1.0"
host: "localhost:3000"
basePath: "/api/v1"
schemes:
  - "http"
  - "https"
```

## 最佳实践

### 1. 文档即代码
- 在编写代码时同步更新注解
- 使用 CI/CD 自动生成文档

### 2. 完整性检查
```bash
# 检查未文档化的端点
swag init --parseVendor --parseDependency
```

### 3. 版本控制
- 为不同版本生成独立文档
- 保留历史版本文档

### 4. 测试覆盖
- 使用 Swagger UI 测试所有端点
- 验证请求/响应格式

### 5. 安全注意事项
- 生产环境可选择关闭文档访问
- 使用认证保护文档端点
- 不在文档中暴露敏感信息

## 常见问题

### Q: 如何隐藏某些端点？
A: 不添加 Swagger 注解即可

### Q: 如何自定义响应示例？
A: 使用 `example` 标签
```go
// @Success 200 {object} User "success" example({"id":1,"name":"张三"})
```

### Q: 如何支持文件上传？
A: 使用 `formData` 参数类型
```go
// @Param file formData file true "上传文件"
// @Accept multipart/form-data
```

### Q: 如何标记已废弃的接口？
A: 使用 `@Deprecated` 注解
```go
// @Deprecated
// @Summary 旧版用户创建接口（已废弃）
```

## 相关资源

- [Swagger 官方文档](https://swagger.io/docs/)
- [swag GitHub](https://github.com/swaggo/swag)
- [OpenAPI 3.0 规范](https://spec.openapis.org/oas/v3.0.0)
- [gin-swagger](https://github.com/swaggo/gin-swagger)
