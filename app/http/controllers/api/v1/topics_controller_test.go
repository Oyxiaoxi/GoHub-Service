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

func setupTopicsControllerTest() (*TopicsController, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	controller := NewTopicsController()
	router := gin.New()
	return controller, router
}

func TestTopicsController_Index(t *testing.T) {
	controller, router := setupTopicsControllerTest()

	t.Run("获取话题列表-成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/topics", controller.Index)

		req, _ := http.NewRequest("GET", "/api/v1/topics?page=1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response, "data")
	})

	t.Run("获取话题列表-分页参数超出范围", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/topics", controller.Index)

		req, _ := http.NewRequest("GET", "/api/v1/topics?page=0", nil)
		router.ServeHTTP(w, req)

		// 页码小于1应该返回错误
		assert.True(t, w.Code == 422 || w.Code == 400)
	})
}

func TestTopicsController_Show(t *testing.T) {
	controller, router := setupTopicsControllerTest()

	t.Run("获取话题详情-不存在的话题", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/topics/:id", controller.Show)

		req, _ := http.NewRequest("GET", "/api/v1/topics/99999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	t.Run("获取话题详情-无效ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/topics/:id", controller.Show)

		req, _ := http.NewRequest("GET", "/api/v1/topics/invalid-id", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}

func TestTopicsController_Store(t *testing.T) {
	controller, router := setupTopicsControllerTest()

	t.Run("创建话题-缺少标题", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/topics", controller.Store)

		body := map[string]interface{}{
			"body":        "Test body",
			"category_id": "1",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/topics", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("创建话题-标题过短", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/topics", controller.Store)

		body := map[string]interface{}{
			"title":       "ab", // 少于3个字符
			"body":        "Test body",
			"category_id": "1",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/topics", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("创建话题-无效分类ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/topics", controller.Store)

		body := map[string]interface{}{
			"title":       "Test Title",
			"body":        "Test body",
			"category_id": "invalid",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/topics", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})
}

func TestTopicsController_Update(t *testing.T) {
	controller, router := setupTopicsControllerTest()

	t.Run("更新话题-空请求体", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.PUT("/api/v1/topics/:id", controller.Update)

		req, _ := http.NewRequest("PUT", "/api/v1/topics/1", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("更新话题-无效JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.PUT("/api/v1/topics/:id", controller.Update)

		req, _ := http.NewRequest("PUT", "/api/v1/topics/1", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}

func TestTopicsController_Delete(t *testing.T) {
	controller, router := setupTopicsControllerTest()

	t.Run("删除话题-不存在的话题", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.DELETE("/api/v1/topics/:id", controller.Delete)

		req, _ := http.NewRequest("DELETE", "/api/v1/topics/99999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}

func TestTopicsController_IncrementViewCount(t *testing.T) {
	controller, router := setupTopicsControllerTest()

	t.Run("增加浏览量-无效ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/topics/:id/view", controller.IncrementViewCount)

		req, _ := http.NewRequest("POST", "/api/v1/topics/invalid/view", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}
