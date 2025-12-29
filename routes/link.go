// Package routes 友情链接相关路由
package routes

import (
	"GoHub-Service/app/http/controllers/api/v1"
	"github.com/gin-gonic/gin"
)

// RegisterLinkRoutes 注册友情链接相关路由
func RegisterLinkRoutes(rg *gin.RouterGroup, linksCtrl *v1.LinksController) {
	linksGroup := rg.Group("/links")
	{
		linksGroup.GET("", linksCtrl.Index)
	}
}
