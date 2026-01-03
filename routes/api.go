// Package routes 注册路由
package routes

import (
	"GoHub-Service/app/http/controllers/api/v1/auth"
	"GoHub-Service/app/http/middlewares"
	"GoHub-Service/pkg/apiversion"
	"GoHub-Service/pkg/config"
	"GoHub-Service/pkg/metrics"

	controllers "GoHub-Service/app/http/controllers/api/v1"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterAPIRoutes 注册 API 相关路由
func RegisterAPIRoutes(r *gin.Engine) {
	// Prometheus 指标端点
	r.GET("/metrics", metrics.Handler())

	// Swagger 文档端点
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 版本信息端点
	r.GET("/api/versions", apiversion.GetVersionInfo())

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
		apiversion.VersionDeprecated("v1"), // 版本管理中间件
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

	// 数据库监控相关
	RegisterDatabaseStatsRoutes(r)
	
	// 缓存监控相关
	RegisterCacheMonitorRoutes(r)

	// Auth 路由组 - 增强安全限流
	authGroup := v1.Group("/auth", 
		middlewares.LimitIP(config.GetString("limiter.auth_ip_rate", "1000-H")),
		middlewares.RateLimitMiddleware(20), // 每分钟最多20次请求（IP级别自动封禁）
	)
	{
		// 登录接口 - 可选签名验证（兼容旧客户端）
		authGroup.POST("/login/using-phone", 
			middlewares.GuestJWT(), 
			middlewares.OptionalSignatureVerification(),
			loginCtrl.LoginByPhone,
		)
		authGroup.POST("/login/using-password", 
			middlewares.GuestJWT(), 
			middlewares.OptionalSignatureVerification(),
			loginCtrl.LoginByPassword,
		)
		authGroup.POST("/login/refresh-token", middlewares.AuthJWT(), loginCtrl.RefreshToken)

		// 密码重置 - 敏感操作，必须使用签名验证
		sensitiveGroup := authGroup.Group("/password-reset", 
			middlewares.RateLimitMiddleware(5), // 每分钟5次
			middlewares.APISignatureVerification(), // 强制签名验证
		)
		{
			sensitiveGroup.POST("/using-email", middlewares.GuestJWT(), passwordCtrl.ResetByEmail)
			sensitiveGroup.POST("/using-phone", middlewares.GuestJWT(), passwordCtrl.ResetByPhone)
		}

		authGroup.POST("/signup/using-phone", middlewares.GuestJWT(), signupCtrl.SignupUsingPhone)
		authGroup.POST("/signup/using-email", middlewares.GuestJWT(), signupCtrl.SignupUsingEmail)
		authGroup.POST("/signup/phone/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute(config.GetString("limiter.signup_exist_rate", "60-H")), signupCtrl.IsPhoneExist)
		authGroup.POST("/signup/email/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute(config.GetString("limiter.signup_exist_rate", "60-H")), signupCtrl.IsEmailExist)

		// 验证码 - 使用增强限流防止滥用
		authGroup.POST("/verify-codes/phone", middlewares.RateLimitMiddleware(10), middlewares.LimitPerRoute(config.GetString("limiter.verify_phone_rate", "20-H")), verifyCodeCtrl.SendUsingPhone)
		authGroup.POST("/verify-codes/email", middlewares.RateLimitMiddleware(10), middlewares.LimitPerRoute(config.GetString("limiter.verify_email_rate", "20-H")), verifyCodeCtrl.SendUsingEmail)
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

	// 管理后台路由
	RegisterAdminRoutes(r)
}
