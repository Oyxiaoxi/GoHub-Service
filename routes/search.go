package routes

import (
	v1 "GoHub-Service/app/http/controllers/api/v1"
	"GoHub-Service/app/http/middlewares"
	"GoHub-Service/app/repositories"
	"GoHub-Service/app/services"

	"github.com/gin-gonic/gin"
)

// RegisterSearchRoutes 注册搜索相关路由
func RegisterSearchRoutes(r *gin.RouterGroup) {
	repo := repositories.NewSearchRepository()
	service := services.NewSearchService(repo)
	controller := v1.NewSearchController(service)

	// 主题搜索无需登录
	r.GET("/search/topics", controller.Topics)

	// 用户搜索需要登录
	r.GET("/search/users", middlewares.AuthJWT(), controller.Users)
}
