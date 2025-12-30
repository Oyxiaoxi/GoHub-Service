package config

import "GoHub-Service/pkg/config"

func init() {
	config.Add("storage", func() map[string]interface{} {
		return map[string]interface{}{
			// 文件存储基础目录（磁盘路径）
			"base_path": config.Env("STORAGE_BASE_PATH", "public/uploads"),

			// 对外暴露的访问前缀（URL Path）
			"public_prefix": config.Env("STORAGE_PUBLIC_PREFIX", "/uploads"),

			// 允许的最大文件尺寸（单位：MB）
			"max_size_mb": config.Env("STORAGE_MAX_SIZE_MB", 2),

			// 图片最大宽度（超过则等比缩放）
			"max_image_width": config.Env("STORAGE_MAX_IMAGE_WIDTH", 1600),

			// JPEG 编码质量 (1-100)
			"jpeg_quality": config.Env("STORAGE_JPEG_QUALITY", 85),

			// 允许的文件扩展名，逗号分隔
			"allowed_ext": config.Env("STORAGE_ALLOWED_EXT", "jpg,jpeg,png,gif,webp"),
		}
	})
}
