// Package routes 通知相关路由
package routes

import (
	"GoHub-Service/app/http/controllers/api/v1"
	"GoHub-Service/app/http/middlewares"
	"github.com/gin-gonic/gin"
)

// RegisterNotificationRoutes 注册通知路由
func RegisterNotificationRoutes(rg *gin.RouterGroup, ctrl *v1.NotificationsController) {
	notifications := rg.Group("/notifications", middlewares.AuthJWT())
	{
		notifications.GET("", ctrl.Index)
		notifications.POST(":id/read", ctrl.Read)
		notifications.POST("/read-all", ctrl.ReadAll)
	}
}
