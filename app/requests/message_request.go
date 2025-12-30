package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// MessageSendRequest 发送私信请求
type MessageSendRequest struct {
	ReceiverID string `json:"receiver_id" valid:"receiver_id"`
	Body       string `json:"body" valid:"body"`
}

// MessageConversationRequest 获取会话请求
type MessageConversationRequest struct {
	UserID string `form:"user_id" valid:"user_id"`
}

// MessageMarkReadRequest 标记已读请求
type MessageMarkReadRequest struct {
	UserID string `json:"user_id" valid:"user_id"`
}

// MessageSend 验证发送私信
func MessageSend(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"receiver_id": []string{"required", "exists:users,id"},
		"body":        []string{"required", "min_cn:1", "max_cn:500"},
	}
	messages := govalidator.MapData{
		"receiver_id": []string{
			"required:收件人必填",
			"exists:收件人不存在",
		},
		"body": []string{
			"required:消息内容不能为空",
			"min_cn:消息内容不能为空",
			"max_cn:消息内容不能超过500字",
		},
	}
	return validate(data, rules, messages)
}

// MessageConversation 验证会话查询
func MessageConversation(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"user_id": []string{"required", "exists:users,id"},
	}
	messages := govalidator.MapData{
		"user_id": []string{
			"required:会话用户ID必填",
			"exists:用户不存在",
		},
	}
	return validate(data, rules, messages)
}

// MessageMarkRead 验证标记已读
func MessageMarkRead(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"user_id": []string{"required", "exists:users,id"},
	}
	messages := govalidator.MapData{
		"user_id": []string{
			"required:会话用户ID必填",
			"exists:用户不存在",
		},
	}
	return validate(data, rules, messages)
}
