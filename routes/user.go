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
		// 更新资料应用内容安全检查
		usersGroup.PUT("", 
			middlewares.AuthJWT(), 
			middlewares.SensitiveWordFilter(),
			usersCtrl.UpdateProfile,
		)
		// 敏感信息修改应用签名验证（防重放攻击）
		usersGroup.PUT("/email", 
			middlewares.AuthJWT(), 
			middlewares.APISignatureVerification(),
			usersCtrl.UpdateEmail,
		)
		usersGroup.PUT("/phone", 
			middlewares.AuthJWT(), 
			middlewares.APISignatureVerification(),
			usersCtrl.UpdatePhone,
		)
		usersGroup.PUT("/password", 
			middlewares.AuthJWT(), 
			middlewares.APISignatureVerification(),
			usersCtrl.UpdatePassword,
		)
		// 头像上传应用安全检查
		usersGroup.PUT("/avatar", 
			middlewares.AuthJWT(), 
			middlewares.ImageUploadSecurity(),
			usersCtrl.UpdateAvatar,
		)
		usersGroup.POST("/:id/follow", middlewares.AuthJWT(), usersCtrl.Follow)
		usersGroup.POST("/:id/unfollow", middlewares.AuthJWT(), usersCtrl.Unfollow)
	}
}
