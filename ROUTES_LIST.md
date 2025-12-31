# GoHub 完整路由列表

## 📌 管理后台路由 (`/api/v1/admin/*`)

### 认证要求
- 需要有效的 JWT token
- 需要 `admin` 角色

---

## 1️⃣ 仪表盘 (Dashboard)

| 方法 | 路由 | 说明 | 实现 |
|------|------|------|------|
| GET | `/api/v1/admin` | 管理后台根路径 | ✅ |
| GET | `/api/v1/admin/dashboard/overview` | 仪表盘概览 | ✅ |
| GET | `/api/v1/admin/dashboard/recent-users` | 最近新增用户 | ✅ |
| GET | `/api/v1/admin/dashboard/recent-topics` | 最近发布话题 | ✅ |

---

## 2️⃣ 用户管理 (Users)

| 方法 | 路由 | 说明 | 实现 |
|------|------|------|------|
| GET | `/api/v1/admin/users` | 用户列表 | ✅ |
| GET | `/api/v1/admin/users/:id` | 用户详情 | ✅ |
| PUT | `/api/v1/admin/users/:id` | 更新用户 | ✅ |
| DELETE | `/api/v1/admin/users/:id` | 删除用户 | ✅ |
| POST | `/api/v1/admin/users/batch-delete` | 批量删除用户 | ✅ |
| POST | `/api/v1/admin/users/:id/ban` | 封禁用户 | ✅ |
| POST | `/api/v1/admin/users/:id/unban` | 解封用户 | ✅ |
| POST | `/api/v1/admin/users/:id/reset-password` | 重置用户密码 | ✅ |
| POST | `/api/v1/admin/users/:id/assign-role` | 分配角色给用户 | ✅ |

---

## 3️⃣ 话题管理 (Topics)

| 方法 | 路由 | 说明 | 实现 |
|------|------|------|------|
| GET | `/api/v1/admin/topics` | 话题列表 | ✅ |
| GET | `/api/v1/admin/topics/:id` | 话题详情 | ✅ |
| PUT | `/api/v1/admin/topics/:id` | 更新话题 | ✅ |
| DELETE | `/api/v1/admin/topics/:id` | 删除话题 | ✅ |
| POST | `/api/v1/admin/topics/batch-delete` | 批量删除话题 | ✅ |
| POST | `/api/v1/admin/topics/:id/pin` | 话题置顶 | ✅ |
| POST | `/api/v1/admin/topics/:id/unpin` | 取消话题置顶 | ✅ |
| POST | `/api/v1/admin/topics/:id/approve` | 话题审核通过 | ✅ |
| POST | `/api/v1/admin/topics/:id/reject` | 话题审核拒绝 | ✅ |

---

## 4️⃣ 分类管理 (Categories)

| 方法 | 路由 | 说明 | 实现 |
|------|------|------|------|
| GET | `/api/v1/admin/categories` | 分类列表 | ✅ |
| GET | `/api/v1/admin/categories/:id` | 分类详情 | ✅ |
| POST | `/api/v1/admin/categories` | 创建分类 | ✅ |
| PUT | `/api/v1/admin/categories/:id` | 更新分类 | ✅ |
| DELETE | `/api/v1/admin/categories/:id` | 删除分类 | ✅ |
| POST | `/api/v1/admin/categories/sort` | 分类排序 | ✅ |

---

## 5️⃣ 角色管理 (Roles) ⭐ NEW

| 方法 | 路由 | 说明 | 实现 |
|------|------|------|------|
| GET | `/api/v1/admin/roles` | 获取角色列表 | ✅ |
| POST | `/api/v1/admin/roles` | 创建角色 | ✅ |
| GET | `/api/v1/admin/roles/:id` | 获取角色详情 | ✅ |
| PUT | `/api/v1/admin/roles/:id` | 更新角色 | ✅ |
| DELETE | `/api/v1/admin/roles/:id` | 删除角色 | ✅ |
| GET | `/api/v1/admin/roles/:id/permissions` | 获取角色权限 | ✅ |
| POST | `/api/v1/admin/roles/:id/permissions` | 分配权限到角色 | ✅ |

---

## 6️⃣ 权限管理 (Permissions) ⭐ NEW

| 方法 | 路由 | 说明 | 实现 |
|------|------|------|------|
| GET | `/api/v1/admin/permissions` | 获取权限列表 | ✅ |
| POST | `/api/v1/admin/permissions` | 创建权限 | ✅ |
| GET | `/api/v1/admin/permissions/:id` | 获取权限详情 | ✅ |
| PUT | `/api/v1/admin/permissions/:id` | 更新权限 | ✅ |
| DELETE | `/api/v1/admin/permissions/:id` | 删除权限 | ✅ |

---

## 📌 版主路由 (`/api/v1/moderator/*`)

### 认证要求
- 需要有效的 JWT token
- 需要 `moderator` 角色

| 方法 | 路由 | 说明 | 实现 |
|------|------|------|------|
| GET | `/api/v1/moderator/topics` | 话题列表（待审核） | ✅ |
| POST | `/api/v1/moderator/topics/:id/approve` | 话题审核通过 | ✅ |
| POST | `/api/v1/moderator/topics/:id/reject` | 话题审核拒绝 | ✅ |
| DELETE | `/api/v1/moderator/topics/:id` | 删除话题 | ✅ |

---

## 📌 公开 API 路由 (`/api/v1/*`)

### 认证要求
- 大部分接口无需认证，部分接口需要 JWT token

---

### 认证相关 (Auth)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录 | ❌ |
| POST | `/api/v1/auth/signup` | 用户注册 | ❌ |
| POST | `/api/v1/auth/logout` | 用户登出 | ✅ JWT |
| POST | `/api/v1/auth/password/forgot` | 忘记密码 | ❌ |
| POST | `/api/v1/auth/password/reset` | 重置密码 | ❌ |
| POST | `/api/v1/auth/verify-code` | 获取验证码 | ❌ |
| POST | `/api/v1/auth/verify-code/verify` | 验证码校验 | ❌ |

---

### 用户相关 (Users)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/users` | 用户列表（公开） | ❌ |
| PUT | `/api/v1/users` | 更新个人信息 | ✅ JWT |
| PUT | `/api/v1/users/email` | 修改邮箱 | ✅ JWT |
| PUT | `/api/v1/users/phone` | 修改手机号 | ✅ JWT |
| PUT | `/api/v1/users/password` | 修改密码 | ✅ JWT |
| PUT | `/api/v1/users/avatar` | 修改头像 | ✅ JWT |
| POST | `/api/v1/users/:id/follow` | 关注用户 | ✅ JWT |
| POST | `/api/v1/users/:id/unfollow` | 取消关注用户 | ✅ JWT |

---

### 分类相关 (Categories)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/categories` | 分类列表 | ❌ |
| GET | `/api/v1/categories/:id` | 分类详情 | ❌ |

---

### 话题相关 (Topics)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/topics` | 话题列表 | ❌ |
| GET | `/api/v1/topics/:id` | 话题详情 | ❌ |
| POST | `/api/v1/topics` | 创建话题 | ✅ JWT |
| PUT | `/api/v1/topics/:id` | 更新话题 | ✅ JWT |
| DELETE | `/api/v1/topics/:id` | 删除话题 | ✅ JWT |
| POST | `/api/v1/topics/:id/like` | 点赞话题 | ✅ JWT |
| POST | `/api/v1/topics/:id/unlike` | 取消点赞话题 | ✅ JWT |

---

### 评论相关 (Comments)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/topics/:id/comments` | 获取话题评论 | ❌ |
| POST | `/api/v1/topics/:id/comments` | 发表评论 | ✅ JWT |
| PUT | `/api/v1/comments/:id` | 更新评论 | ✅ JWT |
| DELETE | `/api/v1/comments/:id` | 删除评论 | ✅ JWT |
| POST | `/api/v1/comments/:id/like` | 点赞评论 | ✅ JWT |
| POST | `/api/v1/comments/:id/unlike` | 取消点赞评论 | ✅ JWT |

---

### 通知相关 (Notifications)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/notifications` | 通知列表 | ✅ JWT |
| GET | `/api/v1/notifications/unread-count` | 未读通知数 | ✅ JWT |
| PUT | `/api/v1/notifications/:id/read` | 标记通知为已读 | ✅ JWT |
| PUT | `/api/v1/notifications/read-all` | 标记所有通知为已读 | ✅ JWT |
| DELETE | `/api/v1/notifications/:id` | 删除通知 | ✅ JWT |

---

### 链接相关 (Links)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/links` | 链接列表 | ❌ |
| GET | `/api/v1/links/:id` | 链接详情 | ❌ |

---

### 搜索相关 (Search)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/search` | 全局搜索 | ❌ |
| GET | `/api/v1/search/topics` | 搜索话题 | ❌ |
| GET | `/api/v1/search/users` | 搜索用户 | ❌ |
| GET | `/api/v1/search/categories` | 搜索分类 | ❌ |

---

### 消息相关 (Messages)

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/messages` | 私信列表 | ✅ JWT |
| POST | `/api/v1/messages` | 发送私信 | ✅ JWT |
| GET | `/api/v1/messages/:id` | 私信详情 | ✅ JWT |
| DELETE | `/api/v1/messages/:id` | 删除私信 | ✅ JWT |

---

## 📊 路由统计

### 管理后台路由
- 仪表盘: 4 个
- 用户管理: 9 个
- 话题管理: 9 个
- 分类管理: 6 个
- 角色管理: 7 个 ⭐
- 权限管理: 5 个 ⭐
- **小计: 40 个**

### 版主路由
- **小计: 4 个**

### 公开 API 路由
- 认证: 7 个
- 用户: 8 个
- 分类: 2 个
- 话题: 7 个
- 评论: 6 个
- 通知: 5 个
- 链接: 2 个
- 搜索: 4 个
- 消息: 4 个
- **小计: 45 个**

### 其他
- Prometheus 指标: 1 个

---

## 🔗 总计: **90+ 个路由**

---

## 📝 路由命名规范

1. **管理后台**: `/api/v1/admin/{resource}` 
   - 需要 admin 权限
   - RESTful 风格

2. **版主路由**: `/api/v1/moderator/{resource}`
   - 需要 moderator 权限
   - RESTful 风格

3. **公开 API**: `/api/v1/{resource}`
   - 部分需要认证
   - RESTful 风格

---

## 🔐 权限要求总结

| 路由前缀 | 认证 | 权限检查 | 说明 |
|---------|------|---------|------|
| `/api/v1/admin/*` | ✅ JWT | ✅ admin | 管理员专用 |
| `/api/v1/moderator/*` | ✅ JWT | ✅ moderator | 版主专用 |
| `/api/v1/auth/*` | ❌ | ❌ | 公开接口 |
| `/api/v1/users*` | 混合 | 取决于操作 | 公开列表，个人操作需认证 |
| `/api/v1/topics*` | 混合 | 取决于操作 | 公开列表，操作需认证 |
| `/api/v1/categories*` | ❌ | ❌ | 公开接口 |
| `/api/v1/*` | 混合 | 取决于操作 | 大部分公开，部分需认证 |
| `/metrics` | ❌ | ❌ | Prometheus 指标 |

---

## ✨ 新增路由 (本次实现)

- ✅ `GET /api/v1/admin/roles` - 获取角色列表
- ✅ `POST /api/v1/admin/roles` - 创建角色
- ✅ `GET /api/v1/admin/roles/:id` - 获取角色详情
- ✅ `PUT /api/v1/admin/roles/:id` - 更新角色
- ✅ `DELETE /api/v1/admin/roles/:id` - 删除角色
- ✅ `GET /api/v1/admin/roles/:id/permissions` - 获取角色权限
- ✅ `POST /api/v1/admin/roles/:id/permissions` - 分配权限
- ✅ `GET /api/v1/admin/permissions` - 获取权限列表
- ✅ `POST /api/v1/admin/permissions` - 创建权限
- ✅ `GET /api/v1/admin/permissions/:id` - 获取权限详情
- ✅ `PUT /api/v1/admin/permissions/:id` - 更新权限
- ✅ `DELETE /api/v1/admin/permissions/:id` - 删除权限
