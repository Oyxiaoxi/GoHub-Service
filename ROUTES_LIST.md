# GoHub 完整路由列表

## 📌 路由总览

### 管理后台路由 (`/api/v1/admin/*`)
- **认证要求**: JWT Token + Admin 角色
- **端点数量**: 43个

### 公开 API 路由 (`/api/v1/*`)
- **认证**: 部分需要 JWT Token
- **端点数量**: 52个

### 版主路由 (`/api/v1/moderator/*`)
- **认证要求**: JWT Token + Moderator 角色
- **端点数量**: 4个

---

## 🔐 管理后台 API (`/api/v1/admin/*`)

### 身份验证
- 所有管理员路由都需要有效的 JWT Token
- 请求用户必须具有 `admin` 角色

---

### 1️⃣ 仪表盘 (Dashboard)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/admin` | 管理后台根路径 | ✅ Admin |
| GET | `/api/v1/admin/dashboard/overview` | 仪表盘概览统计 | ✅ Admin |
| GET | `/api/v1/admin/dashboard/recent-users` | 最近新增用户 | ✅ Admin |
| GET | `/api/v1/admin/dashboard/recent-topics` | 最近发布话题 | ✅ Admin |

---

### 2️⃣ 用户管理 (Users)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/admin/users` | 用户列表 (分页) | ✅ Admin |
| GET | `/api/v1/admin/users/:id` | 用户详情 | ✅ Admin |
| PUT | `/api/v1/admin/users/:id` | 更新用户信息 | ✅ Admin |
| DELETE | `/api/v1/admin/users/:id` | 删除用户 | ✅ Admin |
| POST | `/api/v1/admin/users/batch-delete` | 批量删除用户 | ✅ Admin |
| POST | `/api/v1/admin/users/:id/ban` | 封禁用户 | ✅ Admin |
| POST | `/api/v1/admin/users/:id/unban` | 解封用户 | ✅ Admin |
| POST | `/api/v1/admin/users/:id/reset-password` | 重置用户密码 | ✅ Admin |
| POST | `/api/v1/admin/users/:id/assign-role` | 分配角色给用户 | ✅ Admin |
| GET | `/api/v1/admin/users/:id/followers` | 用户粉丝列表 | ✅ Admin |
| GET | `/api/v1/admin/users/:id/following` | 用户关注列表 | ✅ Admin |

---

### 3️⃣ 话题管理 (Topics)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/admin/topics` | 话题列表 (分页) | ✅ Admin |
| GET | `/api/v1/admin/topics/:id` | 话题详情 | ✅ Admin |
| PUT | `/api/v1/admin/topics/:id` | 更新话题 | ✅ Admin |
| DELETE | `/api/v1/admin/topics/:id` | 删除话题 | ✅ Admin |
| POST | `/api/v1/admin/topics/batch-delete` | 批量删除话题 | ✅ Admin |
| POST | `/api/v1/admin/topics/:id/pin` | 话题置顶 | ✅ Admin |
| POST | `/api/v1/admin/topics/:id/unpin` | 取消话题置顶 | ✅ Admin |
| POST | `/api/v1/admin/topics/:id/approve` | 话题审核通过 | ✅ Admin |
| POST | `/api/v1/admin/topics/:id/reject` | 话题审核拒绝 | ✅ Admin |

---

### 4️⃣ 分类管理 (Categories)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/admin/categories` | 分类列表 | ✅ Admin |
| GET | `/api/v1/admin/categories/:id` | 分类详情 | ✅ Admin |
| POST | `/api/v1/admin/categories` | 创建分类 | ✅ Admin |
| PUT | `/api/v1/admin/categories/:id` | 更新分类 | ✅ Admin |
| DELETE | `/api/v1/admin/categories/:id` | 删除分类 | ✅ Admin |
| POST | `/api/v1/admin/categories/sort` | 分类排序 | ✅ Admin |

---

### 5️⃣ 角色管理 (Roles)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/admin/roles` | 角色列表 | ✅ Admin |
| POST | `/api/v1/admin/roles` | 创建角色 | ✅ Admin |
| GET | `/api/v1/admin/roles/:id` | 角色详情 | ✅ Admin |
| PUT | `/api/v1/admin/roles/:id` | 更新角色 | ✅ Admin |
| DELETE | `/api/v1/admin/roles/:id` | 删除角色 | ✅ Admin |
| GET | `/api/v1/admin/roles/:id/permissions` | 获取角色权限 | ✅ Admin |
| POST | `/api/v1/admin/roles/:id/permissions` | 分配权限给角色 | ✅ Admin |

---

### 6️⃣ 权限管理 (Permissions)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/admin/permissions` | 权限列表 | ✅ Admin |
| POST | `/api/v1/admin/permissions` | 创建权限 | ✅ Admin |
| GET | `/api/v1/admin/permissions/:id` | 权限详情 | ✅ Admin |
| PUT | `/api/v1/admin/permissions/:id` | 更新权限 | ✅ Admin |
| DELETE | `/api/v1/admin/permissions/:id` | 删除权限 | ✅ Admin |

---

### 7️⃣ 评论管理 (Comments)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/admin/comments` | 评论列表 (支持搜索/筛选) | ✅ Admin |
| GET | `/api/v1/admin/comments/:id` | 评论详情 | ✅ Admin |
| DELETE | `/api/v1/admin/comments/:id` | 删除评论 | ✅ Admin |
| POST | `/api/v1/admin/comments/batch-delete` | 批量删除评论 | ✅ Admin |
| GET | `/api/v1/admin/comments/stats` | 评论统计 | ✅ Admin |

---

### 8️⃣ 关注管理 (Follows)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/admin/follows` | 关注列表 (分页) | ✅ Admin |
| DELETE | `/api/v1/admin/follows/:id` | 删除关注关系 | ✅ Admin |
| GET | `/api/v1/admin/follows/stats` | 关注统计 | ✅ Admin |

---

### 9️⃣ 点赞管理 (Likes)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/admin/likes` | 点赞列表 (支持目标类型筛选) | ✅ Admin |
| DELETE | `/api/v1/admin/likes/:id` | 删除点赞 | ✅ Admin |
| GET | `/api/v1/admin/likes/stats` | 点赞统计 | ✅ Admin |

---

## 🌐 公开 API (`/api/v1/*`)

### 身份验证说明
- 某些端点需要 JWT Token (用户必须登录)
- 某些端点无需认证 (匿名访问)

---

### 1️⃣ 身份验证 (Authentication)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/auth/login/using-phone` | 使用手机号登录 | ❌ |
| POST | `/api/v1/auth/login/using-password` | 使用密码登录 | ❌ |
| POST | `/api/v1/auth/login/refresh-token` | 刷新 Token | ✅ |
| POST | `/api/v1/auth/password-reset/using-email` | 邮箱重置密码 | ❌ |
| POST | `/api/v1/auth/password-reset/using-phone` | 手机重置密码 | ❌ |
| POST | `/api/v1/auth/signup/using-phone` | 使用手机号注册 | ❌ |
| POST | `/api/v1/auth/signup/using-email` | 使用邮箱注册 | ❌ |
| POST | `/api/v1/auth/signup/phone/exist` | 检查手机号是否存在 | ❌ |
| POST | `/api/v1/auth/signup/email/exist` | 检查邮箱是否存在 | ❌ |
| POST | `/api/v1/auth/verify-codes/phone` | 发送短信验证码 | ❌ |
| POST | `/api/v1/auth/verify-codes/email` | 发送邮箱验证码 | ❌ |
| POST | `/api/v1/auth/verify-codes/captcha` | 获取图形验证码 | ❌ |

---

### 2️⃣ 用户 (Users)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/user` | 获取当前登录用户信息 | ✅ |
| GET | `/api/v1/users` | 用户列表 | ❌ |
| PUT | `/api/v1/users` | 更新个人资料 | ✅ |
| PUT | `/api/v1/users/email` | 更新邮箱 | ✅ |
| PUT | `/api/v1/users/phone` | 更新手机号 | ✅ |
| PUT | `/api/v1/users/password` | 更新密码 | ✅ |
| PUT | `/api/v1/users/avatar` | 更新头像 | ✅ |
| POST | `/api/v1/users/:id/follow` | 关注用户 | ✅ |
| POST | `/api/v1/users/:id/unfollow` | 取消关注用户 | ✅ |
| GET | `/api/v1/users/:id/comments` | 用户评论列表 | ❌ |

---

### 3️⃣ 分类 (Categories)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/categories` | 分类列表 | ❌ |
| POST | `/api/v1/categories` | 创建分类 | ✅ |
| PUT | `/api/v1/categories/:id` | 更新分类 | ✅ |
| DELETE | `/api/v1/categories/:id` | 删除分类 | ✅ |

---

### 4️⃣ 话题 (Topics)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/topics` | 话题列表 (分页/搜索) | ❌ |
| POST | `/api/v1/topics` | 创建话题 | ✅ |
| POST | `/api/v1/topics/upload-image` | 上传话题配图 | ✅ |
| GET | `/api/v1/topics/:id` | 话题详情 | ❌ |
| PUT | `/api/v1/topics/:id` | 更新话题 | ✅ |
| DELETE | `/api/v1/topics/:id` | 删除话题 | ✅ |
| POST | `/api/v1/topics/:id/like` | 点赞话题 | ✅ |
| POST | `/api/v1/topics/:id/unlike` | 取消点赞话题 | ✅ |
| POST | `/api/v1/topics/:id/favorite` | 收藏话题 | ✅ |
| POST | `/api/v1/topics/:id/unfavorite` | 取消收藏话题 | ✅ |
| POST | `/api/v1/topics/:id/view` | 增加话题浏览次数 | ❌ |
| GET | `/api/v1/topics/:id/comments` | 话题评论列表 | ❌ |

---

### 5️⃣ 评论 (Comments)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/comments` | 评论列表 | ❌ |
| GET | `/api/v1/comments/:id` | 评论详情 | ❌ |
| POST | `/api/v1/comments` | 创建评论 | ✅ |
| PUT | `/api/v1/comments/:id` | 更新评论 | ✅ |
| DELETE | `/api/v1/comments/:id` | 删除评论 | ✅ |
| POST | `/api/v1/comments/:id/like` | 点赞评论 | ✅ |
| POST | `/api/v1/comments/:id/unlike` | 取消点赞评论 | ✅ |
| GET | `/api/v1/comments/:id/replies` | 获取评论回复 | ❌ |

---

### 6️⃣ 私信 (Messages)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/messages` | 发送私信 | ✅ |
| GET | `/api/v1/messages` | 获取对话列表 | ✅ |
| POST | `/api/v1/messages/read` | 标记私信为已读 | ✅ |
| GET | `/api/v1/messages/unread-count` | 获取未读私信数 | ✅ |

---

### 7️⃣ 通知 (Notifications)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/notifications` | 通知列表 | ✅ |
| POST | `/api/v1/notifications/:id/read` | 标记单条通知为已读 | ✅ |
| POST | `/api/v1/notifications/read-all` | 标记所有通知为已读 | ✅ |

---

### 8️⃣ 友情链接 (Links)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/links` | 友情链接列表 | ❌ |

---

### 9️⃣ 搜索 (Search)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/search/topics` | 搜索话题 | ❌ |
| GET | `/api/v1/search/users` | 搜索用户 | ✅ |

---

## 🎖️ 版主 API (`/api/v1/moderator/*`)

### 身份验证
- 所有版主路由都需要有效的 JWT Token
- 请求用户必须具有 `moderator` 角色

### 可用端点

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/moderator/topics` | 话题列表 (等待审核) | ✅ Moderator |
| POST | `/api/v1/moderator/topics/:id/approve` | 审核通过话题 | ✅ Moderator |
| POST | `/api/v1/moderator/topics/:id/reject` | 审核拒绝话题 | ✅ Moderator |
| DELETE | `/api/v1/moderator/topics/:id` | 删除话题 | ✅ Moderator |

---

## 📊 统计汇总

### 路由总数: 99个

| 分类 | 路由数 | 认证类型 |
|------|--------|---------|
| 管理后台 | 43 | Admin |
| 公开 API | 52 | Mixed |
| 版主路由 | 4 | Moderator |

### 认证类型分布

| 类型 | 数量 |
|------|------|
| 无需认证 | 34 |
| 需要 JWT | 59 |
| Admin 权限 | 43 |
| Moderator 权限 | 4 |

---

## 🔑 认证说明

### JWT Token 获取
通过登录端点获取:
- `POST /api/v1/auth/login/using-phone` - 手机号登录
- `POST /api/v1/auth/login/using-password` - 密码登录

### Token 刷新
- `POST /api/v1/auth/login/refresh-token` - 刷新过期的 Token

### 角色说明
- **admin**: 管理员，拥有所有管理权限
- **moderator**: 版主，只能审核话题
- **user**: 普通用户，默认角色

---

## 📝 请求格式说明

### 公共参数

#### 分页参数
```
?page=1&per_page=20
```

#### 搜索参数
根据不同模块支持:
- `keyword` - 关键词
- `user_id` - 用户ID
- `topic_id` - 话题ID
- `category_id` - 分类ID

---

## ✅ 返回格式说明

### 成功响应 (200)
```json
{
  "code": 0,
  "message": "success",
  "data": {...}
}
```

### 错误响应 (4xx/5xx)
```json
{
  "code": 400,
  "message": "error message"
}
```

### 分页响应
```json
{
  "data": [...],
  "pagination": {
    "total": 100,
    "page": 1,
    "per_page": 20,
    "last_page": 5
  }
}
```

---

## 🔗 关联资源

- [RBAC权限系统文档](docs/RBAC.md)
- [API开发指南](docs/DEVELOPMENT.md)
- [快速开始](docs/QUICKSTART.md)
- [安全说明](docs/SECURITY.md)

---

*最后更新: 2025年12月31日*
