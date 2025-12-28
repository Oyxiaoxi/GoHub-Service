# GoHub-Service

基于 Gin 框架的 Go Web 服务项目

## 功能特性

- ✅ 用户注册验证（手机号/邮箱）
- ✅ 图片验证码
- ✅ 短信验证码（阿里云SMS）
- ✅ Redis 缓存
- ✅ 统一响应处理
- ✅ 日志系统（Zap + Lumberjack）
- ✅ 请求验证中间件
- ✅ 异常恢复中间件

## 环境要求

- Go 1.25.5+
- MySQL 8.0+
- Redis 6.0+

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd GoHub-Service
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境变量

复制配置文件并修改：

```bash
cp .env.example .env
```

编辑 `.env` 文件，配置数据库和其他服务：

```env
# 数据库配置
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=GoHub-Service
DB_USERNAME=root
DB_PASSWORD=your_password

# Redis 配置
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
```

### 4. 启动服务

```bash
# 启动 Redis（如果未启动）
redis-server

# 启动应用
go run main.go
```

服务将在 `http://localhost:3000` 启动

## 开发测试

### 测试验证码功能

由于阿里云短信服务需要账户余额，在开发环境中可以使用测试手机号：

**方法一：使用测试手机号前缀**

配置文件中 `verifycode.debug_phone_prefix` 设置为 `000`，使用以 `000` 开头的手机号：

```json
POST /v1/auth/verify-codes/phone
{
    "phone": "00012345678",
    "captcha_id": "captcha_skip_test",
    "captcha_answer": "123456"
}
```

这样会跳过实际的短信发送，直接使用配置的 `debug_code`（默认 123456）

**方法二：配置阿里云短信服务**

1. 登录[阿里云短信服务控制台](https://dysms.console.aliyun.com/)
2. 获取 AccessKey ID 和 AccessKey Secret
3. 配置短信签名和模板
4. 充值账户余额
5. 在 `.env` 中配置：

```env
SMS_ALIYUN_ACCESS_ID=your_access_key_id
SMS_ALIYUN_ACCESS_SECRET=your_access_key_secret
SMS_ALIYUN_SIGN_NAME=your_sign_name
SMS_ALIYUN_TEMPLATE_CODE=SMS_xxxxxx
```

### API 端点

#### 认证相关

- `POST /v1/auth/signup/phone/exist` - 检查手机号是否已注册
- `POST /v1/auth/signup/email/exist` - 检查邮箱是否已注册
- `POST /v1/auth/verify-codes/captcha` - 获取图片验证码
- `POST /v1/auth/verify-codes/phone` - 发送手机验证码

#### 图片验证码示例

```bash
# 获取图片验证码
curl -X POST http://localhost:3000/v1/auth/verify-codes/captcha

# 响应
{
    "captcha_id": "6MRHuD6MVFgO9cxF4nYQ",
    "captcha_image": "data:image/png;base64,..."
}
```

#### 手机验证码示例

```bash
# 发送验证码（生产环境）
curl -X POST http://localhost:3000/v1/auth/verify-codes/phone \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "18800000000",
    "captcha_id": "6MRHuD6MVFgO9cxF4nYQ",
    "captcha_answer": "673053"
  }'

# 发送验证码（开发测试 - 使用测试手机号）
curl -X POST http://localhost:3000/v1/auth/verify-codes/phone \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "00012345678",
    "captcha_id": "captcha_skip_test",
    "captcha_answer": "123456"
  }'
```

## 项目结构

```
GoHub-Service/
├── app/
│   ├── http/
│   │   ├── controllers/    # 控制器
│   │   └── middlewares/    # 中间件
│   ├── models/             # 数据模型
│   └── requests/           # 请求验证
├── bootstrap/              # 初始化模块
├── config/                 # 配置文件
├── pkg/                    # 公共包
│   ├── app/               # 应用帮助函数
│   ├── captcha/           # 图片验证码
│   ├── config/            # 配置管理
│   ├── database/          # 数据库
│   ├── helpers/           # 工具函数
│   ├── logger/            # 日志
│   ├── redis/             # Redis 客户端
│   ├── response/          # 响应处理
│   ├── sms/               # 短信服务
│   └── verifycode/        # 验证码
├── routes/                # 路由
└── storage/               # 存储（日志等）
```

## 技术栈

- **Web框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **日志**: [Zap](https://github.com/uber-go/zap)
- **配置**: [Viper](https://github.com/spf13/viper)
- **验证**: [govalidator](https://github.com/thedevsaddam/govalidator)
- **Redis**: [go-redis](https://github.com/go-redis/redis)
- **短信**: [阿里云SDK](https://github.com/aliyun/alibaba-cloud-sdk-go)

## 常见问题

### 1. 短信发送失败：账户余额不足

**问题**: 日志显示 `isv.AMOUNT_NOT_ENOUGH - 账户余额不足`

**解决方案**:
- 开发测试：使用测试手机号 `000xxxxxxxx`
- 生产环境：在阿里云控制台充值

### 2. Redis 连接失败

**问题**: `dial tcp 127.0.0.1:6379: connect: connection refused`

**解决方案**:
```bash
# macOS
brew services start redis

# Linux
sudo systemctl start redis

# 或直接运行
redis-server
```

### 3. 数据库连接失败

**问题**: 无法连接到 MySQL

**解决方案**:
- 检查 MySQL 服务是否启动
- 确认 `.env` 中的数据库配置正确
- 确保数据库已创建：`CREATE DATABASE GoHub-Service;`

## License

MIT License
