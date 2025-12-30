<div align="center">

# 🚀 GoHub-Service

**现代化的 Go API 服务框架 · 企业级论坛后端**

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Gin-v1.9-00ADD8?style=flat)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-v1.25-00ADD8?style=flat)](https://gorm.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat)](LICENSE)

基于 **Gin** 的高性能论坛后端服务  
采用 **Service/Repository/Cache** 三层架构  
内置 **RBAC 权限系统** · **Redis 缓存** · **结构化日志**

[快速开始](./docs/QUICKSTART.md) · [系统架构](./docs/ARCHITECTURE.md) · [在线文档](./docs/README.md)

</div>

---

## ✨ 核心特性

<table>
<tr>
<td width="50%">

### 🔐 安全可靠
- **JWT 认证** - 令牌加密，自动续期
- **RBAC 权限** - 角色权限分离，灵活配置
- **防护机制** - SQL注入/XSS/CSRF防护
- **速率限制** - 可配置的API限流
- **审计日志** - 完整的操作记录

</td>
<td width="50%">

### ⚡ 高性能
- **多层缓存** - 内存 + Redis 双层缓存
- **查询优化** - 预加载、索引、批量操作
- **连接池** - 数据库连接复用
- **Gzip 压缩** - 响应自动压缩
- **异步处理** - 耗时任务后台执行

</td>
</tr>
<tr>
<td width="50%">

### 🏗️ 架构清晰
- **分层设计** - Controller → Service → Repository
- **依赖注入** - 松耦合，易测试
- **接口抽象** - Mock 友好
- **DTO 模式** - 输入输出隔离
- **RESTful API** - 标准化接口

</td>
<td width="50%">

### 🛠️ 开发友好
- **CLI 工具** - 脚手架快速生成
- **Swagger 文档** - API 自动生成文档
- **热重载** - 开发模式自动重启
- **单元测试** - 完整的测试覆盖
- **结构化日志** - Zap 高性能日志

</td>
</tr>
</table>

---

## � 快速开始

### 前置要求

- **Go** 1.20+
- **MySQL** 8.0+ 或 **SQLite** 3.0+
- **Redis** 6.0+

### 30 秒快速启动

```bash
# 1️⃣ 克隆项目
git clone https://github.com/Oyxiaoxi/GoHub-Service.git
cd GoHub-Service

# 2️⃣ 安装依赖
go mod download

# 3️⃣ 配置环境
cp .env.example .env
# 编辑 .env 文件，设置数据库和 Redis 连接信息

# 4️⃣ 初始化数据库
go run main.go migrate

# 5️⃣ 启动服务
go run main.go serve
# 🎉 服务运行在 http://localhost:3000
```

## 🏗️ 系统架构

```
┌─────────────────────────────────────┐
│          HTTP 请求                   │
└────────────────┬────────────────────┘
                 │
┌────────────────▼────────────────────┐
│     Middleware 中间件层              │
│   认证 → RBAC → 限流 → 日志 → CORS  │
└────────────────┬────────────────────┘
                 │
┌────────────────▼────────────────────┐
│      Controller 控制器               │
│    验证请求 → 调用 Service           │
└────────────────┬────────────────────┘
                 │
┌────────────────▼────────────────────┐
│      Service 业务逻辑层              │
│   处理逻辑 → 协调 Repository/Cache  │
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

### 核心模块

| 模块 | 说明 | 关键特性 |
|------|------|---------|
| 🔐 **认证授权** | JWT + RBAC | 令牌加密、角色权限、中间件保护 |
| 👤 **用户管理** | 注册/登录/信息管理 | 邮箱验证、密码重置、个人资料 |
| 💬 **话题讨论** | 论坛核心功能 | CRUD、点赞、收藏、排序 |
| 📝 **评论系统** | 多级评论 | 嵌套回复、点赞、审核 |
| 📂 **分类管理** | 内容分类 | 树形结构、权重排序 |
| 🔗 **友情链接** | 外部链接管理 | 分组、排序、状态管理 |
| 🎛️ **管理后台** | 完整后台系统 | 用户/话题/分类管理、数据统计 |

👉 **详细架构**: [ARCHITECTURE.md](./docs/ARCHITECTURE.md)

---

## 📂 项目结构

```
GoHub-Service/
├── app/                      # 核心应用代码
│   ├── cmd/                 # CLI 命令 (migrate/seed/serve/make)
│   ├── http/                # HTTP 层
│   │   ├── controllers/     # API 控制器
│   │   └── middlewares/     # 中间件 (认证/RBAC/限流/日志)
│   ├── models/              # 数据模型
│   │   ├── user/           # 用户模型
│   │   ├── topic/          # 话题模型
│   │   ├── category/       # 分类模型
│   │   ├── role/           # 角色模型
│   │   └── permission/     # 权限模型
│   ├── services/            # 业务逻辑层
│   ├── repositories/        # 数据访问层
│   ├── requests/            # 请求验证 (DTO)
│   ├── policies/            # 权限策略
│   └── cache/               # 缓存层
├── bootstrap/                # 启动初始化
│   ├── database.go         # 数据库初始化
│   ├── redis.go            # Redis 初始化
│   ├── logger.go           # 日志初始化
│   └── route.go            # 路由初始化
├── config/                   # 配置管理
├── database/                 # 数据库相关
│   ├── migrations/         # 数据库迁移
│   ├── seeders/            # 数据填充
│   └── factories/          # 数据工厂
├── docs/                     # 📚 完整文档
│   ├── README.md           # 文档导航
│   ├── QUICKSTART.md       # 快速开始
│   ├── ARCHITECTURE.md     # 系统架构
│   ├── RBAC.md             # RBAC 权限
│   ├── SECURITY.md         # 安全指南
│   ├── DEVELOPMENT.md      # 开发指南
│   ├── PERFORMANCE.md      # 性能优化
│   └── FAQ.md              # 常见问题
├── pkg/                      # 通用工具包
│   ├── auth/               # 认证工具
│   ├── cache/              # 缓存工具
│   ├── database/           # 数据库工具
│   ├── logger/             # 日志工具
│   ├── response/           # 响应格式
│   └── ...
├── routes/                   # 路由定义
├── storage/                  # 数据存储
│   ├── logs/               # 日志文件
│   └── uploads/            # 上传文件
├── .env.example             # 环境配置示例
├── go.mod                   # Go 依赖管理
└── main.go                  # 应用入口
```

---

## � RBAC 权限系统

### 权限模型

```
User (用户)
  └─→ UserRole (用户角色关联)
        └─→ Role (角色)
              └─→ RolePermission (角色权限关联)
                    └─→ Permission (权限)
```

### 默认角色

| 角色 | 说明 | 典型权限 |
|------|------|---------|
| 🔴 **admin** | 超级管理员 | 所有权限 (用户管理、系统设置、内容审核) |
| 🟡 **moderator** | 版主 | 内容管理、用户封禁、评论审核 |
| 🟢 **user** | 普通用户 | 创建话题、发表评论、个人信息管理 |
| ⚪ **guest** | 访客 | 查看公开内容 |

### 使用示例

```go
// 路由保护 - 要求特定角色
router.GET("/admin/users", 
    middlewares.Authenticate(),           // JWT 认证
    middlewares.RequireRole("admin"),     // 角色检查
    controllers.UserIndex)

// 路由保护 - 要求特定权限
router.DELETE("/topics/:id",
    middlewares.Authenticate(),
    middlewares.RequirePermission("topics.delete"),  // 权限检查
    controllers.TopicDestroy)
```

👉 **完整指南**: [RBAC.md](./docs/RBAC.md)

---

## 🛠️ 常用命令

### 服务管理

```bash
# 启动开发服务器
go run main.go serve
# 服务运行在 http://localhost:3000

# 生产环境构建
go build -o gohub main.go
./gohub serve
```

### 数据库管理

```bash
# 执行数据库迁移
go run main.go migrate

# 回滚迁移
go run main.go migrate:rollback

# 刷新数据库（清空并重新迁移）
go run main.go migrate:refresh

# 查看迁移状态
go run main.go migrate:status

# 导入测试数据
go run main.go seed
```

### 开发工具

```bash
# 生成模型
go run main.go make:model Post

# 生成迁移文件
go run main.go make:migration create_posts_table

# 生成控制器
go run main.go make:controller PostController

# 生成 Service
go run main.go make:service PostService

# 生成 Repository
go run main.go make:repository PostRepository

# 查看所有可用命令
go run main.go --help
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./app/services/...

# 显示测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📚 文档中心

完整的项目文档导航 👉 [docs/README.md](./docs/README.md)

| 文档 | 说明 |
|------|------|
| 🚀 [快速开始](./docs/QUICKSTART.md) | 5分钟搭建开发环境，包含故障排查 |
| 🏗️ [系统架构](./docs/ARCHITECTURE.md) | 理解分层设计和代码组织 |
| 🔐 [RBAC 权限](./docs/RBAC.md) | 权限系统完整实现指南 |
| 🛡️ [API 安全](./docs/SECURITY.md) | 安全最佳实践和检查清单 |
| 💻 [开发指南](./docs/DEVELOPMENT.md) | 编码规范、测试、Git 工作流 |
| ⚡ [性能优化](./docs/PERFORMANCE.md) | 缓存策略和数据库优化 |
| 🎛️ [管理后台 API](./docs/ADMIN_API.md) | 完整的管理后台接口文档 |
| ❓ [常见问题](./docs/FAQ.md) | 26+ 常见问题解答 |

---

## 🌍 环境配置

最小化配置示例（`.env` 文件）：

```env
# 应用配置
APP_NAME=GoHub-Service
APP_ENV=local
APP_KEY=your-random-32-char-key-here
APP_DEBUG=true
APP_PORT=8080

# 数据库配置 (推荐本地开发使用 SQLite)
DB_CONNECTION=sqlite
DB_SQL_FILE=database/database.db

# 生产环境使用 MySQL
# DB_CONNECTION=mysql
# DB_HOST=127.0.0.1
# DB_PORT=3306
# DB_DATABASE=gohub
# DB_USERNAME=root
# DB_PASSWORD=password

# Redis 配置
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 配置
JWT_SECRET=your-jwt-secret-key
JWT_EXPIRE_TIME=7200
JWT_MAX_REFRESH_TIME=604800

# 日志配置
LOG_LEVEL=debug
LOG_TYPE=daily
```

完整配置说明请查看 [QUICKSTART.md](./docs/QUICKSTART.md)

---

## 📞 联系我们

- 🐛 问题报告: [GitHub Issues](https://github.com/Oyxiaoxi/GoHub-Service/issues)
- 💬 讨论: [GitHub Discussions](https://github.com/Oyxiaoxi/GoHub-Service/discussions)

## 📄 License

MIT License - 详见 [LICENSE](LICENSE) 文件

---
