// Package routes 注册路由
package routes

import (
	"GoHub-Service/app/http/controllers/api/v1/auth"
	"GoHub-Service/app/http/middlewares"
	"GoHub-Service/pkg/config"
	"GoHub-Service/pkg/metrics"

	controllers "GoHub-Service/app/http/controllers/api/v1"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册 API 相关路由
func RegisterAPIRoutes(r *gin.Engine) {
	// Prometheus 指标端点
	r.GET("/metrics", metrics.Handler())

	// 测试一个 v1 的路由组，我们所有的 v1 版本的路由都将存放到这里
	var v1 *gin.RouterGroup
	if len(config.Get("app.api_domain")) == 0 {
		v1 = r.Group("/api/v1")
	} else {
		v1 = r.Group("/v1")
	}

	v1.Use(
		middlewares.RequestID(),
		middlewares.LimitIP(config.GetString("limiter.global_ip_rate", "200-H")),
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
	commentsCtrl := controllers.NewCommentsController()
	notificationsCtrl := controllers.NewNotificationsController()

	// 搜索相关
	RegisterSearchRoutes(v1)

	// Auth 路由组
	authGroup := v1.Group("/auth", middlewares.LimitIP(config.GetString("limiter.auth_ip_rate", "1000-H")))
	{
		authGroup.POST("/login/using-phone", middlewares.GuestJWT(), loginCtrl.LoginByPhone)
		authGroup.POST("/login/using-password", middlewares.GuestJWT(), loginCtrl.LoginByPassword)
		authGroup.POST("/login/refresh-token", middlewares.AuthJWT(), loginCtrl.RefreshToken)

		authGroup.POST("/password-reset/using-email", middlewares.GuestJWT(), passwordCtrl.ResetByEmail)
		authGroup.POST("/password-reset/using-phone", middlewares.GuestJWT(), passwordCtrl.ResetByPhone)

		authGroup.POST("/signup/using-phone", middlewares.GuestJWT(), signupCtrl.SignupUsingPhone)
		authGroup.POST("/signup/using-email", middlewares.GuestJWT(), signupCtrl.SignupUsingEmail)
		authGroup.POST("/signup/phone/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute(config.GetString("limiter.signup_exist_rate", "60-H")), signupCtrl.IsPhoneExist)
		authGroup.POST("/signup/email/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute(config.GetString("limiter.signup_exist_rate", "60-H")), signupCtrl.IsEmailExist)

		authGroup.POST("/verify-codes/phone", middlewares.LimitPerRoute(config.GetString("limiter.verify_phone_rate", "20-H")), verifyCodeCtrl.SendUsingPhone)
		authGroup.POST("/verify-codes/email", middlewares.LimitPerRoute(config.GetString("limiter.verify_email_rate", "20-H")), verifyCodeCtrl.SendUsingEmail)
		authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute(config.GetString("limiter.verify_captcha_rate", "50-H")), verifyCodeCtrl.ShowCaptcha)
	}

	// 用户相关
	v1.GET("/user", middlewares.AuthJWT(), usersCtrl.CurrentUser)
	RegisterUserRoutes(v1, usersCtrl)

	// 分类相关
	RegisterCategoryRoutes(v1, categoriesCtrl)

	// 话题相关
	RegisterTopicRoutes(v1, topicsCtrl)

	// 友情链接
	RegisterLinkRoutes(v1, linksCtrl)

	// 评论相关
	RegisterCommentRoutes(v1, commentsCtrl)

	// 通知相关
	RegisterNotificationRoutes(v1, notificationsCtrl)

	// 私信相关
	RegisterMessageRoutes(v1)
}
