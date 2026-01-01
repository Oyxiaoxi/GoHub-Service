# 🏗️ 系统架构设计

**最后更新**: 2026年1月1日 | **版本**: v2.0

---

## 📖 目录

1. [架构概览](#架构概览)
2. [分层架构](#分层架构)
3. [核心模块](#核心模块)
4. [数据流](#数据流)
5. [扩展性设计](#扩展性设计)

---

## 🎯 架构概览

GoHub-Service 是一个基于Go语言的现代化社区论坛系统，采用**分层架构 + 微服务就绪**的设计。

### 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| **Web框架** | Gin v1.9+ | 高性能HTTP框架 |
| **ORM** | GORM | 数据库操作抽象层 |
| **数据库** | MySQL 8.0+ | 关系型数据存储 |
| **缓存** | Redis 7.0+ | 分布式缓存 |
| **搜索** | Elasticsearch 8.5+ | 全文检索引擎 |
| **认证** | JWT + bcrypt | 身份验证与加密 |
| **日志** | 自定义Logger | 结构化日志系统 |

### 架构图

```
┌─────────────────────────────────────────┐
│         HTTP客户端 / API调用者            │
└─────────────┬───────────────────────────┘
              │
┌─────────────v───────────────────────────┐
│    HTTP层 (Gin Router + Middleware)     │
│  - 路由管理                             │
│  - 中间件 (认证、限流、CORS)             │
│  - 错误处理                             │
└─────────────┬───────────────────────────┘
              │
┌─────────────v───────────────────────────┐
│    控制器层 (Controllers)                │
│  - 请求参数解析                         │
│  - 业务调用                             │
│  - 响应格式化                           │
└─────────────┬───────────────────────────┘
              │
┌─────────────v───────────────────────────┐
│    服务层 (Services)                    │
│  - 业务逻辑处理                         │
│  - 事务管理                             │
│  - 权限检查                             │
└──────┬─────┬─────┬──────┬────────────────┘
       │     │     │      │
   ┌───v─┐ ┌─v──┐ │ ┌────v─┐
   │缓存 │ │日志│ │ │搜索  │
   │服务 │ │服务│ │ │服务  │
   └─┬───┘ └────┘ │ └──────┘
     │            │
   ┌─v────────────v──────────────────────┐
   │   数据访问层 (Repositories)          │
   │  - 数据库查询                        │
   │  - 数据映射                          │
   │  - 缓存策略                          │
   └─┬─────┬─────┬──────────────────────┘
     │     │     │
 ┌───v──┐ │  ┌──v────┐
 │MySQL │ │  │Redis  │
 │数据库 │ │  │缓存   │
 └──────┘ │  └───────┘
          │
      ┌───v────────────┐
      │Elasticsearch  │
      │搜索引擎        │
      └────────────────┘
```

---

## 🏛️ 分层架构

### 1. HTTP层 (Presentation Layer)

**职责**: HTTP请求处理

```
gin.Engine
  ├── 路由管理
  ├── 中间件栈
  │   ├── 日志中间件
  │   ├── 认证中间件
  │   ├── 限流中间件
  │   ├── CORS中间件
  │   └── 错误处理中间件
  └── 请求处理
```

### 2. 控制器层 (Controller Layer)

**职责**: 请求解析与响应格式化

```
type UserController struct {
    userService *services.UserService
}

func (uc *UserController) GetByID(c *gin.Context) {
    // 1. 解析参数
    // 2. 调用服务
    // 3. 格式化响应
}
```

### 3. 服务层 (Service Layer)

**职责**: 核心业务逻辑

```
type TopicService struct {
    repo        *TopicRepository
    cacheService *CacheService
    logger      *Logger
}

// 处理业务规则、事务、权限等
func (s *TopicService) Create(ctx context.Context, req *CreateTopicRequest) error
```

### 4. 数据访问层 (Repository Layer)

**职责**: 数据库操作抽象

```
type UserRepository struct {
    DB *gorm.DB
}

// CRUD操作、查询优化、缓存同步
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error)
```

### 5. 基础设施层 (Infrastructure Layer)

**职责**: 外部服务集成

```
- 数据库连接池
- Redis缓存客户端
- Elasticsearch客户端
- 日志系统
- JWT令牌管理
```

---

## 🧩 核心模块

### 用户模块 (User Module)

```
User
├── Profile (个人资料)
├── Authentication (认证)
└── Authorization (授权)
    ├── Roles (角色)
    └── Permissions (权限)
```

### 话题模块 (Topic Module)

```
Topic
├── Content (内容)
├── Metadata (元数据)
│   ├── ViewCount
│   ├── LikeCount
│   └── CommentCount
├── Category (分类)
└── Tags (标签)
```

### 评论模块 (Comment Module)

```
Comment
├── Content (评论内容)
├── Thread (评论树)
│   ├── ParentID (父评论)
│   └── Replies (回复)
└── Interaction (交互)
    └── Likes (点赞)
```

### 权限模块 (RBAC Module)

```
Role ──一对多──> Permission
  ↓
User ──一对多──> Role

权限检查流程:
User → Roles → Permissions → 访问控制
```

---

## 🔄 数据流

### 创建话题流程

```
1. 前端发送请求
   POST /api/topics
   Content: { title, body, category_id }

2. HTTP层
   ├─ 路由匹配
   ├─ 中间件处理 (认证、验证)
   └─ TopicController.Create()

3. 控制器层
   ├─ 解析请求体
   ├─ 参数验证
   └─ 调用 TopicService.Create()

4. 服务层
   ├─ 业务验证
   │  ├─ 检查权限
   │  ├─ 敏感词过滤
   │  └─ 内容清理
   ├─ 调用 TopicRepository.Create()
   ├─ 更新缓存
   └─ 发送事件

5. 数据访问层
   ├─ 开启事务
   ├─ 数据库INSERT
   ├─ 返回创建的对象
   └─ 提交事务

6. 缓存层
   ├─ 设置话题缓存
   ├─ 更新热点缓存
   └─ 清除列表缓存

7. 搜索层
   └─ 索引新话题到Elasticsearch

8. 响应
   ├─ 格式化响应
   └─ 返回 201 Created
```

### 搜索话题流程

```
1. 前端查询
   GET /api/search/topics?q=golang

2. 搜索控制器
   ├─ 参数验证
   └─ 调用搜索服务

3. 搜索服务
   ├─ 先查缓存
   ├─ 查询Elasticsearch (主路径)
   │  ├─ 全文搜索
   │  ├─ 过滤
   │  ├─ 聚合
   │  └─ 排序
   ├─ 如果ES不可用 → 数据库降级查询
   └─ 结果缓存

4. 结果处理
   ├─ 数据丰富化 (用户信息、关联数据)
   ├─ 格式化
   └─ 返回响应
```

---

## 📈 扩展性设计

### 水平扩展

```
┌──────────┐
│ 负载均衡  │
│ (Nginx)  │
└────┬─────┘
     │
  ┌──┴──┬──────┬──────┐
  │     │      │      │
┌─v──┐┌─v──┐┌─v──┐┌─v──┐
│App1││App2││App3││App4│ (Go应用实例)
└──┬─┘└──┬─┘└──┬─┘└──┬─┘
   │     │     │     │
   └─────┴──┬──┴─────┘
       ┌────v────────────┐
       │  共享资源        │
       ├─────────────────┤
       │ MySQL (主从)    │
       │ Redis (集群)    │
       │ Elasticsearch   │
       │ (多节点)        │
       └─────────────────┘
```

### 垂直扩展

```
当前架构可演进为微服务:

┌──────────────┐
│ API Gateway  │
└─────┬────────┘
      │
  ┌───┴───┬────────┬─────────┐
  │       │        │         │
┌─v─────┐│       │        │
│User   ││Topic  │Comment  │
│Service││Service│Service  │
└───────┘│       │        │
  │  ├────────┐ │        │
  │  │Shared  │ │        │
  │  │ Library│ │        │
  │  └────────┘ │        │
```

---

## 🔒 安全架构

```
请求 → 限流检查 → CORS检查 → 认证 → 授权 → 业务逻辑 → 响应
         ↓           ↓         ↓       ↓
       429        403        401    403
```

---

**架构版本**: v2.0  
**最后更新**: 2026年1月1日  
*由GoHub Architecture Team维护* ✨
