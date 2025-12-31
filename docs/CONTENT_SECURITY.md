# 内容安全使用指南

## 功能概述

GoHub-Service 实现了完整的内容安全增强功能，包括：

1. **敏感词过滤** - 自动检测和过滤政治、色情、暴力等敏感词
2. **XSS 防护** - 防止跨站脚本攻击
3. **内容审核** - 自动检查和清理用户输入
4. **图片上传安全** - 验证图片格式和大小

## 配置说明

在 `.env` 文件中添加以下配置：

```env
# 内容安全配置
CONTENT_CHECK_ENABLED=true
SENSITIVE_WORD_FILTER_ENABLED=true
SENSITIVE_WORD_REPLACEMENT=***
XSS_PROTECTION_ENABLED=true
MAX_CONTENT_LENGTH=10000
MAX_TITLE_LENGTH=200
ALLOW_HTML_TAGS=false
IMAGE_CHECK_ENABLED=true
MAX_IMAGE_SIZE=5242880
AUDIT_LOG_ENABLED=true
```

## 使用方式

### 1. 中间件集成

项目已在以下路由中自动集成内容安全中间件：

**话题相关（/topics）：**
- POST /topics - 创建话题（敏感词过滤 + XSS防护）
- PUT /topics/:id - 更新话题（敏感词过滤 + XSS防护）
- POST /topics/upload-image - 上传图片（图片安全检查）

**评论相关（/comments）：**
- POST /comments - 创建评论（敏感词过滤 + XSS防护）
- PUT /comments/:id - 更新评论（敏感词过滤 + XSS防护）

**用户相关（/users）：**
- PUT /users - 更新资料（敏感词过滤 + XSS防护）
- PUT /users/avatar - 上传头像（图片安全检查）

### 2. 中间件说明

#### SensitiveWordFilter - 敏感词过滤中间件
自动过滤文本中的敏感词，将其替换为配置的替换字符（默认 `***`）

```go
middlewares.SensitiveWordFilter()
```

#### XSSProtection - XSS防护中间件
清理 HTML 标签和危险脚本，防止 XSS 攻击

```go
middlewares.XSSProtection()
```

#### ImageUploadSecurity - 图片上传安全中间件
检查上传图片的格式和大小

```go
middlewares.ImageUploadSecurity()
```

#### ContentSecurity - 综合内容安全中间件
包含所有安全检查，不合格内容会被拦截并返回 403

```go
middlewares.ContentSecurity()
```

### 3. 直接使用内容检查器

在控制器或服务中直接使用：

```go
import "GoHub-Service/pkg/security"

func (ctrl *Controller) SomeMethod(c *gin.Context) {
    checker := security.GetContentChecker()
    
    // 检查标题
    titleResult := checker.CheckTitle(title)
    if !titleResult.IsValid {
        response.Abort403(c, titleResult.Message)
        return
    }
    
    // 检查内容
    contentResult := checker.CheckContent(content)
    if !contentResult.IsValid {
        response.Abort403(c, contentResult.Message)
        return
    }
    
    // 使用过滤后的内容
    filteredTitle := titleResult.FilteredText
    filteredContent := contentResult.FilteredText
    
    // 记录发现的敏感词
    if len(contentResult.FoundWords) > 0 {
        logger.Warn("发现敏感词", contentResult.FoundWords)
    }
}
```

### 4. 敏感词管理

```go
import "GoHub-Service/pkg/security"

// 获取敏感词过滤器
filter := security.GetFilter()

// 添加自定义敏感词
filter.AddWords("新敏感词1", "新敏感词2")

// 检查是否包含敏感词
if filter.Contains("测试文本") {
    // 包含敏感词
}

// 查找所有敏感词
words := filter.FindAll("测试文本")

// 过滤敏感词
filtered := filter.Filter("测试文本")

// 删除敏感词
filter.RemoveWords("敏感词1")

// 清空所有敏感词
filter.Clear()

// 获取所有敏感词列表
allWords := filter.GetWords()
```

### 5. XSS 防护使用

```go
import "GoHub-Service/pkg/security"

checker := security.GetContentChecker()

// HTML 转义
escaped := checker.EscapeHTML("<script>alert('xss')</script>")

// 清理 HTML
cleaned := checker.CleanHTML("<script>alert('xss')</script><p>正常内容</p>")

// 验证输入
if !checker.xssFilter.ValidateInput(userInput) {
    // 包含危险内容
}
```

## 技术实现

### 敏感词过滤算法

使用 **前缀树（Trie）** 算法实现高效的敏感词匹配：

- 时间复杂度：O(n)，n 为文本长度
- 空间复杂度：O(m)，m 为敏感词总字符数
- 支持最长匹配策略
- 线程安全（使用 RWMutex）

### XSS 防护策略

- 移除 `<script>`、`<iframe>`、`<style>` 等危险标签
- 过滤 `javascript:`、`vbscript:`、`data:` 等危险协议
- 清除所有 `on*` 事件处理器（onclick、onerror 等）
- HTML 实体转义
- 可配置的标签和属性白名单

### 默认敏感词库

包含以下分类的敏感词：
- 政治敏感词
- 色情淫秽词汇
- 暴力相关词汇
- 赌博欺诈词汇
- 毒品相关词汇

可通过 `AddWords()` 方法扩展自定义敏感词。

## 性能考虑

1. **敏感词过滤器**采用单例模式，全局共享一个实例
2. **前缀树**在启动时构建，运行时只读，性能优秀
3. **读写锁**保证并发安全的同时，允许多个 goroutine 并发读取
4. 中间件只处理 POST/PUT/PATCH 请求，不影响 GET/DELETE 性能

## 安全建议

1. **定期更新敏感词库** - 根据实际情况添加新的敏感词
2. **启用审计日志** - 记录所有触发安全规则的操作
3. **调整过滤策略** - 根据业务需求调整是否拦截还是仅过滤
4. **监控告警** - 对频繁触发安全规则的用户进行监控
5. **人工审核** - 对敏感内容进行二次人工审核

## 常见问题

**Q: 敏感词过滤会影响正常词汇吗？**
A: 使用最长匹配策略，尽量避免误伤。如发现误判，可通过 RemoveWords() 删除。

**Q: 如何自定义敏感词？**
A: 在应用启动时调用 `filter.AddWords()` 添加自定义敏感词。

**Q: XSS 防护会影响富文本编辑吗？**
A: 可以通过配置 `ALLOW_HTML_TAGS=true` 并设置白名单来允许特定的 HTML 标签。

**Q: 如何禁用某些路由的内容安全检查？**
A: 不在路由上添加对应的中间件即可，或设置配置项 `CONTENT_CHECK_ENABLED=false` 全局禁用。

## 未来扩展

- [ ] 接入第三方内容审核 API（如阿里云、腾讯云）
- [ ] 增加图片内容识别（OCR + 图像审核）
- [ ] 支持正则表达式规则
- [ ] 管理后台敏感词管理界面
- [ ] 内容审核日志查询功能
- [ ] 机器学习模型识别恶意内容
