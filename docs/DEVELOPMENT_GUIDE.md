# 开发者指南

## 环境准备
- Go 1.25.5+，MySQL 8.0+，Redis 6.0+
- 复制 `.env.example` 为 `.env`，填写数据库、Redis、JWT/APP_KEY 及短信/邮件密钥
- 确保 `config/*.go` 中的端口/凭据与本地环境一致

## 迁移与种子
- 迁移：`go run main.go migrate up`（或 `down`/`refresh`）
- 种子：`go run main.go seed`
  - 用户 30 / 分类 6 / 话题 40 / 评论 145（含回复）
  - 运行顺序含占位 `SeedInteractions`、`SeedMessages`，如未实现请先从 `database/seeders/seeder.go` 移除或补齐对应 seeder

## 本地运行
- 开发启动：`go run main.go serve`（默认 http://localhost:3000）
- Swagger：`http://localhost:3000/swagger/index.html`
- CLI 帮助：`go run main.go make --help`

## 编码规范
- 分层：Controller -> Service -> Repository/Cache；不要在 Controller 中操作模型
- DTO：请求/响应 DTO 放在 services 层，控制器只做绑定与校验
- 事务：在 Service 层使用 `database.DB.Transaction` 封装多步写操作
- 日志：使用 `pkg/logger`，错误日志请带上下文（RequestID、业务关键字段）

## 测试与质量
- 单元测试：`go test ./...`
- 针对 Service/Repository 编写 testify 断言；模拟外部依赖时使用接口注入或 fake 实现
- 建议为核心业务添加表驱动测试，覆盖成功/失败/边界场景

## 常见问题
- 连接失败：检查 `.env` 中的 DB/Redis 配置与容器/本地端口
- JWT 相关报错：确认 `APP_KEY` 与 `JWT_SECRET` 已设置且非默认值
- 种子报错：确认迁移已执行；若提示缺少 seeder，按需移除运行顺序或补充实现
