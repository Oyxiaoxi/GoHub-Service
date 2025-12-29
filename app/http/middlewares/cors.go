package middlewares

import (
	"strings"
	"time"

	"GoHub-Service/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 配置跨域资源共享中间件
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许的源列表（从配置读取，逗号分隔）
		AllowOrigins: getAllowedOrigins(),

		// 允许的 HTTP 方法
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		// 允许的请求头
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"X-CSRF-Token",
			"X-Requested-With",
			"X-Request-ID",
		},

		// 暴露给客户端的响应头
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
		},

		// 是否允许携带 Cookie
		AllowCredentials: true,

		// 预检请求的缓存时间
		MaxAge: 12 * time.Hour,
	})
}

// CORSPublic 公开 API 的 CORS 配置（允许所有源）
// 仅用于完全公开的只读 API
func CORSPublic() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{
			"GET",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		MaxAge: 12 * time.Hour,
	})
}

// CORSWithOrigins 自定义源的 CORS 配置
// 用于需要特定源配置的场景
func CORSWithOrigins(origins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// getAllowedOrigins 从配置读取 CORS 允许的域名，使用逗号分隔
func getAllowedOrigins() []string {
	raw := config.GetString("cors.allowed_origins")
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
