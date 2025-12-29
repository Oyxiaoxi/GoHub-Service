package config

import "GoHub-Service/pkg/config"

func init() {
	config.Add("limiter", func() map[string]interface{} {
		return map[string]interface{}{
			"default_message": config.Env("LIMITER_DEFAULT_MESSAGE", "请求过于频繁，请稍后再试"),

			// 全局 IP 限流（API 根）
			"global_ip_rate": config.Env("LIMIT_GLOBAL_IP_RATE", "200-H"),
			// Auth 路由组 IP 限流
			"auth_ip_rate": config.Env("LIMIT_AUTH_IP_RATE", "1000-H"),

			// 注册/验证相关限流
			"signup_exist_rate":   config.Env("LIMIT_SIGNUP_EXIST_RATE", "60-H"),
			"verify_phone_rate":   config.Env("LIMIT_VERIFY_PHONE_RATE", "20-H"),
			"verify_email_rate":   config.Env("LIMIT_VERIFY_EMAIL_RATE", "20-H"),
			"verify_captcha_rate": config.Env("LIMIT_VERIFY_CAPTCHA_RATE", "50-H"),
		}
	})
}
