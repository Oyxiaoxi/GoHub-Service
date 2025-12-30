# 🏗️ 系统架构

完整的系统设计和代码组织说明。

## 整体架构

```
┌─────────────────────────────────────┐
│          HTTP 请求                   │
└────────────────┬────────────────────┘
                 │
┌────────────────▼────────────────────┐
│       Middleware 中间件              │
│  认证 → RBAC → 限流 → 日志 → CORS   │
└────────────────┬────────────────────┘
                 │
┌────────────────▼────────────────────┐
│      Controller 控制器               │
│    验证请求 → 调用 Service           │
└────────────────┬────────────────────┘
                 │
┌────────────────▼────────────────────┐
│      Service 业务逻辑层              │
│   处理逻辑 → 协调 Repo 和 Cache     │
└────────────────┬────────────────────┘
                 │
        ┌────────┴────────┐
        │                 │
   ┌────▼─────┐    ┌────▼────────┐
   │Repository│    │Cache(Redis) │
   │数据访问  │    │缓存层       │
   └────┬─────┘    └────┬────────┘
        │                │
        └────────┬───────┘
                 │
        ┌────────▼────────┐
        │   Database      │
        │  MySQL/SQLite   │
        └─────────────────┘
```

## 项目结构

```
app/                    # 核心应用代码
  ├── cmd/             # CLI 命令
  ├── http/            # HTTP 处理
  │   ├── controllers/ # API 控制器
  │   └── middlewares/ # 中间件（认证、RBAC、限流等）
  ├── models/          # 数据模型
  ├── requests/        # 请求验证和 DTO
  ├── services/        # 业务逻辑
  ├── repositories/    # 数据访问
  ├── policies/        # 权限策略
  └── cache/           # 缓存层

bootstrap/              # 启动初始化
  ├── database.go      # 数据库初始化
  ├── redis.go         # Redis 初始化
  ├── logger.go        # 日志初始化
  ├── cache.go         # 缓存初始化
  └── route.go         # 路由初始化

config/                 # 配置定义
  ├── app.go
  ├── database.go
  ├── redis.go
  └── ...

database/               # 数据库管理
  ├── migrations/      # 数据库迁移
  ├── seeders/         # 测试数据
  └── factories/       # 数据工厂

pkg/                    # 通用工具包
  ├── auth/            # 认证（JWT）
  ├── database/        # 数据库工具
  ├── logger/          # 日志工具
  ├── cache/           # 缓存工具
  ├── response/        # API 响应格式
  └── ...

routes/                 # 路由定义

storage/                # 数据存储
  ├── logs/            # 应用日志
  └── uploads/         # 用户上传
```

## 分层设计

### 1. HTTP 层
- **位置**: `app/http/controllers/`
- **职责**: 处理 HTTP 请求，解析参数，返回响应
- **特点**: 不包含业务逻辑，只负责数据转换

### 2. Middleware 中间件
- **位置**: `app/http/middlewares/`
- **职责**: 认证、授权、限流、日志、数据验证等
- **包含**: 
  - 认证中间件（JWT 验证）
  - RBAC 中间件（权限检查）
  - 限流中间件
  - 日志中间件

### 3. Service 业务逻辑层
- **位置**: `app/services/`
- **职责**: 实现核心业务逻辑，协调 Repository 和 Cache
- **特点**: 
  - 不直接操作数据库
  - 通过 Repository 访问数据
  - 处理缓存逻辑
  - 包含业务规则

### 4. Repository 数据访问层
- **位置**: `app/repositories/`
- **职责**: 数据持久化操作，隐藏 SQL 细节
- **特点**:
  - 封装 GORM 操作
  - 定义数据访问接口
  - 支持高级查询（分页、排序、过滤）

### 5. Cache 缓存层
- **位置**: `app/cache/`
- **职责**: 减少数据库查询，提升性能
- **支持**: 多层缓存、失效策略

## 数据流向示例

以"获取话题列表"为例：

```
GET /api/v1/topics?page=1
    ↓
Controller.Index()
    ├─ 解析参数
    └─ 调用 Service.GetPaginated()
        ↓
    Service.GetPaginated()
        ├─ 尝试从缓存获取
        ├─ 缓存未命中
        └─ 调用 Repository.FindAll()
            ↓
        Repository.FindAll()
            └─ 从数据库查询
                ↓
            返回结果给 Service
                ↓
        Service 更新缓存
            ↓
        返回结果给 Controller
            ↓
Controller 返回 JSON 响应
    ↓
HTTP 200 OK
```

## 权限控制流程

```
用户请求
    ↓
认证中间件 (验证 JWT)
    ↓
RBAC 中间件 (检查权限)
    ├─ 有权限 → 继续执行
    └─ 无权限 → 403 Forbidden
    ↓
业务逻辑处理
    ↓
返回响应
```

## 核心模块

| 模块 | 说明 | 关键文件 |
|------|------|---------|
| User | 用户管理、认证 | `models/user/` |
| Topic | 话题讨论 | `models/topic/` |
| Category | 分类管理 | `models/category/` |
| Comment | 评论系统 | `models/comment/` |
| Role | RBAC 角色 | `models/role/` |
| Permission | RBAC 权限 | `models/permission/` |

## 依赖关系

```
Gin         → Web 框架
GORM        → ORM
Redis       → 缓存
Zap         → 日志
JWT-Go      → 令牌认证
```

## 新增功能步骤

1. 创建 Model: `app/models/feature/`
2. 创建 Migration: `database/migrations/`
3. 创建 Repository: `app/repositories/`
4. 创建 Service: `app/services/`
5. 创建 Controller: `app/http/controllers/`
6. 定义 Routes: `routes/`
7. 添加中间件: `app/http/middlewares/`（如需）

---

更多详情请查看 [DEVELOPMENT.md](./DEVELOPMENT.md)
