# 🎛️ 管理后台 API 文档

完整的管理后台 API 接口文档。

## 🔐 权限说明

所有管理后台 API 都需要：
1. **JWT 认证** - 通过 `Authorization: Bearer <token>` 头部传递
2. **管理员角色** - 用户必须拥有 `admin` 角色

部分审核功能支持 `moderator` 角色（版主）。

---

## 📊 仪表盘 API

### 1. 获取系统概览

获取系统统计数据。

**接口**
```
GET /api/v1/admin/dashboard/overview
```

**响应**
```json
{
  "statistics": {
    "total_users": 1523,
    "total_topics": 8934,
    "total_categories": 12,
    "today_users": 45,
    "today_topics": 123,
    "active_users": 892,
    "popular_topics": 234
  },
  "timestamp": "2025-12-31T10:30:00Z"
}
```

### 2. 最近注册用户

获取最近 10 个注册用户。

**接口**
```
GET /api/v1/admin/dashboard/recent-users
```

**响应**
```json
{
  "users": [
    {
      "id": 123,
      "name": "张三",
      "email": "zhangsan@example.com",
      "created_at": "2025-12-31T10:00:00Z"
    }
  ]
}
```

### 3. 最近发布话题

获取最近 10 个发布的话题。

**接口**
```
GET /api/v1/admin/dashboard/recent-topics
```

**响应**
```json
{
  "topics": [
    {
      "id": 456,
      "title": "话题标题",
      "user": {
        "id": 123,
        "name": "张三"
      },
      "category": {
        "id": 1,
        "name": "技术讨论"
      },
      "created_at": "2025-12-31T09:50:00Z"
    }
  ]
}
```

---

## 👥 用户管理 API

### 1. 用户列表

获取用户列表，支持搜索和筛选。

**接口**
```
GET /api/v1/admin/users
```

**参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20 |
| keyword | string | 否 | 搜索关键词（用户名/邮箱） |
| status | string | 否 | 状态筛选 |

**示例**
```bash
GET /api/v1/admin/users?page=1&keyword=张三&status=1
```

**响应**
```json
{
  "users": [
    {
      "id": 123,
      "name": "张三",
      "email": "zhangsan@example.com",
      "phone": "13800138000",
      "status": 1,
      "created_at": "2025-12-20T10:00:00Z"
    }
  ],
  "paging": {
    "current_page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### 2. 用户详情

获取指定用户的详细信息和统计数据。

**接口**
```
GET /api/v1/admin/users/:id
```

**响应**
```json
{
  "user": {
    "id": 123,
    "name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "status": 1,
    "created_at": "2025-12-20T10:00:00Z"
  },
  "statistics": {
    "topic_count": 45,
    "comment_count": 128
  }
}
```

### 3. 更新用户信息

更新用户基本信息。

**接口**
```
PUT /api/v1/admin/users/:id
```

**请求体**
```json
{
  "name": "李四",
  "email": "lisi@example.com",
  "phone": "13900139000",
  "status": "1"
}
```

**响应**
```json
{
  "user": {
    "id": 123,
    "name": "李四",
    "email": "lisi@example.com"
  }
}
```

### 4. 删除用户

软删除用户（可恢复）。

**接口**
```
DELETE /api/v1/admin/users/:id
```

**响应**
```json
{
  "message": "用户已删除"
}
```

### 5. 批量删除用户

批量删除多个用户。

**接口**
```
POST /api/v1/admin/users/batch-delete
```

**请求体**
```json
{
  "ids": [123, 456, 789]
}
```

**响应**
```json
{
  "message": "批量删除成功",
  "count": 3
}
```

### 6. 封禁用户

封禁指定用户。

**接口**
```
POST /api/v1/admin/users/:id/ban
```

**请求体**
```json
{
  "reason": "违反社区规定",
  "days": 7
}
```

**响应**
```json
{
  "message": "用户已封禁",
  "reason": "违反社区规定",
  "days": 7
}
```

### 7. 解封用户

解除用户封禁。

**接口**
```
POST /api/v1/admin/users/:id/unban
```

**响应**
```json
{
  "message": "用户已解封"
}
```

### 8. 重置密码

管理员重置用户密码。

**接口**
```
POST /api/v1/admin/users/:id/reset-password
```

**请求体**
```json
{
  "password": "newpassword123",
  "password_confirmation": "newpassword123"
}
```

**响应**
```json
{
  "message": "密码重置成功"
}
```

### 9. 分配角色

为用户分配角色。

**接口**
```
POST /api/v1/admin/users/:id/assign-role
```

**请求体**
```json
{
  "role_ids": [1, 2]
}
```

**响应**
```json
{
  "message": "角色分配成功"
}
```

---

## 💬 话题管理 API

### 1. 话题列表

获取话题列表，支持搜索和筛选。

**接口**
```
GET /api/v1/admin/topics
```

**参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20 |
| keyword | string | 否 | 搜索关键词（标题/内容） |
| category_id | int | 否 | 分类 ID |
| user_id | int | 否 | 用户 ID |
| status | string | 否 | 状态筛选 |

**响应**
```json
{
  "topics": [
    {
      "id": 456,
      "title": "话题标题",
      "body": "话题内容...",
      "user": {
        "id": 123,
        "name": "张三"
      },
      "category": {
        "id": 1,
        "name": "技术讨论"
      },
      "order": 0,
      "created_at": "2025-12-30T10:00:00Z"
    }
  ],
  "paging": {
    "current_page": 1,
    "per_page": 20,
    "total": 500
  }
}
```

### 2. 话题详情

获取话题详细信息。

**接口**
```
GET /api/v1/admin/topics/:id
```

**响应**
```json
{
  "topic": {
    "id": 456,
    "title": "话题标题",
    "body": "话题内容...",
    "user": {
      "id": 123,
      "name": "张三"
    },
    "category": {
      "id": 1,
      "name": "技术讨论"
    }
  }
}
```

### 3. 更新话题

更新话题信息。

**接口**
```
PUT /api/v1/admin/topics/:id
```

**请求体**
```json
{
  "title": "新标题",
  "body": "新内容...",
  "category_id": 2,
  "status": 1
}
```

**响应**
```json
{
  "topic": {
    "id": 456,
    "title": "新标题"
  }
}
```

### 4. 删除话题

删除指定话题。

**接口**
```
DELETE /api/v1/admin/topics/:id
```

**响应**
```json
{
  "message": "话题已删除"
}
```

### 5. 批量删除话题

批量删除多个话题。

**接口**
```
POST /api/v1/admin/topics/batch-delete
```

**请求体**
```json
{
  "ids": [456, 789, 1011]
}
```

**响应**
```json
{
  "message": "批量删除成功",
  "count": 3
}
```

### 6. 置顶话题

将话题置顶显示。

**接口**
```
POST /api/v1/admin/topics/:id/pin
```

**响应**
```json
{
  "message": "话题已置顶"
}
```

### 7. 取消置顶

取消话题置顶。

**接口**
```
POST /api/v1/admin/topics/:id/unpin
```

**响应**
```json
{
  "message": "已取消置顶"
}
```

### 8. 审核通过

审核通过话题。

**接口**
```
POST /api/v1/admin/topics/:id/approve
```

**响应**
```json
{
  "message": "审核通过"
}
```

### 9. 审核拒绝

拒绝话题审核。

**接口**
```
POST /api/v1/admin/topics/:id/reject
```

**请求体**
```json
{
  "reason": "内容不符合规定"
}
```

**响应**
```json
{
  "message": "已拒绝",
  "reason": "内容不符合规定"
}
```

---

## 📂 分类管理 API

### 1. 分类列表

获取所有分类。

**接口**
```
GET /api/v1/admin/categories
```

**参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 搜索关键词 |

**响应**
```json
{
  "categories": [
    {
      "id": 1,
      "name": "技术讨论",
      "description": "技术相关话题",
      "order": 100
    }
  ]
}
```

### 2. 分类详情

获取分类详情和统计。

**接口**
```
GET /api/v1/admin/categories/:id
```

**响应**
```json
{
  "category": {
    "id": 1,
    "name": "技术讨论",
    "description": "技术相关话题",
    "order": 100
  },
  "topic_count": 234
}
```

### 3. 创建分类

创建新分类。

**接口**
```
POST /api/v1/admin/categories
```

**请求体**
```json
{
  "name": "产品讨论",
  "description": "产品相关话题",
  "order": 90
}
```

**响应**
```json
{
  "category": {
    "id": 5,
    "name": "产品讨论",
    "order": 90
  }
}
```

### 4. 更新分类

更新分类信息。

**接口**
```
PUT /api/v1/admin/categories/:id
```

**请求体**
```json
{
  "name": "技术交流",
  "description": "更新后的描述",
  "order": 95
}
```

**响应**
```json
{
  "category": {
    "id": 1,
    "name": "技术交流"
  }
}
```

### 5. 删除分类

删除分类（需确保分类下无话题）。

**接口**
```
DELETE /api/v1/admin/categories/:id
```

**响应**
```json
{
  "message": "分类已删除"
}
```

**错误响应**
```json
{
  "error": "该分类下还有话题，无法删除"
}
```

### 6. 分类排序

批量更新分类排序。

**接口**
```
POST /api/v1/admin/categories/sort
```

**请求体**
```json
{
  "categories": [
    {"id": 1, "order": 100},
    {"id": 2, "order": 90},
    {"id": 3, "order": 80}
  ]
}
```

**响应**
```json
{
  "message": "排序更新成功"
}
```

---

## 🛡️ 版主 API

版主拥有部分管理权限，主要用于内容审核。

### 权限说明

版主（`moderator` 角色）可以：
- ✅ 查看话题列表
- ✅ 审核话题（通过/拒绝）
- ✅ 删除话题
- ❌ 不能管理用户
- ❌ 不能管理分类

### 版主话题列表

**接口**
```
GET /api/v1/moderator/topics
```

### 版主审核通过

**接口**
```
POST /api/v1/moderator/topics/:id/approve
```

### 版主审核拒绝

**接口**
```
POST /api/v1/moderator/topics/:id/reject
```

### 版主删除话题

**接口**
```
DELETE /api/v1/moderator/topics/:id
```

---

## 🔒 错误响应

所有 API 可能返回的错误响应：

### 401 未认证
```json
{
  "error": "Unauthorized"
}
```

### 403 权限不足
```json
{
  "error": "Forbidden"
}
```

### 404 资源不存在
```json
{
  "error": "用户不存在"
}
```

### 400 参数错误
```json
{
  "error": "参数错误"
}
```

### 500 服务器错误
```json
{
  "error": "服务器错误"
}
```

---

## 📝 使用示例

### 完整请求示例

```bash
# 1. 登录获取 token
curl -X POST http://localhost:8080/api/v1/auth/login/using-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'

# 响应
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 7200
}

# 2. 使用 token 访问管理后台
curl -X GET http://localhost:8080/api/v1/admin/dashboard/overview \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."

# 3. 管理用户
curl -X GET http://localhost:8080/api/v1/admin/users?page=1&keyword=张三 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."

# 4. 封禁用户
curl -X POST http://localhost:8080/api/v1/admin/users/123/ban \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "违反社区规定",
    "days": 7
  }'
```

---

## 🎯 最佳实践

### 1. 安全建议

- ✅ 始终使用 HTTPS
- ✅ Token 应设置合理的过期时间
- ✅ 敏感操作（删除、封禁）需要二次确认
- ✅ 记录所有管理操作的审计日志

### 2. 性能建议

- ✅ 使用分页避免一次加载过多数据
- ✅ 合理使用搜索和筛选减少数据量
- ✅ 批量操作时注意并发控制

### 3. 使用建议

- ✅ 批量删除前先预览要删除的项目
- ✅ 封禁用户时说明原因
- ✅ 定期检查系统统计数据
- ✅ 合理设置分类排序权重

---

## 📞 技术支持

- 📖 完整文档: [docs/README.md](./README.md)
- 🐛 问题反馈: [GitHub Issues](https://github.com/Oyxiaoxi/GoHub-Service/issues)
- 💬 技术讨论: [GitHub Discussions](https://github.com/Oyxiaoxi/GoHub-Service/discussions)
