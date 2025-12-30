# 🚀 快速开始指南

5 分钟快速搭建 GoHub-Service 开发环境。

## 📋 前置要求

- Go 1.25.5+
- 数据库：MySQL 8.0+ 或 SQLite（推荐本地开发）
- Redis 6.0+（可选，非必需）

## ⚙️ 安装步骤

### 1. 克隆项目

```bash
git clone https://github.com/Oyxiaoxi/GoHub-Service.git
cd GoHub-Service
go mod download
```

### 2. 配置环境

```bash
cp .env.example .env
```

**最小化配置** (SQLite + 本地开发):
```env
APP_NAME=GoHub-Service
APP_ENV=local
APP_KEY=your-random-key-here
APP_DEBUG=true
APP_PORT=3000

DB_CONNECTION=sqlite
DB_SQL_FILE=database/database.db

REDIS_HOST=127.0.0.1
REDIS_PORT=6379

JWT_SECRET=your-jwt-secret
JWT_EXPIRE_TIME=120
```

### 3. 初始化数据库

```bash
# 运行迁移
go run main.go migrate up

# 导入示例数据（可选）
go run main.go seed
```

### 4. 启动服务

```bash
go run main.go serve
```

服务将在 `http://localhost:3000` 启动

## 🧪 验证安装

```bash
# 查看 API 文档
curl http://localhost:3000/swagger

# 运行测试
go test ./...
```

## 📚 常见命令

```bash
# 启动服务
go run main.go serve

# 数据库迁移
go run main.go migrate up       # 执行迁移
go run main.go migrate refresh  # 重置并重新迁移

# 数据填充
go run main.go seed             # 导入所有数据
go run main.go seed UserSeeder  # 导入特定数据

# 代码生成
go run main.go make model User  # 生成新模型
```

## 🚨 常见问题

### Redis 连接超时
- Redis 非必需，服务会在超时后继续
- 若需要，请启动 Redis：`brew services start redis`

### 数据库连接失败
- 确认数据库已启动
- 检查 .env 中的数据库配置

### 端口被占用
- 修改 .env 中的 `APP_PORT`
- 或关闭占用该端口的服务

## 📖 下一步

- 阅读 [ARCHITECTURE.md](./ARCHITECTURE.md) 理解系统设计
- 查看 [DEVELOPMENT.md](./DEVELOPMENT.md) 了解开发规范
- 参考 [RBAC.md](./RBAC.md) 实现权限控制
