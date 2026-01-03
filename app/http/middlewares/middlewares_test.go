// Package middlewares 中间件测试
package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestRequestID 测试请求ID中间件
func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("自动生成请求ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(RequestID())
		r.GET("/test", func(c *gin.Context) {
			requestID := c.GetString("request_id")
			assert.NotEmpty(t, requestID)
			c.String(200, requestID)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	})

	t.Run("使用客户端提供的请求ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(RequestID())
		r.GET("/test", func(c *gin.Context) {
			requestID := c.GetString("request_id")
			assert.Equal(t, "custom-request-id", requestID)
			c.String(200, "OK")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "custom-request-id")
		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "custom-request-id", w.Header().Get("X-Request-ID"))
	})
}

// TestCORS 测试CORS中间件
func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("处理预检请求", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(CORS())
		r.OPTIONS("/test", func(c *gin.Context) {
			c.String(200, "OK")
		})

		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
	})

	t.Run("处理正常请求", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(CORS())
		r.GET("/test", func(c *gin.Context) {
			c.String(200, "OK")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

// TestRecovery 测试Recovery中间件
func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("捕获panic", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(Recovery())
		r.GET("/test", func(c *gin.Context) {
			panic("test panic")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req
		r.ServeHTTP(w, req)

		// Recovery 应该返回 500 错误
		assert.Equal(t, 500, w.Code)
	})

	t.Run("正常请求不受影响", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(Recovery())
		r.GET("/test", func(c *gin.Context) {
			c.String(200, "OK")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})
}

// TestForceUA 测试强制User-Agent中间件
func TestForceUA(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("缺少User-Agent被拒绝", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(ForceUA())
		r.GET("/test", func(c *gin.Context) {
			c.String(200, "OK")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		// 不设置 User-Agent
		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("有User-Agent正常通过", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(ForceUA())
		r.GET("/test", func(c *gin.Context) {
			c.String(200, "OK")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}

// TestLimit 测试限流中间件
func TestLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("限流器正常工作", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		// 使用非常小的限流值测试
		r.Use(Limit("1-S"))
		r.GET("/test", func(c *gin.Context) {
			c.String(200, "OK")
		})

		// 第一个请求应该成功
		req1, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req1
		r.ServeHTTP(w, req1)
		assert.Equal(t, 200, w.Code)

		// 注意：实际限流测试需要真实的限流器实现
	})
}

// TestContentSecurity 测试内容安全中间件
func TestContentSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("添加安全头部", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(ContentSecurity())
		r.GET("/test", func(c *gin.Context) {
			c.String(200, "OK")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		// 检查安全头部是否存在
		assert.NotEmpty(t, w.Header().Get("X-Content-Type-Options"))
		assert.NotEmpty(t, w.Header().Get("X-Frame-Options"))
	})
}
