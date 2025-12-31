package middlewares

import (
	"net/http"

	"GoHub-Service/pkg/app"
	"GoHub-Service/pkg/config"
	"GoHub-Service/pkg/limiter"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// LimitConfig 限流配置结构
type LimitConfig struct {
	// 限流规则，如 "5-S" (5次/秒)
	Rate string
	// 自定义错误消息
	Message string
	// 是否返回剩余次数
	ShowRemaining bool
}

// LimitIP 全局限流中间件，针对 IP 进行限流
// limit 为格式化字符串，如 "5-S" ，示例:
//
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
func LimitIP(limit string) gin.HandlerFunc {
	return LimitIPWithConfig(LimitConfig{
		Rate:          limit,
		Message:       config.GetString("limiter.default_message", "请求过于频繁，请稍后再试"),
		ShowRemaining: true,
	})
}

// LimitIPWithConfig 带配置的 IP 限流中间件
func LimitIPWithConfig(config LimitConfig) gin.HandlerFunc {
	if app.IsTesting() {
		config.Rate = "1000000-H"
	}

	return func(c *gin.Context) {
		// 跳过 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 针对 IP 限流
		key := limiter.GetKeyIP(c)
		if ok := limitHandlerWithConfig(c, key, config); !ok {
			return
		}
		c.Next()
	}
}

// LimitPerRoute 限流中间件，用在单独的路由中
func LimitPerRoute(limit string) gin.HandlerFunc {
	return LimitPerRouteWithConfig(LimitConfig{
		Rate:          limit,
		Message:       config.GetString("limiter.default_message", "请求过于频繁，请稍后再试"),
		ShowRemaining: true,
	})
}

// LimitPerRouteWithConfig 带配置的路由限流中间件
func LimitPerRouteWithConfig(config LimitConfig) gin.HandlerFunc {
	if app.IsTesting() {
		config.Rate = "1000000-H"
	}

	return func(c *gin.Context) {
		// 跳过 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 针对单个路由，增加访问次数
		c.Set("limiter-once", false)

		// 针对 IP + 路由进行限流
		key := limiter.GetKeyRouteWithIP(c)
		if ok := limitHandlerWithConfig(c, key, config); !ok {
			return
		}
		c.Next()
	}
}

// LimitByUser 针对用户 ID 的限流中间件
// 需要在认证中间件之后使用
func LimitByUser(limit string) gin.HandlerFunc {
	return LimitByUserWithConfig(LimitConfig{
		Rate:          limit,
		Message:       config.GetString("limiter.default_message", "请求过于频繁，请稍后再试"),
		ShowRemaining: true,
	})
}

// LimitByUserWithConfig 带配置的用户限流中间件
func LimitByUserWithConfig(config LimitConfig) gin.HandlerFunc {
	if app.IsTesting() {
		config.Rate = "1000000-H"
	}

	return func(c *gin.Context) {
		// 跳过 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 获取当前用户 ID
		userID, exists := c.Get("user_id")
		if !exists {
			// 如果没有用户信息，降级为 IP 限流
			key := limiter.GetKeyIP(c)
			if ok := limitHandlerWithConfig(c, key, config); !ok {
				return
			}
		} else {
			// 使用用户 ID + 路由进行限流
			key := "limiter:" + c.FullPath() + ":user:" + cast.ToString(userID)
			if ok := limitHandlerWithConfig(c, key, config); !ok {
				return
			}
		}
		c.Next()
	}
}

// limitHandler 限流处理函数
func limitHandler(c *gin.Context, key string, limit string) bool {
	return limitHandlerWithConfig(c, key, LimitConfig{
		Rate:          limit,
		Message:       config.GetString("limiter.default_message", "请求过于频繁，请稍后再试"),
		ShowRemaining: true,
	})
}

// limitHandlerWithConfig 带配置的限流处理函数
func limitHandlerWithConfig(c *gin.Context, key string, config LimitConfig) bool {
	// 获取超额的情况
	rate, err := limiter.CheckRate(c, key, config.Rate)
	if err != nil {
		logger.LogIf(err)
		response.Abort500(c, "服务器内部错误")
		return false
	}

	// 设置响应头
	c.Header("X-RateLimit-Limit", cast.ToString(rate.Limit))
	c.Header("X-RateLimit-Remaining", cast.ToString(rate.Remaining))
	c.Header("X-RateLimit-Reset", cast.ToString(rate.Reset))

	// 超额
	if rate.Reached {
		// 构建响应数据
		responseData := gin.H{
			"code":    http.StatusTooManyRequests,
			"message": config.Message,
		}

		// 如果配置显示剩余次数
		if config.ShowRemaining {
			responseData["retry_after"] = rate.Reset
		}

		c.JSON(http.StatusTooManyRequests, responseData)
		c.Abort()
		return false
	}

	return true
}
