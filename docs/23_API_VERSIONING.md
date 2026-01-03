# API 版本管理策略

## 概述

GoHub-Service 采用语义化版本管理策略，确保 API 的稳定性和向后兼容性。

## 版本命名规则

### URL 路径版本
- 格式：`/api/v{major}`
- 示例：`/api/v1`, `/api/v2`
- 当前版本：**v1**

### 版本状态
- **active**: 活跃版本，完全支持
- **deprecated**: 已废弃，计划停用
- **sunset**: 已停用，不再支持
- **planned**: 计划中的未来版本

## 版本生命周期

### 1. 发布新版本
```
v1 (active) -> v2 (active)
```
- 新版本发布时，旧版本保持活跃
- 支持多个版本并存

### 2. 废弃旧版本
```
v1 (active) -> v1 (deprecated)
```
- 提前 6 个月通知废弃
- 响应头包含废弃警告：`X-API-Warn: API version v1 is deprecated`
- 文档标注废弃信息

### 3. 停用旧版本
```
v1 (deprecated) -> v1 (sunset)
```
- 废弃 6 个月后正式停用
- 返回 `410 Gone` 状态码
- 提供迁移指南

## 版本管理中间件

### 使用示例

```go
import "GoHub-Service/pkg/apiversion"

// 在路由组中启用版本管理
v1 := r.Group("/api/v1")
v1.Use(apiversion.VersionDeprecated("v1"))
```

### 版本信息端点

获取所有支持的 API 版本：
```bash
GET /api/versions
```

响应示例：
```json
{
  "current_version": "v1",
  "versions": {
    "v1": {
      "version": "v1",
      "status": "active",
      "release_date": "2024-01-01",
      "features": [
        "用户管理",
        "话题管理",
        "评论管理"
      ]
    },
    "v2": {
      "version": "v2",
      "status": "planned",
      "release_date": "2026-06-01",
      "features": [
        "GraphQL 支持",
        "Websocket 实时通知"
      ]
    }
  },
  "api_docs": "/swagger/index.html"
}
```

## 版本请求方式

### 方式1：URL 路径（推荐）
```
GET /api/v1/users
```

### 方式2：请求头
```
GET /api/users
X-API-Version: v1
```

## 破坏性变更指南

### 需要发布新版本的情况
- 删除或重命名 API 端点
- 修改请求/响应数据结构
- 改变认证/授权机制
- 修改错误码定义

### 不需要发布新版本的情况
- 添加新的 API 端点
- 添加可选的请求参数
- 添加响应字段
- 修复 bug
- 性能优化

## 版本迁移清单

### 从 v1 迁移到 v2（示例）

1. **更新基础 URL**
   ```
   https://api.gohub.com/api/v1 -> https://api.gohub.com/api/v2
   ```

2. **检查废弃端点**
   - 查看 v2 文档中的变更列表
   - 更新已废弃的端点

3. **更新数据结构**
   - 检查响应格式变化
   - 更新数据模型

4. **测试验证**
   - 在测试环境验证所有接口
   - 确认业务逻辑正常

## 最佳实践

### 对于 API 提供者
1. **向后兼容优先**：尽量保持向后兼容
2. **提前通知**：废弃前 6 个月通知
3. **完善文档**：提供详细的迁移指南
4. **监控使用**：跟踪旧版本使用情况
5. **延长支持**：根据使用情况延长支持期

### 对于 API 消费者
1. **及时升级**：关注版本公告，及时升级
2. **测试覆盖**：确保测试覆盖所有 API 调用
3. **错误处理**：正确处理版本相关错误
4. **监控告警**：监控废弃警告
5. **缓存清理**：升级后清理相关缓存

## 版本支持时间表

| 版本 | 发布日期 | 废弃日期 | 停用日期 | 当前状态 |
|------|----------|----------|----------|----------|
| v1   | 2024-01-01 | -        | -        | Active   |
| v2   | 2026-06-01 | -        | -        | Planned  |

## 相关资源

- [API 文档](/swagger/index.html)
- [变更日志](../CHANGELOG.md)
- [迁移指南](./24_MIGRATION_GUIDE.md)
