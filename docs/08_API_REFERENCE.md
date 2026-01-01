# 📖 GoHub-Service API 参考手册

**最后更新**: 2026年1月1日 | **版本**: v2.0

---

## 🎯 API 总览

本文档包含GoHub-Service所有REST API的完整参考。

| API类别 | 端点数 | 认证 | 文档 |
|---------|--------|------|------|
| **用户管理** | 15+ | ✅必需 | [用户API](#用户管理api) |
| **话题管理** | 12+ | ⚠️可选 | [话题API](#话题管理api) |
| **评论系统** | 10+ | ✅必需 | [评论API](#评论管理api) |
| **搜索功能** | 3 | ⚠️可选 | [搜索API](#搜索api) |
| **权限管理** | 8+ | 👑管理员 | [权限API](#权限系统api) |
| **系统接口** | 5 | ⚠️可选 | [系统API](#系统接口) |

---

## 🔐 认证方式

所有需要认证的API使用**JWT Bearer Token**：

```http
GET /api/users HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### 获取Token

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "user@example.com"
    }
  }
}
```

### Token有效期

- **有效期**: 24小时
- **刷新**: 调用 `POST /api/auth/refresh` 获取新token
- **撤销**: 调用 `POST /api/auth/logout` 进行登出

---

## 📋 通用响应格式

所有API响应采用统一格式：

### 成功响应（200）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 响应数据
  }
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": {
    "email": "邮箱格式不正确",
    "password": "密码长度不能少于6位"
  }
}
```

### 常见HTTP状态码

| 状态码 | 含义 | 示例 |
|--------|------|------|
| **200** | 请求成功 | 获取用户信息 |
| **201** | 资源创建 | 创建新话题 |
| **204** | 无内容 | 删除成功 |
| **400** | 请求错误 | 参数校验失败 |
| **401** | 未认证 | 缺少Token |
| **403** | 禁止访问 | 权限不足 |
| **404** | 资源不存在 | 用户不存在 |
| **429** | 限流 | 请求过于频繁 |
| **500** | 服务器错误 | 数据库异常 |

---

## 👥 用户管理API

### 1. 用户注册

**端点**: `POST /api/auth/register`  
**认证**: ❌ 不需要  
**限流**: 5次/分钟

**请求体**:
```json
{
  "name": "John Doe",
  "email": "user@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}
```

**响应** (201):
```json
{
  "code": 201,
  "message": "注册成功",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 用户登录

**端点**: `POST /api/auth/login`  
**认证**: ❌ 不需要  
**限流**: 10次/分钟

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "remember": true
}
```

**响应** (200):
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "user@example.com",
      "avatar": "https://...",
      "status": "active"
    }
  }
}
```

### 3. 获取当前用户信息

**端点**: `GET /api/me`  
**认证**: ✅ 需要  
**权限**: -

**响应** (200):
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "user@example.com",
    "avatar": "https://...",
    "bio": "个人简介",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 4. 更新用户资料

**端点**: `PUT /api/users/{id}`  
**认证**: ✅ 需要  
**权限**: 只能修改自己的信息

**请求体**:
```json
{
  "name": "Jane Doe",
  "bio": "新的个人简介",
  "avatar": "https://..."
}
```

**响应** (200):
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 1,
    "name": "Jane Doe",
    "bio": "新的个人简介",
    "avatar": "https://..."
  }
}
```

### 5. 修改密码

**端点**: `POST /api/users/{id}/change-password`  
**认证**: ✅ 需要  
**权限**: 只能修改自己的密码

**请求体**:
```json
{
  "current_password": "old_password",
  "new_password": "new_password",
  "new_password_confirmation": "new_password"
}
```

**响应** (200):
```json
{
  "code": 200,
  "message": "密码修改成功"
}
```

### 6. 用户列表

**端点**: `GET /api/users`  
**认证**: ⚠️ 可选  
**权限**: -  
**分页**: 支持

**查询参数**:
```
?page=1&limit=20&sort=-created_at&keyword=john
```

**响应** (200):
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "John Doe",
        "email": "user@example.com",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "last_page": 5
    }
  }
}
```

### 7. 获取单个用户

**端点**: `GET /api/users/{id}`  
**认证**: ⚠️ 可选  
**权限**: -

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "user@example.com",
    "avatar": "https://...",
    "bio": "个人简介",
    "created_at": "2024-01-01T00:00:00Z",
    "stats": {
      "topics": 42,
      "comments": 128,
      "followers": 50,
      "following": 30
    }
  }
}
```

### 8. 用户关注

**端点**: `POST /api/users/{id}/follow`  
**认证**: ✅ 需要  
**权限**: -

**响应** (200):
```json
{
  "code": 200,
  "message": "关注成功",
  "data": {
    "following": true
  }
}
```

### 9. 用户取消关注

**端点**: `DELETE /api/users/{id}/follow`  
**认证**: ✅ 需要  
**权限**: -

**响应** (204):
```json
{
  "code": 204,
  "message": "取消关注成功"
}
```

### 10. 获取用户粉丝列表

**端点**: `GET /api/users/{id}/followers`  
**认证**: ⚠️ 可选  
**权限**: -

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": 2,
        "name": "Jane Smith",
        "avatar": "https://..."
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 50
    }
  }
}
```

---

## 📝 话题管理API

### 1. 创建话题

**端点**: `POST /api/topics`  
**认证**: ✅ 需要  
**权限**: -

**请求体**:
```json
{
  "title": "如何学习Go语言？",
  "body": "这是话题内容...",
  "category_id": 1,
  "tags": ["golang", "learning"]
}
```

**响应** (201):
```json
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "id": 1,
    "title": "如何学习Go语言？",
    "slug": "how-to-learn-golang",
    "body": "这是话题内容...",
    "user_id": 1,
    "category_id": 1,
    "view_count": 0,
    "like_count": 0,
    "comment_count": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 获取话题列表

**端点**: `GET /api/topics`  
**认证**: ⚠️ 可选  
**权限**: -  
**分页**: 支持

**查询参数**:
```
?page=1&limit=20&sort=-created_at&category_id=1
```

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": 1,
        "title": "如何学习Go语言？",
        "slug": "how-to-learn-golang",
        "excerpt": "这是话题内容摘要...",
        "user": {
          "id": 1,
          "name": "John Doe"
        },
        "view_count": 150,
        "like_count": 25,
        "comment_count": 8,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 500,
      "last_page": 25
    }
  }
}
```

### 3. 获取单个话题

**端点**: `GET /api/topics/{id}`  
**认证**: ⚠️ 可选  
**权限**: -

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "title": "如何学习Go语言？",
    "body": "完整的话题内容...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "avatar": "https://..."
    },
    "category": {
      "id": 1,
      "name": "编程语言"
    },
    "view_count": 150,
    "like_count": 25,
    "comment_count": 8,
    "is_liked": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  }
}
```

### 4. 更新话题

**端点**: `PUT /api/topics/{id}`  
**认证**: ✅ 需要  
**权限**: 只能修改自己的话题

**请求体**:
```json
{
  "title": "更新的标题",
  "body": "更新的内容",
  "category_id": 2
}
```

**响应** (200):
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 1,
    "title": "更新的标题",
    "body": "更新的内容",
    "category_id": 2,
    "updated_at": "2024-01-02T10:00:00Z"
  }
}
```

### 5. 删除话题

**端点**: `DELETE /api/topics/{id}`  
**认证**: ✅ 需要  
**权限**: 只能删除自己的话题或管理员

**响应** (204):
```json
{
  "code": 204,
  "message": "删除成功"
}
```

### 6. 点赞话题

**端点**: `POST /api/topics/{id}/like`  
**认证**: ✅ 需要  
**权限**: -

**响应** (200):
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": {
    "liked": true,
    "like_count": 26
  }
}
```

### 7. 取消点赞

**端点**: `DELETE /api/topics/{id}/like`  
**认证**: ✅ 需要  
**权限**: -

**响应** (200):
```json
{
  "code": 200,
  "message": "取消点赞成功",
  "data": {
    "liked": false,
    "like_count": 25
  }
}
```

### 8. 收藏话题

**端点**: `POST /api/topics/{id}/favorite`  
**认证**: ✅ 需要  
**权限**: -

**响应** (200):
```json
{
  "code": 200,
  "message": "收藏成功",
  "data": {
    "favorited": true
  }
}
```

### 9. 热门话题

**端点**: `GET /api/topics/hot`  
**认证**: ⚠️ 可选  
**权限**: -

**查询参数**:
```
?limit=10&period=week
```

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": 1,
        "title": "如何学习Go语言？",
        "view_count": 5000,
        "like_count": 200,
        "comment_count": 80,
        "score": 95.5
      }
    ]
  }
}
```

---

## 💬 评论管理API

### 1. 创建评论

**端点**: `POST /api/topics/{topic_id}/comments`  
**认证**: ✅ 需要  
**权限**: -

**请求体**:
```json
{
  "body": "这是一条评论",
  "parent_id": null
}
```

**响应** (201):
```json
{
  "code": 201,
  "message": "评论成功",
  "data": {
    "id": 1,
    "body": "这是一条评论",
    "user_id": 1,
    "topic_id": 1,
    "parent_id": null,
    "like_count": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 获取话题的评论列表

**端点**: `GET /api/topics/{topic_id}/comments`  
**认证**: ⚠️ 可选  
**权限**: -  
**分页**: 支持

**查询参数**:
```
?page=1&limit=10&sort=-created_at
```

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": 1,
        "body": "这是一条评论",
        "user": {
          "id": 1,
          "name": "John Doe",
          "avatar": "https://..."
        },
        "like_count": 5,
        "children": [
          {
            "id": 2,
            "body": "这是一条回复",
            "user": {
              "id": 2,
              "name": "Jane Smith"
            }
          }
        ],
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50
    }
  }
}
```

### 3. 更新评论

**端点**: `PUT /api/comments/{id}`  
**认证**: ✅ 需要  
**权限**: 只能修改自己的评论

**请求体**:
```json
{
  "body": "更新的评论内容"
}
```

**响应** (200):
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 1,
    "body": "更新的评论内容",
    "updated_at": "2024-01-02T00:00:00Z"
  }
}
```

### 4. 删除评论

**端点**: `DELETE /api/comments/{id}`  
**认证**: ✅ 需要  
**权限**: 只能删除自己的评论或管理员

**响应** (204):
```json
{
  "code": 204,
  "message": "删除成功"
}
```

### 5. 点赞评论

**端点**: `POST /api/comments/{id}/like`  
**认证**: ✅ 需要  
**权限**: -

**响应** (200):
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": {
    "liked": true,
    "like_count": 6
  }
}
```

---

## 🔍 搜索API

### 1. 搜索话题

**端点**: `GET /api/search/topics`  
**认证**: ⚠️ 可选  
**权限**: -  
**特性**: 支持Elasticsearch全文搜索，自动降级到数据库查询

**查询参数**:
```
?q=golang&page=1&limit=20&category_id=1&sort=-view_count
```

**响应** (200):
```json
{
  "code": 200,
  "message": "搜索成功",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "如何学习Go语言？",
        "excerpt": "这是话题内容摘要...",
        "user": {
          "id": 1,
          "name": "John Doe"
        },
        "category": {
          "id": 1,
          "name": "编程语言"
        },
        "view_count": 150,
        "relevance_score": 0.95,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 128,
      "last_page": 7
    },
    "search_time_ms": 15
  }
}
```

### 2. 获取搜索建议

**端点**: `GET /api/search/suggestions`  
**认证**: ⚠️ 可选  
**权限**: -

**查询参数**:
```
?q=gol&limit=5
```

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "suggestions": [
      {
        "text": "golang",
        "frequency": 125
      },
      {
        "text": "golf",
        "frequency": 48
      }
    ]
  }
}
```

### 3. 获取热点话题

**端点**: `GET /api/search/hot-topics`  
**认证**: ⚠️ 可选  
**权限**: -

**查询参数**:
```
?limit=10&period=week
```

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": 1,
        "title": "2024 Go语言发展趋势",
        "view_count": 5000,
        "like_count": 200,
        "comment_count": 80,
        "trending_score": 95.5
      }
    ]
  }
}
```

---

## 👑 权限系统API

### 1. 获取所有角色

**端点**: `GET /api/admin/roles`  
**认证**: ✅ 需要  
**权限**: 👑 管理员  
**分页**: 支持

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": 1,
        "name": "admin",
        "display_name": "管理员",
        "description": "拥有系统所有权限",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "total": 5
    }
  }
}
```

### 2. 创建角色

**端点**: `POST /api/admin/roles`  
**认证**: ✅ 需要  
**权限**: 👑 管理员

**请求体**:
```json
{
  "name": "moderator",
  "display_name": "版主",
  "description": "负责内容审核"
}
```

**响应** (201):
```json
{
  "code": 201,
  "data": {
    "id": 2,
    "name": "moderator",
    "display_name": "版主",
    "description": "负责内容审核"
  }
}
```

### 3. 分配权限给角色

**端点**: `POST /api/admin/roles/{id}/permissions`  
**认证**: ✅ 需要  
**权限**: 👑 管理员

**请求体**:
```json
{
  "permission_ids": [1, 2, 3]
}
```

**响应** (200):
```json
{
  "code": 200,
  "message": "权限分配成功"
}
```

### 4. 获取角色权限

**端点**: `GET /api/admin/roles/{id}/permissions`  
**认证**: ✅ 需要  
**权限**: 👑 管理员

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "permissions": [
      {
        "id": 1,
        "name": "topics.create",
        "display_name": "创建话题",
        "description": "允许用户创建新话题"
      }
    ]
  }
}
```

### 5. 为用户分配角色

**端点**: `POST /api/admin/users/{id}/roles`  
**认证**: ✅ 需要  
**权限**: 👑 管理员

**请求体**:
```json
{
  "role_ids": [1, 2]
}
```

**响应** (200):
```json
{
  "code": 200,
  "message": "角色分配成功"
}
```

### 6. 获取所有权限

**端点**: `GET /api/admin/permissions`  
**认证**: ✅ 需要  
**权限**: 👑 管理员

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": 1,
        "name": "topics.create",
        "display_name": "创建话题",
        "description": "允许用户创建新话题"
      }
    ]
  }
}
```

---

## 🔧 系统接口

### 1. 健康检查

**端点**: `GET /api/health`  
**认证**: ❌ 不需要  
**用途**: 监控和负载均衡器探针

**响应** (200):
```json
{
  "code": 200,
  "message": "服务正常",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T00:00:00Z",
    "services": {
      "database": "ok",
      "redis": "ok",
      "elasticsearch": "ok"
    }
  }
}
```

### 2. 系统信息

**端点**: `GET /api/system/info`  
**认证**: ❌ 不需要  
**用途**: 获取系统基本信息

**响应** (200):
```json
{
  "code": 200,
  "data": {
    "version": "2.0.0",
    "environment": "production",
    "go_version": "1.21",
    "uptime": "720h30m",
    "database_version": "8.0.35"
  }
}
```

---

## 📊 分页说明

所有分页接口使用统一参数：

```http
GET /api/users?page=1&limit=20&sort=-created_at
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| **page** | 1 | 页码，从1开始 |
| **limit** | 20 | 每页数量，最大100 |
| **sort** | -created_at | 排序字段，`-`表示倒序 |

**分页响应格式**:
```json
{
  "items": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "last_page": 5
  }
}
```

---

## 🚫 错误代码大全

| 错误码 | HTTP状态 | 含义 | 处理方案 |
|--------|---------|------|---------|
| 400 | 400 | 请求参数错误 | 检查请求参数格式 |
| 401 | 401 | 缺少认证Token | 调用登录接口获取Token |
| 403 | 403 | 权限不足 | 检查用户权限 |
| 404 | 404 | 资源不存在 | 检查资源ID是否正确 |
| 422 | 422 | 业务逻辑错误 | 查看错误详情修正 |
| 429 | 429 | 请求过于频繁 | 等待后重试 |
| 500 | 500 | 服务器错误 | 联系技术支持 |

---

## ⚙️ 请求示例集合

### cURL示例

```bash
# 用户登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# 获取话题列表
curl -X GET http://localhost:8080/api/topics?page=1&limit=20 \
  -H "Authorization: Bearer YOUR_TOKEN"

# 创建话题
curl -X POST http://localhost:8080/api/topics \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "新话题",
    "body": "话题内容",
    "category_id": 1
  }'

# 搜索话题
curl -X GET "http://localhost:8080/api/search/topics?q=golang" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### JavaScript示例

```javascript
// 登录
const login = async (email, password) => {
  const response = await fetch('http://localhost:8080/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  return response.json();
};

// 获取话题
const getTopics = async (token, page = 1) => {
  const response = await fetch(
    `http://localhost:8080/api/topics?page=${page}&limit=20`,
    {
      headers: { 'Authorization': `Bearer ${token}` }
    }
  );
  return response.json();
};

// 搜索话题
const searchTopics = async (token, query) => {
  const response = await fetch(
    `http://localhost:8080/api/search/topics?q=${encodeURIComponent(query)}`,
    {
      headers: { 'Authorization': `Bearer ${token}` }
    }
  );
  return response.json();
};
```

---

## 📚 相关文档

- [权限系统详解](03_RBAC.md) - RBAC权限设计和实现
- [开发指南](05_DEVELOPMENT.md) - API开发规范
- [性能指标](07_PERFORMANCE.md) - API性能基准
- [搜索功能](10_ELASTICSEARCH.md) - 全文搜索集成

---
