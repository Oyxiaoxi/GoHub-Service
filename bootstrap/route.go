// Package bootstrap 处理程序初始化逻辑
package bootstrap

import (
    "GoHub-Service/app/http/middlewares"
    "GoHub-Service/pkg/metrics"
    "GoHub-Service/routes"
    "net/http"
    "strings"

    "github.com/gin-contrib/gzip"
    "github.com/gin-gonic/gin"
)

// SetupRoute 路由初始化
func SetupRoute(router *gin.Engine) {

    // 设置资源追踪器（必须在注册中间件之前）
    middlewares.SetResourceTracker(Tracker)

    // 注册全局中间件
    registerGlobalMiddleWare(router)

    //  注册 API 路由
    routes.RegisterAPIRoutes(router)

    //  配置 404 路由
    setup404Handler(router)
}

func registerGlobalMiddleWare(router *gin.Engine) {
    // 推荐注册顺序：CORS → SecureHeaders → EnhancedSecurityValidation → EnhancedSQLInjectionProtection → XSSProtection → RequestID → ResourceTracking → Prometheus → Recovery → Logger → ForceUA → Gzip
    router.Use(
        middlewares.CORS(),                          // 1. CORS 跨域配置，最先处理
        middlewares.SecureHeaders(),                 // 2. 安全响应头
        middlewares.EnhancedSecurityValidation(),    // 3. 增强安全验证（SQL/XSS/路径遍历）
        middlewares.EnhancedSQLInjectionProtection(), // 4. 增强 SQL 注入防护
        middlewares.EnhancedXSSProtection(),         // 5. 增强 XSS 防护
        middlewares.XSSProtection(),                 // 6. 基础 XSS 防护（向后兼容）
        middlewares.RequestID(),                     // 7. 请求唯一ID
        middlewares.ResourceTracking(),              // 8. 资源追踪（在 Prometheus 之前，更准确）
        metrics.PrometheusMiddleware(),              // 9. Prometheus 指标收集
        middlewares.Recovery(),                      // 10. Panic 恢复，保证后续中间件能捕获异常
        middlewares.Logger(),                        // 11. 日志记录，捕获所有请求日志
        middlewares.ForceUA(),                       // 12. 强制 User-Agent，可根据需要调整顺序
        gzip.Gzip(gzip.DefaultCompression),          // 13. 启用 Gzip 压缩
    )
}

func setup404Handler(router *gin.Engine) {
    // 处理 404 请求
    router.NoRoute(func(c *gin.Context) {
        // 获取标头信息的 Accept 信息
        acceptString := c.Request.Header.Get("Accept")
        if strings.Contains(acceptString, "text/html") {
            // 如果是 HTML 的话
            c.String(http.StatusNotFound, "页面返回 404")
        } else {
            // 默认返回 JSON
            c.JSON(http.StatusNotFound, gin.H{
                "error_code":    404,
                "error_message": "路由未定义，请确认 url 和请求方法是否正确。",
            })
        }
    })
}


