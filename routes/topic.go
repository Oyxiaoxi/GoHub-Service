// Package routes 话题相关路由
package routes

import (
	"GoHub-Service/app/http/controllers/api/v1"
	"GoHub-Service/app/http/middlewares"
	"github.com/gin-gonic/gin"
)

// RegisterTopicRoutes 注册话题相关路由
func RegisterTopicRoutes(rg *gin.RouterGroup, topicsCtrl *v1.TopicsController) {
	topicsGroup := rg.Group("/topics")
	{
		topicsGroup.GET("", topicsCtrl.Index)
		topicsGroup.POST("", middlewares.AuthJWT(), topicsCtrl.Store)
		topicsGroup.PUT(":id", middlewares.AuthJWT(), topicsCtrl.Update)
		topicsGroup.DELETE(":id", middlewares.AuthJWT(), topicsCtrl.Delete)
		topicsGroup.GET(":id", topicsCtrl.Show)
		topicsGroup.POST(":id/like", middlewares.AuthJWT(), topicsCtrl.Like)
		topicsGroup.POST(":id/unlike", middlewares.AuthJWT(), topicsCtrl.Unlike)
		topicsGroup.POST(":id/favorite", middlewares.AuthJWT(), topicsCtrl.Favorite)
		topicsGroup.POST(":id/unfavorite", middlewares.AuthJWT(), topicsCtrl.Unfavorite)
		topicsGroup.POST(":id/view", topicsCtrl.AddView)
	}
}
