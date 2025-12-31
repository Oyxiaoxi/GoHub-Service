// Package response 统一错误码定义
package response

// 业务错误码定义
const (
	// 通用错误码 (1xxx)
	CodeSuccess           = 0
	CodeError             = 1000
	CodeInvalidParams     = 1001
	CodeUnauthorized      = 1002
	CodeForbidden         = 1003
	CodeNotFound          = 1004
	CodeInternalError     = 1005
	CodeValidationError   = 1006
	CodeTooManyRequests   = 1007

	// 用户相关错误码 (2xxx)
	CodeUserNotFound      = 2001
	CodeUserExists        = 2002
	CodePasswordError     = 2003
	CodeUserDisabled      = 2004
	CodePhoneExists       = 2005
	CodeEmailExists       = 2006

	// 认证相关错误码 (3xxx)
	CodeTokenInvalid      = 3001
	CodeTokenExpired      = 3002
	CodeTokenMissing      = 3003
	CodeCaptchaError      = 3004
	CodeVerifyCodeError   = 3005

	// 资源相关错误码 (4xxx)
	CodeTopicNotFound     = 4001
	CodeCategoryNotFound  = 4002
	CodeLinkNotFound      = 4003
	
	// 数据库相关错误码 (5xxx)
	CodeDatabaseError     = 5001
	CodeCreateFailed      = 5002
	CodeUpdateFailed      = 5003
	CodeDeleteFailed      = 5004
	
	// 服务器相关错误码 (6xxx)
	CodeServerError       = 6001
)

// 错误码对应的默认消息
var codeMessages = map[int]string{
	CodeSuccess:           "操作成功",
	CodeError:             "操作失败",
	CodeInvalidParams:     "请求参数错误",
	CodeUnauthorized:      "未授权，请先登录",
	CodeForbidden:         "权限不足",
	CodeNotFound:          "资源不存在",
	CodeInternalError:     "服务器内部错误",
	CodeValidationError:   "请求验证失败",
	CodeTooManyRequests:   "请求过于频繁，请稍后再试",

	CodeUserNotFound:      "用户不存在",
	CodeUserExists:        "用户已存在",
	CodePasswordError:     "密码错误",
	CodeUserDisabled:      "用户已被禁用",
	CodePhoneExists:       "手机号已被注册",
	CodeEmailExists:       "邮箱已被注册",

	CodeTokenInvalid:      "Token无效",
	CodeTokenExpired:      "Token已过期",
	CodeTokenMissing:      "Token缺失",
	CodeCaptchaError:      "图形验证码错误",
	CodeVerifyCodeError:   "验证码错误",

	CodeTopicNotFound:     "话题不存在",
	CodeCategoryNotFound:  "分类不存在",
	CodeLinkNotFound:      "链接不存在",

	CodeDatabaseError:     "数据库操作失败",
	CodeCreateFailed:      "创建失败",
	CodeUpdateFailed:      "更新失败",
	CodeDeleteFailed:      "删除失败",
	CodeServerError:       "服务器内部错误",
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return codeMessages[CodeError]
}
