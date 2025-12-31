# 角色和权限管理 API 文档

## 概述

本文档定义了 GoHub 服务的角色和权限管理 API 接口。所有接口均位于 `/api/v1/admin` 路径下，需要管理员权限认证。

## 认证说明

- 所有接口都需要有效的 JWT token
- 所有接口都需要用户拥有 `admin` 角色
- 在请求头中包含: `Authorization: Bearer {token}`

## 响应格式

所有接口遵循统一的 JSON 响应格式：

### 成功响应 (200/201)

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    // 返回数据
  }
}
```

### 错误响应

```json
{
  "code": 1001,
  "message": "错误描述",
  "data": null
}
```

## 角色管理 API

### 1. 获取角色列表

**请求**
```
GET /api/v1/admin/roles
```

**查询参数**
| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认为 1 |
| per_page | int | 否 | 每页数量，默认为 20，最大为 100 |
| keyword | string | 否 | 搜索关键词（名称或显示名称） |

**响应示例**
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "admin",
        "display_name": "管理员",
        "description": "拥有所有权限",
        "created_at": "2025-01-01T10:00:00Z",
        "updated_at": "2025-01-01T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 10,
      "page": 1,
      "per_page": 20,
      "last_page": 1
    }
  }
}
```

---

### 2. 创建角色

**请求**
```
POST /api/v1/admin/roles
```

**请求体**
```json
{
  "name": "editor",
  "display_name": "编辑",
  "description": "可以编辑内容的角色"
}
```

**参数说明**
| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| name | string | 是 | 角色唯一标识，1-50 个字符 |
| display_name | string | 是 | 角色显示名称，1-100 个字符 |
| description | string | 否 | 角色描述，最多 255 个字符 |

**响应示例**
```json
{
  "code": 0,
  "message": "角色创建成功",
  "data": {
    "id": 2,
    "name": "editor",
    "display_name": "编辑",
    "description": "可以编辑内容的角色",
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z"
  }
}
```

---

### 3. 获取角色详情

**请求**
```
GET /api/v1/admin/roles/{id}
```

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 角色 ID |

**响应示例**
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "id": 1,
    "name": "admin",
    "display_name": "管理员",
    "description": "拥有所有权限",
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z"
  }
}
```

---

### 4. 更新角色

**请求**
```
PUT /api/v1/admin/roles/{id}
```

**请求体**
```json
{
  "name": "editor",
  "display_name": "内容编辑",
  "description": "负责编辑和管理内容"
}
```

**参数说明**
| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| name | string | 否 | 角色唯一标识 |
| display_name | string | 否 | 角色显示名称 |
| description | string | 否 | 角色描述 |

**响应示例**
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "id": 2,
    "name": "editor",
    "display_name": "内容编辑",
    "description": "负责编辑和管理内容",
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T12:00:00Z"
  }
}
```

---

### 5. 删除角色

**请求**
```
DELETE /api/v1/admin/roles/{id}
```

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 角色 ID |

**响应**
```
204 No Content
```

---

### 6. 获取角色权限

**请求**
```
GET /api/v1/admin/roles/{id}/permissions
```

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 角色 ID |

**响应示例**
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "role_id": 1,
    "permissions": [
      {
        "id": 1,
        "name": "user.create",
        "display_name": "创建用户",
        "description": "可以创建新用户",
        "group": "user"
      },
      {
        "id": 2,
        "name": "user.delete",
        "display_name": "删除用户",
        "description": "可以删除用户",
        "group": "user"
      }
    ]
  }
}
```

---

### 7. 分配权限到角色

**请求**
```
POST /api/v1/admin/roles/{id}/permissions
```

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 角色 ID |

**请求体**
```json
{
  "permission_ids": [1, 2, 3, 4]
}
```

**参数说明**
| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| permission_ids | array | 是 | 权限 ID 列表 |

**响应示例**
```json
{
  "code": 0,
  "message": "权限分配成功",
  "data": null
}
```

---

## 权限管理 API

### 1. 获取权限列表

**请求**
```
GET /api/v1/admin/permissions
```

**查询参数**
| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认为 1 |
| per_page | int | 否 | 每页数量，默认为 20 |
| group | string | 否 | 权限分组过滤 |

**响应示例**
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "user.create",
        "display_name": "创建用户",
        "description": "可以创建新用户",
        "group": "user",
        "created_at": "2025-01-01T10:00:00Z",
        "updated_at": "2025-01-01T10:00:00Z"
      },
      {
        "id": 2,
        "name": "user.delete",
        "display_name": "删除用户",
        "description": "可以删除用户",
        "group": "user",
        "created_at": "2025-01-01T10:00:00Z",
        "updated_at": "2025-01-01T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 20,
      "page": 1,
      "per_page": 20,
      "last_page": 1
    }
  }
}
```

---

### 2. 创建权限

**请求**
```
POST /api/v1/admin/permissions
```

**请求体**
```json
{
  "name": "topic.publish",
  "display_name": "发布话题",
  "description": "可以发布新话题",
  "group": "topic"
}
```

**参数说明**
| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| name | string | 是 | 权限唯一标识，1-100 个字符 |
| display_name | string | 是 | 权限显示名称，1-100 个字符 |
| description | string | 否 | 权限描述，最多 255 个字符 |
| group | string | 否 | 权限分组，最多 50 个字符 |

**响应示例**
```json
{
  "code": 0,
  "message": "权限创建成功",
  "data": {
    "id": 3,
    "name": "topic.publish",
    "display_name": "发布话题",
    "description": "可以发布新话题",
    "group": "topic",
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z"
  }
}
```

---

### 3. 获取权限详情

**请求**
```
GET /api/v1/admin/permissions/{id}
```

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 权限 ID |

**响应示例**
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "id": 1,
    "name": "user.create",
    "display_name": "创建用户",
    "description": "可以创建新用户",
    "group": "user",
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z"
  }
}
```

---

### 4. 更新权限

**请求**
```
PUT /api/v1/admin/permissions/{id}
```

**请求体**
```json
{
  "name": "user.create",
  "display_name": "创建新用户",
  "description": "允许创建新的系统用户",
  "group": "user"
}
```

**参数说明**
| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| name | string | 否 | 权限唯一标识 |
| display_name | string | 否 | 权限显示名称 |
| description | string | 否 | 权限描述 |
| group | string | 否 | 权限分组 |

**响应示例**
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "id": 1,
    "name": "user.create",
    "display_name": "创建新用户",
    "description": "允许创建新的系统用户",
    "group": "user",
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T12:00:00Z"
  }
}
```

---

### 5. 删除权限

**请求**
```
DELETE /api/v1/admin/permissions/{id}
```

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 权限 ID |

**响应**
```
204 No Content
```

---

## 常见错误码

| 错误码 | HTTP 状态 | 说明 |
|--------|----------|------|
| 0 | 200/201 | 操作成功 |
| 1001 | 422 | 请求验证失败 |
| 1002 | 401 | 未授权，请先登录 |
| 1003 | 403 | 权限不足 |
| 1004 | 404 | 资源不存在 |
| 1006 | 422 | 验证失败 |
| 6001 | 500 | 服务器内部错误 |

---

## 使用示例

### cURL 示例

#### 获取角色列表
```bash
curl -X GET "http://localhost:8000/api/v1/admin/roles?page=1&per_page=20" \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json"
```

#### 创建角色
```bash
curl -X POST "http://localhost:8000/api/v1/admin/roles" \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "moderator",
    "display_name": "版主",
    "description": "内容审核和管理"
  }'
```

#### 分配权限
```bash
curl -X POST "http://localhost:8000/api/v1/admin/roles/1/permissions" \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": [1, 2, 3]
  }'
```

---

## 数据模型

### Role (角色)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint64 | 主键 |
| name | string | 角色唯一标识 |
| display_name | string | 角色显示名称 |
| description | string | 角色描述 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### Permission (权限)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint64 | 主键 |
| name | string | 权限唯一标识 |
| display_name | string | 权限显示名称 |
| description | string | 权限描述 |
| group | string | 权限分组 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### RolePermission (角色权限关系)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint64 | 主键 |
| role_id | uint64 | 角色 ID |
| permission_id | uint64 | 权限 ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

---

## 文件结构

```
app/
├── models/
│   ├── role/              # 角色模型
│   ├── permission/        # 权限模型
│   └── role_permission/   # 角色权限关系模型
├── repositories/
│   ├── role_repository.go           # 角色仓储
│   ├── permission_repository.go     # 权限仓储
│   └── role_permission_repository.go # 角色权限仓储
├── services/
│   ├── role_service.go              # 角色服务
│   ├── permission_service.go        # 权限服务
│   └── role_permission_service.go   # 角色权限服务
├── requests/
│   ├── role_request.go              # 角色请求验证
│   └── permission_request.go        # 权限请求验证
└── http/
    └── controllers/
        └── admin/
            ├── role_controller.go       # 角色控制器
            └── permission_controller.go # 权限控制器

routes/
└── admin.go                # 管理后台路由配置
```

---

## 业务流程

### 角色权限分配流程

1. **创建角色** -> POST /api/v1/admin/roles
2. **创建权限** -> POST /api/v1/admin/permissions
3. **为角色分配权限** -> POST /api/v1/admin/roles/{roleId}/permissions
4. **为用户分配角色** -> POST /api/v1/admin/users/{userId}/assign-role
5. **用户获得角色对应的权限**

---

## 注意事项

1. 所有分页查询默认每页 20 条，最大 100 条
2. 删除角色或权限时，对应的关系记录会自动清除
3. 权限 ID 和角色 ID 都必须是有效的，否则会返回 404 错误
4. 角色名称和权限名称都必须唯一，否则返回验证错误
5. 所有时间戳均为 UTC 格式
