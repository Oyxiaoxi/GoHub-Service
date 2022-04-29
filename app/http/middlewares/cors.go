package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以用*，也可以用你指定的域名
		c.Header("Access-Control-Allow-Origin", "*")
		// 允许头部参数
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		// 允许的方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		method := c.Request.Method
		//放行OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		// 处理请求
		c.Next()
	}
}
