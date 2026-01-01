package config

import "GoHub-Service/pkg/config"

func init() {
	config.Add("elasticsearch", func() map[string]interface{} {
		return map[string]interface{}{
			// 是否启用Elasticsearch
			"enabled": config.Env("ELASTICSEARCH_ENABLED", false),

			// Elasticsearch地址
			"addresses": []string{
				config.Env("ELASTICSEARCH_HOST", "http://localhost:9200").(string),
			},

			// 认证信息
			"username": config.Env("ELASTICSEARCH_USERNAME", ""),
			"password": config.Env("ELASTICSEARCH_PASSWORD", ""),

			// 连接超时时间(毫秒)
			"timeout": config.Env("ELASTICSEARCH_TIMEOUT", 30000),

			// 索引配置
			"index_name":       config.Env("ELASTICSEARCH_INDEX_NAME", "gohub-topics"),
			"shards":           config.Env("ELASTICSEARCH_SHARDS", 3),
			"replicas":         config.Env("ELASTICSEARCH_REPLICAS", 1),
			"refresh_interval": config.Env("ELASTICSEARCH_REFRESH_INTERVAL", "1s"),

			// 性能配置
			"max_retries":      config.Env("ELASTICSEARCH_MAX_RETRIES", 3),
			"bulk_size":        config.Env("ELASTICSEARCH_BULK_SIZE", 1000),
			"flush_interval":   config.Env("ELASTICSEARCH_FLUSH_INTERVAL", 5000),
			"concurrent_bulks": config.Env("ELASTICSEARCH_CONCURRENT_BULKS", 5),
		}
	})
}
