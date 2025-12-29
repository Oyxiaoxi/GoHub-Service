// Package routes 注册路由
package routes

import (
    "GoHub-Service/app/http/controllers/api/v1/auth"
    "GoHub-Service/app/http/middlewares"

    controllers "GoHub-Service/app/http/controllers/api/v1"

    "github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册 API 相关路由

func RegisterAPIRoutes(r *gin.Engine) {
    v1 := r.Group("/v1")
    v1.Use(
        middlewares.RequestID(),
        middlewares.LimitIP("200-H"),
    )

    // 控制器实例化，便于复用
    loginCtrl := new(auth.LoginController)
    passwordCtrl := new(auth.PasswordController)
    signupCtrl := new(auth.SignupController)
    verifyCodeCtrl := new(auth.VerifyCodeController)
    usersCtrl := controllers.NewUsersController()
    categoriesCtrl := controllers.NewCategoriesController()
    topicsCtrl := controllers.NewTopicsController()
    linksCtrl := controllers.NewLinksController()

    // Auth 路由组
    authGroup := v1.Group("/auth", middlewares.LimitIP("1000-H"))
    {
        authGroup.POST("/login/using-phone", middlewares.GuestJWT(), loginCtrl.LoginByPhone)
        authGroup.POST("/login/using-password", middlewares.GuestJWT(), loginCtrl.LoginByPassword)
        authGroup.POST("/login/refresh-token", middlewares.AuthJWT(), loginCtrl.RefreshToken)

        authGroup.POST("/password-reset/using-email", middlewares.GuestJWT(), passwordCtrl.ResetByEmail)
        authGroup.POST("/password-reset/using-phone", middlewares.GuestJWT(), passwordCtrl.ResetByPhone)

        authGroup.POST("/signup/using-phone", middlewares.GuestJWT(), signupCtrl.SignupUsingPhone)
        authGroup.POST("/signup/using-email", middlewares.GuestJWT(), signupCtrl.SignupUsingEmail)
        authGroup.POST("/signup/phone/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), signupCtrl.IsPhoneExist)
        authGroup.POST("/signup/email/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), signupCtrl.IsEmailExist)

        authGroup.POST("/verify-codes/phone", middlewares.LimitPerRoute("20-H"), verifyCodeCtrl.SendUsingPhone)
        authGroup.POST("/verify-codes/email", middlewares.LimitPerRoute("20-H"), verifyCodeCtrl.SendUsingEmail)
        authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute("50-H"), verifyCodeCtrl.ShowCaptcha)
    }

    // 用户相关
    v1.GET("/user", middlewares.AuthJWT(), usersCtrl.CurrentUser)
    usersGroup := v1.Group("/users")
    {
        usersGroup.GET("", usersCtrl.Index)
        usersGroup.PUT("", middlewares.AuthJWT(), usersCtrl.UpdateProfile)
        usersGroup.PUT("/email", middlewares.AuthJWT(), usersCtrl.UpdateEmail)
        usersGroup.PUT("/phone", middlewares.AuthJWT(), usersCtrl.UpdatePhone)
        usersGroup.PUT("/password", middlewares.AuthJWT(), usersCtrl.UpdatePassword)
        usersGroup.PUT("/avatar", middlewares.AuthJWT(), usersCtrl.UpdateAvatar)
    }

    // 分类相关
    categoriesGroup := v1.Group("/categories")
    {
        categoriesGroup.GET("", categoriesCtrl.Index)
        categoriesGroup.POST("", middlewares.AuthJWT(), categoriesCtrl.Store)
        categoriesGroup.PUT(":id", middlewares.AuthJWT(), categoriesCtrl.Update)
        categoriesGroup.DELETE(":id", middlewares.AuthJWT(), categoriesCtrl.Delete)
    }

    // 话题相关
    topicsGroup := v1.Group("/topics")
    {
        topicsGroup.GET("", topicsCtrl.Index)
        topicsGroup.POST("", middlewares.AuthJWT(), topicsCtrl.Store)
        topicsGroup.PUT(":id", middlewares.AuthJWT(), topicsCtrl.Update)
        topicsGroup.DELETE(":id", middlewares.AuthJWT(), topicsCtrl.Delete)
        topicsGroup.GET(":id", topicsCtrl.Show)
    }

    // 友情链接
    linksGroup := v1.Group("/links")
    {
        linksGroup.GET("", linksCtrl.Index)
    }
}
