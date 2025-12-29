# GoHub-Service

基于 Gin 的论坛后端服务，采用 Service/Repository/Cache 分层，内置 Swagger 文档、限流与安全中间件、结构化日志和 Redis 缓存。

## 核心特性

- 领域分层：Service + Repository + Cache，DTO 输入/输出，Mock 友好的接口定义
- 安全与性能：配置化 CORS/限流、XSS/SQL 注入防护、Gzip 压缩、RequestID + 结构化日志
- 数据与缓存：GORM 模型与迁移、Redis 缓存(Topic/Category/Link/User)、分页助手
- 身份与验证：手机号/邮箱校验、图片验证码、短信验证码（阿里云）、JWT 鉴权
- API 文档：Swagger UI 暴露在 `/swagger/index.html`
- 开发工具：Cobra CLI（serve/migrate/seed/make）、脚手架生成器、测试基于 testify

## 环境要求

- Go 1.25.5+
- MySQL 8.0+
- Redis 6.0+

## 快速开始

1) 克隆并安装依赖

```bash

cd GoHub-Service
go mod download
```

2) 配置环境变量

```bash
cp .env.example .env
# 至少设置数据库、Redis、APP_KEY/JWT、短信或邮件等密钥
```

核心配置示例：

```env
APP_ENV=local
APP_PORT=3000
APP_KEY=please_set_a_random_key

DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=gohub
DB_USERNAME=root
DB_PASSWORD=secret

REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=please_set_jwt_secret
```

3) 初始化数据

```bash
# 运行迁移
go run main.go migrate up

# 可选：导入示例数据
go run main.go seed
```

4) 启动服务

```bash
go run main.go serve        # 等同直接 go run main.go
# 服务默认监听 http://localhost:3000
```

## 常用命令

- 运行 Web 服务：`go run main.go serve`
- 数据库迁移：`go run main.go migrate up` / `down` / `refresh`
- 数据填充：`go run main.go seed` (或指定 seeder 名称)
- 生成代码脚手架：`go run main.go make --help`
- 运行测试：`go test ./...`

## 项目结构

```
app/            # 控制器、服务、仓储、请求验证、策略
bootstrap/      # 启动期初始化（DB/Redis/Logger/Cache/Route）
config/         # 配置定义（app, db, redis, jwt, limiter, cors, sms, mail, etc.）
database/       # 迁移、工厂、种子数据
docs/           # 安全、DTO、日志、性能与服务层指南
pkg/            # 通用库（auth, cache, controller helper, paginator, logger, limiter 等）
routes/         # 路由注册
storage/        # 日志与临时文件
```

## 文档与指南

- 代码规范：[CODING_STANDARDS.md](CODING_STANDARDS.md)
- 控制器复用指南：[CONTROLLER_REUSE_GUIDE.md](CONTROLLER_REUSE_GUIDE.md)
- 服务层架构：[docs/SERVICE_LAYER_GUIDE.md](docs/SERVICE_LAYER_GUIDE.md)
- DTO 设计与实现：[docs/DTO_GUIDE.md](docs/DTO_GUIDE.md) | [docs/DTO_IMPLEMENTATION_SUMMARY.md](docs/DTO_IMPLEMENTATION_SUMMARY.md)
- 安全与限流：[docs/API_SECURITY.md](docs/API_SECURITY.md)
- 日志与性能：[docs/LOGGING_GUIDE.md](docs/LOGGING_GUIDE.md) | [docs/PERFORMANCE_OPTIMIZATION.md](docs/PERFORMANCE_OPTIMIZATION.md)
- 优化计划进展：[OPTIMIZATION_PLAN.md](OPTIMIZATION_PLAN.md)

## 开发提示

- Swagger：启动后访问 `http://localhost:3000/swagger/index.html`
- 限流/CORS：通过 config/limiter.go 与 config/cors.go 的 env 配置启用/调优
- Redis 缓存：服务层优先读缓存，写操作会刷新相关键；确保 Redis 可用
- 日志：Zap + Lumberjack，默认输出到 `storage/logs`

## License

MIT
- `POST /v1/auth/verify-codes/phone` - 发送手机验证码
