package config

import "GoHub-Service/pkg/config"

func init() {
	config.Add("cors", func() map[string]interface{} {
		return map[string]interface{}{
			// 允许的跨域源（逗号分隔）。默认允许本地常用端口。
			"allowed_origins": config.Env("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001,http://localhost:3000,http://127.0.0.1:3000,http://127.0.0.1:3001,http://127.0.0.1:8080"),
		}
	})
}
