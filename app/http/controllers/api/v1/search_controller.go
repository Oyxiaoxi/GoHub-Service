package v1

import (
	"GoHub-Service/app/requests"
	"GoHub-Service/app/services"
	"GoHub-Service/pkg/response"

	"github.com/gin-gonic/gin"
)

// SearchController 提供搜索接口
type SearchController struct {
	SearchService services.SearchService
}

// NewSearchController 创建控制器
func NewSearchController(searchService services.SearchService) SearchController {
	return SearchController{SearchService: searchService}
}

// Topics 搜索主题
func (sc *SearchController) Topics(c *gin.Context) {
	request := requests.SearchRequest{}
	if ok := requests.Validate(c, &request, requests.SearchValidation); !ok {
		return
	}

	result, err := sc.SearchService.SearchTopics(c, request.Keyword)
	if err != nil {
		response.Abort500(c)
		return
	}

	response.JSON(c, gin.H{"topics": result})
}

// Users 搜索用户
func (sc *SearchController) Users(c *gin.Context) {
	request := requests.SearchRequest{}
	if ok := requests.Validate(c, &request, requests.SearchValidation); !ok {
		return
	}

	result, err := sc.SearchService.SearchUsers(c, request.Keyword)
	if err != nil {
		response.Abort500(c)
		return
	}

	response.JSON(c, gin.H{"users": result})
}
