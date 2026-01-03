package config

import "GoHub-Service/pkg/config"

func init() {
	config.Add("storage", func() map[string]interface{} {
		return map[string]interface{}{
			// 文件存储基础目录（磁盘路径）
			"base_path": config.Env("STORAGE_BASE_PATH", "public/uploads"),

			// 对外暴露的访问前缀（URL Path）
			"public_prefix": config.Env("STORAGE_PUBLIC_PREFIX", "/uploads"),

			// 注意：文件大小限制、图片尺寸等业务常量已移至 config/app_constants.go
			// 请使用 appconfig.GetMaxImageSizeMB() 等方法访问
		}
	})
}
