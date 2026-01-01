package routes

import (
	v1 "GoHub-Service/app/http/controllers/api/v1"
	"GoHub-Service/app/http/controllers"
	"GoHub-Service/app/http/middlewares"
	"GoHub-Service/app/repositories"
	"GoHub-Service/app/services"
	"GoHub-Service/pkg/config"
	"GoHub-Service/pkg/elasticsearch"

	"github.com/gin-gonic/gin"
)

// RegisterSearchRoutes 注册搜索相关路由
func RegisterSearchRoutes(r *gin.RouterGroup) {
	// 传统搜索服务
	repo := repositories.NewSearchRepository()
	service := services.NewSearchService(repo)
	controller := v1.NewSearchController(service)

	// Elasticsearch 搜索服务
	var esSearchCtrl *controllers.SearchController
	if config.GetBool("elasticsearch.enabled", false) {
		addresses := []string{"http://localhost:9200"}
		if client, err := elasticsearch.NewClient(addresses); err == nil {
			searchService := elasticsearch.NewSearchService(client)
			esSearchCtrl = controllers.NewSearchController(searchService)
		}
	}

	// 主题搜索 - 优先使用Elasticsearch，降级到传统搜索
	r.GET("/search/topics", func(c *gin.Context) {
		if esSearchCtrl != nil {
			esSearchCtrl.SearchTopics(c)
		} else {
			controller.Topics(c)
		}
	})

	// Elasticsearch 特定功能
	if esSearchCtrl != nil {
		r.GET("/search/suggestions", esSearchCtrl.SearchSuggestions)
		r.GET("/search/hot-topics", esSearchCtrl.GetHotTopics)
	}

	// 用户搜索需要登录
	r.GET("/search/users", middlewares.AuthJWT(), controller.Users)
}
