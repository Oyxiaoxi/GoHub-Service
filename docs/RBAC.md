# 🔐 RBAC 权限系统

基于角色的访问控制（Role-Based Access Control）实现指南。

## 概述

RBAC 是一个三层权限模型：

```
User → Role → Permission
 ↓      ↓        ↓
用户  角色  权限
```

## 数据模型

### 1. Role（角色）
```sql
roles
├── id (PK)
├── name (唯一)        -- admin, moderator, user 等
├── display_name       -- "管理员", "版主", "用户"
├── description
├── created_at
└── updated_at
```

### 2. Permission（权限）
```sql
permissions
├── id (PK)
├── name (唯一)        -- topics.create, users.delete 等
├── display_name       -- "创建话题", "删除用户"
├── description
├── created_at
└── updated_at
```

### 3. RolePermission（角色-权限关联）
```sql
role_permissions
├── id (PK)
├── role_id (FK)       -- 所属角色
├── permission_id (FK) -- 所属权限
├── created_at
└── updated_at
```

### 4. UserRole（用户-角色关联）
```sql
user_roles
├── id (PK)
├── user_id (FK)       -- 所属用户
├── role_id (FK)       -- 所属角色
├── created_at
└── updated_at
```

## 中间件使用

### 1. 基本使用

在路由中添加权限检查中间件：

```go
// 需要 topics.view 权限
r.GET("/topics", middlewares.RequirePermission("topics.view"), controllers.TopicIndex)

// 需要 topics.create 权限
r.POST("/topics", middlewares.RequirePermission("topics.create"), controllers.TopicStore)

// 需要 topics.delete 权限
r.DELETE("/topics/:id", middlewares.RequirePermission("topics.delete"), controllers.TopicDestroy)
```

### 2. 中间件实现

文件位置: `app/http/middlewares/rbac.go`

```go
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 验证认证（从 JWT 获取用户 ID）
        userID, ok := c.Get("userID")
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // 2. 检查用户是否有权限
        if !service.UserService.HasPermission(userID, permission) {
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }

        // 3. 继续执行
        c.Next()
    }
}
```

## Service 方法

### 检查权限
```go
// 检查用户是否有特定权限
has := userService.HasPermission(userID, "topics.create")

// 检查用户是否有某个角色
has := userService.HasRole(userID, "admin")

// 获取用户的所有权限
permissions := userService.GetPermissions(userID)

// 获取用户的所有角色
roles := userService.GetRoles(userID)
```

## 常用权限列表

### 话题权限
| 权限名 | 说明 |
|-------|------|
| topics.view | 查看话题 |
| topics.create | 创建话题 |
| topics.update | 编辑话题 |
| topics.delete | 删除话题 |
| topics.pin | 置顶话题 |
| topics.restore | 恢复话题 |

### 评论权限
| 权限名 | 说明 |
|-------|------|
| comments.view | 查看评论 |
| comments.create | 创建评论 |
| comments.update | 编辑评论 |
| comments.delete | 删除评论 |

### 用户权限
| 权限名 | 说明 |
|-------|------|
| users.view | 查看用户信息 |
| users.update | 编辑用户信息 |
| users.delete | 删除用户 |
| users.ban | 封禁用户 |

### 管理权限
| 权限名 | 说明 |
|-------|------|
| admin.settings | 修改系统设置 |
| admin.users | 管理用户 |
| admin.roles | 管理角色 |
| admin.permissions | 管理权限 |

## 常用角色配置

### 超级管理员 (admin)
拥有所有权限

### 版主 (moderator)
- topics.view, topics.create, topics.update, topics.delete
- comments.view, comments.delete
- users.view, users.ban

### 普通用户 (user)
- topics.view, topics.create, topics.update
- comments.view, comments.create, comments.update
- users.view

### 访客 (guest)
- topics.view
- comments.view

## 数据库初始化

### 创建权限
```go
// database/seeders/permissions_seeder.go
permission := &models.Permission{
    Name: "topics.create",
    DisplayName: "创建话题",
    Description: "允许创建新话题",
}
db.Create(permission)
```

### 创建角色
```go
// database/seeders/roles_seeder.go
role := &models.Role{
    Name: "user",
    DisplayName: "普通用户",
    Description: "默认用户角色",
}
db.Create(role)
```

### 分配权限给角色
```go
// 为 user 角色分配权限
rolePermission := &models.RolePermission{
    RoleID: userRole.ID,
    PermissionID: createTopicPerm.ID,
}
db.Create(rolePermission)
```

### 分配角色给用户
```go
// 给用户分配 user 角色
userRole := &models.UserRole{
    UserID: user.ID,
    RoleID: role.ID,
}
db.Create(userRole)
```

## 最佳实践

### 1. 权限粒度
- ✅ 细粒度权限: `topics.create`, `topics.update`, `topics.delete`
- ❌ 粗粒度权限: `can_manage_topics`

### 2. 权限命名
- ✅ 资源.动作: `users.delete`, `posts.edit`, `comments.create`
- ❌ 动词在前: `deleteUsers`, `editPosts`

### 3. 缓存权限
```go
// 缓存用户权限，避免重复查询
permissions := cache.Get("user_permissions:" + userID)
if permissions == nil {
    permissions = db.GetUserPermissions(userID)
    cache.Set("user_permissions:" + userID, permissions, 1*time.Hour)
}
```

### 4. 权限失效
当修改用户角色或权限时，清除缓存：

```go
// 分配角色时
userService.AssignRole(userID, roleID)
cache.Delete("user_permissions:" + userID)
cache.Delete("user_roles:" + userID)
```

### 5. 默认拒绝原则
- ✅ 只有明确授予的权限才能执行
- ❌ 未授予的权限仍然允许执行

## 常见问题

### Q: 如何给特定用户分配权限？
A: 通过角色分配，而不是直接分配权限。用户 → 角色 → 权限

### Q: 如何快速检查权限？
A: 使用中间件在路由层检查，不要在 Controller 中检查

### Q: 权限太多怎么办？
A: 使用权限分类或权限组，避免权限爆炸

### Q: 如何实现动态权限？
A: 存储权限在数据库，不硬编码，启动时加载到内存

---

更多信息请查看 [ARCHITECTURE.md](./ARCHITECTURE.md) 和 [DEVELOPMENT.md](./DEVELOPMENT.md)
