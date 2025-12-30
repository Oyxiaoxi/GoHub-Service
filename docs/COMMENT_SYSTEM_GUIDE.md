# 评论系统文档

## 概述

GoHub-Service 评论系统提供完整的话题评论功能，包括评论发布、回复、点赞、删除等操作，支持多级评论（一级评论和回复）。

## 功能特性

### ✅ 已实现功能

1. **评论管理**
   - 发布评论
   - 编辑评论
   - 删除评论
   - 查看评论详情

2. **回复功能**
   - 回复话题（一级评论）
   - 回复评论（二级回复）
   - 查看评论的回复列表

3. **互动功能**
   - 评论点赞
   - 取消点赞
   - 点赞数统计

4. **列表查询**
   - 获取所有评论
   - 获取话题的评论列表
   - 获取用户的评论列表
   - 获取评论的回复列表

5. **缓存支持**
   - 评论详情缓存
   - 话题评论列表缓存
   - 自动缓存失效机制

6. **权限控制**
   - 用户只能修改/删除自己的评论
   - Policy 层权限验证

---

## 数据库结构

### comments 表

```sql
CREATE TABLE comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    topic_id BIGINT NOT NULL,           -- 话题ID
    user_id BIGINT NOT NULL,            -- 用户ID
    content TEXT NOT NULL,              -- 评论内容
    parent_id BIGINT DEFAULT 0,         -- 父评论ID，0表示顶级评论
    like_count INT DEFAULT 0,           -- 点赞数
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    
    INDEX idx_comments_topic_id (topic_id),
    INDEX idx_comments_user_id (user_id),
    INDEX idx_comments_parent_id (parent_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (topic_id) REFERENCES topics(id)
);
```

---

## API 接口

### 基础配置

**Base URL**: `/api/v1`

**认证方式**: Bearer Token（除查询接口外都需要）

### 接口列表

#### 1. 获取评论列表

```http
GET /comments
```

**查询参数**:
- `page`: 页码，默认 1
- `per_page`: 每页数量，默认 15

**响应示例**:
```json
{
  "data": [
    {
      "id": "1",
      "topic_id": "1",
      "user_id": "1",
      "content": "这是一条评论",
      "parent_id": "0",
      "like_count": 5,
      "created_at": "2025-12-30T12:00:00Z",
      "updated_at": "2025-12-30T12:00:00Z"
    }
  ],
  "pager": {
    "current_page": 1,
    "per_page": 15,
    "total": 100
  }
}
```

---

#### 2. 获取评论详情

```http
GET /comments/:id
```

**响应示例**:
```json
{
  "id": "1",
  "topic_id": "1",
  "user_id": "1",
  "content": "这是一条评论",
  "parent_id": "0",
  "like_count": 5,
  "created_at": "2025-12-30T12:00:00Z",
  "updated_at": "2025-12-30T12:00:00Z"
}
```

---

#### 3. 发布评论

```http
POST /comments
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "topic_id": "1",
  "content": "这是一条评论",
  "parent_id": "0"  // 可选，0或不传表示顶级评论
}
```

**验证规则**:
- `topic_id`: 必填，话题必须存在
- `content`: 必填，1-1000字
- `parent_id`: 可选，数字类型

**响应**: 201 Created，返回新创建的评论对象

---

#### 4. 更新评论

```http
PUT /comments/:id
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "content": "修改后的评论内容"
}
```

**权限**: 只能修改自己的评论

**响应**: 200 OK，返回更新后的评论对象

---

#### 5. 删除评论

```http
DELETE /comments/:id
Authorization: Bearer {token}
```

**权限**: 只能删除自己的评论

**响应**: 200 OK

---

#### 6. 点赞评论

```http
POST /comments/:id/like
Authorization: Bearer {token}
```

**响应**: 200 OK

---

#### 7. 取消点赞

```http
POST /comments/:id/unlike
Authorization: Bearer {token}
```

**响应**: 200 OK

---

#### 8. 获取话题的评论列表

```http
GET /topics/:topic_id/comments
```

**查询参数**: 同接口1

**说明**: 只返回该话题的顶级评论（parent_id = 0）

---

#### 9. 获取用户的评论列表

```http
GET /users/:user_id/comments
```

**查询参数**: 同接口1

**说明**: 返回该用户发布的所有评论

---

#### 10. 获取评论的回复列表

```http
GET /comments/:comment_id/replies
```

**查询参数**: 同接口1

**说明**: 返回该评论的所有回复（parent_id = comment_id）

---

## 代码结构

```
app/
├── models/comment/
│   └── comment_model.go          # 评论模型
├── repositories/
│   └── comment_repository.go     # 数据访问层
├── services/
│   └── comment_service.go        # 业务逻辑层
├── http/controllers/api/v1/
│   └── comments_controller.go    # 控制器
├── requests/
│   └── comment_request.go        # 请求验证
├── policies/
│   └── comment_policy.go         # 权限策略
└── cache/
    └── comment_cache.go          # 缓存层

database/
├── migrations/
│   └── 2025_12_30_add_comments_table.go  # 数据库迁移
├── factories/
│   └── comment_factory.go        # 测试工厂
└── seeders/
    └── comments_seeder.go        # 数据填充

routes/
└── comment.go                    # 路由注册
```

---

## 使用示例

### 1. 发布评论

```bash
curl -X POST http://localhost:3000/api/v1/comments \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "topic_id": "1",
    "content": "这是一条评论"
  }'
```

### 2. 回复评论

```bash
curl -X POST http://localhost:3000/api/v1/comments \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "topic_id": "1",
    "content": "这是一条回复",
    "parent_id": "5"
  }'
```

### 3. 获取话题评论

```bash
curl http://localhost:3000/api/v1/topics/1/comments?page=1&per_page=20
```

### 4. 点赞评论

```bash
curl -X POST http://localhost:3000/api/v1/comments/1/like \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 缓存策略

评论系统实现了多层缓存：

1. **单个评论缓存**
   - Key: `comment:{id}`
   - TTL: 根据配置的分层缓存时间

2. **话题评论列表缓存**
   - Key: `comment:topic:{topic_id}`
   - 失效时机: 话题有新评论时

3. **自动失效**
   - 创建评论时清除话题评论列表缓存
   - 更新评论时清除该评论和话题缓存
   - 删除评论时清除该评论和话题缓存

---

## 性能优化

1. **数据库索引**
   - `idx_comments_topic_id`: 话题评论查询
   - `idx_comments_user_id`: 用户评论查询
   - `idx_comments_parent_id`: 回复查询
   - `idx_comments_created_at`: 时间排序

2. **分页加载**
   - 默认每页 15 条
   - 支持自定义分页大小

3. **预加载关联**
   - 查询时自动预加载用户和话题信息
   - 减少 N+1 查询问题

---

## 数据填充

### 运行迁移

```bash
go run main.go migrate up
```

### 填充测试数据

```bash
# 填充所有数据（包括评论）
go run main.go seed

# 只填充评论数据
go run main.go seed SeedCommentsTable
```

默认会创建 50 条测试评论数据。

---

## 权限控制

评论系统使用 Policy 模式进行权限控制：

```go
// 检查是否可以修改评论
func CanModifyComment(c *gin.Context, comment comment.Comment) bool {
    return auth.CurrentUID(c) == comment.UserID
}
```

在控制器中使用：
```go
if !policies.CanModifyComment(c, *commentModel) {
    response.Abort403(c)
    return
}
```

---

## 错误处理

评论系统使用统一的错误处理机制：

| 错误码 | 说明 |
|--------|------|
| 1004 | 资源不存在 |
| 1003 | 权限不足 |
| 5001 | 数据库操作失败 |
| 1006 | 请求参数验证失败 |

**错误响应示例**:
```json
{
  "message": "评论不存在"
}
```

---

## 扩展建议

### 待实现功能

1. **高级功能**
   - [ ] 敏感词过滤
   - [ ] 评论审核机制
   - [ ] 评论举报
   - [ ] @提及用户
   - [ ] 评论通知

2. **性能优化**
   - [ ] 评论数缓存（话题的评论总数）
   - [ ] 热门评论排序
   - [ ] 评论树形结构优化

3. **管理功能**
   - [ ] 批量删除评论
   - [ ] 管理员删除任意评论
   - [ ] 评论统计分析

---

## 注意事项

1. **删除策略**: 当前删除评论时不会级联删除回复，建议实现软删除或级联删除
2. **点赞去重**: 当前点赞没有做用户级别的去重，建议添加点赞记录表
3. **内容长度**: 评论内容限制 1000 字，根据实际需求调整
4. **评论层级**: 当前只支持两级（评论和回复），如需多级需要修改结构

---

## 测试建议

```bash
# 运行所有测试
go test ./...

# 运行评论相关测试
go test ./app/repositories -run Comment
go test ./app/services -run Comment
```

建议编写的测试用例：
- Repository 层的 CRUD 测试
- Service 层的业务逻辑测试
- Controller 层的 API 测试
- 缓存层的缓存命中测试

---

## 更新日志

### v1.0.0 (2025-12-30)
- ✅ 实现评论基础 CRUD
- ✅ 实现评论回复功能
- ✅ 实现评论点赞功能
- ✅ 实现多种查询接口
- ✅ 添加缓存支持
- ✅ 添加权限控制
- ✅ 完成数据库迁移和种子数据
