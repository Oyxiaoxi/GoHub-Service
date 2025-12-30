package requests

import (
	"fmt"
	"mime/multipart"

	"GoHub-Service/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// TopicImageUploadRequest 话题配图上传请求
type TopicImageUploadRequest struct {
	Image *multipart.FileHeader `valid:"image" form:"image"`
}

// TopicImageUpload 话题配图上传验证
func TopicImageUpload(data interface{}, c *gin.Context) map[string][]string {
	maxSizeMB := config.GetInt64("storage.max_size_mb", 5)
	sizeRule := fmt.Sprintf("size:%d", maxSizeMB*1024*1024)

	rules := govalidator.MapData{
		"file:image": []string{"ext:png,jpg,jpeg,gif,webp", sizeRule, "required"},
	}
	messages := govalidator.MapData{
		"file:image": []string{
			"ext:图片仅支持 png/jpg/jpeg/gif/webp",
			"size:图片最大不能超过 20MB",
			"required:请上传图片文件",
		},
	}

	return validateFile(c, data, rules, messages)
}
