package config

import (
	"GoHub-Service/pkg/config"
)

func init() {

	config.Add("database", func() map[string]interface{} {
		return map[string]interface{}{

			// 默认数据库
			"connection": config.Env("DB_CONNECTION", "mysql"),

			"mysql": map[string]interface{}{

				// 数据库连接信息
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", "3306"),
				"database": config.Env("DB_DATABASE", "GoHub-Service"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"charset":  "utf8mb4",

				// 连接池配置
				// max_idle_connections: 空闲连接池大小，建议设置为 max_open_connections 的 1/4 到 1/2
				"max_idle_connections": config.Env("DB_MAX_IDLE_CONNECTIONS", 25),
				// max_open_connections: 最大打开连接数，根据业务负载调整，建议 100-200
				"max_open_connections": config.Env("DB_MAX_OPEN_CONNECTIONS", 100),
				// max_life_seconds: 连接最大生命周期（秒），建议 5-30 分钟
				"max_life_seconds": config.Env("DB_MAX_LIFE_SECONDS", 10*60),
			},

			"sqlite": map[string]interface{}{
				"database": config.Env("DB_SQL_FILE", "database/database.db"),
			},
		}
	})
}
