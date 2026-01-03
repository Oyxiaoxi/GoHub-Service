# API 签名验证应用说明

## ✅ 已应用的路由

### 1. 用户敏感信息修改
**路径**: `/api/v1/users/*`  
**中间件**: `APISignatureVerification()` - 强制签名验证

- `PUT /api/v1/users/email` - 修改邮箱 ✅
- `PUT /api/v1/users/phone` - 修改手机号 ✅  
- `PUT /api/v1/users/password` - 修改密码 ✅

**防护效果**:
- 防止重放攻击（同一个请求不能被重复使用）
- 防止中间人篡改数据
- 5分钟时间窗口，过期自动失效

---

### 2. 密码重置接口
**路径**: `/api/v1/auth/password-reset/*`  
**中间件**: 
- `RateLimitMiddleware(5)` - 每分钟5次限制
- `APISignatureVerification()` - 强制签名验证

- `POST /api/v1/auth/password-reset/using-email` ✅
- `POST /api/v1/auth/password-reset/using-phone` ✅

**防护效果**:
- 严格限流（5次/分钟）
- 签名验证（防止暴力破解和重放）
- 结合 JWT + 签名双重验证

---

### 3. 用户登录接口
**路径**: `/api/v1/auth/login/*`  
**中间件**: `OptionalSignatureVerification()` - 可选签名验证

- `POST /api/v1/auth/login/using-phone` ✅
- `POST /api/v1/auth/login/using-password` ✅

**防护效果**:
- 兼容模式：旧客户端可正常登录，新客户端建议使用签名
- 提供了签名则强制验证
- 逐步迁移策略

---

### 4. 管理后台危险操作
**路径**: `/api/v1/admin/users/*`  
**中间件**: `APISignatureVerification()` - 强制签名验证

- `DELETE /api/v1/admin/users/:id` - 删除用户 ✅
- `POST /api/v1/admin/users/batch-delete` - 批量删除 ✅
- `POST /api/v1/admin/users/:id/ban` - 封禁用户 ✅
- `POST /api/v1/admin/users/:id/unban` - 解封用户 ✅

**防护效果**:
- 危险操作必须签名
- 防止管理员账号被盗后批量操作
- 审计追踪（Nonce 记录）

---

### 5. 私信发送
**路径**: `/api/v1/messages`  
**中间件**: `APISignatureVerification()` - 强制签名验证

- `POST /api/v1/messages` - 发送私信 ✅

**防护效果**:
- 防止批量发送垃圾私信
- 防止自动化脚本滥用
- 结合 JWT 身份验证

---

## 🔧 配置

### 环境变量
在 `.env` 文件中添加：

```bash
# API 签名密钥（强烈建议使用 32 位以上随机字符串）
SIGNATURE_SECRET=your-very-strong-secret-key-32chars-please-change-in-production
```

如果未设置，将使用 `APP_KEY` 作为后备密钥。

### 配置文件
`config/app.go` 已更新，新增配置项：

```go
"signature_secret": config.Env("SIGNATURE_SECRET", config.Env("APP_KEY", "...")),
```

---

## 📊 安全策略

### 三级防护体系

| 防护级别 | 应用场景 | 中间件组合 | 防护效果 |
|---------|---------|-----------|---------|
| **高危** | 密码重置、账号删除 | RateLimit(5) + Signature + JWT | 严格限流 + 签名验证 + 身份验证 |
| **敏感** | 修改邮箱/手机/密码 | Signature + JWT | 签名验证 + 身份验证 |
| **兼容** | 登录接口 | OptionalSignature + JWT | 可选签名 + 身份验证 |

### 防护效果对比

| 攻击类型 | 无签名 | 可选签名 | 强制签名 |
|---------|-------|---------|---------|
| 重放攻击 | ❌ 无防护 | ⚠️ 部分防护 | ✅ 完全防护 |
| 数据篡改 | ❌ 无防护 | ⚠️ 部分防护 | ✅ 完全防护 |
| 暴力破解 | ⚠️ 限流防护 | ✅ 限流+签名 | ✅ 限流+签名 |
| 自动化脚本 | ❌ 无防护 | ⚠️ 部分防护 | ✅ 完全防护 |

---

## 🎯 后续扩展建议

### 可选应用签名验证的接口

| 接口 | 路径 | 优先级 | 理由 |
|-----|------|-------|------|
| 转账/支付 | `/api/v1/payments/*` | 🔴 高 | 金额相关，防篡改 |
| 敏感数据查询 | `/api/v1/users/:id/sensitive` | 🟡 中 | 个人隐私数据 |
| 批量操作 | `/api/v1/topics/batch-*` | 🟡 中 | 防止滥用 |
| 文件上传 | `/api/v1/uploads` | 🟢 低 | 防止恶意上传 |

### 渐进式迁移策略

**阶段 1（当前）**: 关键接口强制签名
- ✅ 密码重置
- ✅ 账号信息修改
- ✅ 管理后台危险操作
- ✅ 私信发送

**阶段 2（1-2个月后）**: 扩展到更多接口
- 登录接口从 Optional 改为强制
- 支付相关接口
- 批量操作接口

**阶段 3（3-6个月后）**: 全局应用
- 所有写操作强制签名
- 敏感读操作可选签名

---

## 📝 客户端集成

### Go 客户端示例
```go
client := NewAPIClient("http://localhost:3000", "your-secret-key")

// 修改密码（需要签名）
body := `{"old_password":"123456","new_password":"newpass"}`
resp, err := client.Request("PUT", "/api/v1/users/password", body)
```

### JavaScript 客户端示例
```javascript
const client = new APIClient('http://localhost:3000', 'your-secret-key');

// 修改邮箱（需要签名）
await client.request('PUT', '/api/v1/users/email', {
  email: 'newemail@example.com',
  verify_code: '123456'
});
```

完整客户端代码请参考：[docs/28_API_SIGNATURE.md](28_API_SIGNATURE.md)

---

## 🔍 监控和审计

### Redis Nonce 记录
- **Key 格式**: `api:nonce:{nonce}`
- **有效期**: 5分钟
- **用途**: 防重放攻击检测

### 建议监控指标
```
# 签名验证失败次数
api_signature_failures_total{endpoint="/api/v1/users/password",reason="signature_mismatch"}

# 重放攻击检测次数
api_replay_attempts_total{endpoint="/api/v1/auth/password-reset/using-email"}

# 时间戳过期次数
api_timestamp_expired_total
```

---

## ⚠️ 注意事项

1. **密钥管理**
   - 生产环境务必修改 `SIGNATURE_SECRET`
   - 使用强随机字符串（32位以上）
   - 定期轮换密钥

2. **时间同步**
   - 客户端和服务器时间必须同步（NTP）
   - 时间差超过5分钟会导致验证失败

3. **Nonce 生成**
   - 必须使用加密安全的随机数生成器
   - 长度至少16位
   - 不要使用时间戳或自增ID

4. **HTTPS**
   - 生产环境必须使用 HTTPS
   - 防止签名在传输中被窃取

---

## 📚 相关文档

- [API 签名验证完整指南](28_API_SIGNATURE.md)
- [安全加固指南](27_SECURITY_HARDENING.md)
- [客户端示例代码](examples/api_signature_example.go)

---

**更新日期**: 2026年1月3日  
**版本**: v1.0  
**状态**: ✅ 已应用到生产代码
