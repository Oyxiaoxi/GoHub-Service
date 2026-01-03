package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupUsersControllerTest() (*UsersController, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	controller := NewUsersController()
	router := gin.New()
	return controller, router
}

func TestUsersController_CurrentUser(t *testing.T) {
	controller, router := setupUsersControllerTest()

	t.Run("获取当前用户信息-未登录", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/user", controller.CurrentUser)

		req, _ := http.NewRequest("GET", "/api/v1/user", nil)
		router.ServeHTTP(w, req)

		// 未登录应该返回 401
		assert.Equal(t, 401, w.Code)
	})
}

func TestUsersController_Index(t *testing.T) {
	controller, router := setupUsersControllerTest()

	t.Run("获取用户列表-成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/users", controller.Index)

		req, _ := http.NewRequest("GET", "/api/v1/users?page=1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		
		// 验证响应结构
		assert.Contains(t, response, "data")
		assert.Contains(t, response, "pager")
	})

	t.Run("获取用户列表-无效分页参数", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/users", controller.Index)

		req, _ := http.NewRequest("GET", "/api/v1/users?page=invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})
}

func TestUsersController_Show(t *testing.T) {
	controller, router := setupUsersControllerTest()

	t.Run("获取指定用户-不存在的用户", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/users/:id", controller.Show)

		req, _ := http.NewRequest("GET", "/api/v1/users/99999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	t.Run("获取指定用户-无效ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/users/:id", controller.Show)

		req, _ := http.NewRequest("GET", "/api/v1/users/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}

func TestUsersController_UpdateProfile(t *testing.T) {
	controller, router := setupUsersControllerTest()

	t.Run("更新个人资料-缺少必填字段", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.PUT("/api/v1/user", controller.UpdateProfile)

		// 空请求体
		req, _ := http.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("更新个人资料-无效JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.PUT("/api/v1/user", controller.UpdateProfile)

		req, _ := http.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}

func TestUsersController_UpdateAvatar(t *testing.T) {
	controller, router := setupUsersControllerTest()

	t.Run("更新头像-未提供文件", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/user/avatar", controller.UpdateAvatar)

		req, _ := http.NewRequest("POST", "/api/v1/user/avatar", nil)
		req.Header.Set("Content-Type", "multipart/form-data")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})
}

func TestUsersController_UpdateEmail(t *testing.T) {
	controller, router := setupUsersControllerTest()

	t.Run("更新邮箱-无效格式", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.PUT("/api/v1/user/email", controller.UpdateEmail)

		body := map[string]interface{}{
			"email": "invalid-email",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/api/v1/user/email", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})
}

func TestUsersController_UpdatePassword(t *testing.T) {
	controller, router := setupUsersControllerTest()

	t.Run("修改密码-密码过短", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.PUT("/api/v1/user/password", controller.UpdatePassword)

		body := map[string]interface{}{
			"password":              "123",
			"password_confirmation": "123",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/api/v1/user/password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("修改密码-两次密码不一致", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.PUT("/api/v1/user/password", controller.UpdatePassword)

		body := map[string]interface{}{
			"password":              "password123",
			"password_confirmation": "password456",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/api/v1/user/password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})
}
