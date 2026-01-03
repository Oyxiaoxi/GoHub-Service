package routes

import (
	v1 "GoHub-Service/app/http/controllers/api/v1"
	"GoHub-Service/app/http/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterMessageRoutes 注册私信相关路由
func RegisterMessageRoutes(r *gin.RouterGroup) {
	controller := v1.NewMessagesController()

	group := r.Group("/messages", middlewares.AuthJWT())
	{
		// 私信发送应用签名验证（防止批量发送垃圾私信）
		group.POST("", 
			middlewares.APISignatureVerification(),
			controller.Send,
		)
		group.GET("", controller.Conversation)
		group.POST("/read", controller.MarkRead)
		group.GET("/unread-count", controller.UnreadCount)
	}
}
