package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type CommentRequest struct {
	TopicID  string `json:"topic_id,omitempty" valid:"topic_id"`
	Content  string `json:"content,omitempty" valid:"content"`
	ParentID string `json:"parent_id,omitempty" valid:"parent_id"`
}

func CommentSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"topic_id": []string{"required", "exists:topics,id"},
		"content":  []string{"required", "min_cn:1", "max_cn:1000"},
		"parent_id": []string{"numeric"},
	}
	messages := govalidator.MapData{
		"topic_id": []string{
			"required:话题ID为必填项",
			"exists:话题不存在",
		},
		"content": []string{
			"required:评论内容为必填项",
			"min_cn:评论内容不能为空",
			"max_cn:评论内容不能超过1000字",
		},
		"parent_id": []string{
			"numeric:父评论ID必须为数字",
		},
	}
	return validate(data, rules, messages)
}

func CommentUpdate(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"content": []string{"min_cn:1", "max_cn:1000"},
	}
	messages := govalidator.MapData{
		"content": []string{
			"min_cn:评论内容不能为空",
			"max_cn:评论内容不能超过1000字",
		},
	}
	return validate(data, rules, messages)
}
