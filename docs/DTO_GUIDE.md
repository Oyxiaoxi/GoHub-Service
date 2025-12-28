# DTO 层设计指南

## 1. 概述

DTO（Data Transfer Object）数据传输对象是用于在不同层之间传递数据的对象。本指南定义了GoHub-Service项目的DTO设计标准和最佳实践。

## 2. DTO设计原则

### 2.1 职责分离
- **请求DTO**：用于接收客户端请求数据，包含验证规则
- **响应DTO**：用于返回给客户端的数据，隐藏敏感字段
- **内部DTO**：服务间传递数据使用

### 2.2 命名规范
```
{Entity}{Operation}DTO

示例：
- UserCreateDTO      // 创建用户请求DTO
- UserUpdateDTO      // 更新用户请求DTO
- UserResponseDTO    // 用户响应DTO
- TopicListResponseDTO  // 话题列表响应DTO
```

### 2.3 字段设计
- 使用验证标签：`binding:"required"`, `binding:"min=3,max=255"`
- 使用JSON标签：`json:"field_name"`
- 响应DTO只包含必要字段，过滤敏感信息

## 3. DTO分类

### 3.1 请求DTO（Request DTO）

用于接收和验证客户端请求数据。

**特点：**
- 包含验证规则（binding标签）
- 字段与API请求参数对应
- 不包含业务逻辑

**示例：**
```go
type UserCreateDTO struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Phone    string `json:"phone" binding:"required"`
}
```

### 3.2 响应DTO（Response DTO）

用于返回数据给客户端。

**特点：**
- 过滤敏感字段（如密码、token等）
- 可以组合多个模型数据
- 格式化输出字段

**示例：**
```go
type UserResponseDTO struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    // 注意：不包含Password字段
}
```

### 3.3 列表响应DTO

用于返回列表数据。

**示例：**
```go
type UserListResponseDTO struct {
    Users []UserResponseDTO  `json:"users"`
    Paging *paginator.Paging `json:"paging"`
}
```

## 4. DTO转换

### 4.1 Model到ResponseDTO

在Service层进行转换：

```go
func (s *UserService) GetByID(id string) (*UserResponseDTO, error) {
    user, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    return s.toResponseDTO(user), nil
}

func (s *UserService) toResponseDTO(user *user.User) *UserResponseDTO {
    return &UserResponseDTO{
        ID:        user.GetStringID(),
        Name:      user.Name,
        Email:     user.Email,
        Phone:     user.Phone,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
}
```

### 4.2 RequestDTO到Model

在Service层进行转换：

```go
func (s *UserService) Create(dto UserCreateDTO) (*UserResponseDTO, error) {
    userModel := &user.User{
        Name:     dto.Name,
        Email:    dto.Email,
        Password: dto.Password, // 注意：实际应用中需要加密
        Phone:    dto.Phone,
    }
    
    if err := s.repo.Create(userModel); err != nil {
        return nil, err
    }
    
    return s.toResponseDTO(userModel), nil
}
```

## 5. 现有DTO实现

### 5.1 Topic Service
- ✅ TopicCreateDTO - 创建话题请求DTO
- ✅ TopicUpdateDTO - 更新话题请求DTO
- ✅ TopicResponseDTO - 话题响应DTO
- ✅ TopicListResponseDTO - 话题列表响应DTO

### 5.2 Category Service
- ✅ CategoryCreateDTO - 创建分类请求DTO
- ✅ CategoryUpdateDTO - 更新分类请求DTO
- ✅ CategoryResponseDTO - 分类响应DTO
- ✅ CategoryListResponseDTO - 分类列表响应DTO

### 5.3 User Service
- ✅ UserCreateDTO - 创建用户请求DTO
- ✅ UserUpdateDTO - 更新用户请求DTO
- ✅ UserResponseDTO - 用户响应DTO
- ✅ UserListResponseDTO - 用户列表响应DTO

### 5.4 Link Service
- ✅ LinkResponseDTO - 友情链接响应DTO
- ✅ LinkListResponseDTO - 友情链接列表响应DTO

## 6. 最佳实践

### 6.1 验证规则
```go
// 使用Gin的binding标签
type UserCreateDTO struct {
    Name     string `binding:"required,min=3,max=255"`
    Email    string `binding:"required,email"`
    Password string `binding:"required,min=6"`
    Phone    string `binding:"required,phone"` // 自定义验证器
}
```

### 6.2 敏感字段处理
```go
// 响应DTO中永远不要包含敏感字段
type UserResponseDTO struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // 不包含 Password 字段
}
```

### 6.3 字段组合
```go
// 可以组合多个模型的数据
type TopicDetailResponseDTO struct {
    TopicResponseDTO
    Category CategoryResponseDTO `json:"category"`
    User     UserResponseDTO     `json:"user"`
}
```

### 6.4 可选字段
```go
// 更新DTO中的字段通常是可选的
type UserUpdateDTO struct {
    Name  *string `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
    Email *string `json:"email,omitempty" binding:"omitempty,email"`
    Phone *string `json:"phone,omitempty" binding:"omitempty"`
}
```

## 7. Controller中使用DTO

```go
func (ctrl *UserController) Create(c *gin.Context) {
    // 1. 绑定请求DTO
    var dto services.UserCreateDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        response.ValidationError(c, err)
        return
    }
    
    // 2. 调用Service，传入DTO
    userResponse, err := ctrl.userService.Create(dto)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    // 3. 返回响应DTO
    response.Created(c, userResponse)
}
```

## 8. 测试DTO

```go
func TestUserCreateDTO(t *testing.T) {
    dto := services.UserCreateDTO{
        Name:     "testuser",
        Email:    "test@example.com",
        Password: "password123",
        Phone:    "13800138000",
    }
    
    assert.NotEmpty(t, dto.Name)
    assert.NotEmpty(t, dto.Email)
    assert.NotEmpty(t, dto.Password)
}
```

## 9. 注意事项

1. **不要在DTO中包含业务逻辑**
   - DTO只是数据容器
   - 业务逻辑应该在Service层

2. **避免DTO与Model直接耦合**
   - 使用转换函数进行转换
   - 便于后续修改和维护

3. **响应DTO必须过滤敏感信息**
   - 密码、Token、密钥等
   - 内部状态字段

4. **合理使用验证标签**
   - 在DTO层进行基础验证
   - 复杂业务验证在Service层

5. **保持DTO简单**
   - 一个DTO只负责一个场景
   - 避免过度设计

## 10. 相关文档

- [Service层架构指南](SERVICE_LAYER_GUIDE.md)
- [Controller复用指南](../CONTROLLER_REUSE_GUIDE.md)
- [代码规范](../CODING_STANDARDS.md)
