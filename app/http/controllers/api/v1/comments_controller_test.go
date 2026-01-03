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

func setupCommentsControllerTest() (*CommentsController, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	controller := NewCommentsController()
	router := gin.New()
	return controller, router
}

func TestCommentsController_Index(t *testing.T) {
	controller, router := setupCommentsControllerTest()

	t.Run("获取评论列表-成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/topics/:topic_id/comments", controller.Index)

		req, _ := http.NewRequest("GET", "/api/v1/topics/1/comments?page=1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response, "data")
	})

	t.Run("获取评论列表-无效话题ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/topics/:topic_id/comments", controller.Index)

		req, _ := http.NewRequest("GET", "/api/v1/topics/invalid/comments", nil)
		router.ServeHTTP(w, req)

		assert.True(t, w.Code == 404 || w.Code == 422)
	})
}

func TestCommentsController_Store(t *testing.T) {
	controller, router := setupCommentsControllerTest()

	t.Run("创建评论-缺少内容", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/topics/:topic_id/comments", controller.Store)

		body := map[string]interface{}{
			"topic_id": "1",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/topics/1/comments", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("创建评论-内容过短", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/topics/:topic_id/comments", controller.Store)

		body := map[string]interface{}{
			"topic_id": "1",
			"content":  "a", // 少于2个字符
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/topics/1/comments", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("创建评论-内容过长", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/topics/:topic_id/comments", controller.Store)

		longContent := make([]byte, 1100) // 超过1000字符
		for i := range longContent {
			longContent[i] = 'a'
		}

		body := map[string]interface{}{
			"topic_id": "1",
			"content":  string(longContent),
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/topics/1/comments", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("创建评论-无效JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/topics/:topic_id/comments", controller.Store)

		req, _ := http.NewRequest("POST", "/api/v1/topics/1/comments", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}

func TestCommentsController_Update(t *testing.T) {
	controller, router := setupCommentsControllerTest()

	t.Run("更新评论-空内容", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.PUT("/api/v1/comments/:id", controller.Update)

		body := map[string]interface{}{
			"content": "",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/api/v1/comments/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})
}

func TestCommentsController_Delete(t *testing.T) {
	controller, router := setupCommentsControllerTest()

	t.Run("删除评论-不存在的评论", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.DELETE("/api/v1/comments/:id", controller.Delete)

		req, _ := http.NewRequest("DELETE", "/api/v1/comments/99999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	t.Run("删除评论-无效ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.DELETE("/api/v1/comments/:id", controller.Delete)

		req, _ := http.NewRequest("DELETE", "/api/v1/comments/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}

func TestCommentsController_Show(t *testing.T) {
	controller, router := setupCommentsControllerTest()

	t.Run("获取评论详情-不存在的评论", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/comments/:id", controller.Show)

		req, _ := http.NewRequest("GET", "/api/v1/comments/99999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}
