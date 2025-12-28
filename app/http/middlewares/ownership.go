// Package middlewares 授权检查中间件
package middlewares

import (
	"GoHub-Service/pkg/auth"
	"GoHub-Service/pkg/response"

	"github.com/gin-gonic/gin"
)

// CheckOwnership 检查资源所有权中间件
// 使用方法：在路由中添加 CheckOwnership("user_id") 或 CheckOwnership("UserID")
func CheckOwnership(ownerFieldGetter func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserID := auth.CurrentUID(c)
		if currentUserID == "" {
			response.Abort403(c, "未登录")
			return
		}

		ownerID := ownerFieldGetter(c)
		if ownerID == "" {
			response.Abort404(c, "资源不存在")
			return
		}

		if currentUserID != ownerID {
			response.Abort403(c, "无权限操作")
			return
		}

		c.Next()
	}
}

// OwnershipChecker 所有权检查接口
// 实现此接口的模型可以使用通用的所有权检查
type OwnershipChecker interface {
	GetOwnerID() string
}

// CheckModelOwnership 检查模型所有权的通用函数
// 在Controller中使用，用于检查当前用户是否是资源的所有者
func CheckModelOwnership(c *gin.Context, model OwnershipChecker) bool {
	currentUserID := auth.CurrentUID(c)
	if currentUserID == "" {
		response.Abort403(c, "未登录")
		return false
	}

	ownerID := model.GetOwnerID()
	if ownerID == "" {
		response.Abort404(c, "资源不存在")
		return false
	}

	if currentUserID != ownerID {
		response.Abort403(c, "无权限操作")
		return false
	}

	return true
}

// AuthPolicy 授权策略函数类型
type AuthPolicy func(*gin.Context, interface{}) bool

// CheckPolicy 通用策略检查中间件
// 接受一个策略函数和模型获取函数
func CheckPolicy(policy AuthPolicy, modelGetter func(*gin.Context) interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		model := modelGetter(c)
		if model == nil {
			response.Abort404(c)
			return
		}

		if !policy(c, model) {
			response.Abort403(c)
			return
		}

		c.Next()
	}
}
