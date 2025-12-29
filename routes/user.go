// Package routes 用户相关路由
package routes

import (
	"GoHub-Service/app/http/controllers/api/v1"
	"GoHub-Service/app/http/middlewares"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(rg *gin.RouterGroup, usersCtrl *v1.UsersController) {
	usersGroup := rg.Group("/users")
	{
		usersGroup.GET("", usersCtrl.Index)
		usersGroup.PUT("", middlewares.AuthJWT(), usersCtrl.UpdateProfile)
		usersGroup.PUT("/email", middlewares.AuthJWT(), usersCtrl.UpdateEmail)
		usersGroup.PUT("/phone", middlewares.AuthJWT(), usersCtrl.UpdatePhone)
		usersGroup.PUT("/password", middlewares.AuthJWT(), usersCtrl.UpdatePassword)
		usersGroup.PUT("/avatar", middlewares.AuthJWT(), usersCtrl.UpdateAvatar)
	}
}
