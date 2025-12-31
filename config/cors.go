package config

import "GoHub-Service/pkg/config"

func init() {
	config.Add("cors", func() map[string]interface{} {
		return map[string]interface{}{
			// 允许的跨域源（逗号分隔）
			// 开发环境：允许本地常用端口
			// 生产环境：只允许特定域名，如：https://yourdomain.com
			"allowed_origins": config.Env("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000,http://127.0.0.1:3001,http://127.0.0.1:8080"),

			// 允许的 HTTP 方法
			"allowed_methods": config.Env("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),

			// 允许的请求头
			"allowed_headers": config.Env("CORS_ALLOWED_HEADERS", "Authorization,Content-Type,Accept,Origin,User-Agent,DNT,Cache-Control,X-Requested-With"),

			// 暴露给客户端的响应头
			"exposed_headers": config.Env("CORS_EXPOSED_HEADERS", "Content-Length,Content-Type"),

			// 是否允许携带凭证（cookies、HTTP 认证）
			"allow_credentials": config.Env("CORS_ALLOW_CREDENTIALS", true),

			// 预检请求的有效期（秒），减少预检请求次数
			"max_age": config.Env("CORS_MAX_AGE", 12*3600), // 12小时
		}
	})
}
