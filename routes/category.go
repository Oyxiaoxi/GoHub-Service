// Package routes 分类相关路由
package routes

import (
	"GoHub-Service/app/http/controllers/api/v1"
	"GoHub-Service/app/http/middlewares"
	"github.com/gin-gonic/gin"
)

// RegisterCategoryRoutes 注册分类相关路由
func RegisterCategoryRoutes(rg *gin.RouterGroup, categoriesCtrl *v1.CategoriesController) {
	categoriesGroup := rg.Group("/categories")
	{
		categoriesGroup.GET("", categoriesCtrl.Index)
		categoriesGroup.POST("", middlewares.AuthJWT(), categoriesCtrl.Store)
		categoriesGroup.PUT(":id", middlewares.AuthJWT(), categoriesCtrl.Update)
		categoriesGroup.DELETE(":id", middlewares.AuthJWT(), categoriesCtrl.Delete)
	}
}
