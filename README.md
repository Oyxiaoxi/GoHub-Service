# GoHub-Service

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.25.5-blue.svg)
![Framework](https://img.shields.io/badge/Framework-Gin-00ADD8.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Test Coverage](https://img.shields.io/badge/Coverage-88%25-brightgreen.svg)
![Status](https://img.shields.io/badge/Status-Production%20Ready-success.svg)

**企业级 Go 语言论坛社区 API 服务**

[功能特性](#-功能特性) • [快速开始](#-快速开始) • [架构设计](#-架构设计) • [API 文档](#-api-文档) • [性能优化](#-性能优化) • [安全防护](#-安全防护)

</div>

---

## 📖 项目介绍

GoHub-Service 是一个基于 Go 语言和 Gin 框架构建的企业级论坛社区 API 服务。项目采用三层架构（Controller-Service-Repository），集成了完整的用户系统、内容管理、权限控制、搜索引擎等功能模块，经过 **15 项深度优化**，在性能、安全、代码质量方面达到生产环境标准。

### 💡 为什么选择 GoHub-Service？

- ✅ **生产级质量**：88%+ 测试覆盖率，经过完整的性能优化和安全加固
- ✅ **完整的业务功能**：用户、话题、评论、私信、通知、搜索一应俱全
- ✅ **企业级架构**：RBAC 权限、缓存策略、监控告警、API 签名验证
- ✅ **详尽的文档**：20+ 篇专业文档，涵盖开发、部署、优化各个方面
- ✅ **最佳实践**：遵循 Go 语言最佳实践，代码规范，易于维护扩展

---

## ✨ 功能特性

### 核心业务模块

| 模块 | 功能描述 | 状态 |
|------|---------|------|
| 👤 **用户系统** | 注册/登录、JWT 认证、个人资料、关注系统 | ✅ |
| 📝 **话题管理** | 话题 CRUD、分类、点赞、收藏、热门话题 | ✅ |
| 💬 **评论系统** | 多级评论、点赞、举报、审核机制 | ✅ |
| 📬 **私信系统** | 一对一私信、会话列表、未读提醒 | ✅ |
| 🔔 **通知系统** | 系统通知、互动通知、实时推送 | ✅ |
| 🔍 **全文搜索** | Elasticsearch 集成、高级搜索、分词 | ✅ |
| 🔐 **权限管理** | RBAC 角色权限、动态授权、资源级控制 | ✅ |
| 🔗 **友情链接** | 链接管理、排序、审核 | ✅ |

### 技术特性

#### 🚀 性能优化
- **数据库优化**：N+1 查询消除、连接池配置、慢查询监控
- **缓存策略**：Redis 多级缓存、缓存预热、缓存降级
- **并发安全**：Singleflight 防缓存击穿、资源泄漏防护
- **内存优化**：避免大结构体拷贝、内存复用、垃圾回收优化
- **监控指标**：Prometheus 集成、自定义业务指标、Grafana 可视化

#### 🔒 安全防护
- **输入验证**：13 种正则安全检测（SQL 注入、XSS、路径穿越）
- **API 签名**：HMAC-SHA256 签名验证、防重放攻击、时间戳验证
- **限流策略**：IP 限流、路由限流、自动封禁机制
- **密码安全**：Bcrypt 加密、强度评分、历史密码检查
- **安全响应头**：CSP、HSTS、X-Frame-Options 等

#### 📊 监控与日志
- **结构化日志**：Zap 日志库、TraceID 追踪、敏感信息过滤
- **性能监控**：请求耗时、数据库慢查询、缓存命中率
- **资源监控**：连接池状态、Goroutine 数量、内存使用
- **告警机制**：错误率、响应时间、资源异常自动告警

#### 🛠️ 开发工具
- **代码生成**：Model、Controller、Repository 自动生成
- **数据填充**：工厂模式、Faker 数据、批量填充
- **测试框架**：Mock 测试、集成测试、基准测试
- **API 文档**：Swagger/OpenAPI 自动生成、在线调试

---

## 🏗️ 架构设计

### 三层架构

```
┌─────────────────────────────────────────────────────┐
│                   Client (客户端)                    │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP/HTTPS
┌──────────────────────▼──────────────────────────────┐
│               Controller Layer (控制器层)             │
│  ┌──────────────────────────────────────────────┐   │
│  │ ● 参数验证  ● 请求处理  ● 响应封装  ● 错误处理 │   │
│  └──────────────────────────────────────────────┘   │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                Service Layer (业务层)                 │
│  ┌──────────────────────────────────────────────┐   │
│  │ ● 业务逻辑  ● 数据转换  ● 事务管理  ● 缓存策略 │   │
│  └──────────────────────────────────────────────┘   │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│             Repository Layer (数据访问层)             │
│  ┌──────────────────────────────────────────────┐   │
│  │ ● 数据库操作  ● 缓存操作  ● 查询优化  ● ORM 封装 │   │
│  └──────────────────────────────────────────────┘   │
└──────────────────────┬──────────────────────────────┘
                       │
          ┌────────────┴────────────┐
          │                         │
┌─────────▼────────┐    ┌──────────▼─────────┐
│  MySQL/SQLite    │    │      Redis         │
│  (主数据存储)     │    │    (缓存/会话)      │
└──────────────────┘    └────────────────────┘
```

### 技术栈

| 类别 | 技术选型 | 版本 |
|------|---------|------|
| **语言** | Go | 1.25.5 |
| **Web 框架** | Gin | 1.11.0 |
| **ORM** | GORM | 1.31.1 |
| **数据库** | MySQL / SQLite | 8.0 / 3.x |
| **缓存** | Redis | 7.0+ |
| **搜索引擎** | Elasticsearch | 8.x |
| **日志** | Zap | 1.27.1 |
| **监控** | Prometheus | 1.23.2 |
| **配置管理** | Viper | 1.21.0 |
| **命令行** | Cobra | 1.10.2 |
| **测试** | Testify | 1.11.1 |

---

## 🚀 快速开始

### 环境要求

- Go 1.25+ 
- MySQL 8.0+ 或 SQLite 3.x
- Redis 7.0+
- Elasticsearch 8.x (可选)

### 1. 克隆项目

```bash
git clone https://github.com/Oyxiaoxi/GoHub-Service.git
cd GoHub-Service
```

### 2. 配置环境

```bash
# 复制配置文件
cp .env.example .env

# 编辑配置文件，设置数据库、Redis 等
vim .env
```

**关键配置项**：
```bash
# 应用配置
APP_KEY=your-32-character-secret-key-here
APP_DEBUG=true
APP_PORT=3000

# 数据库配置
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=gohub
DB_USERNAME=root
DB_PASSWORD=your-password

# Redis 配置
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=

# API 签名密钥（生产环境必须修改）
SIGNATURE_SECRET=your-signature-secret-key-32chars
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 数据库迁移

```bash
# 运行数据库迁移
go run main.go migrate

# 填充测试数据（可选）
go run main.go seed
```

### 5. 启动服务

```bash
# 开发环境
go run main.go serve

# 或使用热重载（推荐）
air

# 生产环境
go build -o gohub main.go
./gohub serve
```

服务启动后访问：
- API 服务：http://localhost:3000
- Swagger 文档：http://localhost:3000/swagger/index.html
- Prometheus 指标：http://localhost:3000/metrics

---

## 📚 API 文档

### Swagger 文档

访问 http://localhost:3000/swagger/index.html 查看完整的 API 文档和在线调试。

### 主要 API 端点

#### 认证相关
```
POST   /api/v1/auth/signup/using-phone      # 手机注册
POST   /api/v1/auth/signup/using-email      # 邮箱注册
POST   /api/v1/auth/login/using-phone       # 手机登录
POST   /api/v1/auth/login/using-password    # 密码登录
POST   /api/v1/auth/login/refresh-token     # 刷新 Token
POST   /api/v1/auth/password-reset/*        # 密码重置（需签名）
POST   /api/v1/auth/verify-codes/*          # 验证码发送
```

#### 用户相关
```
GET    /api/v1/user                         # 当前用户信息
GET    /api/v1/users                        # 用户列表
PUT    /api/v1/users                        # 更新资料
PUT    /api/v1/users/email                  # 修改邮箱（需签名）
PUT    /api/v1/users/phone                  # 修改手机（需签名）
PUT    /api/v1/users/password               # 修改密码（需签名）
PUT    /api/v1/users/avatar                 # 更新头像
POST   /api/v1/users/:id/follow             # 关注用户
POST   /api/v1/users/:id/unfollow           # 取消关注
```

#### 话题相关
```
GET    /api/v1/topics                       # 话题列表
POST   /api/v1/topics                       # 创建话题
GET    /api/v1/topics/:id                   # 话题详情
PUT    /api/v1/topics/:id                   # 更新话题
DELETE /api/v1/topics/:id                   # 删除话题
POST   /api/v1/topics/:id/like              # 点赞话题
```

#### 评论相关
```
GET    /api/v1/comments                     # 评论列表
POST   /api/v1/comments                     # 发表评论
PUT    /api/v1/comments/:id                 # 更新评论
DELETE /api/v1/comments/:id                 # 删除评论
```

#### 私信相关
```
POST   /api/v1/messages                     # 发送私信（需签名）
GET    /api/v1/messages                     # 会话列表
POST   /api/v1/messages/read                # 标记已读
GET    /api/v1/messages/unread-count        # 未读数量
```

#### 搜索相关
```
GET    /api/v1/search/topics                # 搜索话题
GET    /api/v1/search/users                 # 搜索用户
GET    /api/v1/search/suggest               # 搜索建议
```

#### 管理后台
```
GET    /api/v1/admin/dashboard/overview     # 数据概览
GET    /api/v1/admin/users                  # 用户管理
DELETE /api/v1/admin/users/:id              # 删除用户（需签名）
POST   /api/v1/admin/users/:id/ban          # 封禁用户（需签名）
GET    /api/v1/admin/topics                 # 话题管理
```

### API 签名验证

敏感操作（密码重置、修改账号信息、管理后台操作）需要 API 签名验证：

```bash
# 请求头
X-Timestamp: 1735891200          # Unix 时间戳（秒）
X-Nonce: abc123xyz789mnop        # 随机字符串（≥16位）
X-Signature: a1b2c3d4...         # HMAC-SHA256 签名
```

**签名生成**：
```
签名字符串 = METHOD\nPATH\nTIMESTAMP\nNONCE\nBODY
签名 = HMAC-SHA256(签名字符串, SIGNATURE_SECRET)
```

**代码示例**: 查看 [API 签名验证示例代码](docs/examples/api_signature_example.go) 了解客户端实现细节。

---

## ⚡ 性能优化

GoHub-Service 经过 **15 项系统性优化**，在性能、安全、代码质量方面达到生产级标准：

### 优化成果

| 优化项 | 优化前 | 优化后 | 提升 |
|--------|-------|-------|------|
| 数据库 N+1 查询 | 大量重复查询 | Preload 预加载 | ~90% |
| 缓存命中率 | 无缓存 | 多级缓存 | ~80% |
| 代码重复率 | 高重复 | 泛型 Mapper | -75% |
| 测试覆盖率 | ~30% | **88%+** | +58% |
| 慢查询数量 | 多个 | 索引优化后清零 | 100% |
| 内存分配 | 大量拷贝 | 指针传递 | ~60% |
| 并发安全 | 缓存击穿 | Singleflight | 100% |

### 核心优化技术

1. **数据库优化**
   - N+1 查询消除（Preload/Join）
   - 连接池配置优化（最大连接数、空闲连接）
   - 慢查询监控（>200ms 自动记录）
   - 索引优化（覆盖索引、组合索引）

2. **缓存策略**
   - Redis 多级缓存（热数据、温数据）
   - 缓存预热（WarmupScheduler）
   - 缓存降级（DegradationManager）
   - Singleflight 防缓存击穿

3. **并发优化**
   - Context 超时控制
   - Goroutine 泄漏防护
   - 资源池复用（连接池、对象池）
   - 读写锁优化（sync.RWMutex）

4. **代码质量**
   - 泛型 Repository 基类（减少 75% 重复代码）
   - 统一错误处理（AppError）
   - 结构化日志（Zap + TraceID）
   - 单元测试覆盖率 88%+

**优化示例代码**：
- [数据库优化示例](docs/examples/database_optimization_examples.go)
- [代码去重示例](docs/examples/code_deduplication_examples.go)
- [资源管理示例](docs/examples/resource_management_examples.go)
- [集成优化示例](docs/examples/integrated_optimization_example.go)

---

## 🔐 安全防护

### 多层安全防护体系

```
┌─────────────────────────────────────────────────────┐
│                  应用层防护                          │
│  ● API 签名验证 (HMAC-SHA256)                       │
│  ● JWT 认证 + 刷新令牌机制                          │
│  ● RBAC 权限控制                                    │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                  输入层防护                          │
│  ● SQL 注入检测 (4种模式)                           │
│  ● XSS 攻击检测 (6种模式)                           │
│  ● 路径穿越检测 (3种模式)                           │
│  ● 参数验证 + 类型检查                              │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                  网络层防护                          │
│  ● IP 限流 (全局 200/小时)                          │
│  ● 路由限流 (敏感接口 5-20/分钟)                    │
│  ● 自动封禁机制 (1分钟)                             │
│  ● 重放攻击检测 (Nonce + Redis)                    │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                  传输层防护                          │
│  ● HTTPS 强制跳转                                   │
│  ● HSTS 响应头                                      │
│  ● CSP 内容安全策略                                 │
│  ● CORS 跨域限制                                    │
└─────────────────────────────────────────────────────┘
```

### 安全特性清单

- ✅ **防 SQL 注入**：13 种正则检测 + 参数化查询
- ✅ **防 XSS 攻击**：输入过滤 + 输出转义 + CSP
- ✅ **防重放攻击**：Nonce + 时间戳验证（5分钟）
- ✅ **防暴力破解**：密码重置限流（5次/分钟）
- ✅ **防 CSRF**：Token 验证 + SameSite Cookie
- ✅ **防路径穿越**：路径规范化 + 白名单检查
- ✅ **密码安全**：Bcrypt + 强度评分 + 历史检查
- ✅ **敏感数据**：日志脱敏 + 数据库加密
- ✅ **API 签名**：HMAC-SHA256 + 防篡改验证

**安全实现示例**：
- [API 签名验证示例](docs/examples/api_signature_example.go)
- [Postman 测试集](docs/examples/postman_api_signature_tests.json)

---

## 📊 监控与日志

### Prometheus 指标

访问 http://localhost:3000/metrics 查看实时指标：

```prometheus
# HTTP 请求
gohub_http_requests_total{method="GET",path="/api/v1/topics",status="200"} 1250
gohub_http_request_duration_seconds{method="GET",path="/api/v1/topics"} 0.025

# 数据库
gohub_db_query_duration_seconds{operation="find"} 0.012

# 缓存
gohub_cache_hits_total{cache_type="redis"} 8500
gohub_cache_misses_total{cache_type="redis"} 1500

# API 签名验证
gohub_api_signature_verifications_total{endpoint="/api/v1/users/password",result="success"} 45
gohub_api_signature_failures_total{endpoint="/api/v1/users/email",reason="signature_mismatch"} 3
gohub_replay_attacks_total{endpoint="/api/v1/messages"} 2

# 错误统计
gohub_errors_total{error_type="http_error"} 23
```

### 结构化日志

日志示例（包含 TraceID 追踪）：
```json
{
  "level": "info",
  "time": "2026-01-03T15:04:05.123+08:00",
  "trace_id": "abc123def456",
  "caller": "controllers/topic_controller.go:45",
  "msg": "Create topic success",
  "user_id": 123,
  "topic_id": 456,
  "duration_ms": 125
}
```

### 监控大盘

推荐使用 Grafana 可视化：
- HTTP 请求 QPS、响应时间
- 数据库慢查询、连接池状态
- 缓存命中率、Redis 内存使用
- API 签名验证成功率、重放攻击检测
- 错误率、资源使用情况

访问 `/metrics` 端点获取所有 Prometheus 指标

---

## 🛠️ 开发指南

### CLI 命令

```bash
# 服务管理
go run main.go serve              # 启动 HTTP 服务
go run main.go serve --env=production  # 生产环境启动

# 数据库管理
go run main.go migrate             # 运行数据库迁移
go run main.go migrate:rollback    # 回滚最后一次迁移
go run main.go migrate:refresh     # 重置数据库

# 数据填充
go run main.go seed                # 运行所有 Seeder
go run main.go seed --seeder=UserSeeder  # 运行指定 Seeder

# 代码生成
go run main.go make:model Post     # 生成 Model
go run main.go make:controller PostController  # 生成 Controller
go run main.go make:service PostService        # 生成 Service
go run main.go make:repository PostRepository  # 生成 Repository
go run main.go make:policy PostPolicy          # 生成 Policy
go run main.go make:request PostRequest        # 生成 Request
go run main.go make:seeder PostSeeder          # 生成 Seeder
go run main.go make:factory PostFactory        # 生成 Factory

# 工具命令
go run main.go key:generate        # 生成 APP_KEY
go run main.go cache:clear         # 清空缓存
go run main.go play                # 进入交互式终端
go run main.go slowlog --file=slow.log  # 分析慢查询日志
```

### 项目结构

```
GoHub-Service/
├── app/                      # 应用核心代码
│   ├── cmd/                  # CLI 命令
│   │   ├── cache.go          # 缓存命令
│   │   ├── key.go            # 密钥生成
│   │   ├── make/             # 代码生成器
│   │   ├── migrate.go        # 数据库迁移
│   │   ├── seed.go           # 数据填充
│   │   ├── serve.go          # HTTP 服务
│   │   └── slowlog.go        # 慢查询分析
│   ├── http/                 # HTTP 层
│   │   ├── controllers/      # 控制器
│   │   └── middlewares/      # 中间件
│   ├── models/               # 数据模型
│   ├── policies/             # 授权策略
│   ├── repositories/         # 数据访问层
│   ├── requests/             # 请求验证
│   └── services/             # 业务逻辑层
├── bootstrap/                # 应用初始化
│   ├── cache.go              # 缓存初始化
│   ├── database.go           # 数据库初始化
│   ├── elasticsearch.go      # ES 初始化
│   ├── logger.go             # 日志初始化
│   ├── redis.go              # Redis 初始化
│   └── route.go              # 路由初始化
├── config/                   # 配置文件
├── database/                 # 数据库相关
│   ├── factories/            # 数据工厂
│   ├── migrations/           # 迁移文件
│   └── seeders/              # 数据填充
├── docs/                     # 项目文档
├── pkg/                      # 公共包
│   ├── auth/                 # 认证
│   ├── cache/                # 缓存
│   ├── database/             # 数据库工具
│   ├── elasticsearch/        # ES 工具
│   ├── logger/               # 日志
│   ├── metrics/              # 监控指标
│   ├── redis/                # Redis 工具
│   ├── repository/           # Repository 基类
│   ├── response/             # 响应封装
│   ├── security/             # 安全工具
│   └── ...
├── routes/                   # 路由定义
├── scripts/                  # 运维脚本
├── storage/                  # 存储目录
│   ├── logs/                 # 日志文件
│   └── uploads/              # 上传文件
├── .env.example              # 环境配置示例
├── docker-compose.yml        # Docker 编排
├── go.mod                    # Go 模块
├── main.go                   # 应用入口
└── README.md                 # 项目说明
```

### 添加新模块

以创建 "文章(Article)" 模块为例：

```bash
# 1. 生成代码骨架
go run main.go make:model Article
go run main.go make:controller ArticleController
go run main.go make:service ArticleService
go run main.go make:repository ArticleRepository
go run main.go make:request ArticleRequest
go run main.go make:policy ArticlePolicy
go run main.go make:seeder ArticleSeeder
go run main.go make:factory ArticleFactory

# 2. 编写数据库迁移
go run main.go make:migration create_articles_table

# 3. 注册路由（routes/article.go）
# 4. 运行迁移和填充
go run main.go migrate
go run main.go seed --seeder=ArticleSeeder

# 5. 运行测试
go test ./app/services/article_service_test.go -v
```

查看项目代码了解更多开发规范和最佳实践。

---

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./... -v

# 运行指定包测试
go test ./app/services/... -v

# 运行带覆盖率测试
go test ./... -cover

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 运行基准测试
go test ./... -bench=. -benchmem

# 运行集成测试
go test ./tests/integration/... -v
```

### 测试覆盖率

当前测试覆盖率：**88%+**

| 模块 | 覆盖率 | 测试数量 |
|------|--------|---------|
| Services | 92% | 150+ |
| Repositories | 85% | 80+ |
| Controllers | 88% | 120+ |
| Middlewares | 90% | 40+ |
| Utilities | 95% | 60+ |

### 编写测试

```go
// Service 测试示例
func TestTopicService_Create(t *testing.T) {
    // 创建 mock repository
    mockRepo := new(mocks.TopicRepository)
    service := services.NewTopicService(mockRepo)
    
    // 设置预期
    mockRepo.On("Create", mock.Anything, mock.Anything).
        Return(&models.Topic{ID: 1}, nil)
    
    // 执行测试
    topic, err := service.Create(context.Background(), &dto.CreateTopicRequest{
        Title: "Test Topic",
        Content: "Test Content",
    })
    
    // 断言
    assert.NoError(t, err)
    assert.NotNil(t, topic)
    assert.Equal(t, uint64(1), topic.ID)
    mockRepo.AssertExpectations(t)
}
```

参考项目中的 `*_test.go` 文件了解更多测试示例。

---

## 📦 部署

### Docker 部署

```bash
# 使用 Docker Compose
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

### 生产环境部署

```bash
# 1. 编译二进制
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gohub main.go

# 2. 上传到服务器
scp gohub user@server:/opt/gohub/

# 3. 配置生产环境变量
cp .env.production.example .env
vim .env

# 4. 运行数据库迁移
./gohub migrate

# 5. 使用 systemd 管理服务
sudo cp scripts/gohub.service.example /etc/systemd/system/gohub.service
sudo systemctl enable gohub
sudo systemctl start gohub

# 6. 配置 Nginx 反向代理
sudo cp scripts/nginx.conf.example /etc/nginx/sites-available/gohub
sudo nginx -t && sudo systemctl reload nginx
```

### 生产环境检查清单

使用自动化脚本检查配置：
```bash
bash scripts/pre-deploy-check.sh
```

检查项包括：
- ✅ 环境变量完整性
- ✅ 数据库连接
- ✅ Redis 连接
- ✅ 文件权限
- ✅ 端口可用性
- ✅ 安全配置（HTTPS、密钥强度等）

查看 `scripts/` 目录获取更多部署脚本和配置示例。

---

## 📖 文档与资源

### API 文档

- **Swagger/OpenAPI**: 访问 http://localhost:3000/swagger/index.html
- **API 规范**: [docs/swagger.json](docs/swagger.json) / [docs/swagger.yaml](docs/swagger.yaml)

### 代码示例

| 示例 | 说明 |
|------|------|
| [api_signature_example.go](docs/examples/api_signature_example.go) | API 签名验证完整示例 |
| [database_optimization_examples.go](docs/examples/database_optimization_examples.go) | 数据库优化技巧 |
| [code_deduplication_examples.go](docs/examples/code_deduplication_examples.go) | 代码去重方案 |
| [resource_management_examples.go](docs/examples/resource_management_examples.go) | 资源管理最佳实践 |
| [integrated_optimization_example.go](docs/examples/integrated_optimization_example.go) | 集成优化示例 |
| [postman_api_signature_tests.json](docs/examples/postman_api_signature_tests.json) | Postman 测试集合 |

### 运维脚本

| 脚本 | 说明 |
|------|------|
| [backup-database.sh](scripts/backup-database.sh) | 数据库备份脚本 |
| [backup-files.sh](scripts/backup-files.sh) | 文件备份脚本 |
| [pre-deploy-check.sh](scripts/pre-deploy-check.sh) | 部署前检查脚本 |
| [run-tests.sh](scripts/run-tests.sh) | 测试运行脚本 |
| [test_api_signature.sh](scripts/test_api_signature.sh) | API 签名测试脚本 |
| [gohub.service.example](scripts/gohub.service.example) | Systemd 服务配置 |
| [nginx.conf.example](scripts/nginx.conf.example) | Nginx 配置示例 |

### 学习建议

1. **快速开始**: 按照本 README 的「快速开始」章节操作
2. **理解架构**: 查看「架构设计」章节和项目目录结构
3. **学习代码**: 阅读 `docs/examples/` 中的示例代码
4. **运行测试**: 查看 `*_test.go` 文件了解测试方法
5. **参考实现**: 研究 `app/` 目录下的实际业务代码
6. **部署实践**: 使用 `scripts/` 中的脚本进行部署

---

## 🤝 贡献指南

欢迎贡献代码、报告问题或提出建议！

### 参与方式

1. **报告 Bug**：在 [Issues](https://github.com/Oyxiaoxi/GoHub-Service/issues) 中创建问题
2. **提出功能建议**：在 [Discussions](https://github.com/Oyxiaoxi/GoHub-Service/discussions) 中讨论
3. **提交代码**：Fork 项目 → 创建分支 → 提交 PR

### 开发流程

```bash
# 1. Fork 项目并克隆
git clone https://github.com/YOUR_USERNAME/GoHub-Service.git
cd GoHub-Service

# 2. 创建功能分支
git checkout -b feature/your-feature-name

# 3. 提交更改
git add .
git commit -m "feat: add your feature"

# 4. 推送到 GitHub
git push origin feature/your-feature-name

# 5. 创建 Pull Request
```

### 代码规范

- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 编写单元测试（覆盖率 > 80%）
- 提交前运行 `go test ./...` 和 `golangci-lint run`

---

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

---

---

## 🌟 致谢

感谢以下优秀的开源项目：

- [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- [GORM](https://gorm.io/) - 优秀的 Go ORM 库
- [Cobra](https://github.com/spf13/cobra) - 强大的 CLI 库
- [Zap](https://github.com/uber-go/zap) - 高性能日志库
- [Prometheus](https://prometheus.io/) - 监控系统
- [Elasticsearch](https://www.elastic.co/) - 搜索引擎

---

<div align="center">

**如果这个项目对你有帮助，请给个 ⭐️ Star 支持一下！**

Made with ❤️ by [Oyxiaoxi](https://github.com/Oyxiaoxi)

</div>
