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

func setupCategoriesControllerTest() (*CategoriesController, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	controller := NewCategoriesController()
	router := gin.New()
	return controller, router
}

func TestCategoriesController_Index(t *testing.T) {
	controller, router := setupCategoriesControllerTest()

	t.Run("获取分类列表-成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/categories", controller.Index)

		req, _ := http.NewRequest("GET", "/api/v1/categories", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response, "data")
	})
}

func TestCategoriesController_Show(t *testing.T) {
	controller, router := setupCategoriesControllerTest()

	t.Run("获取分类详情-不存在的分类", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/categories/:id", controller.Show)

		req, _ := http.NewRequest("GET", "/api/v1/categories/99999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	t.Run("获取分类详情-无效ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.GET("/api/v1/categories/:id", controller.Show)

		req, _ := http.NewRequest("GET", "/api/v1/categories/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}

func TestCategoriesController_Store(t *testing.T) {
	controller, router := setupCategoriesControllerTest()

	t.Run("创建分类-缺少名称", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/categories", controller.Store)

		body := map[string]interface{}{
			"description": "Test description",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("创建分类-名称过短", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/categories", controller.Store)

		body := map[string]interface{}{
			"name":        "a",
			"description": "Test description",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})

	t.Run("创建分类-无效JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.POST("/api/v1/categories", controller.Store)

		req, _ := http.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}

func TestCategoriesController_Update(t *testing.T) {
	controller, router := setupCategoriesControllerTest()

	t.Run("更新分类-空请求体", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.PUT("/api/v1/categories/:id", controller.Update)

		req, _ := http.NewRequest("PUT", "/api/v1/categories/1", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 422, w.Code)
	})
}

func TestCategoriesController_Delete(t *testing.T) {
	controller, router := setupCategoriesControllerTest()

	t.Run("删除分类-不存在的分类", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.DELETE("/api/v1/categories/:id", controller.Delete)

		req, _ := http.NewRequest("DELETE", "/api/v1/categories/99999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}
